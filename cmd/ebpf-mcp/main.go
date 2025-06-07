// main.go - Support both stdio and HTTP transports with well-known endpoints
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sameehj/ebpf-mcp/internal/tools"
	_ "github.com/sameehj/ebpf-mcp/internal/tools" // Import to trigger init()
)

func main() {
	var transport string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or http)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or http)")
	flag.Parse()

	// Create MCP server
	mcpServer := server.NewMCPServer(
		"ebpf-mcp",
		"0.1.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Register all your tools with MCP
	tools.RegisterAllWithMCP(mcpServer)

	// Choose transport based on flag
	if transport == "http" {
		// Create HTTP server with additional endpoints
		mux := http.NewServeMux()

		// Add the well-known endpoint for MCP discovery
		mux.HandleFunc("/.well-known/mcp/metadata.json", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"schema_version": "v1",
				"entrypoint_url": "http://localhost:8080/mcp",
				"display_name":   "eBPF MCP Server",
				"description":    "Exposes Linux kernel tools via MCP protocol",
				"tool_filter":    "all",
			}
			json.NewEncoder(w).Encode(response)
		})

		// Create the streamable HTTP server and mount it
		httpServer := server.NewStreamableHTTPServer(mcpServer)

		// Mount the MCP handler at /mcp
		mux.Handle("/mcp", httpServer)

		// Start custom HTTP server with all endpoints
		log.Printf("🔧 ebpf-mcp HTTP server listening on :8080")
		log.Printf("   MCP endpoint: http://localhost:8080/mcp")
		log.Printf("   Discovery: http://localhost:8080/.well-known/mcp/metadata.json")

		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		log.Printf("🔧 ebpf-mcp stdio server starting...")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
