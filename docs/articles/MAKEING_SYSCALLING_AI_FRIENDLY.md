# Making Syscall Tracing AI-Friendly: A Practical Approach

*How structured interfaces can make kernel debugging accessible to AI agents*

## The Problem: AI Agents Struggle with Traditional eBPF Tools

When you ask Claude or GPT-4 to help debug a performance issue, they can suggest using `strace` or `bpftrace`, but they struggle to actually use these tools effectively. The output is unstructured, the command syntax is complex, and error handling is unpredictable.

Here's what typically happens:

```bash
# AI suggests this command
strace -e trace=file -p 1234 2>&1 | grep -E "(read|write)"

# But the output is inconsistent:
read(3, "data...", 4096) = 42
write(4, "response...", 8) = 8
read(3, "", 0) = 0                      # EOF
read(5, 0x7fff12345678, 1024) = -1 EAGAIN (Resource temporarily unavailable)
```

The AI agent has to parse inconsistent text output, handle edge cases in the format, and guess at error conditions. This leads to brittle automations and frustrated users.

## A Better Approach: Structured Syscall Tracing

What if we could provide the same functionality through a schema-driven interface that AI agents can consume reliably?

### Tool Schema
```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "name": "trace_syscalls",
  "description": "Trace system calls for debugging and performance analysis",
  "parameters": {
    "target": {
      "type": "object",
      "properties": {
        "pid": {"type": "integer", "description": "Process ID to trace"},
        "comm": {"type": "string", "description": "Process name pattern"},
        "container": {"type": "string", "description": "Container ID or name"}
      },
      "oneOf": [{"required": ["pid"]}, {"required": ["comm"]}, {"required": ["container"]}]
    },
    "syscalls": {
      "type": "array",
      "items": {"type": "string"},
      "default": ["openat"],
      "description": "System calls to trace"
    },
    "duration": {
      "type": "integer",
      "minimum": 1,
      "maximum": 300,
      "default": 30,
      "description": "Trace duration in seconds"
    },
    "format": {
      "type": "string",
      "enum": ["json", "csv"],
      "default": "json",
      "description": "Output format"
    }
  },
  "required": ["target"],
  "returns": {
    "type": "array",
    "items": {"$ref": "#/definitions/SyscallEvent"}
  },
  "definitions": {
    "SyscallEvent": {
      "type": "object",
      "properties": {
        "timestamp": {"type": "string", "format": "date-time"},
        "pid": {"type": "integer"},
        "tid": {"type": "integer"},
        "comm": {"type": "string"},
        "syscall": {"type": "string"},
        "args": {"type": "object"},
        "return_value": {"type": "integer"},
        "duration_us": {"type": "integer"},
        "error": {"type": ["object", "null"]}
      }
    }
  }
}
```

### Structured Output
```json
{
  "events": [
    {
      "timestamp": "2024-03-15T10:30:45.123456Z",
      "pid": 1234,
      "tid": 1234,
      "comm": "nginx",
      "syscall": "openat",
      "args": {
        "filename": "/var/log/nginx/access.log",
        "flags": "O_WRONLY|O_APPEND"
      },
      "return_value": 3,
      "duration_us": 45,
      "error": null
    },
    {
      "timestamp": "2024-03-15T10:30:45.125000Z",
      "pid": 1234,
      "tid": 1234,
      "comm": "nginx",
      "syscall": "openat",
      "args": {
        "filename": null,
        "flags": "O_RDONLY",
        "extraction_errors": ["filename_read_failed"]
      },
      "return_value": -1,
      "duration_us": 12,
      "error": {
        "errno": 2,
        "message": "ENOENT: No such file or directory"
      }
    }
  ],
  "summary": {
    "total_events": 2,
    "duration_seconds": 30,
    "syscalls_per_second": 0.067,
    "error_rate": 0.5,
    "tool_stats": {
      "events_captured": 2,
      "events_dropped": 0,
      "ring_buffer_fills": 0,
      "filename_resolution_failures": 1,
      "cpu_overhead_percent": 1.2
    }
  }
}
```

## How It Would Work

### Command Line Interface
```bash
# Start simple: trace nginx for file operations
$ trace-syscalls --target=nginx --syscalls=openat --duration=60

# Container support (when ready)
$ trace-syscalls --container=web-server-1 --format=json > debug.json

# Comparison tool for validation
$ compare-trace --target=nginx --duration=10
```

### Comparison with Existing Tools

| Feature | `strace` | `perf trace` | `trace-syscalls` |
|---------|----------|------------|------------------|
| Structured output | ❌ | Partial | ✅ |
| AI compatibility | ❌ | ❌ | ✅ |
| eBPF-based | ❌ | ✅ | ✅ |
| Container awareness | ❌ | Manual | ✅ |
| Schema validation | ❌ | ❌ | ✅ |
| Error handling | Text parsing | Manual | Structured |
| Performance overhead | Very High | Medium | Target: Low |

## Design Decisions and Open Questions

Rather than trying to solve everything upfront, here are our initial approaches and areas where we need community input:

### Ring Buffer Strategy
**Our hunch**: Start with adaptive ring buffers (1MB-100MB based on event rate). Handle overflow by dropping events but tracking drop counts. Simple backpressure: warn if userspace falls behind and suggest reducing scope.

**Reality check needed**: High-frequency apps (1M syscalls/second) will generate 200MB/second of data. We need streaming or binary formats for busy applications.

### Argument Extraction Safety
**Our approach**: Fail gracefully and document failures in the output. Partial information beats no information.

**Known limitations**: Filename resolution can race with file operations, may be expensive in containers, and impossible for anonymous FDs.

### Container Integration  
**Starting point**: Use `/proc/<pid>/cgroup` to map container names to PIDs. Handle lifecycle by refreshing PID lists every 5 seconds.

**Punt for later**: Kubernetes pod names, sidecar containers, init containers, multiple container runtimes.

### Performance Target
**Our goal**: <2% CPU overhead on traced processes, <1ms latency to first event.

**Critical**: We must beat `perf trace` meaningfully or the value proposition weakens. JSON serialization isn't free - we may need binary format with post-processing.

## Implementation Strategy

### Phase 1: Proof of Concept (Weeks 1-2)
**Scope**: Just `openat` syscalls, just PID targeting, just x86_64
- Basic eBPF program with ring buffer
- JSON schema and output format  
- Performance baseline vs `strace`

**Success criteria**: Works reliably, measurable performance advantage

### Phase 2: Validation Framework (Weeks 3-4)
Build the comparison tool:
```bash
compare-trace --target=nginx --duration=10
# Outputs: accuracy comparison, performance metrics, feature gaps
```

**Success criteria**: Side-by-side validation shows we're not missing events

### Phase 3: Container Integration (Weeks 5-8)
Add container name resolution and namespace handling
- Partner with real Kubernetes users
- Document edge cases and limitations clearly
- Handle failures gracefully

**Success criteria**: Works in containerized environments with known limitations

### Phase 4: AI Agent Testing (Weeks 9-12)
Automated tests with real AI agents doing debugging tasks
- Can Claude parse output and answer questions?
- Does schema prevent hallucinations?
- Are error messages actionable?

**Success criteria**: AI agents successfully complete debugging workflows

## Success Metrics (3 Month Checkpoint)

This project succeeds if:
- [ ] Works reliably for `openat` tracing on Ubuntu 22.04
- [ ] <2% CPU overhead on traced processes (verified by benchmarks)
- [ ] 5+ SRE teams report using it for actual debugging work
- [ ] AI agents can successfully use the tool without text parsing
- [ ] Performance benchmarks published and peer-reviewed

If we hit these targets, we expand scope. If not, we iterate on fundamentals.

## Come Build This With Us

These design decisions are educated guesses. We need:

- **SRE teams** to test real debugging scenarios and break our assumptions
- **Performance engineers** to validate overhead claims with production workloads
- **Container operators** to find edge cases we haven't considered
- **AI researchers** to verify structured output actually helps agents
- **Kernel developers** to review eBPF implementation for safety and efficiency

### Current Status
- ✅ JSON schema and output format designed
- ✅ Basic eBPF program structure planned
- 🚧 Ring buffer management and userspace consumer
- 🚧 Container PID mapping implementation
- ❌ Performance benchmarks (critical path)
- ❌ Side-by-side validation tool

### How to Contribute

**Repository**: [github.com/sameehj/ebpf-mcp](https://github.com/sameehj/ebpf-mcp)

**Immediate needs**:
1. **Test the prototype** on your actual debugging workflows
2. **Benchmark against existing tools** on your production workloads  
3. **Break the container integration** and document failure modes
4. **Validate AI agent integration** with real debugging tasks
5. **Review eBPF implementation** for safety and performance

**Discussion**: GitHub Discussions for design decisions, use cases, and coordination.

## The Real Test

The best argument for this approach isn't theoretical—it's practical. Can we build something that:
- **SREs actually adopt** instead of `strace` for container debugging?
- **AI agents can consume** without brittle text parsing?
- **Performs well enough** for production troubleshooting?
- **Handles edge cases** gracefully with useful error messages?

## Critical Success Factors

### Performance Must Hold Up
Our <2% CPU target is aggressive. `strace` adds 100x overhead, `perf trace` adds 5-15%. We need to beat `perf trace` meaningfully or users won't switch.

### Container Integration Is Make-or-Break
Modern debugging happens in containers. Our `/proc/<pid>/cgroup` approach handles basic cases, but Kubernetes environments have pod names, sidecars, init containers, and restart scenarios we need to handle.

### AI Integration Must Be Validated
We're assuming structured output helps AI agents, but we need to prove it with real debugging tasks and measure improvement over text parsing.

### Graceful Degradation Is Essential
When filename resolution fails, containers restart, or buffers overflow, the tool must provide useful partial information rather than crashing or producing garbage.

---

**Stop debating, start building.** 

The Linux debugging community needs better tools. This is a practical approach to building them. The best validation comes from working code and real users solving real problems.

Join us.