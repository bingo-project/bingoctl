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
		t.Errorf("Expected %d services, got %d", len(expected), len(discovered))
	}

	for _, svc := range discovered {
		if !expected[svc] {
			t.Errorf("Unexpected service: %s", svc)
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
