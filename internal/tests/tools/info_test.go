//go:build linux
// +build linux

// internal/tests/tools/info_test.go
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestInfoTool(t *testing.T) {
	payload := []byte(`{
		"jsonrpc": "2.0",
		"method": "tools/call",
		"params": {
			"tool": "info"
		},
		"id": 2
	}`)

	resp, err := http.Post("http://localhost:8080/rpc", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("info tool call failed: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	info, ok := result["result"].(map[string]interface{})
	if !ok {
		t.Fatalf("missing or invalid 'result' field: %+v", result)
	}

	requiredFields := []string{"kernel_version", "btf_enabled", "go_arch", "go_os", "sys_fs_bpf"}
	for _, field := range requiredFields {
		if _, exists := info[field]; !exists {
			t.Errorf("missing expected field '%s' in info: %+v", field, info)
		}
	}
}
