#!/usr/bin/env python3
"""Extract enriched API spec from Postman collections for Go CLI code generation.

Reuses normalization logic from extract_api_outline.py and adds:
- Full HTTP method & URL path
- Path params, query params with descriptions
- Request body field extraction with Go types
- Response codes, descriptions, OAuth scopes, extra headers
- CC version dedup (keep highest version per duplicate)
"""

import json
import glob
import os
import re
from naming_overrides import NAMING_OVERRIDES, GROUP_ALIASES, COLLECTION_MOVES, SUPERSEDED_ROUTES

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))

# Import normalization from extract_api_outline
from extract_api_outline import (
    normalize_folder, resolve_folder, normalize_endpoint,
    dedup_commands, PHRASE_ABBREVIATIONS, VERB_MAP, FILLER,
)


# ── Body parsing ─────────────────────────────────────────────────────

POSTMAN_TYPE_MAP = {
    '<string>': 'string',
    '<number>': 'int64',
    '<integer>': 'int64',
    '<boolean>': 'bool',
    '<uuid>': 'string',
    '<float>': 'float64',
    '<double>': 'float64',
    '<long>': 'int64',
    '<date>': 'string',
    '<dateTime>': 'string',
}


def infer_go_type(value):
    """Map a Postman placeholder or literal to a Go type."""
    if isinstance(value, str):
        v = value.strip()
        if v in POSTMAN_TYPE_MAP:
            return POSTMAN_TYPE_MAP[v]
        # Enum-like: literal string value → string
        return 'string'
    if isinstance(value, bool):
        return 'bool'
    if isinstance(value, int):
        return 'int64'
    if isinstance(value, float):
        return 'float64'
    if isinstance(value, list):
        return '[]string'  # simplified
    if isinstance(value, dict):
        return 'object'
    if value is None:
        return 'string'
    return 'string'


def parse_body_fields(raw_body):
    """Parse Postman body JSON into typed field list.

    Returns (fields, complex_body) where:
    - fields: list of {name, type, description} for flat/shallow fields
    - complex_body: True if body has deep nesting requiring --body/--body-file only
    """
    if not raw_body or not raw_body.strip():
        return [], False

    try:
        body = json.loads(raw_body)
    except json.JSONDecodeError:
        return [], True  # unparseable → complex

    if not isinstance(body, dict):
        return [], True  # top-level array or primitive

    fields = []
    has_complex = False

    for key, value in body.items():
        go_type = infer_go_type(value)

        if go_type == 'object':
            # Nested object → mark complex
            has_complex = True
            continue
        elif isinstance(value, list):
            if value and isinstance(value[0], dict):
                # Array of objects → complex
                has_complex = True
                continue
            else:
                go_type = '[]string'

        fields.append({
            'name': key,
            'type': go_type,
        })

    # If body has ONLY complex fields and no simple ones, mark whole body complex
    if not fields and has_complex:
        return [], True

    return fields, has_complex


# ── URL path normalization ───────────────────────────────────────────

def normalize_url_path(path_segments):
    """Convert Postman path segments to a clean URL path.

    :var → {var} style path parameters.
    {{var}} → {var} (Postman env variables used as path params).
    {{var}}:action → {var}:action (preserves action suffixes like :export).
    Strips {{baseUrl}} prefix.
    """
    parts = []
    for seg in path_segments:
        if not seg:
            continue
        if seg == '{{baseUrl}}':
            continue
        if seg.startswith('{{'):
            # {{var}} or {{var}}:action → {var} or {var}:action
            m = re.match(r'^\{\{(\w+)\}\}(.*)$', seg)
            if m:
                parts.append('{' + m.group(1) + '}' + m.group(2))
            continue
        if seg.startswith(':'):
            parts.append('{' + seg[1:] + '}')
        else:
            parts.append(seg)
    return '/' + '/'.join(parts)


def extract_path_version(path_segments):
    """Extract API version from URL path segments. Returns int (0 = unversioned)."""
    for seg in path_segments:
        m = re.match(r'^v(\d+)$', seg)
        if m:
            return int(m.group(1))
    return 0


# ── Scope extraction ─────────────────────────────────────────────────

def extract_scopes(description):
    """Extract OAuth scopes from description text."""
    if not description:
        return []
    return re.findall(r'spark(?:-admin)?:\w+', description)


# ── Endpoint extraction ──────────────────────────────────────────────

def extract_endpoint(item):
    """Extract full endpoint spec from a Postman item."""
    req = item.get('request', {})
    if not req:
        return None

    method = req.get('method', 'GET').upper()
    url = req.get('url', {})
    path_segments = url.get('path', [])
    url_path = normalize_url_path(path_segments)
    version = extract_path_version(path_segments)

    # Path params
    path_params = []
    for var in url.get('variable', []):
        path_params.append({
            'name': var.get('key', ''),
            'description': var.get('description', ''),
        })

    # Also extract {{var}} path params not listed in url.variable
    # (Postman sometimes uses env variable syntax for path params like {{flowId}})
    existing_param_names = {p['name'] for p in path_params}
    for seg in path_segments:
        if seg.startswith('{{') and seg != '{{baseUrl}}':
            m = re.match(r'^\{\{(\w+)\}\}', seg)
            if m and m.group(1) not in existing_param_names:
                path_params.append({'name': m.group(1), 'description': ''})
                existing_param_names.add(m.group(1))

    # Query params
    query_params = []
    for q in url.get('query', []):
        if q.get('disabled'):
            continue
        query_params.append({
            'name': q.get('key', ''),
            'description': q.get('description', ''),
        })

    # Body
    body_raw = req.get('body', {}).get('raw', '') if req.get('body') else ''
    body_fields, complex_body = parse_body_fields(body_raw)

    # Description
    description = req.get('description', '') or ''

    # Scopes
    scopes = extract_scopes(description)

    # Response code
    responses = item.get('response', [])
    response_code = responses[0].get('code') if responses else None

    # Extra headers (non-standard ones)
    extra_headers = []
    for h in req.get('header', []):
        key = h.get('key', '')
        if key not in ('Content-Type', 'Accept', 'Authorization'):
            extra_headers.append({
                'name': key,
                'description': h.get('description', ''),
            })

    return {
        'original_name': item['name'],
        'method': method,
        'path': url_path,
        'version': version,
        'path_params': path_params,
        'query_params': query_params,
        'body_fields': body_fields,
        'complex_body': complex_body,
        'has_body': bool(body_raw.strip()),
        'description': description,
        'scopes': scopes,
        'response_code': response_code,
        'extra_headers': extra_headers,
    }


# ── CC version dedup ─────────────────────────────────────────────────

def dedup_version_key(ep):
    """Create a key for version deduplication.

    Normalizes path by stripping /vN/ segment to group same endpoints.
    """
    path = ep['path']
    # Remove /v1/, /v2/, /v3/ etc.
    normalized = re.sub(r'/v\d+/', '/', path)
    return f"{ep['method']}:{normalized}"


def dedup_cc_versions(endpoints, folder_name):
    """Deduplicate versioned CC endpoints.

    Default: keep only highest version.
    Exception — Subscriptions: keep both v1 and v2.
    """
    is_subscriptions = 'subscription' in folder_name.lower()

    # Group by dedup key
    groups = {}
    for ep in endpoints:
        key = dedup_version_key(ep)
        if key not in groups:
            groups[key] = []
        groups[key].append(ep)

    result = []
    for key, eps in groups.items():
        if len(eps) == 1:
            result.append(eps[0])
            continue

        # Multiple versions of same endpoint
        eps.sort(key=lambda e: e['version'])

        if is_subscriptions:
            # Keep all versions, label with version suffix
            for ep in eps:
                v = ep['version']
                ep['_version_suffix'] = f'-v{v}'
                result.append(ep)
        else:
            # Keep only highest version
            result.append(eps[-1])

    return result


# ── Main extraction ──────────────────────────────────────────────────

def extract_collection(filepath):
    """Extract all groups and endpoints from a Postman collection."""
    with open(filepath) as f:
        data = json.load(f)

    collection_name = data['info']['name']
    is_cc = 'contact center' in collection_name.lower()

    folders = []
    for item in data.get('item', []):
        if 'item' not in item:
            continue

        raw_folder = item['name']
        norm_folder = resolve_folder(normalize_folder(raw_folder))

        endpoints = []
        for child in item['item']:
            if 'request' not in child:
                continue
            ep = extract_endpoint(child)
            if ep:
                endpoints.append(ep)

        # CC version dedup
        if is_cc:
            endpoints = dedup_cc_versions(endpoints, raw_folder)

        # Generate command names
        for ep in endpoints:
            cmd = normalize_endpoint(ep['original_name'], norm_folder)
            # Apply version suffix if present (Subscriptions v1/v2)
            if '_version_suffix' in ep:
                cmd += ep['_version_suffix']
                del ep['_version_suffix']
            ep['command'] = cmd

        dedup_commands(endpoints)

        folders.append({
            'original_folder': raw_folder,
            'group': norm_folder,
            'endpoints': endpoints,
        })

    return collection_name, folders


def merge_folders(folders):
    """Merge folders that map to the same group name."""
    merged = {}
    order = []
    for f in folders:
        group = f['group']
        if group not in merged:
            merged[group] = {
                'original_folders': [],
                'group': group,
                'endpoints': [],
            }
            order.append(group)
        merged[group]['original_folders'].append(f['original_folder'])
        merged[group]['endpoints'].extend(f['endpoints'])

    # Dedup commands after merging
    for g in order:
        dedup_commands(merged[g]['endpoints'])

    return [merged[g] for g in order]


def apply_naming_overrides(collection_name, folders):
    """Resolve explicit naming exceptions and identical routes across merged groups."""
    overrides = NAMING_OVERRIDES.get(collection_name, {})
    shortened_groups = set()
    superseded = SUPERSEDED_ROUTES.get(collection_name, {})
    available = {(e['method'], e['path']) for f in folders for e in f['endpoints']}
    merged = {}
    for folder in folders:
        source_group = folder['group']
        for endpoint in folder['endpoints']:
            ep = dict(endpoint)
            key = (source_group, ep['method'], ep['path'])
            if key in superseded:
                if superseded[key] not in available:
                    raise ValueError(f"Missing replacement for {collection_name} {key}")
                continue
            group, command, aliases = overrides.get(
                key, (source_group, ep['command'], ()))
            ep['command'] = command
            ep['aliases'] = list(aliases)
            if command in ('get-id', 'delete-id', 'patch-id', 'update-id'):
                ep['command'] = command[:-3]
                ep['aliases'] = list(dict.fromkeys([command, *aliases]))
                shortened_groups.add(group)
            target = merged.setdefault(group, {
                'group': group, 'original_folders': [], 'endpoints': [],
                'aliases': GROUP_ALIASES.get(collection_name, {}).get(group, []),
            })
            for original in folder['original_folders']:
                if original not in target['original_folders']:
                    target['original_folders'].append(original)
            # Only collapse identical routes when deliberately merging groups.
            # Other collections may document distinct request contracts at one URL.
            if group in GROUP_ALIASES.get(collection_name, {}):
                existing = next((e for e in target['endpoints']
                                 if (e['method'], e['path']) == (ep['method'], ep['path'])), None)
                if existing is not None:
                    for alias in [command, *aliases]:
                        if alias != existing['command'] and alias not in existing['aliases']:
                            existing['aliases'].append(alias)
                    continue
            target['endpoints'].append(ep)
    changed_groups = {value[0] for value in overrides.values()} | shortened_groups
    for group in merged.values():
        if group['group'] not in changed_groups:
            continue
        names = set()
        for ep in group['endpoints']:
            for name in [ep['command'], *ep['aliases']]:
                if name in names:
                    raise ValueError(f"Duplicate command or alias: {collection_name} {group['group']} {name}")
                names.add(name)
    return list(merged.values())


def main():
    pattern = os.path.join(SCRIPT_DIR, "postman", "*.postman_collection.json")
    files = sorted(glob.glob(pattern))

    if not files:
        print("No Postman collection files found.")
        return

    result = {}
    for filepath in files:
        name, raw_folders = extract_collection(filepath)
        folders = merge_folders(raw_folders)
        folders = apply_naming_overrides(name, folders)
        result[name] = folders

        total_groups = len(folders)
        total_eps = sum(len(f['endpoints']) for f in folders)
        print(f"{name}: {total_groups} groups, {total_eps} endpoints")

    for (source, group), destination in COLLECTION_MOVES.items():
        moved = [g for g in result[source] if g['group'] == group]
        if any(g['group'] == group for g in result[destination]):
            raise ValueError(f"Destination group already exists: {destination} {group}")
        result[source] = [g for g in result[source] if g['group'] != group]
        result[destination].extend(moved)

    output_path = os.path.join(SCRIPT_DIR, "api_spec.json")
    with open(output_path, "w") as f:
        json.dump(result, f, indent=2)
    print(f"\nOutput written to: {output_path}")


if __name__ == "__main__":
    main()
