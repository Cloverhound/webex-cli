"""How --paginate reads each list endpoint.

Each API family pages one way, so the paging of an endpoint is derived from its
parameters and its example response. paging_overrides.json holds only the
endpoints whose live behavior differs from that derivation; each entry was
measured against a live organization.

A paging spec has a mode:
  single  one request returns the whole list
  link    send the size parameter, leave out the start parameter, and follow
          Link rel="next" until there is none
  pages   page index plus size, until a stated total (totalPages or totalItems)
  offset  start index plus size, until a stated total or a page shorter than asked
and the items field ('' for a bare array). omitEmpty means an empty list may leave
the items field out. pageSize, when set, is the page size to send instead of 100.
query holds parameters the paging needs sent. notFoundEnds means a 404 for a page
after the first ends the list.
"""
import functools
import json
import os

MODES = ('single', 'link', 'pages', 'offset')
PAGING_PARAMS = {'page', 'pageSize', 'size', 'start', 'max', 'offset', 'limit', 'startIndex', 'count'}
LINK_START_PARAMS = ('start', 'offset', 'startIndex')
# Customer Journey Data Services, under the Contact Center host.
JOURNEY_PATHS = ('/admin/v1/api/', '/v1/api/')
OVERRIDES_FILE = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'paging_overrides.json')


def is_list(ep):
    """--paginate applies to a GET named as a list or taking paging parameters.
    limit alone does not make a list: an object can page a field with it."""
    if ep['method'] != 'GET':
        return False
    command = ep.get('command', '')
    return command == 'list' or command.startswith('list-') or bool(query_names(ep) & (PAGING_PARAMS - {'limit'}))


def query_names(ep):
    return {q['name'] for q in ep.get('query_params', [])}


def derive(ep, family):
    """The paging an endpoint of an API family has by that family's convention.

    family is 'webex' (webexapis.com and its analytics hosts), 'cc' (Contact
    Center), or 'aiagent' (AI Agent Studio)."""
    if not is_list(ep):
        return None
    q = query_names(ep)
    example_list = ep.get('example_list') or {}
    example = example_list.get('items')
    if family == 'cc':
        if 'page' in q and q & {'pageSize', 'size'}:
            size = 'pageSize' if 'pageSize' in q else 'size'
            spec = {'mode': 'pages', 'page': 'page', 'size': size, 'items': example if example is not None else 'data'}
            # Journey states no total, so a page shorter than asked ends the
            # list; it answers 404 for a page past the last.
            if ep['path'].startswith(JOURNEY_PATHS):
                spec['notFoundEnds'] = True
                return spec
            if 'totalPages' in example_list:
                spec['totalPages'] = example_list['totalPages']
            elif spec['items'] == 'resources':
                spec['totalItems'] = 'totalResources'
            else:
                spec['totalPages'] = 'pageInfo.totalPages' if size == 'size' else 'meta.totalPages'
            # The flow store states its page count only when asked to, and then
            # wraps its items in data.
            if 'includePagination' in q:
                spec['items'] = 'data'
                spec['query'] = {'includePagination': 'true'}
            return spec
        if not q & PAGING_PARAMS:
            return {'mode': 'single', 'items': example}
        return None
    if family == 'aiagent':
        if {'offset', 'limit'} <= q:
            return {'mode': 'offset', 'start': 'offset', 'size': 'limit', 'items': example or 'objects', 'totalItems': 'meta.total_count'}
        if not q & PAGING_PARAMS:
            return {'mode': 'single', 'items': example or 'objects'}
        return None
    # Webex: lists leave out an empty items field, and page by Link headers,
    # except SCIM, which states totalResults.
    if '/identity/scim/' in ep['path'] and {'startIndex', 'count'} <= q:
        return {'mode': 'offset', 'start': 'startIndex', 'first': 1, 'size': 'count', 'items': example or 'Resources', 'totalItems': 'totalResults'}
    for size in ('max', 'count', 'limit'):
        if size in q:
            return with_link_start(ep, {'mode': 'link', 'size': size, 'items': example, 'omitEmpty': True})
    if not q & PAGING_PARAMS:
        return {'mode': 'single', 'items': example, 'omitEmpty': True}
    return None


def with_link_start(ep, spec):
    """A Link spec names the endpoint's index parameter, which --paginate leaves
    out so the list is read from its start. It is a documented parameter, not
    measured behavior, so an override does not repeat it."""
    if spec is None or spec.get('mode') != 'link' or 'start' in spec:
        return spec
    for start in LINK_START_PARAMS:
        if start in query_names(ep):
            return {'mode': 'link', 'start': start, **{k: v for k, v in spec.items() if k != 'mode'}}
    return spec


@functools.lru_cache(maxsize=None)
def load_overrides(path=OVERRIDES_FILE):
    with open(path) as f:
        return {k: v for k, v in json.load(f).items() if not k.startswith('_')}


def paging(ep, family, overrides=None):
    """The paging --paginate uses: the endpoint's override, or its derivation."""
    if overrides is None:
        overrides = load_overrides()
    route = f"{ep['method']} {ep['path']}"
    if route in overrides:
        return with_link_start(ep, overrides[route])
    return derive(ep, family)


def redundant_overrides(endpoints, overrides=None):
    """Routes whose override equals the derivation of every endpoint with that
    route, given (endpoint, family) pairs. A route can appear in more than one
    collection, with different examples."""
    if overrides is None:
        overrides = load_overrides()
    needed = {route: False for route in overrides}
    for ep, family in endpoints:
        route = f"{ep['method']} {ep['path']}"
        if route in needed and with_link_start(ep, overrides[route]) != derive(ep, family):
            needed[route] = True
    return sorted(route for route, used in needed.items() if not used)
