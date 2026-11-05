package config

import (
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
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
func Token() string     { return token }

func SetMaxRetry(n int) { maxRetry = n }
func MaxRetry() int     { return maxRetry }

func SetMaxRetryTimer(secs int) { maxRetryTimer = secs }
func MaxRetryTimer() int        { return maxRetryTimer }

func SetDebug(d bool) { debug = d }
func Debug() bool     { return debug }

func SetPaginate(p bool) { paginate = p }
func Paginate() bool     { return paginate }

func SetDryRun(d bool) { dryRun = d }
func DryRun() bool     { return dryRun }

var changesetPath string

// SetChangesetPath turns on capture mode: write requests are recorded in the
// change set file at path instead of being sent.
func SetChangesetPath(path string) { changesetPath = path }
func ChangesetPath() string        { return changesetPath }

var readOnly bool

// SetReadOnly turns on read-only request checks for this process.
func SetReadOnly(on bool) { readOnly = on }
func ReadOnly() bool      { return readOnly }

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

func OrgID() string       { return orgID }
func OrgIDBase64() string { return orgIDBase64 }

// Regions lists the supported --region values.
var Regions = webexapi.DataRegions

// CCRegions lists the supported --cc-region values: Contact Center data centers.
var CCRegions = webexapi.CCDataCenters

var (
	region, ccRegion string
	hosts, _         = webexapi.NewHosts("", "")
)

// SetRegion stores the org's data region (us, ca, eu, eun, in, au).
// Empty means "use the default host and let Webex route by geography".
// An unknown region keeps the default host; callers validate first.
func SetRegion(r string) {
	region = strings.ToLower(strings.TrimSpace(r))
	updateHosts()
}
func Region() string { return region }

// SetCCRegion stores the org's Contact Center data center (us1, eu1, ...),
// which selects the Contact Center and AI Agent Studio hosts. Empty means us1.
func SetCCRegion(r string) {
	ccRegion = strings.ToLower(strings.TrimSpace(r))
	updateHosts()
}
func CCRegion() string { return ccRegion }

func updateHosts() {
	if h, err := webexapi.NewHosts(ccRegion, region); err == nil {
		hosts = h
	}
}

var ccInterval = webexapi.CCInterval

// SetCCInterval spaces Contact Center requests; zero sends them unpaced.
func SetCCInterval(d time.Duration) { ccInterval = d }
func CCInterval() time.Duration     { return ccInterval }

// Hosts returns the API base URLs for this run's regions.
func Hosts() webexapi.Hosts { return hosts }

func CallingBaseURL() string { return hosts.Webex }
func CcBaseURL() string      { return hosts.CC }

// AnalyticsBaseURL serves the Admin analytics reports. Paths for those
// endpoints already carry their own /v1 segment.
func AnalyticsBaseURL() string { return hosts.Analytics }

// AiAgentBaseURL is AI Agent Studio in the same data center as CcBaseURL.
func AiAgentBaseURL() string { return hosts.AIAgent }

// AnalyticsCallingBaseURL returns the Detailed Call History base URL for the
// configured region. Calls to the default host are routed to the nearest
// region; if that region does not hold the org's data the API returns HTTP 451
// with the correct endpoint.
func AnalyticsCallingBaseURL() string { return hosts.AnalyticsCalling }

// DecodeOrgID converts a base64-encoded Webex org ID (ciscospark://us/ORGANIZATION/<uuid>)
// to the raw UUID. If the input is already a UUID or unrecognized, it is returned as-is.
func DecodeOrgID(id string) string { return webexapi.OrgUUID(id) }

// EncodeOrgID converts a UUID to the base64-encoded Webex org ID format,
// unpadded as Webex issues it. Anything else is returned as-is.
func EncodeOrgID(id string) string { return webexapi.OrgIDFromUUID(id) }
