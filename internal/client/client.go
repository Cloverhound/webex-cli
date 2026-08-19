package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	urlpkg "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/internal/config"
)

// ErrDryRun is returned when a write operation is intercepted by --dry-run mode.
var ErrDryRun = errors.New("dry run: no changes made")

// Do executes an HTTP request.
// On a 401 it refreshes the token and retries once.
// On a 429 it reads Retry-After and retries up to 3 times.
func Do(req *Request) ([]byte, int, error) {
	body, status, headers, err := doOnce(req)

	if status == 401 && config.TokenRefresher != nil {
		newToken, refreshErr := config.TokenRefresher()
		if refreshErr != nil {
			return body, status, err
		}
		config.SetToken(newToken)
		body, status, headers, err = doOnce(req)
	}

	var cumulativeWait time.Duration
	for retries := 0; status == 429; retries++ {
		maxRetry := config.MaxRetry()
		maxTimer := time.Duration(config.MaxRetryTimer()) * time.Second
		wait := retryAfterDuration(headers.Get("Retry-After"))

		if retries >= maxRetry {
			return body, status, fmt.Errorf("rate limited (429): retry limit reached after %d attempts — try again after %s", retries, wait.Round(time.Second))
		}
		if maxTimer > 0 && cumulativeWait+wait > maxTimer {
			return body, status, fmt.Errorf("rate limited (429): retry wait of %s would exceed max-retry-timer of %s — try again after %s", (cumulativeWait + wait).Round(time.Second), maxTimer.Round(time.Second), wait.Round(time.Second))
		}

		fmt.Fprintf(os.Stderr, "Rate limited (429). Retrying in %s (attempt %d/%d)...\n", wait.Round(time.Second), retries+1, maxRetry)
		time.Sleep(wait)
		cumulativeWait += wait
		body, status, headers, err = doOnce(req)
	}

	if status == 451 {
		return body, status, wrongRegionError(body)
	}

	return body, status, err
}

// endpointURLPattern finds the regional endpoint Webex returns in a 451 body.
var endpointURLPattern = regexp.MustCompile(`https?://[a-zA-Z0-9.\-]+`)

// wrongRegionError turns an HTTP 451 (org data lives in another region) into an
// actionable message naming the endpoint the API pointed us at.
func wrongRegionError(body []byte) error {
	msg := fmt.Sprintf("wrong data region (451): this organization's data is not hosted in region %q", regionLabel())
	if m := endpointURLPattern.Find(body); m != nil {
		msg += fmt.Sprintf("; Webex says to use %s", m)
	}
	return fmt.Errorf("%s — rerun with --region <%s> (or: webex config set region <region>)",
		msg, strings.Join(config.Regions, "|"))
}

func regionLabel() string {
	if r := config.Region(); r != "" {
		return r
	}
	return "us (default)"
}

// retryAfterDuration parses the Retry-After header value (seconds integer) and
// returns the duration to wait. Defaults to 5 seconds if the header is absent
// or unparseable.
func retryAfterDuration(header string) time.Duration {
	if secs, err := strconv.Atoi(strings.TrimSpace(header)); err == nil && secs > 0 {
		return time.Duration(secs) * time.Second
	}
	return 5 * time.Second
}

// redactHeader replaces a credential value with its length, which is enough to
// tell a missing or truncated token from a present one without printing it.
func redactHeader(key, value string) string {
	if !strings.EqualFold(key, "Authorization") {
		return value
	}
	if scheme, cred, ok := strings.Cut(value, " "); ok {
		return fmt.Sprintf("%s <%d chars>", scheme, len(cred))
	}
	return fmt.Sprintf("<%d chars>", len(value))
}

// doOnce executes a single HTTP request without retry.
func doOnce(req *Request) ([]byte, int, http.Header, error) {
	// Build URL
	url := req.baseURL + req.path
	for k, v := range req.pathParams {
		url = strings.ReplaceAll(url, "{"+k+"}", urlpkg.PathEscape(v))
	}

	// Add query params
	if len(req.queryParams) > 0 {
		params := urlpkg.Values{}
		for k, v := range req.queryParams {
			if v != "" {
				params.Set(k, v)
			}
		}
		if encoded := params.Encode(); encoded != "" {
			url += "?" + encoded
		}
	}

	// Build body
	var bodyReader io.Reader
	if req.bodyRaw != "" {
		bodyReader = strings.NewReader(req.bodyRaw)
	}

	httpReq, err := http.NewRequest(req.method, url, bodyReader)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("building request: %w", err)
	}

	// Auth
	httpReq.Header.Set("Authorization", "Bearer "+config.Token())

	// Content type
	if req.bodyRaw != "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}
	httpReq.Header.Set("Accept", "application/json")

	// Extra headers
	for k, v := range req.headers {
		httpReq.Header.Set(k, v)
	}

	// Debug
	if config.Debug() {
		fmt.Fprintf(os.Stderr, "DEBUG: %s %s\n", req.method, url)
		for k, v := range httpReq.Header {
			fmt.Fprintf(os.Stderr, "DEBUG:   %s: %s\n", k, redactHeader(k, strings.Join(v, ", ")))
		}
		if req.bodyRaw != "" {
			fmt.Fprintf(os.Stderr, "DEBUG:   Body: %s\n", truncate(req.bodyRaw, 500))
		}
	}

	// Dry-run: intercept write operations before making the HTTP call
	if config.DryRun() && isWriteMethod(req.method) {
		fmt.Fprintf(os.Stderr, "[DRY RUN] %s %s\n", req.method, url)
		for k, v := range httpReq.Header {
			fmt.Fprintf(os.Stderr, "[DRY RUN]   %s: %s\n", k, redactHeader(k, strings.Join(v, ", ")))
		}
		if req.bodyRaw != "" {
			fmt.Fprintf(os.Stderr, "[DRY RUN]   Body: %s\n", req.bodyRaw)
		}
		return nil, 0, nil, ErrDryRun
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, resp.Header, fmt.Errorf("reading response: %w", err)
	}

	if config.Debug() {
		fmt.Fprintf(os.Stderr, "DEBUG: Response %d (%d bytes)\n", resp.StatusCode, len(body))
	}

	if resp.StatusCode >= 400 {
		return body, resp.StatusCode, resp.Header, fmt.Errorf("API error %d: %s", resp.StatusCode, truncate(string(body), 500))
	}

	return body, resp.StatusCode, resp.Header, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// isWriteMethod returns true for HTTP methods that modify data.
func isWriteMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return true
	}
	return false
}
