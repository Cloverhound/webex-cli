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
| `site` | list, create, get, update, delete |
| `entry-point` | list, create, get, update, delete |
| `contact-service-queue` | list, create, get, update, delete |
| `team` | list, create, get, update, delete |
| `users` | list, get, get-along-profile-id, patch, update |
| `agents` | login, logout, state-change, get-activities, get-statistics |
| `flow` | list, export, import, publish |
| `audio-files` | list, get, create, update, delete |
| `global-variables` | list, create, get, update, delete |
| `business-hour` | list, create, get, update, delete |
| `auxiliary-code` | list, create, get, update, delete |
| `desktop-profile` | list, create, get, update, delete |
| `desktop-layout` | list, create, get, update, delete |
| `skill` | list, create, get, update, delete |
| `skill-profile` | list, create, get, update, delete |
| `multimedia-profile` | list, create, get, update, delete |
| `dial-plan` | list, create, get, update, delete |
| `dial-number` | list, list-dialed-mapping, get |
| `outdial-ani` | list, create, get, update, delete |
| `address-book` | list, create, get, update, delete |
| `holiday-list` | list, create, get, update, delete |

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
webex cc site get --id <id>
webex cc entry-point get --id <id>
webex cc contact-service-queue get --id <id>
webex cc team get --id <id>

# Create / update / delete
webex cc site create --body '{...}'
webex cc site update --id <id> --body '{...}'
webex cc site delete --id <id>

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
webex cc users get --id <id>
webex cc users get-along-profile-id --id <id>
webex cc users get-ci-id --ci-id <ciUserId>

# Update user (agent settings, team, profile)
webex cc users patch --id <id> --body '{...}'
webex cc users update --id <id> --body '{...}'

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

## Flow authoring and migration

Current flow commands use the unprefixed Contact Center routes. `flow` is the canonical group; `flows` is an alias. Import/export now use the typed FlowV2 JSON contract, not legacy FDL. Existing automation consuming raw FDL must use `export-legacy` instead.

```bash
webex cc flow list
webex cc flow get --flow-id <id>
webex cc flow export --flow-id <id> --version draft > flow.json
webex cc flow validate --body-file flow.json
webex cc flow import --body-file flow.json --overwrite false --flow-type FLOW
webex cc flow patch-draft --flow-id <id> --body-file patch.json
webex cc flow export-legacy --flow-id <id> > legacy-flow.json
webex cc activities list-definitions
webex cc events list-specifications
webex cc templates list-flow
```

Bodies must be supplied with `--body` or `--body-file`; these generated commands do not read stdin. The CLI uses the fixed project ID internally. `--organization` overrides the login organization.

Cloverhound checks verified reads and an unpublished disposable draft's import, lock/unlock, patch, save, export, and deletion across both URL forms. The draft was deleted and both routes returned 404 afterward. Publish was not live-tested. Templates use `/templates`; the scoped prefixed template list returned 404.

**Import limitations:** `flow import-legacy` and `functions import` are exposed but their generated commands do not implement the required upload body. Do not use them for migration until upload handling is implemented. Current `flow import --body-file` was live-tested successfully.

## Other refreshed Contact Center APIs

| Group | Purpose |
| --- | --- |
| `activities`, `events`, `templates` | Discover activity schemas, event schemas, input choices, and starter flows. |
| `functions` | Manage custom function source, drafts, locks, versions, and publish tags. `test` implicitly publishes the draft before execution. |
| `asset`, `channel` | Configure work-item/custom-messaging resources and their references. |
| `campaign-group` | List campaigns belonging to a campaign group. |
| `search-metadata` | Discover supported Search GraphQL fields and operations. |
| `usage-reports` | Discover resource types, request asynchronous CSV exports, inspect status, and download files. |
| `external-data-updates` | Update numeric global-variable values on tasks completed within the preceding two days. |

```bash
webex cc usage-reports list-resource-types
webex cc usage-reports create --body-file report-request.json
webex cc usage-reports get --report-id <report-id>
webex cc usage-reports download-file --file-id <file-id> --output raw > report.bin
webex cc external-data-updates update-task-global-variables --body-file task-update.json
```

Discover report types before creating a report; the collection does not enumerate the types or CSV columns. Download with `fileId`, not `reportId`. Reports can span available data up to 36 months, may split into monthly files, and expire seven days after completion. Only one report processes per organization at a time.

For task controls, `unhold` (alias `resume`) resumes a held voice call; `pause-digital` and `resume-digital` handle paused non-real-time digital tasks. `append-message` retains `update-2` as an alias. Queue lookups use `list-agent-based`, `list-skill-based`, `list-team-based`, `list-by-skill-profile`, `list-by-dynamic-skills`, and `list-by-user-skill-profile`; pass IDs as flags. Whisper coaching is beta. The refresh also adds conference-participant removal, campaign-time lookup, and contact closure across a campaign chain.

## Audio Files

```bash
# List
webex cc audio-files list [--filter <rsql>] [--search <text>]

# Get with download URL
webex cc audio-files get --id <id> --include-url true

# Create / upload
webex cc audio-files create --body '{"name":"Main Greeting","description":"..."}'
# After creating, upload the WAV via update
webex cc audio-files update --id <id> --body '{...}'

# Delete
webex cc audio-files delete --id <id>
```

## Configuration Resources

All of these (global-variables, business-hour, auxiliary-code, desktop-profile, desktop-layout, skill, skill-profile, multimedia-profile, dial-plan, outdial-ani, address-book, holiday-list) follow the same pattern:

```bash
webex cc <resource> list [--filter <rsql>] [--search <text>] [--page <n>] [--page-size <n>]
webex cc <resource> get --id <id>
webex cc <resource> create --body '{...}'
webex cc <resource> update --id <id> --body '{...}'
webex cc <resource> delete --id <id>
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
webex cc dial-number get --id <id>
```

### Outdial ANI (has entries sub-resource)

```bash
webex cc outdial-ani list
webex cc outdial-ani get --id <id>
webex cc outdial-ani list-entry --id <id>          # list ANI entries under an outdial ANI
webex cc outdial-ani create-entry --body '{...}'
webex cc outdial-ani delete-entry-id --id <entryId>
```

### Address Book (has entries sub-resource)

```bash
webex cc address-book list
webex cc address-book get --id <id>
webex cc address-book list-entry --id <id>         # list contacts
webex cc address-book create-entry --body '{...}'
webex cc address-book bulk-save-entry --body '[{...}]'
webex cc address-book update-entry-id --id <entryId> --body '{...}'
webex cc address-book delete-entry-id --id <entryId>
```

## Key Gotchas

1. **`--orgid=VALUE` syntax** — always use `=` (equals) with `--orgid`, not a space: `--orgid="<uuid>"` not `--orgid "<uuid>"`. Shell quoting issues cause silent failures with the space form.

2. **RSQL values with spaces** — wrap values in double quotes inside the RSQL expression: `--filter='name=="Site A"'`.

3. **`get` not `get`** — most CC resources use `get --id <id>` (not `get --<resource>-id`). The `agents` resource is an exception: it has no `list` or `get` commands, only session control operations.

4. **`flow` has no CRUD** — `flow` only has `list`, `export`, `import`, `publish`. To create or edit flows, use Webex Contact Center Flow Designer in Control Hub.

5. **Agent statistics time window** — `agents get-activities` and `get-statistics` enforce a **24-hour** max window per request. Use multiple calls to cover longer periods.

6. **Purge vs delete** — `delete` soft-deletes. `purge-inactive` permanently removes soft-deleted records for resources that support it.

<!-- codegen:start -->
## Command Reference

> Auto-generated from Postman collections. Run `make codegen` to update.
> Organization defaults to the authenticated account; use `--organization` to override.
> The flow project ID is fixed internally; no project flag is needed.

### agent-wellbeing

| Command | Flags |
|---|---|
| `get-burnout-id` | `--id` *(required)* |
| `list-burnout` | `--filter`, `--attributes`, `--page`, `--page-size` |
| `subscribe-realtime-burnout-events` | `--destination-url`, `--event-types`, `--name`, `--description`, `--secret`, `--org-id`, `--body`, `--body-file` |
| `record-realtime-burnout-events` | `--body`, `--body-file` |
| `update-burnout-id` | `--id` *(required)*, `--agent-inclusion-type`, `--enabled`, `--organization-id`, `--version`, `--wellness-break-reminders`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |

### auto-csat

| Command | Flags |
|---|---|
| `get-mapped-question-id` | `--auto-csat-id` *(required)*, `--id` *(required)* |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list` | `--filter`, `--attributes`, `--page`, `--page-size` |
| `list-mapped-question` | `--auto-csat-id` *(required)*, `--filter`, `--attributes`, `--page`, `--page-size` |
| `create-mapped-question` | `--auto-csat-id` *(required)*, `--question-id`, `--questionnaire-id`, `--organization-id`, `--id`, `--version`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save-mapped-question` | `--auto-csat-id` *(required)*, `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--agent-inclusion-type`, `--enabled`, `--selected-global-variable-id`, `--survey-data-source`, `--organization-id`, `--version`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete-mapped-question-id` | `--auto-csat-id` *(required)*, `--id` *(required)* |

### generated-summaries

| Command | Flags |
|---|---|
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list` | `--filter`, `--attributes`, `--page`, `--page-size` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--organization-id`, `--version`, `--call-drop-summaries-enabled`, `--virtual-agent-transfer-summaries-enabled`, `--consult-transfer-summaries-enabled`, `--agent-inclusion-type`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |

### business-hour

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--sort`, `--include-count`, `--single-object-response` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### holiday-list

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--sort`, `--include-count`, `--single-object-response` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### overrides

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--sort`, `--include-latest-override`, `--single-object-response` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### address-book

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `bulk-export` | `--page`, `--page-size` |
| `get-entry-id` | `--address-book-id` *(required)*, `--id` *(required)* |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list-entry` | `--address-book-id` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--name`, `--parent-type`, `--organization-id`, `--id`, `--version`, `--description`, `--site-id`, `--body`, `--body-file` |
| `create-entry` | `--address-book-id` *(required)*, `--name`, `--number`, `--organization-id`, `--id`, `--version`, `--body`, `--body-file` |
| `bulk-save-entry` | `--address-book-id` *(required)*, `--body`, `--body-file` |
| `update-entry-id` | `--address-book-id` *(required)*, `--id` *(required)*, `--name`, `--number`, `--organization-id`, `--version`, `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--name`, `--parent-type`, `--organization-id`, `--version`, `--description`, `--site-id`, `--body`, `--body-file` |
| `delete-entry-id` | `--address-book-id` *(required)*, `--id` *(required)* |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### audio-files

| Command | Flags |
|---|---|
| `get` (aliases: `get-id`) | `--id` *(required)*, `--include-url` |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `patch` (aliases: `patch-id`) | `--id` *(required)*, `--description`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### auxiliary-code

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--desktop-profile-filter`, `--supervised-user-id` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--active`, `--default-code`, `--name`, `--work-type-code`, `--work-type-id`, `--organization-id`, `--id`, `--version`, `--description`, `--is-system-code`, `--burnout-inclusion`, `--system-default`, `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--active`, `--default-code`, `--name`, `--work-type-code`, `--work-type-id`, `--organization-id`, `--version`, `--description`, `--is-system-code`, `--burnout-inclusion`, `--system-default`, `--body`, `--body-file` |
| `bulk-partial-update` | `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### contact-number

| Command | Flags |
|---|---|
| `list-all` | — |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--number`, `--organization-id`, `--id`, `--version`, `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--number`, `--organization-id`, `--version`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### contact-service-queue

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--desktop-profile-filter`, `--provisioning-view`, `--single-object-response` |
| `bulk-export` | `--type`, `--page`, `--page-size` |
| `list-by-skill-profile` (aliases: `list-skill-csqs-skill-profile`) | `--id` *(required)* |
| `get` (aliases: `get-id`) | `--id` *(required)*, `--agents-updated-info` |
| `list-references` (aliases: `list-csq-references-id`) | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list-agent-based` | `--userid` *(required)*, `--search`, `--page`, `--page-size` |
| `list-skill-based` | `--userid` *(required)*, `--search`, `--page`, `--page-size` |
| `list-team-based` | `--userid` *(required)*, `--search`, `--page`, `--page-size` |
| `list-internal-skill-csqs-profile` | `--id` *(required)* |
| `list-team-csqs-team-id` | `--id` *(required)* |
| `list-agent-csqs-ci-user-id` | `--ci-user-id` *(required)* |
| `list-skill-csqs-ci-user-id` | `--id` *(required)* |
| `create` | `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `delete-references` (aliases: `delete-csq-references`) | `--body`, `--body-file` |
| `list-manually-assignable` (aliases: `list-manually-assignable-csqs`) | `--agent-id`, `--team-id`, `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `create-remove-agents-users-agent` | `--id` *(required)*, `--add`, `--remove`, `--body`, `--body-file` |
| `list-mapping-summary-grouped-assistant-skill` | `--page`, `--page-size`, `--assistant-skill-ids`, `--body`, `--body-file` |
| `list-by-dynamic-skills` (aliases: `list-csqs-skills-profile`) | `--body`, `--body-file` |
| `list-by-user-skill-profile` (aliases: `list-csqs-user-profile`) | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `bulk-partial-update` | `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### desktop-layout

| Command | Flags |
|---|---|
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--single-object-response`, `--provisioning-view` |
| `create` | `--default-json-modified`, `--edited-by`, `--global`, `--json-file-content`, `--json-file-name`, `--name`, `--status`, `--validated`, `--organization-id`, `--id`, `--version`, `--description`, `--validated-time`, `--default-json-modified-time`, `--modified-time`, `--team-ids`, `--system-default`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--default-json-modified`, `--edited-by`, `--global`, `--json-file-content`, `--json-file-name`, `--name`, `--status`, `--validated`, `--organization-id`, `--version`, `--description`, `--validated-time`, `--default-json-modified-time`, `--modified-time`, `--team-ids`, `--system-default`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### desktop-profile

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--single-object-response`, `--provisioning-view` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--access-buddy-team`, `--access-entry-point`, `--access-idle-code`, `--access-queue`, `--access-wrap-up-code`, `--active`, `--agent-dnvalidation`, `--name`, `--parent-type`, `--viewable-statistics`, `--organization-id`, `--id`, `--version`, `--description`, `--site-id`, `--screen-popup`, `--last-agent-routing`, `--schedule-and-manage-call-back`, `--auto-wrap-up`, `--auto-answer`, `--agent-personal-greeting`, `--auto-wrap-after-seconds`, `--agent-available-after-outdial`, `--allow-auto-wrap-up-extension`, `--wrap-up-codes`, `--idle-codes`, `--queues`, `--entry-points`, `--buddy-teams`, `--consult-to-queue`, `--outdial-enabled`, `--outdial-entry-point-id`, `--outdial-aniid`, `--address-book-id`, `--dial-plan-enabled`, `--dial-plans`, `--agent-dnvalidation-criteria`, `--agent-dnvalidation-criterions`, `--login-voice-options`, `--threshold-rules`, `--timeout-desktop-inactivity-custom-enabled`, `--show-user-details-ms`, `--state-synchronization-ms`, `--show-user-details-webex`, `--state-synchronization-webex`, `--manage-channel-availability`, `--timeout-desktop-inactivity-mins`, `--system-default`, `--created-time`, `--last-updated-time`, `--auto-accept-digital-interactions`, `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--access-buddy-team`, `--access-entry-point`, `--access-idle-code`, `--access-queue`, `--access-wrap-up-code`, `--active`, `--agent-dnvalidation`, `--name`, `--parent-type`, `--viewable-statistics`, `--organization-id`, `--version`, `--description`, `--site-id`, `--screen-popup`, `--last-agent-routing`, `--schedule-and-manage-call-back`, `--auto-wrap-up`, `--auto-answer`, `--agent-personal-greeting`, `--auto-wrap-after-seconds`, `--agent-available-after-outdial`, `--allow-auto-wrap-up-extension`, `--wrap-up-codes`, `--idle-codes`, `--queues`, `--entry-points`, `--buddy-teams`, `--consult-to-queue`, `--outdial-enabled`, `--outdial-entry-point-id`, `--outdial-aniid`, `--address-book-id`, `--dial-plan-enabled`, `--dial-plans`, `--agent-dnvalidation-criteria`, `--agent-dnvalidation-criterions`, `--login-voice-options`, `--threshold-rules`, `--timeout-desktop-inactivity-custom-enabled`, `--show-user-details-ms`, `--state-synchronization-ms`, `--show-user-details-webex`, `--state-synchronization-webex`, `--manage-channel-availability`, `--timeout-desktop-inactivity-mins`, `--system-default`, `--created-time`, `--last-updated-time`, `--auto-accept-digital-interactions`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### dial-number

| Command | Flags |
|---|---|
| `list-dialed-mapping` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--include-entry-point-name` |
| `bulk-export-dialed-mapping` | `--page`, `--page-size` |
| `list-dialed-dialed-mapping` | — |
| `get-dialed-mapping-id` | `--id` *(required)* |
| `list-references-dialed-mapping` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create-dialed-mapping` | `--entry-point-id`, `--entry-point-name`, `--organization-id`, `--id`, `--version`, `--dialled-number`, `--extension`, `--routing-prefix`, `--esn`, `--route-point-id`, `--default-ani`, `--location`, `--region-id`, `--created-time`, `--last-updated-time`, `--dialled-number-digits`, `--body`, `--body-file` |
| `bulk-save-dialed-mapping` | `--body`, `--body-file` |
| `update-dialed-mapping-id` | `--id` *(required)*, `--entry-point-id`, `--entry-point-name`, `--organization-id`, `--version`, `--dialled-number`, `--extension`, `--routing-prefix`, `--esn`, `--route-point-id`, `--default-ani`, `--location`, `--region-id`, `--created-time`, `--last-updated-time`, `--dialled-number-digits`, `--body`, `--body-file` |
| `delete-all-dialed-mapping` | — |
| `delete-dialed-mapping-id` | `--id` *(required)* |

### dial-plan

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--active`, `--name`, `--regular-expression`, `--organization-id`, `--id`, `--version`, `--description`, `--prefix`, `--stripped-chars`, `--system-default`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--active`, `--name`, `--regular-expression`, `--organization-id`, `--version`, `--description`, `--prefix`, `--stripped-chars`, `--system-default`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### entry-point

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--desktop-profile-filter`, `--provisioning-view`, `--include-count`, `--single-object-response` |
| `bulk-export` | `--type`, `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)*, `--include-names` |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### global-variables

| Command | Flags |
|---|---|
| `bulk-export` | `--page`, `--page-size` |
| `get-reportable-count` | — |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--active`, `--agent-editable`, `--agent-viewable`, `--default-value`, `--name`, `--reportable`, `--variable-type`, `--organization-id`, `--id`, `--version`, `--description`, `--sensitive`, `--desktop-label`, `--system-default`, `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--active`, `--agent-editable`, `--agent-viewable`, `--default-value`, `--name`, `--reportable`, `--variable-type`, `--organization-id`, `--version`, `--description`, `--sensitive`, `--desktop-label`, `--system-default`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### multimedia-profile

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### outdial-ani

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--single-object-response` |
| `bulk-export` | `--page`, `--page-size` |
| `list-entries` (aliases: `list-entry`) | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `get-entry-id` | `--out-dial-ani-id` *(required)*, `--id` *(required)* |
| `list-entries-for-ani` (aliases: `list-entry-2`) | `--out-dial-ani-id` *(required)*, `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `create-entry` | `--out-dial-ani-id` *(required)*, `--name`, `--number`, `--organization-id`, `--id`, `--version`, `--default-anientry`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save-entries` (aliases: `bulk-save-entry`) | `--out-dial-ani-id` *(required)*, `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `update-entry-id` | `--out-dial-ani-id` *(required)*, `--id` *(required)*, `--name`, `--number`, `--organization-id`, `--version`, `--default-anientry`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |
| `delete-entry-id` | `--out-dial-ani-id` *(required)*, `--id` *(required)* |

### site

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--active`, `--multimedia-profile-id`, `--name`, `--organization-id`, `--id`, `--version`, `--description`, `--system-default`, `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--active`, `--multimedia-profile-id`, `--name`, `--organization-id`, `--version`, `--description`, `--system-default`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### skill

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--single-object-response`, `--include` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `populate-json-attributes-field-skill-id-org` | `--id` *(required)* |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### skill-profile

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--single-object-response` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)*, `--include-skill-details` |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--active-skills`, `--name`, `--organization-id`, `--id`, `--version`, `--description`, `--active-enum-skills`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--active-skills`, `--name`, `--organization-id`, `--version`, `--description`, `--active-enum-skills`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### team

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--supervisor-view`, `--provisioning-view`, `--single-object-response` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `create` | `--active`, `--name`, `--rank-queues-for-team`, `--site-id`, `--team-status`, `--team-type`, `--organization-id`, `--id`, `--version`, `--dialed-number`, `--capacity`, `--desktop-layout-id`, `--skill-profile-id`, `--multi-media-profile-id`, `--user-ids`, `--description`, `--system-default`, `--queue-rankings`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--active`, `--name`, `--rank-queues-for-team`, `--site-id`, `--team-status`, `--team-type`, `--organization-id`, `--version`, `--dialed-number`, `--capacity`, `--desktop-layout-id`, `--skill-profile-id`, `--multi-media-profile-id`, `--user-ids`, `--description`, `--system-default`, `--queue-rankings`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### user-profiles

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)*, `--include-names` |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `get-acl-id` | `--id` *(required)*, `--names` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `create` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### users

| Command | Flags |
|---|---|
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--supervisor-managed-agents-only`, `--single-object-response`, `--buddy-team-agents-only`, `--user-in-queue`, `--queue-id`, `--include-aimapping-count`, `--include-dynamic-skills-limit-reached` |
| `bulk-export` | `--page`, `--page-size` |
| `get-ci-id` | `--id` *(required)*, `--include-user-profile`, `--include-names`, `--include-skill-details` |
| `list-along-profile` | — |
| `get-along-profile-id` | `--id` *(required)* |
| `get` (aliases: `get-id`) | `--id` *(required)*, `--include-count`, `--include-user-profile-type`, `--include-skill-profile-audit`, `--include-reskill-audit-info`, `--include-skill-details`, `--check-if-user-has-dynamic-skill`, `--dynamic-skill-id` |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `get-dynamic-skill-id` | `--skill-id` *(required)*, `--search`, `--page`, `--page-size` |
| `list-call-monitoring-id` | `--id` *(required)*, `--page`, `--page-size` |
| `get-agents-matching-skill-requirements` | `--search`, `--page`, `--page-size`, `--body`, `--body-file` |
| `list-by-ids` (aliases: `get-ids`) | `--page`, `--page-size`, `--user-ids`, `--search`, `--queue-id`, `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `bulk-partial-update` | `--body`, `--body-file` |
| `patch` (aliases: `patch-id`) | `--id` *(required)*, `--value-type`, `--body`, `--body-file` |
| `bulk-partial-update-dynamic-skills` (aliases: `bulk-update-dynamic-skills`) | `--skill-id` *(required)*, `--body`, `--body-file` |
| `reskill-agents` | `--id` *(required)*, `--body`, `--body-file` |

### work-types

| Command | Flags |
|---|---|
| `bulk-export` | `--page`, `--page-size` |
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--active`, `--name`, `--work-type-code`, `--organization-id`, `--id`, `--version`, `--description`, `--system-default`, `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `purge-inactive` | `--next-start-id` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--active`, `--name`, `--work-type-code`, `--organization-id`, `--version`, `--description`, `--system-default`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

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
| `whisper-coach-request` | `--interaction-id` *(required)* |
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
| `unhold` (aliases: `resume`) | `--task-id` *(required)*, `--media-resource-id`, `--body`, `--body-file` |
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
| `append-message` (aliases: `update-2`) | `--task-id` *(required)*, `--body`, `--body-file` |
| `drop-participant-conference` | `--task-id` *(required)*, `--participant-id` *(required)* |
| `pause-digital` | `--task-id` *(required)* |
| `resume-digital` | `--task-id` *(required)* |
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
| `get-valid-times` | `--campaign-id`, `--interaction-id`, `--agent-id`, `--tracking-id` |
| `start-request` | `--body`, `--body-file` |
| `update-request` | `--campaign-id` *(required)*, `--body`, `--body-file` |
| `stop-request` | `--campaign-id` *(required)* |

### captures

| Command | Flags |
|---|---|
| `list` | `--body`, `--body-file` |

### flow

Group aliases: `flows`.

| Command | Flags |
|---|---|
| `export-legacy` | `--flow-id` *(required)*, `--version` |
| `list` | `--flow-type`, `--ids`, `--page`, `--partial-name-search`, `--search-by`, `--size`, `--include-pagination`, `--is-validation` |
| `get` | `--flow-id` *(required)*, `--flow-type` |
| `validate-draft` | `--flow-id` *(required)*, `--version-id`, `--flow-type` |
| `export` | `--flow-id` *(required)*, `--version`, `--flow-type` |
| `search` | `--query`, `--flow-type`, `--page`, `--size`, `--key-value-filter` |
| `import-legacy` | `--overwrite`, `--flow-type` |
| `publish` | `--flow-id` *(required)*, `--skip-validation`, `--flow-type`, `--comment`, `--tag-ids`, `--body`, `--body-file` |
| `lock` | `--flow-id` *(required)*, `--flow-type` |
| `unlock` | `--flow-id` *(required)*, `--flow-type` |
| `validate` | `--body`, `--body-file` |
| `import` | `--overwrite`, `--flow-type`, `--body`, `--body-file` |
| `save-draft` | `--flow-id` *(required)*, `--expected-version`, `--flow-type`, `--body`, `--body-file` |
| `patch-draft` | `--flow-id` *(required)*, `--expected-version`, `--flow-type`, `--body`, `--body-file` |
| `delete` | `--flow-id` *(required)*, `--force`, `--skip-rs-epcheck`, `--flow-type` |

### callbacks

| Command | Flags |
|---|---|
| `get-scheduled` | `--callback-number`, `--assignee-agent`, `--page`, `--page-size`, `--sort-by`, `--sort-order` |
| `get-scheduled-id` | `--id` *(required)* |
| `schedule` | `--customer-name`, `--callback-number`, `--timezone`, `--schedule-date`, `--start-time`, `--end-time`, `--queue-id`, `--callback-reason`, `--source-interaction`, `--assignee-agent`, `--body`, `--body-file` |
| `update-scheduled-id` | `--id` *(required)*, `--customer-name`, `--callback-number`, `--timezone`, `--schedule-date`, `--start-time`, `--end-time`, `--queue-id`, `--callback-reason`, `--source-interaction`, `--assignee-agent`, `--body`, `--body-file` |
| `delete-scheduled-id` | `--id` *(required)* |

### ai-feature

| Command | Flags |
|---|---|
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list` | `--filter`, `--attributes`, `--page`, `--page-size` |
| `get-question-mapped-autocsat-id` | `--id` *(required)* |
| `list-question-mapped-autocsat` | `--filter`, `--attributes`, `--page`, `--page-size` |
| `create-question-mapped-autocsat` | `--question-id`, `--questionnaire-id`, `--organization-id`, `--id`, `--version`, `--created-time`, `--last-updated-time`, `--body`, `--body-file` |
| `bulk-save-question-mapped-autocsat` | `--body`, `--body-file` |
| `patch` (aliases: `patch-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `delete-question-mapped-autocsat-id` | `--id` *(required)* |

### agent-personal-greeting-files

| Command | Flags |
|---|---|
| `get` (aliases: `get-id`, `get-id-v2-api`) | `--id` *(required)*, `--include-url` |
| `list` | `--filter`, `--search`, `--attributes`, `--page`, `--page-size`, `--include-agent-details` |
| `create` (aliases: `create-v2-api`) | `--body`, `--body-file` |
| `delete-references` (aliases: `delete-references-1`) | `--body`, `--body-file` |
| `update` (aliases: `update-id`, `update-id-v2-api`) | `--id` *(required)*, `--body`, `--body-file` |
| `patch` (aliases: `patch-id`, `patch-id-v2-api`) | `--id` *(required)*, `--attribute-tag`, `--greeting-purpose-id`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`, `delete-id-v2-api`) | `--id` *(required)* |

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
| `get` (aliases: `get-id`) | `--id` *(required)* |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size` |
| `create` | `--body`, `--body-file` |
| `update-default` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `bulk-partial-update` | `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### contact-list-management

| Command | Flags |
|---|---|
| `get-within-campaign` | `--campaign-id` *(required)*, `--status`, `--source` |
| `create` | `--campaign-id` *(required)*, `--supported-channels`, `--activation-time-lag-minutes`, `--activation-date-time`, `--body`, `--body-file` |
| `create-within` | `--campaign-id` *(required)*, `--contact-list-id` *(required)*, `--body`, `--body-file` |
| `update-status-within` | `--campaign-id` *(required)*, `--contact-list-id` *(required)*, `--contact-id` *(required)*, `--contact-status`, `--body`, `--body-file` |
| `update-status` | `--campaign-id` *(required)*, `--contact-list-id` *(required)*, `--contact-list-status`, `--body`, `--body-file` |
| `update-status-across-campaign-chain` | `--campaign-id` *(required)*, `--contact-id` *(required)*, `--contact-list-id`, `--fields`, `--contact-status`, `--search-across-the-campaign-chain`, `--body`, `--body-file` |

### agent-summaries

| Command | Flags |
|---|---|
| `list` | `--org-id`, `--interaction-id`, `--search-type`, `--body`, `--body-file` |
| `list-2` | `--org-id`, `--search-type`, `--body`, `--body-file` |

### ai-assistant

| Command | Flags |
|---|---|
| `get-suggestions` | `--body`, `--body-file` |

### activities

| Command | Flags |
|---|---|
| `list-definitions` | — |
| `describe` | `--activity-name` *(required)* |
| `get-input-choices` | `--activity-name` *(required)*, `--input-name` *(required)*, `--search`, `--validate`, `--parent-value`, `--parent-input-name` |

### events

| Command | Flags |
|---|---|
| `list-specifications` | — |

### functions

| Command | Flags |
|---|---|
| `list` | `--is-partial-match`, `--is-case-sensitive`, `--name`, `--language`, `--status`, `--sort-by`, `--page`, `--size`, `--ids`, `--fields`, `--is-validation` |
| `get` | `--id` *(required)*, `--version-or-tag`, `--meta-data-only` |
| `list-options` | `--language`, `--runtime` |
| `create` | `--body`, `--body-file` |
| `import` | `--overwrite`, `--associated-rcs` |
| `unlock` | `--id` *(required)* |
| `lock` | `--id` *(required)* |
| `test` | `--id` *(required)*, `--body`, `--body-file` |
| `publish` | `--id` *(required)*, `--tags`, `--comment`, `--key_0`, `--key_1`, `--key_2`, `--key_3`, `--body`, `--body-file` |
| `export` | `--id` *(required)*, `--version-or-tag` |
| `update` | `--id` *(required)*, `--body`, `--body-file` |
| `delete` | `--id` *(required)*, `--is-force-deletion` |

### templates

| Command | Flags |
|---|---|
| `list-flow` | `--type` |
| `get-flow` | `--id` *(required)* |

### campaign-group

| Command | Flags |
|---|---|
| `list` | `--campaign-group-name` *(required)*, `--page`, `--page-size`, `--campaign-status` |

### search-metadata

| Command | Flags |
|---|---|
| `get` | — |

### asset

| Command | Flags |
|---|---|
| `get` (aliases: `get-id`) | `--id` *(required)*, `--include-channel-name` |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--sort`, `--include-count`, `--single-object-response`, `--include-channel-name`, `--exclude-epassociated` |
| `create` | `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)*, `--body`, `--body-file` |
| `patch` (aliases: `patch-id`) | `--id` *(required)*, `--value-type`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### channel

| Command | Flags |
|---|---|
| `get` (aliases: `get-id`) | `--id` *(required)*, `--include-logo-url-versioned` |
| `list-references` | `--id` *(required)*, `--type`, `--page`, `--page-size` |
| `list` | `--filter`, `--attributes`, `--search`, `--page`, `--page-size`, `--sort`, `--include-count`, `--single-object-response`, `--include-logo-url-versioned` |
| `create` | `--body`, `--body-file` |
| `bulk-save` | `--body`, `--body-file` |
| `update` (aliases: `update-id`) | `--id` *(required)* |
| `patch` (aliases: `patch-id`) | `--id` *(required)*, `--description`, `--logo-type`, `--logo-icon-name`, `--body`, `--body-file` |
| `delete` (aliases: `delete-id`) | `--id` *(required)* |

### usage-reports

| Command | Flags |
|---|---|
| `list-resource-types` | — |
| `list` | `--resource-type`, `--created-after`, `--created-before` |
| `get` | `--report-id` *(required)* |
| `download-file` | `--file-id` *(required)* |
| `create` | `--resource-type`, `--start-date`, `--end-date`, `--body`, `--body-file` |
| `delete` | `--report-id` *(required)* |

### external-data-updates

| Command | Flags |
|---|---|
| `update-task-global-variables` | `--body`, `--body-file` |

<!-- codegen:end -->
