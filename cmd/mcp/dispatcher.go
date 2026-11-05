package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/readonly"
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
					"'resource list <type>' and 'resource get <type> <id>' read objects with their fields named, "+
					"and 'resource list-types' lists the types. "+
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

	// Proposing never writes to Webex, so it is offered in read-only mode too.
	s.AddTool(
		mcplib.NewTool("webex_propose",
			mcplib.WithDescription(
				"Record a Webex CLI write command (create, update, delete) in a change set file instead of running it. "+
					"Nothing is sent to Webex. The CLI reads the object, records only the fields the command changes, "+
					"and keeps the values it saw as preconditions. A create returns a placeholder ID like \"new:sales-team\" "+
					"that later proposals can use in flags or bodies. Start the file with webex_changeset action \"new\", "+
					"then give the file to the person, who can review and apply it in Webex Admin Studio.",
			),
			mcplib.WithString("changeset",
				mcplib.Required(),
				mcplib.Description("Path of the change set file, ending in .json (e.g. /path/to/boulder.changeset.json)"),
			),
			mcplib.WithString("command",
				mcplib.Required(),
				mcplib.Description("Write CLI command without 'webex' prefix (e.g. 'calling hunt-group update', 'cc team create')"),
			),
			mcplib.WithString("flags",
				mcplib.Description(`Optional flags as a JSON object string (e.g. {"location-id": "xxx", "hunt-group-id": "yyy", "body": "{\"name\":\"Front Desk\"}"}). Keys are flag names without '--'.`),
			),
		),
		makeDispatchHandler(lg, modePropose),
	)

	s.AddTool(
		mcplib.NewTool("webex_changeset",
			mcplib.WithDescription(
				"Manage a change set file: \"new\" starts one for the current org, \"show\" summarizes it, "+
					"\"validate\" checks it (live=true also re-reads Webex and reports values that changed since they were recorded), "+
					"and \"intent\" adds a plain-language summary and assigns operations to it.",
			),
			mcplib.WithString("action",
				mcplib.Required(),
				mcplib.Enum("new", "show", "validate", "intent"),
			),
			mcplib.WithString("path",
				mcplib.Required(),
				mcplib.Description("Path of the change set file, ending in .json"),
			),
			mcplib.WithString("title", mcplib.Description("new: a short title")),
			mcplib.WithString("rationale", mcplib.Description("new: why these changes are needed")),
			mcplib.WithString("on_behalf_of", mcplib.Description("new: the person you are working for")),
			mcplib.WithBoolean("live", mcplib.Description("validate: also compare preconditions with Webex")),
			mcplib.WithString("intent_id", mcplib.Description("intent: a short ID such as \"front-desk\"")),
			mcplib.WithString("summary", mcplib.Description("intent: what the operations achieve, in plain words")),
			mcplib.WithString("ops", mcplib.Description("intent: space-separated operation IDs to assign, such as \"op-1 op-2\"")),
		),
		makeChangesetHandler(lg),
	)

	if !readOnly {
		s.AddTool(
			mcplib.NewTool("webex_write",
				mcplib.WithDescription(
					"Execute a Webex CLI command that creates, updates, or deletes a resource. "+
						"Maps to POST, PUT, PATCH, or DELETE API calls — requires explicit approval. "+
						"Provide the command without the 'webex' prefix "+
						"(e.g. 'messaging messages create'). "+
						"To change an object of a type 'resource list-types' lists, prefer "+
						"'resource set <type> <id>' with flags {\"values\": \"{\\\"field\\\": value}\"}: "+
						"it keeps the fields you do not name, where a raw Contact Center update replaces the whole object, "+
						"and it reads the result back. 'resource describe <type>' lists the fields. "+
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
	modePropose
)

var toolNames = map[dispatchMode]string{modeRead: "webex_read", modeWrite: "webex_write", modePropose: "webex_propose"}

// makeDispatchHandler returns the handler for webex_read, webex_write or webex_propose.
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

		switch action := actionToken(args); mode {
		case modeRead:
			if action != "" && !readonly.IsReadAction(action) {
				return mcplib.NewToolResultText(fmt.Sprintf(
					"Error: '%s' is a write operation (POST/PUT/PATCH/DELETE). Use the webex_write tool instead.", action,
				)), nil
			}
		case modePropose:
			if strings.HasPrefix(commandPath(args), "changeset") {
				return mcplib.NewToolResultText("Error: use the webex_changeset tool to manage the file."), nil
			}
			path, errText := changesetPath(req.GetString("changeset", ""))
			if errText != "" {
				return mcplib.NewToolResultText(errText), nil
			}
			args = append(args, "--changeset="+path)
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
				if mode == modePropose && (k == "changeset" || k == "dry-run") {
					return mcplib.NewToolResultText("Error: the changeset parameter names the file; do not pass --" + k + "."), nil
				}
				args = append(args, fmt.Sprintf("--%s=%s", k, v))
			}
		}

		return runCLI(ctx, lg, toolName, command, flagsStr, args, mode == modePropose), nil
	}
}

// runCLI runs the CLI binary with args, logs the call and returns its output.
// allOutput adds stdout to a failure, for commands that report a failed check
// on stdout.
func runCLI(ctx context.Context, lg *usageLogger, toolName, command, flagsStr string, args []string, allOutput bool) *mcplib.CallToolResult {
	exe, err := os.Executable()
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
		if allOutput {
			errText = strings.TrimSpace(stdout.String() + "\n" + errText)
		}
		if errText == "" {
			errText = runErr.Error()
		}
		return mcplib.NewToolResultText(fmt.Sprintf(
			"Error running 'webex %s': %s", command, errText,
		))
	}
	return successResult(stdout.String(), stderr.String())
}

// changesetPath checks a change set path from a tool call and makes it absolute,
// since the CLI runs from the server's working directory.
func changesetPath(path string) (string, string) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", "Error: changeset is required."
	}
	if !strings.HasSuffix(path, ".json") {
		return "", "Error: the change set path must end in .json."
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", "Error: " + err.Error()
	}
	return abs, ""
}

func makeChangesetHandler(lg *usageLogger) func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		action := req.GetString("action", "")
		path, errText := changesetPath(req.GetString("path", ""))
		if errText != "" {
			return mcplib.NewToolResultText(errText), nil
		}
		args := []string{"changeset", action, path}
		switch action {
		case "new":
			for flag, param := range map[string]string{"title": "title", "rationale": "rationale", "on-behalf-of": "on_behalf_of"} {
				if v := req.GetString(param, ""); v != "" {
					args = append(args, "--"+flag+"="+v)
				}
			}
			args = append(args, "--author=webex-mcp")
		case "show":
		case "validate":
			if req.GetBool("live", false) {
				args = append(args, "--live")
			}
		case "intent":
			id, summary := req.GetString("intent_id", ""), req.GetString("summary", "")
			if id == "" || summary == "" {
				return mcplib.NewToolResultText("Error: intent needs intent_id and summary."), nil
			}
			args = append(args, id, summary)
			args = append(args, strings.Fields(req.GetString("ops", ""))...)
		default:
			return mcplib.NewToolResultText(`Error: action must be "new", "show", "validate" or "intent".`), nil
		}
		return runCLI(ctx, lg, "webex_changeset", "changeset "+action, "", args, true), nil
	}
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
