package template

import (
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
