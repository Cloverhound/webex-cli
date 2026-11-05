package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/readonly"
	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// ErrDryRun is returned when a write operation is intercepted by --dry-run mode.
var ErrDryRun = errors.New("dry run: no changes made")

// ErrReadOnly is returned when read-only mode blocks a request.
var ErrReadOnly = errors.New("read-only mode")

// ReadOnlyError reports a request blocked by read-only mode.
func ReadOnlyError(method, url string) error {
	return fmt.Errorf("%w: %s %s not sent — run 'webex login' to leave read-only mode", ErrReadOnly, method, url)
}

// Do executes an HTTP request.
// On a 401 it refreshes the token and retries once.
// On a 429 it reads Retry-After and retries up to 3 times.
func Do(req *Request) ([]byte, int, error) {
	body, status, _, err := doWithHeaders(req)
	return body, status, err
}

// doWithHeaders is Do, also returning the final response's headers.
func doWithHeaders(req *Request) ([]byte, int, http.Header, error) {
	// Capture runs before read-only checks: it never writes to Webex, so a
	// read-only login can propose changes.
	if path := config.ChangesetPath(); path != "" && !readOnlyAllows(req.method, req.path) {
		body, err := capture(req, path)
		if err != nil {
			return nil, 0, nil, err
		}
		return body, http.StatusOK, http.Header{}, nil
	}
	url := requestURL(req)
	api := apiClient(req, url)

	apiReq := webexapi.Request{
		Method:      req.method,
		URL:         url,
		ContentType: req.bodyContentType,
		Header:      http.Header{},
	}
	if req.bodyRaw != "" {
		apiReq.Body = []byte(req.bodyRaw)
	}
	for k, v := range req.headers {
		apiReq.Header.Set(k, v)
	}

	resp, err := api.Send(context.Background(), apiReq)
	if resp == nil {
		return nil, 0, nil, err
	}
	if resp.Status == 451 {
		return resp.Body, resp.Status, resp.Header, wrongRegionError(resp.Body)
	}
	return resp.Body, resp.Status, resp.Header, err
}

// apiClient configures the shared transport from the process-wide CLI settings.
func apiClient(req *Request, url string) *webexapi.Client {
	api := &webexapi.Client{
		HTTP: http.DefaultClient,
		Token: func(context.Context) (string, error) {
			return config.Token(), nil
		},
		MaxRetries:   config.MaxRetry(),
		MaxRetryWait: time.Duration(config.MaxRetryTimer()) * time.Second,
		OnRetry: func(wait time.Duration, attempt, max int) {
			fmt.Fprintf(os.Stderr, "Rate limited (429). Retrying in %s (attempt %d/%d)...\n", wait.Round(time.Second), attempt, max)
		},
		BeforeSend: func(httpReq *http.Request) error {
			return beforeSend(req, url, httpReq)
		},
		AfterResponse: func(status, size int) {
			if config.Debug() {
				fmt.Fprintf(os.Stderr, "DEBUG: Response %d (%d bytes)\n", status, size)
			}
		},
	}
	if config.TokenRefresher != nil {
		api.Refresh = func(context.Context) error {
			newToken, err := config.TokenRefresher()
			if err != nil {
				return err
			}
			config.SetToken(newToken)
			return nil
		}
	}
	return api
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

// redactHeader replaces a credential value with its length, which is enough to
// tell a missing or truncated token from a present one without printing it.
func redactHeader(key, value string) string {
	// An AI agent's runtime token lets anyone talk to that agent.
	if strings.EqualFold(key, "Bot-Access-Token") {
		return fmt.Sprintf("<%d chars>", len(value))
	}
	if !strings.EqualFold(key, "Authorization") {
		return value
	}
	if scheme, cred, ok := strings.Cut(value, " "); ok {
		return fmt.Sprintf("%s <%d chars>", scheme, len(cred))
	}
	return fmt.Sprintf("<%d chars>", len(value))
}

// requestURL expands path parameters and appends the query string.
func requestURL(req *Request) string {
	if req.rawURL != "" {
		return req.rawURL
	}
	url := req.baseURL + req.path
	for k, v := range req.pathParams {
		url = strings.ReplaceAll(url, "{"+k+"}", urlpkg.PathEscape(v))
	}
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
	return url
}

// beforeSend applies the CLI's debug output, read-only mode and dry-run mode to each
// attempt.
func beforeSend(req *Request, url string, httpReq *http.Request) error {
	// Multipart diagnostics describe the file; never print its binary contents.
	logBody := req.bodyRaw
	if req.bodySummary != "" {
		logBody = req.bodySummary
	}

	// Debug
	if config.Debug() {
		fmt.Fprintf(os.Stderr, "DEBUG: %s %s\n", req.method, url)
		for k, v := range httpReq.Header {
			fmt.Fprintf(os.Stderr, "DEBUG:   %s: %s\n", k, redactHeader(k, strings.Join(v, ", ")))
		}
		if req.bodyRaw != "" {
			fmt.Fprintf(os.Stderr, "DEBUG:   Body: %s\n", truncate(logBody, 500))
		}
	}

	if config.ReadOnly() && !readOnlyAllows(req.method, req.path) {
		return ReadOnlyError(req.method, url)
	}

	// Dry-run: intercept write operations before making the HTTP call
	if config.DryRun() && isWriteMethod(req.method) {
		fmt.Fprintf(os.Stderr, "[DRY RUN] %s %s\n", req.method, url)
		for k, v := range httpReq.Header {
			fmt.Fprintf(os.Stderr, "[DRY RUN]   %s: %s\n", k, redactHeader(k, strings.Join(v, ", ")))
		}
		if req.bodyRaw != "" {
			fmt.Fprintf(os.Stderr, "[DRY RUN]   Body: %s\n", logBody)
		}
		return ErrDryRun
	}

	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// readOnlyAllows reports whether a request may be sent in read-only mode.
// path is the unexpanded template, which is how query endpoints are listed.
func readOnlyAllows(method, path string) bool {
	if method == http.MethodPost {
		return readonly.IsQueryPOST(path)
	}
	return !isWriteMethod(method)
}

// isWriteMethod returns true for HTTP methods that modify data.
func isWriteMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return true
	}
	return false
}
