#!/bin/bash
# install.sh - One-click installer for eBPF MCP Server

set -e

# Configuration
REPO="sameehj/ebpf-mcp"  # Update this with your actual repo
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="ebpf-mcp-server"
SERVICE_NAME="ebpf-mcp"

# default port - 8080
# custom port can be set using --port/-p option
CUSTOM_PORT="8080"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

log_header() {
    echo -e "${PURPLE}$1${NC}"
}


# parse --port/-p and optional version
while [[ $# -gt 0 ]]; do
    case "$1" in
        --port|-p)
            CUSTOM_PORT="$2"
            shift 2
            ;;
        --uninstall|--version|--help|-h)
            break
            ;;
        *)
            VERSION_ARG="$1"
            shift
            ;;
    esac
done

# Check if running as root
check_root() {
    if [[ $EUID -ne 0 ]]; then
        log_error "This script must be run as root (use sudo)"
        exit 1
    fi
}

# Detect architecture
detect_arch() {
    local arch=$(uname -m)
    case $arch in
        x86_64)
            echo "amd64"
            ;;
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            log_error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac
}

# Detect OS
detect_os() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    case $os in
        linux)
            echo "linux"
            ;;
        darwin)
            echo "darwin"
            ;;
        *)
            log_error "Unsupported OS: $os"
            exit 1
            ;;
    esac
}

# Get latest release version
get_latest_version() {
    if command -v curl &> /dev/null; then
        curl -s https://api.github.com/repos/${REPO}/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget &> /dev/null; then
        wget -qO- https://api.github.com/repos/${REPO}/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        log_error "Neither curl nor wget found. Please install one of them."
        exit 1
    fi
}

# Download and install binary
install_binary() {
    local os=$(detect_os)
    local arch=$(detect_arch)
    local version=${1:-$(get_latest_version)}
    
    if [[ -z "$version" ]]; then
        log_error "Could not determine latest version"
        exit 1
    fi
    
    local binary_name="${BINARY_NAME}-${os}-${arch}"
    local tarball_name="${binary_name}.tar.gz"
    local download_url="https://github.com/${REPO}/releases/download/${version}/${tarball_name}"
    
    log_info "Downloading eBPF MCP Server ${version} for ${os}-${arch}..."
    
    # Download to temporary location
    local temp_dir="/tmp/ebpf-mcp-install"
    mkdir -p "$temp_dir"
    cd "$temp_dir"
    
    if command -v curl &> /dev/null; then
        curl -fsSL -o "$tarball_name" "$download_url"
    elif command -v wget &> /dev/null; then
        wget -q -O "$tarball_name" "$download_url"
    fi
    
    # Verify download
    if [[ ! -f "$tarball_name" ]]; then
        log_error "Failed to download binary from $download_url"
        exit 1
    fi
    
    # Extract and install
    log_info "Extracting and installing..."
    tar -xzf "$tarball_name"
    
    if [[ ! -f "$binary_name" ]]; then
        log_error "Binary not found in tarball"
        exit 1
    fi
    
    chmod +x "$binary_name"
    mv "$binary_name" "${INSTALL_DIR}/${BINARY_NAME}"
    
    # Cleanup
    cd /
    rm -rf "$temp_dir"
    
    log_success "Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Generate systemd service
create_service() {
    local token=$(openssl rand -hex 32 2>/dev/null || xxd -l 32 -p /dev/urandom | tr -d '\n')
    local service_file="/etc/systemd/system/${SERVICE_NAME}.service"
    
    log_info "Creating systemd service..."

    # detect if port flag is supported
    local exec_opts="-t http"
    if "${INSTALL_DIR}/${BINARY_NAME}" --help 2>&1 | grep -q -- '-port'; then
        log_info "Binary supports custom port. Setting to ${CUSTOM_PORT}"
        exec_opts="${exec_opts} --port ${CUSTOM_PORT}"
    else
        CUSTOM_PORT="8080"
        log_warning "Binary does not support custom port. Defaulting to ${CUSTOM_PORT}"
    fi 

    cat > "$service_file" << EOF
[Unit]
Description=eBPF MCP Server
Documentation=https://github.com/${REPO}
After=network.target
Wants=network.target

[Service]
Type=simple
User=root
Group=root
ExecStart=${INSTALL_DIR}/${BINARY_NAME} ${exec_opts}
Environment=MCP_AUTH_TOKEN=${token}
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=ebpf-mcp

# Security settings
NoNewPrivileges=false
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/tmp /var/log
PrivateTmp=false

[Install]
WantedBy=multi-user.target
EOF

    # Save token for user reference
    echo "$token" > "/etc/ebpf-mcp-token"
    chmod 600 "/etc/ebpf-mcp-token"
    
    log_success "Created systemd service: $service_file"
    log_info "Auth token saved to: /etc/ebpf-mcp-token"
    
    # Reload systemd
    systemctl daemon-reload
    log_success "Systemd configuration reloaded"
}

# Setup Claude CLI integration
setup_claude() {
    local token=$(cat /etc/ebpf-mcp-token 2>/dev/null || echo "TOKEN_NOT_FOUND")

    echo
    log_header "🎯 Claude CLI Integration Setup"
    echo "================================="
    echo
    echo "1. Start the eBPF MCP Server:"
    echo "   sudo systemctl start ebpf-mcp"
    echo "   sudo systemctl enable ebpf-mcp"
    echo
    echo "2. Add to Claude CLI:"
    echo "   claude mcp add ebpf http://localhost:${CUSTOM_PORT}/mcp -t http -H \"Authorization: Bearer $token\""
    echo
    echo "3. Test the integration:"
    echo "   claude --debug"
    echo "   Then try: 'What eBPF capabilities does this system have?'"
    echo
    echo "4. Check service status:"
    echo "   sudo systemctl status ebpf-mcp"
    echo "   sudo journalctl -u ebpf-mcp -f"
    echo
    echo "🔑 Your auth token: $token"
    echo "💾 Token saved to: /etc/ebpf-mcp-token"
    echo
    echo "📚 Documentation: https://github.com/${REPO}/tree/main/docs"
    echo "🚀 Quick Start: https://github.com/${REPO}/blob/main/docs/TUTORIAL.md"
    echo
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check kernel version
    local kernel_version=$(uname -r | cut -d. -f1-2)
    local major=$(echo $kernel_version | cut -d. -f1)
    local minor=$(echo $kernel_version | cut -d. -f2)
    
    if [[ $major -lt 4 ]] || [[ $major -eq 4 && $minor -lt 18 ]]; then
        log_warning "Kernel version $kernel_version detected. eBPF requires 4.18+"
        log_warning "Some features may not work properly"
    else
        log_success "Kernel version $kernel_version supports eBPF"
    fi
    
    # Check for eBPF filesystem
    if mount | grep -q bpf; then
        log_success "BPF filesystem is mounted"
    else
        log_warning "BPF filesystem not mounted. Mounting..."
        mount -t bpf bpf /sys/fs/bpf 2>/dev/null || log_warning "Failed to mount BPF filesystem"
    fi
    
    # Check for systemd
    if command -v systemctl &> /dev/null; then
        log_success "Systemd detected"
    else
        log_error "Systemd is required but not found"
        exit 1
    fi

    # Check for required tools
    local missing_tools=""
    for tool in openssl xxd; do
        if ! command -v $tool &> /dev/null; then
            missing_tools="$missing_tools $tool"
        fi
    done
    
    if [[ -n "$missing_tools" ]]; then
        log_warning "Missing tools:$missing_tools (token generation may be affected)"
    fi
}

# Uninstall function
uninstall() {
    log_info "Uninstalling eBPF MCP Server..."
    
    # Stop and disable service
    systemctl stop ebpf-mcp 2>/dev/null || true
    systemctl disable ebpf-mcp 2>/dev/null || true
    
    # Remove files
    rm -f /etc/systemd/system/ebpf-mcp.service
    rm -f "${INSTALL_DIR}/${BINARY_NAME}"
    rm -f /etc/ebpf-mcp-token
    
    # Reload systemd
    systemctl daemon-reload 2>/dev/null || true
    
    log_success "Uninstalled successfully"
}

# Show version
show_version() {
    if [[ -f "${INSTALL_DIR}/${BINARY_NAME}" ]]; then
        "${INSTALL_DIR}/${BINARY_NAME}" --version 2>/dev/null || echo "Version info not available"
    else
        echo "eBPF MCP Server not installed"
    fi
}

# Main installation flow
main() {
    echo
    log_header "🚀 eBPF MCP Server Installer"
    echo "================================="    
    echo
    
    check_root
    check_prerequisites
    install_binary "$VERSION_ARG"
    create_service
    setup_claude
    
    echo
    log_success "🎉 Installation complete!"
    echo
    log_info "Next steps:"
    echo "  1. sudo systemctl start ebpf-mcp"
    echo "  2. sudo systemctl enable ebpf-mcp"
    echo "  3. Add to Claude CLI using the command shown above"
    echo
    log_info "For help: ${BINARY_NAME} --help"
    log_info "Documentation: https://github.com/${REPO}"
}

# Handle command line arguments
case "${1:-}" in
    --uninstall)
        check_root
        uninstall
        ;;
    --version)
        show_version
        ;;
    --help|-h)
        echo "eBPF MCP Server Installer"
        echo
        echo "Usage: $0 [OPTIONS] [VERSION]"
        echo
        echo "Options:"
        echo "  --uninstall       Remove eBPF MCP Server"
        echo "  --version         Show installed version"
        echo "  --help, -h        Show this help message"
        echo "  --port, -p <PORT> Run server on custom port (default: 8080)"
        echo
        echo "Examples:"
        echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | sudo bash"
        echo "  sudo $0 v1.0.0"
        echo "  sudo $0 --port 9090"
        echo "  sudo $0 --uninstall"
        ;;
    *)
        main "$@"
        ;;
esac