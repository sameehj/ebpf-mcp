from setuptools import setup, find_packages

setup(
    name="ebpf-mcp",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "fastapi>=0.68.0",
        "uvicorn>=0.15.0",
        "bcc>=0.19.0",
        "pyyaml>=5.4.1",
        "python-jose>=3.3.0",
        "pydantic>=1.8.2",
    ],
    entry_points={
        'console_scripts': [
            'ebpf-mcp=ebpfmcp.cli:main',
        ],
    },
)
