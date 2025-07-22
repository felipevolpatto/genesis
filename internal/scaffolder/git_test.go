package scaffolder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) (string, string, string) {
	// Create source repository
	sourceDir := t.TempDir()
	repo, err := git.PlainInit(sourceDir, false)
	require.NoError(t, err)

	// Create and commit a template.toml file
	templatePath := filepath.Join(sourceDir, "template.toml")
	err = os.WriteFile(templatePath, []byte("version = \"1.0\""), 0644)
	require.NoError(t, err)

	w, err := repo.Worktree()
	require.NoError(t, err)

	_, err = w.Add("template.toml")
	require.NoError(t, err)

	commit, err := w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Author",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	// Create a tag
	tagName := "v1.0.0"
	_, err = repo.CreateTag(tagName, commit, &git.CreateTagOptions{
		Message: "Version 1.0.0",
		Tagger: &object.Signature{
			Name:  "Test Author",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	// Create a destination directory
	destDir := t.TempDir()

	return sourceDir, destDir, commit.String()
}

func TestCloneTemplate(t *testing.T) {
	sourceDir, _, commitHash := setupTestRepo(t)

	tests := []struct {
		name        string
		url         string
		version     string
		expectError bool
	}{
		{
			name:        "clone without version",
			url:         sourceDir,
			version:     "",
			expectError: false,
		},
		{
			name:        "clone with commit hash",
			url:         sourceDir,
			version:     commitHash,
			expectError: false,
		},
		{
			name:        "clone with tag",
			url:         sourceDir,
			version:     "v1.0.0",
			expectError: false,
		},
		{
			name:        "invalid repository",
			url:         "/nonexistent/repo",
			version:     "",
			expectError: true,
		},
		{
			name:        "invalid version",
			url:         sourceDir,
			version:     "nonexistent-tag",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := CloneTemplate(tt.url, tt.version)
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			defer CleanupTemplate(tempDir)

			// Verify template.toml exists
			_, err = os.Stat(filepath.Join(tempDir, "template.toml"))
			assert.NoError(t, err)

			// If version specified, verify it was checked out
			if tt.version != "" {
				repo, err := git.PlainOpen(tempDir)
				require.NoError(t, err)

				head, err := repo.Head()
				require.NoError(t, err)

				if tt.version == commitHash {
					assert.Equal(t, tt.version, head.Hash().String())
				} else if tt.version == "v1.0.0" {
					// Get the commit that the tag points to
					ref, err := repo.Tag(tt.version)
					require.NoError(t, err)
					tagObj, err := repo.TagObject(ref.Hash())
					require.NoError(t, err)
					assert.Equal(t, tagObj.Target, head.Hash())
				}
			}
		})
	}
}

func TestCleanupTemplate(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create some files
	files := []string{"file1.txt", "file2.txt"}
	for _, file := range files {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte("test"), 0644)
		require.NoError(t, err)
	}

	// Test cleanup
	err := CleanupTemplate(tempDir)
	assert.NoError(t, err)

	// Verify directory was removed
	_, err = os.Stat(tempDir)
	assert.True(t, os.IsNotExist(err))

	// Test cleanup of non-existent directory
	err = CleanupTemplate("/nonexistent/dir")
	assert.NoError(t, err) // Should not error on non-existent directory
}

func TestCloneTemplateWithoutTemplateToml(t *testing.T) {
	// Create a repository without template.toml
	sourceDir := t.TempDir()
	repo, err := git.PlainInit(sourceDir, false)
	require.NoError(t, err)

	// Create a dummy file
	dummyPath := filepath.Join(sourceDir, "dummy.txt")
	err = os.WriteFile(dummyPath, []byte("test"), 0644)
	require.NoError(t, err)

	w, err := repo.Worktree()
	require.NoError(t, err)

	_, err = w.Add("dummy.txt")
	require.NoError(t, err)

	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Author",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	// Try to clone
	_, err = CloneTemplate(sourceDir, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "template.toml not found")
}

func TestCloneTemplateWithDifferentVersions(t *testing.T) {
	sourceDir, _, commitHash := setupTestRepo(t)

	// Test different version formats
	versions := []struct {
		version     string
		expectError bool
	}{
		{version: commitHash, expectError: false},      // Commit hash
		{version: "v1.0.0", expectError: false},        // Tag
		{version: "nonexistent", expectError: true},    // Invalid ref
		{version: "12345678", expectError: true},       // Invalid hash
		{version: "refs/tags/v1.0.0", expectError: false}, // Full ref
	}

	for _, v := range versions {
		t.Run("version "+v.version, func(t *testing.T) {
			tempDir, err := CloneTemplate(sourceDir, v.version)
			if v.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			defer CleanupTemplate(tempDir)

			// Verify template.toml exists
			_, err = os.Stat(filepath.Join(tempDir, "template.toml"))
			assert.NoError(t, err)

			// Verify version was checked out correctly
			repo, err := git.PlainOpen(tempDir)
			require.NoError(t, err)

			head, err := repo.Head()
			require.NoError(t, err)

			if v.version == commitHash {
				assert.Equal(t, v.version, head.Hash().String())
			} else if v.version == "v1.0.0" || v.version == "refs/tags/v1.0.0" {
				// Get the commit that the tag points to
				ref, err := repo.Tag("v1.0.0")
				require.NoError(t, err)
				tagObj, err := repo.TagObject(ref.Hash())
				require.NoError(t, err)
				assert.Equal(t, tagObj.Target, head.Hash())
			}
		})
	}
} 