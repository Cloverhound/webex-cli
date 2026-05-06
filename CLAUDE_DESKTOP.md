# Webex MCP — Claude Desktop Setup

This guide covers installing and using the Webex MCP server with **Claude Desktop** specifically. For Claude Code (CLI), see the [main README](README.md).

## Requirements

- **Claude Desktop** — latest version
- **Webex CLI** — installed and authenticated on each user's machine
  ```bash
  # Install
  curl -fsSL https://raw.githubusercontent.com/Cloverhound/webex-cli/main/install.sh | sh

  # Authenticate (opens browser for OAuth)
  webex login
  ```
- **Platform** — macOS, Windows, or Linux

The MCP server runs **locally** using stdio transport — no network port is opened and no shared server is required. Each user's own Webex credentials are used automatically.

---

## Installation

### Option 1 — Admin Portal (Team Distribution)

Admins can distribute the extension to the entire team without each user needing to install it manually.

1. Build the extension package:
   ```bash
   git clone https://github.com/Cloverhound/webex-cli.git
   cd webex-cli
   make extension   # produces webex-mcp.dxt
   ```
2. Upload `webex-mcp.dxt` to the [Claude admin portal](https://console.anthropic.com) under **Extensions**
3. Enable it for your team or org
4. Each team member must still install the Webex CLI and run `webex login` on their own machine

### Option 2 — Individual Install (Double-click)

Download or build `webex-mcp.dxt` and double-click it in Finder (macOS) or Explorer (Windows). Claude Desktop installs the extension automatically.

### Option 3 — HTTP Transport (Manual Config)

If you prefer HTTP transport instead of the DXT extension:

1. Start the MCP server:
   ```bash
   webex mcp serve --http          # binds to 127.0.0.1:47890
   webex mcp serve --http --http-addr 127.0.0.1:9000   # custom port
   ```
2. Add to `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS):
   ```json
   {
     "mcpServers": {
       "webex": { "url": "http://localhost:47890/mcp" }
     }
   }
   ```
3. Restart Claude Desktop
4. Stop the server when done: `webex mcp stop`

The HTTP server only binds to loopback (`127.x.x.x`) and cannot be exposed on public interfaces.

---

## Tools

The extension exposes four tools:

| Tool | Maps to | Auto-approvable | Description |
|---|---|---|---|
| `webex_read` | GET | Yes | List, get, download, export, search — read-only operations |
| `webex_write` | POST / PUT / PATCH / DELETE | No | Create, update, delete — always prompts for confirmation |
| `webex_help` | — | Yes | Get help text and flag listings for any command |
| `webex_usage` | — | Yes | View recent MCP command history with timing and status |

`webex_read`, `webex_help`, and `webex_usage` are safe to auto-approve. `webex_write` is intentionally kept separate so write operations always require explicit user confirmation.

---

## Adding the Skill for Better Results

The DXT extension provides Claude Desktop with the **tools** to execute Webex commands, but not the **knowledge** of how to use them optimally. Without the skill, Claude will use `webex_help` to discover commands on demand — functional, but slower and more prone to mistakes on edge cases like CC `--orgid` handling, pagination, and subcommand naming exceptions.

**The fix:** Use a Claude Desktop **Project** and paste the skill content into the Project instructions. This is the Claude Desktop equivalent of the Claude Code skill system.

### Setup

1. In Claude Desktop, create a new **Project** (e.g., "Webex Admin")
2. Open **Project Instructions**
3. Paste the contents of [`skill/SKILL.md`](skill/SKILL.md) — this covers auth, command structure, global flags, and best practices
4. Optionally add the relevant sub-skill files for the areas your team uses most:

| Area | File | When to include |
|---|---|---|
| Admin | [`skill/admin/SKILL.md`](skill/admin/SKILL.md) | Managing people, licenses, orgs |
| Calling | [`skill/calling/SKILL.md`](skill/calling/SKILL.md) | Locations, queues, devices, recordings |
| Contact Center | [`skill/cc/SKILL.md`](skill/cc/SKILL.md) | Sites, flows, agents, entry points |
| Devices | [`skill/device/SKILL.md`](skill/device/SKILL.md) | Device configs, workspaces, xAPI |
| Meetings | [`skill/meetings/SKILL.md`](skill/meetings/SKILL.md) | Recordings, transcripts, scheduling |
| Messaging | [`skill/messaging/SKILL.md`](skill/messaging/SKILL.md) | Rooms, messages, teams |

With the skill in Project instructions, Claude will know upfront about `--paginate`, RSQL filters, the CC `--orgid` auto-population, download patterns, and other behaviours that would otherwise require multiple `webex_help` round-trips to discover.

---

## Authentication

The extension uses the same credentials as the Webex CLI — run `webex login` once and the MCP server picks them up automatically. Token refresh is handled transparently.

```bash
webex auth status          # confirm you are logged in
webex auth list            # list all stored users
webex auth switch <email>  # change default user
webex auth set-org <orgId> # set a persistent org override (partner admins)
```

Multiple Webex accounts are supported. Switch the default user with `webex auth switch` and restart Claude Desktop (or `webex mcp stop` + re-open for HTTP transport).

---

## Troubleshooting

**"webex: command not found"**
The install script adds `~/.local/bin` to PATH. If Claude Desktop launches before your shell profile is sourced, the binary may not be found. Fix by setting an absolute path in the extension config or ensuring `~/.local/bin` is in your system PATH (not just shell PATH).

For the HTTP transport option, edit `claude_desktop_config.json` to use an absolute path:
```json
{
  "mcpServers": {
    "webex": {
      "command": "/Users/yourname/.local/bin/webex",
      "args": ["mcp", "serve"]
    }
  }
}
```

**"No running MCP server found" (HTTP mode)**
Run `webex mcp serve --http` before opening Claude Desktop, or switch to the DXT extension which starts the server automatically via stdio.

**Not authenticated**
Run `webex login` in a terminal, then retry. Tokens are stored in the OS keyring and shared with the MCP server automatically.

**Wrong Webex org**
Use `webex auth set-org <orgId>` to set a persistent org override, or pass `--organization <orgId>` on individual commands via the `flags` parameter.
