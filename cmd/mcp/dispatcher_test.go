package mcp

import (
	"context"
	"errors"
	"strings"
	"testing"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func help(t *testing.T, command string) string {
	t.Helper()
	saved := executable
	executable = func() (string, error) { return "", errors.New("the CLI was run") }
	t.Cleanup(func() { executable = saved })
	var req mcplib.CallToolRequest
	req.Params.Arguments = map[string]any{"command": command}
	res, _ := handleHelp(context.Background(), req)
	return strings.Join(texts(t, res), "\n")
}

// A trailing value flag would take --help as its value, and the command would run.
func TestHelpRefusesFlags(t *testing.T) {
	for _, command := range []string{
		"messaging messages delete --message-id",
		"cc team delete --id x --output",
		"--debug cc team",
	} {
		if got := help(t, command); !strings.HasPrefix(got, "Error: give the command without flags") {
			t.Errorf("help(%q) = %q", command, got)
		}
	}
	if got := help(t, "cc team"); !strings.Contains(got, "the CLI was run") {
		t.Errorf("help(%q) = %q, want the CLI run", "cc team", got)
	}
}

func TestBlockedCommand(t *testing.T) {
	cases := map[string]string{
		"auth export":               "auth export",
		"auth export --yes":         "auth export",
		"--debug=false auth export": "auth export",
		"--output json auth export": "auth export",
		"auth --debug export":       "auth export",
		"login":                     "login",
		"--user someone login":      "login",
		"auth status":               "",
		"calling people list":       "",
		"cc recordings export x":    "",
	}
	for command, want := range cases {
		if got := blockedCommand(strings.Fields(command)); got != want {
			t.Errorf("blockedCommand(%q) = %q, want %q", command, got, want)
		}
	}
}

func TestActionTokenIgnoresFlagValues(t *testing.T) {
	cases := map[string]string{
		"auth status":                                "status",
		"--debug=false auth status":                  "status",
		"--user a --user b --user list config set x": "set",
	}
	for command, want := range cases {
		if got := actionToken(strings.Fields(command)); got != want {
			t.Errorf("actionToken(%q) = %q, want %q", command, got, want)
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
