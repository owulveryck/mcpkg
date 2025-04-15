# MCPKG - Model Context Protocol Knowledge Graph

MCPKG is a Go implementation of a knowledge graph system that is exposed through the Model Context Protocol (MCP). It provides a simple yet powerful way to store, manage, and query semantic information in the form of subject-predicate-object triples.

## Features

- Directed graph implementation for storing structured information
- Support for creating and querying semantic triples
- Persistent storage through serialization
- MCP server interface for programmatic access
- Custom URI format for graph queries

## Components

### Knowledge Graph (KG)

The core data structure that:
- Stores entities as nodes
- Represents relationships as predicates (edges)
- Provides methods for inserting and querying triples
- Supports serialization and deserialization

### MCP Server

An interface layer that:
- Exposes the knowledge graph as an MCP server
- Provides tools for inserting triples
- Supports graph queries via URI format
- Implements a stateless design for reliability

## Usage

### Inserting Information

Use the `insert_triple` tool with subject, predicate, and object parameters to add information to the graph.

### Querying Information

Use the `graph://` URI format to find relationships between entities.

## Dependencies

- Go 1.24+
- github.com/mark3labs/mcp-go
- gonum.org/v1/gonum
- github.com/stretchr/testify (for testing)

## License

MIT

