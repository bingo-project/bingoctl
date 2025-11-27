// ABOUTME: Package name and directory name replacer for template processing
// ABOUTME: Handles module name substitution and directory renaming
package template

import (
	"fmt"
	"io/fs"
	"os"
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
