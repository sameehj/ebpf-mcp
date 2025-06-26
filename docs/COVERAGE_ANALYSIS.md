# 📊 Coverage Analysis of eBPF-MCP Operation Model

## Overview

The `ebpf-mcp` project defines a **5-category operational framework** to organize eBPF tooling into structured, composable units. This is **not a formal algebra** like CRUD or relational models, but rather a practical interface pattern for automation and AI integration.

We use this model to:
- Structure MCP-compatible tools into intuitive types
- Cover the majority of kernel observability use cases
- Enable schema-based reasoning for human and AI agents

---

## 🔧 Framework Categories (Recap)

| Category   | Core Purpose                        |
|------------|--------------------------------------|
| `deploy`   | Load programs, maps, and metadata    |
| `control`  | Manage program lifecycle (attach, detach) |
| `reflect`  | Query kernel/eBPF runtime state      |
| `stream`   | Capture real-time kernel data        |
| `map_op`   | Read/write shared kernel data        |

> Think of this as a structured **"eBPF CRUD"** model for automation — but not a complete algebra in the formal sense.

---

## 📐 How This Maps to Real Tools

### 1. [bpftool](https://man7.org/linux/man-pages/man8/bpftool.8.html)

| Feature              | Mapped Category | Supported in `ebpf-mcp`? |
|----------------------|------------------|----------------------------|
| `bpftool prog load`  | `deploy`         | ✅ `ebpf_load`             |
| `bpftool prog attach`| `control`        | ✅ `ebpf_attach`           |
| `bpftool prog show`  | `reflect`        | ✅ `hooks_inspect`         |
| `bpftool map dump`   | `map_op`         | ✅ `map_dump` (MVP)        |
| `bpftool trace`      | `stream`         | 🧪 `trace_errors`          |
| `bpftool cgroup attach` | `control`    | ❌ Not implemented yet     |

---

### 2. [Tracee](https://github.com/aquasecurity/tracee)

| Feature                  | Mapped Category | Supported in `ebpf-mcp`? |
|--------------------------|------------------|----------------------------|
| Event tracing            | `stream`         | 🧪 `trace_errors`          |
| Signature detection      | external logic   | ❌ Requires orchestration  |
| Rules engine             | external logic   | ❌ Out of scope            |
| Program loading          | `deploy`         | ✅ `ebpf_load`             |

---

### 3. [Cilium](https://cilium.io/)

| Feature                          | Mapped Category | Supported in `ebpf-mcp`? |
|----------------------------------|------------------|----------------------------|
| Endpoint BPF management          | `deploy`         | ✅ Basic                   |
| Dynamic attach (TC/XDP)          | `control`        | ✅ Initial                 |
| Policy injection                 | orchestration    | ❌ Not implemented         |
| Metrics/tracing (hubble)         | `stream`         | ❌ Planned                 |
| Map introspection                | `map_op`         | ✅ MVP                     |

---

## 🎯 Interpretation

- **~80% of foundational operations** in common tools are covered
- What’s missing: orchestration, policy engines, high-level coordination
- Goal is not to replace those systems — but to expose their primitives in a structured, AI-compatible format

---

## ⚠️ What This Model Is *Not*

- ❌ Not equivalent to CRUD, REST, or relational algebra
- ❌ Not mathematically complete
- ❌ Not a general-purpose eBPF orchestrator

It is:
- ✅ A practical structure for safe eBPF automation
- ✅ Designed for agents and MCP tooling
- ✅ Extendable through additional tools and orchestration layers

---

## ✅ Conclusion

This coverage model is **implementation-driven, not theory-driven**. It helps us:

- Group tools cleanly
- Expose capability boundaries
- Identify gaps for future roadmap work

> For the most up-to-date implementation status, see the main [README.md](../README.md#📈-current-tool-coverage).