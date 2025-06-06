//go:build linux
// +build linux

// internal/tests/tools/map_dump_test.go
package tests

import (
	"encoding/json"
	"os/exec"
	"testing"

	"bytes"
	"net/http"
)

func TestMapDump(t *testing.T) {
	t.Run("creates and dumps a BPF map", func(t *testing.T) {
		// Step 1: Create a pinned BPF map
		cmd := exec.Command("bpftool", "map", "create", "/sys/fs/bpf/test_map",
			"type", "hash", "key", "4", "value", "4", "entries", "128", "name", "test_map")
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to create map: %v", err)
		}

		// Step 2: Insert a key/value pair
		cmd = exec.Command("bpftool", "map", "update", "name", "test_map",
			"key", "0x01", "0x00", "0x00", "0x00", "value", "0x2a", "0x00", "0x00", "0x00")
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to update map: %v", err)
		}

		// Step 3: Call map_dump tool
		payload := []byte(`{
			"jsonrpc": "2.0",
			"method": "tools/call",
			"params": {
				"tool": "map_dump",
				"input": {
					"map_name": "test_map",
					"max_entries": 10
				}
			},
			"id": 1
		}`)

		resp, err := http.Post("http://localhost:8080/rpc", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("map_dump tool call failed: %v", err)
		}
		defer resp.Body.Close()

		var decoded map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
			t.Fatalf("invalid JSON response: %v", err)
		}

		result, ok := decoded["result"].(map[string]interface{})
		if !ok {
			t.Fatalf("missing or invalid 'result' field: %+v", decoded)
		}

		entries, ok := result["entries"].([]interface{})
		if !ok {
			t.Fatalf("missing or invalid 'entries' field: %+v", result)
		}
		if len(entries) == 0 {
			t.Fatalf("expected at least 1 entry, got 0")
		}

		// Cleanup
		_ = exec.Command("bpftool", "map", "delete", "name", "test_map",
			"key", "0x01", "0x00", "0x00", "0x00").Run()
		_ = exec.Command("rm", "/sys/fs/bpf/test_map").Run()
	})
}
