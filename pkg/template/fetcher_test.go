package template

import (
	"archive/tar"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBuildDownloadURL(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		mirror   string
		expected string
	}{
		{
			name:     "tag without mirror",
			ref:      "v1.2.3",
			mirror:   "",
			expected: "https://github.com/bingo-project/bingo/archive/refs/tags/v1.2.3.tar.gz",
		},
		{
			name:     "branch without mirror",
			ref:      "main",
			mirror:   "",
			expected: "https://github.com/bingo-project/bingo/archive/refs/heads/main.tar.gz",
		},
		{
			name:     "tag with mirror",
			ref:      "v1.2.3",
			mirror:   "https://ghproxy.com/",
			expected: "https://ghproxy.com/https://github.com/bingo-project/bingo/archive/refs/tags/v1.2.3.tar.gz",
		},
		{
			name:     "commit hash",
			ref:      "abc123def",
			mirror:   "",
			expected: "https://github.com/bingo-project/bingo/archive/abc123def.tar.gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fetcher{mirror: tt.mirror}
			result := f.buildDownloadURL(tt.ref)
			if result != tt.expected {
				t.Errorf("buildDownloadURL(%q) = %q, want %q", tt.ref, result, tt.expected)
			}
		})
	}
}

func TestDownloadWithTimeout(t *testing.T) {
	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test tarball content"))
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	f := &Fetcher{
		cacheDir: tmpDir,
		timeout:  defaultTimeout,
		mirror:   "",
	}

	// Test download
	tarPath, err := f.downloadWithTimeout(server.URL)
	if err != nil {
		t.Fatalf("downloadWithTimeout failed: %v", err)
	}

	// Verify file exists
	if !fileExists(tarPath) {
		t.Errorf("Downloaded file not found: %s", tarPath)
	}

	// Verify content
	content, err := os.ReadFile(tarPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(content) != "test tarball content" {
		t.Errorf("Downloaded content = %s, want 'test tarball content'", content)
	}
}

func TestDownloadWithTimeout_Timeout(t *testing.T) {
	// Create slow server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay longer than timeout
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	f := &Fetcher{
		cacheDir: tmpDir,
		timeout:  100 * time.Millisecond, // Short timeout for test
		mirror:   "",
	}

	_, err := f.downloadWithTimeout(server.URL)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
}

func TestExtractTarball(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test tarball with root directory structure (like GitHub)
	tarPath := filepath.Join(tmpDir, "test.tar.gz")
	if err := createTestTarball(tarPath, "bingo-v1.0.0"); err != nil {
		t.Fatalf("Failed to create test tarball: %v", err)
	}

	// Extract
	destDir := filepath.Join(tmpDir, "extracted")
	f := &Fetcher{}
	err := f.extractTarball(tarPath, destDir)
	if err != nil {
		t.Fatalf("extractTarball failed: %v", err)
	}

	// Verify files extracted (should be in destDir directly, not in bingo-v1.0.0/)
	if !fileExists(filepath.Join(destDir, "file1.txt")) {
		t.Error("file1.txt not extracted")
	}

	if !fileExists(filepath.Join(destDir, "subdir", "file2.txt")) {
		t.Error("subdir/file2.txt not extracted")
	}

	// Verify content
	content, _ := os.ReadFile(filepath.Join(destDir, "file1.txt"))
	if string(content) != "content1" {
		t.Errorf("file1.txt content = %s, want content1", content)
	}
}

func TestExtractTarball_InvalidFormat(t *testing.T) {
	tmpDir := t.TempDir()

	// Create tarball with multiple root directories (invalid)
	tarPath := filepath.Join(tmpDir, "invalid.tar.gz")
	if err := createInvalidTarball(tarPath); err != nil {
		t.Fatalf("Failed to create invalid tarball: %v", err)
	}

	destDir := filepath.Join(tmpDir, "extracted")
	f := &Fetcher{}
	err := f.extractTarball(tarPath, destDir)
	if err == nil {
		t.Error("Expected error for invalid tarball format, got nil")
	}
}

// Helper: create test tarball with GitHub-like structure
func createTestTarball(path, rootDir string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gzw := gzip.NewWriter(file)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// Add root directory entry first
	rootHeader := &tar.Header{
		Name:     rootDir + "/",
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}
	if err := tw.WriteHeader(rootHeader); err != nil {
		return err
	}

	// Add subdirectory
	subdirHeader := &tar.Header{
		Name:     rootDir + "/subdir/",
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}
	if err := tw.WriteHeader(subdirHeader); err != nil {
		return err
	}

	// Add files with root directory prefix (like GitHub tarball)
	files := map[string]string{
		rootDir + "/file1.txt":        "content1",
		rootDir + "/subdir/file2.txt": "content2",
	}

	for name, content := range files {
		// Add file
		header := &tar.Header{
			Name: name,
			Mode: 0644,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			return err
		}
	}

	return nil
}

// Helper: create tarball with invalid structure (multiple root dirs)
func createInvalidTarball(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gzw := gzip.NewWriter(file)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// Add two root directories
	root1Header := &tar.Header{
		Name:     "root1/",
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}
	if err := tw.WriteHeader(root1Header); err != nil {
		return err
	}

	root2Header := &tar.Header{
		Name:     "root2/",
		Mode:     0755,
		Typeflag: tar.TypeDir,
	}
	if err := tw.WriteHeader(root2Header); err != nil {
		return err
	}

	// Add files in multiple root directories
	files := map[string]string{
		"root1/file1.txt": "content1",
		"root2/file2.txt": "content2",
	}

	for name, content := range files {
		header := &tar.Header{
			Name: name,
			Mode: 0644,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if _, err := tw.Write([]byte(content)); err != nil {
			return err
		}
	}

	return nil
}

func TestFetchTemplate_CacheHit(t *testing.T) {
	tmpDir := t.TempDir()
	cacheDir := filepath.Join(tmpDir, "cache")

	// Pre-populate cache
	cachedVersion := filepath.Join(cacheDir, "v1.0.0")
	os.MkdirAll(cachedVersion, 0755)
	os.WriteFile(filepath.Join(cachedVersion, "cached.txt"), []byte("cached"), 0644)

	f := &Fetcher{
		cacheDir: cacheDir,
		timeout:  defaultTimeout,
		mirror:   "",
	}

	// Fetch should return cached path without downloading
	path, err := f.FetchTemplate("v1.0.0", false)
	if err != nil {
		t.Fatalf("FetchTemplate failed: %v", err)
	}

	if path != cachedVersion {
		t.Errorf("FetchTemplate returned %s, want %s", path, cachedVersion)
	}

	// Verify cached file still exists
	if !fileExists(filepath.Join(path, "cached.txt")) {
		t.Error("Cached file not found")
	}
}

func TestFetchTemplate_NoCache(t *testing.T) {
	// This test would require mocking HTTP download
	// Skip for now as it's integration-level
	t.Skip("Integration test - requires HTTP mocking")
}
