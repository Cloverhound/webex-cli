package readonly

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsReadScope(t *testing.T) {
	cases := map[string]bool{
		"spark-admin:people_read":         true,
		"spark:kms":                       true,
		"spark:xapi_statuses":             true,
		"spark:xapi_commands":             false,
		"cjp:user":                        false,
		"analytics:read_all":              true,
		"cjp-hybrid-conn:read":            true,
		"spark:all":                       false,
		"spark-admin:all":                 false,
		"spark-admin:people_write":        false,
		"cjp:config":                      false,
		"identity:people_rw":              false,
		"spark-admin:telephony_read_only": false,
	}
	for scope, want := range cases {
		if got := IsReadScope(scope); got != want {
			t.Errorf("IsReadScope(%q) = %v, want %v", scope, got, want)
		}
	}
}

func TestDefaultScopesAreReadOnly(t *testing.T) {
	if err := ValidateScopes(DefaultScopes); err != nil {
		t.Fatal(err)
	}
}

func TestValidateScopesRejectsWriteScopes(t *testing.T) {
	err := ValidateScopes("spark:people_read spark:all cjp:config_write")
	if err == nil {
		t.Fatal("expected error")
	}
	want := "not read-only scopes: spark:all, cjp:config_write"
	if err.Error() != want {
		t.Errorf("got %q, want %q", err, want)
	}
}

func TestValidateScopesRejectsEmpty(t *testing.T) {
	if err := ValidateScopes("  "); err == nil {
		t.Fatal("expected error")
	}
}

func TestIsQueryPOST(t *testing.T) {
	queries := []string{
		"/search",
		"/recordings/query",
		"/meetings/{meetingId}/registrants/query",
		"/v1/{orgid}/functions/{id}:export",
	}
	writes := []string{
		"/telephony/calls/retrieve",
		"/event",
		"/identity/organizations/{orgId}/actions/getDomainVerificationToken",
		"/rooms",
		"/meetings/{meetingId}/registrants/query/",
	}
	for _, path := range queries {
		if !IsQueryPOST(path) {
			t.Errorf("IsQueryPOST(%q) = false, want true", path)
		}
	}
	for _, path := range writes {
		if IsQueryPOST(path) {
			t.Errorf("IsQueryPOST(%q) = true, want false", path)
		}
	}
}

// Query paths must match a POST request template exactly, or read-only mode
// blocks the command. Regenerated commands can rename path parameters.
func TestQueryPOSTPathsMatchRequests(t *testing.T) {
	var src strings.Builder
	for _, dir := range []string{"../../cmd", "../../internal"} {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return err
			}
			b, err := os.ReadFile(path)
			src.Write(b)
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, p := range queryPOSTPaths {
		if !strings.Contains(src.String(), `"POST", "`+p+`")`) {
			t.Errorf("no POST request uses %q", p)
		}
	}
}

func TestIsReadAction(t *testing.T) {
	cases := map[string]bool{
		"list":                 true,
		"list-by-ids":          true,
		"get-dynamic-settings": true,
		"search":               true,
		"create":               false,
		"update":               false,
		"delete":               false,
		"listen":               false,
		"getaway":              false,
		"assign-phone-number":  false,
	}
	for action, want := range cases {
		if got := IsReadAction(action); got != want {
			t.Errorf("IsReadAction(%q) = %v, want %v", action, got, want)
		}
	}
}
