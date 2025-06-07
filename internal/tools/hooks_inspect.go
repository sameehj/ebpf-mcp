package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterTool(types.Tool{
		ID:          "hooks_inspect",
		Title:       "Inspect eBPF Hooks",
		Description: "Lists all loaded eBPF programs and their attachment points.",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{}, // no params
		},
		Annotations: map[string]interface{}{
			"title":          "List Active Hooks",
			"readOnlyHint":   true,
			"idempotentHint": true,
			"openWorldHint":  false,
		},
		Call: func(input map[string]interface{}) (interface{}, error) {
			return ebpf.InspectHooks()
		},
	})
}
