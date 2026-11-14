"""Regression checks for the operation table the CLI is built from."""
import unittest

from generate_cli import build_operation, generate_area_file, render_operation


def operation(ep, collection='Webex Contact Center', is_calling=False, host='cc', group='example'):
    ep = {'command': 'example', 'method': 'GET', 'path': '/things', **ep}
    return build_operation(ep, collection, is_calling, host, group)


class MessagingAttachmentsTest(unittest.TestCase):
    def test_create_uses_handwritten_override_but_other_messages_are_generated(self):
        create = operation({'command': 'create', 'method': 'POST', 'path': '/messages'},
                           'Webex Messaging', True, 'webex', 'messages')
        get = operation({'command': 'get', 'path': '/messages/{messageId}'},
                        'Webex Messaging', True, 'webex', 'messages')
        self.assertTrue(create['handwritten'])
        self.assertFalse(get['handwritten'])


class GroupDescriptionTest(unittest.TestCase):
    def area(self, group):
        source, _ = generate_area_file('Webex Contact Center',
                                       [{'group': group, 'endpoints': [{'command': 'list', 'method': 'GET', 'path': '/x'}]}],
                                       set())
        return source

    def test_described_group_uses_its_description(self):
        source = self.area('call-monitor')
        self.assertIn('Summary: "Live call monitoring sessions: start, barge in, coach, hold, and end"', source)
        self.assertIn('Aliases: []string{"call-monitoring"}', source)

    def test_other_groups_keep_the_generated_description(self):
        self.assertIn('Summary: "Site commands"', self.area('site'))

    def test_search_is_published_by_handwritten_commands(self):
        self.assertIn('Name: "search", Summary: "Search commands", Handwritten: true', self.area('search'))


class UsageReportNegotiationTest(unittest.TestCase):
    def test_only_report_delete_accepts_non_json_acknowledgement(self):
        for method, path, expected in [
            ('DELETE', '/v1/usage-reports/{reportId}', '*/*'),
            ('GET', '/v1/usage-reports/{reportId}', ''),
            ('DELETE', '/v1/tasks/{taskId}', ''),
        ]:
            with self.subTest(method=method, path=path):
                self.assertEqual(operation({'method': method, 'path': path})['accept'], expected)


class ContactCenterPaginationTest(unittest.TestCase):
    def paging(self, query_params, is_calling=False):
        return operation({'command': 'list', 'query_params': [{'name': n} for n in query_params]},
                         is_calling=is_calling)['paging']

    def test_size_endpoint_names_its_page_size_param(self):
        paging = self.paging(['page', 'size'])
        self.assertEqual(paging['size_param'], 'size')
        self.assertEqual(paging['query'], {})

    def test_include_pagination_is_forced_on(self):
        self.assertEqual(self.paging(['page', 'size', 'includePagination'])['query'], {'includePagination': 'true'})

    def test_page_size_endpoint_keeps_default(self):
        self.assertEqual(self.paging(['page', 'pageSize'])['size_param'], '')

    def test_calling_endpoint_skips_cc_hints(self):
        paging = self.paging(['max', 'start', 'size', 'includePagination'], is_calling=True)
        self.assertEqual((paging['size_param'], paging['query'], paging['offset']), ('', {}, None))


class OffsetPaginationTest(unittest.TestCase):
    def paging(self, command, query_params, is_calling=True):
        return operation({'command': command, 'query_params': [{'name': n} for n in query_params]},
                         'Webex Cloud Calling', is_calling, 'webex')['paging']

    def test_start_and_max_use_the_default(self):
        self.assertIsNone(self.paging('list', ['start', 'max'])['offset'])

    def test_max_only_follows_links(self):
        self.assertEqual(self.paging('list', ['max'])['offset'], ('', 'max', 0))

    def test_offset_and_max(self):
        self.assertEqual(self.paging('list-events', ['offset', 'max'])['offset'], ('offset', 'max', 0))

    def test_scim_start_index_and_count(self):
        self.assertEqual(self.paging('search', ['startIndex', 'count'])['offset'], ('startIndex', 'count', 1))

    def test_list_without_params_follows_links(self):
        self.assertEqual(self.paging('list', [])['offset'], ('', '', 0))

    def test_get_without_paging_params_is_not_paginated(self):
        for is_calling in (True, False):
            with self.subTest(is_calling=is_calling):
                self.assertIsNone(self.paging('get', ['id'], is_calling))

    def test_cc_get_with_page_params_is_paginated(self):
        self.assertIsNotNone(self.paging('get', ['pageSize'], is_calling=False))

    def test_cc_limit_and_offset_page_by_offset(self):
        self.assertEqual(self.paging('list', ['limit', 'offset'], is_calling=False)['offset'], ('offset', 'limit', 0))

    def test_writes_do_not_page(self):
        op = operation({'command': 'list', 'method': 'POST', 'query_params': [{'name': 'max'}]},
                       'Webex Cloud Calling', True, 'webex')
        self.assertIsNone(op['paging'])


class FlagSharingTest(unittest.TestCase):
    def test_contact_center_org_id_becomes_orgid_but_query_keeps_its_api_name(self):
        op = operation({'path': '/organization/{orgId}/team', 'path_params': [{'name': 'orgId'}],
                        'query_params': [{'name': 'orgId'}]})
        self.assertEqual(op['path'], '/organization/{orgid}/team')
        self.assertEqual([(p['name'], p['in'], p['flag']) for p in op['params']],
                         [('orgid', 'path', 'orgid'), ('orgId', 'query', 'orgid')])

    def test_project_id_is_fixed(self):
        op = operation({'path': '/{orgId}/project/{projectId}/flows',
                        'path_params': [{'name': 'orgId'}, {'name': 'projectId'}]})
        self.assertEqual(op['params'][1], {'name': 'projectId', 'in': 'path', 'value': '5e5c9ad6d61f870d6d778c1b'})

    def test_string_body_field_shares_a_parameter_flag_and_other_types_are_dropped(self):
        op = operation({'method': 'PUT', 'has_body': True,
                        'path_params': [{'name': 'id'}], 'query_params': [{'name': 'enabled'}],
                        'body_fields': [{'name': 'id', 'type': 'string'}, {'name': 'enabled', 'type': 'bool'},
                                        {'name': 'maxRings', 'type': 'int64'}]})
        self.assertEqual(op['body']['fields'], [{'name': 'id', 'flag': 'id', 'type': 'string'},
                                                {'name': 'maxRings', 'flag': 'max-rings', 'type': 'int64'}])

    def test_go_keyword_names_share_one_flag(self):
        op = operation({'query_params': [{'name': 'type'}, {'name': 'typeVal'}]})
        self.assertEqual([p['flag'] for p in op['params']], ['type', 'type'])

    def test_from_selects_the_time_format(self):
        self.assertEqual(operation({'path': '/v1/tasks', 'query_params': [{'name': 'from'}]})['time_range'], 'epoch-ms')
        self.assertEqual(operation({'path': '/v1/other', 'query_params': [{'name': 'from'}]})['time_range'], 'iso')


class AiAgentStudioTest(unittest.TestCase):
    """AI Agent Studio routes come from codegen/supplemental and live on their own host."""

    def test_studio_routes_use_the_studio_host_and_redact_secrets(self):
        op = operation({'path': '/api/v2/organizations/{orgId}/agents/', 'path_params': [{'name': 'orgId'}]})
        self.assertEqual((op['host'], op['path'], op['redact'], op['private']),
                         ('ai-agent', '/api/v2/organizations/{orgid}/agents/', True, True))

    def test_other_cc_routes_are_unchanged(self):
        op = operation({'path': '/organization/{orgId}/v2/site', 'path_params': [{'name': 'orgId'}]})
        self.assertEqual((op['host'], op['redact'], op['private']), ('cc', False, False))

    def test_agent_token_header_is_resolved_from_agent_id(self):
        op = operation({'path': '/api/v1/llmagent/export/',
                        'extra_headers': [{'name': 'bot-access-token', 'description': 'AI agent ID.'}]})
        self.assertEqual(op['params'], [{'name': 'bot-access-token', 'in': 'agent-token', 'flag': 'agent-id',
                                         'usage': 'AI agent ID.'}])

    def test_agent_token_header_shares_the_agent_id_path_flag(self):
        op = operation({'path': '/api/v1/bot/{agentId}/', 'path_params': [{'name': 'agentId'}],
                        'extra_headers': [{'name': 'bot-access-token'}]})
        self.assertEqual([p['flag'] for p in op['params']], ['agent-id', 'agent-id'])


class RenderTest(unittest.TestCase):
    def test_text_with_backticks_and_non_ascii_is_a_valid_go_string(self):
        op = operation({'description': 'Use `x` — ✓ "quoted"\nnext'})
        lines = '\n'.join(render_operation('CCExample', 'cc', op))
        self.assertIn('Description: "Use `x` — ✓ \\"quoted\\"\\nnext",', lines)


if __name__ == '__main__':
    unittest.main()
