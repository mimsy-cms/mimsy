package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

const (
	// Parameters for Argon2
	Memory     = 64 * 1024
	Time       = 1
	Threads    = 4
	SaltLength = 16
	HashLength = 32
)

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func HashPassword(password string) (string, error) {
	salt, err := generateSalt(SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, Time, Memory, uint8(Threads), HashLength)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("%s$%s", encodedSalt, encodedHash)
	return encoded, nil
}

func CheckPasswordHash(password, encoded string) error {
	parts := strings.Split(encoded, "$")
	if len(parts) != 2 {
		return errors.New("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return fmt.Errorf("could not decode salt: %v", err)
	}

	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("could not decode hash: %v", err)
	}

	hash := argon2.IDKey([]byte(password), salt, Time, Memory, uint8(Threads), HashLength)

	if !compareHashes(hash, expectedHash) {
		return errors.New("password does not match")
	}

	return nil
}

func compareHashes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var user User
		err := db.QueryRow(`SELECT id, email, password FROM "user" WHERE email = $1`, req.Email).
			Scan(&user.ID, &user.Email, &user.PasswordHash)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		if err := CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		sessionToken, err := generateSessionToken()
		if err != nil {
			http.Error(w, "Failed to generate session", http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().Add(1 * time.Hour) // session valid for 1 hour

		_, err = db.Exec(`
			INSERT INTO session (id, user_id, expires_at)
			VALUES ($1, $2, $3)`, sessionToken, user.ID, expiresAt)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    sessionToken,
			Path:     "/",
			Expires:  expiresAt,
			HttpOnly: true,
			// Secure:   true, TODO: Set to true in production
			Secure:   false, // For local development
			SameSite: http.SameSiteLaxMode,
		})

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	}
}

func LogoutHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "No session found", http.StatusUnauthorized)
			return
		}

		_, err = db.Exec(`DELETE FROM session WHERE id = $1`, cookie.Value)
		if err != nil {
			http.Error(w, "Failed to delete session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			// Secure:   true, TODO: Set to true in production
			Secure:   false, // For local development
			SameSite: http.SameSiteLaxMode,
		})

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
	}
}
