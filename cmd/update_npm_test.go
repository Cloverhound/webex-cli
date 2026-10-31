package cmd

import "testing"

func TestInstalledByNPM(t *testing.T) {
	cases := map[string]bool{
		"/usr/local/lib/node_modules/@cloverhound/webex-cli-darwin-arm64/bin/webex":                  true,
		`C:\Users\a\AppData\Roaming\npm\node_modules\@cloverhound\webex-cli-win32-x64\bin\webex.exe`: true,
		"/home/a/.npm/_npx/1a2b/node_modules/@cloverhound/webex-cli-linux-x64/bin/webex":             true,
		"/home/a/.local/bin/webex":      false,
		"/opt/node_modules/other/webex": false,
	}
	for path, want := range cases {
		if got := installedByNPM(path); got != want {
			t.Errorf("installedByNPM(%q) = %v, want %v", path, got, want)
		}
	}
}
