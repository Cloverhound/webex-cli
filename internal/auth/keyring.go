package auth

import (
	"encoding/json"
	"fmt"
	"time"
)

const serviceName = "webex-cli"

// StoredToken represents OAuth tokens persisted in the token store.
type StoredToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
	IssuedAt     time.Time `json:"issued_at"`
	ClientID     string    `json:"client_id,omitempty"`
	ClientSecret string    `json:"client_secret,omitempty"`
	Scopes       string    `json:"scopes,omitempty"`
	ReadOnly     bool      `json:"read_only,omitempty"`
}

// IsExpired returns true if the access token is expired or within 60s of expiry.
func (t *StoredToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt.Add(-60 * time.Second))
}

// IsRefreshExpired returns true if the refresh token is likely expired (>90 days since issue).
func (t *StoredToken) IsRefreshExpired() bool {
	return time.Now().After(t.IssuedAt.Add(90 * 24 * time.Hour))
}

// SaveToken stores a token keyed by email, in the OS keyring when one is
// available and in the credentials file otherwise.
func SaveToken(email string, tok *StoredToken) error {
	data, err := json.Marshal(tok)
	if err != nil {
		return fmt.Errorf("marshaling token: %w", err)
	}
	return storeSet(email, string(data))
}

// LoadToken retrieves the token stored for email.
func LoadToken(email string) (*StoredToken, error) {
	tok, _, err := LoadTokenWithStore(email)
	return tok, err
}

// LoadTokenWithStore retrieves the token stored for email and names the store
// that held it ("keyring" or "file").
func LoadTokenWithStore(email string) (*StoredToken, string, error) {
	data, store, err := storeGet(email)
	if err != nil {
		return nil, "", fmt.Errorf("loading token for %s: %w", email, err)
	}
	var tok StoredToken
	if err := json.Unmarshal([]byte(data), &tok); err != nil {
		return nil, "", fmt.Errorf("parsing stored token: %w", err)
	}
	return &tok, store, nil
}

// DeleteToken removes the token for email from every store that holds it.
func DeleteToken(email string) error {
	return storeDelete(email)
}
