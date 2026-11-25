package cc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/aiagent"
	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	studio "github.com/Cloverhound/webex-cli/pkg/aiagent"
	"github.com/spf13/cobra"
)

const previewHelp = `Talks to an AI agent through AI Agent Studio's text preview, as the Preview
panel in Studio does. It uses the agent's draft and needs only the regular
Webex login.

The conversation is saved per user, organization and agent, so 'send' continues
it. After the first use, --agent-id defaults to the last agent used.

Preview can run the agent's actions, and their fulfillment can be real: a
confirmed order or transfer in preview may reach real systems. The output is
the runtime's full response: actions are in messageStore.tools, and a handover
sets sendtoagent. --brief shows only the text sent and received, so it hides
both. Custom-event actions fail in text preview (errorCode 67405002), since
Studio runs them only on voice.

` + aiagent.HelpCaution

var aiAgentPreviewCmd = &cobra.Command{
	Use:   "ai-agent-preview",
	Short: "AI Agent Studio text preview conversations (unsupported private API)",
	Long:  previewHelp,
}

func init() {
	cmd.CcCmd.AddCommand(aiAgentPreviewCmd)
	registerPreviewStart()
	registerPreviewSend()
	registerPreviewRun()
}

// previewTemplate is the runtime path, which read-only mode refuses.
const previewTemplate = "/api/v1/webhook/web/"

// newPreview looks up the agent's runtime token and returns a preview of it.
// agentID empty means the user's last preview agent.
func newPreview(agentID string) (*studio.Preview, string, error) {
	aiagent.PrivateAPINotice()
	if agentID == "" {
		last, err := aiagent.LastPreviewAgent(config.User(), config.OrgID())
		if err != nil {
			return nil, "", err
		}
		agentID = last
	}
	runtime, origin, err := studio.PreviewHosts(config.AiAgentBaseURL)
	if err != nil {
		return nil, "", err
	}
	token, err := aiagent.BotToken(agentID)
	if err != nil {
		return nil, "", err
	}
	return &studio.Preview{Runtime: runtime, Origin: origin, BotToken: token}, agentID, nil
}

// previewTurn sends one turn. The client is made per turn so --debug and
// --dry-run show its body.
func previewTurn(p *studio.Preview, c *studio.Conversation, kind, message string) (studio.Turn, error) {
	body, err := c.Body(kind, message)
	if err != nil {
		return studio.Turn{}, err
	}
	api := client.API(previewTemplate, string(body))
	// The runtime takes the agent's token; the Webex login is not sent to it.
	api.Token, api.Refresh = nil, nil
	p.Client = api
	var turn studio.Turn
	if kind == "bot" {
		turn, err = p.Start(context.Background(), c)
	} else {
		turn, err = p.Send(context.Background(), c, message)
	}
	return turn, client.CLIError(err)
}

func printPreview(c studio.Conversation, turns []studio.Turn, brief bool) error {
	var v any = map[string]any{"conversation": c, "turns": turns}
	if brief {
		b, err := aiagent.Brief(turns)
		if err != nil {
			return err
		}
		v = b
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return output.Print(aiagent.RedactSecrets(data), http.StatusOK)
}

func registerPreviewStart() {
	var agentID string
	var brief bool
	c := &cobra.Command{
		Use:   "start",
		Short: "Start a conversation and show the agent's first message",
		Long: previewHelp + `

Examples:
  webex cc ai-agent-preview start --agent-id <agent-id>
  webex cc ai-agent-preview send "Where is my order?"`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			p, agent, err := newPreview(agentID)
			if err != nil {
				return err
			}
			conv, err := studio.NewConversation(agent)
			if err != nil {
				return err
			}
			turn, err := previewTurn(p, &conv, "bot", "hi")
			if err != nil {
				return err
			}
			if err := aiagent.SavePreview(config.User(), config.OrgID(), conv); err != nil {
				return err
			}
			return printPreview(conv, []studio.Turn{turn}, brief)
		},
	}
	c.Flags().StringVar(&agentID, "agent-id", "", "AI agent ID (default: the last agent used)")
	c.Flags().BoolVar(&brief, "brief", false, "Show only the text sent and received")
	aiAgentPreviewCmd.AddCommand(c)
}

func registerPreviewSend() {
	var agentID string
	var brief bool
	c := &cobra.Command{
		Use:   "send <message>",
		Short: "Send a message in the saved conversation",
		Long: previewHelp + `

Examples:
  webex cc ai-agent-preview send "Where is my order?"
  webex cc ai-agent-preview send --agent-id <agent-id> "Talk to a person" --brief`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			conv, err := aiagent.SavedPreview(config.User(), config.OrgID(), agentID)
			if err != nil {
				return err
			}
			p, _, err := newPreview(conv.AgentID)
			if err != nil {
				return err
			}
			turn, err := previewTurn(p, &conv, "human", args[0])
			if err != nil {
				return err
			}
			if err := aiagent.SavePreview(config.User(), config.OrgID(), conv); err != nil {
				return err
			}
			return printPreview(conv, []studio.Turn{turn}, brief)
		},
	}
	c.Flags().StringVar(&agentID, "agent-id", "", "AI agent ID (default: the last agent used)")
	c.Flags().BoolVar(&brief, "brief", false, "Show only the text sent and received")
	aiAgentPreviewCmd.AddCommand(c)
}

func registerPreviewRun() {
	var agentID string
	var messages []string
	var brief bool
	c := &cobra.Command{
		Use:   "run",
		Short: "Start a conversation and send each --message in turn",
		Long: previewHelp + `

Examples:
  webex cc ai-agent-preview run --agent-id <agent-id> --message "Hello" --message "Where is my order?" --brief`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error {
			p, agent, err := newPreview(agentID)
			if err != nil {
				return err
			}
			conv, err := studio.NewConversation(agent)
			if err != nil {
				return err
			}
			// Under --dry-run every turn is shown, and none is sent.
			var dryRun error
			var turns []studio.Turn
			for i, m := range append([]string{"hi"}, messages...) {
				kind := "human"
				if i == 0 {
					kind = "bot"
				}
				turn, err := previewTurn(p, &conv, kind, m)
				if errors.Is(err, client.ErrDryRun) {
					dryRun = err
					continue
				}
				if err != nil {
					return err
				}
				turns = append(turns, turn)
				if err := aiagent.SavePreview(config.User(), config.OrgID(), conv); err != nil {
					return err
				}
			}
			if dryRun != nil {
				return dryRun
			}
			return printPreview(conv, turns, brief)
		},
	}
	c.Flags().StringVar(&agentID, "agent-id", "", "AI agent ID (default: the last agent used)")
	c.Flags().StringArrayVar(&messages, "message", nil, "Message to send; repeat for each turn")
	c.Flags().BoolVar(&brief, "brief", false, "Show only the text sent and received")
	_ = c.MarkFlagRequired("message")
	aiAgentPreviewCmd.AddCommand(c)
}
