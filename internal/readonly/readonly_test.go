package readonly

import "testing"

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
