package catalog

import (
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// Every path parameter of a route must have a value: the object's parent or ID, the
// org, or a fixed value. The operation's method must match the route's use.
func checkRoute(t *testing.T, typ resource.Type, use string, r resource.Route, method string) {
	t.Helper()
	if r.Op.Method != method {
		t.Errorf("%s %s: %s is %s, want %s", typ.Name, use, r.Op.ID, r.Op.Method, method)
	}
	wantHost := webexops.HostWebex
	if typ.Area == resource.ContactCenter {
		wantHost = webexops.HostCC
	}
	if r.Op.Host != wantHost {
		t.Errorf("%s %s: %s is on host %s", typ.Name, use, r.Op.ID, r.Op.Host)
	}
	if r.Parent != "" {
		if _, ok := r.Op.Param(r.Parent); !ok {
			t.Errorf("%s %s: %s has no parameter %s", typ.Name, use, r.Op.ID, r.Parent)
		}
	}
	if r.ID != "" {
		if _, ok := r.Op.Param(r.ID); !ok {
			t.Errorf("%s %s: %s has no parameter %s", typ.Name, use, r.Op.ID, r.ID)
		}
	}
	if _, err := r.Path(resource.Key{Parent: "p", ID: "i"}, "o"); err != nil {
		t.Errorf("%s %s: %v", typ.Name, use, err)
	}
	for name := range r.Query {
		if p, ok := r.Op.Param(name); !ok || p.In != webexops.InQuery {
			t.Errorf("%s %s: %s has no query parameter %s", typ.Name, use, r.Op.ID, name)
		}
	}
}

func TestEveryRouteNamesAMatchingOperation(t *testing.T) {
	for _, typ := range All() {
		r := typ.Routes
		if !r.List.Defined() {
			t.Errorf("%s has no list operation", typ.Name)
			continue
		}
		checkRoute(t, typ, "list", r.List, "GET")
		if r.PerLocation != (r.List.Parent != "") {
			t.Errorf("%s: a per-location list, and only one, binds the location", typ.Name)
		}
		if r.Read.Defined() {
			checkRoute(t, typ, "read", r.Read, "GET")
		}
		if r.Create.Defined() {
			checkRoute(t, typ, "create", r.Create, "POST")
		}
		if r.Update.Defined() {
			method := "PUT"
			if typ.Write == resource.WritePatch {
				method = "PATCH"
			}
			checkRoute(t, typ, "update", r.Update, method)
		}
		if r.Delete.Defined() {
			checkRoute(t, typ, "delete", r.Delete, "DELETE")
		}
		if !typ.Managed {
			if r.Create.Defined() || r.Update.Defined() || r.Delete.Defined() {
				t.Errorf("%s is read-only but has a write operation", typ.Name)
			}
			continue
		}
		if !r.Update.Defined() || (!r.Read.Defined() && !r.ReadByFilter) {
			t.Errorf("%s is managed but cannot be read and updated", typ.Name)
		}
		if typ.NoCreate == r.Create.Defined() || typ.NoDelete == r.Delete.Defined() {
			t.Errorf("%s: NoCreate and NoDelete must match the operations it has", typ.Name)
		}
		for _, item := range []resource.Route{r.Read, r.Update, r.Delete} {
			if item.Defined() && item.ID == "" {
				t.Errorf("%s: %s does not bind the object's ID", typ.Name, item.Op.ID)
			}
			if item.Defined() && (typ.ParentField != "") != (item.Parent != "") {
				t.Errorf("%s: %s must bind a parent exactly when the type has one", typ.Name, item.Op.ID)
			}
		}
	}
}

func TestEveryTypeIsWellFormed(t *testing.T) {
	seen := map[string]bool{}
	for _, typ := range All() {
		if seen[typ.Name] {
			t.Errorf("%s is listed twice", typ.Name)
		}
		seen[typ.Name] = true
		if typ.Title == "" || typ.Plural == "" {
			t.Errorf("%s is incomplete: %+v", typ.Name, typ)
		}
		if typ.Routes.Singleton {
			continue
		}
		if _, ok := typ.Field(typ.LabelField); !ok && typ.LabelField != "" {
			t.Errorf("%s has no label field %q", typ.Name, typ.LabelField)
		}
		if typ.Area == resource.Calling && typ.Routes.ArrayKey == "" {
			t.Errorf("%s needs an array key", typ.Name)
		}
		if typ.ParentField != "" {
			if f, ok := typ.Field(typ.ParentField); !ok || f.Path != "" || !f.Immutable {
				t.Errorf("%s parent field must be an immutable field with no path", typ.Name)
			}
		}
		for _, f := range typ.Fields {
			if (f.Kind == resource.ObjectList) != (len(f.Item) > 0) {
				t.Errorf("%s.%s: only a row list has item fields, and it needs them", typ.Name, f.Name)
			}
			if f.MatchKey != "" {
				if _, ok := (resource.Type{Fields: f.Item}).Field(f.MatchKey); !ok {
					t.Errorf("%s.%s matches rows on unknown field %q", typ.Name, f.Name, f.MatchKey)
				}
			}
			if f.Unordered && f.Kind == resource.ObjectList && f.MatchKey == "" {
				t.Errorf("%s.%s: unordered rows need a match key to sort by", typ.Name, f.Name)
			}
			if f.Unique && (f.Kind != resource.String || f.Immutable || f.Name == typ.LabelField) {
				t.Errorf("%s.%s: only editable text other than the label is marked unique; the rest are never bulk-edited anyway", typ.Name, f.Name)
			}
			if f.Needs != "" {
				if _, ok := typ.Field(f.Needs); !ok {
					t.Errorf("%s.%s needs unknown field %q", typ.Name, f.Name, f.Needs)
				}
			}
			if f.OptionScope != "" {
				if _, ok := typ.Field(f.OptionScope); !ok || f.OptionSource == "" {
					t.Errorf("%s.%s: choices scoped by %q need an option source and a field of that name", typ.Name, f.Name, f.OptionScope)
				}
			}
		}
	}
}

func TestFormatsAreKnownAndOnText(t *testing.T) {
	known := map[string]bool{"": true, resource.Date: true, resource.Time: true, resource.DateTime: true}
	for _, typ := range All() {
		for _, f := range append(append([]resource.Field{}, typ.Fields...), itemFields(typ)...) {
			if !known[f.Format] || (f.Format != "" && (f.Kind != resource.String || f.Options != nil || f.Ref != "")) {
				t.Errorf("%s.%s: format %q needs a plain text field", typ.Name, f.Name, f.Format)
			}
		}
	}
}

func TestReferencesPointAtKnownTypes(t *testing.T) {
	for _, typ := range All() {
		for _, f := range append(append([]resource.Field{}, typ.Fields...), itemFields(typ)...) {
			if f.Ref == "" {
				continue
			}
			if _, err := ByName(f.Ref); err != nil {
				t.Errorf("%s.%s references unknown type %q", typ.Name, f.Name, f.Ref)
			}
		}
	}
}

// Single references to managed objects become Terraform references, so they must not
// form a cycle. List references stay literal IDs and are exempt.
func TestSingleReferencesHaveNoCycle(t *testing.T) {
	edges := map[string][]string{}
	for _, typ := range Managed() {
		for _, f := range typ.Fields {
			if f.Ref != "" && f.Kind == resource.String {
				edges[typ.Name] = append(edges[typ.Name], f.Ref)
			}
		}
	}
	const visiting, done = 1, 2
	state := map[string]int{}
	var visit func(name string, path []string)
	visit = func(name string, path []string) {
		switch state[name] {
		case visiting:
			t.Errorf("reference cycle: %s", strings.Join(append(path, name), " → "))
			return
		case done:
			return
		}
		state[name] = visiting
		for _, next := range edges[name] {
			visit(next, append(path, name))
		}
		state[name] = done
	}
	for name := range edges {
		visit(name, nil)
	}
}

func itemFields(typ resource.Type) []resource.Field {
	var out []resource.Field
	for _, f := range typ.Fields {
		out = append(out, f.Item...)
	}
	return out
}

func TestFieldPathsAreUnique(t *testing.T) {
	for _, typ := range Managed() {
		paths := map[string]bool{}
		for _, f := range typ.Fields {
			if f.Path != "" && paths[f.Path] {
				t.Errorf("%s has two fields at %s", typ.Name, f.Path)
			}
			paths[f.Path] = true
		}
	}
}

func TestContactCenterTypesUseFullWrites(t *testing.T) {
	for _, typ := range Managed() {
		if typ.Area == resource.ContactCenter && typ.Write == resource.WriteDelta {
			t.Errorf("%s: every live CC PUT rejects a partial body", typ.Name)
		}
	}
}

func TestManagedTypesCoverV1(t *testing.T) {
	want := []string{"hunt_group", "call_queue", "auto_attendant", "cc_team", "cc_skill", "cc_business_hours", "cc_queue", "cc_entry_point"}
	got := map[string]bool{}
	for _, typ := range Managed() {
		got[typ.Name] = true
	}
	for _, name := range want {
		if !got[name] {
			t.Errorf("missing %s", name)
		}
	}
}
