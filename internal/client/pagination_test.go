package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
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

// flowStore mimics /flows: it pages with page/size (default size 10), ignores
// pageSize, and only returns pageInfo when includePagination=true.
func flowStore(total int, pages *[]string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		*pages = append(*pages, q.Get("page"))
		page, _ := strconv.Atoi(q.Get("page"))
		size := 10
		if s := q.Get("size"); s != "" {
			size, _ = strconv.Atoi(s)
		}
		var data []map[string]string
		for i := page * size; i < (page+1)*size && i < total; i++ {
			data = append(data, map[string]string{"name": fmt.Sprintf("f%d", i)})
		}
		if q.Get("includePagination") != "true" {
			json.NewEncoder(w).Encode(data)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": data,
			"pageInfo": map[string]int{
				"currentPage":  page,
				"pageSize":     size,
				"totalRecords": total,
				"totalPages":   (total + size - 1) / size,
			},
		})
	}))
}

func TestDoPaginatedFollowsPageInfoWithSizeParam(t *testing.T) {
	var pages []string
	srv := flowStore(11, &pages)
	defer srv.Close()

	req := NewRequest(srv.URL, "GET", "/flows")
	req.QueryParam("size", "3")
	req.QueryParam("includePagination", "true")
	req.PageSizeParam("size")
	body, _, err := req.DoPaginated(false)
	if err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if got := names(t, body); len(got) != 11 {
		t.Fatalf("got %d items %v, want 11", len(got), got)
	}
	if len(pages) != 1 {
		t.Fatalf("requested pages %v, want one page of 100", pages)
	}
}

func TestDoPaginatedFollowsPageInfoAcrossPages(t *testing.T) {
	var pages []string
	srv := flowStore(250, &pages)
	defer srv.Close()

	req := NewRequest(srv.URL, "GET", "/flows")
	req.QueryParam("includePagination", "true")
	req.PageSizeParam("size")
	body, _, err := req.DoPaginated(false)
	if err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if got := names(t, body); len(got) != 250 || got[249] != "f249" {
		t.Fatalf("got %d items, want 250 ending in f249", len(got))
	}
	if fmt.Sprint(pages) != "[0 1 2]" {
		t.Fatalf("requested pages %v, want [0 1 2]", pages)
	}
}

func TestDoPaginatedAcceptsPlainArray(t *testing.T) {
	var pages []string
	srv := flowStore(11, &pages)
	defer srv.Close()

	req := NewRequest(srv.URL, "GET", "/flows")
	req.PageSizeParam("size")
	body, _, err := req.DoPaginated(false)
	if err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if got := names(t, body); len(got) != 11 {
		t.Fatalf("got %d items %v, want 11", len(got), got)
	}
}

func TestDoPaginatedLeavesNonListResponseAlone(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"s1","name":"Language_EN-US"}`))
	}))
	defer srv.Close()

	body, _, err := NewRequest(srv.URL, "GET", "/skill/s1").DoPaginated(false)
	if err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if string(body) != `{"id":"s1","name":"Language_EN-US"}` {
		t.Fatalf("got %s, want the object unchanged", body)
	}
}

// rooms mimics a core Webex list: it pages only through a Link header and
// ignores start.
func rooms(total int, hits *int) *httptest.Server {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*hits++
		q := r.URL.Query()
		max, _ := strconv.Atoi(q.Get("max"))
		from, _ := strconv.Atoi(q.Get("cursor"))
		var items []map[string]string
		for i := from; i < from+max && i < total; i++ {
			items = append(items, map[string]string{"name": fmt.Sprintf("r%d", i)})
		}
		if from+max < total {
			w.Header().Set("Link", fmt.Sprintf(`<%s/rooms?max=%d&cursor=%d>; rel="next"`, srv.URL, max, from+max))
		}
		json.NewEncoder(w).Encode(map[string]any{"items": items})
	}))
	return srv
}

func TestDoPaginatedFollowsLinkHeader(t *testing.T) {
	var hits int
	srv := rooms(250, &hits)
	defer srv.Close()

	req := NewRequest(srv.URL, "GET", "/rooms")
	req.OffsetPaging("", "max", 0)
	body, _, err := req.DoPaginated(true)
	if err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if got := names(t, body); len(got) != 250 || got[249] != "r249" {
		t.Fatalf("got %d items, want 250 ending in r249", len(got))
	}
	if hits != 3 {
		t.Fatalf("server hits = %d, want 3", hits)
	}
}

func TestDoPaginatedIgnoresLinkToAnotherHost(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Link", `<https://elsewhere.example/rooms?cursor=1>; rel="next"`)
		json.NewEncoder(w).Encode(map[string]any{"items": []map[string]string{{"name": "r0"}}})
	}))
	defer srv.Close()

	req := NewRequest(srv.URL, "GET", "/rooms")
	req.OffsetPaging("", "max", 0)
	if _, _, err := req.DoPaginated(true); err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if hits != 1 {
		t.Fatalf("server hits = %d, want 1", hits)
	}
}

func TestDoPaginatedStopsWithoutLinkWhenNoOffsetParam(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if r.URL.Query().Has("start") {
			t.Errorf("query %q sent start to an endpoint without it", r.URL.RawQuery)
		}
		items := make([]map[string]string, 100)
		for i := range items {
			items[i] = map[string]string{"name": fmt.Sprintf("p%d", i)}
		}
		json.NewEncoder(w).Encode(map[string]any{"items": items})
	}))
	defer srv.Close()

	req := NewRequest(srv.URL, "GET", "/people")
	req.OffsetPaging("", "max", 0)
	body, _, err := req.DoPaginated(true)
	if err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if got := names(t, body); len(got) != 100 || hits != 1 {
		t.Fatalf("got %d items in %d requests, want 100 in 1", len(got), hits)
	}
}

func TestDoPaginatedStopsWhenOffsetIsIgnored(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		items := make([]map[string]string, 100)
		for i := range items {
			items[i] = map[string]string{"name": fmt.Sprintf("p%d", i)}
		}
		json.NewEncoder(w).Encode(map[string]any{"items": items})
	}))
	defer srv.Close()

	body, _, err := NewRequest(srv.URL, "GET", "/people").DoPaginated(true)
	if err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if got := names(t, body); len(got) != 100 || hits != 2 {
		t.Fatalf("got %d items in %d requests, want 100 in 2", len(got), hits)
	}
}

func TestDoPaginatedPagesByStartAndMax(t *testing.T) {
	var starts []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		starts = append(starts, q.Get("start"))
		start, _ := strconv.Atoi(q.Get("start"))
		max, _ := strconv.Atoi(q.Get("max"))
		var items []map[string]string
		for i := start; i < start+max && i < 150; i++ {
			items = append(items, map[string]string{"name": fmt.Sprintf("n%d", i)})
		}
		json.NewEncoder(w).Encode(map[string]any{"numbers": items})
	}))
	defer srv.Close()

	body, _, err := NewRequest(srv.URL, "GET", "/telephony/config/numbers").DoPaginated(true)
	if err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if got := names(t, body); len(got) != 150 || got[149] != "n149" {
		t.Fatalf("got %d items, want 150 ending in n149", len(got))
	}
	if fmt.Sprint(starts) != "[0 100]" {
		t.Fatalf("starts %v, want [0 100]", starts)
	}
}

func TestDoPaginatedPagesSCIM(t *testing.T) {
	var starts []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		starts = append(starts, q.Get("startIndex"))
		start, _ := strconv.Atoi(q.Get("startIndex"))
		count, _ := strconv.Atoi(q.Get("count"))
		var items []map[string]string
		for i := start; i < start+count && i <= 120; i++ {
			items = append(items, map[string]string{"name": fmt.Sprintf("u%d", i)})
		}
		json.NewEncoder(w).Encode(map[string]any{
			"schemas":      []string{"urn:ietf:params:scim:api:messages:2.0:ListResponse"},
			"totalResults": 120, "startIndex": start, "itemsPerPage": len(items),
			"Resources": items,
		})
	}))
	defer srv.Close()

	req := NewRequest(srv.URL, "GET", "/Users")
	req.OffsetPaging("startIndex", "count", 1)
	body, _, err := req.DoPaginated(true)
	if err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if got := names(t, body); len(got) != 120 || got[0] != "u1" || got[119] != "u120" {
		t.Fatalf("got %d items, want u1..u120", len(got))
	}
	if fmt.Sprint(starts) != "[1 101]" {
		t.Fatalf("startIndex %v, want [1 101]", starts)
	}
}

func TestDoPaginatedFollowsMetaWithPageSize(t *testing.T) {
	var pages []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		pages = append(pages, q.Get("page"))
		if q.Get("pageSize") != "100" || q.Has("size") {
			t.Errorf("query %q, want pageSize=100 and no size", r.URL.RawQuery)
		}
		page, _ := strconv.Atoi(q.Get("page"))
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]string{{"name": fmt.Sprintf("s%d", page)}},
			"meta": map[string]int{"page": page, "totalPages": 2},
		})
	}))
	defer srv.Close()

	body, _, err := NewRequest(srv.URL, "GET", "/skill").DoPaginated(false)
	if err != nil {
		t.Fatalf("DoPaginated: %v", err)
	}
	if got := fmt.Sprint(names(t, body)); got != "[s0 s1]" {
		t.Fatalf("got %s, want [s0 s1]", got)
	}
}
