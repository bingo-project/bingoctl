// ABOUTME: Tests for create command functionality
// ABOUTME: Validates service list computation logic
package create

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestCleanupTemplateFiles(t *testing.T) {
	// Create temporary directory
	tmpDir := t.TempDir()

	// Create docs directory and files
	docsDir := tmpDir + "/docs"
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		t.Fatalf("failed to create docs dir: %v", err)
	}
	if err := os.WriteFile(docsDir+"/test.md", []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Create README.md
	readmePath := tmpDir + "/README.md"
	if err := os.WriteFile(readmePath, []byte("# Bingo"), 0644); err != nil {
		t.Fatalf("failed to create README: %v", err)
	}

	// Create CHANGELOG.md
	changelogPath := tmpDir + "/CHANGELOG.md"
	if err := os.WriteFile(changelogPath, []byte("# Changelog"), 0644); err != nil {
		t.Fatalf("failed to create CHANGELOG: %v", err)
	}

	// Create README.zh-CN.md
	readmeZhPath := tmpDir + "/README.zh-CN.md"
	if err := os.WriteFile(readmeZhPath, []byte("# Bingo 中文"), 0644); err != nil {
		t.Fatalf("failed to create README.zh-CN.md: %v", err)
	}

	// Create CreateOptions and call cleanup
	o := &CreateOptions{AppName: "myapp"}
	if err := o.cleanupTemplateFiles(tmpDir); err != nil {
		t.Fatalf("cleanupTemplateFiles failed: %v", err)
	}

	// Verify docs directory is removed
	if _, err := os.Stat(docsDir); !os.IsNotExist(err) {
		t.Error("docs directory should be removed")
	}

	// Verify new README exists
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("failed to read new README: %v", err)
	}

	// Verify README content
	if !strings.Contains(string(content), "bingoctl") {
		t.Error("README should mention bingoctl")
	}
	if !strings.Contains(string(content), "bingo") {
		t.Error("README should mention bingo")
	}
	if !strings.Contains(string(content), "myapp") {
		t.Error("README should contain project name")
	}

	// Verify CHANGELOG is removed
	if _, err := os.Stat(changelogPath); !os.IsNotExist(err) {
		t.Error("CHANGELOG.md should be removed")
	}

	// Verify README.zh-CN.md is removed
	if _, err := os.Stat(readmeZhPath); !os.IsNotExist(err) {
		t.Error("README.zh-CN.md should be removed")
	}

	// Verify README contains documentation URL
	if !strings.Contains(string(content), "bingoctl.dev") {
		t.Error("README should contain bingoctl.dev documentation URL")
	}

	// Verify README contains build instructions
	if !strings.Contains(string(content), "make build") {
		t.Error("README should contain make build instruction")
	}
	if !strings.Contains(string(content), "_output/platforms") {
		t.Error("README should contain _output/platforms path")
	}
}

func TestCopyExampleConfigs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create configs directory with example files
	configsDir := tmpDir + "/configs"
	if err := os.MkdirAll(configsDir, 0755); err != nil {
		t.Fatalf("failed to create configs dir: %v", err)
	}

	// Create example config files
	examples := []string{"app-apiserver.example.yaml", "app-admserver.example.yaml"}
	for _, name := range examples {
		content := "# " + name
		if err := os.WriteFile(configsDir+"/"+name, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create %s: %v", name, err)
		}
	}

	// Call copyExampleConfigs
	o := &CreateOptions{AppName: "myapp"}
	if err := o.copyExampleConfigs(tmpDir); err != nil {
		t.Fatalf("copyExampleConfigs failed: %v", err)
	}

	// Verify configs were copied to root without .example suffix
	expected := []string{"app-apiserver.yaml", "app-admserver.yaml"}
	for _, name := range expected {
		path := tmpDir + "/" + name
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("%s should exist in root directory", name)
		}
	}
}

func TestComputeServiceList(t *testing.T) {
	tests := []struct {
		name        string
		all         bool
		services    []string
		noServices  []string
		addServices []string
		expected    []string
	}{
		{
			name:     "all flag returns all services",
			all:      true,
			expected: []string{"apiserver", "ctl", "admserver", "bot", "scheduler"},
		},
		{
			name:     "explicit services override",
			services: []string{"apiserver", "bot"},
			expected: []string{"apiserver", "bot"},
		},
		{
			name:     "services none",
			services: []string{"none"},
			expected: []string{},
		},
		{
			name:     "no flags uses defaults",
			expected: []string{"apiserver"},
		},
		{
			name:       "exclude service",
			noServices: []string{"apiserver"},
			expected:   []string{},
		},
		{
			name:        "add services",
			addServices: []string{"bot", "scheduler"},
			expected:    []string{"apiserver", "bot", "scheduler"},
		},
		{
			name:        "combined exclude and add",
			noServices:  []string{"apiserver"},
			addServices: []string{"admserver"},
			expected:    []string{"admserver"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &CreateOptions{
				All:         tt.all,
				Services:    tt.services,
				NoServices:  tt.noServices,
				AddServices: tt.addServices,
			}
			result := o.computeServiceList()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("computeServiceList() = %v, want %v", result, tt.expected)
			}
		})
	}
}
