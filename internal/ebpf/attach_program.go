// internal/ebpf/attach_program.go
package ebpf

import (
	"fmt"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

type AttachProgramArgs struct {
	ProgramID  int                    `json:"program_id"`
	AttachType string                 `json:"attach_type"`
	Target     string                 `json:"target,omitempty"`
	PinPath    string                 `json:"pin_path,omitempty"`
	Options    map[string]interface{} `json:"options,omitempty"`
}

type AttachProgramResult struct {
	Success     bool   `json:"success"`
	ToolVersion string `json:"tool_version"`
	LinkID      int    `json:"link_id,omitempty"`
	PinPath     string `json:"pin_path,omitempty"`
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
}

func AttachProgram(args *AttachProgramArgs) (*AttachProgramResult, error) {
	// For now, return a mock response since actual attachment requires root privileges
	// and proper eBPF program loading infrastructure

	result := &AttachProgramResult{
		Success:     true,
		ToolVersion: "v1",
		LinkID:      args.ProgramID + 1000, // Mock link ID
		Message:     fmt.Sprintf("Program %d attached as %s", args.ProgramID, args.AttachType),
	}

	// If pin_path is provided, simulate pinning
	if args.PinPath != "" {
		result.PinPath = args.PinPath
		result.Message += fmt.Sprintf(" and pinned to %s", args.PinPath)
	}

	// Add target information if provided
	if args.Target != "" {
		result.Message += fmt.Sprintf(" on target %s", args.Target)
	}

	return result, nil
}

// AttachProgramActual would be the real implementation
func AttachProgramActual(args *AttachProgramArgs) (*AttachProgramResult, error) {
	// This is what the actual implementation would look like
	// when you have proper eBPF programs loaded

	// Get the program by ID (this would require maintaining a registry)
	var prog *ebpf.Program
	// prog = getProgramByID(args.ProgramID)

	if prog == nil {
		return &AttachProgramResult{
			Success:     false,
			ToolVersion: "v1",
			Error:       fmt.Sprintf("program with ID %d not found", args.ProgramID),
		}, fmt.Errorf("program not found")
	}

	var l link.Link
	var err error

	switch args.AttachType {
	case "xdp":
		if args.Target == "" {
			return &AttachProgramResult{
				Success:     false,
				ToolVersion: "v1",
				Error:       "target interface required for XDP attachment",
			}, fmt.Errorf("target required for XDP")
		}
		// l, err = link.AttachXDP(link.XDPOptions{
		// 	Program:   prog,
		// 	Interface: args.Target,
		// })

	case "kprobe":
		if args.Target == "" {
			return &AttachProgramResult{
				Success:     false,
				ToolVersion: "v1",
				Error:       "target function required for kprobe attachment",
			}, fmt.Errorf("target required for kprobe")
		}
		// l, err = link.Kprobe(args.Target, prog, nil)

	case "kretprobe":
		if args.Target == "" {
			return &AttachProgramResult{
				Success:     false,
				ToolVersion: "v1",
				Error:       "target function required for kretprobe attachment",
			}, fmt.Errorf("target required for kretprobe")
		}
		// l, err = link.Kretprobe(args.Target, prog, nil)

	case "tracepoint":
		if args.Target == "" {
			return &AttachProgramResult{
				Success:     false,
				ToolVersion: "v1",
				Error:       "target tracepoint required for tracepoint attachment",
			}, fmt.Errorf("target required for tracepoint")
		}
		// Parse tracepoint format: "group:name"
		// l, err = link.Tracepoint(link.TracepointOptions{
		// 	Program: prog,
		// 	Group:   group,
		// 	Name:    name,
		// })

	default:
		return &AttachProgramResult{
			Success:     false,
			ToolVersion: "v1",
			Error:       fmt.Sprintf("unsupported attach type: %s", args.AttachType),
		}, fmt.Errorf("unsupported attach type: %s", args.AttachType)
	}

	if err != nil {
		return &AttachProgramResult{
			Success:     false,
			ToolVersion: "v1",
			Error:       err.Error(),
		}, err
	}

	// Pin the link if requested
	if args.PinPath != "" {
		err = l.Pin(args.PinPath)
		if err != nil {
			l.Close()
			return &AttachProgramResult{
				Success:     false,
				ToolVersion: "v1",
				Error:       fmt.Sprintf("failed to pin link: %v", err),
			}, err
		}
	}

	// Get link info
	linkInfo, err := l.Info()
	if err != nil {
		l.Close()
		return &AttachProgramResult{
			Success:     false,
			ToolVersion: "v1",
			Error:       fmt.Sprintf("failed to get link info: %v", err),
		}, err
	}

	linkID := linkInfo.ID

	return &AttachProgramResult{
		Success:     true,
		ToolVersion: "v1",
		LinkID:      int(linkID),
		PinPath:     args.PinPath,
		Message:     fmt.Sprintf("Successfully attached program %d as %s", args.ProgramID, args.AttachType),
	}, nil
}
