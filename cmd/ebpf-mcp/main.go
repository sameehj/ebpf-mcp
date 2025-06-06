package main

import (
    "log"
    "ebpf-mcp/internal/server"
)

func main() {
    if err := server.Start(); err != nil {
        log.Fatalf("failed to start server: %v", err)
    }
}
