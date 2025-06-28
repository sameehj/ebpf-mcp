// internal/tools/attach_program.go
package tools

import (
	"fmt"

	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func AttachProgramTool(input map[string]interface{}) (interface{}, error) {
	// Add debug logging
	fmt.Printf("[DEBUG] AttachProgram input: %+v\n", input)

	// Handle nil or empty input
	if input == nil {
		return &ebpf.AttachProgramResult{
			Success:     false,
			ToolVersion: "v1",
			Message:     "input is nil",
		}, fmt.Errorf("input is nil")
	}

	// Safe parsing with nil checks
	args, err := parseAttachProgramInputSafely(input)
	if err != nil {
		return &ebpf.AttachProgramResult{
			Success:     false,
			ToolVersion: "v1",
			Message:     fmt.Sprintf("parsing error: %v", err),
		}, err
	}

	// Call the actual eBPF attachment function
	return ebpf.AttachProgram(args)
}

func parseAttachProgramInputSafely(input map[string]interface{}) (*ebpf.AttachProgramArgs, error) {
	var args ebpf.AttachProgramArgs

	// Parse program_id
	if programIDRaw, exists := input["program_id"]; exists && programIDRaw != nil {
		if programID, ok := programIDRaw.(float64); ok {
			args.ProgramID = int(programID)
		} else {
			return nil, fmt.Errorf("program_id must be a number")
		}
	} else {
		return nil, fmt.Errorf("program_id is required")
	}

	// Parse attach_type
	if attachTypeRaw, exists := input["attach_type"]; exists && attachTypeRaw != nil {
		if attachType, ok := attachTypeRaw.(string); ok {
			args.AttachType = attachType
		} else {
			return nil, fmt.Errorf("attach_type must be a string")
		}
	} else {
		return nil, fmt.Errorf("attach_type is required")
	}

	// Parse target (optional for some attach types)
	if targetRaw, exists := input["target"]; exists && targetRaw != nil {
		if target, ok := targetRaw.(string); ok {
			args.Target = target
		}
	}

	// Parse pin_path (optional)
	if pinPathRaw, exists := input["pin_path"]; exists && pinPathRaw != nil {
		if pinPath, ok := pinPathRaw.(string); ok {
			args.PinPath = pinPath
		}
	}

	// Parse options (optional)
	if optionsRaw, exists := input["options"]; exists && optionsRaw != nil {
		if options, ok := optionsRaw.(map[string]interface{}); ok {
			args.Options = options
		}
	}

	return &args, nil
}

func init() {
	RegisterTool(types.Tool{
		ID:          "attach_program",
		Title:       "Attach eBPF Program",
		Description: "Attaches a loaded eBPF program to a kernel hook point.",
		InputSchema: map[string]interface{}{
			"type":     "object",
			"required": []string{"program_id", "attach_type"},
			"properties": map[string]interface{}{
				"program_id": map[string]interface{}{
					"type":        "integer",
					"description": "ID of the loaded eBPF program",
				},
				"attach_type": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"xdp", "kprobe", "kretprobe", "tracepoint", "cgroup"},
					"description": "Type of attachment point",
				},
				"target": map[string]interface{}{
					"type":        "string",
					"description": "Target interface, function, or cgroup path",
				},
				"pin_path": map[string]interface{}{
					"type":        "string",
					"description": "Path to pin the link object",
				},
				"options": map[string]interface{}{
					"type":        "object",
					"description": "Additional attachment options",
					"properties": map[string]interface{}{
						"flags": map[string]interface{}{
							"type": "integer",
						},
						"priority": map[string]interface{}{
							"type": "integer",
						},
					},
				},
			},
		},
		OutputSchema: map[string]interface{}{
			"type":     "object",
			"required": []string{"success", "tool_version"},
			"properties": map[string]interface{}{
				"success":      map[string]interface{}{"type": "boolean"},
				"tool_version": map[string]interface{}{"type": "string"},
				"link_id":      map[string]interface{}{"type": "integer"},
				"pin_path":     map[string]interface{}{"type": "string"},
				"message":      map[string]interface{}{"type": "string"},
				"error": map[string]interface{}{
					"$ref": "https://ebpf-mcp.dev/schemas/error.schema.json#/definitions/Error",
				},
			},
		},
		Annotations: map[string]interface{}{
			"title":          "Attach Program",
			"idempotentHint": false,
			"readOnlyHint":   false,
		},
		Call: AttachProgramTool,
	})
}
