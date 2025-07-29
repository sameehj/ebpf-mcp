// main.go - eBPF-MCP server with debug logging
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/server"
	"github.com/sameehj/ebpf-mcp/internal/prompts"
	"github.com/sameehj/ebpf-mcp/internal/tools"
)

var debugMode bool

func main() {
	var transport string
	var port string
	flag.StringVar(&port, "port", "8080", "Port to listen on")
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or http)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or http)")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug logging")
	flag.Parse()

	// Configure logging
	if debugMode {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Printf("[DEBUG] Debug logging enabled")
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// Create MCP server
	log.Printf("[DEBUG] Creating MCP server...")
	mcpServer := server.NewMCPServer(
		"ebpf-mcp",
		"0.1.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	// Register all tools
	log.Printf("[DEBUG] Registering tools...")
	tools.RegisterAllWithMCP(mcpServer)

	// Register all prompts
	log.Printf("[DEBUG] Registering prompts...")
	prompts.RegisterAllWithMCP(mcpServer)

	if transport == "http" {
		token := os.Getenv("MCP_AUTH_TOKEN")
		if token == "" {
			token = generateRandomToken()
			log.Println("🔐 Auto-generated auth token (no MCP_AUTH_TOKEN was set):")
		} else {
			log.Println("🔐 Using MCP_AUTH_TOKEN from environment:")
		}
		log.Printf("    %s\n", token)
		log.Println("💡 Pass this as Authorization: Bearer <token> in Claude or curl headers.")

		mux := http.NewServeMux()

		mux.HandleFunc("/.well-known/mcp/metadata.json", func(w http.ResponseWriter, r *http.Request) {
			if debugMode {
				log.Printf("[DEBUG] Metadata request from %s", r.RemoteAddr)
			}
			w.Header().Set("Content-Type", "application/json")
			response := map[string]interface{}{
				"schema_version": "v1",
				"entrypoint_url": "http://localhost:" + port + "/mcp",
				"display_name":   "eBPF MCP Server",
				"description":    "Exposes Linux kernel tools via MCP protocol",
				"tool_filter":    "all",
			}
			json.NewEncoder(w).Encode(response)
		})

		httpServer := server.NewStreamableHTTPServer(mcpServer)
		authenticated := tokenAuthMiddleware(token, httpServer)
		mux.Handle("/mcp", authenticated)

		log.Printf("\U0001F527 ebpf-mcp HTTP server listening on :%s", port)
		log.Printf("   MCP endpoint: http://localhost:%s/mcp", port)
		log.Printf("   Discovery: http://localhost:%s/.well-known/mcp/metadata.json", port)

		if err := http.ListenAndServe(":" + port, mux); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		log.Printf("\U0001F527 ebpf-mcp stdio server starting...")
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}

func generateRandomToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}
	return hex.EncodeToString(b)
}

func tokenAuthMiddleware(expectedToken string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if debugMode {
			log.Printf("[DEBUG] Auth request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
			log.Printf("[DEBUG] Headers: %v", r.Header)
		}

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			if debugMode {
				log.Printf("[DEBUG] Missing Bearer token in Authorization header: '%s'", authHeader)
			}
			http.Error(w, "Unauthorized: Missing Bearer token", http.StatusUnauthorized)
			return
		}

		providedToken := strings.TrimPrefix(authHeader, "Bearer ")
		if providedToken != expectedToken {
			if debugMode {
				log.Printf("[DEBUG] Invalid token provided: '%s' (expected: '%s')", providedToken, expectedToken)
			}
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		if debugMode {
			log.Printf("[DEBUG] Authentication successful")
		}

		ctx := context.WithValue(r.Context(), "authToken", providedToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
