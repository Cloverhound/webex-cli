package mcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var (
	flagLogPath     string
	flagLogMaxBytes int64
	flagLogMaxFiles int
	flagHTTP        bool
	flagHTTPAddr    string
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start an MCP server exposing Webex tools to AI clients",
	Long: `Start a Model Context Protocol (MCP) server.

Defaults to stdio transport (Claude Code, Codex, Cursor):
  webex mcp serve

Use --http for HTTP transport (Claude Desktop and other HTTP-based clients):
  webex mcp serve --http

The HTTP server binds to loopback only and cannot be exposed on public interfaces.
Use 'webex mcp stop' to stop a running HTTP server.

Authentication is shared with the CLI — run 'webex login' once and the server
uses the same stored credentials with automatic token refresh.

Tools:
  webex_run    Execute any CLI command and return JSON output
  webex_help   Get help text for any command or command group
  webex_usage  Query the MCP usage log (recent commands, timing)

Resources:
  webex://commands  JSON array of all available CLI commands
  webex://usage     Last 50 lines of the raw usage log

HTTP transport flags:
  --http            Enable HTTP transport instead of stdio
  --http-addr       Listen address — must be loopback (default 127.0.0.1:47890)

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

		if flagHTTP {
			host, _, err := net.SplitHostPort(flagHTTPAddr)
			if err != nil {
				return fmt.Errorf("invalid --http-addr %q: %w", flagHTTPAddr, err)
			}
			ip := net.ParseIP(host)
			if ip == nil || !ip.IsLoopback() {
				return fmt.Errorf("--http-addr must bind to a loopback address (127.x.x.x or ::1), got %q", host)
			}

			pidPath := defaultPIDPath()
			if err := writePIDFile(pidPath); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not write PID file: %v\n", err)
			} else {
				defer os.Remove(pidPath)
			}

			hs := server.NewStreamableHTTPServer(s,
				server.WithEndpointPath("/mcp"),
				server.WithHeartbeatInterval(30*time.Second),
			)

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

			errCh := make(chan error, 1)
			go func() { errCh <- hs.Start(flagHTTPAddr) }()

			fmt.Fprintf(os.Stderr, "Webex MCP server ready (HTTP transport on http://%s/mcp)\n", flagHTTPAddr)
			fmt.Fprintf(os.Stderr, "Stop with: webex mcp stop\n")

			select {
			case sig := <-sigCh:
				fmt.Fprintf(os.Stderr, "\nReceived %s, shutting down...\n", sig)
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				return hs.Shutdown(ctx)
			case err := <-errCh:
				return err
			}
		}

		fmt.Fprintln(os.Stderr, "Webex MCP server ready (stdio transport)")
		return server.ServeStdio(s)
	},
}

func defaultPIDPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".webex-mcp", "server.pid")
}

func writePIDFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0644)
}

func init() {
	serveCmd.Flags().BoolVar(&flagHTTP, "http", false, "Start HTTP server instead of stdio transport")
	serveCmd.Flags().StringVar(&flagHTTPAddr, "http-addr", "127.0.0.1:47890", "HTTP listen address — must be loopback (used with --http)")
	serveCmd.Flags().StringVar(&flagLogPath, "log-path", "", "Usage log path (default ~/.webex-mcp/usage.log)")
	serveCmd.Flags().Int64Var(&flagLogMaxBytes, "log-max-size", 5*1024*1024, "Max log file size in bytes before rotation")
	serveCmd.Flags().IntVar(&flagLogMaxFiles, "log-max-files", 3, "Number of rotated log files to keep")
	cmd.McpCmd.AddCommand(serveCmd)
}
