# 🚧 ebpf-mcp Development Roadmap

This document tracks the phased development of the `ebpf-mcp` project, from core operations to AI orchestration and integration with legacy tools.

---

## ✅ Phase 1: Core Tools (MVP)

**Goal:** Validate the 5-category model with practical tools

- ✅ `ebpf_load` for program loading
- ✅ `ebpf_attach` for lifecycle control
- ✅ `info`, `hooks_inspect` for kernel reflection
- ✅ `map_dump` for map state inspection (MVP)
- ✅ `trace_errors` for streaming syscall failures

---

## 🔄 Phase 2: Streaming + Maps

**Goal:** Support richer introspection and live observability

- ⏳ Structured event streaming (`perf`, `ringbuf`)
- ⏳ Multi-key batch operations for maps
- ⏳ User-defined filters and aggregation (MCP compatible)
- ⏳ WebSocket-based streaming endpoint (planned)

---

## 🔐 Phase 3: AI Orchestration + Control

**Goal:** Secure, structured AI control with role enforcement

- ⏳ Role-based access control (RBAC)
- ⏳ LLM safety layers (purpose declaration, token filtering)
- ⏳ Structured logs + audit trails per tool
- ⏳ Claude, Ollama, Cursor AI integration

---

## 🔁 Phase 4: Integration & Ecosystem Expansion

**Goal:** Bridge with existing tools and production systems

- ⏳ Compatibility adapters for bpftool + Tracee
- ⏳ Kubernetes integration (sidecar + agent mode)
- ⏳ Exportable metrics / Prometheus support
- ⏳ CI/CD safe loaders (MCP-only)

---

## 🔬 Experimental Ideas

- Sandbox execution of `user_scripts` (Lua, Python)
- Schema introspection via `reflect`
- Per-tool metadata served via `tools/describe`

---

## 🚧 Known Limitations

- No verifier step-through or debugging tools
- No native tracee signature engine support
- No orchestration language (e.g. for program chaining)
- No map pinning or unpinning (WIP)
- No distributed coordination yet (single-node only)

---

## 📆 Timeline

| Milestone            | Target         |
|----------------------|----------------|
| Core MCP tools       | ✅ Completed   |
| Streaming + map_ops  | Q3 2025        |
| AI RBAC + audit logs | Q3–Q4 2025     |
| Kubernetes adapter   | Q4 2025        |

---

## 🧠 Prioritization Criteria

We prioritize features that:

- Increase AI/agent compatibility
- Improve runtime safety and auditability
- Reduce operational overhead for devs
- Enable meaningful structured introspection

---

> Want to help? Open an issue or PR against [`internal/tools/`](../internal/tools) — or reach out on [GitHub](https://github.com/sameehj/ebpf-mcp).

