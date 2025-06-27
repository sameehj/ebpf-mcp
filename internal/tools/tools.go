// internal/tools/tools.go
package tools

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/sameehj/ebpf-mcp/pkg/types"
)

var (
	toolRegistry = map[string]types.Tool{}
	mu           sync.RWMutex
)

func RegisterTool(t types.Tool) {
	mu.Lock()
	defer mu.Unlock()
	toolRegistry[t.ID] = t
	log.Printf("[DEBUG] Registered tool: %s", t.ID)
}

func GetAllTools() map[string]types.Tool {
	mu.RLock()
	defer mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make(map[string]types.Tool)
	for k, v := range toolRegistry {
		result[k] = v
	}
	return result
}

func GetTool(id string) (types.Tool, bool) {
	mu.RLock()
	defer mu.RUnlock()
	tool, exists := toolRegistry[id]
	return tool, exists
}

func ListToolsWithSchemas() []map[string]interface{} {
	mu.RLock()
	defer mu.RUnlock()

	tools := make([]map[string]interface{}, 0, len(toolRegistry))
	for _, t := range toolRegistry {
		toolInfo := map[string]interface{}{
			"name":        t.ID,
			"description": t.Description,
		}

		// Include input schema if available
		if t.InputSchema != nil {
			toolInfo["inputSchema"] = t.InputSchema
		}

		// Include output schema if available
		if t.OutputSchema != nil {
			toolInfo["outputSchema"] = t.OutputSchema
		}

		// Include annotations if available
		if t.Annotations != nil {
			toolInfo["annotations"] = t.Annotations
		}

		tools = append(tools, toolInfo)
	}

	log.Printf("[DEBUG] Listed %d tools with schemas", len(tools))
	return tools
}

// For debugging - print all registered tools with their schemas
func DebugPrintTools() {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("[DEBUG] ==> Registered Tools Summary:")
	for id, tool := range toolRegistry {
		log.Printf("[DEBUG] Tool: %s", id)
		log.Printf("[DEBUG]   Description: %s", tool.Description)

		if tool.InputSchema != nil {
			if schemaBytes, err := json.MarshalIndent(tool.InputSchema, "    ", "  "); err == nil {
				log.Printf("[DEBUG]   Input Schema:\n    %s", string(schemaBytes))
			}
		} else {
			log.Printf("[DEBUG]   Input Schema: none")
		}

		log.Printf("[DEBUG]   ---")
	}
}

func List(id interface{}) types.RPCResponse {
	mu.RLock()
	defer mu.RUnlock()

	tools := make([]types.ToolMetadata, 0, len(toolRegistry))
	for _, t := range toolRegistry {
		tools = append(tools, t.Metadata())
	}
	return types.NewSuccessResponse(id, map[string]interface{}{"tools": tools})
}

func Call(req types.RPCRequest) types.RPCResponse {
	mu.RLock()
	defer mu.RUnlock()

	rawToolID, ok := req.Params["tool"]
	if !ok {
		return types.NewErrorResponse(req.ID, "Missing 'tool' field")
	}

	toolID, ok := rawToolID.(string)
	if !ok {
		return types.NewErrorResponse(req.ID, "Invalid tool ID")
	}

	tool, exists := toolRegistry[toolID]
	if !exists {
		return types.NewErrorResponse(req.ID, fmt.Sprintf("Tool '%s' not found", toolID))
	}

	if tool.Call == nil {
		return types.NewErrorResponse(req.ID, fmt.Sprintf("Tool '%s' has no callable function", toolID))
	}

	input, _ := req.Params["input"].(map[string]interface{})
	result, err := tool.Call(input)
	if err != nil {
		return types.NewErrorResponse(req.ID, fmt.Sprintf("%s failed: %v", toolID, err))
	}
	return types.NewSuccessResponse(req.ID, result)
}
