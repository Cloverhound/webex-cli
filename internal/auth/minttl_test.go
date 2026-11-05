package auth

import (
	"strings"
	"testing"
	"time"
)

func saveExpiring(t *testing.T, email string, in time.Duration) {
	t.Helper()
	tok := &StoredToken{AccessToken: "old", RefreshToken: "rt-0", ExpiresAt: time.Now().Add(in), IssuedAt: time.Now()}
	if err := SaveToken(email, tok); err != nil {
		t.Fatal(err)
	}
}

func TestEnsureMinTTLKeepsATokenWithEnoughLife(t *testing.T) {
	noKeyring(t)
	calls := fakeTokenEndpoint(t)
	saveExpiring(t, "a@example.com", 2*time.Hour)

	tok, err := EnsureMinTTL("a@example.com", time.Hour, testConfig("a@example.com"))
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "old" || *calls != 0 {
		t.Fatalf("token=%q refreshes=%d", tok.AccessToken, *calls)
	}
}

func TestEnsureMinTTLRefreshesAShortToken(t *testing.T) {
	noKeyring(t)
	calls := fakeTokenEndpoint(t)
	saveExpiring(t, "a@example.com", 10*time.Minute)

	tok, err := EnsureMinTTL("a@example.com", 30*time.Minute, testConfig("a@example.com"))
	if err != nil {
		t.Fatal(err)
	}
	if tok.AccessToken != "at-1" || *calls != 1 {
		t.Fatalf("token=%q refreshes=%d", tok.AccessToken, *calls)
	}
}

func TestEnsureMinTTLFailsWhenRefreshCannotProvideEnoughLife(t *testing.T) {
	noKeyring(t)
	fakeTokenEndpoint(t)
	saveExpiring(t, "a@example.com", 10*time.Minute)

	_, err := EnsureMinTTL("a@example.com", 2*time.Hour, testConfig("a@example.com"))
	if err == nil || !strings.Contains(err.Error(), "shorter than the requested") {
		t.Fatalf("err = %v", err)
	}
}
