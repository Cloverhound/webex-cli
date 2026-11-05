package changeset

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// New starts an empty change set for an org.
func New(orgID string) *ChangeSet {
	var b [4]byte
	_, _ = rand.Read(b[:])
	now := time.Now().UTC()
	return &ChangeSet{
		Format: Format,
		Metadata: Metadata{
			ID:        "cs_" + now.Format("20060102T150405Z") + "_" + hex.EncodeToString(b[:]),
			CreatedAt: now.Format(time.RFC3339),
		},
		Scope:      Scope{OrgID: orgID},
		Operations: []Operation{},
	}
}

// NewIDPrefix marks the placeholder ID a recorded create returns, so later
// commands can name the object before it exists.
const NewIDPrefix = "new:"

// HandleOf returns the handle in a "new:<handle>" placeholder ID.
func HandleOf(id string) (string, bool) {
	if !strings.HasPrefix(id, NewIDPrefix) {
		return "", false
	}
	return strings.TrimPrefix(id, NewIDPrefix), true
}

// ResolvePlaceholders replaces every "new:<handle>" string in v with a
// reference, and fails on a handle no create in the change set has.
func (cs *ChangeSet) ResolvePlaceholders(v any) (any, error) {
	switch t := v.(type) {
	case string:
		handle, ok := HandleOf(t)
		if !ok {
			return t, nil
		}
		if cs.create(handle) == nil {
			return nil, fmt.Errorf("%s does not match an object created in this change set", t)
		}
		return RefTo(handle), nil
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, e := range t {
			r, err := cs.ResolvePlaceholders(e)
			if err != nil {
				return nil, err
			}
			out[k] = r
		}
		return out, nil
	case []any:
		out := make([]any, len(t))
		for i, e := range t {
			r, err := cs.ResolvePlaceholders(e)
			if err != nil {
				return nil, err
			}
			out[i] = r
		}
		return out, nil
	}
	return v, nil
}

// AddUpdate records field changes to an existing object. Changes to an object
// already in the change set merge into its operation: the latest value wins and
// the first precondition is kept. It returns the operation's ID.
func (cs *ChangeSet) AddUpdate(target Target, changes []Change) string {
	if op := cs.find(target.Key()); op != nil && op.Action == Update {
		for _, c := range changes {
			mergeChange(op, c)
		}
		if target.Display != "" {
			op.Target.Display = target.Display
		}
		return op.ID
	}
	op := Operation{ID: cs.nextID(), Action: Update, Target: target, Changes: changes}
	cs.Operations = append(cs.Operations, op)
	return op.ID
}

func mergeChange(op *Operation, c Change) {
	for i := range op.Changes {
		existing := &op.Changes[i]
		switch {
		case existing.Path == c.Path:
			existing.Value = c.Value
			return
		case strings.HasPrefix(c.Path, existing.Path+"."):
			// A change inside an object this operation already sets whole.
			if m, ok := existing.Value.(map[string]any); ok {
				Set(m, strings.TrimPrefix(c.Path, existing.Path+"."), c.Value)
				return
			}
		}
	}
	op.Changes = append(op.Changes, c)
}

// AddCreate records a new object and returns the operation's ID and the
// object's handle. The handle comes from the body's name where it has one.
func (cs *ChangeSet) AddCreate(target Target, body map[string]any) (string, string) {
	handle := cs.uniqueHandle(handleBase(target, body))
	op := Operation{ID: cs.nextID(), Action: Create, Handle: handle, Target: target, Body: body}
	cs.Operations = append(cs.Operations, op)
	return op.ID, handle
}

// UpdateCreate merges changes into a pending create's body.
func (cs *ChangeSet) UpdateCreate(handle string, changes []Change) (string, error) {
	op := cs.create(handle)
	if op == nil {
		return "", fmt.Errorf("%s%s does not match an object created in this change set", NewIDPrefix, handle)
	}
	body, _ := op.Body.(map[string]any)
	if body == nil {
		body = map[string]any{}
	}
	for _, c := range changes {
		Set(body, c.Path, c.Value)
	}
	op.Body = body
	return op.ID, nil
}

// AddDelete records deleting an existing object. It replaces any update to
// that object, whose changes would no longer mean anything.
func (cs *ChangeSet) AddDelete(target Target) string {
	key := target.Key()
	for i := range cs.Operations {
		op := &cs.Operations[i]
		if op.Target.Key() != key {
			continue
		}
		if op.Action == Delete {
			return op.ID
		}
		if op.Action == Update {
			op.Action, op.Changes = Delete, nil
			if target.Display != "" {
				op.Target.Display = target.Display
			}
			return op.ID
		}
	}
	op := Operation{ID: cs.nextID(), Action: Delete, Target: target}
	cs.Operations = append(cs.Operations, op)
	return op.ID
}

// RemoveCreate withdraws a pending create. It fails while another operation
// still refers to the object.
func (cs *ChangeSet) RemoveCreate(handle string) (string, error) {
	index := -1
	for i, op := range cs.Operations {
		if op.Action == Create && op.Handle == handle {
			index = i
		}
	}
	if index < 0 {
		return "", fmt.Errorf("%s%s does not match an object created in this change set", NewIDPrefix, handle)
	}
	for i, op := range cs.Operations {
		if i != index && refersTo(op, handle) {
			return "", fmt.Errorf("%s still refers to %s%s", op.ID, NewIDPrefix, handle)
		}
	}
	id := cs.Operations[index].ID
	cs.Operations = append(cs.Operations[:index], cs.Operations[index+1:]...)
	return id, nil
}

// AddInvoke records a write with no field-level meaning.
func (cs *ChangeSet) AddInvoke(method string, target Target, body any) string {
	op := Operation{ID: cs.nextID(), Action: Invoke, Method: method, Target: target, Body: body}
	cs.Operations = append(cs.Operations, op)
	return op.ID
}

func refersTo(op Operation, handle string) bool {
	for _, h := range Refs(op.Body) {
		if h == handle {
			return true
		}
	}
	for _, c := range op.Changes {
		for _, h := range Refs(c.Value) {
			if h == handle {
				return true
			}
		}
	}
	return false
}

func (cs *ChangeSet) find(key string) *Operation {
	for i := range cs.Operations {
		if cs.Operations[i].Action != Create && cs.Operations[i].Target.Key() == key {
			return &cs.Operations[i]
		}
	}
	return nil
}

func (cs *ChangeSet) create(handle string) *Operation {
	for i := range cs.Operations {
		if cs.Operations[i].Action == Create && cs.Operations[i].Handle == handle {
			return &cs.Operations[i]
		}
	}
	return nil
}

func (cs *ChangeSet) nextID() string {
	highest := 0
	for _, op := range cs.Operations {
		var n int
		if _, err := fmt.Sscanf(op.ID, "op-%d", &n); err == nil && n > highest {
			highest = n
		}
	}
	return fmt.Sprintf("op-%d", highest+1)
}

func (cs *ChangeSet) uniqueHandle(base string) string {
	taken := map[string]bool{}
	for _, op := range cs.Operations {
		if op.Handle != "" {
			taken[op.Handle] = true
		}
	}
	handle := base
	for n := 2; taken[handle]; n++ {
		handle = fmt.Sprintf("%s-%d", base, n)
	}
	return handle
}

var nonSlug = regexp.MustCompile(`[^a-z0-9]+`)

func handleBase(target Target, body map[string]any) string {
	for _, key := range []string{"name", "displayName", "firstName"} {
		if s, ok := body[key].(string); ok {
			if slug := strings.Trim(nonSlug.ReplaceAllString(strings.ToLower(s), "-"), "-"); slug != "" {
				return slug
			}
		}
	}
	segments := strings.Split(strings.Trim(target.Route, "/"), "/")
	for i := len(segments) - 1; i >= 0; i-- {
		if s := segments[i]; s != "" && !strings.HasPrefix(s, "{") {
			return "new-" + strings.Trim(nonSlug.ReplaceAllString(strings.ToLower(s), "-"), "-")
		}
	}
	return "new"
}
