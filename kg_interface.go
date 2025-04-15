package main

import (
	"strings"
)

// InsertTriple creates a new entry in the knowledge graph.
// It checks if the subject and object nodes exist using FindNode and if the predicate exists using FindPredicate.
func (kg *KG) InsertTriple(subject, predicate, object string, caseSensitiveSearch bool) error {
	// Check if the nodes already exist
	var subjectNode, objectNode *Node
	
	// Get or create subject node
	subjectNode = kg.FindNode(subject, caseSensitiveSearch)
	if subjectNode == nil {
		// Create new subject node
		newNode := kg.NewNode().(*Node)
		newNode.Lexical = subject
		subjectNode = newNode
	}
	
	// Get or create object node
	objectNode = kg.FindNode(object, caseSensitiveSearch)
	if objectNode == nil {
		// Create new object node
		newNode := kg.NewNode().(*Node)
		newNode.Lexical = object
		objectNode = newNode
	}
	
	// Create and set the predicate
	pred := &Predicate{
		F:       subjectNode,
		T:       objectNode,
		subject: predicate,
	}
	
	// Add the edge to the graph
	kg.SetEdge(pred)
	
	return nil
}

// FindNode retrieves a node from the knowledge graph by its subject.
// It returns nil if no matching node is found.
func (kg *KG) FindNode(subject string, caseSensitiveSearch bool) *Node {
	if kg.nodes == nil || len(kg.nodes) == 0 {
		return nil
	}
	
	for _, node := range kg.nodes {
		if node == nil || node.Lexical == "" {
			continue
		}
		
		if caseSensitiveSearch {
			if node.Lexical == subject {
				return node
			}
		} else {
			if strings.EqualFold(node.Lexical, subject) {
				return node
			}
		}
	}
	
	return nil
}

// FindPredicate retrieves a predicate from the knowledge graph by its subject.
// It returns nil if no matching predicate is found.
func (kg *KG) FindPredicate(subject string, caseSensitiveSearch bool) *Predicate {
	if kg.from == nil || len(kg.from) == 0 {
		return nil
	}
	
	for _, toMap := range kg.from {
		for _, pred := range toMap {
			if pred == nil || pred.subject == "" {
				continue
			}
			
			if caseSensitiveSearch {
				if pred.subject == subject {
					return pred
				}
			} else {
				if strings.EqualFold(pred.subject, subject) {
					return pred
				}
			}
		}
	}
	
	return nil
}

// ListAllPredicates returns all the subject of all predicated in the knowledge graph
func (kg *KG) ListAllPredicates() []string {
	if kg.from == nil || len(kg.from) == 0 {
		return []string{}
	}
	
	predicates := make(map[string]struct{}) // Use a map to deduplicate predicate subjects
	
	for _, toMap := range kg.from {
		for _, pred := range toMap {
			if pred != nil && pred.subject != "" {
				predicates[pred.subject] = struct{}{}
			}
		}
	}
	
	// Convert the map to a slice
	result := make([]string, 0, len(predicates))
	for pred := range predicates {
		result = append(result, pred)
	}
	
	return result
}

// ListNodes returns all nodes' subjects in the knowledge graph
func (kg *KG) ListNodes() []string {
	if kg.nodes == nil || len(kg.nodes) == 0 {
		return []string{}
	}
	
	nodes := make([]string, 0, len(kg.nodes))
	for _, node := range kg.nodes {
		if node != nil && node.Lexical != "" {
			nodes = append(nodes, node.Lexical)
		}
	}
	
	return nodes
}
