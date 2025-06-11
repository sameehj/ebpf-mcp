// internal/tools/deploy.go
package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterTool(types.Tool{
		ID:          "deploy",
		Title:       "Deploy eBPF Program",
		Description: "Loads and optionally attaches an eBPF .o file from disk or URL.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"program_path": map[string]interface{}{
					"type":        "string",
					"description": "Local path or URL to the eBPF ELF object",
				},
			},
			"required": []string{"program_path"},
		},
		Annotations: map[string]interface{}{
			"title":          "Deploy BPF",
			"idempotentHint": false,
			"readOnlyHint":   false,
		},
		Call: func(input map[string]interface{}) (interface{}, error) {
			args := ebpf.DeployArgs{
				ProgramPath: input["program_path"].(string),
			}
			return ebpf.DeployProgram(args)
		},
	})
}
