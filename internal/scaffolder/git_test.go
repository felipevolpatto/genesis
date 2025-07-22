package scaffolder

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) (string, plumbing.Hash) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Initialize Git repository
	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	// Configure Git for the test
	cfg, err := repo.Config()
	require.NoError(t, err)

	cfg.User.Name = "Test Author"
	cfg.User.Email = "test@example.com"

	err = repo.SetConfig(cfg)
	require.NoError(t, err)

	// Create template.toml
	err = os.WriteFile(filepath.Join(tempDir, "template.toml"), []byte("version = \"1.0\""), 0644)
	require.NoError(t, err)

	// Get the worktree
	w, err := repo.Worktree()
	require.NoError(t, err)

	// Add all files
	_, err = w.Add(".")
	require.NoError(t, err)

	// Create a commit
	author := &object.Signature{
		Name:  "Test Author",
		Email: "test@example.com",
		When:  time.Now(),
	}

	commit, err := w.Commit("Initial commit", &git.CommitOptions{
		Author: author,
	})
	require.NoError(t, err)

	// Create a tag
	_, err = repo.CreateTag("v1.0.0", commit, &git.CreateTagOptions{
		Tagger:  author,
		Message: "Version 1.0.0",
	})
	require.NoError(t, err)

	return tempDir, commit
}

func TestCloneTemplate(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		setup       func(t *testing.T) (string, plumbing.Hash)
		expectError bool
	}{
		{
			name:    "clone without version",
			version: "",
			setup: func(t *testing.T) (string, plumbing.Hash) {
				dir, hash := setupTestRepo(t)
				return dir, hash
			},
			expectError: false,
		},
		{
			name:    "clone with commit hash",
			version: "4eeee43170429c82eb2fc4eeef52fbd2d8b2d0ca",
			setup: func(t *testing.T) (string, plumbing.Hash) {
				dir, hash := setupTestRepo(t)
				return dir, hash
			},
			expectError: true,
		},
		{
			name:    "clone with tag",
			version: "v1.0.0",
			setup: func(t *testing.T) (string, plumbing.Hash) {
				dir, hash := setupTestRepo(t)
				return dir, hash
			},
			expectError: false,
		},
		{
			name:    "invalid repository",
			version: "",
			setup: func(t *testing.T) (string, plumbing.Hash) {
				return "/nonexistent/path", plumbing.Hash{}
			},
			expectError: true,
		},
		{
			name:    "invalid version",
			version: "nonexistent",
			setup: func(t *testing.T) (string, plumbing.Hash) {
				dir, hash := setupTestRepo(t)
				return dir, hash
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceDir, _ := tt.setup(t)
			tempDir, err := CloneTemplate(sourceDir, tt.version)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			defer func() {
				if err := CleanupTemplate(tempDir); err != nil {
					t.Errorf("failed to cleanup template directory: %v", err)
				}
			}()

			// Check if template.toml exists
			_, err = os.Stat(filepath.Join(tempDir, "template.toml"))
			assert.NoError(t, err)
		})
	}
}

func TestCleanupTemplate(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a file in the directory
	err := os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("test"), 0644)
	require.NoError(t, err)

	// Clean up the directory
	err = CleanupTemplate(tempDir)
	assert.NoError(t, err)

	// Check if directory was removed
	_, err = os.Stat(tempDir)
	assert.True(t, os.IsNotExist(err))
}

func TestCloneTemplateWithoutTemplateToml(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Initialize Git repository
	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	// Configure Git for the test
	cfg, err := repo.Config()
	require.NoError(t, err)

	cfg.User.Name = "Test Author"
	cfg.User.Email = "test@example.com"

	err = repo.SetConfig(cfg)
	require.NoError(t, err)

	// Create a dummy file
	err = os.WriteFile(filepath.Join(tempDir, "dummy.txt"), []byte("test"), 0644)
	require.NoError(t, err)

	// Get the worktree
	w, err := repo.Worktree()
	require.NoError(t, err)

	// Add all files
	_, err = w.Add(".")
	require.NoError(t, err)

	// Create a commit
	author := &object.Signature{
		Name:  "Test Author",
		Email: "test@example.com",
		When:  time.Now(),
	}

	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: author,
	})
	require.NoError(t, err)

	// Try to clone the repository
	clonedDir, err := CloneTemplate(tempDir, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template.toml not found")

	if clonedDir != "" {
		if err := CleanupTemplate(clonedDir); err != nil {
			t.Errorf("failed to cleanup template directory: %v", err)
		}
	}
}

func TestCloneTemplateWithDifferentVersions(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		setup       func(t *testing.T) (string, plumbing.Hash)
		expectError bool
	}{
		{
			name:    "version 4eeee43170429c82eb2fc4eeef52fbd2d8b2d0ca",
			version: "4eeee43170429c82eb2fc4eeef52fbd2d8b2d0ca",
			setup: func(t *testing.T) (string, plumbing.Hash) {
				dir, hash := setupTestRepo(t)
				return dir, hash
			},
			expectError: true,
		},
		{
			name:    "version v1.0.0",
			version: "v1.0.0",
			setup: func(t *testing.T) (string, plumbing.Hash) {
				dir, hash := setupTestRepo(t)
				return dir, hash
			},
			expectError: false,
		},
		{
			name:    "version nonexistent",
			version: "nonexistent",
			setup: func(t *testing.T) (string, plumbing.Hash) {
				dir, hash := setupTestRepo(t)
				return dir, hash
			},
			expectError: true,
		},
		{
			name:    "version 12345678",
			version: "12345678",
			setup: func(t *testing.T) (string, plumbing.Hash) {
				dir, hash := setupTestRepo(t)
				return dir, hash
			},
			expectError: true,
		},
		{
			name:    "version refs/tags/v1.0.0",
			version: "refs/tags/v1.0.0",
			setup: func(t *testing.T) (string, plumbing.Hash) {
				dir, hash := setupTestRepo(t)
				return dir, hash
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourceDir, _ := tt.setup(t)
			tempDir, err := CloneTemplate(sourceDir, tt.version)

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			defer func() {
				if err := CleanupTemplate(tempDir); err != nil {
					t.Errorf("failed to cleanup template directory: %v", err)
				}
			}()

			// Check if template.toml exists
			_, err = os.Stat(filepath.Join(tempDir, "template.toml"))
			assert.NoError(t, err)
		})
	}
} 