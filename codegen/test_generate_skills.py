"""Checks for the flag lists in generated skill command references."""
import unittest

from generate_cli import UPLOADS
from generate_skills import format_flags


class AgentTokenFlagTest(unittest.TestCase):
    def test_agent_token_header_is_listed_as_required_agent_id(self):
        ep = {'extra_headers': [{'name': 'bot-access-token'}]}
        self.assertEqual(format_flags(ep), '`--agent-id` *(required)*')

    def test_agent_id_from_path_and_header_is_listed_once(self):
        ep = {'path_params': [{'name': 'agentId'}], 'extra_headers': [{'name': 'bot-access-token'}]}
        self.assertEqual(format_flags(ep), '`--agent-id` *(required)*')


class UploadFlagTest(unittest.TestCase):
    def test_an_upload_lists_its_part_flags_in_place_of_the_collection_body(self):
        ep = {'path_params': [{'name': 'orgid'}, {'name': 'id'}], 'has_body': True,
              'body_fields': [{'name': 'logoFileInfo', 'type': 'string'}]}
        upload = UPLOADS[('Webex Contact Center', 'PUT', '/organization/{orgid}/channel/{id}')]
        self.assertEqual(format_flags(ep, upload),
                         '`--id` *(required)*, `--body` or `--body-file` *(required)*, `--logo-file`')


if __name__ == '__main__':
    unittest.main()
