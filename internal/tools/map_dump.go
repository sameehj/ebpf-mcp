package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterTool(types.Tool{
		ID:          "map_dump",
		Title:       "Dump eBPF Map",
		Description: "Returns contents of a BPF map",
		Parameters: []types.Param{
			{Name: "map_name", Type: "string", Required: true},
			{Name: "max_entries", Type: "int", Required: false},
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
