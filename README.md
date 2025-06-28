# 🐝 ebpf-mcp: AI-Compatible eBPF Control via Model Context Protocol

> A secure, minimal, and schema-enforced MCP server for eBPF — purpose-built for AI integration, kernel introspection, and automation.

[![Version](https://img.shields.io/github/v/release/sameehj/ebpf-mcp?label=version)](https://github.com/sameehj/ebpf-mcp/releases)
[![MCP Compatible](https://img.shields.io/badge/MCP-Compatible-orange)](https://modelcontextprotocol.io)
[![eBPF Support](https://img.shields.io/badge/eBPF-Linux%205.8%2B-green)](https://ebpf.io)
[![License: GPL v2 (eBPF)](https://img.shields.io/badge/License-GPL%20v2-blue.svg)](https://www.gnu.org/licenses/old-licenses/gpl-2.0.html)
[![License: Apache 2.0 (Core)](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://www.apache.org/licenses/LICENSE-2.0)

---

## 🧠 What Is This?

`ebpf-mcp` is a secure **Model Context Protocol (MCP)** server that exposes **a minimal set of structured tools** to interact with eBPF — optimized for safe AI control, automation agents, and human operators.

It enables **loading, attaching, introspecting, and streaming** eBPF programs — all through strict JSON Schema contracts validated at runtime. No REST APIs, no shell escapes, and no bpftool wrappers.

---

## 🚀 Quick Start

### 📦 One-liner Installation

```bash
# Install ebpf-mcp server
curl -fsSL https://raw.githubusercontent.com/sameehj/ebpf-mcp/main/install.sh | sudo bash

# Start the service (runs on port 8080 by default)
sudo systemctl start ebpf-mcp
sudo systemctl enable ebpf-mcp

# Get your auth token
cat /etc/ebpf-mcp-token

# Check service status
sudo systemctl status ebpf-mcp

# View logs if needed
sudo journalctl -u ebpf-mcp -f
```

**For air-gapped or development environments:**
```bash
git clone https://github.com/sameehj/ebpf-mcp.git
cd ebpf-mcp
sudo ./install.sh v1.0.2
```

### 🧪 Test the Installation

```bash
# Run the complete test suite
cd scripts/
chmod +x test-ebpf-mcp-server.sh
./test-ebpf-mcp-server.sh <your-token>
```

If no token is provided, the script will prompt for it interactively.

---

## 🤖 Claude CLI Integration

Once installed, connect Claude to your eBPF server (runs on port 8080):

```bash
# Add MCP server to Claude CLI
claude mcp add ebpf http://localhost:8080/mcp \
  -t http \
  -H "Authorization: Bearer $(cat /etc/ebpf-mcp-token)"

# Start Claude with eBPF tools
claude --debug

# Optional: Test with MCP Inspector (requires Node.js)
npx @modelcontextprotocol/inspector http://localhost:8080/mcp
```

**Example prompts:**
- `> Get system info and kernel version`
- `> Load and attach a kprobe program to monitor sys_execve`
- `> Show me all active eBPF programs and their types`
- `> Stream events from ringbuffer maps for 10 seconds`
- `> Trace kernel errors for the next 5 seconds`

---

## 📥 Install Options

| Method | Command | Use Case |
|--------|---------|----------|
| **One-liner** | `curl ... \| sudo bash` | Production systems |
| **Manual** | `git clone && sudo ./install.sh` | Development/air-gapped |
| **Build from source** | `make build` | Custom modifications |
| **Docker** | *Coming soon* | Containerized environments |

---

## 🔧 Minimal Toolset

Each tool is designed to be schema-validatable, AI-orchestrable, and safe-by-default. They cover 80%+ of real-world observability and control workflows.

| Tool Name        | Status | Description                                     | Capabilities Required                          |
| ---------------- | ------ | ----------------------------------------------- | ---------------------------------------------- |
| `info`           | ✅      | System introspection: kernel, arch, BTF        | `CAP_BPF` or none (read-only)                  |
| `load_program`   | ✅      | Load and validate `.o` files (CO-RE supported)  | `CAP_BPF` or `CAP_SYS_ADMIN`                   |
| `attach_program` | ✅      | Attach program to XDP, kprobe, tracepoint hooks | Depends on type (e.g. `CAP_NET_ADMIN` for XDP) |
| `inspect_state`  | ✅      | List programs, maps, links, and tool metadata   | `CAP_BPF` (read-only)                          |
| `stream_events`  | ✅      | Stream events from ringbuf/perfbuf maps         | `CAP_BPF` (read-only)                          |
| `trace_errors`   | ✅      | Monitor kernel tracepoints for error conditions | `CAP_BPF` (read-only)                          |

> **All tools return structured JSON output** — AI-ready, streaming-compatible, and schema-validated.

> 🔍 See [`docs/TOOL_SPECS.md`](./docs/TOOL_SPECS.md) for full schema definitions.

---

## 🚀 What You Can Do

* ✅ Query kernel version, architecture, and BTF availability
* ✅ Load programs from disk or inline base64 with optional BTF
* ✅ Attach to live systems with type-safe constraints
* ✅ Inspect pinned objects, kernel version, verifier state
* ✅ Stream real-time events with filtering by pid/comm/cpu
* ✅ Trace kernel errors and system anomalies
* ✅ Discover available tools and their schemas
* ✅ Integrate with Claude, Ollama, or MCP-compatible clients

---

## 🛡️ Security Model

| Layer             | Controls                                 |
| ----------------- | ---------------------------------------- |
| eBPF execution    | Kernel verifier + resource caps          |
| Filesystem        | No shell, no exec, path-validated        |
| Runtime isolation | Session-scoped cleanup, strict inputs    |
| AI safety         | Capability-aware schemas + output limits |
| Authentication    | Bearer token + HTTPS ready              |

🧼 All resources are automatically cleaned up when a client disconnects (no manual unload/detach required unless pinned).

---

## 📦 Project Structure

```
.
├── cmd/              # MCP server + CLI client
├── internal/         # Core logic: eBPF, tools, kernel adapters
├── pkg/types/        # JSON schema bindings + shared types
├── docs/             # Tool specs, design notes, schemas
├── scripts/          # Install script + test suite
└── schemas/          # JSON Schema files for each tool
```

---

## 🧠 Advanced Design Notes

### ✅ Lifecycle Management

* 🔒 **No manual detach**: Links are closed automatically unless pinned
* 🧹 **Auto cleanup**: FDs and memory are released on disconnect
* 📎 **Pinning**: Optional pin paths (`/sys/fs/bpf/...`) for maps/programs/links

### 🤖 AI Tooling Compatibility

* All tools are **strictly typed** with published schemas and return **structured JSON output**
* **AI-ready**: No parsing required — direct integration with language models
* **Streaming-compatible**: Real-time data flows for observability workflows
* Responses include:
  * `tool_version`
  * `verifier_log` (for debugging)
  * Structured `error` with `context`

### 🔗 Extensibility

Future optional tools:
* `pin_object` / `unpin_object`
* `detach_link`
* `map_batch_op`

These are omitted from the default for security and simplicity.

---

## 📚 References

* [Linux Kernel eBPF Docs](https://docs.kernel.org/bpf/)
* [Model Context Protocol](https://modelcontextprotocol.io)
* [MCP Inspector Tool](https://github.com/modelcontextprotocol/inspector)
* [JSON Schema Spec (2020-12)](https://json-schema.org/)
* [eBPF Security Best Practices](https://ebpf.io/security/)
* [Cilium for Kubernetes Observability](https://cilium.io/)

🧪 See [`scripts/test-ebpf-mcp-server.sh`](./scripts/test-ebpf-mcp-server.sh) for full validation suite.

**Basic Architecture:**
```
Claude / Ollama / AI Client
          ↓
     MCP JSON-RPC
          ↓
   ebpf-mcp server
          ↓
     Kernel APIs
```

---

## 📜 Licensing

| Component        | License    |
| ---------------- | ---------- |
| `internal/ebpf/` | GPL-2.0    |
| Everything else  | Apache-2.0 |

---

## ✉️ Contact

📬 [GitHub – sameehj/ebpf-mcp](https://github.com/sameehj/ebpf-mcp)
🛠 Contributions, issues, and PRs welcome!

---

> **Structured. Safe. Schema-native.**
> `ebpf-mcp` brings eBPF to the age of AI.