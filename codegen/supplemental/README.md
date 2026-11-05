# Supplemental collections

Hand-maintained Postman collections for APIs that Cisco has not published in
its Postman workspace. `download_collections.py` never writes here, so a
refresh does not remove them.

Each file's `info.name` names the downloaded collection it joins, such as
`Webex Contact Center`. Its folders become groups of that collection's area,
and `make codegen` generates them like any other group. Extraction fails if no
downloaded collection has that name.

When a downloaded collection documents the same method and route, the
downloaded definition wins and extraction prints the route to remove from
here. Routes here keep their trailing slash, which the APIs below require.

## AI Agent Studio

`AI Agent Studio.postman_collection.json` holds the private AI Agent Studio
routes that accept a regular Webex token. They generate the `cc ai-agent`,
`cc ai-agent-session` and `cc ai-agent-transcript` groups. Each route was
exercised against a live organization before it was added. Add only routes
that work with a Webex OAuth token; routes that need Studio's own browser login
belong elsewhere.

- `generate_cli.py` sends routes under `/api/` to `config.AiAgentBaseURL` and
  passes their responses through `aiagent.RedactSecrets`.
- A request header named `bot-access-token` becomes a required `--agent-id`
  flag; the command looks up the agent's runtime token with `aiagent.BotToken`.
- `limit`/`offset` query parameters make `--paginate` page by offset.
