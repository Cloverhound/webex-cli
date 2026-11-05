package catalog

import "github.com/Cloverhound/webex-cli/pkg/resource"

var (
	huntPolicies    = []string{"CIRCULAR", "REGULAR", "SIMULTANEOUS", "UNIFORM", "WEIGHTED"}
	queuePolicies   = []string{"CIRCULAR", "REGULAR", "SIMULTANEOUS", "UNIFORM", "WEIGHTED"}
	queueRouting    = []string{"PRIORITY_BASED", "SKILL_BASED"}
	overflowActions = []string{"PERFORM_BUSY_TREATMENT", "PLAY_RINGING_UNTIL_CALLER_HANGS_UP", "TRANSFER_TO_PHONE_NUMBER"}
	callerIDNames   = []string{"DISPLAY_NAME", "CUSTOM_NAME"}
	lineIDPolicies  = []string{"DIRECT_LINE", "LOCATION_NUMBER", "CUSTOM"}
	dialingScopes   = []string{"ENTERPRISE", "GROUP"}
	teamTypes       = []string{"AGENT", "CAPACITY"}
	teamStatuses    = []string{"IN_SERVICE", "NOT_AVAILABLE"}
)

func number(name, path string) resource.Field {
	return resource.Field{Name: name, Path: path, Kind: resource.Int}
}

func choice(name, path string, options []string) resource.Field {
	return resource.Field{Name: name, Path: path, Options: options}
}

func sourced(name, path, source string) resource.Field {
	return resource.Field{Name: name, Path: path, OptionSource: source}
}

func fields(fs ...resource.Field) []resource.Field { return fs }

func joined(groups ...[]resource.Field) []resource.Field {
	var out []resource.Field
	for _, g := range groups {
		out = append(out, g...)
	}
	return out
}

var huntGroupFields = joined(
	fields(
		locationField, text("name", "name"), text("first_name", "firstName"), text("last_name", "lastName"),
		unique(text("extension", "extension")), unique(text("phone_number", "phoneNumber")), flag("enabled", "enabled"),
		sourced("language_code", "languageCode", resource.Languages), sourced("time_zone", "timeZone", resource.TimeZones),
		text("dial_by_name", "dialByName")),
	fields(
		choice("caller_id_name", "directLineCallerIdName.selection", callerIDNames),
		text("caller_id_custom_name", "directLineCallerIdName.customName"),
		flag("use_for_outgoing_calls", "huntGroupCallerIdForOutgoingCallsEnabled"),
		flag("distinctive_ring", "alternateNumberSettings.distinctiveRingEnabled")),
	fields(
		choice("policy", "callPolicies.policy", huntPolicies),
		flag("waiting_enabled", "callPolicies.waitingEnabled"),
		flag("group_busy", "callPolicies.groupBusyEnabled"),
		flag("members_control_group_busy", "callPolicies.allowMembersToControlGroupBusyEnabled")),
	fields(
		number("no_answer_rings", "callPolicies.noAnswer.numberOfRings"),
		flag("next_agent_enabled", "callPolicies.noAnswer.nextAgentEnabled"),
		number("next_agent_rings", "callPolicies.noAnswer.nextAgentRings"),
		flag("no_answer_forward", "callPolicies.noAnswer.forwardEnabled"),
		text("no_answer_destination", "callPolicies.noAnswer.destination")),
	fields(
		flag("busy_redirect", "callPolicies.busyRedirect.enabled"),
		text("busy_redirect_destination", "callPolicies.busyRedirect.destination"),
		flag("continuity_redirect", "callPolicies.businessContinuityRedirect.enabled"),
		text("continuity_redirect_destination", "callPolicies.businessContinuityRedirect.destination")),
)

var callQueueFields = joined(
	fields(
		locationField, text("name", "name"), text("first_name", "firstName"), text("last_name", "lastName"),
		unique(text("extension", "extension")), unique(text("phone_number", "phoneNumber")), flag("enabled", "enabled"),
		sourced("language_code", "languageCode", resource.Languages), sourced("time_zone", "timeZone", resource.TimeZones),
		text("dial_by_name", "dialByName")),
	fields(
		choice("calling_line_id_policy", "callingLineIdPolicy", lineIDPolicies),
		text("calling_line_id_number", "callingLineIdPhoneNumber"),
		choice("caller_id_name", "directLineCallerIdName.selection", callerIDNames),
		text("caller_id_custom_name", "directLineCallerIdName.customName"),
		flag("use_for_outgoing_calls", "phoneNumberForOutgoingCallsEnabled"),
		flag("distinctive_ring", "alternateNumberSettings.distinctiveRingEnabled")),
	fields(
		choice("policy", "callPolicies.policy", queuePolicies),
		choice("routing_type", "callPolicies.routingType", queueRouting),
		flag("agents_can_join", "allowAgentJoinEnabled"),
		flag("call_waiting_for_agents", "allowCallWaitingForAgentsEnabled")),
	fields(
		flag("call_bounce", "callPolicies.callBounce.callBounceEnabled"),
		number("call_bounce_rings", "callPolicies.callBounce.callBounceMaxRings"),
		flag("bounce_if_unavailable", "callPolicies.callBounce.agentUnavailableEnabled"),
		flag("alert_agent", "callPolicies.callBounce.alertAgentEnabled"),
		number("alert_agent_seconds", "callPolicies.callBounce.alertAgentMaxSeconds"),
		flag("bounce_on_hold", "callPolicies.callBounce.callBounceOnHoldEnabled"),
		number("bounce_on_hold_seconds", "callPolicies.callBounce.callBounceOnHoldMaxSeconds")),
	fields(
		number("queue_size", "queueSettings.queueSize"),
		flag("call_offer_tone", "queueSettings.callOfferToneEnabled"),
		flag("reset_statistics", "queueSettings.resetCallStatisticsEnabled")),
	fields(
		choice("overflow_action", "queueSettings.overflow.action", overflowActions),
		flag("overflow_after_wait", "queueSettings.overflow.overflowAfterWaitEnabled"),
		number("overflow_wait_seconds", "queueSettings.overflow.overflowAfterWaitTime"),
		flag("overflow_greeting", "queueSettings.overflow.playOverflowGreetingEnabled")),
	fields(
		flag("welcome_message", "queueSettings.welcomeMessage.enabled"),
		flag("welcome_always", "queueSettings.welcomeMessage.alwaysEnabled"),
		flag("wait_message", "queueSettings.waitMessage.enabled"),
		flag("comfort_message", "queueSettings.comfortMessage.enabled"),
		number("comfort_message_seconds", "queueSettings.comfortMessage.timeBetweenMessages"),
		flag("whisper_message", "queueSettings.whisperMessage.enabled")),
)

var autoAttendantFields = joined(
	fields(
		locationField, text("name", "name"), text("first_name", "firstName"), text("last_name", "lastName"),
		unique(text("extension", "extension")), unique(text("phone_number", "phoneNumber")),
		sourced("language_code", "languageCode", resource.Languages), sourced("time_zone", "timeZone", resource.TimeZones),
		text("dial_by_name", "dialByName"),
		sourced("business_schedule", "businessSchedule", resource.Schedules)),
	fields(
		choice("extension_dialing", "extensionDialing", dialingScopes),
		choice("name_dialing", "nameDialing", dialingScopes)),
)
