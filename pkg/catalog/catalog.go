// Package catalog lists the Webex object types the SDK can read, and the ones it can
// write. Each type names the operations from pkg/webexops it uses, so a collection
// refresh that drops one fails to build.
//
// Writes on a Managed type were proven on a live tenant before it was marked so.
// Other types are read-only.
package catalog

import (
	"fmt"

	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
)

// on binds an operation to a type. keys name the operation's path parameters that
// take the object's parent and ID, in that order.
func on(op *webexops.Operation, keys ...string) resource.Route {
	r := resource.Route{Op: op}
	if len(keys) > 0 {
		r.Parent = keys[0]
	}
	if len(keys) > 1 {
		r.ID = keys[1]
	}
	return r
}

var huntGroup = resource.Type{
	Name:        "hunt_group",
	Title:       "Hunt Group",
	Plural:      "Hunt Groups",
	Area:        resource.Calling,
	Managed:     true,
	ParentField: "location_id",
	LabelField:  "name",
	Write:       resource.WriteDelta,
	Routes: resource.Routes{
		List:          on(webexops.CallingHuntGroupList),
		Create:        on(webexops.CallingHuntGroupCreate, "locationId"),
		Read:          on(webexops.CallingHuntGroupGet, "locationId", "huntGroupId"),
		Update:        on(webexops.CallingHuntGroupUpdate, "locationId", "huntGroupId"),
		Delete:        on(webexops.CallingHuntGroupDelete, "locationId", "huntGroupId"),
		ArrayKey:      "huntGroups",
		ParentKey:     "locationId",
		ListIsPartial: true,
	},
	Fields: huntGroupFields,
	CreateDefaults: map[string]any{
		"agents": []any{},
		"callPolicies": map[string]any{
			"waitingEnabled":                        false,
			"groupBusyEnabled":                      false,
			"allowMembersToControlGroupBusyEnabled": false,
			"noAnswer": map[string]any{
				"nextAgentEnabled":            false,
				"nextAgentRings":              5,
				"forwardEnabled":              false,
				"destinationVoicemailEnabled": false,
			},
			"busyRedirect":               map[string]any{"enabled": false, "destinationVoicemailEnabled": false},
			"businessContinuityRedirect": map[string]any{"enabled": false, "destinationVoicemailEnabled": false},
		},
	},
}

var callQueue = resource.Type{
	Name:        "call_queue",
	Title:       "Call Queue",
	Plural:      "Call Queues",
	Area:        resource.Calling,
	Managed:     true,
	ParentField: "location_id",
	LabelField:  "name",
	Write:       resource.WriteDelta,
	Routes: resource.Routes{
		List:          on(webexops.CallingCallQueueListCustomerAssist),
		Create:        on(webexops.CallingCallQueueCreateCustomerAssist, "locationId"),
		Read:          on(webexops.CallingCallQueueGetCustomerAssist, "locationId", "queueId"),
		Update:        on(webexops.CallingCallQueueUpdate, "locationId", "queueId"),
		Delete:        on(webexops.CallingCallQueueDelete, "locationId", "queueId"),
		ArrayKey:      "queues",
		ParentKey:     "locationId",
		ListIsPartial: true,
	},
	Fields: callQueueFields,
	CreateDefaults: map[string]any{
		"agents": []any{},
		"callPolicies": map[string]any{
			"callBounce":      map[string]any{"callBounceEnabled": false},
			"distinctiveRing": map[string]any{"enabled": false},
		},
		"queueSettings": map[string]any{},
	},
}

var autoAttendant = resource.Type{
	Name:        "auto_attendant",
	Title:       "Auto Attendant",
	Plural:      "Auto Attendants",
	Area:        resource.Calling,
	Managed:     true,
	NoCreate:    true,
	NoDelete:    true,
	ParentField: "location_id",
	LabelField:  "name",
	Write:       resource.WriteDelta,
	Routes: resource.Routes{
		List:          on(webexops.CallingAutoAttendantList),
		Read:          on(webexops.CallingAutoAttendantGet, "locationId", "autoAttendantId"),
		Update:        on(webexops.CallingAutoAttendantUpdate, "locationId", "autoAttendantId"),
		ArrayKey:      "autoAttendants",
		ParentKey:     "locationId",
		ListIsPartial: true,
	},
	Fields: autoAttendantFields,
}

var location = resource.Type{
	Name:       "location",
	Title:      "Location",
	Plural:     "Locations",
	Area:       resource.Calling,
	LabelField: "name",
	Routes: resource.Routes{
		List:     on(webexops.CallingLocationsList),
		Read:     on(webexops.CallingLocationsGet, "", "locationId"),
		ArrayKey: "items",
	},
	Fields: []resource.Field{
		{Name: "name", Path: "name"},
		{Name: "time_zone", Path: "timeZone"},
		{Name: "city", Path: "address.city"},
	},
}

// All lists every type in a stable order: Calling first, then Contact Center, each
// grouped by subject. Webex Admin Studio's sidebar follows it.
func All() []resource.Type {
	var out []resource.Type
	for _, group := range [][]resource.Type{
		{location}, callingOrganization,
		{huntGroup, callQueue, autoAttendant}, callingRouting,
		callingSchedules, callingPSTN, callingDevices, callingLocationSettings, callingOrgSettings,
		centerTypes,
	} {
		out = append(out, group...)
	}
	return out
}

// Managed lists the types that can be written.
func Managed() []resource.Type {
	var out []resource.Type
	for _, t := range All() {
		if t.Managed {
			out = append(out, t)
		}
	}
	return out
}

// ByName finds a type.
func ByName(name string) (resource.Type, error) {
	for _, t := range All() {
		if t.Name == name {
			return t, nil
		}
	}
	return resource.Type{}, fmt.Errorf("unknown object type %q", name)
}
