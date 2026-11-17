package webexapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// statusServer answers each request with the next status, repeating the last one.
func statusServer(t *testing.T, statuses []int, body string, header http.Header) (*httptest.Server, *[]*http.Request) {
	t.Helper()
	var seen []*http.Request
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Clone(context.Background()))
		status := statuses[min(len(seen)-1, len(statuses)-1)]
		for k, v := range header {
			w.Header()[k] = v
		}
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv, &seen
}

func noSleep(context.Context, time.Duration) error { return nil }

func TestA401RefreshesTheTokenAndRetriesOnce(t *testing.T) {
	srv, seen := statusServer(t, []int{401, 401}, `{}`, nil)
	token := "old"
	refreshes := 0
	c := &Client{
		Token:   func(context.Context) (string, error) { return token, nil },
		Refresh: func(context.Context) error { refreshes++; token = "new"; return nil },
	}
	resp, err := c.Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Status != 401 || resp.Status != 401 {
		t.Fatalf("resp %v, err %v; want the second 401", resp, err)
	}
	if refreshes != 1 || len(*seen) != 2 {
		t.Fatalf("refreshes %d, requests %d; want 1 and 2", refreshes, len(*seen))
	}
	if got := (*seen)[1].Header.Get("Authorization"); got != "Bearer new" {
		t.Fatalf("retry sent %q", got)
	}
}

func TestAFailedRefreshReturnsThe401(t *testing.T) {
	srv, seen := statusServer(t, []int{401}, `{"message":"expired"}`, nil)
	c := &Client{Refresh: func(context.Context) error { return errors.New("no refresh token") }}
	_, err := c.Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	if err == nil || err.Error() != `API error 401: {"message":"expired"}` || len(*seen) != 1 {
		t.Fatalf("err %v after %d requests", err, len(*seen))
	}
}

func TestA429WaitsForRetryAfterWithinItsLimit(t *testing.T) {
	srv, seen := statusServer(t, []int{429, 429, 200}, `{}`, http.Header{"Retry-After": {"2"}})
	var waits []time.Duration
	var retries []int
	c := &Client{
		MaxRetries: 3,
		Sleep:      func(_ context.Context, d time.Duration) error { waits = append(waits, d); return nil },
		OnRetry:    func(_ time.Duration, attempt, max int) { retries = append(retries, attempt) },
	}
	resp, err := c.Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	if err != nil || resp.Status != 200 || len(*seen) != 3 {
		t.Fatalf("resp %v, err %v, requests %d", resp, err, len(*seen))
	}
	if len(waits) != 2 || waits[0] != 2*time.Second || len(retries) != 2 || retries[1] != 2 {
		t.Fatalf("waits %v, retries %v", waits, retries)
	}
}

func TestA429StopsAtItsRetryLimits(t *testing.T) {
	srv, _ := statusServer(t, []int{429}, `{}`, http.Header{"Retry-After": {"10"}})
	_, err := (&Client{MaxRetries: 1, Sleep: noSleep}).Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	if err == nil || err.Error() != "rate limited (429): retry limit reached after 1 attempts — try again after 10s" {
		t.Fatalf("err %v", err)
	}
	_, err = (&Client{MaxRetries: 5, MaxRetryWait: 15 * time.Second, Sleep: noSleep}).Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	if err == nil || !strings.Contains(err.Error(), "retry wait of 20s would exceed max-retry-timer of 15s") {
		t.Fatalf("err %v", err)
	}
}

func TestA451NamesTheRegionalEndpoint(t *testing.T) {
	srv, _ := statusServer(t, []int{451}, `{"message":"use https://analytics-calling-eu.webexapis.com instead"}`, nil)
	resp, err := (&Client{}).Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	var wrong *WrongRegionError
	if !errors.As(err, &wrong) || wrong.Endpoint != "https://analytics-calling-eu.webexapis.com" || resp.Status != 451 {
		t.Fatalf("err %v", err)
	}
}

func TestAnErrorRepeatsTheRedactedBodyAndTheResponseKeepsIt(t *testing.T) {
	body := `{"token":"secret","pad":"` + strings.Repeat("x", 600) + `"}`
	srv, _ := statusServer(t, []int{400}, body, nil)
	redact := func(b []byte) []byte { return []byte(strings.ReplaceAll(string(b), "secret", "<redacted>")) }
	resp, err := (&Client{}).Send(context.Background(), Request{Method: "GET", URL: srv.URL, Redact: redact})
	if err == nil || strings.Contains(err.Error(), "secret") || !strings.HasSuffix(err.Error(), "...") {
		t.Fatalf("err %v", err)
	}
	if string(resp.Body) != body {
		t.Fatal("the response body was changed")
	}
}

func TestBeforeSendCanStopARequest(t *testing.T) {
	srv, seen := statusServer(t, []int{200}, `{}`, nil)
	stop := errors.New("dry run")
	var sent *http.Request
	c := &Client{BeforeSend: func(r *http.Request) error { sent = r; return stop }}
	resp, err := c.Send(context.Background(), Request{Method: "DELETE", URL: srv.URL})
	if !errors.Is(err, stop) || resp != nil || len(*seen) != 0 {
		t.Fatalf("resp %v, err %v, requests %d", resp, err, len(*seen))
	}
	if sent.Header.Get("Authorization") != "Bearer " {
		t.Fatalf("BeforeSend saw %v", sent.Header)
	}
}

func TestRequestHeaders(t *testing.T) {
	srv, seen := statusServer(t, []int{200}, `{}`, nil)
	c := &Client{Token: func(context.Context) (string, error) { return "t", nil }}
	_, err := c.Send(context.Background(), Request{
		Method: "POST", URL: srv.URL, Body: []byte(`{}`),
		Header: http.Header{"Accept": {"*/*"}, "Trackingid": {"x"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	h := (*seen)[0].Header
	if h.Get("Authorization") != "Bearer t" || h.Get("Content-Type") != "application/json" ||
		h.Get("Accept") != "*/*" || h.Get("Trackingid") != "x" {
		t.Fatalf("headers %v", h)
	}
	_, _ = c.Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	if ct := (*seen)[1].Header.Get("Content-Type"); ct != "" {
		t.Fatalf("a request with no body sent Content-Type %q", ct)
	}
}

func TestHostsForARegion(t *testing.T) {
	h := HostsFor("EU", "eu2")
	if h.AnalyticsCalling != "https://analytics-calling-eu.webexapis.com/v1" || h.CC != "https://api.wxcc-eu2.cisco.com" ||
		h.AIAgent != "https://studio.aiagent-eu2.cisco.com" || h.Webex != "https://webexapis.com/v1" {
		t.Fatalf("%+v", h)
	}
	if d := HostsFor("", ""); d.CC != "https://api.wxcc-us1.cisco.com" || d.AnalyticsCalling != "https://analytics-calling.webexapis.com/v1" {
		t.Fatalf("%+v", d)
	}
}

func TestOrgIDConversion(t *testing.T) {
	const uuid = "6ab1d2c4-0000-4000-8000-00000000c1d0"
	b64 := EncodeOrgID(uuid)
	if strings.HasSuffix(b64, "=") || DecodeOrgID(b64) != uuid || DecodeOrgID(b64+"=") != uuid {
		t.Fatalf("%q", b64)
	}
	if EncodeOrgID(b64) != b64 || DecodeOrgID(uuid) != uuid || DecodeOrgID("other") != "other" {
		t.Fatal("non-matching values are returned as they are")
	}
}
