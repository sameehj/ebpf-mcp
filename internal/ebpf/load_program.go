package ebpf

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cilium/ebpf"
)

type LoadProgramArgs struct {
	Source struct {
		Type     string `json:"type"`
		Path     string `json:"path,omitempty"`
		Blob     string `json:"blob,omitempty"`
		Checksum string `json:"checksum,omitempty"`
	} `json:"source"`
	ProgramType string `json:"program_type"`
	Section     string `json:"section,omitempty"`
	BTFPath     string `json:"btf_path,omitempty"`
	Constraints struct {
		MaxInstructions int      `json:"max_instructions,omitempty"`
		AllowedHelpers  []string `json:"allowed_helpers,omitempty"`
		VerifyOnly      bool     `json:"verify_only,omitempty"`
	} `json:"constraints,omitempty"`
}

type MapInfo struct {
	Name       string `json:"name"`
	FD         int    `json:"fd"`
	ID         int    `json:"id"`
	Type       string `json:"type"`
	KeySize    int    `json:"key_size,omitempty"`
	ValueSize  int    `json:"value_size,omitempty"`
	MaxEntries int    `json:"max_entries,omitempty"`
	PinPath    string `json:"pin_path,omitempty"`
}

type LoadProgramResult struct {
	Success      bool      `json:"success"`
	ToolVersion  string    `json:"tool_version"`
	ProgramFD    int       `json:"program_fd,omitempty"`
	ProgramID    int       `json:"program_id,omitempty"`
	Maps         []MapInfo `json:"maps,omitempty"`
	VerifierLog  string    `json:"verifier_log,omitempty"`
	ErrorMessage string    `json:"error,omitempty"`
}

func ParseLoadProgramArgs(input map[string]interface{}) (LoadProgramArgs, error) {
	var args LoadProgramArgs

	// Basic validation here is optional, since JSON schema does the real check
	src := input["source"].(map[string]interface{})
	args.Source.Type = src["type"].(string)
	if args.Source.Type == "file" {
		args.Source.Path = src["path"].(string)
	} else if args.Source.Type == "data" {
		args.Source.Blob = src["blob"].(string)
	}

	args.ProgramType = input["program_type"].(string)

	if sec, ok := input["section"].(string); ok {
		args.Section = sec
	}
	if btf, ok := input["btf_path"].(string); ok {
		args.BTFPath = btf
	}

	if rawConstraints, ok := input["constraints"].(map[string]interface{}); ok {
		if max, ok := rawConstraints["max_instructions"].(float64); ok {
			args.Constraints.MaxInstructions = int(max)
		}
		if verify, ok := rawConstraints["verify_only"].(bool); ok {
			args.Constraints.VerifyOnly = verify
		}
		if helpers, ok := rawConstraints["allowed_helpers"].([]interface{}); ok {
			for _, h := range helpers {
				if hs, ok := h.(string); ok {
					args.Constraints.AllowedHelpers = append(args.Constraints.AllowedHelpers, hs)
				}
			}
		}
	}

	return args, nil
}

func LoadProgram(args LoadProgramArgs) (*LoadProgramResult, error) {
	var spec *ebpf.CollectionSpec
	var err error

	switch args.Source.Type {
	case "file":
		path := args.Source.Path
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			tmpfile, err := os.CreateTemp("", "ebpf-*.o")
			if err != nil {
				return &LoadProgramResult{Success: false, ErrorMessage: "temp file error"}, err
			}
			defer tmpfile.Close()
			resp, err := http.Get(path)
			if err != nil || resp.StatusCode != http.StatusOK {
				return &LoadProgramResult{Success: false, ErrorMessage: "download failed"}, err
			}
			io.Copy(tmpfile, resp.Body)
			path = tmpfile.Name()
		}
		spec, err = ebpf.LoadCollectionSpec(path)
	case "data":
		blob, err := base64.StdEncoding.DecodeString(args.Source.Blob)
		if err != nil {
			return &LoadProgramResult{Success: false, ErrorMessage: "invalid base64"}, err
		}
		spec, err = ebpf.LoadCollectionSpecFromReader(bytes.NewReader(blob))
	default:
		return &LoadProgramResult{Success: false, ErrorMessage: "invalid source type"}, errors.New("invalid source type")
	}

	if err != nil {
		return &LoadProgramResult{Success: false, ErrorMessage: err.Error()}, err
	}

	coll, err := ebpf.NewCollection(spec)
	if err != nil {
		return &LoadProgramResult{Success: false, ErrorMessage: err.Error()}, err
	}
	defer coll.Close()

	maps := make([]MapInfo, 0)
	for name, m := range coll.Maps {
		info, err := m.Info()
		if err != nil {
			continue
		}
		mid, _ := info.ID()
		maps = append(maps, MapInfo{
			Name:       name,
			FD:         m.FD(),
			ID:         int(mid),
			Type:       info.Type.String(),
			KeySize:    int(info.KeySize),
			ValueSize:  int(info.ValueSize),
			MaxEntries: int(info.MaxEntries),
		})
	}

	// Assuming only one program loaded for now
	for _, prog := range coll.Programs {
		info, err := prog.Info()
		if err != nil {
			continue
		}
		pid, _ := info.ID()
		return &LoadProgramResult{
			Success:     true,
			ToolVersion: "1.0.0",
			ProgramFD:   prog.FD(),
			ProgramID:   int(pid),
			Maps:        maps,
		}, nil
	}

	return &LoadProgramResult{Success: false, ErrorMessage: "no programs found"}, errors.New("no programs found")
}
