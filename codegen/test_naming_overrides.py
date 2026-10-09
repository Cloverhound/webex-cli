"""Regression checks for route-based names and compatibility across refreshes."""
import copy
import unittest
import json
from pathlib import Path

from extract_api_spec import extract_collection, merge_folders, apply_naming_overrides
from naming_overrides import NAMING_OVERRIDES


class NamingOverridesTest(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.raw = {}
        for path in (Path(__file__).parent / 'postman').glob('*.json'):
            name, folders = extract_collection(path)
            cls.raw[name] = merge_folders(folders)

    def resolved(self, area):
        return apply_naming_overrides(area, self.raw[area])

    def command(self, area, group, name):
        folder = next(g for g in self.resolved(area) if g['group'] == group)
        return next(e for e in folder['endpoints']
                    if name in [e['command'], *e.get('aliases', [])])

    def test_every_override_matches_a_collection_endpoint(self):
        for area, overrides in NAMING_OVERRIDES.items():
            identities = {(g['group'], e['method'], e['path'])
                          for g in self.raw[area] for e in g['endpoints']}
            self.assertFalse(set(overrides) - identities, area)

    def test_names_do_not_depend_on_numeric_collision_suffixes(self):
        for area, overrides in NAMING_OVERRIDES.items():
            folders = copy.deepcopy(self.raw[area])
            # Simulate Postman reordering/name normalization changes. Explicit
            # routes must keep their public names regardless of generated names.
            for g in folders:
                for e in g['endpoints']:
                    if (g['group'], e['method'], e['path']) in overrides:
                        e['command'] = 'upstream-renamed-99'
            self.assertEqual(self.resolved(area), apply_naming_overrides(area, folders))

    def test_voice_resume_keeps_its_original_meaning(self):
        area = 'Webex Contact Center'
        self.assertEqual(self.command(area, 'tasks', 'resume')['path'], '/v1/tasks/{taskId}/unhold')
        self.assertEqual(self.command(area, 'tasks', 'resume-digital')['path'], '/v1/tasks/{taskId}/resume')
        self.assertEqual(self.command(area, 'tasks', 'resume-recording')['path'], '/v1/tasks/{taskId}/record/resume')

    def test_redundant_id_names_are_aliases_in_every_area(self):
        for area in self.raw:
            for group in self.resolved(area):
                for ep in group['endpoints']:
                    self.assertNotIn(ep['command'], ('get-id', 'delete-id', 'patch-id', 'update-id'))
            for group in self.raw[area]:
                for ep in group['endpoints']:
                    if ep['command'] in ('get-id', 'delete-id', 'patch-id', 'update-id'):
                        result = self.command(area, group['group'], ep['command'])
                        self.assertEqual(result['path'], ep['path'])
                        self.assertEqual(result['command'], ep['command'][:-3])

    def test_shortening_rejects_a_collision_instead_of_numbering(self):
        folders = [{'group': 'example', 'original_folders': ['Example'], 'endpoints': [
            {'command': 'get-id', 'method': 'GET', 'path': '/example/{id}'},
            {'command': 'get', 'method': 'GET', 'path': '/different'},
        ]}]
        with self.assertRaisesRegex(ValueError, 'Duplicate command or alias'):
            apply_naming_overrides('Webex Admin', folders)

    def test_flow_versions_and_routes_remain_distinct(self):
        area = 'Webex Contact Center'
        self.assertEqual(self.command(area, 'flow', 'import-legacy')['path'],
                         '/flow-store/{orgId}/project/{projectId}/flows:import')
        self.assertEqual(self.command(area, 'flow', 'import')['path'],
                         '/{orgId}/project/{projectId}/v2/flows:import')

    def test_current_flow_routes_are_consolidated(self):
        for group in self.resolved('Webex Contact Center'):
            if group['group'] not in ('flow', 'activities', 'events', 'templates'):
                continue
            for ep in group['endpoints']:
                self.assertFalse(ep['command'].endswith('-direct'))
                if not ep['command'].endswith('-legacy'):
                    self.assertFalse(ep['path'].startswith('/flow-store/'))
        choices = self.command('Webex Contact Center', 'activities', 'get-input-choices')
        self.assertIn('search', [p['name'] for p in choices['query_params']])
        self.assertNotIn('mode', [p['name'] for p in choices['query_params']])
        self.assertEqual(self.command('Webex Contact Center', 'templates', 'list-flow')['path'], '/templates')

    def test_missing_replacement_route_fails_extraction(self):
        area = 'Webex Contact Center'
        folders = copy.deepcopy(self.raw[area])
        for folder in folders:
            folder['endpoints'] = [e for e in folder['endpoints'] if e['path'] != '/templates']
        with self.assertRaisesRegex(ValueError, 'Missing replacement'):
            apply_naming_overrides(area, folders)

    def test_hds_duplicates_are_merged_but_distinct_routes_survive(self):
        groups = self.resolved('Webex Messaging')
        self.assertNotIn('hybrid-data-security', [g['group'] for g in groups])
        hds = next(g for g in groups if g['group'] == 'hds')
        self.assertEqual(len(hds['endpoints']), 13)
        self.assertEqual(len({(e['method'], e['path']) for e in hds['endpoints']}), 13)
        self.assertIn('hybrid-data-security', hds['aliases'])

    def test_new_scoped_call_controls(self):
        for name in ['dial', 'answer', 'hangup', 'list', 'get']:
            ep = self.command('Webex Cloud Calling', 'call-controls-for-me', name)
            self.assertIn('/members/me/', ep['path'])

    def test_metrics_move_preserves_analytics_host(self):
        from generate_cli import resolve_base_url
        spec = json.loads((Path(__file__).parent / 'api_spec.json').read_text())
        self.assertFalse(any(g['group'] in ('metrics', 'calling-metrics') for g in spec['Webex Admin']))
        group = next(g for g in spec['Webex Cloud Calling'] if g['group'] == 'metrics')
        ep = group['endpoints'][0]
        self.assertEqual(ep['command'], 'get-call-quality-stats')
        self.assertEqual(resolve_base_url('Webex Cloud Calling', ep['path'], 'wrong-host'),
                         'config.AnalyticsBaseURL')

    def test_queue_names_keep_selection_method_without_parameter_suffixes(self):
        for name in ['list-agent-based', 'list-skill-based', 'list-team-based',
                     'list-by-skill-profile', 'list-by-dynamic-skills', 'list-by-user-skill-profile']:
            self.command('Webex Contact Center', 'contact-service-queue', name)


if __name__ == '__main__':
    unittest.main()
