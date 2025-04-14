package main

import (
	"fmt"
)

func main() {
	// Create a new knowledge graph
	kg := NewKG()
	
	// Create nodes
	node1 := kg.NewNode().(*Node)
	node1.Lexical = "Person"
	
	node2 := kg.NewNode().(*Node)
	node2.Lexical = "City"
	
	// Create an edge between nodes
	edge := kg.NewEdge(node1, node2)
	kg.SetEdge(edge)
	
	// Print graph information
	fmt.Println("Nodes:")
	nodes := kg.Nodes()
	nodes.Reset() // Make sure we start at the beginning
	for nodes.Next() {
		node := nodes.Node().(*Node)
		fmt.Printf("  Node ID: %d, Lexical: %s\n", node.ID(), node.Lexical)
	}
	
	// Check connectivity
	fmt.Printf("\nEdge from %d to %d exists: %v\n", 
		node1.ID(), node2.ID(), 
		kg.HasEdgeFromTo(node1.ID(), node2.ID()))
		
	fmt.Printf("Edge between %d and %d exists: %v\n", 
		node1.ID(), node2.ID(), 
		kg.HasEdgeBetween(node1.ID(), node2.ID()))
}
