package ebpf

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cilium/ebpf"
)

type DeployArgs struct {
	ProgramPath string `json:"program_path"`
}

type ProgramInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type DeployResult struct {
	Status   string        `json:"status"`
	Output   string        `json:"output"`
	Programs []ProgramInfo `json:"programs"`
}

// DeployProgram loads an eBPF object file into the kernel (local or remote).
func DeployProgram(args DeployArgs) (*DeployResult, error) {
	path := args.ProgramPath

	// Optional: fetch remote files
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		tmpfile, err := os.CreateTemp("", "ebpf-*.o")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp file: %w", err)
		}
		defer tmpfile.Close()
		resp, err := http.Get(path)
		if err != nil {
			return nil, fmt.Errorf("failed to download: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("bad HTTP response: %s", resp.Status)
		}
		_, err = io.Copy(tmpfile, resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to write to file: %w", err)
		}
		path = tmpfile.Name()
	}

	spec, err := ebpf.LoadCollectionSpec(path)
	if err != nil {
		return &DeployResult{Status: "error"}, fmt.Errorf("failed to load spec: %w", err)
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		return &DeployResult{Status: "error"}, fmt.Errorf("failed to create collection: %w", err)
	}
	defer coll.Close()

	programs := make([]ProgramInfo, 0)
	for name, prog := range coll.Programs {
		programs = append(programs, ProgramInfo{
			Name: name,
			Type: prog.Type().String(),
		})
		// Optional: attach here later
	}

	return &DeployResult{
		Status:   "ok",
		Output:   fmt.Sprintf("Loaded %d programs and %d maps", len(coll.Programs), len(coll.Maps)),
		Programs: programs,
	}, nil
}
