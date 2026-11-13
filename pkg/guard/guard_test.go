package guard

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Cloverhound/webex-cli/pkg/catalog"
	"github.com/Cloverhound/webex-cli/pkg/fakeapi"
	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webex"
	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

const hgPath = "/huntGroups"

type fixture struct {
	guard *Guard
	api   *webex.Client
	fake  *fakeapi.Server
	hg    resource.Type
	team  resource.Type
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	fake := fakeapi.New()
	srv := httptest.NewServer(fake)
	t.Cleanup(srv.Close)
	api := &webex.Client{
		API:     &webexapi.Client{Token: func(context.Context) (string, error) { return "tok", nil }},
		Hosts:   fakeapi.Hosts(srv.URL),
		OrgUUID: fakeapi.OrgUUID,
	}
	hg, err := catalog.ByName("hunt_group")
	if err != nil {
		t.Fatal(err)
	}
	team, err := catalog.ByName("cc_team")
	if err != nil {
		t.Fatal(err)
	}
	return &fixture{guard: &Guard{API: api, Writable: true, VerifyDelay: time.Millisecond}, api: api, fake: fake, hg: hg, team: team}
}

var hg1 = resource.Key{Parent: "loc1", ID: "hg1"}

func (f *fixture) prior(t *testing.T) resource.Values {
	t.Helper()
	obj, err := f.api.Read(context.Background(), f.hg, hg1)
	if err != nil {
		t.Fatal(err)
	}
	return obj.Values
}

func with(v resource.Values, field string, value any) resource.Values {
	out := resource.Values{}
	for k, x := range v {
		out[k] = x
	}
	out[field] = value
	return out
}

func failureKind(t *testing.T, err error) Kind {
	t.Helper()
	var f *Failure
	if !errors.As(err, &f) {
		t.Fatalf("err = %v, want a Failure", err)
	}
	return f.Kind
}

// stopDuringRead stands in for another write failing while this one reads.
type stopDuringRead struct {
	API
	guard *Guard
}

func (s stopDuringRead) Read(ctx context.Context, t resource.Type, key resource.Key) (webex.Object, error) {
	s.guard.stop("another write failed")
	return s.API.Read(ctx, t, key)
}

func TestAFailureDuringTheReadStopsTheWrite(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	f.guard.API = stopDuringRead{f.api, f.guard}
	if _, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "no_answer_rings", 7)); failureKind(t, err) != Stopped {
		t.Fatalf("update err = %v", err)
	}
	f.guard.stopped = ""
	if err := f.guard.Delete(context.Background(), f.hg, hg1, prior); failureKind(t, err) != Stopped {
		t.Fatalf("delete err = %v", err)
	}
	if m := f.fake.Mutations(); len(m) != 0 {
		t.Fatalf("sent %v", m)
	}
}

// failingUpdate fails every update and, on the read that observes the failure,
// records whether another write would still be admitted.
type failingUpdate struct {
	API
	guard    *Guard
	reads    int
	admitted *bool
}

func (f *failingUpdate) Read(ctx context.Context, t resource.Type, key resource.Key) (webex.Object, error) {
	f.reads++
	if f.reads == 2 {
		*f.admitted = f.guard.admit(ctx) == nil
	}
	return f.API.Read(ctx, t, key)
}

func (f *failingUpdate) Update(context.Context, resource.Type, resource.Key, map[string]any, resource.Values) error {
	return errors.New("connection reset")
}

func TestAFailedUpdateStopsOtherWritesBeforeItIsObserved(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	var admitted bool
	f.guard.API = &failingUpdate{API: f.api, guard: f.guard, admitted: &admitted}
	if _, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "no_answer_rings", 7)); failureKind(t, err) != OutcomeUnknown {
		t.Fatalf("err = %v", err)
	}
	if admitted {
		t.Fatal("another write was admitted while the failed update was being observed")
	}
}

func TestUpdateWritesOnlyChangedFieldsAndVerifies(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	got, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "no_answer_rings", 7))
	if err != nil {
		t.Fatal(err)
	}
	if got["no_answer_rings"] != 7.0 {
		t.Fatalf("got %v", got)
	}
	if len(f.fake.Mutations()) != 1 {
		t.Fatalf("mutations = %v", f.fake.Mutations())
	}
}

func TestUpdateRefusedWithoutWritePermission(t *testing.T) {
	f := newFixture(t)
	f.guard.Writable = false
	prior := f.prior(t)
	_, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "name", "x"))
	if failureKind(t, err) != Refused {
		t.Fatalf("err = %v", err)
	}
	if len(f.fake.Mutations()) != 0 {
		t.Fatal("a write was sent")
	}
}

func TestConcurrentEditOfManagedFieldStopsWithoutWriting(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	f.fake.External("hg", "hg1", map[string]any{"callPolicies": map[string]any{"noAnswer": map[string]any{"numberOfRings": 11}}})

	got, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "no_answer_rings", 7))
	if failureKind(t, err) != Conflict {
		t.Fatalf("err = %v", err)
	}
	if got["no_answer_rings"] != 11.0 {
		t.Fatalf("returned state must be the real value, got %v", got["no_answer_rings"])
	}
	if len(f.fake.Mutations()) != 0 {
		t.Fatal("the other admin's change was overwritten")
	}
}

func TestChangeToUnmanagedFieldIsNotAConflict(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	f.fake.External("hg", "hg1", map[string]any{"esn": "999"})
	if _, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "name", "x")); err != nil {
		t.Fatal(err)
	}
}

func TestIgnoredWriteIsReportedWithTheRealValue(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	got, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "first_name", ""))
	if failureKind(t, err) != Divergent {
		t.Fatalf("err = %v", err)
	}
	if got["first_name"] != "zz" {
		t.Fatalf("state must keep the real value, got %v", got["first_name"])
	}
}

func TestLostResponseIsUnknownButObservedIntended(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	f.fake.AddFault(fakeapi.Fault{Method: "PUT", Path: hgPath, Mode: fakeapi.LoseResponse})
	got, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "name", "new"))
	if failureKind(t, err) != OutcomeUnknown {
		t.Fatalf("err = %v", err)
	}
	if got["name"] != "new" {
		t.Fatalf("state must show the landed value, got %v", got["name"])
	}
	if puts := f.fake.Mutations(); len(puts) != 1 {
		t.Fatalf("a write must not be retried: %v", puts)
	}
}

func TestServerErrorAfterCommitIsUnknown(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	f.fake.AddFault(fakeapi.Fault{Method: "PUT", Path: hgPath, Mode: fakeapi.Fail500Commit})
	_, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "name", "new"))
	if failureKind(t, err) != OutcomeUnknown {
		t.Fatalf("err = %v", err)
	}
}

func TestClientErrorIsARejection(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	f.fake.AddFault(fakeapi.Fault{Method: "PUT", Path: hgPath, Mode: fakeapi.Fail401})
	_, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "name", "new"))
	if failureKind(t, err) != Rejected {
		t.Fatalf("err = %v", err)
	}
}

func TestWritesStopAfterTheFirstFailure(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	f.fake.AddFault(fakeapi.Fault{Method: "PUT", Path: hgPath, Mode: fakeapi.Fail500NoCommit})
	if _, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "name", "a")); err == nil {
		t.Fatal("expected a failure")
	}
	before := len(f.fake.Mutations())

	if _, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "name", "b")); failureKind(t, err) != Stopped {
		t.Fatalf("update err = %v", err)
	}
	if _, _, err := f.guard.Create(context.Background(), f.team, resource.Values{"name": "n"}); failureKind(t, err) != Stopped {
		t.Fatalf("create err = %v", err)
	}
	if err := f.guard.Delete(context.Background(), f.team, resource.Key{ID: "team1"}, resource.Values{}); failureKind(t, err) != Stopped {
		t.Fatalf("delete err = %v", err)
	}
	if len(f.fake.Mutations()) != before {
		t.Fatal("a write was sent after the stop")
	}
}

func TestCancelledRunSendsNothing(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := f.guard.Update(ctx, f.hg, hg1, prior, with(prior, "name", "x"))
	if failureKind(t, err) != Stopped {
		t.Fatalf("err = %v", err)
	}
	if len(f.fake.Mutations()) != 0 {
		t.Fatal("a write was sent")
	}
}

func TestFailedVerificationReadIsUnverified(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	f.fake.AddFault(fakeapi.Fault{Method: "GET", Path: hgPath, Mode: fakeapi.Fail500NoCommit, After: "PUT", Times: -1})
	_, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "name", "x"))
	if failureKind(t, err) != Unverified {
		t.Fatalf("err = %v", err)
	}
}

func TestCreateRunsSecondStepForIgnoredFields(t *testing.T) {
	f := newFixture(t)
	desired := resource.Values{"location_id": "loc1", "name": "new", "extension": "1777", "enabled": false}
	key, got, err := f.guard.Create(context.Background(), f.hg, desired)
	if err != nil {
		t.Fatal(err)
	}
	if got["enabled"] != false || f.fake.HuntGroup(key.ID)["enabled"] != false {
		t.Fatalf("got %v", got)
	}
	if m := f.fake.Mutations(); len(m) != 2 || m[0].Method != "POST" || m[1].Method != "PUT" {
		t.Fatalf("mutations = %v", m)
	}
}

func TestFailedSecondStepReturnsTheKeySoStateHoldsTheObject(t *testing.T) {
	f := newFixture(t)
	f.fake.AddFault(fakeapi.Fault{Method: "PUT", Path: hgPath, Mode: fakeapi.Fail500NoCommit})
	desired := resource.Values{"location_id": "loc1", "name": "new", "extension": "1777", "enabled": false}
	key, got, err := f.guard.Create(context.Background(), f.hg, desired)
	if err == nil || key.ID == "" {
		t.Fatalf("key=%+v err=%v", key, err)
	}
	if got["enabled"] != true {
		t.Fatalf("must report the group as still enabled, got %v", got["enabled"])
	}
}

func TestLostCreateIsUnknownAndNotRetried(t *testing.T) {
	f := newFixture(t)
	f.fake.AddFault(fakeapi.Fault{Method: "POST", Path: "/team$", Mode: fakeapi.LoseResponse})
	desired := resource.Values{"name": "lost", "site_id": "site1", "team_type": "AGENT", "team_status": "IN_SERVICE", "active": true, "description": "d"}
	key, _, err := f.guard.Create(context.Background(), f.team, desired)
	if failureKind(t, err) != OutcomeUnknown || key.ID != "" {
		t.Fatalf("key=%+v err=%v", key, err)
	}
	if len(f.fake.Mutations()) != 1 {
		t.Fatalf("mutations = %v", f.fake.Mutations())
	}
}

func TestDeleteChecksForConflictsAndVerifies(t *testing.T) {
	f := newFixture(t)
	key := resource.Key{ID: "team1"}
	obj, err := f.api.Read(context.Background(), f.team, key)
	if err != nil {
		t.Fatal(err)
	}
	f.fake.External("team", "team1", map[string]any{"description": "changed"})
	if err := f.guard.Delete(context.Background(), f.team, key, obj.Values); failureKind(t, err) != Conflict {
		t.Fatalf("err = %v", err)
	}
	if f.fake.Team("team1") == nil {
		t.Fatal("team was deleted despite the conflict")
	}

	f = newFixture(t)
	obj, _ = f.api.Read(context.Background(), f.team, key)
	if err := f.guard.Delete(context.Background(), f.team, key, obj.Values); err != nil {
		t.Fatal(err)
	}
	if f.fake.Team("team1") != nil {
		t.Fatal("team still exists")
	}
}

func TestDeleteOfAbsentObjectSucceeds(t *testing.T) {
	f := newFixture(t)
	if err := f.guard.Delete(context.Background(), f.team, resource.Key{ID: "gone"}, resource.Values{}); err != nil {
		t.Fatal(err)
	}
	if len(f.fake.Mutations()) != 0 {
		t.Fatal("no request was needed")
	}
}

func TestUpdateOfDeletedObjectIsAConflict(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	f.fake.Delete("hg", "hg1")
	_, err := f.guard.Update(context.Background(), f.hg, hg1, prior, with(prior, "name", "x"))
	if failureKind(t, err) != Conflict {
		t.Fatalf("err = %v", err)
	}
}

func TestUnknownPlannedValuesAreNotTreatedAsChanges(t *testing.T) {
	f := newFixture(t)
	prior := f.prior(t)
	planned := with(prior, "name", "x")
	planned["time_zone"] = nil
	if _, err := f.guard.Update(context.Background(), f.hg, hg1, prior, planned); err != nil {
		t.Fatal(err)
	}
	if f.fake.HuntGroup("hg1")["timeZone"] != "America/New_York" {
		t.Fatal("an unknown value must not clear the field")
	}
}
