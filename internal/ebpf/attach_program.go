// internal/ebpf/attach_program.go
package ebpf

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

type AttachProgramArgs struct {
	ProgramID  int `json:"program_id"`
	Attachment struct {
		Type   string                 `json:"type"`
		Params map[string]interface{} `json:"params"`
	} `json:"attachment"`
}

type AttachResult struct {
	Success         bool         `json:"success"`
	ToolVersion     string       `json:"tool_version"`
	AttachmentPoint string       `json:"attachment_point,omitempty"`
	PinPath         string       `json:"pin_path,omitempty"`
	Error           *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func ParseAttachProgramArgs(input map[string]interface{}) (*AttachProgramArgs, error) {
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}
	var args AttachProgramArgs
	if err := json.Unmarshal(bytes, &args); err != nil {
		return nil, fmt.Errorf("invalid arguments: %w", err)
	}
	return &args, nil
}

func AttachProgram(args *AttachProgramArgs) (*AttachResult, error) {
	prog, err := ebpf.NewProgramFromID(ebpf.ProgramID(args.ProgramID))
	if err != nil {
		return &AttachResult{Success: false, Error: &ErrorDetail{Type: "PROGRAM_NOT_FOUND", Message: err.Error()}}, nil
	}

	switch args.Attachment.Type {
	case "xdp":
		iface, _ := args.Attachment.Params["interface"].(string)
		mode := args.Attachment.Params["mode"].(string)
		pinPath, _ := args.Attachment.Params["link_pin_path"].(string)

		ifaceIdx, err := ifaceToIndex(iface)
		if err != nil {
			return &AttachResult{Success: false, Error: &ErrorDetail{Type: "ATTACHMENT_FAILED", Message: err.Error()}}, nil
		}

		linkObj, err := link.AttachXDP(link.XDPOptions{
			Program:   prog,
			Interface: ifaceIdx,
			Flags:     xdpModeToFlag(mode),
		})
		if err != nil {
			return &AttachResult{Success: false, Error: &ErrorDetail{Type: "ATTACHMENT_FAILED", Message: err.Error()}}, nil
		}
		if pinPath != "" {
			_ = os.MkdirAll(pinPath, 0755)
			_ = linkObj.Pin(pinPath)
		}
		return &AttachResult{
			Success:         true,
			ToolVersion:     "v1",
			AttachmentPoint: iface,
			PinPath:         pinPath,
		}, nil

	default:
		return &AttachResult{Success: false, Error: &ErrorDetail{Type: "VALIDATION_ERROR", Message: "unsupported attachment type"}}, nil
	}
}

func ifaceToIndex(name string) (int, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return 0, err
	}
	return iface.Index, nil
}

func xdpModeToFlag(mode string) link.XDPAttachFlags {
	switch mode {
	case "GENERIC":
		return link.XDPGenericMode
	case "OFFLOAD":
		return link.XDPOffloadMode
	default:
		return link.XDPDriverMode
	}
}
