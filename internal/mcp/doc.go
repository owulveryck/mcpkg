// Package mcp implements a Model Context Protocol (MCP) server for knowledge graph operations.
//
// The package provides functionalities to create, interact with, and query knowledge graphs
// using the Model Context Protocol. It includes tools for inserting triples (subject-predicate-object)
// into the knowledge graph and retrieving relationships between entities.
//
// The package exposes:
//   - A stateless MCP server that opens the knowledge graph file on each query
//   - An "insert_triple" tool for adding knowledge to the graph
//   - Resource templates for querying relationships using the graph:// URI format
//
// URI Format: graph://{knowledge_graph_path}?from={from_subject}&to={to_subject}
package mcp
