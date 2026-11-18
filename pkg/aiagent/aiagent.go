// Package aiagent supports Webex AI Agent Studio, whose private API Cisco does not
// document or support. Its v1 calls select an agent by the agent's runtime token,
// which BotToken looks up, and its responses carry credentials, which
// RedactSecrets removes.
package aiagent

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// agentsPath lists an organization's agents, with their runtime tokens.
const agentsPath = "/api/v2/organizations/{orgid}/agents/"

// agentsPaging reads the agent list to its end.
var agentsPaging = webexapi.Paging{Mode: webexapi.Offset, Start: "offset", Size: "limit", Items: "objects", TotalItems: "meta.total_count"}

// BotToken returns the runtime token of an AI agent. studio is AI Agent Studio's
// base URL in the organization's Contact Center data center (webexapi.Hosts.AIAgent).
// The v2 agent list is the only place a regular Webex login can read the token.
func BotToken(ctx context.Context, c *webexapi.Client, studio, orgUUID, agentID string) (string, error) {
	if orgUUID == "" {
		return "", fmt.Errorf("no organization to look up AI agent %s in", agentID)
	}
	list := webexapi.List{
		Method: "GET",
		URL:    webexapi.FillPath(studio+agentsPath, map[string]string{"orgid": orgUUID}),
		Redact: RedactSecrets,
		Name:   agentsPath,
	}
	items, err := c.All(ctx, list, agentsPaging)
	if err != nil {
		return "", fmt.Errorf("looking up AI agent %s: %w", agentID, err)
	}
	for _, item := range items {
		var a struct {
			ID    string `json:"id"`
			Token string `json:"bot_access_token"`
		}
		if err := json.Unmarshal(item, &a); err != nil {
			return "", fmt.Errorf("reading the AI agent list: %w", err)
		}
		if a.ID != agentID {
			continue
		}
		if a.Token == "" {
			return "", fmt.Errorf("AI agent %s has no runtime token", agentID)
		}
		return a.Token, nil
	}
	return "", fmt.Errorf("AI agent %s not found in organization %s", agentID, orgUUID)
}

// secretField matches a JSON string value of a Studio credential field.
var secretField = regexp.MustCompile(`("(?:bot_access_token|auth_token|user_access_token|ci_access_token)"\s*:\s*)"(?:[^"\\]|\\.)*"`)

// RedactSecrets replaces Studio credentials in a JSON response, keeping the
// rest of the document byte for byte.
func RedactSecrets(body []byte) []byte {
	return secretField.ReplaceAll(body, []byte(`$1"<redacted>"`))
}
