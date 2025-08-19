package sync

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/mimsy-cms/mimsy/internal/migrations"
	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
	"github.com/mimsy-cms/mimsy/pkg/schema_diff"
	"github.com/mimsy-cms/mimsy/pkg/schema_generator"
	pgroll_migrations "github.com/xataio/pgroll/pkg/migrations"
)

type Migrator struct {
	generator schema_generator.SchemaGenerator
}

func NewMigrator() *Migrator {
	return &Migrator{
		generator: schema_generator.New(),
	}
}

func (m *Migrator) Migrate(ctx context.Context, activeSync *SyncStatus, newSchema *mimsy_schema.Schema, commitName string, commitHash string) error {
	// Decrypt the manifest from the previous activeMigration
	var activeSql *schema_generator.SqlSchema
	if activeSync == nil {
		activeSql = &schema_generator.SqlSchema{}
	} else if err := json.Unmarshal([]byte(activeSync.AppliedMigration), &activeSql); err != nil {
		return fmt.Errorf("Failed to unmarshal active schema: %w", err)
	}

	newSql, err := m.generator.GenerateSqlSchema(newSchema)
	if err != nil {
		return fmt.Errorf("Failed to generate new schema: %w", err)
	}

	// Make the diff operation
	operations := schema_diff.Diff(*activeSql, newSql)
	unrunMigrations := []*pgroll_migrations.Migration{
		{
			Name:       fmt.Sprintf("%s (hash:%s)", commitName, commitHash[:8]),
			Operations: operations,
		},
	}

	data, err := json.Marshal(unrunMigrations)
	if err != nil {
		return fmt.Errorf("Failed to marshal unrun migrations: %w", err)
	}

	slog.Info("Unrun migrations", "data", string(data))

	runConfig := migrations.NewRunConfig(
		migrations.WithStateSchema("mimsy_collections_roll"),
		migrations.WithSchema("mimsy_collections"),
		migrations.WithSearchPath("mimsy_internal"),
		migrations.WithPgURL(getPgURL()),
		migrations.WithUnappliedMigrations(unrunMigrations),
	)

	count, err := migrations.Run(ctx, runConfig)

	if err != nil {
		return fmt.Errorf("Failed to run migrations: %w", err)
	} else if count == 0 {
		return fmt.Errorf("No migrations were run")
	}

	return nil
}

// TODO(Red): Remove this, and modify the migration system to take this as a config.
func getPgURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		cmp.Or(os.Getenv("POSTGRES_PORT"), "5432"),
		os.Getenv("POSTGRES_DATABASE"))
}
