// Package webex reads and writes Webex objects from their resource definitions. It
// knows routes and body shapes and decides nothing about whether a write is safe;
// that is the guard package's job.
package webex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webexapi"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// ErrNotFound means the object does not exist.
var ErrNotFound = errors.New("not found")

// maxParallelReads matches OpenTofu's default parallelism, which sends the same reads
// during a refresh.
const maxParallelReads = 10

// Outcome classifies a failed write.
type Outcome int

const (
	// Rejected means Webex refused the write and nothing changed.
	Rejected Outcome = iota

	// Unknown means the write may or may not have committed.
	Unknown
)

// Classify says what a write error means. Only a 4xx is a definite refusal: a 5xx or
// a broken connection can follow a committed write.
func Classify(err error) Outcome {
	var apiErr *webexapi.APIError
	if errors.As(err, &apiErr) && apiErr.Status >= 400 && apiErr.Status < 500 {
		return Rejected
	}
	return Unknown
}

// Object is one Webex object.
type Object struct {
	Type   string          `json:"type"`
	Key    resource.Key    `json:"key"`
	Label  string          `json:"label"`
	Values resource.Values `json:"values"`

	// Raw is the object as Webex returned it, including fields outside the model.
	Raw map[string]any `json:"raw"`
}

// Client is the Webex object API for one org.
type Client struct {
	API   *webexapi.Client
	Hosts webexapi.Hosts

	// OrgUUID addresses Contact Center routes.
	OrgUUID string

	// OrgID, when set, targets Webex routes at a specific org, as partner admins
	// need.
	OrgID string

	locationsOnce sync.Mutex
	locations     []string
	locationNames map[string]string
}

// callingLocations lists locations with Calling service; /locations also holds
// locations without it, whose Calling routes fail.
var callingLocations = resource.Route{Op: webexops.CallingLocationCallListWebexCalling}

// url builds a route's URL for one object, with query added to the route's own.
func (c *Client) url(r resource.Route, key resource.Key, query map[string]string) (string, error) {
	path, err := r.Path(key, c.OrgUUID)
	if err != nil {
		return "", err
	}
	q := url.Values{}
	for k, v := range r.Query {
		q.Set(k, v)
	}
	for k, v := range query {
		q.Set(k, v)
	}
	if r.Op.Host == webexops.HostWebex && c.OrgID != "" {
		q.Set("orgId", c.OrgID)
	}
	u := r.Op.Host.URL(c.Hosts) + path
	if encoded := q.Encode(); encoded != "" {
		u += "?" + encoded
	}
	return u, nil
}

func (c *Client) object(t resource.Type, key resource.Key, raw map[string]any) Object {
	values := t.Project(key.Parent, raw)
	label := t.Label(values)
	if label == "" && t.FallbackLabelField != "" {
		label, _ = values[t.FallbackLabelField].(string)
	}
	if label == "" && t.Routes.IDKey != "" {
		label = key.ID
	}
	if label == "" {
		label = t.Title
	}
	return Object{Type: t.Name, Key: key, Label: label, Values: values, Raw: raw}
}

// Known holds objects from an earlier read, so a list can keep one in place of an item
// read when its list entry shows the same change time.
type Known struct {
	// Objects are keyed by type name and ID, as "type/id".
	Objects map[string]Object

	// Before is when the earlier read began, less a margin. Change times are whole
	// seconds, and a second edit within the same second as a read leaves the time
	// unchanged, so objects changed this recently are read again. The margin also
	// covers the difference between this clock and Webex's.
	Before time.Time
}

// unchanged returns the earlier read of an object whose list entry shows no change since.
func (k *Known) unchanged(t resource.Type, key resource.Key, listed map[string]any) (map[string]any, bool) {
	if k == nil || t.Routes.Updated == "" {
		return nil, false
	}
	stamp, ok := listed[t.Routes.Updated].(float64)
	if !ok || stamp >= float64(k.Before.UnixMilli()) {
		return nil, false
	}
	earlier, ok := k.Objects[t.Name+"/"+key.ID]
	if !ok || earlier.Key != key {
		return nil, false
	}
	if was, ok := earlier.Raw[t.Routes.Updated].(float64); !ok || was != stamp {
		return nil, false
	}
	return earlier.Raw, true
}

// List reads every object of a type. Any failure fails the whole list, because a
// partial list would look like deleted objects.
func (c *Client) List(ctx context.Context, t resource.Type) ([]Object, error) {
	return c.ListKnown(ctx, t, nil)
}

// ListKnown is List, keeping known objects that have not changed instead of reading
// them again.
func (c *Client) ListKnown(ctx context.Context, t resource.Type, known *Known) ([]Object, error) {
	if !t.Routes.PerLocation {
		return c.listUnder(ctx, t, "", known)
	}
	locations, err := c.locationIDs(ctx)
	if err != nil {
		return nil, err
	}
	var out []Object
	for _, location := range locations {
		found, err := c.listUnder(ctx, t, location, known)
		if IsUnavailable(err) {
			continue
		}
		if err != nil {
			return nil, err
		}
		out = append(out, found...)
	}
	return out, nil
}

func (c *Client) listUnder(ctx context.Context, t resource.Type, location string, known *Known) ([]Object, error) {
	listURL, err := c.url(t.Routes.List, resource.Key{Parent: location}, nil)
	if err != nil {
		return nil, err
	}
	if t.Routes.Singleton {
		return c.readSingleton(ctx, t, listURL, location)
	}
	var items []map[string]any
	visit := func(page []json.RawMessage) error {
		for _, item := range page {
			var raw map[string]any
			if err := json.Unmarshal(item, &raw); err != nil {
				return fmt.Errorf("decoding %s: %w", t.Name, err)
			}
			if location != "" && t.Routes.ParentKey == "" {
				raw[perLocationKey] = location
			}
			items = append(items, raw)
		}
		return nil
	}
	if err := c.listPages(ctx, t, listURL, visit); err != nil {
		return nil, fmt.Errorf("listing %s: %w", t.Plural, err)
	}
	return c.resolve(ctx, t, items, known)
}

// listPages reads a type's list to its end. Outside Contact Center it pages with
// start and max, unless the operation table records other parameters, Link headers
// only, or none at all.
func (c *Client) listPages(ctx context.Context, t resource.Type, listURL string, visit func([]json.RawMessage) error) error {
	if t.Area == resource.ContactCenter {
		return c.API.ListCC(ctx, listURL, visit)
	}
	p := webexapi.Paging{Start: "start", Size: "max", ArrayKey: t.Routes.ArrayKey}
	if pg := t.Routes.List.Op.Paging; pg != nil && pg.Offset != nil {
		p.Start, p.Size, p.First = pg.Offset.Start, pg.Offset.Size, pg.Offset.First
	}
	return c.API.List(ctx, listURL, p, visit)
}

// perLocationKey carries the location a per-location list was read from, for lists
// whose items do not name it. It is removed before the raw object is kept.
const perLocationKey = "__studioLocation"

// readSingleton reads a settings object. An org-level one is keyed by its type name,
// which is stable, and a per-location one by its location.
func (c *Client) readSingleton(ctx context.Context, t resource.Type, target, location string) ([]Object, error) {
	resp, err := c.API.Send(ctx, webexapi.Request{Method: http.MethodGet, URL: target})
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", t.Plural, err)
	}
	var raw map[string]any
	if err := json.Unmarshal(resp.Body, &raw); err != nil {
		return nil, fmt.Errorf("decoding %s: %w", t.Name, err)
	}
	key := resource.Key{ID: t.Name}
	if location != "" {
		key = resource.Key{Parent: location, ID: t.Name + "-" + location}
	}
	obj := c.object(t, key, raw)
	obj.Label = t.Title
	if location != "" {
		obj.Label = c.locationNames[location]
	}
	return []Object{obj}, nil
}

// locationIDs lists calling locations once per client, because every per-location type
// needs them. The calling route is used because /locations also holds locations with
// no calling service, whose calling routes fail.
func (c *Client) locationIDs(ctx context.Context) ([]string, error) {
	c.locationsOnce.Lock()
	defer c.locationsOnce.Unlock()
	if c.locations != nil {
		return c.locations, nil
	}
	ids := []string{}
	names := map[string]string{}
	listURL, err := c.url(callingLocations, resource.Key{}, nil)
	if err != nil {
		return nil, err
	}
	err = c.API.ListOffset(ctx, listURL, "locations", func(page []json.RawMessage) error {
		for _, item := range page {
			var head struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
			if err := json.Unmarshal(item, &head); err != nil || head.ID == "" {
				return fmt.Errorf("location list item has no id")
			}
			ids = append(ids, head.ID)
			names[head.ID] = head.Name
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("listing calling locations: %w", err)
	}
	c.locations, c.locationNames = ids, names
	return ids, nil
}

// IsUnavailable reports an error meaning this org or login cannot read a route, as
// with a missing license, rather than a failed read.
func IsUnavailable(err error) bool {
	var apiErr *webexapi.APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	return apiErr.Status == http.StatusBadRequest || apiErr.Status == http.StatusForbidden || apiErr.Status == http.StatusNotFound
}

// resolve turns list items into objects, keeping list order so snapshots are stable.
// Item reads run in parallel, because partial lists need one read per object and a
// large org would otherwise take most of an hour.
func (c *Client) resolve(ctx context.Context, t resource.Type, items []map[string]any, known *Known) ([]Object, error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	out := make([]Object, len(items))
	slots := make(chan struct{}, maxParallelReads)
	var wg sync.WaitGroup
	for i, raw := range items {
		slots <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() { <-slots }()
			obj, err := c.fromListItem(ctx, t, raw, known)
			if err != nil {
				cancel(err)
				return
			}
			out[i] = obj
		}()
	}
	wg.Wait()
	if err := context.Cause(ctx); err != nil {
		return nil, err
	}
	return out, nil
}

// ListKeys reads only list routes, without the item reads List may need.
func (c *Client) ListKeys(ctx context.Context, t resource.Type) ([]resource.Key, error) {
	if t.Routes.PerLocation || t.Routes.Singleton {
		objects, err := c.List(ctx, t)
		if err != nil {
			return nil, err
		}
		keys := make([]resource.Key, len(objects))
		for i, o := range objects {
			keys[i] = o.Key
		}
		return keys, nil
	}
	listURL, err := c.url(t.Routes.List, resource.Key{}, nil)
	if err != nil {
		return nil, err
	}
	keys := []resource.Key{}
	visit := func(page []json.RawMessage) error {
		for _, item := range page {
			var raw map[string]any
			if err := json.Unmarshal(item, &raw); err != nil {
				return fmt.Errorf("decoding %s: %w", t.Name, err)
			}
			key, err := itemKey(t, raw)
			if err != nil {
				return err
			}
			keys = append(keys, key)
		}
		return nil
	}
	if err = c.listPages(ctx, t, listURL, visit); err != nil {
		return nil, fmt.Errorf("listing %s: %w", t.Plural, err)
	}
	return keys, nil
}

func itemKey(t resource.Type, raw map[string]any) (resource.Key, error) {
	idKey := t.Routes.IDKey
	if idKey == "" {
		idKey = "id"
	}
	id, _ := raw[idKey].(string)
	if id == "" && t.Routes.FallbackIDKey != "" {
		id, _ = raw[t.Routes.FallbackIDKey].(string)
	}
	if id == "" {
		return resource.Key{}, fmt.Errorf("%s list item has no %s", t.Name, idKey)
	}
	key := resource.Key{ID: id}
	if t.Routes.ParentKey != "" {
		key.Parent, _ = raw[t.Routes.ParentKey].(string)
	}
	return key, nil
}

func (c *Client) fromListItem(ctx context.Context, t resource.Type, raw map[string]any, known *Known) (Object, error) {
	key, err := itemKey(t, raw)
	if err != nil {
		return Object{}, err
	}
	if location, ok := raw[perLocationKey].(string); ok {
		key.Parent = location
		delete(raw, perLocationKey)
	}
	if t.Routes.ListIsPartial {
		// Projected again from the raw read, so a newer catalog applies to it.
		if earlier, ok := known.unchanged(t, key, raw); ok {
			return c.object(t, key, earlier), nil
		}
		return c.Read(ctx, t, key)
	}
	return c.object(t, key, raw), nil
}

// Read fetches one object.
func (c *Client) Read(ctx context.Context, t resource.Type, key resource.Key) (Object, error) {
	if t.Routes.ReadByFilter {
		return c.readByFilter(ctx, t, key)
	}
	if !t.Routes.Read.Defined() {
		return Object{}, fmt.Errorf("%s cannot be read one at a time", t.Plural)
	}
	target, err := c.url(t.Routes.Read, key, nil)
	if err != nil {
		return Object{}, err
	}
	resp, err := c.API.Send(ctx, webexapi.Request{Method: http.MethodGet, URL: target})
	if err != nil {
		if resp != nil && resp.Status == http.StatusNotFound {
			return Object{}, ErrNotFound
		}
		return Object{}, err
	}
	var raw map[string]any
	if err := json.Unmarshal(resp.Body, &raw); err != nil {
		return Object{}, fmt.Errorf("decoding %s: %w", t.Name, err)
	}
	return c.object(t, key, raw), nil
}

func (c *Client) readByFilter(ctx context.Context, t resource.Type, key resource.Key) (Object, error) {
	target, err := c.url(t.Routes.List, key, map[string]string{"filter": `id=="` + key.ID + `"`})
	if err != nil {
		return Object{}, err
	}
	var found map[string]any
	err = c.API.ListCC(ctx, target, func(page []json.RawMessage) error {
		for _, item := range page {
			var raw map[string]any
			if err := json.Unmarshal(item, &raw); err != nil {
				return fmt.Errorf("decoding %s: %w", t.Name, err)
			}
			if raw["id"] == key.ID {
				found = raw
			}
		}
		return nil
	})
	if err != nil {
		return Object{}, err
	}
	if found == nil {
		return Object{}, ErrNotFound
	}
	delete(found, "links")
	return c.object(t, key, found), nil
}

// Create makes an object and returns its key.
func (c *Client) Create(ctx context.Context, t resource.Type, desired resource.Values) (resource.Key, error) {
	key := resource.Key{}
	if t.ParentField != "" {
		key.Parent, _ = desired[t.ParentField].(string)
	}
	if !t.Routes.Create.Defined() {
		return resource.Key{}, fmt.Errorf("%s cannot be created through the API", t.Plural)
	}
	body, err := c.CreateBody(ctx, t, desired)
	if err != nil {
		return resource.Key{}, err
	}
	target, err := c.url(t.Routes.Create, key, nil)
	if err != nil {
		return resource.Key{}, err
	}
	resp, err := c.send(ctx, http.MethodPost, target, body)
	if err != nil {
		return resource.Key{}, err
	}
	var created struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(resp.Body, &created); err != nil || created.ID == "" {
		return resource.Key{}, fmt.Errorf("create response has no id: %.200s", resp.Body)
	}
	key.ID = created.ID
	return key, nil
}

// CreateBody returns the body Create sends for desired, including what the type
// copies from an existing object.
func (c *Client) CreateBody(ctx context.Context, t resource.Type, desired resource.Values) (map[string]any, error) {
	body := t.CreateBody(desired)
	if t.CopyOnCreate != nil {
		if err := c.copyInto(ctx, t, body); err != nil {
			return nil, err
		}
	}
	return body, nil
}

// copyInto fills body's missing CopyOnCreate keys from the first matching object, read
// in full because lists can leave those keys out.
func (c *Client) copyInto(ctx context.Context, t resource.Type, body map[string]any) error {
	objects, err := c.List(ctx, t)
	if err != nil {
		return fmt.Errorf("finding a %s to copy: %w", t.Title, err)
	}
	for _, o := range objects {
		if !matches(o.Raw, t.CopyOnCreate.Match) {
			continue
		}
		source, err := c.Read(ctx, t, o.Key)
		if err != nil {
			return fmt.Errorf("reading %s to copy: %w", o.Label, err)
		}
		for _, k := range t.CopyOnCreate.Keys {
			if _, set := body[k]; !set {
				body[k] = source.Raw[k]
			}
		}
		return nil
	}
	return fmt.Errorf("there is no %s to copy from", t.Title)
}

func matches(raw, want map[string]any) bool {
	for k, v := range want {
		if raw[k] != v {
			return false
		}
	}
	return true
}

// Update writes changed fields. raw is a fresh read, which routes that replace the
// whole object need.
func (c *Client) Update(ctx context.Context, t resource.Type, key resource.Key, raw map[string]any, delta resource.Values) error {
	if !t.Routes.Update.Defined() {
		return fmt.Errorf("%s cannot be changed through the API", t.Plural)
	}
	body := t.PatchBody(delta)
	if t.Write == resource.WriteFull {
		body = t.FullBody(raw, delta)
	}
	target, err := c.url(t.Routes.Update, key, nil)
	if err != nil {
		return err
	}
	_, err = c.send(ctx, t.Routes.Update.Op.Method, target, body)
	return err
}

// Delete removes an object.
func (c *Client) Delete(ctx context.Context, t resource.Type, key resource.Key) error {
	if !t.Routes.Delete.Defined() {
		return fmt.Errorf("%s cannot be deleted through the API", t.Plural)
	}
	target, err := c.url(t.Routes.Delete, key, nil)
	if err != nil {
		return err
	}
	_, err = c.API.Send(ctx, webexapi.Request{Method: http.MethodDelete, URL: target})
	return err
}

func (c *Client) send(ctx context.Context, method, target string, body map[string]any) (*webexapi.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encoding body: %w", err)
	}
	return c.API.Send(ctx, webexapi.Request{Method: method, URL: target, Body: data})
}
