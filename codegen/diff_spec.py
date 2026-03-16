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
    """Load spec and return {collection: {group: {command: method}}}."""
    with open(path) as f:
        raw = json.load(f)

    result = {}
    for coll, groups in raw.items():
        result[coll] = {}
        for g in groups:
            group = g["group"]
            result[coll][group] = {}
            for ep in g["endpoints"]:
                result[coll][group][ep["command"]] = ep["method"]
    return result


def cli_path(coll, group, command=None):
    sub = COLLECTION_CLI.get(coll, coll.lower())
    if command:
        return f"webex {sub} {group} {command}"
    return f"webex {sub} {group}"


def diff_specs(old, new):
    sections = []

    all_colls = sorted(set(list(old.keys()) + list(new.keys())))
    for coll in all_colls:
        old_groups = old.get(coll, {})
        new_groups = new.get(coll, {})

        added_groups = sorted(set(new_groups) - set(old_groups))
        removed_groups = sorted(set(old_groups) - set(new_groups))
        common_groups = sorted(set(old_groups) & set(new_groups))

        lines = []

        for g in added_groups:
            cmds = new_groups[g]
            lines.append(f"- **Added group** `{cli_path(coll, g)}` ({len(cmds)} commands)")
            for cmd, method in sorted(cmds.items()):
                lines.append(f"  - `{cli_path(coll, g, cmd)}` ({method})")

        for g in removed_groups:
            cmds = old_groups[g]
            lines.append(f"- **Removed group** `{cli_path(coll, g)}` ({len(cmds)} commands)")

        for g in common_groups:
            added_cmds = sorted(set(new_groups[g]) - set(old_groups[g]))
            removed_cmds = sorted(set(old_groups[g]) - set(new_groups[g]))
            if not added_cmds and not removed_cmds:
                continue
            lines.append(f"- `{cli_path(coll, g)}`")
            for cmd in added_cmds:
                lines.append(f"  - **Added** `{cmd}` ({new_groups[g][cmd]})")
            for cmd in removed_cmds:
                lines.append(f"  - **Removed** `{cmd}` ({old_groups[g][cmd]})")

        if lines:
            sub = COLLECTION_CLI.get(coll, coll)
            sections.append(f"### {coll} (`{sub}`)\n")
            sections.extend(lines)
            sections.append("")

    return "\n".join(sections)


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
