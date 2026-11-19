package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// searchServer answers search queries with handle, given the GraphQL query.
func searchServer(t *testing.T, handle func(query string) (int, string)) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := io.ReadAll(r.Body)
		status, body := handle(graphQL(t, string(data)))
		w.WriteHeader(status)
		io.WriteString(w, body)
	}))
	t.Cleanup(srv.Close)
	saved := config.CcBaseURL
	config.CcBaseURL = srv.URL
	t.Cleanup(func() { config.CcBaseURL = saved })
}

var flowPage = regexp.MustCompile(`pageSize: (\d+), currentPage: (\d+)`)

// tracePage is a page of trace events as Webex returns it: pages start at 0,
// and totalPages is the record count divided by the page size, rounded down.
func tracePage(t *testing.T, query string, records int) string {
	t.Helper()
	m := flowPage.FindStringSubmatch(query)
	if m == nil {
		t.Fatalf("query has no page: %s", query)
	}
	size, _ := strconv.Atoi(m[1])
	page, _ := strconv.Atoi(m[2])
	traces := []any{}
	for i := page * size; i < min((page+1)*size, records); i++ {
		traces = append(traces, map[string]any{"activityProcessId": fmt.Sprint(i)})
	}
	body, _ := json.Marshal(map[string]any{"data": map[string]any{"flowTraceEvents": map[string]any{
		"traces":   traces,
		"pageInfo": map[string]any{"pageSize": size, "currentPage": page, "totalPages": records / size, "totalRecords": records},
	}}})
	return string(body)
}

func TestFlowPagingReadsEveryRecord(t *testing.T) {
	var pages []string
	searchServer(t, func(q string) (int, string) {
		pages = append(pages, flowPage.FindStringSubmatch(q)[2])
		return 200, tracePage(t, q, 5)
	})

	out, _, err := PaginateGraphQL(QueryFlowTraceEvents, Params{From: 1, To: 2, PageSize: 2}, "org", "")
	if err != nil {
		t.Fatal(err)
	}
	var items []map[string]string
	if err := json.Unmarshal(out, &items); err != nil {
		t.Fatal(err)
	}
	if len(items) != 5 || items[0]["activityProcessId"] != "0" || items[4]["activityProcessId"] != "4" {
		t.Errorf("items = %v, want records 0 to 4", items)
	}
	if fmt.Sprint(pages) != "[0 1 2]" {
		t.Errorf("pages read = %v, want [0 1 2]", pages)
	}
}

func TestFlowPagingFailsOnAnEmptyPageBeforeTheTotal(t *testing.T) {
	searchServer(t, func(q string) (int, string) {
		if flowPage.FindStringSubmatch(q)[2] == "1" {
			return 200, `{"data":{"flowTraceEvents":{"traces":[],"pageInfo":{"totalRecords":5}}}}`
		}
		return 200, tracePage(t, q, 5)
	})
	_, _, err := PaginateGraphQL(QueryFlowTraceEvents, Params{From: 1, To: 2, PageSize: 2}, "org", "")
	if !errors.Is(err, webexapi.ErrIncompleteList) {
		t.Errorf("err = %v, want incomplete list", err)
	}
}

func TestPagingFailsWhenALaterPageFails(t *testing.T) {
	searchServer(t, func(q string) (int, string) {
		if flowPage.FindStringSubmatch(q)[2] == "1" {
			return 500, `{"error":"down"}`
		}
		return 200, tracePage(t, q, 5)
	})
	_, _, err := PaginateGraphQL(QueryFlowTraceEvents, Params{From: 1, To: 2, PageSize: 2}, "org", "")
	if !errors.Is(err, webexapi.ErrIncompleteList) {
		t.Errorf("err = %v, want incomplete list", err)
	}
}

func TestPagingFailsOnGraphQLErrors(t *testing.T) {
	searchServer(t, func(string) (int, string) {
		return 200, `{"data":null,"errors":[{"message":"Field 'x' is undefined"}]}`
	})
	_, _, err := PaginateGraphQL(QueryTask, Params{From: 1, To: 2}, "org", "")
	if err == nil {
		t.Fatal("GraphQL errors read as an empty list")
	}
}

func TestCursorPagingFailsWithoutACursor(t *testing.T) {
	searchServer(t, func(string) (int, string) {
		return 200, `{"data":{"task":{"tasks":[{"id":"t1"}],"pageInfo":{"hasNextPage":true,"endCursor":""}}}}`
	})
	_, _, err := PaginateGraphQL(QueryTask, Params{From: 1, To: 2}, "org", "")
	if !errors.Is(err, webexapi.ErrIncompleteList) {
		t.Errorf("err = %v, want incomplete list", err)
	}
}

func TestCursorPagingFailsWithoutHasNextPage(t *testing.T) {
	searchServer(t, func(string) (int, string) {
		return 200, `{"data":{"task":{"tasks":[{"id":"t1"}],"pageInfo":{"endCursor":"c1"}}}}`
	})
	_, _, err := PaginateGraphQL(QueryTask, Params{From: 1, To: 2}, "org", "")
	if !errors.Is(err, webexapi.ErrIncompleteList) {
		t.Errorf("err = %v, want incomplete list", err)
	}
}

func TestFlowPagingRefusesANegativePageSize(t *testing.T) {
	var requests int
	searchServer(t, func(q string) (int, string) {
		requests++
		return 200, tracePage(t, q, 5)
	})
	if _, _, err := PaginateGraphQL(QueryFlowTraceEvents, Params{From: 1, To: 2, PageSize: -1}, "org", ""); err == nil {
		t.Error("a negative page size was accepted")
	}
	if requests != 0 {
		t.Errorf("%d requests sent", requests)
	}
}

func TestCursorPagingFollowsTheCursor(t *testing.T) {
	cursor := regexp.MustCompile(`cursor: "([^"]*)"`)
	searchServer(t, func(q string) (int, string) {
		if m := cursor.FindStringSubmatch(q); m != nil && m[1] == "c1" {
			return 200, `{"data":{"task":{"tasks":[{"id":"t2"}],"pageInfo":{"hasNextPage":false}}}}`
		}
		return 200, `{"data":{"task":{"tasks":[{"id":"t1"}],"pageInfo":{"hasNextPage":true,"endCursor":"c1"}}}}`
	})
	out, _, err := PaginateGraphQL(QueryTask, Params{From: 1, To: 2}, "org", "")
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `[{"id":"t1"},{"id":"t2"}]` {
		t.Errorf("items = %s", out)
	}
}
