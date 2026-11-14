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

## Contact Center lists

`Contact Center lists.postman_collection.json` holds older Contact Center list
routes that Webex still serves but Cisco's collection replaced with newer
versions. They join their resource's group as `list-v1` (`list-v2` for queues;
see `naming_overrides.py`). The SDK catalog (`pkg/catalog`) reads through them,
because their items carry `lastUpdatedTime` and, for most types, every field.
Each was read on a live organization before it was added.

## AI Agent Studio

`AI Agent Studio.postman_collection.json` holds the private AI Agent Studio
routes that accept a regular Webex token. They generate the `cc ai-agent`,
`cc ai-agent-session` and `cc ai-agent-transcript` groups. Each route was
exercised against a live organization before it was added. Add only routes
that work with a Webex OAuth token; routes that need Studio's own browser login
belong elsewhere.

- Cisco does not document or support these routes. `generate_cli.py` lists
  them in `PRIVATE_API_ROUTES`, which marks their operations `Private`: every
  command carries that caution in `--help` and prints it once per user when
  first run. A route added here gets both automatically.
- `generate_cli.py` puts routes under `/api/` on the AI Agent Studio host of
  the selected Contact Center data center and marks them `Redact`, so responses
  pass through `pkg/aiagent.RedactSecrets`.
- A request header named `bot-access-token` becomes a required `--agent-id`
  flag; the agent's runtime token is looked up with `pkg/aiagent.BotToken`.
- `limit`/`offset` query parameters make `--paginate` page by offset.
