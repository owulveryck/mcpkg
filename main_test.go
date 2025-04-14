package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMain tests the main function execution
func TestMain(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run main function
	main()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Verify output
	assert := assert.New(t)
	assert.Contains(output, "Nodes:", "Output should include 'Nodes:' section")
	assert.Contains(output, "Node ID:", "Output should include node information")
	assert.Contains(output, "Edge from", "Output should include edge information")
	assert.Contains(output, "Edge between", "Output should include edge connectivity information")
	assert.Contains(output, "exists: true", "Output should show edges exist")
}