package mcp

import (
	"fmt"
	"net/http"
	"os"
	"time"

	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a stdio MCP server exposing Webex tools to AI clients",
	Long: `Start a Model Context Protocol (MCP) server over stdio.

AI clients such as Claude Code can connect to this server to query and
manage your Webex environment using natural language.

Register with Claude Code:
  claude mcp add webex -- webex mcp serve

Authentication is handled by the CLI's existing keyring — run 'webex login' first.`,
	RunE: func(c *cobra.Command, args []string) error {
		// Set a per-request timeout for all HTTP calls made during tool handling.
		// http.DefaultClient has no timeout by default; 30s matches the prototype.
		http.DefaultClient.Timeout = 30 * time.Second

		s := server.NewMCPServer(
			"Webex MCP",
			"2.0.0",
			server.WithToolCapabilities(false),
		)
		registerTools(s)

		fmt.Fprintln(os.Stderr, "Webex MCP server ready (stdio transport)")
		return server.ServeStdio(s)
	},
}

func init() {
	cmd.McpCmd.AddCommand(serveCmd)
}
