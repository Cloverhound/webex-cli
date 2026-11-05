// Package aiagent supports the Webex AI Agent Studio operations in pkg/webexops
// (those marked Private). Studio's API is not documented or supported by Cisco and
// can change without notice; these routes accept a regular Webex token.
//
// Studio's v1 API selects the agent of a call by its runtime token, sent as the
// bot-access-token header. Callers name the agent by ID, and BotToken looks the
// token up. Responses can carry credentials, which RedactSecrets removes.
package aiagent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// agentList is the v2 agent list, the only place a regular Webex login can read an
// agent's runtime token.
var agentList = webexops.CCAiAgentList

// errFound stops the list once the agent is found.
var errFound = errors.New("found")

// BotToken returns the runtime token of an AI agent in an org.
func BotToken(ctx context.Context, c *webexapi.Client, hosts webexapi.Hosts, orgUUID, agentID string) (string, error) {
	if orgUUID == "" {
		return "", fmt.Errorf("no organization: log in or pass --organization")
	}
	var token string
	var found bool
	err := webexops.List(ctx, c, hosts, agentList, webexops.Args{Path: map[string]string{"orgid": orgUUID}},
		func(page []json.RawMessage) error {
			for _, raw := range page {
				var a struct {
					ID    string `json:"id"`
					Token string `json:"bot_access_token"`
				}
				if json.Unmarshal(raw, &a) == nil && a.ID == agentID {
					token, found = a.Token, true
					return errFound
				}
			}
			return nil
		})
	switch {
	case errors.Is(err, errFound):
	case err != nil:
		return "", fmt.Errorf("looking up AI agent %s: %w", agentID, err)
	case !found:
		return "", fmt.Errorf("AI agent %s not found in organization %s", agentID, orgUUID)
	}
	if token == "" {
		return "", fmt.Errorf("AI agent %s has no runtime token", agentID)
	}
	return token, nil
}

// secretField matches a JSON string value of a Studio credential field.
var secretField = regexp.MustCompile(`("(?:bot_access_token|auth_token|user_access_token|ci_access_token)"\s*:\s*)"(?:[^"\\]|\\.)*"`)

// RedactSecrets replaces Studio credentials in a JSON response, keeping the rest of
// the document byte for byte.
func RedactSecrets(body []byte) []byte {
	return secretField.ReplaceAll(body, []byte(`$1"<redacted>"`))
}

// Call sends a Studio operation. An operation that selects its agent by token gets
// the token of agentID, and a response that can carry credentials comes back with
// them redacted.
func Call(ctx context.Context, c *webexapi.Client, hosts webexapi.Hosts, orgUUID string, op *webexops.Operation, args webexops.Args, agentID string) (*webexapi.Response, error) {
	for _, p := range op.Params {
		if p.In != webexops.InAgentToken {
			continue
		}
		token, err := BotToken(ctx, c, hosts, orgUUID, agentID)
		if err != nil {
			return nil, err
		}
		header := map[string]string{}
		for k, v := range args.Header {
			header[k] = v
		}
		header[p.Name] = token
		args.Header = header
	}
	resp, err := webexops.Call(ctx, c, hosts, op, args)
	if resp != nil && op.Redact {
		resp.Body = RedactSecrets(resp.Body)
	}
	return resp, err
}
