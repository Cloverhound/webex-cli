package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zalando/go-keyring"
)

const serviceName = "webex-cli"

// StoredToken represents OAuth tokens persisted in the OS keyring.
type StoredToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
	IssuedAt     time.Time `json:"issued_at"`
}

// IsExpired returns true if the access token is expired or within 60s of expiry.
func (t *StoredToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt.Add(-60 * time.Second))
}

// IsRefreshExpired returns true if the refresh token is likely expired (>90 days since issue).
func (t *StoredToken) IsRefreshExpired() bool {
	return time.Now().After(t.IssuedAt.Add(90 * 24 * time.Hour))
}

// SaveToken stores a token in the OS keyring keyed by email.
func SaveToken(email string, tok *StoredToken) error {
	data, err := json.Marshal(tok)
	if err != nil {
		return fmt.Errorf("marshaling token: %w", err)
	}
	return keyring.Set(serviceName, email, string(data))
}

// LoadToken retrieves a token from the OS keyring for the given email.
func LoadToken(email string) (*StoredToken, error) {
	data, err := keyring.Get(serviceName, email)
	if err != nil {
		return nil, fmt.Errorf("loading token for %s: %w", email, err)
	}
	var tok StoredToken
	if err := json.Unmarshal([]byte(data), &tok); err != nil {
		return nil, fmt.Errorf("parsing stored token: %w", err)
	}
	return &tok, nil
}

// DeleteToken removes a token from the OS keyring.
func DeleteToken(email string) error {
	err := keyring.Delete(serviceName, email)
	if err == keyring.ErrNotFound {
		return nil
	}
	return err
}
