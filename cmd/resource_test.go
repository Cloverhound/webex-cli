package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/changeset"
	"github.com/Cloverhound/webex-cli/pkg/fakeapi"
	"github.com/Cloverhound/webex-cli/pkg/webexapi"
	"github.com/spf13/pflag"
)

// toFake sends every request to the fake Webex, whichever API host it names; the
// fake serves Webex routes under /v1 and Contact Center routes at its root.
type toFake struct{ target *url.URL }

func (f toFake) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.URL.Scheme, r.URL.Host, r.Host = f.target.Scheme, f.target.Host, f.target.Host
	return http.DefaultTransport.(*http.Transport).RoundTrip(r)
}

func fakeWebex(t *testing.T) *fakeapi.Server {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	for _, env := range []string{"WEBEX_TOKEN", "WEBEX_REFRESH_TOKEN", "WEBEX_USER", "WEBEX_REGION", "WEBEX_CC_REGION", "WEBEX_CHANGESET", readOnlyEnv} {
		t.Setenv(env, "")
	}
	fake := fakeapi.New()
	srv := httptest.NewServer(fake)
	t.Cleanup(srv.Close)
	target, _ := url.Parse(srv.URL)
	saved := http.DefaultClient.Transport
	http.DefaultClient.Transport = toFake{target}
	t.Cleanup(func() { http.DefaultClient.Transport = saved })
	config.SetCCInterval(0)
	t.Cleanup(func() { config.SetCCInterval(webexapi.CCInterval) })
	return fake
}

func resetCommandFlags() {
	reset := func(f *pflag.Flag) { f.Value.Set(f.DefValue); f.Changed = false }
	rootCmd.PersistentFlags().VisitAll(reset)
	for _, c := range resourceCmd.Commands() {
		c.Flags().VisitAll(reset)
	}
}

func runResource(t *testing.T, args ...string) (string, error) {
	t.Helper()
	resetCommandFlags()
	var err error
	out := captureStdout(t, func() {
		rootCmd.SetArgs(append(append([]string{"resource"}, args...), "--token", "tok", "--organization", fakeapi.OrgUUID))
		err = rootCmd.Execute()
	})
	return out, err
}

func decode(t *testing.T, out string, v any) {
	t.Helper()
	if err := json.Unmarshal([]byte(out), v); err != nil {
		t.Fatalf("output is not JSON: %v\n%s", err, out)
	}
}

func TestResourceListAndGetReadThroughTheCatalog(t *testing.T) {
	fakeWebex(t)
	out, err := runResource(t, "list", "cc_team")
	if err != nil {
		t.Fatal(err)
	}
	var teams []objectView
	decode(t, out, &teams)
	if len(teams) != 2 || teams[0].Values["name"] != "t-team1" {
		t.Fatalf("teams = %+v", teams)
	}

	out, err = runResource(t, "get", "hunt_group", "hg2")
	if err != nil {
		t.Fatal(err)
	}
	var hg objectView
	decode(t, out, &hg)
	if hg.Key.Parent != "loc1" || hg.Key.ID != "hg2" || hg.Label != "g-hg2" || hg.Raw != nil {
		t.Fatalf("hunt group = %+v", hg)
	}
}

func TestResourceSetKeepsFieldsTheModelDoesNotCover(t *testing.T) {
	fakeWebex(t)
	out, err := runResource(t, "set", "cc_team", "team1", "name=zz-renamed", "capacity=4")
	if err != nil {
		t.Fatal(err)
	}
	var res writeResult
	decode(t, out, &res)
	if strings.Join(res.Changed, ",") != "capacity,name" || res.Values["name"] != "zz-renamed" {
		t.Fatalf("result = %+v", res)
	}

	// A CC PUT replaces the object, so fields outside the model survive only
	// because set sent the whole object back.
	out, err = runResource(t, "get", "cc_team", "team1", "--raw")
	if err != nil {
		t.Fatal(err)
	}
	var team objectView
	decode(t, out, &team)
	if team.Raw["futureField"] != "keep-me" || team.Raw["name"] != "zz-renamed" {
		t.Fatalf("raw team = %+v", team.Raw)
	}
}

func TestResourceSetTakesValuesAsJSON(t *testing.T) {
	fakeWebex(t)
	out, err := runResource(t, "set", "cc_team", "team1", "--values", `{"name":"Front Desk","capacity":4,"active":false}`)
	if err != nil {
		t.Fatal(err)
	}
	var res writeResult
	decode(t, out, &res)
	if res.Values["name"] != "Front Desk" || res.Values["capacity"] != 4.0 || res.Values["active"] != false {
		t.Fatalf("result = %+v", res)
	}
	if _, err := runResource(t, "set", "cc_team", "team1", "--values", `{"capacity":"four"}`); err == nil {
		t.Fatal("text accepted for a number")
	}
	if _, err := runResource(t, "set", "cc_team", "team1"); err == nil || !strings.Contains(err.Error(), "at least one field") {
		t.Fatalf("no fields: %v", err)
	}
}

func TestResourceSetFindsTheParentOfABareID(t *testing.T) {
	fakeWebex(t)
	out, err := runResource(t, "set", "hunt_group", "hg1", "no_answer_rings=7")
	if err != nil {
		t.Fatal(err)
	}
	var res writeResult
	decode(t, out, &res)
	if res.Key.Parent != "loc1" || res.Values["no_answer_rings"] != 7.0 {
		t.Fatalf("result = %+v", res)
	}
}

func TestResourceSetRefusesWhatItCannotDoSafely(t *testing.T) {
	fakeWebex(t)
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"set", "cc_team", "team1", "nickname=x"}, `has no field "nickname"`},
		{[]string{"set", "cc_skill", "skill1", "skill_type=TEXT"}, "cannot be changed after the object is created"},
		{[]string{"set", "cc_team", "team1", "team_type=BIG"}, `"BIG" is not one of AGENT, CAPACITY`},
		{[]string{"set", "cc_team", "team1", "capacity=many"}, "want a whole number"},
		{[]string{"set", "location", "loc1", "name=x"}, "Locations are read-only"},
		{[]string{"create", "cc_user", "email=x@example.com"}, "Users cannot be created"},
		{[]string{"create", "cc_skill", "name=zz"}, "Skill needs skill_type"},
	} {
		if _, err := runResource(t, tc.args...); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%v: err = %v, want %q", tc.args, err, tc.want)
		}
	}
}

func TestResourceCreateThenDelete(t *testing.T) {
	fakeWebex(t)
	out, err := runResource(t, "create", "cc_skill", "name=zz-cli-skill", "skill_type=BOOLEAN")
	if err != nil {
		t.Fatal(err)
	}
	var res writeResult
	decode(t, out, &res)
	if res.Key.ID == "" || res.Values["name"] != "zz-cli-skill" || res.Values["active"] != true {
		t.Fatalf("result = %+v", res)
	}
	if _, err := runResource(t, "delete", "cc_skill", res.Key.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := runResource(t, "get", "cc_skill", res.Key.ID); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("after delete: %v", err)
	}
}

func TestResourceSetIsRecordedInAChangeSet(t *testing.T) {
	fakeWebex(t)
	path := filepath.Join(t.TempDir(), "proposal.changeset.json")
	if _, err := runResource(t, "set", "cc_team", "team1", "name=zz-proposed", "--changeset", path); err != nil {
		t.Fatal(err)
	}
	cs, err := changeset.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cs.Operations) != 1 || cs.Operations[0].Action != changeset.Update {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	changes := cs.Operations[0].Changes
	if len(changes) != 1 || changes[0].Path != "name" || changes[0].Value != "zz-proposed" {
		t.Fatalf("changes = %+v", changes)
	}

	// Nothing was sent.
	out, err := runResource(t, "get", "cc_team", "team1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"t-team1"`) {
		t.Fatalf("team changed: %s", out)
	}
}

func TestResourceSetOfARecordedCreateFoldsIntoIt(t *testing.T) {
	fakeWebex(t)
	path := filepath.Join(t.TempDir(), "proposal.changeset.json")
	out, err := runResource(t, "create", "cc_skill", "name=zz-review", "skill_type=BOOLEAN", "--changeset", path)
	if err != nil {
		t.Fatal(err)
	}
	var created struct{ ID string }
	decode(t, out, &created)
	if created.ID != "new:zz-review" {
		t.Fatalf("create output = %s", out)
	}
	if _, err := runResource(t, "set", "cc_skill", created.ID, "name=zz-revised", "--changeset", path); err != nil {
		t.Fatal(err)
	}
	cs, err := changeset.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cs.Operations) != 1 || cs.Operations[0].Body.(map[string]any)["name"] != "zz-revised" {
		t.Fatalf("operations = %+v", cs.Operations)
	}
}

func TestResourceSetBackToTheLiveValueWithdrawsOnlyThatChange(t *testing.T) {
	fakeWebex(t)
	path := filepath.Join(t.TempDir(), "proposal.changeset.json")
	for _, args := range [][]string{{"name=zz-proposed"}, {"capacity=4"}, {"name=t-team1"}} {
		if _, err := runResource(t, append(append([]string{"set", "cc_team", "team1"}, args...), "--changeset", path)...); err != nil {
			t.Fatal(err)
		}
	}
	cs, err := changeset.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cs.Operations) != 1 {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	if changes := cs.Operations[0].Changes; len(changes) != 1 || changes[0].Path != "capacity" {
		t.Fatalf("changes = %+v", changes)
	}
}

func TestResourceCreateRecordsTheBodyTheSDKWouldSend(t *testing.T) {
	fake := fakeWebex(t)
	for _, o := range fake.Objects("cc_desktop_layout") {
		o["global"] = false
		fake.Put("cc_desktop_layout", o)
	}
	fake.Put("cc_desktop_layout", map[string]any{"id": "global1", "name": "Global", "global": true, "jsonFileContent": `{"layout":1}`})
	path := filepath.Join(t.TempDir(), "proposal.changeset.json")
	if _, err := runResource(t, "create", "cc_desktop_layout", "name=zz-layout", "--changeset", path); err != nil {
		t.Fatal(err)
	}
	cs, err := changeset.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if body := cs.Operations[0].Body.(map[string]any); body["jsonFileContent"] != `{"layout":1}` {
		t.Fatalf("recorded body = %v", body)
	}
}

// Read-only mode refuses --token, so this checks the guard the write commands call
// and the transport underneath it.
func TestResourceWritesAreRefusedInReadOnlyMode(t *testing.T) {
	fakeWebex(t)
	config.SetReadOnly(true)
	t.Cleanup(func() { config.SetReadOnly(false) })
	if err := refuseReadOnly(); err == nil {
		t.Fatal("read-only mode allowed a write")
	}
	req, _ := http.NewRequest(http.MethodPut, "https://api.wxcc-us1.cisco.com/organization/o/team/t", nil)
	if _, err := client.SDK().Send(req.Context(), webexapi.Request{Method: req.Method, URL: req.URL.String()}); err == nil ||
		!strings.Contains(err.Error(), "read-only mode") {
		t.Fatalf("SDK transport sent a write in read-only mode: %v", err)
	}
}
