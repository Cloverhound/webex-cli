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
webex calling call-queue list

# Contact Center
webex cc site list
webex cc team list
webex cc agents list
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

45 resource groups including auto-attendants, call queues, hunt groups, call controls, call routing (dial plans, route groups, trunks), DECT devices, emergency services, locations, numbers, paging groups, people, workspaces, voicemail, recordings, and more.

### Contact Center (`webex cc`)

54 resource groups including agents, queues, entry points, flows, skills, desktop layouts, campaigns, callbacks, realtime stats, AI assistant, journey analytics, subscriptions, and more.

### Admin (`webex admin`)

39 resource groups including people, licenses, organizations, roles, groups, events, reports, recordings, SCIM 2.0 (users/groups/schemas), hybrid clusters/connectors, security audit, service apps, and more.

### Devices (`webex device`)

9 resource groups including devices, device configurations, workspaces, workspace locations/metrics/personalization, hot-desking, and xAPI (execute commands, query status).

### Meetings (`webex meetings`)

22 resource groups including meetings, participants, recordings, transcripts, summaries, polls, Q&A, chats, invitees, preferences, session types, tracking codes, video mesh, and more.

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
webex cc agents list --output csv > agents.csv
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

## Coding Agent Skill

A skill file is included in `skill/SKILL.md` that enables AI coding agents (Claude Code, Claude Cowork, OpenAI Codex, Cursor) to query and manage your Webex environment.

The installer and `webex post-install` command will offer to install the skill automatically. Skills are also kept up to date when you run `webex update`.

See the [docs](https://cloverhound.github.io/webex-cli/agent-skill/) for manual setup instructions.

## Development

See [CLAUDE.md](CLAUDE.md) for project structure and development workflow.

Commands in `cmd/calling/` and `cmd/cc/` are **generated** from Postman collections — do not edit by hand. See the [code generation pipeline](CLAUDE.md#code-generation-pipeline) for details.

## License

MIT
