package auth

import (
	"fmt"
	"os"
	"sync"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
)

// TokenSource indicates where the token was resolved from.
type TokenSource int

const (
	SourceNone       TokenSource = iota
	SourceFlag                   // --token flag
	SourceEnv                    // $WEBEX_TOKEN
	SourceStored                 // stored login (OS keyring or credentials file)
	SourceEnvRefresh             // $WEBEX_REFRESH_TOKEN
)

func (s TokenSource) String() string {
	switch s {
	case SourceFlag:
		return "flag"
	case SourceEnv:
		return "environment"
	case SourceStored:
		return "stored login"
	case SourceEnvRefresh:
		return "environment refresh token"
	default:
		return "none"
	}
}

// EnvCredentials returns the client credentials for $WEBEX_REFRESH_TOKEN:
// $WEBEX_CLIENT_ID and $WEBEX_CLIENT_SECRET when set, else the configured ones.
func EnvCredentials(cfg *appconfig.Config) (clientID, clientSecret string) {
	clientID, clientSecret = os.Getenv(ClientIDEnv), os.Getenv(ClientSecretEnv)
	if clientID == "" {
		return cfg.EffectiveClientID(), cfg.EffectiveClientSecret()
	}
	return clientID, clientSecret
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
// 3. $WEBEX_REFRESH_TOKEN env
// 4. Stored login (resolve user from --user flag → $WEBEX_USER → config default)
//
// In read-only mode only stored tokens from `webex login --read-only` are accepted.
func ResolveToken(flagToken, envToken, userFlag, envUser string, cfg *appconfig.Config, readOnly bool) (*ResolveResult, error) {
	envRefresh := os.Getenv(RefreshTokenEnv)

	// Webex has no public token introspection, so the scopes of a token the CLI
	// did not obtain itself cannot be checked.
	if readOnly && (flagToken != "" || envToken != "" || envRefresh != "") {
		return nil, fmt.Errorf("read-only mode does not accept --token, $WEBEX_TOKEN or $%s; use a stored read-only login", RefreshTokenEnv)
	}

	// 1. Explicit --token flag
	if flagToken != "" {
		return &ResolveResult{Token: flagToken, Source: SourceFlag}, nil
	}

	// 2. Environment variable
	if envToken != "" {
		return &ResolveResult{Token: envToken, Source: SourceEnv}, nil
	}

	// 3. Refresh token from the environment
	if envRefresh != "" {
		clientID, clientSecret := EnvCredentials(cfg)
		token, email, orgID, err := EnvRefreshAccessToken(envRefresh, clientID, clientSecret, "")
		if err != nil {
			return nil, err
		}
		return &ResolveResult{Token: token, Source: SourceEnvRefresh, UserEmail: email, OrgID: orgID}, nil
	}

	// 4. Stored login
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

	loginCmd := "webex login"
	if readOnly {
		loginCmd = "webex login --read-only"
	}
	tok, err := LoadToken(email)
	if err != nil {
		return nil, fmt.Errorf("no token for %s — run: %s", email, loginCmd)
	}
	if readOnly && !tok.ReadOnly {
		return nil, fmt.Errorf("read-only mode: %s has no read-only login — run: %s", email, loginCmd)
	}

	// Auto-refresh if expired
	if tok.IsExpired() {
		refreshed, err := refreshStored(email, cfg, tok.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("token expired for %s and refresh failed: %w\nRun: webex login", email, err)
		}
		tok = refreshed
	}

	// Determine org ID from config
	orgID := ""
	if userInfo, ok := cfg.Users[email]; ok {
		orgID = userInfo.OrgID
	}

	return &ResolveResult{
		Token:     tok.AccessToken,
		Source:    SourceStored,
		UserEmail: email,
		OrgID:     orgID,
	}, nil
}

// MakeRefresher returns a callback that refreshes the token for the given user.
// current is the access token in use; the callback is called after the API rejects it.
func MakeRefresher(email, current string, cfg *appconfig.Config) func() (string, error) {
	var mu sync.Mutex
	return func() (string, error) {
		mu.Lock()
		defer mu.Unlock()
		refreshed, err := refreshStored(email, cfg, current)
		if err != nil {
			return "", err
		}
		current = refreshed.AccessToken
		return current, nil
	}
}

// MakeEnvRefresher is MakeRefresher for a token minted from $WEBEX_REFRESH_TOKEN.
func MakeEnvRefresher(current string, cfg *appconfig.Config) func() (string, error) {
	var mu sync.Mutex
	return func() (string, error) {
		mu.Lock()
		defer mu.Unlock()
		clientID, clientSecret := EnvCredentials(cfg)
		token, _, _, err := EnvRefreshAccessToken(os.Getenv(RefreshTokenEnv), clientID, clientSecret, current)
		if err != nil {
			return "", err
		}
		current = token
		return current, nil
	}
}

// refreshStored refreshes the stored token for email under a lock shared by all
// CLI processes. Webex refresh tokens may be single-use, so when another process
// has already refreshed (the stored access token is no longer stale), its
// result is reused instead of refreshing again.
func refreshStored(email string, cfg *appconfig.Config, stale string) (*StoredToken, error) {
	var out *StoredToken
	err := withFileLock("refresh.lock", func() error {
		tok, err := LoadToken(email)
		if err != nil {
			return err
		}
		if !tok.IsExpired() && tok.AccessToken != stale {
			out = tok
			return nil
		}
		clientID, clientSecret := effectiveCredentials(tok, cfg)
		refreshed, err := RefreshAccessToken(clientID, clientSecret, tok)
		if err != nil {
			return err
		}
		if saveErr := SaveToken(email, refreshed); saveErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not save refreshed token: %v\n", saveErr)
		}
		out = refreshed
		return nil
	})
	return out, err
}

// effectiveCredentials returns the client ID and secret to use for token refresh.
// Per-token credentials (stored with the user account) take priority over global config.
func effectiveCredentials(tok *StoredToken, cfg *appconfig.Config) (clientID, clientSecret string) {
	if tok.ClientID != "" {
		return tok.ClientID, tok.ClientSecret
	}
	return cfg.EffectiveClientID(), cfg.EffectiveClientSecret()
}
