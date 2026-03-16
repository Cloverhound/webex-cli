#!/usr/bin/env python3
"""Download Postman collections from the Webex Public Workspace.

Uses Postman's public bifrost gateway — no API key or authentication required.
Downloads collections in Postman's internal format and converts them to the
standard Collection v2.1 format expected by extract_api_spec.py.
"""

import argparse
import json
import os
import ssl
import sys
import time
import urllib.error
import urllib.request

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
OUTPUT_DIR = os.path.join(SCRIPT_DIR, "postman")

# Owner ID for the webexdev team on Postman
OWNER_ID = "46560387"

# Collection _postman_id values from the Webex Public Workspace
COLLECTIONS = {
    "Webex Admin": "f825adb1-7a62-434e-a9c3-509952cd3b73",
    "Webex Cloud Calling": "07753680-e684-4ec5-b4e0-58972cd75731",
    "Webex Contact Center": "2d3506aa-a473-41cf-a638-c9dd22fb4948",
    "Webex Device": "9765729d-69c1-423b-a937-5f117bceae9e",
    "Webex Meetings": "9a95c130-f005-4963-ad39-01f65fd14e18",
    "Webex Messaging": "f2747234-f76f-4227-8943-99bfc298c90d",
}

BIFROST_BASE = "https://bifrost-public-https-v4.gw.postman.com/collection"
MAX_RETRIES = 3
RETRY_BACKOFF = 2  # seconds, doubles each retry
REQUEST_DELAY = 1.0  # seconds between requests


def _make_ssl_context():
    """Create an SSL context, falling back to system cert bundle on macOS."""
    ctx = ssl.create_default_context()
    if not ctx.get_ca_certs():
        system_certs = "/etc/ssl/cert.pem"
        if os.path.exists(system_certs):
            ctx.load_verify_locations(system_certs)
    return ctx


_SSL_CONTEXT = _make_ssl_context()


def download_collection(uuid):
    """Download a single collection from bifrost, retrying on transient errors."""
    uid = f"{OWNER_ID}-{uuid}"
    url = f"{BIFROST_BASE}/{uid}?populate=true"
    req = urllib.request.Request(url, headers={"x-entity-team-id": "0"})

    for attempt in range(MAX_RETRIES):
        try:
            with urllib.request.urlopen(req, timeout=120, context=_SSL_CONTEXT) as resp:
                return json.loads(resp.read().decode("utf-8"))
        except urllib.error.HTTPError as e:
            if e.code in (429, 500, 502, 503, 504) and attempt < MAX_RETRIES - 1:
                wait = RETRY_BACKOFF * (2 ** attempt)
                print(f"  HTTP {e.code}, retrying in {wait}s...")
                time.sleep(wait)
                continue
            raise
        except urllib.error.URLError as e:
            if attempt < MAX_RETRIES - 1:
                wait = RETRY_BACKOFF * (2 ** attempt)
                print(f"  Network error ({e.reason}), retrying in {wait}s...")
                time.sleep(wait)
                continue
            raise


# ---------------------------------------------------------------------------
# Bifrost internal format → Postman Collection v2.1 conversion
# ---------------------------------------------------------------------------

def _convert_url(url_str, query_params=None, path_variable_data=None):
    """Convert a URL string + params to v2.1 url object."""
    url_obj = {"raw": url_str}

    # Parse host and path from the raw URL (strip query string for path parsing)
    # e.g. "{{baseUrl}}/telephony/config/locations/:locationId?orgId=<string>"
    base_url = url_str.split("?")[0]
    parts = base_url.split("/")
    if parts:
        url_obj["host"] = [parts[0]]
        url_obj["path"] = [p for p in parts[1:] if p]

    if query_params:
        url_obj["query"] = [
            {k: p[k] for k in ("key", "value", "description", "disabled") if k in p}
            for p in query_params
        ]

    if path_variable_data:
        url_obj["variable"] = [
            {k: p[k] for k in ("key", "value", "description") if k in p}
            for p in path_variable_data
        ]

    return url_obj


def _convert_headers(header_data):
    """Convert headerData array to v2.1 header array."""
    if not header_data:
        return []
    return [
        {k: h[k] for k in ("key", "value", "description", "disabled") if k in h}
        for h in header_data
    ]


def _convert_body(request):
    """Convert bifrost body fields to v2.1 body object."""
    mode = request.get("dataMode")
    if not mode or mode == "null":
        return None

    body = {"mode": mode}
    if mode == "raw":
        body["raw"] = request.get("rawModeData") or ""
        # Add options if content type suggests JSON
        headers_str = request.get("headers", "")
        if "application/json" in headers_str:
            body["options"] = {"raw": {"headerFamily": "json", "language": "json"}}
    elif mode == "urlencoded":
        body["urlencoded"] = request.get("data") or []
    elif mode == "formdata":
        body["formdata"] = request.get("data") or []

    return body


def _convert_response(resp):
    """Convert a bifrost response to v2.1 format."""
    result = {
        "name": resp.get("name", ""),
        "status": resp.get("status") or (resp.get("responseCode") or {}).get("name", ""),
        "code": (resp.get("responseCode") or {}).get("code"),
        "_postman_previewlanguage": resp.get("language", ""),
        "header": resp.get("headers") if isinstance(resp.get("headers"), list) else [],
        "cookie": [],
        "body": resp.get("text", ""),
    }

    # Convert requestObject (stored as JSON string in bifrost)
    req_obj = resp.get("requestObject")
    if req_obj:
        if isinstance(req_obj, str):
            try:
                req_obj = json.loads(req_obj)
            except json.JSONDecodeError:
                req_obj = None
        if req_obj:
            result["originalRequest"] = {
                "method": req_obj.get("method", ""),
                "header": _convert_headers(req_obj.get("headerData")),
                "url": _convert_url(
                    req_obj.get("url", ""),
                    req_obj.get("queryParams"),
                    req_obj.get("pathVariableData"),
                ),
            }
            body = _convert_body(req_obj)
            if body:
                result["originalRequest"]["body"] = body

    return result


def _convert_request(req):
    """Convert a single bifrost request to a v2.1 item."""
    v21_request = {
        "method": req.get("method", "GET"),
        "header": _convert_headers(req.get("headerData")),
        "url": _convert_url(
            req.get("url", ""),
            req.get("queryParams"),
            req.get("pathVariableData"),
        ),
    }

    if req.get("description"):
        v21_request["description"] = req["description"]

    body = _convert_body(req)
    if body:
        v21_request["body"] = body

    if req.get("auth"):
        v21_request["auth"] = req["auth"]

    item = {
        "name": req.get("name", ""),
        "request": v21_request,
        "response": [_convert_response(r) for r in (req.get("responses") or [])],
    }

    return item


def bifrost_to_v21(data):
    """Convert bifrost internal format to Postman Collection v2.1."""
    # Build lookup maps
    folders_by_id = {f["id"]: f for f in data.get("folders", [])}
    requests_by_id = {r["id"]: r for r in data.get("requests", [])}

    # Group requests by folder
    folder_requests = {}
    for req in data.get("requests", []):
        fid = req.get("folder")
        folder_requests.setdefault(fid, []).append(req)

    # Group subfolders by parent folder
    folder_children = {}
    for folder in data.get("folders", []):
        parent = folder.get("folder")
        folder_children.setdefault(parent, []).append(folder)

    def build_folder_item(folder):
        """Recursively build a v2.1 folder item."""
        item = {"name": folder["name"], "item": []}

        if folder.get("description"):
            item["description"] = folder["description"]

        # Add sub-folders in order
        sub_folder_ids = folder.get("folders_order", [])
        sub_folders = folder_children.get(folder["id"], [])
        sub_folder_map = {f["id"]: f for f in sub_folders}
        for fid in sub_folder_ids:
            if fid in sub_folder_map:
                item["item"].append(build_folder_item(sub_folder_map[fid]))
        # Add any remaining sub-folders not in folders_order
        for f in sub_folders:
            if f["id"] not in sub_folder_ids:
                item["item"].append(build_folder_item(f))

        # Add requests in order
        req_ids = folder.get("order", [])
        reqs = folder_requests.get(folder["id"], [])
        req_map = {r["id"]: r for r in reqs}
        for rid in req_ids:
            if rid in req_map:
                item["item"].append(_convert_request(req_map[rid]))
        # Add any remaining requests not in order
        for r in reqs:
            if r["id"] not in req_ids:
                item["item"].append(_convert_request(r))

        return item

    # Build top-level items from folders_order
    items = []
    top_folder_ids = data.get("folders_order", [])
    top_folders = folder_children.get(None, [])
    top_folder_map = {f["id"]: f for f in top_folders}
    for fid in top_folder_ids:
        if fid in top_folder_map:
            items.append(build_folder_item(top_folder_map[fid]))
    for f in top_folders:
        if f["id"] not in top_folder_ids:
            items.append(build_folder_item(f))

    # Add top-level requests (not in any folder)
    top_req_ids = data.get("order", [])
    top_reqs = folder_requests.get(None, [])
    top_req_map = {r["id"]: r for r in top_reqs}
    for rid in top_req_ids:
        if rid in top_req_map:
            items.append(_convert_request(top_req_map[rid]))
    for r in top_reqs:
        if r["id"] not in top_req_ids:
            items.append(_convert_request(r))

    collection = {
        "info": {
            "_postman_id": data.get("id", "").replace(f"{OWNER_ID}-", ""),
            "name": data.get("name", ""),
            "description": data.get("description", ""),
            "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
        },
        "item": items,
    }

    if data.get("variables"):
        collection["variable"] = data["variables"]

    if data.get("auth"):
        collection["auth"] = data["auth"]

    return collection


def save_collection(name, data):
    """Save collection JSON and report whether it changed."""
    filename = f"{name}.postman_collection.json"
    filepath = os.path.join(OUTPUT_DIR, filename)
    new_content = json.dumps(data, indent=2, sort_keys=True) + "\n"

    if os.path.exists(filepath):
        with open(filepath, "r") as f:
            old_content = f.read()
        if old_content == new_content:
            return False

    with open(filepath, "w") as f:
        f.write(new_content)
    return True


def main():
    parser = argparse.ArgumentParser(
        description="Download Postman collections from the Webex Public Workspace"
    )
    parser.add_argument(
        "--collection",
        help="Download only a specific collection by name (e.g. 'Webex Admin')",
    )
    args = parser.parse_args()

    os.makedirs(OUTPUT_DIR, exist_ok=True)

    if args.collection:
        if args.collection not in COLLECTIONS:
            print(
                f"Error: Unknown collection '{args.collection}'.\n"
                f"Available: {', '.join(sorted(COLLECTIONS))}",
                file=sys.stderr,
            )
            sys.exit(1)
        targets = {args.collection: COLLECTIONS[args.collection]}
    else:
        targets = COLLECTIONS

    errors = []
    changed = []
    unchanged = []

    for i, (name, uuid) in enumerate(targets.items()):
        if i > 0:
            time.sleep(REQUEST_DELAY)

        print(f"Downloading {name}...")
        try:
            response = download_collection(uuid)
        except Exception as e:
            print(f"  ERROR: {e}")
            errors.append((name, str(e)))
            continue

        collection_data = bifrost_to_v21(response["data"])

        if save_collection(name, collection_data):
            changed.append(name)
            print(f"  Updated")
        else:
            unchanged.append(name)
            print(f"  Unchanged")

    # Summary
    print()
    if changed:
        print(f"Changed ({len(changed)}): {', '.join(changed)}")
    if unchanged:
        print(f"Unchanged ({len(unchanged)}): {', '.join(unchanged)}")
    if errors:
        print(f"Errors ({len(errors)}):")
        for name, err in errors:
            print(f"  {name}: {err}")
        sys.exit(1)


if __name__ == "__main__":
    main()
