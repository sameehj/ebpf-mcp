// internal/ebpf/inspect_state.go
package ebpf

import (
	"encoding/json"
	"time"

	"github.com/sameehj/ebpf-mcp/pkg/types"
)

type InspectStateArgs struct {
	Fields  []string               `json:"fields"`
	Filters map[string]interface{} `json:"filters"`
}

type InspectStateResult struct {
	Success     bool              `json:"success"`
	ToolVersion string            `json:"tool_version"`
	Programs    []Program         `json:"programs,omitempty"`
	Maps        []Map             `json:"maps,omitempty"`
	Links       []Link            `json:"links,omitempty"`
	System      *SystemInfo       `json:"system,omitempty"`
	Tools       []InspectToolMeta `json:"tools,omitempty"`
	Error       *ErrorDetail      `json:"error,omitempty"`
}

type Program struct {
	ID           int      `json:"id"`
	FD           int      `json:"fd"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	AttachedTo   []string `json:"attached_to"`
	Instructions int      `json:"instructions"`
	MemoryKB     int      `json:"memory_kb"`
	LoadTime     string   `json:"load_time"`
	PinPath      string   `json:"pin_path"`
}

type Map struct {
	ID             int    `json:"id"`
	FD             int    `json:"fd"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	CurrentEntries int    `json:"current_entries"`
	MaxEntries     int    `json:"max_entries"`
	MemoryKB       int    `json:"memory_kb"`
	PinPath        string `json:"pin_path"`
}

type Link struct {
	ID        int    `json:"id"`
	ProgramID int    `json:"program_id"`
	Type      string `json:"type"`
	Target    string `json:"target"`
	PinPath   string `json:"pin_path"`
}

type SystemInfo struct {
	KernelVersion         string   `json:"kernel_version"`
	BTFEnabled            bool     `json:"btf_enabled"`
	TotalPrograms         int      `json:"total_programs"`
	TotalMaps             int      `json:"total_maps"`
	MemoryUsageKB         int      `json:"memory_usage_kb"`
	SupportedProgramTypes []string `json:"supported_program_types"`
}

type InspectToolMeta struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	InputSchema  string `json:"input_schema"`
	OutputSchema string `json:"output_schema"`
}

func ParseInspectStateArgs(input map[string]interface{}) (*InspectStateArgs, error) {
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var args InspectStateArgs
	if err := json.Unmarshal(bytes, &args); err != nil {
		return nil, err
	}
	if len(args.Fields) == 0 {
		args.Fields = []string{"programs", "maps", "links", "system"}
	}
	return &args, nil
}

func InspectState(args *InspectStateArgs) (*InspectStateResult, error) {
	// TODO: Replace with actual introspection of loaded programs, maps, etc.
	demoTools := []types.Tool{
		{
			ID:          "load_program",
			Description: "Load an eBPF program",
		},
		{
			ID:          "attach_program",
			Description: "Attach a loaded program to a hook",
		},
	}

	var toolList []InspectToolMeta
	for _, t := range demoTools {
		toolList = append(toolList, InspectToolMeta{
			Name:         t.ID,
			Version:      "v1",
			Description:  t.Description,
			InputSchema:  "https://ebpf-mcp.dev/schemas/" + t.ID + "_input.json",
			OutputSchema: "https://ebpf-mcp.dev/schemas/" + t.ID + "_output.json",
		})
	}

	return &InspectStateResult{
		Success:     true,
		ToolVersion: "v1",
		Programs: []Program{{
			ID:           123,
			FD:           3,
			Type:         "XDP",
			Name:         "example",
			AttachedTo:   []string{"eth0"},
			Instructions: 256,
			MemoryKB:     64,
			LoadTime:     time.Now().Format(time.RFC3339),
			PinPath:      "/sys/fs/bpf/example",
		}},
		System: &SystemInfo{
			KernelVersion:         "6.9.0",
			BTFEnabled:            true,
			TotalPrograms:         1,
			TotalMaps:             1,
			MemoryUsageKB:         128,
			SupportedProgramTypes: []string{"XDP", "KPROBE", "TRACEPOINT"},
		},
		Tools: toolList,
	}, nil
}
