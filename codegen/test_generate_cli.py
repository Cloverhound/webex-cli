"""Regression checks for service-specific operations in the table."""
import json
import os
import re
import tempfile
import unittest

from generate_cli import generate_area_file, load_names, operation, resolve_base_url


class OperationNamesFileTest(unittest.TestCase):
    def test_a_missing_names_file_stops_codegen(self):
        with tempfile.TemporaryDirectory() as d:
            with self.assertRaises(SystemExit):
                load_names(os.path.join(d, 'operation_names.json'))

    def test_the_names_file_is_read(self):
        with tempfile.TemporaryDirectory() as d:
            path = os.path.join(d, 'operation_names.json')
            with open(path, 'w') as f:
                json.dump({'CallingGetAB': 'GET /a-b'}, f)
            self.assertEqual(load_names(path), {'CallingGetAB': 'GET /a-b'})


def area_file(area, collection, host, is_calling, group, endpoints):
    source, _, _ = generate_area_file(area, collection, host, is_calling, [{'group': group, 'endpoints': endpoints}],
                                      paging_overrides={})
    return source


class MessagingAttachmentsTest(unittest.TestCase):
    def test_create_uses_handwritten_override_but_other_messages_are_generated(self):
        endpoints = [
            {'command': 'create', 'method': 'POST', 'path': '/messages'},
            {'command': 'get', 'method': 'GET', 'path': '/messages/{messageId}'},
        ]
        code = area_file('messaging', 'Webex Messaging', 'webex', True, 'messages', endpoints)
        self.assertNotIn('Command: "create"', code)
        self.assertIn('Command: "get"', code)


class GroupDescriptionTest(unittest.TestCase):
    endpoints = [{'command': 'list', 'method': 'GET', 'path': '/things'}]

    def test_described_group_uses_its_description(self):
        code = area_file('cc', 'Webex Contact Center', 'cc', False, 'call-monitor', self.endpoints)
        self.assertIn('Aliases: []string{"call-monitoring"}', code)
        self.assertIn('Summary: "Live call monitoring sessions: start, barge in, coach, hold, and end"}', code)

    def test_other_groups_keep_the_generated_description(self):
        code = area_file('cc', 'Webex Contact Center', 'cc', False, 'site', self.endpoints)
        self.assertIn('Summary: "Site commands"}', code)


class UsageReportNegotiationTest(unittest.TestCase):
    def test_only_report_delete_accepts_non_json_acknowledgement(self):
        for method, path, expected in [
            ('DELETE', '/v1/usage-reports/{reportId}', '*/*'),
            ('GET', '/v1/usage-reports/{reportId}', ''),
            ('DELETE', '/v1/tasks/{taskId}', ''),
        ]:
            with self.subTest(method=method, path=path):
                endpoint = {'command': 'example', 'method': method, 'path': path}
                self.assertEqual(operation(endpoint, 'Webex Contact Center', 'cc', False)['accept'], expected)


class PaginationTest(unittest.TestCase):
    def paging(self, host, command, query_params, example=None, path='/things'):
        endpoint = {
            'command': command, 'method': 'GET', 'path': path,
            'query_params': [{'name': n} for n in query_params], 'example_list': example,
        }
        return operation(endpoint, 'Webex Test', host, host != 'cc', paging_overrides={})['paging']

    def test_contact_center_list(self):
        self.assertEqual(self.paging('cc', 'list', ['page', 'pageSize'], {'items': 'data'}),
                         '&webexapi.Paging{Mode: webexapi.Pages, Page: "page", Size: "pageSize", Items: "data", TotalPages: "meta.totalPages"}')

    def test_webex_list(self):
        self.assertEqual(self.paging('webex', 'list', ['start', 'max'], {'items': 'huntGroups'}),
                         '&webexapi.Paging{Mode: webexapi.Link, Start: "start", Size: "max", Items: "huntGroups", OmitEmpty: true}')

    def test_scim_list(self):
        self.assertEqual(self.paging('webex', 'search', ['startIndex', 'count'], {'items': 'Resources'}, path='/identity/scim/{orgId}/v2/Users'),
                         '&webexapi.Paging{Mode: webexapi.Offset, Start: "startIndex", First: 1, Size: "count", Items: "Resources", TotalItems: "totalResults"}')

    def test_list_without_paging_params_is_one_request(self):
        self.assertEqual(self.paging('cc', 'list-all', [], {'items': ''}), '&webexapi.Paging{Mode: webexapi.Single}')

    def test_query_the_paging_needs_is_sent(self):
        self.assertEqual(self.paging('cc', 'list', ['page', 'size', 'includePagination'], {'items': 'data'}),
                         '&webexapi.Paging{Mode: webexapi.Pages, Page: "page", Size: "size", Items: "data", TotalPages: "pageInfo.totalPages", Query: map[string]string{"includePagination": "true"}}')

    def test_an_override_is_used(self):
        endpoint = {'command': 'list-jobs', 'method': 'GET', 'path': '/jobs', 'query_params': [{'name': 'start'}, {'name': 'max'}], 'example_list': {'items': 'items'}}
        op = operation(endpoint, 'Webex Test', 'webex', True,
                       paging_overrides={'GET /jobs': {'mode': 'offset', 'start': 'start', 'size': 'max', 'items': 'items', 'omitEmpty': True, 'pageSize': 50}})
        self.assertEqual(op['paging'], '&webexapi.Paging{Mode: webexapi.Offset, Start: "start", Size: "max", PageSize: 50, Items: "items", OmitEmpty: true}')

    def test_get_that_is_not_a_list_is_not_paginated(self):
        self.assertIsNone(self.paging('webex', 'get', ['id']))

    def test_list_without_a_known_item_field_is_not_paginated(self):
        self.assertIsNone(self.paging('webex', 'list', ['max']))


class OperationIDTest(unittest.TestCase):
    park = '/telephony/config/locations/{locationId}/callParks/{callParkId}'
    settings = '/telephony/config/locations/{locationId}/callParks/settings'

    def ids(self, endpoints):
        """Every declared operation's Go name, with its method and path."""
        source = area_file('calling', 'Webex Cloud Calling', 'webex', True, 'call-park', endpoints)
        found = re.findall(r'^var (\w+) = &Operation\{.*?Method:\s+"(\w+)", Host: \w+, Path:\s+"([^"]+)"', source, re.S | re.M)
        return {name: (method, path) for name, method, path in found}

    def exported(self, endpoints):
        return {name: route for name, route in self.ids(endpoints).items() if name[0].isupper()}

    def test_an_operation_is_named_by_its_method_and_path_whatever_its_command(self):
        # A reordered collection makes extraction give "update" and "update-2" to
        # the other route; the names do not move.
        before = self.exported([{'command': 'update', 'method': 'PUT', 'path': self.park},
                                {'command': 'update-2', 'method': 'PUT', 'path': self.settings}])
        after = self.exported([{'command': 'update', 'method': 'PUT', 'path': self.settings},
                               {'command': 'update-2', 'method': 'PUT', 'path': self.park}])
        self.assertEqual(before, after)
        self.assertEqual(before, {
            'CallingPutTelephonyConfigLocationsByLocationIdCallParksByCallParkId': ('PUT', self.park),
            'CallingPutTelephonyConfigLocationsByLocationIdCallParksSettings': ('PUT', self.settings)})

    def test_a_route_listed_twice_with_one_request_is_named_once(self):
        listing = {'method': 'GET', 'path': '/workspaces/{workspaceId}', 'path_params': [{'name': 'workspaceId'}]}
        ids = self.ids([dict(listing, command='get'), dict(listing, command='get-2')])
        self.assertEqual(sorted(ids), ['CallingGetWorkspacesByWorkspaceId', 'callingGetWorkspacesByWorkspaceId_2'])

    def test_a_route_whose_listings_send_different_requests_is_named_once_with_every_parameter(self):
        plain = {'command': 'get', 'method': 'GET', 'path': '/workspaces/{workspaceId}',
                 'path_params': [{'name': 'workspaceId'}], 'query_params': [{'name': 'max'}]}
        detailed = {'command': 'get-2', 'method': 'GET', 'path': '/workspaces/{workspaceId}',
                    'path_params': [{'name': 'workspaceId'}], 'query_params': [{'name': 'includeCapabilities'}]}
        source = area_file('calling', 'Webex Cloud Calling', 'webex', True, 'workspaces', [plain, detailed])
        self.assertEqual(sorted(self.ids([plain, detailed])),
                         ['CallingGetWorkspacesByWorkspaceId', 'callingGetWorkspacesByWorkspaceId_1',
                          'callingGetWorkspacesByWorkspaceId_2'])
        combined = source.split('var CallingGetWorkspacesByWorkspaceId = ')[1].split('\n}\n')[0]
        self.assertEqual(re.findall(r'\{Name: "(\w+)", In: (\w+)', combined),
                         [('workspaceId', 'InPath'), ('max', 'InQuery'), ('includeCapabilities', 'InQuery')])
        self.assertIn('Command: ""', combined)
        operations = source.split('var callingOperations = []*Operation{')[1].split('}')[0]
        self.assertNotIn('CallingGetWorkspacesByWorkspaceId,', operations)

    def test_a_route_that_puts_one_parameter_in_two_places_is_not_named(self):
        query = {'command': 'get', 'method': 'GET', 'path': '/workspaces', 'query_params': [{'name': 'max'}]}
        header = {'command': 'get-2', 'method': 'GET', 'path': '/workspaces', 'extra_headers': [{'name': 'max'}]}
        self.assertEqual(self.exported([query, header]), {})
        self.assertEqual(len(self.ids([query, header])), 2)

    def test_a_name_never_moves_to_another_route_in_a_later_snapshot(self):
        names = {}
        first = generate_area_file('calling', 'Webex Cloud Calling', 'webex', True,
                                   [{'group': 'g', 'endpoints': [{'command': 'get', 'method': 'GET', 'path': '/a-b'}]}],
                                   paging_overrides={}, names=names)[0]
        self.assertIn('var CallingGetAB = ', first)
        self.assertEqual(names, {'CallingGetAB': 'GET /a-b'})
        later = generate_area_file('calling', 'Webex Cloud Calling', 'webex', True,
                                   [{'group': 'g', 'endpoints': [{'command': 'get', 'method': 'GET', 'path': '/a/b'}]}],
                                   paging_overrides={}, names=names)[0]
        self.assertNotIn('var CallingGetAB = ', later)
        both = generate_area_file('calling', 'Webex Cloud Calling', 'webex', True,
                                  [{'group': 'g', 'endpoints': [{'command': 'get', 'method': 'GET', 'path': '/a/b'},
                                                                {'command': 'get-2', 'method': 'GET', 'path': '/a-b'}]}],
                                  paging_overrides={}, names=names)[0]
        self.assertRegex(both, r'var CallingGetAB = &Operation\{\n\tID: "CallingGetAB", Area: "calling", Group: "g", Command: "get-2"')

    def test_two_routes_with_one_name_are_not_named(self):
        ids = self.ids([{'command': 'get', 'method': 'GET', 'path': '/a-b'},
                        {'command': 'get-2', 'method': 'GET', 'path': '/a/b'}])
        self.assertEqual(len(ids), 2)
        self.assertFalse(any(name[0].isupper() for name in ids))


class FlagTest(unittest.TestCase):
    def test_a_query_parameter_shares_the_flag_of_a_path_parameter_with_its_name(self):
        op = operation({'command': 'get', 'method': 'GET', 'path': '/things/{thingId}',
                        'path_params': [{'name': 'thingId', 'description': 'The thing.'}],
                        'query_params': [{'name': 'thingId'}]}, 'Webex Test', 'webex', True)
        self.assertEqual([(p['in'], p['flag'], p.get('required', False)) for p in op['params']],
                         [('InPath', 'thing-id', True), ('InQuery', 'thing-id', False)])

    def test_a_body_field_named_like_a_parameter_sends_its_value_only_when_a_string(self):
        op = operation({'command': 'update', 'method': 'PUT', 'path': '/things/{id}', 'has_body': True,
                        'path_params': [{'name': 'id'}],
                        'body_fields': [{'name': 'id', 'type': 'string'}, {'name': 'ID', 'type': 'int64'},
                                        {'name': 'count', 'type': 'int64'}]}, 'Webex Test', 'webex', True)
        self.assertEqual(op['body'], [{'name': 'id', 'flag': 'id', 'type': 'string'},
                                      {'name': 'count', 'flag': 'count', 'type': 'int64'}])

    def test_a_fixed_path_parameter_has_no_flag(self):
        op = operation({'command': 'get', 'method': 'GET', 'path': '/projects/{projectId}',
                        'path_params': [{'name': 'projectId'}]}, 'Webex Test', 'webex', True)
        self.assertEqual(op['params'], [{'name': 'projectId', 'in': 'InPath', 'value': '5e5c9ad6d61f870d6d778c1b'}])

    def test_a_contact_center_org_is_given_by_orgid(self):
        op = operation({'command': 'list', 'method': 'GET', 'path': '/organization/{orgId}/team',
                        'path_params': [{'name': 'orgId'}], 'query_params': [{'name': 'orgId'}]},
                       'Webex Contact Center', 'cc', False)
        self.assertEqual(op['path'], '/organization/{orgId}/team')
        self.assertEqual([(p['name'], p['flag'], p.get('orgUUID')) for p in op['params']],
                         [('orgId', 'orgid', True), ('orgId', 'orgid', True)])

    def test_contact_center_organization_id_and_header_take_the_uuid(self):
        op = operation({'command': 'create', 'method': 'POST', 'path': '/organization/{orgId}/team',
                        'path_params': [{'name': 'orgId'}], 'query_params': [{'name': 'organizationId'}],
                        'extra_headers': [{'name': 'X-ORGANIZATION-ID'}]}, 'Webex Contact Center', 'cc', False)
        self.assertEqual([(p['name'], p.get('orgUUID')) for p in op['params']],
                         [('orgId', True), ('organizationId', True), ('X-ORGANIZATION-ID', True)])

    def test_an_org_parameter_no_rule_covers_stops_codegen(self):
        for collection, host, is_calling, name in (('Webex Admin', 'webex', True, 'partnerOrgId'),
                                                   ('Webex Contact Center', 'cc', False, 'tenantOrganizationId'),
                                                   ('Webex Meetings', 'webex', True, 'org_id')):
            with self.subTest(name=name), self.assertRaisesRegex(ValueError, name):
                operation({'command': 'list', 'method': 'GET', 'path': '/things',
                           'query_params': [{'name': name}]}, collection, host, is_calling)

    def test_webexapis_org_parameters_take_the_webex_org_id(self):
        op = operation({'command': 'list', 'method': 'GET', 'path': '/telephony/config/people',
                        'query_params': [{'name': 'orgId'}, {'name': 'organizationId'},
                                         {'name': 'onBehalfOfSubPartnerOrgId'}]},
                       'Webex Cloud Calling', 'webex', True)
        self.assertEqual([(p['flag'], p['orgUUID']) for p in op['params']],
                         [('org-id', False), ('organization-id', False), ('on-behalf-of-sub-partner-org-id', False)])

    def test_slidos_session_org_takes_the_uuid(self):
        op = operation({'command': 'list-compliance-events', 'method': 'GET', 'path': '/slido/compliance/events',
                        'query_params': [{'name': 'sessionOrgId'}, {'name': 'max'}]}, 'Webex Meetings', 'webex', True)
        self.assertEqual([(p['name'], p['orgUUID']) for p in op['params']], [('sessionOrgId', True), ('max', False)])

    def test_from_without_to_gets_no_time_range(self):
        op = operation({'command': 'list', 'method': 'GET', 'path': '/v1/events', 'query_params': [{'name': 'from'}]},
                       'Webex Test', 'cc', False)
        self.assertIsNone(op['timeRange'])

    def test_from_gets_a_time_range(self):
        for path, expected in [('/v1/tasks', 'TimeEpochMs'), ('/v1/events', 'TimeISO')]:
            with self.subTest(path=path):
                op = operation({'command': 'list', 'method': 'GET', 'path': path,
                                'query_params': [{'name': 'from'}, {'name': 'to'}]}, 'Webex Test', 'cc', False)
                self.assertEqual(op['timeRange'], expected)


class IdentityTest(unittest.TestCase):
    def op(self, path, collection='Webex Admin', **extra):
        endpoint = {'command': 'example', 'method': 'GET', 'path': path, **extra}
        return operation(endpoint, collection, resolve_base_url(collection, path, 'webex'), True)

    def org_param(self, op):
        return next(p for p in op['params'] if p['name'] == 'orgId')

    def test_scim_routes_use_the_identity_host_and_the_org_uuid(self):
        op = self.op('/identity/scim/{orgId}/v2/Users', path_params=[{'name': 'orgId'}])
        self.assertEqual(op['host'], 'identity')
        self.assertTrue(self.org_param(op)['orgUUID'])
        self.assertEqual(self.org_param(op)['flag'], 'org-id')

    def test_scim_route_without_an_org_uses_the_identity_host(self):
        self.assertEqual(self.op('/identity/scim/v2/Users/me')['host'], 'identity')

    def test_archived_users_use_the_identity_host_and_the_org_uuid(self):
        op = self.op('/identity/organizations/{orgId}/v1/ArchivedUser', path_params=[{'name': 'orgId'}])
        self.assertEqual(op['host'], 'identity')
        self.assertTrue(self.org_param(op)['orgUUID'])

    def test_identity_org_routes_keep_the_v1_host_and_use_the_org_uuid(self):
        op = self.op('/identity/organizations/{orgId}', path_params=[{'name': 'orgId'}])
        self.assertEqual(op['host'], 'webex')
        self.assertTrue(self.org_param(op)['orgUUID'])

    def test_calling_identity_org_routes_use_the_org_uuid(self):
        op = self.op('/identity/organizations/{orgId}/jobs/sendActivationEmails', collection='Webex Cloud Calling',
                     path_params=[{'name': 'orgId'}])
        self.assertEqual(op['host'], 'webex')
        self.assertTrue(self.org_param(op)['orgUUID'])

    def test_other_admin_routes_are_unchanged(self):
        op = self.op('/organizations/{orgId}', path_params=[{'name': 'orgId'}])
        self.assertEqual(op['host'], 'webex')
        self.assertFalse(self.org_param(op)['orgUUID'])


class AiAgentStudioTest(unittest.TestCase):
    """AI Agent Studio routes come from codegen/supplemental and live on their own host."""

    def op(self, endpoint):
        host = resolve_base_url('Webex Contact Center', endpoint['path'], 'cc')
        return operation({'method': 'GET', **endpoint}, 'Webex Contact Center', host, False)

    def test_studio_routes_use_the_studio_host_and_redact_secrets(self):
        op = self.op({'command': 'list', 'path': '/api/v2/organizations/{orgId}/agents/', 'path_params': [{'name': 'orgId'}]})
        self.assertEqual((op['host'], op['path'], op['redact']), ('ai-agent', '/api/v2/organizations/{orgId}/agents/', True))

    def test_other_cc_routes_are_unchanged(self):
        op = self.op({'command': 'list', 'path': '/organization/{orgId}/v2/site', 'path_params': [{'name': 'orgId'}],
                      'description': 'Lists sites.'})
        self.assertEqual((op['host'], op['redact'], op['private'], op['description']), ('cc', False, False, 'Lists sites.'))

    def test_private_routes_are_marked(self):
        op = self.op({'command': 'export', 'path': '/api/v1/llmagent/export/', 'description': 'Exports the agent.'})
        self.assertEqual((op['private'], op['description']), (True, 'Exports the agent.'))

    def test_agent_token_header_is_resolved_from_agent_id(self):
        op = self.op({'command': 'export', 'path': '/api/v1/llmagent/export/',
                      'extra_headers': [{'name': 'bot-access-token', 'description': 'AI agent ID.'}]})
        self.assertEqual(op['params'], [{'name': 'bot-access-token', 'in': 'InAgentToken', 'flag': 'agent-id',
                                         'usage': 'AI agent ID.', 'required': True}])

    def test_agent_token_header_shares_the_agent_id_path_flag(self):
        op = self.op({'command': 'get', 'path': '/api/v1/bot/{agentId}/',
                      'path_params': [{'name': 'agentId', 'description': 'AI agent ID.'}],
                      'extra_headers': [{'name': 'bot-access-token', 'description': 'AI agent ID.'}]})
        self.assertEqual([(p['name'], p['in'], p['flag']) for p in op['params']],
                         [('agentId', 'InPath', 'agent-id'), ('bot-access-token', 'InAgentToken', 'agent-id')])
