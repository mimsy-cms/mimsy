package main

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mimsy-cms/mimsy/pkg/github_fetcher"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <pem_path> <app_id> <repo_path>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s /path/to/key.pem 123456 owner/repo\n", os.Args[0])
		os.Exit(1)
	}

	pemPath := os.Args[1]
	appIDStr := os.Args[2]
	repoPath := os.Args[3]

	// Parse app ID
	appID, err := strconv.ParseInt(appIDStr, 10, 64)
	if err != nil {
		log.Fatalf("Invalid app ID '%s': %v", appIDStr, err)
	}

	// Read PEM file
	pemData, err := os.ReadFile(pemPath)
	if err != nil {
		log.Fatalf("Failed to read PEM file '%s': %v", pemPath, err)
	}

	// Create github fetcher
	fetcher, err := github_fetcher.New(appID, pemData)
	if err != nil {
		log.Fatalf("Failed to create github fetcher: %v", err)
	}

	ctx := context.Background()

	// Check if the app is installed for the repository
	fmt.Printf("Checking if app is installed for repository %s...\n", repoPath)
	isInstalled := fetcher.IsInstalled(ctx, repoPath)
	if !isInstalled {
		log.Fatalf("GitHub App is not installed for repository %s", repoPath)
	}
	fmt.Printf("✓ App is installed for repository %s\n", repoPath)

	// Fetch repository contents
	fmt.Printf("Fetching repository contents for %s...\n", repoPath)
	contents, err := fetcher.GetRepositoryContents(ctx, repoPath)
	if err != nil {
		log.Fatalf("Failed to fetch repository contents: %v", err)
	}
	fmt.Printf("✓ Successfully fetched repository contents (commit: %s)\n", contents.LatestCommitHash)

	// Look for README files
	readmeFiles := findReadmeFiles(contents)
	if len(readmeFiles) == 0 {
		fmt.Println("No README files found in the repository")
		return
	}

	fmt.Printf("Found %d README file(s):\n", len(readmeFiles))
	for _, file := range readmeFiles {
		fmt.Printf("  - %s\n", file.Name)
	}

	// Display the first README file
	readme := readmeFiles[0]
	fmt.Printf("\n=== Content of %s ===\n", readme.Name)
	
	rc, err := readme.Open()
	if err != nil {
		log.Fatalf("Failed to open README file: %v", err)
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		log.Fatalf("Failed to read README content: %v", err)
	}

	fmt.Println(string(content))
}

func findReadmeFiles(contents *github_fetcher.RepositoryContents) []*zip.File {
	var readmeFiles []*zip.File
	
	for _, file := range contents.Reader.File {
		fileName := filepath.Base(file.Name)
		lowerName := strings.ToLower(fileName)
		
		// Check for README files (with common extensions)
		if strings.HasPrefix(lowerName, "readme") {
			readmeFiles = append(readmeFiles, file)
		}
	}
	
	return readmeFiles
}