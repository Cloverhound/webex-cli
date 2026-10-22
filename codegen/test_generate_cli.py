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
