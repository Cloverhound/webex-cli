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
)

func SetToken(t string) { token = t }
func Token() string      { return token }

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
)

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
