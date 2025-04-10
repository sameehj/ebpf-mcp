import os
import json
import time
from typing import Dict, Any
from pathlib import Path

class MCPContextGenerator:
    """Generate and manage MCP context files."""
    
    def __init__(self, base_dir: str = ".well-known/mcp"):
        self.base_dir = Path(base_dir)
        self.maps_dir = self.base_dir / "maps"
        self.traces_dir = self.base_dir / "traces"
        
        # Create directory structure
        self.maps_dir.mkdir(parents=True, exist_ok=True)
        self.traces_dir.mkdir(parents=True, exist_ok=True)
        
    def update_metadata(self, metadata: Dict[str, Any]) -> None:
        """Update metadata.json file."""
        metadata_path = self.base_dir / "metadata.json"
        metadata.update({
            "last_updated": time.time(),
            "version": "0.1.0"
        })
        
        with open(metadata_path, 'w') as f:
            json.dump(metadata, f, indent=2)
            
    def update_trace(self, name: str, data: str) -> None:
        """Update a trace file."""
        trace_path = self.traces_dir / f"{name}.txt"
        with open(trace_path, 'w') as f:
            f.write(data)
            
    def update_map(self, name: str, data: Dict[str, Any]) -> None:
        """Update a map file."""
        map_path = self.maps_dir / f"{name}.json"
        with open(map_path, 'w') as f:
            json.dump(data, f, indent=2)
