"""Checks for merging hand-maintained supplemental collections into downloaded ones."""
import json
import os
import tempfile
import unittest

from extract_api_spec import example_list, extract_endpoint, load_collections


class HeaderTest(unittest.TestCase):
    def test_headers_the_transport_owns_are_never_parameters(self):
        headers = [{'key': k} for k in ('authorization', 'Cookie', 'HOST', 'content-length', 'Proxy-Authorization',
                                        'content-type', 'accept', 'TrackingID', 'password')]
        ep = extract_endpoint({'name': 'Get', 'request': {'method': 'GET', 'header': headers,
                                                          'url': {'path': ['meetings']}}})
        self.assertEqual([h['name'] for h in ep['extra_headers']], ['TrackingID', 'password'])


def collection(name, folders):
    return {'info': {'name': name}, 'item': [
        {'name': folder, 'item': [
            {'name': title, 'request': {'method': method, 'url': {'path': path.strip('/').split('/') + ([''] if path.endswith('/') else [])}}}
            for title, method, path in requests
        ]} for folder, requests in folders.items()
    ]}


class SupplementalCollectionTest(unittest.TestCase):
    def setUp(self):
        self.dir = tempfile.TemporaryDirectory()
        self.downloaded = os.path.join(self.dir.name, 'postman')
        self.supplemental = os.path.join(self.dir.name, 'supplemental')
        os.makedirs(self.downloaded)
        os.makedirs(self.supplemental)

    def tearDown(self):
        self.dir.cleanup()

    def write(self, folder, filename, data):
        with open(os.path.join(folder, filename), 'w') as f:
            json.dump(data, f)

    def test_supplemental_folders_join_their_collection(self):
        self.write(self.downloaded, 'cc.postman_collection.json', collection('Webex Contact Center', {
            'Site': [('List Sites', 'GET', '/organization/:orgId/site')]}))
        self.write(self.supplemental, 'ai.postman_collection.json', collection('Webex Contact Center', {
            'AI Agent': [('Export AI agent', 'GET', '/api/v1/llmagent/export/')]}))
        folders = load_collections(self.downloaded, self.supplemental)['Webex Contact Center']
        self.assertEqual([f['group'] for f in folders], ['site', 'ai-agent'])
        self.assertEqual(folders[1]['endpoints'][0]['path'], '/api/v1/llmagent/export/')

    def test_a_downloaded_route_replaces_the_supplemental_one(self):
        self.write(self.downloaded, 'cc.postman_collection.json', collection('Webex Contact Center', {
            'AI Agent': [('Export AI agent', 'GET', '/api/v1/llmagent/export')]}))
        self.write(self.supplemental, 'ai.postman_collection.json', collection('Webex Contact Center', {
            'AI Agent': [('Export AI agent', 'GET', '/api/v1/llmagent/export/'),
                         ('Publish AI agent', 'POST', '/api/v1/llmagent/makelive/')]}))
        folders = load_collections(self.downloaded, self.supplemental)['Webex Contact Center']
        paths = [e['path'] for f in folders for e in f['endpoints']]
        self.assertEqual(paths, ['/api/v1/llmagent/export', '/api/v1/llmagent/makelive/'])

    def test_supplemental_collection_needs_a_downloaded_one(self):
        self.write(self.supplemental, 'ai.postman_collection.json', collection('Webex AI Agent', {
            'AI Agent': [('Export AI agent', 'GET', '/api/v1/llmagent/export/')]}))
        with self.assertRaises(ValueError):
            load_collections(self.downloaded, self.supplemental)


class ExampleListTest(unittest.TestCase):
    """The item array and total that an endpoint's example response shows."""

    def shape(self, body, code=200):
        return example_list([{'code': code, 'body': json.dumps(body) if not isinstance(body, str) else body}])

    def test_bare_array(self):
        self.assertEqual(self.shape([{'id': 'a'}]), {'items': ''})

    def test_named_array(self):
        self.assertEqual(self.shape({'items': [{'id': 'a'}], 'notFoundIds': []}), {'items': 'items'})
        self.assertEqual(self.shape({'trunks': [{'id': 'a'}]}), {'items': 'trunks'})

    def test_page_count(self):
        self.assertEqual(self.shape({'data': [], 'meta': {'page': 0, 'totalPages': 3}}),
                         {'items': 'data', 'totalPages': 'meta.totalPages'})
        self.assertEqual(self.shape({'data': [], 'pageInfo': {'totalPages': 3}}),
                         {'items': 'data', 'totalPages': 'pageInfo.totalPages'})

    def test_placeholder_totals(self):
        self.assertEqual(self.shape({'campaigns': [], 'pagination': {'page': '<integer>', 'totalPages': '<integer>'}}),
                         {'items': 'campaigns', 'totalPages': 'pagination.totalPages'})
        self.assertEqual(self.shape({'resources': [], 'totalResources': '<long>'}),
                         {'items': 'resources', 'totalItems': 'totalResources'})

    def test_item_count(self):
        self.assertEqual(self.shape({'schemas': ['urn'], 'totalResults': 2, 'Resources': []}),
                         {'items': 'Resources', 'totalItems': 'totalResults'})
        self.assertEqual(self.shape({'resources': [], 'totalResources': 9}),
                         {'items': 'resources', 'totalItems': 'totalResources'})
        self.assertEqual(self.shape({'objects': [], 'meta': {'total_count': 9}}),
                         {'items': 'objects', 'totalItems': 'meta.total_count'})

    def test_no_list(self):
        self.assertIsNone(self.shape({'id': 'a', 'name': 'b'}))
        self.assertIsNone(self.shape({'a': [], 'b': []}))
        self.assertIsNone(self.shape('not json'))
        self.assertIsNone(example_list([]))

    def test_first_success_example_is_used(self):
        self.assertEqual(example_list([{'code': 400, 'body': '{"errors": []}'},
                                       {'code': 200, 'body': '{"items": []}'}]), {'items': 'items'})


if __name__ == '__main__':
    unittest.main()
