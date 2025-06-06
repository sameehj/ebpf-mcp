package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterTool(types.Tool{
		ID:          "hooks_inspect",
		Title:       "Inspect eBPF Hooks",
		Description: "List all loaded eBPF programs and their hook points",
		Parameters:  []types.Param{},
		Call: func(input map[string]interface{}) (interface{}, error) {
			return ebpf.InspectHooks()
		},
	})
}
