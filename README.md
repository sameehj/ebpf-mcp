## 🐝 ebpf-mcp: MCP-Compatible AI Server for Linux eBPF Control
[![License: GPL v2](https://img.shields.io/badge/License-GPL%20v2-blue.svg)](https://www.gnu.org/licenses/old-licenses/gpl-2.0.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/sameehj/ebpf-mcp)](https://goreportcard.com/report/github.com/sameehj/ebpf-mcp)

> 🧠 Turn AI agents into kernel-native observability tools.  
> 🛡️ Structured. Secure. AI-Ready.  
> 🔬 Build the eBPF Agent Army.

**License:** GPL-2.0 — all code is currently licensed under the GNU General Public License v2.0
**Status:** Pre-release  **AI-Ready:** Yes

`ebpf-mcp` is a local **MCP server** that exposes Linux eBPF observability, instrumentation, and program control through a structured, AI-compatible interface. It follows the official [Model Context Protocol (MCP)](https://github.com/modelcontextprotocol/spec), enabling AI assistants (like Claude, LLaMA, GPT) to safely and intelligently invoke kernel-level tools.

---

## ⚔️ The Vision: eBPF Agent Army

We're building the first AI-compatible **Agent Layer for the Linux Kernel**.

Imagine this:
- You chat with your server.
- It understands and invokes kernel-level eBPF tools.
- It traces, debugs, monitors, and adapts — in real-time.

These agents aren’t just observers.  
They’re **doers** — inside your kernel.

We call this the **eBPF Agent Army**:  
A growing ecosystem of AI-guided agents capable of low-level, high-trust observability and control.

> "Into the soul of the kernel." 🧬

---

## 📚 Table of Contents

* [Current Status](#current-status)
* [The Problem](#the-problem)
* [The Solution](#the-solution)
* [What Is MCP?](#what-is-mcp)
* [Why eBPF?](#why-ebpf)
* [Who Should Use This](#who-should-use-this)
* [Project Purpose](#project-purpose)
* [Real-World Scenarios](#real-world-scenarios)
* [What Makes This Project Special](#what-makes-this-project-special)
* [MCP Protocol Compliance](#mcp-protocol-compliance)
* [Architecture](#architecture)
* [Quick Start](#quick-start)
* [Available Tools](#available-tools)
* [Example JSON-RPC Calls](#example-json-rpc-calls)
* [Roadmap](#roadmap)
* [License](#license)
* [Contributing](#contributing)

---

## 📌 Current Status

✅ MVP under development as an **MCP-compatible Go server**
✅ `tools/list` and `tools/call` implemented
🚧 More tools being ported from prototype CLI

---

## ❓ The Problem

AI assistants today can't reason about low-level Linux kernel behavior — there's no structured way for them to:

* Monitor system-level traffic or syscall activity
* Load or control eBPF programs
* Interpret BPF map data
* Use observability tools in a secure, machine-readable way

Existing tooling (e.g. `bpftool`, `bpftrace`) isn't designed for automated or AI-driven use.

---

## ✅ The Solution

`ebpf-mcp` provides an **MCP-compatible server** that:

* Exposes safe eBPF tools as **MCP `tools`**
* Uses **JSON-RPC 2.0**, schemas, and structured responses
* Enables AI agents to deploy, observe, trace, and reason about kernel behavior
* Bridges DevOps, security, and AI observability

---

## 🧠 What Is MCP?

The **Model Context Protocol (MCP)** is a standardized way for AI assistants to interact with tools and data services.

* JSON-RPC 2.0-based
* Defines `tools/list`, `tools/call`, `resources/*` methods
* Enables AI agents to discover, invoke, and reason with tools securely and predictably

For full spec, see: [modelcontextprotocol/spec](https://github.com/modelcontextprotocol/spec)

---

## 🧪 Why eBPF?

[eBPF](https://ebpf.io/) enables safe, efficient, programmable observability inside the Linux kernel.
With `ebpf-mcp`, AI assistants gain:

* Live monitoring of network traffic, syscalls, errors
* Control over program load/attach/unload lifecycle
* Access to structured BPF map data
* Compatibility with XDP, kprobes, tracepoints, and more

---

## 👥 Who Should Use This

* 🤖 **AI/LLM developers** building intelligent infrastructure tools
* 🛡️ **Security engineers** needing automated threat detection
* ⚡ **SREs/DevOps** wanting AI-assisted performance debugging
* 🔬 **System developers** debugging kernel-level issues
* 🏢 **Platform teams** building observability-as-a-service

---

## 🎯 Project Purpose

To bridge advanced Linux kernel observability with LLMs and agents by exposing eBPF control via a **structured, discoverable, AI-native protocol** (MCP).

Use `ebpf-mcp` to:

* Deploy & remove eBPF programs
* Query live map data
* Trace syscalls
* Monitor traffic per interface or container
* Let agents reason about low-level system behavior

---

## 🌟 Real-World Scenarios

### 🤖 AI-Driven Incident Response

Ask: *"Why is CPU spiking on production servers?"*
→ AI deploys CPU profilers, traces network + system usage, reports Redis overload + suggests tuning

### 🎮 Interactive Kernel Debugging

Ask: *"Why is my kernel module crashing?"*
→ AI deploys kprobes, catches crash location, analyzes cause, and suggests fix

### ⚡ Zero-Downtime Performance Optimization

Ask: *"Why is the DB 50% slower today?"*
→ AI traces syscalls + I/O, detects cache thrashing, recommends sysctl tweaks

### 🛡️ Real-time Threat Hunting

Ask: *"Scan for privilege escalation attempts"*
→ AI monitors setuid/setgid, traces ancestry, flags abuse patterns

### 🔍 Security Analysis

Ask: *"Is there any suspicious network activity on this server?"*
→ AI deploys eBPF network probes, analyzes patterns, identifies anomalies

### 🚨 Performance Debugging

Ask: *"Why is my application making so many syscalls?"*
→ AI traces your app, correlates syscall patterns, suggests optimizations

### 🧰 Container Monitoring

Ask: *"Which containers are using the most network bandwidth?"*
→ AI monitors traffic per namespace, provides ranked analysis

---

## 🚀 What Makes This Project Special

This project sits at the intersection of three trends:

1. **AI automation** — LLMs want to control infrastructure
2. **Observability revolution** — eBPF is becoming the standard
3. **Structured protocols** — MCP enables safe AI tool usage

---

## 📦 MCP Protocol Compliance

This project fully adheres to the [Model Context Protocol](https://github.com/modelcontextprotocol/spec):

* ✅ Supports `tools/list`, `tools/call`
* ✅ Uses standard JSON-RPC 2.0 message format
* ✅ Clearly defined inputs/outputs for each tool
* ✅ No custom or invalid fields

---

## 🧱 Architecture

`ebpf-mcp` sits between AI agents and the Linux kernel, exposing a structured interface to low-level observability tools.

```text
        [ User / AI Assistant / LLM (Claude, LLaMA, GPT) ]
                              ↓
                [ MCP JSON-RPC Client (e.g. ollama-chat) ]
                              ↓
                     ┌────────────────────┐
                     │     ebpf-mcp       │
                     └────────────────────┘
                    ↙          ↓           ↘
           trace_errors   map_dump   hooks_inspect
                ↓            ↓           ↓
      Linux Kernel / eBPF Subsystem (XDP, kprobes, maps)


---

## ⚡ Quick Start

```bash
git clone https://github.com/sameehj/ebpf-mcp.git
cd ebpf-mcp
go build -o ebpf-mcp-server .
./ebpf-mcp-server
```

Then POST valid JSON-RPC 2.0 requests to `localhost:8080/rpc`

---

## 🧰 Available Tools (Sample)

| Tool Name              | Description                                      |
| ---------------------- | ------------------------------------------------ |
| `ebpf.deploy`          | Load a compiled BPF program to interface or hook |
| `ebpf.map_dump`        | Dump contents of a named BPF map                 |
| `ebpf.info`            | Return kernel, distro, and BPF support status    |
| `ebpf.trace_errors`    | Trace failing syscalls (e.g., EPERM)             |
| `ebpf.monitor.traffic` | Count packets per interface/port via XDP         |

---

## 📡 Example JSON-RPC Calls

### 🧠 List Tools

```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "id": 1
}
```

### 🧠 Response to `tools/list`

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "tools": [
      {
        "name": "ebpf.deploy",
        "description": "Load a compiled BPF program to interface or hook",
        "inputSchema": {
          "type": "object",
          "properties": {
            "program": { "type": "string" },
            "interface": { "type": "string" }
          },
          "required": ["program", "interface"]
        }
      }
    ]
  }
}
```

### 🚀 Call a Tool

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "tool": "ebpf.deploy",
    "input": {
      "program": "xdp_pass",
      "interface": "eth0"
    }
  }
}
```

### ✅ Response to `tools/call`

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "XDP program 'xdp_pass' successfully attached to eth0"
      }
    ]
  }
}
```

---

## 🛣️ Roadmap

Here’s where the eBPF Agent Army is headed.

> Want to contribute? Open an issue or PR — we welcome collaborators! 🤝

### ✅ Phase 1: Minimal Viable Agent Layer

* [x] MCP-compliant JSON-RPC server
* [x] Tool system: `tools/list`, `tools/call`
* [x] LLM integration via `ollama-chat` CLI
* [x] Core tools: `info`, `hooks_inspect`, `map_dump`, `trace_errors`

### 🚧 Phase 2: Observability + Interactivity

* [ ] Structured map schema support
* [ ] Event streaming & log-follow (`watch` tool output over time)
* [ ] Tool plugin interface (`/plugins/*.so` or Go modules)
* [ ] Agent-authenticated tool execution

### 🔜 Phase 3: Full MCP Agent Runtime

* [ ] Support for `resources/list`, `resources/read`
* [ ] Prompt memory / stateful sessions
* [ ] AI agent scaffolding (for auto-responders & watchdogs)
* [ ] Secure CLI/API for remote tool invocation

### 🧪 Experimental / Stretch Goals

* [ ] Live chat UI with embedded tool visualizations
* [ ] Cross-node tool coordination (multi-host eBPF agents)
* [ ] Automatic map discovery + introspection
* [ ] Integration with Kubernetes operators for agent injection

---

## 🎯 Bonus: Optional Labels

If you track issues in GitHub, consider labeling them:

* `type:tool`
* `type:agent-feature`
* `status:help-wanted`
* `status:experimental`

---

## 🪧 License

GPL-2.0 — see [LICENSE](./LICENSE)

---

## 🤝 Contributing

* 📥 Fork & submit PRs
* 💡 Suggest new tools or use cases
* 🧪 Share testing feedback
* ✨ Help extend MCP support for resource discovery and streaming

> `ebpf-mcp` is the AI-ready interface to Linux kernel observability. Let's build it together.

---

## 🧙 Join the Kernel-Aware AI Movement

This is the future of observability.  
It’s open. It’s structured. It’s agent-ready.

We’re assembling an army of open-source hackers, SREs, and kernel fans to build the next layer of AI-native infrastructure.

🛠️ Star the repo  
🧠 Join the discussion  
💥 Contribute a tool or LLM integration

> GitHub: https://github.com/sameehj/ebpf-mcp  
> Demos: coming soon  
> Let's awaken the agents. 🐝
