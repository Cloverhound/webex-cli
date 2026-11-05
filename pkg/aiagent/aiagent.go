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
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// agentList is the v2 agent list, the only place a regular Webex login can read an
// agent's runtime token.
var agentList = webexops.CCAiAgentList

const pageSize = 100

// BotToken returns the runtime token of an AI agent in an org.
func BotToken(ctx context.Context, c *webexapi.Client, hosts webexapi.Hosts, orgUUID, agentID string) (string, error) {
	if orgUUID == "" {
		return "", fmt.Errorf("no organization: log in or pass --organization")
	}
	for offset := 0; ; offset += pageSize {
		u, err := agentList.URL(hosts, webexops.Args{
			Path:  map[string]string{"orgid": orgUUID},
			Query: map[string]string{"limit": strconv.Itoa(pageSize), "offset": strconv.Itoa(offset)},
		})
		if err != nil {
			return "", err
		}
		resp, err := c.Send(ctx, webexapi.Request{Method: http.MethodGet, URL: u})
		if err != nil {
			return "", fmt.Errorf("looking up AI agent %s: %w", agentID, err)
		}
		agents, err := agentPage(resp.Body)
		if err != nil {
			return "", fmt.Errorf("reading the AI agent list: %w", err)
		}
		for _, a := range agents {
			if a.ID != agentID {
				continue
			}
			if a.Token == "" {
				return "", fmt.Errorf("AI agent %s has no runtime token", agentID)
			}
			return a.Token, nil
		}
		if len(agents) < pageSize {
			return "", fmt.Errorf("AI agent %s not found in organization %s", agentID, orgUUID)
		}
	}
}

type agent struct {
	ID    string `json:"id"`
	Token string `json:"bot_access_token"`
}

// agentPage reads one page, which Studio sends under "objects".
func agentPage(body []byte) ([]agent, error) {
	var bare []agent
	if err := json.Unmarshal(body, &bare); err == nil {
		return bare, nil
	}
	var page struct {
		Objects []agent `json:"objects"`
	}
	if err := json.Unmarshal(body, &page); err != nil {
		return nil, err
	}
	return page.Objects, nil
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
