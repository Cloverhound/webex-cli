"""Regression checks for service-specific request generation."""
import unittest

from generate_cli import generate_command


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
