package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestTemplate(t *testing.T) string {
	// Create a temporary directory for the template
	templateDir := t.TempDir()

	// Initialize Git repository
	repo, err := git.PlainInit(templateDir, false)
	require.NoError(t, err)

	// Create template.toml
	templateConfig := `version = "1.0"

[vars]
  name = { prompt = "Enter name:", default = "test" }
  description = { prompt = "Enter description:", default = "A test project" }

[hooks]
  post = ["echo 'test' > post-hook.txt"]`

	err = os.WriteFile(filepath.Join(templateDir, "template.toml"), []byte(templateConfig), 0644)
	require.NoError(t, err)

	// Create a template file
	mainTemplate := `package main

func main() {
	println("Hello, {{ .name }}!")
}`

	err = os.WriteFile(filepath.Join(templateDir, "main.go.tmpl"), []byte(mainTemplate), 0644)
	require.NoError(t, err)

	// Add and commit files
	w, err := repo.Worktree()
	require.NoError(t, err)

	_, err = w.Add(".")
	require.NoError(t, err)

	commit, err := w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test Author",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	// Create a tag
	_, err = repo.CreateTag("v1.0.0", commit, &git.CreateTagOptions{
		Message: "Version 1.0.0",
		Tagger: &object.Signature{
			Name:  "Test Author",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	return templateDir
}

func TestNewCommand(t *testing.T) {
	templateDir := setupTestTemplate(t)
	projectDir := t.TempDir()

	tests := []struct {
		name        string
		args        []string
		expectError bool
		setup       func(t *testing.T) error
	}{
		{
			name: "basic new project",
			args: []string{
				"new",
				"test-project",
				"--template", templateDir,
				"--yes", // Skip prompts
			},
			expectError: false,
			setup:      func(t *testing.T) error { return nil },
		},
		{
			name: "missing project name",
			args: []string{
				"new",
				"--template", templateDir,
			},
			expectError: true,
			setup:      func(t *testing.T) error { return nil },
		},
		{
			name: "missing template",
			args: []string{
				"new",
				"test-project",
			},
			expectError: true,
			setup:      func(t *testing.T) error {
				templateURL = "" // Reset template URL
				return nil
			},
		},
		{
			name: "invalid template path",
			args: []string{
				"new",
				"test-project",
				"--template", "/nonexistent/path",
			},
			expectError: true,
			setup:      func(t *testing.T) error { return nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Change to the test directory
			currentDir, err := os.Getwd()
			require.NoError(t, err)
			defer os.Chdir(currentDir)

			err = os.Chdir(projectDir)
			require.NoError(t, err)

			// Run setup if any
			err = tt.setup(t)
			require.NoError(t, err)

			// Create a buffer to capture output
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			// Execute command
			rootCmd.SetArgs(tt.args)
			err = rootCmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Check if project was created
			projectPath := filepath.Join(projectDir, "test-project")
			_, err = os.Stat(projectPath)
			assert.NoError(t, err)

			// Check if files were created
			files := []string{
				"main.go",
				"genesis.toml",
			}

			for _, file := range files {
				_, err = os.Stat(filepath.Join(projectPath, file))
				assert.NoError(t, err)
			}

			// Check if post-hook was executed
			_, err = os.Stat(filepath.Join(projectPath, "post-hook.txt"))
			assert.NoError(t, err)
		})
	}
}

func TestNewCommandWithVersion(t *testing.T) {
	templateDir := setupTestTemplate(t)
	projectDir := t.TempDir()

	// Test with version flag
	args := []string{
		"new",
		"test-project",
		"--template", templateDir,
		"--version", "v1.0.0",
		"--yes",
	}

	// Change to the test directory
	currentDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(currentDir)

	err = os.Chdir(projectDir)
	require.NoError(t, err)

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute command
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()

	// Since we're using a local directory, version flag should be ignored
	assert.NoError(t, err)

	// Check if project was created
	projectPath := filepath.Join(projectDir, "test-project")
	_, err = os.Stat(projectPath)
	assert.NoError(t, err)
} 