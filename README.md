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

47 resource groups including auto-attendants, call queues (CxE), hunt groups, call controls, call routing (dial plans, route groups, trunks), DECT devices, emergency services, locations, numbers, paging groups, people, workspaces, voicemail, converged recordings, and more.

### Contact Center (`webex cc`)

54 resource groups including sites, queues, entry points, teams, flows, skills, desktop layouts, global variables, business hours, auxiliary codes, campaigns, callbacks, realtime stats, AI assistant, journey analytics, subscriptions, and more.

### Admin (`webex admin`)

39 resource groups including people, licenses, organizations, roles, groups, events, reports, recordings, SCIM 2.0 (users/groups/schemas), hybrid clusters/connectors, security audit, service apps, and more.

### Devices (`webex device`)

10 resource groups including devices, device configurations, workspaces, workspace locations/metrics/personalization, hot-desking, and xAPI (execute commands, query status).

### Meetings (`webex meetings`)

23 resource groups including meetings, participants, recordings, transcripts, summaries, polls, Q&A, chats, invitees, preferences, session types, tracking codes, video mesh, and more.

### Messaging (`webex messaging`)

12 resource groups including rooms, messages, memberships, teams, team memberships, webhooks, events, attachment actions, room tabs, and more.

## Authentication

- **OAuth PKCE flow** — `webex login` opens a browser, no client secret needed on the user side
- **OS keyring storage** — tokens stored securely in macOS Keychain / Linux keyring / Windows Credential Manager
- **Auto-refresh** — expired tokens are refreshed automatically
- **Multi-user** — log in with multiple Webex accounts and switch between them

```bash
webex login                    # Login (opens browser)
webex logout                   # Remove stored tokens
webex auth status              # Show current user and token status
webex auth list                # List all authenticated users
webex auth switch <email>      # Switch default user
webex auth set-org <org-id>    # Set a persistent org override
webex auth clear-org           # Clear the org override
```

Token resolution order: `--token` flag > `$WEBEX_TOKEN` env var > OS keyring.

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
webex config get client-id               # View current value
```

Config is stored in `~/.webex-cli/config.json`.

## MCP Server

`webex mcp serve` starts a [Model Context Protocol](https://modelcontextprotocol.io/) server over stdio, letting AI clients query and manage your Webex environment directly.

```bash
# Register with Claude Code
claude mcp add webex -- webex mcp serve

# Or register with a specific binary path
claude mcp add webex -- /path/to/webex mcp serve
```

The server exposes 10 tools:

| Tool | Description |
|---|---|
| `webex_list_spaces` | List spaces/rooms with optional keyword and type filters |
| `webex_get_messages` | Retrieve recent messages from a space by ID or name |
| `webex_search_messages` | Fan-out keyword search across recently-active spaces |
| `webex_send_message` | Send a message to a space or person (the only write tool) |
| `webex_list_direct_messages` | List DM conversations with a last-message preview |
| `webex_get_person` | Look up a user by email or display name |
| `webex_list_meetings` | List upcoming and recent meetings by date range |
| `webex_get_meeting_transcript` | Retrieve the transcript for a completed meeting |
| `webex_get_meeting_summary` | Retrieve the AI-generated summary for a completed meeting |
| `webex_list_teams` | List teams the authenticated user belongs to |

Auth is shared with the CLI — run `webex login` once and the MCP server uses the same stored credentials. Token refresh is handled automatically.

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

`webex post-install` offers to install all skill files automatically. They are also updated by `webex update`.

See the [docs](https://cloverhound.github.io/webex-cli/agent-skill/) for manual setup instructions.

## Development

See [CLAUDE.md](CLAUDE.md) for project structure and development workflow.

Commands in `cmd/calling/`, `cmd/cc/`, and the other area packages are **generated** from Postman collections — do not edit by hand. Run `make refresh` to re-download collections, regenerate Go files, update skill documentation, and rebuild. See the [code generation pipeline](CLAUDE.md#code-generation-pipeline) for details.

## License

MIT
