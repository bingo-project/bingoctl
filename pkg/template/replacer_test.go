package template

import (
	"os"
	"path/filepath"
	"strings"
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
		{"proto file", "apiserver.proto", true},

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

func TestReplaceInProtoFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test proto file
	testFile := filepath.Join(tmpDir, "apiserver.proto")
	content := `syntax = "proto3";

package v1;

import "google/protobuf/timestamp.proto";

option go_package = "bingo/internal/apiserver/grpc/proto/v1/pb";

service ApiServer {
  rpc Healthz (HealthzRequest) returns (HealthzReply) {}
}
`
	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Replace
	r := NewReplacer(tmpDir, "bingo", "github.com/mycompany/demo", "demo")
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

	// Check go_package replacement
	expected := `option go_package = "github.com/mycompany/demo/internal/apiserver/grpc/proto/v1/pb";`
	if !strings.Contains(resultStr, expected) {
		t.Errorf("go_package not replaced correctly.\nExpected: %s\nGot: %s", expected, resultStr)
	}
}
