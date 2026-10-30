package cc

import (
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/config"
)

// agentGreetingServer serves one greeting purpose page per query page number,
// fails any page beyond those, and records the agentPersonalGreetingInfo part
// of each upload.
func agentGreetingServer(t *testing.T, pages ...string) (string, *[]map[string]any) {
	t.Helper()
	var infos []map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && r.URL.Path == "/organization/org1/v2/greeting-purpose":
			page, _ := strconv.Atoi(r.URL.Query().Get("page"))
			if page >= len(pages) {
				t.Errorf("unexpected greeting purpose page %d", page)
				w.WriteHeader(500)
				return
			}
			w.Write([]byte(pages[page]))
		case r.Method == "POST" && r.URL.Path == "/organization/org1/v2/agent-personal-greeting":
			_, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
			body, _ := io.ReadAll(r.Body)
			form, err := multipart.NewReader(strings.NewReader(string(body)), params["boundary"]).ReadForm(1 << 20)
			if err != nil || len(form.Value["agentPersonalGreetingInfo"]) != 1 {
				t.Errorf("form = %v (err=%v)", form, err)
				return
			}
			var info map[string]any
			json.Unmarshal([]byte(form.Value["agentPersonalGreetingInfo"][0]), &info)
			infos = append(infos, info)
			w.Write([]byte(`{}`))
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	token := config.Token()
	config.SetToken("test-token")
	t.Cleanup(func() {
		srv.Close()
		config.SetToken(token)
	})
	return srv.URL, &infos
}

func writeGreeting(t *testing.T) string {
	t.Helper()
	file := filepath.Join(t.TempDir(), "greeting.wav")
	if err := os.WriteFile(file, []byte("RIFFwave"), 0600); err != nil {
		t.Fatal(err)
	}
	return file
}

// The v2 endpoint rejects greetings without a purpose, so the org's system
// default purpose is used when none is given.
func TestAgentGreetingUploadDefaultsToSystemPurpose(t *testing.T) {
	url, infos := agentGreetingServer(t, `{"data":[
		{"id":"custom","name":"Sales","systemDefault":false},
		{"id":"default-purpose","name":"Default","systemDefault":true}]}`)
	if err := uploadAgentGreeting(url, "org1", writeGreeting(t), "", "agent1", ""); err != nil {
		t.Fatal(err)
	}
	if len(*infos) != 1 || (*infos)[0]["greetingPurposeId"] != "default-purpose" || (*infos)[0]["agentId"] != "agent1" {
		t.Errorf("info = %v", *infos)
	}
}

// The lookup pages until it finds the default and reads no further.
func TestAgentGreetingUploadFindsDefaultPurposeOnLaterPage(t *testing.T) {
	url, infos := agentGreetingServer(t,
		`{"data":[{"id":"custom","systemDefault":false}],"meta":{"page":0,"totalPages":3}}`,
		`{"data":[{"id":"default-purpose","systemDefault":true}],"meta":{"page":1,"totalPages":3}}`)
	if err := uploadAgentGreeting(url, "org1", writeGreeting(t), "", "agent1", ""); err != nil {
		t.Fatal(err)
	}
	if len(*infos) != 1 || (*infos)[0]["greetingPurposeId"] != "default-purpose" {
		t.Errorf("info = %v", *infos)
	}
}

func TestAgentGreetingUploadKeepsExplicitPurpose(t *testing.T) {
	url, infos := agentGreetingServer(t, `{"data":[]}`)
	if err := uploadAgentGreeting(url, "org1", writeGreeting(t), "Hi", "agent1", "custom"); err != nil {
		t.Fatal(err)
	}
	if len(*infos) != 1 || (*infos)[0]["greetingPurposeId"] != "custom" || (*infos)[0]["name"] != "Hi.wav" {
		t.Errorf("info = %v", *infos)
	}
}

func TestAgentGreetingUploadWithoutDefaultPurposeSendsNothing(t *testing.T) {
	url, infos := agentGreetingServer(t, `{"data":[{"id":"custom","name":"Sales","systemDefault":false}]}`)
	err := uploadAgentGreeting(url, "org1", writeGreeting(t), "", "agent1", "")
	if err == nil || !strings.Contains(err.Error(), "--greeting-purpose-id") || len(*infos) != 0 {
		t.Errorf("err=%v uploads=%d", err, len(*infos))
	}
}
