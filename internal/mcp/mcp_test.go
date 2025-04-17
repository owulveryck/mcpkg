package mcp_test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func compileServer(outputPath string) error {
	cmd := exec.Command(
		"go",
		"build",
		"-o",
		outputPath,
		"./testdata/",
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %v\nOutput: %s", err, output)
	}
	return nil
}

func TestServer(t *testing.T) {
	mockServerPath := "testdata/server"
	err := compileServer(mockServerPath)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer os.Remove(mockServerPath)

	client, err := client.NewStdioMCPClient(mockServerPath, []string{})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	ctx := context.Background()
	request := mcp.InitializeRequest{}
	request.Params.ProtocolVersion = "1.0"
	request.Params.ClientInfo = mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}
	request.Params.Capabilities = mcp.ClientCapabilities{
		Roots: &struct {
			ListChanged bool `json:"listChanged,omitempty"`
		}{
			ListChanged: true,
		},
	}

	_, err = client.Initialize(ctx, request)
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "mcp-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after the test.
	kgfile := filepath.Join(tempDir, "testkg.kg")

	subject := "the sky"
	object := "blue"
	callToolRequest := mcp.CallToolRequest{}
	callToolRequest.Params.Name = "insert_triple"
	callToolRequest.Params.Arguments = map[string]interface{}{
		"knowledge_graph_path": kgfile,
		"subject":              subject,
		"predicate":            "is",
		"object":               object,
	}

	callToolResult, err := client.CallTool(ctx, callToolRequest)
	if err != nil {
		t.Errorf("CallTool failed: %v", err)
	}

	if len(callToolResult.Content) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(callToolResult.Content))
	}

	readRequest := mcp.ReadResourceRequest{}
	template := "graph://%s?from=%s&to=%s" // Pre-encoded template

	// URI-encode the parameters
	encodedPath := url.QueryEscape(kgfile)
	encodedFrom := url.QueryEscape(subject)
	encodedTo := url.QueryEscape(object)
	// Replace plus signs with %20
	encodedFrom = strings.ReplaceAll(encodedFrom, "+", "%20")
	encodedTo = strings.ReplaceAll(encodedTo, "+", "%20")

	// Substitute the values into the template
	uri := fmt.Sprintf(template, encodedPath, encodedFrom, encodedTo)
	readRequest.Params.URI = uri

	resultRequest, err := client.ReadResource(ctx, readRequest)
	if err != nil {
		t.Fatalf("ReadResource failed: %v", err)
	}

	if len(resultRequest.Contents) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(resultRequest.Contents))
	}
}

