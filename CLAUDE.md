# Webex CLI

A CLI tool for Webex Calling and Contact Center APIs.

## Project Layout

- Go CLI at repo root (`main.go`, `cmd/`, `internal/`)
- `codegen/` — Python scripts that generate `cmd/calling/*.go` and `cmd/cc/*.go` from Postman collections
- `site/` — Astro Starlight docs site for GitHub Pages
- `skill/` — Sample Claude Code skill for using the CLI

## Code Generation Pipeline

The files in `cmd/calling/` and `cmd/cc/` are **generated** — do not edit them by hand.

To regenerate after modifying codegen scripts or Postman collections:

```bash
cd codegen
python3 extract_api_spec.py    # reads postman/*.json → api_spec.json
python3 generate_cli.py        # reads api_spec.json → writes cmd/calling/*.go + cmd/cc/*.go
cd ..
go build -o webex .
```

## Build

```bash
go build -o webex .
```

For release builds with embedded credentials (via goreleaser):
```bash
goreleaser release --snapshot --clean
```

## Module Path

`github.com/Cloverhound/webex-cli` — all Go imports use this path.

## Key Files

- `cmd/root.go` — Root cobra command, global flags
- `cmd/auth.go` — Auth subcommand (status, list, switch)
- `cmd/login.go` / `cmd/logout.go` — OAuth flow
- `internal/auth/` — Token storage (keyring), OAuth, org resolution
- `internal/client/` — HTTP client, request builder, pagination
- `internal/config/` — Config file (~/.webex-cli.yaml)
- `internal/output/` — JSON/table/raw output formatting
- `codegen/extract_api_outline.py` — Postman → normalized CLI names
- `codegen/extract_api_spec.py` — Postman → enriched API spec (params, body fields)
- `codegen/generate_cli.py` — API spec → Go cobra command files
