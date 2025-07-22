package testutil

import (
	"bytes"
	"io"
	"os"
)

// CaptureOutput captures stdout and stderr output from a function
func CaptureOutput(f func()) string {
	// Save original stdout
	oldStdout := os.Stdout
	oldStderr := os.Stderr

	// Create a pipe
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	// Set stdout and stderr to pipe writer
	os.Stdout = w
	os.Stderr = w

	// Run the function
	f()

	// Close the writer
	w.Close()

	// Restore original stdout and stderr
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Read the output
	var buf bytes.Buffer
	io.Copy(&buf, r)

	return buf.String()
} 