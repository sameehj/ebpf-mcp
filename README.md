# 🐝 ebpf-mcp: AI-Compatible eBPF Control via Model Context Protocol

> A secure, minimal, and schema-enforced MCP server for eBPF — purpose-built for AI integration, kernel introspection, and automation.

[![MCP Compatible](https://img.shields.io/badge/MCP-Compatible-orange)](https://modelcontextprotocol.io)
[![eBPF Support](https://img.shields.io/badge/eBPF-Linux%205.8%2B-green)](https://ebpf.io)
[![License: GPL v2 (eBPF)](https://img.shields.io/badge/License-GPL%20v2-blue.svg)](https://www.gnu.org/licenses/old-licenses/gpl-2.0.html)
[![License: Apache 2.0 (Core)](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://www.apache.org/licenses/LICENSE-2.0)

---

## 🧠 What Is This?

`ebpf-mcp` is a secure **Model Context Protocol (MCP)** server that exposes **a minimal set of structured tools** to interact with eBPF — optimized for safe AI control, automation agents, and human operators.

It enables **loading, attaching, introspecting, and streaming** eBPF programs — all through strict JSON Schema contracts validated at runtime. No REST APIs, no shell escapes, and no bpftool wrappers.

---

## 🔧 Minimal Toolset

Each tool is designed to be schema-validatable, AI-orchestrable, and safe-by-default. They cover 80%+ of real-world observability and control workflows.

| Tool Name        | Description                                     | Capabilities Required                          |
| ---------------- | ----------------------------------------------- | ---------------------------------------------- |
| `load_program`   | Load and validate `.o` files (CO-RE supported)  | `CAP_BPF` or `CAP_SYS_ADMIN`                   |
| `attach_program` | Attach program to XDP, kprobe, tracepoint hooks | Depends on type (e.g. `CAP_NET_ADMIN` for XDP) |
| `inspect_state`  | List programs, maps, links, and tool metadata   | `CAP_BPF` (read-only)                          |
| `stream_events`  | Stream events from ringbuf/perfbuf maps         | `CAP_BPF` (read-only)                          |

> 🔍 See [`docs/TOOL_SPECS.md`](./docs/TOOL_SPECS.md) for full schema definitions.

---

## 🚀 What You Can Do

* ✅ Load programs from disk or inline base64 with optional BTF
* ✅ Attach to live systems with type-safe constraints
* ✅ Inspect pinned objects, kernel version, verifier state
* ✅ Stream real-time events with filtering by pid/comm/cpu
* ✅ Discover available tools and their schemas
* ✅ Integrate with Claude, Ollama, or MCP-compatible clients

---

## ⚡ Quick Start

```bash
# Clone + build
git clone https://github.com/sameehj/ebpf-mcp.git
cd ebpf-mcp
make build
```

```bash
# Run locally with MCP Inspector
npx @modelcontextprotocol/inspector ./bin/ebpf-mcp-server
```

```jsonc
// ~/.config/Claude/claude_desktop_config.json
{
  "mcpServers": {
    "ebpf": {
      "command": "/absolute/path/to/ebpf-mcp-server",
      "args": ["-t", "stdio"]
    }
  }
}
```

---

## 🛡️ Security Model

| Layer             | Controls                                 |
| ----------------- | ---------------------------------------- |
| eBPF execution    | Kernel verifier + resource caps          |
| Filesystem        | No shell, no exec, path-validated        |
| Runtime isolation | Session-scoped cleanup, strict inputs    |
| AI safety         | Capability-aware schemas + output limits |

🧼 All resources are automatically cleaned up when a client disconnects (no manual unload/detach required unless pinned).

---

## 📦 Project Structure

```
.
├── cmd/              # MCP server + CLI client
├── internal/         # Core logic: eBPF, tools, kernel adapters
├── pkg/types/        # JSON schema bindings + shared types
├── docs/             # Tool specs, design notes, schemas
└── schemas/          # JSON Schema files for each tool
```

---

## 📈 Tool Spec Coverage

| Tool             | Status | Notes                                    |
| ---------------- | ------ | ---------------------------------------- |
| `load_program`   | ✅      | Supports CO-RE, verify-only mode         |
| `attach_program` | ✅      | Supports XDP, kprobe, tracepoint         |
| `inspect_state`  | ✅      | Introspects maps, programs, links, tools |
| `stream_events`  | ✅      | Streams ringbuf/perfbuf with filters     |

---

## 🧠 Advanced Design Notes

### ✅ Lifecycle Management

* 🔒 **No manual detach**: Links are closed automatically unless pinned
* 🧹 **Auto cleanup**: FDs and memory are released on disconnect
* 📎 **Pinning**: Optional pin paths (`/sys/fs/bpf/...`) for maps/programs/links

### 🤖 AI Tooling Compatibility

* All tools are **strictly typed** with published schemas
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
* [JSON Schema Spec (2020-12)](https://json-schema.org/)
* [eBPF Security Best Practices](https://ebpf.io/security/)
* [Cilium for Kubernetes Observability](https://cilium.io/)

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
