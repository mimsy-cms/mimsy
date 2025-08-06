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

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"github.com/mimsy-cms/mimsy/internal/auth"
	"github.com/mimsy-cms/mimsy/internal/collection"
	"github.com/mimsy-cms/mimsy/internal/config"
	"github.com/mimsy-cms/mimsy/internal/logger"
	"github.com/mimsy-cms/mimsy/internal/media"
	"github.com/mimsy-cms/mimsy/internal/migrations"
	"github.com/mimsy-cms/mimsy/internal/storage"
	"github.com/mimsy-cms/mimsy/internal/util"
)

func main() {
	initLogger()
	storage := initStorage()

	ctx := context.Background()

	if err := storage.Authenticate(ctx); err != nil {
		slog.Error("Failed to authenticate storage", "error", err)
	}

	runConfig := migrations.NewRunConfig(
		migrations.WithMigrationsDir("./migrations"),
		migrations.WithPgURL(getPgURL()),
	)

	// NOTE: Migrations should not be run like this in production.
	migrationCount, err := migrations.Run(ctx, runConfig)
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

	authRepository := auth.NewRepository()
	authService := auth.NewAuthService(authRepository)
	authHandler := auth.NewHandler(authService)

	if err := authService.CreateAdminUser(config.ContextWithDB(ctx, db)); err != nil {
		fmt.Println("Failed to create admin user:", err)
		return
	}

	collectionRepository := collection.NewRepository(db)
	collectionService := collection.NewService(collectionRepository)
	collectionHandler := collection.NewHandler(collectionService)

	mediaRepository := media.NewRepository()
	mediaService := media.NewService(storage, mediaRepository)
	mediaHandler := media.NewHandler(mediaService)

	mux := http.NewServeMux()
	v1 := http.NewServeMux()

	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	v1.HandleFunc("POST /auth/login", authHandler.Login)
	v1.HandleFunc("POST /auth/logout", authHandler.Logout)
	v1.HandleFunc("POST /auth/password", authHandler.ChangePassword)
	v1.HandleFunc("POST /auth/register", authHandler.Register)
	v1.HandleFunc("GET /auth/me", authHandler.Me)
	v1.HandleFunc("GET /collections", collectionHandler.List)
	v1.HandleFunc("GET /collections/{collectionSlug}/definition", collectionHandler.Definition)
	v1.HandleFunc("POST /media", mediaHandler.Upload)
	v1.HandleFunc("GET /media", mediaHandler.FindAll)
	v1.HandleFunc("GET /media/{id}", mediaHandler.GetById)
	v1.HandleFunc("DELETE /media/{id}", mediaHandler.Delete)
	v1.HandleFunc("GET /users", authHandler.GetUsers)

	handler := util.ApplyMiddlewares(
		config.WithDB(db),
		auth.WithUser(authService),
	)

	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", cmp.Or(os.Getenv("APP_PORT"), "3000")),
		Handler: handler(mux),
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
			storage.WithSwiftSecretKey(os.Getenv("SWIFT_SECRET_KEY")),
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
