module github.com/sameehj/ebpf-mcp

go 1.23.0

toolchain go1.23.10

require (
	github.com/cilium/ebpf v0.18.0
	github.com/mark3labs/mcp-go v0.0.0
	golang.org/x/sys v0.30.0
)

replace github.com/mark3labs/mcp-go => github.com/mark3labs/mcp-go v0.32.0

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
)

replace github.com/sameehj/ebpf-mcp => .
