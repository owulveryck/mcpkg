package mcp

import (
	"io"
	"os"
	"sync"

	"github.com/owulveryck/mcpkg/internal/kg"
)

// fileLockManager provides a mechanism to synchronize access to knowledge graph files
// to prevent concurrent read/write operations from causing data corruption.
type fileLockManager struct {
	mu    sync.Mutex
	locks map[string]*sync.RWMutex
}

// global file lock manager
var fileManager = &fileLockManager{
	locks: make(map[string]*sync.RWMutex),
}

// getFileLock returns a read-write mutex for the specified file path,
// creating one if it doesn't exist.
func (fm *fileLockManager) getFileLock(path string) *sync.RWMutex {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	if lock, exists := fm.locks[path]; exists {
		return lock
	}

	// Create a new lock for this file
	lock := &sync.RWMutex{}
	fm.locks[path] = lock
	return lock
}

// ReadKnowledgeGraph safely reads a knowledge graph from a file.
// It uses a file-level read lock to allow concurrent reads but prevent
// reads during writes.
func ReadKnowledgeGraph(path string) (*kg.KG, error) {
	// Get the file lock
	lock := fileManager.getFileLock(path)

	// Acquire a read lock
	lock.RLock()
	defer lock.RUnlock()

	// Open the file
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the knowledge graph
	graph, err := kg.ReadFrom(f)
	if err != nil {
		if err == io.EOF {
			// Empty file, return a new KG
			return kg.NewKG(""), nil
		}
		return nil, err
	}

	return graph, nil
}

// WriteKnowledgeGraph safely writes a knowledge graph to a file.
// It uses a file-level write lock to prevent concurrent writes and
// reads during the write operation.
func WriteKnowledgeGraph(path string, graph *kg.KG) error {
	// Get the file lock
	lock := fileManager.getFileLock(path)

	// Acquire a write lock
	lock.Lock()
	defer lock.Unlock()

	// Open the file for writing
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write the knowledge graph
	return kg.WriteTo(f, graph)
}

// ModifyKnowledgeGraph safely modifies a knowledge graph and writes it back to the file.
// It reads the file, applies a modification function, and writes the result back.
// The entire operation is protected by a file-level write lock.
func ModifyKnowledgeGraph(path string, modifier func(*kg.KG) error) error {
	// Get the file lock
	lock := fileManager.getFileLock(path)

	// Acquire a write lock
	lock.Lock()
	defer lock.Unlock()

	// Open the file for reading and writing
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read the knowledge graph
	graph, err := kg.ReadFrom(f)
	if err != nil {
		if err != io.EOF {
			return err
		}
		graph = kg.NewKG("")
	}

	// Apply the modification
	if err := modifier(graph); err != nil {
		return err
	}

	// Truncate the file
	if err := f.Truncate(0); err != nil {
		return err
	}

	// Reset the file pointer to the beginning
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	// Write the modified knowledge graph
	return kg.WriteTo(f, graph)
}
