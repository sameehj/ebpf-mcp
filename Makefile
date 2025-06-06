# Output binary name
BINARY_NAME=ebpf-mcp-server
CMD_DIR=./cmd/ebpf-mcp
OUT_DIR=./bin

# Check for debug flag
ifeq ($(DEBUG),1)
    GO_BUILD_FLAGS=-x
    VERBOSE=1
else
    GO_BUILD_FLAGS=
    VERBOSE=0
endif

.PHONY: all build run clean test fmt

all: build

build:
	@echo "🔧 Building $(BINARY_NAME)..."
	@mkdir -p $(OUT_DIR)
	GO111MODULE=on go mod tidy
	go build $(GO_BUILD_FLAGS) -o $(OUT_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✅ Build complete: $(OUT_DIR)/$(BINARY_NAME)"

run: build
	@echo "🚀 Running server..."
	@$(OUT_DIR)/$(BINARY_NAME)

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(OUT_DIR)

test:
	@echo "🧪 Running tests..."
	@go test ./...

fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

test-map-dump:
	./scripts/test_map_dump.sh
