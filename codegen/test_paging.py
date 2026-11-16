"""How --paginate reads each list endpoint: derived from the collection, or overridden."""
import unittest

import paging


def ep(path, query=(), example=None, method='GET', command='list'):
    return {'method': method, 'path': path, 'command': command,
            'query_params': [{'name': q} for q in query], 'example_list': example}


class DeriveTest(unittest.TestCase):
    def test_contact_center_pages_end_at_meta_total_pages(self):
        self.assertEqual(paging.derive(ep('/organization/{orgid}/v2/site', ['page', 'pageSize'], {'items': 'data'}), 'cc'),
                         {'mode': 'pages', 'page': 'page', 'size': 'pageSize', 'items': 'data', 'totalPages': 'meta.totalPages'})

    def test_contact_center_uses_the_total_its_example_states(self):
        self.assertEqual(paging.derive(ep('/v3/campaign-management/campaign-groups/{campaignGroupName}/campaigns', ['page', 'pageSize'],
                                          {'items': 'campaigns', 'totalPages': 'pagination.totalPages'}), 'cc'),
                         {'mode': 'pages', 'page': 'page', 'size': 'pageSize', 'items': 'campaigns', 'totalPages': 'pagination.totalPages'})

    def test_contact_center_size_param_ends_at_page_info(self):
        self.assertEqual(paging.derive(ep('/functions', ['page', 'size']), 'cc'),
                         {'mode': 'pages', 'page': 'page', 'size': 'size', 'items': 'data', 'totalPages': 'pageInfo.totalPages'})

    def test_contact_center_include_pagination_is_sent(self):
        # The example is the bare array the flow store returns without it.
        self.assertEqual(paging.derive(ep('/flows', ['page', 'size', 'includePagination'], {'items': ''}), 'cc'),
                         {'mode': 'pages', 'page': 'page', 'size': 'size', 'items': 'data', 'totalPages': 'pageInfo.totalPages',
                          'query': {'includePagination': 'true'}})

    def test_journey_lists_end_at_a_short_page(self):
        for path in ('/admin/v1/api/person/workspace-id/{workspaceId}', '/v1/api/events/workspace-id/{workspaceId}'):
            with self.subTest(path=path):
                self.assertEqual(paging.derive(ep(path, ['page', 'pageSize'], {'items': 'data'}), 'cc'),
                                 {'mode': 'pages', 'page': 'page', 'size': 'pageSize', 'items': 'data', 'notFoundEnds': True})

    def test_contact_center_bulk_export_ends_at_total_resources(self):
        self.assertEqual(paging.derive(ep('/organization/{orgid}/site/bulk-export', ['page', 'pageSize'], {'items': 'resources'}), 'cc'),
                         {'mode': 'pages', 'page': 'page', 'size': 'pageSize', 'items': 'resources', 'totalItems': 'totalResources'})

    def test_webex_lists_follow_links(self):
        self.assertEqual(paging.derive(ep('/people', ['max'], {'items': 'items'}), 'webex'),
                         {'mode': 'link', 'size': 'max', 'items': 'items', 'omitEmpty': True})

    def test_link_names_the_index_parameter_it_leaves_out(self):
        self.assertEqual(paging.derive(ep('/telephony/config/locations/{locationId}/huntGroups', ['start', 'max'], {'items': 'huntGroups'}), 'webex'),
                         {'mode': 'link', 'start': 'start', 'size': 'max', 'items': 'huntGroups', 'omitEmpty': True})
        self.assertEqual(paging.derive(ep('/groups/{groupId}/members', ['startIndex', 'count'], {'items': 'members'}), 'webex'),
                         {'mode': 'link', 'start': 'startIndex', 'size': 'count', 'items': 'members', 'omitEmpty': True})
        self.assertEqual(paging.derive(ep('/adminAudit/events', ['from', 'to', 'max', 'offset'], {'items': 'items'}), 'webex'),
                         {'mode': 'link', 'start': 'offset', 'size': 'max', 'items': 'items', 'omitEmpty': True})

    def test_scim_reads_to_total_results(self):
        self.assertEqual(paging.derive(ep('/identity/scim/{orgId}/v2/Users', ['startIndex', 'count'], {'items': 'Resources'}), 'webex'),
                         {'mode': 'offset', 'start': 'startIndex', 'first': 1, 'size': 'count', 'items': 'Resources', 'totalItems': 'totalResults'})

    def test_ai_agent_studio_reads_to_total_count(self):
        self.assertEqual(paging.derive(ep('/api/v2/organizations/{orgid}/agents/', ['limit', 'offset']), 'aiagent'),
                         {'mode': 'offset', 'start': 'offset', 'size': 'limit', 'items': 'objects', 'totalItems': 'meta.total_count'})

    def test_list_without_paging_params_is_one_request(self):
        self.assertEqual(paging.derive(ep('/licenses', [], {'items': 'items'}), 'webex'), {'mode': 'single', 'items': 'items', 'omitEmpty': True})
        self.assertEqual(paging.derive(ep('/organization/{orgid}/user/with-user-profile', [], {'items': ''}), 'cc'), {'mode': 'single', 'items': ''})

    def test_not_a_list(self):
        self.assertIsNone(paging.derive(ep('/people/{personId}', [], command='get'), 'webex'))
        # One license, with a page of the users assigned it.
        self.assertIsNone(paging.derive(ep('/licenses/{licenseId}', ['includeAssignedTo', 'next', 'limit'], {'items': 'users'}, command='get-license'), 'webex'))
        self.assertIsNone(paging.derive(ep('/people', ['max'], {'items': 'items'}, method='POST'), 'webex'))


class OverrideTest(unittest.TestCase):
    def test_an_override_replaces_the_derived_paging(self):
        overrides = {'GET /telephony/config/jobs/numbers/manageNumbers': {'mode': 'offset', 'start': 'start', 'size': 'max', 'items': 'items', 'pageSize': 50}}
        got = paging.paging(ep('/telephony/config/jobs/numbers/manageNumbers', ['start', 'max'], {'items': 'items'}), 'webex', overrides)
        self.assertEqual(got, {'mode': 'offset', 'start': 'start', 'size': 'max', 'items': 'items', 'pageSize': 50})

    def test_a_link_override_leaves_out_the_endpoints_index_parameter(self):
        overrides = {'GET /telephony/config/jobs/updateRoutingPrefix': {'mode': 'link', 'size': 'max', 'items': 'items', 'omitEmpty': True}}
        got = paging.paging(ep('/telephony/config/jobs/updateRoutingPrefix', ['start'], {'items': 'items'}), 'webex', overrides)
        self.assertEqual(got, {'mode': 'link', 'start': 'start', 'size': 'max', 'items': 'items', 'omitEmpty': True})
        self.assertEqual(overrides['GET /telephony/config/jobs/updateRoutingPrefix'],
                         {'mode': 'link', 'size': 'max', 'items': 'items', 'omitEmpty': True})

    def test_an_override_that_matches_every_derivation_is_redundant(self):
        override = {'mode': 'single', 'items': 'devices', 'omitEmpty': True}
        overrides = {'GET /people': {'mode': 'link', 'size': 'max', 'items': 'items', 'omitEmpty': True},
                     'GET /telephony/config/supportedDevices': override}
        endpoints = [(ep('/people', ['max'], {'items': 'items'}), 'webex'),
                     # One collection's example names the items and another's does not.
                     (ep('/telephony/config/supportedDevices', [], {'items': 'devices'}), 'webex'),
                     (ep('/telephony/config/supportedDevices', [], None), 'webex')]
        self.assertEqual(paging.redundant_overrides(endpoints, overrides), ['GET /people'])

    def test_a_null_override_turns_paging_off(self):
        got = paging.paging(ep('/deviceConfigurations', ['deviceId'], {'items': 'items'}), 'webex', {'GET /deviceConfigurations': None})
        self.assertIsNone(got)

    def test_every_override_names_a_mode_and_its_items(self):
        for route, spec in paging.load_overrides().items():
            if spec is None:
                continue
            with self.subTest(route=route):
                self.assertIn(spec.get('mode'), paging.MODES)
                self.assertIsInstance(spec.get('items'), str)


if __name__ == '__main__':
    unittest.main()
