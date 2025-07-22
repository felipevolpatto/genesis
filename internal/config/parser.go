package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Variable represents a template variable with its prompt and default value
type Variable struct {
	Prompt   string `toml:"prompt"`
	Default  string `toml:"default"`
	Regex    string `toml:"regex,omitempty"`
}

// TemplateConfig represents the structure of a template.toml file
type TemplateConfig struct {
	Version string              `toml:"version"`
	Vars    map[string]Variable `toml:"vars"`
	Hooks   struct {
		Pre  []string `toml:"pre,omitempty"`
		Post []string `toml:"post,omitempty"`
	} `toml:"hooks"`
}

// Task represents a runnable task in genesis.toml
type Task struct {
	Description string `toml:"description"`
	Cmd        string `toml:"cmd"`
}

// ProjectConfig represents the structure of a genesis.toml file
type ProjectConfig struct {
	Version string `toml:"version"`
	Project struct {
		TemplateURL     string `toml:"template_url"`
		TemplateVersion string `toml:"template_version"`
	} `toml:"project"`
	Tasks map[string]Task `toml:"tasks"`
}

// ParseTemplateConfig reads and parses a template.toml file
func ParseTemplateConfig(path string) (*TemplateConfig, error) {
	var config TemplateConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to parse template config: %w", err)
	}
	return &config, nil
}

// ParseProjectConfig reads and parses a genesis.toml file
func ParseProjectConfig(path string) (*ProjectConfig, error) {
	var config ProjectConfig
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("failed to parse project config: %w", err)
	}
	return &config, nil
}

// FindProjectConfig searches for genesis.toml in the current directory and parent directories
func FindProjectConfig() (*ProjectConfig, string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get working directory: %w", err)
	}

	for {
		configPath := filepath.Join(dir, "genesis.toml")
		if _, err := os.Stat(configPath); err == nil {
			config, err := ParseProjectConfig(configPath)
			return config, dir, err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return nil, "", fmt.Errorf("no genesis.toml found in current directory or any parent directories")
} 