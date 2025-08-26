package sync

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

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

	// Make the diff operation and check for skipped alterations
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

	if len(operations) == 0 {
		slog.Info("No migrations operations to run", "commit", commitHash)
		return nil
	}

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
	}

	if count == 0 && len(skippedAlterations) == 0 {
		return fmt.Errorf("No migrations were run")
	}

	return nil
}

type SkippedAlterationsError struct {
	Message     string              `json:"message"`
	Alterations []SkippedAlteration `json:"alterations"`
	CommitHash  string              `json:"commit_hash"`
}

func (e *SkippedAlterationsError) Error() string {
	return fmt.Sprintf("%s %s", e.Message, formatSkippedAlterations(e.Alterations))
}

type SkippedAlteration struct {
	Table      string `json:"table"`
	Column     string `json:"column"`
	ChangeType string `json:"change_type"`
	OldValue   string `json:"old_value"`
	NewValue   string `json:"new_value"`
	Reason     string `json:"reason"`
}

func checkForSkippedAlterations(oldSchema, newSchema schema_generator.SqlSchema) []SkippedAlteration {
	var skipped []SkippedAlteration

	for _, newTable := range newSchema.Tables {
		oldTable, exists := oldSchema.GetTable(newTable.Name)
		if !exists {
			continue
		}

		for _, newColumn := range newTable.Columns {
			oldColumn, exists := oldTable.GetColumn(newColumn.Name)
			if !exists {
				continue
			}

			if newColumn.Type != oldColumn.Type {
				skipped = append(skipped, SkippedAlteration{
					Table:      newTable.Name,
					Column:     newColumn.Name,
					ChangeType: "type",
					OldValue:   oldColumn.Type,
					NewValue:   newColumn.Type,
					Reason:     "Type changes are disabled to prevent data loss",
				})
			}

			if newColumn.IsNotNull != oldColumn.IsNotNull {
				skipped = append(skipped, SkippedAlteration{
					Table:      newTable.Name,
					Column:     newColumn.Name,
					ChangeType: "nullability",
					OldValue:   fmt.Sprintf("nullable=%t", !oldColumn.IsNotNull),
					NewValue:   fmt.Sprintf("nullable=%t", !newColumn.IsNotNull),
					Reason:     "Nullability changes are disabled to prevent constraint violations",
				})
			}

			if newColumn.DefaultValue != oldColumn.DefaultValue {
				skipped = append(skipped, SkippedAlteration{
					Table:      newTable.Name,
					Column:     newColumn.Name,
					ChangeType: "default",
					OldValue:   oldColumn.DefaultValue,
					NewValue:   newColumn.DefaultValue,
					Reason:     "Default value changes are disabled",
				})
			}
		}
	}

	return skipped
}

func formatSkippedAlterations(skipped []SkippedAlteration) string {
	if len(skipped) == 0 {
		return ""
	}

	summary := make(map[string]int)
	for _, alt := range skipped {
		summary[alt.ChangeType]++
	}

	var parts []string
	for changeType, count := range summary {
		parts = append(parts, fmt.Sprintf("%d %s changes", count, changeType))
	}

	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
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
