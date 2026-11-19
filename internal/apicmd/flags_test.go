package apicmd

import (
	"testing"

	"github.com/Cloverhound/webex-cli/pkg/webexops"
	"github.com/spf13/pflag"
)

func TestAFlagWithCodeInItsUsageShowsItsType(t *testing.T) {
	op := &webexops.Operation{
		Command: "get-input-choices",
		Method:  "GET",
		Path:    "/choices",
		Params: []webexops.Param{{
			Name: "search", In: webexops.InQuery, Flag: "search",
			Usage: "Search text. When `validate=true`, this is the single value to look up in `search`.",
		}},
	}
	flag := NewCommand(op).Flags().Lookup("search")

	name, usage := pflag.UnquoteUsage(flag)
	if name != "string" {
		t.Errorf("help shows --search %s, want --search string", name)
	}
	if want := "Search text. When validate=true, this is the single value to look up in search."; usage != want {
		t.Errorf("usage = %q, want %q", usage, want)
	}
}

// Every API command is built at startup, so flag help must stay cheap to make.
func TestFlagUsageAllocatesOnlyItsResult(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() { flagUsage("A `code` span\nover two lines.") })
	if allocs > 1 {
		t.Errorf("flagUsage made %v allocations, want at most 1", allocs)
	}
}
