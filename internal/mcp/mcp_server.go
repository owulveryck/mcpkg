package mcp

import (
	"github.com/mark3labs/mcp-go/server"
)

// NewMCPServer exposing the knowledge graph backend.
// The server is stateless; it opens the knowledge graph file on each query.
// Therefore, it does not implement the subscribe or listChanged resources capabilities (https://modelcontextprotocol.io/specification/2024-11-05/server/resources).
func NewMCPServer() *server.MCPServer {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Knowledge Graph",
		"1.0.0",
		server.WithInstructions(`This is a Knowledge Graph service that allows you to store and retrieve information in the form of subject-predicate-object triples.

You can use this service to:
1. Insert new knowledge into the graph using the insert_triple tool
2. Query relationships between entities using the graph:// URI format

To insert information:
- Use the insert_triple tool with a knowledge_graph_path, subject, predicate, and object
- Example: When inserting "The sky is blue", use subject="sky", predicate="is", object="blue"

To query information:
- Use the graph://{knowledge_graph_path}?from={from_subject}&to={to_subject} URI format
- This will return all relationships (predicates) between the from_subject and to_subject
- Example: graph:///path/to/graph.kg?from=sky&to=blue will return ["is"]

This knowledge graph persists your data across sessions in files specified by knowledge_graph_path.`),
		server.WithResourceCapabilities(false, false),
		server.WithLogging(),
		server.WithRecovery(),
	)

	s.AddResourceTemplate(GetRelationFromTo(), GetRelationFromToHandler)
	s.AddTool(InsertTriple(), InsertTripleHandler)

	return s
}
