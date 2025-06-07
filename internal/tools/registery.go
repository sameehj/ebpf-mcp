// internal/tools/registry.go
package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/pkg/mcp"
)

func RegisterAll(server *mcp.Server) {
	mu.RLock()
	defer mu.RUnlock()

	for _, t := range toolRegistry {
		tCopy := t // avoid closure bug

		server.AddTool(&mcp.Tool{
			Name:        tCopy.ID,
			Description: tCopy.Description,
			InputSchema: tCopy.InputSchema,
			Annotations: tCopy.Annotations,
			Handler: func(ctx context.Context, input map[string]interface{}) (*mcp.Result, error) {
				result, err := tCopy.Call(input)
				if err != nil {
					return &mcp.Result{
						IsError: true,
						Content: []mcp.ContentItem{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
					}, nil
				}

				return &mcp.Result{
					Content: []mcp.ContentItem{{Type: "text", Text: fmt.Sprintf("%v", result)}},
				}, nil
			},
		})
	}
}
