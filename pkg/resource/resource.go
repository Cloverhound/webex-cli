// Package resource describes Webex object types as data: their fields, the JSON
// paths those fields live at, and the operations that read and write them. Reading
// and writing (package webex), the write protocol (package guard), the CLI's
// resource commands and Webex Admin Studio all work from these definitions, so
// adding a type means adding a definition, not new code paths.
//
// Definitions hold API facts only. How a program presents a type, such as its
// sidebar folder or form sections, belongs to that program.
package resource

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// Values holds an object's modelled fields by field name.
type Values = map[string]any

// Kind is a field's Terraform type.
type Kind int

const (
	String Kind = iota
	Bool
	Int
	StringList

	// ObjectList is a list of rows, each holding the fields in Item.
	ObjectList
)

// Field maps one Terraform attribute to one JSON location in the Webex object.
type Field struct {
	Name string

	// Path is a dotted JSON path. Empty means the value is the object's parent ID,
	// which Webex carries in the URL, not the body. A "[]" step, as in
	// "agentGroups[].teamId", maps a StringList onto one key of each array element.
	Path string

	Kind Kind

	// Immutable fields can be set at create but never changed, because changing
	// them would mean deleting and recreating the object.
	Immutable bool

	// Ref names the Type this field's ID points at.
	Ref string

	// Options lists the values a field accepts, when they are fixed.
	Options []string

	// OptionSource names where the accepted values come from when they are not fixed,
	// such as the time zones Webex knows.
	OptionSource string

	// OptionScope names the field whose value the OptionSource's choices depend on, as a
	// flow's versions depend on the flow. Without it, choices depend on the parent.
	OptionScope string

	// Format is the layout of a String holding a date or time. Webex stores these as
	// text, not timestamps.
	Format string

	// Required fields must be set to create an object, because Webex has no default
	// for them.
	Required bool

	// Needs names a field Webex requires whenever this one is set, as an entry point's
	// flow needs music on hold.
	Needs string

	// Unique fields identify one object, such as an extension, so they are never set
	// on several objects at once.
	Unique bool

	// Unordered lists are compared as sets, because Webex does not keep their order.
	Unordered bool

	// Item holds an ObjectList row's fields, with paths relative to the element.
	Item []Field

	// MatchKey names the Item field that identifies a row across writes. A row whose
	// key matches an element Webex holds is written over that element, keeping its ID:
	// Webex rejects a list update that re-adds an existing element without it.
	MatchKey string
}

// Sources of accepted values that are not fixed.
const (
	TimeZones = "timezones"
	Languages = "languages"
	Schedules = "schedules"
	FlowTags  = "flow_tags"
)

// Formats of date and time text, as Webex writes them.
const (
	Date     = "date"     // 2026-12-25
	Time     = "time"     // 08:00
	DateTime = "datetime" // 2026-12-25T08:00
)

// Area selects the Webex API host a Type lives on.
type Area int

const (
	Calling Area = iota
	ContactCenter
)

// WriteMode says how a PUT behaves on a Type's route.
type WriteMode int

const (
	// WriteDelta sends only changed fields, for routes that merge the body.
	WriteDelta WriteMode = iota

	// WriteFull sends the whole object from a fresh read, for routes that replace
	// it, so fields the model does not cover are kept.
	WriteFull

	// WritePatch sends only changed fields with PATCH, for routes whose PUT needs a
	// file upload or a different body.
	WritePatch
)

// Route is the operation a type uses for one action. Parent and ID name the
// operation's path parameters that take the object's parent and ID. The client fills
// a Contact Center operation's org parameter and adds the org to Webex routes.
type Route struct {
	Op     *webexops.Operation
	Parent string
	ID     string

	// Query is sent on every request through this route.
	Query map[string]string
}

// Defined reports whether the type has this route.
func (r Route) Defined() bool { return r.Op != nil }

// ccOrgPrefix starts every Contact Center path that names the org first.
const ccOrgPrefix = "/organization/{orgid}"

// Template is the route's path with {parent} and {id} placeholders, relative to its
// API's base URL, the form change set files and the fake API use. Fixed parameters
// are filled in. A Contact Center path is relative to /organization/{org}; one that
// places the org elsewhere, as flows do, is relative to the host and keeps {org}.
func (r Route) Template() string {
	if r.Op == nil {
		return ""
	}
	path := r.Op.Path
	for _, p := range r.Op.Params {
		if p.In != webexops.InPath {
			continue
		}
		switch {
		case p.Value != "":
			path = strings.ReplaceAll(path, "{"+p.Name+"}", p.Value)
		case p.Name == r.Parent:
			path = strings.ReplaceAll(path, "{"+p.Name+"}", "{parent}")
		case p.Name == r.ID:
			path = strings.ReplaceAll(path, "{"+p.Name+"}", "{id}")
		}
	}
	if r.Op.Host == webexops.HostCC {
		if rest, ok := strings.CutPrefix(path, ccOrgPrefix); ok {
			return rest
		}
		path = strings.ReplaceAll(path, "{orgid}", "{org}")
	}
	return path
}

// Path fills the route's path parameters. A Contact Center org parameter takes
// orgUUID. It fails on a parameter it has no value for.
func (r Route) Path(key Key, orgUUID string) (string, error) {
	path := r.Op.Path
	for _, p := range r.Op.Params {
		if p.In != webexops.InPath {
			continue
		}
		var v string
		switch {
		case p.Value != "":
			v = p.Value
		case p.Name == r.Parent:
			v = key.Parent
		case p.Name == r.ID:
			v = key.ID
		case p.Name == "orgid":
			v = orgUUID
		}
		if v == "" {
			return "", fmt.Errorf("%s: no value for path parameter %s", r.Op.ID, p.Name)
		}
		path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(v))
	}
	return path, nil
}

// Routes are the operations a type is read and written through.
type Routes struct {
	List   Route
	Create Route
	Read   Route
	Update Route
	Delete Route

	// ArrayKey names the array in a Calling list response.
	ArrayKey string

	// ParentKey names the parent ID in a list item when items are grouped under a
	// parent.
	ParentKey string

	// ListIsPartial means list items lack fields, so each needs an item read.
	ListIsPartial bool

	// Updated names a list item field holding the object's last change time, in
	// milliseconds. Webex moves it for changes to fields the list leaves out too, so an
	// unchanged value lets a sync keep its earlier item read.
	Updated string

	// IDKey names the ID field in a list item, when it is not "id".
	IDKey string

	// FallbackIDKey is used for items without IDKey, as for extension-only numbers.
	FallbackIDKey string

	// PerLocation means List holds {parent} and is read once for each calling location.
	PerLocation bool

	// Singleton means List returns one settings object, not a list. With PerLocation
	// there is one per location, keyed by the location.
	Singleton bool

	// ReadByFilter reads one object through List with an id filter, for types whose
	// item route accepts writes but refuses reads.
	ReadByFilter bool
}

// Type defines one Webex object type.
type Type struct {
	Name   string
	Title  string
	Plural string
	Area   Area

	// Managed types can be written: their writes are proven on a live tenant. Other
	// types are read-only, read for display and to label references.
	Managed bool

	// NoCreate types cannot be created here, because Webex needs something the model
	// does not cover, such as an uploaded file.
	NoCreate bool

	// NoDelete types cannot be deleted here.
	NoDelete bool

	ParentField string
	LabelField  string

	// FallbackLabelField labels objects whose LabelField is empty, as a dial number
	// with only an extension.
	FallbackLabelField string

	Fields []Field
	Routes Routes
	Write  WriteMode

	// CreateDefaults are merged under the create body, for fields Webex requires
	// that the model does not cover.
	CreateDefaults map[string]any

	// CopyOnCreate fills a create body from an existing object of the same type, for
	// content Webex requires that has no sensible fixed default, such as a layout file.
	CopyOnCreate *Copy
}

// Copy names the object a create copies from and what it copies.
type Copy struct {
	// Match picks the first object whose raw fields hold these values.
	Match map[string]any

	// Keys are copied when the create body does not set them.
	Keys []string
}

// Key identifies a Webex object.
type Key struct {
	Parent string `json:"parent,omitempty"`
	ID     string `json:"id"`
}

// ImportID is the string Terraform import uses for this object.
func (k Key) ImportID() string {
	if k.Parent != "" {
		return k.Parent + "/" + k.ID
	}
	return k.ID
}

// ParseImportID reverses Key.ImportID for a Type.
func (t Type) ParseImportID(s string) (Key, error) {
	if t.ParentField == "" {
		return Key{ID: s}, nil
	}
	parent, id, ok := strings.Cut(s, "/")
	if !ok || parent == "" || id == "" {
		return Key{}, fmt.Errorf("import ID %q must be <parent>/<id>", s)
	}
	return Key{Parent: parent, ID: id}, nil
}

// TerraformType is the type's resource name in Terraform config, as Webex Admin
// Studio writes it.
func (t Type) TerraformType() string { return "webex_" + t.Name }

// ResourceName derives the Terraform name from the Webex ID, so a rename in Webex is
// an edit to the same resource and never a replacement.
func ResourceName(id string) string {
	sum := sha256.Sum256([]byte(id))
	return "id_" + hex.EncodeToString(sum[:6])
}

// Address is the full Terraform address for an object of this Type.
func (t Type) Address(id string) string {
	return t.TerraformType() + "." + ResourceName(id)
}

// FieldNames lists the Type's field names in definition order.
func (t Type) FieldNames() []string {
	names := make([]string, len(t.Fields))
	for i, f := range t.Fields {
		names[i] = f.Name
	}
	return names
}

// Field looks up a field by name.
func (t Type) Field(name string) (Field, bool) {
	for _, f := range t.Fields {
		if f.Name == name {
			return f, true
		}
	}
	return Field{}, false
}

// Project extracts the modelled fields from a raw Webex object.
func (t Type) Project(parent string, raw map[string]any) Values {
	out := Values{}
	for _, f := range t.Fields {
		if f.Path == "" {
			out[f.Name] = parent
			continue
		}
		out[f.Name] = f.project(raw)
	}
	return out
}

func (f Field) project(raw map[string]any) any {
	v := lookup(raw, f.Path)
	if f.Kind != ObjectList {
		return f.Normalize(v)
	}
	elements, ok := v.([]any)
	if !ok {
		return nil
	}
	rows := make([]any, 0, len(elements))
	for _, e := range elements {
		element, _ := e.(map[string]any)
		row := map[string]any{}
		for _, sub := range f.Item {
			row[sub.Name] = lookup(element, sub.Path)
		}
		rows = append(rows, row)
	}
	return f.Normalize(rows)
}

// Normalize puts a value in a canonical order: unordered lists are sorted, at any depth,
// so two values that Webex treats as equal compare equal.
func (f Field) Normalize(v any) any {
	list, ok := v.([]any)
	if !ok {
		return v
	}
	out := make([]any, len(list))
	for i, e := range list {
		row, isRow := e.(map[string]any)
		if f.Kind != ObjectList || !isRow {
			out[i] = e
			continue
		}
		normal := make(map[string]any, len(row))
		for k, x := range row {
			normal[k] = x
		}
		for _, sub := range f.Item {
			if x, ok := normal[sub.Name]; ok {
				normal[sub.Name] = sub.Normalize(x)
			}
		}
		out[i] = normal
	}
	if f.Unordered {
		sortKey := func(v any) string {
			if row, ok := v.(map[string]any); ok && f.MatchKey != "" {
				return canonical(row[f.MatchKey]) + "\x00" + canonical(v)
			}
			return canonical(v)
		}
		sort.SliceStable(out, func(i, j int) bool { return sortKey(out[i]) < sortKey(out[j]) })
	}
	return out
}

// Label is the display name of an object.
func (t Type) Label(values Values) string {
	if s, ok := values[t.LabelField].(string); ok {
		return s
	}
	return ""
}

// PatchBody builds a nested body from changed fields. Fields with no path are
// skipped, since they travel in the URL.
func (t Type) PatchBody(delta Values) map[string]any {
	body := map[string]any{}
	t.apply(body, delta)
	return body
}

// FullBody overlays changed fields on a fresh read of the object.
func (t Type) FullBody(raw map[string]any, delta Values) map[string]any {
	body := clone(raw)
	delete(body, "_links")
	delete(body, "links")
	t.apply(body, delta)
	return body
}

// CreateBody builds a create request from desired values over CreateDefaults. Unknown
// values are left out so Webex applies its own defaults.
func (t Type) CreateBody(desired Values) map[string]any {
	body := clone(t.CreateDefaults)
	if body == nil {
		body = map[string]any{}
	}
	known := Values{}
	for k, v := range desired {
		if v != nil {
			known[k] = v
		}
	}
	t.apply(body, known)
	return body
}

func (t Type) apply(body map[string]any, values Values) {
	for _, f := range t.Fields {
		v, ok := values[f.Name]
		if !ok || f.Path == "" {
			continue
		}
		if f.Kind == ObjectList {
			v = f.elements(lookup(body, f.Path), v)
		}
		set(body, f.Path, v)
	}
}

// elements builds an ObjectList's Webex array from rows. A row whose MatchKey matches
// an element in have is written over a copy of that element.
func (f Field) elements(have, rows any) any {
	list, ok := rows.([]any)
	if !ok {
		return rows
	}
	existing, _ := have.([]any)
	key, hasKey := f.item(f.MatchKey)
	out := make([]any, 0, len(list))
	for _, r := range list {
		row, _ := r.(map[string]any)
		element := map[string]any{}
		if hasKey && row[key.Name] != nil {
			for _, e := range existing {
				candidate, ok := e.(map[string]any)
				if ok && canonical(lookup(candidate, key.Path)) == canonical(row[key.Name]) {
					element = clone(candidate)
					break
				}
			}
		}
		for _, sub := range f.Item {
			if x := row[sub.Name]; x != nil {
				set(element, sub.Path, x)
			}
		}
		out = append(out, element)
	}
	return out
}

func (f Field) item(name string) (Field, bool) {
	for _, sub := range f.Item {
		if name != "" && sub.Name == name {
			return sub, true
		}
	}
	return Field{}, false
}

// Changed lists fields whose values differ, comparing canonical JSON so a Terraform
// int64 equals a decoded JSON float64. A nil list and an empty list are equal, because
// Webex omits empty lists on read.
func Changed(a, b Values, fields []string) []string {
	var out []string
	for _, f := range fields {
		if canonical(a[f]) != canonical(b[f]) {
			out = append(out, f)
		}
	}
	sort.Strings(out)
	return out
}

// Unmet lists fields whose planned value is known and not yet what have holds. A nil
// planned value, also inside a row, means Terraform left it for Webex to fill.
func (t Type) Unmet(have, planned Values) []string {
	var out []string
	for _, f := range t.Fields {
		p := planned[f.Name]
		if p == nil {
			continue
		}
		if f.Kind == ObjectList {
			if !rowsMeet(have[f.Name], p) {
				out = append(out, f.Name)
			}
			continue
		}
		if canonical(have[f.Name]) != canonical(p) {
			out = append(out, f.Name)
		}
	}
	sort.Strings(out)
	return out
}

// rowsMeet compares normalized rows by position, ignoring planned nils.
func rowsMeet(have, planned any) bool {
	h, _ := have.([]any)
	p, _ := planned.([]any)
	if len(h) != len(p) {
		return false
	}
	for i := range p {
		want, _ := p[i].(map[string]any)
		got, _ := h[i].(map[string]any)
		for k, v := range want {
			if v != nil && canonical(v) != canonical(got[k]) {
				return false
			}
		}
	}
	return true
}

func canonical(v any) string {
	if l, ok := v.([]any); ok && len(l) == 0 {
		return "null"
	}
	b, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

func lookup(m map[string]any, path string) any {
	if head, rest, ok := strings.Cut(path, "[]."); ok {
		elements, ok := lookup(m, head).([]any)
		if !ok {
			return nil
		}
		out := make([]any, 0, len(elements))
		for _, e := range elements {
			element, _ := e.(map[string]any)
			out = append(out, lookup(element, rest))
		}
		return out
	}
	var cur any = m
	for _, part := range strings.Split(path, ".") {
		obj, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = obj[part]
	}
	return cur
}

func set(m map[string]any, path string, v any) {
	if head, rest, ok := strings.Cut(path, "[]."); ok {
		values, _ := v.([]any)
		elements := make([]any, len(values))
		for i, x := range values {
			element := map[string]any{}
			set(element, rest, x)
			elements[i] = element
		}
		set(m, head, elements)
		return
	}
	parts := strings.Split(path, ".")
	cur := m
	for _, part := range parts[:len(parts)-1] {
		next, ok := cur[part].(map[string]any)
		if !ok {
			next = map[string]any{}
			cur[part] = next
		}
		cur = next
	}
	cur[parts[len(parts)-1]] = v
}

func clone(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	var out map[string]any
	if json.Unmarshal(b, &out) != nil {
		return nil
	}
	return out
}
