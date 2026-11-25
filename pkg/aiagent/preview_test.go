package aiagent

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

func TestPreviewHosts(t *testing.T) {
	runtime, origin, err := PreviewHosts("https://studio.aiagent-us1.cisco.com")
	if err != nil {
		t.Fatal(err)
	}
	if runtime != "https://bot-appserver.produs1.ciscoccservice.com" || origin != "https://preview.produs1.ciscoccservice.com" {
		t.Fatalf("runtime %q, origin %q", runtime, origin)
	}
	if _, _, err := PreviewHosts("https://studio.example.com"); err == nil {
		t.Fatal("derived a runtime from a host that is not AI Agent Studio")
	}
}

type previewRequest struct {
	Consumer struct {
		UID         string            `json:"uid"`
		ExtraParams map[string]string `json:"extra_params"`
	} `json:"consumer"`
	Type     string `json:"type"`
	Msg      string `json:"msg"`
	Platform string `json:"platform"`
	Merged   *bool  `json:"transparency_msg_merged"`
	IsTest   bool   `json:"is_test"`
}

func previewServer(t *testing.T, reply func(n int) string) (string, *[]previewRequest) {
	t.Helper()
	var got []previewRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/webhook/web/" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Bot-Access-Token") != "bot-secret" {
			t.Errorf("bot-access-token %q", r.Header.Get("Bot-Access-Token"))
		}
		if r.Header.Get("Origin") != "https://preview.example" {
			t.Errorf("origin %q", r.Header.Get("Origin"))
		}
		if _, ok := r.Header["Authorization"]; ok {
			t.Error("sent an Authorization header to the preview runtime")
		}
		var p previewRequest
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			t.Error(err)
		}
		got = append(got, p)
		w.Write([]byte(reply(len(got))))
	}))
	t.Cleanup(srv.Close)
	return srv.URL, &got
}

func TestPreviewStartAndSendKeepTheConversation(t *testing.T) {
	url, got := previewServer(t, func(n int) string {
		return `{"status":"success","generated_msg":[{"text":"reply"}],"room":{"id":"room-1"}}`
	})
	p := &Preview{Client: &webexapi.Client{}, Runtime: url, Origin: "https://preview.example", BotToken: "bot-secret"}
	conv, err := NewConversation("agent-1")
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`).MatchString(conv.ConsumerUID) {
		t.Fatalf("consumer UID is not a UUIDv4: %q", conv.ConsumerUID)
	}
	if !regexp.MustCompile(`^[A-Za-z0-9]{15}$`).MatchString(conv.SocketKey) {
		t.Fatalf("socket key %q", conv.SocketKey)
	}

	first, err := p.Start(context.Background(), &conv)
	if err != nil {
		t.Fatal(err)
	}
	if first.Type != "bot" || first.Message != "hi" || conv.RoomID != "room-1" {
		t.Fatalf("turn %+v, conversation %+v", first, conv)
	}
	second, err := p.Send(context.Background(), &conv, "hello")
	if err != nil {
		t.Fatal(err)
	}
	if second.Type != "human" || second.Message != "hello" || !strings.Contains(string(second.Response), "reply") {
		t.Fatalf("turn %+v", second)
	}

	if len(*got) != 2 {
		t.Fatalf("%d requests", len(*got))
	}
	open, human := (*got)[0], (*got)[1]
	if open.Type != "bot" || open.Msg != "hi" || human.Type != "human" || human.Msg != "hello" {
		t.Fatalf("requests %+v", *got)
	}
	for _, r := range *got {
		if r.Consumer.UID != conv.ConsumerUID || r.Consumer.ExtraParams["socket_key"] != conv.SocketKey ||
			r.Platform != "web" || !r.IsTest || r.Merged == nil || *r.Merged {
			t.Fatalf("request %+v", r)
		}
	}
}

func TestPreviewSendRefusesAnEmptyMessage(t *testing.T) {
	p := &Preview{Client: &webexapi.Client{}, Runtime: "https://runtime.invalid", BotToken: "bot-secret"}
	conv, _ := NewConversation("agent-1")
	if _, err := p.Send(context.Background(), &conv, "  "); err == nil {
		t.Fatal("sent an empty message")
	}
}

func TestPreviewReportsAFailedStatus(t *testing.T) {
	url, _ := previewServer(t, func(int) string { return `{"status":"error","bot_access_token":"leak"}` })
	p := &Preview{Client: &webexapi.Client{}, Runtime: url, Origin: "https://preview.example", BotToken: "bot-secret"}
	conv, _ := NewConversation("agent-1")
	if _, err := p.Start(context.Background(), &conv); err == nil || !strings.Contains(err.Error(), `"error"`) {
		t.Fatalf("err %v", err)
	}
}
