#!/usr/bin/env python3
"""Generate the Webex API operation table from api_spec.json.

Reads the enriched API spec and writes one pkg/webexops/zz_<area>.go file per
collection, plus zz_areas.go, which lists them. internal/apicmd builds the CLI's
commands from the table.
"""

import json
import os
import re
import subprocess
import paging
from naming_overrides import GROUP_ALIASES, GROUP_DESCRIPTIONS

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
CLI_DIR = os.path.dirname(SCRIPT_DIR)  # repo root (parent of codegen/)
SPEC_PATH = os.path.join(SCRIPT_DIR, "api_spec.json")
# Every exported operation name and the route it was given to; see operation_ids.
# Entries are only ever added.
NAMES_PATH = os.path.join(SCRIPT_DIR, "operation_names.json")
OPS_DIR = os.path.join(CLI_DIR, "pkg", "webexops")

# Collection name → (CLI area, default host, is_calling)
COLLECTIONS = {
    "Webex Cloud Calling": ("calling", "webex", True),
    "Webex Contact Center": ("cc", "cc", False),
    "Webex Admin":          ("admin", "webex", True),
    "Webex Device":         ("device", "webex", True),
    "Webex Meetings":       ("meetings", "webex", True),
    "Webex Messaging":      ("messaging", "webex", True),
}

# webexops host constants, by the host names used here.
HOSTS_GO = {
    'webex': 'HostWebex', 'identity': 'HostIdentity', 'analytics': 'HostAnalytics',
    'analytics-calling': 'HostAnalyticsCalling', 'cc': 'HostCC', 'ai-agent': 'HostAIAgent',
}

# Endpoints served from a different host than the rest of their collection.
# The Postman collections use "{{baseUrl}}" for every request, so the real host
# is not in the source data and has to be recorded here.
# Collection name → URL path prefix → host.
BASE_URL_OVERRIDES = {
    "Webex Cloud Calling": {
        "/v1/analytics/": "analytics",
        # Detailed Call History (CDR) — region-specific analytics FQDN.
        "/cdr_feed": "analytics-calling",
        "/cdr_stream": "analytics-calling",
    },
    "Webex Admin": {
        # Analytics reports. These paths already include /v1, so the host
        # deliberately has no version segment.
        "/v1/analytics/": "analytics",
        # SCIM and archived users are served without the /v1 segment.
        "/identity/scim/": "identity",
        "/identity/organizations/{orgId}/v1/": "identity",
    },
    "Webex Contact Center": {
        # AI Agent Studio, from codegen/supplemental.
        "/api/": "ai-agent",
    },
}

# Routes whose responses carry AI agent runtime tokens, which are credentials
# and must not reach the terminal or an agent transcript.
# Collection name → URL path prefixes.
REDACTED_ROUTES = {
    "Webex Contact Center": ("/api/",),
}

# Routes on APIs that Cisco neither documents nor supports. The CLI cautions
# about them in --help and once when run.
# Collection name → URL path prefixes.
PRIVATE_API_ROUTES = {
    "Webex Contact Center": ("/api/",),
}

# AI Agent Studio selects the agent of a v1 call by the agent's runtime token.
# The CLI looks the token up from --agent-id, so users never pass or see it.
AGENT_TOKEN_HEADER = 'bot-access-token'


def redacts_secrets(collection_name, path):
    return any(path.startswith(p) for p in REDACTED_ROUTES.get(collection_name, ()))


def is_private_api(collection_name, path):
    return any(path.startswith(p) for p in PRIVATE_API_ROUTES.get(collection_name, ()))


def uses_agent_token(ep):
    return any(h['name'] == AGENT_TOKEN_HEADER for h in ep.get('extra_headers', []))


def resolve_base_url(collection_name, path, default):
    """Return the host of an endpoint (longest prefix wins)."""
    overrides = BASE_URL_OVERRIDES.get(collection_name, {})
    match = ""
    for prefix in overrides:
        if path.startswith(prefix) and len(prefix) > len(match):
            match = prefix
    return overrides[match] if match else default


# How each org ID parameter is sent: 'uuid' is converted to the org UUID, and
# 'webex' is sent as given, which is the base64 Webex org ID. Collection name →
# parameter name → form. A parameter whose name looks like an org ID and that is
# not listed stops codegen, because sending the wrong form fails only at run time.
ORG_ID_FORMS = {
    "Webex Contact Center": {"orgId": "uuid", "orgid": "uuid", "organizationId": "uuid", "X-ORGANIZATION-ID": "uuid"},
    "Webex Meetings": {"orgId": "webex", "sessionOrgId": "uuid"},
    "Webex Admin": {"orgId": "webex", "personOrgId": "webex", "onBehalfOfSubPartnerOrgId": "webex"},
    "Webex Cloud Calling": {"orgId": "webex", "organizationId": "webex", "onBehalfOfSubPartnerOrgId": "webex"},
    "Webex Device": {"orgId": "webex"},
    "Webex Messaging": {"orgId": "webex", "organizationId": "webex"},
}

ORG_ID_NAME = re.compile(r'(?i)org(anization)?[-_]?id')

# Routes that take the org UUID as orgId in any collection. The contacts note in
# the Postman collection reads "orgId used in path are the org UUIDs".
UUID_ORG_ROUTE_PREFIXES = ('/identity/scim/', '/identity/organizations/')
UUID_ORG_ROUTE_PART = '/contacts/organizations/'


def takes_org_uuid(collection_name, path, name):
    """Whether a parameter takes the org UUID. Raises ValueError for an org ID
    parameter that ORG_ID_FORMS does not list."""
    if not ORG_ID_NAME.search(name):
        return False
    if name == 'orgId' and (path.startswith(UUID_ORG_ROUTE_PREFIXES) or UUID_ORG_ROUTE_PART in path):
        return True
    form = ORG_ID_FORMS.get(collection_name, {}).get(name)
    if form is None:
        raise ValueError(f'{collection_name} {path}: org ID parameter {name} has no form; '
                         'add it to ORG_ID_FORMS in codegen/generate_cli.py')
    return form == 'uuid'


# Path params that are always the same value — fixed instead of given by a flag.
HARDCODED_PATH_PARAMS = {
    "projectId": "5e5c9ad6d61f870d6d778c1b",
}

# Groups left out of the table — hand-written replacements exist (custom_*.go).
SKIP_GROUPS = {
    "Webex Contact Center": {"search"},
}

# Per-command skip list, keyed by collection → group → set of command names.
# These require multipart/form-data file uploads that the table cannot describe
# (it sends JSON or empty POST/PUT bodies). Each is replaced by a hand-written
# custom_*.go implementation in the area's command package.
SKIP_COMMANDS = {
    # Same JSON flags, but files may require a multipart upload.
    # Replaced by: cmd/messaging/custom_messages.go
    "Webex Messaging": {"messages": {"create"}},
    "Webex Cloud Calling": {
        # Replaced by: cmd/calling/custom_announcements.go
        "announcement-repository": {
            "upload-binary-greeting",       # POST org-level announcement (create)
            "upload-binary-greeting-2",     # POST location-level announcement (create)
            "update-binary-greeting",       # PUT org-level announcement (update)
            "update-binary-greeting-2",     # PUT location-level announcement (update)
        },
        # Replaced by: cmd/calling/custom_voicemail_greetings.go
        "user-call": {
            "update-intercept-greeting-person",           # call intercept announcement
            "update-busy-voicemail-greeting-person",      # busy voicemail greeting
            "update-no-answer-voicemail-greeting-person",  # no-answer voicemail greeting
        },
        # Replaced by: cmd/calling/custom_voicemail_greetings.go
        "call-settings-for-me": {
            "upload-voicemail-busy-greeting",      # self busy voicemail greeting
            "upload-voicemail-no-answer-greeting",  # self no-answer voicemail greeting
        },
        # Replaced by: cmd/calling/custom_voicemail_greetings.go
        "virtual-line-call": {
            "update-intercept-greeting",          # call intercept announcement
            "update-busy-voicemail-greeting",     # busy voicemail greeting
            "update-no-answer-voicemail-greeting",  # no-answer voicemail greeting
        },
        # Replaced by: cmd/calling/custom_voicemail_greetings.go
        "workspace-call": {
            "upload-intercept-announcement-file",     # call intercept announcement
            "update-busy-voicemail-greeting-place",   # busy voicemail greeting
            "update-no-answer-voicemail-greeting-place",  # no-answer voicemail greeting
        },
    },
}


def camel_to_kebab(name):
    """Convert camelCase to kebab-case."""
    s = re.sub(r'([a-z0-9])([A-Z])', r'\1-\2', name)
    return s.lower()


def kebab_to_camel(name):
    """Convert kebab-case to camelCase."""
    parts = name.split('-')
    return parts[0] + ''.join(p.capitalize() for p in parts[1:])


def kebab_to_go_var(name):
    """The key that decides which flags share a value: two parameters whose
    names give the same key share the first one's flag."""
    parts = name.split('-')
    var = parts[0] + ''.join(p.capitalize() for p in parts[1:])
    if var in ('type', 'range', 'map', 'func', 'var', 'return', 'select',
               'default', 'switch', 'case', 'break', 'continue', 'interface',
               'package', 'import', 'go', 'defer', 'chan', 'struct'):
        var += 'Val'
    return var


def kebab_to_pascal(name):
    """Convert kebab-case to PascalCase."""
    return ''.join(p.capitalize() for p in name.split('-'))


def go_string(s):
    """A Go double-quoted string literal."""
    return json.dumps(s, ensure_ascii=False)


def go_identifier(*names):
    """A Go exported identifier from kebab-case names."""
    words = []
    for name in names:
        words.extend(w for w in re.split(r'[^A-Za-z0-9]+', name) if w)
    return ''.join(w[0].upper() + w[1:] for w in words)


# CC API paths where --from/--to expect epoch milliseconds (not ISO 8601).
EPOCH_MS_PATHS = {
    "/v1/queues/statistics",
    "/v1/agents/statistics",
    "/v1/agents/activities",
    "/v1/tasks",
}

# Paging families by host; every other host pages as the main Webex APIs do.
PAGING_FAMILY = {'cc': 'cc', 'ai-agent': 'aiagent'}
PAGING_MODES_GO = {'single': 'webexapi.Single', 'link': 'webexapi.Link', 'pages': 'webexapi.Pages', 'offset': 'webexapi.Offset'}
# Field order of the generated webexapi.Paging literal.
PAGING_FIELDS_GO = (('page', 'Page'), ('start', 'Start'), ('first', 'First'), ('size', 'Size'), ('pageSize', 'PageSize'),
                    ('items', 'Items'), ('omitEmpty', 'OmitEmpty'), ('totalPages', 'TotalPages'),
                    ('totalItems', 'TotalItems'), ('linkQueryOnly', 'LinkQueryOnly'), ('query', 'Query'),
                    ('notFoundEnds', 'NotFoundEnds'))
UNPAGED_LISTS = []
PAGED_ENDPOINTS = []


def paging_literal(ep, host, overrides):
    """The webexapi.Paging literal a list is read with, or None when it is not a
    list or its item field is unknown."""
    family = PAGING_FAMILY.get(host, 'webex')
    PAGED_ENDPOINTS.append((ep, family))
    spec = paging.paging(ep, family, overrides)
    if spec is None:
        return None
    if spec.get('items') is None:
        UNPAGED_LISTS.append(f"{ep['method']} {ep['path']}")
        return None
    fields = [f"Mode: {PAGING_MODES_GO[spec['mode']]}"]
    for key, name in PAGING_FIELDS_GO:
        value = spec.get(key)
        if value in (None, '', 0, False, {}):
            continue
        if value is True:
            fields.append(f'{name}: true')
        elif isinstance(value, int):
            fields.append(f'{name}: {value}')
        elif isinstance(value, dict):
            pairs = ', '.join(f'{go_string(k)}: {go_string(v)}' for k, v in value.items())
            fields.append(f'{name}: map[string]string{{{pairs}}}')
        else:
            fields.append(f'{name}: {go_string(value)}')
    return '&webexapi.Paging{' + ', '.join(fields) + '}'


def operation(ep, collection_name, host, is_calling, paging_overrides=None):
    """The operation-table entry for one endpoint, as a dict."""
    path = ep['path']
    path_params = ep.get('path_params', [])
    query_params = ep.get('query_params', [])
    body_fields = ep.get('body_fields', [])
    complex_body = ep.get('complex_body', False)
    method = ep['method']

    # The org parameter of Contact Center and contacts routes has the flag --orgid,
    # which the CLI fills with the UUID; the parameter keeps its API name. Identity
    # routes keep --org-id, filled with the base64 ID, and convert it.
    needs_uuid_orgid = not is_calling or UUID_ORG_ROUTE_PART in path
    if needs_uuid_orgid:
        for p in path_params + query_params:
            if p['name'] == 'orgId':
                p['_api_name'] = 'orgId'
                p['name'] = 'orgid'

    # Parameters whose names give the same key share the first one's flag.
    flags = {}

    def flag_for(name, flag):
        key = kebab_to_go_var(camel_to_kebab(name))
        return flags.setdefault(key, flag)

    def org_uuid(p):
        return takes_org_uuid(collection_name, path, p.get('_api_name', p['name']))

    params = []
    for p in path_params:
        if p['name'] in HARDCODED_PATH_PARAMS:
            params.append({'name': p['name'], 'in': 'InPath', 'value': HARDCODED_PATH_PARAMS[p['name']]})
            continue
        params.append({
            'name': p.get('_api_name', p['name']), 'in': 'InPath', 'flag': flag_for(p['name'], camel_to_kebab(p['name'])),
            'usage': p.get('description', ''), 'required': True, 'orgUUID': org_uuid(p),
        })
    for p in query_params:
        params.append({
            'name': p.get('_api_name', p['name']), 'in': 'InQuery',
            'flag': flag_for(p['name'], camel_to_kebab(p['name'])), 'usage': p.get('description', ''),
            'orgUUID': org_uuid(p),
        })
    for h in ep.get('extra_headers', []):
        if h['name'] == AGENT_TOKEN_HEADER:
            # Shares the --agent-id flag with an {agentId} path parameter when there is one.
            params.append({'name': h['name'], 'in': 'InAgentToken', 'flag': flags.setdefault('agentId', 'agent-id'),
                           'usage': h.get('description', ''), 'required': True})
            continue
        params.append({'name': h['name'], 'in': 'InHeader', 'flag': flag_for(h['name'], camel_to_kebab(h['name'])),
                       'usage': h.get('description', ''), 'orgUUID': org_uuid(h)})

    body = None
    if ep.get('has_body', False):
        fields = []
        if not complex_body:
            for f in body_fields:
                key = kebab_to_go_var(camel_to_kebab(f['name']))
                if key in flags:
                    # A field that repeats a parameter's name sends that parameter's
                    # value, when both are strings.
                    if f['type'] == 'string':
                        fields.append({'name': f['name'], 'flag': flags[key], 'type': 'string'})
                    continue
                flags[key] = camel_to_kebab(f['name'])
                fields.append({'name': f['name'], 'flag': flags[key], 'type': f['type']})
        body = fields

    time_range = None
    # --last sets both ends of a range, so it needs both parameters.
    if {'from', 'to'} <= {p['name'] for p in query_params}:
        time_range = 'TimeEpochMs' if path in EPOCH_MS_PATHS else 'TimeISO'

    return {
        'command': ep['command'], 'aliases': ep.get('aliases') or [],
        'summary': ep.get('original_name', ep['command']), 'description': ep.get('description', ''),
        'method': method, 'host': host, 'path': path, 'params': params, 'body': body,
        'paging': paging_literal(ep, host, paging_overrides), 'timeRange': time_range,
        # This endpoint returns a non-JSON deletion acknowledgement. Restricting
        # Accept to application/json causes the service to reject it with HTTP 406.
        'accept': '*/*' if method == 'DELETE' and path == '/v1/usage-reports/{reportId}' else '',
        'redact': redacts_secrets(collection_name, ep['path']),
        'private': is_private_api(collection_name, ep['path']),
    }


def route_word(segment):
    """A path segment as part of a Go name: {callParkId} is ByCallParkId."""
    if segment.startswith('{'):
        return 'By' + go_identifier(segment.strip('{}'))
    return go_identifier(segment)


def route_name(area, method, path):
    """The Go name of an operation: its area, method and path, such as
    CcGetOrganizationByOrgidV2Team."""
    return go_identifier(area) + method.capitalize() + ''.join(route_word(s) for s in path.strip('/').split('/'))


def request_shape(ep):
    """What a listing says a request carries: its parameter and body field names."""
    return (tuple(sorted(p['name'] for p in ep.get('path_params', []))),
            tuple(sorted(p['name'] for p in ep.get('query_params', []))),
            tuple(sorted(h['name'] for h in ep.get('extra_headers', []))),
            tuple(sorted((f['name'], f['type']) for f in ep.get('body_fields', []))),
            bool(ep.get('has_body')), bool(ep.get('complex_body')))


def operation_ids(area, endpoints, names=None):
    """The Go names of an area's listings, in order, and the routes to export as
    one operation combining their listings: exported name → listing indexes.

    A name is made from the operation's route alone, so it means the same request
    whatever the CLI commands are called and in whatever order the collection lists
    them. A route the collection lists more than once is exported once, on its first
    listing, when every listing sends the same request; otherwise it is exported as
    an operation combining them, which is not a CLI command. The other listings are
    unexported, for the CLI's commands.

    names maps each name ever exported to its route ("GET /path"). A name given to
    one route is never exported for another, since programs built against it would
    call the new route. Routes spelled alike that no name was given to are not
    exported, since neither name would say which route it is."""
    names = names if names is not None else {}
    routes = {}
    for i, ep in enumerate(endpoints):
        routes.setdefault((ep['method'], ep['path']), []).append(i)
    spellings = {}
    for method, path in routes:
        spellings.setdefault(route_name(area, method, path), []).append((method, path))
    ids = [None] * len(endpoints)
    combined = {}
    for (method, path), listings in routes.items():
        name = route_name(area, method, path)
        unexported = f'{name[0].lower()}{name[1:]}'
        owner = names.get(name)
        blocked = None
        if owner is None and len(spellings[name]) > 1:
            blocked = f"{' and '.join(f'{m} {p}' for m, p in spellings[name])} are all spelled {name}; none is exported"
        elif owner is not None and owner != f'{method} {path}':
            blocked = f'{method} {path} is spelled {name}, which {owner} was given; it is not exported'
        if blocked:
            print(f'  WARNING: {blocked}')
            route = spellings[name].index((method, path)) + 1
            for n, i in enumerate(listings, 1):
                ids[i] = f'{unexported}_route{route}_{n}'
            continue
        one_request = len({request_shape(endpoints[i]) for i in listings}) == 1
        for n, i in enumerate(listings, 1):
            ids[i] = name if one_request and n == 1 else f'{unexported}_{n}'
        if not one_request:
            combined[name] = listings
    return ids, combined


def combine_operations(ops):
    """One operation with every listing's parameters and body fields, or None
    when two listings put a parameter in different places or give a body field
    different types."""
    first = ops[0]
    params, places = [], {}
    for op in ops:
        for p in op['params']:
            if p['name'] in places:
                if places[p['name']] != p['in']:
                    return None
                continue
            places[p['name']] = p['in']
            params.append(dict(p, required=p.get('required', False) and all(
                any(q['name'] == p['name'] and q.get('required') for q in other['params']) for other in ops)))
    body, types = None, {}
    for op in ops:
        if op['body'] is None:
            continue
        body = body if body is not None else []
        for f in op['body']:
            if types.setdefault(f['name'], f['type']) != f['type']:
                return None
            if not any(b['name'] == f['name'] for b in body):
                body.append(f)
    pagings = {op['paging'] for op in ops if op['paging']}
    ranges = {op['timeRange'] for op in ops if op['timeRange']}
    if len(pagings) > 1 or len(ranges) > 1:
        return None
    return dict(first, command='', aliases=[], params=params, body=body,
                paging=next(iter(pagings), None), timeRange=next(iter(ranges), None),
                redact=any(op['redact'] for op in ops), private=any(op['private'] for op in ops))


def operation_literal(op_id, area, group, op):
    """Go source declaring one operation."""
    lines = [f'// {op_id} is {go_string(op["summary"])}: {op["method"]} {op["path"]}.',
             f'var {op_id} = &Operation{{',
             f'\tID: {go_string(op_id)}, Area: {go_string(area)}, Group: {go_string(group)}, Command: {go_string(op["command"])},']
    if op['aliases']:
        lines.append('\tAliases: []string{' + ', '.join(go_string(a) for a in op['aliases']) + '},')
    lines.append(f'\tSummary: {go_string(op["summary"])},')
    if op['description']:
        lines.append(f'\tDescription: {go_string(op["description"])},')
    lines.append(f'\tMethod: {go_string(op["method"])}, Host: {HOSTS_GO[op["host"]]}, Path: {go_string(op["path"])},')
    if op['params']:
        lines.append('\tParams: []Param{')
        for p in op['params']:
            fields = [f'Name: {go_string(p["name"])}', f'In: {p["in"]}']
            if p.get('value'):
                fields.append(f'Value: {go_string(p["value"])}')
            else:
                fields.append(f'Flag: {go_string(p["flag"])}')
            if p.get('required'):
                fields.append('Required: true')
            if p.get('orgUUID'):
                fields.append('OrgUUID: true')
            if p.get('usage'):
                fields.append(f'Usage: {go_string(p["usage"])}')
            lines.append('\t\t{' + ', '.join(fields) + '},')
        lines.append('\t},')
    if op['body'] is not None:
        if op['body']:
            lines.append('\tBody: &Body{Fields: []BodyField{')
            for f in op['body']:
                lines.append(f'\t\t{{Name: {go_string(f["name"])}, Flag: {go_string(f["flag"])}, Type: {go_string(f["type"])}}},')
            lines.append('\t}},')
        else:
            lines.append('\tBody: &Body{},')
    if op['paging']:
        lines.append(f'\tPaging: {op["paging"]},')
    if op['timeRange']:
        lines.append(f'\tTimeRange: {op["timeRange"]},')
    if op['accept']:
        lines.append(f'\tAccept: {go_string(op["accept"])},')
    if op['redact']:
        lines.append('\tRedact: true,')
    if op['private']:
        lines.append('\tPrivate: true,')
    lines.append('}')
    return lines


def generate_area_file(area, collection_name, host, is_calling, groups, paging_overrides=None, names=None):
    """Go source for one area's groups and operations, and its counts. names is
    operation_ids' record of exported names, and gains this area's."""
    names = names if names is not None else {}
    group_lines, kept = [], []
    skip_groups = SKIP_GROUPS.get(collection_name, set())
    for group_data in groups:
        group = group_data['group']
        endpoints = group_data['endpoints']
        if not endpoints:
            continue
        if group in skip_groups:
            print(f"  Skipping group '{group}' (hand-written override)")
            continue
        aliases = GROUP_ALIASES.get(collection_name, {}).get(group, [])
        summary = GROUP_DESCRIPTIONS.get(collection_name, {}).get(group, f'{kebab_to_pascal(group)} commands')
        fields = [f'Area: {go_string(area)}', f'Name: {go_string(group)}']
        if aliases:
            fields.append('Aliases: []string{' + ', '.join(go_string(a) for a in aliases) + '}')
        fields.append(f'Summary: {go_string(summary)}')
        group_lines.append('\t{' + ', '.join(fields) + '},')

        skip_cmds = SKIP_COMMANDS.get(collection_name, {}).get(group, set())
        for ep in endpoints:
            if ep['command'] in skip_cmds:
                print(f"    Skipping command '{ep['command']}' (custom override)")
            else:
                kept.append((group, ep))

    op_ids, combined = operation_ids(area, [ep for _, ep in kept], names)
    decls, ops = [], []
    for (group, ep), op_id in zip(kept, op_ids):
        host_of = resolve_base_url(collection_name, ep['path'], host)
        op = operation(ep, collection_name, host_of, is_calling, paging_overrides)
        ops.append(op)
        if op_id[0].isupper():
            names[op_id] = f"{op['method']} {op['path']}"
        decls.append('')
        decls.extend(operation_literal(op_id, area, group, op))
    for name, listings in combined.items():
        op = combine_operations([ops[i] for i in listings])
        if op is None:
            print(f"  WARNING: the listings of {ops[listings[0]]['method']} {ops[listings[0]]['path']} "
                  "disagree on a parameter's place or a field's type; it is not exported")
            continue
        names[name] = f"{op['method']} {op['path']}"
        literal = operation_literal(name, area, kept[listings[0]][0], op)
        literal.insert(1, '// It has the parameters of every listing of the route, and is not a CLI command.')
        decls.append('')
        decls.extend(literal)

    lines = ['// Code generated by codegen/generate_cli.py; DO NOT EDIT.', '',
             'package webexops', '', 'import "github.com/Cloverhound/webex-cli/pkg/webexapi"', '',
             'var _ = webexapi.Single', '',
             f'// {area}Groups are the {collection_name} command groups.',
             f'var {area}Groups = []Group{{', *group_lines, '}', '',
             f'// {area}Operations are the {collection_name} operations, in collection order.',
             f'var {area}Operations = []*Operation{{', *[f'\t{i},' for i in op_ids], '}']
    lines.extend(decls)
    return '\n'.join(lines) + '\n', len(group_lines), len(op_ids)


def generate_areas_file(areas):
    """Go source listing every area's groups and operations."""
    lines = ['// Code generated by codegen/generate_cli.py; DO NOT EDIT.', '', 'package webexops', '',
             'var areas = map[string]area{']
    for area in areas:
        lines.append(f'\t{go_string(area)}: {{{area}Groups, {area}Operations}},')
    lines.append('}')
    return '\n'.join(lines) + '\n'


def load_names(path):
    """The record of exported names. It cannot be rebuilt from the collections,
    so a missing file stops codegen rather than starting a new record."""
    if not os.path.exists(path):
        raise SystemExit(f"{path} is missing; restore it from git before running codegen.")
    with open(path) as f:
        return json.load(f)


def main():
    with open(SPEC_PATH) as f:
        spec = json.load(f)
    names = load_names(NAMES_PATH)

    for f in os.listdir(OPS_DIR):
        if f.startswith('zz_') and f.endswith('.go'):
            os.remove(os.path.join(OPS_DIR, f))

    areas = []
    for collection_name, groups in spec.items():
        if collection_name not in COLLECTIONS:
            print(f"Skipping unknown collection: {collection_name}")
            continue
        area, host, is_calling = COLLECTIONS[collection_name]
        source, group_count, op_count = generate_area_file(area, collection_name, host, is_calling, groups,
                                                           names=names)
        with open(os.path.join(OPS_DIR, f'zz_{area}.go'), 'w') as f:
            f.write(source)
        areas.append(area)
        print(f"{collection_name}: {group_count} groups, {op_count} operations → pkg/webexops/zz_{area}.go")

    with open(os.path.join(OPS_DIR, 'zz_areas.go'), 'w') as f:
        f.write(generate_areas_file(sorted(areas)))
    with open(NAMES_PATH, 'w') as f:
        json.dump(dict(sorted(names.items())), f, indent=1)
        f.write('\n')

    print("\nRunning gofmt...")
    result = subprocess.run(['gofmt', '-w', OPS_DIR], capture_output=True, text=True)
    if result.returncode != 0:
        print(f"gofmt error: {result.stderr}")
    else:
        print("  gofmt pkg/webexops/ OK")

    if UNPAGED_LISTS:
        print("\nWARNING: --paginate is off for these lists, whose item field neither the collection nor")
        print("codegen/paging_overrides.json gives. Measure them and add an override:")
        for route in sorted(set(UNPAGED_LISTS)):
            print(f"  {route}")
    redundant = paging.redundant_overrides(PAGED_ENDPOINTS)
    if redundant:
        print("\nWARNING: these codegen/paging_overrides.json entries match the derivation, or no")
        print("command uses them. Remove them:")
        for route in redundant:
            print(f"  {route}")

    print("\nDone. Run 'go build ./...' to verify.")


if __name__ == "__main__":
    main()
