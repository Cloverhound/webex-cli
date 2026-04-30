package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"strings"

	cmd "github.com/Cloverhound/webex-cli/cmd"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

type commandEntry struct {
	Path  string `json:"path"`
	Short string `json:"short"`
}

func registerResources(s *server.MCPServer, logPath string) {
	s.AddResource(
		mcplib.NewResource(
			"webex://commands",
			"Webex CLI Command Tree",
			mcplib.WithResourceDescription("JSON array of all available CLI commands with their paths and short descriptions"),
			mcplib.WithMIMEType("application/json"),
		),
		func(_ context.Context, _ mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
			cmds := walkCommands(cmd.RootCommand(), "")
			b, err := json.MarshalIndent(cmds, "", "  ")
			if err != nil {
				b = []byte("[]")
			}
			return []mcplib.ResourceContents{
				mcplib.TextResourceContents{
					URI:      "webex://commands",
					MIMEType: "application/json",
					Text:     string(b),
				},
			}, nil
		},
	)

	s.AddResource(
		mcplib.NewResource(
			"webex://usage",
			"Webex MCP Usage Log",
			mcplib.WithResourceDescription("Last 50 raw JSONL lines from the MCP usage log"),
			mcplib.WithMIMEType("text/plain"),
		),
		func(_ context.Context, _ mcplib.ReadResourceRequest) ([]mcplib.ResourceContents, error) {
			text := tailLog(logPath, 50)
			return []mcplib.ResourceContents{
				mcplib.TextResourceContents{
					URI:      "webex://usage",
					MIMEType: "text/plain",
					Text:     text,
				},
			}, nil
		},
	)
}

func walkCommands(c *cobra.Command, parentPath string) []commandEntry {
	var results []commandEntry
	for _, child := range c.Commands() {
		if child.Hidden || child.Name() == "completion" || child.Name() == "help" || child.Name() == "mcp" {
			continue
		}
		name := cmdUseName(child)
		path := name
		if parentPath != "" {
			path = parentPath + " " + name
		}
		var visibleSubs int
		for _, sc := range child.Commands() {
			if !sc.Hidden {
				visibleSubs++
			}
		}
		if visibleSubs == 0 {
			results = append(results, commandEntry{Path: path, Short: child.Short})
		} else {
			results = append(results, walkCommands(child, path)...)
		}
	}
	return results
}

// cmdUseName returns the first word of a cobra command's Use field (strips positional args).
func cmdUseName(c *cobra.Command) string {
	use := c.Use
	if i := strings.IndexByte(use, ' '); i != -1 {
		return use[:i]
	}
	return use
}

func tailLog(path string, n int) string {
	if path == "" {
		return ""
	}
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			lines = append(lines, line)
			if len(lines) > n {
				lines = lines[1:]
			}
		}
	}
	return strings.Join(lines, "\n")
}
