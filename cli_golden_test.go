package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"mime"
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
// change with: make golden

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

// maxRequestsPerRun stops a command that would never stop reading pages.
const maxRequestsPerRun = 20

// goldenBoundary replaces a multipart body's random boundary.
const goldenBoundary = "golden-boundary"

func (r *recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := recordedRequest{Method: req.Method, URL: canonicalURL(req), Headers: map[string]string{}}
	boundary := ""
	if _, params, err := mime.ParseMediaType(req.Header.Get("Content-Type")); err == nil {
		boundary = params["boundary"]
	}
	for k, v := range req.Header {
		rec.Headers[k] = strings.Join(v, ", ")
		if boundary != "" {
			rec.Headers[k] = strings.ReplaceAll(rec.Headers[k], boundary, goldenBoundary)
		}
	}
	if req.Body != nil {
		data, _ := io.ReadAll(req.Body)
		rec.Body = string(data)
		if boundary != "" {
			rec.Body = strings.ReplaceAll(rec.Body, boundary, goldenBoundary)
		}
	}
	r.requests = append(r.requests, rec)
	if len(r.requests) > maxRequestsPerRun {
		return nil, fmt.Errorf("golden: more than %d requests in one run", maxRequestsPerRun)
	}

	header := http.Header{"Content-Type": {"application/json"}}
	// A first page links to a second one on a different path, so a list that
	// applies only the link's query is told apart from one that follows it.
	if !strings.Contains(req.URL.RawQuery, "golden_page=") {
		header.Set("Link", fmt.Sprintf(`<%s://%s/golden-next?golden_page=2>; rel="next"`, req.URL.Scheme, req.URL.Host))
	}
	return &http.Response{
		StatusCode: 200,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(fakeResponse)),
		Request:    req,
	}, nil
}

// fakeResponse answers every request. It holds one item under each item field
// the list commands read, and totals of two, so every list reads a second page.
// It also holds the download links the download commands follow, and the agent
// AI Agent Studio commands look up.
var fakeResponse = func() string {
	body := map[string]any{
		"objects":        []any{map[string]any{"id": "s-agent-id", "bot_access_token": "golden-bot-token"}},
		"meta":           map[string]any{"total_count": 2, "totalPages": 2},
		"pageInfo":       map[string]any{"totalPages": 2},
		"pagination":     map[string]any{"totalPages": 2},
		"totalResources": 2,
		"totalResults":   2,
		"url":            "https://download.golden.example/audio?sig=golden",
		"temporaryDirectDownloadLinks": map[string]any{
			"audioDownloadLink":      "https://download.golden.example/audio?sig=golden",
			"recordingDownloadLink":  "https://download.golden.example/recording?sig=golden",
			"transcriptDownloadLink": "https://download.golden.example/transcript?sig=golden",
		},
	}
	for _, field := range listItemFields {
		body[field] = []any{map[string]any{"id": "golden-1"}}
	}
	data, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return string(data)
}()

// listItemFields are the item fields of the paged lists. A list that reads its
// items from another field gets one page.
var listItemFields = []string{
	"agents", "announcements", "assistants", "autoAttendants", "backgroundImages", "baseStations",
	"callParkExtensions", "callParks", "callPickups", "campaigns", "clusters", "data", "dectNetworks",
	"devices", "dialPatterns", "dialPlans", "directories", "endpoints", "eventCategories",
	"eventThresholds", "features", "groups", "handsets", "hosts", "huntGroups", "items",
	"langsAndRuntimes", "languages", "lineKeyTemplates", "locationPaging", "locations", "members",
	"nodes", "numbers", "operatingModes", "phoneNumbers", "playlists", "queues", "resources",
	"Resources", "result", "routeGroups", "routeIdentities", "routeLists", "schedules",
	"specifications", "supervisors", "tenants", "translationPatterns", "trunks", "trunkTypes",
	"virtualExtensionRanges", "virtualExtensions", "virtualLines", "voicemailGroups", "voices",
	"wrapupReasons",
}

// goldenAudioFile is the upload commands' sample --file, made in each run's
// working directory.
const goldenAudioFile = "golden.wav"

func canonicalURL(req *http.Request) string {
	u := *req.URL
	u.RawQuery = u.Query().Encode()
	return u.String()
}

func goldenEnv(t *testing.T) *recorder {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	for _, env := range []string{"WEBEX_TOKEN", "WEBEX_REFRESH_TOKEN", "WEBEX_USER", "WEBEX_REGION", "WEBEX_CC_REGION", "WEBEX_READ_ONLY", "WEBEX_LOGIN_MODE"} {
		t.Setenv(env, "")
	}
	// The OS keyring is shared with real logins; the file store lives in HOME.
	t.Setenv("WEBEX_TOKEN_STORE", "file")
	t.Setenv("WEBEX_CLI_NONINTERACTIVE", "1")
	work := t.TempDir()
	t.Chdir(work)
	if err := os.WriteFile(filepath.Join(work, goldenAudioFile), []byte("RIFF golden audio"), 0o644); err != nil {
		t.Fatal(err)
	}

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
// depends on the clock, body flags replace the per-field flags, and help would
// print help instead of sending a request.
var skippedFlags = map[string]bool{"orgid": true, "org-id": true, "last": true, "body": true, "body-file": true, "help": true}

func sampleValue(f *pflag.Flag) string {
	switch f.Name {
	case "from":
		return "2026-01-01T00:00:00.000Z"
	case "to":
		return "2026-01-02T00:00:00.000Z"
	case "file":
		return goldenAudioFile
	// The Contact Center search commands parse these before building a request.
	case "filter":
		return `{"golden":true}`
	case "aggregate":
		return "s-field:count:s-name"
	case "interval":
		return "DAILY"
	case "time-comparator":
		return "createdTime"
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

// TestGoldenFilesExerciseEachKindOfRequest checks that the recorded runs reach
// the requests they are meant to pin, rather than stopping at a fixture error.
func TestGoldenFilesExerciseEachKindOfRequest(t *testing.T) {
	var runs []goldenRun
	for _, area := range apiAreas {
		data, err := os.ReadFile(filepath.Join(packageDir, "testdata", "golden", "requests-"+area+".json"))
		if err != nil {
			t.Fatal(err)
		}
		var areaRuns []goldenRun
		if err := json.Unmarshal(data, &areaRuns); err != nil {
			t.Fatal(err)
		}
		runs = append(runs, areaRuns...)
	}

	var uploads, downloads, searches, secondPages, linkQueryOnly int
	for _, r := range runs {
		hasArg := func(prefix string) bool {
			for _, a := range r.Args {
				if strings.HasPrefix(a, prefix) {
					return true
				}
			}
			return false
		}
		switch {
		case r.Variant == "all-flags" && hasArg("--file="):
			uploads++
			if !sentMultipart(r) {
				t.Errorf("%s: upload sent no multipart request (error %q)", r.Command, r.Error)
			}
		case r.Variant == "all-flags" && hasArg("--output="):
			downloads++
			last := r.Requests[len(r.Requests)-1]
			if !strings.HasPrefix(last.URL, "https://download.golden.example/") {
				t.Errorf("%s: download did not fetch the file (error %q)", r.Command, r.Error)
			} else if last.Headers["Authorization"] != "" {
				t.Errorf("%s: download sent the Webex token to the file's URL", r.Command)
			}
		case r.Variant == "all-flags" && strings.HasPrefix(r.Command, "cc search ") && hasArg("--filter="):
			searches++
			if len(r.Requests) == 0 || r.Requests[0].Method != http.MethodPost || r.Error != "" {
				t.Errorf("%s: search sent no query (error %q)", r.Command, r.Error)
			}
		case r.Variant == "paginate" && len(r.Requests) > 1:
			secondPages++
			if r.Command == "admin resource-memberships list-group-v2" {
				linkQueryOnly++
				if strings.Contains(r.Requests[1].URL, "/golden-next") {
					t.Errorf("%s: its next page should keep the request path", r.Command)
				}
			}
		}
	}
	if uploads == 0 || downloads == 0 || searches == 0 || secondPages == 0 || linkQueryOnly == 0 {
		t.Errorf("uploads=%d downloads=%d searches=%d second pages=%d link-query-only=%d; each must be recorded",
			uploads, downloads, searches, secondPages, linkQueryOnly)
	}
}

func sentMultipart(r goldenRun) bool {
	for _, req := range r.Requests {
		if strings.HasPrefix(req.Headers["Content-Type"], "multipart/form-data") && strings.Contains(req.Body, goldenBoundary) {
			return true
		}
	}
	return false
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
	Use     string       `json:"use"`
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
		// Cobra adds --help when a command first runs; adding it here keeps the
		// tree the same whichever commands ran before.
		c.InitDefaultHelpFlag()
		gc := goldenCommand{Path: c.CommandPath(), Use: c.Use, Aliases: c.Aliases, Short: c.Short, Long: c.Long, Hidden: c.Hidden}
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
	full := filepath.Join(packageDir, path)
	if *updateGolden {
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, buf.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(full)
	if err != nil {
		t.Fatalf("%v (run make golden to create it)", err)
	}
	if !bytes.Equal(want, buf.Bytes()) {
		t.Errorf("%s differs; if the change is intended, run make golden and review the diff:\n%s", path, firstDiff(string(want), buf.String()))
	}
}

// packageDir is read before the tests change the working directory.
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
