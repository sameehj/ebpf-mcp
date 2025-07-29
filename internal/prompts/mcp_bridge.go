package prompts

import (
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAllWithMCP registers all prompts with the MCP server
func RegisterAllWithMCP(s *server.MCPServer) {
	log.Printf("[DEBUG] Registering %d prompts", len(promptRegistry))

	for _, p := range promptRegistry {
		mcpPrompt := mcp.NewPrompt(p.ID,
			mcp.WithPromptDescription(p.Description),
		)

		s.AddPrompt(mcpPrompt, p.Handler)
		log.Printf("[DEBUG] Prompt '%s' registered", p.ID)
	}
}
