// ABOUTME: Template fetcher for downloading and caching bingo project from GitHub
// ABOUTME: Handles tarball download, extraction, and local caching with file locking
package template

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
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
func NewFetcher() (*Fetcher, error) {
	// Get cache directory from user home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, ".bingoctl", "templates")

	// Read mirror from environment variable
	mirror := os.Getenv("BINGOCTL_TEMPLATE_MIRROR")

	return &Fetcher{
		cacheDir: cacheDir,
		timeout:  defaultTimeout,
		mirror:   mirror,
	}, nil
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

// extractTarball extracts tarball to destDir
// Handles GitHub tarball root directory:
//  1. Extract all files to temporary directory
//  2. Detect root directory (should have exactly one directory)
//  3. Move root directory content to destDir
//  4. Return error if format is invalid (no root dir or multiple root dirs)
func (f *Fetcher) extractTarball(tarPath, destDir string) error {
	// Open tarball
	file, err := os.Open(tarPath)
	if err != nil {
		return fmt.Errorf("failed to open tarball: %w", err)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	// Extract to temporary directory first to detect root
	tmpExtractDir := filepath.Join(os.TempDir(), fmt.Sprintf("bingoctl-extract-%d", time.Now().UnixNano()))
	defer os.RemoveAll(tmpExtractDir)

	rootDirs := make(map[string]bool)

	// Extract all files
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Get root directory name - only count TypeDir entries that are root level
		// GitHub tarball format: bingo-main/, bingo-main/file, bingo-main/dir/, etc.
		if header.Typeflag == tar.TypeDir {
			parts := strings.Split(strings.TrimSuffix(header.Name, "/"), "/")
			if len(parts) == 1 && parts[0] != "" {
				rootDirs[parts[0]] = true
			}
		}

		target := filepath.Join(tmpExtractDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		case tar.TypeReg:
			// Create parent directory
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// Create file
			outFile, err := os.Create(target)
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file: %w", err)
			}
			outFile.Close()
		}
	}

	// Verify exactly one root directory
	if len(rootDirs) != 1 {
		return fmt.Errorf("invalid tarball format: expected 1 root directory, found %d", len(rootDirs))
	}

	// Get the single root directory name
	var rootDir string
	for dir := range rootDirs {
		rootDir = dir
		break
	}

	// Move root directory content to destDir
	srcRoot := filepath.Join(tmpExtractDir, rootDir)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create dest directory: %w", err)
	}

	// Copy files from srcRoot to destDir
	return filepath.Walk(srcRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			return err
		}

		return os.Chmod(destPath, info.Mode())
	})
}

// FetchTemplate downloads template to cache (if not exists), returns cache path
// Execution steps:
// 1. Check cache directory exists, create if not (permission 0755)
// 2. Check cache directory is writable, return friendly error if not
// 3. Check cache hit (unless noCache=true)
// 4. If need to download, acquire file lock, download and extract to cache
// 5. Return cache path
func (f *Fetcher) FetchTemplate(ref string, noCache bool) (string, error) {
	// Step 1 & 2: Ensure cache directory exists and is writable
	if err := f.ensureCacheDir(); err != nil {
		return "", err
	}

	// Step 3: Check cache
	cachePath := filepath.Join(f.cacheDir, ref)
	if !noCache && fileExists(cachePath) {
		return cachePath, nil
	}

	// Step 4: Download and cache
	if err := f.downloadAndCache(ref); err != nil {
		return "", err
	}

	// Step 5: Return cache path
	return cachePath, nil
}

// ensureCacheDir ensures cache directory exists and is writable
func (f *Fetcher) ensureCacheDir() error {
	// Create if not exists
	if err := os.MkdirAll(f.cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Check writable
	testFile := filepath.Join(f.cacheDir, ".write-test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("cache directory is not writable: %w", err)
	}
	os.Remove(testFile)

	return nil
}

// downloadAndCache downloads template and extracts to cache
func (f *Fetcher) downloadAndCache(ref string) error {
	// Acquire file lock for concurrent safety
	lockPath := filepath.Join(f.cacheDir, ref+".lock")
	fileLock := flock.New(lockPath)

	locked, err := fileLock.TryLock()
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !locked {
		// Another process is downloading, wait for it
		if err := fileLock.Lock(); err != nil {
			return fmt.Errorf("failed to wait for lock: %w", err)
		}
	}
	defer fileLock.Unlock()
	defer os.Remove(lockPath)

	// Check again if cache exists (may be created by another process)
	cachePath := filepath.Join(f.cacheDir, ref)
	if fileExists(cachePath) {
		return nil
	}

	// Download
	url := f.buildDownloadURL(ref)
	tarPath, err := f.downloadWithTimeout(url)
	if err != nil {
		return fmt.Errorf("failed to download template: %w", err)
	}
	defer os.Remove(tarPath)

	// Extract
	if err := f.extractTarball(tarPath, cachePath); err != nil {
		os.RemoveAll(cachePath) // Cleanup on error
		return fmt.Errorf("failed to extract template: %w", err)
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
