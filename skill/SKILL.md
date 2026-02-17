---
name: webex-cli
description: "Webex CLI: query and manage Webex Calling and Contact Center via the `webex` command-line tool. Use for listing resources, checking configurations, debugging API calls, and administering Webex environments."
argument-hint: "[command or resource-name]"
allowed-tools: Bash, Read, Grep, Glob
user-invocable: true
---

# Webex CLI Skill

This skill uses the `webex` CLI tool to interact with Webex Calling and Contact Center APIs.

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
webex login                # OAuth login (opens browser)
webex logout [email]       # Remove stored credentials
```

If not logged in, use `--token <TOKEN>` or set `$WEBEX_TOKEN`.

## Command Structure

```
webex calling <resource> <action> [flags]     # Webex Calling APIs
webex cc <resource> <action> [flags]          # Contact Center APIs
```

### Global Flags
- `--token <token>` — Override authentication
- `--user <email>` — Use a specific authenticated user
- `--organization <orgId>` — Override org ID for this command
- `--output json|table|raw` — Output format (default: json)
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
   webex cc <resource> --help
   ```

## Calling API Examples

Calling commands are simpler — no orgid required (it's inferred from the token):

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

## Filtering and Pagination (CC)

CC list endpoints support RSQL filtering:
```bash
webex cc site list --filter='name=="Site A"'
webex cc team list --search="Sales"
```

Use `--paginate` to auto-paginate and get all results:
```bash
webex cc entry-point list --paginate
```

Or paginate manually:
```bash
webex cc entry-point list --page=0 --page-size=100
webex cc entry-point list --page=1 --page-size=100
```

## Output Handling

```bash
# Pretty JSON (default)
webex cc site list

# Table format
webex cc site list --output=table

# Raw JSON (no formatting)
webex cc site list --output=raw

# Save to file for processing
webex cc site list > /tmp/sites.json
```

## When Answering Questions

1. **Check auth first** with `webex auth status` to confirm the user is logged in
2. **Use `--help`** on any command you're not sure about before running it
3. **Write to temp files** when gathering data, then read the files to analyze
4. **Use `--paginate`** for list commands when you need all results
