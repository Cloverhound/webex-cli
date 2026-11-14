package changeset

import (
	"encoding/json"
	"reflect"
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
