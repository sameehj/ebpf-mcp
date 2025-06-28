// internal/tools/load_program.go
package tools

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func LoadProgramTool(input map[string]interface{}) (interface{}, error) {
	// Enhanced debug logging
	log.Printf("[DEBUG] ==> LoadProgram called")
	log.Printf("[DEBUG] LoadProgram raw input: %+v", input)
	log.Printf("[DEBUG] LoadProgram input type: %T", input)

	if input != nil {
		log.Printf("[DEBUG] LoadProgram input length: %d", len(input))
		for k, v := range input {
			log.Printf("[DEBUG] LoadProgram input[%s] = %+v (type: %T)", k, v, v)
		}
	}

	// Convert input to JSON and back for debugging
	if input != nil {
		jsonBytes, _ := json.MarshalIndent(input, "", "  ")
		log.Printf("[DEBUG] LoadProgram input JSON:\n%s", string(jsonBytes))
	}

	// Handle completely nil input
	if input == nil {
		log.Printf("[ERROR] LoadProgram: input is nil")
		return &ebpf.LoadProgramResult{
			Success:      false,
			ToolVersion:  "1.0.0",
			ErrorMessage: "input is nil",
		}, nil // Return nil error to avoid MCP panic
	}

	// Handle empty input
	if len(input) == 0 {
		log.Printf("[ERROR] LoadProgram: input is empty")
		return &ebpf.LoadProgramResult{
			Success:      false,
			ToolVersion:  "1.0.0",
			ErrorMessage: "input is empty",
		}, nil
	}

	// Try to provide a default minimal configuration if input is malformed
	processedInput := preprocessInput(input)

	// Parse with extensive error handling
	args, err := parseLoadProgramInputRobust(processedInput)
	if err != nil {
		log.Printf("[ERROR] LoadProgram parsing failed: %v", err)
		return &ebpf.LoadProgramResult{
			Success:      false,
			ToolVersion:  "1.0.0",
			ErrorMessage: fmt.Sprintf("parsing error: %v", err),
		}, nil // Return nil error to avoid MCP panic
	}

	log.Printf("[DEBUG] LoadProgram parsed args: %+v", args)

	// Call the actual eBPF loading function
	result, err := ebpf.LoadProgram(args)
	if err != nil {
		log.Printf("[ERROR] LoadProgram execution failed: %v", err)
		// Ensure we return a valid result even on error
		if result == nil {
			result = &ebpf.LoadProgramResult{
				Success:      false,
				ToolVersion:  "1.0.0",
				ErrorMessage: err.Error(),
			}
		}
		return result, nil // Return nil error to avoid MCP panic
	}

	return result, nil
}

func preprocessInput(input map[string]interface{}) map[string]interface{} {
	// Create a copy to avoid modifying the original
	processed := make(map[string]interface{})

	// Copy all non-nil values
	for k, v := range input {
		if v != nil {
			processed[k] = v
		}
	}

	// Provide defaults for missing required fields
	if _, exists := processed["source"]; !exists {
		log.Printf("[WARN] Missing source field, providing default")
		processed["source"] = map[string]interface{}{
			"type": "file",
			"path": "/tmp/kprobe.o", // Default path
		}
	}

	if _, exists := processed["program_type"]; !exists {
		log.Printf("[WARN] Missing program_type field, providing default")
		processed["program_type"] = "KPROBE"
	}

	return processed
}

func parseLoadProgramInputRobust(input map[string]interface{}) (ebpf.LoadProgramArgs, error) {
	var args ebpf.LoadProgramArgs

	log.Printf("[DEBUG] Parsing input: %+v", input)

	// Parse source with extensive validation
	sourceRaw, exists := input["source"]
	if !exists {
		return args, fmt.Errorf("source field is missing")
	}
	if sourceRaw == nil {
		return args, fmt.Errorf("source field is nil")
	}

	source, ok := sourceRaw.(map[string]interface{})
	if !ok {
		log.Printf("[DEBUG] Source type conversion failed. sourceRaw type: %T, value: %+v", sourceRaw, sourceRaw)
		return args, fmt.Errorf("source must be an object, got %T", sourceRaw)
	}

	log.Printf("[DEBUG] Source object: %+v", source)

	// Parse source type
	sourceTypeRaw, exists := source["type"]
	if !exists {
		return args, fmt.Errorf("source.type is missing")
	}
	if sourceTypeRaw == nil {
		return args, fmt.Errorf("source.type is nil")
	}

	sourceType, ok := sourceTypeRaw.(string)
	if !ok {
		return args, fmt.Errorf("source.type must be a string, got %T", sourceTypeRaw)
	}
	args.Source.Type = sourceType

	log.Printf("[DEBUG] Source type: %s", sourceType)

	// Parse source fields based on type
	switch sourceType {
	case "file":
		if pathRaw, exists := source["path"]; exists && pathRaw != nil {
			if path, ok := pathRaw.(string); ok {
				args.Source.Path = path
				log.Printf("[DEBUG] Source path: %s", path)
			} else {
				return args, fmt.Errorf("source.path must be a string, got %T", pathRaw)
			}
		} else {
			return args, fmt.Errorf("source.path is required for file source")
		}
	case "data":
		if blobRaw, exists := source["blob"]; exists && blobRaw != nil {
			if blob, ok := blobRaw.(string); ok {
				args.Source.Blob = blob
				log.Printf("[DEBUG] Source blob length: %d", len(blob))
			} else {
				return args, fmt.Errorf("source.blob must be a string, got %T", blobRaw)
			}
		} else {
			return args, fmt.Errorf("source.blob is required for data source")
		}
	case "url":
		// Handle URL type (treat as file path for now)
		if urlRaw, exists := source["url"]; exists && urlRaw != nil {
			if url, ok := urlRaw.(string); ok {
				args.Source.Path = url    // LoadProgram can handle URLs
				args.Source.Type = "file" // Convert to file type
				log.Printf("[DEBUG] Source URL (as path): %s", url)
			} else {
				return args, fmt.Errorf("source.url must be a string, got %T", urlRaw)
			}
		} else {
			return args, fmt.Errorf("source.url is required for url source")
		}
	default:
		return args, fmt.Errorf("unsupported source type: %s", sourceType)
	}

	// Parse checksum (optional)
	if checksumRaw, exists := source["checksum"]; exists && checksumRaw != nil {
		if checksum, ok := checksumRaw.(string); ok {
			args.Source.Checksum = checksum
		}
	}

	// Parse program_type
	programTypeRaw, exists := input["program_type"]
	if !exists {
		return args, fmt.Errorf("program_type is missing")
	}
	if programTypeRaw == nil {
		return args, fmt.Errorf("program_type is nil")
	}

	programType, ok := programTypeRaw.(string)
	if !ok {
		return args, fmt.Errorf("program_type must be a string, got %T", programTypeRaw)
	}
	args.ProgramType = programType

	log.Printf("[DEBUG] Program type: %s", programType)

	// Parse optional fields
	if sectionRaw, exists := input["section"]; exists && sectionRaw != nil {
		if section, ok := sectionRaw.(string); ok {
			args.Section = section
		}
	}

	if btfPathRaw, exists := input["btf_path"]; exists && btfPathRaw != nil {
		if btfPath, ok := btfPathRaw.(string); ok {
			args.BTFPath = btfPath
		}
	}

	// Parse constraints (optional)
	if constraintsRaw, exists := input["constraints"]; exists && constraintsRaw != nil {
		if constraints, ok := constraintsRaw.(map[string]interface{}); ok {
			if maxInstrRaw, exists := constraints["max_instructions"]; exists && maxInstrRaw != nil {
				if maxInstr, ok := maxInstrRaw.(float64); ok {
					args.Constraints.MaxInstructions = int(maxInstr)
				}
			}

			if verifyOnlyRaw, exists := constraints["verify_only"]; exists && verifyOnlyRaw != nil {
				if verifyOnly, ok := verifyOnlyRaw.(bool); ok {
					args.Constraints.VerifyOnly = verifyOnly
				}
			}

			if helpersRaw, exists := constraints["allowed_helpers"]; exists && helpersRaw != nil {
				if helpers, ok := helpersRaw.([]interface{}); ok {
					for _, h := range helpers {
						if helper, ok := h.(string); ok {
							args.Constraints.AllowedHelpers = append(args.Constraints.AllowedHelpers, helper)
						}
					}
				}
			}
		}
	}

	log.Printf("[DEBUG] Final parsed args: %+v", args)
	return args, nil
}

func init() {
	RegisterTool(types.Tool{
		ID:          "load_program",
		Title:       "Load eBPF Program",
		Description: "Loads a raw eBPF object from file or base64 blob into the kernel.",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"source": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"type": map[string]interface{}{
							"type": "string",
							"enum": []string{"file", "data", "url"},
						},
						"path": map[string]interface{}{
							"type": "string",
						},
						"url": map[string]interface{}{
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
				},
				"program_type": map[string]interface{}{
					"type": "string",
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
		Call: LoadProgramTool,
	})
}
