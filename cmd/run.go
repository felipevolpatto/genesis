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
		Long: `Run a task defined in the genesis.toml configuration file.

Tasks are shell commands that can be executed in the context of your project.
The genesis.toml file should be in the current directory or any parent directory.

Example:
  genesis run test     # Run the test task
  genesis run build    # Run the build task`,
		Args: cobra.ExactArgs(1),
		RunE: runTask,
	}

	rootCmd.AddCommand(runCmd)
}

func runTask(cmd *cobra.Command, args []string) error {
	taskName := args[0]

	// Find and parse genesis.toml
	projectConfig, projectDir, err := config.FindProjectConfig()
	if err != nil {
		return err
	}

	// Create runner
	r := runner.New(projectConfig, projectDir)

	// List tasks if requested
	if taskName == "list" {
		tasks := r.ListTasks()
		if len(tasks) == 0 {
			fmt.Println("No tasks defined in genesis.toml")
			return nil
		}

		fmt.Println("Available tasks:")
		for _, task := range tasks {
			fmt.Printf("  %s\n", task)
		}
		return nil
	}

	// Run the requested task
	return r.RunTask(taskName)
} 