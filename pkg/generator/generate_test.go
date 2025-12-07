// ABOUTME: Tests for code generation utilities including service discovery.
// ABOUTME: Provides test coverage for scanning cmd/ directory and extracting service names.
package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverServices(t *testing.T) {
	// Create temporary cmd directory structure
	tmpDir := t.TempDir()
	cmdDir := filepath.Join(tmpDir, "cmd")
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		t.Fatalf("Failed to create cmd dir: %v", err)
	}

	// Create service directories
	services := []string{"myapp-apiserver", "myapp-admserver", "myappctl"}
	for _, svc := range services {
		if err := os.MkdirAll(filepath.Join(cmdDir, svc), 0755); err != nil {
			t.Fatalf("Failed to create service dir %s: %v", svc, err)
		}
	}

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Test discovery
	discovered, err := discoverServices()
	if err != nil {
		t.Fatalf("discoverServices failed: %v", err)
	}

	expected := map[string]bool{
		"apiserver": true,
		"admserver": true,
		"ctl":       true,
	}

	if len(discovered) != len(expected) {
		t.Errorf("Expected %d services, got %d: %v", len(expected), len(discovered), discovered)
	}

	// Check for unexpected services
	seen := make(map[string]bool)
	for _, svc := range discovered {
		if !expected[svc] {
			t.Errorf("Unexpected service: %s", svc)
		}
		if seen[svc] {
			t.Errorf("Duplicate service found: %s", svc)
		}
		seen[svc] = true
	}

	// Check for missing services
	for svc := range expected {
		if !seen[svc] {
			t.Errorf("Missing expected service: %s", svc)
		}
	}
}

func TestDiscoverServices_NoCmd(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	_, err := discoverServices()
	if err == nil {
		t.Error("Expected error when cmd/ doesn't exist, got nil")
	}
}

func TestDiscoverServices_EdgeCases(t *testing.T) {
	// Create temporary cmd directory structure
	tmpDir := t.TempDir()
	cmdDir := filepath.Join(tmpDir, "cmd")
	if err := os.MkdirAll(cmdDir, 0755); err != nil {
		t.Fatalf("Failed to create cmd dir: %v", err)
	}

	// Create various edge cases:
	// 1. Directory without hyphen or ctl suffix (should be skipped)
	if err := os.MkdirAll(filepath.Join(cmdDir, "randomdir"), 0755); err != nil {
		t.Fatalf("Failed to create randomdir: %v", err)
	}

	// 2. File instead of directory
	testFile := filepath.Join(cmdDir, "README.md")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 3. Valid service with hyphen
	if err := os.MkdirAll(filepath.Join(cmdDir, "myapp-worker"), 0755); err != nil {
		t.Fatalf("Failed to create myapp-worker: %v", err)
	}

	// 4. Valid ctl service
	if err := os.MkdirAll(filepath.Join(cmdDir, "testctl"), 0755); err != nil {
		t.Fatalf("Failed to create testctl: %v", err)
	}

	// 5. Directory starting with dot (should be skipped by logic)
	if err := os.MkdirAll(filepath.Join(cmdDir, ".hidden"), 0755); err != nil {
		t.Fatalf("Failed to create .hidden: %v", err)
	}

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Test discovery
	discovered, err := discoverServices()
	if err != nil {
		t.Fatalf("discoverServices failed: %v", err)
	}

	// Should only find "worker" and "ctl"
	// "randomdir" has no hyphen and doesn't end with ctl, so it's skipped
	// README.md is a file, not a directory
	// .hidden has no hyphen and doesn't end with ctl
	expected := map[string]bool{
		"worker": true,
		"ctl":    true,
	}

	if len(discovered) != len(expected) {
		t.Errorf("Expected %d services, got %d: %v", len(expected), len(discovered), discovered)
	}

	// Check for unexpected services
	seen := make(map[string]bool)
	for _, svc := range discovered {
		if !expected[svc] {
			t.Errorf("Unexpected service: %s", svc)
		}
		if seen[svc] {
			t.Errorf("Duplicate service found: %s", svc)
		}
		seen[svc] = true
	}

	// Check for missing services
	for svc := range expected {
		if !seen[svc] {
			t.Errorf("Missing expected service: %s", svc)
		}
	}
}

func TestExtractSuffix(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "internal with service",
			path:     "internal/apiserver/model",
			expected: "model",
		},
		{
			name:     "internal with nested path",
			path:     "internal/apiserver/handler/v1",
			expected: "handler/v1",
		},
		{
			name:     "pkg prefix",
			path:     "pkg/api/v1",
			expected: "v1",
		},
		{
			name:     "no known prefix",
			path:     "custom/path/model",
			expected: "model",
		},
		{
			name:     "single segment",
			path:     "model",
			expected: "model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSuffix(tt.path)
			if result != tt.expected {
				t.Errorf("extractSuffix(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestInferDirectoryForService(t *testing.T) {
	// Setup temp directory with cmd structure
	tmpDir := t.TempDir()
	cmdDir := filepath.Join(tmpDir, "cmd")
	os.MkdirAll(filepath.Join(cmdDir, "myapp-apiserver"), 0755)
	os.MkdirAll(filepath.Join(cmdDir, "myapp-admserver"), 0755)
	os.MkdirAll(filepath.Join(cmdDir, "myapp-api"), 0755)

	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	tests := []struct {
		name        string
		baseDir     string
		serviceName string
		expected    string
		expectError bool
	}{
		{
			name:        "empty service returns base",
			baseDir:     "internal/apiserver/model",
			serviceName: "",
			expected:    "internal/apiserver/model",
		},
		{
			name:        "smart replacement",
			baseDir:     "internal/apiserver/model",
			serviceName: "admserver",
			expected:    "internal/admserver/model",
		},
		{
			name:        "smart replacement with nested path",
			baseDir:     "internal/apiserver/handler/v1",
			serviceName: "admserver",
			expected:    "internal/admserver/handler/v1",
		},
		{
			name:        "fallback pattern",
			baseDir:     "internal/pkg/model",
			serviceName: "admserver",
			expected:    "internal/admserver/model",
		},
		{
			name:        "substring collision - api vs apiserver",
			baseDir:     "internal/apiserver/model",
			serviceName: "api",
			expected:    "internal/api/model",
		},
		{
			name:        "pkg prefix paths with known service",
			baseDir:     "pkg/api/v1",
			serviceName: "admserver",
			expected:    "pkg/admserver/v1",
		},
		{
			name:        "invalid service name with path separator",
			baseDir:     "internal/apiserver/model",
			serviceName: "adm/server",
			expectError: true,
		},
		{
			name:        "invalid service name with parent directory reference",
			baseDir:     "internal/apiserver/model",
			serviceName: "../admserver",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Options{}
			result, err := o.InferDirectoryForService(tt.baseDir, tt.serviceName)
			if tt.expectError {
				if err == nil {
					t.Errorf("InferDirectoryForService(%q, %q) expected error, got nil",
						tt.baseDir, tt.serviceName)
				}
				return
			}
			if err != nil {
				t.Fatalf("InferDirectoryForService failed: %v", err)
			}
			if result != tt.expected {
				t.Errorf("InferDirectoryForService(%q, %q) = %q, want %q",
					tt.baseDir, tt.serviceName, result, tt.expected)
			}
		})
	}
}

func TestInferDirectoryForService_DiscoveryFailure(t *testing.T) {
	// Test behavior when cmd/ directory doesn't exist
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	o := &Options{}
	result, err := o.InferDirectoryForService("internal/unknown/model", "admserver")
	if err != nil {
		t.Fatalf("InferDirectoryForService failed: %v", err)
	}

	// Should fall back to pattern-based inference
	expected := "internal/admserver/model"
	if result != expected {
		t.Errorf("InferDirectoryForService with no cmd/ = %q, want %q", result, expected)
	}
}
