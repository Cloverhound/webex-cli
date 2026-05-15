---
name: webex-cli/device
description: "Webex Device commands: devices, workspaces, device configurations, xAPI, workspace locations, and hot desking."
---

# Webex Devices

Commands: `webex device <resource> <action> [flags]`  
Alias: `webex devices <resource> <action>`

## Resources

| Resource | Operations |
|---|---|
| `devices` | list, get, create-mac-address, create-activation-code, update-tags, delete |
| `workspaces` | list, get, get-2, get-capabilities, create, update, delete |
| `device-configurations` | list, update |
| `xapi` | execute-command, query-status |
| `workspace-locations` | list, get, create, update, delete, list-floors, get-floor, create-floor, update-floor, delete-floor |
| `hot-desk` | list-sessions, delete-session |
| `workspace-metrics` | get |
| `workspace-personalization` | get, update |
| `device-call` | (calling settings for device-registered lines) |

## Devices

### List devices
```bash
webex device devices list [flags]
  --product <name>           # e.g. "Cisco Board Pro 55"
  --display-name <name>      # filter by display name
  --person-id <id>           # devices assigned to a person
  --workspace-id <id>        # devices in a workspace
  --location-id <id>         # devices in a location
  --mac <address>            # by MAC address (e.g. 001122334455)
  --serial <num>             # by serial number
  --connection-status <s>    # connected, disconnected, connected_with_issues
  --type <type>              # device type
  --software <version>       # by software version
  --upgrade-channel <c>      # upgrade channel
  --tag <tag>                # filter by tag (repeat for multiple, logical AND)
  --max <n>
  --paginate
```

### Get / delete
```bash
webex device devices get --device-id <id>
webex device devices delete --device-id <id>
```

### Create devices
```bash
# Add a device by MAC address
webex device devices create-mac-address --body '{
  "mac": "00:11:22:33:44:55",
  "workspaceId": "<workspaceId>",
  "model": "Cisco 8841"
}'

# Generate an activation code for a workspace device
webex device devices create-activation-code --body '{"workspaceId":"<id>"}'

# Modify tags on a device
webex device devices update-tags --device-id <id> --body '{"op":"add","tags":["lobby"]}'
```

## Workspaces

### List workspaces
```bash
webex device workspaces list [flags]
  --display-name <name>
  --location-id <id>         # from /locations API (preferred over --workspace-location-id)
  --floor-id <id>
  --type focus|huddle|meetingRoom|open|desk|other
  --calling mpls|webexCalling|thirdPartySipCalling|none
  --calendar microsoft|google|none
  --capacity <n>             # -1 = no capacity set
  --device-platform <p>      # roomdesk, phone, dx, boards, etc.
  --supported-devices <s>    # collaborationDevices, phones
  --health-level <l>
  --include-devices true     # include associated devices in response
  --include-capabilities true
  --max <n>
  --paginate
```

### Get / create / update / delete
```bash
webex device workspaces get --workspace-id <id>
webex device workspaces get-capabilities --workspace-id <id>

webex device workspaces create --body '{
  "displayName": "Conf Room A",
  "type": "meetingRoom",
  "locationId": "<locationId>",
  "floorId": "<floorId>",
  "capacity": 10,
  "calling": {"type": "webexCalling"},
  "calendar": {"type": "microsoft"}
}'

webex device workspaces update --workspace-id <id> --body '{...}'
webex device workspaces delete --workspace-id <id>
```

## Device Configurations

Read and write xConfiguration parameters for a specific device.

```bash
# List configurations (all or filtered)
webex device device-configurations list --device-id <id>
webex device device-configurations list --device-id <id> --key "Audio.Ultrasound.*"
webex device device-configurations list --device-id <id> --key "Conference.MaxReceiveCallRate"

# Update configurations
webex device device-configurations update --device-id <id> --body '{
  "items": [
    {"key": "Audio.Ultrasound.MaxVolume", "value": "70"}
  ]
}'
```

**Key path syntax:**
- Absolute: `Conference.MaxReceiveCallRate`
- Wildcard: `Audio.Ultrasound.*` (all Audio Ultrasound configs)
- Range: `FacilityService.Service[1].Name` (first only), `FacilityService.Service[*].Name` (all)
- URL-encode `[` → `%5B`, `]` → `%5D` when using ranges in flag values

## xAPI

Execute RoomOS commands and query device status remotely.

```bash
# Execute a command
webex device xapi execute-command \
  --device-id <id> \
  --command-name "Dial" \
  --body '{"Dial":{"Number":"12345@example.com","Protocol":"SIP"}}'

# Common commands
webex device xapi execute-command --device-id <id> --command-name "Standby.Deactivate" --body '{}'
webex device xapi execute-command --device-id <id> --command-name "SystemUnit.Boot" --body '{"Action":"Restart"}'
webex device xapi execute-command --device-id <id> --command-name "Audio.Volume.Set" --body '{"Audio":{"Volume":{"Level":50}}}'

# Query status
webex device xapi query-status --device-id <id>
webex device xapi query-status --device-id <id> --name "Audio.Volume"
webex device xapi query-status --device-id <id> --name "Network.IPv4.Address"
```

The `--name` flag accepts up to 10 status expressions per request.

## Workspace Locations

Workspace locations are a legacy concept (superseded by `/locations` for Calling). Use for workspace physical mapping.

```bash
webex device workspace-locations list [--display-name <name>] [--address <text>]
webex device workspace-locations get --workspace-location-id <id>
webex device workspace-locations create --body '{"displayName":"HQ Floor 2","address":"123 Main St"}'
webex device workspace-locations update --workspace-location-id <id> --body '{...}'
webex device workspace-locations delete --workspace-location-id <id>

# Floors
webex device workspace-locations list-floors --workspace-location-id <id>
webex device workspace-locations create-floor --workspace-location-id <id> --body '{"floorNumber":1,"displayName":"Lobby"}'
webex device workspace-locations update-floor --workspace-location-id <id> --floor-id <id> --body '{...}'
webex device workspace-locations delete-floor --workspace-location-id <id> --floor-id <id>
```

## Hot Desking

```bash
# List active hot desk sessions
webex device hot-desk list-sessions [--workspace-id <id>] [--person-id <id>]

# End a session
webex device hot-desk delete-session --session-id <id>
```

## Key Gotchas

1. **`--location-id` vs `--workspace-location-id`** — when listing workspaces, prefer `--location-id` (from the `/locations` API). `--workspace-location-id` is deprecated but still accepted.

2. **xAPI `--command-name` is the RoomOS path** — use the RoomOS command tree path, e.g. `Dial`, `Audio.Volume.Set`, `Standby.Deactivate`. The body JSON structure must match the xAPI call body for that command.

3. **Device configuration key ranges** — URL-encode `[` and `]` in key paths when passing as a flag value: `FacilityService.Service%5B1%5D.Name` instead of `FacilityService.Service[1].Name`.

4. **`create-activation-code` vs `create-mac-address`** — activation codes are used for zero-touch provisioning (device calls home and registers). MAC address creation is for adding a known device by hardware address.

5. **`workspaces get` vs `get-2`** — `get-2` uses a different API path that may return additional attributes. Both accept `--workspace-id`.

<!-- codegen:start -->
## Command Reference

> Auto-generated from Postman collections. Run `make codegen` to update.

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

### device-configurations

| Command | Flags |
|---|---|
| `list` | `--device-id`, `--key` |
| `update` | `--device-id`, `--op`, `--path`, `--body`, `--body-file` |

### devices

| Command | Flags |
|---|---|
| `list` | `--max`, `--start`, `--display-name`, `--person-id`, `--workspace-id`, `--org-id`, `--connection-status`, `--product`, `--type`, `--serial`, `--tag`, `--software`, `--upgrade-channel`, `--error-code`, `--capability`, `--permission`, `--location-id`, `--workspace-location-id`, `--mac`, `--device-platform`, `--planned-maintenance` |
| `get` | `--device-id` *(required)*, `--org-id` |
| `create-mac-address` | `--org-id`, `--mac`, `--model`, `--workspace-id`, `--person-id`, `--password`, `--body`, `--body-file` |
| `create-activation-code` | `--org-id`, `--workspace-id`, `--person-id`, `--model`, `--body`, `--body-file` |
| `update-tags` | `--device-id` *(required)*, `--org-id`, `--op`, `--path`, `--value`, `--body`, `--body-file` |
| `delete` | `--device-id` *(required)*, `--org-id` |

### hot-desk

| Command | Flags |
|---|---|
| `list-sessions` | `--org-id`, `--person-id`, `--workspace-id` |
| `delete-session` | `--session-id` *(required)* |

### workspace-locations

| Command | Flags |
|---|---|
| `list` | `--org-id`, `--display-name`, `--address`, `--country-code`, `--city-name` |
| `get` | `--location-id` *(required)* |
| `list-floors` | `--location-id` *(required)* |
| `get-floor` | `--location-id` *(required)*, `--floor-id` *(required)* |
| `create` | `--display-name`, `--address`, `--country-code`, `--latitude`, `--longitude`, `--city-name`, `--notes`, `--body`, `--body-file` |
| `create-floor` | `--location-id` *(required)*, `--floor-number`, `--display-name`, `--body`, `--body-file` |
| `update` | `--location-id` *(required)*, `--display-name`, `--address`, `--country-code`, `--latitude`, `--longitude`, `--id`, `--city-name`, `--notes`, `--body`, `--body-file` |
| `update-floor` | `--location-id` *(required)*, `--floor-id` *(required)*, `--floor-number`, `--display-name`, `--body`, `--body-file` |
| `delete` | `--location-id` *(required)* |
| `delete-floor` | `--location-id` *(required)*, `--floor-id` *(required)* |

### workspace-metrics

| Command | Flags |
|---|---|
| `workspace-metrics` | `--workspace-id`, `--metric-name`, `--aggregation`, `--from`, `--to`, `--unit`, `--sort-by`, `--last` |
| `duration` | `--workspace-id`, `--aggregation`, `--measurement`, `--from`, `--to`, `--last` |

### workspace-personalization

| Command | Flags |
|---|---|
| `get-task` | `--workspace-id` *(required)* |
| `personalize` | `--workspace-id` *(required)*, `--email`, `--body`, `--body-file` |

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

### xapi

| Command | Flags |
|---|---|
| `query-status` | `--device-id`, `--name` |
| `query-schema` | `--device-id`, `--status`, `--command` |
| `execute-command` | `--command-name` *(required)*, `--body`, `--body-file` |

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
| `delete` | `--location-id` *(required)*, `--org-id` |
| `delete-floor` | `--location-id` *(required)*, `--floor-id` *(required)* |

<!-- codegen:end -->
