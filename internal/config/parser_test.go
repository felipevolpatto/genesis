package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseTemplateConfig(t *testing.T) {
	// Create a temporary template.toml file
	content := `version = "1.0"

[vars]
  description = { prompt = "Enter description:", default = "A test project" }
  author = { prompt = "Enter author:", default = "Test Author" }

[hooks]
  pre = ["echo 'pre-hook'"]
  post = ["echo 'post-hook'"]`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "template.toml")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	// Test parsing
	config, err := ParseTemplateConfig(tmpFile)
	require.NoError(t, err)
	assert.Equal(t, "1.0", config.Version)

	// Test variables
	assert.Len(t, config.Vars, 2)
	assert.Equal(t, "Enter description:", config.Vars["description"].Prompt)
	assert.Equal(t, "A test project", config.Vars["description"].Default)
	assert.Equal(t, "Enter author:", config.Vars["author"].Prompt)
	assert.Equal(t, "Test Author", config.Vars["author"].Default)

	// Test hooks
	assert.Len(t, config.Hooks.Pre, 1)
	assert.Equal(t, "echo 'pre-hook'", config.Hooks.Pre[0])
	assert.Len(t, config.Hooks.Post, 1)
	assert.Equal(t, "echo 'post-hook'", config.Hooks.Post[0])
}

func TestParseProjectConfig(t *testing.T) {
	// Create a temporary genesis.toml file
	content := `version = "1.0"

[project]
  template_url = "https://github.com/example/template"
  template_version = "v1.0.0"

[tasks]
  test = { description = "Run tests", cmd = "go test ./..." }
  build = { description = "Build binary", cmd = "go build" }`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "genesis.toml")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)

	// Test parsing
	config, err := ParseProjectConfig(tmpFile)
	require.NoError(t, err)
	assert.Equal(t, "1.0", config.Version)

	// Test project section
	assert.Equal(t, "https://github.com/example/template", config.Project.TemplateURL)
	assert.Equal(t, "v1.0.0", config.Project.TemplateVersion)

	// Test tasks
	assert.Len(t, config.Tasks, 2)
	assert.Equal(t, "Run tests", config.Tasks["test"].Description)
	assert.Equal(t, "go test ./...", config.Tasks["test"].Cmd)
	assert.Equal(t, "Build binary", config.Tasks["build"].Description)
	assert.Equal(t, "go build", config.Tasks["build"].Cmd)
}

func TestFindProjectConfig(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	content := `version = "1.0"
[project]
  template_url = "https://github.com/example/template"
  template_version = "v1.0.0"
[tasks]
  test = { description = "Run tests", cmd = "go test ./..." }`

	configPath := filepath.Join(tmpDir, "genesis.toml")
	err = os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err)

	// Test finding config from subdirectory
	currentDir, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(currentDir)

	err = os.Chdir(subDir)
	require.NoError(t, err)

	config, dir, err := FindProjectConfig()
	require.NoError(t, err)

	// On macOS, /var/folders is a symlink to /private/var/folders
	// We need to evaluate the symlinks to get the real path
	realTmpDir, err := filepath.EvalSymlinks(tmpDir)
	require.NoError(t, err)
	realDir, err := filepath.EvalSymlinks(dir)
	require.NoError(t, err)

	assert.Equal(t, realTmpDir, realDir)
	assert.Equal(t, "1.0", config.Version)
	assert.Equal(t, "https://github.com/example/template", config.Project.TemplateURL)

	// Test error when no config found
	emptyDir := t.TempDir()
	err = os.Chdir(emptyDir)
	require.NoError(t, err)

	_, _, err = FindProjectConfig()
	assert.Error(t, err)
}

func TestParseConfigErrors(t *testing.T) {
	tests := []struct {
		name    string
		content string
		isValid bool
	}{
		{
			name: "invalid TOML syntax",
			content: `version = 1.0
			[invalid
			`,
			isValid: false,
		},
		{
			name: "missing version",
			content: `[tasks]
			test = { description = "Run tests", cmd = "go test ./..." }`,
			isValid: true, // Version is optional
		},
		{
			name: "invalid task format",
			content: `version = "1.0"
			[tasks]
			test = "invalid"`,
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "genesis.toml")
			err := os.WriteFile(tmpFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			_, err = ParseProjectConfig(tmpFile)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
} 