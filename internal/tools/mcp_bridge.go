// internal/tools/mcp_bridge.go - FIXED VERSION
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

// RegisterAllWithMCP converts your existing tools to MCP format using the proper schema API
func RegisterAllWithMCP(s *server.MCPServer) {
	mu.RLock()
	defer mu.RUnlock()

	log.Printf("[DEBUG] Registering %d tools with MCP server", len(toolRegistry))

	for _, tool := range toolRegistry {
		log.Printf("[DEBUG] Converting tool: %s", tool.ID)

		// Convert your tool schema to MCP tool
		mcpTool, err := convertToMCPTool(tool)
		if err != nil {
			log.Printf("[ERROR] Failed to convert tool %s: %v", tool.ID, err)
			continue
		}

		// Capture tool for the handler
		toolCopy := tool

		// Register with MCP server - note: AddTool expects mcp.Tool, not *mcp.Tool
		s.AddTool(mcpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleToolCall(toolCopy, request)
		})

		log.Printf("[DEBUG] Tool %s registered successfully", tool.ID)
	}

	log.Printf("[DEBUG] All %d tools registered with MCP server", len(toolRegistry))
}

// convertToMCPTool converts your tool definition to MCP format
func convertToMCPTool(tool types.Tool) (mcp.Tool, error) {
	// Start with basic tool options
	options := []mcp.ToolOption{
		mcp.WithDescription(tool.Description),
	}

	// Convert input schema if it exists
	if tool.InputSchema != nil {
		schemaOptions, err := convertSchemaToMCPOptions(tool.InputSchema)
		if err != nil {
			return mcp.Tool{}, fmt.Errorf("failed to convert schema for tool %s: %v", tool.ID, err)
		}
		options = append(options, schemaOptions...)
	}

	// Return mcp.Tool (not *mcp.Tool)
	return mcp.NewTool(tool.ID, options...), nil
}

// convertSchemaToMCPOptions converts a JSON schema to MCP tool options
func convertSchemaToMCPOptions(schema map[string]interface{}) ([]mcp.ToolOption, error) {
	var options []mcp.ToolOption

	// Get properties from the schema
	properties, ok := schema["properties"].(map[string]interface{})
	if !ok {
		log.Printf("[DEBUG] No properties found in schema")
		return options, nil
	}

	// Get required fields
	requiredFields := make(map[string]bool)
	if required, ok := schema["required"].([]interface{}); ok {
		for _, field := range required {
			if fieldName, ok := field.(string); ok {
				requiredFields[fieldName] = true
			}
		}
	}

	// Convert each property
	for propName, propDef := range properties {
		propMap, ok := propDef.(map[string]interface{})
		if !ok {
			continue
		}

		propOptions, err := convertPropertyToMCPOptions(propName, propMap, requiredFields[propName])
		if err != nil {
			log.Printf("[WARN] Failed to convert property %s: %v", propName, err)
			continue
		}

		options = append(options, propOptions...)
	}

	return options, nil
}

// convertPropertyToMCPOptions converts a single property to MCP options
func convertPropertyToMCPOptions(propName string, propDef map[string]interface{}, isRequired bool) ([]mcp.ToolOption, error) {
	propType, ok := propDef["type"].(string)
	if !ok {
		return nil, fmt.Errorf("property %s has no type", propName)
	}

	var baseOptions []mcp.PropertyOption

	// Add description if available
	if desc, ok := propDef["description"].(string); ok {
		baseOptions = append(baseOptions, mcp.Description(desc))
	}

	// Add required if needed
	if isRequired {
		baseOptions = append(baseOptions, mcp.Required())
	}

	// Add enum if available
	if enumValues, ok := propDef["enum"].([]interface{}); ok {
		enumStrings := make([]string, 0, len(enumValues))
		for _, val := range enumValues {
			if str, ok := val.(string); ok {
				enumStrings = append(enumStrings, str)
			}
		}
		if len(enumStrings) > 0 {
			baseOptions = append(baseOptions, mcp.Enum(enumStrings...))
		}
	}

	// Handle different property types
	switch propType {
	case "string":
		// Add default value if available
		if defaultVal, ok := propDef["default"].(string); ok {
			baseOptions = append(baseOptions, mcp.DefaultString(defaultVal))
		}
		return []mcp.ToolOption{mcp.WithString(propName, baseOptions...)}, nil

	case "integer", "number":
		// Add default value if available
		if defaultVal, ok := propDef["default"].(float64); ok {
			baseOptions = append(baseOptions, mcp.DefaultNumber(defaultVal))
		}
		return []mcp.ToolOption{mcp.WithNumber(propName, baseOptions...)}, nil

	case "boolean":
		// Add default value if available
		if defaultVal, ok := propDef["default"].(bool); ok {
			baseOptions = append(baseOptions, mcp.DefaultBool(defaultVal))
		}
		return []mcp.ToolOption{mcp.WithBoolean(propName, baseOptions...)}, nil

	case "array":
		// For arrays, we'll use a simpler approach since WithArrayItems doesn't exist
		// Just create a basic array with description
		return []mcp.ToolOption{mcp.WithArray(propName, baseOptions...)}, nil

	case "object":
		// For nested objects, convert recursively
		objectOptions := make([]mcp.PropertyOption, len(baseOptions))
		copy(objectOptions, baseOptions)

		if nestedProps, ok := propDef["properties"].(map[string]interface{}); ok {
			// Get required fields for this object
			nestedRequired := make(map[string]bool)
			if required, ok := propDef["required"].([]interface{}); ok {
				for _, field := range required {
					if fieldName, ok := field.(string); ok {
						nestedRequired[fieldName] = true
					}
				}
			}

			// Convert nested properties - we need to create individual With* calls
			// But since we can't easily nest them in the current API, we'll create a flattened object
			// This is a limitation of the current conversion approach
			log.Printf("[DEBUG] Object property %s has %d nested properties (flattening not fully supported)", propName, len(nestedProps))
		}

		return []mcp.ToolOption{mcp.WithObject(propName, objectOptions...)}, nil

	default:
		log.Printf("[WARN] Unknown property type %s for property %s, treating as string", propType, propName)
		return []mcp.ToolOption{mcp.WithString(propName, baseOptions...)}, nil
	}
}

// handleToolCall handles the actual tool execution
func handleToolCall(tool types.Tool, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	log.Printf("[DEBUG] ==> Tool %s called", tool.ID)

	// Get all arguments as a map
	input := request.GetArguments()

	// Enhanced debug logging
	log.Printf("[DEBUG] Tool %s raw arguments: %+v", tool.ID, input)
	log.Printf("[DEBUG] Tool %s arguments type: %T", tool.ID, input)

	if input != nil {
		log.Printf("[DEBUG] Tool %s arguments length: %d", tool.ID, len(input))
		for k, v := range input {
			log.Printf("[DEBUG] Tool %s arg[%s] = %+v (type: %T)", tool.ID, k, v, v)

			// Deep inspection for nested objects
			if nested, ok := v.(map[string]interface{}); ok {
				log.Printf("[DEBUG] Tool %s arg[%s] is nested object with %d keys", tool.ID, k, len(nested))
				for nk, nv := range nested {
					log.Printf("[DEBUG] Tool %s arg[%s][%s] = %+v (type: %T)", tool.ID, k, nk, nv, nv)
				}
			}
		}
	} else {
		log.Printf("[DEBUG] Tool %s arguments are nil", tool.ID)
	}

	// Call your existing tool with error recovery
	var result interface{}
	var err error

	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[ERROR] Tool %s panicked: %v", tool.ID, r)
				err = fmt.Errorf("tool panicked: %v", r)
			}
		}()

		result, err = tool.Call(input)
	}()

	if err != nil {
		log.Printf("[ERROR] Tool %s failed: %v", tool.ID, err)
		return mcp.NewToolResultError(fmt.Sprintf("%s failed: %v", tool.ID, err)), nil
	}

	// Log the result for debugging
	if result != nil {
		if resultBytes, err := json.MarshalIndent(result, "", "  "); err == nil {
			log.Printf("[DEBUG] Tool %s result:\n%s", tool.ID, string(resultBytes))
		} else {
			log.Printf("[DEBUG] Tool %s result (raw): %+v", tool.ID, result)
		}
	} else {
		log.Printf("[DEBUG] Tool %s returned nil result", tool.ID)
	}

	// Return result based on type
	switch v := result.(type) {
	case nil:
		return mcp.NewToolResultText("null"), nil
	case string:
		return mcp.NewToolResultText(v), nil
	case map[string]interface{}, []interface{}:
		// Convert complex types to JSON string
		if jsonBytes, err := json.MarshalIndent(v, "", "  "); err == nil {
			return mcp.NewToolResultText(string(jsonBytes)), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%v", v)), nil
	default:
		return mcp.NewToolResultText(fmt.Sprintf("%v", v)), nil
	}
}
