package fakeapi

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/pkg/catalog"
	"github.com/Cloverhound/webex-cli/pkg/resource"
)

// collection is an object type served generically, shaped by what the live routes do.
type collection struct {
	name string

	// listPath serves lists. key names the array; empty means a bare array.
	listPath string
	key      string

	// paged wraps the array in {data, meta}, as CC queue lists do.
	paged bool

	// item matches item and collection paths, with the ID as the last group.
	item *regexp.Regexp

	// deletePath, when set, is the only route that accepts DELETE.
	deletePath *regexp.Regexp

	// hidden are fields only item reads return.
	hidden []string

	// summary limits list items to these fields, for lists that lack detail.
	summary []string

	// required makes PUT replace the object, refusing a body missing these fields.
	// Without it, PUT merges, as Calling routes do.
	required []string
}

var collections = []*collection{
	{
		name: "call_queue", listPath: "/v1/telephony/config/queues", key: "queues",
		item:    regexp.MustCompile(`^/v1/telephony/config/locations/[^/]+/queues(?:/([^/]+))?$`),
		summary: []string{"id", "name", "locationId", "enabled"},
	},
	{
		name: "auto_attendant", listPath: "/v1/telephony/config/autoAttendants", key: "autoAttendants",
		item:    regexp.MustCompile(`^/v1/telephony/config/locations/[^/]+/autoAttendants(?:/([^/]+))?$`),
		summary: []string{"id", "name", "locationId"},
	},
	{
		name: "cc_skill", listPath: "/organization/" + OrgUUID + "/skill",
		item:     regexp.MustCompile(`^/organization/[^/]+/skill(?:/([^/]+))?$`),
		required: []string{"name", "skillType", "active"},
	},
	{
		name: "cc_business_hours", listPath: "/organization/" + OrgUUID + "/business-hours",
		item:     regexp.MustCompile(`^/organization/[^/]+/business-hours(?:/([^/]+))?$`),
		hidden:   []string{"workingHours"},
		required: []string{"name", "timezone", "workingHours"},
	},
	{
		name: "cc_queue", listPath: "/organization/" + OrgUUID + "/v2/contact-service-queue", paged: true,
		item:       regexp.MustCompile(`^/organization/[^/]+/v2/contact-service-queue(?:/([^/]+))?$`),
		deletePath: regexp.MustCompile(`^/organization/[^/]+/contact-service-queue/([^/]+)$`),
		hidden:     []string{"callDistributionGroups"},
		required:   []string{"name", "channelType", "callDistributionGroups"},
	},
	{
		name: "cc_entry_point", listPath: "/organization/" + OrgUUID + "/entry-point",
		item:     regexp.MustCompile(`^/organization/[^/]+/entry-point(?:/([^/]+))?$`),
		required: []string{"name", "channelType"},
	},
	{
		name: "cc_site", listPath: "/organization/" + OrgUUID + "/site",
		item:     regexp.MustCompile(`^/organization/[^/]+/site(?:/([^/]+))?$`),
		required: []string{"name"},
	},
	{
		name: "cc_user", listPath: "/organization/" + OrgUUID + "/user",
		item:     regexp.MustCompile(`^/organization/[^/]+/user(?:/([^/]+))?$`),
		hidden:   []string{"teamIds"},
		required: []string{"email"},
	},
}

// The other managed CC types are served from their catalog routes, each seeded with
// one sample object. Their PUT replaces the object, as every live CC PUT does.
func init() {
	have := map[string]bool{"cc_team": true}
	for _, c := range collections {
		have[c.name] = true
	}
	for _, t := range catalog.Managed() {
		if t.Area != resource.ContactCenter || have[t.Name] {
			continue
		}
		item := t.Routes.Read.Template()
		if item == "" {
			item = t.Routes.Update.Template()
		}
		base := t.Routes.Create.Template()
		if base == "" {
			base = strings.TrimSuffix(item, "/{id}")
		}
		c := &collection{
			name: t.Name, listPath: "/organization/" + OrgUUID + t.Routes.List.Template(), paged: true,
			item:     regexp.MustCompile(`^/organization/[^/]+` + regexp.QuoteMeta(base) + `(?:/([^/]+))?$`),
			required: []string{"id"},
		}
		if del := t.Routes.Delete.Template(); del != "" && del != item {
			c.deletePath = regexp.MustCompile(`^/organization/[^/]+` + regexp.QuoteMeta(strings.TrimSuffix(del, "/{id}")) + `/([^/]+)$`)
		}
		collections = append(collections, c)
		generated = append(generated, t)
	}
}

// generated lists the catalog types served by init, for seeding.
var generated []resource.Type

func seedCollections() map[string]map[string]map[string]any {
	seeds := handwrittenSeeds()
	for _, t := range generated {
		object := sampleObject(t)
		seeds[t.Name] = map[string]map[string]any{object["id"].(string): object}
	}
	for name, objects := range seeds {
		for _, o := range objects {
			if _, ok := o["lastUpdatedTime"]; !ok && strings.HasPrefix(name, "cc_") {
				o["lastUpdatedTime"] = 1000.0
			}
		}
	}
	return seeds
}

func handwrittenSeeds() map[string]map[string]map[string]any {
	return map[string]map[string]map[string]any{
		"cc_site": {"site1": {"id": "site1", "name": "Main Site", "active": true}},
		"cc_user": {"user1": {"id": "user1", "firstName": "Ada", "lastName": "Agent", "email": "ada@example.com", "teamIds": []any{"team1"}}},
		"call_queue": {"cq1": {
			"id": "cq1", "name": "q-cq1", "firstName": "zz", "lastName": "queue", "extension": "2001", "languageCode": "en_us",
			"timeZone": "America/New_York", "enabled": true, "locationId": "loc1",
			"callPolicies":  map[string]any{"policy": "CIRCULAR", "callBounce": map[string]any{"callBounceMaxRings": 8.0}},
			"queueSettings": map[string]any{"queueSize": 5.0, "overflow": map[string]any{"action": "PERFORM_BUSY_TREATMENT"}},
			"agents":        []any{map[string]any{"id": "agentA", "firstName": "Ada", "lastName": "Agent"}},
		}},
		"auto_attendant": {"aa1": {
			"id": "aa1", "name": "a-aa1", "firstName": "zz", "lastName": "aa", "extension": "2002", "languageCode": "en_us",
			"timeZone": "America/New_York", "businessSchedule": "Hours", "locationId": "loc1",
			"businessHoursMenu": map[string]any{"keyConfigurations": []any{map[string]any{"key": "0", "action": "EXIT"}}},
		}},
		"cc_skill": {"skill1": {
			"id": "skill1", "name": "s-skill1", "description": "seed", "active": true, "skillType": "PROFICIENCY",
			"serviceLevelThreshold": 30.0, "dynamicSkill": false,
		}},
		"cc_business_hours": {"hours1": {
			"id": "hours1", "name": "h-hours1", "description": "seed", "timezone": "America/New_York",
			"workingHours": []any{map[string]any{"name": "8-6", "days": []any{"MON"}, "startTime": "08:00", "endTime": "18:00"}},
		}},
		"cc_queue": {"queue1": {
			"id": "queue1", "name": "q-queue1", "description": "seed", "active": true, "timezone": "America/New_York",
			"channelType": "TELEPHONY", "maxActiveContacts": 0.0, "maxTimeInQueue": 600.0, "serviceLevelThreshold": 120.0,
			"callDistributionGroups": []any{map[string]any{"order": 1.0}},
		}},
		"cc_entry_point": {"ep1": {
			"id": "ep1", "name": "e-ep1", "description": "seed", "active": true, "timezone": "America/New_York",
			"channelType": "TELEPHONY", "maximumActiveContacts": 0.0, "serviceLevelThreshold": 1.0,
			"flowId": "cc_flow1", "flowTagId": "Live", "controlFlowScriptUrl": "https://example.com/flow",
		}},
	}
}

// match finds the collection and object ID for a path. A delete-only route matches
// only a DELETE, and a delete at the ordinary item route of such a collection is a 405.
func (s *Server) match(method, path string) (*collection, string, bool) {
	for _, c := range collections {
		if c.deletePath != nil {
			if m := c.deletePath.FindStringSubmatch(path); m != nil {
				return c, m[1], true
			}
		}
		if m := c.item.FindStringSubmatch(path); m != nil {
			return c, m[1], true
		}
	}
	return nil, "", false
}

func (s *Server) collectionOp(c *collection, method, path, id string, body map[string]any) (result, bool) {
	objects := s.objects[c.name]
	if method == http.MethodPost {
		for _, o := range objects {
			if o["name"] == body["name"] {
				return problem(409, 0, "duplicate name"), false
			}
		}
		s.next++
		nid := fmt.Sprintf("%s%d", c.name, s.next)
		o := clone(body)
		o["id"] = nid
		touch(c.name, o)
		objects[nid] = o
		return result{201, map[string]any{"id": nid}}, true
	}
	o, ok := objects[id]
	if !ok {
		return problem(404, 0, "not found"), false
	}
	switch method {
	case http.MethodGet:
		return result{200, clone(o)}, false
	case http.MethodPut:
		if len(c.required) == 0 {
			merge(o, body)
			return result{204, nil}, true
		}
		for _, field := range c.required {
			if _, present := body[field]; !present {
				return problem(400, 0, field+": should not be null or blank."), false
			}
		}
		replaced := clone(body)
		replaced["id"] = id
		touch(c.name, replaced)
		objects[id] = replaced
		return result{200, clone(replaced)}, true
	case http.MethodDelete:
		if c.deletePath != nil && !c.deletePath.MatchString(path) {
			return problem(405, 0, "method"), false
		}
		delete(objects, id)
		return result{204, nil}, true
	}
	return problem(405, 0, "method"), false
}

func (s *Server) collectionList(w http.ResponseWriter, r *http.Request) bool {
	for _, c := range collections {
		if r.URL.Path != c.listPath {
			continue
		}
		items := []any{}
		for _, id := range sortedIDs(s.objects[c.name]) {
			items = append(items, listItem(c, s.objects[c.name][id]))
		}
		if filter := r.URL.Query().Get("filter"); strings.HasPrefix(filter, `id=="`) {
			want := strings.TrimSuffix(strings.TrimPrefix(filter, `id=="`), `"`)
			matched := []any{}
			for _, item := range items {
				if item.(map[string]any)["id"] == want {
					matched = append(matched, item)
				}
			}
			items = matched
		}
		switch {
		case c.key != "":
			start, size := queryInt(r, "start", 0), queryInt(r, "max", 100)
			write(w, result{200, map[string]any{c.key: window(items, start, size)}})
		case c.paged:
			page, size := queryInt(r, "page", 0), queryInt(r, "pageSize", 100)
			total := (len(items) + size - 1) / size
			write(w, result{200, map[string]any{"data": window(items, page*size, size), "meta": map[string]any{"totalPages": total}}})
		default:
			page, size := queryInt(r, "page", 0), queryInt(r, "pageSize", 100)
			write(w, result{200, window(items, page*size, size)})
		}
		return true
	}
	return false
}

func listItem(c *collection, o map[string]any) map[string]any {
	item := clone(o)
	for _, field := range c.hidden {
		delete(item, field)
	}
	if c.summary != nil {
		short := map[string]any{}
		for _, field := range c.summary {
			short[field] = o[field]
		}
		return short
	}
	return item
}

// Object returns a copy of a stored object from a generic collection.
func (s *Server) Object(collection, id string) map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return clone(s.objects[collection][id])
}

// Put stores an object in a generic collection, as another admin would.
func (s *Server) Put(collection string, object map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	o := clone(object)
	touch(collection, o)
	s.objects[collection][object["id"].(string)] = o
}

// touch records a change time as Contact Center does on every write, including writes
// to fields its lists leave out.
func touch(collection string, o map[string]any) {
	if strings.HasPrefix(collection, "cc_") {
		o["lastUpdatedTime"] = float64(time.Now().UnixMilli())
	}
}

// Objects returns copies of every object in a generic collection.
func (s *Server) Objects(collection string) []map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []map[string]any
	for _, id := range sortedIDs(s.objects[collection]) {
		out = append(out, clone(s.objects[collection][id]))
	}
	return out
}
