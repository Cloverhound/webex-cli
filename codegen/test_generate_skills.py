"""Checks for the flag lists in generated skill command references."""
import unittest

from generate_skills import format_flags


class AgentTokenFlagTest(unittest.TestCase):
    def test_agent_token_header_is_listed_as_required_agent_id(self):
        ep = {'extra_headers': [{'name': 'bot-access-token'}]}
        self.assertEqual(format_flags(ep), '`--agent-id` *(required)*')

    def test_agent_id_from_path_and_header_is_listed_once(self):
        ep = {'path_params': [{'name': 'agentId'}], 'extra_headers': [{'name': 'bot-access-token'}]}
        self.assertEqual(format_flags(ep), '`--agent-id` *(required)*')


class BodyFlagTest(unittest.TestCase):
    def test_a_nested_body_points_to_the_example_in_help(self):
        ep = {'has_body': True, 'complex_body': True, 'body_example': '{"a": []}'}
        self.assertEqual(format_flags(ep), '`--body` *(example in `--help`)*, `--body-file`')

    def test_a_flat_body_lists_its_field_flags(self):
        ep = {'has_body': True, 'body_fields': [{'name': 'displayName', 'type': 'string'}]}
        self.assertEqual(format_flags(ep), '`--display-name`, `--body`, `--body-file`')


if __name__ == '__main__':
    unittest.main()
