package aiagent

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// previewPath is the text preview runtime's endpoint for web chat.
const previewPath = "/api/v1/webhook/web/"

var studioHost = regexp.MustCompile(`(?i)^studio\.aiagent-([a-z0-9-]+)\.cisco\.com$`)

// PreviewHosts returns the text preview runtime for AI Agent Studio's base URL,
// and the Origin the runtime expects, which is Studio's preview page.
func PreviewHosts(studio string) (runtime, origin string, err error) {
	u, err := url.Parse(studio)
	if err != nil {
		return "", "", err
	}
	m := studioHost.FindStringSubmatch(u.Hostname())
	if m == nil {
		return "", "", fmt.Errorf("no text preview runtime is known for %s", u.Hostname())
	}
	dc := "prod" + strings.ToLower(m[1])
	return "https://bot-appserver." + dc + ".ciscoccservice.com", "https://preview." + dc + ".ciscoccservice.com", nil
}

// Conversation identifies a text preview conversation. Every turn sends the same
// consumer UID and socket key; a new UID starts a new Studio session.
type Conversation struct {
	AgentID     string `json:"agent_id"`
	ConsumerUID string `json:"consumer_uid"`
	SocketKey   string `json:"socket_key"`
	RoomID      string `json:"room_id,omitempty"`
}

// NewConversation returns a conversation with a new consumer UID and socket key.
func NewConversation(agentID string) (Conversation, error) {
	uid := make([]byte, 16)
	if _, err := rand.Read(uid); err != nil {
		return Conversation{}, err
	}
	uid[6] = uid[6]&0x0f | 0x40
	uid[8] = uid[8]&0x3f | 0x80
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	key := make([]byte, 15)
	for i := range key {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return Conversation{}, err
		}
		key[i] = alphabet[n.Int64()]
	}
	return Conversation{
		AgentID:     agentID,
		ConsumerUID: fmt.Sprintf("%x-%x-%x-%x-%x", uid[0:4], uid[4:6], uid[6:8], uid[8:10], uid[10:16]),
		SocketKey:   string(key),
	}, nil
}

// Turn is one request of a conversation and the runtime's response.
type Turn struct {
	// Type is "bot" for the request that opens a conversation and "human" for a
	// message.
	Type     string          `json:"type"`
	Message  string          `json:"message"`
	Response json.RawMessage `json:"response"`
}

// Preview sends text preview requests to one agent. The runtime authenticates the
// agent's runtime token (BotToken), not a Webex login, so Client should have no
// Token.
type Preview struct {
	Client          *webexapi.Client
	Runtime, Origin string
	BotToken        string
}

// Start opens the conversation, as Studio's preview page does, and returns the
// agent's first message. Its "hi" is a placeholder the agent does not see.
func (p *Preview) Start(ctx context.Context, c *Conversation) (Turn, error) {
	return p.turn(ctx, c, "bot", "hi")
}

// Send sends a message in the conversation.
func (p *Preview) Send(ctx context.Context, c *Conversation, message string) (Turn, error) {
	if strings.TrimSpace(message) == "" {
		return Turn{}, fmt.Errorf("the message is empty")
	}
	return p.turn(ctx, c, "human", message)
}

// Body is the request body of a turn: kind "bot" opens the conversation, and
// "human" sends message.
func (c Conversation) Body(kind, message string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"consumer":                map[string]any{"uid": c.ConsumerUID, "extra_params": map[string]string{"socket_key": c.SocketKey}},
		"type":                    kind,
		"msg":                     message,
		"platform":                "web",
		"transparency_msg_merged": false,
		"is_test":                 true,
	})
}

func (p *Preview) turn(ctx context.Context, c *Conversation, kind, message string) (Turn, error) {
	body, err := c.Body(kind, message)
	if err != nil {
		return Turn{}, err
	}
	header := http.Header{"Bot-Access-Token": {p.BotToken}}
	if p.Origin != "" {
		header.Set("Origin", p.Origin)
	}
	resp, err := p.Client.Send(ctx, webexapi.Request{
		Method: http.MethodPost,
		URL:    strings.TrimSuffix(p.Runtime, "/") + previewPath,
		Body:   body,
		Header: header,
		Redact: RedactSecrets,
	})
	if err != nil {
		return Turn{}, err
	}
	var result struct {
		Status string `json:"status"`
		Room   struct {
			ID string `json:"id"`
		} `json:"room"`
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return Turn{}, fmt.Errorf("reading the preview response: %w", err)
	}
	if result.Status != "" && result.Status != "success" {
		return Turn{}, fmt.Errorf("the preview runtime reported status %q", result.Status)
	}
	if result.Room.ID != "" {
		c.RoomID = result.Room.ID
	}
	return Turn{Type: kind, Message: message, Response: resp.Body}, nil
}
