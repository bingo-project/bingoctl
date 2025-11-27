package template

import (
	"testing"
)

func TestBuildDownloadURL(t *testing.T) {
	tests := []struct {
		name     string
		ref      string
		mirror   string
		expected string
	}{
		{
			name:     "tag without mirror",
			ref:      "v1.2.3",
			mirror:   "",
			expected: "https://github.com/bingo-project/bingo/archive/refs/tags/v1.2.3.tar.gz",
		},
		{
			name:     "branch without mirror",
			ref:      "main",
			mirror:   "",
			expected: "https://github.com/bingo-project/bingo/archive/refs/heads/main.tar.gz",
		},
		{
			name:     "tag with mirror",
			ref:      "v1.2.3",
			mirror:   "https://ghproxy.com/",
			expected: "https://ghproxy.com/https://github.com/bingo-project/bingo/archive/refs/tags/v1.2.3.tar.gz",
		},
		{
			name:     "commit hash",
			ref:      "abc123def",
			mirror:   "",
			expected: "https://github.com/bingo-project/bingo/archive/abc123def.tar.gz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fetcher{mirror: tt.mirror}
			result := f.buildDownloadURL(tt.ref)
			if result != tt.expected {
				t.Errorf("buildDownloadURL(%q) = %q, want %q", tt.ref, result, tt.expected)
			}
		})
	}
}
