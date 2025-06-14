package kg

import (
	"sync"

	"gonum.org/v1/gonum/graph"
)

// KG holds the knowledge graph structure
// It is safe for concurrent use as all operations are protected by a mutex.
type KG struct {
	SystemPrompt string
	nodes        map[int64]*Node
	from         map[int64]map[int64]*Predicate
	to           map[int64]map[int64]*Predicate

	currentID int64
	mu        sync.RWMutex // protects concurrent access to the graph
}

// NewKG creates and initializes a new empty knowledge graph.
// It returns a pointer to the new KG with initialized maps for nodes, from, and to relationships.
func NewKG(systemPrompt string) *KG {
	return &KG{
		SystemPrompt: systemPrompt,
		nodes:        make(map[int64]*Node),
		from:         make(map[int64]map[int64]*Predicate),
		to:           make(map[int64]map[int64]*Predicate),
	}
}

// Node returns the node with the given ID if it exists
// in the graph, and nil otherwise.
func (kg *KG) Node(id int64) graph.Node {
	kg.mu.RLock()
	defer kg.mu.RUnlock()
	return kg.nodes[id]
}

// NodeList implements the graph.Nodes interface
type NodeList struct {
	nodes []*Node
	pos   int
}

// NewNodeList creates a new NodeList from the provided slice of nodes.
// It initializes the position to -1, so Next() must be called before the first Node() access.
func NewNodeList(nodes []*Node) *NodeList {
	return &NodeList{
		nodes: nodes,
		pos:   -1,
	}
}

func (n *NodeList) Len() int {
	return len(n.nodes)
}

func (n *NodeList) Next() bool {
	n.pos++
	return n.pos < len(n.nodes)
}

func (n *NodeList) Reset() {
	n.pos = -1
}

func (n *NodeList) Node() graph.Node {
	if n.pos >= len(n.nodes) || n.pos < 0 {
		return nil
	}
	return n.nodes[n.pos]
}

// Nodes returns all the nodes in the graph.
//
// Nodes must not return nil.
func (kg *KG) Nodes() graph.Nodes {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	if len(kg.nodes) == 0 {
		return NewNodeList(nil)
	}

	nodes := make([]*Node, 0, len(kg.nodes))
	for _, node := range kg.nodes {
		nodes = append(nodes, node)
	}
	return NewNodeList(nodes)
}

// From returns all nodes that can be reached directly
// from the node with the given ID.
//
// From must not return nil.
func (kg *KG) From(id int64) graph.Nodes {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	if kg.from[id] == nil {
		return NewNodeList(nil)
	}

	nodes := make([]*Node, 0, len(kg.from[id]))
	for nid := range kg.from[id] {
		if node, ok := kg.nodes[nid]; ok {
			nodes = append(nodes, node)
		}
	}
	return NewNodeList(nodes)
}

// HasEdgeBetween returns whether an edge exists between
// nodes with IDs xid and yid without considering direction.
func (kg *KG) HasEdgeBetween(xid int64, yid int64) bool {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	// Check if there's an edge from xid to yid
	if _, ok := kg.from[xid]; ok {
		if _, ok := kg.from[xid][yid]; ok {
			return true
		}
	}

	// Check if there's an edge from yid to xid
	if _, ok := kg.from[yid]; ok {
		if _, ok := kg.from[yid][xid]; ok {
			return true
		}
	}

	return false
}

// Edge returns the edge from u to v, with IDs uid and vid,
// if such an edge exists and nil otherwise. The node v
// must be directly reachable from u as defined by the
// From method.
func (kg *KG) Edge(uid int64, vid int64) graph.Edge {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	if _, ok := kg.from[uid]; !ok {
		return nil
	}

	pred, ok := kg.from[uid][vid]
	if !ok {
		return nil
	}

	return pred
}

// HasEdgeFromTo returns whether an edge exists
// in the graph from u to v with IDs uid and vid.
func (kg *KG) HasEdgeFromTo(uid int64, vid int64) bool {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	if _, ok := kg.from[uid]; !ok {
		return false
	}

	_, ok := kg.from[uid][vid]
	return ok
}

// To returns all nodes that can reach directly
// to the node with the given ID.
//
// To must not return nil.
func (kg *KG) To(id int64) graph.Nodes {
	kg.mu.RLock()
	defer kg.mu.RUnlock()

	if kg.to[id] == nil {
		return NewNodeList(nil)
	}

	nodes := make([]*Node, 0, len(kg.to[id]))
	for nid := range kg.to[id] {
		if node, ok := kg.nodes[nid]; ok {
			nodes = append(nodes, node)
		}
	}
	return NewNodeList(nodes)
}
