package scaffolder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CloneTemplate clones a Git repository to a temporary directory
func CloneTemplate(url string, version string) (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "genesis-template-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Clone options
	options := &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}

	// Clone the repository
	repo, err := git.PlainClone(tempDir, false, options)
	if err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("failed to clone repository: %w", err)
	}

	// If a specific version is requested, check it out
	if version != "" {
		w, err := repo.Worktree()
		if err != nil {
			os.RemoveAll(tempDir)
			return "", fmt.Errorf("failed to get worktree: %w", err)
		}

		// Try as a commit hash first
		hash := plumbing.NewHash(version)
		if len(version) == 40 { // SHA-1 hash length
			err = w.Checkout(&git.CheckoutOptions{
				Hash: hash,
			})
			if err != nil {
				os.RemoveAll(tempDir)
				return "", fmt.Errorf("failed to checkout version %s: %w", version, err)
			}
		} else {
			// Try as a tag first
			var refName plumbing.ReferenceName
			if strings.HasPrefix(version, "refs/") {
				refName = plumbing.ReferenceName(version)
			} else {
				refName = plumbing.ReferenceName("refs/tags/" + version)
			}

			err = w.Checkout(&git.CheckoutOptions{
				Branch: refName,
			})
			if err != nil {
				// Try as a branch
				refName = plumbing.ReferenceName("refs/heads/" + version)
				err = w.Checkout(&git.CheckoutOptions{
					Branch: refName,
				})
				if err != nil {
					os.RemoveAll(tempDir)
					return "", fmt.Errorf("failed to checkout version %s: %w", version, err)
				}
			}
		}
	}

	// Verify template.toml exists
	if _, err := os.Stat(filepath.Join(tempDir, "template.toml")); err != nil {
		os.RemoveAll(tempDir)
		return "", fmt.Errorf("invalid template: template.toml not found")
	}

	return tempDir, nil
}

// CleanupTemplate removes the temporary directory
func CleanupTemplate(dir string) error {
	return os.RemoveAll(dir)
} 