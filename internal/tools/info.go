package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterTool(types.Tool{
		ID:          "info",
		Title:       "System Info",
		Description: "Returns kernel and environment info for eBPF",
		Parameters:  []types.Param{},
		Call: func(input map[string]interface{}) (interface{}, error) {
			return ebpf.InspectSystemInfo()
		},
	})
}
