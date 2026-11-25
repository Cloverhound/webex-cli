package aiagent

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/pkg/aiagent"
)

const previewStoreVersion = 1

// previewStore holds each user's text preview conversations, so a later command
// can continue one. It holds no credentials.
type previewStore struct {
	Version int `json:"version"`
	// Last is the agent last used, by user and organization.
	Last          map[string]string               `json:"last"`
	Conversations map[string]aiagent.Conversation `json:"conversations"`
}

func previewStorePath() string { return filepath.Join(appconfig.Dir(), "ai-agent-preview.json") }

func previewKey(user, org string) string {
	return strings.ToLower(strings.TrimSpace(user)) + "::" + strings.ToLower(strings.TrimSpace(org))
}

func loadPreviewStore() (*previewStore, error) {
	s := &previewStore{Version: previewStoreVersion, Last: map[string]string{}, Conversations: map[string]aiagent.Conversation{}}
	data, err := os.ReadFile(previewStorePath())
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, s); err != nil || s.Version != previewStoreVersion {
		return nil, fmt.Errorf("cannot read the saved preview conversations in %s; delete the file to start over", previewStorePath())
	}
	if s.Last == nil {
		s.Last = map[string]string{}
	}
	if s.Conversations == nil {
		s.Conversations = map[string]aiagent.Conversation{}
	}
	return s, nil
}

// SavePreview keeps a conversation for SavedPreview, and makes its agent the
// user's last for the organization.
func SavePreview(user, org string, c aiagent.Conversation) error {
	s, err := loadPreviewStore()
	if err != nil {
		return err
	}
	key := previewKey(user, org)
	s.Last[key] = c.AgentID
	s.Conversations[key+"::"+c.AgentID] = c
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	path := previewStorePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	tmp := fmt.Sprintf("%s.%d.tmp", path, os.Getpid())
	if err := os.WriteFile(tmp, append(data, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// LastPreviewAgent returns the agent of the user's last preview conversation in
// the organization.
func LastPreviewAgent(user, org string) (string, error) {
	s, err := loadPreviewStore()
	if err != nil {
		return "", err
	}
	if agent := s.Last[previewKey(user, org)]; agent != "" {
		return agent, nil
	}
	return "", fmt.Errorf("no earlier preview conversation in this organization; pass --agent-id")
}

// SavedPreview returns the user's conversation with an agent, or with the last
// agent when agentID is empty.
func SavedPreview(user, org, agentID string) (aiagent.Conversation, error) {
	if agentID == "" {
		agent, err := LastPreviewAgent(user, org)
		if err != nil {
			return aiagent.Conversation{}, err
		}
		agentID = agent
	}
	s, err := loadPreviewStore()
	if err != nil {
		return aiagent.Conversation{}, err
	}
	c, ok := s.Conversations[previewKey(user, org)+"::"+agentID]
	if !ok {
		return aiagent.Conversation{}, fmt.Errorf("no preview conversation with agent %s; run 'webex cc ai-agent-preview start' first", agentID)
	}
	return c, nil
}

// BriefTurn is a turn's visible text: the message sent and the agent's reply.
type BriefTurn struct {
	In  string `json:"in"`
	Out string `json:"out"`
}

// Brief returns the visible text of each turn. The opening request's
// placeholder message is left out.
func Brief(turns []aiagent.Turn) ([]BriefTurn, error) {
	out := make([]BriefTurn, 0, len(turns))
	for _, t := range turns {
		var r struct {
			Generated []struct {
				Text    string `json:"text"`
				Message string `json:"message"`
			} `json:"generated_msg"`
			BotMsg string `json:"bot_msg"`
		}
		if err := json.Unmarshal(t.Response, &r); err != nil {
			return nil, fmt.Errorf("reading the preview response: %w", err)
		}
		var texts []string
		for _, g := range r.Generated {
			if text := cmp.Or(g.Text, g.Message); text != "" {
				texts = append(texts, text)
			}
		}
		if len(texts) == 0 && r.BotMsg != "" {
			texts = append(texts, r.BotMsg)
		}
		in := t.Message
		if t.Type == "bot" {
			in = ""
		}
		out = append(out, BriefTurn{In: in, Out: strings.Join(texts, "\n")})
	}
	return out, nil
}
