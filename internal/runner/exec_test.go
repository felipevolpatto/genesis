package runner

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/felipevolpatto/genesis/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunner(t *testing.T) {
	// Create a temporary directory
	dir := t.TempDir()

	// Create a test config
	cfg := &config.ProjectConfig{
		Version: "1.0",
		Tasks: map[string]config.Task{
			"echo": {
				Description: "Echo test",
				Cmd:        "echo 'test'",
			},
			"pwd": {
				Description: "Print working directory",
				Cmd:        "pwd",
			},
		},
	}

	// Create runner
	r := New(cfg, dir)

	// Test running a task
	err := r.RunTask("echo")
	assert.NoError(t, err)

	// Test running a non-existent task
	err = r.RunTask("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRunHooks(t *testing.T) {
	dir := t.TempDir()

	// Create a test file to verify hook execution
	testFile := filepath.Join(dir, "test.txt")
	hooks := []string{
		"echo 'test' > " + testFile,
	}

	// Run hooks
	err := RunHooks(hooks, dir)
	require.NoError(t, err)

	// Verify hook execution
	content, err := os.ReadFile(testFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "test")

	// Test invalid hook
	hooks = []string{"invalid_command"}
	err = RunHooks(hooks, dir)
	assert.Error(t, err)
}

func TestListTasks(t *testing.T) {
	cfg := &config.ProjectConfig{
		Tasks: map[string]config.Task{
			"short": {
				Description: "Short description",
				Cmd:        "echo 'short'",
			},
			"longer_name": {
				Description: "Longer description",
				Cmd:        "echo 'longer'",
			},
		},
	}

	r := New(cfg, "")
	tasks := r.ListTasks()

	// Tasks should be formatted with proper alignment
	assert.Len(t, tasks, 2)
	for _, task := range tasks {
		// Each task should contain both the name and description
		assert.True(t, strings.Contains(task, "short") || strings.Contains(task, "longer_name"))
		assert.True(t, strings.Contains(task, "Short description") || strings.Contains(task, "Longer description"))
	}
}

func TestGetShellAndArg(t *testing.T) {
	shell, arg := getShellAndArg()

	if runtime.GOOS == "windows" {
		assert.Equal(t, "cmd", shell)
		assert.Equal(t, "/C", arg)
	} else {
		assert.Equal(t, "/bin/sh", shell)
		assert.Equal(t, "-c", arg)
	}
}

func TestRunTaskWithEnvironment(t *testing.T) {
	dir := t.TempDir()
	testEnvVar := "TEST_ENV_VAR"
	testEnvValue := "test_value"

	// Set environment variable
	os.Setenv(testEnvVar, testEnvValue)
	defer os.Unsetenv(testEnvVar)

	// Create config that uses the environment variable
	cfg := &config.ProjectConfig{
		Tasks: map[string]config.Task{
			"env": {
				Description: "Test environment variables",
				Cmd:        "echo $" + testEnvVar,
			},
		},
	}

	r := New(cfg, dir)
	err := r.RunTask("env")
	assert.NoError(t, err)
}

func TestRunTaskWithWorkingDirectory(t *testing.T) {
	// Create a temporary directory with a unique file
	dir := t.TempDir()
	uniqueFile := filepath.Join(dir, "unique.txt")
	err := os.WriteFile(uniqueFile, []byte("test"), 0644)
	require.NoError(t, err)

	// Create config that lists directory contents
	cfg := &config.ProjectConfig{
		Tasks: map[string]config.Task{
			"ls": {
				Description: "List directory",
				Cmd:        "ls",
			},
		},
	}

	r := New(cfg, dir)
	err = r.RunTask("ls")
	assert.NoError(t, err)
} 