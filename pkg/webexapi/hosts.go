package webexapi

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Hosts are the base URLs of the Webex APIs one org uses.
type Hosts struct {
	// Webex serves Admin, Calling, Devices, Meetings and Messaging.
	Webex string

	// CC is Contact Center in the org's data center.
	CC string

	// Analytics serves the Admin analytics reports, whose paths carry their own /v1.
	Analytics string

	// AnalyticsCalling serves Detailed Call History in the org's data region.
	AnalyticsCalling string

	// AIAgent is AI Agent Studio, which sits in the same data center as CC.
	AIAgent string
}

// CCDataCenters lists the Contact Center data centers, each of which has its own
// API and AI Agent Studio host.
var CCDataCenters = []string{"us1", "eu1", "eu2", "anz1", "ca1", "jp1", "sg1"}

// DataRegions lists the Webex data regions that have their own analytics host.
var DataRegions = []string{"us", "ca", "eu", "eun", "in", "au"}

// analyticsCallingHosts maps a data region to its Detailed Call History host. The
// default host routes to the nearest region and answers 451 with the right one when
// that region does not hold the org's data.
var analyticsCallingHosts = map[string]string{
	"us":  "https://analytics-calling.webexapis.com",
	"ca":  "https://analytics-calling.webexapis.com",
	"eu":  "https://analytics-calling-eu.webexapis.com",
	"eun": "https://analytics-calling-eu.webexapis.com",
	"in":  "https://analytics-calling-in.webexapis.com",
	"au":  "https://analytics-calling-au.webexapis.com",
}

// NewHosts returns the hosts for a Contact Center data center and a data region.
// Empty values mean us1 and us.
func NewHosts(ccDataCenter, dataRegion string) (Hosts, error) {
	dc := strings.ToLower(strings.TrimSpace(ccDataCenter))
	if dc == "" {
		dc = "us1"
	}
	if !contains(CCDataCenters, dc) {
		return Hosts{}, fmt.Errorf("unknown Contact Center data center %q (valid: %s)", ccDataCenter, strings.Join(CCDataCenters, ", "))
	}
	region := strings.ToLower(strings.TrimSpace(dataRegion))
	if region == "" {
		region = "us"
	}
	analytics, ok := analyticsCallingHosts[region]
	if !ok {
		return Hosts{}, fmt.Errorf("unknown data region %q (valid: %s)", dataRegion, strings.Join(DataRegions, ", "))
	}
	return Hosts{
		Webex:            "https://webexapis.com/v1",
		CC:               "https://api.wxcc-" + dc + ".cisco.com",
		Analytics:        "https://analytics.webexapis.com",
		AnalyticsCalling: analytics + "/v1",
		AIAgent:          "https://studio.aiagent-" + dc + ".cisco.com",
	}, nil
}

func contains(list []string, s string) bool {
	for _, x := range list {
		if x == s {
			return true
		}
	}
	return false
}

// OrgUUID returns the UUID inside a Webex org ID (base64 of
// ciscospark://us/ORGANIZATION/<uuid>), which Contact Center paths use. A UUID, or
// anything it does not recognize, is returned unchanged.
func OrgUUID(id string) string {
	if id == "" || isUUID(id) {
		return id
	}
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.RawStdEncoding, base64.URLEncoding, base64.RawURLEncoding} {
		decoded, err := enc.DecodeString(id)
		if err != nil {
			continue
		}
		s := string(decoded)
		if strings.HasPrefix(s, "ciscospark://") {
			if uuid := s[strings.LastIndex(s, "/")+1:]; uuid != "" {
				return uuid
			}
		}
	}
	return id
}

// OrgIDFromUUID returns the Webex org ID for an org UUID, unpadded as Webex issues
// it: Calling multipart uploads never answer a padded orgId. Anything that is not a
// UUID is returned unchanged.
func OrgIDFromUUID(uuid string) string {
	if !isUUID(uuid) {
		return uuid
	}
	return base64.RawStdEncoding.EncodeToString([]byte("ciscospark://us/ORGANIZATION/" + uuid))
}

func isUUID(s string) bool {
	return strings.Count(s, "-") == 4 && len(s) == 36
}
