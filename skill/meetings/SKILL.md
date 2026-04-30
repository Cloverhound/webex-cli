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

<!-- codegen:start -->
## Command Reference

> Auto-generated from Postman collections. Run `make codegen` to update.

### chats

| Command | Flags |
|---|---|
| `list-meeting` | `--meeting-id`, `--max`, `--offset` |
| `delete-meeting` | `--meeting-id` |

### closed-captions

| Command | Flags |
|---|---|
| `list-meeting` | `--meeting-id` |
| `list-meeting-snippets` | `--closed-caption-id` *(required)*, `--meeting-id` |
| `download-meeting-snippets` | `--closed-caption-id` *(required)*, `--format`, `--meeting-id` |

### invitees

| Command | Flags |
|---|---|
| `list-meeting` | `--meeting-id`, `--max`, `--host-email`, `--panelist` |
| `get-meeting` | `--meeting-invitee-id` *(required)*, `--host-email` |
| `create-meeting` | `--meeting-id`, `--email`, `--display-name`, `--co-host`, `--host-email`, `--send-email`, `--panelist`, `--body`, `--body-file` |
| `create-meeting-2` | `--body`, `--body-file` |
| `update-meeting` | `--meeting-invitee-id` *(required)*, `--email`, `--display-name`, `--co-host`, `--host-email`, `--send-email`, `--panelist`, `--body`, `--body-file` |
| `delete-meeting` | `--meeting-invitee-id` *(required)*, `--host-email`, `--send-email` |

### messages

| Command | Flags |
|---|---|
| `delete-meeting` | `--meeting-message-id` *(required)* |

### participants

| Command | Flags |
|---|---|
| `list-meeting` | `--max`, `--meeting-id`, `--breakout-session-id`, `--meeting-start-time-from`, `--meeting-start-time-to`, `--host-email`, `--join-time-from`, `--join-time-to` |
| `get-meeting` | `--participant-id` *(required)*, `--host-email` |
| `query-meeting-email` | `--meeting-id`, `--meeting-start-time-from`, `--meeting-start-time-to`, `--host-email`, `--emails`, `--join-time-from`, `--join-time-to`, `--body`, `--body-file` |
| `admit` | `--body`, `--body-file` |
| `call-out-sip` | `--address`, `--display-name`, `--meeting-id`, `--meeting-number`, `--address-type`, `--invitation-correlation-id`, `--body`, `--body-file` |
| `cancel-calling-out-sip` | `--participant-id`, `--body`, `--body-file` |
| `update` | `--participant-id` *(required)*, `--muted`, `--admit`, `--expel`, `--body`, `--body-file` |

### meeting-polls

| Command | Flags |
|---|---|
| `list` | `--meeting-id` |
| `get-pollresults` | `--meeting-id`, `--max` |
| `list-respondents-question` | `--poll-id` *(required)*, `--question-id` *(required)*, `--meeting-id`, `--max` |

### preferences

| Command | Flags |
|---|---|
| `get-meeting` | `--user-email`, `--site-url` |
| `get-personal-meeting-room-options` | `--user-email`, `--site-url` |
| `get-audio-options` | `--user-email`, `--site-url` |
| `get-video-options` | `--user-email`, `--site-url` |
| `get-scheduling-options` | `--user-email`, `--site-url` |
| `get-site-list` | `--user-email`, `--site-url` |
| `insert-delegate-emails` | `--user-email`, `--site-url`, `--emails`, `--body`, `--body-file` |
| `delete-delegate-emails` | `--user-email`, `--site-url`, `--emails`, `--body`, `--body-file` |
| `batch-refresh-personal-meeting-room-id` | `--body`, `--body-file` |
| `update-personal-meeting-room-options` | `--user-email`, `--site-url`, `--body`, `--body-file` |
| `update-audio-options` | `--user-email`, `--site-url`, `--body`, `--body-file` |
| `update-video-options` | `--user-email`, `--site-url`, `--body`, `--body-file` |
| `update-scheduling-options` | `--user-email`, `--site-url`, `--enabled-join-before-host`, `--join-before-host-minutes`, `--enabled-auto-share-recording`, `--enabled-webex-assistant-by-default`, `--delegate-emails`, `--body`, `--body-file` |
| `update-default-site` | `--default-site`, `--user-email`, `--site-url`, `--body`, `--body-file` |

### meeting-qa

| Command | Flags |
|---|---|
| `list-q` | `--meeting-id`, `--max` |
| `list-answers-question` | `--question-id` *(required)*, `--meeting-id`, `--max` |

### meeting-qualities

| Command | Flags |
|---|---|
| `get` | `--meeting-id`, `--max`, `--offset` |

### session-types

| Command | Flags |
|---|---|
| `list-site` | `--site-url` |
| `list-user` | `--site-url`, `--person-id` |
| `update-user` | `--site-url`, `--session-type-ids`, `--person-id`, `--email`, `--body`, `--body-file` |

### transcripts

| Command | Flags |
|---|---|
| `list-meeting` | `--max`, `--from`, `--to`, `--meeting-id`, `--host-email`, `--site-url`, `--last` |
| `list-meeting-compliance-officer` | `--from`, `--to`, `--max`, `--site-url`, `--last` |
| `download-meeting` | `--transcript-id` *(required)*, `--format`, `--host-email` |
| `list-snippets-meeting` | `--transcript-id` *(required)*, `--max` |
| `get-snippet` | `--transcript-id` *(required)*, `--snippet-id` *(required)* |
| `update-snippet` | `--transcript-id` *(required)*, `--snippet-id` *(required)*, `--text`, `--reason`, `--body`, `--body-file` |
| `delete` | `--transcript-id` *(required)*, `--reason`, `--comment`, `--body`, `--body-file` |

### meeting-summary

| Command | Flags |
|---|---|
| `list-usage-reports` | `--site-url`, `--service-type`, `--from`, `--to`, `--max`, `--last` |
| `list-attendee-reports` | `--site-url`, `--from`, `--to`, `--max`, `--meeting-id`, `--meeting-number`, `--meeting-title`, `--last` |

### meetings

| Command | Flags |
|---|---|
| `get-admin` | `--meeting-id` *(required)*, `--current` |
| `list-admin` | `--meeting-number`, `--web-link`, `--current` |
| `list` | `--meeting-number`, `--web-link`, `--room-id`, `--meeting-series-id`, `--max`, `--from`, `--to`, `--meeting-type`, `--state`, `--scheduled-type`, `--is-modified`, `--has-chat`, `--has-recording`, `--has-transcription`, `--has-summary`, `--has-closed-caption`, `--has-polls`, `--has-qa`, `--has-slido`, `--current`, `--host-email`, `--site-url`, `--integration-tag`, `--last` |
| `get` | `--meeting-id` *(required)*, `--current`, `--host-email` |
| `list-templates` | `--template-type`, `--locale`, `--is-default`, `--is-standard`, `--host-email`, `--site-url` |
| `get-template` | `--template-id` *(required)*, `--host-email` |
| `get-control-status` | `--meeting-id` |
| `list-session-types` | `--host-email`, `--site-url` |
| `get-session-type` | `--session-type-id` *(required)*, `--host-email`, `--site-url` |
| `get-registration-form` | `--meeting-id` *(required)*, `--current`, `--host-email` |
| `list-registrants` | `--meeting-id` *(required)*, `--max`, `--host-email`, `--current`, `--email`, `--registration-time-from`, `--registration-time-to` |
| `get-detailed-information-registrant` | `--meeting-id` *(required)*, `--registrant-id` *(required)*, `--current`, `--host-email` |
| `list-interpreters` | `--meeting-id` *(required)*, `--host-email` |
| `get-interpreter` | `--meeting-id` *(required)*, `--interpreter-id` *(required)*, `--host-email` |
| `list-breakout-sessions` | `--meeting-id` *(required)* |
| `get-survey` | `--meeting-id` *(required)* |
| `list-survey-results` | `--meeting-id` *(required)*, `--meeting-start-time-from`, `--meeting-start-time-to`, `--max` |
| `list-invitation-sources` | `--meeting-id` *(required)* |
| `list-tracking-codes` | `--site-url`, `--service`, `--host-email` |
| `create` | `--body`, `--body-file` |
| `join` | `--meeting-id`, `--meeting-number`, `--web-link`, `--join-directly`, `--email`, `--display-name`, `--password`, `--expiration-minutes`, `--registration-id`, `--host-email`, `--create-join-link-as-web-link`, `--create-start-link-as-web-link`, `--body`, `--body-file` |
| `register-registrant` | `--meeting-id` *(required)*, `--current`, `--host-email`, `--body`, `--body-file` |
| `batch-register-registrants` | `--meeting-id` *(required)*, `--current`, `--host-email`, `--body`, `--body-file` |
| `query-registrants` | `--meeting-id` *(required)*, `--max`, `--current`, `--host-email`, `--emails`, `--status`, `--order-type`, `--order-by`, `--body`, `--body-file` |
| `batch-update-registrants-status` | `--meeting-id` *(required)*, `--status-op-type` *(required)*, `--current`, `--host-email`, `--body`, `--body-file` |
| `create-interpreter` | `--meeting-id` *(required)*, `--language-code1`, `--language-code2`, `--email`, `--display-name`, `--host-email`, `--send-email`, `--body`, `--body-file` |
| `get-survey-links` | `--meeting-id` *(required)*, `--host-email`, `--meeting-start-time-from`, `--meeting-start-time-to`, `--emails`, `--body`, `--body-file` |
| `create-invitation-sources` | `--meeting-id` *(required)*, `--body`, `--body-file` |
| `reassign-host` | `--host-email`, `--meeting-ids`, `--body`, `--body-file` |
| `end` | `--meeting-id` *(required)*, `--reason`, `--body`, `--body-file` |
| `batch-approve-registrants` | `--meeting-id` *(required)*, `--current`, `--host-email`, `--body`, `--body-file` |
| `batch-reject-registrants` | `--meeting-id` *(required)*, `--current`, `--host-email`, `--body`, `--body-file` |
| `batch-cancel-registrants` | `--meeting-id` *(required)*, `--current`, `--host-email`, `--body`, `--body-file` |
| `batch-delete-registrants` | `--meeting-id` *(required)*, `--current`, `--host-email`, `--body`, `--body-file` |
| `update` | `--meeting-id` *(required)*, `--body`, `--body-file` |
| `update-control-status` | `--meeting-id`, `--recording-started`, `--recording-paused`, `--locked`, `--body`, `--body-file` |
| `update-registration-form` | `--meeting-id` *(required)*, `--body`, `--body-file` |
| `update-simultaneous-interpretation` | `--meeting-id` *(required)*, `--body`, `--body-file` |
| `update-interpreter` | `--meeting-id` *(required)*, `--interpreter-id` *(required)*, `--language-code1`, `--language-code2`, `--email`, `--display-name`, `--host-email`, `--send-email`, `--body`, `--body-file` |
| `update-breakout-sessions` | `--meeting-id` *(required)*, `--body`, `--body-file` |
| `patch` | `--meeting-id` *(required)*, `--body`, `--body-file` |
| `delete` | `--meeting-id` *(required)*, `--host-email`, `--send-email` |
| `delete-registration-form` | `--meeting-id` *(required)* |
| `delete-registrant` | `--meeting-id` *(required)*, `--registrant-id` *(required)*, `--current`, `--host-email` |
| `delete-interpreter` | `--meeting-id` *(required)*, `--interpreter-id` *(required)*, `--host-email`, `--send-email` |
| `delete-breakout-sessions` | `--meeting-id` *(required)*, `--send-email` |

### tracking-codes

| Command | Flags |
|---|---|
| `list` | `--site-url` |
| `get` | `--tracking-code-id` *(required)*, `--site-url` |
| `get-user` | `--site-url`, `--person-id` |
| `create` | `--body`, `--body-file` |
| `update` | `--tracking-code-id` *(required)*, `--body`, `--body-file` |
| `update-user` | `--body`, `--body-file` |
| `delete` | `--tracking-code-id` *(required)*, `--site-url` |

### people

| Command | Flags |
|---|---|
| `list` | `--email`, `--display-name`, `--id`, `--org-id`, `--roles`, `--calling-data`, `--location-id`, `--max`, `--exclude-status` |
| `get-person` | `--person-id` *(required)*, `--calling-data` |
| `get-my-own` | `--calling-data` |
| `create-person` | `--calling-data`, `--min-response`, `--body`, `--body-file` |
| `update-person` | `--person-id` *(required)*, `--calling-data`, `--show-all-types`, `--min-response`, `--body`, `--body-file` |
| `delete-person` | `--person-id` *(required)* |

### recording-report

| Command | Flags |
|---|---|
| `list-audit-summaries` | `--max`, `--from`, `--to`, `--host-email`, `--site-url`, `--last` |
| `get-audit` | `--recording-id`, `--host-email`, `--max` |
| `list-meeting-archive-summaries` | `--max`, `--from`, `--to`, `--site-url`, `--last` |
| `get-meeting-archive` | `--archive-id` *(required)* |

### recordings

| Command | Flags |
|---|---|
| `list` | `--max`, `--from`, `--to`, `--meeting-id`, `--host-email`, `--site-url`, `--integration-tag`, `--topic`, `--format`, `--service-type`, `--status`, `--last` |
| `list-admin-compliance-officer` | `--max`, `--from`, `--to`, `--meeting-id`, `--site-url`, `--integration-tag`, `--topic`, `--format`, `--service-type`, `--status`, `--last` |
| `get` | `--recording-id` *(required)*, `--host-email` |
| `list-group` | `--person-id`, `--max`, `--from`, `--to`, `--site-url`, `--integration-tag`, `--topic`, `--format`, `--service-type`, `--last` |
| `get-group` | `--recording-id` *(required)*, `--person-id` |
| `move-recycle-bin` | `--host-email`, `--recording-ids`, `--site-url`, `--body`, `--body-file` |
| `restore-recycle-bin` | `--host-email`, `--restore-all`, `--recording-ids`, `--site-url`, `--body`, `--body-file` |
| `purge-recycle-bin` | `--host-email`, `--purge-all`, `--recording-ids`, `--site-url`, `--body`, `--body-file` |
| `share` | `--recording-id` *(required)*, `--host-email`, `--add-emails`, `--remove-emails`, `--send-email`, `--body`, `--body-file` |
| `share-link` | `--host-email`, `--web-share-link`, `--add-emails`, `--remove-emails`, `--send-email`, `--body`, `--body-file` |
| `delete-admin` | `--recording-id` *(required)* |
| `delete` | `--recording-id` *(required)*, `--host-email`, `--reason`, `--comment`, `--body`, `--body-file` |

### site

| Command | Flags |
|---|---|
| `get-meeting-common-settings-configuration` | `--site-url` |
| `update-meeting-common-settings-configuration` | `--body`, `--body-file` |

### slido

| Command | Flags |
|---|---|
| `list-compliance-events` | `--session-org-id`, `--session-id`, `--start` |

### video-mesh

| Command | Flags |
|---|---|
| `list-clusters-availability` | `--from`, `--to`, `--org-id`, `--last` |
| `get-cluster-availability` | `--cluster-id` *(required)*, `--from`, `--to`, `--last` |
| `list-node-availability` | `--from`, `--to`, `--cluster-id`, `--last` |
| `get-node-availability` | `--node-id` *(required)*, `--from`, `--to`, `--last` |
| `list-media-health-monitoring-tool-test-results-v2` | `--org-id`, `--trigger-type`, `--from`, `--to`, `--last` |
| `get-media-health-monitoring-tool-test-results-clusters-v2` | `--cluster-id`, `--trigger-type`, `--from`, `--to`, `--last` |
| `get-media-health-monitoring-tool-test-results-node-v2` | `--node-id`, `--trigger-type`, `--from`, `--to`, `--last` |
| `list-overflow-cloud` | `--from`, `--to`, `--org-id`, `--last` |
| `list-cluster-redirect` | `--from`, `--to`, `--org-id`, `--last` |
| `get-cluster-redirect` | `--from`, `--to`, `--cluster-id`, `--last` |
| `list-clusters-utilization` | `--from`, `--to`, `--org-id`, `--last` |
| `get-cluster-utilization` | `--from`, `--to`, `--cluster-id`, `--last` |
| `list-reachability-test-results-v2` | `--org-id`, `--trigger-type`, `--from`, `--to`, `--last` |
| `get-reachability-test-results-cluster-v2` | `--cluster-id`, `--trigger-type`, `--from`, `--to`, `--last` |
| `get-reachability-test-results-node-v2` | `--node-id`, `--trigger-type`, `--from`, `--to`, `--last` |
| `list-cluster` | `--org-id` |
| `get-cluster` | `--cluster-id` *(required)* |
| `get-triggered-test-status` | `--command-id` |
| `get-triggered-test-results` | `--command-id` |
| `list-network-test-results` | `--org-id`, `--trigger-type`, `--from`, `--to`, `--last` |
| `get-network-test-results-cluster` | `--cluster-id`, `--trigger-type`, `--from`, `--to`, `--last` |
| `get-network-test-results-node` | `--node-id`, `--trigger-type`, `--from`, `--to`, `--last` |
| `list-cluster-client-type-distribution` | `--org-id`, `--from`, `--to`, `--device-type`, `--last` |
| `get-cluster-client-type-distribution` | `--cluster-id`, `--from`, `--to`, `--device-type`, `--last` |
| `list-event-threshold-configuration` | `--event-name`, `--cluster-id`, `--org-id`, `--event-scope` |
| `get-event-threshold-configuration` | `--event-threshold-id` *(required)* |
| `trigger-on-demand-test-cluster` | `--cluster-id` *(required)*, `--type`, `--nodes`, `--body`, `--body-file` |
| `trigger-on-demand-test-node` | `--node-id` *(required)*, `--type`, `--body`, `--body-file` |
| `reset-event-threshold-configuration` | `--event-threshold-ids`, `--body`, `--body-file` |
| `update-event-threshold-configuration` | `--body`, `--body-file` |

### webhooks

| Command | Flags |
|---|---|
| `list` | `--max`, `--owned-by` |
| `get` | `--webhook-id` *(required)* |
| `create` | `--name`, `--target-url`, `--resource`, `--event`, `--filter`, `--secret`, `--owned-by`, `--body`, `--body-file` |
| `update` | `--webhook-id` *(required)*, `--name`, `--target-url`, `--secret`, `--owned-by`, `--status`, `--body`, `--body-file` |
| `delete` | `--webhook-id` *(required)* |

### summaries

| Command | Flags |
|---|---|
| `get-meeting-id` | `--meeting-id` |
| `get-compliance-officer` | `--meeting-id` |
| `delete` | `--summary-id` *(required)* |

### slidosecurepremium

| Command | Flags |
|---|---|
| `list-compliance-events` | `--session-org-id`, `--session-id`, `--start` |

<!-- codegen:end -->
