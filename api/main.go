package main

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mimsy-cms/mimsy/internal/logger"
	"github.com/mimsy-cms/mimsy/internal/migrations"
)

func main() {
	initLogger()

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

	mux := http.NewServeMux()
	v1 := http.NewServeMux()

	mux.Handle("/v1/", http.StripPrefix("/v1", v1))

	v1.HandleFunc("POST /auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:    net.JoinHostPort("localhost", cmp.Or(os.Getenv("APP_PORT"), "3000")),
		Handler: mux,
	}

	slog.Info("Starting server", "address", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Failed to start server", "error", err)
		return
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

func getPgURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		cmp.Or(os.Getenv("POSTGRES_PORT"), "5432"),
		os.Getenv("POSTGRES_DATABASE"))
}
