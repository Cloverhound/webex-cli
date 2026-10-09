"""Stable CLI names for ambiguous or renamed Postman requests.

Like BASE_URL_OVERRIDES, these exceptions use HTTP routes, not collection order.
Current flow APIs use the unprefixed service routes verified in Cloverhound.
Legacy FDL import/export keep their original /flow-store/ routes.
Do not add aliases for unreleased numeric collision names.
"""

# collection -> (source group, method, path) -> (group, command, aliases)
NAMING_OVERRIDES = {
    'Webex Admin': {
        ('calling-metrics', 'GET', '/v1/analytics/callQualityStats'):
            ('metrics', 'get-call-quality-stats', ()),
    },
    'Webex Cloud Calling': {
        ('call-controls', 'POST', '/telephony/calls/members/me/dial'):
            ('call-controls-for-me', 'dial', ()),
        ('call-controls', 'POST', '/telephony/calls/members/me/answer'):
            ('call-controls-for-me', 'answer', ()),
        ('call-controls', 'POST', '/telephony/calls/members/me/hangup'):
            ('call-controls-for-me', 'hangup', ()),
        ('call-controls', 'GET', '/telephony/calls/members/me/calls'):
            ('call-controls-for-me', 'list', ()),
        ('call-controls', 'GET', '/telephony/calls/members/me/calls/{callId}'):
            ('call-controls-for-me', 'get', ()),
    },
    'Webex Contact Center': {
        ('flow', 'GET', '/flow-store/{orgId}/project/{projectId}/flows/{flowId}:export'): ('flow', 'export-legacy', ()),
        ('flow', 'POST', '/flow-store/{orgId}/project/{projectId}/flows:import'): ('flow', 'import-legacy', ()),
        ('contact-service-queue', 'GET', '/organization/{orgid}/contact-service-queue/by-skill-profile-id/{id}'):
            ('contact-service-queue', 'list-by-skill-profile', ('list-skill-csqs-skill-profile',)),
        ('contact-service-queue', 'POST', '/organization/{orgid}/contact-service-queue/delete-reference'):
            ('contact-service-queue', 'delete-references', ('delete-csq-references',)),
        ('contact-service-queue', 'POST', '/organization/{orgid}/contact-service-queue/fetch-manually-assignable-queues'):
            ('contact-service-queue', 'list-manually-assignable', ('list-manually-assignable-csqs',)),
        ('contact-service-queue', 'GET', '/organization/{orgid}/contact-service-queue/{id}/incoming-references'):
            ('contact-service-queue', 'list-references', ('list-csq-references-id',)),
        ('contact-service-queue', 'GET', '/organization/{orgid}/v2/contact-service-queue/by-user-id/{userid}/agent-based-queues'):
            ('contact-service-queue', 'list-agent-based', ()),
        ('contact-service-queue', 'GET', '/organization/{orgid}/v2/contact-service-queue/by-user-id/{userid}/skill-based-queues'):
            ('contact-service-queue', 'list-skill-based', ()),
        ('contact-service-queue', 'GET', '/organization/{orgid}/v2/contact-service-queue/by-user-id/{userid}/team-based-queues'):
            ('contact-service-queue', 'list-team-based', ()),
        ('contact-service-queue', 'POST', '/organization/{orgid}/contact-service-queue/fetch-by-dynamic-skills-and-skillProfile'):
            ('contact-service-queue', 'list-by-dynamic-skills', ('list-csqs-skills-profile',)),
        ('contact-service-queue', 'POST', '/organization/{orgid}/contact-service-queue/fetch-by-userId-skillProfileId'):
            ('contact-service-queue', 'list-by-user-skill-profile', ('list-csqs-user-profile',)),
        ('outdial-ani', 'GET', '/organization/{orgid}/outdial-ani/entry'):
            ('outdial-ani', 'list-entries', ('list-entry',)),
        ('outdial-ani', 'POST', '/organization/{orgid}/outdial-ani/{outDialAniId}/entry/bulk'):
            ('outdial-ani', 'bulk-save-entries', ('bulk-save-entry',)),
        ('outdial-ani', 'GET', '/organization/{orgid}/v2/outdial-ani/{outDialAniId}/entry'):
            ('outdial-ani', 'list-entries-for-ani', ('list-entry-2',)),
        ('users', 'POST', '/organization/{orgid}/user/fetch-user-details-by-ids'):
            ('users', 'list-by-ids', ('get-ids',)),
        ('users', 'PATCH', '/organization/{orgid}/user/bulk/update-dynamic-skill/{skillId}'):
            ('users', 'bulk-partial-update-dynamic-skills', ('bulk-update-dynamic-skills',)),
        ('tasks', 'POST', '/v1/tasks/{taskId}/unhold'):
            ('tasks', 'unhold', ('resume',)),
        ('tasks', 'POST', '/v2/tasks/{taskId}/messages'):
            ('tasks', 'append-message', ('update-2',)),
        ('tasks', 'POST', '/v1/tasks/{taskId}/pause'):
            ('tasks', 'pause-digital', ()),
        ('tasks', 'POST', '/v1/tasks/{taskId}/resume'):
            ('tasks', 'resume-digital', ()),
        ('agent-personal-greeting-files', 'POST', '/organization/{orgid}/v2/agent-personal-greeting'):
            ('agent-personal-greeting-files', 'create', ('create-v2-api',)),
        ('agent-personal-greeting-files', 'POST', '/organization/{orgid}/agent-personal-greeting/delete-reference'):
            ('agent-personal-greeting-files', 'delete-references', ('delete-references-1',)),
        ('agent-personal-greeting-files', 'GET', '/organization/{orgid}/v2/agent-personal-greeting/{id}'):
            ('agent-personal-greeting-files', 'get-id', ('get-id-v2-api',)),
        ('agent-personal-greeting-files', 'PUT', '/organization/{orgid}/v2/agent-personal-greeting/{id}'):
            ('agent-personal-greeting-files', 'update-id', ('update-id-v2-api',)),
        ('agent-personal-greeting-files', 'DELETE', '/organization/{orgid}/v2/agent-personal-greeting/{id}'):
            ('agent-personal-greeting-files', 'delete-id', ('delete-id-v2-api',)),
        ('agent-personal-greeting-files', 'PATCH', '/organization/{orgid}/v2/agent-personal-greeting/{id}'):
            ('agent-personal-greeting-files', 'patch-id', ('patch-id-v2-api',)),
        ('activities', 'GET', '/{orgId}/project/{projectId}/v2/activities'):
            ('activities', 'list-definitions', ()),
        ('activities', 'GET', '/{orgId}/project/{projectId}/v2/activities/{activityName}'):
            ('activities', 'describe', ()),
        ('activities', 'GET', '/{orgId}/project/{projectId}/v2/activities/{activityName}/inputs/{inputName}/choices'):
            ('activities', 'get-input-choices', ()),
        ('events', 'GET', '/{orgId}/project/{projectId}/v2/event-specifications'):
            ('events', 'list-specifications', ()),
        ('flows', 'GET', '/flow-store/{orgId}/project/{projectId}/flows'):
            ('flow', 'list', ()),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/flows/{flowId}:publish'):
            ('flow', 'publish', ()),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/flows/{flowId}:lock'):
            ('flow', 'lock', ()),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/flows/{flowId}:unlock'):
            ('flow', 'unlock', ()),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/v2/flows:validate'):
            ('flow', 'validate', ()),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/v2/flows:import'):
            ('flow', 'import', ()),
        ('flows', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/flows/{flowId}'):
            ('flow', 'get', ()),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/v2/flows/{flowId}'):
            ('flow', 'save-draft', ()),
        ('flows', 'PATCH', '/flow-store/{orgId}/project/{projectId}/v2/flows/{flowId}'):
            ('flow', 'patch-draft', ()),
        ('flows', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/flows/{flowId}:validate'):
            ('flow', 'validate-draft', ()),
        ('flows', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/flows/{flowId}:export'):
            ('flow', 'export', ()),
        ('flows', 'GET', '/{orgId}/project/{projectId}/flows'):
            ('flow', 'list', ()),
        ('flows', 'POST', '/{orgId}/project/{projectId}/flows/{flowId}:publish'):
            ('flow', 'publish', ()),
        ('flows', 'POST', '/{orgId}/project/{projectId}/flows/{flowId}:lock'):
            ('flow', 'lock', ()),
        ('flows', 'POST', '/{orgId}/project/{projectId}/flows/{flowId}:unlock'):
            ('flow', 'unlock', ()),
        ('flows', 'POST', '/{orgId}/project/{projectId}/v2/flows:validate'):
            ('flow', 'validate', ()),
        ('flows', 'POST', '/{orgId}/project/{projectId}/v2/flows:import'):
            ('flow', 'import', ()),
        ('flows', 'GET', '/{orgId}/project/{projectId}/v2/flows/{flowId}'):
            ('flow', 'get', ()),
        ('flows', 'POST', '/{orgId}/project/{projectId}/v2/flows/{flowId}'):
            ('flow', 'save-draft', ()),
        ('flows', 'PATCH', '/{orgId}/project/{projectId}/v2/flows/{flowId}'):
            ('flow', 'patch-draft', ()),
        ('flows', 'GET', '/{orgId}/project/{projectId}/v2/flows/{flowId}:validate'):
            ('flow', 'validate-draft', ()),
        ('flows', 'GET', '/{orgId}/project/{projectId}/v2/flows/{flowId}:export'):
            ('flow', 'export', ()),
        ('flows', 'GET', '/{orgId}/project/{projectId}/flows:search'):
            ('flow', 'search', ()),
        ('flows', 'DELETE', '/{orgId}/project/{projectId}/flows/{flowId}'):
            ('flow', 'delete', ()),
        ('functions', 'GET', '/v1/{orgId}/functions'):
            ('functions', 'list', ()),
        ('functions', 'POST', '/v1/{orgId}/functions'):
            ('functions', 'create', ()),
        ('functions', 'GET', '/v1/{orgId}/functions/{id}'):
            ('functions', 'get', ()),
        ('functions', 'PUT', '/v1/{orgId}/functions/{id}'):
            ('functions', 'update', ()),
        ('functions', 'DELETE', '/v1/{orgId}/functions/{id}'):
            ('functions', 'delete', ()),
        ('functions', 'POST', '/v1/{orgId}/functions:import'):
            ('functions', 'import', ()),
        ('functions', 'POST', '/v1/{orgId}/functions/{id}:unlock'):
            ('functions', 'unlock', ()),
        ('functions', 'POST', '/v1/{orgId}/functions/{id}:lock'):
            ('functions', 'lock', ()),
        ('functions', 'POST', '/v1/{orgId}/functions/{id}:test'):
            ('functions', 'test', ()),
        ('functions', 'POST', '/v1/{orgId}/functions/{id}:publish'):
            ('functions', 'publish', ()),
        ('functions', 'POST', '/v1/{orgId}/functions/{id}:export'):
            ('functions', 'export', ()),
        ('legacy-flows', 'GET', '/flow-store/{orgId}/project/{projectId}/flows/{flowId}:export'):
            ('flow', 'export-legacy', ()),
        ('legacy-flows', 'POST', '/flow-store/{orgId}/project/{projectId}/flows:import'):
            ('flow', 'import-legacy', ()),
        ('templates', 'GET', '/templates'):
            ('templates', 'list-flow', ()),
        ('templates', 'GET', '/templates/{id}'):
            ('templates', 'get-flow', ()),
        ('usage-reports', 'GET', '/v1/usage-reports/resource-types'):
            ('usage-reports', 'list-resource-types', ()),
        ('activities', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/activities'):
            ('activities', 'list-definitions', ()),
        ('activities', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/activities/{activityName}'):
            ('activities', 'describe', ()),
        ('activities', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/activities/{activityName}/inputs/{inputName}/choices'):
            ('activities', 'get-input-choices', ()),
        ('events', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/event-specifications'):
            ('events', 'list-specifications', ()),
        ('functions', 'GET', '/v1/{orgId}/functions/options'):
            ('functions', 'list-options', ()),
        ('templates', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/templates'):
            ('templates', 'list-flow', ()),
        ('templates', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/templates/{id}'):
            ('templates', 'get-flow', ()),
        ('usage-reports', 'POST', '/v1/usage-reports'):
            ('usage-reports', 'create', ()),
        ('usage-reports', 'GET', '/v1/usage-reports'):
            ('usage-reports', 'list', ()),
        ('usage-reports', 'GET', '/v1/usage-reports/{reportId}'):
            ('usage-reports', 'get', ()),
        ('usage-reports', 'DELETE', '/v1/usage-reports/{reportId}'):
            ('usage-reports', 'delete', ()),
        ('usage-reports', 'GET', '/v1/usage-reports/{fileId}/download'):
            ('usage-reports', 'download-file', ()),
    },
    'Webex Messaging': {
        ('hybrid-data-security', 'GET', '/hds/organizations/{organizationId}'):
            ('hds', 'get-org', ()),
        ('hybrid-data-security', 'GET', '/hds/organizations/{organizationId}/clusters'):
            ('hds', 'list-clusters-org', ()),
        ('hybrid-data-security', 'GET', '/hds/clusters/{clusterId}'):
            ('hds', 'get-cluster', ()),
        ('hybrid-data-security', 'GET', '/hds/clusters/{clusterId}/nodes'):
            ('hds', 'list-nodes-cluster', ()),
        ('hybrid-data-security', 'GET', '/hds/nodes/{nodeId}'):
            ('hds', 'get-node', ()),
        ('hybrid-data-security', 'GET', '/hds/organizations/{organizationId}/database'):
            ('hds', 'get-database-config-org', ()),
        ('hybrid-data-security', 'GET', '/hds/organizations/{organizationId}/tenants'):
            ('hds', 'list-tenants-org', ()),
        ('hybrid-data-security', 'GET', '/hds/nodes/{nodeId}/alarms'):
            ('hds', 'get-alarms-node', ()),
        ('hybrid-data-security', 'GET', '/hds/testResults/nodes/{nodeId}/networkTest'):
            ('hds', 'get-test-results-node', ()),
        ('hybrid-data-security', 'GET', '/hds/nodes/{nodeId}/resourceUsage'):
            ('hds', 'get-usage-node', ()),
        ('hybrid-data-security', 'GET', '/hds/clusters/{clusterId}/availability'):
            ('hds', 'get-availability-cluster', ()),
        ('hds', 'GET', '/hds/testResults/nodes/{nodeId}/networkTest'):
            ('hds', 'get-test-results-node', ()),
        ('hds', 'GET', '/hds/clusters/{clusterId}'):
            ('hds', 'get-cluster', ()),
        ('hds', 'GET', '/hds/organizations/{organizationId}'):
            ('hds', 'get-org', ()),
        ('hds', 'GET', '/hds/nodes/{nodeId}'):
            ('hds', 'get-node', ()),
        ('hds', 'GET', '/hds/organizations/{organizationId}/database/details'):
            ('hds', 'get-database-org', ()),
        ('hds', 'GET', '/hds/organizations/{organizationId}/multiTenant'):
            ('hds', 'get-multi-tenant-org', ()),
        ('hds', 'GET', '/hds/clusters/{clusterId}/availability'):
            ('hds', 'get-availability-cluster', ()),
        ('hds', 'GET', '/hds/organizations/{organizationId}/database'):
            ('hds', 'get-database-config-org', ('get-database-org-2',)),
        ('hds', 'GET', '/hds/organizations/{organizationId}/tenants'):
            ('hds', 'list-tenants-org', ('get-multi-tenant-org-2',)),
        ('hds', 'GET', '/hds/clusters/{clusterId}/nodes'):
            ('hds', 'list-nodes-cluster', ()),
        ('hds', 'GET', '/hds/organizations/{organizationId}/clusters'):
            ('hds', 'list-clusters-org', ()),
        ('hds', 'GET', '/hds/nodes/{nodeId}/alarms'):
            ('hds', 'get-alarms-node', ()),
        ('hds', 'GET', '/hds/nodes/{nodeId}/resourceUsage'):
            ('hds', 'get-usage-node', ()),
    },
}

GROUP_ALIASES = {'Webex Contact Center': {'flow': ['flows']}, 'Webex Messaging': {'hds': ['hybrid-data-security']}}

# Place user-facing commands by product while retaining upstream collections.
COLLECTION_MOVES = {("Webex Admin", "metrics"): "Webex Cloud Calling"}

# Drop these duplicate source requests in favor of the listed unprefixed request.
# Preserve the replacement request's body and query contract (not just its URL).
# Fail extraction if a replacement disappears in a future collection refresh.
SUPERSEDED_ROUTES = {
    "Webex Contact Center": {
        ('flow', 'GET', '/flow-store/{orgId}/project/{projectId}/flows'):
            ('GET', '/{orgId}/project/{projectId}/flows'),
        ('flow', 'POST', '/flow-store/{orgId}/project/{projectId}/flows/{flowId}:publish'):
            ('POST', '/{orgId}/project/{projectId}/flows/{flowId}:publish'),
        ('activities', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/activities'):
            ('GET', '/{orgId}/project/{projectId}/v2/activities'),
        ('activities', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/activities/{activityName}'):
            ('GET', '/{orgId}/project/{projectId}/v2/activities/{activityName}'),
        ('activities', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/activities/{activityName}/inputs/{inputName}/choices'):
            ('GET', '/{orgId}/project/{projectId}/v2/activities/{activityName}/inputs/{inputName}/choices'),
        ('events', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/event-specifications'):
            ('GET', '/{orgId}/project/{projectId}/v2/event-specifications'),
        ('flows', 'GET', '/flow-store/{orgId}/project/{projectId}/flows'):
            ('GET', '/{orgId}/project/{projectId}/flows'),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/flows/{flowId}:publish'):
            ('POST', '/{orgId}/project/{projectId}/flows/{flowId}:publish'),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/flows/{flowId}:lock'):
            ('POST', '/{orgId}/project/{projectId}/flows/{flowId}:lock'),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/flows/{flowId}:unlock'):
            ('POST', '/{orgId}/project/{projectId}/flows/{flowId}:unlock'),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/v2/flows:validate'):
            ('POST', '/{orgId}/project/{projectId}/v2/flows:validate'),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/v2/flows:import'):
            ('POST', '/{orgId}/project/{projectId}/v2/flows:import'),
        ('flows', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/flows/{flowId}'):
            ('GET', '/{orgId}/project/{projectId}/v2/flows/{flowId}'),
        ('flows', 'POST', '/flow-store/{orgId}/project/{projectId}/v2/flows/{flowId}'):
            ('POST', '/{orgId}/project/{projectId}/v2/flows/{flowId}'),
        ('flows', 'PATCH', '/flow-store/{orgId}/project/{projectId}/v2/flows/{flowId}'):
            ('PATCH', '/{orgId}/project/{projectId}/v2/flows/{flowId}'),
        ('flows', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/flows/{flowId}:validate'):
            ('GET', '/{orgId}/project/{projectId}/v2/flows/{flowId}:validate'),
        ('flows', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/flows/{flowId}:export'):
            ('GET', '/{orgId}/project/{projectId}/v2/flows/{flowId}:export'),
        ('templates', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/templates'):
            ('GET', '/templates'),
        ('templates', 'GET', '/flow-store/{orgId}/project/{projectId}/v2/templates/{id}'):
            ('GET', '/templates/{id}'),
    },
}
