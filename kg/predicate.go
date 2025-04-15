package kg

import (
	"gonum.org/v1/gonum/graph"
)

// Predicate is an edge of the graph. It fulfils the graph.Edge interface
type Predicate struct {
	F, T    graph.Node
	subject string
}

// From returns the from node of the edge.
func (predicate *Predicate) From() graph.Node {
	return predicate.F
}

// To returns the to node of the edge.
func (predicate *Predicate) To() graph.Node {
	return predicate.T
}

// ReversedEdge returns the edge reversal of the receiver
// if a reversal is valid for the data type.
// When a reversal is valid an edge of the same type as
// the receiver with nodes of the receiver swapped should
// be returned, otherwise the receiver should be returned
// unaltered.
func (predicate *Predicate) ReversedEdge() graph.Edge {
	return &Predicate{
		F: predicate.T,
		T: predicate.F,
	}
}
