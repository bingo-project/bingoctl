package template

import "testing"

func TestDefaultTemplateVersion(t *testing.T) {
	if DefaultTemplateVersion == "" {
		t.Error("DefaultTemplateVersion should not be empty")
	}

	// Should be a valid semver tag
	if DefaultTemplateVersion[0] != 'v' {
		t.Errorf("DefaultTemplateVersion should start with 'v', got: %s", DefaultTemplateVersion)
	}
}

func TestIsValidRef(t *testing.T) {
	tests := []struct {
		name  string
		ref   string
		valid bool
	}{
		{"valid tag", "v1.2.3", true},
		{"valid branch", "main", true},
		{"valid commit", "abc123def", true},
		{"empty string", "", false},
		{"invalid chars", "v1.2.3@#$", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidRef(tt.ref)
			if result != tt.valid {
				t.Errorf("isValidRef(%q) = %v, want %v", tt.ref, result, tt.valid)
			}
		})
	}
}

func TestRefType(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		expected string
	}{
		{"semver tag", "v1.2.3", "tag"},
		{"tag with prefix", "v0.1.0", "tag"},
		{"branch", "main", "branch"},
		{"branch", "develop", "branch"},
		{"commit hash", "abc123def", "commit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := refType(tt.ref)
			if result != tt.expected {
				t.Errorf("refType(%q) = %v, want %v", tt.ref, result, tt.expected)
			}
		})
	}
}
