---
name: webex-cli/meetings
description: "Webex Meetings commands: scheduling, recordings, transcripts, participants, invitees, preferences, and site settings."
---

# Webex Meetings

Commands: `webex meetings <resource> <action> [flags]`  
Alias: `webex meeting <resource> <action>`

## Resources

| Resource | Common operations |
|---|---|
| `meetings` | list, get, create, update, delete |
| `recordings` | list, list-admin-compliance-officer, get, delete, move-recycle-bin, restore-recycle-bin, purge-recycle-bin |
| `transcripts` | list-meeting, download-meeting, list-snippets-meeting, get-snippet, update-snippet, delete |
| `participants` | list-meeting, get-meeting, update, admit |
| `invitees` | list-meeting, get-meeting, create-meeting, create-meeting-2, update-meeting, delete-meeting |
| `preferences` | get-site-list, get-meeting, get-personal-meeting-room-options, get-scheduling-options, get-audio-options |
| `session-types` | list-site, list-user |
| `tracking-codes` | list, get, create, update, delete |
| `summaries` | get-meeting-id, get-compliance-officer, delete |
| `site` | get-meeting-common-settings-configuration |

## Meetings

### List meetings
```bash
webex meetings meetings list [flags]
  --from <ISO8601>           # start date (default: 7 days before --to or now)
  --to <ISO8601>             # end date
  --last <duration>          # shorthand: 1h, 24h, 7d — sets --from automatically
  --state scheduled|ready|lobby|inProgress|ended|missed|expired
  --meeting-type meetingSeries|scheduledMeeting|meeting
  --meeting-number <num>     # look up a specific meeting by number
  --meeting-series-id <id>   # all occurrences of a series
  --room-id <id>             # meetings associated with a Webex space
  --host-email <email>       # admin only: list for a different host
  --site-url <url>           # specific Webex site
  --max <n>                  # max 100 (default 10)
  --has-recording true|false
  --has-transcription true|false
```

### Get / create / delete
```bash
webex meetings meetings get --meeting-id <id>
webex meetings meetings create --body '{
  "title": "Weekly Sync",
  "start": "2025-01-15T10:00:00Z",
  "end": "2025-01-15T11:00:00Z",
  "invitees": [{"email":"user@example.com"}],
  "enabledJoinBeforeHost": true,
  "password": "abc123"
}'
webex meetings meetings update --meeting-id <id> --body '{...}'
webex meetings meetings delete --meeting-id <id>
```

**Common meeting body fields:** `title`, `start`, `end`, `timezone`, `agenda`, `password`, `enabledJoinBeforeHost`, `joinBeforeHostMinutes`, `invitees` (array), `recurrence`, `siteUrl`, `scheduledType`

## Recordings

### List own recordings
```bash
webex meetings recordings list [flags]
  --from <ISO8601>
  --to <ISO8601>
  --last <duration>          # e.g. 7d, 720h
  --meeting-id <id>          # recordings for a specific meeting
  --topic <text>             # filter by topic (case-insensitive)
  --status available|deleted|purged
  --format <format>          # filter by file format
  --site-url <url>
  --max <n>                  # max 100 per page
```

### List org recordings (admin / compliance officer)
```bash
webex meetings recordings list-admin-compliance-officer [same flags]
# Max 30-day window per request. Requires admin-level recording scope.
```

### Get, delete, download
```bash
webex meetings recordings get --recording-id <id>
webex meetings recordings delete --recording-id <id>
webex meetings recordings delete-admin --recording-id <id>   # admin delete

# Recycle bin
webex meetings recordings move-recycle-bin --body '{"recordingIds":["<id>"]}'
webex meetings recordings restore-recycle-bin --body '{"recordingIds":["<id>"]}'
webex meetings recordings purge-recycle-bin --body '{"recordingIds":["<id>"]}'
```

Note: the `recordings` resource does not have a built-in `download` command. Use the `temporaryDirectDownloadLinks` from the GET response to download via curl or similar.

## Transcripts

```bash
# List transcripts for all meetings or a specific meeting
webex meetings transcripts list-meeting [flags]
  --meeting-id <id>          # filter to a specific meeting instance
  --from <ISO8601>
  --to <ISO8601>
  --last <duration>
  --site-url <url>
  --max <n>                  # max 100 per page

# Admin / compliance officer
webex meetings transcripts list-meeting-compliance-officer [same flags + --host-email]

# Download full transcript file
webex meetings transcripts download-meeting --transcript-id <id>

# Transcript snippets (line-by-line)
webex meetings transcripts list-snippets-meeting --transcript-id <id>
webex meetings transcripts get-snippet --transcript-id <id> --snippet-id <id>
webex meetings transcripts update-snippet --transcript-id <id> --snippet-id <id> --body '{"text":"Corrected text"}'

# Delete
webex meetings transcripts delete --transcript-id <id>
```

## Participants

```bash
# List participants for an ended meeting
webex meetings participants list-meeting [flags]
  --meeting-id <id>          # required
  --join-time-from <ISO8601>
  --join-time-to <ISO8601>
  --meeting-start-time-from <ISO8601>
  --meeting-start-time-to <ISO8601>
  --host-email <email>       # admin only
  --max <n>                  # max 100 per page

webex meetings participants get-meeting --meeting-id <id> --participant-id <id>

# Live meeting controls
webex meetings participants admit --body '{"items":[{"id":"<participantId>"}]}'
webex meetings participants update --meeting-id <id> --participant-id <id> --body '{"muted":true}'
```

## Invitees

```bash
# List invitees for a meeting
webex meetings invitees list-meeting --meeting-id <id>

# Get / create / update / delete
webex meetings invitees get-meeting --meeting-id <id> --invitee-id <id>
webex meetings invitees create-meeting --body '{"meetingId":"<id>","email":"user@example.com","displayName":"User","coHost":false}'
webex meetings invitees create-meeting-2 --body '{"meetingId":"<id>","items":[{"email":"u1@example.com"},{"email":"u2@example.com"}]}'
webex meetings invitees update-meeting --meeting-id <id> --invitee-id <id> --body '{"coHost":true}'
webex meetings invitees delete-meeting --meeting-id <id> --invitee-id <id>
```

## Preferences

```bash
# Site list (user's available Webex sites)
webex meetings preferences get-site-list

# Set default site
webex meetings preferences update-default-site --body '{"defaultSite":"company.webex.com"}'

# Meeting preferences
webex meetings preferences get-meeting
webex meetings preferences get-scheduling-options
webex meetings preferences update-scheduling-options --body '{...}'

# Personal Meeting Room
webex meetings preferences get-personal-meeting-room-options
webex meetings preferences update-personal-meeting-room-options --body '{...}'

# Audio / video options
webex meetings preferences get-audio-options
webex meetings preferences update-audio-options --body '{...}'
```

## Session Types and Tracking Codes

```bash
# Session types
webex meetings session-types list-site --site-url <url>
webex meetings session-types list-user

# Tracking codes (for meeting analytics/reporting)
webex meetings tracking-codes list [--site-url <url>]
webex meetings tracking-codes create --body '{"siteUrl":"company.webex.com","name":"Department","inputMode":"select","hostProfileCode":true}'
webex meetings tracking-codes get --tracking-code-id <id>
webex meetings tracking-codes update --tracking-code-id <id> --body '{...}'
webex meetings tracking-codes delete --tracking-code-id <id>

# Per-user tracking code values
webex meetings tracking-codes get-user --site-url <url>
webex meetings tracking-codes update-user --body '{...}'
```

## Meeting Summaries and Q&A

```bash
# AI-generated meeting summaries
webex meetings summaries get-meeting-id --meeting-id <id>
webex meetings summaries get-compliance-officer --meeting-id <id>   # admin
webex meetings summaries delete --meeting-id <id>

# Meeting Q&A
webex meetings meeting-qa list --meeting-id <id>

# Polls
webex meetings meeting-polls list --meeting-id <id>
```

## Key Gotchas

1. **`--last` vs `--from`/`--to`** — `--last 7d` sets `--from` automatically (shorthand). Do not combine `--last` with `--from` — they conflict.

2. **Recording `list-admin-compliance-officer`** — enforces a **30-day** maximum window per request (`from` to `to`). Use multiple requests to cover longer periods.

3. **Meeting type defaults** — `meetings list` returns `meetingSeries` by default. To list specific scheduled occurrences use `--meeting-type scheduledMeeting`; to list ended instances use `--meeting-type meeting`.

4. **Ended meeting detail delay** — details of an `ended` meeting are only available 15 minutes after it ends. The API may return limited data for `inProgress` meetings.

5. **Transcript vs summary** — `transcripts` are verbatim (speaker + text). `summaries` are AI-generated highlights. Both require the meeting to have the feature enabled and be ended.

6. **Personal Room meetings** — the meeting ID of a scheduled personal room meeting is not supported by the Transcripts or Participants APIs.

7. **Admin scope for list** — `list-admin-compliance-officer` (recordings and transcripts) requires `spark-admin:meetings_read` or compliance officer scope. `--host-email` on `participants list-meeting` also requires admin scope.
