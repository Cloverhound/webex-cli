package mcp

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running HTTP MCP server",
	Long: `Stop a running 'webex mcp serve --http' server.

Reads the PID from ~/.webex-mcp/server.pid, sends SIGTERM, and the server
shuts down gracefully. The PID file is removed automatically on shutdown.`,
	RunE: func(c *cobra.Command, args []string) error {
		pidPath := defaultPIDPath()

		data, err := os.ReadFile(pidPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("no running MCP server found (PID file not at %s)", pidPath)
			}
			return fmt.Errorf("reading PID file: %w", err)
		}

		pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			os.Remove(pidPath)
			return fmt.Errorf("invalid PID file contents — removed stale file")
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			os.Remove(pidPath)
			return fmt.Errorf("process %d not found: %w", pid, err)
		}

		// Signal 0 checks process existence without side effects.
		if err := proc.Signal(syscall.Signal(0)); err != nil {
			os.Remove(pidPath)
			return fmt.Errorf("MCP server (PID %d) is not running — removed stale PID file", pid)
		}

		if err := proc.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("stopping MCP server (PID %d): %w", pid, err)
		}

		fmt.Printf("Stopped MCP server (PID %d)\n", pid)
		return nil
	},
}

func init() {
	cmd.McpCmd.AddCommand(stopCmd)
}
