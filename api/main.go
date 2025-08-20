package main

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
	"github.com/lib/pq"
	pgroll_migrations "github.com/xataio/pgroll/pkg/migrations"
	"github.com/xataio/pgroll/pkg/roll"
	"github.com/xataio/pgroll/pkg/state"

	"github.com/mimsy-cms/mimsy/internal/auth"
	"github.com/mimsy-cms/mimsy/internal/collection"
	"github.com/mimsy-cms/mimsy/internal/config"
	"github.com/mimsy-cms/mimsy/internal/cron"
	"github.com/mimsy-cms/mimsy/internal/logger"
	"github.com/mimsy-cms/mimsy/internal/media"
	"github.com/mimsy-cms/mimsy/internal/migrations"
	"github.com/mimsy-cms/mimsy/internal/storage"
	"github.com/mimsy-cms/mimsy/internal/sync"
	"github.com/mimsy-cms/mimsy/internal/util"
)

func main() {
	ctx := context.Background()

	initLogger()

	storage, err := initStorage(ctx)
	if err != nil {
		slog.Error("Failed to initialize storage", "error", err)
		return
	}

	internalRoll, err := initRoll(ctx, "mimsy_internal", "mimsy_internal_roll")
	if err != nil {
		slog.Error("Failed to initialize roll", "error", err)
		return
	}

	_, err = initRoll(ctx, "mimsy_collections", "mimsy_collections_roll")
	if err != nil {
		slog.Error("Failed to initialize roll", "error", err)
		return
	}

	migrationsCount, err := runMigrations(ctx, internalRoll)
	if err != nil {
		slog.Error("Failed to run migrations", "error", err)
		return
	}

	slog.Info("Successfully ran migrations", "count", migrationsCount)

	db, err := sql.Open("postgres", getPgURL()+"&search_path=mimsy_internal,mimsy_collections")
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	authRepository := auth.NewRepository()
	authService := auth.NewService(authRepository)
	authHandler := auth.NewHandler(authService)

	cronService := initCron(db)
	collectionRepository := collection.NewRepository()

	initSync(db, cronService)
	syncHandler := sync.NewHandler(cronService)

	// Start the cron scheduler
	if err := cronService.Start(ctx); err != nil {
		slog.Error("Failed to start cron service", "error", err)
		return
	}

	if err := authService.CreateAdminUser(config.ContextWithDB(ctx, db)); err != nil {
		fmt.Println("Failed to create admin user:", err)
		return
	}

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
	v1.HandleFunc("GET /collections", collectionHandler.FindAll)
	v1.HandleFunc("GET /collections/{slug}", collectionHandler.GetResources)
	v1.HandleFunc("GET /collections/{slug}/{resourceSlug}", collectionHandler.GetResource)
	v1.HandleFunc("PUT /collections/{slug}/{resourceSlug}", collectionHandler.UpdateResource)
	v1.HandleFunc("POST /collections/{slug}", collectionHandler.CreateResource)
	v1.HandleFunc("GET /collections/{slug}/definition", collectionHandler.Definition)
	v1.HandleFunc("DELETE /collections/{slug}/{resourceSlug}", collectionHandler.DeleteResource)
	v1.HandleFunc("GET /collections/globals", collectionHandler.FindAllGlobals)
	v1.HandleFunc("POST /media", mediaHandler.Upload)
	v1.HandleFunc("GET /media", mediaHandler.FindAll)
	v1.HandleFunc("GET /media/{id}", mediaHandler.GetById)
	v1.HandleFunc("DELETE /media/{id}", mediaHandler.Delete)
	v1.HandleFunc("GET /users", authHandler.GetUsers)
	v1.HandleFunc("GET /sync/status", syncHandler.Status)
	v1.HandleFunc("GET /sync/jobs", syncHandler.Jobs)
	v1.HandleFunc("GET /sync/active-migration", syncHandler.ActiveMigration)

	handler := util.ApplyMiddlewares(
		config.WithDB(db),
		auth.WithRequestUser(authService),
	)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cmp.Or(os.Getenv("APP_PORT"), "3000")),
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

func initStorage(ctx context.Context) (storage.Storage, error) {
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
		return nil, fmt.Errorf("unsupported storage type: %s", os.Getenv("STORAGE"))
	}

	if err := s.Authenticate(ctx); err != nil {
		return nil, fmt.Errorf("error during storage authentication: %w", err)
	}

	return s, nil
}

func initCron(db *sql.DB) cron.CronService {
	cronService, err := cron.NewCronService(db)

	if err != nil {
		slog.Error("Failed to initialize cron service", "error", err)
		panic("Error during setup.")
	}

	return cronService
}

func initSync(db *sql.DB, cronService cron.CronService) sync.SyncProvider {
	slog.Info("Initializing sync service")

	var pemKey string
	if keyPath := os.Getenv("GH_PEM_KEY_FILE"); keyPath != "" {
		// Read the file, and output that
		pemKeyTemp, err := os.ReadFile(keyPath)
		if err != nil {
			slog.Error("Failed to read GitHub PEM key file", "error", err)
			panic("Error during setup.")
		}

		pemKey = string(pemKeyTemp)
	} else {
		pemKey = os.Getenv("GH_PEM_KEY")
		if len(pemKey) == 0 {
			slog.Error("GitHub PEM key not set")
			panic("Error during setup.")
		}
	}

	appId, err := strconv.ParseInt(os.Getenv("GH_APP_ID"), 10, 64)
	if err != nil {
		slog.Error("Failed to parse GitHub App ID", "error", err)
		panic("Error during setup.")
	}

	syncService, err := sync.New(
		db,
		string(pemKey),
		appId,
		os.Getenv("GH_REPO"),
	)
	if err != nil {
		slog.Error("Failed to initialize sync service", "error", err)
		panic("Error during setup.")
	}

	// Register the job
	syncService.RegisterSyncJobs(cronService)

	slog.Info("Sync service initialized")

	return syncService
}

// initRoll initializes pgroll states
//
// The roll cli commands that it replaces are:
// - pgroll init --postgres-url "<postgres-url>" --schema mimsy --pgroll-schema mimsy_internal
// - pgroll init --postgres-url "<postgres-url>" --schema mimsy --pgroll-schema mimsy_collections
func initRoll(ctx context.Context, schema string, pgrollSchema string) (*roll.Roll, error) {
	internalState, err := state.New(ctx, getPgURL(), pgrollSchema)
	if err != nil {
		slog.Error("Failed to create internal state", "error", err)
		return nil, err
	}

	db, err := sql.Open("postgres", getPgURL())
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		return nil, err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", pq.QuoteIdentifier(schema)))
	if err != nil {
		slog.Error("Failed to create schema", "error", err)
		return nil, err
	}

	if isInitialized, _ := internalState.IsInitialized(ctx); !isInitialized {
		slog.Info("Initializing internal state")
	}

	if err := internalState.Init(ctx); err != nil {
		slog.Error("Failed to initialize state", "error", err)
	}

	return roll.New(ctx, getPgURL(), schema, internalState)
}

func runMigrations(ctx context.Context, m *roll.Roll) (int, error) {
	rawMigs, err := m.UnappliedMigrations(ctx, os.DirFS("./migrations"))
	if err != nil {
		return 0, fmt.Errorf("failed to get unapplied migrations: %w", err)
	}

	unappliedMigrations := make([]*pgroll_migrations.Migration, 0, len(rawMigs))
	for _, rawMig := range rawMigs {
		mig, err := pgroll_migrations.ParseMigration(rawMig)
		if err != nil {
			return 0, fmt.Errorf("failed to parse migration %q: %w", rawMig.Name, err)
		}
		unappliedMigrations = append(unappliedMigrations, mig)
	}

	runConfig := migrations.NewRunConfig(
		migrations.WithStateSchema("mimsy_internal_roll"),
		migrations.WithSchema("mimsy_internal"),
		migrations.WithUnappliedMigrations(unappliedMigrations),
		migrations.WithPgURL(getPgURL()),
	)

	// NOTE: Migrations should not be run like this in production.
	migrationCount, err := migrations.Run(ctx, runConfig)
	if err != nil {
		return 0, fmt.Errorf("failed to run migrations: %w", err)
	}

	return migrationCount, nil
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
