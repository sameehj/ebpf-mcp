# eBPF MCP Server Makefile
.PHONY: all build clean test fmt release help install dev

# Build configuration
BINARY_NAME := ebpf-mcp-server
CMD_DIR := ./cmd/ebpf-mcp
OUT_DIR := ./bin
RELEASE_DIR := ./releases

# Version information
VERSION := $(shell cat VERSION 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)
BUILD_FLAGS := -ldflags="$(LDFLAGS)"

# Debug configuration
ifeq ($(DEBUG),1)
    GO_BUILD_FLAGS = -x
    VERBOSE = 1
else
    GO_BUILD_FLAGS = 
    VERBOSE = 0
endif

# Default target
help: ## Show this help message
	@echo "eBPF MCP Server Build System"
	@echo "============================="
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

all: build ## Build all binaries

build: server build-chat ## Build eBPF MCP server and chat tool

server: ## Build the main eBPF MCP server
	@echo "🔧 Building $(BINARY_NAME)..."
	@mkdir -p $(OUT_DIR)
	GO111MODULE=on go mod tidy
	go build $(GO_BUILD_FLAGS) $(BUILD_FLAGS) -o $(OUT_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✅ Build complete: $(OUT_DIR)/$(BINARY_NAME)"

build-chat: ## Build the Ollama chat tool
	@echo "🔧 Building ollama-chat..."
	@mkdir -p $(OUT_DIR)
	go build -o $(OUT_DIR)/ollama-chat ./cmd/ollama-chat
	@echo "✅ Chat tool complete: $(OUT_DIR)/ollama-chat"

dev: build ## Start development server with debug logging
	@echo "🚀 Starting development server..."
	@sudo $(OUT_DIR)/$(BINARY_NAME) -t http -debug

run: build ## Run the server (deprecated, use 'dev')
	@echo "🚀 Running server..."
	@$(OUT_DIR)/$(BINARY_NAME)

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(OUT_DIR)
	@rm -rf $(RELEASE_DIR)
	@echo "✅ Clean complete"

test: ## Run tests
	@echo "🧪 Running tests..."
	@sudo go test ./internal/tests/tools
	@echo "✅ Tests complete"

test-all: ## Run all tests including integration tests
	@echo "🧪 Running all tests..."
	@sudo go test -v ./...
	@./scripts/test-ebpf-mcp-server.sh || true
	@echo "✅ All tests complete"

fmt: ## Format Go code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

lint: ## Run linters
	@echo "🔍 Running linters..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, skipping..."; \
	fi

# Release targets
release: clean release-all ## Build release binaries for all platforms
	@echo "🚀 Building release $(VERSION)..."
	@mkdir -p $(RELEASE_DIR)/$(VERSION)

	@echo "  📦 Building Linux x86_64..."
	@GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(RELEASE_DIR)/$(VERSION)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	@tar -czf $(RELEASE_DIR)/$(VERSION)/$(BINARY_NAME)-linux-amd64.tar.gz -C $(RELEASE_DIR)/$(VERSION) $(BINARY_NAME)-linux-amd64

	@echo "  📦 Building Linux ARM64..."
	@GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(RELEASE_DIR)/$(VERSION)/$(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	@tar -czf $(RELEASE_DIR)/$(VERSION)/$(BINARY_NAME)-linux-arm64.tar.gz -C $(RELEASE_DIR)/$(VERSION) $(BINARY_NAME)-linux-arm64

	@echo "  📦 Building macOS x86_64..."
	@GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(RELEASE_DIR)/$(VERSION)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	@tar -czf $(RELEASE_DIR)/$(VERSION)/$(BINARY_NAME)-darwin-amd64.tar.gz -C $(RELEASE_DIR)/$(VERSION) $(BINARY_NAME)-darwin-amd64

	@echo "  📦 Building macOS ARM64..."
	@GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(RELEASE_DIR)/$(VERSION)/$(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	@tar -czf $(RELEASE_DIR)/$(VERSION)/$(BINARY_NAME)-darwin-arm64.tar.gz -C $(RELEASE_DIR)/$(VERSION) $(BINARY_NAME)-darwin-arm64

	@echo "  🔐 Generating checksums..."
	@cd $(RELEASE_DIR)/$(VERSION) && sha256sum *.tar.gz *.tar.gz > checksums.txt

	@echo "  📋 Creating release info..."
	@echo "Version: $(VERSION)" > $(RELEASE_DIR)/$(VERSION)/release-info.txt
	@echo "Commit: $(COMMIT)" >> $(RELEASE_DIR)/$(VERSION)/release-info.txt
	@echo "Build Time: $(BUILD_TIME)" >> $(RELEASE_DIR)/$(VERSION)/release-info.txt

	@echo "✅ Release $(VERSION) complete in $(RELEASE_DIR)/$(VERSION)"
	@ls -la $(RELEASE_DIR)/$(VERSION)


release-all: ## Internal target for building all platforms
	@echo "Building release binaries..."

install: build ## Install the binary system-wide
	@echo "📦 Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(OUT_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Installation complete"

uninstall: ## Remove installed binary
	@echo "🗑️  Removing $(BINARY_NAME) from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstall complete"

# Development helpers
test-map-dump: ## Test map dump functionality
	@./scripts/test_map_dump.sh

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t ebpf-mcp:$(VERSION) .
	@echo "✅ Docker image built: ebpf-mcp:$(VERSION)"

docker-run: docker-build ## Run in Docker container
	@echo "🐳 Running in Docker..."
	@docker run --rm --privileged -p 8080:8080 ebpf-mcp:$(VERSION)

# CI/CD helpers
ci-test: ## Run tests in CI environment
	@echo "🤖 Running CI tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Built:   $(BUILD_TIME)"

# Dependency management
deps: ## Download and tidy dependencies
	@echo "📦 Managing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies updated"

# Documentation
docs: ## Generate documentation
	@echo "📚 Generating documentation..."
	@if command -v godoc &> /dev/null; then \
		echo "Starting godoc server at http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found, install with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi