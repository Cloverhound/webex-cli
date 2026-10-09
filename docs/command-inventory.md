# Postman refresh inventory

Changes introduced by the Postman collection refresh. Added names include renamed/moved operations, not only new Cisco APIs. See [command-migration.md](command-migration.md) for migration details.

## Change classification

**178 added**, **82 modified**, **97 renamed**, **0 deleted**. Added counts include renamed commands. Renamed entries retain aliases; deleted entries would mean the CLI spelling is no longer available. This refresh has no deleted spellings after compatibility aliases are accounted for. Modified entries include method, route, parameters, body, headers, and documentation changes.

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

### Modified commands

- **Modified** `webex admin data-sources delete` — path parameters: updated `dataSourceId`; help/description updated
- **Modified** `webex admin data-sources get` — path parameters: updated `dataSourceId`; help/description updated
- **Modified** `webex admin data-sources get-all` — help/description updated
- **Modified** `webex admin data-sources get-schema` — help/description updated
- **Modified** `webex admin data-sources get-schemas` — help/description updated
- **Modified** `webex admin data-sources register` — help/description updated
- **Modified** `webex admin data-sources update` — help/description updated
- **Modified** `webex admin people update-person` — help/description updated
- **Modified** `webex admin recordings list` — help/description updated
- **Modified** `webex admin recordings list-admin-compliance-officer` — help/description updated
- **Modified** `webex admin recordings list-group` — query parameters: removed `topic`
- **Modified** `webex admin reports create` — body fields: added `timeZone`
- **Modified** `webex admin scim-2-users update-patch` — help/description updated
- **Modified** `webex admin scim-2-users update-put` — help/description updated


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
| `webex cc address-book delete` | `DELETE /organization/{orgid}/v3/address-book/{id}` | Delete an existing Address Book by ID in a given organization. |
| `webex cc address-book get` | `GET /organization/{orgid}/v3/address-book/{id}` | Retrieve an existing Address Book by ID in a given organization. |
| `webex cc address-book update` | `PUT /organization/{orgid}/v3/address-book/{id}` | Update an existing Address Book by ID in a given organization. |
| `webex cc agent-personal-greeting-files create` | `POST /organization/{orgid}/v2/agent-personal-greeting` | Create a new Greeting File in a given organization. |
| `webex cc agent-personal-greeting-files delete` | `DELETE /organization/{orgid}/v2/agent-personal-greeting/{id}` | Delete an existing Greeting File by ID in a given organization. |
| `webex cc agent-personal-greeting-files delete-references` | `POST /organization/{orgid}/agent-personal-greeting/delete-reference` | Removes all references of the specified agent (ciUserId) from greeting files in the given organization |
| `webex cc agent-personal-greeting-files get` | `GET /organization/{orgid}/v2/agent-personal-greeting/{id}` | Retrieve an existing Greeting File by ID in a given organization. |
| `webex cc agent-personal-greeting-files patch` | `PATCH /organization/{orgid}/v2/agent-personal-greeting/{id}` | Partially update Greeting File by ID in a given organization. |
| `webex cc agent-personal-greeting-files update` | `PUT /organization/{orgid}/v2/agent-personal-greeting/{id}` | Update an existing Greeting File by ID in a given organization. |
| `webex cc ai-feature get` | `GET /organization/{orgid}/ai-feature/{id}` | Retrieve an existing AI Feature resource by ID in a given organization. |
| `webex cc ai-feature patch` | `PATCH /organization/{orgid}/ai-feature/{id}` | Partially update AI Feature resource by ID in a given organization. |
| `webex cc asset bulk-save` | `POST /organization/{orgid}/asset/bulk` | Create, Update or delete Assets in bulk in a given organization |
| `webex cc asset create` | `POST /organization/{orgid}/asset` | Create a new Asset in a given organization |
| `webex cc asset delete` | `DELETE /organization/{orgid}/asset/{id}` | Delete an existing Asset by ID in a given organization |
| `webex cc asset get` | `GET /organization/{orgid}/asset/{id}` | Retrieve an existing Asset by ID in a given organization |
| `webex cc asset list` | `GET /organization/{orgid}/v2/asset` | Retrieve a list of Assets in a given organization |
| `webex cc asset list-references` | `GET /organization/{orgid}/asset/{id}/incoming-references` | Retrieve a list of all entities that have reference to an existing Asset by ID in a given organization |
| `webex cc asset patch` | `PATCH /organization/{orgid}/asset/{id}` | Partially update Asset by ID in a given organization |
| `webex cc asset update` | `PUT /organization/{orgid}/asset/{id}` | Update an existing Asset by ID in a given organization |
| `webex cc audio-files delete` | `DELETE /organization/{orgid}/audio-file/{id}` | Delete an existing Audio File by ID in a given organization. |
| `webex cc audio-files get` | `GET /organization/{orgid}/audio-file/{id}` | Retrieve an existing Audio File by ID in a given organization. |
| `webex cc audio-files patch` | `PATCH /organization/{orgid}/audio-file/{id}` | Partially update Audio File by ID in a given organization. |
| `webex cc audio-files update` | `PUT /organization/{orgid}/audio-file/{id}` | Update an existing Audio File by ID in a given organization. |
| `webex cc auto-csat get` | `GET /organization/{orgid}/auto-csat/{id}` | Retrieve an existing Auto CSAT resource by ID in a given organization |
| `webex cc auto-csat update` | `PUT /organization/{orgid}/auto-csat/{id}` | Update an existing Auto CSAT resource by ID in a given organization |
| `webex cc auxiliary-code delete` | `DELETE /organization/{orgid}/auxiliary-code/{id}` | Delete an existing Auxiliary Code by ID in a given organization. |
| `webex cc auxiliary-code get` | `GET /organization/{orgid}/auxiliary-code/{id}` | Retrieve an existing Auxiliary Code by ID in a given organization. |
| `webex cc auxiliary-code update` | `PUT /organization/{orgid}/auxiliary-code/{id}` | Update an existing Auxiliary Code by ID in a given organization. |
| `webex cc business-hour delete` | `DELETE /organization/{orgid}/business-hours/{id}` | Delete an existing Business Hours resource by ID in a given organization. |
| `webex cc business-hour get` | `GET /organization/{orgid}/business-hours/{id}` | Retrieve an existing Business Hours resource by ID in a given organization. |
| `webex cc business-hour update` | `PUT /organization/{orgid}/business-hours/{id}` | Update an existing Business Hours resource by ID in a given organization. |
| `webex cc call-monitoring whisper-coach-request` | `POST /v1/monitor/{interactionId}/coach` | This feature is currently in Beta |
| `webex cc campaign-group list` | `GET /v3/campaign-management/campaign-groups/{campaignGroupName}/campaigns` | Retrieves the list of campaigns that belong to the specified campaign group |
| `webex cc campaign-manager get-valid-times` | `GET /v1/organization/{orgId}/getValidCampaignTimes` | Gets valid campaign times for a campaign and agent |
| `webex cc channel bulk-save` | `POST /organization/{orgid}/channel/bulk` | Create, Update or delete Channels in bulk in a given organization. |
| `webex cc channel create` | `POST /organization/{orgid}/channel` | Create a new Channel in a given organization |
| `webex cc channel delete` | `DELETE /organization/{orgid}/channel/{id}` | Delete an existing Channel by ID in a given organization. |
| `webex cc channel get` | `GET /organization/{orgid}/channel/{id}` | Retrieve an existing Channel by ID in a given organization |
| `webex cc channel list` | `GET /organization/{orgid}/v2/channel` | Retrieve a list of Channels in a given organization |
| `webex cc channel list-references` | `GET /organization/{orgid}/channel/{id}/incoming-references` | Retrieve a list of all entities that have reference to an existing Channel by ID in a given organization |
| `webex cc channel patch` | `PATCH /organization/{orgid}/channel/{id}` | Partially update a channel by ID |
| `webex cc channel update` | `PUT /organization/{orgid}/channel/{id}` | Update an existing Channel by ID in a given organization. |
| `webex cc contact-list-management update-status-across-campaign-chain` | `PATCH /v3/campaign-management/campaigns/{campaignId}/contacts/{contactId}` | Synchronously closes the specified contacts and returns the outcome in the same response |
| `webex cc contact-number delete` | `DELETE /organization/{orgid}/contact-number/{id}` | Delete an existing Contact Number by ID in a given organization. |
| `webex cc contact-number get` | `GET /organization/{orgid}/contact-number/{id}` | Retrieve an existing Contact Number by ID in a given organization. |
| `webex cc contact-number update` | `PUT /organization/{orgid}/contact-number/{id}` | Update an existing Contact Number by ID in a given organization. |
| `webex cc contact-service-queue delete` | `DELETE /organization/{orgid}/contact-service-queue/{id}` | Delete an existing Contact Service Queue by ID in a given organization. |
| `webex cc contact-service-queue delete-references` | `POST /organization/{orgid}/contact-service-queue/delete-reference` | Removes the references to the specified entities (such as teams, sites, or agents) from Contact Service Queues for a given organization. |
| `webex cc contact-service-queue get` | `GET /organization/{orgid}/v2/contact-service-queue/{id}` | Retrieve an existing Contact Service Queue by ID in a given organization. |
| `webex cc contact-service-queue list-by-dynamic-skills` | `POST /organization/{orgid}/contact-service-queue/fetch-by-dynamic-skills-and-skillProfile` | Retrieve a list of skill-based Contact Service Queues that match the given dynamic skills and skill profile criteria in a given organization. |
| `webex cc contact-service-queue list-by-skill-profile` | `GET /organization/{orgid}/contact-service-queue/by-skill-profile-id/{id}` | Retrieve a list of skill-based Contact Service Queues associated with a given skill profile ID, accessible to authorized clients in a given organization. |
| `webex cc contact-service-queue list-by-user-skill-profile` | `POST /organization/{orgid}/contact-service-queue/fetch-by-userId-skillProfileId` | Retrieve a list of skill-based Contact Service Queues associated with the given skill profile ID and user ID combination in a given organization. |
| `webex cc contact-service-queue list-manually-assignable` | `POST /organization/{orgid}/contact-service-queue/fetch-manually-assignable-queues` | Retrieve a list of Contact Service Queues that are eligible for manual contact assignment based on the provided criteria in a given organization. |
| `webex cc contact-service-queue list-references` | `GET /organization/{orgid}/contact-service-queue/{id}/incoming-references` | Retrieve a list of all entities that have reference to an existing Contact Service Queue by ID in a given organization. |
| `webex cc contact-service-queue update` | `PUT /organization/{orgid}/v2/contact-service-queue/{id}` | Update an existing Contact Service Queue by ID in a given organization. |
| `webex cc desktop-layout delete` | `DELETE /organization/{orgid}/desktop-layout/{id}` | Delete an existing Desktop Layout by ID in a given organization. |
| `webex cc desktop-layout get` | `GET /organization/{orgid}/desktop-layout/{id}` | Retrieve an existing Desktop Layout by ID in a given organization. |
| `webex cc desktop-layout update` | `PUT /organization/{orgid}/desktop-layout/{id}` | Update an existing Desktop Layout by ID in a given organization. |
| `webex cc desktop-profile delete` | `DELETE /organization/{orgid}/agent-profile/{id}` | Delete an existing Desktop Profile by ID in a given organization. |
| `webex cc desktop-profile get` | `GET /organization/{orgid}/agent-profile/{id}` | Retrieve an existing Desktop Profile by ID in a given organization. |
| `webex cc desktop-profile update` | `PUT /organization/{orgid}/agent-profile/{id}` | Update an existing Desktop Profile by ID in a given organization. |
| `webex cc dial-plan delete` | `DELETE /organization/{orgid}/dial-plan/{id}` | Delete an existing Dial Plan by ID in a given organization. |
| `webex cc dial-plan get` | `GET /organization/{orgid}/dial-plan/{id}` | Retrieve an existing Dial Plan by ID in a given organization. |
| `webex cc dial-plan update` | `PUT /organization/{orgid}/dial-plan/{id}` | Update an existing Dial Plan by ID in a given organization. |
| `webex cc entry-point delete` | `DELETE /organization/{orgid}/entry-point/{id}` | Delete an existing Entry Point by ID in a given organization. |
| `webex cc entry-point get` | `GET /organization/{orgid}/entry-point/{id}` | Retrieve an existing Entry Point by ID in a given organization. |
| `webex cc entry-point update` | `PUT /organization/{orgid}/entry-point/{id}` | Update an existing Entry Point by ID in a given organization. |
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
| `webex cc generated-summaries get` | `GET /organization/{orgid}/generated-summaries/{id}` | Retrieve an existing Generated Summaries resource by ID in a given organization |
| `webex cc generated-summaries update` | `PUT /organization/{orgid}/generated-summaries/{id}` | Update an existing Generated Summaries resource by ID in a given organization |
| `webex cc global-variables delete` | `DELETE /organization/{orgid}/cad-variable/{id}` | Delete an existing Global Variable by ID in a given organization. |
| `webex cc global-variables get` | `GET /organization/{orgid}/cad-variable/{id}` | Retrieve an existing Global Variable by ID in a given organization. |
| `webex cc global-variables update` | `PUT /organization/{orgid}/cad-variable/{id}` | Update an existing Global Variable by ID in a given organization |
| `webex cc holiday-list delete` | `DELETE /organization/{orgid}/holiday-list/{id}` | Delete an existing Holiday List by ID in a given organization. |
| `webex cc holiday-list get` | `GET /organization/{orgid}/holiday-list/{id}` | Retrieve an existing Holiday List by ID in a given organization. |
| `webex cc holiday-list update` | `PUT /organization/{orgid}/holiday-list/{id}` | Update an existing Holiday List by ID in a given organization. |
| `webex cc multimedia-profile delete` | `DELETE /organization/{orgid}/multimedia-profile/{id}` | Delete an existing Multimedia Profile by ID in a given organization. |
| `webex cc multimedia-profile get` | `GET /organization/{orgid}/multimedia-profile/{id}` | Retrieve an existing Multimedia Profile by ID in a given organization. |
| `webex cc multimedia-profile update` | `PUT /organization/{orgid}/multimedia-profile/{id}` | Update an existing Multimedia Profile by ID in a given organization. |
| `webex cc outdial-ani bulk-save-entries` | `POST /organization/{orgid}/outdial-ani/{outDialAniId}/entry/bulk` | Create, Update or delete Outdial ANI Entries in bulk for an Address Book in a given organization. |
| `webex cc outdial-ani delete` | `DELETE /organization/{orgid}/outdial-ani/{id}` | Delete an existing Outdial ANI by ID in a given organization. |
| `webex cc outdial-ani get` | `GET /organization/{orgid}/outdial-ani/{id}` | Retrieve an existing Outdial ANI by ID in a given organization. |
| `webex cc outdial-ani list-entries` | `GET /organization/{orgid}/outdial-ani/entry` | Retrieve a list of Outdial ANI Entries in a given organization. |
| `webex cc outdial-ani list-entries-for-ani` | `GET /organization/{orgid}/v2/outdial-ani/{outDialAniId}/entry` | Retrieve a list of Outdial ANI Entries in a given organization. |
| `webex cc outdial-ani update` | `PUT /organization/{orgid}/outdial-ani/{id}` | Update an existing Outdial ANI by ID in a given organization. |
| `webex cc overrides delete` | `DELETE /organization/{orgid}/overrides/{id}` | Delete an existing Overrides resource by ID in a given organization. |
| `webex cc overrides get` | `GET /organization/{orgid}/overrides/{id}` | Retrieve an existing Overrides resource by ID in a given organization. |
| `webex cc overrides update` | `PUT /organization/{orgid}/overrides/{id}` | Update an existing Overrides resource by ID in a given organization. |
| `webex cc resource-collection delete` | `DELETE /organization/{orgid}/resource-collection/{id}` | Delete an existing resource collection by ID in a given organization. |
| `webex cc resource-collection get` | `GET /organization/{orgid}/resource-collection/{id}` | Retrieve an existing Resource Collection by ID in a given organization. |
| `webex cc resource-collection update` | `PUT /organization/{orgid}/resource-collection/{id}` | Update an existing resource collection by ID in a given organization. |
| `webex cc search-metadata get` | `GET /search/v2/meta` | Returns schema metadata for the Webex Contact Center Search GraphQL API, describing supported queries, fields, nested structures, data types, filters, sorting options, grouping capabilities, and aggregation operations. |
| `webex cc site delete` | `DELETE /organization/{orgid}/site/{id}` | Delete an existing Site by ID in a given organization. |
| `webex cc site get` | `GET /organization/{orgid}/site/{id}` | Retrieve an existing Site by ID in a given organization. |
| `webex cc site update` | `PUT /organization/{orgid}/site/{id}` | Update an existing Site by ID in a given organization. |
| `webex cc skill delete` | `DELETE /organization/{orgid}/skill/{id}` | Delete an existing Skill by ID in a given organization. |
| `webex cc skill get` | `GET /organization/{orgid}/skill/{id}` | Retrieve an existing Skill by ID in a given organization. |
| `webex cc skill update` | `PUT /organization/{orgid}/skill/{id}` | Update an existing Skill by ID in a given organization. |
| `webex cc skill-profile delete` | `DELETE /organization/{orgid}/skill-profile/{id}` | Delete an existing Skill Profile by ID in a given organization. |
| `webex cc skill-profile get` | `GET /organization/{orgid}/skill-profile/{id}` | Retrieve an existing Skill Profile by ID in a given organization. |
| `webex cc skill-profile update` | `PUT /organization/{orgid}/skill-profile/{id}` | Update an existing Skill Profile by ID in a given organization. |
| `webex cc tasks append-message` | `POST /v2/tasks/{taskId}/messages` | This feature is currently in Beta |
| `webex cc tasks drop-participant-conference` | `POST /v1/tasks/{taskId}/conference/participants/{participantId}/drop` | Access this endpoint when the user needs to drop a specific participant from an active conference associated with a task |
| `webex cc tasks pause-digital` | `POST /v1/tasks/{taskId}/pause` | Access this endpoint when users such as administrators, supervisors, or agents with an agent license need to pause a task that cannot be handled immediately |
| `webex cc tasks resume-digital` | `POST /v1/tasks/{taskId}/resume` | Access this endpoint when users such as administrators, supervisors, or agents with an agent license need to resume a previously paused task |
| `webex cc tasks unhold` | `POST /v1/tasks/{taskId}/unhold` | Access this endpoint when the user has to resume a call from hold |
| `webex cc team delete` | `DELETE /organization/{orgid}/team/{id}` | Delete an existing Team by ID in a given organization. |
| `webex cc team get` | `GET /organization/{orgid}/team/{id}` | Retrieve an existing Team by ID in a given organization. |
| `webex cc team update` | `PUT /organization/{orgid}/team/{id}` | Update an existing Team by ID in a given organization. |
| `webex cc templates get-flow` | `GET /templates/{id}` | Retrieve a specific flow template by its ID |
| `webex cc templates list-flow` | `GET /templates` | List available flow templates that can be used to create new flows. |
| `webex cc usage-reports create` | `POST /v1/usage-reports` | Creates a new usage report for the specified organization, resource type, |
| `webex cc usage-reports delete` | `DELETE /v1/usage-reports/{reportId}` | Deletes the specified usage report and its associated file |
| `webex cc usage-reports download-file` | `GET /v1/usage-reports/{fileId}/download` | Downloads one completed report file as a binary stream. |
| `webex cc usage-reports get` | `GET /v1/usage-reports/{reportId}` | Retrieves detailed information for a specific usage report. |
| `webex cc usage-reports list` | `GET /v1/usage-reports` | Retrieves a list of usage reports for the organization |
| `webex cc usage-reports list-resource-types` | `GET /v1/usage-reports/resource-types` | Returns the list of available resource types that can be used for report generation, along with available data dates. |
| `webex cc user-profiles delete` | `DELETE /organization/{orgid}/v3/user-profile/{id}` | Delete an existing user profile by ID in a given organization. |
| `webex cc user-profiles get` | `GET /organization/{orgid}/v3/user-profile/{id}` | Retrieve an existing user profile by ID in a given organization. |
| `webex cc user-profiles update` | `PUT /organization/{orgid}/v3/user-profile/{id}` | Update an existing user profile by ID in a given organization. |
| `webex cc users bulk-partial-update-dynamic-skills` | `PATCH /organization/{orgid}/user/bulk/update-dynamic-skill/{skillId}` | Assign or unassign a dynamic skill to/from multiple users in bulk for a given organization. |
| `webex cc users get` | `GET /organization/{orgid}/user/{id}` | Retrieve an existing Users by ID in a given organization. |
| `webex cc users list-by-ids` | `POST /organization/{orgid}/user/fetch-user-details-by-ids` | Retrieve an existing User's first name, last name and email by list of IDs in a given organization. |
| `webex cc users list-call-monitoring-id` | `GET /organization/{orgid}/user/by-call-monitoring-id/{id}` | Fetch paginated users associated to the selected call monitoring team filters while enforcing team ACL. |
| `webex cc users patch` | `PATCH /organization/{orgid}/user/{id}` | Partially update User by ID in a given organization. |
| `webex cc users update` | `PUT /organization/{orgid}/user/{id}` | Update an existing User by ID in a given organization. |
| `webex cc work-types delete` | `DELETE /organization/{orgid}/work-type/{id}` | Delete an existing Work Type by ID in a given organization. |
| `webex cc work-types get` | `GET /organization/{orgid}/work-type/{id}` | Retrieve an existing Work Type by ID in a given organization. |
| `webex cc work-types update` | `PUT /organization/{orgid}/work-type/{id}` | Update an existing Work Type by ID in a given organization. |

### Renamed or deleted canonical names

| Previous command | Status |
| --- | --- |
| `webex cc address-book delete-id` | Alias for `webex cc address-book delete` |
| `webex cc address-book get-id` | Alias for `webex cc address-book get` |
| `webex cc address-book update-id` | Alias for `webex cc address-book update` |
| `webex cc agent-personal-greeting-files create-v2-api` | Alias for `webex cc agent-personal-greeting-files create` |
| `webex cc agent-personal-greeting-files delete-id-v2-api` | Alias for `webex cc agent-personal-greeting-files delete` |
| `webex cc agent-personal-greeting-files delete-references-1` | Alias for `webex cc agent-personal-greeting-files delete-references` |
| `webex cc agent-personal-greeting-files get-id-v2-api` | Alias for `webex cc agent-personal-greeting-files get` |
| `webex cc agent-personal-greeting-files patch-id-v2-api` | Alias for `webex cc agent-personal-greeting-files patch` |
| `webex cc agent-personal-greeting-files update-id-v2-api` | Alias for `webex cc agent-personal-greeting-files update` |
| `webex cc ai-feature get-id` | Alias for `webex cc ai-feature get` |
| `webex cc ai-feature patch-id` | Alias for `webex cc ai-feature patch` |
| `webex cc audio-files delete-id` | Alias for `webex cc audio-files delete` |
| `webex cc audio-files get-id` | Alias for `webex cc audio-files get` |
| `webex cc audio-files patch-id` | Alias for `webex cc audio-files patch` |
| `webex cc audio-files update-id` | Alias for `webex cc audio-files update` |
| `webex cc auto-csat get-id` | Alias for `webex cc auto-csat get` |
| `webex cc auto-csat update-id` | Alias for `webex cc auto-csat update` |
| `webex cc auxiliary-code delete-id` | Alias for `webex cc auxiliary-code delete` |
| `webex cc auxiliary-code get-id` | Alias for `webex cc auxiliary-code get` |
| `webex cc auxiliary-code update-id` | Alias for `webex cc auxiliary-code update` |
| `webex cc business-hour delete-id` | Alias for `webex cc business-hour delete` |
| `webex cc business-hour get-id` | Alias for `webex cc business-hour get` |
| `webex cc business-hour update-id` | Alias for `webex cc business-hour update` |
| `webex cc contact-number delete-id` | Alias for `webex cc contact-number delete` |
| `webex cc contact-number get-id` | Alias for `webex cc contact-number get` |
| `webex cc contact-number update-id` | Alias for `webex cc contact-number update` |
| `webex cc contact-service-queue delete-csq-references` | Alias for `webex cc contact-service-queue delete-references` |
| `webex cc contact-service-queue delete-id` | Alias for `webex cc contact-service-queue delete` |
| `webex cc contact-service-queue get-id` | Alias for `webex cc contact-service-queue get` |
| `webex cc contact-service-queue list-csq-references-id` | Alias for `webex cc contact-service-queue list-references` |
| `webex cc contact-service-queue list-csqs-skills-profile` | Alias for `webex cc contact-service-queue list-by-dynamic-skills` |
| `webex cc contact-service-queue list-csqs-user-profile` | Alias for `webex cc contact-service-queue list-by-user-skill-profile` |
| `webex cc contact-service-queue list-manually-assignable-csqs` | Alias for `webex cc contact-service-queue list-manually-assignable` |
| `webex cc contact-service-queue list-skill-csqs-skill-profile` | Alias for `webex cc contact-service-queue list-by-skill-profile` |
| `webex cc contact-service-queue update-id` | Alias for `webex cc contact-service-queue update` |
| `webex cc desktop-layout delete-id` | Alias for `webex cc desktop-layout delete` |
| `webex cc desktop-layout get-id` | Alias for `webex cc desktop-layout get` |
| `webex cc desktop-layout update-id` | Alias for `webex cc desktop-layout update` |
| `webex cc desktop-profile delete-id` | Alias for `webex cc desktop-profile delete` |
| `webex cc desktop-profile get-id` | Alias for `webex cc desktop-profile get` |
| `webex cc desktop-profile update-id` | Alias for `webex cc desktop-profile update` |
| `webex cc dial-plan delete-id` | Alias for `webex cc dial-plan delete` |
| `webex cc dial-plan get-id` | Alias for `webex cc dial-plan get` |
| `webex cc dial-plan update-id` | Alias for `webex cc dial-plan update` |
| `webex cc entry-point delete-id` | Alias for `webex cc entry-point delete` |
| `webex cc entry-point get-id` | Alias for `webex cc entry-point get` |
| `webex cc entry-point update-id` | Alias for `webex cc entry-point update` |
| `webex cc generated-summaries get-id` | Alias for `webex cc generated-summaries get` |
| `webex cc generated-summaries update-id` | Alias for `webex cc generated-summaries update` |
| `webex cc global-variables delete-id` | Alias for `webex cc global-variables delete` |
| `webex cc global-variables get-id` | Alias for `webex cc global-variables get` |
| `webex cc global-variables update-id` | Alias for `webex cc global-variables update` |
| `webex cc holiday-list delete-id` | Alias for `webex cc holiday-list delete` |
| `webex cc holiday-list get-id` | Alias for `webex cc holiday-list get` |
| `webex cc holiday-list update-id` | Alias for `webex cc holiday-list update` |
| `webex cc multimedia-profile delete-id` | Alias for `webex cc multimedia-profile delete` |
| `webex cc multimedia-profile get-id` | Alias for `webex cc multimedia-profile get` |
| `webex cc multimedia-profile update-id` | Alias for `webex cc multimedia-profile update` |
| `webex cc outdial-ani bulk-save-entry` | Alias for `webex cc outdial-ani bulk-save-entries` |
| `webex cc outdial-ani delete-id` | Alias for `webex cc outdial-ani delete` |
| `webex cc outdial-ani get-id` | Alias for `webex cc outdial-ani get` |
| `webex cc outdial-ani list-entry` | Alias for `webex cc outdial-ani list-entries` |
| `webex cc outdial-ani list-entry-2` | Alias for `webex cc outdial-ani list-entries-for-ani` |
| `webex cc outdial-ani update-id` | Alias for `webex cc outdial-ani update` |
| `webex cc overrides delete-id` | Alias for `webex cc overrides delete` |
| `webex cc overrides get-id` | Alias for `webex cc overrides get` |
| `webex cc overrides update-id` | Alias for `webex cc overrides update` |
| `webex cc resource-collection delete-id` | Alias for `webex cc resource-collection delete` |
| `webex cc resource-collection get-id` | Alias for `webex cc resource-collection get` |
| `webex cc resource-collection update-id` | Alias for `webex cc resource-collection update` |
| `webex cc site delete-id` | Alias for `webex cc site delete` |
| `webex cc site get-id` | Alias for `webex cc site get` |
| `webex cc site update-id` | Alias for `webex cc site update` |
| `webex cc skill delete-id` | Alias for `webex cc skill delete` |
| `webex cc skill get-id` | Alias for `webex cc skill get` |
| `webex cc skill update-id` | Alias for `webex cc skill update` |
| `webex cc skill-profile delete-id` | Alias for `webex cc skill-profile delete` |
| `webex cc skill-profile get-id` | Alias for `webex cc skill-profile get` |
| `webex cc skill-profile update-id` | Alias for `webex cc skill-profile update` |
| `webex cc tasks resume` | Alias for `webex cc tasks unhold` |
| `webex cc tasks update-2` | Alias for `webex cc tasks append-message` |
| `webex cc team delete-id` | Alias for `webex cc team delete` |
| `webex cc team get-id` | Alias for `webex cc team get` |
| `webex cc team update-id` | Alias for `webex cc team update` |
| `webex cc user-profiles delete-id` | Alias for `webex cc user-profiles delete` |
| `webex cc user-profiles get-id` | Alias for `webex cc user-profiles get` |
| `webex cc user-profiles update-id` | Alias for `webex cc user-profiles update` |
| `webex cc users bulk-update-dynamic-skills` | Alias for `webex cc users bulk-partial-update-dynamic-skills` |
| `webex cc users get-id` | Alias for `webex cc users get` |
| `webex cc users get-ids` | Alias for `webex cc users list-by-ids` |
| `webex cc users patch-id` | Alias for `webex cc users patch` |
| `webex cc users update-id` | Alias for `webex cc users update` |
| `webex cc work-types delete-id` | Alias for `webex cc work-types delete` |
| `webex cc work-types get-id` | Alias for `webex cc work-types get` |
| `webex cc work-types update-id` | Alias for `webex cc work-types update` |

### Modified commands

- **Modified** `webex cc agent-personal-greeting-files list` — query parameters: updated `filter`, `includeAgentDetails`; help/description updated
- **Modified** `webex cc agent-summaries list-2` — body fields: removed `agentCiUserId`
- **Modified** `webex cc business-hour bulk-save` — help/description updated
- **Modified** `webex cc business-hour list` — query parameters: updated `includeCount`, `singleObjectResponse`; help/description updated
- **Modified** `webex cc campaign-manager update-request` — JSON body handling changed
- **Modified** `webex cc captures list` — help/description updated
- **Modified** `webex cc contact-list-management get-within-campaign` — path: `/v3/campaign-management/campaigns/{campaignId}/contact-lists` → `/v4/campaign-management/campaigns/{campaignId}/contact-lists`; query parameters: updated `source`; help/description updated
- **Modified** `webex cc contact-list-management update-status-within` — path parameters: updated `contactId`
- **Modified** `webex cc contact-service-queue bulk-partial-update` — help/description updated
- **Modified** `webex cc contact-service-queue bulk-save` — help/description updated
- **Modified** `webex cc contact-service-queue create-remove-agents-users-agent` — help/description updated
- **Modified** `webex cc contact-service-queue list` — query parameters: updated `provisioningView`, `singleObjectResponse`; help/description updated
- **Modified** `webex cc contact-service-queue list-agent-based` — help/description updated
- **Modified** `webex cc contact-service-queue list-skill-based` — help/description updated
- **Modified** `webex cc contact-service-queue list-team-based` — help/description updated
- **Modified** `webex cc contact-service-queue purge-inactive` — query parameters: updated `nextStartId`; help/description updated
- **Modified** `webex cc data-sources delete` — path parameters: updated `dataSourceId`; help/description updated
- **Modified** `webex cc data-sources get` — path parameters: updated `dataSourceId`; help/description updated
- **Modified** `webex cc data-sources get-all` — help/description updated
- **Modified** `webex cc data-sources get-schema` — help/description updated
- **Modified** `webex cc data-sources get-schemas` — help/description updated
- **Modified** `webex cc data-sources register` — help/description updated
- **Modified** `webex cc data-sources update` — help/description updated
- **Modified** `webex cc desktop-layout list` — help/description updated
- **Modified** `webex cc desktop-profile bulk-save` — help/description updated
- **Modified** `webex cc desktop-profile create` — query parameters: added `accessBuddyTeam`, `accessEntryPoint`, `accessIdleCode`, `accessQueue`, `accessWrapUpCode`, `active`, `addressBookId`, `agentAvailableAfterOutdial`, `agentDNValidation`, `agentDNValidationCriteria`, `agentDNValidationCriterions`, `agentPersonalGreeting`, `allowAutoWrapUpExtension`, `autoAcceptDigitalInteractions`, `autoAnswer`, `autoWrapAfterSeconds`, `autoWrapUp`, `buddyTeams`, `consultToQueue`, `createdTime`, `description`, `dialPlanEnabled`, `dialPlans`, `entryPoints`, `id`, `idleCodes`, `lastAgentRouting`, `lastUpdatedTime`, `loginVoiceOptions`, `manageChannelAvailability`, `name`, `organizationId`, `outdialANIId`, `outdialEnabled`, `outdialEntryPointId`, `parentType`, `queues`, `scheduleAndManageCallBack`, `screenPopup`, `showUserDetailsMS`, `showUserDetailsWebex`, `siteId`, `stateSynchronizationMS`, `stateSynchronizationWebex`, `systemDefault`, `thresholdRules`, `timeoutDesktopInactivityCustomEnabled`, `timeoutDesktopInactivityMins`, `version`, `viewableStatistics`, `wrapUpCodes`
- **Modified** `webex cc desktop-profile list` — query parameters: added `provisioningView`; updated `attributes`, `filter`, `singleObjectResponse`; help/description updated
- **Modified** `webex cc desktop-profile purge-inactive` — query parameters: updated `nextStartId`; help/description updated
- **Modified** `webex cc dial-plan bulk-save` — help/description updated
- **Modified** `webex cc dial-plan create` — body fields: added `createdTime`, `lastUpdatedTime`; help/description updated
- **Modified** `webex cc dial-plan list` — query parameters: updated `attributes`, `filter`; help/description updated
- **Modified** `webex cc dial-plan list-references` — help/description updated
- **Modified** `webex cc flow export` — path: `/flow-store/{orgId}/project/{projectId}/flows/{flowId}:export` → `/{orgId}/project/{projectId}/v2/flows/{flowId}:export`; query parameters: added `flowType`; updated `version`; documented response code changed; help/description updated
- **Modified** `webex cc flow import` — path: `/flow-store/{orgId}/project/{projectId}/flows:import` → `/{orgId}/project/{projectId}/v2/flows:import`; query parameters: updated `overwrite`; body fields: added `contactType`, `description`, `flowName`, `flowType`, `status`, `version`; headers: removed `Content-Length`; request body support changed; JSON body handling changed; documented response code changed; help/description updated
- **Modified** `webex cc flow list` — path: `/flow-store/{orgId}/project/{projectId}/flows` → `/{orgId}/project/{projectId}/flows`; query parameters: added `isValidation`, `searchBy`; updated `includePagination`; documented response code changed; help/description updated
- **Modified** `webex cc flow publish` — path: `/flow-store/{orgId}/project/{projectId}/flows/{flowId}:publish` → `/{orgId}/project/{projectId}/flows/{flowId}:publish`; query parameters: added `flowType`, `skipValidation`; documented response code changed; help/description updated
- **Modified** `webex cc holiday-list bulk-save` — help/description updated
- **Modified** `webex cc holiday-list list` — query parameters: updated `includeCount`, `singleObjectResponse`; help/description updated
- **Modified** `webex cc notification subscribe` — help/description updated
- **Modified** `webex cc outdial-ani bulk-save` — help/description updated
- **Modified** `webex cc outdial-ani list` — query parameters: updated `singleObjectResponse`; help/description updated
- **Modified** `webex cc overrides list` — help/description updated
- **Modified** `webex cc realtime subscribe-notification` — help/description updated
- **Modified** `webex cc skill bulk-save` — help/description updated
- **Modified** `webex cc skill create` — query parameters: removed `active`, `createdTime`, `description`, `dynamicSkill`, `enumSkillValues`, `id`, `lastUpdatedTime`, `name`, `organizationId`, `serviceLevelThreshold`, `skillType`, `version`; body fields: added `systemDefault`
- **Modified** `webex cc skill list` — query parameters: added `include`; updated `singleObjectResponse`; help/description updated
- **Modified** `webex cc skill purge-inactive` — query parameters: updated `nextStartId`; help/description updated
- **Modified** `webex cc tasks create` — body fields: added `channelType`; removed `eventTime`, `mediaMgr`, `mediaType`, `orgId`, `trackingId`; help/description updated
- **Modified** `webex cc team bulk-save` — help/description updated
- **Modified** `webex cc team list` — query parameters: updated `filter`, `provisioningView`, `singleObjectResponse`, `supervisorView`; help/description updated
- **Modified** `webex cc team purge-inactive` — query parameters: updated `nextStartId`; help/description updated
- **Modified** `webex cc user-profiles list` — help/description updated
- **Modified** `webex cc users get-agents-matching-skill-requirements` — body fields: removed `condition`, `createdTime`, `dynamicSkill`, `id`, `lastUpdatedTime`, `organizationId`, `skillId`, `skillName`, `skillType`, `skillValue`, `version`, `weight`; JSON body handling changed
- **Modified** `webex cc users get-along-profile-id` — help/description updated
- **Modified** `webex cc users get-ci-id` — query parameters: added `includeSkillDetails`
- **Modified** `webex cc users list` — query parameters: updated `includeAIMappingCount`, `queueId`, `singleObjectResponse`; help/description updated
- **Modified** `webex cc users list-along-profile` — help/description updated


## meetings

### Added canonical names

| Full command | Method and path | Purpose |
| --- | --- | --- |
| `webex meetings meetings list-group` | `GET /group/meetings` | List group meetings with a specified meeting number or web link by a service app which has group meeting access |
| `webex meetings meetings patch-group` | `PATCH /group/meetings/{meetingId}` | Patches details for a group meeting with a specified meeting ID by a service app which has group meeting access |
| `webex meetings meetings update-group-control-status` | `POST /group/meetings/controls` | Update meeting recording control status by a service app which has group meeting access |
| `webex meetings recordings query` | `POST /recordings/query` | Queries recordings with filters in the request body |
| `webex meetings recordings query-admin-compliance-officer` | `POST /admin/recordings/query` | Queries recordings for an admin or compliance officer with filters in the request body |

### Modified commands

- **Modified** `webex meetings meetings delete` — help/description updated
- **Modified** `webex meetings meetings patch` — help/description updated
- **Modified** `webex meetings meetings update` — help/description updated
- **Modified** `webex meetings people update-person` — help/description updated
- **Modified** `webex meetings recordings list` — help/description updated
- **Modified** `webex meetings recordings list-admin-compliance-officer` — help/description updated
- **Modified** `webex meetings recordings list-group` — query parameters: removed `topic`


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

### Renamed or deleted canonical names

| Previous command | Status |
| --- | --- |
| `webex messaging hds get-database-org-2` | Alias for `webex messaging hds get-database-config-org` |
| `webex messaging hds get-multi-tenant-org-2` | Alias for `webex messaging hds list-tenants-org` |

### Modified commands

- **Modified** `webex messaging hds get-org` — help/description updated
- **Modified** `webex messaging people update-person` — help/description updated


## device

### Modified commands

- **Modified** `webex device device-call get-layout-id` — help/description updated
- **Modified** `webex device device-call update-layout-id` — help/description updated

## Limitations

`cc functions import` lacks multipart upload handling. `cc flow import-legacy` lacks legacy upload body handling. Use current-format `flow import --body-file` for typed flow documents; legacy imports require an upload implementation.
