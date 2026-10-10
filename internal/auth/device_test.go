package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// fakeDeviceToken serves the device token endpoint, answering with responses
// in order and repeating the last one.
func fakeDeviceToken(t *testing.T, responses []func(http.ResponseWriter)) (*int32, func()) {
	t.Helper()
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "cid" || pass != "secret" {
			t.Errorf("basic auth = %q/%q/%v, want cid/secret", user, pass, ok)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("grant_type") != "urn:ietf:params:oauth:grant-type:device_code" || r.Form.Get("device_code") != "dev" {
			t.Errorf("unexpected form %v", r.Form)
		}
		n := atomic.AddInt32(&calls, 1)
		i := int(n) - 1
		if i >= len(responses) {
			i = len(responses) - 1
		}
		responses[i](w)
	}))
	old := DeviceTokenURL
	DeviceTokenURL = srv.URL
	return &calls, func() { DeviceTokenURL = old; srv.Close() }
}

func respond(status int, body string) func(http.ResponseWriter) {
	return func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}
}

var granted = respond(200, `{"access_token":"at","refresh_token":"rt","expires_in":3600,"token_type":"Bearer","scope":"spark:all"}`)

func testCode(expiresIn int) *DeviceCode {
	return &DeviceCode{DeviceCode: "dev", UserCode: "123456", ExpiresIn: expiresIn, Interval: 1}
}

func TestPollDeviceTokenStopsStalledRequestAtExpiry(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-release:
		case <-time.After(5 * time.Second):
		}
	}))
	defer srv.Close()
	defer close(release)
	old := DeviceTokenURL
	DeviceTokenURL = srv.URL
	defer func() { DeviceTokenURL = old }()

	start := time.Now()
	_, err := PollDeviceToken(context.Background(), "cid", "secret", testCode(100), "spark:all")
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Errorf("err = %v, want the code to expire", err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("polling took %s; the stalled request was not cut off at expiry", elapsed)
	}
}

func TestPollDeviceTokenPendingThenGranted(t *testing.T) {
	calls, done := fakeDeviceToken(t, []func(http.ResponseWriter){
		respond(428, `{"message":"authorization pending"}`),
		respond(400, `{"error":"authorization_pending"}`),
		granted,
	})
	defer done()

	tok, err := PollDeviceToken(context.Background(), "cid", "secret", testCode(1000), "requested")
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "at" || tok.RefreshToken != "rt" || tok.Scopes != "spark:all" {
		t.Errorf("token = %+v", tok)
	}
	if *calls != 3 {
		t.Errorf("calls = %d, want 3", *calls)
	}
}

func TestPollDeviceTokenSlowDownAddsFiveIntervals(t *testing.T) {
	_, done := fakeDeviceToken(t, []func(http.ResponseWriter){
		respond(400, `{"error":"slow_down"}`),
		granted,
	})
	defer done()

	// With a 1-unit interval and a 5-unit expiry, the poll after slow_down
	// (at 1+6 units) lands past the deadline.
	_, err := PollDeviceToken(context.Background(), "cid", "secret", testCode(5), "")
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("err = %v, want expiry after slow_down", err)
	}
}

func TestPollDeviceTokenDenied(t *testing.T) {
	_, done := fakeDeviceToken(t, []func(http.ResponseWriter){respond(400, `{"error":"access_denied"}`)})
	defer done()

	_, err := PollDeviceToken(context.Background(), "cid", "secret", testCode(1000), "")
	if err == nil || !strings.Contains(err.Error(), "denied") {
		t.Fatalf("err = %v, want denied", err)
	}
}

func TestPollDeviceTokenExpiredGrant(t *testing.T) {
	_, done := fakeDeviceToken(t, []func(http.ResponseWriter){respond(400, `{"error":"expired_token"}`)})
	defer done()

	_, err := PollDeviceToken(context.Background(), "cid", "secret", testCode(1000), "")
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("err = %v, want expired", err)
	}
}

func TestPollDeviceTokenTimesOut(t *testing.T) {
	calls, done := fakeDeviceToken(t, []func(http.ResponseWriter){respond(428, `{}`)})
	defer done()

	_, err := PollDeviceToken(context.Background(), "cid", "secret", testCode(20), "")
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("err = %v, want expiry", err)
	}
	if *calls == 0 {
		t.Error("expected at least one poll before expiry")
	}
}

func TestRequestDeviceCodeRedirectMismatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte(`{"message":"redirect_uri_mismatch"}`))
	}))
	defer srv.Close()
	old := DeviceAuthorizeURL
	DeviceAuthorizeURL = srv.URL
	defer func() { DeviceAuthorizeURL = old }()

	_, err := RequestDeviceCode("cid", "spark:all")
	if err == nil || !strings.Contains(err.Error(), "oauth-helper-") {
		t.Fatalf("err = %v, want a hint naming the helper redirect URIs", err)
	}
}

func TestRequestDeviceCodeDefaults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		if r.Form.Get("client_id") != "cid" || r.Form.Get("scope") != "spark:all" {
			t.Errorf("form = %v", r.Form)
		}
		w.Write([]byte(`{"device_code":"dev","user_code":"123456","verification_uri":"https://v"}`))
	}))
	defer srv.Close()
	old := DeviceAuthorizeURL
	DeviceAuthorizeURL = srv.URL
	defer func() { DeviceAuthorizeURL = old }()

	dc, err := RequestDeviceCode("cid", "spark:all")
	if err != nil {
		t.Fatal(err)
	}
	if dc.ExpiresIn != 300 || dc.Interval != 5 {
		t.Errorf("defaults = %d/%d, want 300/5", dc.ExpiresIn, dc.Interval)
	}
}

func TestHeadlessReason(t *testing.T) {
	env := func(kv ...string) func(string) string {
		m := map[string]string{}
		for i := 0; i < len(kv); i += 2 {
			m[kv[i]] = kv[i+1]
		}
		return func(k string) string { return m[k] }
	}
	cases := []struct {
		name   string
		getenv func(string) string
		goos   string
		want   bool
	}{
		{"mac desktop", env(), "darwin", false},
		{"ssh", env("SSH_CONNECTION", "1 2 3 4"), "darwin", true},
		{"ssh tty", env("SSH_TTY", "/dev/pts/0"), "darwin", true},
		{"ci", env("CI", "true"), "windows", true},
		{"linux no display", env(), "linux", true},
		{"linux x11", env("DISPLAY", ":0"), "linux", false},
		{"linux wayland", env("WAYLAND_DISPLAY", "wayland-0"), "linux", false},
	}
	for _, c := range cases {
		if got := HeadlessReason(c.getenv, c.goos) != ""; got != c.want {
			t.Errorf("%s: headless = %v, want %v", c.name, got, c.want)
		}
	}
}
