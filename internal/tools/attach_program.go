// internal/tools/attach_program.go
package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterTool(types.Tool{
		ID:          "attach_program",
		Title:       "Attach eBPF Program",
		Description: "Attaches a previously loaded eBPF program to a kernel hook (XDP, tracepoint, kprobe).",
		InputSchema: map[string]interface{}{
			"$schema":     "https://json-schema.org/draft/2020-12/schema",
			"$id":         "https://ebpf-mcp.dev/schemas/attach_program_input.json",
			"title":       "Attach Program Input",
			"description": "Schema for attach_program tool input parameters",
			"type":        "object",
			"required":    []string{"program_id", "attachment"},
			"properties": map[string]interface{}{
				"program_id": map[string]interface{}{"type": "integer"},
				"attachment": map[string]interface{}{
					"oneOf": []interface{}{
						map[string]interface{}{
							"type":     "object",
							"required": []string{"type", "params"},
							"properties": map[string]interface{}{
								"type": map[string]interface{}{"const": "xdp"},
								"params": map[string]interface{}{
									"type":     "object",
									"required": []string{"interface"},
									"properties": map[string]interface{}{
										"interface":     map[string]interface{}{"type": "string"},
										"mode":          map[string]interface{}{"enum": []string{"NATIVE", "GENERIC", "OFFLOAD"}, "default": "NATIVE"},
										"link_pin_path": map[string]interface{}{"type": "string"},
									},
								},
							},
						},
						map[string]interface{}{
							"type":     "object",
							"required": []string{"type", "params"},
							"properties": map[string]interface{}{
								"type": map[string]interface{}{"const": "tracepoint"},
								"params": map[string]interface{}{
									"type":     "object",
									"required": []string{"group", "name"},
									"properties": map[string]interface{}{
										"group":         map[string]interface{}{"type": "string"},
										"name":          map[string]interface{}{"type": "string"},
										"link_pin_path": map[string]interface{}{"type": "string"},
									},
								},
							},
						},
						map[string]interface{}{
							"type":     "object",
							"required": []string{"type", "params"},
							"properties": map[string]interface{}{
								"type": map[string]interface{}{"const": "kprobe"},
								"params": map[string]interface{}{
									"type":     "object",
									"required": []string{"function"},
									"properties": map[string]interface{}{
										"function":      map[string]interface{}{"type": "string"},
										"retprobe":      map[string]interface{}{"type": "boolean", "default": false},
										"offset":        map[string]interface{}{"type": "integer", "default": 0},
										"link_pin_path": map[string]interface{}{"type": "string"},
									},
								},
							},
						},
					},
				},
			},
		},
		OutputSchema: map[string]interface{}{
			"$schema":     "https://json-schema.org/draft/2020-12/schema",
			"$id":         "https://ebpf-mcp.dev/schemas/attach_program_output.json",
			"title":       "Attach Program Output",
			"description": "Schema for attach_program tool output",
			"type":        "object",
			"required":    []string{"success", "tool_version"},
			"properties": map[string]interface{}{
				"success":          map[string]interface{}{"type": "boolean"},
				"tool_version":     map[string]interface{}{"type": "string"},
				"link_id":          map[string]interface{}{"type": "integer"},
				"attachment_point": map[string]interface{}{"type": "string"},
				"pin_path":         map[string]interface{}{"type": "string"},
				"error":            map[string]interface{}{"$ref": "error.schema.json#/definitions/Error"},
			},
		},
		Annotations: map[string]interface{}{
			"title":          "Attach Program",
			"idempotentHint": false,
			"readOnlyHint":   false,
		},
		Call: func(input map[string]interface{}) (interface{}, error) {
			args, err := ebpf.ParseAttachProgramArgs(input)
			if err != nil {
				return nil, err
			}
			return ebpf.AttachProgram(args)
		},
	})
}
