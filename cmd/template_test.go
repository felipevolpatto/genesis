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
		contains    []string
	}{
		{
			name: "list templates",
			args: []string{"template", "list"},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: false,
			contains: []string{
				"Available templates:",
				"go-cli",
				"URL: https://github.com/genesis/template-go-cli",
				"Description: A template for Go CLI applications using cobra",
			},
		},
		{
			name: "validate valid template",
			args: []string{"template", "validate", "."},
			setup: func(t *testing.T) string {
				// Create a temporary directory
				tempDir := t.TempDir()

				// Create template.toml
				templateConfig := `version = "1.0"

[vars]
  name = { prompt = "Enter name:", default = "test" }
  description = { prompt = "Enter description:", default = "A test project" }`

				err := os.WriteFile(filepath.Join(tempDir, "template.toml"), []byte(templateConfig), 0644)
				require.NoError(t, err)

				return tempDir
			},
			expectError: false,
			contains: []string{
				"Template in . is valid",
				"Variables defined:",
				"name",
				"Prompt: Enter name:",
				"Default: test",
			},
		},
		{
			name: "validate invalid template",
			args: []string{"template", "validate", "."},
			setup: func(t *testing.T) string {
				// Create a temporary directory
				tempDir := t.TempDir()

				// Create invalid template.toml
				templateConfig := `version = "1.0"
[invalid`

				err := os.WriteFile(filepath.Join(tempDir, "template.toml"), []byte(templateConfig), 0644)
				require.NoError(t, err)

				return tempDir
			},
			expectError: true,
			contains: []string{
				"invalid template.toml",
			},
		},
		{
			name: "validate nonexistent directory",
			args: []string{"template", "validate", "nonexistent"},
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			expectError: true,
			contains: []string{
				"template.toml not found in nonexistent",
			},
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
				for _, s := range tt.contains {
					assert.Contains(t, err.Error(), s)
				}
				return
			}

			assert.NoError(t, err)
			output := buf.String()
			for _, s := range tt.contains {
				assert.Contains(t, output, s)
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

	// Create invalid template.toml
	templateConfig := `version = "1.0"
[invalid`

	err = os.WriteFile("template.toml", []byte(templateConfig), 0644)
	require.NoError(t, err)

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute command
	rootCmd.SetArgs([]string{"template", "validate", "."})
	err = rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid template.toml")
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

	// Create template.toml
	templateConfig := `version = "1.0"
[vars]
  name = { prompt = "Enter name:", default = "test" }`

	err = os.WriteFile("template.toml", []byte(templateConfig), 0644)
	require.NoError(t, err)

	// Create invalid template file
	invalidTemplate := `{{ .name ` // Missing closing brace
	err = os.WriteFile("invalid.tmpl", []byte(invalidTemplate), 0644)
	require.NoError(t, err)

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute command
	rootCmd.SetArgs([]string{"template", "validate", "."})
	err = rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid template syntax")
} 