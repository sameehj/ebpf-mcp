from typing import Dict, Optional
import os
from bcc import BPF
from ..core.exceptions import BPFLoadError

class BPFProgramLoader:
    """Load and manage BPF programs."""
    
    def __init__(self, programs_dir: str = "bpf_programs"):
        self.programs_dir = programs_dir
        self.loaded_programs: Dict[str, BPF] = {}
        
    def load_program(self, name: str) -> BPF:
        """Load a BPF program by name."""
        try:
            # Check if program is already loaded
            if name in self.loaded_programs:
                return self.loaded_programs[name]
            
            # Find program file
            program_path = os.path.join(self.programs_dir, f"{name}.c")
            if not os.path.exists(program_path):
                raise FileNotFoundError(f"Program {name} not found at {program_path}")
            
            # Load program
            with open(program_path) as f:
                program_text = f.read()
                
            bpf = BPF(text=program_text)
            self.loaded_programs[name] = bpf
            return bpf
            
        except Exception as e:
            raise BPFLoadError(f"Failed to load {name}: {str(e)}")
    
    def unload_program(self, name: str) -> None:
        """Unload a BPF program."""
        if name in self.loaded_programs:
            # BCC handles cleanup on object destruction
            del self.loaded_programs[name]
