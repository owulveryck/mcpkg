package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createQueryTestGraph creates a knowledge graph for testing query functions
func createQueryTestGraph() *KG {
	kg := NewKG()

	// Create nodes
	person1 := kg.NewNode().(*Node)
	person1.Lexical = "Alice"

	person2 := kg.NewNode().(*Node)
	person2.Lexical = "Bob"

	person3 := kg.NewNode().(*Node)
	person3.Lexical = "Charlie"

	item1 := kg.NewNode().(*Node)
	item1.Lexical = "Book"

	item2 := kg.NewNode().(*Node)
	item2.Lexical = "Car"

	// Create relationships
	knows := kg.NewEdge(person1, person2).(*Predicate)
	knows.subject = "knows"
	kg.SetEdge(knows)

	friendOf := kg.NewEdge(person1, person3).(*Predicate)
	friendOf.subject = "friendOf"
	kg.SetEdge(friendOf)

	owns1 := kg.NewEdge(person1, item1).(*Predicate)
	owns1.subject = "owns"
	kg.SetEdge(owns1)

	owns2 := kg.NewEdge(person2, item2).(*Predicate)
	owns2.subject = "owns"
	kg.SetEdge(owns2)

	likes := kg.NewEdge(person3, item1).(*Predicate)
	likes.subject = "likes"
	kg.SetEdge(likes)

	return kg
}

func TestListPredicatesFromNode(t *testing.T) {
	kg := createQueryTestGraph()
	assert := assert.New(t)

	// Test for existing node with predicates
	predicates := kg.ListPredicatesFromNode("Alice", true)
	assert.NotNil(predicates, "Predicates should not be nil for existing node")
	assert.Equal(3, len(predicates), "Alice should have 3 outgoing predicates")
	
	// Check the predicate subjects
	subjects := make([]string, 0, len(predicates))
	for _, pred := range predicates {
		subjects = append(subjects, pred.subject)
	}
	sort.Strings(subjects)
	assert.Equal([]string{"friendOf", "knows", "owns"}, subjects, "Predicate subjects should match")

	// Test for existing node with no predicates
	noPredicates := kg.ListPredicatesFromNode("Car", true)
	assert.NotNil(noPredicates, "Result should not be nil for node with no predicates")
	assert.Empty(noPredicates, "Car should have no outgoing predicates")

	// Test for non-existing node
	nilPredicates := kg.ListPredicatesFromNode("NonExistent", true)
	assert.Nil(nilPredicates, "Result should be nil for non-existent node")

	// Test case-insensitive search
	caseInsensitivePredicates := kg.ListPredicatesFromNode("alice", false)
	assert.NotNil(caseInsensitivePredicates, "Predicates should be found with case-insensitive search")
	assert.Equal(3, len(caseInsensitivePredicates), "alice (case-insensitive) should have 3 outgoing predicates")
}

func TestListPredicatesToNode(t *testing.T) {
	kg := createQueryTestGraph()
	assert := assert.New(t)

	// Test for existing node with incoming predicates
	predicates := kg.ListPredicatesToNode("Book", true)
	assert.NotNil(predicates, "Predicates should not be nil for existing node")
	assert.Equal(2, len(predicates), "Book should have 2 incoming predicates")
	
	// Check the predicate subjects
	subjects := make([]string, 0, len(predicates))
	for _, pred := range predicates {
		subjects = append(subjects, pred.subject)
	}
	sort.Strings(subjects)
	assert.Equal([]string{"likes", "owns"}, subjects, "Predicate subjects should match")

	// Test for existing node with no incoming predicates
	noPredicates := kg.ListPredicatesToNode("Alice", true)
	assert.NotNil(noPredicates, "Result should not be nil for node with no incoming predicates")
	assert.Empty(noPredicates, "Alice should have no incoming predicates")

	// Test for non-existing node
	nilPredicates := kg.ListPredicatesToNode("NonExistent", true)
	assert.Nil(nilPredicates, "Result should be nil for non-existent node")

	// Test case-insensitive search
	caseInsensitivePredicates := kg.ListPredicatesToNode("book", false)
	assert.NotNil(caseInsensitivePredicates, "Predicates should be found with case-insensitive search")
	assert.Equal(2, len(caseInsensitivePredicates), "book (case-insensitive) should have 2 incoming predicates")
}

func TestPredicatesFromTo(t *testing.T) {
	kg := createQueryTestGraph()
	assert := assert.New(t)

	// Test for existing connection
	predicates := kg.PredicatesFromTo("Alice", "Bob", true)
	assert.NotNil(predicates, "Predicates should not be nil for existing connection")
	assert.Equal(1, len(predicates), "Should find 1 predicate from Alice to Bob")
	assert.Equal("knows", predicates[0].subject, "Predicate subject should be 'knows'")

	// Test for non-existing connection
	nilPredicates := kg.PredicatesFromTo("Bob", "Alice", true)
	assert.Nil(nilPredicates, "Result should be nil for non-existent connection")

	// Test for non-existing nodes
	assert.Nil(kg.PredicatesFromTo("NonExistent", "Alice", true), "Result should be nil if from node doesn't exist")
	assert.Nil(kg.PredicatesFromTo("Alice", "NonExistent", true), "Result should be nil if to node doesn't exist")

	// Test case-insensitive search
	caseInsensitivePredicates := kg.PredicatesFromTo("alice", "bob", false)
	assert.NotNil(caseInsensitivePredicates, "Predicates should be found with case-insensitive search")
	assert.Equal(1, len(caseInsensitivePredicates), "Should find 1 predicate from Alice to Bob with case-insensitive search")
}

func TestQueryBySubject(t *testing.T) {
	kg := createQueryTestGraph()
	assert := assert.New(t)

	// Test for existing subject with predicates
	results := kg.QueryBySubject("Alice", true)
	assert.NotNil(results, "Results should not be nil for existing subject")
	assert.Equal(3, len(results), "Alice should have 3 different predicate types")
	
	// Check specific predicate-object pairs
	assert.Contains(results, "knows", "Alice should have 'knows' predicate")
	assert.Contains(results["knows"], "Bob", "Alice knows Bob")
	
	assert.Contains(results, "friendOf", "Alice should have 'friendOf' predicate")
	assert.Contains(results["friendOf"], "Charlie", "Alice is a friend of Charlie")
	
	assert.Contains(results, "owns", "Alice should have 'owns' predicate")
	assert.Contains(results["owns"], "Book", "Alice owns a Book")

	// Test for existing subject with no predicates
	emptyResults := kg.QueryBySubject("Car", true)
	assert.NotNil(emptyResults, "Results should not be nil for subject with no predicates")
	assert.Empty(emptyResults, "Car should have no predicates")

	// Test for non-existing subject
	nilResults := kg.QueryBySubject("NonExistent", true)
	assert.Nil(nilResults, "Results should be nil for non-existent subject")

	// Test case-insensitive search
	caseInsensitiveResults := kg.QueryBySubject("alice", false)
	assert.NotNil(caseInsensitiveResults, "Results should be found with case-insensitive search")
	assert.Equal(3, len(caseInsensitiveResults), "alice (case-insensitive) should have 3 different predicate types")
}

func TestQueryByObject(t *testing.T) {
	kg := createQueryTestGraph()
	assert := assert.New(t)

	// Test for existing object with incoming predicates
	results := kg.QueryByObject("Book", true)
	assert.NotNil(results, "Results should not be nil for existing object")
	assert.Equal(2, len(results), "Book should have 2 different predicate types pointing to it")
	
	// Check specific predicate-subject pairs
	assert.Contains(results, "owns", "Book should have 'owns' predicate")
	assert.Contains(results["owns"], "Alice", "Alice owns the Book")
	
	assert.Contains(results, "likes", "Book should have 'likes' predicate")
	assert.Contains(results["likes"], "Charlie", "Charlie likes the Book")

	// Test for existing object with no incoming predicates
	emptyResults := kg.QueryByObject("Alice", true)
	assert.NotNil(emptyResults, "Results should not be nil for object with no incoming predicates")
	assert.Empty(emptyResults, "Alice should have no incoming predicates")

	// Test for non-existing object
	nilResults := kg.QueryByObject("NonExistent", true)
	assert.Nil(nilResults, "Results should be nil for non-existent object")

	// Test case-insensitive search
	caseInsensitiveResults := kg.QueryByObject("book", false)
	assert.NotNil(caseInsensitiveResults, "Results should be found with case-insensitive search")
	assert.Equal(2, len(caseInsensitiveResults), "book (case-insensitive) should have 2 different predicate types pointing to it")
}

func TestQueryByPredicate(t *testing.T) {
	kg := createQueryTestGraph()
	assert := assert.New(t)

	// Test for existing predicate
	results := kg.QueryByPredicate("owns", true)
	assert.NotNil(results, "Results should not be nil for existing predicate")
	assert.Equal(2, len(results), "There should be 2 'owns' relationships")
	
	// Find all subject-object pairs
	pairs := make(map[string]string)
	for _, pair := range results {
		pairs[pair[0]] = pair[1]
	}
	
	assert.Contains(pairs, "Alice", "Alice should be a subject of 'owns'")
	assert.Equal("Book", pairs["Alice"], "Alice owns a Book")
	
	assert.Contains(pairs, "Bob", "Bob should be a subject of 'owns'")
	assert.Equal("Car", pairs["Bob"], "Bob owns a Car")

	// Test for non-existing predicate
	nilResults := kg.QueryByPredicate("NonExistent", true)
	assert.Nil(nilResults, "Results should be nil for non-existent predicate")

	// Test for empty graph
	emptyKG := NewKG()
	emptyResults := emptyKG.QueryByPredicate("owns", true)
	assert.Nil(emptyResults, "Results should be nil for empty graph")

	// Test case-insensitive search
	caseInsensitiveResults := kg.QueryByPredicate("OWNS", false)
	assert.NotNil(caseInsensitiveResults, "Results should be found with case-insensitive search")
	assert.Equal(2, len(caseInsensitiveResults), "There should be 2 'owns' relationships with case-insensitive search")
}

func TestFindTriples(t *testing.T) {
	kg := createQueryTestGraph()
	assert := assert.New(t)

	// Test finding specific triple
	specificTriple := kg.FindTriples("Alice", "knows", "Bob", true)
	assert.Equal(1, len(specificTriple), "Should find exactly 1 triple matching all criteria")
	assert.Equal([3]string{"Alice", "knows", "Bob"}, specificTriple[0], "The triple should match exactly")

	// Test with only subject specified
	aliceTriples := kg.FindTriples("Alice", "", "", true)
	assert.Equal(3, len(aliceTriples), "Should find 3 triples with Alice as subject")
	
	// Test with only predicate specified
	ownsTriples := kg.FindTriples("", "owns", "", true)
	assert.Equal(2, len(ownsTriples), "Should find 2 triples with owns as predicate")
	
	// Test with only object specified
	bookTriples := kg.FindTriples("", "", "Book", true)
	assert.Equal(2, len(bookTriples), "Should find 2 triples with Book as object")
	
	// Test with subject and predicate specified
	aliceOwnsTriples := kg.FindTriples("Alice", "owns", "", true)
	assert.Equal(1, len(aliceOwnsTriples), "Should find 1 triple with Alice as subject and owns as predicate")
	assert.Equal([3]string{"Alice", "owns", "Book"}, aliceOwnsTriples[0], "The triple should match Alice owns Book")
	
	// Test with non-matching criteria
	nonMatchingTriples := kg.FindTriples("Alice", "knows", "Charlie", true)
	assert.Empty(nonMatchingTriples, "Should find no triples with non-matching criteria")
	
	// Test with empty graph
	emptyKG := NewKG()
	emptyTriples := emptyKG.FindTriples("", "", "", true)
	assert.Empty(emptyTriples, "Should find no triples in empty graph")
	
	// Test case-insensitive search
	caseInsensitiveTriples := kg.FindTriples("alice", "knows", "bob", false)
	assert.Equal(1, len(caseInsensitiveTriples), "Should find 1 triple with case-insensitive search")
	assert.Equal([3]string{"Alice", "knows", "Bob"}, caseInsensitiveTriples[0], "The triple should match exactly with original case preserved")
}

// Test edge cases and error conditions
func TestQueryEdgeCases(t *testing.T) {
	kg := createQueryTestGraph()
	assert := assert.New(t)
	
	// Test with empty subject/object/predicate
	emptySubject := kg.QueryBySubject("", true)
	assert.Nil(emptySubject, "Empty subject should return nil")
	
	emptyObject := kg.QueryByObject("", true)
	assert.Nil(emptyObject, "Empty object should return nil")
	
	// Test with nil graph
	var nilKG *KG
	assert.Nil(nilKG.ListPredicatesFromNode("Alice", true), "Nil graph should handle method calls safely")
	assert.Nil(nilKG.ListPredicatesToNode("Alice", true), "Nil graph should handle method calls safely")
	assert.Nil(nilKG.PredicatesFromTo("Alice", "Bob", true), "Nil graph should handle method calls safely")
	assert.Nil(nilKG.QueryBySubject("Alice", true), "Nil graph should handle method calls safely")
	assert.Nil(nilKG.QueryByObject("Book", true), "Nil graph should handle method calls safely")
	assert.Nil(nilKG.QueryByPredicate("owns", true), "Nil graph should handle method calls safely")
	assert.Nil(nilKG.FindTriples("Alice", "owns", "Book", true), "Nil graph should handle method calls safely")
}