---
name: webex-cli
description: "Webex CLI: query and manage Webex Admin, Calling, Contact Center, Devices, Meetings, and Messaging APIs via the `webex` command-line tool. Use for listing resources, checking configurations, debugging API calls, and administering Webex environments."
argument-hint: "[command or resource-name]"
allowed-tools: Bash, Read, Grep, Glob
user-invocable: true
---

# Webex CLI Skill

This skill uses the `webex` CLI tool to interact with Webex APIs — Admin, Calling, Contact Center, Devices, Meetings, and Messaging.

**Setup:** Install the CLI via `curl -fsSL https://raw.githubusercontent.com/Cloverhound/webex-cli/main/install.sh | sh`, or set the path to your local build below.

**Binary path** (update to match your installation):
```bash
webex
```

## Authentication

The CLI supports OAuth login with tokens stored in the OS keyring.

```bash
webex auth status                        # Show current user, org, token expiry
webex auth list                          # List all stored users
webex auth switch <email>                # Change default user
webex auth set-folder-default <email>    # Set default user for current folder
webex auth clear-folder-default          # Remove folder default
webex auth set-org <orgId>               # Set a persistent org override (partner admins)
webex auth clear-org                     # Clear the org override
webex login                              # OAuth login (opens browser)
webex logout [email]                     # Remove stored credentials
```

If not logged in, use `--token <TOKEN>` or set `$WEBEX_TOKEN`.

**Per-folder defaults:** Different folders can be associated with different Webex users via `auth set-folder-default`. When a folder default is set, that user's credentials are used automatically when running commands from that directory. This is useful when different project folders connect to different Webex orgs.

Token resolution order: `--token` flag > `WEBEX_TOKEN` env var > `--user` flag > `WEBEX_USER` env var > folder default (`.webex-cli/config.json`) > global default > OS keyring.

Org resolution order: `--organization` flag > `auth set-org` override > login user's home org.

## Command Structure

```
webex admin <resource> <action> [flags]      # Admin APIs (people, orgs, licenses, roles)
webex calling <resource> <action> [flags]    # Webex Calling APIs
webex cc <resource> <action> [flags]         # Contact Center APIs
webex device <resource> <action> [flags]     # Device APIs (devices, workspaces, xAPI)
webex meetings <resource> <action> [flags]   # Meetings APIs (scheduling, recordings)
webex messaging <resource> <action> [flags]  # Messaging APIs (rooms, messages, teams)
```

Aliases: `devices` for `device`, `meeting` for `meetings`, `msg` for `messaging`, `contact-center` for `cc`.

### Global Flags
- `--token <token>` — Override authentication
- `--user <email>` — Use a specific authenticated user
- `--organization <orgId>` — Override org ID for this command
- `--output json|table|csv|raw` — Output format (default: json)
- `--debug` — Show HTTP request/response details
- `--paginate` — Auto-paginate list results
- `--dry-run` — Print write requests (POST/PUT/DELETE/PATCH) without executing them; read operations still run normally

## Contact Center `--orgid` Handling

CC commands require `--orgid` (a per-command flag). **This is auto-populated from the logged-in user's org**, so you typically don't need to pass it manually.

The CLI auto-decodes base64-encoded Webex org IDs to UUID format. Both of these work:
```bash
webex cc site list                                    # Auto-populated from login
webex cc site list --orgid="4ebc486d-ff5f-..."        # Explicit UUID
webex cc site list --orgid="Y2lzY29zcGFyazovL3Vz..."  # Base64 — auto-decoded
```

If you need to override the org, use `--organization <orgId>` (global flag) which feeds into `--orgid` automatically. Both base64 and UUID formats are accepted.

## CC Subcommand Names

Most CC resources use a consistent `list` subcommand. A few exceptions remain:

| Resource | List command | Notes |
|---|---|---|
| `dial-number` | `list-dialed-mapping` | Exception — descriptive name |
| `agents` | (none) | Agent operations: login, logout, state-change, etc. |
| `flow` | (none) | Only export/import/publish |

All other CC resources (site, team, users, global-variables, business-hour, audio-files, work-types, etc.) use `list`.

**Always check `--help` first** if a command fails with "unknown command".

## Shell Usage Best Practices

1. **Always redirect stderr separately** when capturing JSON output:
   ```bash
   webex cc site list --orgid="$ORG" > /tmp/result.json 2>/tmp/error.log
   ```

2. **Use `--orgid=VALUE` syntax** (with `=`) for CC commands to avoid shell quoting issues. Do NOT use `--orgid "$VAR"` with a space — use `--orgid="$VAR"`.

3. **Write output to temp files first**, then read/analyze. Do NOT pipe webex output directly into python or jq in a single shell command — complex pipes can cause issues with the binary output.

4. **Check `--help` before guessing** subcommand names:
   ```bash
   webex <api> <resource> --help
   ```

## Sub-Skills

For detailed flags, body schemas, and usage examples, Read the sub-skill file before working in that area:

| Area | Installed path |
|---|---|
| Admin | `~/.claude/skills/webex-cli/admin/SKILL.md` |
| Calling | `~/.claude/skills/webex-cli/calling/SKILL.md` |
| Contact Center | `~/.claude/skills/webex-cli/cc/SKILL.md` |
| Devices | `~/.claude/skills/webex-cli/device/SKILL.md` |
| Meetings | `~/.claude/skills/webex-cli/meetings/SKILL.md` |
| Messaging | `~/.claude/skills/webex-cli/messaging/SKILL.md` |

## Quick Reference

```bash
# Admin — people, orgs, licenses
webex admin people list --max 10
webex admin people get-my-own
webex admin licenses list
webex admin organizations list

# Calling — locations, devices, queues, recordings
webex calling locations list
webex calling devices list
webex calling call-queue list
webex calling converged-recordings list --last 720h
webex calling converged-recordings download --recording-id <id> --output call.mp3

# Contact Center — sites, queues, entry points, flows
webex cc site list
webex cc contact-service-queue list
webex cc entry-point list
webex cc flow export --id <id> --output flow.json
webex cc audio-files download --id <id> --output prompt.wav
webex cc audio-files upload --file prompt.wav --name "Main Greeting"

# Device — devices, workspaces, xAPI
webex device devices list
webex device workspaces list
webex device xapi execute-command --device-id <id> --body '{"command":"..."}'

# Meetings — schedule, recordings, transcripts
webex meetings meetings list
webex meetings recordings list
webex meetings recordings download --recording-id <id> --output meeting.mp4 --type recording
webex meetings recordings download --recording-id <id> --output meeting.vtt --type transcript
webex meetings transcripts list

# Messaging — see sub-skill for full reference
webex messaging rooms list --type group --last 24h
webex messaging messages list --room-id <roomId>
webex messaging messages create --body '{"roomId":"<roomId>","text":"Hello!"}'
```

## Filtering and Pagination

CC list endpoints support RSQL filtering:
```bash
webex cc site list --filter='name=="Site A"'
webex cc team list --search="Sales"
```

Use `--paginate` to auto-paginate and get all results:
```bash
webex cc entry-point list --paginate
webex admin people list --paginate
webex meetings recordings list --paginate
```

Or paginate manually (CC):
```bash
webex cc entry-point list --page=0 --page-size=100
webex cc entry-point list --page=1 --page-size=100
```

Or paginate manually (Calling/Admin/Device/Meetings/Messaging — offset/max style):
```bash
webex admin people list --max=100
webex messaging rooms list --max=100
```

## Output Handling

```bash
# Pretty JSON (default)
webex admin people list

# Table format
webex admin people list --output=table

# Raw JSON (no formatting)
webex admin people list --output=raw

# Save to file for processing
webex admin people list > /tmp/people.json
```

## Converged Recordings: Admin vs Non-Admin Endpoints

The converged recordings API has two list endpoints with different access:

- **`list`** (`GET /convergedRecordings`) — Returns only the calling user's own recordings.
- **`list-admin-compliance-officer`** (`GET /admin/convergedRecordings`) — Returns recordings across the entire org. Requires `spark-admin:recordings_read` scope. Max 30-day time range.

The single-recording GET (`GET /convergedRecordings/{id}`) returns metadata for any recording the token can access. The `temporaryDirectDownloadLinks` field (needed for binary download) is included in the response only when the token has sufficient authorization — a user integration token may not receive download links for recordings owned by other users, while a service app token typically does.

```bash
# List recordings across the org (admin, max 30 days)
webex calling converged-recordings list-admin-compliance-officer --last 720h

# List only your own recordings
webex calling converged-recordings list --last 720h

# Download (requires temporaryDirectDownloadLinks in GET response)
webex calling converged-recordings download --recording-id <id> --output call.mp3
```

## When Answering Questions

1. **Check auth first** with `webex auth status` to confirm the user is logged in
2. **Use `--help`** on any command you're not sure about before running it
3. **Write to temp files** when gathering data, then read the files to analyze
4. **Use `--paginate`** for list commands when you need all results
