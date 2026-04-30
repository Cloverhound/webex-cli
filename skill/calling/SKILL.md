---
name: webex-cli/calling
description: "Webex Calling commands: locations, people, devices, call queues, hunt groups, auto attendants, recordings, and telephony configuration."
---

# Webex Calling

Commands: `webex calling <resource> <action> [flags]`

No `--orgid` needed — org is inferred from the login token.

## Key Resources

| Resource | Common operations |
|---|---|
| `locations` | list, get, create, update, delete, list-floors, create-floor |
| `people` | list, get-person, get-my-own, create-person, update-person, delete-person |
| `devices` | list, get, create-mac-address, create-activation-code, delete |
| `call-queue` | list-cxe, get-cxe, create-cxe, delete |
| `hunt-group` | list, get, create, update, delete |
| `auto-attendant` | list, get, create, update, delete |
| `numbers` | get-phone-org, create-phone-location, validate-phone, manage-state-location |
| `converged-recordings` | list, list-admin-compliance-officer, get, download, delete |
| `announcement-repository` | upload-binary-greeting, update-binary-greeting, get-binary-greeting, get-usage |
| `user-call` | (per-user calling settings — forwarding, voicemail, schedules) |
| `call-settings-for-me` | (self-service settings for the authenticated user) |
| `location-voicemail` | get, update, list-voicemailgroup, get-group |
| `call-recording` | get, update, list-jobs |
| `call-routing` | create-dial-plan, list-dial-plans, create-trunk, list-trunks |
| `customer-experience-essentials` | (CxE queue/agent management) |

## Locations

```bash
# List
webex calling locations list [--name <name>] [--id <id>] [--max <n>] [--paginate]

# Get
webex calling locations get --location-id <id>

# Create / update
webex calling locations create --body '{"name":"Main Office","timeZone":"America/New_York","preferredLanguage":"en_US","address":{...},"announcementLanguage":"en_US"}'
webex calling locations update --location-id <id> --body '{...}'

# Floors
webex calling locations list-floors --location-id <id>
webex calling locations create-floor --location-id <id> --floor-number 1 --display-name "Lobby"
webex calling locations delete --location-id <id>
```

## People

```bash
# List — use --paginate to get all results
webex calling people list [flags]
  --email <email>            # exact email match
  --display-name <name>      # starts-with name match
  --location-id <id>         # filter by calling location
  --id <id1,id2,...>         # up to 85 IDs
  --max <n>                  # max results (≤100 with --calling-data)
  --calling-data true        # include calling configuration
  --roles <roleId,...>       # filter by role IDs

webex calling people get-person --person-id <id>
webex calling people get-my-own
webex calling people create-person --body '{"emails":["user@example.com"],"displayName":"..."}'
webex calling people update-person --person-id <id> --body '{...}'
webex calling people delete-person --person-id <id>
```

## Devices

```bash
# List
webex calling devices list [flags]
  --product <name>           # e.g. "Cisco 8841"
  --person-id <id>           # devices assigned to a person
  --workspace-id <id>        # devices in a workspace
  --location-id <id>
  --mac <address>            # by MAC address
  --display-name <name>
  --connection-status <s>    # e.g. connected, disconnected
  --tag <tag>                # filter by tag (repeatable)
  --max <n>

webex calling devices get --device-id <id>
webex calling devices create-mac-address --body '{"mac":"00:11:22:33:44:55","workspaceId":"<id>"}'
webex calling devices create-activation-code --body '{"workspaceId":"<id>"}'
webex calling devices delete --device-id <id>
```

## Call Queues (Customer Experience)

```bash
# List (use list-cxe for CxE queues)
webex calling call-queue list-cxe [flags]
  --location-id <id>
  --name <name>
  --phone-number <num>
  --max <n>
  --has-cx-essentials true   # only CxE queues

webex calling call-queue get-cxe --location-id <id> --queue-id <id>
webex calling call-queue create-cxe --body '{...}'
webex calling call-queue delete --location-id <id> --queue-id <id>

# Supervisor management
webex calling call-queue create-supervisor-cxe --body '{...}'
webex calling call-queue assign-unassign-agents-supervisor-cxe --body '{...}'
```

## Hunt Groups

```bash
webex calling hunt-group list [--location-id <id>] [--name <name>] [--phone-number <num>] [--max <n>]
webex calling hunt-group get --location-id <id> --hunt-group-id <id>
webex calling hunt-group create --body '{...}'
webex calling hunt-group update --location-id <id> --hunt-group-id <id> --body '{...}'
webex calling hunt-group delete --location-id <id> --hunt-group-id <id>

# Call forwarding
webex calling hunt-group get-call-forward --location-id <id> --hunt-group-id <id>
webex calling hunt-group update-call-forward --location-id <id> --hunt-group-id <id> --body '{...}'
```

## Auto Attendants

```bash
webex calling auto-attendant list [--location-id <id>] [--name <name>] [--phone-number <num>] [--max <n>]
webex calling auto-attendant get --location-id <id> --auto-attendant-id <id>
webex calling auto-attendant create --body '{...}'
webex calling auto-attendant update --location-id <id> --auto-attendant-id <id> --body '{...}'
webex calling auto-attendant delete --location-id <id> --auto-attendant-id <id>
```

## Numbers

```bash
# List numbers in the org
webex calling numbers get-phone-org [--location-id <id>] [--phone-number <num>] [--state <state>] [--paginate]

# Add/remove numbers from a location
webex calling numbers create-phone-location --location-id <id> --body '{"phoneNumbers":[{"number":"+15555551234"}]}'
webex calling numbers delete-phone-location --location-id <id> --body '{...}'

# Validate before adding
webex calling numbers validate-phone --body '{"phoneNumbers":["+15555551234"]}'

# Manage state (active, inactive)
webex calling numbers manage-state-location --location-id <id> --body '{...}'
```

## Converged Recordings (Call Recordings)

```bash
# List own recordings
webex calling converged-recordings list [flags]
  --last <duration>          # e.g. 24h, 7d, 720h
  --from <ISO8601>
  --to <ISO8601>
  --location-id <id>
  --topic <text>
  --status available|deleted
  --service-type calling|customerAssist
  --max <n>                  # max 100 per page

# List ALL org recordings (admin / compliance officer, max 30-day window)
webex calling converged-recordings list-admin-compliance-officer [same flags + --owner-id --owner-email]

# Get / download
webex calling converged-recordings get --recording-id <id>
webex calling converged-recordings download --recording-id <id> --output recording.mp3
webex calling converged-recordings download --recording-id <id> --output recording.vtt --type transcript

# Recycle bin
webex calling converged-recordings move-recycle-bin --body '{"recordingIds":["<id>"]}'
webex calling converged-recordings restore-recycle-bin --body '{"recordingIds":["<id>"]}'
webex calling converged-recordings purge-recycle-bin --body '{"recordingIds":["<id>"]}'
webex calling converged-recordings delete --recording-id <id>
```

## Announcement Repository

```bash
# Upload greeting (org level)
webex calling announcement-repository upload-binary-greeting --file greeting.wav --name "Main Greeting"

# Upload greeting (location level)
webex calling announcement-repository upload-binary-greeting-2 --location-id <id> --file greeting.wav --name "Lobby"

# Update existing
webex calling announcement-repository update-binary-greeting --announcement-id <id> --file updated.wav --name "Updated"
webex calling announcement-repository update-binary-greeting-2 --location-id <id> --announcement-id <id> --file updated.wav

# Get / list
webex calling announcement-repository get-binary-greeting --announcement-id <id>
webex calling announcement-repository list-greetings
webex calling announcement-repository get-usage
webex calling announcement-repository get-usage-location --location-id <id>
```

## User Call Settings

```bash
# For a specific person (admin)
webex calling user-call get-call-forwarding-person --person-id <id>
webex calling user-call update-call-forwarding-person --person-id <id> --body '{...}'
webex calling user-call get-voicemail-person --person-id <id>
webex calling user-call update-voicemail-person --person-id <id> --body '{...}'
webex calling user-call update-busy-voicemail-greeting-person --person-id <id> --file greeting.wav

# For yourself (user token)
webex calling call-settings-for-me get-call-forwarding-settings
webex calling call-settings-for-me update-call-forwarding-settings --body '{...}'
webex calling call-settings-for-me upload-voicemail-busy-greeting --file greeting.wav
webex calling call-settings-for-me upload-voicemail-no-answer-greeting --file greeting.wav
```

## Key Gotchas

1. **Call queue naming** — use `call-queue list-cxe` and `call-queue get-cxe` for Customer Experience queues. The plain `list` and `get` subcommands target the legacy call queue API.

2. **Converged recordings: user vs admin** — `list` returns only the calling user's recordings. `list-admin-compliance-officer` returns org-wide recordings but requires admin scope and enforces a **30-day** max time window per request.

3. **Download requires `temporaryDirectDownloadLinks`** — the `download` command needs this field populated in the recording GET response. A user integration token may not receive download links for recordings owned by other users; a service app token typically does.

4. **Locations and calling names** — location names for Webex Calling must be ≤80 characters even though the API accepts up to 256.

5. **`--paginate` on people list** — with `--calling-data=true`, max per page is capped at 100 by the API. Use `--paginate` to traverse large directories.
