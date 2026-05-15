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

<!-- codegen:start -->
## Command Reference

> Auto-generated from Postman collections. Run `make codegen` to update.

### agent-wellbeing

| Command | Flags |
|---|---|
| `get-burnout-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-burnout` | `--orgid` *(required)*, `--filter`, `--attributes`, `--page`, `--page-size` |
| `subscribe-realtime-burnout-events` | `--destination-url`, `--event-types`, `--name`, `--description`, `--secret`, `--org-id`, `--body`, `--body-file` |
| `record-realtime-burnout-events` | `--body`, `--body-file` |
| `update-burnout-id` | `--orgid` *(required)*, `--id` *(required)*, `--agent-inclusion-type`, `--enabled`, `--organization-id`, `--version`, `--wellness-break-reminders`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |

### auto-csat

| Command | Flags |
|---|---|
| `get-mapped-question-id` | `--orgid` *(required)*, `--auto-csat-id` *(required)*, `--id` *(required)* |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--page`, `--page-size` |
| `list-mapped-question` | `--orgid` *(required)*, `--auto-csat-id` *(required)*, `--filter`, `--attributes`, `--page`, `--page-size` |
| `create-mapped-question` | `--orgid` *(required)*, `--auto-csat-id` *(required)*, `--question-id`, `--questionnaire-id`, `--organization-id`, `--id`, `--version`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save-mapped-question` | `--orgid` *(required)*, `--auto-csat-id` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--agent-inclusion-type`, `--enabled`, `--selected-global-variable-id`, `--survey-data-source`, `--organization-id`, `--version`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete-mapped-question-id` | `--orgid` *(required)*, `--auto-csat-id` *(required)*, `--id` *(required)* |

### generated-summaries

| Command | Flags |
|---|---|
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--page`, `--page-size` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--organization-id`, `--version`, `--call-drop-summaries-enabled`, `--virtual-agent-transfer-summaries-enabled`, `--consult-transfer-summaries-enabled`, `--agent-inclusion-type`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |

### business-hour

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--sort`, `--include-count`, `--single-object-response` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### holiday-list

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--sort`, `--include-count`, `--single-object-response` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### overrides

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--sort`, `--include-latest-override`, `--single-object-response` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### address-book

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-entry-id` | `--orgid` *(required)*, `--address-book-id` *(required)*, `--id` *(required)* |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list-entry` | `--orgid` *(required)*, `--address-book-id` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--name`, `--parent-type`, `--organization-id`, `--id`, `--version`, `--description`, `--site-id`, `--body`, `--body-file` |
| `create-entry` | `--orgid` *(required)*, `--address-book-id` *(required)*, `--name`, `--number`, `--organization-id`, `--id`, `--version`, `--body`, `--body-file` |
| `bulk-save-entry` | `--orgid` *(required)*, `--address-book-id` *(required)*, `--body`, `--body-file` |
| `update-entry-id` | `--orgid` *(required)*, `--address-book-id` *(required)*, `--id` *(required)*, `--name`, `--number`, `--organization-id`, `--version`, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--name`, `--parent-type`, `--organization-id`, `--version`, `--description`, `--site-id`, `--body`, `--body-file` |
| `delete-entry-id` | `--orgid` *(required)*, `--address-book-id` *(required)*, `--id` *(required)* |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### audio-files

| Command | Flags |
|---|---|
| `get-id` | `--orgid` *(required)*, `--id` *(required)*, `--include-url` |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `patch-id` | `--orgid` *(required)*, `--id` *(required)*, `--description`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### auxiliary-code

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--desktop-profile-filter`, `--supervised-user-id` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--active`, `--default-code`, `--name`, `--work-type-code`, `--work-type-id`, `--organization-id`, `--id`, `--version`, `--description`, `--is-system-code`, `--burnout-inclusion`, `--system-default`, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--active`, `--default-code`, `--name`, `--work-type-code`, `--work-type-id`, `--organization-id`, `--version`, `--description`, `--is-system-code`, `--burnout-inclusion`, `--system-default`, `--body`, `--body-file` |
| `bulk-partial-update` | `--orgid` *(required)*, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### contact-number

| Command | Flags |
|---|---|
| `list-all` | `--orgid` *(required)* |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--number`, `--organization-id`, `--id`, `--version`, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--number`, `--organization-id`, `--version`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### contact-service-queue

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--desktop-profile-filter`, `--provisioning-view`, `--single-object-response` |
| `bulk-export` | `--orgid` *(required)*, `--type`, `--page`, `--page-size` |
| `list-skill-csqs-skill-profile` | `--orgid` *(required)*, `--id` *(required)* |
| `get-id` | `--orgid` *(required)*, `--id` *(required)*, `--agents-updated-info` |
| `list-csq-references-id` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list-agent-based` | `--orgid` *(required)*, `--userid` *(required)*, `--search`, `--page`, `--page-size` |
| `list-skill-based` | `--orgid` *(required)*, `--userid` *(required)*, `--search`, `--page`, `--page-size` |
| `list-team-based` | `--orgid` *(required)*, `--userid` *(required)*, `--search`, `--page`, `--page-size` |
| `list-internal-skill-csqs-profile` | `--orgid` *(required)*, `--id` *(required)* |
| `list-team-csqs-team-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-agent-csqs-ci-user-id` | `--orgid` *(required)*, `--ci-user-id` *(required)* |
| `list-skill-csqs-ci-user-id` | `--orgid` *(required)*, `--id` *(required)* |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `delete-csq-references` | `--orgid` *(required)*, `--body`, `--body-file` |
| `list-manually-assignable-csqs` | `--orgid` *(required)*, `--agent-id`, `--team-id`, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `create-remove-agents-users-agent` | `--orgid` *(required)*, `--id` *(required)*, `--add`, `--remove`, `--body`, `--body-file` |
| `list-mapping-summary-grouped-assistant-skill` | `--orgid` *(required)*, `--page`, `--page-size`, `--assistant-skill-ids`, `--body`, `--body-file` |
| `list-csqs-skills-profile` | `--orgid` *(required)*, `--body`, `--body-file` |
| `list-csqs-user-profile` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `bulk-partial-update` | `--orgid` *(required)*, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### desktop-layout

| Command | Flags |
|---|---|
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--single-object-response`, `--provisioning-view` |
| `create` | `--orgid` *(required)*, `--default-json-modified`, `--edited-by`, `--global`, `--json-file-content`, `--json-file-name`, `--name`, `--status`, `--validated`, `--organization-id`, `--id`, `--version`, `--description`, `--validated-time`, `--default-json-modified-time`, `--modified-time`, `--team-ids`, `--system-default`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--default-json-modified`, `--edited-by`, `--global`, `--json-file-content`, `--json-file-name`, `--name`, `--status`, `--validated`, `--organization-id`, `--version`, `--description`, `--validated-time`, `--default-json-modified-time`, `--modified-time`, `--team-ids`, `--system-default`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### desktop-profile

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--single-object-response` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### dial-number

| Command | Flags |
|---|---|
| `list-dialed-mapping` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--include-entry-point-name` |
| `bulk-export-dialed-mapping` | `--orgid` *(required)*, `--page`, `--page-size` |
| `list-dialed-dialed-mapping` | `--orgid` *(required)* |
| `get-dialed-mapping-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references-dialed-mapping` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create-dialed-mapping` | `--orgid` *(required)*, `--entry-point-id`, `--entry-point-name`, `--organization-id`, `--id`, `--version`, `--dialled-number`, `--extension`, `--routing-prefix`, `--esn`, `--route-point-id`, `--default-ani`, `--location`, `--region-id`, `--created-time`, `--last-updated-time`, `--dialled-number-digits`, `--body`, `--body-file` |
| `bulk-save-dialed-mapping` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-dialed-mapping-id` | `--orgid` *(required)*, `--id` *(required)*, `--entry-point-id`, `--entry-point-name`, `--organization-id`, `--version`, `--dialled-number`, `--extension`, `--routing-prefix`, `--esn`, `--route-point-id`, `--default-ani`, `--location`, `--region-id`, `--created-time`, `--last-updated-time`, `--dialled-number-digits`, `--body`, `--body-file` |
| `delete-all-dialed-mapping` | `--orgid` *(required)* |
| `delete-dialed-mapping-id` | `--orgid` *(required)*, `--id` *(required)* |

### dial-plan

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--active`, `--name`, `--regular-expression`, `--organization-id`, `--id`, `--version`, `--description`, `--prefix`, `--stripped-chars`, `--system-default`, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--active`, `--name`, `--regular-expression`, `--organization-id`, `--version`, `--description`, `--prefix`, `--stripped-chars`, `--system-default`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### entry-point

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--desktop-profile-filter`, `--provisioning-view`, `--include-count`, `--single-object-response` |
| `bulk-export` | `--orgid` *(required)*, `--type`, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)*, `--include-names` |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### global-variables

| Command | Flags |
|---|---|
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-reportable-count` | `--orgid` *(required)* |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--active`, `--agent-editable`, `--agent-viewable`, `--default-value`, `--name`, `--reportable`, `--variable-type`, `--organization-id`, `--id`, `--version`, `--description`, `--sensitive`, `--desktop-label`, `--system-default`, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--active`, `--agent-editable`, `--agent-viewable`, `--default-value`, `--name`, `--reportable`, `--variable-type`, `--organization-id`, `--version`, `--description`, `--sensitive`, `--desktop-label`, `--system-default`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### multimedia-profile

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### outdial-ani

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--single-object-response` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `list-entry` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `get-entry-id` | `--orgid` *(required)*, `--out-dial-ani-id` *(required)*, `--id` *(required)* |
| `list-entry-2` | `--orgid` *(required)*, `--out-dial-ani-id` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `create-entry` | `--orgid` *(required)*, `--out-dial-ani-id` *(required)*, `--name`, `--number`, `--organization-id`, `--id`, `--version`, `--default-anientry`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save-entry` | `--orgid` *(required)*, `--out-dial-ani-id` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `update-entry-id` | `--orgid` *(required)*, `--out-dial-ani-id` *(required)*, `--id` *(required)*, `--name`, `--number`, `--organization-id`, `--version`, `--default-anientry`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |
| `delete-entry-id` | `--orgid` *(required)*, `--out-dial-ani-id` *(required)*, `--id` *(required)* |

### site

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--active`, `--multimedia-profile-id`, `--name`, `--organization-id`, `--id`, `--version`, `--description`, `--system-default`, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--active`, `--multimedia-profile-id`, `--name`, `--organization-id`, `--version`, `--description`, `--system-default`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### skill

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--single-object-response` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--active`, `--name`, `--service-level-threshold`, `--skill-type`, `--organization-id`, `--id`, `--version`, `--description`, `--enum-skill-values`, `--dynamic-skill`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `populate-json-attributes-field-skill-id-org` | `--orgid` *(required)*, `--id` *(required)* |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--active`, `--name`, `--service-level-threshold`, `--skill-type`, `--organization-id`, `--version`, `--description`, `--enum-skill-values`, `--dynamic-skill`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### skill-profile

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--single-object-response` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)*, `--include-skill-details` |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--active-skills`, `--name`, `--organization-id`, `--id`, `--version`, `--description`, `--active-enum-skills`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--active-skills`, `--name`, `--organization-id`, `--version`, `--description`, `--active-enum-skills`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### team

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--supervisor-view`, `--provisioning-view`, `--single-object-response` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--active`, `--name`, `--rank-queues-for-team`, `--site-id`, `--team-status`, `--team-type`, `--organization-id`, `--id`, `--version`, `--dialed-number`, `--capacity`, `--desktop-layout-id`, `--skill-profile-id`, `--multi-media-profile-id`, `--user-ids`, `--description`, `--system-default`, `--queue-rankings`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--active`, `--name`, `--rank-queues-for-team`, `--site-id`, `--team-status`, `--team-type`, `--organization-id`, `--version`, `--dialed-number`, `--capacity`, `--desktop-layout-id`, `--skill-profile-id`, `--multi-media-profile-id`, `--user-ids`, `--description`, `--system-default`, `--queue-rankings`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### user-profiles

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)*, `--include-names` |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `get-acl-id` | `--orgid` *(required)*, `--id` *(required)*, `--names` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### users

| Command | Flags |
|---|---|
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--supervisor-managed-agents-only`, `--single-object-response`, `--buddy-team-agents-only`, `--user-in-queue`, `--queue-id`, `--include-aimapping-count`, `--include-dynamic-skills-limit-reached` |
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-ci-id` | `--orgid` *(required)*, `--id` *(required)*, `--include-user-profile`, `--include-names` |
| `list-along-profile` | `--orgid` *(required)* |
| `get-along-profile-id` | `--orgid` *(required)*, `--id` *(required)* |
| `get-id` | `--orgid` *(required)*, `--id` *(required)*, `--include-count`, `--include-user-profile-type`, `--include-skill-profile-audit`, `--include-reskill-audit-info`, `--include-skill-details` |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `get-dynamic-skill-id` | `--orgid` *(required)*, `--skill-id` *(required)*, `--search`, `--page`, `--page-size` |
| `get-agents-matching-skill-requirements` | `--orgid` *(required)*, `--search`, `--page`, `--page-size`, `--condition`, `--skill-id`, `--skill-value`, `--organization-id`, `--id`, `--version`, `--skill-name`, `--skill-type`, `--weight`, `--dynamic-skill`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `get-ids` | `--orgid` *(required)*, `--page`, `--page-size`, `--user-ids`, `--search`, `--queue-id`, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `bulk-partial-update` | `--orgid` *(required)*, `--body`, `--body-file` |
| `patch-id` | `--orgid` *(required)*, `--id` *(required)*, `--value-type`, `--body`, `--body-file` |
| `bulk-update-dynamic-skills` | `--orgid` *(required)*, `--skill-id` *(required)*, `--body`, `--body-file` |
| `reskill-agents` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |

### work-types

| Command | Flags |
|---|---|
| `bulk-export` | `--orgid` *(required)*, `--page`, `--page-size` |
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--active`, `--name`, `--work-type-code`, `--organization-id`, `--id`, `--version`, `--description`, `--system-default`, `--body`, `--body-file` |
| `bulk-save` | `--orgid` *(required)*, `--body`, `--body-file` |
| `purge-inactive` | `--orgid` *(required)*, `--next-start-id` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--active`, `--name`, `--work-type-code`, `--organization-id`, `--version`, `--description`, `--system-default`, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### data-sources

| Command | Flags |
|---|---|
| `get-all` | — |
| `get-schemas` | — |
| `get-schema` | `--schema-id` *(required)* |
| `get` | `--data-source-id` *(required)* |
| `register` | `--audience`, `--nonce`, `--schema-id`, `--subject`, `--token-lifetime-minutes`, `--url`, `--body`, `--body-file` |
| `update` | `--data-source-id` *(required)*, `--audience`, `--error-message`, `--nonce`, `--schema-id`, `--status`, `--subject`, `--token-lifetime-minutes`, `--url`, `--body`, `--body-file` |
| `delete` | `--data-source-id` *(required)* |

### estimated-wait-time

| Command | Flags |
|---|---|
| `get` | `--queue-id`, `--lookback-minutes`, `--max-cv`, `--min-valid-samples`, `--org-id` |

### notification

| Command | Flags |
|---|---|
| `subscribe` | `--is-keep-alive-enabled`, `--client-type`, `--allow-multi-login`, `--force`, `--body`, `--body-file` |

### queues

| Command | Flags |
|---|---|
| `get-statistics` | `--from`, `--to`, `--interval`, `--queue-ids`, `--org-id`, `--last` |

### realtime

| Command | Flags |
|---|---|
| `subscribe-notification` | `--is-keep-alive-enabled`, `--client-type`, `--allow-multi-login`, `--force`, `--body`, `--body-file` |

### subscriptions

| Command | Flags |
|---|---|
| `list-v1` | `--org-id` |
| `list-v2` | `--org-id` |
| `get-v1` | `--id` *(required)*, `--org-id` |
| `get-v2` | `--id` *(required)*, `--org-id` |
| `list-event-types-v1` | `--org-id` |
| `list-event-types-v2` | `--org-id` |
| `register-v1` | `--destination-url`, `--event-types`, `--name`, `--description`, `--secret`, `--org-id`, `--body`, `--body-file` |
| `register-v2` | `--destination-url`, `--event-types`, `--name`, `--resource-version`, `--description`, `--secret`, `--org-id`, `--body`, `--body-file` |
| `update-v1` | `--id` *(required)*, `--description`, `--event-types`, `--destination-url`, `--status`, `--secret`, `--org-id`, `--body`, `--body-file` |
| `update-v2` | `--id` *(required)*, `--resource-version`, `--description`, `--event-types`, `--destination-url`, `--status`, `--secret`, `--org-id`, `--body`, `--body-file` |
| `delete-v1` | `--id` *(required)*, `--org-id` |
| `delete-v2` | `--id` *(required)*, `--org-id` |

### agents

| Command | Flags |
|---|---|
| `get-activities` | `--agent-ids`, `--team-ids`, `--channel-types`, `--from`, `--to`, `--page-size`, `--page`, `--org-id`, `--last` |
| `get-statistics` | `--from`, `--to`, `--interval`, `--agent-ids`, `--org-id`, `--last` |
| `login` | `--dial-number`, `--roles`, `--team-id`, `--is-extension`, `--device-type`, `--device-id`, `--body`, `--body-file` |
| `reload` | — |
| `buddy-list` | `--agent-profile-id`, `--media-type`, `--state`, `--body`, `--body-file` |
| `logout` | `--logout-reason`, `--agent-id`, `--body`, `--body-file` |
| `state-change` | `--channel-type`, `--state`, `--aux-code-id`, `--reason`, `--agent-id`, `--body`, `--body-file` |

### call-monitoring

| Command | Flags |
|---|---|
| `get-sessions` | — |
| `create-request` | `--id`, `--monitor-type`, `--task-id`, `--queue-ids`, `--teams`, `--sites`, `--agents`, `--tracking-id`, `--invisible-mode`, `--body`, `--body-file` |
| `barge-in-request` | `--task-id` *(required)* |
| `end-request` | `--task-id` *(required)* |
| `hold-request` | `--task-id` *(required)* |
| `unhold-request` | `--task-id` *(required)* |
| `delete-request` | `--request-id` *(required)* |

### tasks

| Command | Flags |
|---|---|
| `get` | `--channel-types`, `--from`, `--to`, `--page-size`, `--org-id`, `--last` |
| `create` | `--body`, `--body-file` |
| `accept` | `--task-id` *(required)* |
| `end` | `--task-id` *(required)* |
| `wrap-up` | `--task-id` *(required)*, `--wrap-up-reason`, `--aux-code-id`, `--body`, `--body-file` |
| `hold` | `--task-id` *(required)*, `--media-resource-id`, `--body`, `--body-file` |
| `resume` | `--task-id` *(required)*, `--media-resource-id`, `--body`, `--body-file` |
| `reject` | `--task-id` *(required)*, `--media-resource-id`, `--body`, `--body-file` |
| `pause-recording` | `--task-id` *(required)* |
| `resume-recording` | `--task-id` *(required)*, `--auto-resumed`, `--body`, `--body-file` |
| `transfer` | `--task-id` *(required)*, `--to`, `--destination-type`, `--body`, `--body-file` |
| `consult-task` | `--task-id` *(required)*, `--to`, `--destination-type`, `--hold-participants`, `--body`, `--body-file` |
| `consult-conference` | `--task-id` *(required)*, `--to`, `--agent-id`, `--destination-type`, `--body`, `--body-file` |
| `consult-transfer` | `--task-id` *(required)*, `--to`, `--destination-type`, `--body`, `--body-file` |
| `consult-accept` | `--task-id` *(required)* |
| `assign` | `--task-id` *(required)* |
| `consult-end` | `--task-id` *(required)*, `--queue-id`, `--body`, `--body-file` |
| `exit-conference` | `--task-id` *(required)* |
| `accept-preview` | `--campaign-id` *(required)*, `--task-id` *(required)* |
| `skip-preview` | `--campaign-id` *(required)*, `--task-id` *(required)* |
| `delete-preview` | `--campaign-id` *(required)*, `--task-id` *(required)* |
| `update-2` | `--task-id` *(required)*, `--body`, `--body-file` |
| `update` | `--task-id` *(required)*, `--body`, `--body-file` |

### journey

| Command | Flags |
|---|---|
| `get-workspace` | `--workspace-id` *(required)* |
| `get-template-searched-template-id` | `--workspace-id` *(required)*, `--template-id` *(required)* |
| `get-wxcc-subscription` | `--workspace-id` *(required)* |
| `get-all-workspaces` | `--filter`, `--sort-by`, `--sort`, `--page`, `--page-size` |
| `get-all-template` | `--workspace-id` *(required)*, `--filter`, `--sort`, `--sort-by`, `--page`, `--page-size` |
| `get-all-person` | `--workspace-id` *(required)*, `--person-id`, `--filter`, `--sort-by`, `--sort`, `--page`, `--page-size` |
| `get-all-actions` | `--workspace-id` *(required)*, `--sort-by`, `--sort`, `--page`, `--page-size` |
| `get-template-searched-template-name` | `--workspace-id` *(required)*, `--template-name` *(required)* |
| `search-identity-aliases` | `--workspace-id` *(required)*, `--aliases` *(required)*, `--sort-by`, `--sort`, `--page`, `--page-size` |
| `get-all-actions-template` | `--workspace-id` *(required)*, `--template-id` *(required)* |
| `get-action-name` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--action-name` *(required)* |
| `get-action-actionid` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--action-id` *(required)* |
| `get-historic-profile-view-template-name` | `--workspace-id` *(required)*, `--person-id` *(required)*, `--template-name` *(required)* |
| `get-historic-profile-view` | `--workspace-id` *(required)*, `--person-id` *(required)*, `--template-id` *(required)* |
| `get-historic-profile-view-identity-template-name` | `--workspace-id` *(required)*, `--identity` *(required)*, `--template-name` *(required)* |
| `get-historic-profile-view-identity-template-id` | `--workspace-id` *(required)*, `--identity` *(required)*, `--template-id` *(required)* |
| `stream-profile-views-template-name` | `--workspace-id` *(required)*, `--identity` *(required)*, `--template-name` *(required)* |
| `stream-profile-views` | `--workspace-id` *(required)*, `--identity` *(required)*, `--template-id` *(required)* |
| `get-historic-events` | `--workspace-id` *(required)*, `--identity`, `--sort-by`, `--sort`, `--filter`, `--data`, `--page`, `--page-size` |
| `stream-events-identity` | `--workspace-id` *(required)*, `--identity` *(required)*, `--filter`, `--data` |
| `create-wxcc-subscription` | `--workspace-id` *(required)* |
| `create-workspace` | `--organization-id`, `--description`, `--name`, `--body`, `--body-file` |
| `create-template` | `--workspace-id` *(required)*, `--body`, `--body-file` |
| `create-person` | `--workspace-id` *(required)*, `--first-name`, `--last-name`, `--phone`, `--email`, `--temporary-id`, `--customer-id`, `--body`, `--body-file` |
| `merges-identities-primary-identity` | `--workspace-id` *(required)*, `--primary-person-id` *(required)*, `--person-ids-to-merge`, `--body`, `--body-file` |
| `creates-merges-aliases-individual-jds` | `--workspace-id` *(required)*, `--first-name`, `--last-name`, `--phone`, `--email`, `--temporary-id`, `--customer-id`, `--body`, `--body-file` |
| `create-action` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--body`, `--body-file` |
| `event-posting` | `--workspace-id`, `--body`, `--body-file` |
| `update-workspace` | `--workspace-id` *(required)*, `--description`, `--name`, `--body`, `--body-file` |
| `update-profileviewtemplate` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--body`, `--body-file` |
| `update-action` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--action-id` *(required)*, `--body`, `--body-file` |
| `create-remove-replace-person` | `--workspace-id` *(required)*, `--person-id` *(required)*, `--body`, `--body-file` |
| `create-one-more-identities-person` | `--workspace-id` *(required)*, `--person-id` *(required)*, `--phone`, `--email`, `--temporary-id`, `--customer-id`, `--body`, `--body-file` |
| `delete-one-more-identities-person` | `--workspace-id` *(required)*, `--person-id` *(required)*, `--body`, `--body-file` |
| `delete-workspace` | `--workspace-id` *(required)* |
| `delete-template-template-id` | `--workspace-id` *(required)*, `--template-id` *(required)* |
| `delete-person-id` | `--workspace-id` *(required)*, `--person-id` *(required)* |
| `delete-wxcc-subscription` | `--workspace-id` *(required)* |
| `delete-action-configuration-actionid` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--action-id` *(required)* |

### campaign-manager

| Command | Flags |
|---|---|
| `start-request` | `--body`, `--body-file` |
| `update-request` | `--campaign-id` *(required)*, `--dialing-rate`, `--dialing-list-fetch-url`, `--outdial-ani`, `--campaign-name`, `--auth-token`, `--no-answer-ring-limit`, `--max-dialing-rate`, `--reservation-percentage`, `--preview-offer-timeout`, `--preview-offer-timeout-auto-action`, `--preview-actions-disabled`, `--body`, `--body-file` |
| `stop-request` | `--campaign-id` *(required)* |

### captures

| Command | Flags |
|---|---|
| `list` | `--body`, `--body-file` |

### flow

| Command | Flags |
|---|---|
| `list` | `--org-id` *(required)*, `--project-id` *(required)*, `--flow-type`, `--ids`, `--page`, `--partial-name-search`, `--size`, `--include-pagination` |
| `export` | `--org-id` *(required)*, `--project-id` *(required)*, `--flow-id` *(required)*, `--version` |
| `publish` | `--org-id` *(required)*, `--project-id` *(required)*, `--flow-id` *(required)*, `--comment`, `--tag-ids`, `--body`, `--body-file` |
| `import` | `--org-id` *(required)*, `--project-id` *(required)*, `--overwrite`, `--flow-type` |

### callbacks

| Command | Flags |
|---|---|
| `get-scheduled` | `--org-id` *(required)*, `--callback-number`, `--assignee-agent`, `--page`, `--page-size`, `--sort-by`, `--sort-order` |
| `get-scheduled-id` | `--org-id` *(required)*, `--id` *(required)* |
| `schedule` | `--org-id` *(required)*, `--customer-name`, `--callback-number`, `--timezone`, `--schedule-date`, `--start-time`, `--end-time`, `--queue-id`, `--callback-reason`, `--source-interaction`, `--assignee-agent`, `--body`, `--body-file` |
| `update-scheduled-id` | `--org-id` *(required)*, `--id` *(required)*, `--customer-name`, `--callback-number`, `--timezone`, `--schedule-date`, `--start-time`, `--end-time`, `--queue-id`, `--callback-reason`, `--source-interaction`, `--assignee-agent`, `--body`, `--body-file` |
| `delete-scheduled-id` | `--org-id` *(required)*, `--id` *(required)* |

### ai-feature

| Command | Flags |
|---|---|
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--page`, `--page-size` |
| `get-question-mapped-autocsat-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-question-mapped-autocsat` | `--orgid` *(required)*, `--filter`, `--attributes`, `--page`, `--page-size` |
| `create-question-mapped-autocsat` | `--orgid` *(required)*, `--question-id`, `--questionnaire-id`, `--organization-id`, `--id`, `--version`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save-question-mapped-autocsat` | `--orgid` *(required)*, `--body`, `--body-file` |
| `patch-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `delete-question-mapped-autocsat-id` | `--orgid` *(required)*, `--id` *(required)* |

### agent-personal-greeting-files

| Command | Flags |
|---|---|
| `get-id-v2-api` | `--orgid` *(required)*, `--id` *(required)*, `--include-url` |
| `list` | `--orgid` *(required)*, `--filter`, `--search`, `--attributes`, `--page`, `--page-size`, `--include-agent-details` |
| `create-v2-api` | `--orgid` *(required)*, `--body`, `--body-file` |
| `delete-references-1` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id-v2-api` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `patch-id-v2-api` | `--orgid` *(required)*, `--id` *(required)*, `--attribute-tag`, `--greeting-purpose-id`, `--body`, `--body-file` |
| `delete-id-v2-api` | `--orgid` *(required)*, `--id` *(required)* |

### journey-customer-identification

| Command | Flags |
|---|---|
| `get-all-person` | `--workspace-id` *(required)*, `--person-id`, `--filter`, `--sort-by`, `--sort`, `--page`, `--page-size` |
| `search-identity-aliases` | `--workspace-id` *(required)*, `--aliases` *(required)*, `--sort-by`, `--sort`, `--page`, `--page-size` |
| `create-person` | `--workspace-id` *(required)*, `--first-name`, `--last-name`, `--phone`, `--email`, `--temporary-id`, `--customer-id`, `--body`, `--body-file` |
| `merges-identities-primary-identity` | `--workspace-id` *(required)*, `--primary-person-id` *(required)*, `--person-ids-to-merge`, `--body`, `--body-file` |
| `creates-merges-aliases-individual-jds` | `--workspace-id` *(required)*, `--override`, `--first-name`, `--last-name`, `--phone`, `--email`, `--temporary-id`, `--customer-id`, `--social-id`, `--body`, `--body-file` |
| `create-remove-replace-person` | `--workspace-id` *(required)*, `--person-id` *(required)*, `--body`, `--body-file` |
| `create-one-more-identities-person` | `--workspace-id` *(required)*, `--person-id` *(required)*, `--phone`, `--email`, `--temporary-id`, `--customer-id`, `--body`, `--body-file` |
| `delete-one-more-identities-person` | `--workspace-id` *(required)*, `--person-id` *(required)*, `--body`, `--body-file` |
| `delete-person-id` | `--workspace-id` *(required)*, `--person-id` *(required)* |

### journey-data-ingestion

| Command | Flags |
|---|---|
| `event-posting` | `--workspace-id`, `--body`, `--body-file` |

### journey-profile-creation-insights

| Command | Flags |
|---|---|
| `get-template-searched-template-id` | `--workspace-id` *(required)*, `--template-id` *(required)* |
| `get-all-template` | `--workspace-id` *(required)*, `--filter`, `--sort`, `--sort-by`, `--page`, `--page-size` |
| `get-template-searched-template-name` | `--workspace-id` *(required)*, `--template-name` *(required)* |
| `get-historic-view-template-name` | `--workspace-id` *(required)*, `--person-id` *(required)*, `--template-name` *(required)* |
| `get-historic-view` | `--workspace-id` *(required)*, `--person-id` *(required)*, `--template-id` *(required)* |
| `get-historic-view-template-name-2` | `--workspace-id` *(required)*, `--identity` *(required)*, `--template-name` *(required)* |
| `get-historic-view-template-id` | `--workspace-id` *(required)*, `--identity` *(required)*, `--template-id` *(required)* |
| `stream-views-template-name` | `--workspace-id` *(required)*, `--identity` *(required)*, `--template-name` *(required)* |
| `stream-views-template-id` | `--workspace-id` *(required)*, `--identity` *(required)*, `--template-id` *(required)* |
| `get-historic-events` | `--workspace-id` *(required)*, `--identity`, `--sort-by`, `--sort`, `--filter`, `--data`, `--page`, `--page-size` |
| `stream-events-identity` | `--workspace-id` *(required)*, `--identity` *(required)*, `--filter`, `--data` |
| `create-template` | `--workspace-id` *(required)*, `--body`, `--body-file` |
| `update-profileviewtemplate` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--body`, `--body-file` |
| `delete-template-template-id` | `--workspace-id` *(required)*, `--template-id` *(required)* |

### journey-subscription

| Command | Flags |
|---|---|
| `get-wxcc` | `--workspace-id` *(required)* |
| `create-wxcc` | `--workspace-id` *(required)* |
| `delete-wxcc` | `--workspace-id` *(required)* |

### journey-trigger-actions

| Command | Flags |
|---|---|
| `get-all` | `--workspace-id` *(required)*, `--sort-by`, `--sort`, `--page`, `--page-size` |
| `get-all-template` | `--workspace-id` *(required)*, `--template-id` *(required)* |
| `get-name` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--action-name` *(required)* |
| `get-actionid` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--action-id` *(required)* |
| `create` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--body`, `--body-file` |
| `update` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--action-id` *(required)*, `--body`, `--body-file` |
| `delete-configuration-actionid` | `--workspace-id` *(required)*, `--template-id` *(required)*, `--action-id` *(required)* |

### journey-workspace-management

| Command | Flags |
|---|---|
| `get` | `--workspace-id` *(required)* |
| `get-all` | `--filter`, `--sort-by`, `--sort`, `--page`, `--page-size` |
| `create` | `--description`, `--name`, `--body`, `--body-file` |
| `update` | `--workspace-id` *(required)*, `--description`, `--name`, `--body`, `--body-file` |
| `delete` | `--workspace-id` *(required)* |

### dnc-management

| Command | Flags |
|---|---|
| `get-phone-number-list` | `--dnc-list-name` *(required)*, `--phone-number` *(required)* |
| `create-phone-number-list` | `--dnc-list-name` *(required)*, `--phone-number`, `--source`, `--reason`, `--body`, `--body-file` |
| `delete-phone-number-list` | `--dnc-list-name` *(required)*, `--phone-number` *(required)* |

### resource-collection

| Command | Flags |
|---|---|
| `get-id` | `--orgid` *(required)*, `--id` *(required)* |
| `list-references` | `--orgid` *(required)*, `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--orgid` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-default` | `--orgid` *(required)*, `--body`, `--body-file` |
| `update-id` | `--orgid` *(required)*, `--id` *(required)*, `--body`, `--body-file` |
| `bulk-partial-update` | `--orgid` *(required)*, `--body`, `--body-file` |
| `delete-id` | `--orgid` *(required)*, `--id` *(required)* |

### contact-list-management

| Command | Flags |
|---|---|
| `get-within-campaign` | `--campaign-id` *(required)*, `--status`, `--source` |
| `create` | `--campaign-id` *(required)*, `--supported-channels`, `--activation-time-lag-minutes`, `--activation-date-time`, `--body`, `--body-file` |
| `create-within` | `--campaign-id` *(required)*, `--contact-list-id` *(required)*, `--body`, `--body-file` |
| `update-status-within` | `--campaign-id` *(required)*, `--contact-list-id` *(required)*, `--contact-id` *(required)*, `--contact-status`, `--body`, `--body-file` |
| `update-status` | `--campaign-id` *(required)*, `--contact-list-id` *(required)*, `--contact-list-status`, `--body`, `--body-file` |

### agent-summaries

| Command | Flags |
|---|---|
| `list` | `--org-id`, `--interaction-id`, `--search-type`, `--body`, `--body-file` |
| `list-2` | `--org-id`, `--agent-ci-user-id`, `--search-type`, `--body`, `--body-file` |

### ai-assistant

| Command | Flags |
|---|---|
| `get-suggestions` | `--body`, `--body-file` |

<!-- codegen:end -->
