package aiagent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/aiagent"
)

func TestPreviewConversationsAreKeptPerUserOrgAndAgent(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	a := aiagent.Conversation{AgentID: "a1", ConsumerUID: "u1", SocketKey: "k1", RoomID: "r1"}
	b := aiagent.Conversation{AgentID: "a2", ConsumerUID: "u2", SocketKey: "k2"}
	if err := SavePreview("User@Example.com", "org-1", a); err != nil {
		t.Fatal(err)
	}
	if err := SavePreview("user@example.com", "org-1", b); err != nil {
		t.Fatal(err)
	}

	for agent, want := range map[string]aiagent.Conversation{"a1": a, "a2": b, "": b} {
		got, err := SavedPreview("user@example.com", "ORG-1", agent)
		if err != nil || got != want {
			t.Errorf("agent %q: got %+v, err %v", agent, got, err)
		}
	}
	if agent, err := LastPreviewAgent("user@example.com", "org-1"); err != nil || agent != "a2" {
		t.Errorf("last agent %q, err %v", agent, err)
	}
	if _, err := SavedPreview("other@example.com", "org-1", ""); err == nil || !strings.Contains(err.Error(), "--agent-id") {
		t.Errorf("another user's conversation: err %v", err)
	}
	if _, err := SavedPreview("user@example.com", "org-2", "a1"); err == nil {
		t.Error("found a conversation in another organization")
	}

	info, err := os.Stat(previewStorePath())
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("store mode %o", info.Mode().Perm())
	}
}

func TestPreviewStoreRefusesAnUnknownFormat(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	path := previewStorePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(`{"version":9}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := SavePreview("u", "o", aiagent.Conversation{AgentID: "a"}); err == nil {
		t.Fatal("overwrote a store it cannot read")
	}
}

func TestBriefHidesTheOpeningRequest(t *testing.T) {
	turns := []aiagent.Turn{
		{Type: "bot", Message: "hi", Response: json.RawMessage(`{"generated_msg":[{"text":"Welcome"}]}`)},
		{Type: "human", Message: "Hello", Response: json.RawMessage(`{"generated_msg":[{"text":"Hi there"},{"text":"How can I help?"}]}`)},
		{Type: "human", Message: "Bye", Response: json.RawMessage(`{"generated_msg":[],"bot_msg":"Goodbye"}`)},
	}
	got, err := Brief(turns)
	if err != nil {
		t.Fatal(err)
	}
	want := []BriefTurn{{In: "", Out: "Welcome"}, {In: "Hello", Out: "Hi there\nHow can I help?"}, {In: "Bye", Out: "Goodbye"}}
	if len(got) != len(want) {
		t.Fatalf("got %+v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("turn %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}
