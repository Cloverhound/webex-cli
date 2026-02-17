#!/usr/bin/env python3
"""Extract API calls from Postman collections and map to clean CLI names."""

import json
import glob
import os
import re

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))

# ── Folder normalization ──────────────────────────────────────────────

def normalize_folder(name):
    """Convert Postman folder name to a clean kebab-case CLI group name."""
    s = name

    # Strip "Beta ... With ..." wrappers — keep the core subject
    s = re.sub(r'^Beta\s+', '', s)
    s = re.sub(r'\s+With\s+.*$', '', s, flags=re.IGNORECASE)

    # Strip "Features: " or "Features:  " prefix
    s = re.sub(r'^Features:\s+', '', s)

    # Strip "Location Call Settings: " prefix → keep sub-topic
    s = re.sub(r'^Location Call Settings:\s+', 'Location ', s)

    # Strip trailing " Settings" (but keep if the whole name IS "Settings")
    s = re.sub(r'\s+Settings$', '', s)

    # Strip "(1/2)", "(2/2)" splits
    s = re.sub(r'\s*\(\d+/\d+\)$', '', s)

    # Strip trailing " API"
    s = re.sub(r'\s+API$', '', s)

    # "Journey - Xxx" → "Journey Xxx"
    s = re.sub(r'^Journey\s*-\s*', 'Journey ', s)

    # "Partner Reports/Templates" → "Partner Reports"
    s = s.replace('/', ' ')

    # kebab-case
    s = s.strip()
    s = re.sub(r'[^a-zA-Z0-9]+', '-', s)
    s = re.sub(r'-+', '-', s).strip('-').lower()

    return s


# Merge table for split/duplicate folders that should combine
FOLDER_MERGES = {
    # Calling
    'user-call': 'user-call',
    'workspace-call': 'workspace-call',
    'device-call': 'device-call',
    'call-for-me': 'call-for-me',
    'call-queue-with-playlist': 'call-queue',
    'settings-features-for-barge-in': 'call-settings-for-me',
}

def resolve_folder(normalized):
    """Resolve merged/deduplicated folder names."""
    for prefix, target in FOLDER_MERGES.items():
        if normalized.startswith(prefix):
            return target
    return normalized


# ── Phrase abbreviations applied before tokenization ──────────────────

PHRASE_ABBREVIATIONS = [
    ('Customer Experience Essentials', 'CXE'),
    ('Selective Call Forwarding', 'Selective Forward'),
    ('Selective Call Rejection', 'Selective Reject'),
    ('Selective Call Accept', 'Selective Accept'),
    ('Available Phone Numbers', 'Available Numbers'),
    ('Outgoing Calling Permissions', 'Outgoing Permissions'),
    ('Outgoing Permission Digit Pattern', 'Outgoing Digit Pattern'),
    ('Outgoing Permission Access Code', 'Outgoing Access Code'),
    ('Outgoing Permission Auto Transfer Number', 'Outgoing Auto Transfer'),
    ('Emergency Callback Number', 'Emergency Callback'),
    ('Call Forwarding Settings', 'Call Forward'),
    ('Digit Pattern Category Control Settings', 'Digit Pattern Control'),
    ('Calling Services List', 'Calling Services'),
    ('Preferred Answer Endpoint', 'Preferred Endpoint'),
    ('Call Bridge Warning Tone Settings', 'Bridge Warning Tone'),
    ('Compliance Announcement Setting', 'Compliance Announcement'),
    ('Progressive Profile View', 'Profile View'),
    ('Progressive profile Views', 'Profile Views'),
    ('Secondary Line Owner', 'Secondary Owner'),
    ('Caller ID Information', 'Caller ID'),
    ('Caller ID Settings', 'Caller ID'),
    ('Selected Caller ID Settings', 'Selected Caller ID'),
    ('Call Recording Settings', 'Call Recording'),
    ('Voicemail Settings', 'Voicemail'),
    ('Call Center Settings', 'Call Center'),
    ('Call Forwarding Settings', 'Call Forward'),
    ('Call Park Settings', 'Call Park'),
    ('Call Pickup Group Settings', 'Call Pickup Group'),
    ('Feature Access Codes', 'Access Codes'),
    ('Base Station', 'Base Station'),
    ('DECT Network', 'Network'),
    ('DECT Serviceability Password', 'Service Password'),
    ('Announcement Greeting', 'Greeting'),
    ('Announcement Repository', 'Repository'),
    ('Virtual Line', 'VLine'),
    ('Virtual Extension', 'VExt'),
    ('Fax Message Available', 'Fax Available'),
    ('Call Intercept Available', 'Intercept Available'),
    ('Call Forward Available', 'Forward Available'),
    ('Primary Available', 'Primary'),
    ('Alternate Available', 'Alternate'),
    ('Local Gateway', 'LGW'),
    ('Organization', 'Org'),
    ('RedSky', 'Redsky'),
    ('Building Address and Alert Email', 'Building Address'),
    ('Agent Based Queue', 'Agent Queue'),
    ('Contact Service Queue', 'Queue'),
    ('Service Settings', 'Settings'),
]

# ── Endpoint normalization ────────────────────────────────────────────

VERB_MAP = {
    'get': 'get',
    'read': 'get',
    'retrieve': 'get',
    'fetch': 'get',
    'list': 'list',
    'create': 'create',
    'add': 'create',
    'update': 'update',
    'modify': 'update',
    'configure': 'update',
    'put': 'update',
    'set': 'update',
    'partially': 'patch',
    'delete': 'delete',
    'remove': 'delete',
    'purge': 'purge',
    'bulk': 'bulk',
    'search': 'search',
    'validate': 'validate',
    'upload': 'upload',
    'assign': 'assign',
    'subscribe': 'subscribe',
    'unsubscribe': 'unsubscribe',
    'stream': 'stream',
    'switch': 'switch',
    'historic': 'get-historic',
    'safe': 'check',
    'emergency': None,  # not a verb, keep as noun
    'change': 'update',
    'reload': 'reload',
    'reset': 'reset',
    'enable': 'enable',
    'disable': 'disable',
    'generate': 'generate',
    'accept': 'accept',
    'reject': 'reject',
    'transfer': 'transfer',
    'pause': 'pause',
    'resume': 'resume',
    'stop': 'stop',
    'start': 'start',
    'state': None,
    'answer': 'answer',
    'barge': 'barge',
    'bargein': 'barge-in',
}

FILLER = {
    'a', 'an', 'the', 'of', 'for', 'from', 'to', 'in', 'on', 'at',
    'with', 'by', 'its', 'their', 'this', 'that', 'as', 'or', 'and',
    'specific', 'existing', 'new', 'details', 'certain',
    'given', 'criteria', 'level', 'using',
    'before', 'after', 'only', 'if', 'via', 'into',
    'resource', 'resources',
}

def simple_stem(word):
    """Minimal English stemmer for singular/plural matching."""
    w = word.lower()
    if w.endswith('ies') and len(w) > 4:
        return w[:-3] + 'y'       # entries → entry
    if w.endswith('ses') or w.endswith('xes') or w.endswith('zes'):
        return w[:-2]              # addresses → address
    if w.endswith('s') and not w.endswith('ss') and len(w) > 2:
        return w[:-1]              # users → user, files → file
    return w

# Words from the group name to strip from endpoint names
def group_words(group_name):
    """Get set of words from the group name for dedup."""
    return set(group_name.replace('-', ' ').lower().split())


def normalize_endpoint(name, group=''):
    """Convert Postman endpoint name to a clean kebab-case CLI command name."""
    s = name.strip().rstrip('.')

    # Handle "GET List of..." (uppercase method prefix leaked into name)
    s = re.sub(r'^(GET|POST|PUT|DELETE|PATCH)\s+', '', s)

    # Collapse "Get/Read/Fetch [the/a] List [of]" → "List"
    s = re.sub(r'^(?:Get|Read|Fetch|Retrieve)\s+(?:the\s+|a\s+)?List\s+(?:of\s+)?', 'List ', s, flags=re.IGNORECASE)

    # Strip trailing scope qualifiers (query param context, not resource identity)
    # "for an organization", "for a location", "at organization level", etc.
    s = re.sub(r'\s+(?:for|at|in)\s+(?:an?|the)\s+(?:organization|org)(?:\s+level)?$', '', s, flags=re.IGNORECASE)
    s = re.sub(r'\s+(?:at|on)\s+(?:location\s+and\s+)?(?:organization|org)\s+level$', '', s, flags=re.IGNORECASE)
    s = re.sub(r'\s+(?:at|on)\s+(?:the\s+)?location\s+level$', '', s, flags=re.IGNORECASE)

    # Apply phrase abbreviations
    for long, short in PHRASE_ABBREVIATIONS:
        s = re.sub(re.escape(long), short, s, flags=re.IGNORECASE)

    # Tokenize
    words = re.split(r'[\s/]+', s)
    if not words:
        return name.lower().replace(' ', '-')

    # Extract and normalize verb
    first = words[0].lower()

    # Handle compound verbs
    if first == 'bulk' and len(words) > 1:
        second = words[1].lower()
        verb = f"bulk-{second}"
        rest = words[2:]
    elif first == 'consult' and len(words) > 1:
        second = words[1].lower()
        verb = f"consult-{second}"
        rest = words[2:]
    elif first == 'partially' and len(words) > 1:
        verb = 'patch'
        rest = words[2:]  # skip "Update"
    elif first in VERB_MAP:
        mapped = VERB_MAP[first]
        if mapped is None:
            # Not a verb — treat whole thing as noun
            verb = ''
            rest = words
        elif '-' in mapped:
            # compound like 'get-historic'
            verb = mapped
            rest = words[1:]
        else:
            verb = mapped
            rest = words[1:]
    else:
        verb = ''
        rest = words

    # Filter filler words
    noun_words = [w for w in rest if w.lower() not in FILLER]

    # Strip possessives
    noun_words = [re.sub(r"'s$", '', w) for w in noun_words]

    # Remove parenthetical noise
    cleaned = []
    paren_depth = 0
    for w in noun_words:
        if '(' in w:
            paren_depth += 1
            w = re.sub(r'\(.*', '', w)
            if w:
                cleaned.append(w)
        elif ')' in w:
            paren_depth = max(0, paren_depth - 1)
        elif paren_depth == 0:
            cleaned.append(w)
    noun_words = cleaned

    # Re-apply filler filter after parenthetical stripping (e.g. "resource(s)" → "resource")
    noun_words = [w for w in noun_words if w.lower() not in FILLER]

    # Strip words that duplicate the group name (using stems for plural matching)
    if group:
        gw_stems = {simple_stem(w) for w in group_words(group)}
        # Only strip if we'd still have something left
        filtered = [w for w in noun_words if simple_stem(w) not in gw_stems]
        if filtered or verb:
            noun_words = filtered

    # Build result
    noun = '-'.join(noun_words)
    if verb and noun:
        result = f"{verb}-{noun}"
    elif verb:
        result = verb
    else:
        result = noun

    # Final kebab-case cleanup
    result = re.sub(r'[^a-zA-Z0-9-]', '', result)
    result = re.sub(r'-+', '-', result).strip('-').lower()

    return result


# ── Dedup commands within a group ─────────────────────────────────────

def dedup_commands(endpoints):
    """Add numeric suffix to duplicate command names within a group."""
    seen = {}
    for ep in endpoints:
        cmd = ep['command']
        if cmd in seen:
            seen[cmd] += 1
            ep['command'] = f"{cmd}-{seen[cmd]}"
        else:
            seen[cmd] = 1
    # Go back and suffix the first occurrence if there were dupes
    seen2 = {}
    for ep in endpoints:
        base = re.sub(r'-\d+$', '', ep['command'])
        if base in seen and seen[base] > 1 and base not in seen2:
            seen2[base] = True
            # Don't rename the first one — it keeps the base name
    return endpoints


# ── Main extraction ───────────────────────────────────────────────────

def extract_items(items):
    """Extract folder→endpoint mappings with original and normalized names."""
    folders = []
    for item in items:
        if "item" not in item:
            continue
        raw_folder = item["name"]
        norm_folder = resolve_folder(normalize_folder(raw_folder))

        endpoints = []
        for child in item["item"]:
            if "request" not in child:
                continue
            raw_ep = child["name"]
            norm_ep = normalize_endpoint(raw_ep, norm_folder)
            endpoints.append({
                "original": raw_ep,
                "command": norm_ep,
            })

        dedup_commands(endpoints)

        folders.append({
            "original_folder": raw_folder,
            "group": norm_folder,
            "endpoints": endpoints,
        })
    return folders


def merge_folders(folders):
    """Merge folders that map to the same group name."""
    merged = {}
    order = []
    for f in folders:
        group = f["group"]
        if group not in merged:
            merged[group] = {
                "original_folders": [],
                "group": group,
                "endpoints": [],
            }
            order.append(group)
        merged[group]["original_folders"].append(f["original_folder"])
        merged[group]["endpoints"].extend(f["endpoints"])

    # Dedup again after merging
    for g in order:
        dedup_commands(merged[g]["endpoints"])

    return [merged[g] for g in order]


def main():
    pattern = os.path.join(SCRIPT_DIR, "postman", "*.postman_collection.json")
    files = sorted(glob.glob(pattern))

    if not files:
        print("No Postman collection files found.")
        return

    result = {}
    for filepath in files:
        with open(filepath) as f:
            data = json.load(f)
        name = data["info"]["name"]
        raw_folders = extract_items(data.get("item", []))
        folders = merge_folders(raw_folders)
        result[name] = folders

        total_groups = len(folders)
        total_eps = sum(len(f["endpoints"]) for f in folders)
        print(f"{name}: {total_groups} groups, {total_eps} endpoints")

    output_path = os.path.join(SCRIPT_DIR, "api_outline.json")
    with open(output_path, "w") as f:
        json.dump(result, f, indent=2)
    print(f"\nOutput written to: {output_path}")


if __name__ == "__main__":
    main()
