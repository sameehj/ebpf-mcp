package tools

import (
	"fmt"
	"log"

	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func StreamEvents(input map[string]interface{}, emit func(interface{})) error {
	// Add defensive programming and logging
	log.Printf("[DEBUG] StreamEvents called with input: %+v", input)

	if input == nil {
		return fmt.Errorf("input is nil")
	}

	if emit == nil {
		return fmt.Errorf("emit function is nil")
	}

	// Parse arguments with error handling
	args, err := ebpf.ParseStreamEventsArgs(input)
	if err != nil {
		log.Printf("[ERROR] Failed to parse stream events args: %v", err)
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	log.Printf("[DEBUG] Parsed args: %+v", args)

	// Call the streaming function with proper error handling
	if err := ebpf.StreamEvents(args, emit); err != nil {
		log.Printf("[ERROR] StreamEvents failed: %v", err)
		return fmt.Errorf("streaming failed: %w", err)
	}

	return nil
}

func init() {
	RegisterTool(types.Tool{
		ID:          "stream_events",
		Title:       "Stream Kernel Events",
		Description: "Streams real-time eBPF events using perf or ringbuf.",
		InputSchema: map[string]interface{}{
			"type":     "object",
			"required": []string{"source"},
			"properties": map[string]interface{}{
				"source": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"program_id": map[string]interface{}{
							"type":        "integer",
							"description": "ID of the eBPF program to stream from",
						},
						"link_id": map[string]interface{}{
							"type":        "integer",
							"description": "ID of the attached link to stream from",
						},
						"map_id": map[string]interface{}{
							"type":        "integer",
							"description": "Direct map ID to stream from",
						},
						"type": map[string]interface{}{
							"type":        "string",
							"enum":        []string{"ringbuf", "perfbuf"},
							"default":     "ringbuf",
							"description": "Type of buffer to stream from",
						},
					},
					"anyOf": []map[string]interface{}{
						{"required": []string{"program_id"}},
						{"required": []string{"link_id"}},
						{"required": []string{"map_id"}},
					},
				},
				"duration": map[string]interface{}{
					"type":        "integer",
					"description": "Duration in seconds",
					"default":     5,
					"minimum":     1,
					"maximum":     300,
				},
				"duration_ms": map[string]interface{}{
					"type":        "integer",
					"description": "Duration in milliseconds (takes precedence over duration)",
					"minimum":     100,
					"maximum":     300000,
				},
				"max_events": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of events to capture",
					"default":     100,
					"minimum":     1,
					"maximum":     10000,
				},
				"format": map[string]interface{}{
					"type":        "string",
					"enum":        []string{"json", "raw", "base64"},
					"default":     "json",
					"description": "Output format for events",
				},
				"filters": map[string]interface{}{
					"type":        "object",
					"description": "Event filtering options",
					"properties": map[string]interface{}{
						"event_types": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "Filter by event types (e.g., kprobe, tracepoint)",
						},
						"target_pids": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "integer"},
							"description": "Filter by process IDs",
						},
						"min_timestamp": map[string]interface{}{
							"type":        "integer",
							"description": "Minimum timestamp (nanoseconds since epoch)",
						},
						"max_events": map[string]interface{}{
							"type":        "integer",
							"description": "Maximum events to return",
							"minimum":     1,
							"maximum":     10000,
						},
					},
				},
			},
		},
		OutputSchema: map[string]interface{}{
			"type":     "object",
			"required": []string{"success", "tool_version"},
			"properties": map[string]interface{}{
				"success": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether the operation succeeded",
				},
				"tool_version": map[string]interface{}{
					"type":        "string",
					"description": "Version of the tool",
				},
				"session_id": map[string]interface{}{
					"type":        "string",
					"description": "Unique session identifier for this stream",
				},
				"events": map[string]interface{}{
					"type":        "array",
					"items":       map[string]interface{}{"type": "object"},
					"description": "Array of captured events",
				},
				"stats": map[string]interface{}{
					"type":        "object",
					"description": "Statistics about the streaming session",
					"properties": map[string]interface{}{
						"events_received":   map[string]interface{}{"type": "integer"},
						"events_dropped":    map[string]interface{}{"type": "integer"},
						"duration_ms":       map[string]interface{}{"type": "integer"},
						"events_per_second": map[string]interface{}{"type": "number"},
					},
				},
				"complete": map[string]interface{}{
					"type":        "boolean",
					"description": "Whether the stream completed successfully",
				},
				"message": map[string]interface{}{
					"type":        "string",
					"description": "Human-readable status message",
				},
				"error": map[string]interface{}{
					"type":        "object",
					"description": "Error information if the operation failed",
					"properties": map[string]interface{}{
						"code":    map[string]interface{}{"type": "string"},
						"message": map[string]interface{}{"type": "string"},
						"details": map[string]interface{}{"type": "object"},
					},
				},
			},
		},
		Annotations: map[string]interface{}{
			"title":          "Stream Events",
			"streamHint":     true,
			"idempotentHint": false,
			"readOnlyHint":   true,
		},
		Call:   nil,
		Stream: StreamEvents,
	})
}
