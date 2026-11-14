package aiagent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/config"
)

func studio(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/organizations/org-1/agents/" {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"meta": map[string]any{"total_count": 2},
			"objects": []map[string]string{
				{"id": "a1", "name": "One", "bot_access_token": "token-1"},
				{"id": "a2", "name": "Two"},
			},
		})
	}))
	t.Cleanup(srv.Close)
	old := studioURL
	studioURL = func() string { return srv.URL }
	t.Cleanup(func() { studioURL = old })
	config.SetOrgID("org-1")
	t.Cleanup(func() { config.SetOrgID("") })
	return srv
}

func TestBotTokenFindsTheAgent(t *testing.T) {
	studio(t)
	got, err := BotToken("a1")
	if err != nil || got != "token-1" {
		t.Fatalf("BotToken = %q, %v", got, err)
	}
}

func TestBotTokenUsesTheSelectedRegionsStudio(t *testing.T) {
	config.SetCCRegion("eu1")
	t.Cleanup(func() { config.SetCCRegion("") })
	if got := studioURL(); got != "https://studio.aiagent-eu1.cisco.com" {
		t.Errorf("BotToken's Studio host = %s, want eu1's", got)
	}
}

func TestBotTokenNamesAMissingAgent(t *testing.T) {
	studio(t)
	if _, err := BotToken("zz"); err == nil || !strings.Contains(err.Error(), "zz") {
		t.Fatalf("err = %v", err)
	}
	if _, err := BotToken("a2"); err == nil {
		t.Fatal("an agent without a token was accepted")
	}
}

func TestRedactSecretsKeepsTheRest(t *testing.T) {
	in := `{"objects":[{"id":"a1","bot_access_token":"tok\"en","meta":{"auth_token":"x"}}],"user_access_token":null}`
	got := string(RedactSecrets([]byte(in)))
	want := `{"objects":[{"id":"a1","bot_access_token":"<redacted>","meta":{"auth_token":"<redacted>"}}],"user_access_token":null}`
	if got != want {
		t.Fatalf("got  %s\nwant %s", got, want)
	}
}
