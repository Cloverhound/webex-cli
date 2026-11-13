package changeset

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func doc(t *testing.T, s string) map[string]any {
	t.Helper()
	var v map[string]any
	if err := decode([]byte(s), &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func paths(changes []Change) []string {
	var out []string
	for _, c := range changes {
		out = append(out, c.Path)
	}
	return out
}

func huntGroup(id string) Target {
	return Target{API: APIWebex, Route: "/telephony/config/locations/{locationId}/huntGroups/{huntGroupId}",
		Params: map[string]string{"locationId": "loc1", "huntGroupId": id}}
}

func TestDiffReportsOnlyChangedPathsWithBefore(t *testing.T) {
	current := doc(t, `{"name":"Desk","enabled":true,"callPolicies":{"noAnswer":{"numberOfRings":3,"nextAgentEnabled":true}},"agents":[{"id":"a"}]}`)
	desired := doc(t, `{"name":"Desk","enabled":true,"callPolicies":{"noAnswer":{"numberOfRings":5.0,"nextAgentEnabled":true}},"agents":[{"id":"a"},{"id":"b"}],"extension":"4100"}`)

	got := Diff(current, desired)
	if want := []string{"agents", "callPolicies.noAnswer.numberOfRings", "extension"}; !reflect.DeepEqual(paths(got), want) {
		t.Fatalf("paths = %v, want %v", paths(got), want)
	}
	rings := got[1]
	if !rings.HasBefore || !Equal(rings.Before, json.Number("3")) || !Equal(rings.Value, json.Number("5")) {
		t.Errorf("rings change = %+v", rings)
	}
	if ext := got[2]; !ext.HasBefore || ext.Before != nil {
		t.Errorf("a field the object lacks should have a null before, got %+v", ext)
	}
}

func TestDiffIgnoresNumberSpelling(t *testing.T) {
	if changes := Diff(doc(t, `{"n":20}`), doc(t, `{"n":2e1}`)); len(changes) != 0 {
		t.Errorf("changes = %+v", changes)
	}
}

func TestLeavesFlattensObjectsWithoutPreconditions(t *testing.T) {
	got := Leaves(doc(t, `{"a":{"b":1,"c":{"d":"x"}},"e":[1,2]}`))
	if want := []string{"a.b", "a.c.d", "e"}; !reflect.DeepEqual(paths(got), want) {
		t.Fatalf("paths = %v", paths(got))
	}
	for _, c := range got {
		if c.HasBefore {
			t.Errorf("%s has a precondition", c.Path)
		}
	}
}

func TestGetAndSet(t *testing.T) {
	d := map[string]any{}
	Set(d, "a.b.c", 1)
	Set(d, "a.x", "y")
	if v, ok := Get(d, "a.b.c"); !ok || v != 1 {
		t.Errorf("Get a.b.c = %v %v", v, ok)
	}
	if _, ok := Get(d, "a.b.missing"); ok {
		t.Error("Get found a missing path")
	}
}

func TestRepeatUpdatesMergeKeepingTheFirstBefore(t *testing.T) {
	cs := New("org")
	first, _ := cs.AddUpdate(huntGroup("hg1"), []Change{{Path: "name", Value: "B", Before: "A", HasBefore: true}})
	second, _ := cs.AddUpdate(huntGroup("hg1"), []Change{{Path: "name", Value: "C", Before: "B", HasBefore: true}, {Path: "enabled", Value: false}})
	if first != second || len(cs.Operations) != 1 {
		t.Fatalf("expected one operation, got %+v", cs.Operations)
	}
	changes := cs.Operations[0].Changes
	if changes[0].Value != "C" || changes[0].Before != "A" || len(changes) != 2 {
		t.Errorf("changes = %+v", changes)
	}
}

func TestPlaceholderNamesDoNotSplitTargets(t *testing.T) {
	cs := New("org")
	a, _ := cs.AddUpdate(Target{API: APIContactCenter, Route: "/team/{id}", Params: map[string]string{"id": "t1"}}, []Change{{Path: "name", Value: "x"}})
	b, _ := cs.AddUpdate(Target{API: APIContactCenter, Route: "/team/{teamId}", Params: map[string]string{"teamId": "t1"}}, []Change{{Path: "capacity", Value: 2}})
	if a != b {
		t.Errorf("same object recorded twice: %s, %s", a, b)
	}
}

func TestTargetsThatDifferOnlyInTheirQueryAreDifferentObjects(t *testing.T) {
	cs := New("org")
	audio := func(user string) Target {
		return Target{API: APIWebex, Route: "/meetingPreferences/audio", Query: map[string]string{"userEmail": user}}
	}
	a, _ := cs.AddUpdate(audio("a@example.com"), []Change{{Path: "defaultAudioType", Value: "VOIP"}})
	b, _ := cs.AddUpdate(audio("b@example.com"), []Change{{Path: "defaultAudioType", Value: "PSTN"}})
	if a == b || len(cs.Operations) != 2 || cs.Operations[1].Target.Query["userEmail"] != "b@example.com" {
		t.Fatalf("operations = %+v", cs.Operations)
	}
}

func TestSaveKeepsChangeSetsOwnerOnly(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.changeset.json")
	if err := Save(path, New("org")); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o600 {
		t.Fatalf("new file mode = %v", info.Mode().Perm())
	}
	if err := os.Chmod(path, 0o640); err != nil {
		t.Fatal(err)
	}
	if err := Save(path, New("org")); err != nil {
		t.Fatal(err)
	}
	if info, _ := os.Stat(path); info.Mode().Perm() != 0o640 {
		t.Fatalf("an existing file's mode became %v", info.Mode().Perm())
	}
}

func TestDeleteReplacesUpdate(t *testing.T) {
	cs := New("org")
	cs.AddUpdate(huntGroup("hg1"), []Change{{Path: "name", Value: "x"}})
	cs.AddDelete(huntGroup("hg1"))
	if len(cs.Operations) != 1 || cs.Operations[0].Action != Delete || cs.Operations[0].Changes != nil {
		t.Errorf("operations = %+v", cs.Operations)
	}
}

func TestAnUpdateToADeletedObjectIsRefused(t *testing.T) {
	cs := New("org")
	cs.AddDelete(huntGroup("hg1"))
	if _, err := cs.AddUpdate(huntGroup("hg1"), []Change{{Path: "name", Value: "x"}}); err == nil {
		t.Error("recorded an update to an object the change set deletes")
	}
	if len(cs.Operations) != 1 {
		t.Errorf("operations = %+v", cs.Operations)
	}
}

func TestAnUpdateMergesPastAnInvokeOnTheSameObject(t *testing.T) {
	cs := New("org")
	cs.Operations = append(cs.Operations, Operation{ID: "op-1", Action: Invoke, Method: "PUT", Target: huntGroup("hg1"), Body: "raw"})
	first, _ := cs.AddUpdate(huntGroup("hg1"), []Change{{Path: "name", Value: "x"}})
	second, _ := cs.AddUpdate(huntGroup("hg1"), []Change{{Path: "enabled", Value: true}})
	if first != second || len(cs.Operations) != 2 {
		t.Errorf("operations = %+v", cs.Operations)
	}
}

// A queue is updated and deleted through different routes; an identity that knows
// that makes the delete replace the pending update.
func TestAnIdentityJoinsTwoRoutesToOneObject(t *testing.T) {
	cs := New("org")
	cs.SetIdentity(func(t Target) string { return t.Params["id"] })
	update := Target{API: APIContactCenter, Route: "/v2/contact-service-queue/{id}", Params: map[string]string{"id": "q1"}}
	del := Target{API: APIContactCenter, Route: "/contact-service-queue/{id}", Params: map[string]string{"id": "q1"}}
	cs.AddUpdate(update, []Change{{Path: "name", Value: "x"}})
	cs.AddDelete(del)
	if len(cs.Operations) != 1 || cs.Operations[0].Action != Delete {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	if _, err := cs.AddUpdate(update, []Change{{Path: "name", Value: "y"}}); err == nil {
		t.Error("an update to a deleted queue was recorded")
	}
}

func TestCreatesGetHandlesThatReferencesResolveTo(t *testing.T) {
	cs := New("org")
	team := Target{API: APIContactCenter, Route: "/team"}
	_, h1 := cs.AddCreate(team, map[string]any{"name": "Boulder Desk"})
	_, h2 := cs.AddCreate(team, map[string]any{"name": "Boulder desk!"})
	_, h3 := cs.AddCreate(Target{API: APIContactCenter, Route: "/skill"}, map[string]any{})
	if h1 != "boulder-desk" || h2 != "boulder-desk-2" || h3 != "new-skill" {
		t.Fatalf("handles = %q %q %q", h1, h2, h3)
	}

	resolved, err := cs.ResolvePlaceholders(map[string]any{"teamIds": []any{"t0", "new:boulder-desk"}})
	if err != nil {
		t.Fatal(err)
	}
	ids := resolved.(map[string]any)["teamIds"].([]any)
	if ids[0] != "t0" || !reflect.DeepEqual(ids[1], RefTo("boulder-desk")) {
		t.Errorf("resolved = %v", resolved)
	}
	if _, err := cs.ResolvePlaceholders("new:nope"); err == nil {
		t.Error("an unknown handle resolved")
	}
}

func TestUpdatesToANewObjectFoldIntoItsCreate(t *testing.T) {
	cs := New("org")
	_, h := cs.AddCreate(Target{API: APIContactCenter, Route: "/team"}, map[string]any{"name": "T"})
	if _, err := cs.UpdateCreate(h, []Change{{Path: "capacity", Value: 3}, {Path: "a.b", Value: true}}); err != nil {
		t.Fatal(err)
	}
	body := cs.Operations[0].Body.(map[string]any)
	if body["capacity"] != 3 || body["a"].(map[string]any)["b"] != true {
		t.Errorf("body = %v", body)
	}
}

func TestAnUpdateThatComesToReferALaterCreateMovesAfterIt(t *testing.T) {
	cs := New("org")
	ep := Target{API: APIContactCenter, Route: "/entry-point/{id}", Params: map[string]string{"id": "e"}}
	cs.AddUpdate(ep, []Change{{Path: "name", Value: "Main"}})
	_, h := cs.AddCreate(Target{API: APIContactCenter, Route: "/team"}, map[string]any{"name": "T"})
	cs.AddUpdate(Target{API: APIContactCenter, Route: "/skill/{id}", Params: map[string]string{"id": "s"}},
		[]Change{{Path: "name", Value: "S"}})
	id, _ := cs.AddUpdate(ep, []Change{{Path: "teamIds", Value: []any{RefTo(h)}}})

	if order := []string{cs.Operations[0].ID, cs.Operations[1].ID, cs.Operations[2].ID}; id != "op-1" ||
		!reflect.DeepEqual(order, []string{"op-2", "op-1", "op-3"}) {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	if p := Validate(cs); len(p.Errors) > 0 {
		t.Errorf("errors = %v", p.Errors)
	}
}

func TestACreateThatComesToReferALaterCreateMovesAfterIt(t *testing.T) {
	cs := New("org")
	team := Target{API: APIContactCenter, Route: "/team"}
	_, a := cs.AddCreate(team, map[string]any{"name": "A"})
	_, b := cs.AddCreate(team, map[string]any{"name": "B"})
	if _, err := cs.UpdateCreate(a, []Change{{Path: "backupTeamId", Value: RefTo(b)}}); err != nil {
		t.Fatal(err)
	}
	if cs.Operations[0].Handle != b || cs.Operations[1].Handle != a {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	if p := Validate(cs); len(p.Errors) > 0 {
		t.Errorf("errors = %v", p.Errors)
	}

	// c cannot move after d while an operation between them still refers to c.
	_, c := cs.AddCreate(team, map[string]any{"name": "C"})
	cs.AddUpdate(Target{API: APIContactCenter, Route: "/entry-point/{id}", Params: map[string]string{"id": "e"}},
		[]Change{{Path: "teamIds", Value: []any{RefTo(c)}}})
	_, d := cs.AddCreate(team, map[string]any{"name": "D"})
	if _, err := cs.UpdateCreate(c, []Change{{Path: "backupTeamId", Value: RefTo(d)}}); err == nil {
		t.Error("recorded a reference no order of operations can satisfy")
	}
	if p := Validate(cs); len(p.Errors) > 0 {
		t.Errorf("a refused change was left in the change set: %v", p.Errors)
	}
}

func TestWithdrawRemovesChangesAndEmptiedUpdates(t *testing.T) {
	cs := New("org")
	cs.AddUpdate(huntGroup("hg1"), []Change{{Path: "name", Value: "Front"}, {Path: "enabled", Value: false}})
	if got := cs.Withdraw(huntGroup("hg1"), []string{"name", "extension"}); !reflect.DeepEqual(got, []string{"name"}) {
		t.Errorf("withdrawn = %v", got)
	}
	if got := paths(cs.Operations[0].Changes); !reflect.DeepEqual(got, []string{"enabled"}) {
		t.Errorf("changes = %v", got)
	}
	cs.Withdraw(huntGroup("hg1"), []string{"enabled"})
	if len(cs.Operations) != 0 {
		t.Errorf("operations = %+v", cs.Operations)
	}
	if got := cs.Withdraw(huntGroup("hg2"), []string{"name"}); got != nil {
		t.Errorf("withdrew from an object with no update: %v", got)
	}
}

func TestRemoveCreateRefusesWhileReferenced(t *testing.T) {
	cs := New("org")
	_, h := cs.AddCreate(Target{API: APIContactCenter, Route: "/team"}, map[string]any{"name": "T"})
	cs.AddUpdate(Target{API: APIContactCenter, Route: "/entry-point/{id}", Params: map[string]string{"id": "e"}},
		[]Change{{Path: "teamIds", Value: []any{RefTo(h)}}})
	if _, err := cs.RemoveCreate(h); err == nil {
		t.Fatal("removed a create still referred to")
	}
	cs.Operations = cs.Operations[:1]
	if _, err := cs.RemoveCreate(h); err != nil || len(cs.Operations) != 0 {
		t.Errorf("err = %v, operations = %+v", err, cs.Operations)
	}
}

func TestValidate(t *testing.T) {
	cs := New("org")
	cs.Intents = []Intent{{ID: "i", Summary: "s"}}
	_, h := cs.AddCreate(Target{API: APIContactCenter, Route: "/team"}, map[string]any{"name": "T"})
	cs.AddUpdate(huntGroup("hg1"), []Change{{Path: "name", Value: RefTo(h)}})
	cs.Operations[1].Intent = "i"
	if p := Validate(cs); len(p.Errors) > 0 {
		t.Fatalf("valid change set: %v", p.Errors)
	}

	bad := *cs
	bad.Operations = []Operation{
		{ID: "op-1", Action: Update, Target: huntGroup("hg1"), Changes: []Change{{Path: "a..b", Value: 1}, {Path: "x", Value: RefTo("later")}}, Intent: "nope"},
		{ID: "op-1", Action: Delete, Target: Target{API: "calling", Route: "/team"}},
		{ID: "op-3", Action: Create, Handle: "later", Target: Target{API: APIContactCenter, Route: "/team"}, Body: map[string]any{}},
		{ID: "op-4", Action: Invoke, Method: "POST", Target: Target{API: APIContactCenter, Route: "/team/bulk"}},
		{ID: "op-5", Action: Update, Target: huntGroup("new:later"), Changes: []Change{{Path: "a", Value: 1}}},
		{ID: "op-6", Action: Delete, Target: huntGroup("hg1")},
	}
	p := Validate(&bad)
	for _, want := range []string{
		`bad change path "a..b"`, `$ref "later" does not name an earlier create`, `intent "nope" is not defined`,
		"op-1: id appears twice", `target.api is "calling"`, "route of a delete must end with the object's ID",
		"is a placeholder ID", "op-6: targets the same object as op-1",
	} {
		if !strings.Contains(strings.Join(p.Errors, "\n"), want) {
			t.Errorf("missing error %q in:\n%s", want, strings.Join(p.Errors, "\n"))
		}
	}
	if len(p.Warnings) != 1 {
		t.Errorf("warnings = %v", p.Warnings)
	}
}

func TestFileRoundTripKeepsNullBeforeAndExactNumbers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.changeset.json")
	cs := New("org")
	cs.AddUpdate(huntGroup("hg1"), []Change{
		{Path: "a", Value: json.Number("12345678901234567890"), Before: nil, HasBefore: true},
		{Path: "b", Value: nil},
	})
	if err := Save(path, cs); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	if !strings.Contains(string(data), `"before": null`) || strings.Count(string(data), `"before"`) != 1 {
		t.Errorf("file:\n%s", data)
	}
	back, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	changes := back.Operations[0].Changes
	if !changes[0].HasBefore || changes[1].HasBefore || changes[0].Value != json.Number("12345678901234567890") {
		t.Errorf("changes = %+v", changes)
	}
}

func TestParseRejectsUnknownKeys(t *testing.T) {
	if _, err := Parse([]byte(`{"format":"webex-changeset/v1","metdata":{}}`)); err == nil {
		t.Error("a misspelt key was accepted")
	}
}

func TestConcurrentUpdatesKeepEveryOperation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "x.changeset.json")
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := Edit(path, func() (*ChangeSet, error) { return New("org"), nil }, func(cs *ChangeSet) error {
				cs.AddDelete(Target{API: APIContactCenter, Route: "/skill/{id}", Params: map[string]string{"id": string(rune('a' + i))}})
				return nil
			})
			if err != nil {
				t.Error(err)
			}
		}(i)
	}
	wg.Wait()
	cs, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cs.Operations) != 20 {
		t.Errorf("operations = %d", len(cs.Operations))
	}
}

func TestSchemaIsJSON(t *testing.T) {
	var v map[string]any
	if err := json.Unmarshal(Schema, &v); err != nil {
		t.Fatal(err)
	}
}
