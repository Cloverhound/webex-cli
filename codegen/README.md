# Code generation

Run `make codegen` from the repository root to extract API specifications, generate Go commands, and refresh skill command references. Run `make refresh` to download collections first.

Generated files in the six command area packages must not be edited by hand. Change `generate_cli.py` for generation behavior or use `custom_*.go` with the generator skip lists for handwritten request implementations.

## Naming overrides

Define route-specific group names, command names, and compatibility aliases in `naming_overrides.py`. Overrides are keyed by collection, source group, HTTP method, and route so collection reordering does not affect names. Use meaningful resource/action names rather than numeric collision suffixes. IDs belong in flags; exact `get-id`, `delete-id`, `patch-id`, and `update-id` names are shortened centrally with compatibility aliases and collision checks.

Preserve aliases when a name changes without changing the operation. Do not reuse an alias for a different operation or request format. Keep legacy FDL import/export distinct from current typed flow documents. Duplicate route exclusions must identify a replacement; extraction fails if that replacement disappears.

After changes, run `make codegen`, `make check`, and `python3 -m unittest discover -s codegen -p 'test_*.py'`. Update the [command migration reference](../docs/command-migration.md), [command inventory](../docs/command-inventory.md), and affected documentation/skills when public commands change.
