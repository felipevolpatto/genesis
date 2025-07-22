package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		setup       func(t *testing.T) string
		expectError bool
		contains    string
	}{
		{
			name: "list templates",
			args: []string{"template", "list"},
			setup: func(t *testing.T) string {
				// Create a temporary directory
				tempDir := t.TempDir()

				// Create a template directory
				templateDir := filepath.Join(tempDir, "templates")
				err := os.MkdirAll(templateDir, 0755)
				require.NoError(t, err)

				// Create a template
				err = os.MkdirAll(filepath.Join(templateDir, "go-cli"), 0755)
				require.NoError(t, err)

				// Create template.toml
				templateConfig := `version = "1.0"

[vars]
  name = { prompt = "Enter name:", default = "test" }
  description = { prompt = "Enter description:", default = "A test project" }`

				err = os.WriteFile(filepath.Join(templateDir, "go-cli", "template.toml"), []byte(templateConfig), 0644)
				require.NoError(t, err)

				return tempDir
			},
			expectError: false,
			contains:    "go-cli",
		},
		{
			name: "validate valid template",
			args: []string{"template", "validate", "go-cli"},
			setup: func(t *testing.T) string {
				// Create a temporary directory
				tempDir := t.TempDir()

				// Create a template directory
				templateDir := filepath.Join(tempDir, "templates")
				err := os.MkdirAll(templateDir, 0755)
				require.NoError(t, err)

				// Create a template
				err = os.MkdirAll(filepath.Join(templateDir, "go-cli"), 0755)
				require.NoError(t, err)

				// Create template.toml
				templateConfig := `version = "1.0"

[vars]
  name = { prompt = "Enter name:", default = "test" }
  description = { prompt = "Enter description:", default = "A test project" }`

				err = os.WriteFile(filepath.Join(templateDir, "go-cli", "template.toml"), []byte(templateConfig), 0644)
				require.NoError(t, err)

				return tempDir
			},
			expectError: false,
			contains:    "Template is valid",
		},
		{
			name: "validate invalid template",
			args: []string{"template", "validate", "go-cli"},
			setup: func(t *testing.T) string {
				// Create a temporary directory
				tempDir := t.TempDir()

				// Create a template directory
				templateDir := filepath.Join(tempDir, "templates")
				err := os.MkdirAll(templateDir, 0755)
				require.NoError(t, err)

				// Create a template
				err = os.MkdirAll(filepath.Join(templateDir, "go-cli"), 0755)
				require.NoError(t, err)

				// Create invalid template.toml
				templateConfig := `version = "1.0"
[invalid`

				err = os.WriteFile(filepath.Join(templateDir, "go-cli", "template.toml"), []byte(templateConfig), 0644)
				require.NoError(t, err)

				return tempDir
			},
			expectError: true,
			contains:    "failed to parse template config",
		},
		{
			name: "validate nonexistent directory",
			args: []string{"template", "validate", "nonexistent"},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: true,
			contains:    "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save current directory
			currentDir, err := os.Getwd()
			require.NoError(t, err)

			// Create a deferred function to change back to the original directory
			defer func() {
				if err := os.Chdir(currentDir); err != nil {
					t.Errorf("failed to change back to original directory: %v", err)
				}
			}()

			// Change to the test directory
			testDir := tt.setup(t)
			err = os.Chdir(testDir)
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
				if tt.contains != "" {
					assert.Contains(t, err.Error(), tt.contains)
				}
				return
			}

			assert.NoError(t, err)
			if tt.contains != "" {
				assert.Contains(t, buf.String(), tt.contains)
			}
		})
	}
}

func TestTemplateValidateWithInvalidConfig(t *testing.T) {
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

	// Change to the test directory
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Create a template directory
	templateDir := filepath.Join(tempDir, "templates", "go-cli")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)

	// Create invalid template.toml
	templateConfig := `version = "1.0"
[invalid`

	err = os.WriteFile(filepath.Join(templateDir, "template.toml"), []byte(templateConfig), 0644)
	require.NoError(t, err)

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute command
	rootCmd.SetArgs([]string{"template", "validate", "go-cli"})
	err = rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse template config")
}

func TestTemplateValidateWithInvalidTemplateFile(t *testing.T) {
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

	// Change to the test directory
	err = os.Chdir(tempDir)
	require.NoError(t, err)

	// Create a template directory
	templateDir := filepath.Join(tempDir, "templates", "go-cli")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute command
	rootCmd.SetArgs([]string{"template", "validate", "go-cli"})
	err = rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
} 