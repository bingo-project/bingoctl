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
