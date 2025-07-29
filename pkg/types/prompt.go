package types

import (
	"context"
	
	"github.com/mark3labs/mcp-go/mcp"
)

type Prompt struct {
	ID          string
	Description string
	Arguments   []mcp.PromptOption
	Handler     func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)
}

func (p Prompt) Metadata() PromptMetadata {
	return PromptMetadata{
		Name:        p.ID,
		Description: p.Description,
	}
}

type PromptMetadata struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
}