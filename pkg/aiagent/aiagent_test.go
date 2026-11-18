package aiagent

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

func studio(t *testing.T) string {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/organizations/org-1/agents/" {
			http.NotFound(w, r)
			return
		}
		// Two pages of one agent each.
		if r.URL.Query().Get("offset") == "0" {
			fmt.Fprint(w, `{"objects":[{"id":"a1","bot_access_token":"t1"}],"meta":{"total_count":2}}`)
			return
		}
		fmt.Fprint(w, `{"objects":[{"id":"a2","bot_access_token":""}],"meta":{"total_count":2}}`)
	}))
	t.Cleanup(srv.Close)
	return srv.URL
}

func TestBotToken(t *testing.T) {
	url := studio(t)
	c := &webexapi.Client{Origins: []string{url}}
	if token, err := BotToken(context.Background(), c, url, "org-1", "a1"); err != nil || token != "t1" {
		t.Fatalf("token %q, err %v", token, err)
	}
	for agent, want := range map[string]string{"a2": "has no runtime token", "a3": "not found in organization org-1"} {
		if _, err := BotToken(context.Background(), c, url, "org-1", agent); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: err %v, want %q", agent, err, want)
		}
	}
	if _, err := BotToken(context.Background(), c, url, "", "a1"); err == nil {
		t.Error("looked up an agent with no organization")
	}
}

func TestRedactSecretsKeepsTheRestOfTheDocument(t *testing.T) {
	in := `{"id":"a","bot_access_token":"x\"y","nested":{"auth_token": "z"}}`
	want := `{"id":"a","bot_access_token":"<redacted>","nested":{"auth_token": "<redacted>"}}`
	if got := string(RedactSecrets([]byte(in))); got != want {
		t.Fatalf("got %s", got)
	}
}

func TestRedactSecretsHidesAnyValueInJSONHeldInAString(t *testing.T) {
	for _, secret := range []string{`plain`, "new\nline\ttab", `quote"back\slash\\`, "\u0001ctl é ✓ </script>", `\"`, `\\"`} {
		inner, _ := json.Marshal(map[string]string{"clientSecret": secret, "url": "u"})
		outer, _ := json.Marshal(map[string]string{"config": string(inner)})
		var got struct{ Config string }
		if err := json.Unmarshal(RedactSecrets(outer), &got); err != nil {
			t.Fatalf("%q: the redacted document is not JSON: %v", secret, err)
		}
		var config map[string]string
		if err := json.Unmarshal([]byte(got.Config), &config); err != nil || config["clientSecret"] != "<redacted>" || config["url"] != "u" {
			t.Errorf("%q: config %v, err %v", secret, config, err)
		}
	}
}

func TestRedactSecretsFindsEveryCredentialKey(t *testing.T) {
	for in, want := range map[string]string{
		`{"clientSecret":"s","SIPPassword":"p","refresh_token":"r","tokenLifetimeMinutes":30}`: `{"clientSecret":"<redacted>","SIPPassword":"<redacted>","refresh_token":"<redacted>","tokenLifetimeMinutes":30}`,
		// A tool's configuration is JSON held in a string.
		`{"config":"{\"api_token\":\"a\\\"b\",\"url\":\"u\"}"}`:             `{"config":"{\"api_token\":\"<redacted>\",\"url\":\"u\"}"}`,
		`{"config":"{\"clientSecret\":\"a\\nb\\tc\\\\d\\u0041\",\"u\":1}"}`: `{"config":"{\"clientSecret\":\"<redacted>\",\"u\":1}"}`,
		`{"description":"send the token","name":"n"}`:                       `{"description":"send the token","name":"n"}`,
	} {
		if got := string(RedactSecrets([]byte(in))); got != want {
			t.Errorf("%s\n got %s\nwant %s", in, got, want)
		}
	}
}
