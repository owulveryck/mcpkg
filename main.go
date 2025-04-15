package main

import (
	"fmt"

	"github.com/owulveryck/mcpkg/internal/kg"
)

func main() {
	// Create a new knowledge graph
	graph := kg.NewKG()

	// Create nodes
	node1 := graph.NewNode().(*kg.Node)
	node1.Lexical = "Person"

	node2 := graph.NewNode().(*kg.Node)
	node2.Lexical = "City"

	// Create an edge between nodes
	edge := graph.NewEdge(node1, node2)
	graph.SetEdge(edge)

	// Print graph information
	fmt.Println("Nodes:")
	nodes := graph.Nodes()
	nodes.Reset() // Make sure we start at the beginning
	for nodes.Next() {
		node := nodes.Node().(*kg.Node)
		fmt.Printf("  Node ID: %d, Lexical: %s\n", node.ID(), node.Lexical)
	}

	// Check connectivity
	fmt.Printf("\nEdge from %d to %d exists: %v\n",
		node1.ID(), node2.ID(),
		graph.HasEdgeFromTo(node1.ID(), node2.ID()))

	fmt.Printf("Edge between %d and %d exists: %v\n",
		node1.ID(), node2.ID(),
		graph.HasEdgeBetween(node1.ID(), node2.ID()))
}
