## 🛠️ Tutorial: Add `ebpf-mcp` to Claude CLI for Local Development

### ✅ Prerequisites

* You have `ebpf-mcp` compiled and running on a remote machine (e.g., EC2).
* You have [`claude`](https://docs.anthropic.com/claude/docs/cli-overview) CLI installed locally.
* You want to connect to your remote MCP server via `localhost:8080`.

---

### 🖥️ 1. SSH Port Forwarding

Forward the MCP server’s port `8080` from EC2 to your local machine:

```bash
ssh -L 8080:localhost:8080 aws-t3-small
```

This maps your local `localhost:8080` to the remote `ebpf-mcp` server.

---

### 🧠 2. Start the `ebpf-mcp` server remotely

On the remote EC2 machine, run:

```bash
cd ebpf-mcp
make build
sudo ./bin/ebpf-mcp-server -t http
```

You’ll see output like:

```text
🔐 Auto-generated auth token (no MCP_AUTH_TOKEN was set):
    75158a9db7ce4cfbdda112efdb352708c5b10fe8dd0297b5eeb1b43a93672eb5
💡 Pass this as Authorization: Bearer <token> in Claude or curl headers.
```

📌 **Note the token** above — you’ll need it in the next step.

---

### 🤖 3. Add the MCP Server to Claude CLI

Run this on **your local machine**:

```bash
claude mcp remove ebpf  # Optional: if it already exists

claude mcp add ebpf http://localhost:8080/mcp \
  -t http \
  -H "Authorization: Bearer 75158a9db7ce4cfbdda112efdb352708c5b10fe8dd0297b5eeb1b43a93672eb5"
```

---

### 🧪 4. Test it in Claude REPL

Now launch Claude:

```bash
claude
```

And test a command like:

```text
> run ebpf mcp info
```

You should see a response from the `ebpf-mcp` server 🎉

---

### 🧩 Extra Tips

* If the token changes, re-run `claude mcp add ...` with the new one.
* To make the token static, set it manually via `export MCP_AUTH_TOKEN=...` on the remote server **before** launching `ebpf-mcp`.
