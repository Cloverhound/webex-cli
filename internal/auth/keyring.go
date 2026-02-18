package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/zalando/go-keyring"
)

const serviceName = "webex-cli"

var fallbackWarningOnce sync.Once

// isKeyringUnavailable returns true if the error indicates the OS keyring
// service itself is missing (not just that a key wasn't found).
func isKeyringUnavailable(err error) bool {
	if err == nil || err == keyring.ErrNotFound {
		return false
	}
	if isExecNotFound(err) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "org.freedesktop.secrets") ||
		strings.Contains(msg, "org.freedesktop.Secret") ||
		strings.Contains(msg, "secret service") ||
		strings.Contains(msg, "Secret Service") ||
		strings.Contains(msg, "not provided")
}

// isExecNotFound checks whether the error chain contains exec.ErrNotFound.
func isExecNotFound(err error) bool {
	for err != nil {
		if err == exec.ErrNotFound {
			return true
		}
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			break
		}
		err = unwrapped
	}
	return false
}

func printFallbackWarning() {
	fallbackWarningOnce.Do(func() {
		fmt.Fprintln(os.Stderr, "Warning: OS keyring unavailable, storing tokens in ~/.webex-cli/tokens.json")
	})
}

func tokensFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}
	return filepath.Join(home, ".webex-cli", "tokens.json"), nil
}

func readTokensFile() (map[string]StoredToken, error) {
	p, err := tokensFilePath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]StoredToken), nil
		}
		return nil, fmt.Errorf("reading tokens file: %w", err)
	}
	var tokens map[string]StoredToken
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, fmt.Errorf("parsing tokens file: %w", err)
	}
	if tokens == nil {
		tokens = make(map[string]StoredToken)
	}
	return tokens, nil
}

func writeTokensFile(tokens map[string]StoredToken) error {
	p, err := tokensFilePath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating tokens dir: %w", err)
	}
	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling tokens: %w", err)
	}
	if err := os.WriteFile(p, data, 0600); err != nil {
		return fmt.Errorf("writing tokens file: %w", err)
	}
	return nil
}

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
// Falls back to file-based storage if the keyring is unavailable.
func SaveToken(email string, tok *StoredToken) error {
	data, err := json.Marshal(tok)
	if err != nil {
		return fmt.Errorf("marshaling token: %w", err)
	}
	err = keyring.Set(serviceName, email, string(data))
	if err == nil {
		return nil
	}
	if !isKeyringUnavailable(err) {
		return err
	}
	printFallbackWarning()
	tokens, err := readTokensFile()
	if err != nil {
		return err
	}
	tokens[email] = *tok
	return writeTokensFile(tokens)
}

// LoadToken retrieves a token from the OS keyring for the given email.
// Falls back to file-based storage if the keyring is unavailable.
func LoadToken(email string) (*StoredToken, error) {
	data, err := keyring.Get(serviceName, email)
	if err == nil {
		var tok StoredToken
		if err := json.Unmarshal([]byte(data), &tok); err != nil {
			return nil, fmt.Errorf("parsing stored token: %w", err)
		}
		return &tok, nil
	}
	if !isKeyringUnavailable(err) {
		return nil, fmt.Errorf("loading token for %s: %w", email, err)
	}
	printFallbackWarning()
	tokens, err := readTokensFile()
	if err != nil {
		return nil, err
	}
	tok, ok := tokens[email]
	if !ok {
		return nil, fmt.Errorf("loading token for %s: %w", email, keyring.ErrNotFound)
	}
	return &tok, nil
}

// DeleteToken removes a token from the OS keyring.
// Falls back to file-based storage if the keyring is unavailable.
func DeleteToken(email string) error {
	err := keyring.Delete(serviceName, email)
	if err == nil || err == keyring.ErrNotFound {
		return nil
	}
	if !isKeyringUnavailable(err) {
		return err
	}
	printFallbackWarning()
	tokens, err := readTokensFile()
	if err != nil {
		return err
	}
	delete(tokens, email)
	return writeTokensFile(tokens)
}
