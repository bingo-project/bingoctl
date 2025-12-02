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

func TestRenameConfigFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create configs directory
	configsDir := filepath.Join(tmpDir, "configs")
	err := os.MkdirAll(configsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create configs directory: %v", err)
	}

	// Create test config files (simulating bingo template)
	testFiles := []string{
		"bingo-apiserver.example.yaml",
		"bingo-admserver.example.yaml",
		"bingo-bot.example.yaml",
		"bingo-scheduler.example.yaml",
		"bingoctl.example.yaml",
		"promtail.example.yaml", // Should not be renamed
	}

	for _, f := range testFiles {
		filePath := filepath.Join(configsDir, f)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", f, err)
		}
	}

	// Rename config files
	r := NewReplacer(tmpDir, "bingo", "github.com/mycompany/demo", "demo")
	err = r.RenameConfigFiles()
	if err != nil {
		t.Fatalf("RenameConfigFiles failed: %v", err)
	}

	// Verify renamed files
	expectedFiles := []string{
		"demo-apiserver.example.yaml",
		"demo-admserver.example.yaml",
		"demo-bot.example.yaml",
		"demo-scheduler.example.yaml",
		"democtl.example.yaml",
		"promtail.example.yaml", // Should remain unchanged
	}

	for _, f := range expectedFiles {
		filePath := filepath.Join(configsDir, f)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", f)
		}
	}

	// Verify old files are gone
	oldFiles := []string{
		"bingo-apiserver.example.yaml",
		"bingo-admserver.example.yaml",
		"bingo-bot.example.yaml",
		"bingo-scheduler.example.yaml",
		"bingoctl.example.yaml",
	}

	for _, f := range oldFiles {
		filePath := filepath.Join(configsDir, f)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Errorf("Old file %s should not exist", f)
		}
	}
}

func TestReplaceAppName(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test Dockerfile
	dockerDir := filepath.Join(tmpDir, "build", "docker", "demo-apiserver")
	err := os.MkdirAll(dockerDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create docker directory: %v", err)
	}

	dockerfileContent := `FROM alpine:3.22
WORKDIR /opt/bingo
COPY bingo-apiserver bin/
ENTRYPOINT ["/opt/bingo/bin/bingo-apiserver"]
`
	err = os.WriteFile(filepath.Join(dockerDir, "Dockerfile"), []byte(dockerfileContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create Dockerfile: %v", err)
	}

	// Create test docker-compose.yaml
	deploymentsDir := filepath.Join(tmpDir, "deployments", "docker")
	err = os.MkdirAll(deploymentsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create deployments directory: %v", err)
	}

	composeContent := `networks:
  bingo:
    driver: bridge

services:
  apiserver:
    build:
      dockerfile: build/docker/bingo-apiserver/Dockerfile
    networks:
      - bingo
    volumes:
      - ${DATA_PATH_HOST}/config:/etc/bingo
      - ${DATA_PATH_HOST}/data/bingo:/opt/bingo/storage/public
`
	err = os.WriteFile(filepath.Join(deploymentsDir, "docker-compose.yaml"), []byte(composeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create docker-compose.yaml: %v", err)
	}

	// Create test .env.example
	envContent := `REGISTRY_PREFIX=bingo
APP_NAME=bingo
MYSQL_DATABASE=bingo
`
	err = os.WriteFile(filepath.Join(deploymentsDir, ".env.example"), []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env.example: %v", err)
	}

	// Replace app name
	r := NewReplacer(tmpDir, "bingo", "github.com/mycompany/demo", "demo")
	err = r.ReplaceAppName()
	if err != nil {
		t.Fatalf("ReplaceAppName failed: %v", err)
	}

	// Verify Dockerfile
	dockerfileResult, _ := os.ReadFile(filepath.Join(dockerDir, "Dockerfile"))
	dockerfileStr := string(dockerfileResult)
	if !strings.Contains(dockerfileStr, "/opt/demo") {
		t.Error("Dockerfile: /opt/bingo should be replaced with /opt/demo")
	}
	if !strings.Contains(dockerfileStr, "demo-apiserver") {
		t.Error("Dockerfile: bingo-apiserver should be replaced with demo-apiserver")
	}

	// Verify docker-compose.yaml
	composeResult, _ := os.ReadFile(filepath.Join(deploymentsDir, "docker-compose.yaml"))
	composeStr := string(composeResult)
	if strings.Contains(composeStr, "bingo:") {
		t.Error("docker-compose.yaml: network name 'bingo:' should be replaced")
	}
	if strings.Contains(composeStr, "/etc/bingo") {
		t.Error("docker-compose.yaml: /etc/bingo should be replaced with /etc/demo")
	}
	if strings.Contains(composeStr, "build/docker/bingo-apiserver") {
		t.Error("docker-compose.yaml: build/docker/bingo-apiserver should be replaced")
	}

	// Verify .env.example
	envResult, _ := os.ReadFile(filepath.Join(deploymentsDir, ".env.example"))
	envStr := string(envResult)
	if strings.Contains(envStr, "REGISTRY_PREFIX=bingo") {
		t.Error(".env.example: REGISTRY_PREFIX=bingo should be replaced")
	}
	if strings.Contains(envStr, "APP_NAME=bingo") {
		t.Error(".env.example: APP_NAME=bingo should be replaced")
	}
	if strings.Contains(envStr, "MYSQL_DATABASE=bingo") {
		t.Error(".env.example: MYSQL_DATABASE=bingo should be replaced")
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
