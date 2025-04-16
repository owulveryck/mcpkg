package mcp_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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

	readRequest := mcp.ReadResourceRequest{}
	readRequest.Params.URI = "graph://resource"

	resultRequest, err := client.ReadResource(ctx, readRequest)
	if err != nil {
		t.Fatalf("ReadResource failed: %v", err)
	}

	if len(resultRequest.Contents) != 1 {
		t.Fatalf("Expected 1 content item, got %d", len(resultRequest.Contents))
	}
}
