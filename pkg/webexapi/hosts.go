package webexapi

import (
	"encoding/base64"
	"fmt"
	"slices"
	"strings"
)

// Hosts are the base URLs of the Webex APIs for one organization's data region
// and Contact Center data center.
type Hosts struct {
	// Webex serves most APIs, under /v1.
	Webex string

	// Identity serves SCIM and archived users, whose paths have no /v1 segment.
	Identity string

	// Analytics serves the Admin analytics reports, whose paths carry their own
	// /v1 segment.
	Analytics string

	// AnalyticsCalling serves Detailed Call History in the data region.
	AnalyticsCalling string

	// CC is Contact Center in its data center, and AIAgent is AI Agent Studio in
	// the same data center.
	CC, AIAgent string
}

// ccRegions are the Contact Center data centers HostsFor accepts; a list callers
// cannot change.
var ccRegions = []string{"us1", "eu1", "eu2", "anz1", "ca1", "jp1", "sg1", "in1"}

// Regions are the data regions, and CCRegions the Contact Center data centers, for
// display. HostsFor checks against its own copy.
var (
	Regions   = []string{"us", "ca", "eu", "eun", "in", "au"}
	CCRegions = slices.Clone(ccRegions)
)

// DefaultCCRegion is the Contact Center data center used when none is given.
const DefaultCCRegion = "us1"

// analyticsCallingHosts maps a data region to its Detailed Call History host.
// The default host routes a call to the nearest region; if that region does not
// hold the organization's data, the API answers 451 naming the right host.
var analyticsCallingHosts = map[string]string{
	"us":  "https://analytics-calling.webexapis.com",
	"ca":  "https://analytics-calling.webexapis.com",
	"eu":  "https://analytics-calling-eu.webexapis.com",
	"eun": "https://analytics-calling-eu.webexapis.com",
	"in":  "https://analytics-calling-in.webexapis.com",
	"au":  "https://analytics-calling-au.webexapis.com",
}

// HostsFor returns the hosts for a data region from Regions and a Contact Center
// data center from CCRegions. An empty or unknown region uses the US hosts, and an
// empty ccRegion DefaultCCRegion. A ccRegion not in CCRegions is an error, since it
// would be part of the host the token is sent to.
func HostsFor(region, ccRegion string) (Hosts, error) {
	if ccRegion == "" {
		ccRegion = DefaultCCRegion
	}
	if !slices.Contains(ccRegions, ccRegion) {
		return Hosts{}, fmt.Errorf("unknown Contact Center data center %q; use one of %s", ccRegion, strings.Join(ccRegions, ", "))
	}
	return Hosts{
		Webex:            "https://webexapis.com/v1",
		Identity:         "https://webexapis.com",
		Analytics:        "https://analytics.webexapis.com",
		AnalyticsCalling: AnalyticsCallingHost(region),
		CC:               "https://api.wxcc-" + ccRegion + ".cisco.com",
		AIAgent:          "https://studio.aiagent-" + ccRegion + ".cisco.com",
	}, nil
}

// AnalyticsCallingHost is Detailed Call History's base URL for a data region.
func AnalyticsCallingHost(region string) string {
	host, ok := analyticsCallingHosts[strings.ToLower(strings.TrimSpace(region))]
	if !ok {
		host = analyticsCallingHosts["us"]
	}
	return host + "/v1"
}

// DecodeOrgID converts a base64 Webex organization ID
// (ciscospark://us/ORGANIZATION/<uuid>) to its UUID. A UUID or an unrecognized
// value is returned as it is.
func DecodeOrgID(id string) string {
	if id == "" || isUUID(id) {
		return id
	}
	// Webex IDs come with and without padding.
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
		decoded, err := enc.DecodeString(id)
		if err != nil {
			continue
		}
		s := string(decoded)
		if strings.HasPrefix(s, "ciscospark://") {
			parts := strings.Split(s, "/")
			if uuid := parts[len(parts)-1]; uuid != "" {
				return uuid
			}
		}
	}
	return id
}

// EncodeOrgID converts an organization UUID to the base64 Webex form, unpadded as
// Webex issues it; Calling's multipart uploads never answer a padded one. Any
// other value is returned as it is.
func EncodeOrgID(id string) string {
	if !isUUID(id) {
		return id
	}
	return base64.RawStdEncoding.EncodeToString([]byte("ciscospark://us/ORGANIZATION/" + id))
}

func isUUID(id string) bool { return strings.Count(id, "-") == 4 && len(id) == 36 }
