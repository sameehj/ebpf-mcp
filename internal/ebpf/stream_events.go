package ebpf

import (
	"errors"
	"fmt"
	"time"

	"github.com/sameehj/ebpf-mcp/pkg/types"
)

type StreamEventsArgs struct {
	Source struct {
		ProgramID int    `json:"program_id,omitempty"`
		LinkID    int    `json:"link_id,omitempty"`
		MapID     int    `json:"map_id,omitempty"`
		Type      string `json:"type,omitempty"`
	} `json:"source"`
	Duration   int                    `json:"duration,omitempty"`    // seconds
	DurationMs int                    `json:"duration_ms,omitempty"` // milliseconds
	MaxEvents  int                    `json:"max_events,omitempty"`
	Filters    map[string]interface{} `json:"filters,omitempty"`
	Format     string                 `json:"format,omitempty"`
}

type StreamEventsResult struct {
	Events    []map[string]interface{} `json:"events"`
	Stats     map[string]interface{}   `json:"stats"`
	Complete  bool                     `json:"complete"`
	SessionID string                   `json:"session_id"`
}

func ParseStreamEventsArgs(input map[string]interface{}) (*StreamEventsArgs, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	var args StreamEventsArgs
	if err := types.StrictUnmarshal(input, &args); err != nil {
		return nil, fmt.Errorf("failed to parse stream events args: %w", err)
	}
	return &args, nil
}

func StreamEvents(args *StreamEventsArgs, emit func(any)) error {
	// Validate inputs
	if args == nil {
		return errors.New("args cannot be nil")
	}
	if emit == nil {
		return errors.New("emit function cannot be nil")
	}

	// Set default type if not provided
	sourceType := args.Source.Type
	if sourceType == "" {
		sourceType = "ringbuf" // default
	}

	// Validate type
	if sourceType != "ringbuf" && sourceType != "perfbuf" {
		return fmt.Errorf("unsupported map type: %s", sourceType)
	}

	// Determine the source
	var mapID int
	var sourceDesc string

	if args.Source.MapID != 0 {
		mapID = args.Source.MapID
		sourceDesc = fmt.Sprintf("map %d", mapID)
	} else if args.Source.ProgramID != 0 {
		// For now, use the program ID as map ID
		// (in real implementation, you'd look up the actual map)
		mapID = args.Source.ProgramID
		sourceDesc = fmt.Sprintf("program %d", args.Source.ProgramID)
	} else if args.Source.LinkID != 0 {
		// For now, use the link ID as map ID
		// (in real implementation, you'd look up the actual map)
		mapID = args.Source.LinkID
		sourceDesc = fmt.Sprintf("link %d", args.Source.LinkID)
	} else {
		return errors.New("must specify program_id, link_id, or map_id in source")
	}

	// Set duration (prefer duration_ms, fall back to duration * 1000)
	durationMs := args.DurationMs
	if durationMs == 0 && args.Duration > 0 {
		durationMs = args.Duration * 1000
	}
	if durationMs == 0 {
		durationMs = 5000 // default 5 seconds
	}

	// Set max events
	maxEvents := args.MaxEvents
	if maxEvents == 0 {
		maxEvents = 100 // default
	}

	// Set format
	format := args.Format
	if format == "" {
		format = "json"
	}

	// Validate format
	if format != "json" && format != "raw" && format != "base64" {
		return fmt.Errorf("unsupported format: %s", format)
	}

	// Initialize session
	sessionID := fmt.Sprintf("stream-session-%d", time.Now().Unix())

	// Start streaming simulation
	start := time.Now()
	eventsGenerated := 0

	// Emit initial status
	emit(map[string]interface{}{
		"type":       "status",
		"message":    fmt.Sprintf("Starting event stream from %s for %dms", sourceDesc, durationMs),
		"session_id": sessionID,
		"format":     format,
	})

	// Generate events over the duration
	eventInterval := time.Duration(durationMs) * time.Millisecond / time.Duration(maxEvents)
	if eventInterval < time.Millisecond {
		eventInterval = time.Millisecond
	}

	ticker := time.NewTicker(eventInterval)
	defer ticker.Stop()

	timeout := time.After(time.Duration(durationMs) * time.Millisecond)

eventLoop:
	for eventsGenerated < maxEvents {
		select {
		case <-timeout:
			break eventLoop
		case <-ticker.C:
			// Generate a mock event
			event := generateMockEvent(eventsGenerated, mapID, format)

			// Apply filters if any
			if shouldIncludeEvent(event, args.Filters) {
				emit(map[string]interface{}{
					"type":  "event",
					"event": event,
				})
				eventsGenerated++
			}
		}
	}

	// Calculate final stats
	actualDuration := time.Since(start).Milliseconds()
	stats := map[string]interface{}{
		"events_received":   eventsGenerated,
		"events_dropped":    0,
		"duration_ms":       actualDuration,
		"events_per_second": float64(eventsGenerated) / (float64(actualDuration) / 1000.0),
		"source":            sourceDesc,
		"format":            format,
	}

	// Emit final result
	result := map[string]interface{}{
		"success":      true,
		"tool_version": "1.0.0",
		"session_id":   sessionID,
		"stats":        stats,
		"complete":     true,
		"message":      fmt.Sprintf("Stream completed: %d events in %dms", eventsGenerated, actualDuration),
	}

	emit(result)
	return nil
}

func generateMockEvent(index, mapID int, format string) map[string]interface{} {
	timestamp := time.Now().UnixNano()

	baseEvent := map[string]interface{}{
		"timestamp_ns": timestamp,
		"cpu":          index % 4,
		"pid":          1000 + index,
		"tid":          1000 + index,
		"comm":         fmt.Sprintf("task-%d", index),
		"map_id":       mapID,
		"index":        index,
	}

	switch format {
	case "json":
		baseEvent["data"] = map[string]interface{}{
			"syscall":   "execve",
			"args":      []string{"/bin/ls", "-la"},
			"exit_code": 0,
			"message":   fmt.Sprintf("kprobe event %d from map %d", index, mapID),
		}
	case "raw":
		baseEvent["raw_data"] = fmt.Sprintf("raw_event_%d_map_%d_ts_%d", index, mapID, timestamp)
	case "base64":
		// Simulate base64 encoded data
		data := fmt.Sprintf("event_%d_map_%d", index, mapID)
		baseEvent["data_base64"] = data // In real implementation, this would be base64 encoded
	}

	return baseEvent
}

func shouldIncludeEvent(event map[string]interface{}, filters map[string]interface{}) bool {
	if filters == nil {
		return true
	}

	// Apply event type filter
	if eventTypes, ok := filters["event_types"].([]interface{}); ok {
		if len(eventTypes) > 0 {
			// For mock events, we'll assume they're all "kprobe" type
			hasKprobe := false
			for _, et := range eventTypes {
				if et == "kprobe" {
					hasKprobe = true
					break
				}
			}
			if !hasKprobe {
				return false
			}
		}
	}

	// Apply PID filter
	if targetPids, ok := filters["target_pids"].([]interface{}); ok {
		if len(targetPids) > 0 {
			eventPid, hasPid := event["pid"].(int)
			if hasPid {
				pidMatch := false
				for _, pid := range targetPids {
					if pidInt, ok := pid.(int); ok && pidInt == eventPid {
						pidMatch = true
						break
					}
				}
				if !pidMatch {
					return false
				}
			}
		}
	}

	// Apply timestamp filter
	if minTimestamp, ok := filters["min_timestamp"].(float64); ok {
		if eventTs, hasTs := event["timestamp_ns"].(int64); hasTs {
			if float64(eventTs) < minTimestamp {
				return false
			}
		}
	}

	// Apply max events filter (handled in main loop)
	return true
}
