package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Vars rather than consts so tests can point them at a fake server.
var (
	DeviceAuthorizeURL = "https://webexapis.com/v1/device/authorize"
	DeviceTokenURL     = "https://webexapis.com/v1/device/token"

	pollUnit = time.Second
)

const deviceRequestTimeout = 30 * time.Second

var errDeviceCodeExpired = errors.New("the login code expired before it was approved — run webex login again")

// DeviceHelperRedirectURIs must be registered on the integration: Webex's device
// approval page redirects to one of them, and which one depends on the region.
const DeviceHelperRedirectURIs = "https://oauth-helper-{a,r,k,d}.wbx2.com/helperservice/v1/actions/device/callback"

// DeviceCode is the response to a device authorization request.
type DeviceCode struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// RequestDeviceCode starts the OAuth device grant.
func RequestDeviceCode(clientID, scopes string) (*DeviceCode, error) {
	if clientID == "" {
		return nil, fmt.Errorf("client ID not configured — run: webex config set client-id <YOUR_CLIENT_ID>")
	}
	ctx, cancel := context.WithTimeout(context.Background(), deviceRequestTimeout)
	defer cancel()
	form := url.Values{"client_id": {clientID}, "scope": {scopes}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, DeviceAuthorizeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("device authorization request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if strings.Contains(string(body), "redirect_uri") {
			return nil, fmt.Errorf("device login failed (%d): the integration is missing a device-grant redirect URI; add %s to it.\n%s",
				resp.StatusCode, DeviceHelperRedirectURIs, string(body))
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			return nil, fmt.Errorf("device login is rate limited (retry after %s seconds)", resp.Header.Get("Retry-After"))
		}
		return nil, fmt.Errorf("device authorization failed (%d): %s", resp.StatusCode, string(body))
	}

	var dc DeviceCode
	if err := json.Unmarshal(body, &dc); err != nil {
		return nil, fmt.Errorf("parsing device authorization response: %w", err)
	}
	if dc.DeviceCode == "" || dc.UserCode == "" {
		return nil, fmt.Errorf("device authorization response is missing the device or user code: %s", string(body))
	}
	if dc.ExpiresIn <= 0 {
		dc.ExpiresIn = 300
	}
	if dc.Interval <= 0 {
		dc.Interval = 5
	}
	return &dc, nil
}

// PollDeviceToken polls the token endpoint until the user approves the code,
// denies it, or the code expires. It follows RFC 8628: authorization_pending
// keeps polling and slow_down adds 5 seconds to the interval.
func PollDeviceToken(ctx context.Context, clientID, clientSecret string, dc *DeviceCode, scopes string) (*StoredToken, error) {
	interval := time.Duration(dc.Interval) * pollUnit
	parent := ctx
	// Bounding every request by the code's expiry keeps a stalled connection from
	// outliving the code.
	ctx, cancel := context.WithTimeout(ctx, time.Duration(dc.ExpiresIn)*pollUnit)
	defer cancel()
	expired := func() error {
		if parent.Err() != nil {
			return parent.Err()
		}
		return errDeviceCodeExpired
	}

	for {
		select {
		case <-ctx.Done():
			return nil, expired()
		case <-time.After(interval):
		}

		tok, status, err := requestDeviceToken(ctx, clientID, clientSecret, dc.DeviceCode, scopes)
		if err != nil {
			if ctx.Err() != nil {
				return nil, expired()
			}
			return nil, err
		}
		switch status {
		case "":
			return tok, nil
		case "authorization_pending":
		case "slow_down":
			interval += 5 * pollUnit
		case "access_denied":
			return nil, fmt.Errorf("the login request was denied")
		case "expired_token", "invalid_grant":
			return nil, fmt.Errorf("the login code expired or is no longer valid — run webex login again")
		default:
			return nil, fmt.Errorf("device login failed: %s", status)
		}
	}
}

// requestDeviceToken makes one token request. It returns the token on success,
// or the OAuth error code for a response the caller can act on.
func requestDeviceToken(ctx context.Context, clientID, clientSecret, deviceCode, scopes string) (*StoredToken, string, error) {
	form := url.Values{
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"device_code": {deviceCode},
		"client_id":   {clientID},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, DeviceTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("device token request: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		var tr tokenResponse
		if err := json.Unmarshal(body, &tr); err != nil {
			return nil, "", fmt.Errorf("parsing device token response: %w", err)
		}
		return tr.storedToken(tr.grantedScopes(scopes)), "", nil
	}

	var oauthErr struct {
		Error       string `json:"error"`
		Description string `json:"error_description"`
		Message     string `json:"message"`
	}
	_ = json.Unmarshal(body, &oauthErr)

	switch {
	case resp.StatusCode == http.StatusPreconditionRequired:
		// Webex answers 428 while the user has not approved yet.
		return nil, "authorization_pending", nil
	case resp.StatusCode == http.StatusTooManyRequests:
		return nil, "slow_down", nil
	case oauthErr.Error != "":
		return nil, oauthErr.Error, nil
	}
	return nil, "", fmt.Errorf("device token request failed (%d): %s", resp.StatusCode, string(body))
}

// DeviceLogin runs the device grant. show is called once with the code for the
// user to approve; the call then blocks until approval, denial or expiry.
func DeviceLogin(clientID, clientSecret, scopes string, show func(*DeviceCode)) (*LoginResult, error) {
	if clientSecret == "" {
		return nil, fmt.Errorf("client secret not configured — run: webex config set client-secret <YOUR_CLIENT_SECRET>")
	}
	dc, err := RequestDeviceCode(clientID, scopes)
	if err != nil {
		return nil, err
	}
	show(dc)

	tok, err := PollDeviceToken(context.Background(), clientID, clientSecret, dc, scopes)
	if err != nil {
		return nil, err
	}
	return completeLogin(tok)
}
