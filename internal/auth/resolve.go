package auth

import (
	"fmt"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
)

// TokenSource indicates where the token was resolved from.
type TokenSource int

const (
	SourceNone    TokenSource = iota
	SourceFlag                // --token flag
	SourceEnv                 // $WEBEX_TOKEN
	SourceKeyring             // OS keyring
)

func (s TokenSource) String() string {
	switch s {
	case SourceFlag:
		return "flag"
	case SourceEnv:
		return "environment"
	case SourceKeyring:
		return "keyring"
	default:
		return "none"
	}
}

// ResolveResult contains the resolved token and metadata.
type ResolveResult struct {
	Token     string
	Source    TokenSource
	UserEmail string
	OrgID     string
}

// ResolveToken determines the access token using the priority chain:
// 1. --token flag
// 2. $WEBEX_TOKEN env
// 3. Keyring (resolve user from --user flag → $WEBEX_USER → config default)
func ResolveToken(flagToken, envToken, userFlag, envUser string, cfg *appconfig.Config) (*ResolveResult, error) {
	// 1. Explicit --token flag
	if flagToken != "" {
		return &ResolveResult{Token: flagToken, Source: SourceFlag}, nil
	}

	// 2. Environment variable
	if envToken != "" {
		return &ResolveResult{Token: envToken, Source: SourceEnv}, nil
	}

	// 3. Keyring lookup
	email := userFlag
	if email == "" {
		email = envUser
	}
	if email == "" {
		email = cfg.DefaultUser
	}
	if email == "" {
		return nil, fmt.Errorf("no authenticated user — run: webex login")
	}

	tok, err := LoadToken(email)
	if err != nil {
		return nil, fmt.Errorf("no token for %s — run: webex login", email)
	}

	// Auto-refresh if expired
	if tok.IsExpired() {
		clientID := cfg.EffectiveClientID()
		clientSecret := cfg.EffectiveClientSecret()
		refreshed, err := RefreshAccessToken(clientID, clientSecret, tok)
		if err != nil {
			return nil, fmt.Errorf("token expired for %s and refresh failed: %w\nRun: webex login", email, err)
		}
		tok = refreshed
		// Persist the refreshed token
		if saveErr := SaveToken(email, tok); saveErr != nil {
			fmt.Printf("Warning: could not save refreshed token: %v\n", saveErr)
		}
	}

	// Determine org ID from config
	orgID := ""
	if userInfo, ok := cfg.Users[email]; ok {
		orgID = userInfo.OrgID
	}

	return &ResolveResult{
		Token:     tok.AccessToken,
		Source:    SourceKeyring,
		UserEmail: email,
		OrgID:     orgID,
	}, nil
}

// MakeRefresher returns a callback that refreshes the token for the given user.
func MakeRefresher(email string, cfg *appconfig.Config) func() (string, error) {
	return func() (string, error) {
		tok, err := LoadToken(email)
		if err != nil {
			return "", err
		}
		refreshed, err := RefreshAccessToken(cfg.EffectiveClientID(), cfg.EffectiveClientSecret(), tok)
		if err != nil {
			return "", err
		}
		if saveErr := SaveToken(email, refreshed); saveErr != nil {
			fmt.Printf("Warning: could not save refreshed token: %v\n", saveErr)
		}
		return refreshed.AccessToken, nil
	}
}
