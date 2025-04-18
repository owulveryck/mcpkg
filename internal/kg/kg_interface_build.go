package kg

import "strings"

// InsertTriple creates a new entry in the knowledge graph represented as a triple.
// It checks if the subject and object nodes exist using FindNode. If they don't exist,
// it creates new nodes for them. Then it creates a predicate connecting these nodes.
// The caseSensitiveSearch parameter determines if node matching is case-sensitive.
func (kg *KG) InsertTriple(subject, predicate, object string, caseSensitiveSearch bool) error {
	kg.mu.Lock()
	defer kg.mu.Unlock()
	
	// Check if the nodes already exist
	var subjectNode, objectNode *Node

	// We need to search for nodes without holding locks since we already have a write lock
	// Local implementation of FindNode to avoid lock reacquisition
	findNode := func(lexical string) *Node {
		if kg.nodes == nil || len(kg.nodes) == 0 {
			return nil
		}

		for _, node := range kg.nodes {
			if node == nil || node.Lexical == "" {
				continue
			}

			if caseSensitiveSearch {
				if node.Lexical == lexical {
					return node
				}
			} else {
				if strings.EqualFold(node.Lexical, lexical) {
					return node
				}
			}
		}

		return nil
	}

	// Get or create subject node
	subjectNode = findNode(subject)
	if subjectNode == nil {
		// Create new subject node without calling kg.NewNode() to avoid lock reacquisition
		subjectNode = &Node{
			Identifier: kg.currentID,
			Lexical:    subject,
		}
		kg.nodes[kg.currentID] = subjectNode
		kg.currentID++
	}

	// Get or create object node
	objectNode = findNode(object)
	if objectNode == nil {
		// Create new object node without calling kg.NewNode() to avoid lock reacquisition
		objectNode = &Node{
			Identifier: kg.currentID,
			Lexical:    object,
		}
		kg.nodes[kg.currentID] = objectNode
		kg.currentID++
	}

	// Create and set the predicate
	pred := &Predicate{
		F:       subjectNode,
		T:       objectNode,
		Subject: predicate,
	}

	// Add the edge to the graph without calling kg.SetEdge() to avoid lock reacquisition
	// Initialize maps if they don't exist
	if kg.from[subjectNode.ID()] == nil {
		kg.from[subjectNode.ID()] = make(map[int64]*Predicate)
	}
	if kg.to[objectNode.ID()] == nil {
		kg.to[objectNode.ID()] = make(map[int64]*Predicate)
	}

	// Set the edge in both maps
	kg.from[subjectNode.ID()][objectNode.ID()] = pred
	kg.to[objectNode.ID()][subjectNode.ID()] = pred

	return nil
}

// FindNode retrieves a node from the knowledge graph by its lexical value.
// It searches through all nodes and compares their Lexical field with the provided subject.
// The caseSensitiveSearch parameter determines if the comparison is case-sensitive.
// It returns nil if no matching node is found.
func (kg *KG) FindNode(subject string, caseSensitiveSearch bool) *Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()
	
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

// FindPredicate retrieves a predicate from the knowledge graph by its subject value.
// It searches through all predicates in the graph and compares their subject field.
// The caseSensitiveSearch parameter determines if the comparison is case-sensitive.
// It returns nil if no matching predicate is found.
func (kg *KG) FindPredicate(subject string, caseSensitiveSearch bool) *Predicate {
	kg.mu.RLock()
	defer kg.mu.RUnlock()
	
	if kg.from == nil || len(kg.from) == 0 {
		return nil
	}

	for _, toMap := range kg.from {
		for _, pred := range toMap {
			if pred == nil || pred.Subject == "" {
				continue
			}

			if caseSensitiveSearch {
				if pred.Subject == subject {
					return pred
				}
			} else {
				if strings.EqualFold(pred.Subject, subject) {
					return pred
				}
			}
		}
	}

	return nil
}

// ListAllPredicates returns all unique predicate subjects in the knowledge graph.
// It returns an empty slice if there are no predicates in the graph.
func (kg *KG) ListAllPredicates() []string {
	kg.mu.RLock()
	defer kg.mu.RUnlock()
	
	if kg.from == nil || len(kg.from) == 0 {
		return []string{}
	}

	predicates := make(map[string]struct{}) // Use a map to deduplicate predicate subjects

	for _, toMap := range kg.from {
		for _, pred := range toMap {
			if pred != nil && pred.Subject != "" {
				predicates[pred.Subject] = struct{}{}
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

// ListNodes returns the lexical values of all nodes in the knowledge graph.
// It returns an empty slice if there are no nodes in the graph.
func (kg *KG) ListNodes() []string {
	kg.mu.RLock()
	defer kg.mu.RUnlock()
	
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

// RemoveTriple removes a triple from the knowledge graph based on the provided subject, predicate, and object values.
// The caseSensitiveSearch parameter determines if the node and predicate matching is case-sensitive.
// It returns true if the triple was found and successfully removed, false otherwise.
func (kg *KG) RemoveTriple(subject, predicate, object string, caseSensitiveSearch bool) bool {
	// Check for nil graph
	if kg == nil {
		return false
	}

	kg.mu.Lock()
	defer kg.mu.Unlock()
	
	// Local implementation of FindNode to avoid lock reacquisition
	findNode := func(lexical string) *Node {
		if kg.nodes == nil || len(kg.nodes) == 0 {
			return nil
		}

		for _, node := range kg.nodes {
			if node == nil || node.Lexical == "" {
				continue
			}

			if caseSensitiveSearch {
				if node.Lexical == lexical {
					return node
				}
			} else {
				if strings.EqualFold(node.Lexical, lexical) {
					return node
				}
			}
		}

		return nil
	}

	// Find the subject and object nodes
	subjectNode := findNode(subject)
	if subjectNode == nil {
		return false
	}

	objectNode := findNode(object)
	if objectNode == nil {
		return false
	}

	// Check if a predicate exists between these nodes
	if kg.from[subjectNode.ID()] == nil {
		return false
	}

	pred := kg.from[subjectNode.ID()][objectNode.ID()]
	if pred == nil {
		return false
	}

	// Check if the predicate matches
	matches := false
	if caseSensitiveSearch {
		matches = pred.Subject == predicate
	} else {
		matches = strings.EqualFold(pred.Subject, predicate)
	}

	if !matches {
		return false
	}

	// Remove the predicate from both maps
	delete(kg.from[subjectNode.ID()], objectNode.ID())
	delete(kg.to[objectNode.ID()], subjectNode.ID())

	// If the subject node has no more outgoing edges, clean up the empty map
	if len(kg.from[subjectNode.ID()]) == 0 {
		delete(kg.from, subjectNode.ID())
	}

	// If the object node has no more incoming edges, clean up the empty map
	if len(kg.to[objectNode.ID()]) == 0 {
		delete(kg.to, objectNode.ID())
	}

	return true
}
