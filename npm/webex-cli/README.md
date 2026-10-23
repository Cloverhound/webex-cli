# @cloverhound/webex-cli

The [Webex CLI](https://github.com/Cloverhound/webex-cli) packaged for npm. It manages Webex Admin, Calling, Contact Center, Devices, Meetings, and Messaging APIs from the command line, and runs as an MCP server for AI agents.

```bash
npx -y @cloverhound/webex-cli login --device
npx -y @cloverhound/webex-cli calling people list
```

Or install it globally:

```bash
npm install -g @cloverhound/webex-cli
webex --version
```

npm installs only the prebuilt binary for your platform (macOS, Linux, or Windows on x64 or arm64) through an optional dependency. No install script runs and nothing is downloaded from GitHub, so the package works in sandboxes that allow the npm registry but block other hosts.

To run it as an MCP server, use `npx -y @cloverhound/webex-cli@<version> mcp serve` as the server command.

See the [documentation](https://cloverhound.github.io/webex-cli/) for login options, including device login for machines without a browser.
