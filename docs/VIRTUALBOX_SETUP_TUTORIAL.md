
## 🧪 Running eBPF-MCP on VirtualBox

This short guide shows how to get started with eBPF-MCP using VirtualBox and an Ubuntu VM.

> ✅ _This setup has been tested on macOS (Apple M2 Pro) running Ubuntu 25.04 (ARM64, CLI only) inside VirtualBox._

## 🖥️ System Setup (Example)
- Host OS: macOS (Apple M2 Pro)
- VM Tool: VirtualBox (ARM build)
- Guest OS: Ubuntu 25.04 (ARM64, CLI only)
- Architecture: ARM64
## 📦 Step 1: Install VirtualBox
Follow this tutorial to install VirtualBox and set up Ubuntu 
> _This tutorial shows setup for Apple Silicon with Ubutu ARM 64, you may follow the same steps but with different OS._

📺 [YouTube Setup Guide (macOS ARM)](https://www.youtube.com/watch?v=LjL_N0OZxvY)

Make sure you:
- Download the correct platform package depending on your OS `(e.g., macOS / Apple Silicon hosts)`
- Create a VM with Ubuntu `e.g., 25.04 ARM64 ISO for macOS` - you may install other versions depending on your setup.
## 🔧 Step 2: Boot the VM and Update Packages
```bash
sudo apt update && sudo apt upgrade -y
```
## ⚙️ Step 3: Install eBPF-MCP
Run the official one-liner from the [Quick Start](https://github.com/sameehj/ebpf-mcp?tab=readme-ov-file#-quick-start):
```bash
curl -fsSL https://raw.githubusercontent.com/sameehj/ebpf-mcp/main/install.sh | sudo bash
```
✅ This installs everything needed, sets up the systemd service, and starts the eBPF MCP server automatically.
## 🚀 Step 4: Start the MCP Server
```bash
sudo systemctl start ebpf-mcp
sudo systemctl enable ebpf-mcp
sudo systemctl status ebpf-mcp
```
The server runs on http://localhost:8080 and provides rich eBPF monitoring via an MCP-compatible API.
## ✅ Verification
Check the server is running with:
```bash
curl -H "Authorization: Bearer $(cat /etc/ebpf-mcp-token)" http://localhost:8080/.well-known/mcp/metadata.json
```
## 💡 Notes
- Works seamlessly with Ubuntu 25.04 (ARM64) CLI in VirtualBox.

## 🙌 Contribution
 🛠 Contributions, issues, and PRs welcome!

If you've tested this setup on another architecture or VM solution (e.g., VMware, QEMU, x86_64), feel free to open a PR and expand this section.