// Package readonly classifies OAuth scopes and CLI actions for read-only mode.
package readonly

import (
	"fmt"
	"strings"
)

// DefaultScopes are requested by `webex login --read-only` when the build does
// not inject appconfig.DefaultReadOnlyScopes. Every entry must be registered on
// the OAuth integration, or Webex rejects the login with invalid_scope.
//
// Webex also rejects some valid scope sets unless they follow the order of the
// integration's own scope list (the authorize URL shown in the developer
// portal), so keep this list in that order rather than grouped or sorted.
var DefaultScopes = strings.Join([]string{
	"meeting:admin_preferences_read",
	"spark:people_read",
	"analytics:read_all",
	"meeting:admin_participants_read",
	"spark:messages_read",
	"spark-admin:places_read",
	"spark-admin:workspace_metrics_read",
	"spark:rooms_read",
	"identity:groups_read",
	"spark-admin:recordings_read",
	"spark-admin:telephony_pstn_read",
	"spark-admin:workspace_locations_read",
	"cjp:config_read",
	"spark-admin:call_qualities_read",
	"spark:kms",
	"spark-admin:messages_read",
	"spark-admin:reports_read",
	"meeting:admin_config_read",
	"spark-admin:people_read",
	"meeting:admin_transcripts_read",
	"spark-admin:resource_groups_read",
	"meeting:recordings_read",
	"spark-admin:locations_read",
	"meeting:participants_read",
	"spark-admin:organizations_read",
	"meeting:admin_recordings_read",
	"meeting:transcripts_read",
	"spark:xapi_statuses",
	"spark:memberships_read",
	"spark-admin:calling_cdr_read",
	"identity:organizations_read",
	"spark-admin:devices_read",
	"spark-admin:hybrid_clusters_read",
	"spark-admin:telephony_config_read",
	"meeting:admin_schedule_read",
	"spark:team_memberships_read",
	"meeting:schedules_read",
	"spark-admin:roles_read",
	"meeting:preferences_read",
	"identity:people_read",
	"spark-admin:workspaces_read",
	"spark-admin:resource_group_memberships_read",
	"audit:events_read",
	"spark-admin:hybrid_connectors_read",
	"spark:teams_read",
	"spark-admin:licenses_read",
}, " ")

// Read-only scopes whose names lack a _read suffix. spark:kms grants the
// encryption keys that reading message content requires; spark:xapi_statuses
// reads device status. Neither allows changes on its own.
var readScopeExceptions = map[string]bool{
	"spark:kms":           true,
	"spark:xapi_statuses": true,
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
