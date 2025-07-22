package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTemplateConfig(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a template.toml file
	templateConfig := `version = "1.0"

[vars]
  name = { prompt = "Enter name:", default = "test" }
  description = { prompt = "Enter description:", default = "A test project" }

[hooks]
  pre = ["echo 'pre-hook'"]
  post = ["echo 'post-hook'"]`

	err := os.WriteFile(filepath.Join(tempDir, "template.toml"), []byte(templateConfig), 0644)
	require.NoError(t, err)

	// Parse the config
	config, err := ParseTemplateConfig(filepath.Join(tempDir, "template.toml"))
	require.NoError(t, err)

	// Verify the config
	assert.Equal(t, "1.0", config.Version)
	assert.Len(t, config.Vars, 2)
	assert.Equal(t, "Enter name:", config.Vars["name"].Prompt)
	assert.Equal(t, "test", config.Vars["name"].Default)
	assert.Equal(t, "Enter description:", config.Vars["description"].Prompt)
	assert.Equal(t, "A test project", config.Vars["description"].Default)
	assert.Equal(t, []string{"echo 'pre-hook'"}, config.Hooks.Pre)
	assert.Equal(t, []string{"echo 'post-hook'"}, config.Hooks.Post)
}

func TestParseProjectConfig(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a genesis.toml file
	projectConfig := `version = "1.0"

[project]
  template_url = "https://github.com/example/template"
  template_version = "v1.0.0"

[tasks]
  test = { description = "Run tests", cmd = "go test ./..." }
  build = { description = "Build binary", cmd = "go build" }`

	err := os.WriteFile(filepath.Join(tempDir, "genesis.toml"), []byte(projectConfig), 0644)
	require.NoError(t, err)

	// Parse the config
	config, err := ParseProjectConfig(filepath.Join(tempDir, "genesis.toml"))
	require.NoError(t, err)

	// Verify the config
	assert.Equal(t, "1.0", config.Version)
	assert.Equal(t, "https://github.com/example/template", config.Project.TemplateURL)
	assert.Equal(t, "v1.0.0", config.Project.TemplateVersion)
	assert.Len(t, config.Tasks, 2)
	assert.Equal(t, "Run tests", config.Tasks["test"].Description)
	assert.Equal(t, "go test ./...", config.Tasks["test"].Cmd)
	assert.Equal(t, "Build binary", config.Tasks["build"].Description)
	assert.Equal(t, "go build", config.Tasks["build"].Cmd)
}

func TestFindProjectConfig(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	// Create a genesis.toml file in the root directory
	projectConfig := `version = "1.0"

[project]
  template_url = "https://github.com/example/template"
  template_version = "v1.0.0"

[tasks]
  test = { description = "Run tests", cmd = "go test ./..." }`

	configPath := filepath.Join(tempDir, "genesis.toml")
	err = os.WriteFile(configPath, []byte(projectConfig), 0644)
	require.NoError(t, err)

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err)

	// Change to the subdirectory
	err = os.Chdir(subDir)
	require.NoError(t, err)

	// Create a deferred function to change back to the original directory
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Errorf("failed to change back to original directory: %v", err)
		}
	}()

	// Find the config from the subdirectory
	foundPath, err := FindProjectConfig()
	require.NoError(t, err)

	// On macOS, /var/folders is a symlink to /private/var/folders
	// We need to evaluate the symlinks to get the real paths
	realConfigPath, err := filepath.EvalSymlinks(configPath)
	require.NoError(t, err)
	realFoundPath, err := filepath.EvalSymlinks(foundPath)
	require.NoError(t, err)

	// Verify the config was found
	assert.Equal(t, realConfigPath, realFoundPath)
}

func TestParseConfigErrors(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		testFunc func(t *testing.T, content string)
	}{
		{
			name: "invalid TOML syntax",
			content: `version = "1.0"
[invalid
`,
			testFunc: func(t *testing.T, content string) {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "template.toml")
				err := os.WriteFile(configPath, []byte(content), 0644)
				require.NoError(t, err)

				_, err = ParseTemplateConfig(configPath)
				assert.Error(t, err)
			},
		},
		{
			name:    "missing version",
			content: `[vars]`,
			testFunc: func(t *testing.T, content string) {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "template.toml")
				err := os.WriteFile(configPath, []byte(content), 0644)
				require.NoError(t, err)

				_, err = ParseTemplateConfig(configPath)
				assert.Error(t, err)
			},
		},
		{
			name: "invalid task format",
			content: `version = "1.0"
[tasks]
  test = "invalid"`,
			testFunc: func(t *testing.T, content string) {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "genesis.toml")
				err := os.WriteFile(configPath, []byte(content), 0644)
				require.NoError(t, err)

				_, err = ParseProjectConfig(configPath)
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t, tt.content)
		})
	}
} 