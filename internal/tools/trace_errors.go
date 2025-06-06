package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterTool(types.Tool{
		ID:          "trace_errors",
		Title:       "Trace Syscall Errors",
		Description: "Attaches to sched_switch to simulate eBPF trace for 2s",
		Parameters:  []types.Param{},
		Call: func(input map[string]interface{}) (interface{}, error) {
			return ebpf.RunTraceErrors()
		},
	})
}
