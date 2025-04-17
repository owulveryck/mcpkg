package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/owulveryck/mcpkg/internal/mcp"
)

func main() {
	s := mcp.NewMCPServer()
	// Start the stdio server
	logFile, err := os.OpenFile("/tmp/mylog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Create a new logger that writes to the log file.
	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("hello")

	if err := server.ServeStdio(s, server.WithErrorLogger(logger)); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
