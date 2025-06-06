// internal/ebpf/map_dump.go
package ebpf

import (
	"fmt"
	"time"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"log"
)

type DumpResult struct {
	Entries        []MapEntry `json:"entries"`
	TotalEntries   int        `json:"total_entries"`
	ExecutionTimeMs int       `json:"execution_time_ms"`
}

type MapEntry struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// DumpPinnedMap opens a pinned eBPF map and returns its contents.
func DumpPinnedMap(path string, maxEntries int) (*DumpResult, error) {
	start := time.Now()

	// Raise rlimit so we can access pinned eBPF maps.
	if err := rlimit.RemoveMemlock(); err != nil {
		return nil, fmt.Errorf("failed to set rlimit: %w", err)
	}

	m, err := ebpf.LoadPinnedMap(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load pinned map: %w", err)
	}
	defer m.Close()

	entries := []MapEntry{}
	it := m.Iterate()
	var k, v []byte
	count := 0

	for it.Next(&k, &v) {
		entries = append(entries, MapEntry{
			Key:   fmt.Sprintf("0x%x", k),
			Value: fmt.Sprintf("0x%x", v),
		})
		count++
		if count >= maxEntries {
			break
		}
	}

	if err := it.Err(); err != nil {
		log.Printf("iteration error: %v", err)
	}

	duration := time.Since(start)
	return &DumpResult{
		Entries:        entries,
		TotalEntries:   count,
		ExecutionTimeMs: int(duration.Milliseconds()),
	}, nil
}
