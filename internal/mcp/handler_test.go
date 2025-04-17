package mcp

import (
	"context"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandlerDirectly(t *testing.T) {
	// Create test context
	ctx := context.Background()
	kgPath := "/tmp/mytest.kg"
	defer os.Remove(kgPath)

	// First test InsertTripleHandler
	insertRequest := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "tools/call",
			Params: struct {
				Meta *struct {
					ProgressToken mcp.ProgressToken "json:\"progressToken,omitempty\""
				} "json:\"_meta,omitempty\""
			}{
				Meta: (*struct {
					ProgressToken mcp.ProgressToken "json:\"progressToken,omitempty\""
				})(nil),
			},
		},
		Params: struct {
			Name      string                 "json:\"name\""
			Arguments map[string]interface{} "json:\"arguments,omitempty\""
			Meta      *struct {
				ProgressToken mcp.ProgressToken "json:\"progressToken,omitempty\""
			} "json:\"_meta,omitempty\""
		}{
			Name: "insert_triple",
			Arguments: map[string]interface{}{
				"knowledge_graph_path": kgPath,
				"subject":              "the sky",
				"predicate":            "is",
				"object":               "blue",
			},
			Meta: (*struct {
				ProgressToken mcp.ProgressToken "json:\"progressToken,omitempty\""
			})(nil),
		},
	}

	// Call InsertTripleHandler
	insertResult, err := InsertTripleHandler(ctx, insertRequest)
	if err != nil {
		t.Fatalf("InsertTripleHandler failed: %v", err)
	}

	if len(insertResult.Content) != 1 {
		t.Fatalf("Expected 1 content item in insert result, got %d", len(insertResult.Content))
	}

	if insertResult.IsError {
		t.Fatalf("InsertTripleHandler returned error: %s", insertResult.Content[0].(mcp.TextContent).Text)
	}

	// Now test GetRelationFromToHandler
	readRequest := mcp.ReadResourceRequest{
		Request: mcp.Request{
			Method: "resources/read",
			Params: struct {
				Meta *struct {
					ProgressToken mcp.ProgressToken "json:\"progressToken,omitempty\""
				} "json:\"_meta,omitempty\""
			}{
				Meta: (*struct {
					ProgressToken mcp.ProgressToken "json:\"progressToken,omitempty\""
				})(nil),
			},
		},
		Params: struct {
			URI       string                 "json:\"uri\""
			Arguments map[string]interface{} "json:\"arguments,omitempty\""
		}{
			URI: "graph://%2Ftmp%2Fmytest.kg?from=the%20sky&to=blue",
			Arguments: map[string]interface{}{
				"from_subject":         []string{"the sky"},
				"knowledge_graph_path": []string{kgPath},
				"to_subject":           []string{"blue"},
			},
		},
	}

	// Call GetRelationFromToHandler
	readResult, err := GetRelationFromToHandler(ctx, readRequest)
	if err != nil {
		t.Fatalf("GetRelationFromToHandler failed: %v", err)
	}

	// There should be one predicate ("is") in the result
	if readResult == nil {
		t.Fatalf("Expected non-nil result")
	}

	if len(readResult) != 1 {
		t.Fatalf("Expected 1 content item in read result, got %d", len(readResult))
	}

	// Print the type and content for debugging
	t.Logf("Result type: %T", readResult[0])
	t.Logf("Result value: %v", readResult[0])

	// Most simple test - just verify we have a result
	if readResult[0] == nil {
		t.Fatalf("Expected non-nil content")
	}
}
