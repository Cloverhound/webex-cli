package changeset

import (
	"fmt"
	"strings"
)

// Problems are what Validate finds. Errors make the file unusable; warnings
// are operations an engine may refuse.
type Problems struct {
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

func (p *Problems) errorf(format string, args ...any) {
	p.Errors = append(p.Errors, fmt.Sprintf(format, args...))
}

func (p *Problems) warnf(format string, args ...any) {
	p.Warnings = append(p.Warnings, fmt.Sprintf(format, args...))
}

// Validate checks a change set's structure without reading Webex.
func Validate(cs *ChangeSet) Problems {
	var p Problems
	if cs.Format != Format {
		p.errorf("format is %q, expected %q", cs.Format, Format)
	}
	if cs.Scope.OrgID == "" {
		p.errorf("scope.orgId is missing")
	}

	intents := map[string]bool{}
	for _, in := range cs.Intents {
		if in.ID == "" {
			p.errorf("an intent has no id")
		}
		if intents[in.ID] {
			p.errorf("intent %q appears twice", in.ID)
		}
		intents[in.ID] = true
	}

	ids := map[string]bool{}
	handles := map[string]bool{}
	targets := map[string]string{}
	for i, op := range cs.Operations {
		name := op.ID
		if name == "" {
			name = fmt.Sprintf("operation %d", i+1)
			p.errorf("%s has no id", name)
		} else if ids[op.ID] {
			p.errorf("%s: id appears twice", name)
		}
		ids[op.ID] = true

		if op.Intent != "" && !intents[op.Intent] {
			p.errorf("%s: intent %q is not defined", name, op.Intent)
		}
		validateTarget(&p, name, op)

		switch op.Action {
		case Create:
			if op.Handle == "" {
				p.errorf("%s: a create needs a handle", name)
			} else if handles[op.Handle] {
				p.errorf("%s: handle %q appears twice", name, op.Handle)
			}
			if _, ok := op.Body.(map[string]any); !ok {
				p.errorf("%s: a create needs a body object", name)
			}
			if len(op.Changes) > 0 {
				p.errorf("%s: a create carries its fields in body, not changes", name)
			}
		case Update:
			if len(op.Changes) == 0 {
				p.errorf("%s: an update needs changes", name)
			}
			paths := map[string]bool{}
			for _, c := range op.Changes {
				if c.Path == "" || strings.HasPrefix(c.Path, ".") || strings.HasSuffix(c.Path, ".") || strings.Contains(c.Path, "..") {
					p.errorf("%s: bad change path %q", name, c.Path)
				}
				if paths[c.Path] {
					p.errorf("%s: path %q appears twice", name, c.Path)
				}
				paths[c.Path] = true
			}
			if op.Body != nil {
				p.errorf("%s: an update carries changes, not a body", name)
			}
		case Delete:
			if len(op.Changes) > 0 || op.Body != nil {
				p.errorf("%s: a delete carries no changes or body", name)
			}
		case Invoke:
			if op.Method == "" {
				p.errorf("%s: an invoke needs a method", name)
			}
			p.warnf("%s: %s %s is not a create, update or delete, so it cannot be reviewed field by field", name, op.Method, op.Target.Route)
		default:
			p.errorf("%s: unknown action %q", name, op.Action)
		}

		if op.Action == Update || op.Action == Delete {
			key := op.Target.Key()
			if other, ok := targets[key]; ok {
				p.errorf("%s: targets the same object as %s", name, other)
			}
			targets[key] = name
		}

		for _, handle := range opRefs(op) {
			if !handles[handle] {
				p.errorf("%s: $ref %q does not name an earlier create", name, handle)
			}
		}
		if op.Action == Create && op.Handle != "" {
			handles[op.Handle] = true
		}
	}
	return p
}

func validateTarget(p *Problems, name string, op Operation) {
	t := op.Target
	if t.API != APIWebex && t.API != APIContactCenter {
		p.errorf("%s: target.api is %q, expected %q or %q", name, t.API, APIWebex, APIContactCenter)
	}
	if !strings.HasPrefix(t.Route, "/") {
		p.errorf("%s: target.route must start with /", name)
	}
	for _, param := range Placeholders(t.Route) {
		value := t.Params[param]
		if value == "" {
			p.errorf("%s: target.params has no %q", name, param)
		} else if _, ok := HandleOf(value); ok {
			p.errorf("%s: target.params.%s is a placeholder ID; fold the change into the create", name, param)
		}
	}
	// An update may target a settings route that has no ID.
	if op.Action == Delete && t.ID() == "" {
		p.errorf("%s: the target route of a %s must end with the object's ID", name, op.Action)
	}
}

func opRefs(op Operation) []string {
	refs := Refs(op.Body)
	for _, c := range op.Changes {
		refs = append(refs, Refs(c.Value)...)
	}
	return refs
}
