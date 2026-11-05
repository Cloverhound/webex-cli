package webexapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func staticToken(token string) TokenSource {
	return func(context.Context) (string, error) { return token, nil }
}

func noSleep(context.Context, time.Duration) error { return nil }

func TestSendReturnsResponseAndSetsHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer tok" || r.Header.Get("Accept") != "application/json" {
			t.Errorf("headers = %v", r.Header)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("content type = %q", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != `{"a":1}` {
			t.Errorf("body = %q", body)
		}
		w.Header().Set("Trackingid", "T1")
		w.WriteHeader(201)
		io.WriteString(w, `{"id":"x"}`)
	}))
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	resp, err := c.Send(context.Background(), Request{Method: "POST", URL: srv.URL, Body: []byte(`{"a":1}`)})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != 201 || string(resp.Body) != `{"id":"x"}` || resp.Header.Get("Trackingid") != "T1" {
		t.Fatalf("resp = %+v", resp)
	}
}

func TestErrorStatusIsStructuredAndKeepsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Trackingid", "T2")
		w.WriteHeader(404)
		io.WriteString(w, "missing")
	}))
	defer srv.Close()

	c := &Client{Token: staticToken("tok")}
	resp, err := c.Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %v", err)
	}
	if apiErr.Status != 404 || apiErr.TrackingID != "T2" || string(apiErr.Body) != "missing" {
		t.Fatalf("apiErr = %+v", apiErr)
	}
	if err.Error() != "API error 404: missing" {
		t.Fatalf("message = %q", err)
	}
	if resp == nil || resp.Status != 404 || string(resp.Body) != "missing" {
		t.Fatalf("resp = %+v", resp)
	}
}

func TestRefreshesOnceOn401ThenHandles429(t *testing.T) {
	statuses := []int{401, 429, 200}
	var auth []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = append(auth, r.Header.Get("Authorization"))
		status := statuses[len(auth)-1]
		if status == 429 {
			w.Header().Set("Retry-After", "2")
		}
		w.WriteHeader(status)
	}))
	defer srv.Close()

	token := "old"
	var retried []string
	c := &Client{
		Token:      func(context.Context) (string, error) { return token, nil },
		Refresh:    func(context.Context) error { token = "new"; return nil },
		MaxRetries: 3,
		Sleep:      noSleep,
		OnRetry: func(wait time.Duration, attempt, max int) {
			retried = append(retried, wait.String())
		},
	}
	resp, err := c.Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	if err != nil || resp.Status != 200 {
		t.Fatalf("resp=%+v err=%v", resp, err)
	}
	if strings.Join(auth, ",") != "Bearer old,Bearer new,Bearer new" {
		t.Fatalf("auth = %v", auth)
	}
	if strings.Join(retried, ",") != "2s" {
		t.Fatalf("retried = %v", retried)
	}
}

func TestRateLimitLimits(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(429)
	}))
	defer srv.Close()

	c := &Client{Token: staticToken("tok"), MaxRetries: 2, Sleep: noSleep}
	_, err := c.Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	if err == nil || err.Error() != "rate limited (429): retry limit reached after 2 attempts — try again after 5s" {
		t.Fatalf("err = %v", err)
	}
	if calls != 3 {
		t.Fatalf("calls = %d", calls)
	}

	c.MaxRetryWait = time.Second
	_, err = c.Send(context.Background(), Request{Method: "GET", URL: srv.URL})
	if err == nil || err.Error() != "rate limited (429): retry wait of 5s would exceed max-retry-timer of 1s — try again after 5s" {
		t.Fatalf("err = %v", err)
	}
}

func TestServerErrorIsNeverRetried(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(500)
	}))
	defer srv.Close()

	c := &Client{Token: staticToken("tok"), MaxRetries: 3, Sleep: noSleep}
	if _, err := c.Send(context.Background(), Request{Method: "PUT", URL: srv.URL}); err == nil {
		t.Fatal("expected an error")
	}
	if calls != 1 {
		t.Fatalf("calls = %d", calls)
	}
}

func TestBeforeSendErrorIsReturnedUnchangedAndNothingIsSent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("request was sent")
	}))
	defer srv.Close()

	blocked := errors.New("blocked")
	c := &Client{Token: staticToken("tok"), BeforeSend: func(*http.Request) error { return blocked }}
	resp, err := c.Send(context.Background(), Request{Method: "POST", URL: srv.URL})
	if err != blocked || resp != nil {
		t.Fatalf("resp=%+v err=%v", resp, err)
	}
}

func TestTransportFailureIsNotAnAPIError(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close()

	c := &Client{Token: staticToken("tok")}
	resp, err := c.Send(context.Background(), Request{Method: "PUT", URL: url})
	var apiErr *APIError
	if err == nil || errors.As(err, &apiErr) || resp != nil {
		t.Fatalf("resp=%+v err=%v", resp, err)
	}
	if !strings.HasPrefix(err.Error(), "executing request: ") {
		t.Fatalf("err = %v", err)
	}
}

func TestCancelledContextInterruptsRetryWait(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(429)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{Token: staticToken("tok"), MaxRetries: 3, OnRetry: func(time.Duration, int, int) { cancel() }}
	start := time.Now()
	_, err := c.Send(ctx, Request{Method: "GET", URL: srv.URL})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatal("retry wait was not interrupted")
	}
}

func TestDefaultHTTPClientHasATimeout(t *testing.T) {
	if defaultHTTP.Timeout <= 0 {
		t.Fatal("library callers must not hang forever")
	}
}

func TestNextLinkOnlyFollowsTheSameHost(t *testing.T) {
	h := http.Header{}
	h.Add("Link", `<https://webexapis.com/v1/x?cursor=2>; rel="next"`)
	if got := NextLink(h, "https://webexapis.com/v1"); got != "https://webexapis.com/v1/x?cursor=2" {
		t.Fatalf("got %q", got)
	}
	if got := NextLink(h, "https://other.example.com"); got != "" {
		t.Fatalf("got %q", got)
	}
}
