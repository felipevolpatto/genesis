package scaffolder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/felipevolpatto/genesis/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScaffolder(t *testing.T) {
	// Create a temporary template directory
	templateDir := t.TempDir()
	targetDir := t.TempDir()

	// Create template files
	files := map[string]string{
		"main.go.tmpl": `package main

func main() {
	println("Hello, {{ .name }}!")
}`,
		"README.md": "# Static File",
		"config.json.tmpl": `{
	"name": "{{ .name }}",
	"description": "{{ .description }}"
}`,
	}

	for name, content := range files {
		err := os.WriteFile(filepath.Join(templateDir, name), []byte(content), 0644)
		require.NoError(t, err)
	}

	// Create template config
	templateConfig := &config.TemplateConfig{
		Version: "1.0",
		Vars: map[string]config.Variable{
			"name":        {Prompt: "Enter name:", Default: "World"},
			"description": {Prompt: "Enter description:", Default: "A test project"},
		},
	}

	// Create variables
	variables := map[string]string{
		"name":        "Test",
		"description": "A test description",
	}

	// Create scaffolder
	s := New(templateDir, targetDir, variables, templateConfig)

	// Test scaffolding
	err := s.Scaffold()
	require.NoError(t, err)

	// Verify files were created correctly
	files = map[string]string{
		"main.go": `package main

func main() {
	println("Hello, Test!")
}`,
		"README.md": "# Static File",
		"config.json": `{
	"name": "Test",
	"description": "A test description"
}`,
	}

	for name, expectedContent := range files {
		content, err := os.ReadFile(filepath.Join(targetDir, name))
		require.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	}
}

func TestCreateGenesisConfig(t *testing.T) {
	targetDir := t.TempDir()
	templateConfig := &config.TemplateConfig{Version: "1.0"}
	s := New("", targetDir, nil, templateConfig)

	templateURL := "https://github.com/example/template"
	templateVersion := "v1.0.0"

	err := s.CreateGenesisConfig(templateURL, templateVersion)
	require.NoError(t, err)

	// Verify genesis.toml was created correctly
	content, err := os.ReadFile(filepath.Join(targetDir, "genesis.toml"))
	require.NoError(t, err)

	expectedContent := `# The version of the genesis config spec
version = "1.0"

[project]
  template_url = "https://github.com/example/template"
  template_version = "v1.0.0"

# [tasks] defines the commands that can be run with 'genesis run <task-name>'
[tasks]
`

	assert.Equal(t, expectedContent, string(content))
}

func TestScaffolderErrors(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) (*Scaffolder, error)
		expectError bool
	}{
		{
			name: "invalid template syntax",
			setupFunc: func(t *testing.T) (*Scaffolder, error) {
				templateDir := t.TempDir()
				targetDir := t.TempDir()

				// Create invalid template
				err := os.WriteFile(filepath.Join(templateDir, "invalid.tmpl"), []byte("{{ .invalid"), 0644)
				if err != nil {
					return nil, err
				}

				templateConfig := &config.TemplateConfig{Version: "1.0"}
				return New(templateDir, targetDir, nil, templateConfig), nil
			},
			expectError: true,
		},
		{
			name: "target directory creation error",
			setupFunc: func(t *testing.T) (*Scaffolder, error) {
				// Create a file where the target directory should be
				targetDir := filepath.Join(t.TempDir(), "target")
				err := os.WriteFile(targetDir, []byte(""), 0644)
				if err != nil {
					return nil, err
				}

				templateConfig := &config.TemplateConfig{Version: "1.0"}
				return New("", targetDir, nil, templateConfig), nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.setupFunc(t)
			require.NoError(t, err)

			err = s.Scaffold()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
} 