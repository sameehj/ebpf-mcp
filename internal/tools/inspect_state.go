// internal/tools/inspect_state.go
package tools

import (
	"runtime"

	"github.com/sameehj/ebpf-mcp/pkg/types"
)

type InspectStateInput struct {
	Fields  []string               `json:"fields"`
	Filters map[string]interface{} `json:"filters"`
}

type InspectStateOutput struct {
	Success     bool               `json:"success"`
	ToolVersion string             `json:"tool_version"`
	Programs    []map[string]any   `json:"programs,omitempty"`
	Maps        []map[string]any   `json:"maps,omitempty"`
	Links       []map[string]any   `json:"links,omitempty"`
	System      map[string]any     `json:"system,omitempty"`
	Error       *types.ErrorDetail `json:"error,omitempty"`
}

func InspectState(input map[string]interface{}) (interface{}, error) {
	var args InspectStateInput
	if err := types.StrictUnmarshal(input, &args); err != nil {
		return nil, err
	}

	if len(args.Fields) == 0 {
		args.Fields = []string{"programs", "maps", "links", "system"}
	}

	out := InspectStateOutput{
		Success:     true,
		ToolVersion: "v1",
	}

	for _, field := range args.Fields {
		switch field {
		case "system":
			out.System = map[string]any{
				"kernel_version":          runtime.GOARCH, // TODO: replace with actual version
				"btf_enabled":             true,
				"total_programs":          0,
				"total_maps":              0,
				"memory_usage_kb":         0,
				"supported_program_types": []string{"XDP", "KPROBE", "TRACEPOINT", "CGROUP_SKB"},
			}
		case "programs":
			out.Programs = append(out.Programs, map[string]any{
				"id":   123,
				"type": "XDP",
				"name": "demo_prog",
			})
		case "maps":
			out.Maps = append(out.Maps, map[string]any{
				"id":   1,
				"name": "demo_map",
				"type": "hash",
			})
		case "links":
			out.Links = append(out.Links, map[string]any{
				"id":         7,
				"program_id": 123,
				"type":       "xdp",
				"target":     "eth0",
			})
		}
	}

	return out, nil
}

func init() {
	RegisterTool(types.Tool{
		ID:          "inspect_state",
		Title:       "Inspect eBPF State",
		Description: "Returns current state of loaded eBPF programs, maps, links, and tools.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"fields": map[string]any{
					"type": "array",
					"items": map[string]any{
						"enum": []string{"programs", "maps", "links", "system"},
					},
					"default": []string{"programs", "maps", "links", "system"},
				},
				"filters": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"program_type": map[string]any{"type": "string"},
						"interface":    map[string]any{"type": "string"},
						"name_pattern": map[string]any{"type": "string"},
						"program_id":   map[string]any{"type": "integer"},
						"map_id":       map[string]any{"type": "integer"},
						"pinned_only":  map[string]any{"type": "boolean"},
					},
				},
			},
		},
		OutputSchema: map[string]any{
			"type":     "object",
			"required": []string{"success", "tool_version"},
			"properties": map[string]any{
				"success":      map[string]any{"type": "boolean"},
				"tool_version": map[string]any{"type": "string"},
				"programs":     map[string]any{"type": "array", "items": map[string]any{"type": "object"}},
				"maps":         map[string]any{"type": "array", "items": map[string]any{"type": "object"}},
				"links":        map[string]any{"type": "array", "items": map[string]any{"type": "object"}},
				"system":       map[string]any{"type": "object"},
				"error": map[string]any{
					"$ref": "https://ebpf-mcp.dev/schemas/error.schema.json#/definitions/Error",
				},
			},
		},
		Annotations: map[string]any{
			"title":          "Inspect State",
			"idempotentHint": true,
			"readOnlyHint":   true,
		},
		Call: InspectState,
	})
}
