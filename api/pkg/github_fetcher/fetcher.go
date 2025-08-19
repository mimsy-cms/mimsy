package github_fetcher

import (
	"archive/zip"
	"context"
	"fmt"
)

type RepositoryContents struct {
	Reader       *zip.Reader
	LatestCommit Commit
}

type GithubProvider interface {
	// IsInstalled checks if the given app token is installed for the provided repository.
	IsInstalled(ctx context.Context, repository string) bool
	// GetLastCommit retrieves the latest commit from the specified repository.
	GetLastCommit(ctx context.Context, repository string) (*Commit, error)
	// GetContents retrieves the contents of the specified repository, and exposes them as a zip.Reader.
	GetContents(ctx context.Context, repository, ref string) (*zip.Reader, error)
	// GetRepositoryContents retrieves the contents of the specified repository, and exposes them as a zip.Reader.
	GetRepositoryContents(ctx context.Context, repository string) (*RepositoryContents, error)
	// GetFileContent retrieves the content of a specific file from the repository at the given ref and path.
	GetFileContent(ctx context.Context, repository, ref, path string) ([]byte, error)
	// CreateCommitStatus creates a commit status for the specified repository and commit.
	CreateCommitStatus(ctx context.Context, repository, commitSHA, state, description, targetURL string) error
}

type githubProvider struct {
	appID       int64
	authManager *authManager
}

func New(appID int64, privateKeyPEM []byte) (GithubProvider, error) {
	if appID <= 0 {
		return nil, fmt.Errorf("invalid app ID: must be positive")
	}

	if len(privateKeyPEM) == 0 {
		return nil, fmt.Errorf("private key cannot be empty")
	}

	authManager, err := newAuthManager(appID, privateKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth manager: %w", err)
	}

	return &githubProvider{
		appID:       appID,
		authManager: authManager,
	}, nil
}

func (f *githubProvider) IsInstalled(ctx context.Context, repository string) bool {
	owner, repo, err := parseRepository(repository)
	if err != nil {
		return false
	}

	_, err = f.authManager.getInstallationID(ctx, owner, repo)
	return err == nil
}

func (f *githubProvider) GetLastCommit(ctx context.Context, repository string) (*Commit, error) {
	owner, repo, err := parseRepository(repository)
	if err != nil {
		return nil, err
	}

	installationID, err := f.authManager.getInstallationID(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation ID: %w", err)
	}

	token, err := f.authManager.getInstallationToken(ctx, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation token: %w", err)
	}

	client := createInstallationClient(ctx, token)

	latestCommit, err := fetchLatestCommit(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest commit: %w", err)
	}

	return latestCommit, nil
}

func (f *githubProvider) GetContents(ctx context.Context, repository, ref string) (*zip.Reader, error) {
	owner, repo, err := parseRepository(repository)
	if err != nil {
		return nil, err
	}

	installationID, err := f.authManager.getInstallationID(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation ID: %w", err)
	}

	token, err := f.authManager.getInstallationToken(ctx, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation token: %w", err)
	}

	client := createInstallationClient(ctx, token)

	zipReader, err := downloadZipArchive(ctx, client, owner, repo, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to download repository: %w", err)
	}

	return zipReader, nil
}

func (f *githubProvider) GetRepositoryContents(ctx context.Context, repository string) (*RepositoryContents, error) {
	latestCommit, err := f.GetLastCommit(ctx, repository)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest commit: %w", err)
	}

	zipReader, err := f.GetContents(ctx, repository, latestCommit.Sha)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository contents: %w", err)
	}

	return &RepositoryContents{
		Reader:       zipReader,
		LatestCommit: *latestCommit,
	}, nil
}

func (f *githubProvider) GetFileContent(ctx context.Context, repository, ref, path string) ([]byte, error) {
	owner, repo, err := parseRepository(repository)
	if err != nil {
		return nil, err
	}

	installationID, err := f.authManager.getInstallationID(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation ID: %w", err)
	}

	token, err := f.authManager.getInstallationToken(ctx, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation token: %w", err)
	}

	client := createInstallationClient(ctx, token)

	return getFileContent(ctx, client, owner, repo, ref, path)
}

func (f *githubProvider) CreateCommitStatus(ctx context.Context, repository, commitSHA, state, description, targetURL string) error {
	owner, repo, err := parseRepository(repository)
	if err != nil {
		return err
	}

	installationID, err := f.authManager.getInstallationID(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get installation ID: %w", err)
	}

	token, err := f.authManager.getInstallationToken(ctx, installationID)
	if err != nil {
		return fmt.Errorf("failed to get installation token: %w", err)
	}

	client := createInstallationClient(ctx, token)

	return CreateCommitStatus(ctx, client, repository, commitSHA, state, description, targetURL)
}
