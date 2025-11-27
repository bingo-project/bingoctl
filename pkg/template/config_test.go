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
