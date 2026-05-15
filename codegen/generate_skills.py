#!/usr/bin/env python3
"""Generate Command Reference sections in skill/<area>/SKILL.md from api_spec.json.

Run after extract_api_spec.py + generate_cli.py as part of `make codegen`.
Adds/updates a sentinel-bounded block at the end of each skill file without
touching the hand-written content above it.
"""

import json
import os
import re
import sys

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
CLI_DIR = os.path.dirname(SCRIPT_DIR)
SPEC_PATH = os.path.join(SCRIPT_DIR, "api_spec.json")
SKILL_DIR = os.path.join(CLI_DIR, "skill")

SENTINEL_START = "<!-- codegen:start -->"
SENTINEL_END = "<!-- codegen:end -->"

COLLECTIONS = {
    "Webex Cloud Calling": "calling",
    "Webex Contact Center": "cc",
    "Webex Admin":          "admin",
    "Webex Device":         "device",
    "Webex Meetings":       "meetings",
    "Webex Messaging":      "messaging",
}

# Commands replaced by hand-written custom_*.go — omit from generated reference
# since their actual flags differ from the Postman spec.
SKIP_COMMANDS = {
    "Webex Cloud Calling": {
        "announcement-repository": {
            "upload-binary-greeting",
            "upload-binary-greeting-2",
            "update-binary-greeting",
            "update-binary-greeting-2",
        },
        "user-call": {
            "update-intercept-greeting-person",
            "update-busy-voicemail-greeting-person",
            "update-no-answer-voicemail-greeting-person",
        },
        "call-settings-for-me": {
            "upload-voicemail-busy-greeting",
            "upload-voicemail-no-answer-greeting",
        },
        "virtual-line-call": {
            "update-intercept-greeting",
            "update-busy-voicemail-greeting",
            "update-no-answer-voicemail-greeting",
        },
        "workspace-call": {
            "upload-intercept-announcement-file",
            "update-busy-voicemail-greeting-place",
            "update-no-answer-voicemail-greeting-place",
        },
    },
}

# Groups skipped entirely during Go codegen — still document them here since
# hand-written custom implementations exist and need to be discoverable.
SKIP_GROUPS_CODEGEN = {
    "Webex Contact Center": {"search"},
}


def camel_to_kebab(name):
    s = re.sub(r'([a-z0-9])([A-Z])', r'\1-\2', name)
    return s.lower()


def method_sort_key(method):
    return {'GET': 0, 'POST': 1, 'PUT': 2, 'PATCH': 3, 'DELETE': 4}.get(method, 5)


def format_flags(ep):
    """Return a markdown string listing all flags for one endpoint."""
    parts = []
    seen = set()

    for p in ep.get('path_params', []):
        flag = camel_to_kebab(p['name'])
        if flag not in seen:
            seen.add(flag)
            parts.append(f"`--{flag}` *(required)*")

    for p in ep.get('query_params', []):
        flag = camel_to_kebab(p['name'])
        if flag not in seen:
            seen.add(flag)
            parts.append(f"`--{flag}`")

    # --last is added automatically whenever 'from' is a query param
    if any(p['name'] == 'from' for p in ep.get('query_params', [])):
        parts.append("`--last`")

    if ep.get('has_body'):
        if not ep.get('complex_body'):
            for f in ep.get('body_fields', []):
                flag = camel_to_kebab(f['name'])
                if flag not in seen:
                    seen.add(flag)
                    parts.append(f"`--{flag}`")
        parts.append("`--body`")
        parts.append("`--body-file`")

    return ", ".join(parts) if parts else "—"


def generate_command_reference(area_cmd, groups, collection_name):
    """Build the full ## Command Reference markdown block for one area."""
    skip_groups = SKIP_GROUPS_CODEGEN.get(collection_name, set())
    skip_cmds = SKIP_COMMANDS.get(collection_name, {})

    lines = [
        "## Command Reference",
        "",
        "> Auto-generated from Postman collections. Run `make codegen` to update.",
    ]

    for group_data in groups:
        group = group_data['group']
        endpoints = group_data.get('endpoints', [])

        if not endpoints or group in skip_groups:
            continue

        group_skip = skip_cmds.get(group, set())
        visible = [ep for ep in endpoints if ep['command'] not in group_skip]
        if not visible:
            continue

        visible.sort(key=lambda ep: method_sort_key(ep['method']))

        lines += [
            "",
            f"### {group}",
            "",
            "| Command | Flags |",
            "|---|---|",
        ]
        for ep in visible:
            lines.append(f"| `{ep['command']}` | {format_flags(ep)} |")

    return "\n".join(lines) + "\n"


def update_skill_file(skill_path, generated_block):
    """Replace or append the sentinel block in an existing skill file.

    Returns True if the file was modified, False if content was unchanged.
    """
    if not os.path.exists(skill_path):
        print(f"  WARNING: {skill_path} not found — skipping")
        return False

    with open(skill_path) as f:
        content = f.read()

    block = f"{SENTINEL_START}\n{generated_block}\n{SENTINEL_END}"

    start_idx = content.find(SENTINEL_START)
    end_idx = content.find(SENTINEL_END)

    if start_idx != -1 and end_idx != -1:
        new_content = content[:start_idx] + block + content[end_idx + len(SENTINEL_END):]
    else:
        new_content = content.rstrip("\n") + "\n\n" + block + "\n"

    if new_content == content:
        return False

    with open(skill_path, 'w') as f:
        f.write(new_content)
    return True


def main():
    if not os.path.exists(SPEC_PATH):
        print("ERROR: api_spec.json not found — run 'make download' first.")
        sys.exit(1)

    with open(SPEC_PATH) as f:
        spec = json.load(f)

    for collection_name, area_cmd in COLLECTIONS.items():
        if collection_name not in spec:
            print(f"  {collection_name}: not in spec — skipping")
            continue

        groups = spec[collection_name]
        block = generate_command_reference(area_cmd, groups, collection_name)

        skill_path = os.path.join(SKILL_DIR, area_cmd, "SKILL.md")
        changed = update_skill_file(skill_path, block)

        cmd_count = sum(len(g.get('endpoints', [])) for g in groups)
        status = "updated" if changed else "unchanged"
        print(f"  {area_cmd}: {len(groups)} groups, {cmd_count} commands → {status}")

    print("\nDone.")


if __name__ == "__main__":
    main()
