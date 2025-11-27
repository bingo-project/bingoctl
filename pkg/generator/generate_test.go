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
