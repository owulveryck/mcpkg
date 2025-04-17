package mcp

import (
	"context"
	"io"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/owulveryck/mcpkg/internal/kg"
)

func InsertTriple() mcp.Tool {
	return mcp.NewTool(
		"insert_triple",
		mcp.WithDescription("Insert a triple in the knowledge graph in the form subject predicate object"),
		mcp.WithString("knowledge_graph_path",
			mcp.Required(),
			mcp.Description("the path of the knowledge graph to interact with"),
		),
		mcp.WithString("subject",
			mcp.Required(),
			mcp.Description("the subject of the triple"),
		),
		mcp.WithString("predicate",
			mcp.Required(),
			mcp.Description("the predicate of the triple"),
		),
		mcp.WithString("object",
			mcp.Required(),
			mcp.Description("the object of the triple"),
		),
	)
}

func InsertTripleHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	graphPath := request.Params.Arguments["knowledge_graph_path"].(string)
	subject := request.Params.Arguments["subject"].(string)
	predicate := request.Params.Arguments["predicate"].(string)
	object := request.Params.Arguments["object"].(string)

	f, err := os.OpenFile(graphPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	g, err := kg.ReadFrom(f)
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
		g = kg.NewKG()
	}
	err = g.InsertTriple(subject, predicate, object, false)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: err.Error(),
				},
			},
			IsError: true,
		}, nil
	}
	// save the kg
	err = kg.WriteTo(f, g)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "success",
			},
		},
		IsError: false,
	}, nil
}
