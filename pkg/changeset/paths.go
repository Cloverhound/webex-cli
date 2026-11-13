package changeset

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
)

// Get returns the value at a dotted path. Only objects are walked; a list is a
// leaf.
func Get(doc any, path string) (any, bool) {
	cur := doc
	for _, key := range strings.Split(path, ".") {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, false
		}
		if cur, ok = m[key]; !ok {
			return nil, false
		}
	}
	return cur, true
}

// Set stores v at a dotted path, creating objects on the way and replacing
// anything that is not one.
func Set(doc map[string]any, path string, v any) {
	keys := strings.Split(path, ".")
	cur := doc
	for _, key := range keys[:len(keys)-1] {
		next, ok := cur[key].(map[string]any)
		if !ok {
			next = map[string]any{}
			cur[key] = next
		}
		cur = next
	}
	cur[keys[len(keys)-1]] = v
}

// Diff returns the changes that turn current into desired, for the paths
// desired sets. A path desired leaves out is not a change, because a full-object
// PUT names fields it does not mean to change. Objects are compared key by key
// and lists whole.
func Diff(current, desired map[string]any) []Change {
	var out []Change
	var walk func(prefix string, cur, want map[string]any)
	walk = func(prefix string, cur, want map[string]any) {
		for _, key := range sortedKeys(want) {
			path := prefix + key
			wantValue := want[key]
			curValue, present := cur[key]
			wantMap, wantIsMap := wantValue.(map[string]any)
			if _, isRef := Ref(wantValue); wantIsMap && !isRef && len(wantMap) > 0 {
				curMap, _ := curValue.(map[string]any)
				walk(path+".", curMap, wantMap)
				continue
			}
			// An absent field already reads as null.
			if (present && Equal(curValue, wantValue)) || (!present && wantValue == nil) {
				continue
			}
			out = append(out, Change{Path: path, Value: wantValue, Before: curValue, HasBefore: true})
		}
	}
	walk("", current, desired)
	return out
}

// Leaves lists every path a body sets, with no precondition.
func Leaves(body map[string]any) []Change {
	out := Diff(map[string]any{}, body)
	for i := range out {
		out[i].Before, out[i].HasBefore = nil, false
	}
	return out
}

// Equal compares two decoded JSON values. Numbers compare by value whether they
// were decoded as json.Number or float64.
func Equal(a, b any) bool {
	return reflect.DeepEqual(canonical(a), canonical(b))
}

func canonical(v any) any {
	switch t := v.(type) {
	case json.Number:
		if f, err := t.Float64(); err == nil {
			return f
		}
		return string(t)
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, e := range t {
			out[k] = canonical(e)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, e := range t {
			out[i] = canonical(e)
		}
		return out
	}
	return v
}

// Clone deep-copies a decoded JSON value.
func Clone(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, e := range t {
			out[k] = Clone(e)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, e := range t {
			out[i] = Clone(e)
		}
		return out
	}
	return v
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
