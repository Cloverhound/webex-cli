package cmd

import "testing"

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
