
## 🚀 Running eBPF-MCP-server with Cursor IDE

This guide explains how to run the ebpf-mcp-server and connect it to Cursor IDE for tool-based AI interaction.

## ✅ Prerequisites

- Go installed on your system or virtual machine.
- Cursor IDE installed.
- Terminal access to build and run Go applications.

> ⚠️ **Note:** If you're using **VirtualBox**, make sure port `8080` is forwarded to your host so that Cursor can access the server.  
> You can do this under:  
> `VM Settings → Network → Port Forwarding`  
> Example configuration:
> -  **Protocol:** TCP  
> -  **Host Port:** 8080  
> -  **Guest IP:** Your VM’s IP  
> -  **Guest Port:** 8080

## 🏃 Run the Server

1. **Clone the repository**
```bash
git clone https://github.com/sameehj/ebpf-mcp.git
cd ebpf-mcp
```
2. **Build the server**
```bash
make build
```
3. **Run the server in debug HTTP mode**
```bash
sudo ./bin/ebpf-mcp-server --debug --transport=http
```
After running, the server will print and authentication token in the logs.
## ⚙️ Configure MCP in Cursor IDE
1. Open **Cursor IDE**
2. Go to:
`Settings -> Tools & Integration -> MCP Tools`
3. Add a new MCP server with the following configuration:
```json
{
  "mcpServers": {
    "ebpf-mcp": {
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer <your-token>"
      }
    }
  }
}
```
4. Switch Cursor to an AI model that supports MCP & Tools
`e.g., Clause 3.5 Sonnet`
## 💬 Example Prompt
Try asking Cursor:
```
Can you get the system info and kernel version?
```
Cursor will now use the `ebpf-mcp` server and tools like `info` to provide accurate system-level information.
## 📝 Notes
- The setup works locally on in a VM as long as the server is reachable on port `8080`.
- Make sure the server is running before using MCP in Cursor.
- Make sure to choose a valid AI model that supports MCP & Tools.
## Need Help with VirtualBox?
This guide was written and tested using VirtualBox.
If you're setting this up on VirtualBox and need help with networking, SSH, or port forwarding - feel free to leave an issue.
## 🤝 Contribution
 🛠 Contributions, issues, and PRs welcome!