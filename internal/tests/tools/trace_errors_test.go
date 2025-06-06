//go:build linux
// +build linux

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTraceErrors(t *testing.T) {
	payload := []byte(`{
        "jsonrpc": "2.0",
        "method": "tools/call",
        "params": {
            "tool": "trace_errors"
        },
        "id": 1
    }`)

	resp, err := http.Post("http://localhost:8080/rpc", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("trace_errors tool call failed: %v", err)
	}
	defer resp.Body.Close()

	var decoded map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if decoded["error"] != nil {
		t.Errorf("trace_errors returned error: %+v", decoded["error"])
		return
	}

	result, ok := decoded["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing or invalid 'result' field: %+v", decoded)
	}

	if _, ok := result["traced_events"]; !ok {
		t.Errorf("expected 'traced_events' field in result: %+v", result)
	}
}
