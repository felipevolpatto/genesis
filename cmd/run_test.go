package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/felipevolpatto/genesis/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestProject(t *testing.T) string {
	// Create a temporary directory
	projectDir := t.TempDir()

	// Create genesis.toml
	config := `version = "1.0"

[project]
  template_url = "https://github.com/example/template"
  template_version = "v1.0.0"

[tasks]
  echo = { description = "Echo test", cmd = "echo 'test'" }
  invalid = { description = "Invalid command", cmd = "invalidcommand" }
  list = { description = "List files", cmd = "ls -la" }`

	err := os.WriteFile(filepath.Join(projectDir, "genesis.toml"), []byte(config), 0644)
	require.NoError(t, err)

	return projectDir
}

func TestRunCommand(t *testing.T) {
	projectDir := setupTestProject(t)

	tests := []struct {
		name        string
		args        []string
		expectError bool
		contains    string
	}{
		{
			name:        "valid task",
			args:        []string{"run", "echo"},
			expectError: false,
			contains:    "test",
		},
		{
			name:        "list task",
			args:        []string{"run", "list"},
			expectError: false,
			contains:    "Available tasks:",
		},
		{
			name:        "invalid task",
			args:        []string{"run", "invalid"},
			expectError: true,
			contains:    "command not found",
		},
		{
			name:        "nonexistent task",
			args:        []string{"run", "nonexistent"},
			expectError: true,
			contains:    "not found",
		},
		{
			name:        "missing task name",
			args:        []string{"run"},
			expectError: true,
			contains:    "",
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

			// Create a buffer to capture output
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			// Execute command and capture output
			var output string
			if tt.contains != "" {
				output = testutil.CaptureOutput(func() {
					rootCmd.SetArgs(tt.args)
					err = rootCmd.Execute()
				})
			} else {
				rootCmd.SetArgs(tt.args)
				err = rootCmd.Execute()
			}

			if tt.expectError {
				assert.Error(t, err)
				if tt.contains != "" {
					assert.Contains(t, output+buf.String(), tt.contains)
				}
				return
			}

			assert.NoError(t, err)

			// Check output based on task
			if tt.contains != "" {
				assert.Contains(t, output+buf.String(), tt.contains)
			}
		})
	}
}

func TestRunCommandNoConfig(t *testing.T) {
	// Create a temporary directory without genesis.toml
	projectDir := t.TempDir()

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
	rootCmd.SetArgs([]string{"run", "test"})
	err = rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no genesis.toml found")
}

func TestRunCommandListTasks(t *testing.T) {
	projectDir := setupTestProject(t)

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

	// Execute command and capture output
	output := testutil.CaptureOutput(func() {
		rootCmd.SetArgs([]string{"run", "list"})
		err = rootCmd.Execute()
	})

	assert.NoError(t, err)
	assert.Contains(t, output+buf.String(), "echo")
	assert.Contains(t, output+buf.String(), "Echo test")
	assert.Contains(t, output+buf.String(), "invalid")
	assert.Contains(t, output+buf.String(), "Invalid command")
} 