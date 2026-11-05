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

func items(from, n int) string {
	out := "["
	for i := 0; i < n; i++ {
		if i > 0 {
			out += ","
		}
		out += fmt.Sprintf(`{"id":"%d"}`, from+i)
	}
	return out + "]"
}

func collect(t *testing.T, run func(visit func([]json.RawMessage) error) error) ([]json.RawMessage, error) {
	t.Helper()
	var all []json.RawMessage
	err := run(func(page []json.RawMessage) error {
		all = append(all, page...)
		return nil
	})
	return all, err
}

func ccServer(handler func(page, size int) (int, string)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		size, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
		status, body := handler(page, size)
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
}

func TestCCPagerFollowsMetaTotalPages(t *testing.T) {
	srv := ccServer(func(page, size int) (int, string) {
		return 200, fmt.Sprintf(`{"data":%s,"meta":{"totalPages":3}}`, items(page*10, 10))
	})
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	all, err := collect(t, func(v func([]json.RawMessage) error) error { return c.ListCC(context.Background(), srv.URL, v) })
	if err != nil || len(all) != 30 {
		t.Fatalf("len=%d err=%v", len(all), err)
	}
}

func TestCCPagerReadsPastAFullBareArray(t *testing.T) {
	srv := ccServer(func(page, size int) (int, string) {
		if page == 0 {
			return 200, items(0, size)
		}
		return 200, items(size, 3)
	})
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	all, err := collect(t, func(v func([]json.RawMessage) error) error { return c.ListCC(context.Background(), srv.URL, v) })
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != ccPageSize+3 {
		t.Fatalf("len = %d", len(all))
	}
}

func TestCCPagerStopsOnAShortBareArray(t *testing.T) {
	calls := 0
	srv := ccServer(func(page, size int) (int, string) {
		calls++
		return 200, items(0, 8)
	})
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	all, err := collect(t, func(v func([]json.RawMessage) error) error { return c.ListCC(context.Background(), srv.URL, v) })
	if err != nil || len(all) != 8 || calls != 1 {
		t.Fatalf("len=%d calls=%d err=%v", len(all), calls, err)
	}
}

func TestCCPagerFailsWhenAFullPageRepeats(t *testing.T) {
	srv := ccServer(func(page, size int) (int, string) { return 200, items(0, size) })
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	_, err := collect(t, func(v func([]json.RawMessage) error) error { return c.ListCC(context.Background(), srv.URL, v) })
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v", err)
	}
}

func TestCCPagerFailsOnALaterPageWithoutData(t *testing.T) {
	srv := ccServer(func(page, size int) (int, string) {
		if page == 0 {
			return 200, fmt.Sprintf(`{"data":%s,"meta":{"totalPages":2}}`, items(0, 2))
		}
		return 200, `{"message":"try later"}`
	})
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	_, err := collect(t, func(v func([]json.RawMessage) error) error { return c.ListCC(context.Background(), srv.URL, v) })
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v", err)
	}
}

func TestCCPagerReadsPastFullPageWithoutMeta(t *testing.T) {
	srv := ccServer(func(page, size int) (int, string) {
		if page == 0 {
			return 200, fmt.Sprintf(`{"data":%s}`, items(0, size))
		}
		return 200, `{"data":[]}`
	})
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	all, err := collect(t, func(v func([]json.RawMessage) error) error { return c.ListCC(context.Background(), srv.URL, v) })
	if err != nil || len(all) != ccPageSize {
		t.Fatalf("len=%d err=%v", len(all), err)
	}
}

func TestCCPagerStopsWhenVisitFails(t *testing.T) {
	srv := ccServer(func(page, size int) (int, string) {
		return 200, fmt.Sprintf(`{"data":%s,"meta":{"totalPages":5}}`, items(0, 1))
	})
	defer srv.Close()

	stop := errors.New("stop")
	c := &Client{Token: staticToken("tok")}
	err := c.ListCC(context.Background(), srv.URL, func([]json.RawMessage) error { return stop })
	if err != stop {
		t.Fatalf("err = %v", err)
	}
}

func TestOffsetPagerAdvancesUntilAShortPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start, _ := strconv.Atoi(r.URL.Query().Get("start"))
		size, _ := strconv.Atoi(r.URL.Query().Get("max"))
		n := size
		if start > 0 {
			n = 4
		}
		fmt.Fprintf(w, `{"huntGroups":%s}`, items(start, n))
	}))
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	all, err := collect(t, func(v func([]json.RawMessage) error) error {
		return c.ListOffset(context.Background(), srv.URL+"?orgId=o", "huntGroups", v)
	})
	if err != nil || len(all) != offsetPageSize+4 {
		t.Fatalf("len=%d err=%v", len(all), err)
	}
}

func TestOffsetPagerFollowsLinkHeaders(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("cursor") == "" {
			w.Header().Set("Link", fmt.Sprintf(`<%s/x?cursor=2>; rel="next"`, srv.URL))
			fmt.Fprintf(w, `{"items":%s}`, items(0, 2))
			return
		}
		fmt.Fprintf(w, `{"items":%s}`, items(2, 2))
	}))
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	all, err := collect(t, func(v func([]json.RawMessage) error) error {
		return c.ListOffset(context.Background(), srv.URL+"/x", "items", v)
	})
	if err != nil || len(all) != 4 {
		t.Fatalf("len=%d err=%v", len(all), err)
	}
}

func TestOffsetPagerTreatsAnEmptyFirstPageObjectAsNoItems(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{}`)
	}))
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	all, err := collect(t, func(v func([]json.RawMessage) error) error {
		return c.ListOffset(context.Background(), srv.URL, "operatingModes", v)
	})
	if err != nil || len(all) != 0 {
		t.Fatalf("len=%d err=%v", len(all), err)
	}
}

func TestOffsetPagerFailsOnAnEmptyLaterPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("start") != "0" {
			fmt.Fprint(w, `{}`)
			return
		}
		size, _ := strconv.Atoi(r.URL.Query().Get("max"))
		fmt.Fprintf(w, `{"items":%s}`, items(0, size))
	}))
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	_, err := collect(t, func(v func([]json.RawMessage) error) error {
		return c.ListOffset(context.Background(), srv.URL, "items", v)
	})
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v", err)
	}
}

func TestOffsetPagerFailsWithoutTheNamedArray(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"somethingElse":[]}`)
	}))
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	_, err := collect(t, func(v func([]json.RawMessage) error) error {
		return c.ListOffset(context.Background(), srv.URL, "huntGroups", v)
	})
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v", err)
	}
}

func TestOffsetPagerFailsWhenAFullPageRepeats(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		size, _ := strconv.Atoi(r.URL.Query().Get("max"))
		fmt.Fprintf(w, `{"items":%s}`, items(0, size))
	}))
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	_, err := collect(t, func(v func([]json.RawMessage) error) error {
		return c.ListOffset(context.Background(), srv.URL, "items", v)
	})
	if !errors.Is(err, ErrIncompleteList) {
		t.Fatalf("err = %v", err)
	}
}
