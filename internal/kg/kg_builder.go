package kg

import "gonum.org/v1/gonum/graph"

// NewNode returns a new Node with a unique arbitrary ID.
// This method satisfies the graph.NodeAdder interface.
// It increments the currentID counter to ensure unique IDs.
func (kg *KG) NewNode() graph.Node {
	kg.mu.Lock()
	defer kg.mu.Unlock()
	
	n := &Node{
		Identifier: kg.currentID,
	}
	kg.nodes[kg.currentID] = n
	kg.currentID++
	return n
}

// AddNode adds a node to the graph. AddNode panics if
// the added node ID matches an existing node ID.
func (kg *KG) AddNode(n graph.Node) {
	kg.mu.Lock()
	defer kg.mu.Unlock()
	
	if _, exists := kg.nodes[n.ID()]; exists {
		panic("graph: AddNode: node ID collision")
	}
	node, ok := n.(*Node)
	if !ok {
		node = &Node{Identifier: n.ID()}
	}
	kg.nodes[n.ID()] = node
}

// NewEdge returns a new Edge (Predicate) from the source to the destination node.
// This method satisfies the graph.EdgeAdder interface.
func (kg *KG) NewEdge(from graph.Node, to graph.Node) graph.Edge {
	return &Predicate{
		F: from,
		T: to,
	}
}

// SetEdge adds an edge from one node to another.
// If the graph supports node addition the nodes
// will be added if they do not exist, otherwise
// SetEdge will panic.
// The behavior of an EdgeAdder when the IDs
// returned by e.From() and e.To() are equal is
// implementation-dependent.
// Whether e, e.From() and e.To() are stored
// within the graph is implementation dependent.
func (kg *KG) SetEdge(e graph.Edge) {
	kg.mu.Lock()
	defer kg.mu.Unlock()
	
	from := e.From()
	to := e.To()

	// Add nodes if they don't exist
	if _, exists := kg.nodes[from.ID()]; !exists {
		// We're already holding the lock, so don't call AddNode which would try to acquire it again
		node, ok := from.(*Node)
		if !ok {
			node = &Node{Identifier: from.ID()}
		}
		kg.nodes[from.ID()] = node
	}

	if _, exists := kg.nodes[to.ID()]; !exists {
		// We're already holding the lock, so don't call AddNode which would try to acquire it again
		node, ok := to.(*Node)
		if !ok {
			node = &Node{Identifier: to.ID()}
		}
		kg.nodes[to.ID()] = node
	}

	// Create the predicate
	pred, ok := e.(*Predicate)
	if !ok {
		pred = &Predicate{
			F: from,
			T: to,
		}
	}

	// Initialize maps if they don't exist
	if kg.from[from.ID()] == nil {
		kg.from[from.ID()] = make(map[int64]*Predicate)
	}
	if kg.to[to.ID()] == nil {
		kg.to[to.ID()] = make(map[int64]*Predicate)
	}

	// Set the edge in both maps
	kg.from[from.ID()][to.ID()] = pred
	kg.to[to.ID()][from.ID()] = pred
}
