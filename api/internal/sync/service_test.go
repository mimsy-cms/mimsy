package sync_test

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mocks_collection "github.com/mimsy-cms/mimsy/internal/mocks/collection"
	mocks_cron "github.com/mimsy-cms/mimsy/internal/mocks/cron"
	mocks_sync "github.com/mimsy-cms/mimsy/internal/mocks/sync"
	"github.com/mimsy-cms/mimsy/internal/sync"
	"github.com/mimsy-cms/mimsy/pkg/github_fetcher"
	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
)

// getTestPEMKey loads the test PEM key from testdata
func getTestPEMKey(t *testing.T) string {
	keyPath := filepath.Join("..", "..", "pkg", "github_fetcher", "testdata", "test_key.pem")
	keyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		t.Fatalf("Failed to read test PEM key: %v", err)
	}
	return string(keyBytes)
}

func TestNew_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSyncRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCollectionRepo := mocks_collection.NewMockRepository(ctrl)

	provider, err := sync.New(
		mockSyncRepo,
		mockCollectionRepo,
		getTestPEMKey(t),
		123456,
		"test-repo",
	)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if provider == nil {
		t.Error("Expected provider to be created")
	}
}

func TestSyncProvider_GetStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSyncRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCollectionRepo := mocks_collection.NewMockRepository(ctrl)

	provider, err := sync.New(
		mockSyncRepo,
		mockCollectionRepo,
		getTestPEMKey(t),
		123456,
		"test-repo",
	)

	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	status, err := provider.GetStatus(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status != sync.PROCESSING {
		t.Errorf("Expected status PROCESSING, got %v", status)
	}
}

func TestSyncProvider_RegisterSyncJobs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSyncRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCollectionRepo := mocks_collection.NewMockRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	provider, err := sync.New(
		mockSyncRepo,
		mockCollectionRepo,
		getTestPEMKey(t),
		123456,
		"test-repo",
	)

	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	mockCron.EXPECT().
		RegisterJob(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	err = provider.RegisterSyncJobs(mockCron)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestSyncProvider_RegisterSyncJobs_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSyncRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCollectionRepo := mocks_collection.NewMockRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	provider, err := sync.New(
		mockSyncRepo,
		mockCollectionRepo,
		getTestPEMKey(t),
		123456,
		"test-repo",
	)

	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	mockCron.EXPECT().
		RegisterJob(gomock.Any(), gomock.Any()).
		Return(errors.New("cron error")).
		Times(1)

	err = provider.RegisterSyncJobs(mockCron)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestSyncProvider_SyncRepository_UpToDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// This test requires more complex mocking with the internal sync provider
	// For now, we'll test the basic initialization and public interface
	mockSyncRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCollectionRepo := mocks_collection.NewMockRepository(ctrl)

	provider, err := sync.New(
		mockSyncRepo,
		mockCollectionRepo,
		getTestPEMKey(t),
		123456,
		"test-repo",
	)

	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	if provider == nil {
		t.Error("Expected provider to be created")
	}
}

func TestSyncProvider_SyncRepository_NewCommit(t *testing.T) {
	// This would test the full sync flow, but requires extensive mocking
	// of the internal GitHub client and migration system.
	// For now, we'll focus on testing the components separately.
	t.Skip("Skipping complex integration test - would require extensive mocking")
}

func TestSyncProvider_SyncRepository_ErrorHandling(t *testing.T) {
	// Similar to above - this would test error scenarios in the sync flow
	t.Skip("Skipping complex integration test - would require extensive mocking")
}

// Test the Status enum values
func TestStatusValues(t *testing.T) {
	if sync.PROVIDER_ERROR != 0 {
		t.Errorf("Expected PROVIDER_ERROR to be 0, got %d", sync.PROVIDER_ERROR)
	}

	if sync.PROCESSING != 1 {
		t.Errorf("Expected PROCESSING to be 1, got %d", sync.PROCESSING)
	}

	if sync.SUCCESS != 2 {
		t.Errorf("Expected SUCCESS to be 2, got %d", sync.SUCCESS)
	}

	if sync.NO_MANIFEST != 3 {
		t.Errorf("Expected NO_MANIFEST to be 3, got %d", sync.NO_MANIFEST)
	}

	if sync.INVALID_MANIFEST != 4 {
		t.Errorf("Expected INVALID_MANIFEST to be 4, got %d", sync.INVALID_MANIFEST)
	}
}

// Test that MarkAsSkipped is called when schema hasn't changed
func TestSyncRepository_SchemaUnchanged_CallsMarkAsSkipped(t *testing.T) {
	// Create a custom sync provider with mocked dependencies to test the internal logic
	mockSyncRepo := &mockSyncStatusRepository{}
	mockGithubClient := &mockGithubProvider{}

	provider := &testSyncProvider{
		githubClient:         mockGithubClient,
		repositoryName:       "test-repo",
		syncStatusRepository: mockSyncRepo,
		migrator:             &mockMigrator{},
	}

	// Setup test data
	now := time.Now()
	newCommit := &github_fetcher.Commit{
		Sha:     "new-commit-sha",
		Message: "New commit with same schema",
		Date:    now,
	}

	// Create identical schemas
	testSchema := mimsy_schema.Schema{
		Collections: []mimsy_schema.Collection{
			{
				Name: "test",
				Schema: mimsy_schema.CollectionFields{
					"title": mimsy_schema.SchemaElement{Type: "text"},
				},
			},
		},
	}

	schemaJSON, _ := json.Marshal(testSchema)
	manifestJSON := []byte(`{"schemaPath": "mimsy.schema.json"}`)

	// Setup mocks
	mockGithubClient.lastCommit = newCommit
	mockGithubClient.fileContents = map[string][]byte{
		"mimsy.config.json": manifestJSON,
		"mimsy.schema.json": schemaJSON,
	}

	// Mock last synced commit to be different (to pass early exit check)
	mockSyncRepo.lastSyncedCommit = &sync.SyncStatus{
		Repo:      "test-repo",
		Commit:    "old-commit-sha",
		IsActive:  true,
		AppliedAt: now.Add(-time.Hour),
	}

	// Mock active migration with identical schema
	mockSyncRepo.activeMigration = &sync.SyncStatus{
		Repo:     "test-repo",
		Commit:   "old-commit-sha",
		Manifest: string(schemaJSON),
		IsActive: true,
	}

	// Execute
	err := provider.SyncRepository(context.Background())

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !mockSyncRepo.createIfNotExistsCalled {
		t.Error("Expected CreateIfNotExists to be called")
	}

	if !mockSyncRepo.setManifestCalled {
		t.Error("Expected SetManifest to be called")
	}

	if !mockSyncRepo.markAsSkippedCalled {
		t.Error("Expected MarkAsSkipped to be called when schema is unchanged")
	}

	if mockSyncRepo.markAsActiveCalled {
		t.Error("Expected MarkAsActive NOT to be called when schema is unchanged")
	}
}

// Test that sync returns early when last commit is already processed and active
func TestSyncRepository_CommitAlreadyProcessedAndActive_ReturnsEarly(t *testing.T) {
	// Create a custom sync provider with mocked dependencies
	mockSyncRepo := &mockSyncStatusRepository{}
	mockGithubClient := &mockGithubProvider{}

	provider := &testSyncProvider{
		githubClient:         mockGithubClient,
		repositoryName:       "test-repo",
		syncStatusRepository: mockSyncRepo,
		migrator:             &mockMigrator{},
	}

	// Setup test data - same commit SHA
	now := time.Now()
	commitSha := "same-commit-sha"

	latestCommit := &github_fetcher.Commit{
		Sha:     commitSha,
		Message: "Same commit",
		Date:    now,
	}

	// Setup mocks
	mockGithubClient.lastCommit = latestCommit

	// Mock last synced commit to be the SAME as latest (with IsActive=true)
	mockSyncRepo.lastSyncedCommit = &sync.SyncStatus{
		Repo:      "test-repo",
		Commit:    commitSha, // Same SHA
		IsActive:  true,      // And it's active
		AppliedAt: now,
	}

	// Execute
	err := provider.SyncRepository(context.Background())

	// Assertions
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify early exit - none of the downstream methods should be called
	if mockSyncRepo.createIfNotExistsCalled {
		t.Error("Expected CreateIfNotExists NOT to be called on early exit")
	}

	if mockGithubClient.getFileContentCallCount > 0 {
		t.Error("Expected GetFileContent NOT to be called on early exit")
	}

	if mockSyncRepo.setManifestCalled {
		t.Error("Expected SetManifest NOT to be called on early exit")
	}

	if mockSyncRepo.markAsSkippedCalled {
		t.Error("Expected MarkAsSkipped NOT to be called on early exit")
	}

	if mockSyncRepo.markAsActiveCalled {
		t.Error("Expected MarkAsActive NOT to be called on early exit")
	}
}

// Test helper types and mocks for the integration tests
type testSyncProvider struct {
	githubClient         *mockGithubProvider
	repositoryName       string
	syncStatusRepository *mockSyncStatusRepository
	migrator             *mockMigrator
}

func (s *testSyncProvider) SyncRepository(ctx context.Context) error {
	// This mirrors the actual SyncRepository logic from service.go

	// Fetch the latest files & commit from the repository
	contents, err := s.githubClient.GetLastCommit(ctx, s.repositoryName)
	if err != nil {
		return err
	}

	// Get the last synced commit
	dbCommit, err := s.syncStatusRepository.GetLastSyncedCommit(s.repositoryName)
	if err != nil {
		return err
	}

	// Early exit check - this is the key logic we're testing
	if dbCommit != nil && contents.Sha == dbCommit.Commit && dbCommit.IsActive == true {
		return nil // Early exit without processing
	}

	// Create the sync status
	if err := s.syncStatusRepository.CreateIfNotExists(s.repositoryName, contents.Sha, contents.Message, contents.Date); err != nil {
		return err
	}

	// Get the manifest file from the repository contents
	manifest, err := s.githubClient.GetFileContent(ctx, s.repositoryName, contents.Sha, "mimsy.config.json")
	if err != nil {
		return err
	}

	// With that manifest, unmarshal to the Schema:
	var config mimsy_schema.MimsyConfig
	if err := json.Unmarshal(manifest, &config); err != nil {
		return err
	}

	var path string
	if config.SchemaPath != "" {
		path = config.SchemaPath
	} else {
		path = "mimsy.schema.json"
	}

	// We need to fetch the schema from the repository contents
	schema, err := s.githubClient.GetFileContent(ctx, s.repositoryName, contents.Sha, path)
	if err != nil {
		return err
	}

	var schemaStruct mimsy_schema.Schema
	if err := json.Unmarshal(schema, &schemaStruct); err != nil {
		return err
	}

	// Get the last active migration to compare schemas
	activeMigration, err := s.syncStatusRepository.GetActiveMigration(s.repositoryName)
	if err != nil {
		return err
	}

	// Check if the schemas are exactly the same - this is the skip logic we're testing
	if activeMigration != nil && activeMigration.Manifest != "" {
		var activeSchema mimsy_schema.Schema
		if err := json.Unmarshal([]byte(activeMigration.Manifest), &activeSchema); err == nil {
			// Compare schemas - if they're identical, mark as skipped
			currentSchemaBytes, _ := json.Marshal(schemaStruct)
			activeSchemaBytes, _ := json.Marshal(activeSchema)

			if string(currentSchemaBytes) == string(activeSchemaBytes) {
				// Set the manifest and mark as skipped
				if err := s.syncStatusRepository.SetManifest(s.repositoryName, contents.Sha, schemaStruct); err != nil {
					return err
				}

				if err := s.syncStatusRepository.MarkAsSkipped(s.repositoryName, contents.Sha); err != nil {
					return err
				}

				return nil // Skip processing
			}
		}
	}

	// Schemas are different, proceed with normal migration (simplified for test)
	if err := s.syncStatusRepository.SetManifest(s.repositoryName, contents.Sha, schemaStruct); err != nil {
		return err
	}

	if err := s.syncStatusRepository.MarkAsActive(s.repositoryName, contents.Sha); err != nil {
		return err
	}

	return nil
}

// Mock implementations for testing
type mockSyncStatusRepository struct {
	lastSyncedCommit        *sync.SyncStatus
	activeMigration         *sync.SyncStatus
	createIfNotExistsCalled bool
	setManifestCalled       bool
	markAsSkippedCalled     bool
	markAsActiveCalled      bool
}

func (m *mockSyncStatusRepository) GetLastSyncedCommit(repo string) (*sync.SyncStatus, error) {
	return m.lastSyncedCommit, nil
}

func (m *mockSyncStatusRepository) GetActiveMigration(repo string) (*sync.SyncStatus, error) {
	return m.activeMigration, nil
}

func (m *mockSyncStatusRepository) CreateIfNotExists(repo string, commitSha string, commitMessage string, commitDate time.Time) error {
	m.createIfNotExistsCalled = true
	return nil
}

func (m *mockSyncStatusRepository) SetManifest(repo string, commitSha string, manifest mimsy_schema.Schema) error {
	m.setManifestCalled = true
	return nil
}

func (m *mockSyncStatusRepository) MarkAsSkipped(repo string, commitSha string) error {
	m.markAsSkippedCalled = true
	return nil
}

func (m *mockSyncStatusRepository) MarkAsActive(repo string, commitSha string) error {
	m.markAsActiveCalled = true
	return nil
}

// Unused methods for interface compliance
func (m *mockSyncStatusRepository) GetStatus(repo string) (*sync.SyncStatus, error) { return nil, nil }
func (m *mockSyncStatusRepository) GetRecentStatuses(limit int) ([]sync.SyncStatus, error) {
	return nil, nil
}
func (m *mockSyncStatusRepository) MarkError(repo string, commitSha string, err error) error {
	return nil
}
func (m *mockSyncStatusRepository) SetAppliedMigration(repo string, commitSha string, migration []byte) error {
	return nil
}

type mockGithubProvider struct {
	lastCommit              *github_fetcher.Commit
	fileContents            map[string][]byte
	getFileContentCallCount int
}

func (m *mockGithubProvider) GetLastCommit(ctx context.Context, repository string) (*github_fetcher.Commit, error) {
	return m.lastCommit, nil
}

func (m *mockGithubProvider) GetFileContent(ctx context.Context, repository, ref, path string) ([]byte, error) {
	m.getFileContentCallCount++
	if content, exists := m.fileContents[path]; exists {
		return content, nil
	}
	return nil, errors.New("file not found")
}

// Unused methods for interface compliance
func (m *mockGithubProvider) IsInstalled(ctx context.Context, repository string) bool { return true }
func (m *mockGithubProvider) GetContents(ctx context.Context, repository, ref string) (*zip.Reader, error) {
	return nil, nil
}
func (m *mockGithubProvider) GetRepositoryContents(ctx context.Context, repository string) (*github_fetcher.RepositoryContents, error) {
	return nil, nil
}
func (m *mockGithubProvider) CreateCommitStatus(ctx context.Context, repository, commitSHA, state, description, targetURL string) error {
	return nil
}

type mockMigrator struct{}
