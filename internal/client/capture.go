package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/catalog"
	"github.com/Cloverhound/webex-cli/pkg/changeset"
	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webex"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// ccOrgPrefix is the org segment every Contact Center route starts with.
// Change set routes are relative to it.
var ccOrgPrefix = regexp.MustCompile(`^/organization/\{([^}]+)\}`)

// readTimeout bounds the read of the object a write changes.
const readTimeout = 30 * time.Second

// capture records a write request in the change set at path instead of sending
// it, and returns the JSON the command prints. Only creates, updates and deletes of
// the types Webex Admin Studio applies are recorded, because Studio refuses a whole
// file over one operation it cannot apply.
func capture(req *Request, path string) ([]byte, error) {
	if req.bodySummary != "" {
		return nil, errors.New("file uploads cannot be recorded in a change set")
	}
	if req.rawURL != "" {
		return nil, errors.New("this request cannot be recorded in a change set")
	}
	w, err := captureWrite(req)
	if err != nil {
		return nil, err
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
	fields, isObject := body.(map[string]any)
	switch {
	case w.Action != catalog.Delete && !isObject:
		return nil, fmt.Errorf("a %s %s needs its fields as a JSON object in --body", w.Type.Title, w.Action)
	case w.Action == catalog.Delete && body != nil:
		return nil, errors.New("a delete with a body cannot be recorded in a change set")
	}
	// A body can name the org too, as a person's create does.
	for k, v := range fields {
		s, isText := v.(string)
		if org := config.DecodeOrgID(s); isText && s != "" && (strings.EqualFold(k, "orgId") || strings.EqualFold(k, "organizationId")) && org != config.OrgID() {
			return nil, fmt.Errorf("the request body names org %s, but change sets record changes for the --organization org (%s)", org, config.OrgID())
		}
	}

	key := resource.Key{Parent: req.pathParams[w.Route.Parent], ID: req.pathParams[w.Route.ID]}
	handle, isNew := changeset.HandleOf(key.ID)
	for name, value := range target.Params {
		if _, ok := changeset.HandleOf(value); ok && (!isNew || value != key.ID) {
			return nil, fmt.Errorf("--%s %s: requests under an object this change set creates cannot be recorded", name, value)
		}
	}

	// Read the object before taking the file lock, which other recorders wait on.
	var current webex.Object
	if !isNew && w.Action != catalog.Create {
		if current, err = readCurrent(w.Type, key); err != nil {
			return nil, err
		}
		target.Display = current.Label
	}

	var result map[string]any
	err = changeset.Edit(path, newChangeSet, func(cs *changeset.ChangeSet) error {
		if org := config.OrgID(); org != "" && cs.Scope.OrgID != org {
			return fmt.Errorf("%s is for org %s, but this command is for org %s", path, cs.Scope.OrgID, org)
		}
		cs.SetIdentity(objectIdentity)
		resolved, err := cs.ResolvePlaceholders(body)
		if err != nil {
			return err
		}
		if isNew {
			if err := sameCreate(cs, handle, w, key); err != nil {
				return err
			}
		}
		result, err = record(cs, w, target, resolved, current.Raw, handle, isNew)
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

// captureWrite finds the catalog type and action a request writes, and refuses
// what Studio cannot apply.
func captureWrite(req *Request) (catalog.Write, error) {
	var host webexops.Host
	switch req.baseURL {
	case config.CallingBaseURL():
		host = webexops.HostWebex
	case config.CcBaseURL():
		host = webexops.HostCC
	default:
		return catalog.Write{}, fmt.Errorf("requests to %s cannot be recorded in a change set", req.baseURL)
	}
	w, ok := catalog.WriteFor(host, req.method, req.path)
	switch {
	case !ok || !w.Type.Managed:
		return w, fmt.Errorf("%s %s is not a create, update or delete of a type Webex Admin Studio can apply, so it cannot be recorded; see 'webex changeset types'", req.method, req.path)
	case w.Action == catalog.Create && w.Type.NoCreate:
		return w, fmt.Errorf("%s cannot be created in a change set", w.Type.Plural)
	case w.Action == catalog.Delete && w.Type.NoDelete:
		return w, fmt.Errorf("%s cannot be deleted in a change set", w.Type.Plural)
	}
	return w, nil
}

// readCurrent reads the object a write changes, through the catalog's reader, which
// knows each type's read route, including types read by filter.
func readCurrent(t resource.Type, key resource.Key) (webex.Object, error) {
	o, err := readCurrentObject(t, key)
	switch {
	case errors.Is(err, webex.ErrNotFound):
		return o, fmt.Errorf("%s %s is not in Webex, so nothing was recorded", t.Title, key.ID)
	case err != nil:
		return o, fmt.Errorf("reading %s %s: %w; nothing was recorded", t.Title, key.ID, err)
	}
	return o, nil
}

func readCurrentObject(t resource.Type, key resource.Key) (webex.Object, error) {
	ctx, cancel := context.WithTimeout(context.Background(), readTimeout)
	defer cancel()
	c := &webex.Client{API: SDK(), Hosts: config.Hosts(), OrgUUID: config.OrgID(), OrgID: config.OrgIDBase64()}
	return c.Read(ctx, t, key)
}

// sameCreate checks that a new:<handle> ID names a pending create of the type the
// request writes, under the same parent.
func sameCreate(cs *changeset.ChangeSet, handle string, w catalog.Write, key resource.Key) error {
	op, ok := cs.CreateOf(handle)
	if !ok {
		return fmt.Errorf("%s%s does not match an object created in this change set", changeset.NewIDPrefix, handle)
	}
	create := w.Type.Routes.Create
	if changeset.NormalizeRoute(op.Target.Route) != changeset.NormalizeRoute(create.Template()) {
		return fmt.Errorf("%s%s is not a new %s", changeset.NewIDPrefix, handle, w.Type.Title)
	}
	if create.Parent != "" && op.Target.Params[create.Parent] != key.Parent {
		return fmt.Errorf("%s%s is created under %s, not %s", changeset.NewIDPrefix, handle, op.Target.Params[create.Parent], key.Parent)
	}
	return nil
}

// objectIdentity names the object a target addresses by its catalog type and key,
// so writes to one object through different routes merge.
func objectIdentity(t changeset.Target) string {
	area := resource.Calling
	if t.API == changeset.APIContactCenter {
		area = resource.ContactCenter
	}
	typ, r, ok := catalog.ResolveTarget(area, t.Route, false)
	if !ok {
		return t.Key()
	}
	key := catalog.KeyOf(r, t.Route, t.Params)
	return typ.Name + "\x00" + key.Parent + "\x00" + key.ID
}

// record adds one write to cs. current is the object as Webex holds it, for an
// update or delete of an existing object.
func record(cs *changeset.ChangeSet, w catalog.Write, target changeset.Target, body any, current map[string]any,
	handle string, isNew bool) (map[string]any, error) {
	fields, _ := body.(map[string]any)

	switch {
	case w.Action == catalog.Delete && isNew:
		id, err := cs.RemoveCreate(handle)
		return map[string]any{"recorded": id, "action": "removed create"}, err

	case w.Action == catalog.Delete:
		return map[string]any{"recorded": cs.AddDelete(target), "action": changeset.Delete}, nil

	case w.Action == catalog.Update && isNew:
		id, err := cs.UpdateCreate(handle, changeset.Leaves(fields))
		return map[string]any{"recorded": id, "action": changeset.Create, "id": changeset.NewIDPrefix + handle}, err

	case w.Action == catalog.Update:
		owned, ignored := w.Type.Split(fields)
		changes := changeset.Diff(current, owned)
		var withdrawn []string
		// A body that holds only the fields being set can set one back, which withdraws
		// its pending change. A whole-object body repeats every field, so it cannot.
		if w.Type.Write != resource.WriteFull {
			withdrawn = cs.Withdraw(target, unchangedPaths(paths(changeset.Leaves(owned)), changes))
		}
		result := map[string]any{"recorded": nil}
		if len(ignored) > 0 {
			result["ignored"] = ignored
		}
		if len(withdrawn) > 0 {
			result["withdrawn"] = withdrawn
		}
		if len(changes) == 0 {
			if len(withdrawn) == 0 {
				result["unchanged"] = true
			}
			return result, nil
		}
		id, err := cs.AddUpdate(target, changes)
		if err != nil {
			return nil, err
		}
		result["recorded"], result["action"], result["changes"] = id, changeset.Update, paths(changes)
		return result, nil
	}

	if target.Display == "" {
		target.Display = label(fields)
	}
	id, handle := cs.AddCreate(target, fields)
	return map[string]any{"recorded": id, "action": changeset.Create, "id": changeset.NewIDPrefix + handle}, nil
}

func newChangeSet() (*changeset.ChangeSet, error) {
	if config.OrgID() == "" {
		return nil, errors.New("no organization to record changes for: log in or pass --organization")
	}
	return changeset.New(config.OrgID()), nil
}

// captureTarget names the request's object in change set terms. Query parameters
// are refused, except the org and a repeat of a path parameter: Studio applies a
// change from the target's route and the recorded fields, so a query would be lost.
func captureTarget(req *Request) (changeset.Target, error) {
	t := changeset.Target{API: changeset.APIWebex, Route: req.path}
	if req.baseURL == config.CcBaseURL() {
		t.API = changeset.APIContactCenter
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
	for k, v := range req.queryParams {
		if v == "" || strings.EqualFold(k, "orgId") || req.pathParams[k] == v {
			continue
		}
		return t, fmt.Errorf("query parameter %s cannot be recorded in a change set; put the fields in --body", k)
	}
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
	return t, nil
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
