package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// These tests pin the CLI's published surface: every command's name, aliases,
// help text and flags, and the exact request each API command sends for a fixed
// set of flags. Scripts and agents depend on both, so a change to how commands
// are built must leave these files unchanged. Regenerate after an intended
// change with: go test -run Golden -update .

var updateGolden = flag.Bool("update", false, "rewrite the golden files")

const goldenOrg = "6ab1d2c4-0000-4000-8000-00000000c1d0"

var apiAreas = []string{"admin", "calling", "cc", "device", "meetings", "messaging"}

// recordedRequest is one HTTP request as the server would see it.
type recordedRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    string            `json:"body,omitempty"`
}

type goldenRun struct {
	Command  string            `json:"command"`
	Variant  string            `json:"variant"`
	Args     []string          `json:"args"`
	Requests []recordedRequest `json:"requests"`
	Error    string            `json:"error,omitempty"`
}

// recorder answers every request without touching the network.
type recorder struct {
	requests []recordedRequest
}

func (r *recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := recordedRequest{Method: req.Method, URL: canonicalURL(req), Headers: map[string]string{}}
	for k, v := range req.Header {
		rec.Headers[k] = strings.Join(v, ", ")
	}
	if req.Body != nil {
		data, _ := io.ReadAll(req.Body)
		rec.Body = string(data)
	}
	r.requests = append(r.requests, rec)

	body := "{}"
	// AI Agent Studio commands look up the agent's runtime token first.
	if strings.HasSuffix(req.URL.Path, "/agents/") {
		body = `[{"id":"s-agent-id","bot_access_token":"golden-bot-token"}]`
	}
	return &http.Response{
		StatusCode: 200,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
		Request:    req,
	}, nil
}

func canonicalURL(req *http.Request) string {
	u := *req.URL
	u.RawQuery = u.Query().Encode()
	return u.String()
}

func goldenEnv(t *testing.T) *recorder {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, env := range []string{"WEBEX_TOKEN", "WEBEX_REFRESH_TOKEN", "WEBEX_USER", "WEBEX_REGION", "WEBEX_CHANGESET", "WEBEX_READ_ONLY", "WEBEX_ORG", "WEBEX_CC_REGION"} {
		t.Setenv(env, "")
	}
	t.Chdir(t.TempDir())

	rec := &recorder{}
	saved := http.DefaultTransport
	http.DefaultTransport = rec
	t.Cleanup(func() { http.DefaultTransport = saved })

	devnull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	stdout, stderr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = devnull, devnull
	t.Cleanup(func() { os.Stdout, os.Stderr = stdout, stderr; devnull.Close() })

	root := cmd.RootCommand()
	root.SetOut(io.Discard)
	root.SetErr(io.Discard)
	return rec
}

// leafCommands lists runnable commands under the API areas, by path.
func leafCommands(root *cobra.Command) []*cobra.Command {
	var out []*cobra.Command
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		if c.Runnable() && !c.HasSubCommands() {
			out = append(out, c)
		}
		for _, child := range c.Commands() {
			walk(child)
		}
	}
	for _, area := range apiAreas {
		for _, c := range root.Commands() {
			if c.Name() == area {
				walk(c)
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CommandPath() < out[j].CommandPath() })
	return out
}

// skippedFlags are left unset: orgid and org-id come from --organization, last
// depends on the clock, and body flags replace the per-field flags.
var skippedFlags = map[string]bool{"orgid": true, "org-id": true, "last": true, "body": true, "body-file": true}

func sampleValue(f *pflag.Flag) string {
	switch f.Name {
	case "from":
		return "2026-01-01T00:00:00.000Z"
	case "to":
		return "2026-01-02T00:00:00.000Z"
	}
	switch f.Value.Type() {
	case "bool":
		return "true"
	case "int", "int64", "int32":
		return "7"
	case "float64", "float32":
		return "1.5"
	case "stringSlice", "stringArray":
		return "s-" + f.Name + "-1,s-" + f.Name + "-2"
	case "duration":
		return "90s"
	default:
		return "s-" + f.Name
	}
}

func isRequired(f *pflag.Flag) bool {
	v, ok := f.Annotations[cobra.BashCompOneRequiredFlag]
	return ok && len(v) > 0 && v[0] == "true"
}

// resetFlags returns every flag to its default, since cobra keeps values set by an
// earlier run of the same command.
func resetFlags(c *cobra.Command) {
	reset := func(f *pflag.Flag) {
		if sv, ok := f.Value.(pflag.SliceValue); ok {
			sv.Replace(nil)
		} else {
			f.Value.Set(f.DefValue)
		}
		f.Changed = false
	}
	for x := c; x != nil; x = x.Parent() {
		x.Flags().VisitAll(reset)
		x.PersistentFlags().VisitAll(reset)
	}
}

func run(rec *recorder, c *cobra.Command, variant string, args []string) goldenRun {
	resetFlags(c)
	rec.requests = nil
	path := strings.Fields(c.CommandPath())[1:]
	full := append(append([]string{}, path...), args...)
	full = append(full, "--token", "golden-token", "--organization", goldenOrg)
	root := cmd.RootCommand()
	root.SetArgs(full)
	err := cmd.Execute()
	out := goldenRun{Command: strings.Join(path, " "), Variant: variant, Args: args, Requests: rec.requests}
	if out.Requests == nil {
		out.Requests = []recordedRequest{}
	}
	if err != nil {
		out.Error = err.Error()
	}
	return out
}

func TestGoldenRequests(t *testing.T) {
	rec := goldenEnv(t)
	root := cmd.RootCommand()
	// Execute strips the required mark from org flags, which decides "required" below.
	root.SetArgs([]string{"version"})
	cmd.Execute()

	byArea := map[string][]goldenRun{}
	for _, c := range leafCommands(root) {
		area := strings.Fields(c.CommandPath())[1]
		var all, required []string
		c.Flags().VisitAll(func(f *pflag.Flag) {
			if skippedFlags[f.Name] || c.InheritedFlags().Lookup(f.Name) != nil {
				return
			}
			arg := "--" + f.Name + "=" + sampleValue(f)
			all = append(all, arg)
			if isRequired(f) {
				required = append(required, arg)
			}
		})
		runs := []goldenRun{run(rec, c, "all-flags", all), run(rec, c, "required", required)}
		if c.Flags().Lookup("body") != nil {
			runs = append(runs, run(rec, c, "body", append(append([]string{}, required...), `--body={"golden":true}`)))
		}
		if hasGET(runs[1]) {
			runs = append(runs, run(rec, c, "paginate", append(append([]string{}, required...), "--paginate")))
		}
		byArea[area] = append(byArea[area], runs...)
	}
	for _, area := range apiAreas {
		compareGolden(t, filepath.Join("testdata", "golden", "requests-"+area+".json"), byArea[area])
	}
}

func hasGET(r goldenRun) bool {
	for _, req := range r.Requests {
		if req.Method == http.MethodGet {
			return true
		}
	}
	return false
}

type goldenFlag struct {
	Name      string `json:"name"`
	Shorthand string `json:"shorthand,omitempty"`
	Type      string `json:"type"`
	Default   string `json:"default,omitempty"`
	Usage     string `json:"usage,omitempty"`
	Required  bool   `json:"required,omitempty"`
	Hidden    bool   `json:"hidden,omitempty"`
}

type goldenCommand struct {
	Path    string       `json:"path"`
	Aliases []string     `json:"aliases,omitempty"`
	Short   string       `json:"short,omitempty"`
	Long    string       `json:"long,omitempty"`
	Hidden  bool         `json:"hidden,omitempty"`
	Flags   []goldenFlag `json:"flags,omitempty"`
}

func TestGoldenCommandTree(t *testing.T) {
	goldenEnv(t)
	root := cmd.RootCommand()
	root.SetArgs([]string{"version"})
	cmd.Execute()

	var out []goldenCommand
	var walk func(c *cobra.Command)
	walk = func(c *cobra.Command) {
		gc := goldenCommand{Path: c.CommandPath(), Aliases: c.Aliases, Short: c.Short, Long: c.Long, Hidden: c.Hidden}
		c.LocalFlags().VisitAll(func(f *pflag.Flag) {
			gc.Flags = append(gc.Flags, goldenFlag{
				Name: f.Name, Shorthand: f.Shorthand, Type: f.Value.Type(), Default: f.DefValue,
				Usage: f.Usage, Required: isRequired(f), Hidden: f.Hidden,
			})
		})
		out = append(out, gc)
		children := c.Commands()
		sort.Slice(children, func(i, j int) bool { return children[i].Name() < children[j].Name() })
		for _, child := range children {
			walk(child)
		}
	}
	walk(root)
	compareGolden(t, filepath.Join("testdata", "golden", "commands.json"), out)
}

func compareGolden(t *testing.T, path string, v any) {
	t.Helper()
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	if err := enc.Encode(v); err != nil {
		t.Fatal(err)
	}
	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(goldenPath(t, path)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(goldenPath(t, path), buf.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(goldenPath(t, path))
	if err != nil {
		t.Fatalf("%v (run go test -run Golden -update . to create it)", err)
	}
	if !bytes.Equal(want, buf.Bytes()) {
		t.Errorf("%s differs:\n%s", path, firstDiff(string(want), buf.String()))
	}
}

// goldenPath resolves a testdata path against the package directory, since the
// tests change the working directory.
func goldenPath(t *testing.T, rel string) string {
	return filepath.Join(packageDir, rel)
}

var packageDir, _ = os.Getwd()

func firstDiff(want, got string) string {
	w, g := strings.Split(want, "\n"), strings.Split(got, "\n")
	for i := 0; i < len(w) || i < len(g); i++ {
		var a, b string
		if i < len(w) {
			a = w[i]
		}
		if i < len(g) {
			b = g[i]
		}
		if a != b {
			start := max(0, i-3)
			return fmt.Sprintf("line %d\n want: %s\n got:  %s\n context: %s", i+1, a, b, strings.Join(g[start:min(len(g), i+3)], "\n"))
		}
	}
	return "(no line differs)"
}
