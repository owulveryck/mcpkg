package mcp

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
)

// NewMCPServer exposing the knowledge graph backend.
// The server is stateless; it opens the knowledge graph file on each query.
// Therefore, it does not implement the subscribe or listChanged resources capabilities (https://modelcontextprotocol.io/specification/2024-11-05/server/resources).
func NewMCPServer() *server.MCPServer {
	log.Println("HELLO")
	// Create a new MCP server
	s := server.NewMCPServer(
		"Knowledge Graph",
		"1.0.0",
		server.WithInstructions(`# Knowledge Graph MCP Service

This service enables you to build, query, and explore a knowledge graph using subject-predicate-object triples. The graph stores semantic relationships between entities that can be used for data organization, information retrieval, recommendation systems, and much more.

## Core Capabilities

### 1. Creating and Managing Knowledge

#### Insert Triples

insert_triple(
  knowledge_graph_path="/Users/username/knowledge.kg", 
  subject="Python", 
  predicate="is_a", 
  object="Programming Language"
)


#### Remove Incorrect Information

remove_triple(
  knowledge_graph_path="/Users/username/knowledge.kg", 
  subject="Python", 
  predicate="created_by", 
  object="Microsoft"
)


### 2. Querying the Knowledge Graph

#### Find Facts About a Specific Entity

find_triples(
  knowledge_graph_path="/Users/username/knowledge.kg", 
  subject="Python"
)

→ Returns all facts about Python

#### Find Relationships Between Entities

graph:///Users/username/knowledge.kg?from=Python&to=Django

→ Returns all predicates connecting Python to Django

#### Find All Entities with a Specific Relationship

find_triples(
  knowledge_graph_path="/Users/username/knowledge.kg", 
  predicate="is_creator_of"
)

→ Returns all creator-creation relationships

#### Find All Triples Matching a Pattern

find_triples(
  knowledge_graph_path="/Users/username/knowledge.kg", 
  predicate="is_a", 
  object="Programming Language"
)

→ Returns all programming languages in the graph

### 3. Exploring the Knowledge Graph

#### Get Complete Context for an Entity

describe_entity(
  knowledge_graph_path="/Users/username/knowledge.kg", 
  entity="Python"
)

→ Returns all relationships where Python appears (both as subject and object)

## Real-World Examples

### Building a Technology Knowledge Base

insert_triple(path="/Users/username/tech.kg", subject="Python", predicate="created_by", object="Guido van Rossum")
insert_triple(path="/Users/username/tech.kg", subject="Python", predicate="first_released", object="1991")
insert_triple(path="/Users/username/tech.kg", subject="Django", predicate="written_in", object="Python")
insert_triple(path="/Users/username/tech.kg", subject="Instagram", predicate="built_with", object="Django")


→ Then query: describe_entity(path="/Users/username/tech.kg", entity="Python")

### Mapping Geographical Information

insert_triple(path="/Users/username/geo.kg", subject="Paris", predicate="is_capital_of", object="France")
insert_triple(path="/Users/username/geo.kg", subject="France", predicate="is_part_of", object="European Union")
insert_triple(path="/Users/username/geo.kg", subject="Eiffel Tower", predicate="located_in", object="Paris")


→ Then query: find_triples(path="/Users/username/geo.kg", object="Paris")

### Academic Citation Network

insert_triple(path="/Users/username/citations.kg", subject="Paper A", predicate="cites", object="Paper B")
insert_triple(path="/Users/username/citations.kg", subject="Paper C", predicate="cites", object="Paper A")
insert_triple(path="/Users/username/citations.kg", subject="Paper A", predicate="authored_by", object="Researcher X")


→ Then query: find_triples(path="/Users/username/citations.kg", predicate="cites", object="Paper A")

## Usage Notes

- Create a new .kg file or use an existing one by specifying the appropriate path
- For best results, be consistent with naming and predicates
- The knowledge graph persists your data across sessions in the files you specify
- You can build multiple specialized knowledge graphs for different domains
- Use wildcards in find_triples by omitting parameters to get broader results`),
		server.WithResourceCapabilities(false, false),
		server.WithLogging(),
		server.WithRecovery(),
	)

	s.AddResourceTemplate(GetRelationFromTo(), GetRelationFromToHandler)
	s.AddTool(InsertTriple(), InsertTripleHandler)
	s.AddTool(RemoveTriple(), RemoveTripleHandler)
	s.AddTool(FindTriples(), FindTriplesHandler)
	s.AddTool(DescribeEntity(), DescribeEntityHandler)

	return s
}
