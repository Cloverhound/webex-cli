package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Environment variables for non-interactive auth from a refresh token.
const (
	RefreshTokenEnv = "WEBEX_REFRESH_TOKEN"
	ClientIDEnv     = "WEBEX_CLIENT_ID"
	ClientSecretEnv = "WEBEX_CLIENT_SECRET"
)

// envRefreshEntry caches the access token minted from $WEBEX_REFRESH_TOKEN, and
// the rotated refresh token, so each CLI call does not mint a new one.
type envRefreshEntry struct {
	Token StoredToken `json:"token"`
	Email string      `json:"email,omitempty"`
	OrgID string      `json:"org_id,omitempty"`
}

const envRefreshKeyPrefix = "env-refresh:"

// The cache is keyed by a hash of the env value so a changed secret starts fresh
// and the original value never lands on disk.
func envRefreshKey(refreshToken string) string {
	sum := sha256.Sum256([]byte(refreshToken))
	return envRefreshKeyPrefix + hex.EncodeToString(sum[:8])
}

// DeleteEnvRefreshCache removes every cached token minted from $WEBEX_REFRESH_TOKEN.
func DeleteEnvRefreshCache() error {
	return fileDeletePrefix(envRefreshKeyPrefix)
}

// EnvRefreshAccessToken returns an access token for the refresh token in
// $WEBEX_REFRESH_TOKEN. A cached access token is reused until it expires or
// equals stale, the token the API just rejected. The cache lives in the
// credentials file whatever $WEBEX_TOKEN_STORE says, since CI and containers
// rarely have a keyring.
func EnvRefreshAccessToken(envRefresh, clientID, clientSecret, stale string) (token, email, orgID string, err error) {
	key := envRefreshKey(envRefresh)
	var entry *envRefreshEntry

	lockErr := withFileLock("refresh.lock", func() error {
		cached := loadEnvRefreshEntry(key)
		if cached != nil && !cached.Token.IsExpired() && cached.Token.AccessToken != stale {
			entry = cached
			return nil
		}

		refreshToken := envRefresh
		if cached != nil && cached.Token.RefreshToken != "" {
			refreshToken = cached.Token.RefreshToken
		}
		tok, err := RefreshAccessToken(clientID, clientSecret, &StoredToken{RefreshToken: refreshToken, IssuedAt: time.Now()})
		if err != nil && refreshToken != envRefresh {
			// The cached refresh token may have been revoked while the secret was
			// updated in place; the env value is the source of truth.
			tok, err = RefreshAccessToken(clientID, clientSecret, &StoredToken{RefreshToken: envRefresh, IssuedAt: time.Now()})
		}
		if err != nil {
			return fmt.Errorf("$%s refresh failed: %w", RefreshTokenEnv, err)
		}
		if tok.RefreshToken == "" {
			tok.RefreshToken = refreshToken
		}

		entry = &envRefreshEntry{Token: *tok}
		if cached != nil {
			entry.Email, entry.OrgID = cached.Email, cached.OrgID
		}
		if entry.OrgID == "" {
			if id, err := fetchIdentity(tok.AccessToken); err == nil {
				entry.Email, entry.OrgID = id.email, id.orgID
			}
		}
		if data, err := json.Marshal(entry); err == nil {
			if err := fileSet(key, string(data)); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not cache the token from $%s (%v); the rotated refresh token is lost when this process exits\n", RefreshTokenEnv, err)
			}
		}
		return nil
	})
	if lockErr != nil {
		return "", "", "", lockErr
	}
	return entry.Token.AccessToken, entry.Email, entry.OrgID, nil
}

func loadEnvRefreshEntry(key string) *envRefreshEntry {
	data, err := fileGet(key)
	if err != nil {
		return nil
	}
	var e envRefreshEntry
	if json.Unmarshal([]byte(data), &e) != nil {
		return nil
	}
	return &e
}
