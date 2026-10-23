package mcp

import (
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

func TestBlockedCommand(t *testing.T) {
	cases := map[string]string{
		"auth export":            "auth export",
		"auth export --yes":      "auth export",
		"login":                  "login",
		"auth status":            "",
		"calling people list":    "",
		"cc recordings export x": "",
	}
	for command, want := range cases {
		if got := blockedCommand(strings.Fields(command)); got != want {
			t.Errorf("blockedCommand(%q) = %q, want %q", command, got, want)
		}
	}
}

func TestReadOnlyServerOmitsWriteTool(t *testing.T) {
	for _, readOnly := range []bool{false, true} {
		s := server.NewMCPServer("test", "0", server.WithToolCapabilities(false))
		registerTools(s, nil, readOnly)

		if s.GetTool("webex_read") == nil {
			t.Errorf("readOnly=%v: webex_read missing", readOnly)
		}
		if hasWrite := s.GetTool("webex_write") != nil; hasWrite == readOnly {
			t.Errorf("readOnly=%v: webex_write registered = %v", readOnly, hasWrite)
		}
	}
}
