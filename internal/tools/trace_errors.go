package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterTool(types.Tool{
		ID:          "trace_errors",
		Title:       "Trace Syscall Errors",
		Description: "Attaches to sched_switch to simulate an eBPF trace for 2s",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
		Annotations: map[string]interface{}{
			"title":          "Trace Errors (sched_switch)",
			"readOnlyHint":   true,
			"idempotentHint": true,
			"openWorldHint":  false,
		},
		Call: func(input map[string]interface{}) (interface{}, error) {
			return ebpf.RunTraceErrors()
		},
	})
}
