package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
	"github.com/zalando/go-keyring"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	fn()
	os.Stdout = old
	w.Close()
	out, _ := io.ReadAll(r)
	return string(out)
}

func TestAuthStatusReportsEnvTokenWithoutStoredLogin(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv(readOnlyEnv, "")
	t.Setenv("WEBEX_TOKEN", "")
	t.Setenv(auth.RefreshTokenEnv, "rt")

	out := captureStdout(t, func() {
		if err := authStatusCmd.RunE(authStatusCmd, nil); err != nil {
			t.Error(err)
		}
	})
	if !strings.Contains(out, "$"+auth.RefreshTokenEnv+" is set") {
		t.Errorf("status output does not mention the env token:\n%s", out)
	}
	if strings.Contains(out, "No authenticated user") {
		t.Errorf("status claims no user although commands authenticate from the env:\n%s", out)
	}
}

func TestAuthStatusShowsTheRequestedUser(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv(readOnlyEnv, "")
	t.Setenv("WEBEX_TOKEN", "")
	t.Setenv(auth.RefreshTokenEnv, "")
	keyring.MockInit()

	cfg, err := appconfig.Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg.DefaultUser = "default@example.com"
	cfg.Users = map[string]appconfig.UserInfo{
		"default@example.com": {OrgID: "default-org", OrgName: "Default Org"},
		"other@example.com":   {OrgID: "other-org", OrgName: "Other Org"},
	}
	if err := cfg.Save(); err != nil {
		t.Fatal(err)
	}
	// Expired tokens with a live refresh token keep status from calling Webex.
	for email := range cfg.Users {
		tok := &auth.StoredToken{AccessToken: "at", RefreshToken: "rt",
			ExpiresAt: time.Now().Add(-time.Hour), IssuedAt: time.Now().Add(-2 * time.Hour)}
		if err := auth.SaveToken(email, tok); err != nil {
			t.Fatal(err)
		}
	}

	status := func(t *testing.T) string {
		return captureStdout(t, func() {
			if err := authStatusCmd.RunE(authStatusCmd, nil); err != nil {
				t.Error(err)
			}
		})
	}

	t.Run("flag", func(t *testing.T) {
		t.Setenv("WEBEX_USER", "")
		if err := authStatusCmd.ParseFlags([]string{"--user", "other@example.com"}); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { authStatusCmd.Flags().Set("user", "") })
		if out := status(t); !strings.Contains(out, "User:    other@example.com") || !strings.Contains(out, "Other Org") {
			t.Errorf("status does not show the --user login:\n%s", out)
		}
	})
	t.Run("env", func(t *testing.T) {
		t.Setenv("WEBEX_USER", "other@example.com")
		if out := status(t); !strings.Contains(out, "User:    other@example.com") {
			t.Errorf("status does not show the $WEBEX_USER login:\n%s", out)
		}
	})
	t.Run("default", func(t *testing.T) {
		t.Setenv("WEBEX_USER", "")
		if out := status(t); !strings.Contains(out, "User:    default@example.com") {
			t.Errorf("status does not show the default login:\n%s", out)
		}
	})
}
