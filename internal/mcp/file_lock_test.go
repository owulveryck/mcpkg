package mcp

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	
	"github.com/owulveryck/mcpkg/internal/kg"
)

func TestConcurrentFileOperations(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "file-lock-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	kgPath := filepath.Join(tempDir, "concurrent.kg")
	
	// Test concurrent reads and writes
	var wg sync.WaitGroup
	const numReaders = 5
	const numWriters = 3
	
	// Start readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := ReadKnowledgeGraph(kgPath)
			// File might not exist yet, so EOF or file not found errors are acceptable
			if err != nil && !os.IsNotExist(err) {
				t.Errorf("Read error: %v", err)
			}
		}()
	}
	
	// Start writers (adding triples)
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			subject := "Entity"
			predicate := "property"
			object := "value"
			
			err := ModifyKnowledgeGraph(kgPath, func(g *kg.KG) error {
				return g.InsertTriple(subject, predicate, object, false)
			})
			
			if err != nil {
				t.Errorf("Write error: %v", err)
			}
		}(i)
	}
	
	// Wait for all operations to complete
	wg.Wait()
	
	// Verify the file contains a valid knowledge graph with triples
	g, err := ReadKnowledgeGraph(kgPath)
	if err != nil {
		t.Fatalf("Failed to read final knowledge graph: %v", err)
	}
	
	// Check there are triples in the graph
	triples := g.FindTriples("", "", "", false)
	if len(triples) == 0 {
		t.Errorf("Expected triples in the graph, but found none")
	}
}

func TestConcurrentModifications(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "concurrent-modify-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	kgPath := filepath.Join(tempDir, "concurrent-modify.kg")
	
	// Create initial graph with a single triple
	err = ModifyKnowledgeGraph(kgPath, func(g *kg.KG) error {
		return g.InsertTriple("Base", "has", "Value", false)
	})
	if err != nil {
		t.Fatalf("Failed to create initial knowledge graph: %v", err)
	}
	
	// Test concurrent modifications to the same file
	var wg sync.WaitGroup
	const numConcurrentOps = 10
	
	// Start concurrent modify operations
	for i := 0; i < numConcurrentOps; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			subject := "Base"
			predicate := "property"
			object := "value"
			
			err := ModifyKnowledgeGraph(kgPath, func(g *kg.KG) error {
				// Insert a new triple
				err := g.InsertTriple(subject, predicate+string(rune(id+'0')), object+string(rune(id+'0')), false)
				return err
			})
			
			if err != nil {
				t.Errorf("Modification error: %v", err)
			}
		}(i)
	}
	
	// Wait for all operations to complete
	wg.Wait()
	
	// Verify the file contains a valid knowledge graph with all triples
	g, err := ReadKnowledgeGraph(kgPath)
	if err != nil {
		t.Fatalf("Failed to read final knowledge graph: %v", err)
	}
	
	// Check there are at least numConcurrentOps+1 triples in the graph (initial + added ones)
	triples := g.FindTriples("", "", "", false)
	if len(triples) < numConcurrentOps+1 {
		t.Errorf("Expected at least %d triples in the graph, but found %d", numConcurrentOps+1, len(triples))
	}
}