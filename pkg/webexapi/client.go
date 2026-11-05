// Package webexapi is the HTTP transport shared by the webex CLI and other programs
// that call Webex APIs. A Client holds no process-global state, so one program can
// run several at once, and every call takes a context.
package webexapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// defaultHTTP bounds each request. Callers that stream large downloads supply their
// own http.Client.
var defaultHTTP = &http.Client{Timeout: 60 * time.Second}

const defaultRetryAfter = 5 * time.Second

// TokenSource returns the bearer token. It is called for every attempt, so a token
// changed by Refresh takes effect on the retry.
type TokenSource func(ctx context.Context) (string, error)

// Client sends requests. The zero value makes no retries and has no token hooks.
type Client struct {
	// HTTP defaults to a client with a timeout.
	HTTP *http.Client

	Token TokenSource

	// Refresh is called at most once per Send, after a 401.
	Refresh func(ctx context.Context) error

	// MaxRetries limits 429 retries. MaxRetryWait, when positive, limits their
	// total wait.
	MaxRetries   int
	MaxRetryWait time.Duration

	// BeforeSend sees each attempt's request. An error it returns is passed back
	// unchanged and nothing is sent.
	BeforeSend func(*http.Request) error

	// AfterResponse is called with each attempt's status and body size.
	AfterResponse func(status, size int)

	// OnRetry is called before waiting out a 429.
	OnRetry func(wait time.Duration, attempt, max int)

	// Sleep defaults to a wait that ends early when ctx is cancelled.
	Sleep func(ctx context.Context, d time.Duration) error
}

// Request is one API call. URL is complete, including any query string.
type Request struct {
	Method string
	URL    string
	Body   []byte

	// ContentType applies when Body is set and defaults to application/json.
	ContentType string

	// Header entries override the defaults.
	Header http.Header
}

// Response is the final attempt's result.
type Response struct {
	Status int
	Header http.Header
	Body   []byte
}

// Send performs the request. A 401 refreshes the token and retries once, and a 429
// waits out Retry-After and retries within the limits. Nothing else is retried: a 5xx
// or a broken connection may follow a committed write, so repeating it is unsafe.
//
// A nil Response means no status was received. A non-nil Response with an error
// carries the failed status and body.
func (c *Client) Send(ctx context.Context, r Request) (*Response, error) {
	resp, err := c.attempt(ctx, r)
	if resp != nil && resp.Status == http.StatusUnauthorized && c.Refresh != nil {
		if refreshErr := c.Refresh(ctx); refreshErr != nil {
			return resp, err
		}
		resp, err = c.attempt(ctx, r)
	}

	var waited time.Duration
	for retries := 0; resp != nil && resp.Status == http.StatusTooManyRequests; retries++ {
		wait := RetryAfter(resp.Header.Get("Retry-After"))
		if retries >= c.MaxRetries {
			return resp, rateLimited(resp, fmt.Sprintf("retry limit reached after %d attempts — try again after %s",
				retries, wait.Round(time.Second)))
		}
		if c.MaxRetryWait > 0 && waited+wait > c.MaxRetryWait {
			return resp, rateLimited(resp, fmt.Sprintf("retry wait of %s would exceed max-retry-timer of %s — try again after %s",
				(waited+wait).Round(time.Second), c.MaxRetryWait.Round(time.Second), wait.Round(time.Second)))
		}
		if c.OnRetry != nil {
			c.OnRetry(wait, retries+1, c.MaxRetries)
		}
		if sleepErr := c.sleep(ctx, wait); sleepErr != nil {
			return resp, sleepErr
		}
		waited += wait
		resp, err = c.attempt(ctx, r)
	}
	return resp, err
}

func (c *Client) attempt(ctx context.Context, r Request) (*Response, error) {
	var body io.Reader
	if len(r.Body) > 0 {
		body = bytes.NewReader(r.Body)
	}
	httpReq, err := http.NewRequestWithContext(ctx, r.Method, r.URL, body)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	token := ""
	if c.Token != nil {
		if token, err = c.Token(ctx); err != nil {
			return nil, fmt.Errorf("getting token: %w", err)
		}
	}
	httpReq.Header.Set("Authorization", "Bearer "+token)
	if len(r.Body) > 0 {
		contentType := r.ContentType
		if contentType == "" {
			contentType = "application/json"
		}
		httpReq.Header.Set("Content-Type", contentType)
	}
	httpReq.Header.Set("Accept", "application/json")
	for k, values := range r.Header {
		httpReq.Header[http.CanonicalHeaderKey(k)] = values
	}

	if c.BeforeSend != nil {
		if err := c.BeforeSend(httpReq); err != nil {
			return nil, err
		}
	}

	client := c.HTTP
	if client == nil {
		client = defaultHTTP
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return &Response{Status: httpResp.StatusCode, Header: httpResp.Header}, fmt.Errorf("reading response: %w", err)
	}
	resp := &Response{Status: httpResp.StatusCode, Header: httpResp.Header, Body: data}
	if c.AfterResponse != nil {
		c.AfterResponse(resp.Status, len(data))
	}
	if resp.Status >= 400 {
		return resp, newAPIError(resp)
	}
	return resp, nil
}

func (c *Client) sleep(ctx context.Context, d time.Duration) error {
	if c.Sleep != nil {
		return c.Sleep(ctx, d)
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// RetryAfter parses a Retry-After value given in seconds. An absent or unusable
// value means five seconds.
func RetryAfter(header string) time.Duration {
	if secs, err := strconv.Atoi(strings.TrimSpace(header)); err == nil && secs > 0 {
		return time.Duration(secs) * time.Second
	}
	return defaultRetryAfter
}

// NextLink returns the rel="next" URL from a Link header, but only when it points at
// the same host as baseURL, so the bearer token is never sent elsewhere.
func NextLink(h http.Header, baseURL string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	for _, value := range h.Values("Link") {
		for _, part := range strings.Split(value, ",") {
			target, params, ok := strings.Cut(part, ";")
			if !ok || !strings.Contains(strings.ReplaceAll(params, " ", ""), `rel="next"`) {
				continue
			}
			target = strings.Trim(strings.TrimSpace(target), "<>")
			u, err := url.Parse(target)
			if err != nil || u.Scheme != base.Scheme || u.Host != base.Host {
				return ""
			}
			return target
		}
	}
	return ""
}
