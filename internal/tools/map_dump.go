package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterTool(types.Tool{
		ID:          "map_dump",
		Title:       "Dump eBPF Map",
		Description: "Returns key-value contents of a pinned BPF map by name.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"map_name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the pinned BPF map under /sys/fs/bpf/",
				},
				"max_entries": map[string]interface{}{
					"type":    "integer",
					"default": 1000,
				},
			},
			"required": []string{"map_name"},
		},
		Annotations: map[string]interface{}{
			"title":          "Dump Map",
			"readOnlyHint":   true,
			"idempotentHint": true,
			"openWorldHint":  false,
		},
		Call: func(input map[string]interface{}) (interface{}, error) {
			mapName, _ := input["map_name"].(string)
			maxEntries := 1000
			if v, ok := input["max_entries"].(float64); ok {
				maxEntries = int(v)
			}
			return ebpf.DumpPinnedMap("/sys/fs/bpf/"+mapName, maxEntries)
		},
	})
}
