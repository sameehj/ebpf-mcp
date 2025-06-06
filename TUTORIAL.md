## 📘 TUTORIAL.md

````markdown
# 🐝 ebpf-mcp: Getting Started Tutorial

This guide walks you through setting up, running, and testing the `ebpf-mcp` server — a lightweight, AI-compatible eBPF server that speaks the [Model Context Protocol (MCP)](https://github.com/modelcontextprotocol/spec).

---

## ⚙️ Prerequisites

Before you begin:

- Go 1.20+ installed
- Linux or macOS with terminal access
- Optional: `curl` for testing JSON-RPC requests

---

## 📦 Step 1: Clone the Project

```bash
git clone https://github.com/sameehj/ebpf-mcp.git
cd ebpf-mcp
````

---

## 🛠️ Step 2: Build the Server

### Standard build:

```bash
make build
```

### Verbose (debug) build:

```bash
make build DEBUG=1
```

This compiles the server binary to `./bin/ebpf-mcp-server`.

---

## 🚀 Step 3: Run the Server

```bash
make run
```

Or directly:

```bash
./bin/ebpf-mcp-server
```

The server starts an HTTP endpoint at:

```
http://localhost:8080/rpc
```

---

## 🧪 Step 4: Test `tools/list` via `curl`

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "id": 1
  }'
```

✅ Expected response:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [
      {
        "name": "map_dump",
        "description": "Returns contents of a BPF map"
      }
    ]
  }
}
```

---

## 📬 Step 5: Test `tools/call` (Stub)

```bash
curl -X POST http://localhost:8080/rpc \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "tool": "map_dump",
      "input": {
        "map_name": "dummy_map"
      }
    },
    "id": 2
  }'
```

✅ Expected stub response:

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Executed tool: map_dump"
      }
    ]
  }
}
```

---

## 🧼 Step 6: Clean the Build

```bash
make clean
```

---

## 🔍 Next Steps

* Implement real `map_dump` logic using `github.com/cilium/ebpf`
* Add tools like `trace`, `deploy`, and `info`
* Integrate with Prometheus, seccomp, and namespace isolation
* Expand MCP method support: `resources/list`, `prompts/ask`

---

## 🧠 Need Help?

Feel free to open an issue or discussion on [GitHub](https://github.com/sameehj/ebpf-mcp) or contribute new tools and improvements.

Let’s build the AI-native observability layer for Linux 🐧🚀

```