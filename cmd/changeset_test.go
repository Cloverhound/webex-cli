package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/changeset"
	"github.com/spf13/pflag"
)

// runChangeset runs a changeset subcommand, none of which needs a login.
func runChangeset(t *testing.T, args ...string) (string, error) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	for _, env := range []string{"WEBEX_TOKEN", "WEBEX_REFRESH_TOKEN", "WEBEX_USER", "WEBEX_CHANGESET"} {
		t.Setenv(env, "")
	}
	reset := func(f *pflag.Flag) {
		if sv, ok := f.Value.(pflag.SliceValue); ok {
			sv.Replace(nil)
		} else {
			f.Value.Set(f.DefValue)
		}
		f.Changed = false
	}
	rootCmd.PersistentFlags().VisitAll(reset)
	for _, c := range changesetCmd.Commands() {
		c.Flags().VisitAll(reset)
	}
	var err error
	out := captureStdout(t, func() {
		rootCmd.SetArgs(append([]string{"changeset"}, args...))
		err = rootCmd.Execute()
	})
	return out, err
}

func sampleChangeSet(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "p.changeset.json")
	cs := changeset.New("org1")
	cs.Metadata.Title, cs.Metadata.Rationale = "Desk", "Ticket 1"
	cs.Metadata.Links = map[string]string{"ticket": "https://example.com/1"}
	team := changeset.Target{API: changeset.APIContactCenter, Route: "/team/{id}", Params: map[string]string{"id": "t1"}}
	cs.AddUpdate(team, []changeset.Change{{Path: "name", Value: "A"}, {Path: "capacity", Value: 4}})
	cs.AddDelete(changeset.Target{API: changeset.APIContactCenter, Route: "/skill/{id}", Params: map[string]string{"id": "s1"}})
	if err := changeset.Save(path, cs); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestChangesetDropRevertsOperationsAndChanges(t *testing.T) {
	path := sampleChangeSet(t)
	if _, err := runChangeset(t, "drop", path, "op-1", "capacity"); err != nil {
		t.Fatal(err)
	}
	if _, err := runChangeset(t, "drop", path, "op-2"); err != nil {
		t.Fatal(err)
	}
	cs, _ := changeset.Load(path)
	if len(cs.Operations) != 1 || len(cs.Operations[0].Changes) != 1 || cs.Operations[0].Changes[0].Path != "name" {
		t.Fatalf("operations = %+v", cs.Operations)
	}
	if _, err := runChangeset(t, "drop", path, "op-7"); err == nil || !strings.Contains(err.Error(), "no operation op-7") {
		t.Errorf("err = %v", err)
	}
}

func TestChangesetMetaChangesOnlyTheFlagsGiven(t *testing.T) {
	path := sampleChangeSet(t)
	if _, err := runChangeset(t, "meta", path, "--title", "Front desk", "--on-behalf-of", "a@example.com", "--link", "ticket=", "--link", "doc=https://example.com/d"); err != nil {
		t.Fatal(err)
	}
	cs, _ := changeset.Load(path)
	m := cs.Metadata
	if m.Title != "Front desk" || m.Rationale != "Ticket 1" || m.Author == nil || m.Author.OnBehalfOf != "a@example.com" ||
		m.Links["ticket"] != "" || m.Links["doc"] != "https://example.com/d" {
		t.Fatalf("metadata = %+v, author %+v", m, m.Author)
	}
	if _, err := runChangeset(t, "meta", path, "--rationale", ""); err != nil {
		t.Fatal(err)
	}
	if cs, _ := changeset.Load(path); cs.Metadata.Rationale != "" || cs.Metadata.Title != "Front desk" {
		t.Errorf("metadata = %+v", cs.Metadata)
	}
}

func TestChangesetTypesListsWhatCanBeRecorded(t *testing.T) {
	out, err := runChangeset(t, "types")
	if err != nil {
		t.Fatal(err)
	}
	var types []recordableType
	if err := json.Unmarshal([]byte(out), &types); err != nil {
		t.Fatalf("%v: %s", err, out)
	}
	byName := map[string]recordableType{}
	for _, rt := range types {
		byName[rt.Name] = rt
	}
	if byName["cc_team"].Commands["update"] != "webex cc team update" {
		t.Errorf("cc_team = %+v", byName["cc_team"])
	}
	if aa := byName["auto_attendant"]; aa.Commands["create"] != "" || aa.Commands["update"] == "" {
		t.Errorf("auto attendants can only be updated: %+v", aa.Commands)
	}
	if _, readOnly := byName["person"]; readOnly {
		t.Error("a read-only type is listed")
	}

	out, err = runChangeset(t, "types", "cc_skill")
	if err != nil {
		t.Fatal(err)
	}
	var skill recordableType
	json.Unmarshal([]byte(out), &skill)
	fixed := map[string]bool{}
	for _, f := range skill.Fields {
		fixed[f.Path] = f.Fixed
	}
	if skill.WriteMode != "whole object" || !fixed["name"] || !fixed["skillType"] || fixed["description"] {
		t.Errorf("cc_skill = %s", out)
	}
	if _, err := runChangeset(t, "types", "person"); err == nil {
		t.Error("listed a read-only type")
	}
}
