package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/changeset"
)

const testOrg = "11111111-2222-3333-4444-555555555555"

// fakeWebex answers GETs from objects keyed by path and counts every other
// method, which capture must never send.
type fakeWebex struct {
	mu      sync.Mutex
	objects map[string]string
	writes  int
}

func (f *fakeWebex) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if r.Method != http.MethodGet {
		f.writes++
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Objects are keyed by path, and by path and query where the query selects one;
	// the org is the same for every request.
	q := r.URL.Query()
	q.Del("orgId")
	body, ok := f.objects[r.URL.Path+"?"+q.Encode()]
	if !ok {
		body, ok = f.objects[r.URL.Path]
	}
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.WriteString(w, body)
}

// redirect sends requests for the real Webex hosts to srv.
type redirect struct{ srv *url.URL }

func (r redirect) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.URL.Scheme, req.URL.Host = r.srv.Scheme, r.srv.Host
	return http.DefaultTransport.RoundTrip(req)
}

func setupCapture(t *testing.T, objects map[string]string) (*fakeWebex, string) {
	t.Helper()
	resetClientConfig(t)
	fake := &fakeWebex{objects: objects}
	srv := httptest.NewServer(fake)
	t.Cleanup(srv.Close)
	u, _ := url.Parse(srv.URL)
	saved := http.DefaultClient.Transport
	http.DefaultClient.Transport = redirect{u}
	t.Cleanup(func() { http.DefaultClient.Transport = saved })

	path := filepath.Join(t.TempDir(), "t.changeset.json")
	config.SetChangesetPath(path)
	config.SetOrgID(testOrg)
	t.Cleanup(func() { config.SetChangesetPath(""); config.SetOrgID("") })
	return fake, path
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

func huntGroupUpdate(body string) *Request {
	req := NewRequest(config.CallingBaseURL(), "PUT", "/telephony/config/locations/{locationId}/huntGroups/{huntGroupId}")
	req.PathParam("locationId", "loc1")
	req.PathParam("huntGroupId", "hg1")
	req.QueryParam("orgId", config.EncodeOrgID(testOrg))
	req.SetBodyRaw(body)
	return req
}

func TestCaptureRecordsOnlyChangedFieldsAndSendsNothing(t *testing.T) {
	fake, path := setupCapture(t, map[string]string{
		"/v1/telephony/config/locations/loc1/huntGroups/hg1": `{"name":"Desk","enabled":true,"callPolicies":{"noAnswer":{"numberOfRings":3}}}`,
	})

	out := run(t, huntGroupUpdate(`{"name":"Desk","callPolicies":{"noAnswer":{"numberOfRings":5}}}`))
	if out["recorded"] != "op-1" || !reflect.DeepEqual(out["changes"], []any{"callPolicies.noAnswer.numberOfRings"}) {
		t.Errorf("output = %v", out)
	}
	if fake.writes != 0 {
		t.Fatalf("%d writes reached Webex", fake.writes)
	}

	cs, err := changeset.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cs.Scope.OrgID != testOrg || len(cs.Operations) != 1 {
		t.Fatalf("change set = %+v", cs)
	}
	op := cs.Operations[0]
	if op.Target.API != changeset.APIWebex || op.Target.Display != "Desk" || op.Target.Query != nil ||
		op.Target.Params["huntGroupId"] != "hg1" {
		t.Errorf("target = %+v", op.Target)
	}
	c := op.Changes[0]
	if !c.HasBefore || !changeset.Equal(c.Before, 3) || !changeset.Equal(c.Value, 5) {
		t.Errorf("change = %+v", c)
	}
}

func TestCaptureWorksInReadOnlyMode(t *testing.T) {
	_, path := setupCapture(t, map[string]string{
		"/v1/telephony/config/locations/loc1/huntGroups/hg1": `{"name":"Desk"}`,
	})
	config.SetReadOnly(true)

	run(t, huntGroupUpdate(`{"name":"Front"}`))
	if cs, err := changeset.Load(path); err != nil || len(cs.Operations) != 1 {
		t.Fatalf("cs = %+v, err = %v", cs, err)
	}
}

func TestCaptureCreateReturnsAPlaceholderLaterCommandsCanUse(t *testing.T) {
	_, path := setupCapture(t, map[string]string{
		"/organization/" + testOrg + "/entry-point/ep1": `{"name":"Main","teamIds":["t0"]}`,
	})

	create := NewRequest(config.CcBaseURL(), "POST", "/organization/{orgid}/team")
	create.PathParam("orgid", testOrg)
	create.SetBodyRaw(`{"name":"Boulder Desk","siteId":"s1"}`)
	out := run(t, create)
	if out["id"] != "new:boulder-desk" || out["action"] != "create" {
		t.Fatalf("output = %v", out)
	}

	rename := NewRequest(config.CcBaseURL(), "PUT", "/organization/{orgid}/team/{id}")
	rename.PathParam("orgid", testOrg)
	rename.PathParam("id", "new:boulder-desk")
	rename.SetBodyRaw(`{"capacity":4}`)
	run(t, rename)

	ep := NewRequest(config.CcBaseURL(), "PUT", "/organization/{orgid}/entry-point/{id}")
	ep.PathParam("orgid", testOrg)
	ep.PathParam("id", "ep1")
	ep.SetBodyRaw(`{"name":"Main","teamIds":["t0","new:boulder-desk"]}`)
	run(t, ep)

	cs, err := changeset.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cs.Operations) != 2 {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	team := cs.Operations[0]
	body := team.Body.(map[string]any)
	if team.Target.Route != "/team" || team.Target.Params != nil || team.Handle != "boulder-desk" ||
		!changeset.Equal(body["capacity"], 4) {
		t.Errorf("create = %+v", team)
	}
	change := cs.Operations[1].Changes
	if cs.Operations[1].Target.Route != "/entry-point/{id}" || len(change) != 1 ||
		!reflect.DeepEqual(change[0].Value, []any{"t0", changeset.RefTo("boulder-desk")}) {
		t.Errorf("entry point = %+v", cs.Operations[1])
	}
	if p := changeset.Validate(cs); len(p.Errors) > 0 {
		t.Errorf("recorded change set is invalid: %v", p.Errors)
	}
}

func TestCaptureDeleteOfANewObjectWithdrawsItsCreate(t *testing.T) {
	_, path := setupCapture(t, nil)
	create := NewRequest(config.CcBaseURL(), "POST", "/organization/{orgid}/skill")
	create.PathParam("orgid", testOrg)
	create.SetBodyRaw(`{"name":"Spanish"}`)
	run(t, create)

	del := NewRequest(config.CcBaseURL(), "DELETE", "/organization/{orgid}/skill/{id}")
	del.PathParam("orgid", testOrg)
	del.PathParam("id", "new:spanish")
	run(t, del)

	if cs, _ := changeset.Load(path); len(cs.Operations) != 0 {
		t.Errorf("operations = %+v", cs.Operations)
	}
}

func TestCaptureActionPostIsAnInvoke(t *testing.T) {
	_, path := setupCapture(t, nil)
	bulk := NewRequest(config.CcBaseURL(), "POST", "/organization/{orgid}/team/bulk")
	bulk.PathParam("orgid", testOrg)
	bulk.SetBodyRaw(`{"items":[]}`)
	if out := run(t, bulk); out["action"] != "invoke" {
		t.Errorf("output = %v", out)
	}
	cs, _ := changeset.Load(path)
	if op := cs.Operations[0]; op.Method != "POST" || op.Target.Route != "/team/bulk" {
		t.Errorf("operation = %+v", op)
	}
}

// A query parameter can select the object, as userEmail selects whose preferences.
func TestCaptureAndLiveCheckReadTheObjectTheQuerySelects(t *testing.T) {
	objects := map[string]string{
		"/v1/meetingPreferences/audio":                           `{"defaultAudioType":"PSTN"}`,
		"/v1/meetingPreferences/audio?userEmail=b%40example.com": `{"defaultAudioType":"VOIP"}`,
	}
	_, path := setupCapture(t, objects)
	req := NewRequest(config.CallingBaseURL(), "PUT", "/meetingPreferences/audio")
	req.QueryParam("userEmail", "b@example.com")
	req.SetBodyRaw(`{"defaultAudioType":"PSTN"}`)
	if out := run(t, req); out["action"] != "update" {
		t.Fatalf("output = %v", out)
	}
	cs, _ := changeset.Load(path)
	if c := cs.Operations[0].Changes[0]; !changeset.Equal(c.Before, "VOIP") {
		t.Fatalf("change = %+v", c)
	}
	if r := CheckLive(cs); !r.OK() {
		t.Fatalf("live check read another object: %+v", r)
	}
	objects["/v1/meetingPreferences/audio?userEmail=b%40example.com"] = `{"defaultAudioType":"WEBEX"}`
	if r := CheckLive(cs); len(r.Drift) != 1 {
		t.Fatalf("drift on the selected object went unseen: %+v", r)
	}
}

// Call control and other verbs under /telephony/calls answer a POST without making
// an object that can be read back.
func TestCapturePostToARouteWithNoItemsIsAnInvoke(t *testing.T) {
	_, path := setupCapture(t, nil)
	answer := NewRequest(config.CallingBaseURL(), "POST", "/telephony/calls/answer")
	answer.SetBodyRaw(`{"callId":"c1"}`)
	if out := run(t, answer); out["action"] != "invoke" {
		t.Errorf("answer = %v", out)
	}
	hg := NewRequest(config.CallingBaseURL(), "POST", "/telephony/config/locations/{locationId}/huntGroups")
	hg.PathParam("locationId", "loc1")
	hg.SetBodyRaw(`{"name":"Desk"}`)
	if out := run(t, hg); out["action"] != "create" {
		t.Errorf("hunt group = %v", out)
	}
	cs, _ := changeset.Load(path)
	if len(cs.Operations) != 2 {
		t.Fatalf("operations = %+v", cs.Operations)
	}
}

func TestCaptureChecksTheOrgABodyNames(t *testing.T) {
	_, path := setupCapture(t, nil)
	person := func(org string) *Request {
		req := NewRequest(config.CallingBaseURL(), "POST", "/people")
		req.SetBodyRaw(`{"emails":["a@example.com"],"orgId":"` + org + `"}`)
		return req
	}
	if _, _, err := person(config.EncodeOrgID("99999999-2222-3333-4444-555555555555")).Do(); err == nil {
		t.Fatal("a body naming another org was recorded")
	}
	if out := run(t, person(config.EncodeOrgID(testOrg))); out["action"] != "create" {
		t.Fatalf("the change set's own org was refused: %v", out)
	}
	if cs, _ := changeset.Load(path); len(cs.Operations) != 1 {
		t.Fatalf("operations = %+v", cs.Operations)
	}
}

// Some deletes carry a reason in their body, which the delete action has no room for.
func TestCaptureDeleteWithABodyIsAnInvoke(t *testing.T) {
	_, path := setupCapture(t, nil)
	del := NewRequest(config.CallingBaseURL(), "DELETE", "/convergedRecordings/{recordingId}")
	del.PathParam("recordingId", "r1")
	del.SetBodyRaw(`{"reason":"audit","comment":"duplicate"}`)
	if out := run(t, del); out["action"] != "invoke" {
		t.Fatalf("output = %v", out)
	}
	cs, _ := changeset.Load(path)
	if op := cs.Operations[0]; op.Method != "DELETE" || op.Body.(map[string]any)["comment"] != "duplicate" {
		t.Fatalf("operation = %+v", op)
	}
}

func TestCaptureRefusesAnotherOrgsFile(t *testing.T) {
	_, path := setupCapture(t, map[string]string{
		"/v1/telephony/config/locations/loc1/huntGroups/hg1": `{"name":"Desk"}`,
	})
	if err := changeset.Save(path, changeset.New("other-org")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := huntGroupUpdate(`{"name":"x"}`).Do(); err == nil {
		t.Error("recorded into another org's change set")
	}
}

func TestCaptureRefusesARequestForAnotherOrg(t *testing.T) {
	const other = "99999999-8888-7777-6666-555555555555"
	_, path := setupCapture(t, map[string]string{
		"/v1/telephony/config/locations/loc1/huntGroups/hg1": `{"name":"Desk"}`,
	})
	create := NewRequest(config.CcBaseURL(), "POST", "/organization/{orgid}/skill")
	create.PathParam("orgid", other)
	create.SetBodyRaw(`{"name":"Spanish"}`)
	if _, _, err := create.Do(); err == nil {
		t.Error("recorded a Contact Center create for another org")
	}
	update := huntGroupUpdate(`{"name":"Front"}`)
	update.QueryParam("orgId", config.EncodeOrgID(other))
	if _, _, err := update.Do(); err == nil {
		t.Error("recorded a Webex update for another org")
	}
	if cs, err := changeset.Load(path); err == nil && len(cs.Operations) != 0 {
		t.Errorf("operations = %+v", cs.Operations)
	}
}

func TestCaptureRefusesAnotherOrgInAWebexPath(t *testing.T) {
	setupCapture(t, nil)
	req := NewRequest(config.CallingBaseURL(), "PUT", "/contacts/organizations/{orgid}/contacts/{contactId}")
	req.PathParam("orgid", "99999999-8888-7777-6666-555555555555")
	req.PathParam("contactId", "c1")
	req.SetBodyRaw(`{"displayName":"x"}`)
	if _, _, err := req.Do(); err == nil {
		t.Error("recorded a Webex update for another org")
	}
}

func TestCaptureBulkDeleteIsAnInvokeWithItsBody(t *testing.T) {
	_, path := setupCapture(t, nil)
	for _, number := range []string{"+15550100", "+15550101"} {
		del := NewRequest(config.CallingBaseURL(), "DELETE", "/telephony/config/locations/{locationId}/numbers")
		del.PathParam("locationId", "loc1")
		del.SetBodyRaw(`{"phoneNumbers":["` + number + `"]}`)
		if out := run(t, del); out["action"] != "invoke" {
			t.Fatalf("output = %v", out)
		}
	}
	cs, _ := changeset.Load(path)
	if len(cs.Operations) != 2 || cs.Operations[1].Method != "DELETE" ||
		!reflect.DeepEqual(cs.Operations[1].Body, map[string]any{"phoneNumbers": []any{"+15550101"}}) {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	if p := changeset.Validate(cs); len(p.Errors) > 0 {
		t.Errorf("errors = %v", p.Errors)
	}
}

func TestDryRunWithAChangeSetRecordsNothing(t *testing.T) {
	fake, path := setupCapture(t, map[string]string{
		"/v1/telephony/config/locations/loc1/huntGroups/hg1": `{"name":"Desk"}`,
	})
	config.SetDryRun(true)
	t.Cleanup(func() { config.SetDryRun(false) })
	huntGroupUpdate(`{"name":"Front"}`).Do()
	if _, err := changeset.Load(path); err == nil || fake.writes != 0 {
		t.Errorf("dry run recorded or sent a change (writes %d)", fake.writes)
	}
}

func TestCaptureRefusesContactCenterRoutesOutsideTheOrgPath(t *testing.T) {
	setupCapture(t, nil)
	patch := NewRequest(config.CcBaseURL(), "PATCH", "/{orgid}/project/{projectId}/v2/flows/{flowId}")
	patch.PathParam("orgid", testOrg)
	patch.PathParam("projectId", "p1")
	patch.PathParam("flowId", "f1")
	patch.SetBodyRaw(`{"name":"x"}`)
	if _, _, err := patch.Do(); err == nil {
		t.Error("recorded a route the change set format cannot name")
	}
}

func TestCaptureSettingAFieldBackWithdrawsItsPendingChange(t *testing.T) {
	_, path := setupCapture(t, map[string]string{
		"/v1/telephony/config/locations/loc1/huntGroups/hg1": `{"name":"Desk","enabled":true}`,
	})
	run(t, huntGroupUpdate(`{"name":"Front","enabled":false}`))
	out := run(t, huntGroupUpdate(`{"name":"Desk"}`))
	if !reflect.DeepEqual(out["withdrawn"], []any{"name"}) {
		t.Errorf("output = %v", out)
	}
	cs, _ := changeset.Load(path)
	if len(cs.Operations) != 1 || len(cs.Operations[0].Changes) != 1 || cs.Operations[0].Changes[0].Path != "enabled" {
		t.Fatalf("operations = %+v", cs.Operations)
	}

	run(t, huntGroupUpdate(`{"enabled":true}`))
	if cs, _ := changeset.Load(path); len(cs.Operations) != 0 {
		t.Errorf("an update with no changes left was kept: %+v", cs.Operations)
	}
}

// A Contact Center PUT carries the whole object as read, so fields it repeats from
// Webex do not undo changes recorded earlier.
func TestCaptureWholeObjectPutKeepsEarlierChanges(t *testing.T) {
	_, path := setupCapture(t, map[string]string{
		"/organization/" + testOrg + "/team/t1": `{"name":"Desk","capacity":2}`,
	})
	put := func(body string) {
		req := NewRequest(config.CcBaseURL(), "PUT", "/organization/{orgid}/team/{id}")
		req.PathParam("orgid", testOrg)
		req.PathParam("id", "t1")
		req.SetBodyRaw(body)
		run(t, req)
	}
	put(`{"name":"Front","capacity":2}`)
	put(`{"name":"Desk","capacity":4}`)
	cs, _ := changeset.Load(path)
	if got := paths(cs.Operations[0].Changes); !reflect.DeepEqual(got, []string{"name", "capacity"}) {
		t.Fatalf("changes = %v", got)
	}
}

func TestLiveCheckIsNotOKWhenAnObjectCannotBeRead(t *testing.T) {
	if (LiveReport{Unreadable: []string{"op-1: HTTP 500"}}).OK() {
		t.Error("a check that could not read an object passed")
	}
}

func TestCaptureUnchangedUpdateRecordsNothing(t *testing.T) {
	_, path := setupCapture(t, map[string]string{
		"/v1/telephony/config/locations/loc1/huntGroups/hg1": `{"name":"Desk"}`,
	})
	if out := run(t, huntGroupUpdate(`{"name":"Desk"}`)); out["unchanged"] != true {
		t.Errorf("output = %v", out)
	}
	if cs, _ := changeset.Load(path); len(cs.Operations) != 0 {
		t.Errorf("operations = %+v", cs.Operations)
	}
}
