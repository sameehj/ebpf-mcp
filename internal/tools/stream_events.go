package tools

import (
	"github.com/sameehj/ebpf-mcp/internal/ebpf"
	"github.com/sameehj/ebpf-mcp/pkg/types"
)

func StreamEvents(input map[string]interface{}, emit func(interface{})) error {
	args, err := ebpf.ParseStreamEventsArgs(input)
	if err != nil {
		return err
	}
	return ebpf.StreamEvents(args, emit)
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
					"type":     "object",
					"required": []string{"map_id", "type"},
					"properties": map[string]interface{}{
						"map_id": map[string]interface{}{"type": "integer"},
						"type":   map[string]interface{}{"enum": []string{"ringbuf", "perfbuf"}},
					},
				},
				"session_id":  map[string]interface{}{"type": "string"},
				"duration_ms": map[string]interface{}{"type": "integer", "default": 10000},
				"max_events":  map[string]interface{}{"type": "integer", "default": 1000},
				"filters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"pid":              map[string]interface{}{"type": "integer"},
						"comm":             map[string]interface{}{"type": "string"},
						"cpu":              map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "integer"}},
						"min_timestamp_ns": map[string]interface{}{"type": "integer"},
					},
				},
				"format": map[string]interface{}{
					"enum":    []string{"json", "raw", "base64"},
					"default": "json",
				},
			},
		},
		OutputSchema: map[string]interface{}{
			"type":     "object",
			"required": []string{"success", "tool_version"},
			"properties": map[string]interface{}{
				"success":      map[string]interface{}{"type": "boolean"},
				"tool_version": map[string]interface{}{"type": "string"},
				"session_id":   map[string]interface{}{"type": "string"},
				"events":       map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "object"}},
				"stats":        map[string]interface{}{"type": "object"},
				"complete":     map[string]interface{}{"type": "boolean"},
				"error": map[string]interface{}{
					"$ref": "https://ebpf-mcp.dev/schemas/error.schema.json#/definitions/Error",
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
