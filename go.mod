module github.com/sameehj/ebpf-mcp

go 1.23.0

toolchain go1.23.10

require (
	github.com/cilium/ebpf v0.18.0
	github.com/gorilla/mux v1.8.0
)

require golang.org/x/sys v0.30.0 // indirect

replace github.com/sameehj/ebpf-mcp => .
