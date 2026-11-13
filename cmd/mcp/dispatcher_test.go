package mcp

import (
	"context"
	"errors"
	"path/filepath"
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
		for _, name := range []string{"webex_propose", "webex_changeset"} {
			if s.GetTool(name) == nil {
				t.Errorf("readOnly=%v: %s missing", readOnly, name)
			}
		}
	}
}

func propose(t *testing.T, command string) string {
	t.Helper()
	saved := executable
	executable = func() (string, error) { return "", errors.New("test: the command was not refused") }
	t.Cleanup(func() { executable = saved })
	var req mcplib.CallToolRequest
	req.Params.Arguments = map[string]any{"command": command, "changeset": filepath.Join(t.TempDir(), "p.changeset.json")}
	res, err := makeDispatchHandler(nil, modePropose)(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	return res.Content[0].(mcplib.TextContent).Text
}

func TestProposeRefusesWhatItCannotRecord(t *testing.T) {
	for _, command := range []string{
		"cc team delete --id=team1 --", // flags after -- would be arguments, so nothing turns capture on
		"auth token",
		"auth set-folder-default someone",
		"calling people list",
		"calling", // resolves to no command
		"resource set cc_team team1 name=x",
		"admin people update --person-id=p1", // no type Studio applies lives outside calling and cc
	} {
		if got := propose(t, command); !strings.HasPrefix(got, "Error:") || strings.Contains(got, "was not refused") {
			t.Errorf("%q was run: %s", command, got)
		}
	}
	for _, command := range []string{"cc team delete --id=team1", "calling hunt-group update --location-id=l1 --hunt-group-id=h1"} {
		if got := propose(t, command); !strings.Contains(got, "was not refused") {
			t.Errorf("%q was refused: %s", command, got)
		}
	}
}

func TestChangesetPath(t *testing.T) {
	if _, errText := changesetPath("notes.txt"); errText == "" {
		t.Error("accepted a path that is not .json")
	}
	if path, errText := changesetPath("x.changeset.json"); errText != "" || !filepath.IsAbs(path) {
		t.Errorf("path = %q, err = %q", path, errText)
	}
}
