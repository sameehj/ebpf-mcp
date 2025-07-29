package prompts

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

var (
	promptRegistry = map[string]types.Prompt{}
	mu             sync.RWMutex
)

func RegisterPrompt(p types.Prompt) {
	mu.Lock()
	defer mu.Unlock()
	promptRegistry[p.ID] = p
	log.Printf("[DEBUG] Registered prompt: %s", p.ID)
}

func GetPromptByID(id string) (types.Prompt, bool) {
	mu.RLock()
	defer mu.RUnlock()
	p, exists := promptRegistry[id]
	return p, exists
}

func List(id interface{}) types.RPCResponse {
	mu.RLock()
	defer mu.RUnlock()

	prompts := make([]types.PromptMetadata, 0, len(promptRegistry))
	for _, p := range promptRegistry {
		prompts = append(prompts, p.Metadata())
	}
	return types.NewSuccessResponse(id, map[string]interface{}{"prompts": prompts})
}

// Get handles the prompts/get method
func Get(req types.RPCRequest) types.RPCResponse {
	params := req.Params
	if params == nil {
		return types.NewErrorResponse(req.ID, "Missing parameters")
	}

	promptName, ok := params["name"].(string)
	if !ok {
		return types.NewErrorResponse(req.ID, "Missing or invalid prompt name")
	}

	// Get the prompt from registry
	prompt, exists := GetPromptByID(promptName)
	if !exists {
		return types.NewErrorResponseWithCode(req.ID, -32602, "Prompt not found")
	}

	// Parse arguments if provided and convert to map[string]string
	arguments := make(map[string]string)
	if args, ok := params["arguments"].(map[string]interface{}); ok {
		for key, value := range args {
			// Convert each argument value to string
			if strValue, ok := value.(string); ok {
				arguments[key] = strValue
			} else {
				// Convert non-string values to string representation
				arguments[key] = fmt.Sprintf("%v", value)
			}
		}
	}

	// Create the MCP request
	mcpReq := mcp.GetPromptRequest{
		Params: mcp.GetPromptParams{
			Name:      promptName,
			Arguments: arguments,
		},
	}

	// Execute the prompt handler
	ctx := context.Background()
	result, err := prompt.Handler(ctx, mcpReq)
	if err != nil {
		return types.NewErrorResponse(req.ID, err.Error())
	}

	// Convert the result to the expected format
	response := map[string]interface{}{
		"description": result.Description,
		"messages":    convertPromptMessages(result.Messages),
	}

	return types.NewSuccessResponse(req.ID, response)
}

// Helper function to convert MCP prompt messages to the expected format
func convertPromptMessages(messages []mcp.PromptMessage) []map[string]interface{} {
	result := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		result[i] = map[string]interface{}{
			"role":    string(msg.Role),
			"content": convertContent(msg.Content),
		}
	}
	return result
}

// Helper function to convert content
func convertContent(content mcp.Content) interface{} {
	switch c := content.(type) {
	case *mcp.TextContent:
		return map[string]interface{}{
			"type": "text",
			"text": c.Text,
		}
	default:
		// Handle other content types as needed
		return map[string]interface{}{
			"type": "text",
			"text": "",
		}
	}
}
