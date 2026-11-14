package mcp

import (
	"context"
	"errors"
	"strings"
	"testing"

	_ "github.com/Cloverhound/webex-cli/cmd/calling"
	_ "github.com/Cloverhound/webex-cli/cmd/cc"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

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

// stubExecutable makes any attempt to run the CLI fail visibly, so a test sees
// whether a tool call was refused before it ran anything.
func stubExecutable(t *testing.T) {
	t.Helper()
	saved := executable
	executable = func() (string, error) { return "", errors.New("test: the command was not refused") }
	t.Cleanup(func() { executable = saved })
}

func TestReadRefusesWritesAndDoubleDash(t *testing.T) {
	stubExecutable(t)
	for _, command := range []string{"cc team delete --id=team1", "calling people list --"} {
		var req mcplib.CallToolRequest
		req.Params.Arguments = map[string]any{"command": command}
		res, err := makeDispatchHandler(nil, modeRead)(context.Background(), req)
		if err != nil {
			t.Fatal(err)
		}
		if got := res.Content[0].(mcplib.TextContent).Text; !strings.HasPrefix(got, "Error:") || strings.Contains(got, "was not refused") {
			t.Errorf("%q was run: %s", command, got)
		}
	}
}

// A trailing value flag would swallow --help and run the command.
func TestHelpRefusesFlags(t *testing.T) {
	stubExecutable(t)
	var req mcplib.CallToolRequest
	req.Params.Arguments = map[string]any{"command": "cc team delete --id x --output"}
	res, _ := handleHelp(context.Background(), req)
	if got := res.Content[0].(mcplib.TextContent).Text; !strings.HasPrefix(got, "Error: give the command without flags") {
		t.Errorf("help = %q", got)
	}
}
