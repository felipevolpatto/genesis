package scaffolder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/felipevolpatto/genesis/internal/config"
)

// Scaffolder handles the project scaffolding process
type Scaffolder struct {
	templateDir string
	targetDir   string
	variables   map[string]string
	config      *config.TemplateConfig
}

// New creates a new Scaffolder instance
func New(templateDir, targetDir string, variables map[string]string, config *config.TemplateConfig) *Scaffolder {
	return &Scaffolder{
		templateDir: templateDir,
		targetDir:   targetDir,
		variables:   variables,
		config:      config,
	}
}

// Scaffold processes the template and creates the new project
func (s *Scaffolder) Scaffold() error {
	// Create target directory if it doesn't exist
	if err := os.MkdirAll(s.targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Walk through the template directory
	return filepath.Walk(s.templateDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Skip template.toml
		if info.Name() == "template.toml" {
			return nil
		}

		// Calculate relative path from template root
		relPath, err := filepath.Rel(s.templateDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Create target path
		targetPath := filepath.Join(s.targetDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		// Process or copy the file
		if strings.HasSuffix(path, ".tmpl") {
			return s.processTemplate(path, strings.TrimSuffix(targetPath, ".tmpl"))
		}

		return s.copyFile(path, targetPath)
	})
}

// processTemplate processes a template file and writes the result
func (s *Scaffolder) processTemplate(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	tmpl, err := template.New(filepath.Base(src)).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	if err := tmpl.Execute(out, s.variables); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func (s *Scaffolder) copyFile(src, dst string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	return os.WriteFile(dst, content, 0644)
}

// CreateGenesisConfig creates a genesis.toml file in the target directory
func (s *Scaffolder) CreateGenesisConfig(templateURL, templateVersion string) error {
	config := fmt.Sprintf(`# The version of the genesis config spec
version = "1.0"

[project]
  template_url = %q
  template_version = %q

# [tasks] defines the commands that can be run with 'genesis run <task-name>'
[tasks]
`, templateURL, templateVersion)

	return os.WriteFile(filepath.Join(s.targetDir, "genesis.toml"), []byte(config), 0644)
} 