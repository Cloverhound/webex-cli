// Package readonly classifies OAuth scopes and CLI actions for read-only mode.
package readonly

import (
	"fmt"
	"strings"
)

// DefaultScopes are requested by `webex login --read-only`. Every entry must be
// registered on the OAuth integration, or Webex rejects the login with invalid_scope.
var DefaultScopes = strings.Join([]string{
	"spark:kms",
	"spark:people_read",
	"spark:rooms_read",
	"spark:memberships_read",
	"spark:messages_read",
	"spark:teams_read",
	"spark:team_memberships_read",
	"spark:devices_read",
	"spark:recordings_read",
	"spark:telephony_config_read",
	"spark:workspaces_read",
	"spark-admin:people_read",
	"spark-admin:organizations_read",
	"spark-admin:licenses_read",
	"spark-admin:roles_read",
	"spark-admin:locations_read",
	"spark-admin:devices_read",
	"spark-admin:workspaces_read",
	"spark-admin:workspace_locations_read",
	"spark-admin:telephony_config_read",
	"spark-admin:telephony_pstn_read",
	"spark-admin:calling_cdr_read",
	"spark-admin:recordings_read",
	"spark-admin:reports_read",
	"spark-admin:hybrid_clusters_read",
	"spark-admin:hybrid_connectors_read",
	"meeting:schedules_read",
	"meeting:recordings_read",
	"meeting:transcripts_read",
	"meeting:participants_read",
	"meeting:preferences_read",
	"meeting:admin_schedule_read",
	"meeting:admin_recordings_read",
	"meeting:admin_participants_read",
	"meeting:admin_transcripts_read",
	"cjp:config_read",
	"analytics:read_all",
	"audit:events_read",
}, " ")

// spark:kms grants access to message encryption keys, which reading message
// content requires; it allows no changes on its own.
var readScopeExceptions = map[string]bool{
	"spark:kms": true,
}

// IsReadScope reports whether an OAuth scope grants read access only.
func IsReadScope(scope string) bool {
	if readScopeExceptions[scope] {
		return true
	}
	return strings.HasSuffix(scope, "_read") ||
		strings.HasSuffix(scope, ":read") ||
		strings.HasSuffix(scope, ":read_all")
}

// ValidateScopes returns an error naming every scope in a space-separated list
// that is not read-only.
func ValidateScopes(scopes string) error {
	fields := strings.Fields(scopes)
	if len(fields) == 0 {
		return fmt.Errorf("no read-only scopes configured")
	}
	var bad []string
	for _, s := range fields {
		if !IsReadScope(s) {
			bad = append(bad, s)
		}
	}
	if len(bad) > 0 {
		return fmt.Errorf("not read-only scopes: %s", strings.Join(bad, ", "))
	}
	return nil
}

// readActions are command verbs (and their hyphenated variants) that only read data.
var readActions = []string{
	"list",
	"get",
	"download",
	"export",
	"search",
	"status",
	"describe",
	"query",
	"show",
	"fetch",
}

// IsReadAction reports whether a command verb such as "list" or "get-details" only reads data.
func IsReadAction(action string) bool {
	for _, p := range readActions {
		if action == p || strings.HasPrefix(action, p+"-") {
			return true
		}
	}
	return false
}
