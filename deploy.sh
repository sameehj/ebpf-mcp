#!/bin/bash
set -e

# Colors for prettier output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}eBPF-MCP Deployment Script${NC}"
echo "This script will install and deploy the eBPF-MCP service."

# Check if running as root
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Please run as root or with sudo${NC}"
  exit 1
fi

# Check kernel version
KERNEL_VERSION=$(uname -r | cut -d. -f1,2)
KERNEL_MAJOR=$(echo $KERNEL_VERSION | cut -d. -f1)
KERNEL_MINOR=$(echo $KERNEL_VERSION | cut -d. -f2)

if [ "$KERNEL_MAJOR" -lt 5 ] || ([ "$KERNEL_MAJOR" -eq 5 ] && [ "$KERNEL_MINOR" -lt 8 ]); then
  echo -e "${YELLOW}Warning: Your kernel version ($KERNEL_VERSION) may have limited eBPF support.${NC}"
  echo -e "${YELLOW}For best results, kernel 5.8+ is recommended.${NC}"
  read -p "Continue anyway? (y/n) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
  fi
fi

echo -e "${GREEN}Installing system dependencies...${NC}"
if [ -f /etc/debian_version ]; then
  # Debian/Ubuntu
  apt-get update
  apt-get install -y build-essential linux-headers-$(uname -r) \
    python3-dev python3-pip git clang llvm libelf-dev \
    python3-bcc bpfcc-tools libbpfcc libbpfcc-dev nginx
elif [ -f /etc/redhat-release ]; then
  # CentOS/RHEL/Fedora
  yum install -y gcc gcc-c++ make kernel-devel elfutils-libelf-devel \
    python3-devel python3-pip git clang llvm bcc bcc-devel bcc-tools nginx
elif [ -f /etc/os-release ]; then
  # Check for Amazon Linux
  source /etc/os-release
  if [[ "$ID" == "amzn" ]]; then
    # Amazon Linux
    yum install -y gcc gcc-c++ make kernel-devel elfutils-libelf-devel \
      python3-devel python3-pip git clang llvm nginx
    # Install BCC from source on Amazon Linux
    echo -e "${YELLOW}Installing BCC from source (this may take a while)...${NC}"
    yum install -y bison cmake3 flex git libstdc++-static python3-netaddr
    git clone https://github.com/iovisor/bcc.git /tmp/bcc
    mkdir /tmp/bcc/build
    cd /tmp/bcc/build
    cmake3 -DCMAKE_INSTALL_PREFIX=/usr ..
    make -j$(nproc)
    make install
    cd -
  else
    echo -e "${RED}Unsupported distribution. Please install dependencies manually.${NC}"
    exit 1
  fi
else
  echo -e "${RED}Unsupported distribution. Please install dependencies manually.${NC}"
  exit 1
fi

echo -e "${GREEN}Installing Python dependencies...${NC}"
pip3 install --upgrade pip
pip3 install bcc fastapi uvicorn pydantic click pyyaml requests

# Install the eBPF-MCP package
echo -e "${GREEN}Installing eBPF-MCP package...${NC}"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
pip3 install -e "$SCRIPT_DIR"

# Create necessary directories
echo -e "${GREEN}Creating directories...${NC}"
mkdir -p /usr/local/share/ebpf-mcp/bpf_programs
mkdir -p /var/lib/ebpf-mcp/mcp/.well-known/mcp/{maps,traces,tools}

# Copy BPF programs
echo -e "${GREEN}Copying BPF programs...${NC}"
cp "$SCRIPT_DIR/bpf_programs/"* /usr/local/share/ebpf-mcp/bpf_programs/

# Create config file
echo -e "${GREEN}Creating configuration...${NC}"
mkdir -p /etc/ebpf-mcp
cat > /etc/ebpf-mcp/config.yaml << 'EOF'
server:
  host: 0.0.0.0
  port: 8080
  debug: false
ebpf:
  programs_dir: /usr/local/share/ebpf-mcp/bpf_programs
mcp:
  context_dir: /var/lib/ebpf-mcp/mcp
logging:
  level: info
  file: /var/log/ebpf-mcp.log
EOF

# Create systemd service
echo -e "${GREEN}Setting up systemd service...${NC}"
cat > /etc/systemd/system/ebpf-mcp.service << 'EOF'
[Unit]
Description=eBPF Model Context Protocol Service
After=network.target

[Service]
ExecStart=/usr/local/bin/ebpf-mcp serve
Restart=on-failure
User=root
Group=root
Environment=PYTHONUNBUFFERED=1

[Install]
WantedBy=multi-user.target
EOF

# Set up Nginx as a reverse proxy
echo -e "${GREEN}Configuring Nginx...${NC}"
cat > /etc/nginx/sites-available/ebpf-mcp << 'EOF'
server {
    listen 80;
    server_name _;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
EOF

# Enable the Nginx site
if [ -d /etc/nginx/sites-enabled ]; then
  ln -sf /etc/nginx/sites-available/ebpf-mcp /etc/nginx/sites-enabled/
  rm -f /etc/nginx/sites-enabled/default
else
  # For distributions that don't use sites-enabled directory
  mv /etc/nginx/sites-available/ebpf-mcp /etc/nginx/conf.d/ebpf-mcp.conf
fi

# Create welcome page
mkdir -p /var/www/html
cat > /var/www/html/index.html << 'EOF'
<!DOCTYPE html>
<html>
<head>
    <title>eBPF-MCP Deployment</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; padding: 20px; max-width: 800px; margin: 0 auto; }
        h1 { color: #333; border-bottom: 1px solid #eee; padding-bottom: 10px; }
        .info { background: #f8f8f8; padding: 15px; border-radius: 5px; border-left: 4px solid #5cb85c; }
        code { background: #f1f1f1; padding: 2px 5px; border-radius: 3px; font-family: monospace; }
        pre { background: #f1f1f1; padding: 10px; border-radius: 5px; overflow-x: auto; }
    </style>
</head>
<body>
    <h1>eBPF-MCP Deployment Successful!</h1>
    <div class="info">
        <p><strong>Server Status:</strong> <span id="status">Checking...</span></p>
        <p><strong>Server IP:</strong> <span id="ip">Detecting...</span></p>
    </div>
    
    <h2>MCP Endpoints</h2>
    <ul>
        <li>MCP Metadata: <a href="/.well-known