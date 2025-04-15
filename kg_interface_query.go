package main

import "strings"

// ListPredicatesFromNode returns all predicates from the node identified by subject
// If node is not found, it returns nil
func (kg *KG) ListPredicatesFromNode(subject string, caseSensitiveSearch bool) []*Predicate {
	// Check for nil graph
	if kg == nil {
		return nil
	}
	
	// Find the node
	node := kg.FindNode(subject, caseSensitiveSearch)
	if node == nil {
		return nil
	}

	// Check if there are outgoing edges
	fromEdges := kg.from[node.ID()]
	if fromEdges == nil || len(fromEdges) == 0 {
		return []*Predicate{} // Return empty array instead of nil if node has no outgoing edges
	}

	// Collect all predicates
	predicates := make([]*Predicate, 0, len(fromEdges))
	for _, pred := range fromEdges {
		if pred != nil {
			predicates = append(predicates, pred)
		}
	}

	return predicates
}

// ListPredicatesToNode returns all predicates to the node identified by subject
// If node is not found, it returns nil
func (kg *KG) ListPredicatesToNode(subject string, caseSensitiveSearch bool) []*Predicate {
	// Check for nil graph
	if kg == nil {
		return nil
	}
	
	// Find the node
	node := kg.FindNode(subject, caseSensitiveSearch)
	if node == nil {
		return nil
	}

	// Check if there are incoming edges
	toEdges := kg.to[node.ID()]
	if toEdges == nil || len(toEdges) == 0 {
		return []*Predicate{} // Return empty array instead of nil if node has no incoming edges
	}

	// Collect all predicates
	predicates := make([]*Predicate, 0, len(toEdges))
	for _, pred := range toEdges {
		if pred != nil {
			predicates = append(predicates, pred)
		}
	}

	return predicates
}

// PredicatesFromTo returns all the predicates that links the node identified by fromSubject to the node identified by toSubject
// returns nil if no link is found or if fromSubject or toSubject does not exists
func (kg *KG) PredicatesFromTo(fromSubject, toSubject string, caseSensitiveSearch bool) []*Predicate {
	// Check for nil graph
	if kg == nil {
		return nil
	}
	
	// Find the nodes
	fromNode := kg.FindNode(fromSubject, caseSensitiveSearch)
	if fromNode == nil {
		return nil
	}

	toNode := kg.FindNode(toSubject, caseSensitiveSearch)
	if toNode == nil {
		return nil
	}

	// Check if there are edges from fromNode
	fromEdges := kg.from[fromNode.ID()]
	if fromEdges == nil {
		return nil
	}

	// Get the predicate from fromNode to toNode
	pred := fromEdges[toNode.ID()]
	if pred == nil {
		return nil
	}

	// We need to return it as a slice since the function signature requires []*Predicate
	return []*Predicate{pred}
}

// QueryBySubject returns all predicates and objects for a given subject.
// It returns nil if the subject is not found.
func (kg *KG) QueryBySubject(subject string, caseSensitiveSearch bool) map[string][]string {
	// Check for nil graph
	if kg == nil {
		return nil
	}
	
	subjectNode := kg.FindNode(subject, caseSensitiveSearch)
	if subjectNode == nil {
		return nil
	}

	// Get all edges starting from this node
	result := make(map[string][]string)
	fromEdges := kg.from[subjectNode.ID()]
	if fromEdges == nil {
		return result // Return empty map instead of nil if node exists but has no outgoing edges
	}

	// Process all outgoing edges
	for _, pred := range fromEdges {
		if pred == nil || pred.subject == "" {
			continue
		}

		toNode := pred.T.(*Node)
		if toNode == nil || toNode.Lexical == "" {
			continue
		}

		// Initialize the slice if needed
		if result[pred.subject] == nil {
			result[pred.subject] = make([]string, 0)
		}
		
		// Add the object to the list
		result[pred.subject] = append(result[pred.subject], toNode.Lexical)
	}

	return result
}

// QueryByObject returns all subjects and predicates pointing to a given object.
// It returns nil if the object is not found.
func (kg *KG) QueryByObject(object string, caseSensitiveSearch bool) map[string][]string {
	// Check for nil graph
	if kg == nil {
		return nil
	}
	
	objectNode := kg.FindNode(object, caseSensitiveSearch)
	if objectNode == nil {
		return nil
	}

	// Get all edges ending at this node
	result := make(map[string][]string)
	toEdges := kg.to[objectNode.ID()]
	if toEdges == nil {
		return result // Return empty map instead of nil if node exists but has no incoming edges
	}

	// Process all incoming edges
	for _, pred := range toEdges {
		if pred == nil || pred.subject == "" {
			continue
		}

		fromNode := pred.F.(*Node)
		if fromNode == nil || fromNode.Lexical == "" {
			continue
		}

		// Initialize the slice if needed
		if result[pred.subject] == nil {
			result[pred.subject] = make([]string, 0)
		}
		
		// Add the subject to the list
		result[pred.subject] = append(result[pred.subject], fromNode.Lexical)
	}

	return result
}

// QueryByPredicate returns all subjects and objects connected by a given predicate.
// It returns nil if the predicate is not found.
func (kg *KG) QueryByPredicate(predicate string, caseSensitiveSearch bool) [][2]string {
	// Check for nil graph
	if kg == nil {
		return nil
	}
	
	var result [][2]string

	// Check if there are any edges
	if kg.from == nil || len(kg.from) == 0 {
		return nil
	}

	foundPredicate := false
	// Iterate through all edges
	for _, toMap := range kg.from {
		for _, pred := range toMap {
			if pred == nil || pred.subject == "" {
				continue
			}

			// Check if this is the predicate we're looking for
			matches := false
			if caseSensitiveSearch {
				matches = pred.subject == predicate
			} else {
				matches = strings.EqualFold(pred.subject, predicate)
			}

			if matches {
				foundPredicate = true
				fromNode := pred.F.(*Node)
				toNode := pred.T.(*Node)
				
				if fromNode != nil && toNode != nil && fromNode.Lexical != "" && toNode.Lexical != "" {
					// Add the subject-object pair to the result
					result = append(result, [2]string{fromNode.Lexical, toNode.Lexical})
				}
			}
		}
	}

	if !foundPredicate {
		return nil
	}

	return result
}

// FindTriples returns all triples in the knowledge graph that match the given pattern.
// Any of the parameters can be empty, which means "match any value".
func (kg *KG) FindTriples(subject, predicate, object string, caseSensitiveSearch bool) [][3]string {
	// Check for nil graph
	if kg == nil {
		return nil
	}
	
	var result [][3]string

	// If no edges exist, return empty result
	if kg.from == nil || len(kg.from) == 0 {
		return result
	}

	// Helper function to check if a string matches a pattern (empty pattern matches anything)
	matchesPattern := func(value, pattern string) bool {
		if pattern == "" {
			return true // Empty pattern matches any value
		}
		
		if caseSensitiveSearch {
			return value == pattern
		}
		
		return strings.EqualFold(value, pattern)
	}

	// Iterate through all edges
	for _, toMap := range kg.from {
		for _, pred := range toMap {
			if pred == nil || pred.subject == "" {
				continue
			}

			fromNode := pred.F.(*Node)
			toNode := pred.T.(*Node)
			
			if fromNode == nil || toNode == nil || fromNode.Lexical == "" || toNode.Lexical == "" {
				continue
			}
			
			// Check if this triple matches the pattern
			if matchesPattern(fromNode.Lexical, subject) && 
			   matchesPattern(pred.subject, predicate) && 
			   matchesPattern(toNode.Lexical, object) {
				// Add the matching triple to the result
				result = append(result, [3]string{fromNode.Lexical, pred.subject, toNode.Lexical})
			}
		}
	}

	return result
}
