package template

import (
	"net/http"
	"net/http/httptest"
	"os"
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

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
