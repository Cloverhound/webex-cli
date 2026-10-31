package cmd

import (
	"testing"

	"github.com/Cloverhound/webex-cli/internal/auth"
)

func TestMarkReadOnly(t *testing.T) {
	tok := auth.StoredToken{Scopes: "spark:people_read spark:kms"}
	if err := markReadOnly(&tok, "ro-client", "ro-secret"); err != nil {
		t.Fatal(err)
	}
	if !tok.ReadOnly || tok.ClientID != "ro-client" || tok.ClientSecret != "ro-secret" {
		t.Errorf("token = %+v, want read-only and tied to the client that issued it", tok)
	}

	write := auth.StoredToken{Scopes: "spark:all"}
	if err := markReadOnly(&write, "ro-client", "ro-secret"); err == nil || write.ReadOnly {
		t.Errorf("write scope accepted: err = %v, read-only = %v", err, write.ReadOnly)
	}
}

func TestChooseLoginMode(t *testing.T) {
	cases := []struct {
		name            string
		device, browser bool
		env, headless   string
		want            string
		wantErr         bool
	}{
		{name: "desktop default", want: "auto"},
		{name: "headless default", headless: "SSH session", want: "device"},
		{name: "device flag", device: true, want: "device"},
		{name: "browser flag beats headless", browser: true, headless: "SSH session", want: "browser"},
		{name: "env device", env: "device", want: "device"},
		{name: "env browser beats headless", env: "Browser", headless: "CI environment", want: "browser"},
		{name: "flag beats env", device: true, env: "browser", want: "device"},
		{name: "both flags", device: true, browser: true, wantErr: true},
		{name: "bad env", env: "qr", wantErr: true},
	}
	for _, c := range cases {
		got, _, err := chooseLoginMode(c.device, c.browser, c.env, c.headless)
		if (err != nil) != c.wantErr {
			t.Errorf("%s: err = %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s: mode = %q, want %q", c.name, got, c.want)
		}
	}
}
