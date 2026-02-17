package config

import (
	"encoding/base64"
	"strings"
)

var (
	token          string
	debug          bool
	paginate       bool
	orgID          string
	TokenRefresher func() (string, error)
)

func SetToken(t string) { token = t }
func Token() string      { return token }

func SetDebug(d bool) { debug = d }
func Debug() bool     { return debug }

func SetPaginate(p bool) { paginate = p }
func Paginate() bool     { return paginate }

// SetOrgID stores the org ID, auto-decoding base64 Webex IDs to UUID.
func SetOrgID(id string) { orgID = DecodeOrgID(id) }
func OrgID() string      { return orgID }

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
