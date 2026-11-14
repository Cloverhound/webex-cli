package webexapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func list(t *testing.T, url string, p Paging) ([]json.RawMessage, error) {
	t.Helper()
	c := &Client{Token: staticToken("tok")}
	return collect(t, func(v func([]json.RawMessage) error) error { return c.List(context.Background(), url, p, v) })
}

func TestCCListUsesItsPageSizeParameter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("size") != "100" || r.URL.Query().Has("pageSize") {
			w.WriteHeader(400)
			return
		}
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		fmt.Fprintf(w, `{"data":%s,"pageInfo":{"totalPages":2}}`, items(page*100, 100))
	}))
	defer srv.Close()
	all, err := list(t, srv.URL, Paging{CC: true, SizeParam: "size"})
	if err != nil || len(all) != 200 {
		t.Fatalf("len=%d err=%v", len(all), err)
	}
}

func TestOffsetListNamesItsParametersAndFindsTheOnlyArray(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		n := min(limit, 150-offset)
		fmt.Fprintf(w, `{"meta":{"total_count":150},"objects":%s}`, items(offset, n))
	}))
	defer srv.Close()
	all, err := list(t, srv.URL, Paging{Start: "offset", Size: "limit"})
	if err != nil || len(all) != 150 {
		t.Fatalf("len=%d err=%v", len(all), err)
	}
}

func TestSCIMListStartsAtOneAndStopsAtTotalResults(t *testing.T) {
	var starts []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start, _ := strconv.Atoi(r.URL.Query().Get("startIndex"))
		starts = append(starts, r.URL.Query().Get("startIndex"))
		fmt.Fprintf(w, `{"totalResults":200,"Resources":%s}`, items(start, 100))
	}))
	defer srv.Close()
	all, err := list(t, srv.URL, Paging{Start: "startIndex", Size: "count", First: 1})
	if err != nil || len(all) != 200 || fmt.Sprint(starts) != "[1 101]" {
		t.Fatalf("len=%d starts=%v err=%v", len(all), starts, err)
	}
}

func TestLinkOnlyListReadsOnePageWithoutALink(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		fmt.Fprintf(w, `{"items":%s}`, items(0, 100))
	}))
	defer srv.Close()
	all, err := list(t, srv.URL, Paging{})
	if err != nil || len(all) != 100 || calls != 1 {
		t.Fatalf("len=%d calls=%d err=%v", len(all), calls, err)
	}
}

func TestOffsetListRefusesAmbiguousOrMissingArrays(t *testing.T) {
	for _, body := range []string{`{"a":[],"b":[]}`, `{"count":3}`} {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, body) }))
		_, err := list(t, srv.URL, Paging{Start: "start", Size: "max", First: 0})
		srv.Close()
		if !errors.Is(err, ErrIncompleteList) {
			t.Errorf("%s: err = %v", body, err)
		}
	}
	if _, err := list(t, "http://unused.example", Paging{Start: "offset"}); err == nil {
		t.Error("an offset without a page size was accepted")
	}
}
