package main

import (
	"cmp"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"database/sql"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"github.com/mimsy-cms/mimsy/internal/migrations"
	"github.com/mimsy-cms/mimsy/internal/auth"
)

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

	mux := http.NewServeMux()
	v1 := http.NewServeMux()

	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	v1.HandleFunc("POST /auth/login", auth.LoginHandler(db))

	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", cmp.Or(os.Getenv("APP_PORT"), "3000")),
		Handler: mux,
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
