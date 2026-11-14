package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/readonly"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

var findMu sync.Mutex

var blockedPrefixes = []string{
	"login",
	"logout",
	"auth export",
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

func registerTools(s *server.MCPServer, lg *usageLogger, readOnly bool) {
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
		makeDispatchHandler(lg, modeRead),
	)

	if !readOnly {
		s.AddTool(
			mcplib.NewTool("webex_write",
				mcplib.WithDescription(
					"Execute a Webex CLI command that creates, updates, or deletes a resource. "+
						"Maps to POST, PUT, PATCH, or DELETE API calls — requires explicit approval. "+
						"Provide the command without the 'webex' prefix "+
						"(e.g. 'messaging messages create'). "+
						"A Contact Center update replaces the whole object: fields you leave out are cleared. "+
						"Get the object first and send the full body. "+
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
			makeDispatchHandler(lg, modeWrite),
		)
	}

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

type dispatchMode int

const (
	modeRead dispatchMode = iota
	modeWrite
)

var toolNames = map[dispatchMode]string{modeRead: "webex_read", modeWrite: "webex_write"}

// makeDispatchHandler returns the handler for webex_read or webex_write.
func makeDispatchHandler(lg *usageLogger, mode dispatchMode) func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	toolName := toolNames[mode]

	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		command := strings.TrimSpace(req.GetString("command", ""))
		if command == "" {
			return mcplib.NewToolResultText("Error: command is required."), nil
		}

		args := strings.Fields(command)

		if blocked := blockedCommand(args); blocked != "" {
			return mcplib.NewToolResultText(fmt.Sprintf(
				"Error: command '%s' is not permitted via the MCP server.", blocked,
			)), nil
		}

		// Injection guard: reject tokens containing shell-special characters, and "--",
		// after which the flags the server adds would be read as arguments.
		for _, arg := range args {
			if arg == "--" {
				return mcplib.NewToolResultText("Error: command contains disallowed characters."), nil
			}
			for _, ch := range dangerousChars {
				if strings.Contains(arg, ch) {
					return mcplib.NewToolResultText(
						"Error: command contains disallowed characters.",
					), nil
				}
			}
		}

		if action := actionToken(args); mode == modeRead && action != "" && !readonly.IsReadAction(action) {
			return mcplib.NewToolResultText(fmt.Sprintf(
				"Error: '%s' is a write operation (POST/PUT/PATCH/DELETE). Use the webex_write tool instead.", action,
			)), nil
		}
		// A command's verb can read like a read while its request writes, as
		// 'calling call-controls get' POSTs to retrieve a parked call, so a generated
		// command is judged by its method.
		if op, ok := operationOf(args); mode == modeRead && ok && op.Method != http.MethodGet && !readonly.IsQueryPOST(op.Path) {
			return mcplib.NewToolResultText(fmt.Sprintf(
				"Error: '%s' sends a %s request. Use the webex_write tool instead.", commandPath(args), op.Method,
			)), nil
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

		return runCLI(ctx, lg, toolName, command, flagsStr, args), nil
	}
}

// executable is the CLI binary the tools run. Tests replace it: under go test
// os.Executable is the test binary, which would run the tests again.
var executable = os.Executable

// runCLI runs the CLI binary with args, logs the call and returns its output.
func runCLI(ctx context.Context, lg *usageLogger, toolName, command, flagsStr string, args []string) *mcplib.CallToolResult {
	exe, err := executable()
	if err != nil {
		return mcplib.NewToolResultText("Error: cannot resolve binary path: " + err.Error())
	}

	runCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var stdout, stderr bytes.Buffer
	//nolint:gosec // args are validated by the callers; exec.Command does not invoke a shell
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
		))
	}
	return successResult(stdout.String(), stderr.String())
}

// successResult returns a command's output, and its stderr as a second text
// item when there is any. Notices such as the private-API caution are printed
// there once, so dropping stderr would leave agents never seeing them.
func successResult(stdout, stderr string) *mcplib.CallToolResult {
	result := mcplib.NewToolResultText(strings.TrimSpace(stdout))
	if note := strings.TrimSpace(stderr); note != "" {
		result.Content = append(result.Content, mcplib.NewTextContent(note))
	}
	return result
}

// commandPath returns the command the CLI will run for args, without the binary
// name, or "" if args name none, in which case the CLI rejects them too. Cobra's
// own lookup is used so the checks below cannot disagree with it about which
// tokens are flag values. The mutex is there because Find merges flag sets and
// tool calls can run concurrently.
func commandPath(args []string) string {
	findMu.Lock()
	defer findMu.Unlock()
	root := cmd.RootCommand()
	c, _, err := root.Find(args)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(strings.TrimPrefix(c.CommandPath(), root.Name()), " ")
}

// blockedCommand returns the blocked prefix that args resolve to, or "".
func blockedCommand(args []string) string {
	path := commandPath(args)
	for _, blocked := range blockedPrefixes {
		if strings.HasPrefix(path, blocked) {
			return blocked
		}
	}
	return ""
}

// operationOf finds the API operation a generated command sends.
func operationOf(args []string) (*webexops.Operation, bool) {
	words := strings.Fields(commandPath(args))
	if len(words) != 3 {
		return nil, false
	}
	return webexops.Find(words[0], words[1], words[2])
}

// actionToken returns the action verb from the command args: "list" for
// "calling people list" and "status" for "auth status".
func actionToken(args []string) string {
	words := strings.Fields(commandPath(args))
	if len(words) == 0 {
		return ""
	}
	return strings.ToLower(words[len(words)-1])
}

func handleHelp(_ context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	command := strings.TrimSpace(req.GetString("command", ""))

	// Only command words: a trailing value flag would swallow --help, and the
	// command would run.
	args := strings.Fields(command)
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			return mcplib.NewToolResultText("Error: give the command without flags, such as 'cc team update'."), nil
		}
	}
	args = append(args, "--help")

	exe, err := executable()
	if err != nil {
		return mcplib.NewToolResultText("Error: cannot resolve binary path: " + err.Error()), nil
	}

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
