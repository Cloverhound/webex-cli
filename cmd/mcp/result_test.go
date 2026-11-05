package mcp

import (
	"testing"

	mcplib "github.com/mark3labs/mcp-go/mcp"
)

func texts(t *testing.T, r *mcplib.CallToolResult) []string {
	t.Helper()
	var out []string
	for _, c := range r.Content {
		tc, ok := c.(mcplib.TextContent)
		if !ok {
			t.Fatalf("content %T is not text", c)
		}
		out = append(out, tc.Text)
	}
	return out
}

func TestSuccessfulResultCarriesStderrSeparately(t *testing.T) {
	got := texts(t, successResult(`{"a":1}`+"\n", "Caution: private API.\n"))
	if len(got) != 2 || got[0] != `{"a":1}` || got[1] != "Caution: private API." {
		t.Fatalf("got %q", got)
	}
}

func TestSuccessfulResultWithoutStderrIsUnchanged(t *testing.T) {
	got := texts(t, successResult(`{"a":1}`, "  \n"))
	if len(got) != 1 || got[0] != `{"a":1}` {
		t.Fatalf("got %q", got)
	}
}
