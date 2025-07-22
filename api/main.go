package main

import (
	"cmp"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/argon2"

	"github.com/mimsy-cms/mimsy/internal/auth"
	"github.com/mimsy-cms/mimsy/internal/migrations"
)

const (
	adminEmail    = "admin@example.com"
	adminPassword = "admin123"

	// Argon2id parameters
	memory     = 64 * 1024
	time       = 1
	threads    = 4
	saltLength = 16
	keyLength  = 32
)

func createAdminUser(ctx context.Context, db *sql.DB) error {
	var exists bool
	err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM "user" WHERE email = $1)`, adminEmail).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if admin user exists: %w", err)
	}
	if exists {
		return nil // Admin user already exists
	}

	hash, err := generatePasswordHash(adminPassword)
	if err != nil {
		return fmt.Errorf("failed to create admin user hash: %w", err)
	}

	_, err = db.ExecContext(ctx, `INSERT INTO "user" (email, password, must_change_password) VALUES ($1, $2, $3)`, adminEmail, hash, true)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	log.Printf("Initial admin user %s created with temporary password %s", adminEmail, adminPassword)
	return nil
}

func generatePasswordHash(password string) (string, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, time, memory, uint8(threads), keyLength)
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	hashB64 := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s$%s", saltB64, hashB64), nil
}

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "http://localhost:5173" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-SvelteKit-Action")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		}

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	runConfig := migrations.NewRunConfig(
		migrations.WithMigrationsDir("./migrations"),
		migrations.WithPgURL(getPgURL()),
	)

	// NOTE: Migrations should not be run like this in production.
	migrationCount, err := migrations.Run(context.Background(), runConfig)
	if err != nil {
		fmt.Println("Failed to run migrations:", err)
	} else {
		fmt.Printf("Successfully ran %d migrations\n", migrationCount)
	}

	db, err := sql.Open("postgres", getPgURL())
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	if err := createAdminUser(context.Background(), db); err != nil {
		fmt.Println("Failed to create admin user:", err)
		return
	}

	mux := http.NewServeMux()
	v1 := http.NewServeMux()

	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	v1.HandleFunc("POST /auth/login", auth.LoginHandler(db))
	v1.HandleFunc("POST /auth/logout", auth.LogoutHandler(db))
	v1.HandleFunc("POST /auth/password", auth.ChangePasswordHandler(db))

	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", cmp.Or(os.Getenv("APP_PORT"), "3000")),
		Handler: WithCORS(mux),
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Failed to start server:", err)
		return
	}
}

func getPgURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		cmp.Or(os.Getenv("POSTGRES_PORT"), "5432"),
		os.Getenv("POSTGRES_DATABASE"))
}
