package ebpf

import (
	"errors"
	"fmt"
	"time"

	"github.com/sameehj/ebpf-mcp/pkg/types"
)

type StreamEventsArgs struct {
	Source struct {
		MapID int    `json:"map_id"`
		Type  string `json:"type"`
	} `json:"source"`
	SessionID  string                 `json:"session_id"`
	DurationMs int                    `json:"duration_ms"`
	MaxEvents  int                    `json:"max_events"`
	Filters    map[string]interface{} `json:"filters"`
	Format     string                 `json:"format"`
}

type StreamEventsResult struct {
	Events    []map[string]interface{}
	Stats     map[string]interface{}
	Complete  bool
	SessionID string
}

func ParseStreamEventsArgs(input map[string]interface{}) (*StreamEventsArgs, error) {
	var args StreamEventsArgs
	if err := types.StrictUnmarshal(input, &args); err != nil {
		return nil, err
	}
	return &args, nil
}

func StreamEvents(args *StreamEventsArgs, emit func(any)) error {
	if args.Source.Type != "ringbuf" && args.Source.Type != "perfbuf" {
		return errors.New("unsupported map type")
	}

	// Simulate event streaming by emitting dummy data
	start := time.Now()
	for i := 0; i < args.MaxEvents; i++ {
		event := map[string]interface{}{
			"timestamp_ns": time.Now().UnixNano(),
			"cpu":          i % 4,
			"pid":          1000 + i,
			"comm":         fmt.Sprintf("task-%d", i),
			"data": map[string]interface{}{
				"msg": fmt.Sprintf("event %d", i),
			},
		}
		emit(event)
		time.Sleep(time.Millisecond * 10) // simulate delay
	}

	duration := time.Since(start).Milliseconds()
	stats := map[string]interface{}{
		"events_received": args.MaxEvents,
		"events_dropped":  0,
		"duration_ms":     duration,
	}

	emit(map[string]interface{}{
		"success":      true,
		"tool_version": "1.0.0",
		"session_id":   args.SessionID,
		"events":       []interface{}{},
		"stats":        stats,
		"complete":     true,
	})

	return nil
}
