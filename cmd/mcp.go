package cmd

import "github.com/spf13/cobra"

var McpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "MCP server integration",
	Long:  "Commands for integrating with MCP-compatible AI clients (Claude Code, etc.).",
}

func init() {
	rootCmd.AddCommand(McpCmd)
}
