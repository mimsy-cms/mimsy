package github_fetcher

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/go-github/v74/github"
)

func parseRepository(repository string) (owner, repo string, err error) {
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: expected 'owner/repo', got '%s'", repository)
	}

	owner = strings.TrimSpace(parts[0])
	repo = strings.TrimSpace(parts[1])

	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("invalid repository format: owner and repo cannot be empty")
	}

	return owner, repo, nil
}

func downloadZipArchive(ctx context.Context, client *github.Client, owner, repo string, ref string) (*zip.Reader, error) {
	opts := &github.RepositoryContentGetOptions{}
	if ref != "" {
		opts.Ref = ref
	}

	url, _, err := client.Repositories.GetArchiveLink(ctx, owner, repo, github.Zipball, opts, 3)
	if err != nil {
		return nil, fmt.Errorf("failed to get archive link: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	transport := client.Client().Transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	httpClient := &http.Client{Transport: transport}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download archive: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read archive data: %w", err)
	}

	reader := bytes.NewReader(data)
	zipReader, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to create zip reader: %w", err)
	}

	return zipReader, nil
}

func fetchLatestCommit(ctx context.Context, client *github.Client, owner, repo string) (string, error) {
	// Get the default branch first
	repository, _, err := client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return "", fmt.Errorf("failed to get repository info: %w", err)
	}

	defaultBranch := repository.GetDefaultBranch()
	if defaultBranch == "" {
		defaultBranch = "main"
	}

	// Get the latest commit from the default branch
	commits, _, err := client.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{
		SHA:         defaultBranch,
		ListOptions: github.ListOptions{PerPage: 1},
	})
	if err != nil {
		return "", fmt.Errorf("failed to list commits: %w", err)
	}

	if len(commits) == 0 {
		return "", fmt.Errorf("no commits found in repository")
	}

	return commits[0].GetSHA(), nil
}

func createInstallationClient(_ context.Context, token string) *github.Client {
	httpClient := &http.Client{
		Transport: &tokenTransport{
			token: token,
			base:  http.DefaultTransport,
		},
	}
	return github.NewClient(httpClient)
}

type tokenTransport struct {
	token string
	base  http.RoundTripper
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "token "+t.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	return t.base.RoundTrip(req)
}

// CreateCommitStatus creates a commit status for the specified repository and commit
func CreateCommitStatus(ctx context.Context, client *github.Client, repository, commitSHA, state, description, targetURL string) error {
	owner, repo, err := parseRepository(repository)
	if err != nil {
		return fmt.Errorf("failed to parse repository: %w", err)
	}

	status := &github.RepoStatus{
		State:       &state,
		Description: &description,
		Context:     github.String("mimsy/db-migration"),
	}

	if targetURL != "" {
		status.TargetURL = &targetURL
	}

	_, _, err = client.Repositories.CreateStatus(ctx, owner, repo, commitSHA, status)
	if err != nil {
		return fmt.Errorf("failed to create commit status: %w", err)
	}

	return nil
}
