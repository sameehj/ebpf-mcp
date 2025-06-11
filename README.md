## 🐝 ebpf-mcp: Kernel-Level Observability for AI Agents  
[![License: GPL v2 (eBPF)](https://img.shields.io/badge/License-GPL%20v2-blue.svg)](https://www.gnu.org/licenses/old-licenses/gpl-2.0.html)
[![License: Apache 2.0 (Core)](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://www.apache.org/licenses/LICENSE-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/sameehj/ebpf-mcp)](https://goreportcard.com/report/github.com/sameehj/ebpf-mcp)

<p align="center">
  <img src="assets/logo.png" alt="ebpf-mcp logo" width="360"/>
</p>

---

## 🧠 AI-Ready Observability for Linux

`ebpf-mcp` is a lightweight MCP-compatible server that exposes **structured, AI-safe access to Linux kernel observability tools**, built on top of eBPF.

It wraps powerful tools like `bpftool` and the Cilium eBPF library into JSON-RPC endpoints that can be called by AI agents or CLI clients, enabling:

- ✅ Safe eBPF program deployment (from disk or remote URL)
- ✅ Structured inspection of attached kernel hooks
- ✅ BPF map introspection
- ✅ Error tracing of failing syscalls
- ✅ Kernel capability discovery

---

## ✅ What It Actually Delivers

These features are **implemented, tested, and available today**:

### 🔍 System Introspection

```bash
curl -X POST localhost:8080/mcp -d '{
  "jsonrpc": "2.0", "id": 1,
  "method": "tools/call",
  "params": { "tool": "info", "input": {} }
}'
````

✔ Detects kernel version, BTF support, cgroup v2
✔ Returns structured JSON for AI agents to reason over

### 🧪 Hook Inspection (bpftool wrapped in JSON)

```bash
curl -X POST localhost:8080/mcp -d '{
  "jsonrpc": "2.0", "id": 2,
  "method": "tools/call",
  "params": { "tool": "hooks_inspect", "input": {} }
}'
```

Returns:

```json
{
  "programs": [
    {
      "id": 14,
      "type": "tracepoint",
      "name": "handle_syscall_error",
      "attached_to": "sys_enter",
      "pinned": false
    }
  ]
}
```

### 🚀 eBPF Deployment with Remote Support

```json
{
  "tool": "deploy",
  "args": {
    "program_path": "https://example.com/xdp_prog.o"
  }
}
```

✔ Uses Cilium's Go library
✔ Supports loading from URL or local path
✔ Returns structured success or error output
✔ Prints how many programs/maps were loaded

---

## ⚙️ System Requirements

| Requirement         | Why It Matters                  |
| ------------------- | ------------------------------- |
| **Linux 5.8+**      | For modern eBPF support         |
| **BTF Enabled**     | Required for many bpftool ops   |
| **bpftool in PATH** | Used by inspection tools        |
| **cgroup v2**       | Required for some program types |
| **Clang/LLVM**      | Needed only if compiling `.c`   |

---

## 🔐 Security & Privilege Requirements

`ebpf-mcp` must run with sufficient privileges to interact with the kernel:

* ✅ `CAP_BPF` and `CAP_SYS_ADMIN` usually required
* ✅ XDP and tracepoints need elevated rights
* ⚠️ Always audit `.o` files before loading
* 🧪 `deploy` validates programs via kernel verifier

---

## ❌ Failure Modes to Expect

| Condition               | Behavior                         |
| ----------------------- | -------------------------------- |
| Missing `bpftool`       | `hooks_inspect` fails gracefully |
| Invalid `.o` program    | `deploy` returns error via MCP   |
| Insufficient privileges | Kernel rejects program load      |
| No BTF support          | Some introspection may fail      |

---

## 📡 MCP Protocol Support

| Feature       | Status     |
| ------------- | ---------- |
| `tools/list`  | ✅          |
| `tools/call`  | ✅          |
| `resources/*` | 🚧 Planned |
| Streaming     | 🚧 Planned |

---

## 🔮 Roadmap

> These are **not yet implemented**, but planned:

### 🧠 Claude / MCP Agent Integration

* Claude CLI can call `tools/call`, but doesn’t fully interpret streamed output yet
* Working on improved Claude and Ollama support via `ollama-chat` CLI
* MCP compliance is prioritized for LLM compatibility

### 🧰 Cursor AI (IDE Integration)

* We're exploring ways for Cursor AI to call local MCP endpoints (currently not supported natively)
* Early experiments with `ollama + ebpf-mcp` are promising for kernel debugging inside the dev environment

---

## ⚡ Quick Start

```bash
git clone https://github.com/sameehj/ebpf-mcp.git
cd ebpf-mcp
make build
sudo ./bin/ebpf-mcp-server -t http
```

Then call it using your favorite JSON-RPC client or the included [ollama-chat CLI](./cmd/ollama-chat).

---

## 🔐 Dual Licensing

`ebpf-mcp` uses a dual-license model to balance kernel compatibility with integration flexibility:

* 🧬 **GPL-2.0** for all code under `internal/ebpf/`  
  - Covers eBPF program loading and kernel-level interactions  
  - eBPF programs run in kernel space and may link with GPL-licensed kernel helpers  
  - Ensures compliance and compatibility with the Linux kernel and existing GPL eBPF code

* 🧠 **Apache-2.0** for all other components  
  - Covers the MCP server, protocol layer, tool registry, and client CLI  
  - Allows integration with proprietary or commercial AI agents, dev tools, and infrastructure  
  - Encourages broader adoption and contribution outside the kernel ecosystem

This model keeps kernel code legally compatible while enabling wide, flexible usage in AI-first systems and enterprise automation.

---

## 🧙 Join the eBPF Agent Army

We’re building the first structured agent layer over the Linux kernel — and we need your help:

* ⭐ Star this repo
* 🛠️ Contribute a tool (`internal/tools/`)
* 🧪 File bug reports or integration ideas
* 🤖 Test it with LLMs and share feedback

> Contact: [sameeh.j@gmail.com](mailto:sameeh[dot]j@gmail.com)
> GitHub: [github.com/sameehj/ebpf-mcp](https://github.com/sameehj/ebpf-mcp)
