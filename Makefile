BINARY_NAME=ebpf-mcp-server
PKG=github.com/sameehj/ebpf-mcp
CMD_DIR=./cmd/ebpf-mcp
OUT_DIR=./bin

.PHONY: all build run clean test fmt

all: build

build:
	go build -o $(OUT_DIR)/$(BINARY_NAME) $(CMD_DIR)

run:
	go run $(CMD_DIR)

clean:
	rm -rf $(OUT_DIR)

test:
	go test ./...

fmt:
	go fmt ./...
