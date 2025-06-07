package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mark3labs/mcp-go/pkg/mcp"
)

func main() {
	// Initialize server metadata
	server := mcp.NewServer("ebpf-mcp", "0.1.0")

	// Register your tool
	server.AddTool(&mcp.Tool{
		Name:        "trace_errors",
		Description: "Trace sched_switch for 2 seconds",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		Annotations: map[string]interface{}{
			"title":          "Trace Scheduler Switch",
			"readOnlyHint":   true,
			"openWorldHint":  false,
			"idempotentHint": true,
		},
		Handler: func(ctx context.Context, input map[string]interface{}) (*mcp.Result, error) {
			// Call your existing logic
			result := traceErrorsMock() // replace with real one if desired

			return &mcp.Result{
				Content: []mcp.ContentItem{
					{Type: "text", Text: fmt.Sprintf("Trace result: %s", result.Status)},
				},
			}, nil
		},
	})

	// Serve MCP at /rpc and metadata at /.well-known/mcp/metadata.json
	log.Println("Starting MCP server on :8080")
	http.ListenAndServe(":8080", server)
}
