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
