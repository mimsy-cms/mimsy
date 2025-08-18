package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

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
}

func New(syncStatusRepository SyncStatusRepository, pemKey string, appId int64, repositoryName string) (SyncProvider, error) {
	githubClient, err := github_fetcher.New(appId, []byte(pemKey))
	if err != nil {
		return nil, err
	}

	return &syncProvider{
		githubClient:         githubClient,
		repositoryName:       repositoryName,
		pathToProject:        "",
		syncStatusRepository: syncStatusRepository,
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

	if dbCommit != nil && contents.Sha == dbCommit.Commit {
		slog.Info("Repository is up to date, queueing next sync", "repository", s.repositoryName)
		return nil
	}

	// Create the sync status
	if err := s.syncStatusRepository.CreateStatus(s.repositoryName, contents.Sha, contents.Message, contents.Date); err != nil {
		return fmt.Errorf("failed to create sync status for repository %s: %w", s.repositoryName, err)
	}

	// Get the manifest file from the repository contents
	manifest, err := s.githubClient.GetFileContent(ctx, s.repositoryName, contents.Sha, s.pathToProject+"mimsy.config.json")
	if err != nil {
		// Mark the error to the sync status repository
		if err := s.syncStatusRepository.MarkError(s.repositoryName, contents.Sha, err); err != nil {
			return fmt.Errorf("failed to mark error for repository %s: %w", s.repositoryName, err)
		}
		return fmt.Errorf("failed to open manifest file for repository %s: %w", s.repositoryName, err)
	}

	// With that manifest, unmarshall to the Schema:
	var config mimsy_schema.MimsyConfig
	if err := json.Unmarshal(manifest, &config); err != nil {
		return fmt.Errorf("failed to unmarshal config file for repository %s: %w", s.repositoryName, err)
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
		return fmt.Errorf("failed to fetch schema for repository %s: %w", s.repositoryName, err)
	}

	var schemaStruct mimsy_schema.Schema
	if err := json.Unmarshal(schema, &schemaStruct); err != nil {
		return fmt.Errorf("failed to unmarshal schema file for repository %s: %w", s.repositoryName, err)
	}

	// once unmarshalled, set the manifest to the sync status repository
	if err := s.syncStatusRepository.SetManifest(s.repositoryName, contents.Sha, schemaStruct); err != nil {
		return fmt.Errorf("failed to set manifest for repository %s: %w", s.repositoryName, err)
	}

	slog.Info("Completed sync for repository", "repository", s.repositoryName)
	return nil
}
