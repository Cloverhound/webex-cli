// Package readonly classifies OAuth scopes and CLI actions for read-only mode.
package readonly

import (
	"fmt"
	"slices"
	"strings"
)

// DefaultScopes are requested by `webex login --read-only` when the build does
// not inject appconfig.DefaultReadOnlyScopes. Every entry must be registered on
// the OAuth integration, or Webex rejects the login with invalid_scope.
// spark:{devices,recordings,telephony_config,workspaces}_read are not registered
// on the release integration; their spark-admin counterparts cover those reads.
var DefaultScopes = strings.Join([]string{
	"spark:kms",
	"spark:xapi_statuses",
	"spark:people_read",
	"spark:rooms_read",
	"spark:memberships_read",
	"spark:messages_read",
	"spark:teams_read",
	"spark:team_memberships_read",
	"spark-admin:people_read",
	"spark-admin:organizations_read",
	"spark-admin:licenses_read",
	"spark-admin:roles_read",
	"spark-admin:locations_read",
	"spark-admin:places_read",
	"spark-admin:devices_read",
	"spark-admin:workspaces_read",
	"spark-admin:workspace_locations_read",
	"spark-admin:workspace_metrics_read",
	"spark-admin:telephony_config_read",
	"spark-admin:telephony_pstn_read",
	"spark-admin:calling_cdr_read",
	"spark-admin:call_qualities_read",
	"spark-admin:recordings_read",
	"spark-admin:reports_read",
	"spark-admin:messages_read",
	"spark-admin:hybrid_clusters_read",
	"spark-admin:hybrid_connectors_read",
	"spark-admin:resource_groups_read",
	"spark-admin:resource_group_memberships_read",
	"meeting:schedules_read",
	"meeting:recordings_read",
	"meeting:transcripts_read",
	"meeting:participants_read",
	"meeting:preferences_read",
	"meeting:admin_schedule_read",
	"meeting:admin_recordings_read",
	"meeting:admin_participants_read",
	"meeting:admin_transcripts_read",
	"meeting:admin_preferences_read",
	"meeting:admin_config_read",
	"identity:people_read",
	"identity:organizations_read",
	"identity:groups_read",
	"cjp:config_read",
	"analytics:read_all",
	"audit:events_read",
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

// queryPOSTPaths are POST endpoints that only read: Webex takes their filters
// in the request body. Read-only mode blocks every other POST, because a
// command's verb does not prove an endpoint is safe (`calling call-controls
// get` POSTs to /telephony/calls/retrieve, which picks up a parked call).
var queryPOSTPaths = []string{
	"/admin/recordings/query",
	"/recordings/query",
	"/generated-summaries/search",
	"/livemonitoring/liveMeetingsByCountry",
	"/meetingParticipants/query",
	"/meetings/{meetingId}/registrants/query",
	"/meetings/{meetingId}/surveyLinks",
	"/organization/{orgid}/contact-service-queue/fetch-by-dynamic-skills-and-skillProfile",
	"/organization/{orgid}/contact-service-queue/fetch-by-userId-skillProfileId",
	"/organization/{orgid}/contact-service-queue/fetch-manually-assignable-queues",
	"/organization/{orgid}/user/fetch-by-skill-requirements",
	"/organization/{orgid}/user/fetch-user-details-by-ids",
	"/organization/{orgid}/v2/contact-service-queue/fetch-by-grouped-assistant-skill",
	"/search",
	"/summary/list",
	"/telephony/config/lists/devices/dynamicSettings/actions/getSettings/invoke",
	"/telephony/config/lists/devices/{deviceId}/dynamicSettings/actions/getSettings/invoke",
	"/telephony/config/lists/locations/{locationId}/devices/dynamicSettings/actions/getSettings/invoke",
	"/v1/captures/query",
	"/v1/{orgid}/functions/{id}:export",
}

// IsQueryPOST reports whether a POST to the path template only reads data.
func IsQueryPOST(path string) bool {
	return slices.Contains(queryPOSTPaths, path)
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
