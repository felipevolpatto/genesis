package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/felipevolpatto/genesis/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	templateCmd := &cobra.Command{
		Use:   "template",
		Short: "Manage project templates",
		Long:  `Manage project templates, including listing available templates and validating template structure.`,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available templates",
		Long:  `List available official and curated templates.`,
		RunE:  listTemplates,
	}

	validateCmd := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate a template",
		Long: `Validate the structure and syntax of a template.
The template should contain a valid template.toml file and properly formatted template files.`,
		Args: cobra.MaximumNArgs(1),
		RunE: validateTemplate,
	}

	templateCmd.AddCommand(listCmd, validateCmd)
	rootCmd.AddCommand(templateCmd)
}

func listTemplates(cmd *cobra.Command, args []string) error {
	// For V1, we'll just show a simple hardcoded list of example templates
	templates := []struct {
		Name        string
		URL         string
		Description string
	}{
		{
			Name:        "go-cli",
			URL:         "https://github.com/genesis/template-go-cli",
			Description: "A template for Go CLI applications using cobra",
		},
		{
			Name:        "node-express",
			URL:         "https://github.com/genesis/template-node-express",
			Description: "A template for Node.js Express applications",
		},
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Available templates:")
	for _, t := range templates {
		fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", t.Name)
		fmt.Fprintf(cmd.OutOrStdout(), "  URL: %s\n", t.URL)
		fmt.Fprintf(cmd.OutOrStdout(), "  Description: %s\n", t.Description)
	}

	return nil
}

func validateTemplate(cmd *cobra.Command, args []string) error {
	// Get template path
	templatePath := "."
	if len(args) > 0 {
		templatePath = args[0]
	}

	// Check if template.toml exists
	configPath := filepath.Join(templatePath, "template.toml")
	if _, err := os.Stat(configPath); err != nil {
		return fmt.Errorf("template.toml not found in %s", templatePath)
	}

	// Parse template.toml
	templateConfig, err := config.ParseTemplateConfig(configPath)
	if err != nil {
		return fmt.Errorf("invalid template.toml: %w", err)
	}

	// Validate template files
	err = filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-.tmpl files
		if info.IsDir() || filepath.Ext(path) != ".tmpl" {
			return nil
		}

		// Try to parse the template file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		_, err = template.New(filepath.Base(path)).Parse(string(content))
		if err != nil {
			return fmt.Errorf("invalid template syntax in %s: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Template in %s is valid\n", templatePath)
	fmt.Fprintf(cmd.OutOrStdout(), "Variables defined:\n")
	for name, v := range templateConfig.Vars {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s\n    Prompt: %s\n    Default: %s\n", name, v.Prompt, v.Default)
		if v.Regex != "" {
			fmt.Fprintf(cmd.OutOrStdout(), "    Validation: %s\n", v.Regex)
		}
	}

	return nil
} 