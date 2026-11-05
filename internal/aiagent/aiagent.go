// Package aiagent supports the Contact Center commands for Webex AI Agent
// Studio, whose private API is described in codegen/supplemental.
package aiagent

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
)

// baseURL is a variable so tests can point it at a fake Studio.
var baseURL = config.AiAgentBaseURL

// BotToken returns the runtime token of an AI agent in the current
// organization. Studio's v1 API selects the agent of a call by this token, and
// the v2 agent list is the only place a regular Webex login can read it.
func BotToken(agentID string) (string, error) {
	orgID := config.OrgID()
	if orgID == "" {
		return "", fmt.Errorf("no organization: log in or pass --organization")
	}
	req := client.NewRequest(baseURL(), "GET", "/api/v2/organizations/{orgid}/agents/")
	req.PathParam("orgid", orgID)
	req.OffsetPaging("offset", "limit", 0)
	body, _, err := req.DoPaginated(false)
	if err != nil {
		return "", fmt.Errorf("looking up AI agent %s: %w", agentID, err)
	}
	var agents []struct {
		ID    string `json:"id"`
		Token string `json:"bot_access_token"`
	}
	if err := json.Unmarshal(body, &agents); err != nil {
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
	return "", fmt.Errorf("AI agent %s not found in organization %s", agentID, orgID)
}

// secretField matches a JSON string value of a Studio credential field.
var secretField = regexp.MustCompile(`("(?:bot_access_token|auth_token|user_access_token|ci_access_token)"\s*:\s*)"(?:[^"\\]|\\.)*"`)

// RedactSecrets replaces Studio credentials in a JSON response, keeping the
// rest of the document byte for byte.
func RedactSecrets(body []byte) []byte {
	return secretField.ReplaceAll(body, []byte(`$1"<redacted>"`))
}
