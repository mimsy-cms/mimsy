package main

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"github.com/mimsy-cms/mimsy/internal/auth"
	"github.com/mimsy-cms/mimsy/internal/collection"
	"github.com/mimsy-cms/mimsy/internal/logger"
	"github.com/mimsy-cms/mimsy/internal/migrations"
	"github.com/mimsy-cms/mimsy/internal/storage"
)

func main() {
	initLogger()
	storage := initStorage()

	runConfig := migrations.NewRunConfig(
		migrations.WithMigrationsDir("./migrations"),
		migrations.WithPgURL(getPgURL()),
	)

	// NOTE: Migrations should not be run like this in production.
	migrationCount, err := migrations.Run(context.Background(), runConfig)
	if err != nil {
		slog.Error("Failed to run migrations", "error", err)
	} else {
		slog.Info("Successfully ran migrations", "count", migrationCount)
	}

	db, err := sql.Open("postgres", getPgURL())
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	if err := auth.CreateAdminUser(context.Background(), db); err != nil {
		fmt.Println("Failed to create admin user:", err)
		return
	}

	router := chi.NewRouter()

	router.Group(func(r chi.Router) {
		r.Use(auth.WithUser(db))

		r.Route("/v1", func(v1 chi.Router) {
			v1.Post("/auth/login", auth.LoginHandler(db))
			v1.Post("/auth/logout", auth.LogoutHandler(db))
			v1.Post("/auth/password", auth.ChangePasswordHandler(db))
			v1.Post("/auth/register", auth.RegisterHandler(db))
			v1.Get("/auth/me", auth.MeHandler(db))
			v1.Route("/collections", func(c chi.Router) {
				c.Get("/{collectionSlug}/definition", collection.DefinitionHandler(db))
				c.Get("/{resourceSlug}/items", collection.ItemsHandler(db))

				c.Post("/media", func(w http.ResponseWriter, r *http.Request) {
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
						http.Error(w, "Failed to generate uuid", http.StatusInternalServerError)
						return
					}

					if err := storage.Upload(r.Context(), id.String(), file, contentType); err != nil {
						http.Error(w, "Failed to upload file", http.StatusInternalServerError)
						return
					}

					w.WriteHeader(http.StatusCreated)
				})
			})
		})
	})

	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", cmp.Or(os.Getenv("APP_PORT"), "3000")),
		Handler: router,
	}

	slog.Info("Starting server", "address", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Failed to start server", "error", err)
	}
}

// initLogger initializes the logger with the specified format and level.
// It defaults to text format and info level if not specified in the environment.
// Supported log formats are "text" and "json".
// Supported log levels are "debug", "info", "warn", and "error".
func initLogger() {
	logFormat := cmp.Or(os.Getenv("LOG_FORMAT"), logger.LogFormatText)
	level := logger.LevelToSlogLevel(cmp.Or(os.Getenv("LOG_LEVEL"), "info"))

	options := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if logFormat == logger.LogFormatJSON {
		handler = slog.NewJSONHandler(os.Stdout, options)
	} else {
		handler = slog.NewTextHandler(os.Stdout, options)
	}

	slog.SetDefault(slog.New(handler))
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
