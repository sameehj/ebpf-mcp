# ebpf-mcp Design Document

## 🧠 Executive Summary

`ebpf-mcp` is a lightweight, AI-compatible, JSON-RPC 2.0 server designed to provide structured access to Linux eBPF observability and control tools. Inspired by the Model Context Protocol (MCP), it allows AI agents to request operations via well-defined `tools`, `resources`, and `prompts` abstractions. This enables programmatic kernel insight, automation, and experimentation in a secure, extensible fashion.

## 🏗️ Architecture Design

```
+--------------------+       +-----------------------+       +----------------------+
|   MCP Client (AI)  | <---> |   ebpf-mcp JSON-RPC   | <---> |   Tool Manager       |
+--------------------+       +-----------------------+       +----------+-----------+
                                                      |
                                                      v
                                            +----------------------+
                                            |   eBPF Manager       |
                                            | - Program Loader     |
                                            | - Map Inspector      |
                                            | - Event Monitor      |
                                            +----------------------+
                                                      |
                                                      v
                                            +----------------------+
                                            |   Security Manager   |
                                            | - Policy Enforcer    |
                                            | - Bytecode Validator |
                                            | - Sandbox Isolation  |
                                            +----------------------+
                                                      |
                                                      v
                                            +----------------------+
                                            |   Resource Manager   |
                                            | - Limit Enforcer     |
                                            | - Cleanup Daemon     |
                                            | - Usage Tracker      |
                                            +----------------------+
```

### Component Responsibilities

* **MCP Handler**: Parses, validates, and dispatches JSON-RPC requests with comprehensive error handling.
* **Tool Manager**: Manages built-in and user-defined tools, handles parameters, validation, and execution.
* **eBPF Manager**: Loads/verifies programs, inspects maps, hooks to tracepoints with BTF/CO-RE support.
* **Security Manager**: Validates bytecode, applies policies, enforces isolation via namespaces and seccomp.
* **Resource Manager**: Tracks resource usage, enforces limits, and ensures proper cleanup of eBPF objects.

## 📁 Directory Structure

```
ebpf-mcp/
├── cmd/                # CLI entry point and server main
├── internal/
│   ├── core/           # MCP protocol logic (JSON-RPC 2.0)
│   ├── tools/          # Tool definitions and execution logic
│   ├── ebpf/           # Program loading, map handling, events
│   ├── security/       # Bytecode validation, policies, sandboxing
│   ├── resources/      # Resource tracking, limits, cleanup
│   ├── config/         # Configuration management and validation
│   ├── server/         # Main loop, API handlers, middleware
│   └── monitoring/     # Metrics, health checks, observability
├── pkg/
│   ├── types/          # Shared type definitions
│   └── utils/          # Common utilities and helpers
├── examples/           # Sample tools and configurations
├── scripts/            # Install / build / test scripts
├── configs/            # Default and reference YAML configs
├── tests/              # Integration and end-to-end tests
└── docs/               # Documentation and usage examples
```

## ⚙️ Core Components

### Server

Main runtime entry point responsible for configuration loading, starting the MCP protocol handler, and managing lifecycle, metrics, and graceful shutdown.

### MCPHandler

Handles incoming JSON-RPC 2.0 requests, routes them to appropriate methods (`tools/call`, `tools/list`, etc.), validates schemas, and formats responses with proper error handling.

### ToolManager

Manages the registry of built-in tools. Responsible for exposing metadata, parameter validation, security checks, sandboxing, and execution coordination with other managers.

### EBPFManager

Wraps eBPF-related operations such as program loading, map inspection, and event monitoring. Abstracts away low-level kernel interactions while providing BTF/CO-RE compatibility.

### SecurityManager

Enforces per-tool policies, bytecode validation, user-based restrictions, and runtime containment via seccomp/cgroups/namespaces. Supports audit logging and integration with AppArmor/SELinux.

### ResourceManager

Tracks active eBPF programs and maps, enforces resource limits (memory, CPU, program count), implements cleanup mechanisms, and provides usage metrics for monitoring.

## 🤩 MCP Protocol Implementation

The server implements a subset of the MCP (Model Context Protocol) with support for:

* `tools/list`: Return all supported tool IDs and metadata
* `tools/call`: Execute a tool with parameters and resource validation
* `resources/list` *(future)*: List active maps/programs with usage statistics
* `resources/read` *(future)*: Read map contents or program info with access control
* `prompts/ask` *(future)*: Interactive AI requests (not implemented in v1)

Follows strict JSON-RPC 2.0 spec with comprehensive schema validation and structured error responses.

### Error Handling

```go
type EBPFError struct {
    Category    ErrorCategory `json:"category"`
    Code        string       `json:"code"`
    Message     string       `json:"message"`
    Recoverable bool         `json:"recoverable"`
    Context     interface{}  `json:"context"`
    Timestamp   time.Time    `json:"timestamp"`
}

type ErrorCategory string

const (
    ErrorValidation  ErrorCategory = "validation"
    ErrorPermission  ErrorCategory = "permission"
    ErrorResource    ErrorCategory = "resource"
    ErrorKernel      ErrorCategory = "kernel"
    ErrorInternal    ErrorCategory = "internal"
)
```

## 🔬 eBPF Integration Layer

Built using the [cilium/ebpf](https://github.com/cilium/ebpf) Go library with comprehensive error handling and resource management.

* Program types: XDP, tracepoints, kprobes, uprobes (expanded support planned)
* Map types: Array, Hash, Per-CPU, RingBuffer, LRU Hash
* Event streaming via polling, perf buffers, and ring buffers

### Compatibility

* **CO-RE** (Compile Once – Run Everywhere): supports portable programs across kernel versions
* **BTF** (BPF Type Format): enables introspection, verifier support, and automatic structure decoding

### Resource Management

```go
type ResourceLimits struct {
    MaxPrograms      int           `yaml:"max_programs"`
    MaxMaps          int           `yaml:"max_maps"`
    MaxInstructions  int           `yaml:"max_instructions"`
    MaxMemoryPages   int           `yaml:"max_memory_pages"`
    MaxCPUTime       time.Duration `yaml:"max_cpu_time"`
    CleanupInterval  time.Duration `yaml:"cleanup_interval"`
    MaxEventRate     int           `yaml:"max_events_per_second"`
}

type ResourceTracker struct {
    ActivePrograms   map[string]*ProgramHandle
    ActiveMaps       map[string]*MapHandle
    UsageMetrics     *ResourceMetrics
    CleanupHandlers  []CleanupFunc
}
```

## 🛠 Tool Definitions

Each `tool` is a declarative Go struct that defines execution parameters, security requirements, and resource needs:

```go
type Tool struct {
    ID           string              `json:"id"`
    Title        string              `json:"title"`
    Description  string              `json:"description"`
    Parameters   []Param             `json:"parameters"`
    Security     SecurityRequirements `json:"security"`
    Limits       ResourceLimits      `json:"limits"`
    Run          ToolExecutor        `json:"-"`
}

type SecurityRequirements struct {
    RequiredCapabilities []string      `json:"required_capabilities"`
    AllowedNamespaces   []string      `json:"allowed_namespaces"`
    MaxExecutionTime    time.Duration `json:"max_execution_time"`
    AllowedAttachPoints []string      `json:"allowed_attach_points"`
}
```

### Example: `map_dump`

```json
{
  "id": "map_dump",
  "title": "Dump Map Contents",
  "description": "Reads and returns all key-value pairs from a given eBPF map.",
  "parameters": [
    {"name": "map_name", "type": "string", "required": true},
    {"name": "max_entries", "type": "integer", "default": 1000}
  ],
  "security": {
    "required_capabilities": ["CAP_BPF"],
    "max_execution_time": "30s"
  }
}
```

### Registration and Execution

All tools are registered at init time and listed via `tools/list`. Each tool undergoes security validation, resource checking, and sandboxed execution.

### Usage Example (JSON-RPC)

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "id": "map_dump",
    "args": {
      "map_name": "xdp_stats",
      "max_entries": 500
    }
  },
  "id": 1
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "result": {
    "entries": [
      {"key": "eth0", "value": {"packets": 1000, "bytes": 64000}},
      {"key": "eth1", "value": {"packets": 500, "bytes": 32000}}
    ],
    "total_entries": 2,
    "execution_time_ms": 15
  },
  "id": 1
}
```

## 🛡️ Security Model

**Multi-layered enforcement with defense in depth:**

### 1. Request Validation
- JSON schema validation for all parameters
- User identity verification and authorization
- Tool availability and permission checks
- Input sanitization and bounds checking

### 2. Access Control
- UID/GID-based permissions with group membership
- cgroup and namespace-based restrictions
- MAC label enforcement (AppArmor/SELinux integration)
- Per-tool capability requirements

### 3. Bytecode Validation
```go
type BytecodeValidator struct {
    MaxInstructions  int         `yaml:"max_instructions"`
    AllowedOpcodes   []BPFOpcode `yaml:"allowed_opcodes"`
    BlockedSymbols   []string    `yaml:"blocked_symbols"`
    BlockedHelpers   []string    `yaml:"blocked_helpers"`
}
```

### 4. Runtime Isolation
- Mount, user, and network namespace containment
- seccomp filters for syscall restriction
- cgroup limits for CPU, memory, and I/O
- Process isolation with restricted capabilities

### 5. Audit and Monitoring
```go
type AuditLogger struct {
    LogFile     string
    Structured  bool
    Fields      []string
}

type AuditEvent struct {
    Timestamp   time.Time   `json:"timestamp"`
    UserID      int         `json:"user_id"`
    ToolID      string      `json:"tool_id"`
    Parameters  interface{} `json:"parameters"`
    Result      string      `json:"result"`
    Error       string      `json:"error,omitempty"`
    Duration    int64       `json:"duration_ms"`
}
```

### 6. Resource Protection
- Per-session program and map limits
- Memory and CPU time enforcement
- Automatic cleanup of orphaned resources
- Rate limiting for tool execution

### 7. Future Enhancements
- Cryptographic signature verification for eBPF programs
- Hardware security module (HSM) integration
- Remote attestation for trusted execution
- Integration with cloud security services

## 🚀 Implementation Plan

### Phase 1: Foundation & Core Security (Weeks 1–4)

**Week 1-2: Core Infrastructure**
* Implement MCP protocol handler with JSON-RPC 2.0
* Build configuration management with YAML validation
* Set up structured logging and error categorization
* Create basic tool registry framework

**Week 3-4: Security Foundation**
* Implement bytecode validator with instruction filtering
* Add basic sandboxing with namespaces and seccomp
* Create resource tracking and limit enforcement
* Add audit logging infrastructure

### Phase 2: eBPF Integration & Tool Development (Weeks 5–8)

**Week 5-6: eBPF Core**
* Create eBPF program loader with BTF/CO-RE support
* Implement map manager and event monitoring
* Add kernel compatibility detection
* Build resource cleanup mechanisms

**Week 7-8: Tool Implementation**
* Implement core tools: `deploy`, `map_dump`, `trace`, `info`
* Add tool parameter validation and execution sandboxing
* Create tool-specific resource limits
* Implement comprehensive error handling

### Phase 3: Production Hardening & Observability (Weeks 9–12)

**Week 9-10: Production Features**
* Integrate Prometheus metrics and health endpoints
* Implement advanced security features (AppArmor/SELinux)
* Add configuration hot-reload capability
* Create comprehensive monitoring dashboards

**Week 11-12: Deployment & Testing**
* Package as binary, Docker image, and system packages
* Implement chaos and penetration testing harness
* Create production deployment guides
* Finalize documentation and examples

## 🧪 Testing Strategy

### Comprehensive Test Coverage

**Unit Tests (>90% coverage)**
* All core modules with mocked dependencies
* Security validation logic with edge cases
* Error handling and recovery mechanisms
* Configuration parsing and validation

**Integration Tests**
* End-to-end eBPF program loading and execution
* Multi-tool orchestration and resource sharing
* Security boundary validation and bypass attempts
* Performance testing under load

**Security Testing**
* Malicious bytecode injection attempts
* Privilege escalation via crafted parameters
* Resource exhaustion and DoS scenarios
* Audit log integrity and completeness

**Chaos Engineering**
* Process crashes during eBPF program execution
* Network interface failures during tracing
* Kernel module loading/unloading during operation
* Memory pressure and OOM scenarios

**Regression Testing**
* Tool behavior consistency across kernel versions
* Performance benchmarks and latency tracking
* API contract validation and schema compliance

**Production Simulation**
* Load testing with concurrent AI agents
* Long-running stability tests (24+ hours)
* Resource leak detection and cleanup verification
* Graceful shutdown and recovery testing

## 🚚 Deployment Guide

### Prerequisites

* Linux system with modern kernel (>= 5.8 recommended, 4.18+ minimum)
* `clang`, `llvm`, `bpftool`, `iproute2`, and kernel headers
* Go >= 1.20 for source builds
* Docker (optional for containerized deployment)
* Appropriate user permissions or sudo access

### Deployment Options

#### 1. Binary Installation

```bash
curl -sSL https://raw.githubusercontent.com/your-org/ebpf-mcp/main/scripts/install.sh | bash
```

Installs to `/usr/local/bin/ebpf-mcp`, creates systemd service, and sets up default configuration.

#### 2. Docker Deployment

```bash
docker run -d --name ebpf-mcp \
  --privileged --network host \
  -v /sys:/sys:ro \
  -v /lib/modules:/lib/modules:ro \
  -v /etc/ebpf-mcp:/etc/ebpf-mcp:ro \
  -p 8080:8080 \
  your-org/ebpf-mcp:latest
```

#### 3. Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: ebpf-mcp
spec:
  selector:
    matchLabels:
      app: ebpf-mcp
  template:
    spec:
      hostPID: true
      hostNetwork: true
      containers:
      - name: ebpf-mcp
        image: your-org/ebpf-mcp:latest
        securityContext:
          privileged: true
        volumeMounts:
        - name: sys
          mountPath: /sys
          readOnly: true
        - name: modules
          mountPath: /lib/modules
          readOnly: true
      volumes:
      - name: sys
        hostPath:
          path: /sys
      - name: modules
        hostPath:
          path: /lib/modules
```

#### 4. Build from Source

```bash
git clone https://github.com/your-org/ebpf-mcp.git
cd ebpf-mcp
make build
sudo ./bin/ebpf-mcp --config ./configs/production.yaml
```

### Production Configuration

```yaml
server:
  bind_address: "127.0.0.1:8080"
  max_concurrent_requests: 100
  request_timeout: "30s"
  tls:
    enabled: true
    cert_file: "/etc/ebpf-mcp/server.crt"
    key_file: "/etc/ebpf-mcp/server.key"

security:
  enable_sandboxing: true
  enable_audit: true
  audit_log: "/var/log/ebpf-mcp/audit.log"
  
  bytecode_validation:
    max_instructions: 4096
    allowed_opcodes:
      - "BPF_ALU"
      - "BPF_LD"
      - "BPF_ST"
      - "BPF_JMP"
    blocked_symbols:
      - "kallsyms_lookup_name"
      - "kernel_text_address"
    blocked_helpers:
      - "bpf_probe_write_user"
      - "bpf_override_return"

  access_control:
    permitted_users: ["ebpf-user", "monitoring"]
    permitted_groups: ["ebpf", "admin"]
    required_capabilities: ["CAP_BPF", "CAP_PERFMON"]

resources:
  max_programs: 10
  max_maps: 50
  max_memory_mb: 64
  max_cpu_time_ms: 1000
  cleanup_interval: "5m"
  max_events_per_second: 10000

ebpf:
  enable_btf: true
  enable_core: true
  program_timeout: "30s"
  map_cleanup_interval: "10m"

observability:
  metrics:
    enabled: true
    prometheus_endpoint: ":9090/metrics"
    collection_interval: "15s"
  
  logging:
    level: "info"
    format: "json"
    file: "/var/log/ebpf-mcp/server.log"
  
  health:
    endpoint: "/health"
    check_interval: "30s"

tools:
  map_dump:
    enabled: true
    max_entries: 1000
    timeout: "10s"
  
  trace:
    enabled: true
    max_duration: "60s"
    max_events: 10000
  
  deploy:
    enabled: false  # Disabled by default in production
    allowed_types: ["xdp", "tc"]
```

## ✅ Next Steps

### Immediate Implementation (Weeks 1-4)
1. **Set up development environment** with eBPF toolchain and Go 1.20+
2. **Implement JSON-RPC 2.0 server** with MCP protocol compliance
3. **Build configuration management** with YAML parsing and validation
4. **Create structured logging** with categorized error handling
5. **Implement basic security framework** with bytecode validation
6. **Add resource tracking** with limits and cleanup mechanisms

### Core Development (Weeks 5-8)
7. **Develop eBPF integration layer** with cilium/ebpf library
8. **Implement core tools** (map_dump, trace, info, deploy)
9. **Add comprehensive testing** including unit, integration, and security tests
10. **Create sandboxing infrastructure** with namespaces and seccomp

### Production Readiness (Weeks 9-12)
11. **Integrate monitoring and metrics** with Prometheus endpoints
12. **Add advanced security features** (AppArmor/SELinux integration)
13. **Create deployment packages** (binary, Docker, Kubernetes manifests)
14. **Develop comprehensive documentation** including API reference and operational guides
15. **Implement chaos and performance testing** for production validation