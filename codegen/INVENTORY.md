# Postman refresh inventory

Compared with PR base `85b6a523894c82bb6b8440a103b6f05d39762312`. Added names include renamed/moved operations, not only new Cisco APIs. See [NAMING.md](NAMING.md) for migration details and live route checks.

## Coverage

| Area | Groups | Spec operations |
| --- | ---: | ---: |
| `admin` | 39 | 149 |
| `calling` | 49 | 1111 |
| `cc` | 64 | 498 |
| `device` | 10 | 113 |
| `meetings` | 23 | 173 |
| `messaging` | 12 | 65 |

Counts cover the normalized spec; handwritten commands and codegen exclusions can differ. Devices has no new endpoints.

## admin

### Added canonical names

| Full command | Method and path | Purpose |
| --- | --- | --- |
| `webex admin archive-users query` | `GET /identity/organizations/{orgId}/v1/ArchivedUser` | Retrieves archived user information for the specified organization |
| `webex admin identity-org update-authentication-configuration-settings` | `PATCH /identity/organizations/{orgId}/authenticationConfig` | Update the authentication configuration details, by organizationID. |
| `webex admin identity-org update-password-policy` | `PATCH /identity/organizations/{orgId}/passwordPolicy` | Update Organization Password Policy, by organizationID. |
| `webex admin recordings query` | `POST /recordings/query` | Queries recordings with filters in the request body |
| `webex admin recordings query-admin-compliance-officer` | `POST /admin/recordings/query` | Queries recordings for an admin or compliance officer with filters in the request body |

## calling

### Added canonical names

| Full command | Method and path | Purpose |
| --- | --- | --- |
| `webex calling call-controls-for-me answer` | `POST /telephony/calls/members/me/answer` | Answer an incoming call |
| `webex calling call-controls-for-me dial` | `POST /telephony/calls/members/me/dial` | Initiate an outbound call to a specified destination |
| `webex calling call-controls-for-me get` | `GET /telephony/calls/members/me/calls/{callId}` | Get the details of the specified active call for the user. |
| `webex calling call-controls-for-me hangup` | `POST /telephony/calls/members/me/hangup` | Hangup a call |
| `webex calling call-controls-for-me list` | `GET /telephony/calls/members/me/calls` | Get the list of details for all active calls associated with the user. |
| `webex calling metrics get-call-quality-stats` | `GET /v1/analytics/callQualityStats` | Returns aggregated Webex Calling call and media-quality statistics for the authenticated organization |

## cc

### Added canonical names

| Full command | Method and path | Purpose |
| --- | --- | --- |
| `webex cc activities describe` | `GET /{orgId}/project/{projectId}/v2/activities/{activityName}` | Retrieve a single activity's full definition, including inputs, outputs, named ports, and the JSON Schema for its inputs — all returned inline. |
| `webex cc activities get-input-choices` | `GET /{orgId}/project/{projectId}/v2/activities/{activityName}/inputs/{inputName}/choices` | Resolve choices for an activity input |
| `webex cc activities list-definitions` | `GET /{orgId}/project/{projectId}/v2/activities` | List all available activity types with their full schemas, including inputs, outputs, ports, and descriptions |
| `webex cc agent-personal-greeting-files create` | `POST /organization/{orgid}/v2/agent-personal-greeting` | Create a new Greeting File in a given organization. |
| `webex cc agent-personal-greeting-files delete-id` | `DELETE /organization/{orgid}/v2/agent-personal-greeting/{id}` | Delete an existing Greeting File by ID in a given organization. |
| `webex cc agent-personal-greeting-files delete-references` | `POST /organization/{orgid}/agent-personal-greeting/delete-reference` | Removes all references of the specified agent (ciUserId) from greeting files in the given organization |
| `webex cc agent-personal-greeting-files get-id` | `GET /organization/{orgid}/v2/agent-personal-greeting/{id}` | Retrieve an existing Greeting File by ID in a given organization. |
| `webex cc agent-personal-greeting-files patch-id` | `PATCH /organization/{orgid}/v2/agent-personal-greeting/{id}` | Partially update Greeting File by ID in a given organization. |
| `webex cc agent-personal-greeting-files update-id` | `PUT /organization/{orgid}/v2/agent-personal-greeting/{id}` | Update an existing Greeting File by ID in a given organization. |
| `webex cc asset bulk-save` | `POST /organization/{orgid}/asset/bulk` | Create, Update or delete Assets in bulk in a given organization |
| `webex cc asset create` | `POST /organization/{orgid}/asset` | Create a new Asset in a given organization |
| `webex cc asset delete-id` | `DELETE /organization/{orgid}/asset/{id}` | Delete an existing Asset by ID in a given organization |
| `webex cc asset get-id` | `GET /organization/{orgid}/asset/{id}` | Retrieve an existing Asset by ID in a given organization |
| `webex cc asset list` | `GET /organization/{orgid}/v2/asset` | Retrieve a list of Assets in a given organization |
| `webex cc asset list-references` | `GET /organization/{orgid}/asset/{id}/incoming-references` | Retrieve a list of all entities that have reference to an existing Asset by ID in a given organization |
| `webex cc asset patch-id` | `PATCH /organization/{orgid}/asset/{id}` | Partially update Asset by ID in a given organization |
| `webex cc asset update-id` | `PUT /organization/{orgid}/asset/{id}` | Update an existing Asset by ID in a given organization |
| `webex cc call-monitoring whisper-coach-request` | `POST /v1/monitor/{interactionId}/coach` | This feature is currently in Beta |
| `webex cc campaign-group list` | `GET /v3/campaign-management/campaign-groups/{campaignGroupName}/campaigns` | Retrieves the list of campaigns that belong to the specified campaign group |
| `webex cc campaign-manager get-valid-times` | `GET /v1/organization/{orgId}/getValidCampaignTimes` | Gets valid campaign times for a campaign and agent |
| `webex cc channel bulk-save` | `POST /organization/{orgid}/channel/bulk` | Create, Update or delete Channels in bulk in a given organization. |
| `webex cc channel create` | `POST /organization/{orgid}/channel` | Create a new Channel in a given organization |
| `webex cc channel delete-id` | `DELETE /organization/{orgid}/channel/{id}` | Delete an existing Channel by ID in a given organization. |
| `webex cc channel get-id` | `GET /organization/{orgid}/channel/{id}` | Retrieve an existing Channel by ID in a given organization |
| `webex cc channel list` | `GET /organization/{orgid}/v2/channel` | Retrieve a list of Channels in a given organization |
| `webex cc channel list-references` | `GET /organization/{orgid}/channel/{id}/incoming-references` | Retrieve a list of all entities that have reference to an existing Channel by ID in a given organization |
| `webex cc channel patch-id` | `PATCH /organization/{orgid}/channel/{id}` | Partially update a channel by ID |
| `webex cc channel update-id` | `PUT /organization/{orgid}/channel/{id}` | Update an existing Channel by ID in a given organization. |
| `webex cc contact-list-management update-status-across-campaign-chain` | `PATCH /v3/campaign-management/campaigns/{campaignId}/contacts/{contactId}` | Synchronously closes the specified contacts and returns the outcome in the same response |
| `webex cc contact-service-queue delete-references` | `POST /organization/{orgid}/contact-service-queue/delete-reference` | Removes the references to the specified entities (such as teams, sites, or agents) from Contact Service Queues for a given organization. |
| `webex cc contact-service-queue list-by-dynamic-skills` | `POST /organization/{orgid}/contact-service-queue/fetch-by-dynamic-skills-and-skillProfile` | Retrieve a list of skill-based Contact Service Queues that match the given dynamic skills and skill profile criteria in a given organization. |
| `webex cc contact-service-queue list-by-skill-profile` | `GET /organization/{orgid}/contact-service-queue/by-skill-profile-id/{id}` | Retrieve a list of skill-based Contact Service Queues associated with a given skill profile ID, accessible to authorized clients in a given organization. |
| `webex cc contact-service-queue list-by-user-skill-profile` | `POST /organization/{orgid}/contact-service-queue/fetch-by-userId-skillProfileId` | Retrieve a list of skill-based Contact Service Queues associated with the given skill profile ID and user ID combination in a given organization. |
| `webex cc contact-service-queue list-manually-assignable` | `POST /organization/{orgid}/contact-service-queue/fetch-manually-assignable-queues` | Retrieve a list of Contact Service Queues that are eligible for manual contact assignment based on the provided criteria in a given organization. |
| `webex cc contact-service-queue list-references` | `GET /organization/{orgid}/contact-service-queue/{id}/incoming-references` | Retrieve a list of all entities that have reference to an existing Contact Service Queue by ID in a given organization. |
| `webex cc events list-specifications` | `GET /{orgId}/project/{projectId}/v2/event-specifications` | List the event types available for use in `event_flows[]` |
| `webex cc external-data-updates update-task-global-variables` | `PUT /v1/data/updateExternal` | Updates numeric global variable values associated with a completed Webex Contact Center task |
| `webex cc flow delete` | `DELETE /{orgId}/project/{projectId}/flows/{flowId}` | Permanently deletes a flow or subflow |
| `webex cc flow export-legacy` | `GET /flow-store/{orgId}/project/{projectId}/flows/{flowId}:export` | Returns the exported flow/subflow in response. |
| `webex cc flow get` | `GET /{orgId}/project/{projectId}/v2/flows/{flowId}` | Retrieve the current draft of a flow/subflow as a flow document |
| `webex cc flow import-legacy` | `POST /flow-store/{orgId}/project/{projectId}/flows:import` | Returns the imported flow/subflow in response. |
| `webex cc flow lock` | `POST /{orgId}/project/{projectId}/flows/{flowId}:lock` | Lock a flow to prevent concurrent edits by other users |
| `webex cc flow patch-draft` | `PATCH /{orgId}/project/{projectId}/v2/flows/{flowId}` | Apply partial updates to an existing flow/subflow draft without replacing the entire document |
| `webex cc flow save-draft` | `POST /{orgId}/project/{projectId}/v2/flows/{flowId}` | Save a complete flow/subflow document as the current draft, replacing the existing draft |
| `webex cc flow search` | `GET /{orgId}/project/{projectId}/flows:search` | Returns a list of flows in response |
| `webex cc flow unlock` | `POST /{orgId}/project/{projectId}/flows/{flowId}:unlock` | Release the lock on a flow to allow other users to edit it. |
| `webex cc flow validate` | `POST /{orgId}/project/{projectId}/v2/flows:validate` | Dry-run validate a flow/subflow definition without persisting it |
| `webex cc flow validate-draft` | `GET /{orgId}/project/{projectId}/v2/flows/{flowId}:validate` | Validate the current draft of an existing flow/subflow (read-only operation) |
| `webex cc functions create` | `POST /v1/{orgId}/functions` | Create a new custom function |
| `webex cc functions delete` | `DELETE /v1/{orgId}/functions/{id}` | Delete a custom function by ID |
| `webex cc functions export` | `POST /v1/{orgId}/functions/{id}:export` | Export a custom function for the given version or publish tag |
| `webex cc functions get` | `GET /v1/{orgId}/functions/{id}` | Retrieve a custom function by its ID |
| `webex cc functions import` | `POST /v1/{orgId}/functions:import` | Import a custom function from a previously exported function-definition JSON file, uploaded as the multipart `file` part (not a zip or base64 envelope) |
| `webex cc functions list` | `GET /v1/{orgId}/functions` | List or search custom functions in the organization |
| `webex cc functions list-options` | `GET /v1/{orgId}/functions/options` | List the supported programming languages, runtimes, and corresponding wrapper code for custom functions |
| `webex cc functions lock` | `POST /v1/{orgId}/functions/{id}:lock` | Acquire an edit lock on a custom function to prevent concurrent writes by other users. |
| `webex cc functions publish` | `POST /v1/{orgId}/functions/{id}:publish` | Publish the latest draft of a custom function under one or more tags (`Dev`, `Test`, `Latest`, `Live`) |
| `webex cc functions test` | `POST /v1/{orgId}/functions/{id}:test` | Execute a custom function with a test payload |
| `webex cc functions unlock` | `POST /v1/{orgId}/functions/{id}:unlock` | Release the edit lock on a custom function so that other users can edit it. |
| `webex cc functions update` | `PUT /v1/{orgId}/functions/{id}` | Update an existing custom function by ID |
| `webex cc outdial-ani bulk-save-entries` | `POST /organization/{orgid}/outdial-ani/{outDialAniId}/entry/bulk` | Create, Update or delete Outdial ANI Entries in bulk for an Address Book in a given organization. |
| `webex cc outdial-ani list-entries` | `GET /organization/{orgid}/outdial-ani/entry` | Retrieve a list of Outdial ANI Entries in a given organization. |
| `webex cc outdial-ani list-entries-for-ani` | `GET /organization/{orgid}/v2/outdial-ani/{outDialAniId}/entry` | Retrieve a list of Outdial ANI Entries in a given organization. |
| `webex cc search-metadata get` | `GET /search/v2/meta` | Returns schema metadata for the Webex Contact Center Search GraphQL API, describing supported queries, fields, nested structures, data types, filters, sorting options, grouping capabilities, and aggregation operations. |
| `webex cc tasks append-message` | `POST /v2/tasks/{taskId}/messages` | This feature is currently in Beta |
| `webex cc tasks drop-participant-conference` | `POST /v1/tasks/{taskId}/conference/participants/{participantId}/drop` | Access this endpoint when the user needs to drop a specific participant from an active conference associated with a task |
| `webex cc tasks pause-digital` | `POST /v1/tasks/{taskId}/pause` | Access this endpoint when users such as administrators, supervisors, or agents with an agent license need to pause a task that cannot be handled immediately |
| `webex cc tasks resume-digital` | `POST /v1/tasks/{taskId}/resume` | Access this endpoint when users such as administrators, supervisors, or agents with an agent license need to resume a previously paused task |
| `webex cc tasks unhold` | `POST /v1/tasks/{taskId}/unhold` | Access this endpoint when the user has to resume a call from hold |
| `webex cc templates get-flow` | `GET /templates/{id}` | Retrieve a specific flow template by its ID |
| `webex cc templates list-flow` | `GET /templates` | List available flow templates that can be used to create new flows. |
| `webex cc usage-reports create` | `POST /v1/usage-reports` | Creates a new usage report for the specified organization, resource type, |
| `webex cc usage-reports delete` | `DELETE /v1/usage-reports/{reportId}` | Deletes the specified usage report and its associated file |
| `webex cc usage-reports download-file` | `GET /v1/usage-reports/{fileId}/download` | Downloads one completed report file as a binary stream. |
| `webex cc usage-reports get` | `GET /v1/usage-reports/{reportId}` | Retrieves detailed information for a specific usage report. |
| `webex cc usage-reports list` | `GET /v1/usage-reports` | Retrieves a list of usage reports for the organization |
| `webex cc usage-reports list-resource-types` | `GET /v1/usage-reports/resource-types` | Returns the list of available resource types that can be used for report generation, along with available data dates. |
| `webex cc users bulk-partial-update-dynamic-skills` | `PATCH /organization/{orgid}/user/bulk/update-dynamic-skill/{skillId}` | Assign or unassign a dynamic skill to/from multiple users in bulk for a given organization. |
| `webex cc users list-by-ids` | `POST /organization/{orgid}/user/fetch-user-details-by-ids` | Retrieve an existing User's first name, last name and email by list of IDs in a given organization. |
| `webex cc users list-call-monitoring-id` | `GET /organization/{orgid}/user/by-call-monitoring-id/{id}` | Fetch paginated users associated to the selected call monitoring team filters while enforcing team ACL. |

### Existing names with changed routes

| Full command | Previous route | Current route |
| --- | --- | --- |
| `webex cc contact-list-management get-within-campaign` | `GET /v3/campaign-management/campaigns/{campaignId}/contact-lists` | `GET /v4/campaign-management/campaigns/{campaignId}/contact-lists` |
| `webex cc flow export` | `GET /flow-store/{orgId}/project/{projectId}/flows/{flowId}:export` | `GET /{orgId}/project/{projectId}/v2/flows/{flowId}:export` |
| `webex cc flow import` | `POST /flow-store/{orgId}/project/{projectId}/flows:import` | `POST /{orgId}/project/{projectId}/v2/flows:import` |
| `webex cc flow list` | `GET /flow-store/{orgId}/project/{projectId}/flows` | `GET /{orgId}/project/{projectId}/flows` |
| `webex cc flow publish` | `POST /flow-store/{orgId}/project/{projectId}/flows/{flowId}:publish` | `POST /{orgId}/project/{projectId}/flows/{flowId}:publish` |

### Previous canonical names

| Previous command | Status |
| --- | --- |
| `webex cc agent-personal-greeting-files create-v2-api` | Alias for `webex cc agent-personal-greeting-files create` |
| `webex cc agent-personal-greeting-files delete-id-v2-api` | Alias for `webex cc agent-personal-greeting-files delete-id` |
| `webex cc agent-personal-greeting-files delete-references-1` | Alias for `webex cc agent-personal-greeting-files delete-references` |
| `webex cc agent-personal-greeting-files get-id-v2-api` | Alias for `webex cc agent-personal-greeting-files get-id` |
| `webex cc agent-personal-greeting-files patch-id-v2-api` | Alias for `webex cc agent-personal-greeting-files patch-id` |
| `webex cc agent-personal-greeting-files update-id-v2-api` | Alias for `webex cc agent-personal-greeting-files update-id` |
| `webex cc contact-service-queue delete-csq-references` | Alias for `webex cc contact-service-queue delete-references` |
| `webex cc contact-service-queue list-csq-references-id` | Alias for `webex cc contact-service-queue list-references` |
| `webex cc contact-service-queue list-csqs-skills-profile` | Alias for `webex cc contact-service-queue list-by-dynamic-skills` |
| `webex cc contact-service-queue list-csqs-user-profile` | Alias for `webex cc contact-service-queue list-by-user-skill-profile` |
| `webex cc contact-service-queue list-manually-assignable-csqs` | Alias for `webex cc contact-service-queue list-manually-assignable` |
| `webex cc contact-service-queue list-skill-csqs-skill-profile` | Alias for `webex cc contact-service-queue list-by-skill-profile` |
| `webex cc outdial-ani bulk-save-entry` | Alias for `webex cc outdial-ani bulk-save-entries` |
| `webex cc outdial-ani list-entry` | Alias for `webex cc outdial-ani list-entries` |
| `webex cc outdial-ani list-entry-2` | Alias for `webex cc outdial-ani list-entries-for-ani` |
| `webex cc tasks resume` | Alias for `webex cc tasks unhold` |
| `webex cc tasks update-2` | Alias for `webex cc tasks append-message` |
| `webex cc users bulk-update-dynamic-skills` | Alias for `webex cc users bulk-partial-update-dynamic-skills` |
| `webex cc users get-ids` | Alias for `webex cc users list-by-ids` |

## meetings

### Added canonical names

| Full command | Method and path | Purpose |
| --- | --- | --- |
| `webex meetings meetings list-group` | `GET /group/meetings` | List group meetings with a specified meeting number or web link by a service app which has group meeting access |
| `webex meetings meetings patch-group` | `PATCH /group/meetings/{meetingId}` | Patches details for a group meeting with a specified meeting ID by a service app which has group meeting access |
| `webex meetings meetings update-group-control-status` | `POST /group/meetings/controls` | Update meeting recording control status by a service app which has group meeting access |
| `webex meetings recordings query` | `POST /recordings/query` | Queries recordings with filters in the request body |
| `webex meetings recordings query-admin-compliance-officer` | `POST /admin/recordings/query` | Queries recordings for an admin or compliance officer with filters in the request body |

## messaging

### Added canonical names

| Full command | Method and path | Purpose |
| --- | --- | --- |
| `webex messaging hds get-alarms-node` | `GET /hds/nodes/{nodeId}/alarms` | Returns the alarm details for a single HDS node for the provided time range (last 24 hours). |
| `webex messaging hds get-database-config-org` | `GET /hds/organizations/{organizationId}/database` | Retrieve details of database information for an HDS organization, such as database type and version used. |
| `webex messaging hds get-usage-node` | `GET /hds/nodes/{nodeId}/resourceUsage` | Retrieve CPU, memory, and disk resource usage details for a specific HDS node over the requested time range. |
| `webex messaging hds list-clusters-org` | `GET /hds/organizations/{organizationId}/clusters` | Retrieve a list of all clusters for a specific HDS organization, including cluster status, release channel, and upgrade schedule details. |
| `webex messaging hds list-nodes-cluster` | `GET /hds/clusters/{clusterId}/nodes` | Retrieve a list of all nodes for a specific HDS cluster, including availability, proxy details, deployment type, and release version. |
| `webex messaging hds list-tenants-org` | `GET /hds/organizations/{organizationId}/tenants` | Retrieve details of Multi-Tenant HDS organization such as Organization Name and ID, CMK state and state of Tenants Organizations. |

### Previous canonical names

| Previous command | Status |
| --- | --- |
| `webex messaging hds get-database-org-2` | Alias for `webex messaging hds get-database-config-org` |
| `webex messaging hds get-multi-tenant-org-2` | Alias for `webex messaging hds list-tenants-org` |

## Limitations

`cc functions import` lacks multipart upload handling. `cc flow import-legacy` lacks legacy upload body handling. Current-format flow import was live-tested; publish was not. Usage report types and completed-task variable updates were not tenant-tested.
