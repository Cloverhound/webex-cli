---
name: webex-cli/admin
description: "Webex Admin commands: people, organizations, licenses, roles, groups, reports, events, recordings, and SCIM 2.0 user/group management."
---

# Webex Admin

Commands: `webex admin <resource> <action> [flags]`

Requires admin-level token scopes (`spark-admin:*`). Partner admins managing multiple orgs use `--organization <orgId>` to target a specific org.

## Resources

| Resource | Operations |
|---|---|
| `people` | list, get-person, get-my-own, create-person, update-person, delete-person |
| `organizations` | list, get, delete |
| `licenses` | list, get-license, assign-users |
| `roles` | list, get |
| `groups` | list-search, get, create, update, delete, get-members |
| `reports` | list, get, create, delete |
| `report-templates` | list |
| `events` | list, get |
| `recordings` | list, list-admin-compliance-officer, get, delete, delete-admin |
| `scim-2-users` | search, get, get-me, create, update-put, update-patch, delete |
| `scim-2-groups` | search, get, create, update-put, update-patch, delete, get-members |
| `hybrid-clusters` | list, get |
| `hybrid-connectors` | list, get |
| `authorizations` | list |
| `service-apps` | list, get, authorize |

## People

```bash
# List people
webex admin people list [flags]
  --email <email>            # exact email match
  --display-name <name>      # starts-with name search
  --location-id <id>         # people in a calling location
  --id <id1,id2,...>         # up to 85 person IDs
  --max <n>                  # max per page (≤100 with --calling-data)
  --calling-data true        # include calling configuration
  --roles <roleId,...>       # filter by role IDs
  --exclude-status true      # omit presence/status (improves performance)
  --paginate

# Get / manage
webex admin people get-person --person-id <id>
webex admin people get-my-own

webex admin people create-person --body '{
  "emails": ["user@example.com"],
  "displayName": "Jane Doe",
  "firstName": "Jane",
  "lastName": "Doe",
  "orgId": "<orgId>",
  "roles": ["<roleId>"],
  "licenses": ["<licenseId>"]
}'
webex admin people create-person --calling-data true --body '{...}'

webex admin people update-person --person-id <id> --body '{...}'
webex admin people delete-person --person-id <id>
```

## Organizations

```bash
webex admin organizations list
webex admin organizations get --org-id <id>
webex admin organizations delete --org-id <id>
```

## Licenses

```bash
# List all licenses in the org
webex admin licenses list

# Get details of a specific license
webex admin licenses get-license --license-id <id>

# Assign licenses to users
webex admin licenses assign-users --body '{
  "email": "user@example.com",
  "skuAssignments": [{"skuName": "Webex Calling Professional", "quantity": 1}]
}'
```

## Roles

```bash
webex admin roles list
webex admin roles get --role-id <id>
```

## Groups

```bash
# List and search groups
webex admin groups list-search [flags]
  --filter <expr>            # filter by displayName: eq (equal) or sw (starts with)
  --count <n>                # results per page
  --start-index <n>          # 1-based pagination offset
  --sort-by displayName
  --sort-order ascending|descending
  --include-members true     # include up to 500 members in response

webex admin groups get --group-id <id>
webex admin groups get-members --group-id <id>

webex admin groups create --body '{"displayName":"Admins","members":[{"value":"<personId>"}]}'
webex admin groups update --group-id <id> --body '{...}'
webex admin groups delete --group-id <id>
```

## Reports

```bash
# List available reports
webex admin reports list

# Get a report (may need polling until status = done)
webex admin reports get --report-id <id>

# Create a report (use template ID from report-templates list)
webex admin reports create --body '{"templateId":"<id>","siteList":"company.webex.com"}'

# List available report templates
webex admin report-templates list

# Delete
webex admin reports delete --report-id <id>
```

## Events (Admin Audit)

```bash
webex admin events list [flags]
  --resource <resource>      # e.g. messages, rooms, memberships
  --type <type>              # e.g. created, updated, deleted
  --actor-id <personId>      # events performed by this person
  --from <ISO8601>
  --to <ISO8601>
  --last <duration>          # e.g. 24h, 7d
  --max <n>                  # max 1000 per page
  --paginate

webex admin events get --event-id <id>
```

**Note:** `webex admin events` returns compliance/admin audit events. For messaging-specific events (messages sent/deleted/etc.), use `webex messaging events list`.

## Recordings (Admin)

```bash
# List own recordings
webex admin recordings list [flags]
  --from <ISO8601>
  --to <ISO8601>
  --last <duration>
  --meeting-id <id>
  --topic <text>
  --status available|deleted|purged
  --max <n>                  # max 100 per page

# List ALL org recordings (admin/compliance, max 30-day window)
webex admin recordings list-admin-compliance-officer [same flags]

webex admin recordings get --recording-id <id>
webex admin recordings delete --recording-id <id>
webex admin recordings delete-admin --recording-id <id>    # admin-level delete

# Recycle bin management
webex admin recordings move-recycle-bin --body '{"recordingIds":["<id>"]}'
webex admin recordings restore-recycle-bin --body '{"recordingIds":["<id>"]}'
webex admin recordings purge-recycle-bin --body '{"recordingIds":["<id>"]}'
```

## SCIM 2.0 Users

For automated user provisioning with SCIM-compliant identity providers.

```bash
# Search users
webex admin scim-2-users search [flags]
  --filter <scim-filter>     # e.g. 'userName eq "user@example.com"'
  --count <n>                # results per page (default 100)
  --start-index <n>          # 1-based offset (default 1)
  --sort-by userName|id|meta.lastModified
  --sort-order ascending|descending
  --attributes <list>        # comma-separated attribute list to return
  --return-groups true       # include group membership

webex admin scim-2-users get --user-id <id>
webex admin scim-2-users get-me

# Create
webex admin scim-2-users create --body '{
  "schemas": ["urn:ietf:params:scim:schemas:core:2.0:User"],
  "userName": "user@example.com",
  "name": {"givenName": "Jane", "familyName": "Doe"},
  "emails": [{"value": "user@example.com", "type": "work", "primary": true}]
}'

# Full replace (PUT) vs partial update (PATCH)
webex admin scim-2-users update-put --user-id <id> --body '{...}'
webex admin scim-2-users update-patch --user-id <id> --body '{
  "schemas": ["urn:ietf:params:scim:api:messages:2.0:PatchOp"],
  "Operations": [{"op": "replace", "path": "displayName", "value": "Jane Smith"}]
}'

webex admin scim-2-users delete --user-id <id>
```

## SCIM 2.0 Groups

```bash
webex admin scim-2-groups search [--filter <expr>] [--count <n>] [--start-index <n>]
webex admin scim-2-groups get --group-id <id>
webex admin scim-2-groups get-members --group-id <id>
webex admin scim-2-groups create --body '{"displayName":"IT Admins","members":[{"value":"<userId>"}]}'
webex admin scim-2-groups update-put --group-id <id> --body '{...}'
webex admin scim-2-groups update-patch --group-id <id> --body '{...}'
webex admin scim-2-groups delete --group-id <id>
```

## Hybrid and Infrastructure

```bash
# Hybrid clusters and connectors
webex admin hybrid-clusters list
webex admin hybrid-clusters get --cluster-id <id>
webex admin hybrid-connectors list [--cluster-id <id>]
webex admin hybrid-connectors get --connector-id <id>

# Authorizations (OAuth integrations)
webex admin authorizations list [--person-id <id>]

# Service apps
webex admin service-apps list
webex admin service-apps get --service-app-id <id>
```

## Key Gotchas

1. **`admin people list` vs `calling people list`** — both exist but `admin people list` operates on the Webex user directory (all users); `calling people list` returns users specifically provisioned for Calling. Use `--calling-data true` on admin list to get calling details.

2. **Admin recordings 30-day limit** — `recordings list-admin-compliance-officer` enforces a **30-day** window. Use multiple requests with overlapping ranges to cover longer periods.

3. **SCIM `update-put` replaces the entire user** — omit any attribute to unset it. Use `update-patch` to make surgical changes without affecting unspecified attributes.

4. **`events list` vs `admin-audit list`** — `webex admin events list` returns admin/compliance events. `webex admin admin-audit list` (if present) is a separate endpoint. Always verify the `--resource` and `--type` values against the actual event types in your org.

5. **License assignment** — `licenses assign-users` uses SKU names (not license IDs). Get the exact SKU name from `licenses list`.

6. **Partner admins** — use `--organization <orgId>` global flag to manage a customer org. The org ID can be base64 or UUID format; the CLI auto-decodes.
