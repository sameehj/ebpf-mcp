## 🧪 Tutorial: Running `ollama-chat` with a Live eBPF MCP Server

This walkthrough shows how to:

1. Run your **eBPF MCP server** remotely on an EC2 instance
2. Tunnel the MCP port to your local machine
3. Connect `ollama-chat` to the server
4. Interact with it using **natural language**

---

### ✅ 1. Launch the eBPF MCP Server on EC2

```bash
ssh ubuntu@your-ec2-host
git clone https://github.com/sameehj/ebpf-mcp.git
cd ebpf-mcp
go build -o ./bin/ebpf-mcp-server ./cmd/ebpf-mcp
./bin/ebpf-mcp-server
```

By default, the server listens on: `http://localhost:8080/mcp`

---

### 🔒 2. Tunnel the MCP Port to Your Local Machine

In another terminal (on your local machine):

```bash
ssh -N -L 8080:localhost:8080 ubuntu@your-ec2-host
```

This forwards the remote MCP port (`8080`) to your local port (`8080`). Now your local apps can talk to the MCP server as if it’s running locally.

---

### 🧠 3. Start `ollama-chat`

Make sure [Ollama](https://ollama.com/) is running locally with your model of choice (e.g. `llama3`):

```bash
ollama run llama3
```

Then run:

```bash
go build -o ./bin/ollama-chat ./cmd/ollama-chat
./bin/ollama-chat
```

> The chat client will connect to `localhost:8080/mcp`, which is now your **remote EC2 server**.

---

### 🎤 4. Example Interactive Session

Here’s what a real run looks like:

```
🔬 Welcome to the Simple eBPF Chat!
Type 'exit' to quit, 'list tools' to see available tools.

🔗 Connecting to eBPF MCP server...
✅ Connected! Found 4 eBPF tools available.

You 🧠: show me kernel info

AI 🤖: The `info` tool provides kernel and eBPF environment details...

🔬 Calling eBPF tool: info
🔬 eBPF Tool 'info' output:
{
  "btf_enabled": true,
  "kernel_version": "6.1.134-amzn2023",
  ...
}

You 🧠: trace some syscall errors

AI 🤖: Running `trace_errors` for 2 seconds...
🔬 eBPF Tool 'trace_errors' output:
Tracepoint attached and ran for 2s
```

You can ask things like:

* `what eBPF programs are running?`
* `dump the conntrack map`
* `monitor the ingress hook`

---

### 🧠 How It Works

* Ollama generates a prompt →
* `ollama-chat` builds an MCP `tools/call` request →
* MCP server executes an eBPF tool on EC2 →
* The result is streamed back and shown locally

---

## 📌 Requirements

* Go 1.21+ on both local and remote
* Open `ssh` access to your EC2 instance
* Ollama running locally with a model like `llama3`
* MCP server accessible on port `8080` (via SSH tunnel)

---

## 🚀 Want to Go Further?

* Set up `ssh-mcp` for AI command execution over SSH
* Try deploying eBPF programs via `ebpf.deploy` tool
* Write your own tool plugin and register it under `internal/tools/`

