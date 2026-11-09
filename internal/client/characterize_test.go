package client

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Cloverhound/webex-cli/internal/config"
)

// These tests pin user-visible CLI behaviour: scripts and agents match this error
// and stderr text, so it must survive changes to the transport underneath.

func captureStderr(t *testing.T, run func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	saved := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = saved }()
	run()
	w.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func resetClientConfig(t *testing.T) {
	t.Helper()
	reset := func() {
		config.SetToken("tok-1")
		config.TokenRefresher = nil
		refreshedAt = time.Time{}
		config.SetMaxRetry(3)
		config.SetMaxRetryTimer(0)
		config.SetDryRun(false)
		config.SetReadOnly(false)
		config.SetDebug(false)
	}
	reset()
	t.Cleanup(reset)
}

func TestRefreshThenRateLimitOrderAndStderr(t *testing.T) {
	resetClientConfig(t)
	var seen []string
	statuses := []int{401, 429, 200}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Header.Get("Authorization"))
		status := statuses[len(seen)-1]
		if status == 429 {
			w.Header().Set("Retry-After", "1")
		}
		w.WriteHeader(status)
		io.WriteString(w, `{"ok":true}`)
	}))
	defer srv.Close()

	refreshes := 0
	config.TokenRefresher = func() (string, error) {
		refreshes++
		return "tok-2", nil
	}

	var body []byte
	var status int
	var err error
	stderr := captureStderr(t, func() {
		body, status, err = NewRequest(srv.URL, "GET", "/x").Do()
	})
	if err != nil || status != 200 || string(body) != `{"ok":true}` {
		t.Fatalf("got body=%q status=%d err=%v", body, status, err)
	}
	if refreshes != 1 {
		t.Fatalf("refreshes = %d", refreshes)
	}
	want := []string{"Bearer tok-1", "Bearer tok-2", "Bearer tok-2"}
	if strings.Join(seen, ",") != strings.Join(want, ",") {
		t.Fatalf("authorization headers = %v", seen)
	}
	if stderr != "Rate limited (429). Retrying in 1s (attempt 1/3)...\n" {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestFailedRefreshReturnsTheOriginal401(t *testing.T) {
	resetClientConfig(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		io.WriteString(w, "denied")
	}))
	defer srv.Close()
	config.TokenRefresher = func() (string, error) { return "", errors.New("no refresh") }

	body, status, err := NewRequest(srv.URL, "GET", "/x").Do()
	if status != 401 || string(body) != "denied" {
		t.Fatalf("status=%d body=%q", status, body)
	}
	if err == nil || err.Error() != "API error 401: denied" {
		t.Fatalf("err = %v", err)
	}
}

func TestErrorKeepsBodyAndTruncatesMessage(t *testing.T) {
	resetClientConfig(t)
	long := strings.Repeat("e", 600)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		io.WriteString(w, long)
	}))
	defer srv.Close()

	body, status, err := NewRequest(srv.URL, "GET", "/x").Do()
	if status != 404 || string(body) != long {
		t.Fatalf("status=%d len(body)=%d", status, len(body))
	}
	if err == nil || err.Error() != "API error 404: "+long[:500]+"..." {
		t.Fatalf("err = %v", err)
	}
}

func TestRateLimitMessages(t *testing.T) {
	resetClientConfig(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
	}))
	defer srv.Close()

	config.SetMaxRetry(0)
	_, status, err := NewRequest(srv.URL, "GET", "/x").Do()
	want := "rate limited (429): retry limit reached after 0 attempts — try again after 5s"
	if status != 429 || err == nil || err.Error() != want {
		t.Fatalf("status=%d err=%v", status, err)
	}

	config.SetMaxRetry(3)
	config.SetMaxRetryTimer(1)
	_, status, err = NewRequest(srv.URL, "GET", "/x").Do()
	want = "rate limited (429): retry wait of 5s would exceed max-retry-timer of 1s — try again after 5s"
	if status != 429 || err == nil || err.Error() != want {
		t.Fatalf("status=%d err=%v", status, err)
	}
}

func TestDryRunAndReadOnlySendNothing(t *testing.T) {
	resetClientConfig(t)
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { hits++ }))
	defer srv.Close()

	config.SetDryRun(true)
	var status int
	var err error
	stderr := captureStderr(t, func() {
		req := NewRequest(srv.URL, "POST", "/x")
		req.SetBodyRaw(`{"a":1}`)
		_, status, err = req.Do()
	})
	if !errors.Is(err, ErrDryRun) || status != 0 {
		t.Fatalf("status=%d err=%v", status, err)
	}
	if !strings.Contains(stderr, "[DRY RUN] POST "+srv.URL+"/x\n") || !strings.Contains(stderr, `[DRY RUN]   Body: {"a":1}`) {
		t.Fatalf("stderr = %q", stderr)
	}

	config.SetDryRun(false)
	config.SetReadOnly(true)
	_, status, err = NewRequest(srv.URL, "DELETE", "/x").Do()
	if !errors.Is(err, ErrReadOnly) || status != 0 {
		t.Fatalf("status=%d err=%v", status, err)
	}
	if hits != 0 {
		t.Fatalf("server received %d requests", hits)
	}
}

func TestWrongRegionMessage(t *testing.T) {
	resetClientConfig(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(451)
		io.WriteString(w, `use https://analytics-eu.webexapis.com please`)
	}))
	defer srv.Close()

	_, status, err := NewRequest(srv.URL, "GET", "/x").Do()
	if status != 451 || err == nil {
		t.Fatalf("status=%d err=%v", status, err)
	}
	for _, want := range []string{"wrong data region (451)", "Webex says to use https://analytics-eu.webexapis.com", "rerun with --region"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("err = %v, missing %q", err, want)
		}
	}
}

func TestTransportFailureIsWrapped(t *testing.T) {
	resetClientConfig(t)
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close()

	_, status, err := NewRequest(url, "GET", "/x").Do()
	if status != 0 || err == nil || !strings.HasPrefix(err.Error(), "executing request: ") {
		t.Fatalf("status=%d err=%v", status, err)
	}
}

type recordingTransport struct{ used bool }

func (r *recordingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	r.used = true
	return &http.Response{StatusCode: 200, Header: http.Header{}, Body: http.NoBody}, nil
}

func TestUsesDefaultClientSoLongDownloadsHaveNoTimeout(t *testing.T) {
	resetClientConfig(t)
	saved := http.DefaultClient.Transport
	rec := &recordingTransport{}
	http.DefaultClient.Transport = rec
	defer func() { http.DefaultClient.Transport = saved }()

	if _, status, err := NewRequest("http://example.invalid", "GET", "/x").Do(); err != nil || status != 200 {
		t.Fatalf("status=%d err=%v", status, err)
	}
	if !rec.used {
		t.Fatal("http.DefaultClient was not used")
	}
	if http.DefaultClient.Timeout != 0 {
		t.Fatalf("DefaultClient timeout = %s", http.DefaultClient.Timeout)
	}
}
