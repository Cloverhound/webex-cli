package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var blockedPrefixes = []string{
	"login",
	"logout",
	"auth switch",
	"auth set-org",
	"auth clear-org",
	"config set",
	"post-install",
	"update",
	"mcp",
}

// dangerousChars are rejected in any command token to prevent shell-style injection.
// exec.Command does not shell-eval, but these characters indicate misuse.
var dangerousChars = []string{"$", "`", "|", ";", "&&"}

// readPrefixes are action verbs (and their hyphenated variants) that map to GET API calls.
var readPrefixes = []string{
	"list",
	"get",
	"download",
	"export",
	"search",
	"status",
	"describe",
	"query",
	"show",
	"fetch",
}

func registerTools(s *server.MCPServer, lg *usageLogger) {
	s.AddTool(
		mcplib.NewTool("webex_read",
			mcplib.WithDescription(
				"Execute a read-only Webex CLI command (list, get, download, export, search, status). "+
					"Maps to GET API calls — safe to auto-approve. "+
					"Provide the command without the 'webex' prefix "+
					"(e.g. 'calling people list' runs 'webex calling people list'). "+
					"Use webex_write for create/update/delete operations. "+
					"Use webex_help to discover valid commands and flags.",
			),
			mcplib.WithString("command",
				mcplib.Required(),
				mcplib.Description("Read-only CLI command without 'webex' prefix (e.g. 'calling people list', 'cc site list', 'admin people get-my-own')"),
			),
			mcplib.WithString("flags",
				mcplib.Description(`Optional flags as a JSON object string (e.g. {"max": "10", "paginate": "true"}). Keys are flag names without '--'.`),
			),
		),
		makeDispatchHandler(lg, true),
	)

	s.AddTool(
		mcplib.NewTool("webex_write",
			mcplib.WithDescription(
				"Execute a Webex CLI command that creates, updates, or deletes a resource. "+
					"Maps to POST, PUT, PATCH, or DELETE API calls — requires explicit approval. "+
					"Provide the command without the 'webex' prefix "+
					"(e.g. 'messaging messages create'). "+
					"Use webex_read for list/get/download/export operations.",
			),
			mcplib.WithString("command",
				mcplib.Required(),
				mcplib.Description("Write CLI command without 'webex' prefix (e.g. 'messaging messages create', 'cc team update', 'admin people delete')"),
			),
			mcplib.WithString("flags",
				mcplib.Description(`Optional flags as a JSON object string (e.g. {"room-id": "xxx", "text": "Hello"}). Keys are flag names without '--'.`),
			),
		),
		makeDispatchHandler(lg, false),
	)

	s.AddTool(
		mcplib.NewTool("webex_help",
			mcplib.WithDescription(
				"Get help text for any Webex CLI command or command group. "+
					"Returns usage, available flags, and subcommand descriptions. "+
					"Use this to explore what commands exist before calling webex_read or webex_write.",
			),
			mcplib.WithString("command",
				mcplib.Description("Command path to get help for (e.g. 'calling people', 'cc site'). Omit for top-level CLI help."),
			),
		),
		handleHelp,
	)

	s.AddTool(
		mcplib.NewTool("webex_usage",
			mcplib.WithDescription(
				"Query the MCP usage log to see recent commands executed via webex_read or webex_write. "+
					"Returns tool name, command, flags, status (ok/error), and elapsed time for each entry.",
			),
			mcplib.WithNumber("limit",
				mcplib.Description("Max log entries to return (1–100, default 20)"),
			),
			mcplib.WithString("command_filter",
				mcplib.Description("Substring to filter by command (e.g. 'calling', 'cc site')"),
			),
		),
		makeUsageHandler(lg),
	)
}

// makeDispatchHandler returns a handler for either webex_read (readOnly=true) or webex_write (readOnly=false).
func makeDispatchHandler(lg *usageLogger, readOnly bool) func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	toolName := "webex_write"
	if readOnly {
		toolName = "webex_read"
	}

	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		command := strings.TrimSpace(req.GetString("command", ""))
		if command == "" {
			return mcplib.NewToolResultText("Error: command is required."), nil
		}

		args := strings.Fields(command)

		// Blocklist check — join up to 2 tokens to cover two-word prefixes like "auth switch".
		prefix := strings.Join(args[:min(2, len(args))], " ")
		for _, blocked := range blockedPrefixes {
			if strings.HasPrefix(prefix, blocked) {
				return mcplib.NewToolResultText(fmt.Sprintf(
					"Error: command '%s' is not permitted via the MCP server.", blocked,
				)), nil
			}
		}

		// Injection guard: reject tokens containing shell-special characters.
		for _, arg := range args {
			for _, ch := range dangerousChars {
				if strings.Contains(arg, ch) {
					return mcplib.NewToolResultText(
						"Error: command contains disallowed characters.",
					), nil
				}
			}
		}

		// Read-only enforcement for webex_read.
		if readOnly {
			action := actionToken(args)
			if action != "" && !isReadAction(action) {
				return mcplib.NewToolResultText(fmt.Sprintf(
					"Error: '%s' is a write operation (POST/PUT/PATCH/DELETE). Use the webex_write tool instead.", action,
				)), nil
			}
		}

		// Force JSON output.
		args = append(args, "--output=json")

		// Parse and append optional flags.
		flagsStr := strings.TrimSpace(req.GetString("flags", ""))
		if flagsStr != "" {
			var flagMap map[string]string
			if err := json.Unmarshal([]byte(flagsStr), &flagMap); err != nil {
				return mcplib.NewToolResultText(
					`Error: flags must be a JSON object string, e.g. {"max": "10"}`,
				), nil
			}
			for k, v := range flagMap {
				if strings.ContainsAny(k, "= \t\n") || strings.ContainsAny(v, "\n") {
					return mcplib.NewToolResultText(
						"Error: flag key or value contains disallowed characters.",
					), nil
				}
				args = append(args, fmt.Sprintf("--%s=%s", k, v))
			}
		}

		exe, err := os.Executable()
		if err != nil {
			return mcplib.NewToolResultText("Error: cannot resolve binary path: " + err.Error()), nil
		}

		runCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()

		var stdout, stderr bytes.Buffer
		//nolint:gosec // args are validated above; exec.Command does not invoke a shell
		c := exec.CommandContext(runCtx, exe, args...)
		c.Stdout = &stdout
		c.Stderr = &stderr

		start := time.Now()
		runErr := c.Run()
		elapsed := time.Since(start).Milliseconds()

		status := "ok"
		if runErr != nil {
			status = "error"
		}

		_ = lg.Write(usageEntry{
			Time:    time.Now().UTC().Format(time.RFC3339),
			Tool:    toolName,
			Command: command,
			Flags:   flagsStr,
			Status:  status,
			Ms:      elapsed,
		})

		if runErr != nil {
			errText := strings.TrimSpace(stderr.String())
			if errText == "" {
				errText = runErr.Error()
			}
			return mcplib.NewToolResultText(fmt.Sprintf(
				"Error running 'webex %s': %s", command, errText,
			)), nil
		}
		return mcplib.NewToolResultText(strings.TrimSpace(stdout.String())), nil
	}
}

// actionToken returns the action verb from the command args.
// For "calling people list", it returns "list" (positional index 2).
// For two-token commands like "auth status", it returns "status" (positional index 1).
func actionToken(args []string) string {
	var positional []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			positional = append(positional, a)
		}
	}
	switch {
	case len(positional) >= 3:
		return strings.ToLower(positional[2])
	case len(positional) == 2:
		return strings.ToLower(positional[1])
	case len(positional) == 1:
		return strings.ToLower(positional[0])
	default:
		return ""
	}
}

// isReadAction returns true if the verb maps to a GET API call.
func isReadAction(action string) bool {
	for _, p := range readPrefixes {
		if action == p || strings.HasPrefix(action, p+"-") {
			return true
		}
	}
	return false
}

func handleHelp(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	command := strings.TrimSpace(req.GetString("command", ""))

	exe, err := os.Executable()
	if err != nil {
		return mcplib.NewToolResultText("Error: cannot resolve binary path: " + err.Error()), nil
	}

	var args []string
	if command != "" {
		args = append(args, strings.Fields(command)...)
	}
	args = append(args, "--help")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var out bytes.Buffer
	//nolint:gosec // --help invocation only; no shell eval
	c := exec.CommandContext(ctx, exe, args...)
	c.Stdout = &out
	c.Stderr = &out
	_ = c.Run()

	return mcplib.NewToolResultText(out.String()), nil
}

func makeUsageHandler(lg *usageLogger) func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	return func(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		limit := req.GetInt("limit", 20)
		if limit < 1 {
			limit = 1
		}
		if limit > 100 {
			limit = 100
		}
		cmdFilter := req.GetString("command_filter", "")

		var path string
		if lg != nil {
			path = lg.path
		}
		entries := readUsage(path, limit, cmdFilter)
		if entries == nil {
			entries = []usageEntry{}
		}
		return jsonText(map[string]any{
			"count":   len(entries),
			"entries": entries,
		}), nil
	}
}
