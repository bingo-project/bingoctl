# GitHub Template Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Replace embed.FS template system with GitHub-based template fetching for bingoctl create command

**Architecture:** Download bingo project from GitHub as tarball, cache locally, apply service filtering, directory renaming, and package name replacement, then atomically move to target directory

**Tech Stack:**
- Go 1.24
- github.com/gofrs/flock (file locking)
- github.com/schollz/progressbar/v3 (download progress)
- gopkg.in/yaml.v3 (already in dependencies)
- Go stdlib: archive/tar, compress/gzip, net/http

---

## Prerequisites (Must Complete First)

### Task 0: Create .bingoctl.yaml in bingo project

**Context:** This file must exist in the bingo project repository before implementing bingoctl changes.

**Files:**
- Create in bingo repo: `.bingoctl.yaml`

**Content:**

```yaml
version: 1
services:
  apiserver:
    cmd: cmd/bingo-apiserver
    internal: internal/apiserver
    description: API 服务器
  admserver:
    cmd: cmd/bingo-admserver
    internal: internal/admserver
    description: 管理后台服务器
  bot:
    cmd: cmd/bingo-bot
    internal: internal/bot
    description: Bot 服务
  scheduler:
    cmd: cmd/bingo-scheduler
    internal: internal/scheduler
    description: 定时任务调度器
  ctl:
    cmd: cmd/bingoctl
    internal: internal/bingoctl
    description: 命令行工具
```

**Validation:**
1. File created at bingo repo root
2. YAML is valid
3. Committed to main branch

**Note:** If you don't have access to bingo repo, skip this and we'll use mock data for testing

---

## Phase 1: Setup Dependencies

### Task 1: Install required dependencies

**Step 1: Add dependencies to go.mod**

Run:
```bash
go get github.com/gofrs/flock@latest
go get github.com/schollz/progressbar/v3@latest
```

Expected: Dependencies added successfully

**Step 2: Verify dependencies**

Run: `go mod tidy`

Expected: No errors, go.mod and go.sum updated

**Step 3: Commit dependency changes**

```bash
git add go.mod go.sum
git commit -m "deps: add flock and progressbar for GitHub template fetching"
```

---

## Phase 2: Implement Core Modules

### Task 2: Implement version.go

**Files:**
- Create: `pkg/template/version.go`
- Create: `pkg/template/version_test.go`

**Step 1: Write failing test for DefaultTemplateVersion**

Create `pkg/template/version_test.go`:

```go
package template

import "testing"

func TestDefaultTemplateVersion(t *testing.T) {
	if DefaultTemplateVersion == "" {
		t.Error("DefaultTemplateVersion should not be empty")
	}

	// Should be a valid semver tag
	if DefaultTemplateVersion[0] != 'v' {
		t.Errorf("DefaultTemplateVersion should start with 'v', got: %s", DefaultTemplateVersion)
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/template -v`

Expected: FAIL - package template is not in GOROOT

**Step 3: Write minimal implementation**

Create `pkg/template/version.go`:

```go
// ABOUTME: Template version management for GitHub-based template fetching
// ABOUTME: Defines default version and ref validation logic
package template

// DefaultTemplateVersion is the recommended template version
const DefaultTemplateVersion = "v1.0.0"
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/template -v`

Expected: PASS

**Step 5: Write test for isValidRef**

Add to `pkg/template/version_test.go`:

```go
func TestIsValidRef(t *testing.T) {
	tests := []struct {
		name  string
		ref   string
		valid bool
	}{
		{"valid tag", "v1.2.3", true},
		{"valid branch", "main", true},
		{"valid commit", "abc123def", true},
		{"empty string", "", false},
		{"invalid chars", "v1.2.3@#$", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidRef(tt.ref)
			if result != tt.valid {
				t.Errorf("isValidRef(%q) = %v, want %v", tt.ref, result, tt.valid)
			}
		})
	}
}
```

**Step 6: Run test to verify it fails**

Run: `go test ./pkg/template -v`

Expected: FAIL - undefined: isValidRef

**Step 7: Implement isValidRef**

Add to `pkg/template/version.go`:

```go
import "regexp"

var refRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// isValidRef checks if ref format is valid
// Supports: v1.2.3, main, abc123def, etc.
func isValidRef(ref string) bool {
	if ref == "" {
		return false
	}
	return refRegex.MatchString(ref)
}
```

**Step 8: Run test to verify it passes**

Run: `go test ./pkg/template -v`

Expected: PASS

**Step 9: Write test for refType**

Add to `pkg/template/version_test.go`:

```go
func TestRefType(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{"semver tag", "v1.2.3", "tag"},
		{"tag with prefix", "v0.1.0", "tag"},
		{"branch", "main", "branch"},
		{"branch", "develop", "branch"},
		{"commit hash", "abc123def", "commit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := refType(tt.ref)
			if result != tt.expected {
				t.Errorf("refType(%q) = %v, want %v", tt.ref, result, tt.expected)
			}
		})
	}
}
```

**Step 10: Run test to verify it fails**

Run: `go test ./pkg/template -v`

Expected: FAIL - undefined: refType

**Step 11: Implement refType**

Add to `pkg/template/version.go`:

```go
import "strings"

// refType returns ref type: tag, branch, commit
func refType(ref string) string {
	if strings.HasPrefix(ref, "v") && strings.Contains(ref, ".") {
		return "tag"
	}

	// Common branch names
	if ref == "main" || ref == "master" || ref == "develop" || strings.HasPrefix(ref, "release/") || strings.HasPrefix(ref, "feature/") {
		return "branch"
	}

	// Default to commit hash
	return "commit"
}
```

**Step 12: Run test to verify it passes**

Run: `go test ./pkg/template -v`

Expected: PASS

**Step 13: Commit version module**

```bash
git add pkg/template/version.go pkg/template/version_test.go
git commit -m "feat(template): add version management module"
```

---

### Task 3: Implement config.go

**Files:**
- Create: `pkg/template/config.go`
- Create: `pkg/template/config_test.go`

**Step 1: Write test for BingoctlConfig structure**

Create `pkg/template/config_test.go`:

```go
package template

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadBingoctlConfig(t *testing.T) {
	// Create temp file with valid YAML
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".bingoctl.yaml")

	content := `version: 1
services:
  apiserver:
    cmd: cmd/bingo-apiserver
    internal: internal/apiserver
    description: API server
  ctl:
    cmd: cmd/bingoctl
    internal: internal/bingoctl
    description: CLI tool
`

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test loading
	config, err := loadBingoctlConfig(configPath)
	if err != nil {
		t.Fatalf("loadBingoctlConfig failed: %v", err)
	}

	// Verify structure
	if config.Version != 1 {
		t.Errorf("Version = %d, want 1", config.Version)
	}

	if len(config.Services) != 2 {
		t.Errorf("Services count = %d, want 2", len(config.Services))
	}

	apiserver, ok := config.Services["apiserver"]
	if !ok {
		t.Fatal("apiserver service not found")
	}

	if apiserver.Cmd != "cmd/bingo-apiserver" {
		t.Errorf("apiserver.Cmd = %s, want cmd/bingo-apiserver", apiserver.Cmd)
	}

	if apiserver.Internal != "internal/apiserver" {
		t.Errorf("apiserver.Internal = %s, want internal/apiserver", apiserver.Internal)
	}
}

func TestLoadBingoctlConfig_FileNotExists(t *testing.T) {
	_, err := loadBingoctlConfig("/nonexistent/path/.bingoctl.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestLoadBingoctlConfig_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".bingoctl.yaml")

	// Write invalid YAML
	err := os.WriteFile(configPath, []byte("invalid: yaml: content: ["), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = loadBingoctlConfig(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/template -v -run TestLoadBingoctlConfig`

Expected: FAIL - undefined: loadBingoctlConfig

**Step 3: Implement config structures and loader**

Create `pkg/template/config.go`:

```go
// ABOUTME: Configuration file loader for .bingoctl.yaml metadata
// ABOUTME: Reads service mappings from bingo project template
package template

import (
	"os"

	"gopkg.in/yaml.v3"
)

// BingoctlConfig represents .bingoctl.yaml configuration file structure
type BingoctlConfig struct {
	Version  int                    `yaml:"version"`
	Services map[string]ServiceInfo `yaml:"services"`
}

// ServiceInfo describes a service's directory structure
type ServiceInfo struct {
	Cmd         string `yaml:"cmd"`
	Internal    string `yaml:"internal"`
	Description string `yaml:"description"`
}

// loadBingoctlConfig loads and parses .bingoctl.yaml
func loadBingoctlConfig(path string) (*BingoctlConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config BingoctlConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/template -v -run TestLoadBingoctlConfig`

Expected: PASS (3 tests)

**Step 5: Commit config module**

```bash
git add pkg/template/config.go pkg/template/config_test.go
git commit -m "feat(template): add .bingoctl.yaml config loader"
```

---

### Task 4: Implement utility functions in util/file.go

**Files:**
- Modify: `pkg/util/file.go`
- Create: `pkg/util/file_test.go`

**Step 1: Write test for CopyDir**

Create `pkg/util/file_test.go`:

```go
package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyDir(t *testing.T) {
	// Create source directory structure
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	dstDir := filepath.Join(tmpDir, "dst")

	// Create test structure
	err := os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	if err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}

	err = os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	err = os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	// Test copy
	err = CopyDir(srcDir, dstDir)
	if err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	// Verify files copied
	content1, err := os.ReadFile(filepath.Join(dstDir, "file1.txt"))
	if err != nil {
		t.Errorf("file1.txt not copied: %v", err)
	}
	if string(content1) != "content1" {
		t.Errorf("file1.txt content = %s, want content1", content1)
	}

	content2, err := os.ReadFile(filepath.Join(dstDir, "subdir", "file2.txt"))
	if err != nil {
		t.Errorf("subdir/file2.txt not copied: %v", err)
	}
	if string(content2) != "content2" {
		t.Errorf("file2.txt content = %s, want content2", content2)
	}
}

func TestCopyDir_SrcNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	err := CopyDir("/nonexistent/src", filepath.Join(tmpDir, "dst"))
	if err == nil {
		t.Error("Expected error for non-existent source, got nil")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/util -v -run TestCopyDir`

Expected: FAIL - undefined: CopyDir

**Step 3: Implement CopyDir**

Add to `pkg/util/file.go`:

```go
import (
	"io"
	"path/filepath"
)

// CopyDir recursively copies a directory
func CopyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return copyFile(path, dstPath, info.Mode())
	})
}

// copyFile copies a single file
func copyFile(src, dst string, mode os.FileMode) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create parent directory
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return os.Chmod(dst, mode)
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/util -v -run TestCopyDir`

Expected: PASS (2 tests)

**Step 5: Commit utility functions**

```bash
git add pkg/util/file.go pkg/util/file_test.go
git commit -m "feat(util): add CopyDir function for recursive directory copying"
```

---

### Task 5: Implement fetcher.go (Part 1: Core structure and URL building)

**Files:**
- Create: `pkg/template/fetcher.go`
- Create: `pkg/template/fetcher_test.go`

**Step 1: Write test for buildDownloadURL**

Create `pkg/template/fetcher_test.go`:

```go
package template

import (
	"testing"
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/template -v -run TestBuildDownloadURL`

Expected: FAIL - undefined: Fetcher

**Step 3: Implement Fetcher structure and buildDownloadURL**

Create `pkg/template/fetcher.go`:

```go
// ABOUTME: Template fetcher for downloading and caching bingo project from GitHub
// ABOUTME: Handles tarball download, extraction, and local caching with file locking
package template

import (
	"fmt"
	"time"
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/template -v -run TestBuildDownloadURL`

Expected: PASS

**Step 5: Commit fetcher part 1**

```bash
git add pkg/template/fetcher.go pkg/template/fetcher_test.go
git commit -m "feat(template): add Fetcher structure and URL builder"
```

---

### Task 6: Implement fetcher.go (Part 2: Download with timeout and progress)

**Step 1: Write test for downloadWithTimeout (mock)**

Add to `pkg/template/fetcher_test.go`:

```go
import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

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
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/template -v -run TestDownloadWithTimeout`

Expected: FAIL - undefined: downloadWithTimeout

**Step 3: Implement downloadWithTimeout**

Add to `pkg/template/fetcher.go`:

```go
import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

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

	// Create progress bar if content length is known
	var reader io.Reader = resp.Body
	if resp.ContentLength > 0 {
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"Downloading",
		)
		reader = progressbar.NewReader(resp.Body, bar)
	}

	// Download to temp file
	if _, err := io.Copy(tmpFile, reader); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return tmpFile.Name(), nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/template -v -run TestDownloadWithTimeout`

Expected: PASS (2 tests)

**Step 5: Commit fetcher part 2**

```bash
git add pkg/template/fetcher.go pkg/template/fetcher_test.go
git commit -m "feat(template): add download with timeout and progress"
```

---

### Task 7: Implement fetcher.go (Part 3: Tarball extraction)

**Step 1: Write test for extractTarball**

Add to `pkg/template/fetcher_test.go`:

```go
import (
	"archive/tar"
	"compress/gzip"
)

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

	// Add files with root directory prefix (like GitHub tarball)
	files := map[string]string{
		rootDir + "/file1.txt":         "content1",
		rootDir + "/subdir/file2.txt": "content2",
	}

	for name, content := range files {
		// Add directory entry if needed
		if filepath.Dir(name) != rootDir {
			dirHeader := &tar.Header{
				Name:     filepath.Dir(name) + "/",
				Mode:     0755,
				Typeflag: tar.TypeDir,
			}
			if err := tw.WriteHeader(dirHeader); err != nil {
				return err
			}
		}

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
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/template -v -run TestExtractTarball`

Expected: FAIL - undefined: extractTarball

**Step 3: Implement extractTarball**

Add to `pkg/template/fetcher.go`:

```go
import (
	"archive/tar"
	"compress/gzip"
	"path/filepath"
	"strings"
)

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

		// Get root directory name
		parts := strings.Split(header.Name, "/")
		if len(parts) > 0 && parts[0] != "" {
			rootDirs[parts[0]] = true
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
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/template -v -run TestExtractTarball`

Expected: PASS (2 tests)

**Step 5: Commit fetcher part 3**

```bash
git add pkg/template/fetcher.go pkg/template/fetcher_test.go
git commit -m "feat(template): add tarball extraction with root directory handling"
```

---

### Task 8: Implement fetcher.go (Part 4: Cache management and FetchTemplate)

**Step 1: Write test for FetchTemplate**

Add to `pkg/template/fetcher_test.go`:

```go
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
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/template -v -run TestFetchTemplate`

Expected: FAIL - undefined: FetchTemplate

**Step 3: Implement cache directory initialization and FetchTemplate**

Add to `pkg/template/fetcher.go`:

```go
import (
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
)

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
```

**Step 4: Update NewFetcher calls in tests**

Update `pkg/template/fetcher_test.go` to use new constructor:

```go
// In tests that create Fetcher manually, they can still use &Fetcher{} directly
// No changes needed for existing tests
```

**Step 5: Run test to verify it passes**

Run: `go test ./pkg/template -v -run TestFetchTemplate`

Expected: PASS

**Step 6: Commit fetcher part 4**

```bash
git add pkg/template/fetcher.go pkg/template/fetcher_test.go
git commit -m "feat(template): add FetchTemplate with caching and file locking"
```

---

### Task 9: Implement replacer.go (Part 1: shouldReplaceFile)

**Files:**
- Create: `pkg/template/replacer.go`
- Create: `pkg/template/replacer_test.go`

**Step 1: Write test for shouldReplaceFile**

Create `pkg/template/replacer_test.go`:

```go
package template

import (
	"testing"
)

func TestShouldReplaceFile(t *testing.T) {
	r := &Replacer{}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		// Should replace
		{"go file", "main.go", true},
		{"go.mod", "go.mod", true},
		{"Makefile", "Makefile", true},
		{"Dockerfile", "Dockerfile", true},
		{".gitignore", ".gitignore", true},
		{".env", ".env", true},
		{".env.example", ".env.example", true},
		{"yaml file", "config.yaml", true},
		{"shell script", "build.sh", true},

		// Should not replace
		{"binary", "app", false},
		{"image", "logo.png", false},
		{"pdf", "doc.pdf", false},
		{"exe", "app.exe", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := r.shouldReplaceFile(tt.path)
			if result != tt.expected {
				t.Errorf("shouldReplaceFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/template -v -run TestShouldReplaceFile`

Expected: FAIL - undefined: Replacer

**Step 3: Implement Replacer structure and shouldReplaceFile**

Create `pkg/template/replacer.go`:

```go
// ABOUTME: Package name and directory name replacer for template processing
// ABOUTME: Handles module name substitution and directory renaming
package template

import (
	"path/filepath"
	"strings"
)

var replaceableExtensions = map[string]bool{
	// Go related
	".go":  true,
	".mod": true,
	".sum": true,

	// Documentation
	".md":  true,
	".txt": true,

	// Build and scripts
	".mk":   true,
	".sh":   true,
	".bash": true,

	// Config files
	".yaml": true,
	".yml":  true,
	".toml": true,
	".json": true,
	".env":  true,

	// Docker
	".dockerignore": true,

	// Git
	".gitignore": true,
}

var replaceableBasenames = map[string]bool{
	"Makefile":   true,
	"Dockerfile": true,
}

// Replacer handles module name and directory name replacement
type Replacer struct {
	targetDir string // target directory
	oldModule string // "bingo"
	newModule string // "github.com/mycompany/demo"
	appName   string // "demo"
}

// NewReplacer creates a new Replacer instance
func NewReplacer(targetDir, oldModule, newModule, appName string) *Replacer {
	return &Replacer{
		targetDir: targetDir,
		oldModule: oldModule,
		newModule: newModule,
		appName:   appName,
	}
}

// shouldReplaceFile determines if file should be processed for replacement
// Based on file extension whitelist
func (r *Replacer) shouldReplaceFile(path string) bool {
	ext := filepath.Ext(path)
	base := filepath.Base(path)

	// Check basename first
	if replaceableBasenames[base] {
		return true
	}

	// Check extension
	if replaceableExtensions[ext] {
		return true
	}

	// Special case: .env files with extensions like .env.example
	if strings.HasPrefix(base, ".env") {
		return true
	}

	return false
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/template -v -run TestShouldReplaceFile`

Expected: PASS

**Step 5: Commit replacer part 1**

```bash
git add pkg/template/replacer.go pkg/template/replacer_test.go
git commit -m "feat(template): add Replacer with file type whitelist"
```

---

### Task 10: Implement replacer.go (Part 2: ReplaceModuleName)

**Step 1: Write test for replaceInFile**

Add to `pkg/template/replacer_test.go`:

```go
import (
	"os"
	"path/filepath"
)

func TestReplaceInFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tmpDir, "test.go")
	content := `package main

import "bingo/internal/pkg/log"

func main() {
	log.Info("bingo application")
	path := "bingo/cmd/server"
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Replace
	r := NewReplacer(tmpDir, "bingo", "github.com/myapp/myapp", "myapp")
	err = r.replaceInFile(testFile)
	if err != nil {
		t.Fatalf("replaceInFile failed: %v", err)
	}

	// Verify
	result, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	resultStr := string(result)

	// Check replacements
	if !strings.Contains(resultStr, `import "github.com/myapp/myapp/internal/pkg/log"`) {
		t.Error("Import path not replaced")
	}

	if !strings.Contains(resultStr, `log.Info("bingo application")`) {
		t.Error("String literal should not be replaced with full module name")
	}

	if !strings.Contains(resultStr, `"github.com/myapp/myapp/cmd/server"`) {
		t.Error("Path in string not replaced")
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/template -v -run TestReplaceInFile`

Expected: FAIL - undefined: replaceInFile

**Step 3: Implement replaceInFile and ReplaceModuleName**

Add to `pkg/template/replacer.go`:

```go
import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ReplaceModuleName replaces all files with module name
// Traverses target directory, replaces based on file extension
func (r *Replacer) ReplaceModuleName() error {
	return filepath.WalkDir(r.targetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if r.shouldReplaceFile(path) {
			return r.replaceInFile(path)
		}

		return nil
	})
}

// replaceInFile replaces module name in a single file
// Uses string replacement to avoid breaking binary files
func (r *Replacer) replaceInFile(path string) error {
	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	str := string(content)

	// Replace patterns
	// 1. go.mod: module bingo -> module {newModule}
	str = strings.ReplaceAll(str, "module "+r.oldModule, "module "+r.newModule)

	// 2. imports: "bingo/xxx" -> "{newModule}/xxx"
	str = strings.ReplaceAll(str, `"`+r.oldModule+"/", `"`+r.newModule+"/")

	// 3. paths in strings: bingo/ -> {newModule}/
	// Note: This is aggressive but necessary for Makefile, Dockerfile, etc.
	str = strings.ReplaceAll(str, r.oldModule+"/", r.newModule+"/")

	// Write back
	err = os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/template -v -run TestReplaceInFile`

Expected: PASS

**Step 5: Commit replacer part 2**

```bash
git add pkg/template/replacer.go pkg/template/replacer_test.go
git commit -m "feat(template): add module name replacement"
```

---

### Task 11: Implement replacer.go (Part 3: RenameDirs)

**Step 1: Write test for RenameDirs**

Add to `pkg/template/replacer_test.go`:

```go
func TestRenameDirs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure
	dirs := []string{
		"cmd/bingo-apiserver",
		"cmd/bingo-admserver",
		"cmd/bingoctl",
		"internal/apiserver",
		"internal/pkg",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	// Rename
	r := NewReplacer(tmpDir, "bingo", "github.com/myapp/myapp", "myapp")
	err := r.RenameDirs()
	if err != nil {
		t.Fatalf("RenameDirs failed: %v", err)
	}

	// Verify renames
	expected := []string{
		"cmd/myapp-apiserver",
		"cmd/myapp-admserver",
		"cmd/myappctl",
		"internal/apiserver", // Should not be renamed
		"internal/pkg",       // Should not be renamed
	}

	for _, dir := range expected {
		path := filepath.Join(tmpDir, dir)
		if !fileExists(path) {
			t.Errorf("Expected directory not found: %s", dir)
		}
	}

	// Verify old directories removed
	notExpected := []string{
		"cmd/bingo-apiserver",
		"cmd/bingo-admserver",
		"cmd/bingoctl",
	}

	for _, dir := range notExpected {
		path := filepath.Join(tmpDir, dir)
		if fileExists(path) {
			t.Errorf("Old directory still exists: %s", dir)
		}
	}
}
```

**Step 2: Run test to verify it fails**

Run: `go test ./pkg/template -v -run TestRenameDirs`

Expected: FAIL - undefined: RenameDirs

**Step 3: Implement RenameDirs**

Add to `pkg/template/replacer.go`:

```go
// renameRules defines directory rename mappings
// Only these explicitly listed directories will be renamed
var renameRules = map[string]string{
	"cmd/bingo-apiserver":   "cmd/{app}-apiserver",
	"cmd/bingo-admserver":   "cmd/{app}-admserver",
	"cmd/bingo-bot":         "cmd/{app}-bot",
	"cmd/bingo-scheduler":   "cmd/{app}-scheduler",
	"cmd/bingoctl":          "cmd/{app}ctl",
}

// RenameDirs renames directories according to explicit rules
// Only renames directories that still exist (after service filtering)
func (r *Replacer) RenameDirs() error {
	for oldPath, newPathTemplate := range renameRules {
		// Replace {app} placeholder
		newPath := strings.ReplaceAll(newPathTemplate, "{app}", r.appName)

		oldFullPath := filepath.Join(r.targetDir, oldPath)
		newFullPath := filepath.Join(r.targetDir, newPath)

		// Skip if old path doesn't exist (may be filtered out)
		if !fileExists(oldFullPath) {
			continue
		}

		// Rename
		err := os.Rename(oldFullPath, newFullPath)
		if err != nil {
			return fmt.Errorf("failed to rename %s to %s: %w", oldPath, newPath, err)
		}
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
```

**Step 4: Run test to verify it passes**

Run: `go test ./pkg/template -v -run TestRenameDirs`

Expected: PASS

**Step 5: Commit replacer part 3**

```bash
git add pkg/template/replacer.go pkg/template/replacer_test.go
git commit -m "feat(template): add directory renaming"
```

---

## Phase 3: Integrate with Create Command

### Task 12: Refactor create.go to use new template system

**Files:**
- Modify: `pkg/cmd/create/create.go`

**Step 1: Add new fields to CreateOptions**

Update `CreateOptions` struct in `pkg/cmd/create/create.go`:

```go
type CreateOptions struct {
	// Add new fields
	ModuleName   string // Go module name (optional)
	TemplateRef  string // Template version
	NoCache      bool   // Force re-download

	// Keep existing fields
	GoVersion    string
	TemplatePath string
	RootPackage  string
	AppName      string
	AppNameCamel string

	// Service selection (keep existing)
	Services    []string
	NoServices  []string
	AddServices []string
	Interactive bool

	selectedServices []string
}
```

**Step 2: Add new flags to NewCmdCreate**

Add flags after line 75 in `pkg/cmd/create/create.go`:

```go
	cmd.Flags().StringVarP(&o.ModuleName, "module", "m", "",
		"Go module name (e.g., github.com/mycompany/myapp)")
	cmd.Flags().StringVarP(&o.TemplateRef, "ref", "r", "",
		"Template version (tag/branch/commit, default: recommended version)")
	cmd.Flags().BoolVar(&o.NoCache, "no-cache", false,
		"Force re-download template (for branches)")
```

**Step 3: Update Complete method**

Replace the Complete method implementation:

```go
func (o *CreateOptions) Complete(cmd *cobra.Command, args []string) error {
	// 1. Parse template version
	if o.TemplateRef == "" {
		o.TemplateRef = template.DefaultTemplateVersion
		console.Info(fmt.Sprintf("使用推荐版本：%s", o.TemplateRef))
	}

	// 2. Compute service list (keep existing logic)
	o.Interactive = len(o.Services) == 0 && len(o.NoServices) == 0 && len(o.AddServices) == 0

	if o.Interactive {
		console.Info("进入交互模式...")
		selected, err := o.selectServicesInteractively()
		if err != nil {
			return err
		}
		o.selectedServices = selected
	} else {
		o.selectedServices = o.computeServiceList()
	}

	// Warn if no services selected
	if len(o.selectedServices) == 0 {
		console.Warn("未选择任何服务，将创建最小项目骨架")
		prompt := promptui.Prompt{
			Label:     "继续",
			IsConfirm: true,
		}
		_, err := prompt.Run()
		if err != nil {
			console.Exit("已取消创建")
		}
	}

	return nil
}
```

**Step 4: Rewrite Run method**

Replace the Run method implementation:

```go
import (
	"github.com/bingo-project/bingoctl/pkg/template"
	"time"
)

func (o *CreateOptions) Run(args []string) error {
	console.Info(fmt.Sprintf("Creating project %s", o.RootPackage))

	// 1. Fetch template (download or use cache)
	fetcher, err := template.NewFetcher()
	if err != nil {
		return fmt.Errorf("failed to create fetcher: %w", err)
	}

	templatePath, err := fetcher.FetchTemplate(o.TemplateRef, o.NoCache)
	if err != nil {
		return fmt.Errorf("获取模板失败: %w", err)
	}

	// 2. Create temporary directory
	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("bingoctl-%d", time.Now().Unix()))
	defer os.RemoveAll(tmpDir)

	// 3. Copy to temporary directory
	console.Info("复制模板...")
	if err := cmdutil.CopyDir(templatePath, tmpDir); err != nil {
		return fmt.Errorf("复制模板失败: %w", err)
	}

	// 4. Filter services (before renaming, using original directory names)
	if len(o.selectedServices) > 0 {
		console.Info("过滤服务...")
		if err := o.filterServices(tmpDir); err != nil {
			return err
		}
	}

	// 5. Rename directories (always execute, only for remaining directories)
	replacer := template.NewReplacer(tmpDir, "bingo", o.ModuleName, o.AppName)
	console.Info("重命名目录...")
	if err := replacer.RenameDirs(); err != nil {
		return fmt.Errorf("重命名目录失败: %w", err)
	}

	// 6. Replace module name (only if -m specified)
	if o.ModuleName != "" {
		console.Info(fmt.Sprintf("替换模块名: bingo -> %s", o.ModuleName))
		if err := replacer.ReplaceModuleName(); err != nil {
			return fmt.Errorf("替换模块名失败: %w", err)
		}
	}

	// 7. Atomically move to target location
	if err := os.Rename(tmpDir, o.AppName); err != nil {
		return fmt.Errorf("移动项目失败: %w", err)
	}

	// 8. Success message
	console.Success("项目创建成功！")
	if len(o.selectedServices) == 0 {
		console.Info("提示：已删除所有服务，建议运行 'go mod tidy' 清理未使用的依赖")
	}

	return nil
}
```

**Step 5: Implement filterServices method**

Add new method to `pkg/cmd/create/create.go`:

```go
// filterServices deletes unselected service directories
// Reads service mapping from .bingoctl.yaml
func (o *CreateOptions) filterServices(targetDir string) error {
	// Load .bingoctl.yaml
	configPath := filepath.Join(targetDir, ".bingoctl.yaml")
	config, err := template.loadBingoctlConfig(configPath)
	if err != nil {
		return fmt.Errorf("加载 .bingoctl.yaml 失败: %w\n提示：模板版本 %s 可能不包含此配置文件", o.TemplateRef, err)
	}

	allServices := config.Services

	// Mark selected services
	selected := make(map[string]bool)
	for _, svc := range o.selectedServices {
		selected[svc] = true
	}

	// Delete unselected service directories
	for svc, service := range allServices {
		if !selected[svc] {
			// Delete cmd directory
			cmdPath := filepath.Join(targetDir, service.Cmd)
			if cmdutil.Exists(cmdPath) {
				console.Info(fmt.Sprintf("  删除 %s", service.Cmd))
				if err := os.RemoveAll(cmdPath); err != nil {
					return fmt.Errorf("删除 %s 失败: %w", service.Cmd, err)
				}
			}

			// Delete internal directory
			internalPath := filepath.Join(targetDir, service.Internal)
			if cmdutil.Exists(internalPath) {
				console.Info(fmt.Sprintf("  删除 %s", service.Internal))
				if err := os.RemoveAll(internalPath); err != nil {
					return fmt.Errorf("删除 %s 失败: %w", service.Internal, err)
				}
			}
		}
	}

	return nil
}
```

**Step 6: Fix import for loadBingoctlConfig**

Since `loadBingoctlConfig` is not exported, we need to export it in `pkg/template/config.go`:

```go
// LoadBingoctlConfig loads and parses .bingoctl.yaml (exported)
func LoadBingoctlConfig(path string) (*BingoctlConfig, error) {
	return loadBingoctlConfig(path)
}
```

Then update the call in create.go:

```go
config, err := template.LoadBingoctlConfig(configPath)
```

**Step 7: Add missing imports**

Add to imports in `pkg/cmd/create/create.go`:

```go
import (
	// ... existing imports ...
	"path/filepath"
	"time"

	"github.com/bingo-project/bingoctl/pkg/template"
)
```

**Step 8: Build and test compilation**

Run: `go build ./pkg/cmd/create`

Expected: No compilation errors

**Step 9: Commit create.go refactor**

```bash
git add pkg/cmd/create/create.go pkg/template/config.go
git commit -m "feat(create): refactor to use GitHub template system"
```

---

### Task 13: Remove old template system

**Step 1: Remove embed.FS variables**

Delete or comment out lines 22-25 in `pkg/cmd/create/create.go`:

```go
// REMOVE:
// var (
// 	//go:embed tpl
// 	tplFS embed.FS
// 	root  = "tpl"
// )
```

**Step 2: Remove embed import**

Remove `"embed"` from imports in `pkg/cmd/create/create.go`

**Step 3: Delete tpl directory**

Run: `rm -rf pkg/cmd/create/tpl`

**Step 4: Build and verify**

Run: `go build ./cmd/bingoctl`

Expected: No errors

**Step 5: Commit cleanup**

```bash
git add pkg/cmd/create/create.go
git status # Verify tpl/ is deleted
git add -A # Add deletion
git commit -m "refactor(create): remove old embed.FS template system"
```

---

## Phase 4: Testing and Documentation

### Task 14: Manual integration test

**Step 1: Build bingoctl**

Run: `go build -o /tmp/bingoctl-test ./cmd/bingoctl`

Expected: Binary created at /tmp/bingoctl-test

**Step 2: Test default creation**

Run:
```bash
cd /tmp
./bingoctl-test create test-default
```

Expected:
- Downloads template (shows progress)
- Creates test-default/ directory
- Contains default services (apiserver, ctl)

**Step 3: Test with custom module**

Run:
```bash
cd /tmp
./bingoctl-test create test-custom -m github.com/test/custom
```

Expected:
- Creates test-custom/ directory
- go.mod contains "module github.com/test/custom"
- Import paths updated

**Step 4: Test with service selection**

Run:
```bash
cd /tmp
./bingoctl-test create test-services --services apiserver
```

Expected:
- Only apiserver service created
- Other service directories not present

**Step 5: Test cache behavior**

Run:
```bash
cd /tmp
./bingoctl-test create test-cache1
./bingoctl-test create test-cache2
```

Expected:
- First run downloads
- Second run uses cache (no download progress)

**Step 6: Document results**

Create test log in `docs/plans/test-results.md` with output

---

### Task 15: Update README

**Files:**
- Modify: `README.md`

**Step 1: Document new flags**

Add section in README.md:

```markdown
### Create Command Options

#### Template Version

```bash
# Use default recommended version
bingoctl create myapp

# Use specific version
bingoctl create myapp -r v1.2.3

# Use branch (development)
bingoctl create myapp -r main
```

#### Custom Module Name

```bash
# Replace package name
bingoctl create myapp -m github.com/mycompany/myapp
```

#### Service Selection

```bash
# Select specific services
bingoctl create myapp --services apiserver,ctl

# Exclude services
bingoctl create myapp --no-service bot,scheduler

# Add services to defaults
bingoctl create myapp --add-service admserver
```

#### Cache Management

```bash
# Force refresh template (for branches)
bingoctl create myapp -r main --no-cache
```

#### Environment Variables

```bash
# Use mirror for GitHub access
export BINGOCTL_TEMPLATE_MIRROR=https://ghproxy.com/
bingoctl create myapp
```
```

**Step 2: Commit README update**

```bash
git add README.md
git commit -m "docs(readme): add GitHub template system documentation"
```

---

## Summary and Next Steps

Once all tasks are complete:

1. **Run full test suite**: `go test ./...`
2. **Build final binary**: `go build -o bingoctl ./cmd/bingoctl`
3. **Tag release**: Consider updating DefaultTemplateVersion if needed
4. **Monitor for issues**: Watch for cache-related bugs, download failures

## Future Enhancements (Not in This Plan)

- Cache management commands (`bingoctl cache clean`, etc.)
- Progress bar improvements
- Better error messages with suggestions
- Unit tests for edge cases
- Integration tests with real GitHub

---

**Plan complete!** Ready for execution.
