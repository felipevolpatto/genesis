package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "genesis",
	Short: "Genesis - Begin any project, unified.",
	Long: `Genesis is a powerful project scaffolding and task runner tool that helps you
start and manage projects across different tech stacks with a unified interface.

It provides:
  - Project scaffolding from templates
  - Task running and automation
  - Template management and validation`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
} 