package changeset

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func huntGroup(id string) Target {
	return Target{API: APIWebex, Route: "/telephony/config/locations/{locationId}/huntGroups/{huntGroupId}",
		Params: map[string]string{"locationId": "loc1", "huntGroupId": id}}
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

func TestEqualIgnoresNumberSpelling(t *testing.T) {
	if !Equal(json.Number("20"), json.Number("2e1")) || !Equal(json.Number("3"), 3) || Equal("3", 3) {
		t.Error("numbers must compare by value, and only numbers")
	}
}

func TestTargetKeys(t *testing.T) {
	a := Target{API: APIContactCenter, Route: "/team/{id}", Params: map[string]string{"id": "t1"}}
	b := Target{API: APIContactCenter, Route: "/team/{teamId}", Params: map[string]string{"teamId": "t1"}}
	if a.Key() != b.Key() {
		t.Error("placeholder names split one object")
	}
	audio := func(user string) Target {
		return Target{API: APIWebex, Route: "/meetingPreferences/audio", Query: map[string]string{"userEmail": user}}
	}
	if audio("a@example.com").Key() == audio("b@example.com").Key() {
		t.Error("targets that differ in their query are one object")
	}
}

func TestHandleOf(t *testing.T) {
	if h, ok := HandleOf("new:desk"); !ok || h != "desk" {
		t.Errorf("handle = %q %v", h, ok)
	}
	if _, ok := HandleOf("team1"); ok {
		t.Error("an ordinary ID is a placeholder")
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

func TestValidate(t *testing.T) {
	cs := New("org")
	cs.Intents = []Intent{{ID: "i", Summary: "s"}}
	cs.Operations = []Operation{
		{ID: "op-1", Action: Create, Handle: "t", Target: Target{API: APIContactCenter, Route: "/team"}, Body: map[string]any{"name": "T"}},
		{ID: "op-2", Intent: "i", Action: Update, Target: huntGroup("hg1"), Changes: []Change{{Path: "name", Value: RefTo("t")}}},
	}
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
	cs.Operations = []Operation{{ID: "op-1", Action: Update, Target: huntGroup("hg1"), Changes: []Change{
		{Path: "a", Value: json.Number("12345678901234567890"), Before: nil, HasBefore: true},
		{Path: "b", Value: nil},
	}}}
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

func TestSchemaIsJSON(t *testing.T) {
	var v map[string]any
	if err := json.Unmarshal(Schema, &v); err != nil {
		t.Fatal(err)
	}
}
