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

<!-- codegen:start -->
## Command Reference

> Auto-generated from Postman collections. Run `make codegen` to update.

### admin-audit

| Command | Flags |
|---|---|
| `list-events` | `--org-id`, `--from`, `--to`, `--actor-id`, `--max`, `--offset`, `--event-categories`, `--last` |
| `list-event-categories` | — |

### authorizations

| Command | Flags |
|---|---|
| `list-user` | `--person-id`, `--person-email` |
| `get-expiration-status-token` | — |
| `delete-org-client-id` | `--client-id`, `--org-id` |
| `delete` | `--authorization-id` *(required)* |

### classifications

| Command | Flags |
|---|---|
| `list` | — |

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

### domain-management

| Command | Flags |
|---|---|
| `get-verification-token` | `--org-id` *(required)*, `--domain`, `--body`, `--body-file` |
| `verify` | `--org-id` *(required)*, `--domain`, `--claim-domain`, `--reserve-domain`, `--body`, `--body-file` |
| `claim` | `--org-id` *(required)*, `--body`, `--body-file` |
| `unverify` | `--org-id` *(required)*, `--domain`, `--remove-pending`, `--body`, `--body-file` |
| `unclaim` | `--org-id` *(required)*, `--domain`, `--body`, `--body-file` |

### identity-org

| Command | Flags |
|---|---|
| `get` | `--org-id` *(required)* |
| `generate-otp` | `--org-id` *(required)*, `--user-id` *(required)* |
| `update` | `--org-id` *(required)*, `--schemas`, `--display-name`, `--preferred-language`, `--body`, `--body-file` |

### events

| Command | Flags |
|---|---|
| `list` | `--resource`, `--type`, `--actor-id`, `--from`, `--to`, `--max`, `--service-type`, `--last` |
| `get` | `--event-id` *(required)* |

### groups

| Command | Flags |
|---|---|
| `list-search` | `--org-id`, `--filter`, `--attributes`, `--sort-by`, `--sort-order`, `--include-members`, `--start-index`, `--count` |
| `get` | `--group-id` *(required)*, `--include-members` |
| `get-members` | `--group-id` *(required)*, `--start-index`, `--count` |
| `create` | `--body`, `--body-file` |
| `update` | `--group-id` *(required)*, `--body`, `--body-file` |
| `delete` | `--group-id` *(required)* |

### guest-management

| Command | Flags |
|---|---|
| `get-count` | — |
| `create` | `--subject`, `--display-name`, `--body`, `--body-file` |

### analytics

| Command | Flags |
|---|---|
| `historical-data-related-messaging` | `--from`, `--to`, `--last` |
| `historical-data-related-room-devices` | `--from`, `--to`, `--last` |
| `historical-data-related-meetings` | `--site-url`, `--from`, `--to`, `--last` |

### hybrid-clusters

| Command | Flags |
|---|---|
| `list` | `--org-id` |
| `get` | `--hybrid-cluster-id` *(required)*, `--org-id` |

### hybrid-connectors

| Command | Flags |
|---|---|
| `list` | `--org-id` |
| `get` | `--connector-id` *(required)* |

### licenses

| Command | Flags |
|---|---|
| `list` | `--org-id` |
| `get-license` | `--license-id` *(required)*, `--include-assigned-to`, `--next`, `--limit` |
| `assign-users` | `--body`, `--body-file` |

### meeting-qualities

| Command | Flags |
|---|---|
| `get` | `--meeting-id`, `--max`, `--offset` |

### org-contacts

| Command | Flags |
|---|---|
| `get` | `--org-id` *(required)*, `--contact-id` *(required)* |
| `list` | `--org-id` *(required)*, `--keyword`, `--source`, `--limit`, `--group-ids` |
| `create` | `--org-id` *(required)*, `--body`, `--body-file` |
| `bulk-create-update` | `--org-id` *(required)*, `--body`, `--body-file` |
| `bulk-delete` | `--org-id` *(required)*, `--schemas`, `--object-ids`, `--body`, `--body-file` |
| `update` | `--org-id` *(required)*, `--contact-id` *(required)*, `--body`, `--body-file` |
| `delete` | `--org-id` *(required)*, `--contact-id` *(required)* |

### organizations

| Command | Flags |
|---|---|
| `list` | — |
| `get` | `--org-id` *(required)* |
| `delete` | `--org-id` *(required)* |

### partner-administrators

| Command | Flags |
|---|---|
| `get-all-customers-managed-admin` | `--managed-by` |
| `get-all-admins-assigned-customer` | `--org-id` *(required)* |
| `assign-admin-customer` | `--org-id` *(required)*, `--person-id` *(required)* |
| `unassign-admin-customer` | `--org-id` *(required)*, `--person-id` *(required)* |
| `revoke-all-admin-roles-person-id` | `--person-id` *(required)* |

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

### report-templates

| Command | Flags |
|---|---|
| `list` | — |

### reports

| Command | Flags |
|---|---|
| `list` | `--report-id`, `--service`, `--template-id`, `--from`, `--to`, `--last` |
| `get` | `--report-id` *(required)* |
| `create` | `--template-id`, `--start-date`, `--end-date`, `--site-list`, `--body`, `--body-file` |
| `delete` | `--report-id` *(required)* |

### resource-memberships

| Command | Flags |
|---|---|
| `list-group` | `--license-id`, `--person-id`, `--person-org-id`, `--status`, `--max` |
| `get-group` | `--resource-group-membership-id` *(required)* |
| `list-group-v2` | `--license-id`, `--id`, `--org-id`, `--status`, `--type`, `--max` |
| `update-group` | `--resource-group-membership-id` *(required)*, `--resource-group-id`, `--license-id`, `--person-id`, `--person-org-id`, `--status`, `--body`, `--body-file` |

### resource-groups

| Command | Flags |
|---|---|
| `list` | `--org-id` |
| `get` | `--resource-group-id` *(required)* |

### roles

| Command | Flags |
|---|---|
| `list` | — |
| `get` | `--role-id` *(required)* |

### scim-2-bulk

| Command | Flags |
|---|---|
| `user-api` | `--org-id` *(required)*, `--body`, `--body-file` |

### scim-2-groups

| Command | Flags |
|---|---|
| `search` | `--org-id` *(required)*, `--filter`, `--attributes`, `--excluded-attributes`, `--sort-by`, `--sort-order`, `--start-index`, `--count`, `--include-members`, `--member-type` |
| `get` | `--org-id` *(required)*, `--group-id` *(required)*, `--excluded-attributes` |
| `get-members` | `--org-id` *(required)*, `--group-id` *(required)*, `--start-index`, `--count`, `--member-type` |
| `create` | `--org-id` *(required)*, `--body`, `--body-file` |
| `update-put` | `--org-id` *(required)*, `--group-id` *(required)*, `--body`, `--body-file` |
| `update-patch` | `--org-id` *(required)*, `--group-id` *(required)*, `--body`, `--body-file` |
| `delete` | `--org-id` *(required)*, `--group-id` *(required)* |

### scim-2-users

| Command | Flags |
|---|---|
| `search` | `--org-id` *(required)*, `--filter`, `--attributes`, `--excluded-attributes`, `--sort-by`, `--sort-order`, `--start-index`, `--count`, `--return-groups`, `--include-group-details`, `--group-usage-types` |
| `get` | `--org-id` *(required)*, `--user-id` *(required)* |
| `get-me` | — |
| `create` | `--org-id` *(required)*, `--body`, `--body-file` |
| `update-put` | `--org-id` *(required)*, `--user-id` *(required)*, `--body`, `--body-file` |
| `update-patch` | `--org-id` *(required)*, `--user-id` *(required)*, `--body`, `--body-file` |
| `delete` | `--org-id` *(required)*, `--user-id` *(required)* |

### security-audit

| Command | Flags |
|---|---|
| `list-events` | `--org-id`, `--start-time`, `--end-time`, `--actor-id`, `--max`, `--event-categories` |

### settings

| Command | Flags |
|---|---|
| `get-org` | `--org-id` *(required)*, `--setting-key` *(required)* |
| `create-update-org` | `--org-id` *(required)*, `--key`, `--value`, `--body`, `--body-file` |

### partner-tags

| Command | Flags |
|---|---|
| `get-all-customer` | `--type` |
| `get-customer-org` | `--org-id` *(required)* |
| `get-all-customers-set` | `--tags`, `--max` |
| `subscription-list-name-set` | `--tags`, `--max` |
| `get-subscription` | `--org-id` *(required)*, `--subscription-id` *(required)* |
| `create-replace-customer-provided-ones` | `--org-id` *(required)*, `--body`, `--body-file` |
| `create-replace-subscription-provided-ones` | `--org-id` *(required)*, `--subscription-id` *(required)*, `--body`, `--body-file` |

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

### scim-2-schemas

| Command | Flags |
|---|---|
| `get-group` | — |
| `get-user` | — |
| `get-group-id` | `--schema-id` *(required)* |

### send-activation-email

| Command | Flags |
|---|---|
| `get-bulk-resend-job-status` | `--org-id` *(required)*, `--job-id` *(required)* |
| `get-bulk-resend-job-errors` | `--org-id` *(required)*, `--job-id` *(required)*, `--max` |
| `initiate-bulk-resend-job` | `--org-id` *(required)* |

### service-apps

| Command | Flags |
|---|---|
| `create-access-token` | `--application-id` *(required)*, `--client-id`, `--client-secret`, `--target-org-id`, `--body`, `--body-file` |

### live-monitoring

| Command | Flags |
|---|---|
| `get-meeting-metrics-categorized-country` | `--site-ids`, `--site-url`, `--body`, `--body-file` |

### archive-users

| Command | Flags |
|---|---|
| `get` | `--org-id` *(required)*, `--useruuid` *(required)* |

<!-- codegen:end -->
