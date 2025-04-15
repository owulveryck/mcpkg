package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/owulveryck/mcpkg/internal/mcp"
)

func main() {
	s := mcp.NewMCPServer()
	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
