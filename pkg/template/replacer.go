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
	".go":    true,
	".mod":   true,
	".sum":   true,
	".proto": true,

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
	"Makefile":       true,
	"Dockerfile":     true,
	"Bin.Dockerfile": true,
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
// Only performs replacements if newModule is not empty
func (r *Replacer) replaceInFile(path string) error {
	// Skip replacement if no new module name specified
	if r.newModule == "" {
		return nil
	}

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

	// 4. Makefile-style assignments: ROOT_PACKAGE=bingo -> ROOT_PACKAGE={newModule}
	// This handles cases like "ROOT_PACKAGE=bingo" or "REGISTRY_PREFIX = bingo"
	str = strings.ReplaceAll(str, "="+r.oldModule, "="+r.newModule)

	// 5. Makefile-style conditional assignments: REGISTRY_PREFIX ?= bingo -> REGISTRY_PREFIX ?= {newModule}
	str = strings.ReplaceAll(str, "?= "+r.oldModule, "?= "+r.newModule)

	// Write back
	err = os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}

// renameRules defines directory rename mappings
// Only these explicitly listed directories will be renamed
var renameRules = map[string]string{
	// cmd directories
	"cmd/bingo-apiserver": "cmd/{app}-apiserver",
	"cmd/bingo-admserver": "cmd/{app}-admserver",
	"cmd/bingo-bot":       "cmd/{app}-bot",
	"cmd/bingo-scheduler": "cmd/{app}-scheduler",
	"cmd/bingoctl":        "cmd/{app}ctl",
	// internal directories
	"internal/bingoctl": "internal/{app}ctl",
	// build/docker directories
	"build/docker/bingo-apiserver": "build/docker/{app}-apiserver",
	"build/docker/bingo-admserver": "build/docker/{app}-admserver",
	"build/docker/bingo-bot":       "build/docker/{app}-bot",
	"build/docker/bingo-scheduler": "build/docker/{app}-scheduler",
	"build/docker/bingoctl":        "build/docker/{app}ctl",
}

// configFileRenameRules defines config file rename mappings
var configFileRenameRules = map[string]string{
	"configs/bingo-apiserver.example.yaml": "configs/{app}-apiserver.example.yaml",
	"configs/bingo-admserver.example.yaml": "configs/{app}-admserver.example.yaml",
	"configs/bingo-bot.example.yaml":       "configs/{app}-bot.example.yaml",
	"configs/bingo-scheduler.example.yaml": "configs/{app}-scheduler.example.yaml",
	"configs/bingoctl.example.yaml":        "configs/{app}ctl.example.yaml",
}

// RenameDirs renames directories according to explicit rules
// Only renames directories that still exist (after service filtering)
func (r *Replacer) RenameDirs() error {
	for oldPath, newPathTemplate := range renameRules {
		// Replace {app} placeholder
		newPath := strings.ReplaceAll(newPathTemplate, "{app}", r.appName)

		oldFullPath := filepath.Join(r.targetDir, oldPath)
		newFullPath := filepath.Join(r.targetDir, newPath)

		// Skip if old path doesn't exist (may be filtered out)
		if !fileExists(oldFullPath) {
			continue
		}

		// Rename
		err := os.Rename(oldFullPath, newFullPath)
		if err != nil {
			return fmt.Errorf("failed to rename %s to %s: %w", oldPath, newPath, err)
		}
	}

	return nil
}

// RenameConfigFiles renames config files according to explicit rules
// Only renames files that exist (some may be filtered out with services)
func (r *Replacer) RenameConfigFiles() error {
	for oldPath, newPathTemplate := range configFileRenameRules {
		// Replace {app} placeholder
		newPath := strings.ReplaceAll(newPathTemplate, "{app}", r.appName)

		oldFullPath := filepath.Join(r.targetDir, oldPath)
		newFullPath := filepath.Join(r.targetDir, newPath)

		// Skip if old file doesn't exist
		if !fileExists(oldFullPath) {
			continue
		}

		// Rename
		err := os.Rename(oldFullPath, newFullPath)
		if err != nil {
			return fmt.Errorf("failed to rename %s to %s: %w", oldPath, newPath, err)
		}
	}

	return nil
}

// appNameReplacements defines app name replacement patterns
// Order matters: more specific patterns should come first
var appNameReplacements = []struct {
	old string
	new string
}{
	// Service names (more specific, should be first)
	{"bingo-apiserver", "{app}-apiserver"},
	{"bingo-admserver", "{app}-admserver"},
	{"bingo-bot", "{app}-bot"},
	{"bingo-scheduler", "{app}-scheduler"},
	{"bingoctl", "{app}ctl"},
	// Paths
	{"/opt/bingo", "/opt/{app}"},
	{"/etc/bingo", "/etc/{app}"},
	{"/var/log/bingo", "/var/log/{app}"},
	{"/data/bingo", "/data/{app}"},
	// Generic app name (catches env vars, network names, etc.)
	{"=bingo\n", "={app}\n"},
	{"=bingo\r\n", "={app}\r\n"},
	// Network references in yaml
	{"  bingo:\n", "  {app}:\n"},
	{"- bingo\n", "- {app}\n"},
	{"- bingo\r\n", "- {app}\r\n"},
}

// ReplaceBingoConfig replaces values in .bingo.example.yaml
// Replaces rootPackage and database fields with the new module/app name
func (r *Replacer) ReplaceBingoConfig() error {
	configPath := filepath.Join(r.targetDir, ".bingo.example.yaml")
	if !fileExists(configPath) {
		return nil
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", configPath, err)
	}

	str := string(content)

	// Replace rootPackage: bingo -> rootPackage: {newModule}
	if r.newModule != "" {
		str = strings.ReplaceAll(str, "rootPackage: bingo", "rootPackage: "+r.newModule)
	}

	// Replace database: bingo -> database: {appName}
	str = strings.ReplaceAll(str, "database: bingo", "database: "+r.appName)

	err = os.WriteFile(configPath, []byte(str), 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s: %w", configPath, err)
	}

	return nil
}

// ReplaceAppName replaces app name references in files
// Should be called after ReplaceModuleName to avoid conflicts with module paths
func (r *Replacer) ReplaceAppName() error {
	return filepath.WalkDir(r.targetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if r.shouldReplaceFile(path) {
			return r.replaceAppNameInFile(path)
		}

		return nil
	})
}

// protectedPatterns are patterns that should not be modified during replacement
// These are external dependencies that should remain unchanged
var protectedPatterns = []string{
	"github.com/bingo-project/bingoctl",
}

// placeholder is used to temporarily protect patterns from replacement
const placeholder = "___PROTECTED_BINGOCTL___"

// replaceAppNameInFile replaces app name patterns in a single file
func (r *Replacer) replaceAppNameInFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	str := string(content)

	// Step 1: Protect external dependencies with placeholder
	for i, pattern := range protectedPatterns {
		ph := fmt.Sprintf("%s_%d_", placeholder, i)
		str = strings.ReplaceAll(str, pattern, ph)
	}

	// Step 2: Apply all replacement patterns
	for _, repl := range appNameReplacements {
		newPattern := strings.ReplaceAll(repl.new, "{app}", r.appName)
		str = strings.ReplaceAll(str, repl.old, newPattern)
	}

	// Step 3: Restore protected patterns
	for i, pattern := range protectedPatterns {
		ph := fmt.Sprintf("%s_%d_", placeholder, i)
		str = strings.ReplaceAll(str, ph, pattern)
	}

	err = os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}
