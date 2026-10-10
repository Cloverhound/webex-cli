package mcp

import (
	"testing"

	"github.com/mark3labs/mcp-go/server"
)

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
