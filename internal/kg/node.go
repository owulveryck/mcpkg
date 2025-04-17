package kg

// Node represents a vertex in the knowledge graph.
// Each node has a unique ID and a Lexical field that holds its string representation.
type Node struct {
	Identifier int64 `json:"id"` // Using ID_ with json tag for serialization
	Lexical    string
}

// ID returns the unique identifier of the node.
// This method satisfies the graph.Node interface.
func (node *Node) ID() int64 {
	return node.Identifier
}
