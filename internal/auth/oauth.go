package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/browser"
)

// openBrowser is a var so tests can avoid launching a real browser.
var openBrowser = func(u string) error {
	// pkg/browser copies the launcher's output to ours, which would corrupt
	// JSON on stdout.
	browser.Stdout = io.Discard
	browser.Stderr = io.Discard
	return browser.OpenURL(u)
}

const (
	AuthorizeURL = "https://webexapis.com/v1/authorize"
	RedirectURI  = "http://localhost:8085/callback"
)

// Vars rather than consts so tests can point them at a fake server.
var (
	TokenURL    = "https://webexapis.com/v1/access_token"
	peopleMeURL = "https://webexapis.com/v1/people/me"
)

// LoginResult contains the token and identity info from a successful login.
type LoginResult struct {
	Token       StoredToken
	Email       string
	DisplayName string
	OrgID       string
	OrgName     string
}

// ErrNoBrowser reports that the browser flow cannot run here: the callback
// port is unavailable or no browser could be opened.
var ErrNoBrowser = errors.New("browser login unavailable")

// Login runs the OAuth PKCE flow: opens browser, waits for callback, exchanges code, fetches identity.
// With failFast, it returns an error wrapping ErrNoBrowser instead of waiting
// when the browser cannot be opened, so the caller can switch to device login.
func Login(clientID, clientSecret, scopes string, failFast bool) (*LoginResult, error) {
	if clientID == "" {
		return nil, fmt.Errorf("client ID not configured — run: webex config set client-id <YOUR_CLIENT_ID>")
	}

	// Generate PKCE verifier + challenge
	verifier, challenge, err := generatePKCE()
	if err != nil {
		return nil, fmt.Errorf("generating PKCE: %w", err)
	}
	state := generateState()

	// Start local server before opening browser
	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	listener, err := net.Listen("tcp", "127.0.0.1:8085")
	if err != nil {
		return nil, fmt.Errorf("%w: starting callback server: %v", ErrNoBrowser, err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != state {
			errCh <- fmt.Errorf("state mismatch")
			fmt.Fprint(w, "<html><body><h2>Error: state mismatch</h2></body></html>")
			return
		}
		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			errCh <- fmt.Errorf("OAuth error: %s — %s", errMsg, r.URL.Query().Get("error_description"))
			fmt.Fprintf(w, "<html><body><h2>Login failed: %s</h2></body></html>", errMsg)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in callback")
			fmt.Fprint(w, "<html><body><h2>Error: no authorization code</h2></body></html>")
			return
		}
		fmt.Fprint(w, "<html><body><h2>Login successful!</h2><p>You can close this tab.</p></body></html>")
		codeCh <- code
	})

	srv := &http.Server{Handler: mux}
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	defer srv.Shutdown(context.Background())

	// Build authorize URL
	authURL := fmt.Sprintf("%s?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s&code_challenge=%s&code_challenge_method=S256&prompt=select_account",
		AuthorizeURL,
		url.QueryEscape(clientID),
		url.QueryEscape(RedirectURI),
		url.QueryEscape(scopes),
		url.QueryEscape(state),
		url.QueryEscape(challenge),
	)

	if failFast {
		if err := openBrowser(authURL); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrNoBrowser, err)
		}
		fmt.Fprintln(os.Stderr, "Opened browser for Webex login. If it didn't open, visit:")
		fmt.Fprintln(os.Stderr, authURL)
	} else {
		fmt.Fprintln(os.Stderr, "Opening browser for Webex login...")
		fmt.Fprintln(os.Stderr, "If browser doesn't open, visit:")
		fmt.Fprintln(os.Stderr, authURL)
		_ = openBrowser(authURL)
	}

	// Wait for callback
	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("login timed out (5 minutes)")
	}

	// Exchange code for tokens
	tok, err := exchangeCode(clientID, clientSecret, code, verifier, scopes)
	if err != nil {
		return nil, err
	}
	return completeLogin(tok)
}

// completeLogin looks up the identity and org behind a newly issued token.
func completeLogin(tok *StoredToken) (*LoginResult, error) {
	identity, err := fetchIdentity(tok.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("fetching identity: %w", err)
	}

	// Fetch org name
	orgName := ""
	if identity.orgID != "" {
		orgName, _ = FetchOrgName(tok.AccessToken, identity.orgID)
	}

	return &LoginResult{
		Token:       *tok,
		Email:       identity.email,
		DisplayName: identity.displayName,
		OrgID:       identity.orgID,
		OrgName:     orgName,
	}, nil
}

// RefreshAccessToken uses a refresh token to get a new access token.
func RefreshAccessToken(clientID, clientSecret string, tok *StoredToken) (*StoredToken, error) {
	if tok.IsRefreshExpired() {
		return nil, fmt.Errorf("refresh token expired — run: webex login")
	}

	form := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"refresh_token": {tok.RefreshToken},
	}

	resp, err := http.PostForm(TokenURL, form)
	if err != nil {
		return nil, fmt.Errorf("refresh request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("refresh failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("parsing refresh response: %w", err)
	}

	// A refresh keeps the original grant, so the scopes and the per-account
	// credentials carry over unchanged.
	refreshed := tokenResp.storedToken(tok.Scopes)
	refreshed.ClientID = tok.ClientID
	refreshed.ClientSecret = tok.ClientSecret
	refreshed.ReadOnly = tok.ReadOnly
	return refreshed, nil
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

// grantedScopes prefers the scopes Webex reports granting and otherwise
// records the requested ones.
func (tr tokenResponse) grantedScopes(requested string) string {
	if tr.Scope != "" {
		return tr.Scope
	}
	return requested
}

// storedToken converts a token response issued now.
func (tr tokenResponse) storedToken(scopes string) *StoredToken {
	now := time.Now()
	return &StoredToken{
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		ExpiresAt:    now.Add(time.Duration(tr.ExpiresIn) * time.Second),
		TokenType:    tr.TokenType,
		IssuedAt:     now,
		Scopes:       scopes,
	}
}

func exchangeCode(clientID, clientSecret, code, verifier, scopes string) (*StoredToken, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {RedirectURI},
		"code_verifier": {verifier},
	}

	resp, err := http.PostForm(TokenURL, form)
	if err != nil {
		return nil, fmt.Errorf("token exchange: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}
	return tokenResp.storedToken(tokenResp.grantedScopes(scopes)), nil
}

type identityInfo struct {
	email       string
	displayName string
	orgID       string
}

func fetchIdentity(accessToken string) (*identityInfo, error) {
	req, _ := http.NewRequest("GET", peopleMeURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET /people/me failed (%d): %s", resp.StatusCode, string(body))
	}

	var person struct {
		Emails      []string `json:"emails"`
		DisplayName string   `json:"displayName"`
		OrgID       string   `json:"orgId"`
	}
	if err := json.Unmarshal(body, &person); err != nil {
		return nil, err
	}

	email := ""
	if len(person.Emails) > 0 {
		email = person.Emails[0]
	}
	return &identityInfo{
		email:       email,
		displayName: person.DisplayName,
		orgID:       person.OrgID,
	}, nil
}

func FetchOrgName(accessToken, orgID string) (string, error) {
	req, _ := http.NewRequest("GET", "https://webexapis.com/v1/organizations/"+url.PathEscape(orgID), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GET /organizations/%s failed (%d)", orgID, resp.StatusCode)
	}

	var org struct {
		DisplayName string `json:"displayName"`
	}
	if err := json.Unmarshal(body, &org); err != nil {
		return "", err
	}
	return org.DisplayName, nil
}

func generatePKCE() (verifier, challenge string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	verifier = base64.RawURLEncoding.EncodeToString(buf)
	h := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h[:])
	return verifier, challenge, nil
}

func generateState() string {
	buf := make([]byte, 16)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
