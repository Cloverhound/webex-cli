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


if __name__ == '__main__':
    unittest.main()
