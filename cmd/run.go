package cmd

import (
	"fmt"

	"github.com/felipevolpatto/genesis/internal/config"
	"github.com/felipevolpatto/genesis/internal/runner"
	"github.com/spf13/cobra"
)

func init() {
	runCmd := &cobra.Command{
		Use:   "run [task-name]",
		Short: "Run a task defined in genesis.toml",
		Long: `Run a task defined in genesis.toml.

Tasks are defined in the project's genesis.toml file under the [tasks] section.
Use 'run list' to see all available tasks.`,
		Args: cobra.MinimumNArgs(1),
		RunE: runTask,
	}

	rootCmd.AddCommand(runCmd)
}

func runTask(cmd *cobra.Command, args []string) error {
	taskName := args[0]

	// Handle list command
	if taskName == "list" {
		return listTasks(cmd)
	}

	// Find and parse genesis.toml
	configPath, err := config.FindProjectConfig()
	if err != nil {
		return fmt.Errorf("failed to find genesis.toml: %w", err)
	}

	projectConfig, err := config.ParseProjectConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse genesis.toml: %w", err)
	}

	// Find the task
	task, ok := projectConfig.Tasks[taskName]
	if !ok {
		return fmt.Errorf("task %q not found", taskName)
	}

	// Run the task
	r := runner.New()
	if err := r.RunTask(task); err != nil {
		return fmt.Errorf("failed to run task %q: %w", taskName, err)
	}

	return nil
}

func listTasks(cmd *cobra.Command) error {
	// Find and parse genesis.toml
	configPath, err := config.FindProjectConfig()
	if err != nil {
		return fmt.Errorf("failed to find genesis.toml: %w", err)
	}

	projectConfig, err := config.ParseProjectConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to parse genesis.toml: %w", err)
	}

	// List tasks
	fmt.Fprintln(cmd.OutOrStdout(), "Available tasks:")
	for name, task := range projectConfig.Tasks {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s: %s\n", name, task.Description)
	}

	return nil
} 