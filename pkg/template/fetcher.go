// ABOUTME: Template fetcher for downloading and caching bingo project from GitHub
// ABOUTME: Handles tarball download, extraction, and local caching with file locking
package template

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	githubArchiveBase = "https://github.com/bingo-project/bingo/archive"
	defaultTimeout    = 30 * time.Second
)

// Fetcher handles template downloading and caching
type Fetcher struct {
	cacheDir string        // ~/.bingoctl/templates
	timeout  time.Duration // 30s
	mirror   string        // mirror address from env var
}

// NewFetcher creates a new Fetcher instance
func NewFetcher() *Fetcher {
	// TODO: get cache dir from user home
	// TODO: read mirror from env BINGOCTL_TEMPLATE_MIRROR
	return &Fetcher{
		cacheDir: "", // will implement later
		timeout:  defaultTimeout,
		mirror:   "", // will implement later
	}
}

// buildDownloadURL constructs download URL (supports mirror)
// Examples:
//   - tag: https://github.com/.../archive/refs/tags/v1.2.3.tar.gz
//   - branch: https://github.com/.../archive/refs/heads/main.tar.gz
//   - commit: https://github.com/.../archive/{hash}.tar.gz
func (f *Fetcher) buildDownloadURL(ref string) string {
	var url string

	refKind := refType(ref)
	switch refKind {
	case "tag":
		url = fmt.Sprintf("%s/refs/tags/%s.tar.gz", githubArchiveBase, ref)
	case "branch":
		url = fmt.Sprintf("%s/refs/heads/%s.tar.gz", githubArchiveBase, ref)
	case "commit":
		url = fmt.Sprintf("%s/%s.tar.gz", githubArchiveBase, ref)
	}

	if f.mirror != "" {
		return f.mirror + url
	}

	return url
}

// downloadWithTimeout downloads tarball with 30s timeout and shows progress bar
func (f *Fetcher) downloadWithTimeout(url string) (string, error) {
	// Create HTTP client with timeout
	ctx, cancel := context.WithTimeout(context.Background(), f.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create temporary file
	tmpFile, err := os.CreateTemp(f.cacheDir, "bingoctl-*.tar.gz")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	// Download to temp file with progress bar
	if resp.ContentLength > 0 {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"Downloading",
		)
		if _, err := io.Copy(io.MultiWriter(tmpFile, bar), resp.Body); err != nil {
			os.Remove(tmpFile.Name())
			return "", fmt.Errorf("failed to write file: %w", err)
		}
	} else {
		// No content length, copy without progress bar
		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			os.Remove(tmpFile.Name())
			return "", fmt.Errorf("failed to write file: %w", err)
		}
	}

	return tmpFile.Name(), nil
}
