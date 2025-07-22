package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/felipevolpatto/genesis/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunner(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err)

	// Create a deferred function to change back to the original directory
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Errorf("failed to change back to original directory: %v", err)
		}
	}()

	// Change to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temporary directory: %v", err)
	}

	// Create a test file
	err = os.WriteFile("test.txt", []byte("test"), 0644)
	require.NoError(t, err)

	// Create a runner
	r := New()

	// Run a task
	task := config.Task{
		Cmd: "cat test.txt",
	}
	err = r.RunTask(task)
	assert.NoError(t, err)

	// Run a nonexistent command
	task = config.Task{
		Cmd: "nonexistent",
	}
	err = r.RunTask(task)
	assert.Error(t, err)
}

func TestRunHooks(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err)

	// Create a deferred function to change back to the original directory
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Errorf("failed to change back to original directory: %v", err)
		}
	}()

	// Change to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temporary directory: %v", err)
	}

	// Run hooks
	hooks := []string{
		"echo 'test' > test.txt",
		"invalid_command",
	}

	err = RunHooks(hooks, tempDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid_command")

	// Verify the first hook ran successfully
	_, err = os.Stat("test.txt")
	assert.NoError(t, err)
}

func TestRunTaskWithEnvironment(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err)

	// Create a deferred function to change back to the original directory
	defer func() {
		if err := os.Chdir(currentDir); err != nil {
			t.Errorf("failed to change back to original directory: %v", err)
		}
	}()

	// Change to the temporary directory
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temporary directory: %v", err)
	}

	// Create a runner
	r := New()

	// Run a task with environment variables
	task := config.Task{
		Cmd: "echo $TEST_VALUE",
		Env: map[string]string{
			"TEST_VALUE": "test_value",
		},
	}
	err = r.RunTask(task)
	assert.NoError(t, err)
}

func TestRunTaskWithWorkingDirectory(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	// Create a unique file in the subdirectory
	err = os.WriteFile(filepath.Join(subDir, "unique.txt"), []byte("test"), 0644)
	require.NoError(t, err)

	// Create a runner
	r := New()

	// Run a task in the subdirectory
	task := config.Task{
		Cmd: "ls",
		Dir: subDir,
	}
	err = r.RunTask(task)
	assert.NoError(t, err)
} 