package github_fetcher

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/go-github/v74/github"
)

func TestNew(t *testing.T) {
	privateKey, err := os.ReadFile("testdata/test_key.pem")
	if err != nil {
		t.Fatalf("failed to read test key: %v", err)
	}

	tests := []struct {
		name      string
		appID     int64
		key       []byte
		wantError bool
	}{
		{
			name:      "valid configuration",
			appID:     12345,
			key:       privateKey,
			wantError: false,
		},
		{
			name:      "invalid app ID",
			appID:     0,
			key:       privateKey,
			wantError: true,
		},
		{
			name:      "negative app ID",
			appID:     -1,
			key:       privateKey,
			wantError: true,
		},
		{
			name:      "empty private key",
			appID:     12345,
			key:       []byte{},
			wantError: true,
		},
		{
			name:      "invalid private key",
			appID:     12345,
			key:       []byte("not a valid key"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcher, err := New(tt.appID, tt.key)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if fetcher != nil {
					t.Errorf("expected nil fetcher on error, got %v", fetcher)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if fetcher == nil {
					t.Errorf("expected fetcher, got nil")
				}
			}
		})
	}
}

func TestParseRepository(t *testing.T) {
	tests := []struct {
		name       string
		repository string
		wantOwner  string
		wantRepo   string
		wantError  bool
	}{
		{
			name:       "valid repository",
			repository: "owner/repo",
			wantOwner:  "owner",
			wantRepo:   "repo",
			wantError:  false,
		},
		{
			name:       "repository with spaces",
			repository: " owner / repo ",
			wantOwner:  "owner",
			wantRepo:   "repo",
			wantError:  false,
		},
		{
			name:       "missing slash",
			repository: "ownerrepo",
			wantError:  true,
		},
		{
			name:       "too many slashes",
			repository: "owner/repo/extra",
			wantError:  true,
		},
		{
			name:       "empty owner",
			repository: "/repo",
			wantError:  true,
		},
		{
			name:       "empty repo",
			repository: "owner/",
			wantError:  true,
		},
		{
			name:       "empty string",
			repository: "",
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := parseRepository(tt.repository)
			if tt.wantError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if owner != tt.wantOwner {
					t.Errorf("owner = %q, want %q", owner, tt.wantOwner)
				}
				if repo != tt.wantRepo {
					t.Errorf("repo = %q, want %q", repo, tt.wantRepo)
				}
			}
		})
	}
}

func TestGithubFetcherMocked(t *testing.T) {
	privateKey, err := os.ReadFile("testdata/test_key.pem")
	if err != nil {
		t.Fatalf("failed to read test key: %v", err)
	}

	installationID := int64(987654)
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/test-owner/test-repo/installation":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         installationID,
				"account":    map[string]string{"login": "test-owner"},
			})

		case fmt.Sprintf("/app/installations/%d/access_tokens", installationID):
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"token":      "test-installation-token",
				"expires_at": time.Now().Add(time.Hour).Format(time.RFC3339),
			})

		case "/repos/test-owner/test-repo/zipball":
			buf := new(bytes.Buffer)
			zipWriter := zip.NewWriter(buf)
			
			file1, _ := zipWriter.Create("README.md")
			file1.Write([]byte("# Test Repository\n"))
			
			file2, _ := zipWriter.Create("main.go")
			file2.Write([]byte("package main\n\nfunc main() {}\n"))
			
			zipWriter.Close()
			
			w.Header().Set("Content-Type", "application/zip")
			w.Write(buf.Bytes())

		case "/repos/test-owner/not-installed/installation":
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Not Found",
			})

		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	fetcher, err := New(12345, privateKey)
	if err != nil {
		t.Fatalf("failed to create fetcher: %v", err)
	}

	t.Run("IsInstalled - invalid repository format", func(t *testing.T) {
		installed := fetcher.IsInstalled(context.Background(), "invalid-format")
		if installed {
			t.Errorf("expected invalid repository format to return false")
		}
	})

	t.Run("GetRepositoryContents - invalid repository format", func(t *testing.T) {
		_, err := fetcher.GetRepositoryContents(context.Background(), "invalid-format")
		if err == nil {
			t.Error("expected error for invalid repository format")
		}
	})
}

func TestAuthManager(t *testing.T) {
	privateKey, err := os.ReadFile("testdata/test_key.pem")
	if err != nil {
		t.Fatalf("failed to read test key: %v", err)
	}

	t.Run("JWT generation", func(t *testing.T) {
		am, err := newAuthManager(12345, privateKey)
		if err != nil {
			t.Fatalf("failed to create auth manager: %v", err)
		}

		jwt, err := am.generateJWT()
		if err != nil {
			t.Fatalf("failed to generate JWT: %v", err)
		}

		if jwt == "" {
			t.Error("expected non-empty JWT")
		}
	})

	t.Run("Installation token caching", func(t *testing.T) {
		am, err := newAuthManager(12345, privateKey)
		if err != nil {
			t.Fatalf("failed to create auth manager: %v", err)
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/app/installations/123/access_tokens" {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]interface{}{
					"token":      "cached-token",
					"expires_at": time.Now().Add(time.Hour).Format(time.RFC3339),
				})
			}
		}))
		defer server.Close()

		am.installationCache[123] = &installationToken{
			token:     "cached-token",
			expiresAt: time.Now().Add(time.Hour),
		}

		token, err := am.getInstallationToken(context.Background(), 123)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if token != "cached-token" {
			t.Errorf("expected cached token, got %q", token)
		}
	})
}

func TestDownloadZipArchive(t *testing.T) {
	var testServerURL string
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/owner/repo/zipball" {
			w.Header().Set("Location", testServerURL+"/download/archive.zip")
			w.WriteHeader(http.StatusFound)
		} else if r.URL.Path == "/download/archive.zip" {
			buf := new(bytes.Buffer)
			zipWriter := zip.NewWriter(buf)
			
			file, _ := zipWriter.Create("test.txt")
			file.Write([]byte("test content"))
			
			zipWriter.Close()
			
			w.Header().Set("Content-Type", "application/zip")
			w.Write(buf.Bytes())
		}
	}))
	defer testServer.Close()
	testServerURL = testServer.URL

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}
	
	baseURL, _ := url.Parse(testServer.URL + "/")
	githubClient := github.NewClient(client)
	githubClient.BaseURL = baseURL

	zipReader, err := downloadZipArchive(context.Background(), githubClient, "owner", "repo", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if zipReader == nil {
		t.Fatal("expected zip reader, got nil")
	}

	if len(zipReader.File) != 1 {
		t.Errorf("expected 1 file in zip, got %d", len(zipReader.File))
	}

	if zipReader.File[0].Name != "test.txt" {
		t.Errorf("expected file name 'test.txt', got %q", zipReader.File[0].Name)
	}
}