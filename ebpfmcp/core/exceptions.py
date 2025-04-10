class BPFLoadError(Exception):
    """Raised when a BPF program fails to load."""
    pass

class ConfigError(Exception):
    """Raised when configuration is invalid."""
    pass

class MCPError(Exception):
    """Base exception for MCP-related errors."""
    pass
