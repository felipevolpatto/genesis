package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Test help output
	rootCmd.SetArgs([]string{"--help"})
	err := rootCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Genesis is a powerful project scaffolding")
	assert.Contains(t, output, "Available Commands:")
	assert.Contains(t, output, "new")
	assert.Contains(t, output, "run")
	assert.Contains(t, output, "template")
}

func TestRootCommandNoArgs(t *testing.T) {
	// Create a buffer to capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Test with no arguments
	rootCmd.SetArgs([]string{})
	err := rootCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Available Commands:")
} 