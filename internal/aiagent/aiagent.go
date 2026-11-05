// Package aiagent connects the CLI's AI Agent Studio commands to pkg/aiagent, with
// this run's login, org and hosts, and shows the private-API notice.
package aiagent

import (
	"context"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	sdk "github.com/Cloverhound/webex-cli/pkg/aiagent"
)

// BotToken returns the runtime token of an AI agent in the current organization.
func BotToken(agentID string) (string, error) {
	return sdk.BotToken(context.Background(), client.SDK(), config.Hosts(), config.OrgID(), agentID)
}

// RedactSecrets replaces Studio credentials in a JSON response.
func RedactSecrets(body []byte) []byte { return sdk.RedactSecrets(body) }
