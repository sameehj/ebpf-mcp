import click
import uvicorn
from .api.server import app

@click.group()
def cli():
    """eBPF-MCP CLI"""
    pass

@cli.command()
@click.option('--host', default='127.0.0.1', help='Host to bind to')
@click.option('--port', default=8000, help='Port to bind to')
def serve(host: str, port: int):
    """Start the eBPF-MCP server."""
    uvicorn.run(app, host=host, port=port)

@cli.command()
@click.argument('tool')
@click.option('--params', '-p', help='Tool parameters in JSON format')
def run(tool: str, params: str):
    """Run a BPF tool directly."""
    click.echo(f"Running tool: {tool}")
    # Implement tool execution logic

def main():
    cli()

if __name__ == '__main__':
    main()
