package webexops

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
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
	identifier := regexp.MustCompile(`^[A-Z][A-Za-z0-9_]*$`)
	total := 0
	for _, area := range Areas() {
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
}

func TestOrgParametersThatTakeTheUUIDAreMarked(t *testing.T) {
	for _, op := range append(Operations("cc"), AdminOrgContactsList) {
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
	got := op.Request(testHosts, Args{Params: map[string]string{"orgid": testB64, "orgId": testB64}}).URL
	if got != "https://cc/organization/"+testUUID+"/team?orgId="+testUUID {
		t.Fatalf("URL %s", got)
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
	r := op.Request(testHosts, Args{
		Params: map[string]string{"orgId": testB64, "userId": "a b", "max": "5", "TrackingID": "t"},
		Body:   []byte(`{"a":1}`),
	})
	if r.URL != "https://i/identity/scim/"+testUUID+"/v2/Users/a%20b?max=5" {
		t.Fatalf("URL %s", r.URL)
	}
	if r.Method != "PUT" || string(r.Body) != `{"a":1}` || r.Header.Get("Accept") != "*/*" || r.Header.Get("TrackingID") != "t" {
		t.Fatalf("%+v", r)
	}
	if noBody := (&Operation{Method: "GET", Path: "/x"}).Request(testHosts, Args{Body: []byte("{}")}); noBody.Body != nil {
		t.Fatal("an operation without a body sent one")
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
	if err == nil || strings.Contains(err.Error(), "secret") {
		t.Fatalf("err %v", err)
	}
	if _, err := All(context.Background(), &webexapi.Client{}, testHosts, op, Args{}); err == nil {
		t.Fatal("All read an operation that is not a list")
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
