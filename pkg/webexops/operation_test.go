package webexops

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

func TestEveryOperationIsWellFormed(t *testing.T) {
	seen := map[string]bool{}
	for _, area := range Areas {
		groups := map[string]bool{}
		for _, g := range Groups(area) {
			groups[g.Name] = true
		}
		ops := Operations(area)
		if len(ops) == 0 {
			t.Fatalf("%s has no operations", area)
		}
		for _, op := range ops {
			if op.Area != area || !groups[op.Group] {
				t.Errorf("%s: area or group is not listed", op.Name())
			}
			for _, p := range op.Params {
				if p.In == InPath && !strings.Contains(op.Path, "{"+p.Name+"}") {
					t.Errorf("%s: path %s has no {%s}", op.Name(), op.Path, p.Name)
				}
				if (p.Flag == "") == (p.Value == "") {
					t.Errorf("%s: parameter %s needs exactly one of a flag or a fixed value", op.Name(), p.Name)
				}
			}
			seen[op.Name()] = true
		}
	}
	if len(seen) < 1500 {
		t.Fatalf("only %d distinct operations", len(seen))
	}
}

func TestURLFillsPathAndQueryAndRefusesUnknownArguments(t *testing.T) {
	hosts, _ := webexapi.NewHosts("eu1", "")
	op, ok := Find("cc", "flow", "list")
	if !ok {
		t.Fatal("cc flow list is missing")
	}
	u, err := op.URL(hosts, Args{Path: map[string]string{"orgid": "org 1"}, Query: map[string]string{"size": "10", "page": ""}})
	if err != nil {
		t.Fatal(err)
	}
	if u != "https://api.wxcc-eu1.cisco.com/org%201/project/5e5c9ad6d61f870d6d778c1b/flows?size=10" {
		t.Fatalf("url = %s", u)
	}
	if _, err := op.URL(hosts, Args{}); err == nil || !strings.Contains(err.Error(), "orgid is required") {
		t.Fatalf("missing path parameter: %v", err)
	}
	if _, err := op.URL(hosts, Args{Path: map[string]string{"orgid": "o"}, Query: map[string]string{"sise": "1"}}); err == nil {
		t.Fatal("misspelled query parameter accepted")
	}
	if _, err := op.URL(hosts, Args{Path: map[string]string{"orgid": "o"}, Header: map[string]string{"Trakking-Id": "x"}}); err == nil {
		t.Fatal("undeclared header accepted")
	}
	if _, err := MeetingsRecordingsListGroup.URL(hosts, Args{Header: map[string]string{"HostEmail": "a@example.com"}}); err != nil {
		t.Fatalf("a declared header, in another case: %v", err)
	}
}

func TestListPagesAsTheOperationSays(t *testing.T) {
	var queries []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.RawQuery)
		io.WriteString(w, `{"data":[{"id":"f1"}],"pageInfo":{"totalPages":1}}`)
	}))
	defer srv.Close()
	op, _ := Find("cc", "flow", "list")
	var got []json.RawMessage
	err := List(context.Background(), &webexapi.Client{}, webexapi.Hosts{CC: srv.URL}, op,
		Args{Path: map[string]string{"orgid": "o1"}, Query: map[string]string{"flowType": "FLOW"}},
		func(page []json.RawMessage) error { got = append(got, page...); return nil })
	if err != nil || len(got) != 1 {
		t.Fatalf("got %d, err %v", len(got), err)
	}
	if queries[0] != "flowType=FLOW&includePagination=true&page=0&size=100" {
		t.Fatalf("query = %s", queries[0])
	}
	get, _ := Find("cc", "team", "get")
	if err := List(context.Background(), &webexapi.Client{}, webexapi.Hosts{}, get, Args{}, nil); err == nil {
		t.Fatal("a read that does not page was listed")
	}
}

func TestListSendsTheCallersHeadersOnEveryPage(t *testing.T) {
	var hostEmails []string
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hostEmails = append(hostEmails, r.Header.Get("hostEmail"))
		if len(hostEmails) > 3 { // stops a pager that loops
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.URL.Query().Get("cursor") == "" {
			w.Header().Set("Link", `<`+srv.URL+`/group/recordings?cursor=2>; rel="next"`)
		}
		io.WriteString(w, `{"items":[{"id":"`+r.URL.Query().Get("cursor")+`"}]}`)
	}))
	defer srv.Close()
	err := List(context.Background(), &webexapi.Client{}, webexapi.Hosts{Webex: srv.URL}, MeetingsRecordingsListGroup,
		Args{Header: map[string]string{"hostEmail": "a@example.com"}}, func([]json.RawMessage) error { return nil })
	if err != nil || len(hostEmails) != 2 || hostEmails[0] != "a@example.com" || hostEmails[1] != "a@example.com" {
		t.Fatalf("hostEmail per page = %q, err %v", hostEmails, err)
	}
}

func TestCallSendsTheOperation(t *testing.T) {
	var got *http.Request
	var body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		io.WriteString(w, `{"id":"t1"}`)
	}))
	defer srv.Close()

	op, ok := Find("cc", "team", "create")
	if !ok {
		t.Fatal("cc team create is missing")
	}
	hosts := webexapi.Hosts{CC: srv.URL}
	resp, err := Call(context.Background(), &webexapi.Client{}, hosts, op, Args{
		Path: map[string]string{"orgid": "o1"},
		Body: []byte(`{"name":"zz"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Method != http.MethodPost || got.URL.Path != "/organization/o1/team" || body != `{"name":"zz"}` || string(resp.Body) != `{"id":"t1"}` {
		t.Fatalf("%s %s %s → %s", got.Method, got.URL.Path, body, resp.Body)
	}
}
