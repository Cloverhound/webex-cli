package resource

import (
	"encoding/json"
	"reflect"
	"testing"
)

var sample = Type{
	Name:        "thing",
	ParentField: "location_id",
	LabelField:  "name",
	Fields: []Field{
		{Name: "location_id", Immutable: true, Ref: "location"},
		{Name: "name", Path: "name"},
		{Name: "rings", Path: "policy.noAnswer.rings", Kind: Int},
		{Name: "enabled", Path: "enabled", Kind: Bool},
		{Name: "tags", Path: "tags", Kind: StringList},
	},
	CreateDefaults: map[string]any{"policy": map[string]any{"mode": "CIRCULAR"}, "agents": []any{}},
}

func TestProjectReadsNestedPathsAndParent(t *testing.T) {
	raw := map[string]any{"name": "a", "enabled": true, "policy": map[string]any{"noAnswer": map[string]any{"rings": 5.0}}, "tags": []any{"x"}, "hidden": "kept"}
	got := sample.Project("loc1", raw)
	want := Values{"location_id": "loc1", "name": "a", "rings": 5.0, "enabled": true, "tags": []any{"x"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v", got)
	}
}

func TestProjectGivesNilForMissingPath(t *testing.T) {
	got := sample.Project("", map[string]any{"policy": "not an object"})
	if got["rings"] != nil {
		t.Fatalf("rings = %v", got["rings"])
	}
}

func TestPatchBodyIsNestedAndSkipsParent(t *testing.T) {
	got := sample.PatchBody(Values{"rings": 7, "location_id": "loc2"})
	want := map[string]any{"policy": map[string]any{"noAnswer": map[string]any{"rings": 7}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v", got)
	}
}

func TestFullBodyKeepsUnmodelledFieldsAndDropsLinks(t *testing.T) {
	raw := map[string]any{"name": "a", "userIds": []any{"u1"}, "_links": "x", "links": "y"}
	got := sample.FullBody(raw, Values{"name": "b"})
	want := map[string]any{"name": "b", "userIds": []any{"u1"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v", got)
	}
	if raw["name"] != "a" {
		t.Fatal("the fresh read must not be modified")
	}
}

func TestCreateBodyLayersOverDefaultsAndOmitsUnknown(t *testing.T) {
	got := sample.CreateBody(Values{"name": "a", "rings": 3, "enabled": nil})
	want := map[string]any{"name": "a", "agents": []any{}, "policy": map[string]any{"mode": "CIRCULAR", "noAnswer": map[string]any{"rings": 3}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v", got)
	}
	if _, ok := sample.CreateDefaults["name"]; ok {
		t.Fatal("defaults must not be modified")
	}
}

func TestChangedTreatsNumberTypesAndEmptyListsAsEqual(t *testing.T) {
	fields := []string{"rings", "tags", "name"}
	a := Values{"rings": int64(5), "tags": []any{}, "name": "a"}
	b := Values{"rings": 5.0, "tags": nil, "name": "b"}
	if got := Changed(a, b, fields); !reflect.DeepEqual(got, []string{"name"}) {
		t.Fatalf("got %v", got)
	}
}

func TestChangedKeepsListOrder(t *testing.T) {
	a := Values{"tags": []any{"x", "y"}}
	b := Values{"tags": []any{"y", "x"}}
	if got := Changed(a, b, []string{"tags"}); len(got) != 1 {
		t.Fatalf("got %v", got)
	}
}

func TestImportIDRoundTrip(t *testing.T) {
	key := Key{Parent: "loc1", ID: "abc"}
	got, err := sample.ParseImportID(key.ImportID())
	if err != nil || got != key {
		t.Fatalf("got %v err %v", got, err)
	}
	if _, err := sample.ParseImportID("noslash"); err == nil {
		t.Fatal("expected an error")
	}
	flat := Type{}
	got, err = flat.ParseImportID("abc")
	if err != nil || got != (Key{ID: "abc"}) {
		t.Fatalf("got %v err %v", got, err)
	}
}

func TestAddressIsStableAndValidHCL(t *testing.T) {
	a := sample.Address("Y2lzY29zcGFyazovL3VzL0hVTlRfR1JPVVA")
	if a != sample.Address("Y2lzY29zcGFyazovL3VzL0hVTlRfR1JPVVA") {
		t.Fatal("address must be deterministic")
	}
	if ResourceName("a")[:3] != "id_" || ResourceName("a") == ResourceName("b") {
		t.Fatalf("names = %s %s", ResourceName("a"), ResourceName("b"))
	}
}

var withRows = Type{
	Name: "profile",
	Fields: []Field{
		{Name: "name", Path: "name"},
		{Name: "skills", Path: "activeSkills", Kind: ObjectList, MatchKey: "skill_id", Item: []Field{
			{Name: "skill_id", Path: "skillId", Ref: "skill"},
			{Name: "level", Path: "proficiencyValue", Kind: Int},
		}},
		{Name: "hours", Path: "workingHours", Kind: ObjectList, Unordered: true, MatchKey: "name", Item: []Field{
			{Name: "name", Path: "name"},
			{Name: "days", Path: "days", Kind: StringList, Unordered: true},
		}},
		{Name: "groups", Path: "groups", Kind: ObjectList, Item: []Field{
			{Name: "order", Path: "order", Kind: Int},
			{Name: "team_ids", Path: "agentGroups[].teamId", Kind: StringList},
		}},
		{Name: "codes", Path: "codes", Kind: StringList, Unordered: true},
	},
}

func TestProjectReadsRowsAndSortsUnorderedValues(t *testing.T) {
	raw := map[string]any{
		"activeSkills": []any{map[string]any{"id": "e1", "skillId": "s1", "proficiencyValue": 5.0, "skillName": "x"}},
		"workingHours": []any{
			map[string]any{"name": "b", "days": []any{"TUE", "MON"}},
			map[string]any{"name": "a", "days": []any{"WED"}},
		},
		"groups": []any{map[string]any{"order": 1.0, "agentGroups": []any{map[string]any{"teamId": "t1"}, map[string]any{"teamId": "t2"}}}},
		"codes":  []any{"z", "a"},
	}
	got := withRows.Project("", raw)
	want := Values{
		"name":   nil,
		"skills": []any{map[string]any{"skill_id": "s1", "level": 5.0}},
		"hours": []any{
			map[string]any{"name": "a", "days": []any{"WED"}},
			map[string]any{"name": "b", "days": []any{"MON", "TUE"}},
		},
		"groups": []any{map[string]any{"order": 1.0, "team_ids": []any{"t1", "t2"}}},
		"codes":  []any{"a", "z"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got  %v\nwant %v", got, want)
	}
}

func TestFullBodyKeepsUnmodelledPartsOfMatchingRows(t *testing.T) {
	raw := map[string]any{"activeSkills": []any{
		map[string]any{"id": "e1", "skillId": "s1", "proficiencyValue": 5.0},
		map[string]any{"id": "e2", "skillId": "s2", "proficiencyValue": 3.0},
	}}
	got := withRows.FullBody(raw, Values{"skills": []any{
		map[string]any{"skill_id": "s2", "level": int64(9)},
		map[string]any{"skill_id": "s3", "level": int64(1)},
	}})
	want := []any{
		map[string]any{"id": "e2", "skillId": "s2", "proficiencyValue": int64(9)},
		map[string]any{"skillId": "s3", "proficiencyValue": int64(1)},
	}
	if !reflect.DeepEqual(got["activeSkills"], want) {
		t.Fatalf("got %v", got["activeSkills"])
	}
}

func TestBodyBuildsListsUnderArrayPaths(t *testing.T) {
	got := withRows.CreateBody(Values{"groups": []any{map[string]any{"order": int64(1), "team_ids": []any{"t1", "t2"}}}})
	want := map[string]any{"groups": []any{map[string]any{"order": int64(1), "agentGroups": []any{
		map[string]any{"teamId": "t1"}, map[string]any{"teamId": "t2"},
	}}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v", got)
	}
}

func TestNormalizeMakesUnorderedValuesComparable(t *testing.T) {
	f, _ := withRows.Field("hours")
	a := f.Normalize([]any{map[string]any{"name": "b", "days": []any{"TUE", "MON"}}, map[string]any{"name": "a", "days": nil}})
	b := f.Normalize([]any{map[string]any{"name": "a", "days": nil}, map[string]any{"name": "b", "days": []any{"MON", "TUE"}}})
	if got := Changed(Values{"hours": a}, Values{"hours": b}, []string{"hours"}); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
}

func TestUnmetLeavesUnknownValuesToWebex(t *testing.T) {
	have := Values{"name": "a", "skills": []any{map[string]any{"skill_id": "s1", "level": 5.0}}}
	planned := Values{"name": nil, "skills": []any{map[string]any{"skill_id": "s1", "level": nil}}}
	if got := withRows.Unmet(have, planned); len(got) != 0 {
		t.Fatalf("got %v", got)
	}
	planned["skills"] = []any{map[string]any{"skill_id": "s1", "level": int64(6)}}
	planned["name"] = "b"
	if got := withRows.Unmet(have, planned); !reflect.DeepEqual(got, []string{"name", "skills"}) {
		t.Fatalf("got %v", got)
	}
	planned["skills"] = []any{}
	if got := withRows.Unmet(have, planned); !reflect.DeepEqual(got, []string{"name", "skills"}) {
		t.Fatalf("an emptied list must count: got %v", got)
	}
}

func TestProjectedValuesAreJSONSafe(t *testing.T) {
	if _, err := json.Marshal(sample.Project("", map[string]any{})); err != nil {
		t.Fatal(err)
	}
}
