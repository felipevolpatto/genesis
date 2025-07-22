package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/felipevolpatto/genesis/internal/config"
	"github.com/felipevolpatto/genesis/internal/runner"
	"github.com/felipevolpatto/genesis/internal/scaffolder"
	"github.com/felipevolpatto/genesis/internal/tui"
	"github.com/spf13/cobra"
)

var (
	templateURL string
	skipPrompts bool
	version     string
)

func init() {
	newCmd := &cobra.Command{
		Use:   "new [project-name]",
		Short: "Create a new project from a template",
		Long: `Create a new project from a template.

The template should be a Git repository containing a template.toml file and the project structure.
Template files ending in .tmpl will be processed using Go's text/template package.`,
		Args: cobra.ExactArgs(1),
		RunE: runNew,
	}

	newCmd.Flags().StringVarP(&templateURL, "template", "t", "", "Template URL or path (required)")
	newCmd.Flags().BoolVarP(&skipPrompts, "yes", "y", false, "Skip all prompts and use default values")
	newCmd.Flags().StringVarP(&version, "version", "v", "", "Template version (tag, branch, or commit hash)")

	newCmd.MarkFlagRequired("template")
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// Validate template URL
	if templateURL == "" {
		return fmt.Errorf("template URL is required")
	}

	// Clone template repository
	templateDir, err := scaffolder.CloneTemplate(templateURL, version)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	defer scaffolder.CleanupTemplate(templateDir)

	// Parse template config
	templateConfig, err := config.ParseTemplateConfig(filepath.Join(templateDir, "template.toml"))
	if err != nil {
		return fmt.Errorf("failed to parse template config: %w", err)
	}

	// Get variable values
	variables := make(map[string]string)
	if !skipPrompts {
		variables, err = tui.PromptForVariables(templateConfig.Vars)
		if err != nil {
			return fmt.Errorf("failed to get variable values: %w", err)
		}
	} else {
		// Use default values
		for name, v := range templateConfig.Vars {
			variables[name] = v.Default
		}
	}

	// Create project directory
	projectDir := filepath.Join(".", projectName)
	s := scaffolder.New(templateDir, projectDir, variables, templateConfig)

	// Run pre-hooks
	if len(templateConfig.Hooks.Pre) > 0 {
		fmt.Println("Running pre-hooks...")
		if err := runner.RunHooks(templateConfig.Hooks.Pre, projectDir); err != nil {
			return fmt.Errorf("failed to run pre-hooks: %w", err)
		}
	}

	// Scaffold project
	if err := s.Scaffold(); err != nil {
		return fmt.Errorf("failed to scaffold project: %w", err)
	}

	// Create genesis.toml
	if err := s.CreateGenesisConfig(templateURL, version); err != nil {
		return fmt.Errorf("failed to create genesis.toml: %w", err)
	}

	// Run post-hooks
	if len(templateConfig.Hooks.Post) > 0 {
		fmt.Println("Running post-hooks...")
		if err := runner.RunHooks(templateConfig.Hooks.Post, projectDir); err != nil {
			return fmt.Errorf("failed to run post-hooks: %w", err)
		}
	}

	fmt.Printf("\nProject %q created successfully!\n", projectName)
	return nil
} 