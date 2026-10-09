import unittest
from diff_spec import diff_specs, endpoint_changes


class SpecDiffTest(unittest.TestCase):
    def test_route_and_body_changes_are_reported_without_renaming(self):
        old = {'method': 'POST', 'path': '/legacy', 'has_body': False}
        new = {'method': 'POST', 'path': '/v2', 'has_body': True}
        report = diff_specs({'Webex Contact Center': {'flow': {'import': old}}},
                            {'Webex Contact Center': {'flow': {'import': new}}})
        self.assertIn('**Modified** `import`', report)
        self.assertIn('`/legacy` → `/v2`', report)
        self.assertIn('request body support changed', report)

    def test_alias_is_a_rename_not_a_deletion(self):
        old = {'Webex Contact Center': {'site': {'get-id': {'method': 'GET'}}}}
        new = {'Webex Contact Center': {'site': {'get': {'method': 'GET', 'aliases': ['get-id']}}}}
        report = diff_specs(old, new)
        self.assertIn('**Renamed** `get-id` → `get`', report)
        self.assertNotIn('**Deleted**', report)

    def test_deleted_group_lists_every_removed_command(self):
        report = diff_specs({'Webex Contact Center': {'old': {'one': {'method': 'GET'}, 'two': {'method': 'POST'}}}}, {})
        self.assertIn('**Deleted group**', report)
        self.assertIn('**Deleted** `one`', report)
        self.assertIn('**Deleted** `two`', report)

    def test_parameter_add_remove_and_type_change(self):
        old = {'body_fields': [{'name': 'old'}, {'name': 'shared', 'type': 'string'}]}
        new = {'body_fields': [{'name': 'new'}, {'name': 'shared', 'type': 'bool'}]}
        report = '; '.join(endpoint_changes(old, new))
        for expected in ['added `new`', 'removed `old`', 'updated `shared`']:
            self.assertIn(expected, report)
