package apicmd

import (
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

func TestHelpShowsTheExampleOfABodyFlagsCannotGive(t *testing.T) {
	op := &webexops.Operation{
		Command: "create", Summary: "Create a thing", Description: "Create a thing in an organization.",
		Method: "POST", Path: "/things",
		Body: &webexops.Body{Example: "{\n  \"settings\": [\n    \"<string>\"\n  ]\n}"},
	}
	want := "Create a thing in an organization.\n\n" + bodyExampleHeading + "\n\n" +
		"  {\n    \"settings\": [\n      \"<string>\"\n    ]\n  }"
	if got := NewCommand(op).Long; got != want {
		t.Fatalf("help:\n%s\nwant:\n%s", got, want)
	}
}

func TestHelpOfAPrivateOperationEndsWithItsNotice(t *testing.T) {
	op := &webexops.Operation{
		Command: "create", Method: "POST", Path: "/api/v1/things/", Private: true,
		Body: &webexops.Body{Example: "{}"},
	}
	long := NewCommand(op).Long
	if !strings.HasPrefix(long, bodyExampleHeading) || !strings.HasSuffix(long, privateAPINote) {
		t.Fatalf("help:\n%s", long)
	}
}

func TestHelpOfAFlatBodyIsItsDescription(t *testing.T) {
	op := &webexops.Operation{
		Command: "create", Description: "Create a thing.", Method: "POST", Path: "/things",
		Body: &webexops.Body{Fields: []webexops.BodyField{{Name: "name", Flag: "name", Type: "string"}}},
	}
	if got := NewCommand(op).Long; got != "Create a thing." {
		t.Fatalf("help: %q", got)
	}
}
