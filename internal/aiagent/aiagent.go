// Package aiagent supports the Contact Center commands for Webex AI Agent
// Studio, whose private API is described in codegen/supplemental.
package aiagent

import (
	"context"
	"fmt"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/aiagent"
)

// studioURL is read at call time, after the command has chosen the org's data
// center, and is a variable so tests can point it at a fake Studio.
var studioURL = func() string { return config.AiAgentBaseURL }

// BotToken returns the runtime token of an AI agent in the current
// organization. Studio's v1 API selects the agent of a call by this token, and
// the v2 agent list is the only place a regular Webex login can read it.
func BotToken(agentID string) (string, error) {
	orgID := config.OrgID()
	if orgID == "" {
		return "", fmt.Errorf("no organization: log in or pass --organization")
	}
	return aiagent.BotToken(context.Background(), client.API("/api/v2/organizations/{orgid}/agents/", ""), studioURL(), orgID, agentID)
}

// RedactSecrets replaces Studio credentials in a JSON response; see
// aiagent.RedactSecrets.
func RedactSecrets(body []byte) []byte { return aiagent.RedactSecrets(body) }
