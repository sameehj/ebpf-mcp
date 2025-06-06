package core

import (
    "encoding/json"
    "net/http"
    "ebpf-mcp/pkg/types"
    "ebpf-mcp/internal/tools"
)

func HandleMCP(w http.ResponseWriter, r *http.Request) {
    var req types.RPCRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        json.NewEncoder(w).Encode(types.NewErrorResponse(req.ID, "Invalid JSON"))
        return
    }

    switch req.Method {
    case "tools/list":
        json.NewEncoder(w).Encode(tools.List(req.ID))
    case "tools/call":
        json.NewEncoder(w).Encode(tools.Call(req))
    default:
        json.NewEncoder(w).Encode(types.NewErrorResponse(req.ID, "Method not found"))
    }
}
