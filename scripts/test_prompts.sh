#!/bin/bash

# test_prompts.sh
# Comprehensive test script for prompts in eBPF MCP Server

set -e  # Exit on any error

# Configuration
SERVER_URL="http://localhost:8080/mcp"
KPROBE_URL="https://github.com/cilium/ebpf/raw/refs/heads/main/examples/kprobe/bpf_bpfel.o"
TEST_FILE="/tmp/kprobe_test.o"
SESSION_ID=""
TOKEN=""
PROGRAM_ID=""
MAP_ID=""
LINK_ID=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Initialize test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if server is running
check_server() {
    log_info "Checking if eBPF MCP server is running..."
    
    if ! curl -s "$SERVER_URL" > /dev/null 2>&1; then
        log_error "Server is not running at $SERVER_URL"
        log_info "Please start the server with: sudo ./bin/ebpf-mcp-server -t http -debug"
        exit 1
    fi
    
    log_success "Server is running"
}

# Get authentication token from server logs or environment
get_auth_token() {
    log_info "Getting authentication token..."
    
    if [ -n "$MCP_AUTH_TOKEN" ]; then
        TOKEN="$MCP_AUTH_TOKEN"
        log_success "Using token from environment variable"
    elif [ -n "$1" ]; then
        TOKEN="$1"
        log_success "Using token from command line argument"
    else
        log_warning "No MCP_AUTH_TOKEN environment variable found"
        log_info "Please check server logs for the Bearer token or set MCP_AUTH_TOKEN"
        log_info "Or run with: $0 <token>"
        read -p "Enter the Bearer token from server logs: " TOKEN
    fi
    
    if [ -z "$TOKEN" ]; then
        log_error "No authentication token provided"
        exit 1
    fi
}


# Initialize MCP session
init_session() {
    log_info "Initializing MCP session..."
    
    local response=$(curl -s -X POST "$SERVER_URL" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2025-03-26",
                "capabilities": {"prompts": {}},
                "clientInfo": {"name": "ebpf-prompt-test-client", "version": "1.0.0"}
            }
        }')
     
    # Extract session ID from headers
    SESSION_ID=$(curl -s -X POST "$SERVER_URL" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -D - \
        -d '{
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2025-03-26",
                "capabilities": {"prompts": {}},
                "clientInfo": {"name": "ebpf-prompt-test-client", "version": "1.0.0"}
            }
        }' 2>/dev/null | grep -i "mcp-session-id" | cut -d: -f2 | tr -d ' \r\n')
    
    if [ -z "$SESSION_ID" ]; then
        log_warning "No session ID found in headers, continuing without it"
        SESSION_ID="test-session"
    fi
    
    log_success "Session initialized: $SESSION_ID"
}

# Make MCP request
make_mcp_request() {
    local method="$1"
    local params="$2"
    local id="${3:-$(date +%s)}"
    
    local headers=(-H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN")
    if [ -n "$SESSION_ID" ]; then
        headers+=(-H "MCP-Session-ID: $SESSION_ID")
    fi
    
    local payload=$(cat << EOF
{
    "jsonrpc": "2.0",
    "id": $id,
    "method": "$method",
    "params": $params
}
EOF
)
    
    local response=$(curl -s -X POST "$SERVER_URL" "${headers[@]}" -d "$payload")
    
    echo "$response"
}

# List available prompts
test_list_prompts() {
    TESTS_RUN=$((TESTS_RUN + 1))
    log_info "Test 1: Listing available prompts..."

    local response=$(make_mcp_request "prompts/list" "{}")

    if echo "$response" | jq -e '.result.prompts' > /dev/null 2>&1; then
        local prompt_count=$(echo "$response" | jq '.result.prompts | length')
        
        if [[ "$prompt_count" -eq 0 ]]; then
            log_error "No prompts found"
            echo "$response" | jq '.'
            exit 1
        fi

        log_success "Found $prompt_count prompts"
        TESTS_PASSED=$((TESTS_PASSED + 1))
        echo "$response" | jq -r '.result.prompts[].name' | while read prompt; do
            log_info "  - $prompt"
        done
    else
        log_error "Failed to list prompts"
        echo "$response" | jq '.'
        TESTS_FAILED=$((TESTS_FAILED + 1))
        exit 1
    fi
}

test_get_empty_prompt() {
    TESTS_RUN=$((TESTS_RUN + 1))
    log_info "Test 2: Get empty prompt name"

    local response=$(make_mcp_request "prompts/get" "{}")
    local error_msg=$(echo "$response" | jq -r '.error.message // empty')

    if [[ "$error_msg" == *"prompt not found"* ]]; then
        log_success "Received expected error: '$error_msg'"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        log_error "Unexpected response or error: '$error_msg'"
        TESTS_FAILED=$((TESTS_FAILED + 1))
        exit 1
    fi
}

test_get_prompt_without_args() {
    TESTS_RUN=$((TESTS_RUN + 1))
    log_info "Test 3: Get show_system_info prompt"

    local params=$(cat << 'EOF'
{
    "name": "show_system_info",
    "arguments": {}
}
EOF
)

    local response=$(make_mcp_request "prompts/get" "$params")
    echo "$response"

    local result=$(echo "$response" | jq -r '.result')
    local error_msg=$(echo "$response" | jq -r '.error.message // empty')

    if [[ -n "$result" && "$result" != "null" ]]; then
        TESTS_PASSED=$((TESTS_PASSED + 1))
        log_success "Prompt returned result successfully."
    else
        TESTS_FAILED=$((TESTS_FAILED + 1))
        log_error "Unexpected error or missing result: '$error_msg'"
        exit 1
    fi
}

# Run all tests
run_all_tests() {
    log_info "Starting eBPF MCP prompt tests..."
    
    test_list_prompts
    echo
    test_get_empty_prompt
    echo
    test_get_prompt_without_args
    
    # Print summary
    echo
    log_info "Test Summary:"
    log_info "Tests run: $TESTS_RUN"
    log_success "Tests passed: $TESTS_PASSED"
    if [ $TESTS_FAILED -gt 0 ]; then
        log_error "Tests failed: $TESTS_FAILED"
    else
        log_success "Tests failed: $TESTS_FAILED"
    fi
    
    if [ $TESTS_FAILED -gt 0 ]; then
        log_error "Some tests failed!"
        exit 1
    else
        log_success "All tests passed!"
    fi
}

# Main execution
main() {
    echo "=============================================="
    echo "    eBPF MCP Server Comprehensive Test"
    echo "=============================================="
    echo
    
    # Pre-flight checks
    check_server
    get_auth_token "$1"
    init_session
    
    echo
    echo "Running test suite..."
    echo

    run_all_tests
}

# Run main function with all arguments
main "$@"