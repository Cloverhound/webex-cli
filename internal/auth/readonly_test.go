package auth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/zalando/go-keyring"
)

func testConfig(emails ...string) *appconfig.Config {
	cfg := &appconfig.Config{Users: map[string]appconfig.UserInfo{}}
	for _, e := range emails {
		cfg.AddUser(e, e, "org", "Org")
	}
	return cfg
}

func saveTestToken(t *testing.T, email string, readOnly bool) {
	t.Helper()
	tok := &StoredToken{
		AccessToken: "tok-" + email,
		ExpiresAt:   time.Now().Add(time.Hour),
		IssuedAt:    time.Now(),
		ReadOnly:    readOnly,
	}
	if err := SaveToken(email, tok); err != nil {
		t.Fatal(err)
	}
}

func TestPurgeWriteTokens(t *testing.T) {
	keyring.MockInit()
	cfg := testConfig("ro@example.com", "rw@example.com", "rw2@example.com", "gone@example.com")
	cfg.SetDefaultUser("rw@example.com")
	saveTestToken(t, "ro@example.com", true)
	saveTestToken(t, "rw@example.com", false)
	saveTestToken(t, "rw2@example.com", false)

	removed, err := PurgeWriteTokens(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(removed, ",") != "rw2@example.com,rw@example.com" {
		t.Errorf("removed = %v", removed)
	}
	for _, e := range removed {
		if _, err := LoadToken(e); err == nil {
			t.Errorf("token for %s still in keyring", e)
		}
		if _, ok := cfg.Users[e]; ok {
			t.Errorf("%s still in config", e)
		}
	}
	if cfg.DefaultUser != "" {
		t.Errorf("default user = %q, want cleared", cfg.DefaultUser)
	}
	if !HasReadOnlyToken("ro@example.com") {
		t.Error("read-only token was removed")
	}
	if _, ok := cfg.Users["gone@example.com"]; !ok {
		t.Error("user without a token should be left alone")
	}
}

func TestPurgeWriteTokensDeletesUnreadableEntries(t *testing.T) {
	keyring.MockInit()
	cfg := testConfig("corrupt@example.com")
	if err := keyring.Set(serviceName, "corrupt@example.com", "not json"); err != nil {
		t.Fatal(err)
	}

	removed, err := PurgeWriteTokens(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(removed, ",") != "corrupt@example.com" {
		t.Errorf("removed = %v", removed)
	}
	if _, err := keyring.Get(serviceName, "corrupt@example.com"); !errors.Is(err, keyring.ErrNotFound) {
		t.Errorf("unreadable entry still in keyring: %v", err)
	}
}

func TestPurgeWriteTokensChecksEveryStore(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(func() { os.Remove(CredentialsPath()) })
	t.Setenv(TokenStoreEnv, StoreFile)
	if err := keyring.Set(serviceName, "ro@example.com", `{"access_token":"write"}`); err != nil {
		t.Fatal(err)
	}
	saveTestToken(t, "ro@example.com", true)
	cfg := testConfig("ro@example.com")

	if _, err := PurgeWriteTokens(cfg); err != nil {
		t.Fatal(err)
	}
	if _, err := keyring.Get(serviceName, "ro@example.com"); !errors.Is(err, keyring.ErrNotFound) {
		t.Errorf("write-capable keyring copy survived: %v", err)
	}
	if !HasReadOnlyToken("ro@example.com") {
		t.Error("read-only file copy was removed")
	}
	if _, ok := cfg.Users["ro@example.com"]; !ok {
		t.Error("user with a read-only login was removed from config")
	}
}

func TestPurgeWriteTokensDeletesEnvRefreshCache(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(func() { os.Remove(CredentialsPath()) })
	cached := envRefreshKey("env-rt")
	if err := fileSet(cached, `{"token":{"refresh_token":"rotated"}}`); err != nil {
		t.Fatal(err)
	}
	if err := fileSet("file@example.com", `{"read_only":true}`); err != nil {
		t.Fatal(err)
	}

	if _, err := PurgeWriteTokens(testConfig("file@example.com")); err != nil {
		t.Fatal(err)
	}
	if _, err := fileGet(cached); !errors.Is(err, ErrTokenNotFound) {
		t.Errorf("env refresh cache still in the credentials file: %v", err)
	}
	if !HasReadOnlyToken("file@example.com") {
		t.Error("read-only login in the credentials file was removed")
	}
}

func TestPurgeWriteTokensSkipsUnreadableKeyring(t *testing.T) {
	noKeyring(t)
	defer keyring.MockInit()
	if err := fileSet("ro@example.com", `{"read_only":true}`); err != nil {
		t.Fatal(err)
	}
	if err := fileSet("rw@example.com", `{"access_token":"write"}`); err != nil {
		t.Fatal(err)
	}
	cfg := testConfig("ro@example.com", "rw@example.com")

	removed, err := PurgeWriteTokens(cfg)
	if err != nil {
		t.Fatalf("purge failed without a keyring: %v", err)
	}
	if strings.Join(removed, ",") != "rw@example.com" {
		t.Errorf("removed = %v", removed)
	}
	if _, err := fileGet("rw@example.com"); !errors.Is(err, ErrTokenNotFound) {
		t.Errorf("write-capable file copy survived: %v", err)
	}
	if !HasReadOnlyToken("ro@example.com") {
		t.Error("read-only file copy was removed")
	}
}

func TestPurgeWriteTokensRemovesUnlistedLogins(t *testing.T) {
	keyring.MockInit()
	t.Cleanup(func() { os.Remove(CredentialsPath()) })
	if err := keyring.Set(serviceName, "orphan@example.com", `{"access_token":"write"}`); err != nil {
		t.Fatal(err)
	}
	if err := fileSet("orphan@example.com", `{"access_token":"write"}`); err != nil {
		t.Fatal(err)
	}
	saveTestToken(t, "ro@example.com", true)

	if _, err := PurgeWriteTokens(testConfig("ro@example.com")); err != nil {
		t.Fatal(err)
	}
	if _, err := keyring.Get(serviceName, "orphan@example.com"); !errors.Is(err, keyring.ErrNotFound) {
		t.Errorf("unlisted keyring login survived: %v", err)
	}
	if _, err := fileGet("orphan@example.com"); !errors.Is(err, ErrTokenNotFound) {
		t.Errorf("unlisted file login survived: %v", err)
	}
	if !HasReadOnlyToken("ro@example.com") {
		t.Error("read-only login was removed")
	}
}

func TestPurgeWriteTokensWaitsForRefresh(t *testing.T) {
	keyring.MockInit()
	started, release := make(chan struct{}), make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		<-release
		w.Write([]byte(`{"access_token":"renewed","refresh_token":"next","expires_in":3600}`))
	}))
	defer srv.Close()
	old := TokenURL
	TokenURL = srv.URL
	defer func() { TokenURL = old }()
	cfg := testConfig("rw@example.com")
	if err := SaveToken("rw@example.com", &StoredToken{AccessToken: "old", RefreshToken: "rt", IssuedAt: time.Now()}); err != nil {
		t.Fatal(err)
	}

	refreshed := make(chan error, 1)
	go func() { _, err := refreshStored("rw@example.com", cfg, "old"); refreshed <- err }()
	<-started
	purged := make(chan error, 1)
	go func() { _, err := PurgeWriteTokens(cfg); purged <- err }()
	time.Sleep(100 * time.Millisecond)
	close(release)

	if err := <-refreshed; err != nil {
		t.Fatal(err)
	}
	if err := <-purged; err != nil {
		t.Fatal(err)
	}
	if _, err := LoadToken("rw@example.com"); err == nil {
		t.Error("a refresh in progress saved a write-capable token after the purge")
	}
}

func TestKeyringExpected(t *testing.T) {
	runDir := t.TempDir()
	env := func(kv ...string) func(string) string {
		return func(k string) string {
			for i := 0; i < len(kv); i += 2 {
				if kv[i] == k {
					return kv[i+1]
				}
			}
			return ""
		}
	}
	cases := []struct {
		name   string
		goos   string
		getenv func(string) string
		want   bool
	}{
		{"mac", "darwin", env(), true},
		{"windows", "windows", env(), true},
		{"linux sandbox", "linux", env(), false},
		{"linux session bus", "linux", env("DBUS_SESSION_BUS_ADDRESS", "unix:path=/x"), true},
		{"file store", "darwin", env(TokenStoreEnv, "file"), false},
	}
	for _, c := range cases {
		if got := keyringExpected(c.goos, c.getenv, runDir); got != c.want {
			t.Errorf("%s: keyringExpected = %v, want %v", c.name, got, c.want)
		}
	}
	if err := os.WriteFile(filepath.Join(runDir, "bus"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	if !keyringExpected("linux", env(), runDir) {
		t.Error("linux with a systemd user bus: keyringExpected = false")
	}
}

func TestResolveTokenReadOnly(t *testing.T) {
	keyring.MockInit()
	cfg := testConfig("ro@example.com", "rw@example.com")
	cfg.SetDefaultUser("ro@example.com")
	saveTestToken(t, "ro@example.com", true)
	saveTestToken(t, "rw@example.com", false)

	if _, err := ResolveToken("flag-token", "", "", "", cfg, true); err == nil {
		t.Error("--token accepted in read-only mode")
	}
	if _, err := ResolveToken("", "env-token", "", "", cfg, true); err == nil {
		t.Error("$WEBEX_TOKEN accepted in read-only mode")
	}
	if _, err := ResolveToken("", "", "rw@example.com", "", cfg, true); err == nil {
		t.Error("write-capable user accepted in read-only mode")
	}

	res, err := ResolveToken("", "", "", "", cfg, true)
	if err != nil {
		t.Fatal(err)
	}
	if res.Token != "tok-ro@example.com" {
		t.Errorf("token = %q", res.Token)
	}

	if _, err := ResolveToken("", "", "rw@example.com", "", cfg, false); err != nil {
		t.Errorf("write-capable user rejected outside read-only mode: %v", err)
	}
}
