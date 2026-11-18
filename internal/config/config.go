package config

import (
	"net/url"
	"slices"
	"strings"

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

var defaultHosts, _ = webexapi.HostsFor("", "")

var (
	CallingBaseURL = defaultHosts.Webex

	// IdentityBaseURL serves SCIM and archived users, whose paths have no /v1 segment.
	IdentityBaseURL = defaultHosts.Identity

	// AnalyticsBaseURL serves the Admin analytics reports. Paths for those
	// endpoints already carry their own /v1 segment.
	AnalyticsBaseURL = defaultHosts.Analytics
)

// CCRegions lists the Contact Center data centers. An org's Contact Center lives
// in one of them, and each has its own API and AI Agent Studio host.
var CCRegions = webexapi.CCRegions

// DefaultCCRegion is used when no region is given or saved for the org.
const DefaultCCRegion = webexapi.DefaultCCRegion

var (
	ccRegion = DefaultCCRegion

	CcBaseURL = defaultHosts.CC

	// AiAgentBaseURL is AI Agent Studio in the same data center as CcBaseURL.
	AiAgentBaseURL = defaultHosts.AIAgent
)

// SetCCRegion points Contact Center and AI Agent Studio requests at a data center
// from CCRegions. Empty means DefaultCCRegion; any other value keeps the current
// one, so callers validate it first.
func SetCCRegion(r string) {
	if r == "" {
		r = DefaultCCRegion
	}
	hosts, err := webexapi.HostsFor("", r)
	if err != nil {
		return
	}
	ccRegion = r
	CcBaseURL = hosts.CC
	AiAgentBaseURL = hosts.AIAgent
}

// Regions lists the supported --region values.
var Regions = webexapi.Regions

var region string

// SetRegion stores the org's data region (us, ca, eu, eun, in, au).
// Empty means "use the default host and let Webex route by geography".
func SetRegion(r string) { region = strings.ToLower(strings.TrimSpace(r)) }
func Region() string     { return region }

// AnalyticsCallingBaseURL returns the Detailed Call History base URL for the
// configured region, falling back to the US/Canada host.
func AnalyticsCallingBaseURL() string { return webexapi.AnalyticsCallingHost(region) }

// Hosts are this run's API hosts, for its data region and Contact Center data
// center.
func Hosts() webexapi.Hosts {
	hosts, _ := webexapi.HostsFor(region, ccRegion) // ccRegion is always one of CCRegions
	return hosts
}

// trusted are the origins added by TrustOrigin.
var trusted []string

// Origins are the origins the CLI sends its token to: the Webex origins, and any
// a test added with TrustOrigin.
func Origins() []string { return append(webexapi.WebexOrigins(), trusted...) }

// TrustOrigin lets a test send the token to the server at rawURL. It returns a
// function that withdraws it.
func TrustOrigin(rawURL string) (withdraw func()) {
	saved := trusted
	if u, err := url.Parse(rawURL); err == nil {
		trusted = append(slices.Clone(trusted), u.Scheme+"://"+u.Host)
	}
	return func() { trusted = saved }
}

// DecodeOrgID converts a base64 Webex org ID to its UUID; see webexapi.DecodeOrgID.
func DecodeOrgID(id string) string { return webexapi.DecodeOrgID(id) }

// EncodeOrgID converts an org UUID to the base64 Webex form; see webexapi.EncodeOrgID.
func EncodeOrgID(id string) string { return webexapi.EncodeOrgID(id) }
