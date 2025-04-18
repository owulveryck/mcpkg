package mcp

import (
	"context"

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

	// Use the file-safe modifier function
	err := ModifyKnowledgeGraph(graphPath, func(g *kg.KG) error {
		return g.InsertTriple(subject, predicate, object, false)
	})
	
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

func RemoveTriple() mcp.Tool {
	return mcp.NewTool(
		"remove_triple",
		mcp.WithDescription("Remove a triple from the knowledge graph in the form subject predicate object"),
		mcp.WithString("knowledge_graph_path",
			mcp.Required(),
			mcp.Description("the path of the knowledge graph to interact with"),
		),
		mcp.WithString("subject",
			mcp.Required(),
			mcp.Description("the subject of the triple to remove"),
		),
		mcp.WithString("predicate",
			mcp.Required(),
			mcp.Description("the predicate of the triple to remove"),
		),
		mcp.WithString("object",
			mcp.Required(),
			mcp.Description("the object of the triple to remove"),
		),
	)
}

func RemoveTripleHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	graphPath := request.Params.Arguments["knowledge_graph_path"].(string)
	subject := request.Params.Arguments["subject"].(string)
	predicate := request.Params.Arguments["predicate"].(string)
	object := request.Params.Arguments["object"].(string)

	// First read the graph to check if it's empty
	g, err := ReadKnowledgeGraph(graphPath)
	if err != nil {
		return nil, err
	}
	
	// If the graph has no nodes, it's effectively empty
	if len(g.ListNodes()) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "The knowledge graph is empty, nothing to remove.",
				},
			},
			IsError: false,
		}, nil
	}
	
	// Check if the triple exists
	if !g.RemoveTriple(subject, predicate, object, false) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Triple not found.",
				},
			},
			IsError: false,
		}, nil
	}
	
	// If the triple exists and was removed, save the updated graph
	err = WriteKnowledgeGraph(graphPath, g)
	if err != nil {
		return nil, err
	}
	
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "Triple successfully removed.",
			},
		},
		IsError: false,
	}, nil
}

func FindTriples() mcp.Tool {
	return mcp.NewTool(
		"find_triples",
		mcp.WithDescription("Find triples in the knowledge graph by matching subject, predicate, or object (using None as wildcard)"),
		mcp.WithString("knowledge_graph_path",
			mcp.Required(),
			mcp.Description("the path of the knowledge graph to interact with"),
		),
		mcp.WithString("subject",
			mcp.Description("the subject to search for (leave empty to match any subject)"),
		),
		mcp.WithString("predicate",
			mcp.Description("the predicate to search for (leave empty to match any predicate)"),
		),
		mcp.WithString("object",
			mcp.Description("the object to search for (leave empty to match any object)"),
		),
	)
}

func FindTriplesHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	graphPath := request.Params.Arguments["knowledge_graph_path"].(string)
	
	// Extract optional parameters, using empty string as default (wildcard)
	var subject, predicate, object string
	if val, ok := request.Params.Arguments["subject"]; ok && val != nil {
		subject = val.(string)
	}
	if val, ok := request.Params.Arguments["predicate"]; ok && val != nil {
		predicate = val.(string)
	}
	if val, ok := request.Params.Arguments["object"]; ok && val != nil {
		object = val.(string)
	}

	// Read the graph using the thread-safe method
	g, err := ReadKnowledgeGraph(graphPath)
	if err != nil {
		return nil, err
	}
	
	// If the graph has no nodes, it's effectively empty
	if len(g.ListNodes()) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "The knowledge graph is empty, no triples found.",
				},
			},
			IsError: false,
		}, nil
	}

	// Find triples matching the criteria
	triples := g.FindTriples(subject, predicate, object, false)
	
	if len(triples) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "No matching triples found.",
				},
			},
			IsError: false,
		}, nil
	}

	// Format the results
	result := "Found triples:\n"
	for i, triple := range triples {
		result += "- (" + triple[0] + ", " + triple[1] + ", " + triple[2] + ")\n"
		// Add a newline after 10 triples for better readability, but not after the last one
		if i > 0 && i%10 == 0 && i < len(triples)-1 {
			result += "\n"
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
		IsError: false,
	}, nil
}

func DescribeEntity() mcp.Tool {
	return mcp.NewTool(
		"describe_entity",
		mcp.WithDescription("Get all triples involving a specific entity (as either subject or object)"),
		mcp.WithString("knowledge_graph_path",
			mcp.Required(),
			mcp.Description("the path of the knowledge graph to interact with"),
		),
		mcp.WithString("entity",
			mcp.Required(),
			mcp.Description("the entity to describe"),
		),
	)
}

func DescribeEntityHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	graphPath := request.Params.Arguments["knowledge_graph_path"].(string)
	entity := request.Params.Arguments["entity"].(string)

	// Read the graph using the thread-safe method
	g, err := ReadKnowledgeGraph(graphPath)
	if err != nil {
		return nil, err
	}
	
	// If the graph has no nodes, it's effectively empty
	if len(g.ListNodes()) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "The knowledge graph is empty, no entity information found.",
				},
			},
			IsError: false,
		}, nil
	}

	// Get all triples involving the entity
	triples := g.DescribeEntity(entity, false)
	
	if len(triples) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "No information found for entity: " + entity,
				},
			},
			IsError: false,
		}, nil
	}

	// Format the results
	result := "Entity: " + entity + "\n\n"
	
	// Group by "as subject" and "as object" for better organization
	result += "As subject:\n"
	hasSubjectRoles := false
	for _, triple := range triples {
		if triple[0] == entity {
			result += "- " + entity + " " + triple[1] + " " + triple[2] + "\n"
			hasSubjectRoles = true
		}
	}
	
	if !hasSubjectRoles {
		result += "- No triples found where " + entity + " is the subject\n"
	}
	
	result += "\nAs object:\n"
	hasObjectRoles := false
	for _, triple := range triples {
		if triple[2] == entity {
			result += "- " + triple[0] + " " + triple[1] + " " + entity + "\n"
			hasObjectRoles = true
		}
	}
	
	if !hasObjectRoles {
		result += "- No triples found where " + entity + " is the object\n"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
		IsError: false,
	}, nil
}
