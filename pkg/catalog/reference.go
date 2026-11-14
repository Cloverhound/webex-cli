package catalog

import (
	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// Everything here is read-only. A type moves into the managed set once its writes
// are proven on a live tenant.
//
// RedSky settings are left out on purpose, because the response carries a secret.

func text(name, path string) resource.Field { return resource.Field{Name: name, Path: path} }

func flag(name, path string) resource.Field {
	return resource.Field{Name: name, Path: path, Kind: resource.Bool}
}

func list(name, path string) resource.Field {
	return resource.Field{Name: name, Path: path, Kind: resource.StringList}
}

var locationField = resource.Field{Name: "location_id", Immutable: true, Ref: "location"}

// calling defines a read-only Calling list.
func calling(name, title, plural string, op *webexops.Operation, key, label string, fields ...resource.Field) resource.Type {
	return resource.Type{Name: name, Title: title, Plural: plural, Area: resource.Calling,
		LabelField: label, Fields: fields, Routes: resource.Routes{List: on(op), ArrayKey: key}}
}

// withLocation groups a Calling type's objects by the location each list item names.
func withLocation(t resource.Type) resource.Type {
	t.ParentField = "location_id"
	t.Routes.ParentKey = "locationId"
	t.Fields = append([]resource.Field{locationField}, t.Fields...)
	return t
}

// perLocation reads a Calling list once for each location.
func perLocation(t resource.Type) resource.Type {
	t.ParentField = "location_id"
	t.Routes.PerLocation = true
	t.Routes.List.Parent = "locationId"
	t.Fields = append([]resource.Field{locationField}, t.Fields...)
	return t
}

// orgSettings defines an organization-wide Calling settings object.
func orgSettings(name, title string, op *webexops.Operation) resource.Type {
	return resource.Type{Name: name, Title: title, Plural: title, Area: resource.Calling,
		Routes: resource.Routes{List: on(op), Singleton: true}}
}

// locationSettings defines a settings object held once per location.
func locationSettings(name, title string, op *webexops.Operation) resource.Type {
	return resource.Type{Name: name, Title: title, Plural: title, Area: resource.Calling,
		ParentField: "location_id", Fields: []resource.Field{locationField},
		Routes: resource.Routes{List: on(op, "locationId"), Singleton: true, PerLocation: true}}
}

// center defines a read-only Contact Center list.
func center(name, title, plural string, list resource.Route, label string, fields ...resource.Field) resource.Type {
	return resource.Type{Name: name, Title: title, Plural: plural, Area: resource.ContactCenter,
		LabelField: label, Fields: fields, Routes: resource.Routes{List: list}}
}

var callingOrganization = []resource.Type{
	func() resource.Type {
		t := calling("announcement_language", "Language", "Languages", webexops.CallingCallingServiceListAnnouncementLanguages, "languages", "name",
			text("name", "name"), text("code", "code"))
		t.Routes.IDKey = "code"
		return t
	}(),
	calling("person", "Person", "People", webexops.CallingPeopleList, "items", "display_name",
		text("display_name", "displayName"), list("emails", "emails"), text("status", "status")),
	func() resource.Type {
		t := calling("phone_number", "Phone Number", "Phone Numbers", webexops.CallingNumbersGetPhoneOrg, "phoneNumbers", "phone_number",
			text("phone_number", "phoneNumber"), text("extension", "extension"),
			text("location", "location.name"), text("state", "state"))
		t.Routes.IDKey = "phoneNumber"
		t.Routes.FallbackIDKey = "extension"
		return t
	}(),
}

var callingRouting = []resource.Type{
	withLocation(calling("voicemail_group", "Voicemail Group", "Voicemail Groups", webexops.CallingLocationVoicemailListVoicemailgroup, "voicemailGroups", "name",
		text("name", "name"), text("extension", "extension"), flag("enabled", "enabled"))),
	withLocation(calling("paging_group", "Paging Group", "Paging Groups", webexops.CallingPagingGroupList, "locationPaging", "name",
		text("name", "name"), text("extension", "extension"))),
	perLocation(calling("call_park", "Call Park", "Call Parks", webexops.CallingCallParkList, "callParks", "name",
		text("name", "name"))),
	withLocation(calling("call_park_extension", "Call Park Extension", "Call Park Extensions", webexops.CallingCallParkListExtensions, "callParkExtensions", "name",
		text("name", "name"), text("extension", "extension"))),
	perLocation(calling("call_pickup", "Call Pickup", "Call Pickups", webexops.CallingCallPickupList, "callPickups", "name",
		text("name", "name"))),
	calling("virtual_line", "Virtual Line", "Virtual Lines", webexops.CallingVirtualLineCallList, "virtualLines", "last_name",
		text("last_name", "lastName"), text("first_name", "firstName"), text("number", "callerIdNumber")),
	calling("virtual_extension", "Virtual Extension", "Virtual Extensions", webexops.CallingVirtualExtensionsList, "virtualExtensions", "display_name",
		text("display_name", "displayName"), text("extension", "extension")),
	calling("virtual_extension_range", "Virtual Extension Range", "Virtual Extension Ranges", webexops.CallingVirtualExtensionsListRange, "virtualExtensionRanges", "name",
		text("name", "name")),
	calling("wrapup_reason", "Wrap-Up Reason", "Wrap-Up Reasons", webexops.CallingCustomerExperienceEssentialsListWrapUpReasons, "wrapupReasons", "name",
		text("name", "name")),
}

var callingSchedules = []resource.Type{
	perLocation(calling("schedule", "Schedule", "Schedules", webexops.CallingLocationSchedulesList, "schedules", "name",
		text("name", "name"), text("type", "type"))),
	calling("operating_mode", "Operating Mode", "Operating Modes", webexops.CallingOperatingModesList, "operatingModes", "name",
		text("name", "name"), text("type", "type"), text("level", "level")),
	calling("announcement", "Announcement", "Announcements", webexops.CallingAnnouncementRepositoryListGreetings, "announcements", "name",
		text("name", "name"), text("file_name", "fileName"), text("level", "level")),
	calling("announcement_playlist", "Playlist", "Playlists", webexops.CallingAnnouncementPlaylistList, "playlists", "name",
		text("name", "name")),
}

var callingPSTN = []resource.Type{
	calling("trunk", "Trunk", "Trunks", webexops.CallingCallRoutingListTrunks, "trunks", "name",
		text("name", "name"), text("trunk_type", "trunkType"), text("location", "location.name"), flag("in_use", "inUse")),
	calling("route_group", "Route Group", "Route Groups", webexops.CallingCallRoutingListGroups, "routeGroups", "name",
		text("name", "name"), flag("in_use", "inUse")),
	calling("route_list", "Route List", "Route Lists", webexops.CallingCallRoutingListRouteLists, "routeLists", "name",
		text("name", "name")),
	calling("dial_plan", "Dial Plan", "Dial Plans", webexops.CallingCallRoutingListDialPlans, "dialPlans", "name",
		text("name", "name"), text("route_name", "routeName"), text("route_type", "routeType")),
	calling("translation_pattern", "Translation Pattern", "Translation Patterns", webexops.CallingCallRoutingListTranslationPatterns, "translationPatterns", "name",
		text("name", "name"), text("matching_pattern", "matchingPattern"),
		text("replacement_pattern", "replacementPattern"), text("level", "level")),
}

var callingDevices = []resource.Type{
	calling("device", "Device", "Devices", webexops.CallingDevicesList, "items", "display_name",
		text("display_name", "displayName"), text("product", "product"), text("connection_status", "connectionStatus")),
	calling("workspace", "Workspace", "Workspaces", webexops.CallingWorkspacesList, "items", "display_name",
		text("display_name", "displayName"), text("type", "type")),
	calling("line_key_template", "Line Key Template", "Line Key Templates", webexops.CallingDeviceCallListLineKeyTemplates, "lineKeyTemplates", "template_name",
		text("template_name", "templateName"), text("model", "modelDisplayName")),
	calling("dect_network", "DECT Network", "DECT Networks", webexops.CallingDectDevicesListNetworks, "dectNetworks", "name",
		text("name", "name")),
}

var callingLocationSettings = []resource.Type{
	locationSettings("location_voicemail", "Voicemail", webexops.CallingLocationVoicemailGet),
	locationSettings("location_voice_portal", "Voice Portal", webexops.CallingLocationVoicemailGetVoiceportal),
	locationSettings("location_music_on_hold", "Music on Hold", webexops.CallingLocationCallGetMusicHold),
	locationSettings("location_internal_dialing", "Internal Dialing", webexops.CallingLocationCallHandlingGetInternalDialingConfiguration),
	locationSettings("location_intercept", "Call Intercept", webexops.CallingLocationCallHandlingGetIntercept),
	locationSettings("location_emergency_callback", "Emergency Callback", webexops.CallingLocationCallGetEmergencyCallback),
	locationSettings("location_emergency_notification", "Emergency Notification", webexops.CallingEmergencyServicesGetLocationCallNotification),
	locationSettings("location_call_captions", "Call Captions", webexops.CallingLocationCallGetCaptionsSettings),
	locationSettings("location_call_park_settings", "Call Park Settings", webexops.CallingCallParkGet2),
	locationSettings("location_recording_announcement", "Recording Announcement", webexops.CallingCallRecordingGetLocationComplianceAnnouncement),
}

var callingOrgSettings = []resource.Type{
	orgSettings("voicemail_settings", "Voicemail Settings", webexops.CallingCallingServiceGetVoicemail),
	orgSettings("voicemail_rules", "Voicemail Passcode Rules", webexops.CallingCallingServiceGetVoicemailRules),
	orgSettings("moh_settings", "Music on Hold", webexops.CallingCallingServiceGetOrgMusicHoldConfiguration),
	orgSettings("call_recording", "Call Recording", webexops.CallingCallRecordingGet),
	orgSettings("recording_announcement", "Recording Announcement", webexops.CallingCallRecordingGetOrgComplianceAnnouncement),
	orgSettings("emergency_notification", "Emergency Notification", webexops.CallingEmergencyServicesGetOrgCallNotification),
	orgSettings("call_captions", "Call Captions", webexops.CallingCallingServiceGetOrgCallCaptionsSettings),
	orgSettings("caller_reputation", "Caller Reputation", webexops.CallingCallerReputationProviderGetSettings),
	orgSettings("virtual_extension_settings", "Virtual Extension Mode", webexops.CallingVirtualExtensionsGetSettings),
}
