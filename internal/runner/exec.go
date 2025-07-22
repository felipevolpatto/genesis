package runner

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/felipevolpatto/genesis/internal/config"
)

// Runner executes tasks defined in genesis.toml
type Runner struct {
	shell string
	arg   string
}

// New creates a new Runner
func New() *Runner {
	shell, arg := getShellAndArg()
	return &Runner{
		shell: shell,
		arg:   arg,
	}
}

// RunTask executes a task
func (r *Runner) RunTask(task config.Task) error {
	// Set up command environment
	env := os.Environ()
	for k, v := range task.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	// Set up working directory
	dir := task.Dir
	if dir == "" {
		dir = "."
	}

	// Create command
	cmd := exec.Command(r.shell, r.arg, task.Cmd)
	cmd.Env = env
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// RunHooks executes a list of commands
func RunHooks(hooks []string, dir string) error {
	r := New()
	for _, cmd := range hooks {
		task := config.Task{
			Cmd: cmd,
			Dir: dir,
		}
		if err := r.RunTask(task); err != nil {
			return fmt.Errorf("failed to run hook %q: %w", cmd, err)
		}
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