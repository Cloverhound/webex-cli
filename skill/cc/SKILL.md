---
name: webex-cli/cc
description: "Webex Contact Center commands: sites, queues, entry points, flows, agents, audio files, and configuration resources."
---

# Webex Contact Center

Commands: `webex cc <resource> <action> [flags]`  
Alias: `webex contact-center <resource> <action>`

**`--orgid` is auto-populated from the logged-in user's org.** Use `--orgid=<uuid>` (with `=`) only when overriding. Both UUID and base64 org ID formats are accepted.

## Resources

| Resource | Common operations |
|---|---|
| `site` | list, create, get-id, update-id, delete-id |
| `entry-point` | list, create, get-id, update-id, delete-id |
| `contact-service-queue` | list, create, get-id, update-id, delete-id |
| `team` | list, create, get-id, update-id, delete-id |
| `users` | list, get-id, get-along-profile-id, patch-id, update-id |
| `agents` | login, logout, state-change, get-activities, get-statistics |
| `flow` | list, export, import, publish |
| `audio-files` | list, get-id, create, update-id, delete-id |
| `global-variables` | list, create, get-id, update-id, delete-id |
| `business-hour` | list, create, get-id, update-id, delete-id |
| `auxiliary-code` | list, create, get-id, update-id, delete-id |
| `desktop-profile` | list, create, get-id, update-id, delete-id |
| `desktop-layout` | list, create, get-id, update-id, delete-id |
| `skill` | list, create, get-id, update-id, delete-id |
| `skill-profile` | list, create, get-id, update-id, delete-id |
| `multimedia-profile` | list, create, get-id, update-id, delete-id |
| `dial-plan` | list, create, get-id, update-id, delete-id |
| `dial-number` | list, list-dialed-mapping, get-id |
| `outdial-ani` | list, create, get-id, update-id, delete-id |
| `address-book` | list, create, get-id, update-id, delete-id |
| `holiday-list` | list, create, get-id, update-id, delete-id |

## RSQL Filtering (All Config Resources)

All config list commands support `--filter` (RSQL) and `--search` (keyword):

```bash
# Exact match
webex cc site list --filter='name=="Site A"'

# Not equal
webex cc site list --filter='name!="Site A"'

# In list
webex cc site list --filter='id=in=("<id1>","<id2>")'

# Keyword search
webex cc site list --search="Sales"

# Pagination
webex cc site list --page=0 --page-size=100
webex cc site list --paginate   # auto-paginate
```

**Filter value gotcha:** values with spaces must be quoted. Use `--filter=` (with `=`) to avoid shell quoting issues.

## Sites, Entry Points, Queues, Teams

These resources all follow the same pattern:

```bash
# List
webex cc site list [--filter <rsql>] [--search <text>] [--page <n>] [--page-size <n>]
webex cc entry-point list [same flags]
webex cc contact-service-queue list [same flags + --desktop-profile-filter true]
webex cc team list [same flags]

# Get by ID
webex cc site get-id --id <id>
webex cc entry-point get-id --id <id>
webex cc contact-service-queue get-id --id <id>
webex cc team get-id --id <id>

# Create / update / delete
webex cc site create --body '{...}'
webex cc site update-id --id <id> --body '{...}'
webex cc site delete-id --id <id>

# List references (what uses this resource)
webex cc site list-references --id <id>

# Bulk operations
webex cc site bulk-export
webex cc site bulk-save --body '[{...},{...}]'
```

## Users

```bash
# List users (agents / admins)
webex cc users list [flags]
  --filter <rsql>            # e.g. id=="<id>"
  --search <text>            # firstName, lastName, email
  --page <n>
  --page-size <n>
  --queue-id <id>            # filter by queue assignment
  --user-in-queue assigned|unassigned
  --supervisor-managed-agents-only true

# Get user with profile
webex cc users get-id --id <id>
webex cc users get-along-profile-id --id <id>
webex cc users get-ci-id --ci-id <ciUserId>

# Update user (agent settings, team, profile)
webex cc users patch-id --id <id> --body '{...}'
webex cc users update-id --id <id> --body '{...}'

# Bulk
webex cc users bulk-export
webex cc users bulk-partial-update --body '[{...}]'
```

## Agents (Live Session Control)

These operate on active agent desktop sessions — not config data.

```bash
# Login / logout desktop
webex cc agents login --body '{"dialNumber":"<ext>","roles":["AGENT"],"teamId":"<id>"}'
webex cc agents logout --body '{"logoutReason":"End of shift","agentId":"<id>"}'

# Change agent state (Available / Idle)
webex cc agents state-change --body '{"state":"AVAILABLE","channelType":["telephony"]}'
webex cc agents state-change --body '{"state":"IDLE","auxCodeId":"<id>","channelType":["telephony"]}'

# Activity and statistics (max 24-hour window per request)
webex cc agents get-activities --last 8h [--agent-ids <id1,id2>] [--channel-types telephony,chat]
webex cc agents get-statistics --last 8h --interval 15 [--agent-ids <id1,id2>]

# Buddy list
webex cc agents buddy-list --body '{"agentProfileId":"<id>","mediaType":"telephony","state":"AVAILABLE"}'
```

## Flows

```bash
# List flows
webex cc flow list [--partial-name-search <text>] [--flow-type FLOW|SUBFLOW] [--page <n>] [--size <n>]

# Export (download flow JSON)
webex cc flow export --flow-id <id> [--version draft|latest|<versionId>] > flow.json

# Import (upload flow JSON)
webex cc flow import --overwrite yes [--flow-type FLOW]
# Note: import reads the body from stdin or --body-file

# Publish a flow
webex cc flow publish --flow-id <id> --body '{"comment":"Promoting to live"}'
```

## Audio Files

```bash
# List
webex cc audio-files list [--filter <rsql>] [--search <text>]

# Get with download URL
webex cc audio-files get-id --id <id> --include-url true

# Create / upload
webex cc audio-files create --body '{"name":"Main Greeting","description":"..."}'
# After creating, upload the WAV via update
webex cc audio-files update-id --id <id> --body '{...}'

# Delete
webex cc audio-files delete-id --id <id>
```

## Configuration Resources

All of these (global-variables, business-hour, auxiliary-code, desktop-profile, desktop-layout, skill, skill-profile, multimedia-profile, dial-plan, outdial-ani, address-book, holiday-list) follow the same pattern:

```bash
webex cc <resource> list [--filter <rsql>] [--search <text>] [--page <n>] [--page-size <n>]
webex cc <resource> get-id --id <id>
webex cc <resource> create --body '{...}'
webex cc <resource> update-id --id <id> --body '{...}'
webex cc <resource> delete-id --id <id>
webex cc <resource> list-references --id <id>
webex cc <resource> bulk-export
webex cc <resource> bulk-save --body '[{...}]'
webex cc <resource> purge-inactive     # removes soft-deleted records (most resources)
```

### Dial Numbers (exception — no CRUD)

```bash
# List DN→EP mappings
webex cc dial-number list-dialed-mapping [--filter <rsql>]
webex cc dial-number list [--filter <rsql>]
webex cc dial-number get-id --id <id>
```

### Outdial ANI (has entries sub-resource)

```bash
webex cc outdial-ani list
webex cc outdial-ani get-id --id <id>
webex cc outdial-ani list-entry --id <id>          # list ANI entries under an outdial ANI
webex cc outdial-ani create-entry --body '{...}'
webex cc outdial-ani delete-entry-id --id <entryId>
```

### Address Book (has entries sub-resource)

```bash
webex cc address-book list
webex cc address-book get-id --id <id>
webex cc address-book list-entry --id <id>         # list contacts
webex cc address-book create-entry --body '{...}'
webex cc address-book bulk-save-entry --body '[{...}]'
webex cc address-book update-entry-id --id <entryId> --body '{...}'
webex cc address-book delete-entry-id --id <entryId>
```

## Key Gotchas

1. **`--orgid=VALUE` syntax** — always use `=` (equals) with `--orgid`, not a space: `--orgid="<uuid>"` not `--orgid "<uuid>"`. Shell quoting issues cause silent failures with the space form.

2. **RSQL values with spaces** — wrap values in double quotes inside the RSQL expression: `--filter='name=="Site A"'`.

3. **`get-id` not `get`** — most CC resources use `get-id --id <id>` (not `get --<resource>-id`). The `agents` resource is an exception: it has no `list` or `get-id` commands, only session control operations.

4. **`flow` has no CRUD** — `flow` only has `list`, `export`, `import`, `publish`. To create or edit flows, use Webex Contact Center Flow Designer in Control Hub.

5. **Agent statistics time window** — `agents get-activities` and `get-statistics` enforce a **24-hour** max window per request. Use multiple calls to cover longer periods.

6. **Purge vs delete** — `delete-id` soft-deletes. `purge-inactive` permanently removes soft-deleted records for resources that support it.
