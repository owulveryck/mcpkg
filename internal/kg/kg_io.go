package kg

import (
	"encoding/gob"
	"encoding/json"
	"io"
)

// SerializablePredicate represents a serializable version of a Predicate.
// It stores node references as IDs rather than pointers to enable serialization.
type SerializablePredicate struct {
	FromID  int64  // ID of the source node
	ToID    int64  // ID of the target node
	Subject string // Subject of the predicate
}

// SerializableKG is a serializable representation of the knowledge graph.
// It converts the graph structure to a format that can be easily serialized.
type SerializableKG struct {
	Nodes     map[int64]*Node         // All nodes in the graph
	Edges     []SerializablePredicate // All edges in a serializable format
	CurrentID int64                   // The current ID counter for node creation
}

// WriteTo serializes and writes the knowledge graph to the provided writer
// using gob encoding. It converts the KG to a SerializableKG first to ensure
// that the graph structure can be properly encoded.
func WriteTo(w io.Writer, kg *KG) error {
	encoder := gob.NewEncoder(w)

	// Register the Node type with gob
	gob.Register(&Node{})

	// Create a serializable representation of the KG
	serialKG := SerializableKG{
		Nodes:     kg.nodes,
		Edges:     make([]SerializablePredicate, 0),
		CurrentID: kg.currentID,
	}

	// Convert predicates to serializable form
	for fromID, toMap := range kg.from {
		for toID, pred := range toMap {
			serialKG.Edges = append(serialKG.Edges, SerializablePredicate{
				FromID:  fromID,
				ToID:    toID,
				Subject: pred.Subject,
			})
		}
	}

	// Encode the serializable representation
	return encoder.Encode(serialKG)
}

// ReadFrom deserializes a knowledge graph from the provided reader
// using gob encoding. It reads a SerializableKG and converts it back
// to a proper KG structure with all node and predicate relationships.
func ReadFrom(r io.Reader) (*KG, error) {
	decoder := gob.NewDecoder(r)

	// Register the Node type with gob
	gob.Register(&Node{})

	// Create a serializable representation to decode into
	var serialKG SerializableKG

	// Decode into the serializable representation
	err := decoder.Decode(&serialKG)
	if err != nil {
		return nil, err
	}

	// Create a new KG with the decoded data
	kg := &KG{
		nodes:     serialKG.Nodes,
		from:      make(map[int64]map[int64]*Predicate),
		to:        make(map[int64]map[int64]*Predicate),
		currentID: serialKG.CurrentID,
	}

	// Reconstruct predicates
	for _, edge := range serialKG.Edges {
		fromID := edge.FromID
		toID := edge.ToID

		fromNode := kg.nodes[fromID]
		toNode := kg.nodes[toID]

		if fromNode == nil || toNode == nil {
			continue // Skip if nodes don't exist
		}

		// Create the predicate
		pred := &Predicate{
			F:       fromNode,
			T:       toNode,
			Subject: edge.Subject,
		}

		// Initialize maps if needed
		if kg.from[fromID] == nil {
			kg.from[fromID] = make(map[int64]*Predicate)
		}
		if kg.to[toID] == nil {
			kg.to[toID] = make(map[int64]*Predicate)
		}

		// Set the predicate in both maps
		kg.from[fromID][toID] = pred
		kg.to[toID][fromID] = pred
	}

	return kg, nil
}

// SaveToJSON serializes and writes the knowledge graph to the provided writer
// using JSON encoding. It converts the KG to a SerializableKG first to ensure
// that the graph structure can be properly encoded in JSON format.
func SaveToJSON(w io.Writer, kg *KG) error {
	encoder := json.NewEncoder(w)

	// Create a serializable representation of the KG
	serialKG := SerializableKG{
		Nodes:     kg.nodes,
		Edges:     make([]SerializablePredicate, 0),
		CurrentID: kg.currentID,
	}

	// Convert predicates to serializable form
	for fromID, toMap := range kg.from {
		for toID, pred := range toMap {
			serialKG.Edges = append(serialKG.Edges, SerializablePredicate{
				FromID:  fromID,
				ToID:    toID,
				Subject: pred.Subject,
			})
		}
	}

	// Encode the serializable representation
	return encoder.Encode(serialKG)
}

// ReadFromJSON deserializes a knowledge graph from the provided reader
// using JSON encoding. It reads a SerializableKG from JSON and converts it back
// to a proper KG structure with all node and predicate relationships.
func ReadFromJSON(r io.Reader) (*KG, error) {
	decoder := json.NewDecoder(r)

	// Create a serializable representation to decode into
	var serialKG SerializableKG

	// Decode into the serializable representation
	err := decoder.Decode(&serialKG)
	if err != nil {
		return nil, err
	}

	// Create a new KG with the decoded data
	kg := &KG{
		nodes:     serialKG.Nodes,
		from:      make(map[int64]map[int64]*Predicate),
		to:        make(map[int64]map[int64]*Predicate),
		currentID: serialKG.CurrentID,
	}

	// Reconstruct predicates
	for _, edge := range serialKG.Edges {
		fromID := edge.FromID
		toID := edge.ToID

		fromNode := kg.nodes[fromID]
		toNode := kg.nodes[toID]

		if fromNode == nil || toNode == nil {
			continue // Skip if nodes don't exist
		}

		// Create the predicate
		pred := &Predicate{
			F:       fromNode,
			T:       toNode,
			Subject: edge.Subject,
		}

		// Initialize maps if needed
		if kg.from[fromID] == nil {
			kg.from[fromID] = make(map[int64]*Predicate)
		}
		if kg.to[toID] == nil {
			kg.to[toID] = make(map[int64]*Predicate)
		}

		// Set the predicate in both maps
		kg.from[fromID][toID] = pred
		kg.to[toID][fromID] = pred
	}

	return kg, nil
}
