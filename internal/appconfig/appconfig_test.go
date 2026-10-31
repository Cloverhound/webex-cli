package appconfig

import "testing"

func TestEffectiveReadOnlyClient(t *testing.T) {
	oldID, oldSecret, oldROID, oldROSecret := DefaultClientID, DefaultClientSecret, DefaultReadOnlyClientID, DefaultReadOnlyClientSecret
	defer func() {
		DefaultClientID, DefaultClientSecret, DefaultReadOnlyClientID, DefaultReadOnlyClientSecret = oldID, oldSecret, oldROID, oldROSecret
	}()
	DefaultClientID, DefaultClientSecret = "main", "main-secret"

	cases := []struct {
		name               string
		readOnlyID, secret string
		cfg                Config
		wantID, wantSecret string
	}{
		{name: "read-only integration", readOnlyID: "ro", secret: "ro-secret", wantID: "ro", wantSecret: "ro-secret"},
		{name: "build without one", wantID: "main", wantSecret: "main-secret"},
		{name: "configured integration wins", readOnlyID: "ro", secret: "ro-secret",
			cfg: Config{ClientID: "mine", ClientSecret: "my-secret"}, wantID: "mine", wantSecret: "my-secret"},
	}
	for _, c := range cases {
		DefaultReadOnlyClientID, DefaultReadOnlyClientSecret = c.readOnlyID, c.secret
		id, secret := c.cfg.EffectiveReadOnlyClient()
		if id != c.wantID || secret != c.wantSecret {
			t.Errorf("%s: got %q/%q, want %q/%q", c.name, id, secret, c.wantID, c.wantSecret)
		}
	}
}
