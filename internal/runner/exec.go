package runner

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/felipevolpatto/genesis/internal/config"
)

// Runner handles task execution
type Runner struct {
	config *config.ProjectConfig
	dir    string
}

// New creates a new Runner instance
func New(config *config.ProjectConfig, dir string) *Runner {
	return &Runner{
		config: config,
		dir:    dir,
	}
}

// RunTask executes a task by name
func (r *Runner) RunTask(taskName string) error {
	task, ok := r.config.Tasks[taskName]
	if !ok {
		return fmt.Errorf("task %q not found", taskName)
	}

	// Get the shell command and arguments based on the OS
	shell, shellArg := getShellAndArg()

	// Create the command
	cmd := exec.Command(shell, shellArg, task.Cmd)
	cmd.Dir = r.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Set environment variables
	cmd.Env = os.Environ()

	// Run the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("task %q failed: %w", taskName, err)
	}

	return nil
}

// getShellAndArg returns the appropriate shell and argument for the current OS
func getShellAndArg() (string, string) {
	if runtime.GOOS == "windows" {
		return "cmd", "/C"
	}
	return "/bin/sh", "-c"
}

// RunHooks executes a list of hooks
func RunHooks(hooks []string, dir string) error {
	for _, hook := range hooks {
		shell, shellArg := getShellAndArg()
		cmd := exec.Command(shell, shellArg, hook)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook %q failed: %w", hook, err)
		}
	}
	return nil
}

// ListTasks returns a list of available tasks with their descriptions
func (r *Runner) ListTasks() []string {
	var tasks []string
	maxLen := 0

	// Find the longest task name for alignment
	for name := range r.config.Tasks {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	// Format each task with proper alignment
	for name, task := range r.config.Tasks {
		padding := strings.Repeat(" ", maxLen-len(name))
		tasks = append(tasks, fmt.Sprintf("%s%s  %s", name, padding, task.Description))
	}

	return tasks
} 