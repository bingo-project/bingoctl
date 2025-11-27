// ABOUTME: Template fetcher for downloading and caching bingo project from GitHub
// ABOUTME: Handles tarball download, extraction, and local caching with file locking
package template

import (
	"fmt"
	"time"
)

const (
	githubArchiveBase = "https://github.com/bingo-project/bingo/archive"
	defaultTimeout    = 30 * time.Second
)

// Fetcher handles template downloading and caching
type Fetcher struct {
	cacheDir string        // ~/.bingoctl/templates
	timeout  time.Duration // 30s
	mirror   string        // mirror address from env var
}

// NewFetcher creates a new Fetcher instance
func NewFetcher() *Fetcher {
	// TODO: get cache dir from user home
	// TODO: read mirror from env BINGOCTL_TEMPLATE_MIRROR
	return &Fetcher{
		cacheDir: "", // will implement later
		timeout:  defaultTimeout,
		mirror:   "", // will implement later
	}
}

// buildDownloadURL constructs download URL (supports mirror)
// Examples:
//   - tag: https://github.com/.../archive/refs/tags/v1.2.3.tar.gz
//   - branch: https://github.com/.../archive/refs/heads/main.tar.gz
//   - commit: https://github.com/.../archive/{hash}.tar.gz
func (f *Fetcher) buildDownloadURL(ref string) string {
	var url string

	refKind := refType(ref)
	switch refKind {
	case "tag":
		url = fmt.Sprintf("%s/refs/tags/%s.tar.gz", githubArchiveBase, ref)
	case "branch":
		url = fmt.Sprintf("%s/refs/heads/%s.tar.gz", githubArchiveBase, ref)
	case "commit":
		url = fmt.Sprintf("%s/%s.tar.gz", githubArchiveBase, ref)
	}

	if f.mirror != "" {
		return f.mirror + url
	}

	return url
}
