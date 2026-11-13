package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/catalog"
	"github.com/Cloverhound/webex-cli/pkg/changeset"
	"github.com/Cloverhound/webex-cli/pkg/fakeapi"
	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// redirect sends requests for the real Webex hosts to srv.
type redirect struct{ srv *url.URL }

func (r redirect) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme, req.URL.Host = r.srv.Scheme, r.srv.Host
	return http.DefaultTransport.RoundTrip(req)
}

// setupCapture records into a fresh change set file against the fake Webex, which
// serves every catalog type and logs any write that reaches it.
func setupCapture(t *testing.T) (*fakeapi.Server, string) {
	t.Helper()
	resetClientConfig(t)
	fake := fakeapi.New()
	srv := httptest.NewServer(fake)
	t.Cleanup(srv.Close)
	u, _ := url.Parse(srv.URL)
	saved := http.DefaultClient.Transport
	http.DefaultClient.Transport = redirect{u}
	t.Cleanup(func() { http.DefaultClient.Transport = saved })

	path := filepath.Join(t.TempDir(), "t.changeset.json")
	config.SetChangesetPath(path)
	config.SetOrgID(fakeapi.OrgUUID)
	config.SetCCInterval(0)
	t.Cleanup(func() { config.SetChangesetPath(""); config.SetOrgID(""); config.SetCCInterval(fakeCCInterval) })
	return fake, path
}

var fakeCCInterval = config.CCInterval()

// write builds the request a generated command sends for a catalog type's action.
// params holds "parent" and "id".
func write(t *testing.T, typeName string, action catalog.Action, params map[string]string, body string) *Request {
	t.Helper()
	typ, err := catalog.ByName(typeName)
	if err != nil {
		t.Fatal(err)
	}
	r := map[catalog.Action]resource.Route{catalog.Create: typ.Routes.Create, catalog.Update: typ.Routes.Update, catalog.Delete: typ.Routes.Delete}[action]
	if r.Op == nil {
		t.Fatalf("%s has no %s operation", typeName, action)
	}
	base := config.CallingBaseURL()
	if r.Op.Host == webexops.HostCC {
		base = config.CcBaseURL()
	}
	req := NewRequest(base, r.Op.Method, r.Op.Path)
	for _, p := range r.Op.Params {
		if p.In != webexops.InPath {
			continue
		}
		switch p.Name {
		case "orgid":
			req.PathParam(p.Name, fakeapi.OrgUUID)
		case r.Parent:
			req.PathParam(p.Name, params["parent"])
		case r.ID:
			req.PathParam(p.Name, params["id"])
		}
	}
	if body != "" {
		req.SetBodyRaw(body)
	}
	return req
}

func run(t *testing.T, req *Request) map[string]any {
	t.Helper()
	body, status, err := req.Do()
	if err != nil || status != 200 {
		t.Fatalf("status %d, err %v", status, err)
	}
	var out map[string]any
	if err := json.Unmarshal(body, &out); err != nil {
		t.Fatalf("output %s: %v", body, err)
	}
	return out
}

func refused(t *testing.T, req *Request, want string) {
	t.Helper()
	if _, _, err := req.Do(); err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("err = %v, want %q", err, want)
	}
}

func huntGroupUpdate(t *testing.T, body string) *Request {
	return write(t, "hunt_group", catalog.Update, map[string]string{"parent": "loc1", "id": "hg1"}, body)
}

func load(t *testing.T, path string) *changeset.ChangeSet {
	t.Helper()
	cs, err := changeset.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	return cs
}

func TestCaptureRecordsOnlyChangedFieldsAndSendsNothing(t *testing.T) {
	fake, path := setupCapture(t)
	out := run(t, huntGroupUpdate(t, `{"name":"g-hg1","callPolicies":{"noAnswer":{"numberOfRings":7}}}`))
	if out["recorded"] != "op-1" || !reflect.DeepEqual(out["changes"], []any{"callPolicies.noAnswer.numberOfRings"}) {
		t.Errorf("output = %v", out)
	}
	if m := fake.Mutations(); len(m) != 0 {
		t.Fatalf("writes reached Webex: %v", m)
	}
	cs := load(t, path)
	op := cs.Operations[0]
	if cs.Scope.OrgID != fakeapi.OrgUUID || op.Target.Display != "g-hg1" || op.Target.Query != nil || op.Target.Params["huntGroupId"] != "hg1" {
		t.Errorf("operation = %+v", op)
	}
	if c := op.Changes[0]; !c.HasBefore || !changeset.Equal(c.Before, 5) || !changeset.Equal(c.Value, 7) {
		t.Errorf("change = %+v", c)
	}
}

func TestCaptureWorksInReadOnlyMode(t *testing.T) {
	_, path := setupCapture(t)
	config.SetReadOnly(true)
	run(t, huntGroupUpdate(t, `{"name":"Front"}`))
	if cs := load(t, path); len(cs.Operations) != 1 {
		t.Fatalf("operations = %+v", cs.Operations)
	}
}

// Studio refuses a whole file over one operation it cannot apply, so recording
// refuses anything that is not a create, update or delete of a type Studio applies.
func TestCaptureRefusesWhatStudioCannotApply(t *testing.T) {
	fake, path := setupCapture(t)
	answer := NewRequest(config.CallingBaseURL(), "POST", "/telephony/calls/answer")
	answer.SetBodyRaw(`{"callId":"c1"}`)
	person := NewRequest(config.CallingBaseURL(), "POST", "/people")
	person.SetBodyRaw(`{"emails":["a@example.com"]}`)
	prefs := NewRequest(config.CallingBaseURL(), "PUT", "/meetingPreferences/audio")
	prefs.SetBodyRaw(`{"defaultAudioType":"PSTN"}`)
	bulk := NewRequest(config.CcBaseURL(), "POST", "/organization/{orgid}/team/bulk")
	bulk.PathParam("orgid", fakeapi.OrgUUID)
	bulk.SetBodyRaw(`{"items":[]}`)
	flow := NewRequest(config.CcBaseURL(), "PATCH", "/{orgid}/project/{projectId}/v2/flows/{flowId}")
	flow.PathParam("orgid", fakeapi.OrgUUID)
	flow.SetBodyRaw(`{"name":"x"}`)
	for _, req := range []*Request{answer, person, prefs, bulk, flow} {
		refused(t, req, "cannot be recorded")
	}
	refused(t, write(t, "cc_skill", catalog.Delete, map[string]string{"id": "skill1"}, `{"reason":"x"}`), "a delete with a body")
	refused(t, write(t, "cc_team", catalog.Update, map[string]string{"id": "team1"}, ""), "--body")
	if _, err := changeset.Load(path); err == nil {
		t.Error("a refused write started a change set")
	}
	if m := fake.Mutations(); len(m) != 0 {
		t.Fatalf("writes reached Webex: %v", m)
	}
}

// Studio applies a change from its route and fields, so a query that changes what a
// write means would be lost.
func TestCaptureRefusesQueryParameters(t *testing.T) {
	_, path := setupCapture(t)
	queue := write(t, "call_queue", catalog.Create, map[string]string{"parent": "loc1"}, `{"name":"Q","extension":"2100"}`)
	queue.QueryParam("hasCxEssentials", "true")
	refused(t, queue, "query parameter hasCxEssentials")

	// A generated Contact Center update repeats its ID as a query parameter.
	team := write(t, "cc_team", catalog.Update, map[string]string{"id": "team1"}, `{"name":"t-team1","description":"new"}`)
	team.QueryParam("id", "team1")
	run(t, team)
	if cs := load(t, path); cs.Operations[0].Target.Query != nil {
		t.Errorf("target = %+v", cs.Operations[0].Target)
	}
}

func TestCaptureCreateReturnsAPlaceholderLaterCommandsCanUse(t *testing.T) {
	_, path := setupCapture(t)
	out := run(t, write(t, "cc_team", catalog.Create, nil, `{"name":"Boulder Desk","siteId":"site1"}`))
	if out["id"] != "new:boulder-desk" || out["action"] != "create" {
		t.Fatalf("output = %v", out)
	}
	run(t, write(t, "cc_team", catalog.Update, map[string]string{"id": "new:boulder-desk"}, `{"capacity":4}`))
	cs := load(t, path)
	if len(cs.Operations) != 1 {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	team := cs.Operations[0]
	if team.Target.Route != "/team" || team.Handle != "boulder-desk" || !changeset.Equal(team.Body.(map[string]any)["capacity"], 4) {
		t.Errorf("create = %+v", team)
	}
	if p := changeset.Validate(cs); len(p.Errors) > 0 {
		t.Errorf("recorded change set is invalid: %v", p.Errors)
	}
}

// A new: ID must name a create of the type the request writes, under the same parent.
func TestCaptureRefusesAPlaceholderOfAnotherTypeOrParent(t *testing.T) {
	_, path := setupCapture(t)
	run(t, write(t, "cc_team", catalog.Create, nil, `{"name":"Spanish","siteId":"site1"}`))
	refused(t, write(t, "cc_skill", catalog.Update, map[string]string{"id": "new:spanish"}, `{"description":"x"}`), "is not a new Skill")
	refused(t, write(t, "cc_skill", catalog.Delete, map[string]string{"id": "new:spanish"}, ""), "is not a new Skill")

	run(t, write(t, "hunt_group", catalog.Create, map[string]string{"parent": "loc1"}, `{"name":"Desk","extension":"4100"}`))
	refused(t, write(t, "hunt_group", catalog.Update, map[string]string{"parent": "loc2", "id": "new:desk"}, `{"enabled":false}`), "is created under loc1")

	cs := load(t, path)
	if len(cs.Operations) != 2 || cs.Operations[0].Body.(map[string]any)["description"] != nil {
		t.Fatalf("operations = %+v", cs.Operations)
	}
}

func TestCaptureDeleteOfANewObjectWithdrawsItsCreate(t *testing.T) {
	_, path := setupCapture(t)
	run(t, write(t, "cc_skill", catalog.Create, nil, `{"name":"Spanish","skillType":"BOOLEAN"}`))
	run(t, write(t, "cc_skill", catalog.Delete, map[string]string{"id": "new:spanish"}, ""))
	if cs := load(t, path); len(cs.Operations) != 0 {
		t.Errorf("operations = %+v", cs.Operations)
	}
}

// A Contact Center queue is updated and deleted through different routes.
func TestCaptureADeleteReplacesAnUpdateThroughAnotherRoute(t *testing.T) {
	_, path := setupCapture(t)
	queue := map[string]string{"id": "queue1"}
	run(t, write(t, "cc_queue", catalog.Update, queue, `{"name":"q-queue1","description":"new"}`))
	run(t, write(t, "cc_queue", catalog.Delete, queue, ""))
	cs := load(t, path)
	if len(cs.Operations) != 1 || cs.Operations[0].Action != changeset.Delete {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	refused(t, write(t, "cc_queue", catalog.Update, queue, `{"description":"again"}`), "cannot also be updated")
}

// An agent sends back the whole object it read, which carries fields no catalog field
// owns. Studio refuses a change to any of them, so they are ignored, not recorded.
func TestCaptureIgnoresFieldsStudioDoesNotModel(t *testing.T) {
	_, path := setupCapture(t)
	body := `{"id":"team1","name":"t-team1","teamType":"AGENT","teamStatus":"IN_SERVICE","active":true,"siteId":"site1",
		"description":"changed","version":2,"lastUpdatedTime":999,"dbId":"db-team1","futureField":"stale",
		"links":{"self":"https://example.com"},"rankQueuesForTeam":false,"userIds":["user1"]}`
	out := run(t, write(t, "cc_team", catalog.Update, map[string]string{"id": "team1"}, body))
	if !reflect.DeepEqual(out["changes"], []any{"description"}) {
		t.Errorf("changes = %v", out["changes"])
	}
	ignored, _ := out["ignored"].([]any)
	for _, want := range []string{"version", "lastUpdatedTime", "futureField", "links.self"} {
		found := false
		for _, got := range ignored {
			found = found || got == want
		}
		if !found {
			t.Errorf("%s was not ignored: %v", want, ignored)
		}
	}
	if cs := load(t, path); len(cs.Operations[0].Changes) != 1 {
		t.Errorf("changes = %+v", cs.Operations[0].Changes)
	}
}

// Work types have no route to read one, so recording reads them by filter.
func TestCaptureReadsAWorkTypeByFilter(t *testing.T) {
	_, path := setupCapture(t)
	run(t, write(t, "cc_work_type", catalog.Update, map[string]string{"id": "cc_work_type1"}, `{"name":"sample name","description":"changed"}`))
	c := load(t, path).Operations[0].Changes
	if len(c) != 1 || c[0].Path != "description" || c[0].Before != "sample description" {
		t.Fatalf("changes = %+v", c)
	}
}

func TestCaptureRefusesAnObjectItCannotRead(t *testing.T) {
	_, path := setupCapture(t)
	refused(t, write(t, "cc_team", catalog.Update, map[string]string{"id": "nope"}, `{"name":"x"}`), "is not in Webex")
	if _, err := changeset.Load(path); err == nil {
		t.Error("recorded an update to an object that is not there")
	}
}

func TestCaptureRefusesAnotherOrgsFile(t *testing.T) {
	_, path := setupCapture(t)
	if err := changeset.Save(path, changeset.New("other-org")); err != nil {
		t.Fatal(err)
	}
	refused(t, huntGroupUpdate(t, `{"name":"x"}`), "is for org other-org")
}

func TestCaptureRefusesARequestForAnotherOrg(t *testing.T) {
	const other = "99999999-8888-7777-6666-555555555555"
	_, path := setupCapture(t)
	create := write(t, "cc_skill", catalog.Create, nil, `{"name":"Spanish"}`)
	create.PathParam("orgid", other)
	refused(t, create, "is for org "+other)

	update := huntGroupUpdate(t, `{"name":"Front"}`)
	update.QueryParam("orgId", config.EncodeOrgID(other))
	refused(t, update, "is for org "+other)

	body := write(t, "cc_team", catalog.Create, nil, `{"name":"x","organizationId":"`+other+`"}`)
	refused(t, body, "the request body names org "+other)
	run(t, write(t, "cc_team", catalog.Create, nil, `{"name":"y","organizationId":"`+fakeapi.OrgUUID+`"}`))
	if cs := load(t, path); len(cs.Operations) != 1 {
		t.Errorf("operations = %+v", cs.Operations)
	}
}

func TestDryRunWithAChangeSetRecordsNothing(t *testing.T) {
	fake, path := setupCapture(t)
	config.SetDryRun(true)
	huntGroupUpdate(t, `{"name":"Front"}`).Do()
	if _, err := changeset.Load(path); err == nil || len(fake.Mutations()) != 0 {
		t.Errorf("dry run recorded or sent a change")
	}
}

func TestCaptureSettingAFieldBackWithdrawsItsPendingChange(t *testing.T) {
	_, path := setupCapture(t)
	run(t, huntGroupUpdate(t, `{"name":"Front","enabled":false}`))
	out := run(t, huntGroupUpdate(t, `{"name":"g-hg1"}`))
	if !reflect.DeepEqual(out["withdrawn"], []any{"name"}) {
		t.Errorf("output = %v", out)
	}
	cs := load(t, path)
	if len(cs.Operations) != 1 || len(cs.Operations[0].Changes) != 1 || cs.Operations[0].Changes[0].Path != "enabled" {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	run(t, huntGroupUpdate(t, `{"enabled":true}`))
	if cs := load(t, path); len(cs.Operations) != 0 {
		t.Errorf("an update with no changes left was kept: %+v", cs.Operations)
	}
}

// A Contact Center PUT carries the whole object, so fields it repeats from Webex do
// not undo changes recorded earlier.
func TestCaptureWholeObjectPutKeepsEarlierChanges(t *testing.T) {
	_, path := setupCapture(t)
	team := func(body string) *Request {
		return write(t, "cc_team", catalog.Update, map[string]string{"id": "team1"}, body)
	}
	run(t, team(`{"name":"Front","description":"seed"}`))
	run(t, team(`{"name":"t-team1","description":"new"}`))
	if got := paths(load(t, path).Operations[0].Changes); !reflect.DeepEqual(got, []string{"name", "description"}) {
		t.Fatalf("changes = %v", got)
	}
}

func TestCaptureUnchangedUpdateRecordsNothing(t *testing.T) {
	_, path := setupCapture(t)
	if out := run(t, huntGroupUpdate(t, `{"name":"g-hg1"}`)); out["unchanged"] != true {
		t.Errorf("output = %v", out)
	}
	if cs, err := changeset.Load(path); err == nil && len(cs.Operations) != 0 {
		t.Errorf("operations = %+v", cs.Operations)
	}
}

func TestLiveCheckReportsDriftMissingAndUnknownTargets(t *testing.T) {
	fake, path := setupCapture(t)
	run(t, huntGroupUpdate(t, `{"name":"Front"}`))
	run(t, write(t, "cc_team", catalog.Delete, map[string]string{"id": "team2"}, ""))
	cs := load(t, path)
	if r := CheckLive(cs); !r.OK() {
		t.Fatalf("report = %+v", r)
	}

	fake.External("hg", "hg1", map[string]any{"name": "Renamed"})
	fake.Delete("team", "team2")
	cs.Operations = append(cs.Operations, changeset.Operation{ID: "op-9", Action: changeset.Update,
		Target:  changeset.Target{API: changeset.APIWebex, Route: "/meetingPreferences/audio"},
		Changes: []changeset.Change{{Path: "x", Value: 1}}})
	r := CheckLive(cs)
	if len(r.Drift) != 1 || r.Drift[0].Op != "op-1" || !reflect.DeepEqual(r.Missing, []string{"op-2"}) || len(r.Unreadable) != 1 {
		t.Errorf("report = %+v", r)
	}
}

func TestLiveCheckIsNotOKWhenAnObjectCannotBeRead(t *testing.T) {
	if (LiveReport{Unreadable: []string{"op-1: HTTP 500"}}).OK() {
		t.Error("a check that could not read an object passed")
	}
}
