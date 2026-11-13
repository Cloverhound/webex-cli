package resource

import (
	"sort"
	"strings"
)

// Split separates a JSON body into the paths the type's fields own and the dotted
// paths no field owns. A field owns its path and everything under it; a list mapped
// with "[]" owns the whole list. A {"$ref": ...} object is one value. This is the
// rule Webex Admin Studio applies when it imports a change, so a change set that
// records only owned paths imports.
func (t Type) Split(body map[string]any) (owned map[string]any, unowned []string) {
	heads := map[string]bool{}
	for _, f := range t.Fields {
		if f.Path == "" {
			continue
		}
		head, _, _ := strings.Cut(f.Path, "[]")
		heads[head] = true
	}
	var split func(prefix string, m map[string]any) map[string]any
	split = func(prefix string, m map[string]any) map[string]any {
		out := map[string]any{}
		for k, v := range m {
			path := k
			if prefix != "" {
				path = prefix + "." + k
			}
			if heads[path] {
				out[k] = v
				continue
			}
			if sub, ok := v.(map[string]any); ok && len(sub) > 0 && !isRef(sub) {
				if kept := split(path, sub); len(kept) > 0 {
					out[k] = kept
				}
				continue
			}
			unowned = append(unowned, path)
		}
		return out
	}
	owned = split("", body)
	sort.Strings(unowned)
	return owned, unowned
}

func isRef(m map[string]any) bool {
	_, ok := m["$ref"].(string)
	return ok && len(m) == 1
}
