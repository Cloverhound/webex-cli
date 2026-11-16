"""Regression checks for service-specific request generation."""
import unittest

from generate_cli import generate_command, generate_group_file


class MessagingAttachmentsTest(unittest.TestCase):
    def test_create_uses_handwritten_override_but_other_messages_are_generated(self):
        endpoints = [
            {'command': 'create', 'method': 'POST', 'path': '/messages'},
            {'command': 'get', 'method': 'GET', 'path': '/messages/{messageId}'},
        ]
        code = generate_group_file('messages', endpoints, 'messaging',
                                   'cmd.MessagingCmd', 'config.CallingBaseURL',
                                   False, 'Webex Messaging')
        self.assertNotIn('{ // create', code)
        self.assertIn('{ // get', code)


class GroupDescriptionTest(unittest.TestCase):
    def test_described_group_uses_its_description(self):
        code = generate_group_file('call-monitor', [], 'cc', 'cmd.CcCmd',
                                   'config.CcBaseURL', False, 'Webex Contact Center')
        self.assertIn('Short: "Live call monitoring sessions: start, barge in, coach, hold, and end",', code)
        self.assertIn('Aliases: []string{"call-monitoring"},', code)

    def test_other_groups_keep_the_generated_description(self):
        code = generate_group_file('site', [], 'cc', 'cmd.CcCmd',
                                   'config.CcBaseURL', False, 'Webex Contact Center')
        self.assertIn('Short: "Site commands",', code)


class UsageReportNegotiationTest(unittest.TestCase):
    def test_only_report_delete_accepts_non_json_acknowledgement(self):
        for method, path, expected in [
            ('DELETE', '/v1/usage-reports/{reportId}', True),
            ('GET', '/v1/usage-reports/{reportId}', False),
            ('DELETE', '/v1/tasks/{taskId}', False),
        ]:
            with self.subTest(method=method, path=path):
                endpoint = {'command': 'example', 'method': method, 'path': path}
                code = '\n'.join(generate_command(endpoint, 'exampleCmd', 'config.CcBaseURL', False))
                self.assertEqual('req.Header("Accept", "*/*")' in code, expected)


class PaginationTest(unittest.TestCase):
    def paging_call(self, base, command, query_params, example=None, path='/things'):
        endpoint = {
            'command': command, 'method': 'GET', 'path': path,
            'query_params': [{'name': n} for n in query_params], 'example_list': example,
        }
        code = '\n'.join(generate_command(endpoint, 'exampleCmd', base, base != 'config.CcBaseURL', paging_overrides={}))
        calls = [l.strip() for l in code.splitlines() if 'DoPaginated' in l]
        return calls[0] if calls else None

    def test_contact_center_list(self):
        self.assertEqual(self.paging_call('config.CcBaseURL', 'list', ['page', 'pageSize'], {'items': 'data'}),
                         'resp, statusCode, err := req.DoPaginated(client.Paging{Mode: client.Pages, Page: "page", Size: "pageSize", Items: "data", TotalPages: "meta.totalPages"})')

    def test_webex_list(self):
        self.assertEqual(self.paging_call('config.CallingBaseURL', 'list', ['start', 'max'], {'items': 'huntGroups'}),
                         'resp, statusCode, err := req.DoPaginated(client.Paging{Mode: client.Link, Start: "start", Size: "max", Items: "huntGroups", OmitEmpty: true})')

    def test_scim_list(self):
        self.assertEqual(self.paging_call('config.CallingBaseURL', 'search', ['startIndex', 'count'], {'items': 'Resources'}, path='/identity/scim/{orgId}/v2/Users'),
                         'resp, statusCode, err := req.DoPaginated(client.Paging{Mode: client.Offset, Start: "startIndex", First: 1, Size: "count", Items: "Resources", TotalItems: "totalResults"})')

    def test_list_without_paging_params_is_one_request(self):
        self.assertEqual(self.paging_call('config.CcBaseURL', 'list-all', [], {'items': ''}),
                         'resp, statusCode, err := req.DoPaginated(client.Paging{Mode: client.Single})')

    def test_query_the_paging_needs_is_sent(self):
        self.assertEqual(self.paging_call('config.CcBaseURL', 'list', ['page', 'size', 'includePagination'], {'items': 'data'}),
                         'resp, statusCode, err := req.DoPaginated(client.Paging{Mode: client.Pages, Page: "page", Size: "size", Items: "data", TotalPages: "pageInfo.totalPages", Query: map[string]string{"includePagination": "true"}})')

    def test_an_override_is_used(self):
        endpoint = {'command': 'list-jobs', 'method': 'GET', 'path': '/jobs', 'query_params': [{'name': 'start'}, {'name': 'max'}], 'example_list': {'items': 'items'}}
        code = '\n'.join(generate_command(endpoint, 'exampleCmd', 'config.CallingBaseURL', True,
                                          paging_overrides={'GET /jobs': {'mode': 'offset', 'start': 'start', 'size': 'max', 'items': 'items', 'omitEmpty': True, 'pageSize': 50}}))
        self.assertIn('req.DoPaginated(client.Paging{Mode: client.Offset, Start: "start", Size: "max", PageSize: 50, Items: "items", OmitEmpty: true})', code)

    def test_get_that_is_not_a_list_is_not_paginated(self):
        self.assertIsNone(self.paging_call('config.CallingBaseURL', 'get', ['id']))

    def test_list_without_a_known_item_field_is_not_paginated(self):
        self.assertIsNone(self.paging_call('config.CallingBaseURL', 'list', ['max']))


class AiAgentStudioTest(unittest.TestCase):
    """AI Agent Studio routes come from codegen/supplemental and live on their own host."""

    def group(self, endpoints):
        return generate_group_file('ai-agent', endpoints, 'cc', 'cmd.CcCmd',
                                   'config.CcBaseURL', False, 'Webex Contact Center')

    def test_studio_routes_use_the_studio_host_and_redact_secrets(self):
        code = self.group([{'command': 'list', 'method': 'GET', 'path': '/api/v2/organizations/{orgId}/agents/',
                            'path_params': [{'name': 'orgId'}]}])
        self.assertIn('client.NewRequest(config.AiAgentBaseURL, "GET", "/api/v2/organizations/{orgid}/agents/")', code)
        self.assertIn('"github.com/Cloverhound/webex-cli/internal/aiagent"', code)
        self.assertEqual(code.count('resp = aiagent.RedactSecrets(resp)'), 2)

    def test_other_cc_routes_are_unchanged(self):
        code = self.group([{'command': 'list', 'method': 'GET', 'path': '/organization/{orgId}/v2/site',
                            'path_params': [{'name': 'orgId'}],
                            'description': 'Lists sites.'}])
        self.assertIn('client.NewRequest(config.CcBaseURL,', code)
        self.assertNotIn('aiagent', code)
        self.assertNotIn('supported by Cisco', code)

    def test_private_routes_warn_in_help_and_when_run(self):
        code = self.group([{'command': 'export', 'method': 'GET', 'path': '/api/v1/llmagent/export/',
                            'description': 'Exports the agent.'}])
        long_help = code[code.index('Long:'):code.index('RunE:')]
        self.assertIn('Exports the agent.', long_help)
        self.assertIn('not documented or supported by Cisco', long_help)
        run = code[code.index('RunE:'):]
        self.assertEqual(run.splitlines()[1].strip(), 'aiagent.PrivateAPINotice()')

    def test_private_route_without_description_still_warns(self):
        code = self.group([{'command': 'export', 'method': 'GET', 'path': '/api/v1/llmagent/export/'}])
        self.assertIn('not documented or supported by Cisco', code)

    def test_agent_token_header_is_resolved_from_agent_id(self):
        code = self.group([{'command': 'export', 'method': 'GET', 'path': '/api/v1/llmagent/export/',
                            'extra_headers': [{'name': 'bot-access-token', 'description': 'AI agent ID.'}]}])
        self.assertIn('cmd.Flags().StringVar(&agentId, "agent-id", "", "AI agent ID.")', code)
        self.assertIn('cmd.MarkFlagRequired("agent-id")', code)
        self.assertIn('botToken, err := aiagent.BotToken(agentId)', code)
        self.assertIn('req.Header("bot-access-token", botToken)', code)
        self.assertNotIn('"bot-access-token", "",', code)

    def test_agent_token_header_shares_the_agent_id_path_flag(self):
        code = self.group([{'command': 'get', 'method': 'GET', 'path': '/api/v1/bot/{agentId}/',
                            'path_params': [{'name': 'agentId', 'description': 'AI agent ID.'}],
                            'extra_headers': [{'name': 'bot-access-token', 'description': 'AI agent ID.'}]}])
        self.assertEqual(code.count('"agent-id", "",'), 1)
        self.assertIn('req.PathParam("agentId", agentId)', code)
        self.assertIn('aiagent.BotToken(agentId)', code)
