package github_fetcher

import "archive/zip"

type GithubFetcher interface {
	// IsInstalled checks if the given app token is installed for the provided repository.
	IsInstalled(repository string) bool
	// GetRepositoryContents retrieves the contents of the specified repository, and exposes them as a zip.Reader.
	GetRepositoryContents(repository string) (zip.Reader, error)
}
