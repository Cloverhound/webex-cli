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


class ContactCenterPaginationTest(unittest.TestCase):
    def paginated_block(self, query_params, is_calling=False):
        endpoint = {
            'command': 'list', 'method': 'GET', 'path': '/things',
            'query_params': [{'name': n} for n in query_params],
        }
        code = '\n'.join(generate_command(endpoint, 'exampleCmd', 'config.CcBaseURL', is_calling))
        return code[code.index('if config.Paginate()'):code.index('req.DoPaginated')]

    def test_size_endpoint_names_its_page_size_param(self):
        block = self.paginated_block(['page', 'size'])
        self.assertIn('req.PageSizeParam("size")', block)
        self.assertNotIn('includePagination', block)

    def test_include_pagination_is_forced_on(self):
        block = self.paginated_block(['page', 'size', 'includePagination'])
        self.assertIn('req.QueryParam("includePagination", "true")', block)

    def test_page_size_endpoint_keeps_default(self):
        block = self.paginated_block(['page', 'pageSize'])
        self.assertNotIn('PageSizeParam', block)

    def test_calling_endpoint_skips_cc_hints(self):
        block = self.paginated_block(['max', 'start', 'size', 'includePagination'], is_calling=True)
        self.assertNotIn('PageSizeParam', block)
        self.assertNotIn('includePagination', block)


class OffsetPaginationTest(unittest.TestCase):
    def generate(self, command, query_params, is_calling=True):
        endpoint = {
            'command': command, 'method': 'GET', 'path': '/things',
            'query_params': [{'name': n} for n in query_params],
        }
        return '\n'.join(generate_command(endpoint, 'exampleCmd', 'config.CallingBaseURL', is_calling))

    def paging_call(self, command, query_params):
        code = self.generate(command, query_params)
        block = code[code.index('if config.Paginate()'):code.index('req.DoPaginated')]
        calls = [l.strip() for l in block.splitlines() if 'OffsetPaging' in l]
        return calls[0] if calls else None

    def test_start_and_max_use_the_default(self):
        self.assertIsNone(self.paging_call('list', ['start', 'max']))

    def test_max_only_follows_links(self):
        self.assertEqual(self.paging_call('list', ['max']), 'req.OffsetPaging("", "max", 0)')

    def test_offset_and_max(self):
        self.assertEqual(self.paging_call('list-events', ['offset', 'max']), 'req.OffsetPaging("offset", "max", 0)')

    def test_scim_start_index_and_count(self):
        self.assertEqual(self.paging_call('search', ['startIndex', 'count']), 'req.OffsetPaging("startIndex", "count", 1)')

    def test_list_without_params_follows_links(self):
        self.assertEqual(self.paging_call('list', []), 'req.OffsetPaging("", "", 0)')

    def test_get_without_paging_params_is_not_paginated(self):
        for is_calling in (True, False):
            with self.subTest(is_calling=is_calling):
                self.assertNotIn('DoPaginated', self.generate('get', ['id'], is_calling))

    def test_cc_get_with_page_params_is_paginated(self):
        self.assertIn('DoPaginated', self.generate('get', ['pageSize'], is_calling=False))

    def test_cc_limit_and_offset_page_by_offset(self):
        code = self.generate('list', ['limit', 'offset'], is_calling=False)
        block = code[code.index('if config.Paginate()'):code.index('req.DoPaginated')]
        self.assertIn('req.OffsetPaging("offset", "limit", 0)', block)


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
                            'path_params': [{'name': 'orgId'}]}])
        self.assertIn('client.NewRequest(config.CcBaseURL,', code)
        self.assertNotIn('aiagent', code)

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
