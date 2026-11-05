# Code generation

Run `make codegen` from the repository root to extract API specifications, generate the operation table, and refresh skill command references. Run `make refresh` to download collections first.

`generate_cli.py` writes `pkg/webexops/zz_<area>.go`: one exported `*Operation` per API call, holding its route, parameters, body fields, paging and CLI names. The CLI builds its commands from that table at startup (`internal/apicmd`), and other Go programs call the same operations. The generated files must not be edited by hand. Change `generate_cli.py` for what the table holds, `internal/apicmd` for how operations become commands, or use `custom_*.go` with the generator skip lists for handwritten request implementations.

`testdata/golden/` at the repository root records every command's flags and help, and the request each API command sends. Run `make golden` after an intended change and review the diff.

## Supplemental collections

APIs that Cisco has not published in Postman are described by hand in [`supplemental/`](supplemental/README.md). Each one joins the downloaded collection its `info.name` names and is regenerated with it, while `make download` leaves the folder alone. A route that a downloaded collection also documents is taken from the download.

## Naming overrides

Define route-specific group names, command names, and compatibility aliases in `naming_overrides.py`. Overrides are keyed by collection, source group, HTTP method, and route so collection reordering does not affect names. Use meaningful resource/action names rather than numeric collision suffixes. IDs belong in flags; exact `get-id`, `delete-id`, `patch-id`, and `update-id` names are shortened centrally with compatibility aliases and collision checks.

Preserve aliases when a name changes without changing the operation. Do not reuse an alias for a different operation or request format. Keep legacy FDL import/export distinct from current typed flow documents. Duplicate route exclusions must identify a replacement; extraction fails if that replacement disappears.

After changes, run `make codegen`, `make check`, and `python3 -m unittest discover -s codegen -p 'test_*.py'`. Update the [command migration reference](../docs/command-migration.md), [command inventory](../docs/command-inventory.md), and affected documentation/skills when public commands change.
