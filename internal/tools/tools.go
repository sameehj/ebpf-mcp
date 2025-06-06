// internal/tools/tools.go
package tools

import (
	"fmt"
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

var toolRegistry = map[string]types.Tool{
	"map_dump": {
		ID:          "map_dump",
		Title:       "Dump eBPF Map",
		Description: "Returns contents of a BPF map",
		Parameters: []types.Param{
			{Name: "map_name", Type: "string", Required: true},
			{Name: "max_entries", Type: "int", Required: false},
		},
	},
}

func List(id interface{}) types.RPCResponse {
	tools := make([]types.ToolMetadata, 0)
	for _, t := range toolRegistry {
		tools = append(tools, t.Metadata())
	}
	return types.NewSuccessResponse(id, map[string]interface{}{"tools": tools})
}

func Call(req types.RPCRequest) types.RPCResponse {
	rawToolID, ok := req.Params["tool"]
	if !ok {
		return types.NewErrorResponse(req.ID, "Missing 'tool' field")
	}

	toolID, ok := rawToolID.(string)
	if !ok {
		return types.NewErrorResponse(req.ID, "Invalid tool ID")
	}

	_, exists := toolRegistry[toolID]
	if !exists {
		return types.NewErrorResponse(req.ID, fmt.Sprintf("Tool '%s' not found", toolID))
	}

	input, _ := req.Params["input"].(map[string]interface{})

	switch toolID {
	case "map_dump":
		mapName, _ := input["map_name"].(string)
		maxEntriesFloat, ok := input["max_entries"].(float64)
		if !ok {
			maxEntriesFloat = 1000
		}
		result, err := ebpf.DumpPinnedMap("/sys/fs/bpf/"+mapName, int(maxEntriesFloat))
		if err != nil {
			return types.NewErrorResponse(req.ID, fmt.Sprintf("map_dump failed: %v", err))
		}
		return types.NewSuccessResponse(req.ID, result)

	default:
		return types.NewErrorResponse(req.ID, fmt.Sprintf("Tool '%s' is not yet implemented", toolID))
	}
}
