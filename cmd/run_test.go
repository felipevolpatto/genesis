package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		config      string
		expectError bool
		setup       func(t *testing.T) error
	}{
		{
			name: "valid task",
			args: []string{"run", "test"},
			config: `version = "1.0"
[tasks]
  test = { description = "Run tests", cmd = "echo 'test'" }`,
			expectError: false,
			setup:      func(t *testing.T) error { return nil },
		},
		{
			name: "list task",
			args: []string{"run", "list"},
			config: `version = "1.0"
[tasks]
  test = { description = "Run tests", cmd = "echo 'test'" }`,
			expectError: false,
			setup:      func(t *testing.T) error { return nil },
		},
		{
			name: "invalid task",
			args: []string{"run", "invalid"},
			config: `version = "1.0"
[tasks]
  test = { description = "Run tests", cmd = "echo 'test'" }`,
			expectError: true,
			setup:      func(t *testing.T) error { return nil },
		},
		{
			name: "nonexistent task",
			args: []string{"run", "nonexistent"},
			config: `version = "1.0"
[tasks]
  test = { description = "Run tests", cmd = "echo 'test'" }`,
			expectError: true,
			setup:      func(t *testing.T) error { return nil },
		},
		{
			name:        "missing task name",
			args:        []string{"run"},
			config:      `version = "1.0"`,
			expectError: true,
			setup:      func(t *testing.T) error { return nil },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory
			tempDir := t.TempDir()

			// Save current directory
			currentDir, err := os.Getwd()
			require.NoError(t, err)

			// Create a deferred function to change back to the original directory
			defer func() {
				if err := os.Chdir(currentDir); err != nil {
					t.Errorf("failed to change back to original directory: %v", err)
				}
			}()

			// Change to the temporary directory
			if err := os.Chdir(tempDir); err != nil {
				t.Fatalf("failed to change to temporary directory: %v", err)
			}

			// Create genesis.toml
			err = os.WriteFile("genesis.toml", []byte(tt.config), 0644)
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
		})
	}
}

func TestRunCommandNoConfig(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err)

	// Create a deferred function to change back to the original directory
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Errorf("failed to change back to original directory: %v", err)
		}
	}()

	// Change to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temporary directory: %v", err)
	}

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute command
	rootCmd.SetArgs([]string{"run", "test"})
	err = rootCmd.Execute()
	assert.Error(t, err)
}

func TestRunCommandListTasks(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err)

	// Create a deferred function to change back to the original directory
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Errorf("failed to change back to original directory: %v", err)
		}
	}()

	// Change to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temporary directory: %v", err)
	}

	// Create genesis.toml
	config := `version = "1.0"
[tasks]
  test = { description = "Run tests", cmd = "echo 'test'" }
  build = { description = "Build binary", cmd = "go build" }`

	err = os.WriteFile("genesis.toml", []byte(config), 0644)
	require.NoError(t, err)

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute command
	rootCmd.SetArgs([]string{"run", "list"})
	err = rootCmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Available tasks:")
	assert.Contains(t, buf.String(), "test")
	assert.Contains(t, buf.String(), "build")
} 