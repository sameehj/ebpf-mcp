package prompts

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func init() {
	RegisterPrompt(types.Prompt{
		ID:          "hello",
		Description: "Says hello to the user",
		Arguments: []mcp.PromptOption{
			mcp.WithArgument("name",
				mcp.ArgumentDescription("Name of the person to greet"),
				mcp.RequiredArgument()),
		},
		Handler: func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			name := req.Params.Arguments["name"]
			if name == "" {
				return nil, fmt.Errorf("name is required")
			}

			message := fmt.Sprintf("Hello, %s! 👋", name)
			return mcp.NewGetPromptResult("Say hello to the user", []mcp.PromptMessage{
				mcp.NewPromptMessage(
					mcp.RoleUser,
					mcp.NewTextContent("Say hello"),
				),
				mcp.NewPromptMessage(mcp.RoleAssistant,
					mcp.NewTextContent(message)),
			}), nil
		},
	})
}
