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

    # Contact Center paths, and /contacts/organizations/ in other collections
    # (Postman note: "orgId used in path are the org UUIDs"), take the org UUID.
    # Their org parameter's flag is --orgid, which the CLI fills with the UUID;
    # the parameter keeps its API name.
    needs_uuid_orgid = not is_calling or '/contacts/organizations/' in path
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

    params = []
    for p in path_params:
        if p['name'] in HARDCODED_PATH_PARAMS:
            params.append({'name': p['name'], 'in': 'InPath', 'value': HARDCODED_PATH_PARAMS[p['name']]})
            continue
        params.append({
            'name': p.get('_api_name', p['name']), 'in': 'InPath', 'flag': flag_for(p['name'], camel_to_kebab(p['name'])),
            'usage': p.get('description', ''), 'required': True,
            # Identity routes take the org UUID, but --org-id is filled with the base64 org ID.
            'orgUUID': (needs_uuid_orgid and p['name'] == 'orgid') or
                       (p['name'] == 'orgId' and path.startswith(('/identity/scim/', '/identity/organizations/'))),
        })
    for p in query_params:
        params.append({
            'name': p.get('_api_name', p['name']), 'in': 'InQuery',
            'flag': flag_for(p['name'], camel_to_kebab(p['name'])), 'usage': p.get('description', ''),
            'orgUUID': needs_uuid_orgid and p['name'] == 'orgid',
        })
    for h in ep.get('extra_headers', []):
        if h['name'] == AGENT_TOKEN_HEADER:
            # Shares the --agent-id flag with an {agentId} path parameter when there is one.
            params.append({'name': h['name'], 'in': 'InAgentToken', 'flag': flags.setdefault('agentId', 'agent-id'),
                           'usage': h.get('description', ''), 'required': True})
            continue
        params.append({'name': h['name'], 'in': 'InHeader', 'flag': flag_for(h['name'], camel_to_kebab(h['name'])),
                       'usage': h.get('description', '')})

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
    if any(p['name'] == 'from' for p in query_params):
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


def operation_ids(area, group, endpoints):
    """The Go names of a group's operations, in order.

    Extraction numbers commands that share a name (update, update-2) in collection
    order, so a reordered collection can swap their numbers. Each member of such a
    family is named instead by the family's command and the part of its route after
    the family's common path, which keeps the name with the route. A route listed
    twice gets _2, _3: those names reach the same endpoint whichever way round they
    fall."""
    families = {}
    for i, ep in enumerate(endpoints):
        families.setdefault(re.sub(r'(-\d+)+$', '', ep['command']), []).append(i)
    ids = [None] * len(endpoints)
    for base, members in families.items():
        if len(members) == 1:
            ids[members[0]] = go_identifier(area, group, endpoints[members[0]]['command'])
            continue
        paths = [endpoints[i]['path'].strip('/').split('/') for i in members]
        common = 0
        while all(len(p) > common for p in paths) and len({p[common] for p in paths}) == 1:
            common += 1
        rests = {i: ''.join(route_word(s) for s in p[common:]) for i, p in zip(members, paths)}
        methods = {}
        for i in members:
            methods.setdefault(rests[i], set()).add(endpoints[i]['method'])
        routes, seen = {}, {}
        for i in members:
            ep = endpoints[i]
            name = go_identifier(area, group, base) + rests[i]
            if len(methods[rests[i]]) > 1:
                name += ep['method'].capitalize()
            route = (ep['method'], ep['path'])
            if routes.setdefault(name, route) != route:
                raise SystemExit(f"{ep['method']} {ep['path']} and {routes[name][0]} {routes[name][1]} would both be "
                                 f"named {name} in pkg/webexops; extend operation_ids in codegen/generate_cli.py.")
            seen[name] = seen.get(name, 0) + 1
            ids[i] = name if seen[name] == 1 else f'{name}_{seen[name]}'
    return ids


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


def generate_area_file(area, collection_name, host, is_calling, groups, paging_overrides=None):
    """Go source for one area's groups and operations, and its counts."""
    used_ids = set()
    group_lines, op_ids, decls = [], [], []
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
        kept = [ep for ep in endpoints if ep['command'] not in skip_cmds]
        for ep, op_id in zip(kept, operation_ids(area, group, kept)):
            host_of = resolve_base_url(collection_name, ep['path'], host)
            op = operation(ep, collection_name, host_of, is_calling, paging_overrides)
            if op_id in used_ids:
                raise SystemExit(f"Two {collection_name} operations would be named {op_id}; the second is "
                                 f"{ep['method']} {ep['path']}. Extend operation_ids in codegen/generate_cli.py.")
            used_ids.add(op_id)
            op_ids.append(op_id)
            decls.append('')
            decls.extend(operation_literal(op_id, area, group, op))

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


def main():
    with open(SPEC_PATH) as f:
        spec = json.load(f)

    for f in os.listdir(OPS_DIR):
        if f.startswith('zz_') and f.endswith('.go'):
            os.remove(os.path.join(OPS_DIR, f))

    areas = []
    for collection_name, groups in spec.items():
        if collection_name not in COLLECTIONS:
            print(f"Skipping unknown collection: {collection_name}")
            continue
        area, host, is_calling = COLLECTIONS[collection_name]
        source, group_count, op_count = generate_area_file(area, collection_name, host, is_calling, groups)
        with open(os.path.join(OPS_DIR, f'zz_{area}.go'), 'w') as f:
            f.write(source)
        areas.append(area)
        print(f"{collection_name}: {group_count} groups, {op_count} operations → pkg/webexops/zz_{area}.go")

    with open(os.path.join(OPS_DIR, 'zz_areas.go'), 'w') as f:
        f.write(generate_areas_file(sorted(areas)))

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
