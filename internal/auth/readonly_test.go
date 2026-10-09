package auth

import (
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
