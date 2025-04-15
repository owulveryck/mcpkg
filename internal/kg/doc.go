// Package kg implements a directed knowledge graph for storing and querying structured information.
//
// The package provides a graph-based data structure that represents semantic information as triples
// in the form of (subject, predicate, object). This structure allows for storing relationships
// between entities and querying them efficiently.
//
// The knowledge graph is built on top of the gonum/graph package and implements its interfaces
// for graph operations. It offers features such as:
//
// - Creating and managing nodes that represent entities
//
// - Establishing predicates (edges) between nodes to represent relationships
//
// - Querying the graph by subject, predicate, or object
//
// - Finding specific triples or relationships between nodes
//
// - Serialization and deserialization for persistent storage
//
// The package is designed to be used for various knowledge representation tasks, such as
// storing metadata, dependency information, relationships between resources, and more.
//
// Example usage:
//
//	kg := NewKG()
//	kg.InsertTriple("Go", "is", "Programming Language", true)
//	kg.InsertTriple("Go", "created by", "Google", true)
//	results := kg.QueryBySubject("Go", true)
//	// results will contain predicates "is" and "created by" with their objects
package kg
