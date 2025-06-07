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
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		Annotations: map[string]interface{}{
			"title":          "eBPF System Info",
			"readOnlyHint":   true,
			"idempotentHint": true,
			"openWorldHint":  false,
		},
		Call: func(input map[string]interface{}) (interface{}, error) {
			return ebpf.InspectSystemInfo()
		},
	})
}
