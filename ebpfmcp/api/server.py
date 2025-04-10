import os
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from typing import Dict, Any

from ..bpf.loader import BPFProgramLoader
from ..core.context import MCPContextGenerator

app = FastAPI(title="eBPF-MCP Server")

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Initialize components
bpf_loader = BPFProgramLoader()
context_gen = MCPContextGenerator()

@app.get("/.well-known/mcp/metadata.json")
async def get_metadata():
    """Get MCP metadata."""
    return {
        "version": "0.1.0",
        "capabilities": ["syscall_trace", "network_monitor"],
        "kernel_version": os.uname().release
    }

@app.post("/api/tools/{tool_name}")
async def execute_tool(tool_name: str, params: Dict[str, Any]):
    """Execute a BPF tool."""
    try:
        if tool_name == "syscall_trace":
            bpf = bpf_loader.load_program("syscall_trace")
            # Tool-specific logic here
            return {"status": "success", "message": "Tool executed successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
