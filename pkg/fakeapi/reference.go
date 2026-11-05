package fakeapi

import (
	"net/http"
	"strings"

	"github.com/Cloverhound/webex-cli/pkg/catalog"
	"github.com/Cloverhound/webex-cli/pkg/resource"
)

// served lists catalog types this package answers with specific behaviour. Managed CC
// types are served as collections, and every other type is read-only in the app, so
// one sample object in the right envelope is enough.
var served = map[string]bool{
	"location": true, "hunt_group": true, "cc_team": true, "call_queue": true, "auto_attendant": true,
}

var sampleText = map[string]string{resource.Date: "2026-12-25", resource.Time: "08:00", resource.DateTime: "2026-12-25T08:00"}

func sampleValues(fields []resource.Field) resource.Values {
	values := resource.Values{}
	for _, f := range fields {
		switch f.Kind {
		case resource.Bool:
			values[f.Name] = true
		case resource.Int:
			values[f.Name] = 1
		case resource.StringList:
			values[f.Name] = []any{"sample"}
		case resource.ObjectList:
			values[f.Name] = []any{map[string]any(sampleValues(f.Item))}
		default:
			values[f.Name] = sampleText[f.Format]
			if f.Format == "" {
				values[f.Name] = "sample " + f.Name
			}
		}
	}
	return values
}

func sampleObject(t resource.Type) map[string]any {
	object := clone(t.PatchBody(sampleValues(t.Fields)))
	object["id"] = t.Name + "1"
	if t.Routes.IDKey != "" {
		object[t.Routes.IDKey] = t.Name + "1"
	}
	if t.Routes.ParentKey != "" {
		object[t.Routes.ParentKey] = "loc1"
	}
	return object
}

func samplePath(t resource.Type) string {
	route := strings.ReplaceAll(t.Routes.List.Template(), "{parent}", "loc1")
	if strings.Contains(route, "{org}") {
		return strings.ReplaceAll(route, "{org}", OrgUUID)
	}
	if t.Area == resource.ContactCenter {
		return "/organization/" + OrgUUID + route
	}
	return "/v1" + route
}

// referenceList answers read-only catalog routes and the calling location list.
func (s *Server) referenceList(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == "/v1/telephony/config/locations" {
		items := []any{}
		for _, o := range s.location {
			items = append(items, o)
		}
		write(w, result{200, map[string]any{"locations": items}})
		return true
	}
	for _, t := range catalog.All() {
		if served[t.Name] || r.URL.Path != samplePath(t) {
			continue
		}
		object := sampleObject(t)
		switch {
		case t.Routes.Singleton:
			write(w, result{200, object})
		case t.Area == resource.ContactCenter:
			write(w, result{200, map[string]any{"data": []any{object}, "meta": map[string]any{"totalPages": 1}}})
		default:
			write(w, result{200, map[string]any{t.Routes.ArrayKey: []any{object}}})
		}
		return true
	}
	return false
}

// Types the fake seeds, for tests that count objects.
func SeededObjects() int {
	return len(catalog.All()) + 3 // three hunt groups and two teams instead of one each
}

// SeededManaged is the number of seeded objects in the Terraform model.
func SeededManaged() int {
	return len(catalog.Managed()) + 3
}
