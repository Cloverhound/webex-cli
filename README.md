# Webex CLI

A command-line tool for Webex APIs — Calling, Contact Center, Admin, Devices, Meetings, and Messaging.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/Cloverhound/webex-cli/main/install.sh | sh
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
webex admin organizations get --orgId <id>

# Devices
webex device devices list
webex device xapi execute-command --deviceId <id> --commandName <name>

# Meetings
webex meetings meetings list
webex meetings recordings list

# Messaging
webex messaging rooms list
webex messaging messages list --roomId <id>
```

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
```

Token resolution order: `--token` flag > `$WEBEX_TOKEN` env var > OS keyring.

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

## Configuration

```bash
webex config set client-id <id>          # Use custom OAuth client ID
webex config set client-secret <secret>  # Use custom OAuth client secret
webex config set scopes <scopes>         # Override OAuth scopes
webex config get client-id               # View current value
```

Config is stored in `~/.webex-cli/config.json`.

## Documentation

See the [docs site](https://cloverhound.github.io/webex-cli/) for full API reference and guides.

## Claude Code Integration

A sample Claude Code skill is included in `skill/SKILL.md`. See the [docs](https://cloverhound.github.io/webex-cli/claude-skill/) for setup instructions.

## Development

See [CLAUDE.md](CLAUDE.md) for project structure and development workflow.

Commands in `cmd/calling/` and `cmd/cc/` are **generated** from Postman collections — do not edit by hand. See the [code generation pipeline](CLAUDE.md#code-generation-pipeline) for details.

## License

MIT
