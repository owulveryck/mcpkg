package mcp

import (
	"context"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/owulveryck/mcpkg/internal/kg"
)

// GetRelationFromTo returns a ResourceTemplate for retrieving relations between two nodes in a graph.
// The URI format is: graph://{knowledge_graph_path}?from={from_subject}&to={to_subject}
// where:
//   - {knowledge_graph_path} is the path to the knowledge graph file.
//   - {from_subject} is the name of the subject node.
//   - {to_subject} is the name of the object node.
func GetRelationFromTo() mcp.ResourceTemplate {
	return mcp.NewResourceTemplate(
		"graph://{knowledge_graph_path}?from={from_subject}&to={to_subject}",
		"get_predicate_from_to",
		mcp.WithTemplateDescription("Returns all the relations between two elements of the graph."),
	)
}

// GetRelationFromToHandler handles requests for retrieving relations between two nodes in a graph.
// It extracts the graph path, "from" subject, and "to" subject from the request URI,
// reads the graph from the specified file, and returns the predicates (relations) between the two nodes.
func GetRelationFromToHandler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	graphPath := extractGraphPathFromURI(request.Params.URI)
	from := extractFromFromURI(request.Params.URI)
	to := extractToFromURI(request.Params.URI)

	f, err := os.Open(graphPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	graph, err := kg.ReadFrom(f)
	if err != nil {
		return nil, err
	}
	predicates := graph.PredicatesFromTo(from, to, false)

	result := make([]mcp.ResourceContents, len(predicates))
	for i, predicate := range predicates {
		result[i], _ = mcp.AsTextResourceContents(predicate)
	}
	return result, nil
}
