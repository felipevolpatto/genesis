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

func setupTestTemplateDir(t *testing.T) string {
	// Create a temporary directory
	templateDir := t.TempDir()

	// Create a valid template.toml
	validConfig := `version = "1.0"

[vars]
  name = { prompt = "Enter name:", default = "test" }
  description = { prompt = "Enter description:", default = "A test project" }

[hooks]
  post = ["echo 'test'"]`

	err := os.WriteFile(filepath.Join(templateDir, "template.toml"), []byte(validConfig), 0644)
	require.NoError(t, err)

	// Create a valid template file
	validTemplate := `package main

func main() {
	println("Hello, {{ .name }}!")
}`

	err = os.WriteFile(filepath.Join(templateDir, "main.go.tmpl"), []byte(validTemplate), 0644)
	require.NoError(t, err)

	return templateDir
}

func TestTemplateCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		setup       func(t *testing.T) string
		contains    []string
	}{
		{
			name:        "list templates",
			args:        []string{"template", "list"},
			expectError: false,
			setup:       func(t *testing.T) string { return "" },
			contains:    []string{"Available templates:", "go-cli", "node-express"},
		},
		{
			name:        "validate valid template",
			args:        []string{"template", "validate"},
			expectError: false,
			setup:       setupTestTemplateDir,
			contains:    []string{"is valid", "Variables defined:"},
		},
		{
			name: "validate invalid template",
			args: []string{"template", "validate"},
			expectError: true,
			setup: func(t *testing.T) string {
				dir := t.TempDir()
				invalidTemplate := `{{ .invalid }`
				err := os.WriteFile(filepath.Join(dir, "invalid.tmpl"), []byte(invalidTemplate), 0644)
				require.NoError(t, err)
				return dir
			},
			contains: nil,
		},
		{
			name:        "validate nonexistent directory",
			args:        []string{"template", "validate", "/nonexistent/path"},
			expectError: true,
			setup:       func(t *testing.T) string { return "" },
			contains:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up test directory if needed
			dir := tt.setup(t)
			if dir != "" {
				currentDir, err := os.Getwd()
				require.NoError(t, err)
				defer os.Chdir(currentDir)

				err = os.Chdir(dir)
				require.NoError(t, err)
			}

			// Create a buffer to capture output
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			// Execute command and capture output
			var output string
			if tt.contains != nil {
				output = testutil.CaptureOutput(func() {
					rootCmd.SetArgs(tt.args)
					err := rootCmd.Execute()
					if tt.expectError {
						assert.Error(t, err)
						return
					}
					assert.NoError(t, err)
				})
			} else {
				err := rootCmd.Execute()
				if tt.expectError {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
			}

			// Check output
			if tt.contains != nil {
				for _, s := range tt.contains {
					assert.Contains(t, output, s)
				}
			}
		})
	}
}

func TestTemplateValidateWithInvalidConfig(t *testing.T) {
	// Create a temporary directory
	dir := t.TempDir()

	// Create an invalid template.toml
	invalidConfig := `version = "1.0"
[invalid`

	err := os.WriteFile(filepath.Join(dir, "template.toml"), []byte(invalidConfig), 0644)
	require.NoError(t, err)

	// Change to the test directory
	currentDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(currentDir)

	err = os.Chdir(dir)
	require.NoError(t, err)

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute command
	rootCmd.SetArgs([]string{"template", "validate"})
	err = rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid")
}

func TestTemplateValidateWithInvalidTemplateFile(t *testing.T) {
	// Create a temporary directory
	dir := t.TempDir()

	// Create a valid template.toml
	validConfig := `version = "1.0"
[vars]
  name = { prompt = "Enter name:", default = "test" }`

	err := os.WriteFile(filepath.Join(dir, "template.toml"), []byte(validConfig), 0644)
	require.NoError(t, err)

	// Create an invalid template file
	invalidTemplate := `{{ .name ` // Missing closing brace

	err = os.WriteFile(filepath.Join(dir, "invalid.tmpl"), []byte(invalidTemplate), 0644)
	require.NoError(t, err)

	// Change to the test directory
	currentDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(currentDir)

	err = os.Chdir(dir)
	require.NoError(t, err)

	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute command
	rootCmd.SetArgs([]string{"template", "validate"})
	err = rootCmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid template syntax")
} 