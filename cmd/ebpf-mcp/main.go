package main

import (
	"log"
	"net/http"

	"github.com/mark3labs/mcp-go/pkg/mcp"
	"github.com/sameehj/ebpf-mcp/internal/tools"
)

func main() {
	server := mcp.NewServer("ebpf-mcp", "0.1.0")
	tools.RegisterAll(server)

	log.Println("Starting MCP server on :8080")
	http.ListenAndServe(":8080", server)
}
