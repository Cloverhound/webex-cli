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
	_, err := SaveTokenWithStore(email, tok)
	return err
}

// SaveTokenWithStore stores a token like SaveToken and names the store that
// now holds it ("keyring" or "file").
func SaveTokenWithStore(email string, tok *StoredToken) (string, error) {
	data, err := json.Marshal(tok)
	if err != nil {
		return "", fmt.Errorf("marshaling token: %w", err)
	}
	return storeSet(email, string(data))
}

// LoadToken retrieves the token stored for email.
func LoadToken(email string) (*StoredToken, error) {
	tok, _, err := LoadTokenWithStore(email)
	return tok, err
}

// LoadTokenWithStore retrieves the token stored for email and names the store
// that held it ("keyring" or "file"). When both stores hold a copy, as after a
// save fell back to the file while the keyring was down, the newest one wins.
func LoadTokenWithStore(email string) (*StoredToken, string, error) {
	stores, err := selectedStores()
	if err != nil {
		return nil, "", err
	}
	copies, err := readCopies(email, stores)
	if err != nil {
		return nil, "", fmt.Errorf("loading token for %s: %w", email, err)
	}
	var best *StoredToken
	var bestStore string
	var parseErr error
	for _, c := range copies {
		tok, err := decodeToken(c.data)
		if err != nil {
			parseErr = err
			continue
		}
		if best == nil || tok.IssuedAt.After(best.IssuedAt) {
			best, bestStore = tok, c.store
		}
	}
	if best == nil {
		return nil, "", parseErr
	}
	return best, bestStore, nil
}

func decodeToken(data string) (*StoredToken, error) {
	var tok StoredToken
	if err := json.Unmarshal([]byte(data), &tok); err != nil {
		return nil, fmt.Errorf("parsing stored token: %w", err)
	}
	return &tok, nil
}

// DeleteToken removes the token for email from every store that holds it.
func DeleteToken(email string) error {
	return storeDelete(email)
}
