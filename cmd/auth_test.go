package cmd

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/auth"
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
