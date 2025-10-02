package main

import (
	"log"

	"elastic-integration-docs-mcp/internal/mcp"
)

func main() {
	server := mcp.NewServer()
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
