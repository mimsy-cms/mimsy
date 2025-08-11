package github_fetcher

import (
	"archive/zip"
	"context"
	"fmt"
)

type RepositoryContents struct {
	Reader           *zip.Reader
	LatestCommitHash string
}

type GithubFetcher interface {
	// IsInstalled checks if the given app token is installed for the provided repository.
	IsInstalled(ctx context.Context, repository string) bool
	// GetRepositoryContents retrieves the contents of the specified repository, and exposes them as a zip.Reader.
	GetRepositoryContents(ctx context.Context, repository string) (*RepositoryContents, error)
}

type githubFetcher struct {
	appID       int64
	authManager *authManager
}

func New(appID int64, privateKeyPEM []byte) (GithubFetcher, error) {
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

	return &githubFetcher{
		appID:       appID,
		authManager: authManager,
	}, nil
}

func (f *githubFetcher) IsInstalled(ctx context.Context, repository string) bool {
	owner, repo, err := parseRepository(repository)
	if err != nil {
		return false
	}

	_, err = f.authManager.getInstallationID(ctx, owner, repo)
	return err == nil
}

func (f *githubFetcher) GetRepositoryContents(ctx context.Context, repository string) (*RepositoryContents, error) {
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

	// Fetch the latest commit
	latestCommit, err := fetchLatestCommit(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest commit: %w", err)
	}

	zipReader, err := downloadZipArchive(ctx, client, owner, repo, latestCommit)
	if err != nil {
		return nil, fmt.Errorf("failed to download repository: %w", err)
	}

	return &RepositoryContents{
		Reader:           zipReader,
		LatestCommitHash: latestCommit,
	}, nil
}
