# eBPF-MCP: One-Click AWS Deployment

[![Deploy to AWS](https://img.shields.io/badge/Deploy%20to-AWS-%23FF9900?style=for-the-badge&logo=amazon-aws)](https://console.aws.amazon.com/cloudformation/home#/stacks/create/review?templateURL=https://ebpf-mcp-templates.s3.amazonaws.com/ebpf-mcp-stack.yaml)

The eBPF Model Context Protocol (eBPF-MCP) bridges the gap between Linux kernel observability and LLMs (Large Language Models). This project provides a one-click deployment solution for AWS, allowing you to quickly set up an eBPF-MCP server.

## What is eBPF-MCP?

eBPF-MCP provides a standardized way for LLMs to:

- Access real-time system state through eBPF
- Monitor system calls, network traffic, and performance metrics
- Execute safe eBPF tools for system analysis
- Reason about system behavior with full context

The Model Context Protocol (MCP) exposes this data through a structured interface at `/.well-known/mcp/` that LLMs can access and understand.

## One-Click Deployment

To deploy eBPF-MCP on AWS:

1. Click the "Deploy to AWS" button above
2. Log in to your AWS account if prompted
3. Select a key pair for SSH access
4. Choose an instance type (t3.medium recommended)
5. Adjust the IP address range for access if needed
6. Click "Create stack"

The deployment takes about 5-10 minutes. When complete, you'll see:
- **WebsiteURL**: URL for the web interface
- **MCPURL**: URL for the MCP endpoint to use with LLMs
- **SSHCommand**: Command to SSH into the instance

## Using with LLMs

To use your eBPF-MCP server with an LLM like Claude:

1. Deploy the stack and get your MCP URL (e.g., `http://12.34.56.78/.well-known/mcp/`)
2. In your conversation with the LLM, provide the MCP endpoint:

```
You can access system information from my eBPF-MCP server at: http://12.34.56.78/.well-known/mcp/
```

The LLM can then:
- Access metadata: `/.well-known/mcp/metadata.json`
- See available tools: `/.well-known/mcp/tools.json`
- Read system summary: `/.well-known/mcp/llms.txt` 
- Use tools via the API: `/api/tools/{tool_name}`

## Available Tools

The default installation includes:

- **syscall_trace**: Monitor system calls for a process
- **network_monitor**: Analyze network connections

## Architecture

```
+------------------------+
|      LLM Agent         |
| (Claude, GPT, etc.)    |
+-----------+------------+
            |
   MCP Protocol (MCP API)
            |
+-----------v------------+
|     eBPF-MCP Server     |
|  Exposes context files  |
|  and safe tool actions  |
+-----------+-------------+
            |
+------------------------+------------------------+
|                                                 |
+-------------------+                         +-------------------+
| eBPF Program Hooks|                         | System Context    |
| - Syscalls        |                         | - logs/           |
| - Net packets     |                         | - metrics/        |
| - BPF maps        |                         | - trace outputs/  |
+-------------------+                         +-------------------+
```

## Manual Installation

If you prefer to install manually:

```bash
git clone https://github.com/ebpf-mcp/ebpf-mcp.git
cd ebpf-mcp
chmod +x deploy.sh
sudo ./deploy.sh
```

## Security Considerations

- The deployment opens ports 22 (SSH), 80 (HTTP), and 8080 to the specified IP range
- For production use, restrict the access CIDR to your IP address
- Consider using HTTPS for production deployments
- eBPF programs run with root privileges, so only trusted tools should be used

## License

This project is licensed under the GNU General Public License v2.0.