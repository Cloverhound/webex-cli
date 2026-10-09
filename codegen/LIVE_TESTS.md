# Cloverhound CLI live validation

Tested 2026-10-09 with the CLI built from branch `auto/update-postman-collections`, starting at `8f354da`, plus the usage-report deletion fix found during testing. Authentication: `eumansky@cloverhound.com`; organization UUID: `67b6b831-2fb5-448a-b5f5-b06e5e6253c2`; Contact Center region: US.

These are selected live smoke tests, not coverage of every refreshed endpoint. Calls used explicit user/organization selection and disabled 429 retries. Outputs and errors were captured separately. Existing resources were read only; writes were limited to a disposable unpublished/unassigned flow and a disposable usage report.

| Commands tested | Result |
|---|---|
| `webex cc flow list`, `get`, `export --version=draft` | 200; ten existing flows; typed FlowV2 get/export payloads match |
| `webex cc flow export-legacy --version=draft` | 200; legacy process/diagram contract differs from typed export |
| `webex cc activities list-definitions`, `describe --activity-name=queue-contact` | 200; 54 definitions; activity schema retrieved |
| `webex cc activities get-input-choices --activity-name=queue-contact --input-name=channelType` | 200; static choices retrieved |
| `webex cc activities get-input-choices --activity-name=queue-contact --input-name=destination --parent-input-name=channelType --parent-value=TELEPHONY` | 200; dependent queue choices retrieved |
| `webex cc events list-specifications` | 200 |
| `webex cc templates list-flow`, `get-flow --id <returned-id>` | 200; 39 templates; individual template retrieved |
| `webex cc functions list`, `list-options`, `get --id <returned-id> --meta-data-only=true` | 200; existing function metadata and runtime choices retrieved; no execution/publishing |
| `webex cc asset list`, `webex cc channel list` | 200; empty lists (no item-level CRUD tested) |
| `webex cc site list`, `get --id <returned-id>`, `get-id --id <same-id>` | 200; four sites; canonical command and compatibility alias returned identical JSON |
| `webex calling locations list` | 200; baseline Calling read |
| `webex calling metrics get-call-quality-stats --last=24h` | 200; three calls, one call with media, one good-quality media call |
| `webex meetings recordings query --from=2026-10-01T00:00:00Z --to=2026-10-09T00:00:00Z --max=5` | 200; empty result; POST read query worked |
| `webex cc usage-reports list-resource-types`, `list` | 200; five record types; initial report list empty |
| `webex cc usage-reports create --resource-type=ContactSessionRecord --start-date=2026-10-08 --end-date=2026-10-08` | 200; report created asynchronously |
| `webex cc usage-reports get --report-id <test-report-id>` | 200; queued → in progress → completed in approximately 55 seconds |
| `webex cc usage-reports download-file --file-id <returned-file-id> --output=raw` | 200; 1,333-byte gzip; integrity verified; decompressed CSV 4,086 bytes, 218 columns, zero data rows |
| `webex cc usage-reports delete --report-id <test-report-id>` | Initially 406; after generator fix, 204; subsequent get 404 and list empty |
| `webex cc external-data-updates update-task-global-variables --dry-run --body <synthetic-body>` | PUT route/body construction checked; no live update or server-side body validation |
| `webex calling call-controls-for-me list` | 403 Forbidden; this login cannot use the endpoint; cause not established by this test |
| `webex admin archive-users query --filter='username eq "eumansky@cloverhound.com"'` | 404, gateway message “API Definition not found”; not a successful archived-user lookup; route/authentication requirements remain unverified |
| `webex messaging hds list-clusters-org --organization-id <Cloverhound-base64-org-id>` | 404, explicit service message “This organization does not contain any HDS clusters”; tenant prerequisite absent |

## Disposable flow lifecycle

Created `CLI_Smoke_Test_be85fe3abc76` (`6ac944e949c39e7ea3c4ce32`) with only start/disconnect nodes. Validation returned `valid: true` with an activity-description recommendation. Import returned 201. Get/list confirmed Draft status and no routing assignments. Lock, patch-draft, save-draft, export, validate-draft, export-legacy, and unlock returned 200. Readback assertions verified both description changes. Import and draft writes returned a `flow` wrapper with preflight warnings; reads/typed export returned the FlowV2 document directly.

Deletion used `--force=no --skip-rs-epcheck=false`, returned 200, and subsequent get returned 404. Final list contained the same ten original flow IDs. The temporary flow was never published or assigned. Existing flow get/export reads matched; existing flows were not written.

## Fix and limits

The usage-report delete endpoint rejected `Accept: application/json` with 406. Code generation now emits `Accept: */*` for that exact DELETE route, allowing its non-JSON acknowledgement. Live retry returned 204 and cleanup was verified. Other routes retain their existing Accept behavior. A generation regression check guards that scope.

No live task controls, completed-task variable updates, existing configuration changes, flow publishing, legacy flow import, or function import/execution were attempted. Populated usage CSV rows and asset/channel item operations remain untested. Raw response artifacts are local under `/tmp/webex-live-results`; they are not committed because they contain organization data.
