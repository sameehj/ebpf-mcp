package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	OLLAMA_URL = "http://localhost:11434/api/generate"
	MCP_URL    = "http://localhost:8080/mcp"
	MODEL      = "llama3.2:3b-instruct-fp16"
)

type ChatTurn struct {
	User string
	AI   string
}

var SESSION_HISTORY []ChatTurn

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("You 🧠: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading input: %v", err)
		return ""
	}
	return strings.TrimSpace(input)
}

func ask_ollama(prompt string) string {
	fmt.Println("[AI 💡] Thinking...")
	payload := map[string]interface{}{
		"model":  MODEL,
		"prompt": prompt,
		"stream": false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling payload: %v", err)
		return "<no response>"
	}

	req, err := http.NewRequest("POST", OLLAMA_URL, strings.NewReader(string(jsonData)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return "<no response>"
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return "<no response>"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Non-200 response: %d", resp.StatusCode)
		return "<no response>"
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Error decoding response: %v", err)
		return "<no response>"
	}

	if responseText, ok := data["response"].(string); ok {
		return responseText
	}
	log.Printf("Unexpected response: %+v", data)
	return "<no response>"
}

// Simple MCP client that mimics your working test
func connectAndListTools() ([]mcp.Tool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create HTTP transport
	httpTransport, err := transport.NewStreamableHTTP(MCP_URL)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP transport: %v", err)
	}

	// Create client
	c := client.NewClient(httpTransport)

	// Start the client
	if err := c.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start client: %v", err)
	}
	defer c.Close()

	// Initialize the client
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Simple eBPF Chat",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	serverInfo, err := c.Initialize(ctx, initRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize: %v", err)
	}

	// List available tools if the server supports them
	if serverInfo.Capabilities.Tools != nil {
		toolsRequest := mcp.ListToolsRequest{}
		toolsResult, err := c.ListTools(ctx, toolsRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to list tools: %v", err)
		}
		return toolsResult.Tools, nil
	}

	return []mcp.Tool{}, nil
}

// Simple tool execution that mimics your working test
func executeToolOnce(toolName string, args map[string]interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create HTTP transport
	httpTransport, err := transport.NewStreamableHTTP(MCP_URL)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP transport: %v", err)
	}

	// Create client
	c := client.NewClient(httpTransport)

	// Start the client
	if err := c.Start(ctx); err != nil {
		return "", fmt.Errorf("failed to start client: %v", err)
	}
	defer c.Close()

	// Initialize the client
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "Simple eBPF Chat",
		Version: "1.0.0",
	}
	initRequest.Params.Capabilities = mcp.ClientCapabilities{}

	_, err = c.Initialize(ctx, initRequest)
	if err != nil {
		return "", fmt.Errorf("failed to initialize: %v", err)
	}

	// Call the tool
	var callRequest mcp.CallToolRequest
	callRequest.Params.Name = toolName
	if args != nil {
		callRequest.Params.Arguments = args
	}

	result, err := c.CallTool(ctx, callRequest)
	if err != nil {
		return "", fmt.Errorf("failed to call tool: %v", err)
	}

	// Extract text from result - use JSON marshaling to avoid type issues
	if len(result.Content) > 0 {
		// Convert to JSON and back to extract text safely
		resultJSON, err := json.Marshal(result.Content)
		if err != nil {
			return "", fmt.Errorf("failed to marshal result: %v", err)
		}

		var contentArray []map[string]interface{}
		if err := json.Unmarshal(resultJSON, &contentArray); err != nil {
			return "", fmt.Errorf("failed to unmarshal result: %v", err)
		}

		var output strings.Builder
		for _, content := range contentArray {
			if text, ok := content["text"].(string); ok {
				output.WriteString(text)
				output.WriteString("\n")
			}
		}
		return output.String(), nil
	}

	return "Tool executed successfully (no output)", nil
}

func buildContextualPrompt(userInput string, availableTools []mcp.Tool) string {
	var promptBuilder strings.Builder

	promptBuilder.WriteString("You are an AI assistant with access to eBPF monitoring tools. You can help analyze system performance, kernel behavior, and network activity.\n")

	if len(availableTools) > 0 {
		promptBuilder.WriteString("\nAvailable eBPF tools:\n")
		for _, tool := range availableTools {
			promptBuilder.WriteString(fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description))
		}
		promptBuilder.WriteString("\nTo use a tool, clearly mention its name in your response.\n")
	}

	promptBuilder.WriteString("\nConversation history:\n")

	// Include recent conversation history (last 3 turns)
	start := 0
	if len(SESSION_HISTORY) > 3 {
		start = len(SESSION_HISTORY) - 3
	}

	for i := start; i < len(SESSION_HISTORY); i++ {
		turn := SESSION_HISTORY[i]
		promptBuilder.WriteString(fmt.Sprintf("Human: %s\nAssistant: %s\n", turn.User, turn.AI))
	}

	promptBuilder.WriteString(fmt.Sprintf("Human: %s\nAssistant:", userInput))

	return promptBuilder.String()
}

func detectAndCallTools(llmResponse string, userInput string, availableTools []mcp.Tool) string {
	var toolOutput strings.Builder
	lowerResponse := strings.ToLower(llmResponse + " " + userInput)

	// eBPF-specific tool detection
	toolMapping := map[string]string{
		"hooks_inspect": "hooks",
		"info":          "info",
		"map_dump":      "map",
		"trace_errors":  "trace",
	}

	for _, tool := range availableTools {
		keyword := toolMapping[tool.Name]
		if keyword != "" && (strings.Contains(lowerResponse, keyword) || strings.Contains(lowerResponse, tool.Name)) {
			fmt.Printf("🔬 Calling eBPF tool: %s\n", tool.Name)

			var args map[string]interface{}

			// Handle map_dump with simple map name extraction
			if tool.Name == "map_dump" {
				words := strings.Fields(userInput + " " + llmResponse)
				for i, word := range words {
					if strings.Contains(strings.ToLower(word), "map") && i+1 < len(words) {
						args = map[string]interface{}{"map_name": words[i+1]}
						break
					}
				}
			}

			result, err := executeToolOnce(tool.Name, args)
			if err != nil {
				toolOutput.WriteString(fmt.Sprintf("\n🔬 Error calling eBPF tool '%s': %s", tool.Name, err.Error()))
			} else {
				toolOutput.WriteString(fmt.Sprintf("\n🔬 eBPF Tool '%s' output:\n%s", tool.Name, result))
			}
			break // Only call one tool per response
		}
	}

	return toolOutput.String()
}

func displayToolsInfo(tools []mcp.Tool) {
	fmt.Println("\n🔬 Available eBPF MCP Tools:")
	for i, tool := range tools {
		fmt.Printf("  %d. %s - %s\n", i+1, tool.Name, tool.Description)
	}
	fmt.Println("\nExample commands:")
	fmt.Println("  • 'show me system info' → uses info tool")
	fmt.Println("  • 'list eBPF programs' → uses hooks_inspect tool")
	fmt.Println("  • 'trace kernel activity' → uses trace_errors tool")
	fmt.Println("  • 'dump map_name' → uses map_dump tool")
}

func chat() {
	fmt.Println("\n🔬 Welcome to the Simple eBPF Chat!")
	fmt.Println("This chat connects to your eBPF MCP server for kernel monitoring and analysis.")
	fmt.Println("Type 'exit' to quit, 'new chat' to reset, or 'list tools' to see available eBPF tools.\n")

	// Test connection and get tools
	fmt.Println("🔗 Connecting to eBPF MCP server...")
	tools, err := connectAndListTools()
	if err != nil {
		fmt.Printf("⚠️  Failed to connect to eBPF MCP server: %v\n", err)
		fmt.Println("   Make sure your eBPF MCP server is running at http://localhost:8080/mcp")
		return
	}

	fmt.Printf("✅ Connected! Found %d eBPF tools available.\n", len(tools))

	for {
		userInput := getUserInput()

		if strings.TrimSpace(strings.ToLower(userInput)) == "exit" {
			fmt.Println("Bye! 👋")
			break
		} else if strings.TrimSpace(strings.ToLower(userInput)) == "new chat" {
			SESSION_HISTORY = nil
			fmt.Println("\n🔄 Starting a new chat session...\n")
			continue
		} else if strings.TrimSpace(strings.ToLower(userInput)) == "list tools" {
			if len(tools) > 0 {
				displayToolsInfo(tools)
			} else {
				fmt.Println("🔬 No eBPF tools available")
			}
			continue
		}

		if userInput == "" {
			continue
		}

		prompt := buildContextualPrompt(userInput, tools)
		llmResponse := ask_ollama(prompt)
		fmt.Printf("AI 🤖: %s\n", llmResponse)

		// Try to detect and call eBPF tools
		toolOutput := detectAndCallTools(llmResponse, userInput, tools)

		fullResponse := llmResponse
		if toolOutput != "" {
			fullResponse += toolOutput
			fmt.Print(toolOutput)
		}

		SESSION_HISTORY = append(SESSION_HISTORY, ChatTurn{
			User: userInput,
			AI:   fullResponse,
		})
	}
}

func main() {
	chat()
}
