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
webex auth status          # Show current user, org, token expiry
webex auth list            # List all stored users
webex auth switch <email>  # Change default user
webex auth set-org <orgId> # Set a persistent org override (partner admins)
webex auth clear-org       # Clear the org override
webex login                # OAuth login (opens browser)
webex logout [email]       # Remove stored credentials
```

If not logged in, use `--token <TOKEN>` or set `$WEBEX_TOKEN`.

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

## Admin API Examples

Admin commands manage people, organizations, licenses, roles, and other administrative resources:

```bash
# List people in the org
webex admin people list --max 10

# Get a specific person
webex admin people get-person --person-id "PERSON_ID"

# Get my own details
webex admin people get-my-own

# List organizations
webex admin organizations list

# List licenses
webex admin licenses list

# List roles
webex admin roles list

# List events
webex admin events list

# List reports
webex admin reports list

# List recordings
webex admin recordings list
```

## Calling API Examples

Calling commands — no orgid required (it's inferred from the token):

```bash
# List people
webex calling people list --max 10

# Get a specific person
webex calling people get-person --person-id "PERSON_ID"

# List locations
webex calling locations list

# List devices
webex calling devices list

# Get my own details
webex calling people get-my-own
```

## Contact Center API Examples

When logged in, `--orgid` is auto-populated so you can omit it:

```bash
# List sites
webex cc site list

# List queues
webex cc contact-service-queue list

# List entry points
webex cc entry-point list

# List teams
webex cc team list

# List users
webex cc users list

# List global variables
webex cc global-variables list

# List business hours
webex cc business-hour list

# List dial number mappings (exception: uses list-dialed-mapping)
webex cc dial-number list-dialed-mapping

# List skills
webex cc skill list

# List skill profiles
webex cc skill-profile list

# List desktop profiles
webex cc desktop-profile list

# List desktop layouts
webex cc desktop-layout list

# List auxiliary codes (idle + wrap-up)
webex cc auxiliary-code list

# List multimedia profiles
webex cc multimedia-profile list

# List dial plans
webex cc dial-plan list

# List user profiles
webex cc user-profiles list

# List address books
webex cc address-book list

# List holiday lists
webex cc holiday-list list

# List outdial ANIs
webex cc outdial-ani list
```

## Device API Examples

Device commands manage Webex devices, workspaces, and configurations:

```bash
# List devices
webex device devices list

# Get a specific device
webex device devices get --device-id "DEVICE_ID"

# List workspaces
webex device workspaces list

# Get workspace details
webex device workspaces get --workspace-id "WORKSPACE_ID"

# List device configurations
webex device device-configurations list --device-id "DEVICE_ID"

# List workspace locations
webex device workspace-locations list

# Execute xAPI command
webex device xapi execute-command --device-id "DEVICE_ID" --body '{"command":"..."}'

# Query xAPI status
webex device xapi query-status --device-id "DEVICE_ID"
```

## Meetings API Examples

Meetings commands manage meeting scheduling, recordings, participants, and more:

```bash
# List meetings
webex meetings meetings list

# Get a specific meeting
webex meetings meetings get --meeting-id "MEETING_ID"

# Create a meeting
webex meetings meetings create --body '{"title":"Stand-up","start":"...","end":"..."}'

# List meeting participants
webex meetings participants list --meeting-id "MEETING_ID"

# List recordings
webex meetings recordings list

# List transcripts
webex meetings transcripts list

# List meeting preferences
webex meetings preferences list-sites

# List session types
webex meetings session-types list --site-url "SITE_URL"

# List tracking codes
webex meetings meetings list-tracking-codes --site-url "SITE_URL"
```

## Messaging API Examples

Messaging commands manage rooms, messages, teams, and webhooks:

```bash
# List rooms
webex messaging rooms list --max 10

# Get room details
webex messaging rooms get --room-id "ROOM_ID"

# Create a room
webex messaging rooms create --body '{"title":"My Room"}'

# List messages in a room
webex messaging messages list --room-id "ROOM_ID"

# Send a message
webex messaging messages create --body '{"roomId":"ROOM_ID","text":"Hello!"}'

# List teams
webex messaging teams list

# List team memberships
webex messaging team-memberships list --team-id "TEAM_ID"

# List webhooks
webex messaging webhooks list

# List room memberships
webex messaging memberships list --room-id "ROOM_ID"
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

## When Answering Questions

1. **Check auth first** with `webex auth status` to confirm the user is logged in
2. **Use `--help`** on any command you're not sure about before running it
3. **Write to temp files** when gathering data, then read the files to analyze
4. **Use `--paginate`** for list commands when you need all results
