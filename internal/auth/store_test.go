package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zalando/go-keyring"
)

// noKeyring simulates a host without a keyring service, such as Linux without D-Bus.
func noKeyring(t *testing.T) {
	t.Helper()
	keyring.MockInitWithError(errors.New("dbus: no session bus"))
	t.Cleanup(func() { os.Remove(CredentialsPath()) })
}

func TestStoreFallsBackToFile(t *testing.T) {
	noKeyring(t)

	tok := &StoredToken{AccessToken: "at", RefreshToken: "rt", ExpiresAt: time.Now().Add(time.Hour)}
	if err := SaveToken("a@example.com", tok); err != nil {
		t.Fatal(err)
	}
	got, store, err := LoadTokenWithStore("a@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if store != StoreFile || got.AccessToken != "at" {
		t.Errorf("store=%s token=%+v", store, got)
	}

	if runtime.GOOS != "windows" {
		info, err := os.Stat(CredentialsPath())
		if err != nil {
			t.Fatal(err)
		}
		if perm := info.Mode().Perm(); perm != 0600 {
			t.Errorf("credentials mode = %o, want 600", perm)
		}
		dirInfo, _ := os.Stat(CredentialsDir())
		if perm := dirInfo.Mode().Perm(); perm != 0700 {
			t.Errorf("credentials dir mode = %o, want 700", perm)
		}
	}

	if err := DeleteToken("a@example.com"); err != nil {
		t.Fatalf("delete with no keyring: %v", err)
	}
	if _, err := LoadToken("a@example.com"); err == nil {
		t.Error("token still loadable after delete")
	}
}

func TestStoreForcedKeyringDoesNotFallBack(t *testing.T) {
	noKeyring(t)
	t.Setenv(TokenStoreEnv, "keyring")

	if err := SaveToken("a@example.com", &StoredToken{AccessToken: "at"}); err == nil {
		t.Fatal("expected keyring error with WEBEX_TOKEN_STORE=keyring")
	}
	if _, err := os.Stat(CredentialsPath()); !os.IsNotExist(err) {
		t.Error("credentials file written despite WEBEX_TOKEN_STORE=keyring")
	}
}

func TestStoreForcedFileSkipsKeyring(t *testing.T) {
	keyring.MockInit()
	t.Setenv(TokenStoreEnv, "file")
	t.Cleanup(func() { os.Remove(CredentialsPath()) })

	if err := SaveToken("a@example.com", &StoredToken{AccessToken: "at"}); err != nil {
		t.Fatal(err)
	}
	if _, err := keyring.Get(serviceName, "a@example.com"); err == nil {
		t.Error("token written to keyring despite WEBEX_TOKEN_STORE=file")
	}
	if _, store, err := LoadTokenWithStore("a@example.com"); err != nil || store != StoreFile {
		t.Errorf("store=%s err=%v", store, err)
	}
}

func TestStoreRejectsUnknownStore(t *testing.T) {
	t.Setenv(TokenStoreEnv, "vault")
	if err := SaveToken("a@example.com", &StoredToken{}); err == nil {
		t.Error("expected error for unknown WEBEX_TOKEN_STORE")
	}
}

// fakeTokenEndpoint counts refresh requests and answers each with a fresh token.
func fakeTokenEndpoint(t *testing.T) *int32 {
	t.Helper()
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			n := atomic.AddInt32(&calls, 1)
			r.ParseForm()
			if r.Form.Get("grant_type") != "refresh_token" {
				t.Errorf("grant_type = %q", r.Form.Get("grant_type"))
			}
			// Hold the response so parallel callers overlap.
			time.Sleep(50 * time.Millisecond)
			json.NewEncoder(w).Encode(map[string]any{
				"access_token":  "at-" + string(rune('0'+n)),
				"refresh_token": "rt-" + string(rune('0'+n)),
				"expires_in":    3600,
			})
		case "/me":
			w.Write([]byte(`{"emails":["svc@example.com"],"orgId":"org-1"}`))
		}
	}))
	oldToken, oldMe := TokenURL, peopleMeURL
	TokenURL, peopleMeURL = srv.URL+"/token", srv.URL+"/me"
	t.Cleanup(func() { TokenURL, peopleMeURL = oldToken, oldMe; srv.Close() })
	return &calls
}

func TestParallelRefreshRefreshesOnce(t *testing.T) {
	noKeyring(t)
	calls := fakeTokenEndpoint(t)

	expired := &StoredToken{AccessToken: "old", RefreshToken: "rt-0", ExpiresAt: time.Now().Add(-time.Hour), IssuedAt: time.Now()}
	if err := SaveToken("a@example.com", expired); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig("a@example.com")

	var wg sync.WaitGroup
	results := make([]string, 4)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tok, err := refreshStored("a@example.com", cfg, "old")
			if err != nil {
				t.Error(err)
				return
			}
			results[i] = tok.AccessToken
		}(i)
	}
	wg.Wait()

	if *calls != 1 {
		t.Errorf("refresh requests = %d, want 1", *calls)
	}
	for _, r := range results {
		if r != "at-1" {
			t.Errorf("results = %v, want all at-1", results)
			break
		}
	}
	stored, _ := LoadToken("a@example.com")
	if stored.RefreshToken != "rt-1" {
		t.Errorf("stored refresh token = %q, want rt-1", stored.RefreshToken)
	}
}

func TestLoadTokenPrefersNewestCopy(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(func() { os.Remove(CredentialsPath()) })
	encode := func(at string, issued time.Time) string {
		data, _ := json.Marshal(StoredToken{AccessToken: at, IssuedAt: issued})
		return string(data)
	}
	if err := keyring.Set(serviceName, "a@example.com", encode("stale", time.Now().Add(-time.Hour))); err != nil {
		t.Fatal(err)
	}
	if err := fileSet("a@example.com", encode("fresh", time.Now())); err != nil {
		t.Fatal(err)
	}

	got, store, err := LoadTokenWithStore("a@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if got.AccessToken != "fresh" || store != StoreFile {
		t.Errorf("loaded %q from %s, want fresh from file", got.AccessToken, store)
	}
}

func TestEnvRefreshKeepsRotatedTokenWithoutCache(t *testing.T) {
	keyring.MockInit()
	// A regular file where the credentials directory belongs makes every write fail.
	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, nil, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", blocker)

	var mu sync.Mutex
	valid := map[string]bool{"env-rt-nocache": true}
	issued := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/me" {
			w.Write([]byte(`{"emails":["svc@example.com"],"orgId":"org-1"}`))
			return
		}
		r.ParseForm()
		mu.Lock()
		defer mu.Unlock()
		rt := r.Form.Get("refresh_token")
		if !valid[rt] {
			http.Error(w, `{"error":"invalid_grant"}`, http.StatusBadRequest)
			return
		}
		// Webex rotates the refresh token, and the old one stops working.
		delete(valid, rt)
		issued++
		next := fmt.Sprintf("rt-%d", issued)
		valid[next] = true
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": fmt.Sprintf("at-%d", issued), "refresh_token": next, "expires_in": 3600,
		})
	}))
	oldToken, oldMe := TokenURL, peopleMeURL
	TokenURL, peopleMeURL = srv.URL+"/token", srv.URL+"/me"
	t.Cleanup(func() { TokenURL, peopleMeURL = oldToken, oldMe; srv.Close() })

	tok, _, _, err := EnvRefreshAccessToken("env-rt-nocache", "cid", "secret", "")
	if err != nil || tok != "at-1" {
		t.Fatalf("first refresh: tok=%s err=%v", tok, err)
	}
	tok, _, _, err = EnvRefreshAccessToken("env-rt-nocache", "cid", "secret", "at-1")
	if err != nil || tok != "at-2" {
		t.Fatalf("refresh after a 401: tok=%s err=%v; want at-2 from the rotated refresh token", tok, err)
	}
}

func TestEnvRefreshCachesAndRotates(t *testing.T) {
	noKeyring(t)
	calls := fakeTokenEndpoint(t)

	tok, email, org, err := EnvRefreshAccessToken("env-rt", "cid", "secret", "")
	if err != nil {
		t.Fatal(err)
	}
	if tok != "at-1" || email != "svc@example.com" || org != "org-1" {
		t.Errorf("got %s %s %s", tok, email, org)
	}

	// A second process reuses the cached access token.
	tok, _, _, err = EnvRefreshAccessToken("env-rt", "cid", "secret", "")
	if err != nil || tok != "at-1" || *calls != 1 {
		t.Fatalf("second call: tok=%s calls=%d err=%v", tok, *calls, err)
	}

	// After a 401, it refreshes with the rotated token, not the env value.
	tok, _, org, err = EnvRefreshAccessToken("env-rt", "cid", "secret", "at-1")
	if err != nil || tok != "at-2" || org != "org-1" {
		t.Fatalf("forced refresh: tok=%s org=%s err=%v", tok, org, err)
	}
	entry := loadEnvRefreshEntry(envRefreshKey("env-rt"))
	if entry == nil || entry.Token.RefreshToken != "rt-2" {
		t.Errorf("cached entry = %+v, want rotated refresh token rt-2", entry)
	}
}
