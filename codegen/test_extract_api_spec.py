"""Checks for merging hand-maintained supplemental collections into downloaded ones."""
import json
import os
import tempfile
import unittest

from extract_api_spec import load_collections


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


if __name__ == '__main__':
    unittest.main()
