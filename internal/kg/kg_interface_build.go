package kg

import "strings"

// InsertTriple creates a new entry in the knowledge graph represented as a triple.
// It checks if the subject and object nodes exist using FindNode. If they don't exist,
// it creates new nodes for them. Then it creates a predicate connecting these nodes.
// The caseSensitiveSearch parameter determines if node matching is case-sensitive.
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
		Subject: predicate,
	}

	// Add the edge to the graph
	kg.SetEdge(pred)

	return nil
}

// FindNode retrieves a node from the knowledge graph by its lexical value.
// It searches through all nodes and compares their Lexical field with the provided subject.
// The caseSensitiveSearch parameter determines if the comparison is case-sensitive.
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

// FindPredicate retrieves a predicate from the knowledge graph by its subject value.
// It searches through all predicates in the graph and compares their subject field.
// The caseSensitiveSearch parameter determines if the comparison is case-sensitive.
// It returns nil if no matching predicate is found.
func (kg *KG) FindPredicate(subject string, caseSensitiveSearch bool) *Predicate {
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
