package prompts

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sameehj/ebpf-mcp/internal/tools"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterPrompt(types.Prompt{
		ID:          "show_system_info",
		Description: "Shows system and kernel info for eBPF.",
		Arguments:   []mcp.PromptOption{},
		Handler: func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			return tools.CallTool("info", map[string]interface{}{})
		},
	})
}
