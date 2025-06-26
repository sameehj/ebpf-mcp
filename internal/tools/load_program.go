package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterTool(types.Tool{
		ID:          "load_program",
		Title:       "Load eBPF Program",
		Description: "Loads a raw eBPF object from file or base64 blob into the kernel.",
		InputSchema: map[string]interface{}{
			"type":     "object",
			"required": []string{"source", "program_type"},
			"properties": map[string]interface{}{
				"source": map[string]interface{}{
					"type":     "object",
					"required": []string{"type"},
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"enum": []string{"file", "data"},
						},
						"path": map[string]interface{}{
							"type": "string",
						},
						"blob": map[string]interface{}{
							"type":            "string",
							"contentEncoding": "base64",
						},
						"checksum": map[string]interface{}{
							"type":    "string",
							"pattern": "^sha256:[a-f0-9]{64}$",
						},
					},
					"oneOf": []map[string]interface{}{
						{
							"properties": map[string]interface{}{
								"type": map[string]interface{}{"const": "file"},
							},
							"required": []string{"path"},
						},
						{
							"properties": map[string]interface{}{
								"type": map[string]interface{}{"const": "data"},
							},
							"required": []string{"blob"},
						},
					},
				},
				"program_type": map[string]interface{}{
					"enum": []string{"XDP", "KPROBE", "TRACEPOINT", "CGROUP_SKB"},
				},
				"section": map[string]interface{}{
					"type": "string",
				},
				"btf_path": map[string]interface{}{
					"type": "string",
				},
				"constraints": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"max_instructions": map[string]interface{}{"type": "integer"},
						"allowed_helpers": map[string]interface{}{
							"type":  "array",
							"items": map[string]interface{}{"type": "string"},
						},
						"verify_only": map[string]interface{}{"type": "boolean"},
					},
				},
			},
		},
		Annotations: map[string]interface{}{
			"title":          "Load Program",
			"idempotentHint": true,
			"readOnlyHint":   false,
		},
		Call: func(input map[string]interface{}) (interface{}, error) {
			args, err := ebpf.ParseLoadProgramArgs(input)
			if err != nil {
				return nil, err
			}
			return ebpf.LoadProgram(args)
		},
	})
}
