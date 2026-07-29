package config

import (
	"encoding/base64"
	"strings"
)

var (
	token          string
	debug          bool
	paginate       bool
	dryRun         bool
	orgID          string // UUID format
	orgIDBase64    string // base64 Webex format
	TokenRefresher func() (string, error)

	maxRetry      int // max number of 429 retries (0 = no retries)
	maxRetryTimer int // max total seconds allowed across all 429 waits (0 = unlimited)
)

func SetToken(t string) { token = t }
func Token() string      { return token }

func SetMaxRetry(n int)  { maxRetry = n }
func MaxRetry() int      { return maxRetry }

func SetMaxRetryTimer(secs int) { maxRetryTimer = secs }
func MaxRetryTimer() int        { return maxRetryTimer }

func SetDebug(d bool) { debug = d }
func Debug() bool     { return debug }

func SetPaginate(p bool) { paginate = p }
func Paginate() bool     { return paginate }

func SetDryRun(d bool) { dryRun = d }
func DryRun() bool     { return dryRun }

// SetOrgID stores the org ID in both UUID and base64 formats.
// Accepts either format as input and derives the other.
func SetOrgID(id string) {
	if id == "" {
		orgID = ""
		orgIDBase64 = ""
		return
	}
	uuid := DecodeOrgID(id)
	orgID = uuid
	orgIDBase64 = EncodeOrgID(uuid)
}

func OrgID() string      { return orgID }
func OrgIDBase64() string { return orgIDBase64 }

const (
	CallingBaseURL = "https://webexapis.com/v1"
	CcBaseURL      = "https://api.wxcc-us1.cisco.com"

	// AnalyticsBaseURL serves the Admin analytics reports. Paths for those
	// endpoints already carry their own /v1 segment.
	AnalyticsBaseURL = "https://analytics.webexapis.com"
)

// analyticsCallingHosts maps a data region to the Detailed Call History (CDR) FQDN.
// Calls to the default host are routed to the nearest region; if that region does
// not hold the org's data the API returns HTTP 451 with the correct endpoint.
var analyticsCallingHosts = map[string]string{
	"us":  "https://analytics-calling.webexapis.com",
	"ca":  "https://analytics-calling.webexapis.com",
	"eu":  "https://analytics-calling-eu.webexapis.com",
	"eun": "https://analytics-calling-eu.webexapis.com",
	"in":  "https://analytics-calling-in.webexapis.com",
	"au":  "https://analytics-calling-au.webexapis.com",
}

// Regions lists the supported --region values.
var Regions = []string{"us", "ca", "eu", "eun", "in", "au"}

var region string

// SetRegion stores the org's data region (us, ca, eu, eun, in, au).
// Empty means "use the default host and let Webex route by geography".
func SetRegion(r string) { region = strings.ToLower(strings.TrimSpace(r)) }
func Region() string     { return region }

// AnalyticsCallingBaseURL returns the Detailed Call History base URL for the
// configured region, falling back to the US/Canada host.
func AnalyticsCallingBaseURL() string {
	host, ok := analyticsCallingHosts[region]
	if !ok {
		host = analyticsCallingHosts["us"]
	}
	return host + "/v1"
}

// DecodeOrgID converts a base64-encoded Webex org ID (ciscospark://us/ORGANIZATION/<uuid>)
// to the raw UUID. If the input is already a UUID or unrecognized, it is returned as-is.
func DecodeOrgID(id string) string {
	if id == "" {
		return id
	}

	// If it already looks like a UUID, return as-is
	if strings.Count(id, "-") == 4 && len(id) == 36 {
		return id
	}

	// Try all base64 variants (Webex IDs may or may not have padding)
	encodings := []*base64.Encoding{
		base64.StdEncoding,
		base64.RawStdEncoding,
		base64.URLEncoding,
		base64.RawURLEncoding,
	}

	for _, enc := range encodings {
		decoded, err := enc.DecodeString(id)
		if err != nil {
			continue
		}
		s := string(decoded)
		if strings.HasPrefix(s, "ciscospark://") {
			parts := strings.Split(s, "/")
			if len(parts) > 0 {
				uuid := parts[len(parts)-1]
				if uuid != "" {
					return uuid
				}
			}
		}
	}

	return id
}

// EncodeOrgID converts a UUID to the base64-encoded Webex org ID format:
// base64("ciscospark://us/ORGANIZATION/<uuid>").
// If the input is empty, returns empty. If the input is already base64-encoded, returns as-is.
func EncodeOrgID(id string) string {
	if id == "" {
		return ""
	}
	// If it looks like a UUID, encode it
	if strings.Count(id, "-") == 4 && len(id) == 36 {
		uri := "ciscospark://us/ORGANIZATION/" + id
		return base64.StdEncoding.EncodeToString([]byte(uri))
	}
	// Already base64 or unknown format — return as-is
	return id
}
