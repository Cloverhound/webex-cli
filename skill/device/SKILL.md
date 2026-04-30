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
