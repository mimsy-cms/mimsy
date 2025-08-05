package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	auth_interface "github.com/mimsy-cms/mimsy/internal/interfaces/auth"
	"github.com/mimsy-cms/mimsy/internal/util"
	"golang.org/x/crypto/argon2"
)

var (
	// Admin credentials from environment variables
	adminEmail    = os.Getenv("ADMIN_EMAIL")
	adminPassword = os.Getenv("ADMIN_PASSWORD")
)

const (
	// Parameters for Argon2
	Memory     = 64 * 1024
	Time       = 1
	Threads    = 4
	SaltLength = 16
	HashLength = 32
)

func CreateAdminUser(ctx context.Context, db *sql.DB) error {
	var userCount int
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM "user"`).Scan(&userCount)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}
	if userCount > 0 {
		return nil // Admin user already exists
	}

	hash, err := HashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("failed to create admin user hash: %w", err)
	}

	_, err = db.ExecContext(ctx, `INSERT INTO "user" (email, password, must_change_password, is_admin) VALUES ($1, $2, $3, $4)`, adminEmail, hash, true, true)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Printf("Initial admin user %s created with temporary password %s", adminEmail, adminPassword)
	return nil
}

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
	ID                 int64
	Email              string
	PasswordHash       string
	IsAdmin            bool
	MustChangePassword bool
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

type LoginResponse struct {
	Session            string `json:"session"`
	MustChangePassword bool   `json:"mustChangePassword"`
}

func LoginHandler(db auth_interface.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := util.DecodeJSON[LoginRequest](r)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var user User
		if err = db.QueryRow(`SELECT id, email, password, must_change_password FROM "user" WHERE email = $1`, req.Email).
			Scan(&user.ID, &user.Email, &user.PasswordHash, &user.MustChangePassword); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		log.Printf("User: %+v", user)
		log.Printf("Checking password: %s", req.Password)
		if err := CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
			log.Printf("Password check failed: %v", err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// if err := CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		// 	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		// 	return
		// }

		// Clean up expired sessions
		if _, err := db.Exec(`DELETE FROM "session" WHERE expires_at < NOW()`); err != nil {
			http.Error(w, "Failed to clean up expired sessions", http.StatusInternalServerError)
			return
		}

		sessionToken, err := generateSessionToken()
		if err != nil {
			http.Error(w, "Failed to generate session", http.StatusInternalServerError)
			return
		}

		expiresAt := time.Now().Add(7 * 24 * time.Hour) // session valid for 7 days

		_, err = db.Exec(
			`INSERT INTO session (id, user_id, expires_at)
			VALUES ($1, $2, $3)`, sessionToken, user.ID, expiresAt)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		util.JSON(w, http.StatusOK, LoginResponse{
			Session:            sessionToken,
			MustChangePassword: user.MustChangePassword,
		})
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

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully", "session": cookie.Value})
	}
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func ChangePasswordHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := UserFromContext(r.Context())
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		req, err := util.DecodeJSON[ChangePasswordRequest](r)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.OldPassword == "" || req.NewPassword == "" {
			http.Error(w, "Old and new passwords are required", http.StatusBadRequest)
			return
		}

		var currentHash string
		err = db.QueryRow(`SELECT password FROM "user" WHERE id = $1`, user.ID).Scan(&currentHash)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusInternalServerError)
			return
		}

		if err := CheckPasswordHash(req.OldPassword, currentHash); err != nil {
			http.Error(w, "Old password is incorrect", http.StatusUnauthorized)
			return
		}

		newHash, err := HashPassword(req.NewPassword)
		if err != nil {
			http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`UPDATE "user" SET password = $1, must_change_password = FALSE WHERE id = $2`, newHash, user.ID)
		if err != nil {
			http.Error(w, "Failed to update password", http.StatusInternalServerError)
			return
		}

		util.JSON(w, http.StatusOK, struct{}{})
	}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

func RegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := util.DecodeJSON[CreateUserRequest](r)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		email := strings.TrimSpace(req.Email)
		if email == "" || len(req.Password) < 8 {
			http.Error(w, "Invalid email or password too short", http.StatusBadRequest)
			return
		}

		var exists bool
		if err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM "user" WHERE email=$1)`, email).Scan(&exists); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if exists {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		hashed, err := HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(
			`INSERT INTO "user" (email, password, is_admin, must_change_password) VALUES ($1, $2, $3, TRUE)`,
			email, hashed, req.IsAdmin,
		)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		util.JSON(w, http.StatusCreated, struct{}{})
	}
}

type MeResponse struct {
	ID                 int64  `json:"id"`
	Email              string `json:"email"`
	IsAdmin            bool   `json:"is_admin"`
	MustChangePassword bool   `json:"must_change_password"`
}

func MeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		user := UserFromContext(r.Context())
		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		util.JSON(w, http.StatusOK, MeResponse{
			ID:                 user.ID,
			Email:              user.Email,
			IsAdmin:            user.IsAdmin,
			MustChangePassword: user.MustChangePassword,
		})
	}
}
