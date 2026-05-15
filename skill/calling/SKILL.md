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

<!-- codegen:start -->
## Command Reference

> Auto-generated from Postman collections. Run `make codegen` to update.

### dect-devices

| Command | Flags |
|---|---|
| `list-networks` | `--org-id`, `--name`, `--location-id` |
| `get-network` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id` |
| `list-network-base-stations` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id` |
| `get-network-base-station` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--base-station-id` *(required)*, `--org-id` |
| `list-handsets-network-id` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id`, `--basestation-id`, `--member-id` |
| `get-network-handset` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--handset-id` *(required)*, `--org-id` |
| `list-networks-associated-person` | `--person-id` *(required)*, `--org-id` |
| `list-networks-associated-workspace` | `--workspace-id` *(required)*, `--org-id` |
| `search-available-members` | `--org-id`, `--start`, `--max`, `--member-name`, `--phone-number`, `--extension`, `--order`, `--location-id`, `--exclude-virtual-line`, `--usage-type` |
| `get-service-password-status` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id` |
| `create-network` | `--location-id` *(required)*, `--org-id`, `--name`, `--model`, `--default-access-code-enabled`, `--default-access-code`, `--display-name`, `--body`, `--body-file` |
| `create-multiple-base-stations` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id`, `--base-station-macs`, `--body`, `--body-file` |
| `create-handset-network` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id`, `--line1-member-id`, `--custom-display-name`, `--line2-member-id`, `--body`, `--body-file` |
| `create-list-handsets-network` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `generate-enable-service-password` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id` |
| `update-network` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id`, `--name`, `--default-access-code-enabled`, `--default-access-code`, `--display-name`, `--body`, `--body-file` |
| `update-network-handset` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--handset-id` *(required)*, `--org-id`, `--line1-member-id`, `--custom-display-name`, `--line2-member-id`, `--body`, `--body-file` |
| `update-service-password-status` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `delete-network` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id` |
| `delete-bulk-network-base-stations` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id` |
| `delete-network-base-station` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--base-station-id` *(required)*, `--org-id` |
| `delete-network-handset` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--handset-id` *(required)*, `--org-id` |
| `delete-multiple-handsets` | `--location-id` *(required)*, `--dect-network-id` *(required)*, `--org-id`, `--handset-ids`, `--delete-all`, `--body`, `--body-file` |

### call-controls

| Command | Flags |
|---|---|
| `list` | `--line-owner-id` |
| `get-2` | `--call-id` *(required)*, `--line-owner-id` |
| `list-history` | `--type` |
| `list-member-id` | `--member-id` *(required)*, `--org-id` |
| `get-member-id` | `--member-id` *(required)*, `--call-id` *(required)*, `--org-id` |
| `dial` | `--destination`, `--endpoint-id`, `--single-number-reach-phone-number`, `--line-owner-id`, `--body`, `--body-file` |
| `answer` | `--call-id`, `--endpoint-id`, `--line-owner-id`, `--body`, `--body-file` |
| `reject` | `--call-id`, `--action`, `--line-owner-id`, `--body`, `--body-file` |
| `hangup` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `hold` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `resume` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `mute` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `unmute` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `divert` | `--call-id`, `--destination`, `--to-voicemail`, `--line-owner-id`, `--body`, `--body-file` |
| `transfer` | `--call-id1`, `--call-id2`, `--destination`, `--line-owner-id`, `--body`, `--body-file` |
| `park` | `--call-id`, `--destination`, `--is-group-park`, `--line-owner-id`, `--body`, `--body-file` |
| `get` | `--destination`, `--endpoint-id`, `--single-number-reach-phone-number`, `--line-owner-id`, `--body`, `--body-file` |
| `start-recording` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `stop-recording` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `pause-recording` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `resume-recording` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `transmit-dtmf` | `--call-id`, `--dtmf`, `--line-owner-id`, `--body`, `--body-file` |
| `push` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `pickup` | `--target`, `--endpoint-id`, `--single-number-reach-phone-number`, `--line-owner-id`, `--body`, `--body-file` |
| `barge` | `--target`, `--endpoint-id`, `--single-number-reach-phone-number`, `--line-owner-id`, `--body`, `--body-file` |
| `dial-member-id` | `--member-id` *(required)*, `--org-id`, `--destination`, `--endpoint-id`, `--single-number-reach-phone-number`, `--body`, `--body-file` |
| `answer-member-id` | `--member-id` *(required)*, `--org-id`, `--call-id`, `--endpoint-id`, `--body`, `--body-file` |
| `hangup-member-id` | `--member-id` *(required)*, `--org-id`, `--call-id`, `--body`, `--body-file` |
| `pull` | `--endpoint-id`, `--line-owner-id`, `--body`, `--body-file` |

### call-routing

| Command | Flags |
|---|---|
| `get-lgw-dial-plan-usage-trunk` | `--trunk-id` *(required)*, `--org-id`, `--start`, `--max`, `--order`, `--name` |
| `get-locations-lgw-pstn-connection` | `--trunk-id` *(required)*, `--org-id` |
| `get-route-groups-lgw` | `--trunk-id` *(required)*, `--org-id` |
| `get-lgw-usage-count` | `--trunk-id` *(required)*, `--org-id` |
| `list-dial-plans` | `--org-id`, `--dial-plan-name`, `--route-group-name`, `--trunk-name`, `--max`, `--start`, `--order` |
| `get-dial-plan` | `--dial-plan-id` *(required)*, `--org-id` |
| `list-trunks` | `--org-id`, `--name`, `--location-name`, `--trunk-type`, `--max`, `--start`, `--order` |
| `get-trunk` | `--trunk-id` *(required)*, `--org-id` |
| `list-trunk-types` | `--org-id` |
| `list-groups` | `--org-id`, `--name`, `--max`, `--start`, `--order` |
| `get-route-group` | `--route-group-id` *(required)*, `--org-id` |
| `get-usage-group` | `--route-group-id` *(required)*, `--org-id` |
| `get-extension-locations-group` | `--route-group-id` *(required)*, `--org-id`, `--location-name`, `--max`, `--start`, `--order` |
| `get-dial-plan-locations-group` | `--route-group-id` *(required)*, `--org-id`, `--location-name`, `--max`, `--start`, `--order` |
| `get-pstn-connection-locations-group` | `--route-group-id` *(required)*, `--org-id`, `--location-name`, `--max`, `--start`, `--order` |
| `get-route-lists-group` | `--route-group-id` *(required)*, `--org-id`, `--name`, `--max`, `--start`, `--order` |
| `list-route-lists` | `--org-id`, `--start`, `--max`, `--order`, `--name`, `--location-id` |
| `get-route-list` | `--route-list-id` *(required)*, `--org-id` |
| `get-numbers-assigned-route-list` | `--route-list-id` *(required)*, `--org-id`, `--start`, `--max`, `--number`, `--order` |
| `get-lgw-on-premises-extension-usage-trunk` | `--trunk-id` *(required)*, `--org-id`, `--start`, `--max`, `--order`, `--name` |
| `list-translation-patterns` | `--org-id`, `--limit-to-location-id`, `--limit-to-org-level-enabled`, `--max`, `--start`, `--order`, `--name`, `--matching-pattern` |
| `get-translation-pattern` | `--translation-id` *(required)*, `--org-id` |
| `get-translation-pattern-location` | `--location-id` *(required)*, `--translation-id` *(required)*, `--org-id` |
| `test` | `--org-id`, `--originator-id`, `--originator-type`, `--destination`, `--originator-number`, `--include-applied-services`, `--body`, `--body-file` |
| `validate-dial-pattern` | `--org-id`, `--dial-patterns`, `--body`, `--body-file` |
| `create-dial-plan` | `--org-id`, `--name`, `--route-id`, `--route-type`, `--dial-patterns`, `--body`, `--body-file` |
| `validate-lgw-fqdn-domain-trunk` | `--org-id`, `--address`, `--domain`, `--port`, `--body`, `--body-file` |
| `create-trunk` | `--org-id`, `--name`, `--location-id`, `--password`, `--trunk-type`, `--dual-identity-support-enabled`, `--device-type`, `--address`, `--domain`, `--port`, `--max-concurrent-calls`, `--p-charge-info-support-policy`, `--body`, `--body-file` |
| `create-route-group` | `--org-id`, `--body`, `--body-file` |
| `create-route-list` | `--org-id`, `--name`, `--location-id`, `--route-group-id`, `--body`, `--body-file` |
| `create-translation-pattern` | `--org-id`, `--name`, `--matching-pattern`, `--replacement-pattern`, `--body`, `--body-file` |
| `create-translation-pattern-location` | `--location-id` *(required)*, `--org-id`, `--name`, `--matching-pattern`, `--replacement-pattern`, `--body`, `--body-file` |
| `update-dial-patterns` | `--dial-plan-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-dial-plan` | `--dial-plan-id` *(required)*, `--org-id`, `--name`, `--route-id`, `--route-type`, `--body`, `--body-file` |
| `update-trunk` | `--trunk-id` *(required)*, `--org-id`, `--name`, `--password`, `--dual-identity-support-enabled`, `--max-concurrent-calls`, `--p-charge-info-support-policy`, `--body`, `--body-file` |
| `update-route-group` | `--route-group-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-route-list` | `--route-list-id` *(required)*, `--org-id`, `--name`, `--route-group-id`, `--body`, `--body-file` |
| `update-numbers-route-list` | `--route-list-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-translation-pattern` | `--translation-id` *(required)*, `--org-id`, `--name`, `--matching-pattern`, `--replacement-pattern`, `--body`, `--body-file` |
| `update-translation-pattern-location` | `--location-id` *(required)*, `--translation-id` *(required)*, `--org-id`, `--name`, `--matching-pattern`, `--replacement-pattern`, `--body`, `--body-file` |
| `delete-dial-plan` | `--dial-plan-id` *(required)*, `--org-id` |
| `delete-trunk` | `--trunk-id` *(required)*, `--org-id` |
| `delete-route-group-org` | `--route-group-id` *(required)*, `--org-id` |
| `delete-route-list` | `--route-list-id` *(required)*, `--org-id` |
| `delete-translation-pattern` | `--translation-id` *(required)*, `--org-id` |
| `delete-translation-pattern-location` | `--location-id` *(required)*, `--translation-id` *(required)*, `--org-id` |

### location-call

| Command | Flags |
|---|---|
| `list-dial-patterns` | `--dial-plan-id` *(required)*, `--org-id`, `--dial-pattern`, `--max`, `--start`, `--order` |
| `list-webex-calling` | `--org-id`, `--max`, `--start`, `--name`, `--order` |
| `get-webex-calling` | `--location-id` *(required)*, `--org-id` |
| `list-update-routing-prefix-jobs` | `--org-id` |
| `get-job-status-update-routing-prefix-job` | `--job-id` *(required)*, `--org-id` |
| `get-job-errors-update-routing-prefix-job` | `--job-id` *(required)*, `--org-id` |
| `get-emergency-callback` | `--location-id` *(required)*, `--org-id` |
| `get-music-hold` | `--location-id` *(required)*, `--org-id` |
| `get-private-network-connect` | `--location-id` *(required)*, `--org-id` |
| `list-routing-choices` | `--org-id`, `--route-group-name`, `--trunk-name`, `--max`, `--start`, `--order` |
| `list-phone-numbers-available-external-caller-id` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--person-id` |
| `get-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name` |
| `get-webex-go-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-ecbn-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name` |
| `get-intercept-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `list-receptionist-contact-directories` | `--location-id` *(required)*, `--org-id` |
| `get-receptionist-contact-directory` | `--location-id` *(required)*, `--directory-id` *(required)*, `--org-id`, `--search-criteria-mode-or`, `--first-name`, `--last-name`, `--phone-number`, `--extension`, `--person-id` |
| `get-available-charge-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name` |
| `list-disable-calling-jobs` | `--org-id`, `--max`, `--start` |
| `get-errors-disable-calling-job` | `--job-id` *(required)*, `--org-id` |
| `get-disable-calling-job-status` | `--job-id` *(required)*, `--org-id` |
| `get-captions-settings` | `--location-id` *(required)*, `--org-id` |
| `enable-webex-calling` | `--org-id`, `--body`, `--body-file` |
| `update-announcement-language` | `--location-id` *(required)*, `--org-id`, `--announcement-language-code`, `--agent-enabled`, `--service-enabled`, `--body`, `--body-file` |
| `validate-list-extensions` | `--org-id`, `--extensions`, `--body`, `--body-file` |
| `validate-extensions` | `--location-id` *(required)*, `--org-id`, `--extensions`, `--body`, `--body-file` |
| `create-receptionist-contact-directory` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `disable-webex-calling` | `--org-id`, `--location-id`, `--location-name`, `--force-delete`, `--body`, `--body-file` |
| `check-delete-check-disabling-webex-calling` | `--location-id` *(required)*, `--org-id` |
| `pause-disable-calling-job` | `--job-id` *(required)*, `--org-id` |
| `resume-paused-disable-calling-job` | `--job-id` *(required)*, `--org-id` |
| `update-webex-calling` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-emergency-callback` | `--location-id` *(required)*, `--org-id`, `--selected`, `--location-member-id`, `--elin-expiry-time-minutes`, `--body`, `--body-file` |
| `update-music-hold` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-private-network-connect` | `--location-id` *(required)*, `--org-id`, `--network-connection-type`, `--body`, `--body-file` |
| `update-receptionist-contact-directory` | `--location-id` *(required)*, `--directory-id` *(required)*, `--org-id`, `--name`, `--contacts`, `--body`, `--body-file` |
| `update-captions-settings` | `--location-id` *(required)*, `--org-id`, `--location-closed-captions-enabled`, `--location-transcripts-enabled`, `--use-org-settings-enabled`, `--body`, `--body-file` |
| `delete-receptionist-contact-directory` | `--location-id` *(required)*, `--directory-id` *(required)*, `--org-id` |

### call-settings-for-me

| Command | Flags |
|---|---|
| `get-my-personal-assistant` | — |
| `list-available-preferred-endpoints` | — |
| `get-my-secondary-owner-available-preferred-endpoint-list` | `--line-owner-id` *(required)* |
| `get-preferred-endpoint` | — |
| `get-my-secondary-owner-preferred-endpoint` | `--line-owner-id` *(required)* |
| `get-my-webexgooverride` | — |
| `get-my-caller-id` | — |
| `get-my-secondary-owner-caller-id` | `--lineowner-id` *(required)* |
| `get-my-selected-caller-id` | — |
| `get-my-secondary-owner-selected-caller-id` | `--lineowner-id` *(required)* |
| `get-my-available-caller-id-list` | — |
| `get-my-secondary-owner-available-caller-id-list` | `--lineowner-id` *(required)* |
| `list-my-endpoints` | — |
| `get-my-endpoints` | `--endpoint-id` *(required)* |
| `get-my-recording` | — |
| `get-my-secondary-owner-recording` | `--lineowner-id` *(required)* |
| `get-my-own` | — |
| `get-my-access-codes` | — |
| `get-my-access-codes-secondary-owner` | `--lineowner-id` *(required)* |
| `get-my-executive-assigned-assistants` | — |
| `get-my-executive-available-assistants` | — |
| `get-my-executive-assistant` | — |
| `get-my-calling-services` | — |
| `get-my-secondary-owner-calling-services` | `--lineowner-id` *(required)* |
| `get-user-single-number-reach` | — |
| `get-my-forward` | — |
| `get-my-secondary-owner-forward` | `--lineowner-id` *(required)* |
| `get-my-pickup-group` | — |
| `get-my-secondary-owner-pickup-group` | `--lineowner-id` *(required)* |
| `get-my-park` | — |
| `get-my-secondary-owner-park` | `--lineowner-id` *(required)* |
| `get-voicemail-person` | — |
| `get-my-secondary-owner-voicemail` | `--lineowner-id` *(required)* |
| `get-my-block` | — |
| `get-my-block-state-number` | `--phone-number-id` *(required)* |
| `get-my-monitoring` | — |
| `get-my-center` | — |
| `get-my-secondary-owner-center` | `--lineowner-id` *(required)* |
| `get-my-captions` | — |
| `get-policies-user` | — |
| `get-user-executive-screening` | — |
| `get-user-executive-filtering` | — |
| `get-user-executive-filtering-2` | `--id` *(required)* |
| `get-do-not-disturb-user` | — |
| `get-user-executive-alert` | — |
| `get-country-specific-telephony-configuration-requirements` | `--country-code` *(required)* |
| `get-announcement-languages-authenticated-user` | — |
| `get-barge-in` | — |
| `get-priority-alert` | — |
| `get-priority-alert-2` | `--id` *(required)* |
| `get-user-schedules` | — |
| `get-user-schedule` | `--schedule-type` *(required)*, `--schedule-id` *(required)* |
| `get-user-schedule-event` | `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)* |
| `get-user-location-schedule` | `--schedule-type` *(required)*, `--schedule-id` *(required)* |
| `get-notify-user` | — |
| `get-notify` | `--id` *(required)* |
| `get-selective-accept-user` | — |
| `get-selective-accept-user-2` | `--id` *(required)* |
| `get-available-numbers-user-location` | `--max`, `--start`, `--name`, `--phone-number`, `--extension`, `--order` |
| `get-selective-forward` | `--id` *(required)* |
| `get-selective-forward-user` | — |
| `get-selective-reject-user` | — |
| `get-anonymous-rejection-user` | — |
| `get-selective-reject-user-2` | `--id` *(required)* |
| `get-waiting-user` | — |
| `get-sequential-ring-user` | — |
| `get-sequential-ring-user-2` | `--id` *(required)* |
| `get-priority-alert-2` | — |
| `get-priority-alert-2-2` | `--id` *(required)* |
| `get-user-schedules-2` | — |
| `get-user-schedule-2` | `--schedule-type` *(required)*, `--schedule-id` *(required)* |
| `get-user-schedule-event-2` | `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)* |
| `get-user-location-schedule-2` | `--schedule-type` *(required)*, `--schedule-id` *(required)* |
| `get-notify-user-2` | — |
| `get-notify-2` | `--id` *(required)* |
| `get-selective-accept-user-2` | — |
| `get-selective-accept-user-2-2` | `--id` *(required)* |
| `get-available-numbers-user-location-2` | `--max`, `--start`, `--name`, `--phone-number`, `--extension`, `--order` |
| `get-selective-forward-2` | `--id` *(required)* |
| `get-selective-forward-user-2` | — |
| `get-selective-reject-user-2` | — |
| `get-anonymous-rejection-user-2` | — |
| `get-selective-reject-user-2-2` | `--id` *(required)* |
| `get-waiting-user-2` | — |
| `get-sequential-ring-user-2` | — |
| `get-sequential-ring-user-2-2` | `--id` *(required)* |
| `get-my-simultaneous-ring` | — |
| `get-my-simultaneous-ring-2` | `--id` *(required)* |
| `get-my-guest-calling-numbers` | — |
| `get-personal-assistant` | — |
| `get-person-voicemail-rules` | — |
| `get-hoteling-guest` | — |
| `get-available-hoteling-hosts` | `--max`, `--start`, `--name`, `--phone-number` |
| `create-phone-number-user-single-number-reach` | `--phone-number`, `--name`, `--enabled`, `--do-not-forward-calls-enabled`, `--answer-confirmation-enabled`, `--body`, `--body-file` |
| `create-phone-number-user-block-list` | `--phone-number`, `--body`, `--body-file` |
| `create-user-executive-filtering` | `--body`, `--body-file` |
| `create-priority-alert` | `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--notification-enabled`, `--body`, `--body-file` |
| `create-user-schedule` | `--body`, `--body-file` |
| `create-event-user-schedule` | `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--body`, `--body-file` |
| `create-notify` | `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--notification-enabled`, `--body`, `--body-file` |
| `create-user-selective-accept` | `--calls-from`, `--accept-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-selective-forward` | `--forward-to-phone-number`, `--destination-voicemail-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--forward-enabled`, `--body`, `--body-file` |
| `create-user-selective-reject` | `--calls-from`, `--reject-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-user-sequential-ring` | `--calls-from`, `--ring-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-priority-alert-2` | `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--notification-enabled`, `--body`, `--body-file` |
| `create-user-schedule-2` | `--body`, `--body-file` |
| `create-event-user-schedule-2` | `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--body`, `--body-file` |
| `create-notify-2` | `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--notification-enabled`, `--body`, `--body-file` |
| `create-user-selective-accept-2` | `--calls-from`, `--accept-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-selective-forward-2` | `--forward-to-phone-number`, `--destination-voicemail-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--forward-enabled`, `--body`, `--body-file` |
| `create-user-selective-reject-2` | `--calls-from`, `--reject-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-user-sequential-ring-2` | `--calls-from`, `--ring-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-my-simultaneous-ring` | `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--ring-enabled`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `update-my-personal-assistant` | `--enabled`, `--presence`, `--until-date-time`, `--transfer-enabled`, `--transfer-number`, `--alerting`, `--alert-me-first-number-of-rings`, `--body`, `--body-file` |
| `update-preferred-endpoint` | `--preferred-answer-endpoint-id`, `--body`, `--body-file` |
| `update-my-secondary-owner-preferred-endpoint` | `--line-owner-id` *(required)*, `--preferred-answer-endpoint-id`, `--body`, `--body-file` |
| `update-my-webexgooverride` | `--enabled`, `--body`, `--body-file` |
| `update-my-caller-id` | `--calling-line-id-delivery-blocking-enabled`, `--connected-line-identification-restriction-enabled`, `--body`, `--body-file` |
| `update-my-secondary-owner-caller-id` | `--lineowner-id` *(required)*, `--calling-line-id-delivery-blocking-enabled`, `--connected-line-identification-restriction-enabled`, `--body`, `--body-file` |
| `update-my-selected-caller-id` | `--body`, `--body-file` |
| `update-my-secondary-owner-selected-caller-id` | `--lineowner-id` *(required)*, `--body`, `--body-file` |
| `update-my-endpoints` | `--endpoint-id` *(required)*, `--body`, `--body-file` |
| `update-my-executive-assigned-assistants` | `--allow-opt-in-out-enabled`, `--assistant-ids`, `--body`, `--body-file` |
| `update-my-executive-assistant` | `--body`, `--body-file` |
| `update-user-single-number-reach` | `--alert-all-locations-for-click-to-dial-calls-enabled`, `--body`, `--body-file` |
| `update-user-single-number-reach-contact` | `--phone-number-id` *(required)*, `--phone-number`, `--name`, `--enabled`, `--do-not-forward-calls-enabled`, `--answer-confirmation-enabled`, `--body`, `--body-file` |
| `update-my-forward` | `--body`, `--body-file` |
| `update-my-secondary-owner-forward` | `--lineowner-id` *(required)*, `--body`, `--body-file` |
| `update-voicemail-person` | `--body`, `--body-file` |
| `update-my-secondary-owner-voicemail` | `--lineowner-id` *(required)*, `--body`, `--body-file` |
| `update-my-center` | `--body`, `--body-file` |
| `update-my-secondary-owner-center` | `--lineowner-id` *(required)*, `--body`, `--body-file` |
| `update-policies-user` | `--connected-line-id-privacy-on-redirected-calls`, `--body`, `--body-file` |
| `update-user-executive-screening` | `--enabled`, `--alert-type`, `--alert-anywhere-location-enabled`, `--alert-mobility-location-enabled`, `--alert-shared-call-appearance-location-enabled`, `--body`, `--body-file` |
| `update-user-executive-filtering` | `--body`, `--body-file` |
| `update-user-executive-filtering-2` | `--id` *(required)*, `--body`, `--body-file` |
| `update-do-not-disturb-user` | `--webex-go-override-enabled`, `--enabled`, `--ring-splash-enabled`, `--body`, `--body-file` |
| `update-user-executive-alert` | `--alerting-mode`, `--next-assistant-number-of-rings`, `--rollover-enabled`, `--rollover-action`, `--rollover-forward-to-phone-number`, `--rollover-wait-time-in-secs`, `--clid-name-mode`, `--custom-clidname`, `--custom-clidname-in-unicode`, `--clid-phone-number-mode`, `--custom-clidphone-number`, `--body`, `--body-file` |
| `update-barge-in` | `--enabled`, `--tone-enabled`, `--body`, `--body-file` |
| `update-priority-alert-user` | `--enabled`, `--body`, `--body-file` |
| `update-priority-alert` | `--id` *(required)*, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--notification-enabled`, `--body`, `--body-file` |
| `update-user-schedule` | `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--body`, `--body-file` |
| `update-user-schedule-event` | `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)*, `--body`, `--body-file` |
| `update-notify-user` | `--enabled`, `--email-address`, `--body`, `--body-file` |
| `update-notify` | `--id` *(required)*, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--notification-enabled`, `--body`, `--body-file` |
| `update-selective-accept-user` | `--enabled`, `--body`, `--body-file` |
| `update-selective-accept` | `--id` *(required)*, `--calls-from`, `--accept-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `update-selective-forward` | `--id` *(required)*, `--forward-to-phone-number`, `--destination-voicemail-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--forward-enabled`, `--body`, `--body-file` |
| `update-selective-forward-user` | `--enabled`, `--default-phone-number-to-forward`, `--ring-reminder-enabled`, `--destination-voicemail-enabled`, `--body`, `--body-file` |
| `update-selective-reject-user` | `--enabled`, `--body`, `--body-file` |
| `update-anonymous-rejection-user` | `--enabled`, `--body`, `--body-file` |
| `update-selective-reject` | `--id` *(required)*, `--calls-from`, `--reject-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `update-waiting-user` | `--enabled`, `--body`, `--body-file` |
| `update-sequential-ring-user` | `--body`, `--body-file` |
| `update-sequential-ring-user-2` | `--id` *(required)*, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--ring-enabled`, `--body`, `--body-file` |
| `update-priority-alert-user-2` | `--enabled`, `--body`, `--body-file` |
| `update-priority-alert-2` | `--id` *(required)*, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--notification-enabled`, `--body`, `--body-file` |
| `update-user-schedule-2` | `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--body`, `--body-file` |
| `update-user-schedule-event-2` | `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)*, `--body`, `--body-file` |
| `update-notify-user-2` | `--enabled`, `--email-address`, `--body`, `--body-file` |
| `update-notify-2` | `--id` *(required)*, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--notification-enabled`, `--body`, `--body-file` |
| `update-selective-accept-user-2` | `--enabled`, `--body`, `--body-file` |
| `update-selective-accept-2` | `--id` *(required)*, `--calls-from`, `--accept-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `update-selective-forward-2` | `--id` *(required)*, `--forward-to-phone-number`, `--destination-voicemail-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--forward-enabled`, `--body`, `--body-file` |
| `update-selective-forward-user-2` | `--enabled`, `--default-phone-number-to-forward`, `--ring-reminder-enabled`, `--destination-voicemail-enabled`, `--body`, `--body-file` |
| `update-selective-reject-user-2` | `--enabled`, `--body`, `--body-file` |
| `update-anonymous-rejection-user-2` | `--enabled`, `--body`, `--body-file` |
| `update-selective-reject-2` | `--id` *(required)*, `--calls-from`, `--reject-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `update-waiting-user-2` | `--enabled`, `--body`, `--body-file` |
| `update-sequential-ring-user-2` | `--body`, `--body-file` |
| `update-sequential-ring-user-2-2` | `--id` *(required)*, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--ring-enabled`, `--body`, `--body-file` |
| `update-my-simultaneous-ring` | `--body`, `--body-file` |
| `update-my-simultaneous-ring-2` | `--id` *(required)*, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--ring-enabled`, `--body`, `--body-file` |
| `update-personal-assistant` | `--enabled`, `--presence`, `--until-date-time`, `--transfer-enabled`, `--transfer-number`, `--alerting`, `--alert-me-first-number-of-rings`, `--body`, `--body-file` |
| `update-voicemail-pin` | `--passcode`, `--body`, `--body-file` |
| `update-hoteling-guest` | `--enabled`, `--association-limit-enabled`, `--association-limit-hours`, `--host-id`, `--body`, `--body-file` |
| `delete-user-single-number-reach-contact` | `--phone-number-id` *(required)* |
| `delete-user-block-number` | `--phone-number-id` *(required)* |
| `delete-user-executive-filtering` | `--id` *(required)* |
| `delete-priority-alert` | `--id` *(required)* |
| `delete-user-schedule` | `--schedule-type` *(required)*, `--schedule-id` *(required)* |
| `delete-user-schedule-event` | `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)* |
| `delete-notify` | `--id` *(required)* |
| `delete-selective-accept` | `--id` *(required)* |
| `delete-selective-forward` | `--id` *(required)* |
| `delete-selective-reject` | `--id` *(required)* |
| `delete-sequential-ring` | `--id` *(required)* |
| `delete-priority-alert-2` | `--id` *(required)* |
| `delete-user-schedule-2` | `--schedule-type` *(required)*, `--schedule-id` *(required)* |
| `delete-user-schedule-event-2` | `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)* |
| `delete-notify-2` | `--id` *(required)* |
| `delete-selective-accept-2` | `--id` *(required)* |
| `delete-selective-forward-2` | `--id` *(required)* |
| `delete-selective-reject-2` | `--id` *(required)* |
| `delete-sequential-ring-2` | `--id` *(required)* |
| `delete-my-simultaneous-ring` | `--id` *(required)* |

### calling-service

| Command | Flags |
|---|---|
| `list-announcement-languages` | `--tts-language` |
| `get-voicemail` | `--org-id` |
| `get-voicemail-rules` | `--org-id` |
| `get-org-music-hold-configuration` | `--org-id` |
| `get-org-call-captions-settings` | `--org-id` |
| `get-large-org-status` | `--org-id` |
| `update-voicemail` | `--org-id`, `--message-expiry-enabled`, `--number-of-days-for-message-expiry`, `--strict-deletion-enabled`, `--voice-message-forwarding-enabled`, `--body`, `--body-file` |
| `update-voicemail-rules` | `--org-id`, `--body`, `--body-file` |
| `update-org-music-hold-configuration` | `--org-id`, `--default-org-moh`, `--body`, `--body-file` |
| `update-org-call-captions-settings` | `--org-id`, `--org-closed-captions-enabled`, `--org-transcripts-enabled`, `--body`, `--body-file` |

### client-call

| Command | Flags |
|---|---|
| `get-org-ms-teams-settings` | `--org-id` |
| `update-org-ms-teams-setting` | `--org-id`, `--setting-name`, `--value`, `--body`, `--body-file` |

### conference-controls

| Command | Flags |
|---|---|
| `get` | `--line-owner-id` |
| `start` | `--call-ids`, `--line-owner-id`, `--body`, `--body-file` |
| `create-participant` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `mute` | `--call-id`, `--body`, `--body-file` |
| `unmute` | `--call-id`, `--body`, `--body-file` |
| `deafen-participant` | `--call-id`, `--body`, `--body-file` |
| `undeafen-participant` | `--call-id`, `--body`, `--body-file` |
| `hold` | `--line-owner-id` |
| `resume` | `--line-owner-id` |
| `release` | `--line-owner-id` |

### converged-recordings

| Command | Flags |
|---|---|
| `list` | `--max`, `--from`, `--to`, `--status`, `--service-type`, `--format`, `--owner-type`, `--storage-region`, `--location-id`, `--topic`, `--last` |
| `list-admin-compliance-officer` | `--max`, `--from`, `--to`, `--status`, `--service-type`, `--format`, `--owner-id`, `--owner-email`, `--owner-type`, `--storage-region`, `--location-id`, `--topic`, `--last` |
| `get` | `--recording-id` *(required)* |
| `get-metadata` | `--recording-id` *(required)*, `--show-all-types` |
| `reassign` | `--reassign-owner-email`, `--owner-email`, `--owner-id`, `--recording-ids`, `--body`, `--body-file` |
| `move-recycle-bin` | `--trash-all`, `--owner-email`, `--recording-ids`, `--body`, `--body-file` |
| `restore-recycle-bin` | `--restore-all`, `--owner-email`, `--recording-ids`, `--body`, `--body-file` |
| `purge-recycle-bin` | `--purge-all`, `--owner-email`, `--recording-ids`, `--body`, `--body-file` |
| `delete` | `--recording-id` *(required)*, `--reason`, `--comment`, `--body`, `--body-file` |

### device-call

| Command | Flags |
|---|---|
| `get-members` | `--device-id` *(required)*, `--org-id` |
| `search-members` | `--device-id` *(required)*, `--org-id`, `--start`, `--max`, `--member-name`, `--phone-number`, `--location-id`, `--extension`, `--usage-type`, `--order` |
| `get-settings` | `--device-id` *(required)*, `--org-id`, `--device-model` |
| `get-location-settings` | `--location-id` *(required)*, `--org-id` |
| `get-webex-calling` | `--device-id` *(required)*, `--org-id` |
| `get-person` | `--person-id` *(required)*, `--org-id` |
| `get-workspace` | `--workspace-id` *(required)*, `--org-id` |
| `list-supported` | `--org-id`, `--allow-configure-layout-enabled`, `--type` |
| `get-override-settings` | `--org-id` |
| `list-line-key-templates` | `--org-id` |
| `get-line-key-template` | `--template-id` *(required)*, `--org-id` |
| `list-apply-line-key-template-jobs` | `--org-id` |
| `get-job-status-apply-line-key-template-job` | `--job-id` *(required)*, `--org-id` |
| `get-job-errors-apply-line-key-template-job` | `--job-id` *(required)*, `--org-id` |
| `get-dect-type-list-deprecated` | `--org-id` |
| `get-dect-type-list` | `--org-id` |
| `list-change-settings-jobs` | `--org-id`, `--start`, `--max` |
| `get-change-settings-job-status` | `--job-id` *(required)* |
| `list-change-settings-job-errors` | `--job-id` *(required)*, `--org-id`, `--start`, `--max` |
| `get-layout-id` | `--device-id` *(required)*, `--org-id` |
| `list-rebuild-phones-jobs` | `--org-id` |
| `get-job-status-rebuild-phones-job` | `--job-id` *(required)*, `--org-id` |
| `get-job-errors-rebuild-phones-job` | `--job-id` *(required)*, `--org-id` |
| `get-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-settings-workspace` | `--workspace-id` *(required)*, `--org-id` |
| `list-background-images` | `--org-id` |
| `get-user-count` | `--person-id` *(required)*, `--org-id` |
| `get-count-members` | `--device-id` *(required)*, `--org-id`, `--member-name`, `--phone-number`, `--location-id`, `--extension`, `--usage-type` |
| `get-count-available-members` | `--org-id`, `--member-name`, `--phone-number`, `--location-id`, `--extension`, `--usage-type`, `--exclude-virtual-line`, `--device-location-id` |
| `list-supported-2` | `--org-id`, `--allow-configure-layout-enabled`, `--type` |
| `get-validation-schema` | `--org-id`, `--family-or-model-display-name` |
| `get-settings-groups` | `--org-id`, `--family-or-model-display-name`, `--include-settings-type` |
| `list-dynamic-settings-jobs` | `--org-id`, `--start`, `--max` |
| `get-dynamic-settings-job-status` | `--job-id` *(required)* |
| `list-dynamic-settings-job-errors` | `--job-id` *(required)*, `--org-id`, `--start`, `--max` |
| `list-supported-3` | `--org-id`, `--allow-configure-layout-enabled`, `--type` |
| `get-settings-groups-2` | `--org-id`, `--family-or-model-display-name`, `--include-settings-type` |
| `list-dynamic-settings-jobs-2` | `--org-id`, `--start`, `--max` |
| `get-dynamic-settings-job-status-2` | `--job-id` *(required)* |
| `list-dynamic-settings-job-errors-2` | `--job-id` *(required)*, `--org-id`, `--start`, `--max` |
| `apply-changes` | `--device-id` *(required)*, `--org-id` |
| `create-line-key-template` | `--org-id`, `--body`, `--body-file` |
| `preview-apply-line-key-template` | `--org-id`, `--body`, `--body-file` |
| `apply-line-key-template` | `--org-id`, `--body`, `--body-file` |
| `validate-list-mac-address` | `--org-id`, `--macs`, `--body`, `--body-file` |
| `update-settings-across-org-location-job` | `--org-id`, `--body`, `--body-file` |
| `rebuild-phones-configuration` | `--org-id`, `--location-id`, `--body`, `--body-file` |
| `upload-background-image` | `--device-id` *(required)*, `--org-id` |
| `get-customer-dynamic-settings` | `--org-id`, `--family-or-model-display-name`, `--tags`, `--body`, `--body-file` |
| `get-location-dynamic-settings` | `--location-id` *(required)*, `--org-id`, `--family-or-model-display-name`, `--tags`, `--body`, `--body-file` |
| `get-dynamic-settings` | `--device-id` *(required)*, `--org-id`, `--tags`, `--body`, `--body-file` |
| `updates-dynamic-settings-across-org-location` | `--org-id`, `--body`, `--body-file` |
| `get-customer-dynamic-settings-2` | `--org-id`, `--family-or-model-display-name`, `--tags`, `--body`, `--body-file` |
| `get-location-dynamic-settings-2` | `--location-id` *(required)*, `--org-id`, `--family-or-model-display-name`, `--tags`, `--body`, `--body-file` |
| `get-dynamic-settings-2` | `--device-id` *(required)*, `--org-id`, `--tags`, `--body`, `--body-file` |
| `update-dynamic-settings-across-org-location` | `--org-id`, `--body`, `--body-file` |
| `update-members` | `--device-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-settings` | `--device-id` *(required)*, `--org-id`, `--device-model`, `--body`, `--body-file` |
| `update-third-party` | `--device-id` *(required)*, `--org-id`, `--sip-password`, `--body`, `--body-file` |
| `update-hoteling-settings-person-primary` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-workspace` | `--workspace-id` *(required)*, `--org-id`, `--enabled`, `--limit-guest-use`, `--guest-hours-limit`, `--body`, `--body-file` |
| `update-line-key-template` | `--template-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-layout-id` | `--device-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-settings-person` | `--person-id` *(required)*, `--org-id`, `--compression`, `--body`, `--body-file` |
| `update-settings-workspace` | `--workspace-id` *(required)*, `--org-id`, `--compression`, `--body`, `--body-file` |
| `update-specified-settings` | `--device-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-dynamic-settings` | `--device-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `delete-line-key-template` | `--template-id` *(required)*, `--org-id` |
| `delete-background-images` | `--org-id`, `--body`, `--body-file` |

### devices

| Command | Flags |
|---|---|
| `list` | `--max`, `--start`, `--display-name`, `--person-id`, `--workspace-id`, `--org-id`, `--connection-status`, `--product`, `--type`, `--serial`, `--tag`, `--software`, `--upgrade-channel`, `--error-code`, `--capability`, `--permission`, `--location-id`, `--workspace-location-id`, `--mac`, `--device-platform`, `--planned-maintenance` |
| `get` | `--device-id` *(required)*, `--org-id` |
| `create-mac-address` | `--org-id`, `--mac`, `--model`, `--workspace-id`, `--person-id`, `--password`, `--body`, `--body-file` |
| `create-activation-code` | `--org-id`, `--workspace-id`, `--person-id`, `--model`, `--body`, `--body-file` |
| `update-tags` | `--device-id` *(required)*, `--org-id`, `--op`, `--path`, `--value`, `--body`, `--body-file` |
| `delete` | `--device-id` *(required)*, `--org-id` |

### emergency-services

| Command | Flags |
|---|---|
| `get-redsky-account` | `--org-id` |
| `get-org-compliance-status-redsky-account` | `--org-id` |
| `get-org-compliance-status-location-status-list` | `--org-id`, `--start`, `--max`, `--order` |
| `get-location-redsky-calling-parameters` | `--location-id` *(required)*, `--org-id` |
| `get-location-redsky-compliance-status` | `--location-id` *(required)*, `--org-id` |
| `get-org-call-notification` | `--org-id` |
| `get-location-call-notification` | `--location-id` *(required)*, `--org-id` |
| `get-dependencies-hunt-group-callback` | `--hunt-group-id` *(required)*, `--org-id` |
| `get-person-callback` | `--person-id` *(required)*, `--org-id` |
| `get-person-callback-dependencies` | `--person-id` *(required)*, `--org-id` |
| `get-workspace-callback` | `--workspace-id` *(required)*, `--org-id` |
| `get-workspace-callback-dependencies` | `--workspace-id` *(required)*, `--org-id` |
| `get-dependencies-vline-callback` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-vline-callback-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `create-account-admin-redsky` | `--org-id`, `--email`, `--org-prefix`, `--partner-redsky-org-id`, `--body`, `--body-file` |
| `login-redsky-admin-account` | `--org-id`, `--email`, `--password`, `--red-sky-org-id`, `--body`, `--body-file` |
| `create-redsky-building-address-location` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-redsky-settings` | `--org-id`, `--enabled`, `--company-id`, `--secret`, `--external-tenant-enabled`, `--email`, `--password`, `--body`, `--body-file` |
| `update-org-redsky-account-compliance-status` | `--org-id`, `--compliance-status`, `--body`, `--body-file` |
| `update-location-redsky-compliance-status` | `--location-id` *(required)*, `--org-id`, `--compliance-status`, `--body`, `--body-file` |
| `update-redsky-building-address-location` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-org-call-notification` | `--org-id`, `--emergency-call-notification-enabled`, `--allow-email-notification-all-location-enabled`, `--email-address`, `--body`, `--body-file` |
| `update-location-call-notification` | `--location-id` *(required)*, `--org-id`, `--emergency-call-notification-enabled`, `--email-address`, `--body`, `--body-file` |
| `update-person-callback` | `--person-id` *(required)*, `--org-id`, `--selected`, `--location-member-id`, `--elin-enabled`, `--elin-for-webex-app-enabled`, `--body`, `--body-file` |
| `update-workspace-callback` | `--workspace-id` *(required)*, `--org-id`, `--selected`, `--location-member-id`, `--elin-enabled`, `--body`, `--body-file` |
| `update-vline-callback-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--selected`, `--location-member-id`, `--body`, `--body-file` |

### announcement-playlist

| Command | Flags |
|---|---|
| `list` | `--org-id` |
| `get` | `--playlist-id` *(required)*, `--org-id` |
| `list-locations` | `--playlist-id` *(required)*, `--org-id` |
| `create` | `--org-id`, `--name`, `--announcement-ids`, `--body`, `--body-file` |
| `update` | `--playlist-id` *(required)*, `--org-id`, `--name`, `--announcement-ids`, `--body`, `--body-file` |
| `update-locations` | `--playlist-id` *(required)*, `--org-id`, `--location-ids`, `--body`, `--body-file` |
| `delete` | `--playlist-id` *(required)*, `--org-id` |

### announcement-repository

| Command | Flags |
|---|---|
| `list-greetings` | `--org-id`, `--location-id`, `--max`, `--start`, `--order`, `--file-name`, `--file-type`, `--media-file-type`, `--name` |
| `get-usage` | `--org-id` |
| `get-binary-greeting` | `--announcement-id` *(required)*, `--org-id` |
| `get-usage-location` | `--location-id` *(required)*, `--org-id` |
| `get-binary-greeting-2` | `--location-id` *(required)*, `--announcement-id` *(required)*, `--org-id` |
| `get-text-to-speech-usage` | `--org-id` |
| `get-text-to-speech-generation-status` | `--tts-id` *(required)*, `--org-id` |
| `list-text-to-speech-voices` | `--org-id` |
| `generate-text-to-speech-prompt` | `--org-id`, `--voice`, `--text`, `--language-code`, `--body`, `--body-file` |
| `delete-greeting-org` | `--announcement-id` *(required)*, `--org-id` |
| `delete-greeting-location` | `--location-id` *(required)*, `--announcement-id` *(required)*, `--org-id` |

### auto-attendant

| Command | Flags |
|---|---|
| `list` | `--org-id`, `--location-id`, `--max`, `--start`, `--name`, `--phone-number` |
| `get` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--org-id` |
| `get-call-forward` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--org-id` |
| `get-selective-forward-rule` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--rule-id` *(required)*, `--org-id` |
| `get-primary-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-alternate-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-forward-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `list-announcement-files` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--org-id` |
| `create` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `create-selective-forward-rule` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `switch-mode-call-forward` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-call-forward` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-selective-forward-rule` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--rule-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `delete` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--org-id` |
| `delete-selective-forward-rule` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--rule-id` *(required)*, `--org-id` |
| `delete-announcement-file` | `--location-id` *(required)*, `--auto-attendant-id` *(required)*, `--file-name` *(required)*, `--org-id` |

### call-park

| Command | Flags |
|---|---|
| `list` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--order`, `--name` |
| `get` | `--location-id` *(required)*, `--call-park-id` *(required)*, `--org-id` |
| `get-available-agents` | `--location-id` *(required)*, `--org-id`, `--call-park-name`, `--max`, `--start`, `--name`, `--phone-number`, `--order` |
| `get-available-recall-hunt-groups` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--name`, `--order` |
| `get-2` | `--location-id` *(required)*, `--org-id` |
| `list-extensions` | `--org-id`, `--location-id`, `--max`, `--start`, `--extension`, `--location-name`, `--name`, `--order` |
| `get-extension` | `--location-id` *(required)*, `--call-park-extension-id` *(required)*, `--org-id` |
| `create` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `create-extension` | `--location-id` *(required)*, `--org-id`, `--name`, `--extension`, `--body`, `--body-file` |
| `update` | `--location-id` *(required)*, `--call-park-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-2` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-extension` | `--location-id` *(required)*, `--call-park-extension-id` *(required)*, `--org-id`, `--name`, `--extension`, `--body`, `--body-file` |
| `delete` | `--location-id` *(required)*, `--call-park-id` *(required)*, `--org-id` |
| `delete-extension` | `--location-id` *(required)*, `--call-park-extension-id` *(required)*, `--org-id` |

### call-pickup

| Command | Flags |
|---|---|
| `list` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--order`, `--name` |
| `get` | `--location-id` *(required)*, `--call-pickup-id` *(required)*, `--org-id` |
| `get-available-agents` | `--location-id` *(required)*, `--org-id`, `--call-pickup-name`, `--max`, `--start`, `--name`, `--phone-number`, `--order` |
| `create` | `--location-id` *(required)*, `--org-id`, `--name`, `--notification-type`, `--notification-delay-timer-seconds`, `--agents`, `--body`, `--body-file` |
| `update` | `--location-id` *(required)*, `--call-pickup-id` *(required)*, `--org-id`, `--name`, `--notification-type`, `--notification-delay-timer-seconds`, `--agents`, `--body`, `--body-file` |
| `delete` | `--location-id` *(required)*, `--call-pickup-id` *(required)*, `--org-id` |

### call-queue

| Command | Flags |
|---|---|
| `list-customer-assist` | `--org-id`, `--location-id`, `--max`, `--start`, `--name`, `--phone-number`, `--department-id`, `--department-name`, `--has-cx-essentials` |
| `get-customer-assist` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--has-cx-essentials` |
| `list-announcement-files` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id` |
| `get-forward` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id` |
| `get-selective-forward-rule` | `--location-id` *(required)*, `--queue-id` *(required)*, `--rule-id` *(required)*, `--org-id` |
| `get-holiday-service` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id` |
| `get-night-service` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id` |
| `get-forced-forward` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id` |
| `get-stranded` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id` |
| `get-primary-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-alternate-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-forward-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `get-available-agents` | `--location-id`, `--org-id`, `--max`, `--start`, `--name`, `--phone-number`, `--order` |
| `list-supervisors-customer-assist` | `--org-id`, `--max`, `--start`, `--name`, `--phone-number`, `--order`, `--has-cx-essentials` |
| `get-supervisor-detail-customer-assist` | `--supervisor-id` *(required)*, `--org-id`, `--max`, `--start`, `--name`, `--phone-number`, `--order`, `--has-cx-essentials` |
| `list-available-supervisors-customer-assist` | `--org-id`, `--max`, `--start`, `--name`, `--phone-number`, `--order`, `--has-cx-essentials` |
| `list-available-agents-customer-assist` | `--org-id`, `--max`, `--start`, `--name`, `--phone-number`, `--order`, `--has-cx-essentials` |
| `list-agents-customer-assist` | `--org-id`, `--location-id`, `--queue-id`, `--max`, `--start`, `--name`, `--phone-number`, `--join-enabled`, `--has-cx-essentials`, `--order` |
| `get-agent-customer-assist` | `--id` *(required)*, `--org-id`, `--has-cx-essentials`, `--max`, `--start` |
| `get-playlist-usage` | `--play-list-id` *(required)*, `--playlist-usage-type` |
| `create-customer-assist` | `--location-id` *(required)*, `--org-id`, `--has-cx-essentials`, `--body`, `--body-file` |
| `create-selective-forward-rule` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `create-supervisor-customer-assist` | `--org-id`, `--has-cx-essentials`, `--body`, `--body-file` |
| `switch-mode-forward` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-forward` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-selective-forward-rule` | `--location-id` *(required)*, `--queue-id` *(required)*, `--rule-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-holiday-service` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-night-service` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-forced-forward-service` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-stranded-service` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `assign-unassign-agents-supervisor-customer-assist` | `--supervisor-id` *(required)*, `--org-id`, `--has-cx-essentials`, `--body`, `--body-file` |
| `update-agent-settings-one-more-customer-assist` | `--id` *(required)*, `--org-id`, `--has-cx-essentials`, `--body`, `--body-file` |
| `delete` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id` |
| `delete-announcement-file` | `--location-id` *(required)*, `--queue-id` *(required)*, `--file-name` *(required)*, `--org-id` |
| `delete-selective-forward-rule` | `--location-id` *(required)*, `--queue-id` *(required)*, `--rule-id` *(required)*, `--org-id` |
| `delete-bulk-supervisors` | `--org-id`, `--supervisor-ids`, `--delete-all`, `--body`, `--body-file` |
| `delete-supervisor` | `--supervisor-id` *(required)*, `--org-id` |

### call-recording

| Command | Flags |
|---|---|
| `get` | `--org-id` |
| `get-terms-settings` | `--vendor-id` *(required)*, `--org-id` |
| `get-org-compliance-announcement` | `--org-id` |
| `get-location-compliance-announcement` | `--location-id` *(required)*, `--org-id` |
| `get-regions` | `--org-id` |
| `get-vendor-users` | `--org-id`, `--max`, `--start`, `--standard-user-only` |
| `get-location-vendors` | `--location-id` *(required)*, `--org-id` |
| `get-vendor-users-location` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--standard-user-only` |
| `list-jobs` | `--org-id`, `--max`, `--start` |
| `get-job-status-job` | `--job-id` *(required)*, `--org-id` |
| `get-job-errors-job` | `--job-id` *(required)*, `--org-id`, `--max` |
| `get-org-vendors` | `--org-id` |
| `update` | `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-terms-settings` | `--vendor-id` *(required)*, `--org-id`, `--terms-of-service-enabled`, `--body`, `--body-file` |
| `update-org-compliance-announcement` | `--org-id`, `--inbound-pstncalls-enabled`, `--outbound-pstncalls-enabled`, `--outbound-pstncalls-delay-enabled`, `--delay-in-seconds`, `--body`, `--body-file` |
| `update-location-compliance-announcement` | `--location-id` *(required)*, `--org-id`, `--inbound-pstncalls-enabled`, `--use-org-settings-enabled`, `--outbound-pstncalls-enabled`, `--outbound-pstncalls-delay-enabled`, `--delay-in-seconds`, `--body`, `--body-file` |
| `update-vendor-location` | `--location-id` *(required)*, `--org-id`, `--id`, `--org-default-enabled`, `--storage-region`, `--org-storage-region-enabled`, `--failure-behavior`, `--org-failure-behavior-enabled`, `--body`, `--body-file` |
| `update-org-vendor` | `--org-id`, `--vendor-id`, `--storage-region`, `--failure-behavior`, `--body`, `--body-file` |

### customer-experience-essentials

| Command | Flags |
|---|---|
| `list-wrap-up-reasons` | — |
| `get-wrap-up-reason` | `--wrapup-reason-id` *(required)* |
| `get-available-queues` | `--wrapup-reason-id` *(required)* |
| `get-wrap-up-reason-settings` | `--location-id` *(required)*, `--queue-id` *(required)* |
| `get-screen-pop-configuration` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id` |
| `list-available-agents` | `--location-id` *(required)*, `--org-id`, `--has-cx-essentials` |
| `create-wrap-up-reason` | `--name`, `--description`, `--queues`, `--assign-all-queues-enabled`, `--body`, `--body-file` |
| `validate-wrap-up-reason` | `--name`, `--body`, `--body-file` |
| `update-wrap-up-reason` | `--wrapup-reason-id` *(required)*, `--name`, `--description`, `--queues-to-assign`, `--queues-to-unassign`, `--assign-all-queues-enabled`, `--unassign-all-queues-enabled`, `--body`, `--body-file` |
| `update-wrap-up-reason-settings` | `--location-id` *(required)*, `--queue-id` *(required)*, `--wrapup-reasons`, `--default-wrapup-reason-id`, `--wrapup-timer-enabled`, `--wrapup-timer`, `--body`, `--body-file` |
| `update-screen-pop-configuration` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `delete-wrap-up-reason` | `--wrapup-reason-id` *(required)* |

### hot-desking-sign-in-via-voice-portal

| Command | Flags |
|---|---|
| `location` | `--location-id` *(required)*, `--org-id` |
| `user` | `--person-id` *(required)*, `--org-id` |
| `update-location` | `--location-id` *(required)*, `--org-id`, `--voice-portal-hot-desk-sign-in-enabled`, `--body`, `--body-file` |
| `update-user` | `--person-id` *(required)*, `--org-id`, `--voice-portal-hot-desk-sign-in-enabled`, `--body`, `--body-file` |

### hunt-group

| Command | Flags |
|---|---|
| `list` | `--org-id`, `--location-id`, `--max`, `--start`, `--name`, `--phone-number` |
| `get` | `--location-id` *(required)*, `--hunt-group-id` *(required)*, `--org-id` |
| `get-call-forward` | `--location-id` *(required)*, `--hunt-group-id` *(required)*, `--org-id` |
| `get-selective-forward-rule` | `--location-id` *(required)*, `--hunt-group-id` *(required)*, `--rule-id` *(required)*, `--org-id` |
| `get-primary-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-alternate-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-forward-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `create` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `create-selective-forward-rule` | `--location-id` *(required)*, `--hunt-group-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `switch-mode-call-forward` | `--location-id` *(required)*, `--hunt-group-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update` | `--location-id` *(required)*, `--hunt-group-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-call-forward` | `--location-id` *(required)*, `--hunt-group-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-selective-forward-rule` | `--location-id` *(required)*, `--hunt-group-id` *(required)*, `--rule-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `delete` | `--location-id` *(required)*, `--hunt-group-id` *(required)*, `--org-id` |
| `delete-selective-forward-rule` | `--location-id` *(required)*, `--hunt-group-id` *(required)*, `--rule-id` *(required)*, `--org-id` |

### operating-modes

| Command | Flags |
|---|---|
| `list` | `--name`, `--limit-to-location-id`, `--limit-to-org-level-enabled`, `--max`, `--start`, `--order`, `--org-id` |
| `get` | `--mode-id` *(required)*, `--org-id` |
| `get-holiday` | `--mode-id` *(required)*, `--holiday-id` *(required)*, `--org-id` |
| `list-available-location` | `--location-id` *(required)*, `--org-id` |
| `get-forward-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `create` | `--org-id`, `--body`, `--body-file` |
| `create-holiday` | `--mode-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update` | `--mode-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-holiday` | `--mode-id` *(required)*, `--holiday-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `delete` | `--mode-id` *(required)*, `--org-id` |
| `delete-holiday` | `--mode-id` *(required)*, `--holiday-id` *(required)*, `--org-id` |

### paging-group

| Command | Flags |
|---|---|
| `list` | `--org-id`, `--max`, `--start`, `--location-id`, `--name`, `--phone-number` |
| `get` | `--location-id` *(required)*, `--paging-id` *(required)*, `--org-id` |
| `get-primary-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `create` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update` | `--location-id` *(required)*, `--paging-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `delete` | `--location-id` *(required)*, `--paging-id` *(required)*, `--org-id` |

### single-number-reach

| Command | Flags |
|---|---|
| `get-primary` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-settings-person` | `--person-id` *(required)* |
| `create-person` | `--person-id` *(required)*, `--phone-number`, `--enabled`, `--name`, `--do-not-forward-calls-enabled`, `--answer-confirmation-enabled`, `--body`, `--body-file` |
| `update-settings-person` | `--person-id` *(required)*, `--alert-all-numbers-for-click-to-dial-calls-enabled`, `--body`, `--body-file` |
| `update-settings` | `--person-id` *(required)*, `--id` *(required)*, `--phone-number`, `--enabled`, `--name`, `--do-not-forward-calls-enabled`, `--answer-confirmation-enabled`, `--body`, `--body-file` |
| `delete` | `--person-id` *(required)*, `--id` *(required)*, `--org-id` |

### virtual-extensions

| Command | Flags |
|---|---|
| `list` | `--org-id`, `--max`, `--start`, `--order`, `--extension`, `--phone-number`, `--name`, `--location-name`, `--location-id`, `--org-level-only` |
| `get` | `--extension-id` *(required)*, `--org-id` |
| `get-settings` | `--org-id` |
| `list-range` | `--org-id`, `--max`, `--start`, `--order`, `--name`, `--prefix`, `--location-id`, `--org-level-only` |
| `get-range` | `--extension-range-id` *(required)*, `--org-id` |
| `create` | `--org-id`, `--display-name`, `--phone-number`, `--extension`, `--first-name`, `--last-name`, `--location-id`, `--body`, `--body-file` |
| `validate-external-phone-number` | `--org-id`, `--phone-numbers`, `--body`, `--body-file` |
| `create-range` | `--org-id`, `--name`, `--prefix`, `--patterns`, `--location-id`, `--body`, `--body-file` |
| `validate-prefix-pattern-range` | `--org-id`, `--location-id`, `--name`, `--prefix`, `--patterns`, `--range-id`, `--body`, `--body-file` |
| `update` | `--extension-id` *(required)*, `--org-id`, `--first-name`, `--last-name`, `--display-name`, `--phone-number`, `--extension`, `--body`, `--body-file` |
| `update-settings` | `--org-id`, `--mode`, `--body`, `--body-file` |
| `update-range` | `--extension-range-id` *(required)*, `--org-id`, `--name`, `--prefix`, `--patterns`, `--action`, `--body`, `--body-file` |
| `delete` | `--extension-id` *(required)*, `--org-id` |
| `delete-range` | `--extension-range-id` *(required)*, `--org-id` |

### hot-desk

| Command | Flags |
|---|---|
| `list-sessions` | `--org-id`, `--person-id`, `--workspace-id` |
| `delete-session` | `--session-id` *(required)* |

### location-call-handling

| Command | Flags |
|---|---|
| `get-internal-dialing-configuration` | `--location-id` *(required)*, `--org-id` |
| `get-intercept` | `--location-id` *(required)*, `--org-id` |
| `get-outgoing-permission` | `--location-id` *(required)*, `--org-id` |
| `get-outgoing-auto-transfer` | `--location-id` *(required)*, `--org-id` |
| `get-outgoing-permission-access-code` | `--location-id` *(required)*, `--org-id` |
| `get-outgoing-digit-pattern` | `--location-id` *(required)*, `--org-id` |
| `get-outgoing-digit-pattern-2` | `--location-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id` |
| `generate-example-password` | `--location-id` *(required)*, `--org-id`, `--generate`, `--body`, `--body-file` |
| `create-outgoing-permission-access-code-customer` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `create-outgoing-permission-digit-pattern` | `--location-id` *(required)*, `--org-id`, `--name`, `--pattern`, `--action`, `--transfer-enabled`, `--body`, `--body-file` |
| `update-internal-dialing-configuration` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-intercept` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-outgoing-permission` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-outgoing-auto-transfer` | `--location-id` *(required)*, `--org-id`, `--auto-transfer-number1`, `--auto-transfer-number2`, `--auto-transfer-number3`, `--body`, `--body-file` |
| `delete-outgoing-access-code` | `--location-id` *(required)*, `--org-id`, `--delete-codes`, `--body`, `--body-file` |
| `update-outgoing-digit-pattern` | `--location-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id`, `--name`, `--pattern`, `--action`, `--transfer-enabled`, `--body`, `--body-file` |
| `delete-all-outgoing-access-code` | `--location-id` *(required)*, `--org-id` |
| `delete-all-outgoing-digit-patterns` | `--location-id` *(required)*, `--org-id` |
| `delete-outgoing-digit-pattern` | `--location-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id` |

### location-schedules

| Command | Flags |
|---|---|
| `list` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--name`, `--type` |
| `get` | `--location-id` *(required)*, `--type` *(required)*, `--schedule-id` *(required)*, `--org-id` |
| `get-event` | `--location-id` *(required)*, `--type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)*, `--org-id` |
| `create` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `create-event` | `--location-id` *(required)*, `--type` *(required)*, `--schedule-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update` | `--location-id` *(required)*, `--type` *(required)*, `--schedule-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-event` | `--location-id` *(required)*, `--type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `delete` | `--location-id` *(required)*, `--type` *(required)*, `--schedule-id` *(required)*, `--org-id` |
| `delete-event` | `--location-id` *(required)*, `--type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)*, `--org-id` |

### location-voicemail

| Command | Flags |
|---|---|
| `get` | `--location-id` *(required)*, `--org-id` |
| `get-voiceportal` | `--location-id` *(required)*, `--org-id` |
| `get-voiceportal-passcode-rule` | `--location-id` *(required)*, `--org-id` |
| `list-voicemailgroup` | `--org-id`, `--location-id`, `--name`, `--phone-number`, `--max`, `--start` |
| `get-group` | `--location-id` *(required)*, `--voicemail-group-id` *(required)*, `--org-id` |
| `get-group-fax-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-group-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-voiceportal-available-numbers` | `--location-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `create-group` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update` | `--location-id` *(required)*, `--org-id`, `--voicemail-transcription-enabled`, `--body`, `--body-file` |
| `update-voiceportal` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-group` | `--location-id` *(required)*, `--voicemail-group-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `delete-group` | `--location-id` *(required)*, `--voicemail-group-id` *(required)*, `--org-id` |

### locations

| Command | Flags |
|---|---|
| `list` | `--name`, `--id`, `--org-id`, `--max` |
| `get` | `--location-id` *(required)*, `--org-id` |
| `list-floors` | `--location-id` *(required)* |
| `get-floor` | `--location-id` *(required)*, `--floor-id` *(required)* |
| `create` | `--org-id`, `--body`, `--body-file` |
| `create-floor` | `--location-id` *(required)*, `--floor-number`, `--display-name`, `--body`, `--body-file` |
| `update` | `--location-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-floor` | `--location-id` *(required)*, `--floor-id` *(required)*, `--floor-number`, `--display-name`, `--body`, `--body-file` |
| `delete-floor` | `--location-id` *(required)*, `--floor-id` *(required)* |
| `delete` | `--location-id` *(required)*, `--org-id` |

### numbers

| Command | Flags |
|---|---|
| `get-phone-org` | `--org-id`, `--location-id`, `--max`, `--start`, `--phone-number`, `--available`, `--order`, `--owner-name`, `--owner-id`, `--owner-type`, `--extension`, `--number-type`, `--phone-number-type`, `--state`, `--details`, `--toll-free-numbers`, `--restricted-non-geo-numbers`, `--included-telephony-types`, `--service-number` |
| `list-manage-jobs` | `--org-id`, `--start`, `--max` |
| `get-manage-job-status` | `--job-id` *(required)* |
| `list-manage-job-errors` | `--job-id` *(required)*, `--org-id`, `--start`, `--max` |
| `create-phone-location` | `--location-id` *(required)*, `--org-id`, `--phone-numbers`, `--number-type`, `--number-usage-type`, `--state`, `--subscription-id`, `--carrier-id`, `--body`, `--body-file` |
| `validate-phone` | `--org-id`, `--phone-numbers`, `--body`, `--body-file` |
| `initiate-jobs` | `--body`, `--body-file` |
| `pause-manage-job` | `--job-id` *(required)*, `--org-id` |
| `resume-manage-job` | `--job-id` *(required)*, `--org-id` |
| `manage-state-location` | `--location-id` *(required)*, `--org-id`, `--phone-numbers`, `--action`, `--body`, `--body-file` |
| `delete-phone-location` | `--location-id` *(required)*, `--org-id`, `--phone-numbers`, `--body`, `--body-file` |

### partner-reports-templates

| Command | Flags |
|---|---|
| `list` | `--service`, `--template-id`, `--from`, `--to`, `--region-id`, `--on-behalf-of-sub-partner-org-id`, `--last` |
| `get` | `--report-id` *(required)*, `--on-behalf-of-sub-partner-org-id` |
| `list-2` | `--on-behalf-of-sub-partner-org-id` |
| `create` | `--on-behalf-of-sub-partner-org-id`, `--template-id`, `--start-date`, `--end-date`, `--region-id`, `--body`, `--body-file` |
| `delete` | `--report-id` *(required)*, `--on-behalf-of-sub-partner-org-id` |

### people

| Command | Flags |
|---|---|
| `list` | `--email`, `--display-name`, `--id`, `--org-id`, `--roles`, `--calling-data`, `--location-id`, `--max`, `--exclude-status` |
| `get-person` | `--person-id` *(required)*, `--calling-data` |
| `get-my-own` | `--calling-data` |
| `create-person` | `--calling-data`, `--min-response`, `--body`, `--body-file` |
| `update-person` | `--person-id` *(required)*, `--calling-data`, `--show-all-types`, `--min-response`, `--body`, `--body-file` |
| `delete-person` | `--person-id` *(required)* |

### pstn

| Command | Flags |
|---|---|
| `get-connection-options-location` | `--location-id` *(required)*, `--org-id`, `--service-types` |
| `get-connection-location` | `--location-id` *(required)*, `--org-id` |
| `emergency-address-lookup-verify-address-is-valid` | `--location-id` *(required)*, `--org-id`, `--address1`, `--address2`, `--city`, `--state`, `--postal-code`, `--country`, `--body`, `--body-file` |
| `create-emergency-address-location` | `--location-id` *(required)*, `--org-id`, `--address1`, `--address2`, `--city`, `--state`, `--postal-code`, `--country`, `--body`, `--body-file` |
| `setup-connection-location` | `--location-id` *(required)*, `--org-id`, `--id`, `--premise-route-type`, `--premise-route-id`, `--body`, `--body-file` |
| `update-emergency-address-phone-number` | `--phone-number` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-emergency-address-location` | `--location-id` *(required)*, `--address-id` *(required)*, `--org-id`, `--address1`, `--address2`, `--city`, `--state`, `--postal-code`, `--country`, `--body`, `--body-file` |

### recording-report

| Command | Flags |
|---|---|
| `list-audit-summaries` | `--max`, `--from`, `--to`, `--host-email`, `--site-url`, `--last` |
| `get-audit` | `--recording-id`, `--host-email`, `--max` |
| `list-meeting-archive-summaries` | `--max`, `--from`, `--to`, `--site-url`, `--last` |
| `get-meeting-archive` | `--archive-id` *(required)* |

### reports

| Command | Flags |
|---|---|
| `list` | `--report-id`, `--service`, `--template-id`, `--from`, `--to`, `--last` |
| `get` | `--report-id` *(required)* |
| `create` | `--template-id`, `--start-date`, `--end-date`, `--site-list`, `--body`, `--body-file` |
| `delete` | `--report-id` *(required)* |

### reports-detailed-call-history

| Command | Flags |
|---|---|
| `get` | `--start-time`, `--end-time`, `--locations`, `--max` |
| `get-live-stream` | `--start-time`, `--end-time`, `--locations`, `--max` |

### send-activation-email

| Command | Flags |
|---|---|
| `get-bulk-resend-job-status` | `--org-id` *(required)*, `--job-id` *(required)* |
| `get-bulk-resend-job-errors` | `--org-id` *(required)*, `--job-id` *(required)*, `--max` |
| `initiate-bulk-resend-job` | `--org-id` *(required)* |

### user-call

| Command | Flags |
|---|---|
| `get-person-application-services-settings` | `--person-id` *(required)*, `--org-id` |
| `get-barge-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-forwarding-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-intercept-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-recording-person` | `--person-id` *(required)*, `--org-id` |
| `get-waiting-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-caller-id-person` | `--person-id` *(required)*, `--org-id` |
| `get-person-calling-behavior` | `--person-id` *(required)*, `--org-id` |
| `get-do-not-disturb-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-executive-assistant-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-hoteling-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-person-monitoring-settings` | `--person-id` *(required)*, `--org-id` |
| `get-incoming-permission-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-person-outgoing-permissions-settings` | `--person-id` *(required)*, `--org-id` |
| `list-phone-numbers-person` | `--person-id` *(required)*, `--org-id`, `--prefer-e164-format` |
| `get-person-privacy-settings` | `--person-id` *(required)*, `--org-id` |
| `get-push-to-talk-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-receptionist-client-settings-person` | `--person-id` *(required)*, `--org-id` |
| `list-schedules-person` | `--person-id` *(required)*, `--org-id`, `--start`, `--max`, `--name`, `--type` |
| `get-schedule` | `--person-id` *(required)*, `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--org-id` |
| `get-event-person-schedule` | `--person-id` *(required)*, `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)*, `--org-id` |
| `get-voicemail-person` | `--person-id` *(required)*, `--org-id` |
| `list-move-jobs` | `--org-id`, `--start`, `--max` |
| `get-move-job-status` | `--job-id` *(required)*, `--org-id` |
| `list-move-job-errors` | `--job-id` *(required)*, `--org-id`, `--start`, `--max` |
| `get-music-hold-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-access-codes-person` | `--person-id` *(required)*, `--org-id` |
| `get-transfer-numbers-person` | `--person-id` *(required)*, `--org-id` |
| `get-digit-patterns-person` | `--person-id` *(required)*, `--org-id` |
| `get-digit-pattern-person` | `--person-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id` |
| `get-preferred-endpoint` | `--person-id` *(required)*, `--org-id` |
| `search-shared-line-appearance-members` | `--person-id` *(required)*, `--application-id` *(required)*, `--max`, `--start`, `--location`, `--name`, `--number`, `--order`, `--extension` |
| `get-shared-line-appearance-members` | `--person-id` *(required)*, `--application-id` *(required)* |
| `get-message-summary` | — |
| `list-messages` | `--line-owner-id` |
| `get-agent-list-available-caller-ids` | `--person-id` *(required)*, `--org-id` |
| `get-agent-caller-id` | `--person-id` *(required)* |
| `get-bridge-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-person-secondary-available-numbers` | `--person-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-person-fax-available-numbers` | `--person-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-person-forward-available-numbers` | `--person-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `get-person-primary-numbers` | `--org-id`, `--location-id`, `--max`, `--start`, `--phone-number`, `--license-type` |
| `get-person-ecbn-available-numbers` | `--person-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name` |
| `get-person-intercept-available-numbers` | `--person-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `get-person-ms-teams-settings` | `--person-id` *(required)*, `--org-id` |
| `get-personal-assistant` | `--person-id` *(required)*, `--org-id` |
| `list-available-features` | `--person-id` *(required)*, `--name`, `--phone-number`, `--extension`, `--max`, `--start`, `--order`, `--org-id` |
| `list-features-assigned-mode-management` | `--person-id` *(required)*, `--org-id` |
| `get-users-selective-accept-list` | `--person-id` *(required)*, `--org-id` |
| `get-users-selective-accept-service` | `--person-id` *(required)*, `--id` *(required)*, `--org-id` |
| `get-users-selective-reject-listing` | `--person-id` *(required)*, `--org-id` |
| `get-users-selective-reject-service` | `--person-id` *(required)*, `--id` *(required)*, `--org-id` |
| `get-users-selective-forward` | `--person-id` *(required)*, `--org-id` |
| `get-users-selective-forward-service` | `--person-id` *(required)*, `--id` *(required)*, `--org-id` |
| `get-person-application-services-settings-2` | `--person-id` *(required)*, `--org-id` |
| `search-shared-line-appearance-members-2` | `--person-id` *(required)*, `--max`, `--start`, `--order`, `--location`, `--name`, `--phone-number`, `--extension` |
| `get-shared-line-appearance-members-2` | `--person-id` *(required)* |
| `get-captions-settings` | `--person-id` *(required)*, `--org-id` |
| `get-person-executive-filtering-settings` | `--person-id` *(required)*, `--org-id` |
| `get-person-executive-filtering-settings-2` | `--person-id` *(required)*, `--id` *(required)*, `--org-id` |
| `get-person-executive-alert-settings` | `--person-id` *(required)*, `--org-id` |
| `get-person-executive-assigned-assistants` | `--person-id` *(required)*, `--org-id` |
| `get-person-executive-available-assistants` | `--person-id` *(required)*, `--org-id`, `--max`, `--start`, `--name`, `--phone-number` |
| `get-person-executive-assistant-settings` | `--person-id` *(required)*, `--org-id` |
| `get-person-executive-screening-settings` | `--person-id` *(required)*, `--org-id` |
| `get-count-shared-line-appearance-members` | `--person-id` *(required)*, `--org-id`, `--location-id`, `--member-name`, `--phone-number`, `--extension` |
| `get-timezone-announcement-language-settings-person` | `--person-id` *(required)*, `--org-id` |
| `get-country-calling-configuration` | `--country-code` *(required)*, `--org-id` |
| `create-schedule-person` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `create-event-person-schedule` | `--person-id` *(required)*, `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `reset-voicemail-pin` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `validate-initiate-move-job` | `--org-id`, `--body`, `--body-file` |
| `pause-move-job` | `--job-id` *(required)*, `--org-id` |
| `resume-move-job` | `--job-id` *(required)*, `--org-id` |
| `create-access-codes-person` | `--person-id` *(required)*, `--org-id`, `--code`, `--description`, `--body`, `--body-file` |
| `create-digit-patterns-person` | `--person-id` *(required)*, `--org-id`, `--name`, `--pattern`, `--action`, `--transfer-enabled`, `--body`, `--body-file` |
| `mark-read` | `--message-id`, `--line-owner-id`, `--body`, `--body-file` |
| `mark-unread` | `--message-id`, `--line-owner-id`, `--body`, `--body-file` |
| `create-users-selective-accept-service` | `--person-id` *(required)*, `--org-id`, `--calls-from`, `--accept-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-users-selective-reject-service` | `--person-id` *(required)*, `--org-id`, `--calls-from`, `--reject-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-users-selective-forward-service` | `--person-id` *(required)*, `--org-id`, `--forward-to-phone-number`, `--send-to-voicemail-enabled`, `--calls-from`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--forward-enabled`, `--body`, `--body-file` |
| `create-person-executive-filtering` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-person-application-services-settings` | `--person-id` *(required)*, `--org-id`, `--ring-devices-for-click-to-dial-calls-enabled`, `--ring-devices-for-group-page-enabled`, `--ring-devices-for-call-park-enabled`, `--browser-client-enabled`, `--desktop-client-enabled`, `--tablet-client-enabled`, `--mobile-client-enabled`, `--body`, `--body-file` |
| `update-barge-settings-person` | `--person-id` *(required)*, `--org-id`, `--enabled`, `--tone-enabled`, `--body`, `--body-file` |
| `update-forward-person` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-intercept-settings-person` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-recording-person` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-waiting-settings-person` | `--person-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-caller-id-person` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-person-calling-behavior` | `--person-id` *(required)*, `--org-id`, `--behavior-type`, `--profile-id`, `--body`, `--body-file` |
| `update-do-not-disturb-settings-person` | `--person-id` *(required)*, `--org-id`, `--webex-go-override-enabled`, `--enabled`, `--ring-splash-enabled`, `--body`, `--body-file` |
| `update-executive-assistant-settings-person` | `--person-id` *(required)*, `--org-id`, `--type`, `--body`, `--body-file` |
| `update-hoteling-settings-person` | `--person-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-person-monitoring-settings` | `--person-id` *(required)*, `--org-id`, `--enable-call-park-notification`, `--monitored-elements`, `--body`, `--body-file` |
| `update-incoming-permission-settings-person` | `--person-id` *(required)*, `--org-id`, `--use-custom-enabled`, `--external-transfer`, `--internal-calls-enabled`, `--collect-calls-enabled`, `--body`, `--body-file` |
| `update-person-outgoing-permissions-settings` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-person-privacy-settings` | `--person-id` *(required)*, `--org-id`, `--aa-extension-dialing-enabled`, `--aa-naming-dialing-enabled`, `--enable-phone-status-directory-privacy`, `--enable-phone-status-pickup-barge-in-privacy`, `--monitoring-agents`, `--body`, `--body-file` |
| `update-push-to-talk-settings-person` | `--person-id` *(required)*, `--org-id`, `--allow-auto-answer`, `--connection-type`, `--access-type`, `--members`, `--body`, `--body-file` |
| `update-receptionist-client-settings-person` | `--person-id` *(required)*, `--org-id`, `--reception-enabled`, `--monitored-members`, `--body`, `--body-file` |
| `update-schedule` | `--person-id` *(required)*, `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-event-person-schedule` | `--person-id` *(required)*, `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-voicemail-person` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-music-hold-settings-person` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-access-codes-person` | `--person-id` *(required)*, `--org-id`, `--delete-codes`, `--body`, `--body-file` |
| `update-transfer-numbers-person` | `--person-id` *(required)*, `--org-id`, `--use-custom-transfer-numbers`, `--auto-transfer-number1`, `--auto-transfer-number2`, `--auto-transfer-number3`, `--body`, `--body-file` |
| `update-digit-pattern-control-person` | `--person-id` *(required)*, `--org-id`, `--use-custom-digit-patterns`, `--body`, `--body-file` |
| `update-digit-pattern-person` | `--person-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id`, `--name`, `--pattern`, `--action`, `--transfer-enabled`, `--body`, `--body-file` |
| `assign-unassign-numbers-person` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-preferred-endpoint` | `--person-id` *(required)*, `--org-id`, `--preferred-answer-endpoint-id`, `--body`, `--body-file` |
| `update-shared-line-appearance-members` | `--person-id` *(required)*, `--application-id` *(required)*, `--body`, `--body-file` |
| `update-person-voicemail-passcode` | `--person-id` *(required)*, `--org-id`, `--passcode`, `--body`, `--body-file` |
| `update-agent-caller-id` | `--person-id` *(required)*, `--selected-caller-id`, `--body`, `--body-file` |
| `update-bridge-settings-person` | `--person-id` *(required)*, `--org-id`, `--warning-tone-enabled`, `--body`, `--body-file` |
| `update-person-ms-teams-setting` | `--person-id` *(required)*, `--org-id`, `--setting-name`, `--value`, `--body`, `--body-file` |
| `update-personal-assistant` | `--person-id` *(required)*, `--org-id`, `--enabled`, `--presence`, `--until-date-time`, `--transfer-enabled`, `--transfer-number`, `--alerting`, `--alert-me-first-number-of-rings`, `--body`, `--body-file` |
| `assign-list-features-mode-management` | `--person-id` *(required)*, `--org-id`, `--feature-ids`, `--body`, `--body-file` |
| `update-users-selective-accept` | `--person-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-users-selective-accept-service` | `--person-id` *(required)*, `--id` *(required)*, `--org-id`, `--calls-from`, `--accept-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `update-users-selective-reject-list` | `--person-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-users-selective-reject-service` | `--person-id` *(required)*, `--id` *(required)*, `--org-id`, `--calls-from`, `--reject-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `update-users-selective-forward-list` | `--person-id` *(required)*, `--org-id`, `--enabled`, `--default-phone-number-to-forward`, `--ring-reminder-enabled`, `--destination-voicemail-enabled`, `--body`, `--body-file` |
| `update-users-selective-forward-service` | `--person-id` *(required)*, `--id` *(required)*, `--org-id`, `--forward-to-phone-number`, `--send-to-voicemail-enabled`, `--calls-from`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--forward-enabled`, `--body`, `--body-file` |
| `update-shared-line-appearance-members-2` | `--person-id` *(required)*, `--body`, `--body-file` |
| `update-captions-settings` | `--person-id` *(required)*, `--org-id`, `--user-closed-captions-enabled`, `--user-transcripts-enabled`, `--use-location-settings-enabled`, `--body`, `--body-file` |
| `update-person-executive-filtering-settings` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-person-executive-filtering-settings-2` | `--person-id` *(required)*, `--id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-person-executive-alert-settings` | `--person-id` *(required)*, `--org-id`, `--alerting-mode`, `--next-assistant-number-of-rings`, `--rollover-enabled`, `--rollover-action`, `--rollover-forward-to-phone-number`, `--rollover-wait-time-in-secs`, `--clid-name-mode`, `--custom-clidname`, `--custom-clidname-in-unicode`, `--clid-phone-number-mode`, `--custom-clidphone-number`, `--body`, `--body-file` |
| `update-person-executive-assigned-assistants` | `--person-id` *(required)*, `--org-id`, `--assistant-ids`, `--body`, `--body-file` |
| `update-person-executive-assistant-settings` | `--person-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-person-executive-screening-settings` | `--person-id` *(required)*, `--org-id`, `--enabled`, `--alert-type`, `--alert-anywhere-location-enabled`, `--alert-mobility-location-enabled`, `--alert-shared-call-appearance-location-enabled`, `--body`, `--body-file` |
| `update-timezone-announcement-language-settings-person` | `--person-id` *(required)*, `--org-id`, `--announcement-language`, `--time-zone`, `--body`, `--body-file` |
| `delete-schedule` | `--person-id` *(required)*, `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--org-id` |
| `delete-event-person-schedule` | `--person-id` *(required)*, `--schedule-type` *(required)*, `--schedule-id` *(required)*, `--event-id` *(required)*, `--org-id` |
| `delete-access-codes-person` | `--person-id` *(required)*, `--org-id` |
| `delete-all-digit-patterns-person` | `--person-id` *(required)*, `--org-id` |
| `delete-digit-pattern-person` | `--person-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id` |
| `delete-message` | `--message-id` *(required)* |
| `delete-users-selective-accept-service` | `--person-id` *(required)*, `--id` *(required)*, `--org-id` |
| `delete-users-selective-reject-service` | `--person-id` *(required)*, `--id` *(required)*, `--org-id` |
| `delete-users-selective-forward-service` | `--person-id` *(required)*, `--id` *(required)*, `--org-id` |
| `delete-person-executive-filtering` | `--person-id` *(required)*, `--id` *(required)*, `--org-id` |

### virtual-line-call

| Command | Flags |
|---|---|
| `list` | `--org-id`, `--location-id`, `--max`, `--start`, `--id`, `--owner-name`, `--phone-number`, `--location-name`, `--order`, `--has-device-assigned`, `--has-extension-assigned`, `--has-dn-assigned` |
| `get-recording` | `--virtual-line-id` *(required)*, `--org-id` |
| `get` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-phone-number-assigned` | `--virtual-line-id` *(required)*, `--org-id` |
| `list-devices-assigned` | `--virtual-line-id` *(required)*, `--org-id` |
| `list-networks-handsets` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-caller-id` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-waiting-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-forward` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-incoming-permission-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-outgoing-permissions-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-access-codes` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-transfer-numbers` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-digit-patterns-profile` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-specified-digit-pattern-profile` | `--virtual-line-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id` |
| `get-intercept-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-agent-list-available-caller-ids` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-agent-caller-id` | `--virtual-line-id` *(required)* |
| `get-voicemail` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-music-hold-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-push-to-talk-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-bridge-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-barge-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-privacy-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `get-fax-available-numbers` | `--virtual-line-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-forward-available-numbers` | `--virtual-line-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `get-available-numbers` | `--org-id`, `--location-id`, `--max`, `--start`, `--phone-number` |
| `get-ecbn-available-numbers` | `--virtual-line-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name` |
| `get-intercept-available-numbers` | `--virtual-line-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `get-donotdisturb-settings` | `--virtual-line-id` *(required)*, `--org-id` |
| `create` | `--org-id`, `--first-name`, `--last-name`, `--location-id`, `--display-name`, `--phone-number`, `--extension`, `--caller-id-last-name`, `--caller-id-first-name`, `--caller-id-number`, `--body`, `--body-file` |
| `create-access-codes` | `--virtual-line-id` *(required)*, `--org-id`, `--code`, `--description`, `--body`, `--body-file` |
| `create-digit-pattern-profile` | `--virtual-line-id` *(required)*, `--org-id`, `--name`, `--pattern`, `--action`, `--transfer-enabled`, `--body`, `--body-file` |
| `reset-voicemail-pin` | `--virtual-line-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-recording` | `--virtual-line-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update` | `--virtual-line-id` *(required)*, `--org-id`, `--first-name`, `--last-name`, `--display-name`, `--phone-number`, `--extension`, `--announcement-language`, `--caller-id-last-name`, `--caller-id-first-name`, `--caller-id-number`, `--time-zone`, `--body`, `--body-file` |
| `update-directory-search` | `--virtual-line-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-caller-id` | `--virtual-line-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-waiting-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-forward` | `--virtual-line-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-incoming-permission-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--use-custom-enabled`, `--external-transfer`, `--internal-calls-enabled`, `--collect-calls-enabled`, `--body`, `--body-file` |
| `update-outgoing-permissions-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-access-codes` | `--virtual-line-id` *(required)*, `--org-id`, `--delete-codes`, `--body`, `--body-file` |
| `update-transfer-numbers` | `--virtual-line-id` *(required)*, `--org-id`, `--use-custom-transfer-numbers`, `--auto-transfer-number1`, `--auto-transfer-number2`, `--auto-transfer-number3`, `--body`, `--body-file` |
| `update-digit-pattern-control-profile` | `--virtual-line-id` *(required)*, `--org-id`, `--use-custom-digit-patterns`, `--body`, `--body-file` |
| `update-digit-pattern-profile` | `--virtual-line-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id`, `--name`, `--pattern`, `--action`, `--transfer-enabled`, `--body`, `--body-file` |
| `update-intercept-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-agent-caller-id` | `--virtual-line-id` *(required)*, `--selected-caller-id`, `--body`, `--body-file` |
| `update-voicemail` | `--virtual-line-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-voicemail-passcode` | `--virtual-line-id` *(required)*, `--org-id`, `--passcode`, `--body`, `--body-file` |
| `update-music-hold-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-push-to-talk-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--allow-auto-answer`, `--connection-type`, `--access-type`, `--members`, `--body`, `--body-file` |
| `update-bridge-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--warning-tone-enabled`, `--body`, `--body-file` |
| `update-barge-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--enabled`, `--tone-enabled`, `--body`, `--body-file` |
| `update-privacy-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--aa-extension-dialing-enabled`, `--aa-naming-dialing-enabled`, `--enable-phone-status-directory-privacy`, `--enable-phone-status-pickup-barge-in-privacy`, `--monitoring-agents`, `--body`, `--body-file` |
| `update-donotdisturb-settings` | `--virtual-line-id` *(required)*, `--org-id`, `--enabled`, `--ring-splash-enabled`, `--body`, `--body-file` |
| `delete` | `--virtual-line-id` *(required)*, `--org-id` |
| `delete-access-codes` | `--virtual-line-id` *(required)*, `--org-id` |
| `delete-all-digit-patterns-profile` | `--virtual-line-id` *(required)*, `--org-id` |
| `delete-digit-pattern-profile` | `--virtual-line-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id` |

### workspace-call

| Command | Flags |
|---|---|
| `get-forward` | `--workspace-id` *(required)*, `--org-id` |
| `get-waiting-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-caller-id` | `--workspace-id` *(required)*, `--org-id` |
| `get-monitoring-settings` | `--workspace-id` *(required)*, `--org-id` |
| `list-numbers-associated` | `--workspace-id` *(required)*, `--org-id` |
| `get-incoming-permission-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-outgoing-permission-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-access-codes` | `--workspace-id` *(required)*, `--org-id` |
| `get-intercept-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-transfer-numbers-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-music-hold-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-all-digit-patterns` | `--workspace-id` *(required)*, `--org-id` |
| `get-digit-pattern` | `--workspace-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id` |
| `get-recording` | `--workspace-id` *(required)*, `--org-id` |
| `get-available-numbers` | `--org-id`, `--location-id`, `--max`, `--start`, `--phone-number` |
| `get-ecbn-available-numbers` | `--workspace-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name` |
| `get-forward-available-numbers` | `--workspace-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `get-intercept-available-numbers` | `--workspace-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number`, `--owner-name`, `--extension` |
| `get-anonymous-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-barge-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-donotdisturb-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-bridge-warning-tone` | `--workspace-id` *(required)*, `--org-id` |
| `get-push-to-talk-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-privacy-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-voicemail` | `--workspace-id` *(required)*, `--org-id` |
| `get-sequential-ring` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `get-policy-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-sequential-ring-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-simultaneous-ring-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-simultaneous-ring` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `get-selective-reject-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-selective-reject` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `get-selective-accept-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-selective-accept` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `get-priority-alert-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-priority-alert` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `get-selective-forward-settings` | `--workspace-id` *(required)*, `--org-id` |
| `get-selective-forward` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `get-fax-available-numbers` | `--workspace-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `get-secondary-available-numbers` | `--workspace-id` *(required)*, `--org-id`, `--max`, `--start`, `--phone-number` |
| `create-access-codes` | `--workspace-id` *(required)*, `--org-id`, `--code`, `--description`, `--body`, `--body-file` |
| `create-digit-pattern` | `--workspace-id` *(required)*, `--org-id`, `--name`, `--pattern`, `--action`, `--transfer-enabled`, `--body`, `--body-file` |
| `create-sequential-ring` | `--workspace-id` *(required)*, `--org-id`, `--calls-from`, `--ring-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-simultaneous-ring` | `--workspace-id` *(required)*, `--org-id`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--ring-enabled`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-selective-reject` | `--workspace-id` *(required)*, `--org-id`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--reject-enabled`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-selective-accept` | `--workspace-id` *(required)*, `--org-id`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--accept-enabled`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-priority-alert` | `--workspace-id` *(required)*, `--org-id`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--notification-enabled`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--body`, `--body-file` |
| `create-selective-forward` | `--workspace-id` *(required)*, `--org-id`, `--calls-from`, `--forward-to-phone-number`, `--destination-voicemail-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--numbers`, `--forward-enabled`, `--body`, `--body-file` |
| `update-forward` | `--workspace-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-waiting-settings` | `--workspace-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-caller-id` | `--workspace-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-monitoring-settings` | `--workspace-id` *(required)*, `--org-id`, `--enable-call-park-notification`, `--monitored-elements`, `--body`, `--body-file` |
| `update-incoming-permission-settings` | `--workspace-id` *(required)*, `--org-id`, `--use-custom-enabled`, `--external-transfer`, `--internal-calls-enabled`, `--collect-calls-enabled`, `--body`, `--body-file` |
| `update-outgoing-permission-settings` | `--workspace-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-access-codes` | `--workspace-id` *(required)*, `--org-id`, `--delete-codes`, `--body`, `--body-file` |
| `update-intercept-settings` | `--workspace-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-transfer-numbers-settings` | `--workspace-id` *(required)*, `--org-id`, `--use-custom-transfer-numbers`, `--auto-transfer-number1`, `--auto-transfer-number2`, `--auto-transfer-number3`, `--body`, `--body-file` |
| `update-music-hold-settings` | `--workspace-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-digit-pattern-control` | `--workspace-id` *(required)*, `--org-id`, `--use-custom-digit-patterns`, `--body`, `--body-file` |
| `update-digit-pattern` | `--workspace-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id`, `--name`, `--pattern`, `--action`, `--transfer-enabled`, `--body`, `--body-file` |
| `update-recording` | `--workspace-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-anonymous-settings` | `--workspace-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-barge-settings` | `--workspace-id` *(required)*, `--org-id`, `--enabled`, `--tone-enabled`, `--body`, `--body-file` |
| `update-donotdisturb-settings` | `--workspace-id` *(required)*, `--org-id`, `--enabled`, `--ring-splash-enabled`, `--body`, `--body-file` |
| `update-bridge-warning-tone` | `--workspace-id` *(required)*, `--org-id`, `--warning-tone-enabled`, `--body`, `--body-file` |
| `update-push-to-talk-settings` | `--workspace-id` *(required)*, `--org-id`, `--allow-auto-answer`, `--connection-type`, `--access-type`, `--members`, `--body`, `--body-file` |
| `update-privacy-settings` | `--workspace-id` *(required)*, `--org-id`, `--aa-extension-dialing-enabled`, `--aa-naming-dialing-enabled`, `--enable-phone-status-directory-privacy`, `--enable-phone-status-pickup-barge-in-privacy`, `--monitoring-agents`, `--body`, `--body-file` |
| `update-voicemail` | `--workspace-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-voicemail-passcode` | `--place-id` *(required)*, `--org-id`, `--passcode`, `--body`, `--body-file` |
| `update-sequential-ring` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--ring-enabled`, `--body`, `--body-file` |
| `update-policy-settings` | `--workspace-id` *(required)*, `--org-id`, `--connected-line-id-privacy-on-redirected-calls`, `--body`, `--body-file` |
| `update-sequential-ring-settings` | `--workspace-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-simultaneous-ring-settings` | `--workspace-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-simultaneous-ring` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--ring-enabled`, `--body`, `--body-file` |
| `update-selective-reject-settings` | `--workspace-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-selective-reject` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--reject-enabled`, `--body`, `--body-file` |
| `assign-unassign-numbers-associated` | `--workspace-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `update-selective-accept-settings` | `--workspace-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-selective-accept` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--accept-enabled`, `--body`, `--body-file` |
| `update-priority-alert-settings` | `--workspace-id` *(required)*, `--org-id`, `--enabled`, `--body`, `--body-file` |
| `update-priority-alert` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--calls-from`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--phone-numbers`, `--notification-enabled`, `--body`, `--body-file` |
| `update-selective-forward-settings` | `--workspace-id` *(required)*, `--org-id`, `--enabled`, `--default-phone-number-to-forward`, `--ring-reminder-enabled`, `--destination-voicemail-enabled`, `--body`, `--body-file` |
| `update-selective-forward` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id`, `--calls-from`, `--forward-to-phone-number`, `--destination-voicemail-enabled`, `--schedule-name`, `--schedule-type`, `--schedule-level`, `--anonymous-callers-enabled`, `--unavailable-callers-enabled`, `--numbers`, `--forward-enabled`, `--body`, `--body-file` |
| `delete-all-access-codes` | `--workspace-id` *(required)*, `--org-id` |
| `delete-access-code` | `--workspace-id` *(required)*, `--access-code` *(required)*, `--org-id` |
| `delete-all-digit-patterns` | `--workspace-id` *(required)*, `--org-id` |
| `delete-digit-pattern` | `--workspace-id` *(required)*, `--digit-pattern-id` *(required)*, `--org-id` |
| `delete-sequential-ring` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `delete-simultaneous-ring` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `delete-selective-reject` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `delete-selective-accept` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `delete-priority-alert` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |
| `delete-selective-forward` | `--workspace-id` *(required)*, `--id` *(required)*, `--org-id` |

### workspaces

| Command | Flags |
|---|---|
| `list` | `--org-id`, `--location-id`, `--workspace-location-id`, `--floor-id`, `--display-name`, `--capacity`, `--type`, `--start`, `--max`, `--calling`, `--supported-devices`, `--calendar`, `--device-hosted-meetings-enabled`, `--device-platform`, `--health-level`, `--include-devices`, `--include-capabilities`, `--planned-maintenance`, `--custom-attribute` |
| `get` | `--workspace-id` *(required)*, `--include-devices` |
| `get-capabilities` | `--workspace-id` *(required)* |
| `get-2` | `--workspace-id` *(required)*, `--include-devices`, `--include-capabilities` |
| `create` | `--body`, `--body-file` |
| `update` | `--workspace-id` *(required)*, `--body`, `--body-file` |
| `delete` | `--workspace-id` *(required)* |

### external-voicemail

| Command | Flags |
|---|---|
| `list-calls` | `--line-owner-id` |
| `get-call` | `--call-id` *(required)*, `--line-owner-id` |
| `list-call-history` | `--type` |
| `dial` | `--destination`, `--endpoint-id`, `--single-number-reach-phone-number`, `--line-owner-id`, `--body`, `--body-file` |
| `answer` | `--call-id`, `--endpoint-id`, `--line-owner-id`, `--body`, `--body-file` |
| `reject` | `--call-id`, `--action`, `--line-owner-id`, `--body`, `--body-file` |
| `hangup` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `hold` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `resume` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `mute` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `unmute` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `divert` | `--call-id`, `--destination`, `--to-voicemail`, `--line-owner-id`, `--body`, `--body-file` |
| `transfer` | `--call-id1`, `--call-id2`, `--destination`, `--line-owner-id`, `--body`, `--body-file` |
| `park` | `--call-id`, `--destination`, `--is-group-park`, `--line-owner-id`, `--body`, `--body-file` |
| `get` | `--destination`, `--endpoint-id`, `--single-number-reach-phone-number`, `--line-owner-id`, `--body`, `--body-file` |
| `start-recording` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `stop-recording` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `pause-recording` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `resume-recording` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `transmit-dtmf` | `--call-id`, `--dtmf`, `--line-owner-id`, `--body`, `--body-file` |
| `push` | `--call-id`, `--line-owner-id`, `--body`, `--body-file` |
| `pickup` | `--target`, `--endpoint-id`, `--single-number-reach-phone-number`, `--line-owner-id`, `--body`, `--body-file` |
| `barge` | `--target`, `--endpoint-id`, `--single-number-reach-phone-number`, `--line-owner-id`, `--body`, `--body-file` |
| `update-clear-message-waiting-indicator` | `--id`, `--org-id`, `--action`, `--body`, `--body-file` |
| `pull` | `--endpoint-id`, `--line-owner-id`, `--body`, `--body-file` |

### caller-reputation-provider

| Command | Flags |
|---|---|
| `get-settings` | `--organization-id` |
| `get-status` | `--organization-id` |
| `get` | `--organization-id` |
| `unlock` | `--organization-id`, `--id`, `--body`, `--body-file` |
| `update-settings` | `--organization-id`, `--enabled`, `--id`, `--name`, `--client-id`, `--client-secret`, `--call-block-score-threshold`, `--call-allow-score-threshold`, `--body`, `--body-file` |

### mode-management

| Command | Flags |
|---|---|
| `get-features` | — |
| `get-common` | `--feature-ids` |
| `get-feature` | `--feature-id` *(required)* |
| `get-normal-operation` | `--feature-id` *(required)* |
| `get-operating` | `--feature-id` *(required)*, `--mode-id` *(required)* |
| `switch-multiple-features` | `--feature-ids`, `--operating-mode-name`, `--body`, `--body-file` |
| `switch-normal-operation` | `--feature-id` *(required)* |
| `switch-single-feature` | `--feature-id` *(required)*, `--operating-mode-id`, `--is-manual-switchback-enabled`, `--body`, `--body-file` |
| `extend-current-operating-duration` | `--feature-id` *(required)*, `--operating-mode-id`, `--extension-time`, `--body`, `--body-file` |

### customer-assist

| Command | Flags |
|---|---|
| `list-wrap-up-reasons` | — |
| `get-wrap-up-reason` | `--wrapup-reason-id` *(required)* |
| `get-available-queues` | `--wrapup-reason-id` *(required)* |
| `get-wrap-up-reason-settings` | `--location-id` *(required)*, `--queue-id` *(required)* |
| `get-screen-pop-configuration` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id` |
| `list-available-agents` | `--location-id` *(required)*, `--org-id`, `--has-cx-essentials` |
| `create-wrap-up-reason` | `--name`, `--description`, `--queues`, `--assign-all-queues-enabled`, `--body`, `--body-file` |
| `validate-wrap-up-reason` | `--name`, `--body`, `--body-file` |
| `update-wrap-up-reason` | `--wrapup-reason-id` *(required)*, `--name`, `--description`, `--queues-to-assign`, `--queues-to-unassign`, `--assign-all-queues-enabled`, `--unassign-all-queues-enabled`, `--body`, `--body-file` |
| `update-wrap-up-reason-settings` | `--location-id` *(required)*, `--queue-id` *(required)*, `--wrapup-reasons`, `--default-wrapup-reason-id`, `--wrapup-timer-enabled`, `--wrapup-timer`, `--body`, `--body-file` |
| `update-screen-pop-configuration` | `--location-id` *(required)*, `--queue-id` *(required)*, `--org-id`, `--body`, `--body-file` |
| `delete-wrap-up-reason` | `--wrapup-reason-id` *(required)* |

<!-- codegen:end -->
