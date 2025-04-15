package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertTriple(t *testing.T) {
	kg := NewKG()
	assert := assert.New(t)

	// Insert a triple
	err := kg.InsertTriple("subject1", "predicate1", "object1", true)
	assert.NoError(err, "InsertTriple should not return an error")

	// Check that the nodes were created
	subject := kg.FindNode("subject1", true)
	assert.NotNil(subject, "Subject node should exist")
	assert.Equal("subject1", subject.Lexical, "Subject node should have correct lexical value")

	object := kg.FindNode("object1", true)
	assert.NotNil(object, "Object node should exist")
	assert.Equal("object1", object.Lexical, "Object node should have correct lexical value")

	// Check that the predicate was created
	predicate := kg.FindPredicate("predicate1", true)
	assert.NotNil(predicate, "Predicate should exist")
	assert.Equal("predicate1", predicate.subject, "Predicate should have correct subject")

	// Insert a triple with existing nodes
	err = kg.InsertTriple("subject1", "predicate2", "object2", true)
	assert.NoError(err, "InsertTriple should not return an error with existing nodes")

	// Check that the object node was created
	object2 := kg.FindNode("object2", true)
	assert.NotNil(object2, "New object node should exist")
	assert.Equal("object2", object2.Lexical, "New object node should have correct lexical value")

	// Test case-insensitive search (with different case)
	err = kg.InsertTriple("SUBJECT3", "PREDICATE3", "OBJECT3", false)
	assert.NoError(err, "InsertTriple should not return an error")

	// Should find the node with case-insensitive search
	subject3 := kg.FindNode("subject3", false)
	assert.NotNil(subject3, "Subject node should be found case-insensitive")
	assert.Equal("SUBJECT3", subject3.Lexical, "Original case should be preserved")
}

func TestFindNode(t *testing.T) {
	kg := NewKG()
	assert := assert.New(t)

	// Create some nodes
	node1 := kg.NewNode().(*Node)
	node1.Lexical = "Test Node"
	
	node2 := kg.NewNode().(*Node)
	node2.Lexical = "Another Node"
	
	node3 := kg.NewNode().(*Node)
	node3.Lexical = "test node" // Different case

	// Test case-sensitive search
	foundNode := kg.FindNode("Test Node", true)
	assert.NotNil(foundNode, "Node should be found with case-sensitive search")
	assert.Equal(node1.ID(), foundNode.ID(), "Found node should match the original node")
	
	// Test case-insensitive search
	foundNode = kg.FindNode("test node", false)
	assert.NotNil(foundNode, "Node should be found with case-insensitive search")
	
	// Test non-existent node
	foundNode = kg.FindNode("Non-existent Node", true)
	assert.Nil(foundNode, "Non-existent node should not be found")
	
	// Test empty graph
	emptyKG := NewKG()
	foundNode = emptyKG.FindNode("Any Node", true)
	assert.Nil(foundNode, "Empty graph should not find any nodes")
}

func TestFindPredicate(t *testing.T) {
	kg := NewKG()
	assert := assert.New(t)

	// Create nodes and predicates
	node1 := kg.NewNode().(*Node)
	node1.Lexical = "Node 1"
	
	node2 := kg.NewNode().(*Node)
	node2.Lexical = "Node 2"
	
	// Create a predicate
	pred1 := kg.NewEdge(node1, node2).(*Predicate)
	pred1.subject = "Test Predicate"
	kg.SetEdge(pred1)
	
	pred2 := kg.NewEdge(node2, node1).(*Predicate)
	pred2.subject = "test predicate" // Different case
	kg.SetEdge(pred2)
	
	// Test case-sensitive search
	foundPred := kg.FindPredicate("Test Predicate", true)
	assert.NotNil(foundPred, "Predicate should be found with case-sensitive search")
	assert.Equal("Test Predicate", foundPred.subject, "Found predicate should have correct subject")
	
	// Test case-insensitive search
	foundPred = kg.FindPredicate("test predicate", false)
	assert.NotNil(foundPred, "Predicate should be found with case-insensitive search")
	
	// Test non-existent predicate
	foundPred = kg.FindPredicate("Non-existent Predicate", true)
	assert.Nil(foundPred, "Non-existent predicate should not be found")
	
	// Test empty graph
	emptyKG := NewKG()
	foundPred = emptyKG.FindPredicate("Any Predicate", true)
	assert.Nil(foundPred, "Empty graph should not find any predicates")
}

func TestListAllPredicates(t *testing.T) {
	kg := NewKG()
	assert := assert.New(t)

	// Create nodes and predicates
	node1 := kg.NewNode().(*Node)
	node1.Lexical = "Node 1"
	
	node2 := kg.NewNode().(*Node)
	node2.Lexical = "Node 2"
	
	node3 := kg.NewNode().(*Node)
	node3.Lexical = "Node 3"
	
	// Create predicates
	pred1 := kg.NewEdge(node1, node2).(*Predicate)
	pred1.subject = "Predicate 1"
	kg.SetEdge(pred1)
	
	pred2 := kg.NewEdge(node2, node3).(*Predicate)
	pred2.subject = "Predicate 2"
	kg.SetEdge(pred2)
	
	pred3 := kg.NewEdge(node1, node3).(*Predicate)
	pred3.subject = "Predicate 1" // Duplicate to test deduplication
	kg.SetEdge(pred3)
	
	// Test listing predicates
	predicates := kg.ListAllPredicates()
	assert.Len(predicates, 2, "Should list 2 unique predicates")
	assert.Contains(predicates, "Predicate 1", "Should contain Predicate 1")
	assert.Contains(predicates, "Predicate 2", "Should contain Predicate 2")
	
	// Test empty graph
	emptyKG := NewKG()
	predicates = emptyKG.ListAllPredicates()
	assert.Empty(predicates, "Empty graph should return empty predicate list")
}

func TestListNodes(t *testing.T) {
	kg := NewKG()
	assert := assert.New(t)

	// Create some nodes
	node1 := kg.NewNode().(*Node)
	node1.Lexical = "Node 1"
	
	node2 := kg.NewNode().(*Node)
	node2.Lexical = "Node 2"
	
	node3 := kg.NewNode().(*Node)
	node3.Lexical = "Node 3"
	
	// Create a node with empty lexical value (should be filtered out)
	_ = kg.NewNode().(*Node)
	
	// Test listing nodes
	nodes := kg.ListNodes()
	assert.Len(nodes, 3, "Should list 3 nodes with non-empty lexical values")
	assert.Contains(nodes, "Node 1", "Should contain Node 1")
	assert.Contains(nodes, "Node 2", "Should contain Node 2")
	assert.Contains(nodes, "Node 3", "Should contain Node 3")
	assert.NotContains(nodes, "", "Should not contain empty node")
	
	// Test empty graph
	emptyKG := NewKG()
	nodes = emptyKG.ListNodes()
	assert.Empty(nodes, "Empty graph should return empty node list")
}