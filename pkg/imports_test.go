// Package pkg_test checks the rule that the SDK under pkg/ depends only on the
// standard library and itself, so another program can import it without the CLI.
package pkg_test

import (
	"os/exec"
	"strings"
	"testing"
)

func TestTheSDKImportsOnlyTheStandardLibraryAndItself(t *testing.T) {
	out, err := exec.Command("go", "list", "-deps", "-f", "{{if not .Standard}}{{.ImportPath}}{{end}}", "./...").CombinedOutput()
	if err != nil {
		t.Fatalf("go list: %v\n%s", err, out)
	}
	const sdk = "github.com/Cloverhound/webex-cli/pkg"
	for _, path := range strings.Fields(string(out)) {
		if path != sdk && !strings.HasPrefix(path, sdk+"/") {
			t.Errorf("pkg/ depends on %s", path)
		}
	}
}
