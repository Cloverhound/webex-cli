#!/usr/bin/env python3
"""Compare two api_spec.json files and output a markdown changelog.

Usage: python3 diff_spec.py old_spec.json new_spec.json
"""

import json
import sys

# Collection name → CLI subcommand
COLLECTION_CLI = {
    "Webex Cloud Calling": "calling",
    "Webex Contact Center": "cc",
    "Webex Admin": "admin",
    "Webex Device": "device",
    "Webex Meetings": "meetings",
    "Webex Messaging": "messaging",
}


def load_spec(path):
    """Load spec and return {collection: {group: {command: endpoint}}}."""
    with open(path) as f:
        raw = json.load(f)

    result = {}
    for coll, groups in raw.items():
        result[coll] = {}
        for g in groups:
            group = g["group"]
            result[coll][group] = {}
            for ep in g["endpoints"]:
                result[coll][group][ep["command"]] = ep
    return result


def cli_path(coll, group, command=None):
    sub = COLLECTION_CLI.get(coll, coll.lower())
    if command:
        return f"webex {sub} {group} {command}"
    return f"webex {sub} {group}"


def endpoint_changes(old, new):
    """Summarize executable contract changes, independently of request titles."""
    changes = []
    for field in ('method', 'path'):
        if old.get(field) != new.get(field):
            changes.append(f"{field}: `{old.get(field)}` → `{new.get(field)}`")
    for field, label in [('path_params', 'path parameters'), ('query_params', 'query parameters'),
                         ('body_fields', 'body fields'), ('extra_headers', 'headers')]:
        before = {p['name']: p for p in old.get(field, [])}
        after = {p['name']: p for p in new.get(field, [])}
        added, removed = sorted(after.keys() - before.keys()), sorted(before.keys() - after.keys())
        modified = sorted(k for k in before.keys() & after.keys() if before[k] != after[k])
        details = []
        if added:
            details.append('added ' + ', '.join(f'`{k}`' for k in added))
        if removed:
            details.append('removed ' + ', '.join(f'`{k}`' for k in removed))
        if modified:
            details.append('updated ' + ', '.join(f'`{k}`' for k in modified))
        if details:
            changes.append(label + ': ' + '; '.join(details))
    for field, label in [('has_body', 'request body support'), ('complex_body', 'JSON body handling'),
                         ('aliases', 'compatibility aliases'), ('scopes', 'documented scopes'),
                         ('response_code', 'documented response code')]:
        if old.get(field, [] if field in ('aliases', 'scopes') else None) != new.get(field, [] if field in ('aliases', 'scopes') else None):
            changes.append(label + ' changed')
    if old.get('description') != new.get('description'):
        changes.append('help/description updated')
    return changes


def diff_specs(old, new):
    sections = []
    for coll in sorted(set(old) | set(new)):
        old_groups, new_groups = old.get(coll, {}), new.get(coll, {})
        lines = []
        for group in sorted(set(old_groups) | set(new_groups)):
            before, after = old_groups.get(group, {}), new_groups.get(group, {})
            entries = []
            for name in sorted(after.keys() - before.keys()):
                entries.append(f"  - **Added** `{name}` ({after[name]['method']})")
            for name in sorted(before.keys() & after.keys()):
                changes = endpoint_changes(before[name], after[name])
                if changes:
                    entries.append(f"  - **Modified** `{name}` — " + '; '.join(changes))
            for name in sorted(before.keys() - after.keys()):
                aliases = [n for n, ep in after.items() if name in ep.get('aliases', [])]
                if aliases:
                    entries.append(f"  - **Renamed** `{name}` → " + ', '.join(f'`{n}`' for n in aliases) + ' (old spelling retained as an alias)')
                else:
                    entries.append(f"  - **Deleted** `{name}` ({before[name]['method']}; command removed, not necessarily the upstream API)")
            if entries:
                status = '**Added group** ' if group not in old_groups else '**Deleted group** ' if group not in new_groups else ''
                lines.extend([f"- {status}`{cli_path(coll, group)}`", *entries, ''])
        if lines:
            sections.extend([f"### {coll} (`{COLLECTION_CLI.get(coll, coll)}`)", '', *lines])
    return '\n'.join(sections)


def main():
    if len(sys.argv) != 3:
        print(f"Usage: {sys.argv[0]} old_spec.json new_spec.json", file=sys.stderr)
        sys.exit(1)

    old = load_spec(sys.argv[1])
    new = load_spec(sys.argv[2])
    output = diff_specs(old, new)

    if output.strip():
        print(output)
    else:
        print("No changes to API groups or endpoints.")


if __name__ == "__main__":
    main()
