package catalog

import (
	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// Contact Center types. Every CC PUT replaces the whole object, so each type writes a
// fresh read with the changes laid over it. Routes and required create fields were
// checked on a live tenant; a type without a create or delete operation is one Webex
// would not create or delete through this API.

var (
	accessAll         = []string{"ALL", "SPECIFIC"}
	accessAllOrNone   = []string{"ALL", "SPECIFIC", "NONE"}
	workTypeCodes     = []string{"IDLE_CODE", "WRAP_UP_CODE"}
	burnoutInclusion  = []string{"INCLUDED", "EXCLUDED", "NOT_APPLICABLE"}
	variableTypes     = []string{"String", "Boolean", "Integer", "Decimal", "DateTime"}
	blendingModes     = []string{"BLENDED", "BLENDED_REALTIME", "EXCLUSIVE"}
	channelTypes      = []string{"TELEPHONY", "CHAT", "EMAIL", "SOCIAL_CHANNEL"}
	directions        = []string{"INBOUND", "OUTBOUND"}
	queueRoutingTypes = []string{"TEAM_BASED", "SKILL_BASED", "AGENT_BASED"}
	routingTypes      = []string{"LONGEST_AVAILABLE_AGENT", "BEST_AVAILABLE_AGENT", "LINEAR", "CIRCULAR"}
	skillTypes        = []string{"PROFICIENCY", "BOOLEAN", "TEXT"}
	weekDays          = []string{"MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"}
	loginVoice        = []string{"BROWSER", "EXTENSION", "AGENT_DN"}
	dnValidation      = []string{"NONE", "PROVISIONED_VALUE"}
	parentTypes       = []string{"ORGANIZATION", "SITE"}
	profileTypes      = []string{"ADMINISTRATOR", "ADMINISTRATOR_ONLY", "SUPERVISOR", "PREMIUM_AGENT", "STANDARD_AGENT"}
	enabledDisabled   = []string{"ENABLED", "DISABLED"}
	inclusionType     = []string{"ALL", "SPECIFIC"}
	includedExcluded  = []string{"INCLUDED", "EXCLUDED"}
	permissionAccess  = []string{"EDIT", "VIEW", "ENABLED", "DISABLED", "NONE"}

	// The permission names Webex uses in user profiles.
	permissionNames = []string{"additional-supervisory", "address-book", "agent-personalized-greeting", "agent-wellbeing",
		"analyzer", "api-key", "asset", "audio-prompt", "audio-recording-transcript", "auto-csat", "barge-in", "basic-digital",
		"business-rules-engine", "call-monitoring", "call-recording-schedule", "change-agent-state", "channel", "custom-channel",
		"desktop-layout", "desktop-profile", "download", "entry-point", "flow", "flow-settings-override", "function",
		"generated-summary", "global-variable", "holiday-list", "idle-wrapup-code", "manage-agent-queues-assignment",
		"mid-call-monitor", "multimedia", "multimedia-profile", "outdial-ani", "override", "personalized-ai-routing", "play",
		"queue", "quick-reply", "real-time-transcription", "reskill-agent", "resource-collection", "restricted-monitor-only",
		"routing-strategy", "schedule-listening", "send-message", "sentiment-analysis", "sign-out-agent", "site",
		"skill-definition", "skill-profile", "social-channel", "suggested-response", "support-number", "team", "tenant-setting",
		"threshold-rule", "user", "user-profile", "view-blind-monitor-request", "web-chat-asset", "whisper-coach", "working-hour"}

	// The resource types a resource collection lists. Webex expects these lowercase
	// names; uppercase ones fail with a 500. It adds "channel" itself if it is left out.
	resourceTypes = []string{"address-book", "audio-prompt", "cad-variable", "channel", "desktop-layout", "desktop-profile", "entry-point",
		"flow", "function", "holiday-list", "idle-wrapup-code", "multimedia-profile", "outdial-ani", "override", "queue",
		"site", "skill-definition", "skill-profile", "sub-flow", "team", "working-hour"}
)

func ref(name, path, target string) resource.Field {
	return resource.Field{Name: name, Path: path, Ref: target}
}

// refs is an unordered list of IDs. List references stay literal IDs in config, so they
// may point at any type, including one that points back.
func refs(name, path, target string) resource.Field {
	return resource.Field{Name: name, Path: path, Kind: resource.StringList, Ref: target, Unordered: true}
}

func picks(name, path string, options []string) resource.Field {
	return resource.Field{Name: name, Path: path, Kind: resource.StringList, Options: options, Unordered: true}
}

func when(name, path, format string) resource.Field {
	return resource.Field{Name: name, Path: path, Format: format}
}

func required(f resource.Field) resource.Field {
	f.Required = true
	return f
}

func fixed(f resource.Field) resource.Field {
	f.Immutable = true
	return f
}

func unique(f resource.Field) resource.Field {
	f.Unique = true
	return f
}

func matching(f resource.Field, pattern, hint string) resource.Field {
	f.Pattern, f.PatternHint = pattern, hint
	return f
}

func needs(f resource.Field, other string) resource.Field {
	f.Needs = other
	return f
}

func scoped(f resource.Field, by string) resource.Field {
	f.OptionScope = by
	return f
}

func rows(name, path, key string, unordered bool, item ...resource.Field) resource.Field {
	return resource.Field{Name: name, Path: path, Kind: resource.ObjectList, MatchKey: key, Unordered: unordered, Item: item}
}

func basics() []resource.Field {
	return []resource.Field{text("name", "name"), text("description", "description"), flag("active", "active")}
}

// ccRoutes binds a CC type's operations. Item operations take the object's ID as
// {id}. A nil operation is an action Webex does not offer for the type.
func ccRoutes(list, create, read, update, del *webexops.Operation) resource.Routes {
	r := resource.Routes{List: on(list), Updated: "lastUpdatedTime"}
	if create != nil {
		r.Create = on(create)
	}
	if read != nil {
		r.Read = on(read, "", "id")
	}
	if update != nil {
		r.Update = on(update, "", "id")
	}
	if del != nil {
		r.Delete = on(del, "", "id")
	}
	return r
}

// managedCC defines a writable CC type. A type with an active field is created active
// unless the create says otherwise, because Webex requires the value.
func managedCC(name, title, plural string, routes resource.Routes, fields ...[]resource.Field) resource.Type {
	t := resource.Type{Name: name, Title: title, Plural: plural, Area: resource.ContactCenter,
		Managed: true, LabelField: "name", Write: resource.WriteFull, Routes: routes, Fields: joined(fields...)}
	t.NoCreate = !routes.Create.Defined()
	t.NoDelete = !routes.Delete.Defined()
	if _, ok := t.Field("active"); ok {
		t.CreateDefaults = map[string]any{"active": true}
	}
	return t
}

func withDefaults(t resource.Type, defaults map[string]any) resource.Type {
	if t.CreateDefaults == nil {
		t.CreateDefaults = map[string]any{}
	}
	for k, v := range defaults {
		t.CreateDefaults[k] = v
	}
	return t
}

// settingsCC defines an org-wide CC settings object: one per org, edited but never
// created or deleted.
func settingsCC(name, title string, list, read, update *webexops.Operation, fs ...resource.Field) resource.Type {
	t := managedCC(name, title, title, ccRoutes(list, nil, read, update, nil), fs)
	t.LabelField = ""
	t.CreateDefaults = nil
	return t
}

// Agents & Teams

var ccMultimediaProfile = withDefaults(managedCC("cc_multimedia_profile", "Multimedia Profile", "Multimedia Profiles",
	ccRoutes(webexops.CCMultimediaProfileList, webexops.CCMultimediaProfileCreate, webexops.CCMultimediaProfileGet, webexops.CCMultimediaProfileUpdate, webexops.CCMultimediaProfileDelete),
	append(basics(), choice("blending_mode", "blendingMode", blendingModes), flag("blending_enabled", "blendingModeEnabled")),
	fields(number("telephony", "telephony"), number("chat", "chat"), number("email", "email"), number("social", "social")),
	fields(number("manual_chat", "manuallyAssignable.chat"), number("manual_email", "manuallyAssignable.email"),
		number("manual_social", "manuallyAssignable.social")),
), map[string]any{
	"blendingMode": "BLENDED", "blendingModeEnabled": false, "telephony": 1, "chat": 0, "email": 0, "social": 0,
	"manuallyAssignable": map[string]any{"chat": 0, "email": 0, "social": 0, "telephony": 0},
})

var ccSite = managedCC("cc_site", "Site", "Sites",
	ccRoutes(webexops.CCSiteListV1, webexops.CCSiteCreate, webexops.CCSiteGet, webexops.CCSiteUpdate, webexops.CCSiteDelete),
	append(basics(), required(ref("multimedia_profile_id", "multimediaProfileId", "cc_multimedia_profile"))))

// Webex refuses an update that changes a skill's name.
var ccSkill = managedCC("cc_skill", "Skill", "Skills",
	ccRoutes(webexops.CCSkillListV1, webexops.CCSkillCreate, webexops.CCSkillGet, webexops.CCSkillUpdate, webexops.CCSkillDelete),
	append(append([]resource.Field{fixed(text("name", "name"))}, basics()[1:]...),
		required(fixed(choice("skill_type", "skillType", skillTypes))),
		number("service_level_threshold", "serviceLevelThreshold")))

var ccSkillProfile = func() resource.Type {
	t := managedCC("cc_skill_profile", "Skill Profile", "Skill Profiles",
		ccRoutes(webexops.CCSkillProfileList, webexops.CCSkillProfileCreate, webexops.CCSkillProfileGet, webexops.CCSkillProfileUpdate, webexops.CCSkillProfileDelete),
		fields(text("name", "name"), text("description", "description")),
		fields(rows("skills", "activeSkills", "skill_id", true,
			required(ref("skill_id", "skillId", "cc_skill")), number("proficiency", "proficiencyValue"),
			flag("value", "booleanValue"), text("text", "textValue"))))
	t.Routes.ListIsPartial = true
	return t
}()

// A new desktop layout starts as a copy of the org's global layout file, because Webex
// requires a layout file and the model does not cover its content.
var ccDesktopLayout = func() resource.Type {
	t := withDefaults(managedCC("cc_desktop_layout", "Desktop Layout", "Desktop Layouts",
		ccRoutes(webexops.CCDesktopLayoutList, webexops.CCDesktopLayoutCreate, webexops.CCDesktopLayoutGet, webexops.CCDesktopLayoutUpdate, webexops.CCDesktopLayoutDelete),
		fields(text("name", "name"), text("description", "description"), flag("status", "status"),
			fixed(flag("global", "global"))),
		fields(refs("team_ids", "teamIds", "cc_team"))),
		map[string]any{"global": false, "status": true, "validated": false, "defaultJsonModified": false,
			"editedBy": "Webex Admin Studio", "jsonFileName": "Desktop Layout.json", "teamIds": []any{}})
	t.CopyOnCreate = &resource.Copy{Match: map[string]any{"global": true}, Keys: []string{"jsonFileContent"}}
	return t
}()

var ccTeam = withDefaults(managedCC("cc_team", "Team", "Teams",
	ccRoutes(webexops.CCTeamListV1, webexops.CCTeamCreate, webexops.CCTeamGet, webexops.CCTeamUpdate, webexops.CCTeamDelete),
	append(basics(),
		required(ref("site_id", "siteId", "cc_site")),
		choice("team_type", "teamType", teamTypes), choice("team_status", "teamStatus", teamStatuses),
		number("capacity", "capacity"), unique(text("dialed_number", "dialedNumber"))),
	fields(
		flag("rank_queues", "rankQueuesForTeam"),
		ref("desktop_layout_id", "desktopLayoutId", "cc_desktop_layout"),
		ref("multimedia_profile_id", "multiMediaProfileId", "cc_multimedia_profile"),
		ref("skill_profile_id", "skillProfileId", "cc_skill_profile")),
), map[string]any{"rankQueuesForTeam": false, "teamType": "AGENT", "teamStatus": "IN_SERVICE"})

// Each user profile type requires its own permissions, and Webex refuses a profile
// whose permissions duplicate another's. New permissions carry no ID: copying one from
// another profile is refused as a duplicate.
var ccUserProfile = func() resource.Type {
	t := withDefaults(managedCC("cc_user_profile", "User Profile", "User Profiles",
		ccRoutes(webexops.CCUserProfilesList, webexops.CCUserProfilesCreate, webexops.CCUserProfilesGet, webexops.CCUserProfilesUpdate, webexops.CCUserProfilesDelete),
		append(basics(), required(fixed(choice("profile_type", "profileType", profileTypes))),
			choice("permission_access", "permissionAccessLevel", accessAllOrNone),
			choice("resource_access", "resourceAccessLevel", accessAllOrNone)),
		fields(rows("permissions", "permissions", "name", true,
			required(choice("name", "name", permissionNames)), required(choice("access", "access", permissionAccess))))),
		map[string]any{"permissionAccessLevel": "SPECIFIC", "resourceAccessLevel": "ALL",
			"permissions": []any{map[string]any{"name": "basic-digital", "access": "ENABLED"}, map[string]any{"name": "multimedia", "access": "EDIT"}}})
	t.Routes.ListIsPartial = true
	return t
}()

// Webex requires a resource collection to list every resource type, and refuses NONE
// for sites, teams, queues and entry points. The IDs a SPECIFIC row names are kept on
// writes but not edited here.
var ccResourceCollection = func() resource.Type {
	defaults := []any{}
	for _, name := range resourceTypes {
		defaults = append(defaults, map[string]any{"name": name, "accessLevel": "ALL", "ids": []any{}})
	}
	t := withDefaults(managedCC("cc_resource_collection", "Resource Collection", "Resource Collections",
		ccRoutes(webexops.CCResourceCollectionList, webexops.CCResourceCollectionCreate, webexops.CCResourceCollectionGet, webexops.CCResourceCollectionUpdate, webexops.CCResourceCollectionDelete),
		fields(text("name", "name"), text("description", "description")),
		fields(rows("resources", "resources", "resource", true,
			required(choice("resource", "name", resourceTypes)), required(choice("access", "accessLevel", accessAllOrNone))))),
		map[string]any{"resources": defaults})
	t.Routes.ListIsPartial = true
	return t
}()

// Business Hours

var ccHolidayList = func() resource.Type {
	t := managedCC("cc_holiday_list", "Holiday List", "Holiday Lists",
		ccRoutes(webexops.CCHolidayListList, webexops.CCHolidayListCreate, webexops.CCHolidayListGet, webexops.CCHolidayListUpdate, webexops.CCHolidayListDelete),
		fields(text("name", "name"), text("description", "description")),
		fields(rows("holidays", "holidays", "name", false,
			required(text("name", "name")), required(when("start_date", "startDate", resource.Date)), required(when("end_date", "endDate", resource.Date)))))
	t.Routes.ListIsPartial = true
	return t
}()

var ccOverride = func() resource.Type {
	t := managedCC("cc_override", "Override", "Overrides",
		ccRoutes(webexops.CCOverridesList, webexops.CCOverridesCreate, webexops.CCOverridesGet, webexops.CCOverridesUpdate, webexops.CCOverridesDelete),
		fields(text("name", "name"), text("description", "description"),
			required(sourced("timezone", "timezone", resource.TimeZones))),
		fields(rows("overrides", "overrides", "name", false,
			required(text("name", "name")), flag("working_hours", "workingHours"),
			required(when("start", "startDateTime", resource.DateTime)), required(when("end", "endDateTime", resource.DateTime)))))
	t.Routes.ListIsPartial = true
	return t
}()

var ccBusinessHours = func() resource.Type {
	t := managedCC("cc_business_hours", "Business Hours", "Business Hours",
		ccRoutes(webexops.CCBusinessHourListV1, webexops.CCBusinessHourCreate, webexops.CCBusinessHourGet, webexops.CCBusinessHourUpdate, webexops.CCBusinessHourDelete),
		fields(text("name", "name"), text("description", "description"),
			required(sourced("timezone", "timezone", resource.TimeZones)),
			ref("holidays_id", "holidaysId", "cc_holiday_list"), ref("overrides_id", "overridesId", "cc_override")),
		fields(rows("working_hours", "workingHours", "name", true,
			required(text("name", "name")), picks("days", "days", weekDays),
			required(when("start_time", "startTime", resource.Time)), required(when("end_time", "endTime", resource.Time)))))
	t.Routes.ListIsPartial = true
	return t
}()

// Routing

var ccQueue = func() resource.Type {
	t := withDefaults(managedCC("cc_queue", "Queue", "Queues",
		ccRoutes(webexops.CCContactServiceQueueListV2, webexops.CCContactServiceQueueCreate, webexops.CCContactServiceQueueGet, webexops.CCContactServiceQueueUpdate, webexops.CCContactServiceQueueDelete),
		append(basics(),
			fixed(choice("channel_type", "channelType", channelTypes)), fixed(choice("queue_type", "queueType", directions)),
			fixed(choice("queue_routing_type", "queueRoutingType", queueRoutingTypes)),
			choice("routing_type", "routingType", routingTypes),
			required(sourced("timezone", "timezone", resource.TimeZones)),
			required(ref("default_music_id", "defaultMusicInQueueMediaFileId", "cc_audio_file"))),
		fields(
			number("max_active_contacts", "maxActiveContacts"), number("max_time_in_queue", "maxTimeInQueue"),
			number("service_level_threshold", "serviceLevelThreshold"),
			flag("check_agent_availability", "checkAgentAvailability"), flag("manually_assignable", "manuallyAssignable")),
		fields(rows("distribution_groups", "callDistributionGroups", "order", false,
			required(number("order", "order")), number("duration", "duration"),
			resource.Field{Name: "team_ids", Path: "agentGroups[].teamId", Kind: resource.StringList, Ref: "cc_team"})),
		fields(
			flag("monitoring_permitted", "monitoringPermitted"), flag("parking_permitted", "parkingPermitted"),
			flag("recording_permitted", "recordingPermitted"), flag("record_all_calls", "recordingAllCallsPermitted"),
			flag("pause_recording", "pauseRecordingPermitted"), number("recording_pause_seconds", "recordingPauseDuration")),
		fields(
			choice("auto_csat", "queueLevelAutoCsatInclusion", includedExcluded),
			choice("realtime_transcriptions", "queueLevelRealTimeTranscriptionsInclusion", includedExcluded),
			choice("sentiment_analysis", "queueLevelSentimentAnalysisInclusion", includedExcluded),
			choice("generated_summaries", "queueLevelSummariesInclusion", includedExcluded)),
	), map[string]any{
		"queueType": "INBOUND", "channelType": "TELEPHONY", "queueRoutingType": "TEAM_BASED", "routingType": "LONGEST_AVAILABLE_AGENT",
		"checkAgentAvailability": false, "maxActiveContacts": 0, "maxTimeInQueue": 600, "serviceLevelThreshold": 30,
		"monitoringPermitted": false, "parkingPermitted": false, "recordingPermitted": false,
		"recordingAllCallsPermitted": false, "pauseRecordingPermitted": false, "callDistributionGroups": []any{},
	})
	t.Routes.ListIsPartial = true
	return t
}()

// ccFlow lists Flow Designer flows, which entry points run. The project ID is the same
// fixed value for every org, and the list pages by size, not pageSize.
var ccFlow = center("cc_flow", "Flow", "Flows",
	resource.Route{Op: webexops.CCFlowList, Query: map[string]string{"flowType": "FLOW", "includePagination": "true", "size": "100"}}, "name",
	text("name", "name"), text("description", "description"), text("status", "status"))

// An entry point runs one version of a flow, named by one of the flow's tags, such as
// Live or Latest. Webex refuses a tag that points at no published version.
var ccEntryPoint = withDefaults(managedCC("cc_entry_point", "Entry Point", "Entry Points",
	ccRoutes(webexops.CCEntryPointListV1, webexops.CCEntryPointCreate, webexops.CCEntryPointGet, webexops.CCEntryPointUpdate, webexops.CCEntryPointDelete),
	append(basics(),
		fixed(choice("channel_type", "channelType", channelTypes)), fixed(choice("entry_point_type", "entryPointType", directions)),
		required(sourced("timezone", "timezone", resource.TimeZones)),
		number("maximum_active_contacts", "maximumActiveContacts"), number("service_level_threshold", "serviceLevelThreshold"),
		ref("music_on_hold_id", "musicOnHoldId", "cc_audio_file")),
	fields(needs(ref("flow_id", "flowId", "cc_flow"), "music_on_hold_id"), scoped(sourced("flow_version", "flowTagId", resource.FlowTags), "flow_id")),
	fields(
		flag("callback_enabled", "callbackEnabled"), flag("outdial_transfer_to_queue", "outdialTransferToQueueEnabled"),
		ref("outdial_queue_id", "outdialQueueId", "cc_queue")),
), map[string]any{"entryPointType": "INBOUND", "channelType": "TELEPHONY", "serviceLevelThreshold": 30, "maximumActiveContacts": 0})

var ccDialPlan = managedCC("cc_dial_plan", "Dial Plan", "Dial Plans",
	ccRoutes(webexops.CCDialPlanList, webexops.CCDialPlanCreate, webexops.CCDialPlanGet, webexops.CCDialPlanUpdate, webexops.CCDialPlanDelete),
	append(basics(), required(text("regular_expression", "regularExpression")),
		text("prefix", "prefix"), text("stripped_chars", "strippedChars")))

var ccOutdialANI = func() resource.Type {
	t := managedCC("cc_outdial_ani", "Outdial ANI", "Outdial ANIs",
		ccRoutes(webexops.CCOutdialAniList, webexops.CCOutdialAniCreate, webexops.CCOutdialAniGet, webexops.CCOutdialAniUpdate, webexops.CCOutdialAniDelete),
		fields(text("name", "name"), text("description", "description")),
		fields(rows("numbers", "outdialANIEntries", "number", true,
			required(text("name", "name")), required(text("number", "number")), flag("default", "defaultANIEntry"))))
	t.Routes.ListIsPartial = true
	return t
}()

var ccContactNumber = func() resource.Type {
	t := managedCC("cc_contact_number", "Contact Number", "Contact Numbers",
		ccRoutes(webexops.CCContactNumberList, webexops.CCContactNumberCreate, webexops.CCContactNumberGet, webexops.CCContactNumberUpdate, webexops.CCContactNumberDelete),
		fields(required(text("number", "number"))))
	t.LabelField = "number"
	return t
}()

var ccAddressBook = withDefaults(managedCC("cc_address_book", "Address Book", "Address Books",
	ccRoutes(webexops.CCAddressBookList, webexops.CCAddressBookCreate, webexops.CCAddressBookGet, webexops.CCAddressBookUpdate, webexops.CCAddressBookDelete),
	fields(text("name", "name"), text("description", "description"),
		fixed(choice("parent_type", "parentType", parentTypes)), ref("site_id", "siteId", "cc_site"))),
	map[string]any{"parentType": "ORGANIZATION"})

// ccDialNumber maps one number to one entry point, so an entry point's phone numbers
// are the dial numbers that name it. It cannot be created here: Webex needs a number
// the org already owns and its CC location, which come from Control Hub's number
// assignment.
var ccDialNumber = func() resource.Type {
	t := managedCC("cc_dial_number", "Dial Number", "Dial Numbers",
		ccRoutes(webexops.CCDialNumberListDialedMapping, nil, webexops.CCDialNumberGetDialedMappingId, webexops.CCDialNumberUpdateDialedMappingId, webexops.CCDialNumberDeleteDialedMappingId),
		fields(fixed(text("dialled_number", "dialledNumber")), unique(text("extension", "extension")),
			ref("entry_point_id", "entryPointId", "cc_entry_point"), flag("default_ani", "defaultAni")))
	t.LabelField = "dialled_number"
	t.FallbackLabelField = "extension"
	return t
}()

// Desktop & Media

// Work type item reads return 406, so a work type is read through the list with an
// id filter.
var ccWorkType = func() resource.Type {
	t := managedCC("cc_work_type", "Work Type", "Work Types",
		ccRoutes(webexops.CCWorkTypesList, webexops.CCWorkTypesCreate, nil, webexops.CCWorkTypesUpdate, webexops.CCWorkTypesDelete),
		append(basics(), required(fixed(choice("work_type_code", "workTypeCode", workTypeCodes)))))
	t.Routes.ReadByFilter = true
	return t
}()

var ccAuxiliaryCode = withDefaults(managedCC("cc_auxiliary_code", "Idle or Wrap-Up Code", "Idle & Wrap-Up Codes",
	ccRoutes(webexops.CCAuxiliaryCodeList, webexops.CCAuxiliaryCodeCreate, webexops.CCAuxiliaryCodeGet, webexops.CCAuxiliaryCodeUpdate, webexops.CCAuxiliaryCodeDelete),
	append(basics(),
		required(fixed(choice("work_type_code", "workTypeCode", workTypeCodes))),
		required(ref("work_type_id", "workTypeId", "cc_work_type")),
		flag("default_code", "defaultCode"), choice("burnout_inclusion", "burnoutInclusion", burnoutInclusion)),
), map[string]any{"defaultCode": false, "isSystemCode": false})

// Webex refuses a global variable name longer than 80 characters or with any
// character other than an ASCII letter, digit or underscore.
var ccGlobalVariable = withDefaults(managedCC("cc_global_variable", "Global Variable", "Global Variables",
	ccRoutes(webexops.CCGlobalVariablesList, webexops.CCGlobalVariablesCreate, webexops.CCGlobalVariablesGet, webexops.CCGlobalVariablesUpdate, webexops.CCGlobalVariablesDelete),
	append(append([]resource.Field{matching(text("name", "name"), `[A-Za-z0-9_]{1,80}`,
		"Use only letters A–Z, digits and _, up to 80 characters")},
		basics()[1:]...),
		required(fixed(choice("variable_type", "variableType", variableTypes))), text("default_value", "defaultValue"),
		text("desktop_label", "desktopLabel")),
	fields(flag("agent_viewable", "agentViewable"), flag("agent_editable", "agentEditable"),
		flag("reportable", "reportable"), flag("sensitive", "sensitive")),
), map[string]any{"agentViewable": false, "agentEditable": false, "reportable": false})

// ccAudioFile is edited with PATCH: its PUT and POST take a multipart file upload, so
// files are added in Control Hub and only the description is edited here.
var ccAudioFile = func() resource.Type {
	t := managedCC("cc_audio_file", "Audio File", "Audio Files",
		ccRoutes(webexops.CCAudioFilesList, nil, webexops.CCAudioFilesGet, webexops.CCAudioFilesPatch, webexops.CCAudioFilesDelete),
		fields(fixed(text("name", "name")), text("description", "description")))
	t.Write = resource.WritePatch
	return t
}()

// Desktop profiles come after the types their single references point at, so config
// can refer to those objects by address.
var ccDesktopProfile = func() resource.Type {
	t := withDefaults(managedCC("cc_agent_profile", "Desktop Profile", "Desktop Profiles",
		ccRoutes(webexops.CCDesktopProfileList, webexops.CCDesktopProfileCreate, webexops.CCDesktopProfileGet, webexops.CCDesktopProfileUpdate, webexops.CCDesktopProfileDelete),
		basics(),
		fields(
			choice("idle_code_access", "accessIdleCode", accessAll), refs("idle_codes", "idleCodes", "cc_auxiliary_code"),
			choice("wrap_up_code_access", "accessWrapUpCode", accessAll), refs("wrap_up_codes", "wrapUpCodes", "cc_auxiliary_code"),
			choice("queue_access", "accessQueue", accessAllOrNone), refs("queues", "queues", "cc_queue"),
			choice("entry_point_access", "accessEntryPoint", accessAllOrNone), refs("entry_points", "entryPoints", "cc_entry_point"),
			choice("buddy_team_access", "accessBuddyTeam", accessAllOrNone), refs("buddy_teams", "buddyTeams", "cc_team")),
		fields(
			picks("login_voice_options", "loginVoiceOptions", loginVoice),
			flag("auto_answer", "autoAnswer"), flag("auto_wrap_up", "autoWrapUp"),
			number("auto_wrap_after_seconds", "autoWrapAfterSeconds"), flag("allow_auto_wrap_up_extension", "allowAutoWrapUpExtension"),
			flag("screen_popup", "screenPopup"), flag("last_agent_routing", "lastAgentRouting"),
			flag("consult_to_queue", "consultToQueue"), flag("manage_channel_availability", "manageChannelAvailability"),
			flag("auto_accept_digital", "autoAcceptDigitalInteractions"), flag("personal_greeting", "agentPersonalGreeting"),
			flag("schedule_callbacks", "scheduleAndManageCallBack"),
			flag("custom_inactivity_timeout", "timeoutDesktopInactivityCustomEnabled"),
			number("inactivity_timeout_minutes", "timeoutDesktopInactivityMins")),
		fields(
			flag("outdial_enabled", "outdialEnabled"), ref("outdial_entry_point_id", "outdialEntryPointId", "cc_entry_point"),
			ref("outdial_ani_id", "outdialANIId", "cc_outdial_ani"), ref("address_book_id", "addressBookId", "cc_address_book"),
			flag("agent_available_after_outdial", "agentAvailableAfterOutdial"),
			flag("dial_plan_enabled", "dialPlanEnabled"), refs("dial_plans", "dialPlans", "cc_dial_plan"),
			choice("agent_dn_validation", "agentDNValidation", dnValidation)),
		fields(
			flag("agent_statistics", "viewableStatistics.agentStats"),
			choice("queue_statistics", "viewableStatistics.accessQueueStats", accessAllOrNone),
			choice("team_statistics", "viewableStatistics.accessTeamStats", accessAllOrNone),
			flag("logged_in_team_statistics", "viewableStatistics.loggedInTeamStats")),
		fields(
			flag("webex_user_details", "showUserDetailsWebex"), flag("webex_state_sync", "stateSynchronizationWebex"),
			flag("teams_user_details", "showUserDetailsMS"), flag("teams_state_sync", "stateSynchronizationMS")),
	), map[string]any{
		"parentType": "ORGANIZATION", "accessIdleCode": "ALL", "accessWrapUpCode": "ALL", "accessQueue": "ALL",
		"accessEntryPoint": "ALL", "accessBuddyTeam": "ALL", "agentDNValidation": "NONE",
		"loginVoiceOptions": []any{"BROWSER", "EXTENSION", "AGENT_DN"}, "dialPlanEnabled": false, "outdialEnabled": false,
		"viewableStatistics": map[string]any{"agentStats": false, "accessQueueStats": "NONE", "contactServiceQueues": []any{},
			"loggedInTeamStats": false, "accessTeamStats": "NONE", "teams": []any{}},
	})
	t.Routes.ListIsPartial = true
	return t
}()

// ccUser edits a person's Contact Center settings. People are added and removed in
// Control Hub, so users are never created or deleted here.
var ccUser = func() resource.Type {
	t := managedCC("cc_user", "User", "Users",
		ccRoutes(webexops.CCUsersListV1, nil, webexops.CCUsersGet, webexops.CCUsersUpdate, nil),
		fields(fixed(text("email", "email")), fixed(text("first_name", "firstName")), fixed(text("last_name", "lastName")),
			unique(text("external_identifier", "externalIdentifier")), unique(text("default_dialled_number", "deafultDialledNumber"))),
		fields(
			ref("site_id", "siteId", "cc_site"), refs("team_ids", "teamIds", "cc_team"),
			ref("preferred_supervisor_team_id", "preferredSupervisorTeamId", "cc_team"),
			ref("desktop_profile_id", "agentProfileId", "cc_agent_profile"), ref("user_profile_id", "userProfileId", "cc_user_profile"),
			ref("multimedia_profile_id", "multimediaProfileId", "cc_multimedia_profile"),
			ref("skill_profile_id", "skillProfileId", "cc_skill_profile")),
		fields(
			choice("auto_csat", "userLevelAutoCSATInclusion", includedExcluded),
			choice("agent_wellbeing", "userLevelBurnoutInclusion", includedExcluded),
			choice("generated_summaries", "userLevelSummariesInclusion", includedExcluded),
			choice("wellness_break_reminders", "userLevelWellnessBreakReminders", enabledDisabled)))
	t.LabelField = "email"
	t.Routes.ListIsPartial = true
	return t
}()

// AI Features. The combined ai-feature object is left read-only: its route is a PATCH
// whose handling of nested settings is unverified, and the three objects below edit
// the same settings through full PUTs.

var ccAgentWellbeing = settingsCC("cc_agent_wellbeing", "Agent Wellbeing",
	webexops.CCAgentWellbeingListBurnout, webexops.CCAgentWellbeingGetBurnoutId, webexops.CCAgentWellbeingUpdateBurnoutId,
	flag("enabled", "enabled"), choice("agent_inclusion", "agentInclusionType", inclusionType),
	choice("wellness_break_reminders", "wellnessBreakReminders", enabledDisabled))

var ccAutoCSAT = settingsCC("cc_auto_csat", "Auto CSAT",
	webexops.CCAutoCsatList, webexops.CCAutoCsatGet, webexops.CCAutoCsatUpdate,
	flag("enabled", "enabled"), choice("agent_inclusion", "agentInclusionType", inclusionType),
	text("survey_data_source", "surveyDataSource"),
	ref("global_variable_id", "selectedGlobalVariableId", "cc_global_variable"))

var ccGeneratedSummaries = settingsCC("cc_generated_summaries", "Generated Summaries",
	webexops.CCGeneratedSummariesList, webexops.CCGeneratedSummariesGet, webexops.CCGeneratedSummariesUpdate,
	flag("call_drop_summaries", "callDropSummariesEnabled"), flag("virtual_agent_transfer_summaries", "virtualAgentTransferSummariesEnabled"),
	flag("consult_transfer_summaries", "consultTransferSummariesEnabled"), choice("agent_inclusion", "agentInclusionType", inclusionType))

var ccAIFeature = center("cc_ai_feature", "AI Features", "AI Features", on(webexops.CCAiFeatureList), "")

// centerTypes lists the CC types, grouped by subject.
var centerTypes = []resource.Type{
	ccSite, ccUser, ccTeam, ccSkill, ccSkillProfile, ccDesktopProfile, ccMultimediaProfile, ccUserProfile, ccDesktopLayout, ccResourceCollection,
	ccEntryPoint, ccFlow, ccQueue, ccDialNumber, ccDialPlan, ccOutdialANI, ccContactNumber, ccAddressBook,
	ccBusinessHours, ccHolidayList, ccOverride,
	ccAuxiliaryCode, ccWorkType, ccGlobalVariable, ccAudioFile,
	ccAIFeature, ccAgentWellbeing, ccAutoCSAT, ccGeneratedSummaries,
}
