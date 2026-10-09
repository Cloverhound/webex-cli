# Renamed CLI commands in the Postman refresh

Every row gives the full previous command, its current canonical command, what the operation does, and whether the previous spelling still invokes that same operation. Append `--help` to any current command to see its arguments and flags.

**Previous-name source:** “PR #10” means the automatic refresh before our naming overrides. “Before refresh” means an older spelling retained for compatibility. These are different baselines: in particular, `webex cc tasks resume` meant unhold before the refresh, but the unmodified PR assigned it to digital-task resume.

**Flow formats:** `flow import`/`export` use the current typed flow document; `import-legacy`/`export-legacy` use deprecated raw FDL. This intentionally changes the format of the pre-refresh unsuffixed commands. Current flow, activity, event, and template commands use unprefixed URLs. Duplicate prefixed commands and the temporary `-direct` names have been removed.

The refresh tables below contain **94 previous-to-current mappings**, followed by **88 ID-suffix mappings**. Commands whose names and groups did not change are omitted.

## calling call-controls-for-me

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex calling call-controls dial-2` | PR #10 | `webex calling call-controls-for-me dial` | Initiate an outbound call to a specified destination. | No — use the new command |
| `webex calling call-controls answer-2` | PR #10 | `webex calling call-controls-for-me answer` | Answer an incoming call. | No — use the new command |
| `webex calling call-controls hangup-2` | PR #10 | `webex calling call-controls-for-me hangup` | Hangup a call. | No — use the new command |
| `webex calling call-controls list-2` | PR #10 | `webex calling call-controls-for-me list` | Get the list of details for all active calls associated with the user. | No — use the new command |
| `webex calling call-controls get-3` | PR #10 | `webex calling call-controls-for-me get` | Get the details of the specified active call for the user. | No — use the new command |

## calling metrics

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex admin calling-metrics webex-call-quality-stats` | PR #10 | `webex calling metrics get-call-quality-stats` | Retrieve aggregated Calling call and media-quality statistics for the authenticated organization. Requires Pro Pack; query up to seven days within the previous three weeks. | No — use the new command |

## cc activities

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc activities list-definitions-2` | PR #10 | `webex cc activities list-definitions` | List all available activity types with their full schemas, including inputs, outputs, ports, and descriptions. | No — use the new command |
| `webex cc activities describe-2` | PR #10 | `webex cc activities describe` | Retrieve a single activity's full definition, including inputs, outputs, named ports, and the JSON Schema for its inputs — all returned inline. | No — use the new command |
| `webex cc activities get-input-choices-2` | PR #10 | `webex cc activities get-input-choices` | Resolve valid values for an activity input, such as queues, audio files, or agents; supports search, validation, and cascading parent inputs. | No — use the new command |

## cc agent-personal-greeting-files

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc agent-personal-greeting-files create-v2-api` | Before refresh | `webex cc agent-personal-greeting-files create` | Create a new Greeting File in a given organization. | Yes — alias |
| `webex cc agent-personal-greeting-files delete-references-1` | Before refresh | `webex cc agent-personal-greeting-files delete-references` | Removes all references of the specified agent (ciUserId) from greeting files in the given organization. | Yes — alias |
| `webex cc agent-personal-greeting-files get-id-v2-api` | Before refresh | `webex cc agent-personal-greeting-files get` | Retrieve an existing Greeting File by ID in a given organization. | Yes — alias |
| `webex cc agent-personal-greeting-files update-id-v2-api` | Before refresh | `webex cc agent-personal-greeting-files update` | Update an existing Greeting File by ID in a given organization. | Yes — alias |
| `webex cc agent-personal-greeting-files delete-id-v2-api` | Before refresh | `webex cc agent-personal-greeting-files delete` | Delete an existing Greeting File by ID in a given organization. | Yes — alias |
| `webex cc agent-personal-greeting-files patch-id-v2-api` | Before refresh | `webex cc agent-personal-greeting-files patch` | Partially update Greeting File by ID in a given organization. | Yes — alias |

## cc contact-service-queue

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc contact-service-queue list-skill-based-skill-profile-id` | PR #10 | `webex cc contact-service-queue list-by-skill-profile` | Retrieve a list of skill-based Contact Service Queues associated with a given skill profile ID, accessible to authorized clients in a given organization. | No — use the new command |
| `webex cc contact-service-queue list-skill-csqs-skill-profile` | Before refresh | `webex cc contact-service-queue list-by-skill-profile` | Retrieve a list of skill-based Contact Service Queues associated with a given skill profile ID, accessible to authorized clients in a given organization. | Yes — alias |
| `webex cc contact-service-queue delete-csq-references` | Before refresh | `webex cc contact-service-queue delete-references` | Removes the references to the specified entities (such as teams, sites, or agents) from Contact Service Queues for a given organization. | Yes — alias |
| `webex cc contact-service-queue list-manually-assignable-csqs` | Before refresh | `webex cc contact-service-queue list-manually-assignable` | Retrieve a list of Contact Service Queues that are eligible for manual contact assignment based on the provided criteria in a given organization. | Yes — alias |
| `webex cc contact-service-queue list-csq-references-id` | Before refresh | `webex cc contact-service-queue list-references` | Retrieve a list of all entities that have reference to an existing Contact Service Queue by ID in a given organization. | Yes — alias |
| `webex cc contact-service-queue list-agent-based-user-id` | PR #10 | `webex cc contact-service-queue list-agent-based` | Retrieve a list of agent-based Contact Service Queues by user ID in a given organization. | No — use the new command |
| `webex cc contact-service-queue list-skill-based-user-id` | PR #10 | `webex cc contact-service-queue list-skill-based` | Retrieve a list of skill-based Contact Service Queues by user ID in a given organization. | No — use the new command |
| `webex cc contact-service-queue list-team-based-user-id` | PR #10 | `webex cc contact-service-queue list-team-based` | Retrieve a list of team-based Contact Service Queues by user ID in a given organization. | No — use the new command |
| `webex cc contact-service-queue list-skill-based-dynamic-skills-skill-profile` | PR #10 | `webex cc contact-service-queue list-by-dynamic-skills` | Retrieve a list of skill-based Contact Service Queues that match the given dynamic skills and skill profile criteria in a given organization. | No — use the new command |
| `webex cc contact-service-queue list-csqs-skills-profile` | Before refresh | `webex cc contact-service-queue list-by-dynamic-skills` | Retrieve a list of skill-based Contact Service Queues that match the given dynamic skills and skill profile criteria in a given organization. | Yes — alias |
| `webex cc contact-service-queue list-skill-based-skill-profile-id-user-id` | PR #10 | `webex cc contact-service-queue list-by-user-skill-profile` | Retrieve a list of skill-based Contact Service Queues associated with the given skill profile ID and user ID combination in a given organization. | No — use the new command |
| `webex cc contact-service-queue list-csqs-user-profile` | Before refresh | `webex cc contact-service-queue list-by-user-skill-profile` | Retrieve a list of skill-based Contact Service Queues associated with the given skill profile ID and user ID combination in a given organization. | Yes — alias |

## cc events

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc events list-specifications-2` | PR #10 | `webex cc events list-specifications` | List the event types available for use in `event_flows[]`. | No — use the new command |

## cc flow

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc flow export` | PR #10 | `webex cc flow export-legacy` | Export a flow or subflow in legacy raw FDL format. This API is deprecated; the existing command keeps its original route. | **No — now resolves to a different operation; use the new command.** |
| `webex cc flow import` | PR #10 | `webex cc flow import-legacy` | Import a flow or subflow in legacy raw FDL format. This API is deprecated; the existing command keeps its original route. | **No — now resolves to a different operation; use the new command.** |
| `webex cc flows list` | PR #10 | `webex cc flow list` | List flows or subflows, with name, ID, channel, and pagination filters. | Yes — alias |
| `webex cc flows publish` | PR #10 | `webex cc flow publish` | Publish a flow or subflow. | Yes — alias |
| `webex cc flows lock` | PR #10 | `webex cc flow lock` | Lock a flow or subflow against concurrent edits. Locks expire after 15 minutes of inactivity. | Yes — alias |
| `webex cc flows unlock` | PR #10 | `webex cc flow unlock` | Release a flow or subflow edit lock. | Yes — alias |
| `webex cc flows validate` | PR #10 | `webex cc flow validate` | Validate a supplied flow definition without saving it. Passing validation does not guarantee import will succeed. | Yes — alias |
| `webex cc flows import` | PR #10 | `webex cc flow import` | Import a v2 flow definition as a new draft and return the assigned flow metadata. | Yes — alias |
| `webex cc flows get` | PR #10 | `webex cc flow get` | Retrieve the current flow draft. Use export to retrieve a published version. | Yes — alias |
| `webex cc flows save-draft` | PR #10 | `webex cc flow save-draft` | Replace the complete flow draft. Optional expectedVersion prevents overwriting a concurrently changed draft. | Yes — alias |
| `webex cc flows patch-draft` | PR #10 | `webex cc flow patch-draft` | Apply a partial update to a flow draft; validate the merged result before saving it. | Yes — alias |
| `webex cc flows validate-draft` | PR #10 | `webex cc flow validate-draft` | Validate the stored flow draft without changing it. | Yes — alias |
| `webex cc flows export` | PR #10 | `webex cc flow export` | Export a v2 flow document for backup, migration, or version control. | Yes — alias |
| `webex cc flows list-2` | PR #10 | `webex cc flow list` | List flows or subflows, with name, ID, channel, and pagination filters. | No — use the new command |
| `webex cc flows publish-2` | PR #10 | `webex cc flow publish` | Publish a flow or subflow. | No — use the new command |
| `webex cc flows lock-2` | PR #10 | `webex cc flow lock` | Lock a flow or subflow against concurrent edits. Locks expire after 15 minutes of inactivity. | No — use the new command |
| `webex cc flows unlock-2` | PR #10 | `webex cc flow unlock` | Release a flow or subflow edit lock. | No — use the new command |
| `webex cc flows validate-2` | PR #10 | `webex cc flow validate` | Validate a supplied flow definition without saving it. Passing validation does not guarantee import will succeed. | No — use the new command |
| `webex cc flows import-2` | PR #10 | `webex cc flow import` | Import a v2 flow definition as a new draft and return the assigned flow metadata. | No — use the new command |
| `webex cc flows get-2` | PR #10 | `webex cc flow get` | Retrieve the current flow draft. Use export to retrieve a published version. | No — use the new command |
| `webex cc flows save-draft-2` | PR #10 | `webex cc flow save-draft` | Replace the complete flow draft. Optional expectedVersion prevents overwriting a concurrently changed draft. | No — use the new command |
| `webex cc flows patch-draft-2` | PR #10 | `webex cc flow patch-draft` | Apply a partial update to a flow draft; validate the merged result before saving it. | No — use the new command |
| `webex cc flows validate-draft-2` | PR #10 | `webex cc flow validate-draft` | Validate the stored flow draft without changing it. | No — use the new command |
| `webex cc flows export-2` | PR #10 | `webex cc flow export` | Export a v2 flow document for backup, migration, or version control. | No — use the new command |
| `webex cc flows search` | PR #10 | `webex cc flow search` | Search flows or subflows using a case-sensitive query. | No — use the new command |
| `webex cc flows delete` | PR #10 | `webex cc flow delete` | Permanently delete a flow or subflow; deletion cannot be undone. | No — use the new command |
| `webex cc legacy-flows export` | PR #10 | `webex cc flow export-legacy` | Export a flow or subflow in legacy raw FDL format. This API is deprecated; the existing command keeps its original route. | No — use the new command |
| `webex cc legacy-flows import` | PR #10 | `webex cc flow import-legacy` | Import a flow or subflow in legacy raw FDL format. This API is deprecated; the existing command keeps its original route. | No — use the new command |

## cc functions

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc functions list-custom` | PR #10 | `webex cc functions list` | List or search custom functions in the organization. | No — use the new command |
| `webex cc functions create-custom` | PR #10 | `webex cc functions create` | Create a new custom function. | No — use the new command |
| `webex cc functions get-custom` | PR #10 | `webex cc functions get` | Retrieve a function draft, or a specific published version/tag. | No — use the new command |
| `webex cc functions update-custom` | PR #10 | `webex cc functions update` | Update an existing custom function by ID. | No — use the new command |
| `webex cc functions delete-custom` | PR #10 | `webex cc functions delete` | Delete a custom function by ID. | No — use the new command |
| `webex cc functions import-custom` | PR #10 | `webex cc functions import` | Import a function-definition JSON file as a multipart upload. The current generated CLI command still lacks file-upload handling; naming it does not fix that limitation. | No — use the new command |
| `webex cc functions unlock-custom` | PR #10 | `webex cc functions unlock` | Release the edit lock on a custom function so that other users can edit it. | No — use the new command |
| `webex cc functions lock-custom` | PR #10 | `webex cc functions lock` | Acquire an edit lock on a custom function to prevent concurrent writes by other users. | No — use the new command |
| `webex cc functions test-custom` | PR #10 | `webex cc functions test` | Execute a function with a test payload. The API implicitly publishes the latest draft before running it. | No — use the new command |
| `webex cc functions publish-custom` | PR #10 | `webex cc functions publish` | Publish the latest function draft under one or more tags: Dev, Test, Latest, or Live. | No — use the new command |
| `webex cc functions export-custom` | PR #10 | `webex cc functions export` | Export a selected function version or publish tag as function-definition JSON, including source code, inputs, and outputs. | No — use the new command |

## cc outdial-ani

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc outdial-ani list-entry` | Before refresh | `webex cc outdial-ani list-entries` | Retrieve a list of Outdial ANI Entries in a given organization. | Yes — alias |
| `webex cc outdial-ani bulk-save-entry` | Before refresh | `webex cc outdial-ani bulk-save-entries` | Create, Update or delete Outdial ANI Entries in bulk for an Address Book in a given organization. | Yes — alias |
| `webex cc outdial-ani list-entries-2` | PR #10 | `webex cc outdial-ani list-entries-for-ani` | List outbound caller-ID (ANI) entries belonging to a specified outdial ANI ID. | No — use the new command |
| `webex cc outdial-ani list-entry-2` | Before refresh | `webex cc outdial-ani list-entries-for-ani` | List outbound caller-ID (ANI) entries belonging to a specified outdial ANI ID. | Yes — alias |

## cc tasks

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc tasks resume` | Before refresh | `webex cc tasks unhold` | Resume a voice call that is on hold. This preserves the released meaning of `tasks resume`. | Yes — alias |
| `webex cc tasks update-2` | Before refresh | `webex cc tasks append-message` | Append a message to an existing task using the v2 messages endpoint (beta). | Yes — alias |
| `webex cc tasks pause` | PR #10 | `webex cc tasks pause-digital` | Pause a non-real-time digital task that cannot be handled immediately. | No — use the new command |
| `webex cc tasks resume` | PR #10 | `webex cc tasks resume-digital` | Resume a paused non-real-time digital task, such as email or social messaging. | **No — now resolves to a different operation; use the new command.** |

## cc templates

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc templates list-flow-2` | PR #10 | `webex cc templates list-flow` | List available flow templates that can be used to create new flows. | No — use the new command |
| `webex cc templates get-flow-2` | PR #10 | `webex cc templates get-flow` | Retrieve a specific flow template by its ID. | No — use the new command |

## cc usage-reports

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc usage-reports get-available-types` | PR #10 | `webex cc usage-reports list-resource-types` | Returns the list of available resource types that can be used for report generation, along with available data dates. | No — use the new command |

## cc users

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex cc users list-2` | PR #10 | `webex cc users list-by-ids` | Retrieve an existing User's first name, last name and email by list of IDs in a given organization. | No — use the new command |
| `webex cc users get-ids` | Before refresh | `webex cc users list-by-ids` | Retrieve an existing User's first name, last name and email by list of IDs in a given organization. | Yes — alias |
| `webex cc users bulk-update-dynamic-skills` | Before refresh | `webex cc users bulk-partial-update-dynamic-skills` | Assign or unassign a dynamic skill to/from multiple users in bulk for a given organization. | Yes — alias |

## messaging hds

| Previous full command | Source | Current full command | What it does | Previous spelling still works? |
| --- | --- | --- | --- | --- |
| `webex messaging hds get-database-org-2` | PR #10 | `webex messaging hds get-database-config-org` | Retrieve details of database information for an HDS organization, such as database type and version used. | Yes — alias |
| `webex messaging hds get-multi-tenant-org-2` | PR #10 | `webex messaging hds list-tenants-org` | Retrieve tenant organization names/IDs, tenant states, and customer-managed-key state for a multi-tenant HDS organization. | Yes — alias |
| `webex messaging hybrid-data-security get-org` | PR #10 | `webex messaging hds get-org` | Retrieve details for an Hybrid Data Security organization, such as the organization name, type of organization. | Yes — alias |
| `webex messaging hybrid-data-security list-clusters-org` | PR #10 | `webex messaging hds list-clusters-org` | Retrieve a list of all clusters for a specific Hybrid Data Security organization, including cluster status, release channel, and upgrade schedule details. | Yes — alias |
| `webex messaging hybrid-data-security get-cluster` | PR #10 | `webex messaging hds get-cluster` | Retrieve details for a specific Hybrid Data Security cluster, such as the cluster name, cluster status, upgrade schedule, and Hybrid Data Security nodes in the cluster. | Yes — alias |
| `webex messaging hybrid-data-security list-nodes-cluster` | PR #10 | `webex messaging hds list-nodes-cluster` | Retrieve a list of all nodes for a specific Hybrid Data Security cluster, including availability, proxy details, deployment type, and release version. | Yes — alias |
| `webex messaging hybrid-data-security get-node` | PR #10 | `webex messaging hds get-node` | Retrieve details for a specific Hybrid Data Security node, such as host name, release version, proxy details, deployment and build type, availability details, etc. | Yes — alias |
| `webex messaging hybrid-data-security get-database-org` | PR #10 | `webex messaging hds get-database-config-org` | Retrieve details of database information for an Hybrid Data Security organization, such as database type and version used. | **No — now resolves to a different operation; use the new command.** |
| `webex messaging hybrid-data-security get-multi-tenant-org` | PR #10 | `webex messaging hds list-tenants-org` | Retrieve tenant organization names/IDs, tenant states, and customer-managed-key state for a multi-tenant HDS organization. | **No — now resolves to a different operation; use the new command.** |
| `webex messaging hybrid-data-security get-alarms-node` | PR #10 | `webex messaging hds get-alarms-node` | Returns the alarm details for a single Hybrid Data Security node for the provided time range (last 24 hours). | Yes — alias |
| `webex messaging hybrid-data-security get-test-results-node` | PR #10 | `webex messaging hds get-test-results-node` | Retrieve the latest node network-test results, including bandwidth, DNS resolution, and HTTPS connectivity. | Yes — alias |
| `webex messaging hybrid-data-security get-usage-node` | PR #10 | `webex messaging hds get-usage-node` | Retrieve CPU, memory, and disk resource usage details for a specific Hybrid Data Security node over the requested time range. | Yes — alias |
| `webex messaging hybrid-data-security get-availability-cluster` | PR #10 | `webex messaging hds get-availability-cluster` | Get the availability details for an Hybrid Data Security cluster, where each segment specifies the start and end times, as well as the number of online, offline, and total nodes within that segment. | Yes — alias |

## Cloverhound read-only route checks — 2026-09-29

Tested with the Cloverhound organization on `https://api.wxcc-us1.cisco.com` using authenticated GET requests only. No import, publish, lock, save, patch, or delete requests were executed.

| Operation | With `/flow-store/` | Without `/flow-store/` | Comparison |
| --- | --- | --- | --- |
| List flows | 200 | 200 | Identical parsed JSON |
| Get one flow draft | 200 | 200 | Identical parsed JSON |
| Export that flow using the current format | 200 | 200 | Identical parsed JSON |
| List activity definitions | 200 | 200 | Identical parsed JSON |
| List event specifications | 200 | 200 | Identical parsed JSON |
| List templates | 404 at the org/project-scoped v2 template route | 200 at `/templates` | Different routing behavior |

Legacy flow export also returned 200 and an FDL document. For this organization and the tested reads, both flow URL forms work on the same production host. This does not establish equivalence for all operations or tenants. The subsequent draft write checks below support consolidating current operations onto unprefixed routes. `webex cc templates list-flow` now uses `/templates`. Activity input choices retain the unprefixed request’s search/validate/parent-input parameters.

## Cloverhound disposable-flow write checks — 2026-09-29 (Denver)

Created draft `CLI_Route_Test_ce61ca68ca8b` (ID `6abc548b79b443250d35110c`) using the unprefixed import URL with `overwrite=false`. The flow contained only a start activity and a disconnect activity. It was never published or assigned to an entry point/routing strategy.

| Check | Result |
| --- | --- |
| Validate through unprefixed route | Valid; only a recommendation to add activity descriptions |
| Import through unprefixed route | 201; new unassigned draft |
| Read new draft through `/flow-store/` route | 200; correct new flow |
| Lock prefixed, unlock unprefixed | Both 200 |
| Lock unprefixed, unlock prefixed | Both 200 |
| Patch description prefixed, read unprefixed | Both 200; updated description verified |
| Patch description unprefixed, read prefixed | Both 200; updated description verified |
| Save complete draft prefixed, read unprefixed | Both 200; saved description verified |
| Save complete draft unprefixed, read prefixed | Both 200; saved description verified |
| Export draft through both routes | Both 200; correct test flow |
| Delete through unprefixed route, without force or skipped reference checks | 200 |
| Read deleted flow through both routes | Both 404; cleanup confirmed |

These checks establish that both URL forms act on the same flow and support the tested draft lifecycle in Cloverhound. The unprefixed URLs are sufficient for those current-format operations; separate prefixed command variants provide no demonstrated additional capability. Publish and legacy import/export were not exercised in this write test. Keep legacy-format commands distinct from current-format commands regardless of URL-prefix consolidation.

## Redundant ID suffixes

The ID remains a parameter: `get --id`, `delete --id`, `patch --id`, and `update --id`. All six API areas were checked; current matches are in Contact Center. The old names remain compatibility aliases, including earlier versioned aliases. Qualifiers such as `get-ci-id` remain because they distinguish another lookup.

| Previous full command | Current full command | Purpose | Old name supported? |
| --- | --- | --- | --- |
| `webex cc auto-csat get-id` | `webex cc auto-csat get` | Retrieve an existing Auto CSAT resource by ID in a given organization | Yes |
| `webex cc auto-csat update-id` | `webex cc auto-csat update` | Update an existing Auto CSAT resource by ID in a given organization | Yes |
| `webex cc generated-summaries get-id` | `webex cc generated-summaries get` | Retrieve an existing Generated Summaries resource by ID in a given organization | Yes |
| `webex cc generated-summaries update-id` | `webex cc generated-summaries update` | Update an existing Generated Summaries resource by ID in a given organization | Yes |
| `webex cc business-hour get-id` | `webex cc business-hour get` | Retrieve an existing Business Hours resource by ID in a given organization. | Yes |
| `webex cc business-hour update-id` | `webex cc business-hour update` | Update an existing Business Hours resource by ID in a given organization. | Yes |
| `webex cc business-hour delete-id` | `webex cc business-hour delete` | Delete an existing Business Hours resource by ID in a given organization. | Yes |
| `webex cc holiday-list get-id` | `webex cc holiday-list get` | Retrieve an existing Holiday List by ID in a given organization. | Yes |
| `webex cc holiday-list update-id` | `webex cc holiday-list update` | Update an existing Holiday List by ID in a given organization. | Yes |
| `webex cc holiday-list delete-id` | `webex cc holiday-list delete` | Delete an existing Holiday List by ID in a given organization. | Yes |
| `webex cc overrides get-id` | `webex cc overrides get` | Retrieve an existing Overrides resource by ID in a given organization. | Yes |
| `webex cc overrides update-id` | `webex cc overrides update` | Update an existing Overrides resource by ID in a given organization. | Yes |
| `webex cc overrides delete-id` | `webex cc overrides delete` | Delete an existing Overrides resource by ID in a given organization. | Yes |
| `webex cc address-book get-id` | `webex cc address-book get` | Retrieve an existing Address Book by ID in a given organization. | Yes |
| `webex cc address-book update-id` | `webex cc address-book update` | Update an existing Address Book by ID in a given organization. | Yes |
| `webex cc address-book delete-id` | `webex cc address-book delete` | Delete an existing Address Book by ID in a given organization. | Yes |
| `webex cc audio-files get-id` | `webex cc audio-files get` | Retrieve an existing Audio File by ID in a given organization. | Yes |
| `webex cc audio-files update-id` | `webex cc audio-files update` | Update an existing Audio File by ID in a given organization. | Yes |
| `webex cc audio-files delete-id` | `webex cc audio-files delete` | Delete an existing Audio File by ID in a given organization. | Yes |
| `webex cc audio-files patch-id` | `webex cc audio-files patch` | Partially update Audio File by ID in a given organization. | Yes |
| `webex cc auxiliary-code get-id` | `webex cc auxiliary-code get` | Retrieve an existing Auxiliary Code by ID in a given organization. | Yes |
| `webex cc auxiliary-code update-id` | `webex cc auxiliary-code update` | Update an existing Auxiliary Code by ID in a given organization. | Yes |
| `webex cc auxiliary-code delete-id` | `webex cc auxiliary-code delete` | Delete an existing Auxiliary Code by ID in a given organization. | Yes |
| `webex cc contact-number get-id` | `webex cc contact-number get` | Retrieve an existing Contact Number by ID in a given organization. | Yes |
| `webex cc contact-number update-id` | `webex cc contact-number update` | Update an existing Contact Number by ID in a given organization. | Yes |
| `webex cc contact-number delete-id` | `webex cc contact-number delete` | Delete an existing Contact Number by ID in a given organization. | Yes |
| `webex cc contact-service-queue get-id` | `webex cc contact-service-queue get` | Retrieve an existing Contact Service Queue by ID in a given organization. | Yes |
| `webex cc contact-service-queue update-id` | `webex cc contact-service-queue update` | Update an existing Contact Service Queue by ID in a given organization. | Yes |
| `webex cc contact-service-queue delete-id` | `webex cc contact-service-queue delete` | Delete an existing Contact Service Queue by ID in a given organization. | Yes |
| `webex cc desktop-layout get-id` | `webex cc desktop-layout get` | Retrieve an existing Desktop Layout by ID in a given organization. | Yes |
| `webex cc desktop-layout update-id` | `webex cc desktop-layout update` | Update an existing Desktop Layout by ID in a given organization. | Yes |
| `webex cc desktop-layout delete-id` | `webex cc desktop-layout delete` | Delete an existing Desktop Layout by ID in a given organization. | Yes |
| `webex cc desktop-profile get-id` | `webex cc desktop-profile get` | Retrieve an existing Desktop Profile by ID in a given organization. | Yes |
| `webex cc desktop-profile update-id` | `webex cc desktop-profile update` | Update an existing Desktop Profile by ID in a given organization. | Yes |
| `webex cc desktop-profile delete-id` | `webex cc desktop-profile delete` | Delete an existing Desktop Profile by ID in a given organization. | Yes |
| `webex cc dial-plan get-id` | `webex cc dial-plan get` | Retrieve an existing Dial Plan by ID in a given organization. | Yes |
| `webex cc dial-plan update-id` | `webex cc dial-plan update` | Update an existing Dial Plan by ID in a given organization. | Yes |
| `webex cc dial-plan delete-id` | `webex cc dial-plan delete` | Delete an existing Dial Plan by ID in a given organization. | Yes |
| `webex cc entry-point get-id` | `webex cc entry-point get` | Retrieve an existing Entry Point by ID in a given organization. | Yes |
| `webex cc entry-point update-id` | `webex cc entry-point update` | Update an existing Entry Point by ID in a given organization. | Yes |
| `webex cc entry-point delete-id` | `webex cc entry-point delete` | Delete an existing Entry Point by ID in a given organization. | Yes |
| `webex cc global-variables get-id` | `webex cc global-variables get` | Retrieve an existing Global Variable by ID in a given organization. | Yes |
| `webex cc global-variables update-id` | `webex cc global-variables update` | Update an existing Global Variable by ID in a given organization | Yes |
| `webex cc global-variables delete-id` | `webex cc global-variables delete` | Delete an existing Global Variable by ID in a given organization. | Yes |
| `webex cc multimedia-profile get-id` | `webex cc multimedia-profile get` | Retrieve an existing Multimedia Profile by ID in a given organization. | Yes |
| `webex cc multimedia-profile update-id` | `webex cc multimedia-profile update` | Update an existing Multimedia Profile by ID in a given organization. | Yes |
| `webex cc multimedia-profile delete-id` | `webex cc multimedia-profile delete` | Delete an existing Multimedia Profile by ID in a given organization. | Yes |
| `webex cc outdial-ani get-id` | `webex cc outdial-ani get` | Retrieve an existing Outdial ANI by ID in a given organization. | Yes |
| `webex cc outdial-ani update-id` | `webex cc outdial-ani update` | Update an existing Outdial ANI by ID in a given organization. | Yes |
| `webex cc outdial-ani delete-id` | `webex cc outdial-ani delete` | Delete an existing Outdial ANI by ID in a given organization. | Yes |
| `webex cc site get-id` | `webex cc site get` | Retrieve an existing Site by ID in a given organization. | Yes |
| `webex cc site update-id` | `webex cc site update` | Update an existing Site by ID in a given organization. | Yes |
| `webex cc site delete-id` | `webex cc site delete` | Delete an existing Site by ID in a given organization. | Yes |
| `webex cc skill get-id` | `webex cc skill get` | Retrieve an existing Skill by ID in a given organization. | Yes |
| `webex cc skill update-id` | `webex cc skill update` | Update an existing Skill by ID in a given organization. | Yes |
| `webex cc skill delete-id` | `webex cc skill delete` | Delete an existing Skill by ID in a given organization. | Yes |
| `webex cc skill-profile get-id` | `webex cc skill-profile get` | Retrieve an existing Skill Profile by ID in a given organization. | Yes |
| `webex cc skill-profile update-id` | `webex cc skill-profile update` | Update an existing Skill Profile by ID in a given organization. | Yes |
| `webex cc skill-profile delete-id` | `webex cc skill-profile delete` | Delete an existing Skill Profile by ID in a given organization. | Yes |
| `webex cc team get-id` | `webex cc team get` | Retrieve an existing Team by ID in a given organization. | Yes |
| `webex cc team update-id` | `webex cc team update` | Update an existing Team by ID in a given organization. | Yes |
| `webex cc team delete-id` | `webex cc team delete` | Delete an existing Team by ID in a given organization. | Yes |
| `webex cc user-profiles get-id` | `webex cc user-profiles get` | Retrieve an existing user profile by ID in a given organization. | Yes |
| `webex cc user-profiles update-id` | `webex cc user-profiles update` | Update an existing user profile by ID in a given organization. | Yes |
| `webex cc user-profiles delete-id` | `webex cc user-profiles delete` | Delete an existing user profile by ID in a given organization. | Yes |
| `webex cc users get-id` | `webex cc users get` | Retrieve an existing Users by ID in a given organization. | Yes |
| `webex cc users update-id` | `webex cc users update` | Update an existing User by ID in a given organization. | Yes |
| `webex cc users patch-id` | `webex cc users patch` | Partially update User by ID in a given organization. | Yes |
| `webex cc work-types get-id` | `webex cc work-types get` | Retrieve an existing Work Type by ID in a given organization. | Yes |
| `webex cc work-types update-id` | `webex cc work-types update` | Update an existing Work Type by ID in a given organization. | Yes |
| `webex cc work-types delete-id` | `webex cc work-types delete` | Delete an existing Work Type by ID in a given organization. | Yes |
| `webex cc ai-feature get-id` | `webex cc ai-feature get` | Retrieve an existing AI Feature resource by ID in a given organization. | Yes |
| `webex cc ai-feature patch-id` | `webex cc ai-feature patch` | Partially update AI Feature resource by ID in a given organization. | Yes |
| `webex cc agent-personal-greeting-files get-id` | `webex cc agent-personal-greeting-files get` | Retrieve an existing Greeting File by ID in a given organization. | Yes |
| `webex cc agent-personal-greeting-files update-id` | `webex cc agent-personal-greeting-files update` | Update an existing Greeting File by ID in a given organization. | Yes |
| `webex cc agent-personal-greeting-files delete-id` | `webex cc agent-personal-greeting-files delete` | Delete an existing Greeting File by ID in a given organization. | Yes |
| `webex cc agent-personal-greeting-files patch-id` | `webex cc agent-personal-greeting-files patch` | Partially update Greeting File by ID in a given organization. | Yes |
| `webex cc resource-collection get-id` | `webex cc resource-collection get` | Retrieve an existing Resource Collection by ID in a given organization. | Yes |
| `webex cc resource-collection update-id` | `webex cc resource-collection update` | Update an existing resource collection by ID in a given organization. | Yes |
| `webex cc resource-collection delete-id` | `webex cc resource-collection delete` | Delete an existing resource collection by ID in a given organization. | Yes |
| `webex cc asset get-id` | `webex cc asset get` | Retrieve an existing Asset by ID in a given organization | Yes |
| `webex cc asset update-id` | `webex cc asset update` | Update an existing Asset by ID in a given organization | Yes |
| `webex cc asset delete-id` | `webex cc asset delete` | Delete an existing Asset by ID in a given organization | Yes |
| `webex cc asset patch-id` | `webex cc asset patch` | Partially update Asset by ID in a given organization | Yes |
| `webex cc channel get-id` | `webex cc channel get` | Retrieve an existing Channel by ID in a given organization | Yes |
| `webex cc channel update-id` | `webex cc channel update` | Update an existing Channel by ID in a given organization. | Yes |
| `webex cc channel delete-id` | `webex cc channel delete` | Delete an existing Channel by ID in a given organization. | Yes |
| `webex cc channel patch-id` | `webex cc channel patch` | Partially update a channel by ID | Yes |

## Maintaining this reference

Duplicate prefixed entries are excluded through `SUPERSEDED_ROUTES` in `codegen/naming_overrides.py`. Each exclusion identifies a required replacement request; extraction fails if it disappears. The generator keeps the replacement’s complete request contract. Legacy import/export retain their prefixed routes. Publish uses the unprefixed route but was not live-tested.

Exact `get-id`, `delete-id`, `patch-id`, and `update-id` names are shortened centrally by `apply_naming_overrides` in `codegen/extract_api_spec.py`, with collision checks and old-name aliases. Route-specific names and compatibility aliases are defined in `codegen/naming_overrides.py` and applied to the enriched spec before generating Go commands and skill references. Update the override table, run `make codegen`, and update this reference when changing a mapping. Handwritten request implementations continue to use `custom_*.go` and the generator skip lists.
