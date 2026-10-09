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
