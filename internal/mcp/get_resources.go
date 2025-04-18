package mcp

import (
	"context"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
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
	arguments := request.Params.Arguments
	graphPath, err := url.PathUnescape(arguments["knowledge_graph_path"].([]string)[0])
	if err != nil {
		return nil, err
	}
	from := arguments["from_subject"].([]string)[0]
	to := arguments["to_subject"].([]string)[0]

	// Read the graph using the thread-safe method
	graph, err := ReadKnowledgeGraph(graphPath)
	if err != nil {
		return nil, err
	}

	predicates := graph.PredicatesFromTo(from, to, false)

	result := make([]mcp.ResourceContents, len(predicates))
	for i, predicate := range predicates {
		result[i] = mcp.TextResourceContents{
			URI:  request.Params.URI,
			Text: predicate.Subject,
		}
	}
	return result, nil
}
