package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func names(t *testing.T, body []byte) []string {
	t.Helper()
	var items []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &items); err != nil {
		t.Fatalf("decoding merged items: %v (%s)", err, body)
	}
	out := make([]string, len(items))
	for i, it := range items {
		out[i] = it.Name
	}
	return out
}

func named(from, to int) []map[string]string {
	items := []map[string]string{}
	for i := from; i < to; i++ {
		items = append(items, map[string]string{"name": fmt.Sprintf("n%d", i)})
	}
	return items
}

// server serves handle and fails the test after 50 requests, so a pager that
// never stops ends.
func server(t *testing.T, handle func(w http.ResponseWriter, r *http.Request)) (*httptest.Server, *[]string) {
	t.Helper()
	var queries []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.Path+"?"+r.URL.RawQuery)
		if len(queries) > 50 {
			t.Errorf("more than 50 requests")
			http.Error(w, "too many", http.StatusTeapot)
			return
		}
		handle(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv, &queries
}

func atoi(s string) int { n, _ := strconv.Atoi(s); return n }

func TestLinkFollowsNextUntilThereIsNone(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		start := atoi(r.URL.Query().Get("cursor"))
		size := atoi(r.URL.Query().Get("max"))
		if end := start + size; end < 250 {
			w.Header().Set("Link", fmt.Sprintf(`<http://%s/people?max=%d&cursor=%d>; rel="next"`, r.Host, size, end))
		}
		json.NewEncoder(w).Encode(map[string]any{"items": named(start, min(start+size, 250))})
	})
	body, _, err := NewRequest(srv.URL, "GET", "/people").DoPaginated(Paging{Mode: Link, Size: "max", Items: "items", OmitEmpty: true})
	if err != nil {
		t.Fatal(err)
	}
	if got := names(t, body); len(got) != 250 || got[249] != "n249" {
		t.Fatalf("got %d items, want n0..n249", len(got))
	}
	if want := "[/people?max=100 /people?max=100&cursor=100 /people?max=100&cursor=200]"; fmt.Sprint(*queries) != want {
		t.Fatalf("requests %v, want %s", *queries, want)
	}
}

func TestLinkReadsFromTheStartWhateverStartTheCallerSet(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"items": named(0, 3)})
	})
	req := NewRequest(srv.URL, "GET", "/huntGroups")
	req.QueryParam("start", "100")
	req.QueryParam("name", "x")
	if _, _, err := req.DoPaginated(Paging{Mode: Link, Start: "start", Size: "max", Items: "items", OmitEmpty: true}); err != nil {
		t.Fatal(err)
	}
	if want := "[/huntGroups?max=100&name=x]"; fmt.Sprint(*queries) != want {
		t.Fatalf("requests %v, want %s", *queries, want)
	}
}

func TestLinkResolvesARelativeNext(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") == "" {
			w.Header().Set("Link", `</v2/site?page=1>;rel="self",</v2/site?page=1>;rel="next"`)
			json.NewEncoder(w).Encode(map[string]any{"data": named(0, 1)})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"data": named(1, 2)})
	})
	body, _, err := NewRequest(srv.URL, "GET", "/v2/site").DoPaginated(Paging{Mode: Link, Items: "data"})
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprint(names(t, body)); got != "[n0 n1]" {
		t.Fatalf("got %s", got)
	}
}

func TestLinkRefusesANextPageOnAnotherHost(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", `<https://elsewhere.example/people?cursor=x>; rel="next"`)
		json.NewEncoder(w).Encode(map[string]any{"items": named(0, 1)})
	})
	_, _, err := NewRequest(srv.URL, "GET", "/people").DoPaginated(Paging{Mode: Link, Size: "max", Items: "items"})
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v, want ErrIncompleteList", err)
	}
	if len(*queries) != 1 {
		t.Fatalf("%d requests, want 1", len(*queries))
	}
}

func TestLinkFailsWhenANextLinkRepeats(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", fmt.Sprintf(`<http://%s/people?cursor=same>; rel="next"`, r.Host))
		json.NewEncoder(w).Encode(map[string]any{"items": named(0, 1)})
	})
	_, _, err := NewRequest(srv.URL, "GET", "/people").DoPaginated(Paging{Mode: Link, Items: "items"})
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v, want ErrIncompleteList", err)
	}
}

// An endpoint that links past its last page would otherwise be read forever.
func TestLinkFailsOnANextLinkAfterAnEmptyPage(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		start := atoi(r.URL.Query().Get("start"))
		w.Header().Set("Link", fmt.Sprintf(`<http://%s/x?max=1&start=%d>;rel="next"`, r.Host, start+1))
		json.NewEncoder(w).Encode(map[string]any{"items": named(start, min(start+1, 1))})
	})
	_, _, err := NewRequest(srv.URL, "GET", "/x").DoPaginated(Paging{Mode: Link, Size: "max", Items: "items", OmitEmpty: true})
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v, want ErrIncompleteList", err)
	}
	if len(*queries) != 2 {
		t.Fatalf("%d requests, want 2", len(*queries))
	}
}

// The next link of some endpoints has lost its path; its query still pages the
// request path.
func TestLinkQueryOnlyAppliesTheNextQueryToTheRequestPath(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/resourceGroup/memberships/v2" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("cursor") == "" {
			w.Header().Set("Link", fmt.Sprintf(`<http://%s?max=100&cursor=c2>; rel="next"`, r.Host))
			json.NewEncoder(w).Encode(map[string]any{"items": named(0, 1)})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"items": named(1, 2)})
	})
	body, _, err := NewRequest(srv.URL, "GET", "/resourceGroup/memberships/v2").DoPaginated(Paging{Mode: Link, Size: "max", Items: "items", LinkQueryOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprint(names(t, body)); got != "[n0 n1]" {
		t.Fatalf("got %s (requests %v)", got, *queries)
	}
}

func TestAnEmptyListMayLeaveOutItsItemsOnlyWhenOmitEmpty(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{}`)) })
	body, _, err := NewRequest(srv.URL, "GET", "/x").DoPaginated(Paging{Mode: Link, Size: "max", Items: "items", OmitEmpty: true})
	if err != nil || string(body) != "[]" {
		t.Fatalf("OmitEmpty: body %s, err %v; want []", body, err)
	}
	if _, _, err := NewRequest(srv.URL, "GET", "/x").DoPaginated(Paging{Mode: Link, Size: "max", Items: "items"}); !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("without OmitEmpty: err = %v, want ErrIncompleteList", err)
	}
}

func TestItemsThatAreNotAnArrayFail(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"items": {"a": 1}}`)) })
	if _, _, err := NewRequest(srv.URL, "GET", "/x").DoPaginated(Paging{Mode: Single, Items: "items", OmitEmpty: true}); !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v, want ErrIncompleteList", err)
	}
}

func TestPagesEndAtTotalPages(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		page, size := atoi(r.URL.Query().Get("page")), atoi(r.URL.Query().Get("pageSize"))
		json.NewEncoder(w).Encode(map[string]any{
			"data": named(page*size, min((page+1)*size, 250)),
			"meta": map[string]int{"page": page, "totalPages": 3},
		})
	})
	body, _, err := NewRequest(srv.URL, "GET", "/site").DoPaginated(Paging{Mode: Pages, Page: "page", Size: "pageSize", Items: "data", TotalPages: "meta.totalPages"})
	if err != nil {
		t.Fatal(err)
	}
	if got := names(t, body); len(got) != 250 {
		t.Fatalf("got %d items, want 250", len(got))
	}
	if want := "[/site?page=0&pageSize=100 /site?page=1&pageSize=100 /site?page=2&pageSize=100]"; fmt.Sprint(*queries) != want {
		t.Fatalf("requests %v, want %s", *queries, want)
	}
}

func TestPagesFailWithoutTheirTotal(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"data": named(0, 100), "meta": map[string]string{"orgId": "o"}})
	})
	_, _, err := NewRequest(srv.URL, "GET", "/site").DoPaginated(Paging{Mode: Pages, Page: "page", Size: "pageSize", Items: "data", TotalPages: "meta.totalPages"})
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v, want ErrIncompleteList", err)
	}
}

func TestPagesFailOnAnEmptyPageBeforeTheirTotal(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) {
		data := named(0, 100)
		if r.URL.Query().Get("page") != "0" {
			data = nil
		}
		json.NewEncoder(w).Encode(map[string]any{"data": data, "meta": map[string]int{"totalPages": 3}})
	})
	_, _, err := NewRequest(srv.URL, "GET", "/site").DoPaginated(Paging{Mode: Pages, Page: "page", Size: "pageSize", Items: "data", TotalPages: "meta.totalPages"})
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v, want ErrIncompleteList", err)
	}
}

func TestPagesEndAtTotalItems(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		page, size := atoi(r.URL.Query().Get("page")), atoi(r.URL.Query().Get("pageSize"))
		json.NewEncoder(w).Encode(map[string]any{"resources": named(page*size, min((page+1)*size, 150)), "totalResources": 150})
	})
	body, _, err := NewRequest(srv.URL, "GET", "/site/bulk-export").DoPaginated(Paging{Mode: Pages, Page: "page", Size: "pageSize", Items: "resources", TotalItems: "totalResources"})
	if err != nil {
		t.Fatal(err)
	}
	if got := names(t, body); len(got) != 150 || len(*queries) != 2 {
		t.Fatalf("got %d items in %d requests, want 150 in 2", len(got), len(*queries))
	}
}

func TestNotFoundEndsAListAfterItsFirstPage(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		page := atoi(r.URL.Query().Get("page"))
		if page > 1 {
			http.Error(w, `{"message":"Resource Not Found"}`, http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"data": named(page*100, (page+1)*100)})
	})
	p := Paging{Mode: Pages, Page: "page", Size: "pageSize", Items: "data", NotFoundEnds: true}
	body, _, err := NewRequest(srv.URL, "GET", "/person").DoPaginated(p)
	if err != nil {
		t.Fatal(err)
	}
	if got := names(t, body); len(got) != 200 || len(*queries) != 3 {
		t.Fatalf("got %d items in %d requests, want 200 in 3", len(got), len(*queries))
	}
}

func TestNotFoundOnTheFirstPageIsAnError(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"message":"Resource Not Found"}`, http.StatusNotFound)
	})
	p := Paging{Mode: Pages, Page: "page", Size: "pageSize", Items: "data", NotFoundEnds: true}
	if _, _, err := NewRequest(srv.URL, "GET", "/person").DoPaginated(p); err == nil || !strings.Contains(err.Error(), "404") {
		t.Fatalf("err = %v, want the 404", err)
	}
}

func TestOffsetReadsToTotalItemsFromItsFirstIndex(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		start, count := atoi(r.URL.Query().Get("startIndex")), atoi(r.URL.Query().Get("count"))
		count = min(count, 50) // the server's own page size
		json.NewEncoder(w).Encode(map[string]any{"totalResults": 120, "Resources": named(start-1, min(start-1+count, 120))})
	})
	body, _, err := NewRequest(srv.URL, "GET", "/Users").DoPaginated(Paging{Mode: Offset, Start: "startIndex", First: 1, Size: "count", Items: "Resources", TotalItems: "totalResults"})
	if err != nil {
		t.Fatal(err)
	}
	if got := names(t, body); len(got) != 120 || got[119] != "n119" {
		t.Fatalf("got %d items, want n0..n119", len(got))
	}
	if want := "[/Users?count=100&startIndex=1 /Users?count=100&startIndex=51 /Users?count=100&startIndex=101]"; fmt.Sprint(*queries) != want {
		t.Fatalf("requests %v, want %s", *queries, want)
	}
}

func TestOffsetFailsWhenTheListEndsBeforeItsTotal(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) {
		items := named(0, 0)
		if r.URL.Query().Get("offset") == "0" {
			items = named(0, 100)
		}
		json.NewEncoder(w).Encode(map[string]any{"objects": items, "meta": map[string]int{"total_count": 150}})
	})
	_, _, err := NewRequest(srv.URL, "GET", "/agents/").DoPaginated(Paging{Mode: Offset, Start: "offset", Size: "limit", Items: "objects", TotalItems: "meta.total_count"})
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v, want ErrIncompleteList", err)
	}
}

func TestOffsetWithoutATotalEndsAtAShortPage(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		start, size := atoi(r.URL.Query().Get("start")), atoi(r.URL.Query().Get("max"))
		json.NewEncoder(w).Encode(map[string]any{"members": named(start, min(start+size, 57))})
	})
	body, _, err := NewRequest(srv.URL, "GET", "/members").DoPaginated(Paging{Mode: Offset, Start: "start", Size: "max", Items: "members", PageSize: 20})
	if err != nil {
		t.Fatal(err)
	}
	if got := names(t, body); len(got) != 57 || len(*queries) != 3 {
		t.Fatalf("got %d items in %d requests, want 57 in 3", len(got), len(*queries))
	}
}

func TestOffsetFailsWhenAPageRepeats(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"members": named(0, 100)})
	})
	_, _, err := NewRequest(srv.URL, "GET", "/members").DoPaginated(Paging{Mode: Offset, Start: "start", Size: "max", Items: "members"})
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v, want ErrIncompleteList", err)
	}
}

func TestSingleSendsOneRequestWithoutPagingParameters(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(named(0, 100))
	})
	req := NewRequest(srv.URL, "GET", "/numbers")
	req.QueryParam("filter", "x")
	body, _, err := req.DoPaginated(Paging{Mode: Single, Items: ""})
	if err != nil {
		t.Fatal(err)
	}
	if got := names(t, body); len(got) != 100 || fmt.Sprint(*queries) != "[/numbers?filter=x]" {
		t.Fatalf("got %d items from %v", len(got), *queries)
	}
}

func TestItemsCanBeNested(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"locationsStatus": map[string]any{"locations": named(0, 2)}})
	})
	body, _, err := NewRequest(srv.URL, "GET", "/x").DoPaginated(Paging{Mode: Single, Items: "locationsStatus.locations"})
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprint(names(t, body)); got != "[n0 n1]" {
		t.Fatalf("got %s", got)
	}
}

func TestPagingKeepsTheCallersOtherQueryParameters(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"data": named(0, 1), "meta": map[string]int{"totalPages": 1}})
	})
	req := NewRequest(srv.URL, "GET", "/site")
	req.QueryParam("filter", "name==x")
	req.QueryParam("page", "7")
	if _, _, err := req.DoPaginated(Paging{Mode: Pages, Page: "page", Size: "pageSize", Items: "data", TotalPages: "meta.totalPages"}); err != nil {
		t.Fatal(err)
	}
	if got := (*queries)[0]; !strings.Contains(got, "filter=name%3D%3Dx") || !strings.Contains(got, "page=0") {
		t.Fatalf("request %s, want the filter kept and page 0", got)
	}
}

func TestPagingSendsTheQueryItNeeds(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"data": named(0, 1), "pageInfo": map[string]int{"totalPages": 1}})
	})
	p := Paging{Mode: Pages, Page: "page", Size: "size", Items: "data", TotalPages: "pageInfo.totalPages", Query: map[string]string{"includePagination": "true"}}
	if _, _, err := NewRequest(srv.URL, "GET", "/flows").DoPaginated(p); err != nil {
		t.Fatal(err)
	}
	if got := (*queries)[0]; got != "/flows?includePagination=true&page=0&size=100" {
		t.Fatalf("request %s", got)
	}
}

func TestAnAPIErrorIsReturned(t *testing.T) {
	srv, _ := server(t, func(w http.ResponseWriter, r *http.Request) { http.Error(w, `{"message":"no"}`, http.StatusForbidden) })
	if _, _, err := NewRequest(srv.URL, "GET", "/x").DoPaginated(Paging{Mode: Single, Items: "items"}); err == nil || !strings.Contains(err.Error(), "403") {
		t.Fatalf("err = %v, want the 403", err)
	}
}

func TestVisitStopsWhenAskedTo(t *testing.T) {
	srv, queries := server(t, func(w http.ResponseWriter, r *http.Request) {
		page := atoi(r.URL.Query().Get("page"))
		json.NewEncoder(w).Encode(map[string]any{"data": named(page*100, (page+1)*100), "meta": map[string]int{"totalPages": 5}})
	})
	var seen int
	err := NewRequest(srv.URL, "GET", "/greeting-purpose").Visit(Paging{Mode: Pages, Page: "page", Size: "pageSize", Items: "data", TotalPages: "meta.totalPages"},
		func(items []json.RawMessage) bool { seen += len(items); return false })
	if err != nil || seen != 100 || len(*queries) != 1 {
		t.Fatalf("err %v, saw %d items in %d requests; want 100 in 1", err, seen, len(*queries))
	}
}
