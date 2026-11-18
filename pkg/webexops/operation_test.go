package webexops

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

const (
	testUUID = "6ab1d2c4-0000-4000-8000-00000000c1d0"
	testB64  = "Y2lzY29zcGFyazovL3VzL09SR0FOSVpBVElPTi82YWIxZDJjNC0wMDAwLTQwMDAtODAwMC0wMDAwMDAwMGMxZDA"
)

var testHosts = webexapi.Hosts{Webex: "https://w/v1", Identity: "https://i", CC: "https://cc", AIAgent: "https://ai"}

func TestTheTableIsConsistent(t *testing.T) {
	ids := map[string]bool{}
	identifier := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
	total := 0
	exported := map[string]bool{}
	for _, area := range Areas() {
		// An exported name is the operation's area, method and path.
		for _, op := range Operations(area) {
			if op.ID[0] < 'A' || op.ID[0] > 'Z' {
				continue
			}
			route := op.Method + " " + op.Path
			if exported[area+" "+route] {
				t.Errorf("%s: %s is exported twice", area, route)
			}
			exported[area+" "+route] = true
			prefix := strings.ToUpper(area[:1]) + area[1:] + strings.ToUpper(op.Method[:1]) + strings.ToLower(op.Method[1:])
			if !strings.HasPrefix(op.ID, prefix) {
				t.Errorf("%s does not start with %s", op.ID, prefix)
			}
		}
		groups := map[string]bool{}
		for _, g := range Groups(area) {
			if g.Area != area || groups[g.Name] {
				t.Errorf("%s: group %+v", area, g)
			}
			groups[g.Name] = true
		}
		for _, op := range Operations(area) {
			total++
			if ids[op.ID] || !identifier.MatchString(op.ID) || op.Area != area || !groups[op.Group] {
				t.Errorf("%s: operation %s in group %q", area, op.ID, op.Group)
			}
			ids[op.ID] = true
			for _, p := range op.Params {
				if (p.Value == "") == (p.Flag == "") || p.Name == "" {
					t.Errorf("%s: parameter %+v needs a name and either a value or a flag", op.ID, p)
				}
				if p.In == InPath && !strings.Contains(op.Path, "{"+p.Name+"}") {
					t.Errorf("%s: path %s has no {%s}", op.ID, op.Path, p.Name)
				}
			}
		}
	}
	if total < 2000 {
		t.Fatalf("only %d operations", total)
	}
}

func TestListingsAreCopies(t *testing.T) {
	area := Areas()[0]
	ops, groups := Operations(area), Groups(area)
	first, firstGroup := ops[0], groups[0]
	ops[0], groups[0] = nil, Group{}
	if Operations(area)[0] != first || Groups(area)[0].Name != firstGroup.Name {
		t.Fatal("changing a listing changed the table")
	}
	for _, g := range Groups("cc") {
		if g.Name == "flow" {
			g.Aliases[0] = "changed"
		}
	}
	for _, g := range Groups("cc") {
		if g.Name == "flow" && g.Aliases[0] != "flows" {
			t.Fatalf("changing a listed alias changed the table: %v", g.Aliases)
		}
	}
}

func TestParametersHaveTheirAPINames(t *testing.T) {
	var org *Param
	for i, p := range AdminGetContactsOrganizationsByOrgIdContactsByContactId.Params {
		if p.In == InPath && p.Flag == "orgid" {
			org = &AdminGetContactsOrganizationsByOrgIdContactsByContactId.Params[i]
		}
	}
	if org == nil || org.Name != "orgId" || !strings.Contains(AdminGetContactsOrganizationsByOrgIdContactsByContactId.Path, "{orgId}") {
		t.Fatalf("org parameter %+v in %s", org, AdminGetContactsOrganizationsByOrgIdContactsByContactId.Path)
	}
}

func request(t *testing.T, op *Operation, args Args) webexapi.Request {
	t.Helper()
	r, err := op.Request(testHosts, args)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestRequestRefusesArgumentsItCannotSend(t *testing.T) {
	op := &Operation{ID: "Test", Method: "PUT", Path: "/locations/{locationId}/parks/{parkId}", Params: []Param{
		{Name: "locationId", In: InPath, Flag: "location-id"},
		{Name: "parkId", In: InPath, Flag: "park-id"},
		{Name: "projectId", In: InPath, Value: "fixed"},
		{Name: "orgId", In: InQuery, Flag: "org-id"},
	}}
	good := map[string]string{"locationId": "l", "parkId": "p"}
	with := func(k, v string) map[string]string {
		m := map[string]string{"locationId": "l", "parkId": "p"}
		m[k] = v
		return m
	}
	for name, args := range map[string]Args{
		"a misspelled parameter":            {Params: with("orgID", "o")},
		"a query parameter it doesn't have": {Params: good, Query: url.Values{"nope": {"x"}}},
		"a value for a fixed parameter":     {Params: with("projectId", "other")},
		"a missing path value":              {Params: map[string]string{"locationId": "l"}},
	} {
		if _, err := op.Request(testHosts, args); err == nil {
			t.Errorf("%s was accepted", name)
		}
		if _, err := op.List(testHosts, args); err == nil {
			t.Errorf("List: %s was accepted", name)
		}
	}
	r, err := op.Request(testHosts, Args{Params: map[string]string{"locationId": "l"}, AllowEmptyPath: true})
	if err != nil || !strings.HasSuffix(r.URL, "/locations/l/parks/{parkId}") {
		t.Fatalf("AllowEmptyPath: %v, %v", r.URL, err)
	}
}

func TestARouteListedWithDifferentRequestsTakesEveryListingsParameters(t *testing.T) {
	r := request(t, CallingGetWorkspacesByWorkspaceId, Args{Params: map[string]string{
		"workspaceId": "w", "includeDevices": "true", "includeCapabilities": "true",
	}})
	if r.URL != "https://w/v1/workspaces/w?includeCapabilities=true&includeDevices=true" {
		t.Fatalf("URL %s", r.URL)
	}
	for _, op := range Operations("calling") {
		if op == CallingGetWorkspacesByWorkspaceId {
			t.Fatal("the combined operation is listed as a CLI command")
		}
	}
}

func TestCallAndAllRedactAIAgentResponses(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"objects":[{"id":"a","bot_access_token":"secret"}],"meta":{"total_count":1}}`)
	}))
	defer srv.Close()
	op := &Operation{ID: "Test", Method: "GET", Host: HostAIAgent, Path: "/api/v2/agents/", Redact: true,
		Paging: &webexapi.Paging{Mode: webexapi.Offset, Start: "offset", Size: "limit", Items: "objects", TotalItems: "meta.total_count"}}
	hosts := webexapi.Hosts{AIAgent: srv.URL}
	resp, err := Call(context.Background(), &webexapi.Client{}, hosts, op, Args{})
	if err != nil || strings.Contains(string(resp.Body), "secret") {
		t.Fatalf("Call: %s, %v", resp.Body, err)
	}
	items, err := All(context.Background(), &webexapi.Client{}, hosts, op, Args{})
	if err != nil || len(items) != 1 || strings.Contains(string(items[0]), "secret") {
		t.Fatalf("All: %s, %v", items, err)
	}
}

func TestAQueryParameterCanBeSentMoreThanOnce(t *testing.T) {
	op := &Operation{Method: "GET", Path: "/adminAudit/events", Params: []Param{
		{Name: "eventCategories", In: InQuery, Flag: "event-categories"},
		{Name: "max", In: InQuery, Flag: "max"},
	}}
	got := request(t, op, Args{
		Params: map[string]string{"eventCategories": "ignored", "max": "5"},
		Query:  url.Values{"eventCategories": {"a", "b"}},
	}).URL
	if got != "https://w/v1/adminAudit/events?eventCategories=a&eventCategories=b&max=5" {
		t.Fatalf("URL %s", got)
	}
}

func TestOrgParametersThatTakeTheUUIDAreMarked(t *testing.T) {
	for _, op := range append(Operations("cc"), AdminGetContactsOrganizationsByOrgIdContactsSearch) {
		for _, p := range op.Params {
			if strings.EqualFold(p.Name, "orgid") && p.Flag == "orgid" && !p.OrgUUID {
				t.Errorf("%s: %s parameter %s is not marked OrgUUID", op.ID, p.In, p.Name)
			}
		}
	}
	op := &Operation{Method: "GET", Host: HostCC, Path: "/organization/{orgid}/team", Params: []Param{
		{Name: "orgid", In: InPath, Flag: "orgid", OrgUUID: true},
		{Name: "orgId", In: InQuery, Flag: "orgid", OrgUUID: true},
	}}
	got := request(t, op, Args{Params: map[string]string{"orgid": testB64, "orgId": testB64}}).URL
	if got != "https://cc/organization/"+testUUID+"/team?orgId="+testUUID {
		t.Fatalf("URL %s", got)
	}
	query := url.Values{"orgId": {testB64}}
	got = request(t, op, Args{Params: map[string]string{"orgid": testUUID}, Query: query}).URL
	if got != "https://cc/organization/"+testUUID+"/team?orgId="+testUUID || query["orgId"][0] != testB64 {
		t.Fatalf("URL %s, caller's query %v", got, query)
	}
}

func TestRequestFillsTheOperation(t *testing.T) {
	op := &Operation{
		Method: "PUT", Host: HostIdentity, Path: "/identity/scim/{orgId}/v2/Users/{userId}", Accept: "*/*",
		Params: []Param{
			{Name: "orgId", In: InPath, Flag: "org-id", OrgUUID: true},
			{Name: "userId", In: InPath, Flag: "user-id"},
			{Name: "projectId", In: InPath, Value: "fixed"},
			{Name: "max", In: InQuery, Flag: "max"},
			{Name: "empty", In: InQuery, Flag: "empty"},
			{Name: "TrackingID", In: InHeader, Flag: "tracking-id"},
		},
		Body: &Body{},
	}
	r := request(t, op, Args{
		Params: map[string]string{"orgId": testB64, "userId": "a b", "max": "5", "TrackingID": "t"},
		Body:   []byte(`{"a":1}`),
	})
	if r.URL != "https://i/identity/scim/"+testUUID+"/v2/Users/a%20b?max=5" {
		t.Fatalf("URL %s", r.URL)
	}
	if r.Method != "PUT" || string(r.Body) != `{"a":1}` || r.Header.Get("Accept") != "*/*" || r.Header.Get("TrackingID") != "t" {
		t.Fatalf("%+v", r)
	}
	// An upload the collection gives no JSON body for still sends the body given.
	upload := request(t, DevicePostTelephonyConfigDevicesByDeviceIdActionsBackgroundImageUploadInvoke, Args{
		Params: map[string]string{"deviceId": "d"}, Body: []byte("--b\r\n"), ContentType: "multipart/form-data; boundary=b",
	})
	if DevicePostTelephonyConfigDevicesByDeviceIdActionsBackgroundImageUploadInvoke.Body != nil || string(upload.Body) != "--b\r\n" || upload.ContentType != "multipart/form-data; boundary=b" {
		t.Fatalf("upload %+v", upload)
	}
}

func TestCallRedactsTheErrorsOfAIAgentOperations(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("bot-access-token") != "token" {
			t.Errorf("headers %v", r.Header)
		}
		w.WriteHeader(400)
		io.WriteString(w, `{"bot_access_token":"secret"}`)
	}))
	defer srv.Close()
	op := &Operation{ID: "Test", Method: "GET", Host: HostAIAgent, Path: "/api/v1/bot/", Redact: true,
		Params: []Param{{Name: "bot-access-token", In: InAgentToken, Flag: "agent-id"}}}
	_, err := Call(context.Background(), &webexapi.Client{}, webexapi.Hosts{AIAgent: srv.URL}, op,
		Args{Params: map[string]string{"bot-access-token": "token"}})
	var apiErr *webexapi.APIError
	if !errors.As(err, &apiErr) || apiErr.Status != 400 || strings.Contains(err.Error(), "secret") {
		t.Fatalf("err %v", err)
	}
	if _, err := All(context.Background(), &webexapi.Client{}, testHosts, op, Args{}); err == nil {
		t.Fatal("All read an operation that is not a list")
	}
}

func TestCallAndAllRefuseAMissingPathParameter(t *testing.T) {
	sent := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { sent++ }))
	defer srv.Close()
	hosts := webexapi.Hosts{CC: srv.URL}
	_, err := Call(context.Background(), &webexapi.Client{}, hosts, CcGetOrganizationByOrgidV2Team, Args{})
	if err == nil || !strings.Contains(err.Error(), "orgid") {
		t.Fatalf("Call: err %v", err)
	}
	if _, err := All(context.Background(), &webexapi.Client{}, hosts, CcGetOrganizationByOrgidV2Team, Args{}); err == nil {
		t.Fatal("All sent a request without its path parameter")
	}
	if sent != 0 {
		t.Fatalf("%d requests were sent", sent)
	}
}

func TestAllReadsAListToItsEnd(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") == "0" {
			io.WriteString(w, `{"data":[{"id":"1"}],"meta":{"totalPages":2}}`)
			return
		}
		io.WriteString(w, `{"data":[{"id":"2"}],"meta":{"totalPages":2}}`)
	}))
	defer srv.Close()
	op := &Operation{ID: "Test", Method: "GET", Host: HostCC, Path: "/organization/{orgid}/team",
		Params: []Param{{Name: "orgid", In: InPath, Flag: "orgid"}},
		Paging: &webexapi.Paging{Mode: webexapi.Pages, Page: "page", Size: "pageSize", Items: "data", TotalPages: "meta.totalPages"}}
	items, err := All(context.Background(), &webexapi.Client{}, webexapi.Hosts{CC: srv.URL}, op, Args{Params: map[string]string{"orgid": testUUID}})
	if err != nil || len(items) != 2 {
		t.Fatalf("items %s, err %v", items, err)
	}
}
