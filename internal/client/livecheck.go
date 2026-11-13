package client

import (
	"errors"
	"fmt"

	"github.com/Cloverhound/webex-cli/pkg/catalog"
	"github.com/Cloverhound/webex-cli/pkg/changeset"
	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webex"
)

// Drift is a precondition that no longer holds.
type Drift struct {
	Op     string `json:"op"`
	Path   string `json:"path"`
	Before any    `json:"before"`
	Live   any    `json:"live"`
}

// LiveReport compares a change set's targets with Webex now.
type LiveReport struct {
	Drift      []Drift  `json:"drift,omitempty"`
	Missing    []string `json:"missing,omitempty"`
	Unreadable []string `json:"unreadable,omitempty"`
}

// OK reports whether every precondition was checked and holds.
func (r LiveReport) OK() bool {
	return len(r.Drift) == 0 && len(r.Missing) == 0 && len(r.Unreadable) == 0
}

// CheckLive re-reads every object the change set updates or deletes and
// reports preconditions that have drifted and objects that are gone. A delete keeps
// no preconditions, so for one it reports only whether the object is gone.
func CheckLive(cs *changeset.ChangeSet) LiveReport {
	var report LiveReport
	for _, op := range cs.Operations {
		if op.Action != changeset.Update && op.Action != changeset.Delete {
			continue
		}
		live, err := readTarget(op.Target)
		switch {
		case errors.Is(err, webex.ErrNotFound):
			report.Missing = append(report.Missing, op.ID)
			continue
		case err != nil:
			report.Unreadable = append(report.Unreadable, fmt.Sprintf("%s: %v", op.ID, err))
			continue
		}
		for _, c := range op.Changes {
			if !c.HasBefore {
				continue
			}
			value, _ := changeset.Get(live, c.Path)
			if !changeset.Equal(value, c.Before) {
				report.Drift = append(report.Drift, Drift{Op: op.ID, Path: c.Path, Before: c.Before, Live: value})
			}
		}
	}
	return report
}

// readTarget reads the object a target names, through its catalog type, as
// recording does.
func readTarget(t changeset.Target) (map[string]any, error) {
	area := resource.Calling
	if t.API == changeset.APIContactCenter {
		area = resource.ContactCenter
	}
	typ, r, ok := catalog.ResolveTarget(area, t.Route, false)
	if !ok {
		return nil, fmt.Errorf("%s %s is not a route of a type Webex Admin Studio knows", t.API, t.Route)
	}
	o, err := readCurrentObject(typ, catalog.KeyOf(r, t.Route, t.Params))
	return o.Raw, err
}
