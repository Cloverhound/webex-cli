package aiagent

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

func studio(t *testing.T) (webexapi.Hosts, *[]*http.Request) {
	t.Helper()
	var seen []*http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r)
		switch r.URL.Path {
		case "/api/v2/organizations/org-1/agents/":
			json.NewEncoder(w).Encode(map[string]any{
				"meta": map[string]any{"total_count": 2},
				"objects": []map[string]string{
					{"id": "a1", "name": "One", "bot_access_token": "token-1"},
					{"id": "a2", "name": "Two"},
				},
			})
		case "/api/v1/llmagent/export/":
			io := `{"bot":{"bot_access_token":"token-1","name":"One"}}`
			w.Write([]byte(io))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return webexapi.Hosts{AIAgent: srv.URL}, &seen
}

func TestBotTokenFindsTheAgent(t *testing.T) {
	hosts, _ := studio(t)
	got, err := BotToken(context.Background(), &webexapi.Client{}, hosts, "org-1", "a1")
	if err != nil || got != "token-1" {
		t.Fatalf("BotToken = %q, %v", got, err)
	}
}

func TestBotTokenNamesAMissingAgent(t *testing.T) {
	hosts, _ := studio(t)
	c := &webexapi.Client{}
	if _, err := BotToken(context.Background(), c, hosts, "org-1", "zz"); err == nil || !strings.Contains(err.Error(), "zz") {
		t.Fatalf("err = %v", err)
	}
	if _, err := BotToken(context.Background(), c, hosts, "org-1", "a2"); err == nil {
		t.Fatal("an agent without a token was accepted")
	}
	if _, err := BotToken(context.Background(), c, hosts, "", "a1"); err == nil {
		t.Fatal("no organization was accepted")
	}
}

func TestCallSendsTheAgentTokenAndRedactsTheResponse(t *testing.T) {
	hosts, seen := studio(t)
	resp, err := Call(context.Background(), &webexapi.Client{}, hosts, "org-1", webexops.CCAiAgentExport, webexops.Args{}, "a1")
	if err != nil {
		t.Fatal(err)
	}
	last := (*seen)[len(*seen)-1]
	if last.Header.Get("bot-access-token") != "token-1" {
		t.Fatalf("token header = %q", last.Header.Get("bot-access-token"))
	}
	if strings.Contains(string(resp.Body), "token-1") || !strings.Contains(string(resp.Body), `"name":"One"`) {
		t.Fatalf("body = %s", resp.Body)
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
