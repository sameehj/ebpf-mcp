package ebpf

import (
	"fmt"
	"time"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/asm"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/rlimit"
)

type TraceResult struct {
	Status   string   `json:"status"`
	Traced   int      `json:"traced_events"`
	Warnings []string `json:"warnings,omitempty"`
}

func RunTraceErrors() (*TraceResult, error) {
	if err := rlimit.RemoveMemlock(); err != nil {
		return nil, fmt.Errorf("rlimit error: %v", err)
	}

	progSpec := &ebpf.ProgramSpec{
		Name:    "trace_exit",
		Type:    ebpf.TracePoint,
		License: "GPL",
		Instructions: asm.Instructions{
			asm.Mov.Imm(asm.R0, 0),
			asm.Return(),
		},
	}

	prog, err := ebpf.NewProgram(progSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to load program: %v", err)
	}
	defer prog.Close()

	tp, err := link.Tracepoint("sched", "sched_switch", prog, nil)
	if err != nil {
		return nil, fmt.Errorf("tracepoint attach failed: %v", err)
	}
	defer tp.Close()

	time.Sleep(2 * time.Second)

	return &TraceResult{
		Status: "Tracepoint attached and ran for 2s",
		Traced: 0,
	}, nil
}
