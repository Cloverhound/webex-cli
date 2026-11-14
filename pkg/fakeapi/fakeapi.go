// Package fakeapi mimics the Webex behaviour observed on a live tenant, on the same
// routes, for the catalog's types, and adds fault injection a live tenant cannot
// safely provide. Tests use it in-process; Webex Admin Studio serves it to end-to-end
// runs.
package fakeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// OrgUUID is the fake org's Contact Center ID.
const OrgUUID = "org1"

// Hosts are the API hosts of a fake served at baseURL.
func Hosts(baseURL string) webexapi.Hosts {
	return webexapi.Hosts{Webex: baseURL + "/v1", CC: baseURL}
}

// Fault modes.
const (
	LoseResponse    = "lose_response"      // commits, then drops the connection
	Fail500NoCommit = "fail_500_no_commit" // 500, nothing changed
	Fail500Commit   = "fail_500_commit"    // 500 after committing
	Fail401         = "fail_401"           // 401, nothing changed
	Fail403         = "fail_403"           // 403, nothing changed
	DropNoCommit    = "drop_no_commit"     // drops the connection, nothing changed
	DelayCommit     = "delay_commit"       // commits after DelayMS, even if the client left
)

// Fault makes matching requests misbehave.
type Fault struct {
	Method string `json:"method"`
	Path   string `json:"path"` // regexp
	Mode   string `json:"mode"`

	// After holds the fault until a request with this method has been seen.
	After string `json:"after"`

	Skip    int `json:"skip"`
	DelayMS int `json:"delay_ms"`

	// Times limits the fault, 1 by default. Negative means unlimited.
	Times int `json:"times"`

	seen int
}

// LogEntry records a request.
type LogEntry struct {
	Seq       int    `json:"seq"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Committed bool   `json:"committed"`
	Fault     string `json:"fault,omitempty"`
}

// Server is the fake. The zero value is not usable; call New.
type Server struct {
	mu       sync.Mutex
	hg       map[string]map[string]any
	team     map[string]map[string]any
	location []map[string]any
	objects  map[string]map[string]map[string]any
	next     int
	faults   []*Fault
	log      []*LogEntry
	seen     map[string]bool
	latency  time.Duration
}

var (
	hgItem   = regexp.MustCompile(`^/v1/telephony/config/locations/([^/]+)/huntGroups(?:/([^/]+))?$`)
	teamItem = regexp.MustCompile(`^/organization/([^/]+)/team(?:/([^/]+))?$`)
)

// New returns a server holding three hunt groups, two teams and reference data.
func New() *Server {
	s := &Server{}
	s.Reset()
	return s
}

// Reset restores the seeded state and clears faults and the request log.
func (s *Server) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reset()
}

func (s *Server) reset() {
	s.hg, s.team, s.faults, s.log, s.seen, s.next = map[string]map[string]any{}, map[string]map[string]any{}, nil, nil, map[string]bool{}, 100
	s.latency = 0
	for _, id := range []string{"hg1", "hg2", "hg3"} {
		s.hg[id] = seedHG(id)
	}
	for _, id := range []string{"team1", "team2"} {
		s.team[id] = map[string]any{"id": id, "name": "t-" + id, "teamType": "AGENT", "teamStatus": "IN_SERVICE",
			"active": true, "siteId": "site1", "description": "seed", "version": 3.0, "dbId": "db-" + id,
			"lastUpdatedTime": 1000.0, "rankQueuesForTeam": false, "skillProfileId": "unmodelled-sp",
			"userIds": []any{"user1"}, "futureField": "keep-me"}
	}
	s.location = []map[string]any{{"id": "loc1", "name": "Boston", "timeZone": "America/New_York",
		"address": map[string]any{"city": "Boston"}}}
	s.objects = seedCollections()
}

func seedHG(id string) map[string]any {
	return map[string]any{"id": id, "name": "g-" + id, "firstName": "zz", "lastName": "spike", "extension": "1" + id[2:],
		"languageCode": "en_us", "timeZone": "America/New_York", "enabled": true, "esn": "091", "routingPrefix": "091",
		"locationId":              "loc1",
		"alternateNumberSettings": map[string]any{"distinctiveRingEnabled": true},
		"callPolicies":            map[string]any{"policy": "CIRCULAR", "noAnswer": map[string]any{"numberOfRings": 5.0, "systemMaxNumberOfRings": 20.0}},
		"agents":                  []any{map[string]any{"id": "agentA", "firstName": "enriched"}},
		"directoryPin":            "s3cret-unmodelled"}
}

// Seed adds count hunt groups, for scale tests.
func (s *Server) Seed(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < count; i++ {
		id := fmt.Sprintf("hg%d", 1000+i)
		s.hg[id] = seedHG(id)
	}
}

// SetLatency delays every API request.
func (s *Server) SetLatency(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latency = d
}

// AddFault arms a fault.
func (s *Server) AddFault(f Fault) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if f.Times == 0 {
		f.Times = 1
	}
	s.faults = append(s.faults, &f)
}

// External changes an object as another admin would, outside any client request.
func (s *Server) External(kind, id string, patch map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if kind == "hg" {
		merge(s.hg[id], patch)
	} else if objects, ok := s.objects[kind]; ok {
		merge(objects[id], patch)
		touch(kind, objects[id])
	} else {
		merge(s.team[id], patch)
	}
}

// Delete removes an object as another admin would.
func (s *Server) Delete(kind, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if kind == "hg" {
		delete(s.hg, id)
	} else if objects, ok := s.objects[kind]; ok {
		delete(objects, id)
	} else {
		delete(s.team, id)
	}
}

// Log returns a copy of the request log.
func (s *Server) Log() []LogEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]LogEntry, len(s.log))
	for i, e := range s.log {
		out[i] = *e
	}
	return out
}

// Mutations returns logged requests other than GET.
func (s *Server) Mutations() []LogEntry {
	var out []LogEntry
	for _, e := range s.Log() {
		if e.Method != http.MethodGet {
			out = append(out, e)
		}
	}
	return out
}

// HuntGroup returns a copy of a stored hunt group, including hidden fields.
func (s *Server) HuntGroup(id string) map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return clone(s.hg[id])
}

// Team returns a copy of a stored team, including hidden fields.
func (s *Server) Team(id string) map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return clone(s.team[id])
}

func (s *Server) take(method, path string) *Fault {
	for _, f := range s.faults {
		if f.Times == 0 || f.Method != method {
			continue
		}
		if ok, _ := regexp.MatchString(f.Path, path); !ok {
			continue
		}
		if f.After != "" && !s.seen[f.After] {
			continue
		}
		f.seen++
		if f.seen <= f.Skip {
			continue
		}
		if f.Times > 0 {
			f.Times--
		}
		return f
	}
	return nil
}

// merge deep-merges a body the way the live hunt group PUT does. Clearing firstName
// is accepted and ignored, which is real Webex behaviour.
func merge(dst, patch map[string]any) {
	for k, v := range patch {
		if k == "firstName" && (v == nil || v == "") {
			continue
		}
		if m, ok := v.(map[string]any); ok {
			if d, ok := dst[k].(map[string]any); ok {
				merge(d, m)
				continue
			}
		}
		dst[k] = v
	}
}

type result struct {
	status int
	body   any
}

func problem(status, code int, msg string) result {
	return result{status, map[string]any{"errorCode": code, "message": msg}}
}

func clone(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	b, _ := json.Marshal(m)
	var o map[string]any
	_ = json.Unmarshal(b, &o)
	return o
}

func (s *Server) hgOp(method, id string, body map[string]any) (result, bool) {
	if method == http.MethodPost {
		for _, o := range s.hg {
			if o["name"] == body["name"] {
				return problem(400, 4208, "[Error 4208] The name is already in use by a service of this type."), false
			}
		}
		s.next++
		nid := fmt.Sprintf("hg%d", s.next)
		o := seedHG(nid)
		delete(o, "directoryPin")
		merge(o, body)
		o["id"] = nid
		o["enabled"] = true
		s.hg[nid] = o
		return result{201, map[string]any{"id": nid}}, true
	}
	o, ok := s.hg[id]
	if !ok {
		return problem(404, 4003, "not found"), false
	}
	switch method {
	case http.MethodGet:
		out := clone(o)
		delete(out, "directoryPin")
		return result{200, out}, false
	case http.MethodPut:
		merge(o, body)
		return result{204, nil}, true
	case http.MethodDelete:
		delete(s.hg, id)
		return result{204, nil}, true
	}
	return problem(405, 0, "method"), false
}

func (s *Server) teamOp(method, id string, body map[string]any) (result, bool) {
	if method == http.MethodPost {
		for _, o := range s.team {
			if o["name"] == body["name"] && o["siteId"] == body["siteId"] {
				return result{409, map[string]any{"error": map[string]any{"key": "409", "errorKey": "error_duplicate_entry"}}}, false
			}
		}
		s.next++
		nid := fmt.Sprintf("team%d", s.next)
		o := clone(body)
		o["id"], o["version"], o["dbId"], o["lastUpdatedTime"] = nid, 1.0, "db-"+nid, float64(time.Now().UnixMilli())
		s.team[nid] = o
		return result{201, clone(o)}, true
	}
	o, ok := s.team[id]
	if !ok {
		return problem(404, 0, "not found"), false
	}
	switch method {
	case http.MethodGet:
		out := clone(o)
		delete(out, "version")
		delete(out, "dbId")
		out["_links"] = map[string]any{"self": "x"}
		out["siteName"] = "Site-1"
		return result{200, out}, false
	case http.MethodPut:
		for _, f := range []string{"name", "teamType", "teamStatus", "siteId"} {
			if _, ok := body[f]; !ok {
				return problem(400, 0, "missing "+f), false
			}
		}
		n := clone(body)
		n["version"], n["dbId"], n["lastUpdatedTime"] = o["version"].(float64)+1, o["dbId"], float64(time.Now().UnixMilli())
		delete(n, "siteName")
		s.team[id] = n
		return result{200, clone(n)}, true
	case http.MethodDelete:
		delete(s.team, id)
		return result{204, nil}, true
	}
	return problem(405, 0, "method"), false
}

func write(w http.ResponseWriter, r result) {
	w.Header().Set("Trackingid", "FAKE_"+fmt.Sprint(time.Now().UnixNano()))
	if r.body == nil {
		w.WriteHeader(r.status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.status)
	_ = json.NewEncoder(w).Encode(r.body)
}

func drop(w http.ResponseWriter) {
	if h, ok := w.(http.Hijacker); ok {
		if c, _, err := h.Hijack(); err == nil {
			c.Close()
		}
	}
}

func sortedIDs(m map[string]map[string]any) []string {
	ids := make([]string, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	// Numeric-aware order keeps pages stable as seeded counts grow.
	for i := 1; i < len(ids); i++ {
		for j := i; j > 0 && less(ids[j], ids[j-1]); j-- {
			ids[j], ids[j-1] = ids[j-1], ids[j]
		}
	}
	return ids
}

func less(a, b string) bool {
	if len(a) != len(b) {
		return len(a) < len(b)
	}
	return a < b
}

func window[T any](items []T, start, size int) []T {
	if start >= len(items) {
		return []T{}
	}
	end := start + size
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}

func queryInt(r *http.Request, name string, fallback int) int {
	if v, err := strconv.Atoi(r.URL.Query().Get(name)); err == nil {
		return v
	}
	return fallback
}

// list answers list routes. CC lists here return a bare array, as the live tenant's
// team, site and user routes do.
func (s *Server) list(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		return false
	}
	start, size := queryInt(r, "start", 0), queryInt(r, "max", 100)
	page, pageSize := queryInt(r, "page", 0), queryInt(r, "pageSize", 100)
	ccWindow := func(items []any) []any { return window(items, page*pageSize, pageSize) }
	switch {
	case r.URL.Path == "/v1/telephony/config/huntGroups":
		out := []any{}
		for _, id := range sortedIDs(s.hg) {
			o := s.hg[id]
			out = append(out, map[string]any{"id": o["id"], "name": o["name"], "locationId": o["locationId"]})
		}
		write(w, result{200, map[string]any{"huntGroups": window(out, start, size)}})
	case r.URL.Path == "/v1/locations":
		out := []any{}
		for _, o := range s.location {
			out = append(out, o)
		}
		write(w, result{200, map[string]any{"items": window(out, start, size)}})
	case r.URL.Path == "/organization/"+OrgUUID+"/team":
		out := []any{}
		for _, id := range sortedIDs(s.team) {
			o := clone(s.team[id])
			delete(o, "version")
			delete(o, "dbId")
			out = append(out, o)
		}
		write(w, result{200, ccWindow(out)})
	default:
		return false
	}
	return true
}

func (s *Server) serveAPI(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	lat := s.latency
	s.mu.Unlock()
	if lat > 0 {
		time.Sleep(lat)
	}
	raw, _ := io.ReadAll(r.Body)
	var body map[string]any
	_ = json.Unmarshal(raw, &body)

	s.mu.Lock()
	e := &LogEntry{Seq: len(s.log) + 1, Method: r.Method, Path: r.URL.Path}
	f := s.take(r.Method, r.URL.Path)
	s.log = append(s.log, e)
	mode := ""
	if f != nil {
		mode = f.Mode
		e.Fault = mode
	}
	if mode == "" && r.Method == http.MethodGet && (s.list(w, r) || s.collectionList(w, r) || s.referenceList(w, r)) {
		s.mu.Unlock()
		return
	}
	var op func() (result, bool)
	if m := hgItem.FindStringSubmatch(r.URL.Path); m != nil {
		op = func() (result, bool) { return s.hgOp(r.Method, m[2], body) }
	} else if m := teamItem.FindStringSubmatch(r.URL.Path); m != nil {
		op = func() (result, bool) { return s.teamOp(r.Method, m[2], body) }
	} else if c, id, ok := s.match(r.Method, r.URL.Path); ok {
		op = func() (result, bool) { return s.collectionOp(c, r.Method, r.URL.Path, id, body) }
	} else if mode == "" {
		s.mu.Unlock()
		write(w, problem(404, 0, "no route"))
		return
	}
	switch mode {
	case Fail500NoCommit:
		s.seen[r.Method] = true
		s.mu.Unlock()
		write(w, problem(500, 0, "injected"))
		return
	case Fail401:
		s.mu.Unlock()
		write(w, problem(401, 0, "token expired"))
		return
	case Fail403:
		s.mu.Unlock()
		write(w, problem(403, 0, "not licensed"))
		return
	case DropNoCommit:
		s.mu.Unlock()
		drop(w)
		return
	case DelayCommit:
		s.mu.Unlock()
		time.Sleep(time.Duration(f.DelayMS) * time.Millisecond)
		s.mu.Lock()
	}
	if op == nil {
		// A commit-type fault on a route with no item handler has nothing to commit.
		s.mu.Unlock()
		write(w, problem(404, 0, "no route"))
		return
	}
	res, committed := op()
	e.Committed = committed
	s.seen[r.Method] = true
	s.mu.Unlock()
	switch mode {
	case LoseResponse:
		drop(w)
	case Fail500Commit:
		write(w, problem(500, 0, "injected after commit"))
	default:
		write(w, res)
	}
}

func (s *Server) serveControl(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	switch strings.TrimPrefix(r.URL.Path, "/__ctl/") {
	case "reset":
		s.Reset()
	case "latency":
		var v struct {
			MS int `json:"ms"`
		}
		_ = json.Unmarshal(raw, &v)
		s.SetLatency(time.Duration(v.MS) * time.Millisecond)
	case "seed":
		var v struct {
			Count int `json:"count"`
		}
		_ = json.Unmarshal(raw, &v)
		s.Seed(v.Count)
	case "fault":
		var f Fault
		_ = json.Unmarshal(raw, &f)
		s.AddFault(f)
	case "log":
		_ = json.NewEncoder(w).Encode(s.Log())
	case "state":
		s.mu.Lock()
		state := map[string]any{"hg": s.hg, "team": s.team}
		_ = json.NewEncoder(w).Encode(state)
		s.mu.Unlock()
	case "external":
		var v struct {
			Kind  string         `json:"kind"`
			ID    string         `json:"id"`
			Patch map[string]any `json:"patch"`
		}
		_ = json.Unmarshal(raw, &v)
		s.External(v.Kind, v.ID, v.Patch)
	case "delete":
		var v struct {
			Kind string `json:"kind"`
			ID   string `json:"id"`
		}
		_ = json.Unmarshal(raw, &v)
		s.Delete(v.Kind, v.ID)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/__ctl/") {
		s.serveControl(w, r)
		return
	}
	s.serveAPI(w, r)
}
