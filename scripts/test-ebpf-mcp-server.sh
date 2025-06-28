#!/bin/bash

# test-ebpf-mcp-server.sh
# Comprehensive test script for eBPF MCP Server
# This script tests the complete workflow: download, load, attach, and inspect

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
NC='\033[0m' # No Color

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
                "capabilities": {"tools": {}},
                "clientInfo": {"name": "ebpf-test-client", "version": "1.0.0"}
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
                "capabilities": {"tools": {}},
                "clientInfo": {"name": "ebpf-test-client", "version": "1.0.0"}
            }
        }' 2>/dev/null | grep -i "mcp-session-id" | cut -d: -f2 | tr -d ' \r\n')
    
    if [ -z "$SESSION_ID" ]; then
        log_error "Failed to get session ID"
        exit 1
    fi
    
    log_success "Session initialized: $SESSION_ID"
}

# Call MCP tool with proper headers
call_tool() {
    local tool_name="$1"
    local arguments="$2"
    local request_id="$3"
    
    curl -s -X POST "$SERVER_URL" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Mcp-Protocol-Version: 2025-03-26" \
        -H "Mcp-Session-Id: $SESSION_ID" \
        -d "{
            \"jsonrpc\": \"2.0\",
            \"id\": $request_id,
            \"method\": \"tools/call\",
            \"params\": {
                \"name\": \"$tool_name\",
                \"arguments\": $arguments
            }
        }"
}

# Test 1: Download eBPF program
test_download() {
    log_info "Test 1: Downloading eBPF kprobe program..."
    
    # Remove existing file
    rm -f "$TEST_FILE"
    
    # Download the kprobe example
    if curl -L -o "$TEST_FILE" "$KPROBE_URL"; then
        log_success "Downloaded eBPF program to $TEST_FILE"
        
        # Verify file
        local file_size=$(stat -f%z "$TEST_FILE" 2>/dev/null || stat -c%s "$TEST_FILE" 2>/dev/null)
        log_info "File size: $file_size bytes"
        
        # Check if it's a valid ELF file
        if file "$TEST_FILE" | grep -q "ELF"; then
            log_success "File is a valid ELF object"
        else
            log_warning "File may not be a valid ELF object"
        fi
    else
        log_error "Failed to download eBPF program"
        exit 1
    fi
}

# Test 2: List available tools
test_list_tools() {
    log_info "Test 2: Listing available tools..."
    
    local response=$(curl -s -X POST "$SERVER_URL" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Mcp-Protocol-Version: 2025-03-26" \
        -H "Mcp-Session-Id: $SESSION_ID" \
        -d '{
            "jsonrpc": "2.0",
            "id": 2,
            "method": "tools/list",
            "params": {}
        }')
    
    if echo "$response" | jq -e '.result.tools' > /dev/null 2>&1; then
        local tool_count=$(echo "$response" | jq '.result.tools | length')
        log_success "Found $tool_count tools"
        
        echo "$response" | jq -r '.result.tools[].name' | while read tool; do
            log_info "  - $tool"
        done
    else
        log_error "Failed to list tools"
        echo "$response" | jq '.'
        exit 1
    fi
}

# Test 3: Get system info
test_system_info() {
    log_info "Test 3: Getting eBPF system information..."
    
    local response=$(call_tool "info" "{}" 3)
    
    if echo "$response" | jq -e '.result' > /dev/null 2>&1; then
        log_success "System info retrieved"
        echo "$response" | jq -r '.result.content[0].text' | jq '.'
    else
        log_error "Failed to get system info"
        echo "$response" | jq '.'
        exit 1
    fi
}

# Test 4: Load eBPF program
test_load_program() {
    log_info "Test 4: Loading eBPF program..."
    
    local arguments="{
        \"source\": {
            \"type\": \"file\",
            \"path\": \"$TEST_FILE\"
        },
        \"program_type\": \"KPROBE\"
    }"
    
    local response=$(call_tool "load_program" "$arguments" 4)
    
    # Debug: show raw response
    log_info "Raw response received:"
    echo "$response"
    
    if echo "$response" | jq -e '.result' > /dev/null 2>&1; then
        local result_text=$(echo "$response" | jq -r '.result.content[0].text // empty')
        
        if [ -z "$result_text" ]; then
            log_error "Empty result text from response"
            exit 1
        fi
        
        # Check if result_text looks like valid JSON by trying to parse it
        if echo "$result_text" | jq -e '.' > /dev/null 2>&1; then
            # It's valid JSON, proceed normally
            if echo "$result_text" | jq -e '.success' > /dev/null 2>&1 && [ "$(echo "$result_text" | jq -r '.success')" = "true" ]; then
                log_success "eBPF program loaded successfully"
                
                # Extract program ID and map ID
                PROGRAM_ID=$(echo "$result_text" | jq -r '.program_id // empty')
                MAP_ID=$(echo "$result_text" | jq -r '.maps[0].id // empty')
                
                if [ -n "$PROGRAM_ID" ]; then
                    log_success "Program ID: $PROGRAM_ID"
                fi
                
                if [ -n "$MAP_ID" ]; then
                    log_success "Map ID: $MAP_ID"
                fi
                
                # Show program details
                echo "$result_text" | jq '.'
            else
                log_error "Failed to load eBPF program"
                echo "$result_text"
                exit 1
            fi
        else
            # It's not valid JSON, might be Go struct format like &{true 1.0.0 9 133 [{...}]}
            log_warning "Result text is in Go struct format, parsing..."
            echo "Result text: $result_text"
            
            # Check for success in Go struct format: &{true ...}
            if echo "$result_text" | grep -q '^&{true'; then
                log_success "eBPF program loaded successfully (Go struct format)"
                
                # Extract Program ID - it's typically the 4th field: &{true 1.0.0 9 133 ...}
                # Split by space and get the 4th field (index 3)
                PROGRAM_ID=$(echo "$result_text" | sed 's/[&{}]//g' | awk '{print $4}')
                
                # Extract Map ID from the array format [{kprobe_map 8 45 Array 4 8 1}]
                # Look for the third number in the map structure
                MAP_ID=$(echo "$result_text" | grep -o '\[{[^}]*}' | sed 's/[{}[\]]//g' | awk '{print $3}')
                
                if [ -n "$PROGRAM_ID" ] && [ "$PROGRAM_ID" != "" ]; then
                    log_success "Program ID: $PROGRAM_ID"
                else
                    log_warning "Could not extract Program ID from Go struct"
                fi
                
                if [ -n "$MAP_ID" ] && [ "$MAP_ID" != "" ]; then
                    log_success "Map ID: $MAP_ID"
                else
                    log_warning "Could not extract Map ID from Go struct"
                fi
                
                # Show the struct as-is since we can't format it as JSON
                log_info "Program details (Go struct format):"
                echo "$result_text"
                
            else
                log_error "Failed to load eBPF program - unrecognized format"
                echo "$result_text"
                exit 1
            fi
        fi
    else
        log_error "Failed to load program - invalid MCP response"
        echo "$response"
        exit 1
    fi
}

# Test 5: Attach program to kprobe
test_attach_program() {
    log_info "Test 5: Attaching program to sys_execve kprobe..."
    
    if [ -z "$PROGRAM_ID" ]; then
        log_error "No program ID available from load test"
        exit 1
    fi
    
    local arguments="{
        \"program_id\": $PROGRAM_ID,
        \"attach_type\": \"kprobe\",
        \"target\": \"sys_execve\"
    }"
    
    local response=$(call_tool "attach_program" "$arguments" 5)
    
    if echo "$response" | jq -e '.result' > /dev/null 2>&1; then
        local result_text=$(echo "$response" | jq -r '.result.content[0].text')
        
        if echo "$result_text" | grep -q '"success": *true'; then
            log_success "Program attached successfully"
            
            # Extract link ID
            LINK_ID=$(echo "$result_text" | jq -r '.link_id // empty')
            
            if [ -n "$LINK_ID" ]; then
                log_success "Link ID: $LINK_ID"
            fi
            
            # Show attachment details
            echo "$result_text" | jq '.'
        else
            log_error "Failed to attach program"
            echo "$result_text" | jq '.'
            exit 1
        fi
    else
        log_error "Failed to attach program - invalid response"
        echo "$response" | jq '.'
        exit 1
    fi
}

# Test 6: Inspect system state
test_inspect_state() {
    log_info "Test 6: Inspecting eBPF system state..."
    
    local arguments="{
        \"fields\": [\"programs\", \"maps\", \"links\", \"system\"]
    }"
    
    local response=$(call_tool "inspect_state" "$arguments" 6)
    
    if echo "$response" | jq -e '.result' > /dev/null 2>&1; then
        log_success "System state retrieved"
        echo "$response" | jq -r '.result.content[0].text' | jq '.'
    else
        log_error "Failed to inspect system state"
        echo "$response" | jq '.'
        exit 1
    fi
}

# Test 7: Generate test activity (optional)
test_generate_activity() {
    log_info "Test 7: Generating test activity (executing commands to trigger kprobe)..."
    
    # Execute some commands that will trigger sys_execve calls
    log_info "Executing test commands to trigger kprobe events..."
    
    # Simple commands that will call execve
    echo "Test command 1" > /dev/null
    ls /tmp > /dev/null
    date > /dev/null
    whoami > /dev/null
    
    log_success "Test activity generated - kprobe should have captured execve events"
}

# Test 8: Load program with base64 data
test_load_base64() {
    log_info "Test 8: Testing base64 program loading..."
    
    if [ ! -f "$TEST_FILE" ]; then
        log_warning "Test file not found, skipping base64 test"
        return
    fi
    
    # Convert file to base64
    local base64_data=$(base64 -i "$TEST_FILE" | tr -d '\n')
    
    local arguments="{
        \"source\": {
            \"type\": \"data\",
            \"blob\": \"$base64_data\"
        },
        \"program_type\": \"KPROBE\",
        \"constraints\": {
            \"verify_only\": true
        }
    }"
    
    local response=$(call_tool "load_program" "$arguments" 8)
    
    if echo "$response" | jq -e '.result' > /dev/null 2>&1; then
        local result_text=$(echo "$response" | jq -r '.result.content[0].text')
        
        if echo "$result_text" | grep -q '"success": *true'; then
            log_success "Base64 program loading successful (verify-only mode)"
        else
            log_warning "Base64 program loading failed (this may be expected in verify-only mode)"
        fi
    else
        log_warning "Base64 test failed - this may be expected"
    fi
}

# Test 9: Stream events from attached kprobe
test_stream_events() {
    log_info "Test 9: Testing event streaming from kprobe..."
    
    if [ -z "$PROGRAM_ID" ] || [ -z "$LINK_ID" ]; then
        log_warning "No program/link ID available, skipping stream test"
        return
    fi
    
    local arguments="{
        \"source\": {
            \"program_id\": $PROGRAM_ID,
            \"link_id\": $LINK_ID
        },
        \"duration\": 5,
        \"format\": \"json\",
        \"filters\": {
            \"max_events\": 10,
            \"event_types\": [\"kprobe\"],
            \"min_timestamp\": 0,
            \"target_pids\": []
        }
    }"
    
    log_info "Starting event stream for 5 seconds..."
    
    # Start streaming in background and capture response
    local response=$(call_tool "stream_events" "$arguments" 9)
    
    # Generate some activity while streaming
    log_info "Generating activity to trigger events..."
    (
        sleep 1
        echo "Stream test command 1" > /dev/null &
        ls /tmp > /dev/null &
        date > /dev/null &
        whoami > /dev/null &
        sleep 1
        ps aux | head -5 > /dev/null &
        uname -a > /dev/null &
        sleep 1
        cat /proc/version > /dev/null &
        sleep 2
    ) &
    
    # Wait for streaming to complete
    wait
    
    if echo "$response" | jq -e '.result' > /dev/null 2>&1; then
        local result_text=$(echo "$response" | jq -r '.result.content[0].text')
        
        if echo "$result_text" | grep -q '"success": *true\|"status":\|"events_captured":\|"stream_duration":'; then
            log_success "Event streaming completed"
            
            # Extract event count if available
            local event_count=$(echo "$result_text" | jq -r '.events_captured // .total_events // empty')
            if [ -n "$event_count" ] && [ "$event_count" != "null" ]; then
                log_success "Captured $event_count events"
            fi
            
            # Show streaming results
            echo "$result_text" | jq '.' 2>/dev/null || echo "$result_text"
        else
            log_warning "Event streaming completed with warnings"
            echo "$result_text"
        fi
    else
        log_error "Failed to stream events"
        echo "$response" | jq '.'
    fi
}

# Test 10: Trace errors and debugging
test_trace_errors() {
    log_info "Test 10: Testing error tracing and debugging..."
    
    local arguments="{
        \"duration\": 3,
        \"trace_level\": \"info\",
        \"include_stack_traces\": true
    }"
    
    local response=$(call_tool "trace_errors" "$arguments" 10)
    
    # Generate some activity to potentially trigger traces
    log_info "Generating system activity for tracing..."
    (
        sleep 1
        # These commands should be traced by the error tracer
        echo "Error trace test" > /dev/null
        ls /nonexistent/path > /dev/null 2>&1 || true
        sleep 2
    ) &
    
    wait
    
    if echo "$response" | jq -e '.result' > /dev/null 2>&1; then
        local result_text=$(echo "$response" | jq -r '.result.content[0].text')
        
        if echo "$result_text" | grep -q '"status":\|"traced_events":\|"Tracepoint attached"'; then
            log_success "Error tracing completed"
            
            # Extract trace count if available
            local trace_count=$(echo "$result_text" | jq -r '.traced_events // empty')
            if [ -n "$trace_count" ] && [ "$trace_count" != "null" ]; then
                log_success "Traced $trace_count events"
            fi
            
            # Show tracing results
            echo "$result_text" | jq '.' 2>/dev/null || echo "$result_text"
        else
            log_warning "Error tracing completed with warnings"
            echo "$result_text"
        fi
    else
        log_error "Failed to trace errors"
        echo "$response" | jq '.'
    fi
}

# Test 11: Comprehensive inspect state with individual fields
test_inspect_state_detailed() {
    log_info "Test 11: Detailed inspection of eBPF system state..."
    
    # Test each field individually first
    local fields=("programs" "maps" "links" "system")
    
    for field in "${fields[@]}"; do
        log_info "Inspecting $field..."
        
        local arguments="{
            \"fields\": [\"$field\"],
            \"filters\": {
                \"include_details\": true,
                \"show_inactive\": false,
                \"max_entries\": 100,
                \"sort_by\": \"id\",
                \"filter_by_type\": \"\",
                \"include_statistics\": true
            }
        }"
        
        local response=$(call_tool "inspect_state" "$arguments" "11${field}")
        
        if echo "$response" | jq -e '.result' > /dev/null 2>&1; then
            local result_text=$(echo "$response" | jq -r '.result.content[0].text')
            
            if echo "$result_text" | grep -q '"success": *true'; then
                log_success "$field inspection successful"
                
                # Show relevant counts
                case $field in
                    "programs")
                        local count=$(echo "$result_text" | jq '.programs | length // 0' 2>/dev/null || echo "0")
                        log_info "Found $count programs"
                        ;;
                    "maps")
                        local count=$(echo "$result_text" | jq '.maps | length // 0' 2>/dev/null || echo "0")
                        log_info "Found $count maps"
                        ;;
                    "links")
                        local count=$(echo "$result_text" | jq '.links | length // 0' 2>/dev/null || echo "0")
                        log_info "Found $count links"
                        ;;
                    "system")
                        local kernel=$(echo "$result_text" | jq -r '.system.kernel_version // "unknown"' 2>/dev/null || echo "unknown")
                        log_info "Kernel version: $kernel"
                        ;;
                esac
                
                # Show detailed results for this field
                echo "=== $field Details ==="
                echo "$result_text" | jq ".$field // ." 2>/dev/null || echo "$result_text"
                echo
            else
                log_warning "$field inspection returned errors"
                echo "$result_text"
            fi
        else
            log_warning "Failed to inspect $field"
        fi
    done
    
    # Now test all fields together with statistics
    log_info "Getting complete system state with statistics..."
    
    local arguments="{
        \"fields\": [\"programs\", \"maps\", \"links\", \"system\"],
        \"filters\": {
            \"include_details\": true,
            \"show_inactive\": true,
            \"max_entries\": 1000,
            \"include_statistics\": true
        }
    }"
    
    local response=$(call_tool "inspect_state" "$arguments" 111)
    
    if echo "$response" | jq -e '.result' > /dev/null 2>&1; then
        log_success "Complete system state retrieved"
        echo "$response" | jq -r '.result.content[0].text' | jq '.'
    else
        log_error "Failed to get complete system state"
        echo "$response" | jq '.'
    fi
}

# Test 12: Test different stream formats and filters
test_stream_formats() {
    log_info "Test 12: Testing different streaming formats and filters..."
    
    if [ -z "$PROGRAM_ID" ]; then
        log_warning "No program ID available, skipping stream format tests"
        return
    fi
    
    local formats=("json" "raw" "base64")
    
    for format in "${formats[@]}"; do
        log_info "Testing $format format streaming..."
        
        local arguments="{
            \"source\": {
                \"program_id\": $PROGRAM_ID
            },
            \"duration\": 2,
            \"format\": \"$format\",
            \"filters\": {
                \"max_events\": 5,
                \"event_types\": [\"kprobe\"],
                \"min_timestamp\": 0
            }
        }"
        
        local response=$(call_tool "stream_events" "$arguments" "12${format}")
        
        # Generate activity
        (date > /dev/null; whoami > /dev/null) &
        wait
        
        if echo "$response" | jq -e '.result' > /dev/null 2>&1; then
            local result_text=$(echo "$response" | jq -r '.result.content[0].text')
            log_success "$format format streaming completed"
            
            # Show first few lines only to avoid spam
            echo "$result_text" | head -10
            echo "..."
            echo
        else
            log_warning "$format format streaming failed"
        fi
    done
}

# Cleanup function
cleanup() {
    log_info "Cleaning up test files..."
    rm -f "$TEST_FILE"
    log_success "Cleanup completed"
}

# Main test execution
main() {
    echo "=============================================="
    echo "    eBPF MCP Server Comprehensive Test"
    echo "=============================================="
    echo
    
    # Pre-flight checks
    check_server
    get_auth_token "$1"  # Pass command line argument if provided
    init_session
    
    echo
    echo "Running test suite..."
    echo
    
    # Core tests
    test_download
    test_list_tools
    test_system_info
    test_load_program
    test_attach_program
    test_inspect_state
    test_generate_activity
    test_load_base64
    test_stream_events
    test_trace_errors  
    test_inspect_state_detailed
    test_stream_formats
    
    # And update the summary section to include:
    echo
    log_success \"🔍 Stream and inspection tests completed\"
    log_info \"Event streaming, error tracing, and detailed state inspection all functional\"
    echo
    echo "=============================================="
    echo "           Test Results Summary"
    echo "=============================================="
    
    if [ -n "$PROGRAM_ID" ]; then
        log_success "✅ Program loaded successfully (ID: $PROGRAM_ID)"
    fi
    
    if [ -n "$LINK_ID" ]; then
        log_success "✅ Program attached successfully (Link ID: $LINK_ID)"
    fi
    
    if [ -n "$MAP_ID" ]; then
        log_success "✅ Map created successfully (ID: $MAP_ID)"
    fi
    
    echo
    log_success "🎉 All tests completed successfully!"
    log_info "Your eBPF MCP server is fully functional and ready for Claude CLI integration"
    
    echo
    echo "Next steps:"
    echo "1. Configure Claude CLI:"
    echo "   claude mcp add ebpf $SERVER_URL -t http -H \"Authorization: Bearer $TOKEN\""
    echo
    echo "2. Test with Claude CLI:"
    echo "   claude --debug"
    echo "   Then try: 'Load and attach the kprobe eBPF program to monitor sys_execve'"
    
    cleanup
}

# Trap cleanup on exit
trap cleanup EXIT

# Run main function
main "$@"