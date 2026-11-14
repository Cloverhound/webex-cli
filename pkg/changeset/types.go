// Package changeset defines a portable file of proposed Webex configuration
// changes. It names objects by API route and Webex ID and fields by API JSON
// path, so a tool that knows only the Webex REST API can write one, and an
// engine such as Webex Admin Studio can review and apply it.
package changeset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Format is the value of a v1 file's "format" key.
const Format = "webex-changeset/v1"

// API names the base URL a route is relative to.
const (
	// APIWebex is https://webexapis.com/v1: Calling, Admin and the other
	// webexapis.com areas.
	APIWebex = "webex"
	// APIContactCenter is the Contact Center API. Routes are relative to
	// /organization/{orgId}.
	APIContactCenter = "cc"
)

type Action string

const (
	Create Action = "create"
	Update Action = "update"
	Delete Action = "delete"
	// Invoke is any other write, such as a bulk save or a purge. It is kept so
	// nothing an author did is lost, but it has no field-level meaning.
	Invoke Action = "invoke"
)

type ChangeSet struct {
	Format     string      `json:"format"`
	Metadata   Metadata    `json:"metadata"`
	Scope      Scope       `json:"scope"`
	Intents    []Intent    `json:"intents,omitempty"`
	Operations []Operation `json:"operations"`
	Review     *Review     `json:"review,omitempty"`
}

type Metadata struct {
	ID        string            `json:"id"`
	Title     string            `json:"title,omitempty"`
	CreatedAt string            `json:"createdAt,omitempty"`
	Rationale string            `json:"rationale,omitempty"`
	Author    *Author           `json:"author,omitempty"`
	Links     map[string]string `json:"links,omitempty"`
}

type Author struct {
	Type       string `json:"type,omitempty"` // agent or human
	ID         string `json:"id,omitempty"`
	OnBehalfOf string `json:"onBehalfOf,omitempty"`
}

type Scope struct {
	OrgID string `json:"orgId"`
}

// Intent is a plain-language summary that operations can point to.
type Intent struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
}

type Operation struct {
	ID     string `json:"id"`
	Intent string `json:"intent,omitempty"`
	Action Action `json:"action"`
	// Handle names a created object so other operations can refer to its ID
	// with {"$ref": handle}.
	Handle string `json:"handle,omitempty"`
	Target Target `json:"target"`
	// Changes are an update's field changes.
	Changes []Change `json:"changes,omitempty"`
	// Body is a create's or invoke's request body.
	Body any `json:"body,omitempty"`
	// Method is an invoke's HTTP method.
	Method string `json:"method,omitempty"`
}

type Target struct {
	API string `json:"api"`
	// Route is the API route template. For an update or delete it names the
	// object, or for a settings update the settings; for a create it is the
	// collection. Placeholder names do not matter when routes are compared.
	Route   string            `json:"route"`
	Params  map[string]string `json:"params,omitempty"`
	Query   map[string]string `json:"query,omitempty"`
	Display string            `json:"display,omitempty"`
}

// Change sets the value at a dotted path in the object's API JSON. A list is
// always a single value: Webex replaces lists whole.
type Change struct {
	Path  string
	Value any
	// Before is the value the author saw. HasBefore distinguishes a seen null
	// (the field was absent or empty) from no precondition at all.
	Before    any
	HasBefore bool
}

type changeJSON struct {
	Path   string          `json:"path"`
	Before json.RawMessage `json:"before,omitempty"`
	Value  json.RawMessage `json:"value"`
}

func (c Change) MarshalJSON() ([]byte, error) {
	value, err := json.Marshal(c.Value)
	if err != nil {
		return nil, err
	}
	out := changeJSON{Path: c.Path, Value: value}
	if c.HasBefore {
		if out.Before, err = json.Marshal(c.Before); err != nil {
			return nil, err
		}
	}
	return json.Marshal(out)
}

func (c *Change) UnmarshalJSON(data []byte) error {
	var in changeJSON
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}
	if in.Value == nil {
		return fmt.Errorf("change %q has no value", in.Path)
	}
	*c = Change{Path: in.Path}
	if err := decode(in.Value, &c.Value); err != nil {
		return err
	}
	if in.Before != nil {
		c.HasBefore = true
		return decode(in.Before, &c.Before)
	}
	return nil
}

// Review carries facts about an exported change set. Nothing enforces it.
type Review struct {
	Plan      *PlanInfo  `json:"plan,omitempty"`
	Approvals []Approval `json:"approvals,omitempty"`
}

type PlanInfo struct {
	Hash      string `json:"hash"`
	Commit    string `json:"commit,omitempty"`
	PlannedAt string `json:"plannedAt,omitempty"`
}

type Approval struct {
	By      string `json:"by"`
	At      string `json:"at"`
	Hash    string `json:"hash,omitempty"`
	Comment string `json:"comment,omitempty"`
}

// Ref returns the handle when v is {"$ref": handle}.
func Ref(v any) (string, bool) {
	m, ok := v.(map[string]any)
	if !ok || len(m) != 1 {
		return "", false
	}
	handle, ok := m["$ref"].(string)
	return handle, ok
}

// RefTo builds a reference to a created object's ID.
func RefTo(handle string) map[string]any { return map[string]any{"$ref": handle} }

// Refs lists every handle referred to anywhere in v.
func Refs(v any) []string {
	var out []string
	var walk func(any)
	walk = func(v any) {
		if h, ok := Ref(v); ok {
			out = append(out, h)
			return
		}
		switch t := v.(type) {
		case map[string]any:
			keys := make([]string, 0, len(t))
			for k := range t {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				walk(t[k])
			}
		case []any:
			for _, e := range t {
				walk(e)
			}
		}
	}
	walk(v)
	return out
}

var placeholder = regexp.MustCompile(`\{[^}]*\}`)

// NormalizeRoute replaces every placeholder with {} so routes from different
// sources compare equal.
func NormalizeRoute(route string) string {
	return placeholder.ReplaceAllString(strings.TrimRight(route, "/"), "{}")
}

// Placeholders lists a route's placeholder names in order.
func Placeholders(route string) []string {
	var out []string
	for _, m := range placeholder.FindAllString(route, -1) {
		out = append(out, m[1:len(m)-1])
	}
	return out
}

// Key identifies the object a target names: its API, normalized route, parameter
// values in route order, and query, which can select the object as userEmail
// selects whose preferences.
func (t Target) Key() string {
	parts := []string{t.API, NormalizeRoute(t.Route)}
	for _, name := range Placeholders(t.Route) {
		parts = append(parts, t.Params[name])
	}
	names := make([]string, 0, len(t.Query))
	for name := range t.Query {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		parts = append(parts, name+"="+t.Query[name])
	}
	return strings.Join(parts, "\x00")
}

// Path fills the route's placeholders from params.
func (t Target) Path() string {
	return placeholder.ReplaceAllStringFunc(t.Route, func(m string) string {
		return t.Params[m[1:len(m)-1]]
	})
}

// ID is the value of the route's last placeholder, which for an item route is
// the object's Webex ID.
func (t Target) ID() string {
	names := Placeholders(t.Route)
	if len(names) == 0 || !strings.HasSuffix(t.Route, "}") {
		return ""
	}
	return t.Params[names[len(names)-1]]
}

func decode(data []byte, v any) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	return d.Decode(v)
}
