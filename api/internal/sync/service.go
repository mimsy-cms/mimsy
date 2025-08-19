package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/mimsy-cms/mimsy/internal/collection"
	"github.com/mimsy-cms/mimsy/internal/cron"
	"github.com/mimsy-cms/mimsy/pkg/github_fetcher"
	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
)

type Status int

const (
	PROVIDER_ERROR Status = iota
	PROCESSING
	SUCCESS
	NO_MANIFEST
	INVALID_MANIFEST
)

type SyncProvider interface {
	GetStatus(ctx context.Context) (Status, error)
	RegisterSyncJobs(cronService cron.CronService) error
	SyncRepository(ctx context.Context) error
}

type syncProvider struct {
	githubClient         github_fetcher.GithubProvider
	repositoryName       string
	pathToProject        string
	syncStatusRepository SyncStatusRepository
	migrator             Migrator
}

func New(syncStatusRepository SyncStatusRepository, collectionRepository collection.Repository, pemKey string, appId int64, repositoryName string) (SyncProvider, error) {
	githubClient, err := github_fetcher.New(appId, []byte(pemKey))
	if err != nil {
		return nil, err
	}

	return &syncProvider{
		githubClient:         githubClient,
		repositoryName:       repositoryName,
		pathToProject:        "",
		syncStatusRepository: syncStatusRepository,
		migrator:             *NewMigrator(collectionRepository),
	}, nil
}

func (s *syncProvider) GetStatus(ctx context.Context) (Status, error) {
	return PROCESSING, nil
}

func (s *syncProvider) RegisterSyncJobs(cronService cron.CronService) error {
	ctx := context.Background()

	syncJob := cron.Job{
		Name:     fmt.Sprintf("sync-repo-%s", s.repositoryName),
		Schedule: "*/1 * * * *", // Every 5 minutes
		Function: func() error {
			ctx := context.Background()
			slog.Info("Start sync of repository", "repository", s.repositoryName)
			if err := s.SyncRepository(ctx); err != nil {
				slog.Error("Error syncing repository", "repository", s.repositoryName, "error", err)
				return err
			}
			return nil
		},
		Params: []any{},
	}

	if err := cronService.RegisterJob(ctx, syncJob); err != nil {
		return fmt.Errorf("failed to register sync job for repository %s: %w", s.repositoryName, err)
	}

	slog.Info("Successfully registered sync job for repository", "repository", s.repositoryName)
	return nil
}

func (s *syncProvider) markErrorAndReturn(repositoryName, commitSha string, err error, message string) error {
	if markErr := s.syncStatusRepository.MarkError(repositoryName, commitSha, err); markErr != nil {
		return fmt.Errorf("failed to mark error for repository %s: %w", repositoryName, markErr)
	}
	return fmt.Errorf(message+": %w", repositoryName, err)
}

func (s *syncProvider) SyncRepository(ctx context.Context) error {
	slog.Info("Starting sync for repository", "repository", s.repositoryName)

	// Fetch the latest files & commit from the repository
	contents, err := s.githubClient.GetLastCommit(ctx, s.repositoryName)

	if err != nil {
		return fmt.Errorf("failed to fetch latest files and commit for repository %s: %w", s.repositoryName, err)
	}
	// Get the last synced commit
	dbCommit, err := s.syncStatusRepository.GetLastSyncedCommit(s.repositoryName)
	if err != nil {
		return fmt.Errorf("failed to get last synced commit for repository %s: %w", s.repositoryName, err)
	}

	if dbCommit != nil && contents.Sha == dbCommit.Commit && dbCommit.IsActive == true {
		slog.Info("Repository is up to date, queueing next sync", "repository", s.repositoryName)
		return nil
	}

	// Create the sync status
	if err := s.syncStatusRepository.CreateIfNotExists(s.repositoryName, contents.Sha, contents.Message, contents.Date); err != nil {
		return fmt.Errorf("failed to create sync status for repository %s: %w", s.repositoryName, err)
	}

	// Get the manifest file from the repository contents
	manifest, err := s.githubClient.GetFileContent(ctx, s.repositoryName, contents.Sha, s.pathToProject+"mimsy.config.json")
	if err != nil {
		return s.markErrorAndReturn(s.repositoryName, contents.Sha, err, "failed to open manifest file for repository %s")
	}

	// With that manifest, unmarshall to the Schema:
	var config mimsy_schema.MimsyConfig
	if err := json.Unmarshal(manifest, &config); err != nil {
		return s.markErrorAndReturn(s.repositoryName, contents.Sha, err, "failed to unmarshal config file for repository %s")
	}

	var path string
	if config.SchemaPath != "" {
		path = config.SchemaPath
	} else if config.BasePath != "" {
		path = config.BasePath + "/mimsy.schema.json"
	} else {
		path = s.pathToProject + "mimsy.schema.json"
	}

	// We need to fetch the schema from the repository contents
	schema, err := s.githubClient.GetFileContent(ctx, s.repositoryName, contents.Sha, path)
	if err != nil {
		return s.markErrorAndReturn(s.repositoryName, contents.Sha, err, "failed to fetch schema for repository %s")
	}

	var schemaStruct mimsy_schema.Schema
	if err := json.Unmarshal(schema, &schemaStruct); err != nil {
		return s.markErrorAndReturn(s.repositoryName, contents.Sha, err, "failed to unmarshal schema file for repository %s")
	}

	// Get the last active migration to compare schemas
	activeMigration, err := s.syncStatusRepository.GetActiveMigration(s.repositoryName)
	if err != nil {
		return fmt.Errorf("failed to get last active migration for repository %s: %w", s.repositoryName, err)
	}

	// Check if the schemas are exactly the same
	if activeMigration != nil && activeMigration.Manifest != "" {
		var activeSchema mimsy_schema.Schema
		if err := json.Unmarshal([]byte(activeMigration.Manifest), &activeSchema); err == nil {
			// Compare schemas - if they're identical, mark as skipped
			currentSchemaBytes, _ := json.Marshal(schemaStruct)
			activeSchemaBytes, _ := json.Marshal(activeSchema)

			if string(currentSchemaBytes) == string(activeSchemaBytes) {
				slog.Info("Schema is identical to active migration, marking as skipped", "repository", s.repositoryName, "commit", contents.Sha)

				// Set the manifest and mark as skipped
				if err := s.syncStatusRepository.SetManifest(s.repositoryName, contents.Sha, schemaStruct); err != nil {
					return fmt.Errorf("failed to set manifest for repository %s: %w", s.repositoryName, err)
				}

				if err := s.syncStatusRepository.MarkAsSkipped(s.repositoryName, contents.Sha); err != nil {
					return fmt.Errorf("failed to mark as skipped for repository %s: %w", s.repositoryName, err)
				}

				slog.Info("Completed sync (skipped) for repository", "repository", s.repositoryName)
				return nil
			}
		}
	}

	// Schemas are different, proceed with normal migration
	// Set the manifest to the sync status repository
	if err := s.syncStatusRepository.SetManifest(s.repositoryName, contents.Sha, schemaStruct); err != nil {
		return fmt.Errorf("failed to set manifest for repository %s: %w", s.repositoryName, err)
	}

	// Generate the sql migration, and store it
	sqlSchema, err := s.migrator.GenerateSchema(ctx, &schemaStruct)
	if err != nil {
		return s.markErrorAndReturn(s.repositoryName, contents.Sha, err, "failed to generate sql migration for repository %s")
	}
	//Serialize the sql schema
	sqlSchemaBytes, err := json.Marshal(sqlSchema)
	if err != nil {
		return s.markErrorAndReturn(s.repositoryName, contents.Sha, err, "failed to serialize sql schema for repository %s")
	}

	if err := s.syncStatusRepository.SetAppliedMigration(s.repositoryName, contents.Sha, sqlSchemaBytes); err != nil {
		return fmt.Errorf("failed to store sql migration for repository %s: %w", s.repositoryName, err)
	}

	// Run the migration
	if err := s.migrator.Migrate(ctx, activeMigration, sqlSchema, contents.Message, contents.Sha); err != nil {
		return s.markErrorAndReturn(s.repositoryName, contents.Sha, err, "failed to run migration for repository %s")
	}

	// Store the sql migration inside of the

	// Mark the migration as active
	if err := s.syncStatusRepository.MarkAsActive(s.repositoryName, contents.Sha); err != nil {
		return fmt.Errorf("failed to mark migration as active for repository %s: %w", s.repositoryName, err)
	}

	// Generate the diff between the schema and the active migration
	slog.Info("Completed sync for repository", "repository", s.repositoryName)
	return nil
}
