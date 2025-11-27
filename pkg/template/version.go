// ABOUTME: Template version management for GitHub-based template fetching
// ABOUTME: Defines default version and ref validation logic
package template

import (
	"regexp"
	"strings"
)

// DefaultTemplateVersion is the recommended template version
const DefaultTemplateVersion = "develop"

var refRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

// isValidRef checks if ref format is valid
// Supports: v1.2.3, main, abc123def, etc.
func isValidRef(ref string) bool {
	if ref == "" {
		return false
	}
	return refRegex.MatchString(ref)
}

// refType returns ref type: tag, branch, commit
func refType(ref string) string {
	if strings.HasPrefix(ref, "v") && strings.Contains(ref, ".") {
		return "tag"
	}

	// Common branch names
	if ref == "main" || ref == "master" || ref == "develop" || strings.HasPrefix(ref, "release/") || strings.HasPrefix(ref, "feature/") {
		return "branch"
	}

	// Default to commit hash
	return "commit"
}
