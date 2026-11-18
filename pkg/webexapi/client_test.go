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

func TestARefreshCutOffByTheContextReturnsTheContextError(t *testing.T) {
	srv, _ := statusServer(t, []int{401}, `{}`, nil)
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{Refresh: func(ctx context.Context) error { cancel(); return ctx.Err() }}
	resp, err := c.Send(ctx, Request{Method: "GET", URL: srv.URL})
	if !errors.Is(err, context.Canceled) || resp == nil || resp.Status != 401 {
		t.Fatalf("resp %v, err %v", resp, err)
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

func TestRedirectsPassBeforeSendAndLoseTheAgentTokenOnAnotherHost(t *testing.T) {
	var other []*http.Request
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		other = append(other, r.Clone(context.Background()))
	}))
	t.Cleanup(target.Close)
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL+"/elsewhere", http.StatusFound)
	}))
	t.Cleanup(origin.Close)
	header := http.Header{"Bot-Access-Token": {"agent-token"}}

	var seen []string
	refuse := errors.New("not approved")
	c := &Client{BeforeSend: func(r *http.Request) error {
		seen = append(seen, r.URL.Host)
		if r.URL.Host != strings.TrimPrefix(origin.URL, "http://") {
			return refuse
		}
		return nil
	}}
	_, err := c.Send(context.Background(), Request{Method: "GET", URL: origin.URL, Header: header})
	if !errors.Is(err, refuse) || len(seen) != 2 || len(other) != 0 {
		t.Fatalf("err %v, BeforeSend saw %v, other host got %d requests", err, seen, len(other))
	}

	if _, err := (&Client{}).Send(context.Background(), Request{Method: "GET", URL: origin.URL, Header: header}); err != nil {
		t.Fatal(err)
	}
	if len(other) != 1 || other[0].Header.Get("Bot-Access-Token") != "" {
		t.Fatalf("the agent token reached another host: %v", other)
	}
}

func TestARedirectToAnotherPortSendsNoCredentials(t *testing.T) {
	var got http.Header
	// Same hostname, another port: net/http itself would forward Authorization.
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { got = r.Header.Clone() }))
	t.Cleanup(target.Close)
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	t.Cleanup(origin.Close)
	c := &Client{Token: func(context.Context) (string, error) { return "secret", nil }}
	header := http.Header{"Bot-Access-Token": {"agent"}, "Password": {"meeting-password"}}
	if _, err := c.Send(context.Background(), Request{Method: "GET", URL: origin.URL, Header: header}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"Authorization", "Bot-Access-Token", "Password"} {
		if got.Get(name) != "" {
			t.Errorf("%s reached another origin", name)
		}
	}
	if got.Get("Accept") != "application/json" {
		t.Errorf("Accept was dropped: %v", got)
	}
}

func TestARefusedRedirectReturnsItsResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/moved", http.StatusTemporaryRedirect)
	}))
	t.Cleanup(srv.Close)
	refuse := errors.New("refused")
	calls := 0
	c := &Client{BeforeSend: func(*http.Request) error {
		if calls++; calls > 1 {
			return refuse
		}
		return nil
	}}
	resp, err := c.Send(context.Background(), Request{Method: "POST", URL: srv.URL, Body: []byte("{}")})
	if !errors.Is(err, refuse) || resp == nil || resp.Status != http.StatusTemporaryRedirect || resp.Header.Get("Location") != "/moved" {
		t.Fatalf("resp %+v, err %v", resp, err)
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
	h, err := HostsFor("EU", "eu2")
	if err != nil || h.AnalyticsCalling != "https://analytics-calling-eu.webexapis.com/v1" || h.CC != "https://api.wxcc-eu2.cisco.com" ||
		h.AIAgent != "https://studio.aiagent-eu2.cisco.com" || h.Webex != "https://webexapis.com/v1" {
		t.Fatalf("%+v, %v", h, err)
	}
	if d, err := HostsFor("", ""); err != nil || d.CC != "https://api.wxcc-us1.cisco.com" || d.AnalyticsCalling != "https://analytics-calling.webexapis.com/v1" {
		t.Fatalf("%+v, %v", d, err)
	}
}

func TestHostsForRefusesAnUnknownContactCenterRegion(t *testing.T) {
	for _, ccRegion := range []string{"us1.example.com/#", "xx9", "US1"} {
		if h, err := HostsFor("us", ccRegion); err == nil {
			t.Errorf("%q gave %+v", ccRegion, h)
		}
	}
}

func TestARedactorThatReturnsNothingHidesTheWholeBody(t *testing.T) {
	srv, _ := statusServer(t, []int{400}, `{"secret":"x"}`, nil)
	_, err := (&Client{}).Send(context.Background(), Request{Method: "GET", URL: srv.URL, Redact: func([]byte) []byte { return nil }})
	if err == nil || err.Error() != "API error 400: " {
		t.Fatalf("err %v", err)
	}
}

func TestErrorsHoldOnlyTheRedactedBody(t *testing.T) {
	redact := func(b []byte) []byte { return []byte(strings.ReplaceAll(string(b), "secret", "<redacted>")) }
	for _, status := range []int{400, 451} {
		srv, _ := statusServer(t, []int{status}, `{"token":"secret","url":"https://x.webexapis.com"}`, nil)
		resp, err := (&Client{}).Send(context.Background(), Request{Method: "GET", URL: srv.URL, Redact: redact})
		var apiErr *APIError
		var wrong *WrongRegionError
		switch {
		case errors.As(err, &wrong):
			if strings.Contains(string(wrong.Body), "secret") {
				t.Errorf("%d: WrongRegionError.Body holds the secret", status)
			}
		case errors.As(err, &apiErr):
			if strings.Contains(string(apiErr.Body), "secret") {
				t.Errorf("%d: APIError.Body holds the secret", status)
			}
		default:
			t.Errorf("%d: err %v", status, err)
		}
		if !strings.Contains(string(resp.Body), "secret") {
			t.Errorf("%d: the response lost its body", status)
		}
	}
}

func TestARedactorThatEditsItsArgumentLeavesTheResponse(t *testing.T) {
	srv, _ := statusServer(t, []int{400}, `secret`, nil)
	redact := func(b []byte) []byte { copy(b, "XXXXXX"); return b }
	resp, err := (&Client{}).Send(context.Background(), Request{Method: "GET", URL: srv.URL, Redact: redact})
	if string(resp.Body) != "secret" || err.Error() != "API error 400: XXXXXX" {
		t.Fatalf("body %q, err %v", resp.Body, err)
	}
}

func TestA451NamesOnlyWhatTheRedactorLeaves(t *testing.T) {
	srv, _ := statusServer(t, []int{451}, `{"message":"use https://analytics-calling-eu.webexapis.com"}`, nil)
	resp, err := (&Client{}).Send(context.Background(), Request{Method: "GET", URL: srv.URL, Redact: func([]byte) []byte { return nil }})
	var wrong *WrongRegionError
	if !errors.As(err, &wrong) || wrong.Endpoint != "" || strings.Contains(err.Error(), "webexapis") {
		t.Fatalf("err %v", err)
	}
	if !strings.Contains(string(resp.Body), "webexapis") {
		t.Fatal("the response lost its body")
	}
}

func TestA429WhoseBodyCannotBeReadKeepsTheReadError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "100")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("short"))
	}))
	t.Cleanup(srv.Close)
	_, err := (&Client{}).Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	var limited *RateLimitError
	if err == nil || errors.As(err, &limited) || !strings.Contains(err.Error(), "reading response") {
		t.Fatalf("err %v", err)
	}
}

func TestChangingCCRegionsDoesNotChangeWhatHostsForAccepts(t *testing.T) {
	saved := CCRegions[0]
	CCRegions[0] = "changed"
	t.Cleanup(func() { CCRegions[0] = saved })
	if _, err := HostsFor("", ""); err != nil {
		t.Fatalf("the default data center was refused: %v", err)
	}
	if _, err := HostsFor("", "changed"); err == nil {
		t.Fatal("a data center added to CCRegions was accepted")
	}
}

func TestAnExhaustedRateLimitStillUnwrapsToTheAPIError(t *testing.T) {
	srv, _ := statusServer(t, []int{429}, `{"message":"slow down"}`, http.Header{"Retry-After": {"1"}})
	_, err := (&Client{}).Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Status != 429 || string(apiErr.Body) != `{"message":"slow down"}` {
		t.Fatalf("err %v does not unwrap to the 429", err)
	}
	if err.Error() != "rate limited (429): retry limit reached after 0 attempts — try again after 1s" {
		t.Fatalf("message %q", err)
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
