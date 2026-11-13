package catalog

import (
	"regexp"
	"strings"

	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// Action is the change a write makes to one object.
type Action string

const (
	Create Action = "create"
	Update Action = "update"
	Delete Action = "delete"
)

// Write is the operation a type uses for one action.
type Write struct {
	Type   resource.Type
	Route  resource.Route
	Action Action
}

// WriteFor finds the type whose create, update or delete operation is served from
// host at method and path, a path template exactly as the operation table writes
// it. A managed type wins over a read-only type with the same operation.
func WriteFor(host webexops.Host, method, path string) (Write, bool) {
	var found *Write
	for _, t := range All() {
		for _, w := range []Write{{t, t.Routes.Create, Create}, {t, t.Routes.Update, Update}, {t, t.Routes.Delete, Delete}} {
			op := w.Route.Op
			if op == nil || op.Host != host || op.Method != method || op.Path != path {
				continue
			}
			if found == nil || (t.Managed && !found.Type.Managed) {
				w := w
				found = &w
			}
		}
	}
	if found == nil {
		return Write{}, false
	}
	return *found, true
}

// ResolveTarget finds the type a route template in area names: its create route for
// a create, and otherwise its read, update or delete route, which can differ, as a
// Contact Center queue's do. Contact Center routes are relative to the org, as
// Route.Template writes them, and placeholder names do not matter. A managed type
// wins over a read-only type with the same route.
func ResolveTarget(area resource.Area, route string, create bool) (resource.Type, resource.Route, bool) {
	want := normalize(route)
	var found *Write
	for _, t := range All() {
		if t.Area != area {
			continue
		}
		routes := []resource.Route{t.Routes.Read, t.Routes.Update, t.Routes.Delete}
		if create {
			routes = []resource.Route{t.Routes.Create}
		}
		for _, r := range routes {
			if tmpl := r.Template(); tmpl == "" || normalize(tmpl) != want {
				continue
			}
			if found == nil || (t.Managed && !found.Type.Managed) {
				found = &Write{Type: t, Route: r}
			}
		}
	}
	if found == nil {
		return resource.Type{}, resource.Route{}, false
	}
	return found.Type, found.Route, true
}

// KeyOf reads an object's key from a route template and its parameter values,
// pairing the template's placeholders by position with the type's own route.
func KeyOf(r resource.Route, route string, params map[string]string) resource.Key {
	ours := placeholder.FindAllString(r.Template(), -1)
	theirs := placeholder.FindAllString(route, -1)
	var key resource.Key
	for i, name := range ours {
		if i >= len(theirs) {
			break
		}
		value := params[strings.Trim(theirs[i], "{}")]
		switch name {
		case "{parent}":
			key.Parent = value
		case "{id}":
			key.ID = value
		}
	}
	return key
}

var placeholder = regexp.MustCompile(`\{[^}]*\}`)

func normalize(route string) string {
	return placeholder.ReplaceAllString(strings.TrimRight(route, "/"), "{}")
}
