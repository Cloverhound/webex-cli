package mcp

import (
	"fmt"
	"os"
	"path/filepath"

	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var (
	flagLogPath     string
	flagLogMaxBytes int64
	flagLogMaxFiles int
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a stdio MCP server exposing Webex tools to AI clients",
	Long: `Start a Model Context Protocol (MCP) server over stdio.

AI clients such as Claude Code can connect to this server to query and
manage your Webex environment using natural language.

Register with Claude Code:
  claude mcp add webex -- webex mcp serve

Authentication is shared with the CLI — run 'webex login' once and the server
uses the same stored credentials with automatic token refresh.

Tools:
  webex_run    Execute any CLI command and return JSON output
  webex_help   Get help text for any command or command group
  webex_usage  Query the MCP usage log (recent commands, timing)

Resources:
  webex://commands  JSON array of all available CLI commands
  webex://usage     Last 50 lines of the raw usage log

Usage log flags:
  --log-path        Path for the usage log (default ~/.webex-mcp/usage.log)
  --log-max-size    Max log file size in bytes before rotation (default 5MB)
  --log-max-files   Number of rotated log files to keep (default 3)`,
	RunE: func(c *cobra.Command, args []string) error {
		logPath := flagLogPath
		if logPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				home = "."
			}
			logPath = filepath.Join(home, ".webex-mcp", "usage.log")
		}

		lg, err := openLogger(logPath, flagLogMaxBytes, flagLogMaxFiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: usage log unavailable: %v\n", err)
		}
		defer lg.Close()

		s := server.NewMCPServer(
			"Webex MCP",
			"2.0.0",
			server.WithToolCapabilities(false),
		)
		registerResources(s, logPath)
		registerTools(s, lg)

		fmt.Fprintln(os.Stderr, "Webex MCP server ready (stdio transport)")
		return server.ServeStdio(s)
	},
}

func init() {
	serveCmd.Flags().StringVar(&flagLogPath, "log-path", "", "Usage log path (default ~/.webex-mcp/usage.log)")
	serveCmd.Flags().Int64Var(&flagLogMaxBytes, "log-max-size", 5*1024*1024, "Max log file size in bytes before rotation")
	serveCmd.Flags().IntVar(&flagLogMaxFiles, "log-max-files", 3, "Number of rotated log files to keep")
	cmd.McpCmd.AddCommand(serveCmd)
}
