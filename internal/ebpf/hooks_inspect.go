// internal/ebpf/hooks_inspect.go
package ebpf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

type HookInfo struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	AttachedTo string `json:"attached_to"`
	Pinned     bool   `json:"pinned"`
	UID        int    `json:"uid,omitempty"`
}

type HookInspectResult struct {
	Programs        []HookInfo `json:"programs"`
	ExecutionTimeMs int        `json:"execution_time_ms"`
}

// InspectHooks uses bpftool to list attached programs and returns them in structured form.
func InspectHooks() (*HookInspectResult, error) {
	start := time.Now()

	cmd := exec.Command("bpftool", "prog", "show", "--json")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run bpftool: %v\noutput: %s", err, out.String())
	}

	// Parse output
	var raw []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse bpftool JSON: %w", err)
	}

	programs := []HookInfo{}
	for _, p := range raw {
		info := HookInfo{
			ID:         intFromMap(p, "id"),
			Type:       strFromMap(p, "type"),
			Name:       strFromMap(p, "name"),
			Pinned:     hasKey(p, "pinned"),
			AttachedTo: strFromMap(p, "attach_to"),
		}
		programs = append(programs, info)
	}

	duration := time.Since(start)
	return &HookInspectResult{
		Programs:        programs,
		ExecutionTimeMs: int(duration.Milliseconds()),
	}, nil
}

func strFromMap(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func intFromMap(m map[string]interface{}, key string) int {
	if val, ok := m[key]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		}
	}
	return 0
}

func hasKey(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}
