# Webex CLI

A command-line tool for Webex Calling and Contact Center APIs.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/Cloverhound/webex-cli/main/install.sh | sh
```

Or download from [Releases](https://github.com/Cloverhound/webex-cli/releases).

## Quick Start

```bash
# Login (opens browser for OAuth)
webex login

# Webex Calling
webex calling people list --max 10
webex calling locations list

# Contact Center
webex cc site list
webex cc team list
webex cc users list
webex cc entry-point list
```

## Usage

```
webex calling <resource> <action> [flags]     # Webex Calling APIs
webex cc <resource> <action> [flags]          # Contact Center APIs
```

### Global Flags

| Flag | Description |
|------|-------------|
| `--token <token>` | Override authentication |
| `--user <email>` | Use a specific authenticated user |
| `--organization <orgId>` | Override org ID |
| `--output json\|table\|raw` | Output format (default: json) |
| `--debug` | Show HTTP request/response details |
| `--paginate` | Auto-paginate list results |

## Documentation

See the [docs site](https://cloverhound.github.io/webex-cli/) for full documentation.

## Claude Code Integration

A sample Claude Code skill is included in `skill/SKILL.md`. See the [docs](https://cloverhound.github.io/webex-cli/claude-skill/) for setup instructions.

## Development

See [CLAUDE.md](CLAUDE.md) for project structure and development workflow.

## License

MIT
