package core

import (
	"encoding/json"
	"net/http"

	"github.com/sameehj/ebpf-mcp/internal/prompts"
	"github.com/sameehj/ebpf-mcp/internal/tools"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func HandleMCP(w http.ResponseWriter, r *http.Request) {
	var req types.RPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(types.NewErrorResponse(nil, "Invalid JSON"))
		return
	}

	switch req.Method {
	case "initialize":
		json.NewEncoder(w).Encode(types.NewSuccessResponse(req.ID, map[string]interface{}{
			"protocolVersion": "2025-03-26",
			"capabilities": map[string]interface{}{
				"tools": map[string]bool{
					"listChanged": true,
				},
				"prompts": map[string]bool{
					"listChanged": true,
				},
			},
			"serverInfo": map[string]string{
				"name":    "ebpf-mcp",
				"version": "0.1.0",
			},
		}))
	case "notifications/initialized":
		json.NewEncoder(w).Encode(types.NewSuccessResponse(req.ID, "Initialized"))
	case "tools/list":
		json.NewEncoder(w).Encode(tools.List(req.ID))
	case "tools/call":
		json.NewEncoder(w).Encode(tools.Call(req))
	case "tools/execute":
		json.NewEncoder(w).Encode(tools.Call(req))
	case "prompts/list":
		json.NewEncoder(w).Encode(prompts.List(req.ID))
	case "prompts/get":
		json.NewEncoder(w).Encode(prompts.Get(req))
	default:
		json.NewEncoder(w).Encode(types.NewErrorResponseWithCode(req.ID, -32601, "Method not found"))
	}
}