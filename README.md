# Webex CLI

A command-line tool for Webex APIs — Calling, Contact Center, Admin, Devices, Meetings, and Messaging.

See the [docs site](https://cloverhound.github.io/webex-cli/) for full API reference and guides.

## Install

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/Cloverhound/webex-cli/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/Cloverhound/webex-cli/main/install.ps1 | iex
```

Or download from [Releases](https://github.com/Cloverhound/webex-cli/releases).

The installers verify the download against the release's `checksums.txt` and find the latest version without the GitHub API. For CI or sandboxed agents:

- Pin a version with `WEBEX_CLI_VERSION=0.16.0` (or `sh -s -- -v 0.16.0`); `webex update` honors it too.
- Set `INSTALL_DIR` to install somewhere other than `~/.local/bin` (Windows: `%LOCALAPPDATA%\webex-cli`).
- Without a terminal, the installer does not prompt: it prints the PATH line to add and installs the agent skill for detected agents. Run `webex post-install --yes` to accept the defaults (including editing your shell profile), `--no-skills` to skip skills, or set `WEBEX_CLI_NONINTERACTIVE=1` to turn prompts off.

## Quick Start

```bash
# Login (opens browser for OAuth)
webex login

# Webex Calling
webex calling people list --max 10
webex calling locations list
webex calling call-queue list-cxe

# Contact Center
webex cc site list
webex cc team list
webex cc contact-service-queue list
webex cc entry-point list

# Admin
webex admin people list
webex admin licenses list
webex admin organizations get --org-id <id>

# Devices
webex device devices list
webex device xapi execute-command --device-id <id> --command-name <name>

# Meetings
webex meetings meetings list
webex meetings recordings list

# Messaging
webex messaging rooms list
webex messaging messages list --room-id <id>
```

## Audio File & Recording Downloads

Several commands provide streamlined binary download and multipart upload for audio files and recordings:

```bash
# Download a meeting recording (audio, video, or transcript)
webex meetings recordings download --recording-id <id> --output meeting.mp3
webex meetings recordings download --recording-id <id> --output meeting.mp4 --type recording
webex meetings recordings download --recording-id <id> --output meeting.vtt --type transcript

# Download a converged recording (Webex Calling call recording)
webex calling converged-recordings download --recording-id <id> --output call.mp3

# Upload/download Contact Center audio files
webex cc audio-files download --id <id> --output prompt.wav
webex cc audio-files upload --file prompt.wav --name "Main Greeting"

# Upload/download CC agent personal greetings
webex cc agent-personal-greeting-files download --id <id> --output greeting.wav
webex cc agent-personal-greeting-files upload --agent-id <id> --file greeting.wav

# Upload announcement greetings (org or location level)
webex calling announcement-repository upload-binary-greeting --file greeting.wav --name "Greeting"
webex calling announcement-repository upload-binary-greeting-2 --file greeting.wav --name "Greeting" --location-id <id>

# Upload voicemail/intercept greetings (person, virtual line, workspace, or self)
webex calling user-call update-busy-voicemail-greeting-person --person-id <id> --file greeting.wav
webex calling call-settings-for-me upload-voicemail-busy-greeting --file greeting.wav
webex calling virtual-line-call update-busy-voicemail-greeting --virtual-line-id <id> --file greeting.wav
webex calling workspace-call update-busy-voicemail-greeting-place --workspace-id <id> --file greeting.wav
```

All upload commands support `--dry-run` to preview the request without sending it.

## API Coverage

### Calling (`webex calling`)

49 resource groups including auto-attendants, call queues (CxE), hunt groups, call controls, call routing (dial plans, route groups, trunks), DECT devices, emergency services, locations, numbers, paging groups, people, workspaces, voicemail, converged recordings, and more.

### Contact Center (`webex cc`)

64 resource groups including sites, queues, entry points, teams, flows, skills, desktop layouts, global variables, business hours, auxiliary codes, campaigns, callbacks, realtime stats, AI assistant, journey analytics, subscriptions, and more.

### Admin (`webex admin`)

39 resource groups including people, licenses, organizations, roles, groups, events, reports, recordings, SCIM 2.0 (users/groups/schemas), hybrid clusters/connectors, security audit, service apps, and more.

### Devices (`webex device`)

10 resource groups including devices, device configurations, workspaces, workspace locations/metrics/personalization, hot-desking, and xAPI (execute commands, query status).

### Meetings (`webex meetings`)

23 resource groups including meetings, participants, recordings, transcripts, summaries, polls, Q&A, chats, invitees, preferences, session types, tracking codes, video mesh, and more.

### Messaging (`webex messaging`)

12 resource groups including rooms, messages, memberships, teams, team memberships, webhooks, events, attachment actions, room tabs, and more.

## Postman refresh: command changes

Resource operations use `get`, `delete`, `patch`, and `update` with `--id`; previous `get-id`, `delete-id`, `patch-id`, and `update-id` spellings remain aliases.

- Calling adds `call-controls-for-me` and `metrics get-call-quality-stats`.
- Contact Center adds flow activity/event/template discovery, custom functions, assets/channels, usage reports, campaign groups, search metadata, and completed-task variable updates.
- `cc flow import`/`export` now use the current FlowV2 format. Use `export-legacy` for raw FDL; `import-legacy` upload handling remains incomplete. `functions import` also needs multipart upload handling.
- `cc tasks resume` still means voice unhold. Digital tasks use `pause-digital`/`resume-digital`.
- Messaging consolidates HDS operations under `hds`, with `hybrid-data-security` as an alias. Admin and Meetings add body-based recording searches; Meetings also adds group service-app operations.

See the [full refresh inventory](docs/command-inventory.md), [every renamed command](docs/command-migration.md), and [Contact Center skill](skill/cc/SKILL.md) for migration details and examples. Upstream collections remain unchanged; overrides control names, routes, and product placement.

## Authentication

- **OAuth PKCE flow** — `webex login` opens a browser, no client secret needed on the user side
- **Device login** — `webex login --device` prints a URL and code to approve from any browser; used automatically over SSH, in CI, on Linux without a display, or when no browser opens
- **OS keyring storage** — tokens stored securely in macOS Keychain / Linux keyring / Windows Credential Manager, with a plain-text file fallback when no keyring is available
- **Auto-refresh** — expired tokens are refreshed automatically
- **Multi-user** — log in with multiple Webex accounts and switch between them

```bash
webex login                    # Login (opens browser, or device login when headless)
webex login --device           # Login by approving a code from any browser
webex login --browser          # Force the browser flow
webex login --read-only        # Login with read scopes only and turn on read-only mode
webex logout                   # Remove stored tokens
webex auth status              # Show current user, token status, and token store
webex auth list                # List all authenticated users
webex auth switch <email>      # Switch default user
webex auth set-org <org-id>    # Set a persistent org override
webex auth clear-org           # Clear the org override
webex auth export              # Print the refresh token for $WEBEX_REFRESH_TOKEN
```

Token resolution order: `--token` flag > `$WEBEX_TOKEN` env var > `$WEBEX_REFRESH_TOKEN` env var > stored login.

### Headless and Unattended Use

Device login needs no local browser: `webex login` shows `Open <url> and enter <code>` (and a QR code on a terminal), and you approve from a browser on any device. With `--output json`, the URL, code, and expiry are printed as JSON on stdout so an agent can relay them. `WEBEX_LOGIN_MODE=device|browser` picks the flow without a flag. A custom integration needs the redirect URIs `https://oauth-helper-{a,r,k,d}.wbx2.com/helperservice/v1/actions/device/callback` for device login.

Without an OS keyring (containers, SSH sessions, CI), tokens are saved in plain text to `$XDG_CONFIG_HOME/webex-cli/credentials.json` (`~/.config/webex-cli/` by default, `%APPDATA%\webex-cli\` on Windows) with mode `0600`. `WEBEX_TOKEN_STORE=keyring|file` forces a store; `webex auth status` shows which one holds the token.

For CI and unattended agents where nobody can approve a code:

- `WEBEX_TOKEN` takes a personal access token (12 hours) for short jobs.
- `WEBEX_REFRESH_TOKEN` takes a refresh token; the CLI mints access tokens from it as needed. Create one by logging in on your laptop and running `webex auth export`. Set `WEBEX_CLIENT_ID` and `WEBEX_CLIENT_SECRET` if the token came from your own integration. Webex refresh tokens last about 90 days from last use. The CLI caches the access token and any rotated refresh token in the credentials file when it is writable; otherwise each process mints a new access token.
- For organization-level automation, prefer a [service app](https://developer.webex.com/docs/service-apps) or a bot token over a person's login.

### Read-only Mode

`webex login --read-only` requests only read scopes, so Webex itself rejects writes made with the token. Use it when an AI agent should read but not change your Webex data. While read-only mode is on:

- Stored logins with write access are deleted from the keyring, because any program running as you can read it.
- `PUT`, `PATCH`, `DELETE`, and uploads are refused before they are sent. `POST` is allowed only for a fixed list of query endpoints that take their filters in the request body, such as `cc search` and recording queries.
- `--token`, `$WEBEX_TOKEN`, and `$WEBEX_REFRESH_TOKEN` are refused, and `auth switch`, `--user`, and folder defaults accept only users who logged in with `--read-only`.
- `webex mcp serve` does not register `webex_write`.

Leaving read-only mode requires `webex login` from an interactive terminal, with a confirmation prompt. `WEBEX_READ_ONLY=1` turns the same checks on for one process and never turns them off. If your OAuth integration lacks some of the default read scopes, set your own list with `webex config set read-only-scopes "<scopes>"`.

### Organization Override

Partner admins managing customer orgs can set a persistent default org so they don't need `--organization` on every command:

```bash
# Set a default org (accepts UUID or base64 ID, validates via API)
webex auth set-org <org-id>

# All subsequent commands use the override org
webex calling devices list
webex cc site list

# The --organization flag still takes priority for one-off commands
webex calling people list --organization <other-org-id>

# Clear the override to revert to your home org
webex auth clear-org
```

Org resolution order: `--organization` flag > `auth set-org` override > login user's home org.

## Output Formats

Control output with `--output`:

| Format | Description |
|--------|-------------|
| `json` | Pretty-printed JSON (default) |
| `table` | ASCII table with auto-detected columns and terminal-width formatting |
| `csv` | CSV with headers |
| `raw` | Raw API response |

```bash
webex calling people list --output table
webex cc users list --output csv > users.csv
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--token <token>` | Override authentication |
| `--user <email>` | Use a specific authenticated user |
| `--organization <orgId>` | Override org ID |
| `--output json\|table\|csv\|raw` | Output format (default: json) |
| `--debug` | Show HTTP request/response details |
| `--paginate` | Auto-paginate list results |
| `--dry-run` | Print write requests without executing them |

## Configuration

```bash
webex config set client-id <id>          # Use custom OAuth client ID
webex config set client-secret <secret>  # Use custom OAuth client secret
webex config set scopes <scopes>         # Override OAuth scopes
webex config set read-only-scopes <scopes>  # Override scopes for login --read-only
webex config get client-id               # View current value
```

Config is stored in `~/.webex-cli/config.json`.

## MCP Server

> **Claude Desktop users:** see [CLAUDE_DESKTOP.md](CLAUDE_DESKTOP.md) for installation, DXT extension setup, and skill configuration specific to Claude Desktop.


`webex mcp serve` starts a [Model Context Protocol](https://modelcontextprotocol.io/) server over stdio, letting AI clients query and manage your Webex environment directly.

```bash
# Register with Claude Code
claude mcp add webex -- webex mcp serve

# Or register with a specific binary path
claude mcp add webex -- /path/to/webex mcp serve
```

The server exposes 4 tools:

| Tool | API methods | Description |
|---|---|---|
| `webex_read` | GET | Execute a read-only CLI command (list, get, download, export, search) |
| `webex_write` | POST / PUT / PATCH / DELETE | Execute a command that creates, updates, or deletes a resource |
| `webex_help` | — | Get help text for any command or command group |
| `webex_usage` | — | Query the MCP usage log (recent commands, timing, status) |

Read and write operations are split into separate tools so that `webex_read` can be auto-approved in permissions while `webex_write` always prompts for confirmation. In [read-only mode](#read-only-mode) the server does not register `webex_write`.

And 2 MCP resources:

| Resource | Description |
|---|---|
| `webex://commands` | JSON array of all available CLI commands with short descriptions |
| `webex://usage` | Last 50 raw lines of the usage log |

Rather than exposing fixed per-API tools, the dispatcher pattern lets AI clients invoke the full CLI surface via `webex_read` and `webex_write` — the command tree expands automatically as new commands are added.

Auth is shared with the CLI — run `webex login` once and the MCP server uses the same stored credentials. Token refresh is handled automatically.

### Usage Log

All dispatcher invocations are logged to `~/.webex-mcp/usage.log` (JSONL). Configure with serve flags:

```bash
webex mcp serve --log-path /tmp/webex.log --log-max-size 10485760 --log-max-files 5
```

| Flag | Default | Description |
|---|---|---|
| `--log-path` | `~/.webex-mcp/usage.log` | Log file path |
| `--log-max-size` | `5242880` (5 MB) | Max file size before rotation |
| `--log-max-files` | `3` | Number of rotated files to keep |

### Team Distribution via Claude Admin Portal

Package the extension once and upload it to the Claude admin console — Claude Desktop installs it automatically for all team members. Each user's local `webex` installation and credentials are used; no shared server is needed.

```bash
make extension   # produces webex-mcp.dxt
```

1. Upload `webex-mcp.dxt` to the Claude admin portal under **Extensions**
2. Enable it for your team or org
3. Team members must have `webex` installed and authenticated (`webex login`)

Alternatively, team members can install the extension locally by double-clicking `webex-mcp.dxt` in Finder.

## Coding Agent Skill

A set of skill files in `skill/` enables AI coding agents (Claude Code, Claude Cowork, OpenAI Codex, Cursor) to use the CLI via natural language. The root skill (`skill/SKILL.md`) covers auth, command structure, and global flags. Six per-area sub-skills provide comprehensive flag and body-schema documentation for each CLI area:

| Area | Sub-skill |
|---|---|
| Admin | `skill/admin/SKILL.md` |
| Calling | `skill/calling/SKILL.md` |
| Contact Center | `skill/cc/SKILL.md` |
| Devices | `skill/device/SKILL.md` |
| Meetings | `skill/meetings/SKILL.md` |
| Messaging | `skill/messaging/SKILL.md` |

Each sub-skill also contains an auto-generated **Command Reference** section (updated by `make codegen`) that lists every command and its flags, keeping documentation in sync with the API as Postman collections change.

The skill files are embedded in the binary, so the installed skill always matches the installed version. `webex post-install` offers to install them, and `webex update` updates installed copies to match the new binary.

See the [docs](https://cloverhound.github.io/webex-cli/agent-skill/) for manual setup instructions.

## Development

See [AGENTS.md](AGENTS.md) for project structure and development workflow.

Commands in `cmd/calling/`, `cmd/cc/`, and the other area packages are **generated** from Postman collections — do not edit by hand. Run `make refresh` to re-download collections, regenerate Go files, update skill documentation, and rebuild. See the [code generation pipeline](AGENTS.md#code-generation-pipeline) for details.

## License

MIT
