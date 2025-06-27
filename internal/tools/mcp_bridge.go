// internal/tools/mcp_bridge.go
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAllWithMCP converts your existing tools to MCP format
func RegisterAllWithMCP(s *server.MCPServer) {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("[DEBUG] Registering %d tools with MCP server", len(toolRegistry))

	for _, tool := range toolRegistry {
		log.Printf("[DEBUG] Registering tool: %s - %s", tool.ID, tool.Description)

		// Create MCP tool with description
		mcpTool := mcp.NewTool(
			tool.ID,
			mcp.WithDescription(tool.Description),
		)

		// Log schema information (for debugging, but can't register it directly)
		if tool.InputSchema != nil {
			if schemaBytes, err := json.MarshalIndent(tool.InputSchema, "", "  "); err == nil {
				log.Printf("[DEBUG] Tool %s input schema:\n%s", tool.ID, string(schemaBytes))
			}
		} else {
			log.Printf("[DEBUG] Tool %s has no input schema", tool.ID)
		}

		// Capture tool in closure
		toolCopy := tool

		// Register with MCP server
		s.AddTool(mcpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			log.Printf("[DEBUG] ==> Tool %s called", toolCopy.ID)

			// Get all arguments as a map
			input := request.GetArguments()

			// Enhanced debug logging
			log.Printf("[DEBUG] Tool %s raw arguments: %+v", toolCopy.ID, input)
			log.Printf("[DEBUG] Tool %s arguments type: %T", toolCopy.ID, input)

			if input != nil {
				log.Printf("[DEBUG] Tool %s arguments length: %d", toolCopy.ID, len(input))
				for k, v := range input {
					log.Printf("[DEBUG] Tool %s arg[%s] = %+v (type: %T)", toolCopy.ID, k, v, v)

					// Deep inspection for nested objects
					if nested, ok := v.(map[string]interface{}); ok {
						log.Printf("[DEBUG] Tool %s arg[%s] is nested object with %d keys", toolCopy.ID, k, len(nested))
						for nk, nv := range nested {
							log.Printf("[DEBUG] Tool %s arg[%s][%s] = %+v (type: %T)", toolCopy.ID, k, nk, nv, nv)
						}
					}
				}
			} else {
				log.Printf("[DEBUG] Tool %s arguments are nil", toolCopy.ID)
			}

			// Call your existing tool with error recovery
			var result interface{}
			var err error

			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("[ERROR] Tool %s panicked: %v", toolCopy.ID, r)
						err = fmt.Errorf("tool panicked: %v", r)
					}
				}()

				result, err = toolCopy.Call(input)
			}()

			if err != nil {
				log.Printf("[ERROR] Tool %s failed: %v", toolCopy.ID, err)
				return mcp.NewToolResultError(fmt.Sprintf("%s failed: %v", toolCopy.ID, err)), nil
			}

			// Log the result for debugging
			if result != nil {
				if resultBytes, err := json.MarshalIndent(result, "", "  "); err == nil {
					log.Printf("[DEBUG] Tool %s result:\n%s", toolCopy.ID, string(resultBytes))
				} else {
					log.Printf("[DEBUG] Tool %s result (raw): %+v", toolCopy.ID, result)
				}
			} else {
				log.Printf("[DEBUG] Tool %s returned nil result", toolCopy.ID)
			}

			// Return result based on type
			switch v := result.(type) {
			case nil:
				return mcp.NewToolResultText("null"), nil
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

		log.Printf("[DEBUG] Tool %s registered successfully", tool.ID)
	}

	log.Printf("[DEBUG] All %d tools registered with MCP server", len(toolRegistry))
}
