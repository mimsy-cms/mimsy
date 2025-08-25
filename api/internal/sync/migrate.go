package sync

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/mimsy-cms/mimsy/internal/collection"
	"github.com/mimsy-cms/mimsy/internal/migrations"
	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
	"github.com/mimsy-cms/mimsy/pkg/schema_diff"
	"github.com/mimsy-cms/mimsy/pkg/schema_generator"
	pgroll_migrations "github.com/xataio/pgroll/pkg/migrations"
)

type Migrator struct {
	generator            schema_generator.SchemaGenerator
	collectionRepository collection.Repository
}

func NewMigrator(collectionRepository collection.Repository) *Migrator {
	return &Migrator{
		generator:            schema_generator.New(),
		collectionRepository: collectionRepository,
	}
}

func (m *Migrator) GenerateSchema(ctx context.Context, schema *mimsy_schema.Schema) (*schema_generator.SqlSchema, error) {
	newSql, err := m.generator.GenerateSqlSchema(schema)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate new schema: %w", err)
	}

	return &newSql, nil
}

func (m *Migrator) UpdateCollections(ctx context.Context, schema *mimsy_schema.Schema) error {
	for _, collection := range schema.Collections {
		// Convert collection schema to JSON
		fieldsJson, err := json.Marshal(collection.Schema)
		if err != nil {
			return fmt.Errorf("Failed to marshal collection schema for %s: %w", collection.Name, err)
		}

		// Generate slug from collection name (convert to lowercase, replace spaces with underscores)
		slug := collection.Name

		// Check if collection exists
		exists, err := m.collectionRepository.CollectionExists(ctx, slug)
		if err != nil {
			return fmt.Errorf("Failed to check if collection %s exists: %w", slug, err)
		}

		if exists {
			// Update existing collection
			if err := m.collectionRepository.UpdateCollection(ctx, slug, collection.Name, fieldsJson); err != nil {
				return fmt.Errorf("Failed to update collection %s: %w", slug, err)
			}
			slog.Info("Updated collection", "slug", slug, "name", collection.Name)
		} else {
			// Create new collection
			if err := m.collectionRepository.CreateCollection(ctx, slug, collection.Name, fieldsJson, collection.IsGlobal); err != nil {
				return fmt.Errorf("Failed to create collection %s: %w", slug, err)
			}
			slog.Info("Created collection", "slug", slug, "name", collection.Name)
		}
	}

	return nil
}

func (m *Migrator) Migrate(ctx context.Context, activeSync *SyncStatus, newSql *schema_generator.SqlSchema, commitName string, commitHash string) error {
	// Decrypt the manifest from the previous activeMigration
	var activeSql *schema_generator.SqlSchema
	if activeSync == nil {
		activeSql = &schema_generator.SqlSchema{}
	} else if err := json.Unmarshal([]byte(activeSync.AppliedMigration), &activeSql); err != nil {
		return fmt.Errorf("Failed to unmarshal active schema: %w", err)
	}

	// Make the diff operation
	operations := schema_diff.Diff(*activeSql, *newSql)
	unrunMigrations := []*pgroll_migrations.Migration{
		{
			Name: fmt.Sprintf("%s", func() string {
				if len(commitHash) > 8 {
					return commitHash[:8]
				}
				return commitHash
			}()),
			Operations: operations,
		},
	}

	data, err := json.Marshal(unrunMigrations)
	if err != nil {
		return fmt.Errorf("Failed to marshal migrations: %w", err)
	}

	slog.Info("Pending Migrations", "migrations", unrunMigrations)

	runConfig := migrations.NewRunConfig(
		migrations.WithStateSchema("mimsy_collections_roll"),
		migrations.WithSchema("mimsy_collections"),
		migrations.WithSearchPath("mimsy_internal"),
		migrations.WithPgURL(getPgURL()),
		migrations.WithUnappliedMigrations(unrunMigrations),
	)

	count, err := migrations.Run(ctx, runConfig)

	if err != nil {
		return fmt.Errorf("Failed to run migrations: %w (Migrations: %s)", err, string(data))
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
