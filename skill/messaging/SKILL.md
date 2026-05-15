---
name: webex-cli/messaging
description: "Webex Messaging commands: rooms, messages, teams, memberships, webhooks, and related resources."
---

# Webex Messaging

Commands: `webex messaging <resource> <action> [flags]`  
Alias: `webex msg <resource> <action> [flags]`

## Resources

| Resource | Operations |
|---|---|
| `messages` | list, list-direct, get, create, edit, delete |
| `rooms` | list, get, get-meeting, create, update, delete |
| `memberships` | list, get, create, update, delete |
| `teams` | list, get, create, update, delete |
| `team-memberships` | list, get, create, update, delete |
| `webhooks` | list, get, create, update, delete |
| `people` | list, get-person, get-my-own, create-person, update-person, delete-person |
| `attachment-actions` | get, create |
| `room-tabs` | list, get, create, update, delete |
| `events` | list |
| `ecm-folder-linking` | list, get, create, delete |
| `hds` | (status/settings commands) |

## Messages

### List messages in a space
```bash
webex messaging messages list --room-id <roomId> [flags]
```
Flags:
- `--room-id` — target space ID (required for `list`)
- `--max` — max returned (≤100 when combined with `--mentioned-people`)
- `--before` — ISO 8601 timestamp; list messages sent before this time
- `--before-message` — list messages sent before this message ID
- `--mentioned-people` — filter by person ID; use `me` for current user. **Bots must include this flag to list messages in group spaces.**
- `--parent-id` — list threaded replies under a parent message
- `--paginate` — auto-fetch all pages

### List direct (1:1) messages
```bash
webex messaging messages list-direct --person-email <email>
webex messaging messages list-direct --person-id <personId>
webex messaging messages list-direct --person-email <email> --parent-id <msgId>
```
Use `list-direct` for 1:1 DM threads — no `--room-id` needed.

### Send a message (create)
`create` takes `--body` (raw JSON) or `--body-file <path>`. No discrete content flags.

```bash
# To a space
webex messaging messages create --body '{"roomId":"<roomId>","text":"Hello!"}'

# 1:1 by person ID or email (creates/reuses the direct room)
webex messaging messages create --body '{"toPersonId":"<personId>","text":"Hello!"}'
webex messaging messages create --body '{"toPersonEmail":"user@example.com","markdown":"**Hello!**"}'

# Threaded reply
webex messaging messages create --body '{"roomId":"<roomId>","parentId":"<parentMsgId>","text":"Reply"}'

# File attachment (URL)
webex messaging messages create --body '{"roomId":"<roomId>","text":"See attached","files":["https://example.com/file.pdf"]}'
```

**Message body fields:**

| Field | Notes |
|---|---|
| `roomId` | Target space — mutually exclusive with `toPersonId`/`toPersonEmail` |
| `toPersonId` | Send 1:1 by person ID |
| `toPersonEmail` | Send 1:1 by email |
| `text` | Plain text |
| `markdown` | Markdown (renders in clients; takes precedence over text) |
| `html` | HTML subset |
| `files` | Array of file URLs (max 1 currently) |
| `attachments` | Array of Adaptive Card objects |
| `parentId` | Thread reply — parent message ID |

### Edit and delete
```bash
webex messaging messages edit --message-id <id> --body '{"text":"Updated text"}'
webex messaging messages delete --message-id <id>
webex messaging messages get --message-id <id>
```

## Rooms (Spaces)

### List rooms
```bash
webex messaging rooms list [flags]
```
Flags:
- `--max` — max results (1–1000)
- `--type` — `direct` (1:1) or `group`
- `--last` — shorthand duration, e.g. `24h`, `7d` — sets `--from` automatically
- `--from` / `--to` — filter by `madePublic` timestamp
- `--team-id` — rooms belonging to a specific team
- `--sort-by` — sort field; cannot combine with `--org-public-spaces`
- `--org-public-spaces` — list org's public spaces (joined and unjoined)

```bash
webex messaging rooms list --type group --last 48h
webex messaging rooms list --type direct
webex messaging rooms list --team-id <teamId>
webex messaging rooms list --paginate
```

### Create, update, delete
```bash
webex messaging rooms create --title "Project Alpha"
webex messaging rooms create --body '{"title":"Project Alpha","teamId":"<tid>","isLocked":true}'
# Body fields: title, teamId, classificationId, description, isLocked, isAnnouncementOnly, isPublic

webex messaging rooms update --room-id <id> --body '{"title":"New Name"}'
webex messaging rooms delete --room-id <id>
webex messaging rooms get-meeting --room-id <id>   # SIP/PSTN join info
```

## Memberships

```bash
# List members of a space
webex messaging memberships list --room-id <roomId>
webex messaging memberships list --room-id <roomId> --person-email user@example.com

# Compliance officers: query all rooms for a person (no --room-id needed)
webex messaging memberships list --person-id <personId>
webex messaging memberships list --person-email user@example.com

# Add member
webex messaging memberships create --body '{"roomId":"<roomId>","personEmail":"user@example.com"}'
webex messaging memberships create --body '{"roomId":"<roomId>","personId":"<id>","isModerator":true}'

# Update / remove
webex messaging memberships update --membership-id <id> --body '{"isModerator":true}'
webex messaging memberships delete --membership-id <id>
```

## Teams and Team Memberships

```bash
webex messaging teams list [--max <n>]
webex messaging teams create --name "Engineering"
webex messaging teams create --body '{"name":"Engineering","description":"..."}'
webex messaging teams update --team-id <id> --body '{"name":"New Name"}'
webex messaging teams delete --team-id <id>

webex messaging team-memberships list --team-id <teamId>
webex messaging team-memberships create --body '{"teamId":"<id>","personEmail":"user@example.com"}'
webex messaging team-memberships create --body '{"teamId":"<id>","personId":"<id>","isModerator":true}'
webex messaging team-memberships update --team-membership-id <id> --body '{"isModerator":true}'
webex messaging team-memberships delete --team-membership-id <id>
```

## Webhooks

### Create
```bash
webex messaging webhooks create \
  --name "New Messages" \
  --target-url "https://myapp.example.com/hook" \
  --resource messages \
  --event created \
  --filter "roomId=<roomId>"

# With secret for HMAC signature validation
webex messaging webhooks create --body '{
  "name":"New Messages","targetUrl":"https://myapp.example.com/hook",
  "resource":"messages","event":"created","filter":"roomId=<id>","secret":"<s>"
}'

# Org-wide webhook
webex messaging webhooks create --owned-by org --name "Org" \
  --target-url "https://..." --resource messages --event created
```

### Resources and events

| Resource | Events |
|---|---|
| `messages` | `created`, `updated`, `deleted`, `seen` |
| `memberships` | `created`, `updated`, `deleted` |
| `rooms` | `created`, `updated`, `deleted` |
| `teams` | `created`, `updated`, `deleted` |
| `teamMemberships` | `created`, `updated`, `deleted` |
| `attachmentActions` | `created` |
| `videoMesh` | (varies) |

Use `all` as the event to subscribe to all events for a resource.  
Filter syntax: `roomId=<id>`, `personEmail=<email>`. Combine with `&`.

### Manage
```bash
webex messaging webhooks list [--owned-by org]
webex messaging webhooks get --webhook-id <id>
webex messaging webhooks update --webhook-id <id> --body '{"targetUrl":"https://new"}'
webex messaging webhooks delete --webhook-id <id>
```

## Other Resources

```bash
# People
webex messaging people list --email user@example.com
webex messaging people get-my-own
webex messaging people get-person --person-id <id>

# Adaptive Card actions (from card button clicks)
webex messaging attachment-actions get --attachment-action-id <id>

# Room tabs
webex messaging room-tabs list --room-id <roomId>
webex messaging room-tabs create --body '{"roomId":"<id>","displayName":"Wiki","contentUrl":"https://..."}'

# Compliance events
webex messaging events list --resource messages --type created
```

## Key Gotchas

1. **Bots in group spaces** — bots get an empty result from `messages list` without `--mentioned-people me`. Always include this flag for bot tokens listing group space messages.

2. **`list` vs `list-direct`** — `messages list` requires `--room-id` (space messages). `messages list-direct` uses `--person-email` or `--person-id` (1:1 DM thread, no roomId needed).

3. **Threading** — to reply in a thread: include `"parentId": "<messageId>"` in the create body. To fetch replies: `messages list --room-id <roomId> --parent-id <messageId>`.

4. **Rich text** — use `markdown` for formatted messages. `text` is the plain-text fallback. Provide both when targeting mixed clients.

5. **1:1 rooms** — `rooms list --type direct` returns 1:1 spaces. The room `title` is the other person's display name. Send a message with `toPersonEmail` to start a new 1:1 — Webex creates the room automatically.

6. **Compliance officer memberships** — query all rooms a person belongs to org-wide via `memberships list --person-email <email>` (no `--room-id`). Requires compliance officer scope.

<!-- codegen:start -->
## Command Reference

> Auto-generated from Postman collections. Run `make codegen` to update.

### attachment-actions

| Command | Flags |
|---|---|
| `get` | `--id` *(required)* |
| `create` | `--body`, `--body-file` |

### ecm-folder-linking

| Command | Flags |
|---|---|
| `list` | `--room-id` |
| `get` | `--id` *(required)* |
| `create-configuration` | `--room-id`, `--content-url`, `--display-name`, `--drive-id`, `--item-id`, `--default-folder`, `--body`, `--body-file` |
| `update-linked` | `--id` *(required)*, `--room-id`, `--content-url`, `--display-name`, `--drive-id`, `--item-id`, `--default-folder`, `--body`, `--body-file` |
| `unlink-linked` | `--id` *(required)* |

### events

| Command | Flags |
|---|---|
| `list` | `--resource`, `--type`, `--actor-id`, `--from`, `--to`, `--max`, `--service-type`, `--last` |
| `get` | `--event-id` *(required)* |

### memberships

| Command | Flags |
|---|---|
| `list` | `--room-id`, `--person-id`, `--person-email`, `--max` |
| `get` | `--membership-id` *(required)* |
| `create` | `--room-id`, `--person-id`, `--person-email`, `--is-moderator`, `--body`, `--body-file` |
| `update` | `--membership-id` *(required)*, `--is-moderator`, `--is-room-hidden`, `--body`, `--body-file` |
| `delete` | `--membership-id` *(required)* |

### messages

| Command | Flags |
|---|---|
| `list` | `--room-id`, `--parent-id`, `--mentioned-people`, `--before`, `--before-message`, `--max` |
| `list-direct` | `--parent-id`, `--person-id`, `--person-email` |
| `get` | `--message-id` *(required)* |
| `create` | `--body`, `--body-file` |
| `edit` | `--message-id` *(required)*, `--room-id`, `--text`, `--markdown`, `--body`, `--body-file` |
| `delete` | `--message-id` *(required)* |

### people

| Command | Flags |
|---|---|
| `list` | `--email`, `--display-name`, `--id`, `--org-id`, `--roles`, `--calling-data`, `--location-id`, `--max`, `--exclude-status` |
| `get-person` | `--person-id` *(required)*, `--calling-data` |
| `get-my-own` | `--calling-data` |
| `create-person` | `--calling-data`, `--min-response`, `--body`, `--body-file` |
| `update-person` | `--person-id` *(required)*, `--calling-data`, `--show-all-types`, `--min-response`, `--body`, `--body-file` |
| `delete-person` | `--person-id` *(required)* |

### room-tabs

| Command | Flags |
|---|---|
| `list` | `--room-id` |
| `get` | `--id` *(required)* |
| `create` | `--room-id`, `--content-url`, `--display-name`, `--body`, `--body-file` |
| `update` | `--id` *(required)*, `--room-id`, `--content-url`, `--display-name`, `--body`, `--body-file` |
| `delete` | `--id` *(required)* |

### rooms

| Command | Flags |
|---|---|
| `list` | `--team-id`, `--type`, `--org-public-spaces`, `--from`, `--to`, `--sort-by`, `--max`, `--last` |
| `get` | `--room-id` *(required)* |
| `get-meeting` | `--room-id` *(required)* |
| `create` | `--title`, `--team-id`, `--classification-id`, `--is-locked`, `--is-public`, `--description`, `--is-announcement-only`, `--body`, `--body-file` |
| `update` | `--room-id` *(required)*, `--title`, `--classification-id`, `--team-id`, `--is-locked`, `--is-public`, `--description`, `--is-announcement-only`, `--is-read-only`, `--body`, `--body-file` |
| `delete` | `--room-id` *(required)* |

### team-memberships

| Command | Flags |
|---|---|
| `list` | `--team-id`, `--max` |
| `get` | `--membership-id` *(required)* |
| `create` | `--team-id`, `--person-id`, `--person-email`, `--is-moderator`, `--body`, `--body-file` |
| `update` | `--membership-id` *(required)*, `--is-moderator`, `--body`, `--body-file` |
| `delete` | `--membership-id` *(required)* |

### teams

| Command | Flags |
|---|---|
| `list` | `--max` |
| `get` | `--team-id` *(required)*, `--description` |
| `create` | `--name`, `--description`, `--body`, `--body-file` |
| `update` | `--team-id` *(required)*, `--name`, `--description`, `--body`, `--body-file` |
| `delete` | `--team-id` *(required)* |

### webhooks

| Command | Flags |
|---|---|
| `list` | `--max`, `--owned-by` |
| `get` | `--webhook-id` *(required)* |
| `create` | `--name`, `--target-url`, `--resource`, `--event`, `--filter`, `--secret`, `--owned-by`, `--body`, `--body-file` |
| `update` | `--webhook-id` *(required)*, `--name`, `--target-url`, `--secret`, `--owned-by`, `--status`, `--body`, `--body-file` |
| `delete` | `--webhook-id` *(required)* |

### hds

| Command | Flags |
|---|---|
| `get-test-results-node` | `--node-id` *(required)*, `--trigger-type` |
| `get-cluster` | `--cluster-id` *(required)* |
| `get-org` | `--organization-id` *(required)* |
| `get-node` | `--node-id` *(required)* |
| `get-database-org` | `--organization-id` *(required)* |
| `get-multi-tenant-org` | `--organization-id` *(required)* |
| `get-availability-cluster` | `--cluster-id` *(required)*, `--from`, `--to`, `--last` |
| `get-database-org-2` | `--organization-id` *(required)* |
| `get-multi-tenant-org-2` | `--organization-id` *(required)* |

<!-- codegen:end -->
