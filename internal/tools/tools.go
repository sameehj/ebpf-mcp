package tools

import (
	"fmt"
    "github.com/sameehj/ebpf-mcp/pkg/types"
)

var toolRegistry = map[string]types.Tool{
    "map_dump": {
        ID:          "map_dump",
        Title:       "Dump eBPF Map",
        Description: "Returns contents of a BPF map",
        Parameters:  []types.Param{{Name: "map_name", Type: "string", Required: true}},
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
    toolID := req.Params["tool"]
    input := req.Params["input"]

    tool, ok := toolRegistry[toolID.(string)]
    if !ok {
        return types.NewErrorResponse(req.ID, "Tool not found")
    }

	result := map[string]string{
		"message": "Executed tool: " + tool.ID,
		"input": fmt.Sprintf("%v", input), // TODO: Use input in response for now, but we should use it in the tool execution later
	}	

    return types.NewSuccessResponse(req.ID, map[string]interface{}{
        "content": []map[string]string{{"type": "text", "text": result["message"]}},
    })
}
