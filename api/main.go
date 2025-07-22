package main

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mimsy-cms/mimsy/internal/migrations"
	"github.com/mimsy-cms/mimsy/internal/storage"
)

func main() {
	storage := initStorage()

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

	mux := http.NewServeMux()
	v1 := http.NewServeMux()

	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	v1.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	v1.HandleFunc("POST /collections/media", func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(256 * 1024) // 256 MB

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Failed to get file from form", http.StatusBadRequest)
			return
		}
		defer file.Close()

		contentType := header.Header.Get("Content-Type")
		if contentType == "" {
			http.Error(w, "Content-Type header is missing", http.StatusBadRequest)
			return
		}

		id, err := uuid.NewV7()
		if err != nil {
			http.Error(w, "Failed to generated uuid", http.StatusInternalServerError)
			return
		}

		if err := storage.Upload(r.Context(), id.String(), file, contentType); err != nil {
			http.Error(w, "Failed to upload file", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", cmp.Or(os.Getenv("APP_PORT"), "3000")),
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Failed to start server:", err)
		return
	}
}

func initStorage() storage.Storage {
	var s storage.Storage

	switch os.Getenv("STORAGE") {
	case "swift":
		s = storage.NewSwift(
			storage.WithSwiftUsername(os.Getenv("SWIFT_USERNAME")),
			storage.WithSwiftApiKey(os.Getenv("SWIFT_API_KEY")),
			storage.WithSwiftAuthURL(os.Getenv("SWIFT_AUTH_URL")),
			storage.WithSwiftDomain(os.Getenv("SWIFT_DOMAIN")),
			storage.WithSwiftTenant(os.Getenv("SWIFT_TENANT")),
			storage.WithSwiftContainer(os.Getenv("SWIFT_CONTAINER")),
			storage.WithSwiftRegion(os.Getenv("SWIFT_REGION")),
		)

		slog.Info("Using Swift storage backend", "container", os.Getenv("SWIFT_CONTAINER"))
	default:
		slog.Info("No storage backend configured")
	}

	return s
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
