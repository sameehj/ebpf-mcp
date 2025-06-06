package tools

import (
	"fmt"
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
