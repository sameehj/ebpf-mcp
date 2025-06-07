package core

import (
	"encoding/json"
	"net/http"
)

type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"` // Optional for notifications
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *MCPError       `json:"error,omitempty"`
}

type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func HandleMCP(w http.ResponseWriter, r *http.Request) {
	var req MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON-RPC", http.StatusBadRequest)
		return
	}

	switch req.Method {

	case "initialize":
		resp := MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2025-03-26",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{
						"listChanged": true,
					},
				},
				"serverInfo": map[string]string{
					"name":    "ebpf-mcp",
					"version": "0.1.0",
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
		return

	case "notifications/initialized":
		// No response needed for notifications
		w.WriteHeader(http.StatusNoContent)
		return

	case "tools/list":
		handleToolList(w, req.ID)

	case "tools/call":
		handleToolCall(w, req.ID, req.Params)

	default:
		resp := MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}
}
