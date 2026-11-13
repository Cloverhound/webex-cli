package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/changeset"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// ccOrgPrefix is the org segment every Contact Center route starts with.
// Change set routes are relative to it.
var ccOrgPrefix = regexp.MustCompile(`^/organization/\{([^}]+)\}`)

// capture records a write request in the change set at path instead of sending
// it, and returns the JSON the command prints.
func capture(req *Request, path string) ([]byte, error) {
	if req.bodySummary != "" {
		return nil, errors.New("file uploads cannot be recorded in a change set")
	}
	if req.rawURL != "" {
		return nil, errors.New("this request cannot be recorded in a change set")
	}
	target, err := captureTarget(req)
	if err != nil {
		return nil, err
	}

	var body any
	if strings.TrimSpace(req.bodyRaw) != "" {
		d := json.NewDecoder(strings.NewReader(req.bodyRaw))
		d.UseNumber()
		if err := d.Decode(&body); err != nil {
			return nil, fmt.Errorf("the request body is not JSON: %w", err)
		}
	}
	// A body can name the org too, as a person's create does.
	if fields, ok := body.(map[string]any); ok {
		for k, v := range fields {
			s, isText := v.(string)
			if org := config.DecodeOrgID(s); isText && s != "" && (strings.EqualFold(k, "orgId") || strings.EqualFold(k, "organizationId")) && org != config.OrgID() {
				return nil, fmt.Errorf("the request body names org %s, but change sets record changes for the --organization org (%s)", org, config.OrgID())
			}
		}
	}

	handle, isNew := changeset.HandleOf(target.ID())
	for name, value := range target.Params {
		if _, ok := changeset.HandleOf(value); ok && (!isNew || value != target.ID()) {
			return nil, fmt.Errorf("--%s %s: requests under an object this change set creates cannot be recorded", name, value)
		}
	}

	// Read the object before taking the file lock, which other recorders wait on.
	var current map[string]any
	if !isNew && req.method != http.MethodPost {
		if current, err = fetchCurrent(req); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not read %s (%v); recording without preconditions\n", target.Path(), err)
		} else {
			target.Display = label(current)
		}
	}

	var result map[string]any
	err = changeset.Edit(path, newChangeSet, func(cs *changeset.ChangeSet) error {
		if org := config.OrgID(); org != "" && cs.Scope.OrgID != org {
			return fmt.Errorf("%s is for org %s, but this command is for org %s", path, cs.Scope.OrgID, org)
		}
		resolved, err := cs.ResolvePlaceholders(body)
		if err != nil {
			return err
		}
		result, err = record(cs, req.method, target, resolved, current, handle, isNew)
		if err == nil {
			result["changeset"] = path
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

// record adds one write to cs.
func record(cs *changeset.ChangeSet, method string, target changeset.Target, body any, current map[string]any,
	handle string, isNew bool) (map[string]any, error) {
	fields, isObject := body.(map[string]any)

	switch {
	case method == http.MethodDelete && isNew:
		id, err := cs.RemoveCreate(handle)
		return map[string]any{"recorded": id, "action": "removed create"}, err

	// A DELETE to a route that names no object, such as removing numbers from a
	// location, is a bulk action, and a DELETE with a body, such as a recording's
	// reason, carries more than a delete records. Both are recorded with their body
	// below.
	case method == http.MethodDelete && target.ID() != "" && body == nil:
		return map[string]any{"recorded": cs.AddDelete(target), "action": changeset.Delete}, nil

	case (method == http.MethodPut || method == http.MethodPatch) && isObject && isNew:
		id, err := cs.UpdateCreate(handle, changeset.Leaves(fields))
		return map[string]any{"recorded": id, "action": changeset.Create, "id": changeset.NewIDPrefix + handle}, err

	case (method == http.MethodPut || method == http.MethodPatch) && isObject:
		changes := changeset.Leaves(fields)
		var withdrawn []string
		if current != nil {
			changes = changeset.Diff(current, fields)
			// Webex API bodies carry only the fields being set, unlike Contact Center's
			// whole-object PUTs.
			if target.API == changeset.APIWebex {
				withdrawn = cs.Withdraw(target, unchangedPaths(paths(changeset.Leaves(fields)), changes))
			}
		}
		if len(changes) == 0 {
			if len(withdrawn) > 0 {
				return map[string]any{"recorded": nil, "withdrawn": withdrawn}, nil
			}
			return map[string]any{"recorded": nil, "unchanged": true}, nil
		}
		id, err := cs.AddUpdate(target, changes)
		if err != nil {
			return nil, err
		}
		result := map[string]any{"recorded": id, "action": changeset.Update, "changes": paths(changes)}
		if len(withdrawn) > 0 {
			result["withdrawn"] = withdrawn
		}
		return result, nil

	case method == http.MethodPost && isObject && isCollection(target):
		if target.Display == "" {
			target.Display = label(fields)
		}
		id, handle := cs.AddCreate(target, fields)
		return map[string]any{"recorded": id, "action": changeset.Create, "id": changeset.NewIDPrefix + handle}, nil
	}

	if isNew {
		return nil, fmt.Errorf("%s %s cannot be recorded against an object this change set creates", method, target.Route)
	}
	return map[string]any{"recorded": cs.AddInvoke(method, target, body), "action": changeset.Invoke}, nil
}

func newChangeSet() (*changeset.ChangeSet, error) {
	if config.OrgID() == "" {
		return nil, errors.New("no organization to record changes for: log in or pass --organization")
	}
	return changeset.New(config.OrgID()), nil
}

func captureTarget(req *Request) (changeset.Target, error) {
	var t changeset.Target
	switch req.baseURL {
	case config.CallingBaseURL():
		t.API = changeset.APIWebex
	case config.CcBaseURL():
		t.API = changeset.APIContactCenter
	default:
		return t, fmt.Errorf("requests to %s cannot be recorded in a change set", req.baseURL)
	}

	// The change set's scope names the org, so every org the request names, in its
	// path or query, must be that org.
	for _, params := range []map[string]string{req.queryParams, req.pathParams} {
		for k, v := range params {
			if org := config.DecodeOrgID(v); strings.EqualFold(k, "orgId") && org != config.OrgID() {
				return t, fmt.Errorf("this request is for org %s, but change sets record changes for the --organization org (%s)", org, config.OrgID())
			}
		}
	}
	t.Route = req.path
	if t.API == changeset.APIContactCenter {
		m := ccOrgPrefix.FindStringSubmatch(t.Route)
		if m == nil {
			return t, fmt.Errorf("%s is not under /organization/{orgid}, so it cannot be recorded in a change set", req.path)
		}
		t.Route = t.Route[len(m[0]):]
		if t.Route == "" {
			t.Route = "/"
		}
	}
	for _, name := range changeset.Placeholders(t.Route) {
		if t.Params == nil {
			t.Params = map[string]string{}
		}
		t.Params[name] = req.pathParams[name]
	}
	for k, v := range req.queryParams {
		if v == "" || strings.EqualFold(k, "orgId") {
			continue
		}
		if t.Query == nil {
			t.Query = map[string]string{}
		}
		t.Query[k] = v
	}
	return t, nil
}

// itemRoutes holds the normalized routes, keyed by API, that the operation table
// reaches one object at. Contact Center routes are relative to /organization/{orgid}.
var itemRoutes = sync.OnceValue(func() map[string][]string {
	out := map[string][]string{}
	for _, area := range webexops.Areas {
		for _, op := range webexops.Operations(area) {
			path := strings.TrimRight(op.Path, "/")
			if !strings.HasSuffix(path, "}") {
				continue
			}
			switch op.Host {
			case webexops.HostWebex:
				out[changeset.APIWebex] = append(out[changeset.APIWebex], changeset.NormalizeRoute(path))
			case webexops.HostCC:
				if m := ccOrgPrefix.FindStringSubmatch(path); m != nil {
					out[changeset.APIContactCenter] = append(out[changeset.APIContactCenter], changeset.NormalizeRoute(path[len(m[0]):]))
				}
			}
		}
	}
	return out
})

// isCollection reports whether a POST to the target creates an object: whether the
// API has a route to one object under it. A POST that names an action, such as
// answering a call or a bulk save, has none.
func isCollection(t changeset.Target) bool {
	item := changeset.NormalizeRoute(t.Route) + "/{}"
	for _, r := range itemRoutes()[t.API] {
		if r == item || strings.HasPrefix(r, item+"/") {
			return true
		}
	}
	return false
}

// fetchCurrent reads the object a write names, from the same URL. The query goes
// too, since it can select the object, as userEmail selects whose preferences.
func fetchCurrent(req *Request) (map[string]any, error) {
	get := NewRequest(req.baseURL, http.MethodGet, req.path)
	for k, v := range req.pathParams {
		get.PathParam(k, v)
	}
	for k, v := range req.queryParams {
		get.QueryParam(k, v)
	}
	data, status, err := Do(get)
	if err != nil {
		return nil, err
	}
	if status < 200 || status > 299 {
		return nil, fmt.Errorf("HTTP %d", status)
	}
	var current map[string]any
	d := json.NewDecoder(bytes.NewReader(data))
	d.UseNumber()
	if err := d.Decode(&current); err != nil {
		return nil, fmt.Errorf("the response is not a JSON object")
	}
	return current, nil
}

func label(object map[string]any) string {
	for _, key := range []string{"name", "displayName", "email"} {
		if s, ok := object[key].(string); ok && s != "" {
			return s
		}
	}
	return ""
}

// unchangedPaths returns the paths in sets that changes leaves out: the fields a
// write sets to the value Webex already holds.
func unchangedPaths(sets []string, changes []changeset.Change) []string {
	var out []string
	for _, p := range sets {
		if !slices.ContainsFunc(changes, func(c changeset.Change) bool { return c.Path == p }) {
			out = append(out, p)
		}
	}
	return out
}

func paths(changes []changeset.Change) []string {
	out := make([]string, len(changes))
	for i, c := range changes {
		out[i] = c.Path
	}
	return out
}
