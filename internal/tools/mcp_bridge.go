// internal/tools/mcp_bridge.go
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAllWithMCP converts your existing tools to MCP format
func RegisterAllWithMCP(s *server.MCPServer) {
	mu.RLock()
	defer mu.RUnlock()

	for _, tool := range toolRegistry {
		// Create MCP tool with basic info
		mcpTool := mcp.NewTool(tool.ID, mcp.WithDescription(tool.Description))

		// Capture tool in closure
		toolCopy := tool

		// Register with MCP server
		s.AddTool(mcpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			// Get all arguments as a map (your tools expect map[string]interface{})
			input := request.GetArguments()

			// Call your existing tool
			result, err := toolCopy.Call(input)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("%s failed: %v", toolCopy.ID, err)), nil
			}

			// Return result based on type
			switch v := result.(type) {
			case string:
				return mcp.NewToolResultText(v), nil
			case map[string]interface{}, []interface{}:
				// Convert complex types to JSON string
				if jsonBytes, err := json.MarshalIndent(v, "", "  "); err == nil {
					return mcp.NewToolResultText(string(jsonBytes)), nil
				}
				return mcp.NewToolResultText(fmt.Sprintf("%v", v)), nil
			default:
				return mcp.NewToolResultText(fmt.Sprintf("%v", v)), nil
			}
		})
	}
}
