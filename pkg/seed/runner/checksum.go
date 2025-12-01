// ABOUTME: Calculates checksum of seeder files for cache invalidation
// ABOUTME: Uses SHA256 hash of all .go files in seeder directory
package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// CalculateChecksum calculates a SHA256 checksum of all .go files in the directory.
func CalculateChecksum(dir string) (string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return "", nil
	}

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	if len(files) == 0 {
		return "", nil
	}

	sort.Strings(files)

	h := sha256.New()
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(h, f); err != nil {
			f.Close()
			return "", err
		}
		f.Close()
		h.Write([]byte(file))
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// CalculatePathHash calculates a short hash of the project path.
func CalculatePathHash(path string) string {
	h := sha256.Sum256([]byte(path))
	return hex.EncodeToString(h[:])[:8]
}
