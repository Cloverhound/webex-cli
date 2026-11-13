package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/changeset"
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
// reports preconditions that have drifted and objects that are gone.
func CheckLive(cs *changeset.ChangeSet) LiveReport {
	var report LiveReport
	for _, op := range cs.Operations {
		if op.Action != changeset.Update && op.Action != changeset.Delete {
			continue
		}
		live, status, err := readTarget(op.Target)
		switch {
		case status == http.StatusNotFound:
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

func readTarget(t changeset.Target) (map[string]any, int, error) {
	var req *Request
	switch t.API {
	case changeset.APIWebex:
		req = NewRequest(config.CallingBaseURL(), http.MethodGet, t.Route)
		req.QueryParam("orgId", config.OrgIDBase64())
	case changeset.APIContactCenter:
		req = NewRequest(config.CcBaseURL(), http.MethodGet, "/organization/{orgid}"+t.Route)
		req.PathParam("orgid", config.OrgID())
	default:
		return nil, 0, fmt.Errorf("unknown api %q", t.API)
	}
	for k, v := range t.Params {
		req.PathParam(k, v)
	}
	for k, v := range t.Query {
		req.QueryParam(k, v)
	}
	data, status, err := Do(req)
	if err != nil {
		return nil, status, err
	}
	var live map[string]any
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if err := d.Decode(&live); err != nil {
		return nil, status, fmt.Errorf("the response is not a JSON object")
	}
	return live, status, nil
}
