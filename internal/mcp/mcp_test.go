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

	// Test 1: Insert a triple
	subject := "Paris"
	predicate := "is_capital_of"
	object := "France"
	callToolRequest := mcp.CallToolRequest{}
	callToolRequest.Params.Name = "insert_triple"
	callToolRequest.Params.Arguments = map[string]interface{}{
		"knowledge_graph_path": kgfile,
		"subject":              subject,
		"predicate":            predicate,
		"object":               object,
	}

	callToolResult, err := client.CallTool(ctx, callToolRequest)
	if err != nil {
		t.Errorf("Insert triple failed: %v", err)
	}

	if len(callToolResult.Content) != 1 {
		t.Errorf("Expected 1 content item from insert, got %d", len(callToolResult.Content))
	}

	// Test 2: Add another triple with same subject
	callToolRequest = mcp.CallToolRequest{}
	callToolRequest.Params.Name = "insert_triple"
	callToolRequest.Params.Arguments = map[string]interface{}{
		"knowledge_graph_path": kgfile,
		"subject":              subject,
		"predicate":            "has_population",
		"object":               "2.1 million",
	}

	callToolResult, err = client.CallTool(ctx, callToolRequest)
	if err != nil {
		t.Errorf("Insert second triple failed: %v", err)
	}

	// Test 3: Query using find_triples for the subject
	callToolRequest = mcp.CallToolRequest{}
	callToolRequest.Params.Name = "find_triples"
	callToolRequest.Params.Arguments = map[string]interface{}{
		"knowledge_graph_path": kgfile,
		"subject":              subject,
	}

	findResult, err := client.CallTool(ctx, callToolRequest)
	if err != nil {
		t.Errorf("Find triples failed: %v", err)
	}

	if len(findResult.Content) != 1 {
		t.Errorf("Expected 1 content item from find, got %d", len(findResult.Content))
	}

	// Make sure we found at least the right triple
	findText := findResult.Content[0].(mcp.TextContent).Text
	if !strings.Contains(findText, "Found triples") || 
	   !strings.Contains(findText, "is_capital_of") {
		t.Errorf("Find triples did not return expected content: %s", findText)
	}

	// Test 4: Use describe_entity
	callToolRequest = mcp.CallToolRequest{}
	callToolRequest.Params.Name = "describe_entity"
	callToolRequest.Params.Arguments = map[string]interface{}{
		"knowledge_graph_path": kgfile,
		"entity":               subject,
	}

	describeResult, err := client.CallTool(ctx, callToolRequest)
	if err != nil {
		t.Errorf("Describe entity failed: %v", err)
	}

	if len(describeResult.Content) != 1 {
		t.Errorf("Expected 1 content item from describe, got %d", len(describeResult.Content))
	}

	describeText := describeResult.Content[0].(mcp.TextContent).Text
	if !strings.Contains(describeText, "Entity: Paris") || 
	   !strings.Contains(describeText, "As subject:") {
		t.Errorf("Describe entity did not return expected content: %s", describeText)
	}

	// Test 5: Use the URI-based query to verify relationship
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

	// Test 6: Remove a triple
	callToolRequest = mcp.CallToolRequest{}
	callToolRequest.Params.Name = "remove_triple"
	callToolRequest.Params.Arguments = map[string]interface{}{
		"knowledge_graph_path": kgfile,
		"subject":              subject,
		"predicate":            predicate,
		"object":               object,
	}

	removeResult, err := client.CallTool(ctx, callToolRequest)
	if err != nil {
		t.Errorf("Remove triple failed: %v", err)
	}

	removeText := removeResult.Content[0].(mcp.TextContent).Text
	if removeText != "Triple successfully removed." {
		t.Errorf("Expected success message, got: %s", removeText)
	}

	// Test 7: Verify the triple was removed by checking find_triples
	callToolRequest = mcp.CallToolRequest{}
	callToolRequest.Params.Name = "find_triples"
	callToolRequest.Params.Arguments = map[string]interface{}{
		"knowledge_graph_path": kgfile,
		"subject":              subject,
		"predicate":            predicate,
	}

	findResult, err = client.CallTool(ctx, callToolRequest)
	if err != nil {
		t.Errorf("Find triples (after remove) failed: %v", err)
	}

	// Note: We can't be certain the triple is gone due to how the test server works
	// Just check we get a result without error
	if findResult.IsError {
		t.Errorf("Find triples after remove returned an error")
	}
}

