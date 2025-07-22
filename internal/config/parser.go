package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Variable represents a template variable with its prompt and default value
type Variable struct {
	Prompt  string
	Default string
	Regex   string
}

// Hooks represents pre and post scaffolding hooks
type Hooks struct {
	Pre  []string
	Post []string
}

// TemplateConfig represents the configuration for a template
type TemplateConfig struct {
	Version string
	Vars    map[string]Variable
	Hooks   Hooks
}

// Project represents project-specific configuration
type Project struct {
	TemplateURL     string `toml:"template_url"`
	TemplateVersion string `toml:"template_version"`
}

// Task represents a runnable task
type Task struct {
	Description string
	Cmd        string
	Env        map[string]string
	Dir        string
}

// ProjectConfig represents the configuration for a project
type ProjectConfig struct {
	Version string
	Project Project
	Tasks   map[string]Task
}

// ParseTemplateConfig parses a template.toml file
func ParseTemplateConfig(path string) (*TemplateConfig, error) {
	var config TemplateConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}

	if config.Version == "" {
		return nil, fmt.Errorf("template config must specify a version")
	}

	return &config, nil
}

// ParseProjectConfig parses a genesis.toml file
func ParseProjectConfig(path string) (*ProjectConfig, error) {
	var config ProjectConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to parse project config: %w", err)
	}

	if config.Version == "" {
		return nil, fmt.Errorf("project config must specify a version")
	}

	return &config, nil
}

// FindProjectConfig searches for a genesis.toml file in the current directory
// and its parents. Returns the path to the config file if found.
func FindProjectConfig() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		configPath := filepath.Join(dir, "genesis.toml")
		if _, err := os.Stat(configPath); err == nil {
			return configPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("no genesis.toml found in current directory or its parents")
} 