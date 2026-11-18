// Package webexapi sends requests to the Webex APIs: the bearer token, one refresh
// after a 401, waits for a 429's Retry-After, and errors that carry the failed
// response. It also reads lists to their end, names each data region's hosts, and
// converts organization IDs.
//
// A Client holds no process-wide state, so one program can run several at once,
// and every call takes a context.
package webexapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client sends requests. The zero value sends with http.DefaultClient, makes no
// 429 retries, and sends an empty bearer token.
//
// Redirects are never followed: a 3xx returns a *RedirectError. Following one
// would let net/http carry the token and other credentials to another host, and
// turn a write into a GET.
type Client struct {
	// HTTP defaults to http.DefaultClient. Its CheckRedirect is not used.
	HTTP *http.Client

	// Token returns the bearer token. It is called before every attempt, so a
	// token changed by Refresh is used on the retry.
	Token func(ctx context.Context) (string, error)

	// Refresh renews the token after a 401. The request is then sent once more.
	Refresh func(ctx context.Context) error

	// MaxRetries limits the retries after a 429. MaxRetryWait, when positive,
	// limits the total time spent waiting for them.
	MaxRetries   int
	MaxRetryWait time.Duration

	// BeforeSend sees each attempt's request, headers included. An error it
	// returns is returned unchanged, and nothing is sent.
	BeforeSend func(*http.Request) error

	// AfterResponse is called with each attempt's status and body size.
	AfterResponse func(status, size int)

	// OnRetry is called before waiting out a 429.
	OnRetry func(wait time.Duration, attempt, max int)

	// Sleep waits between 429 retries. It defaults to a wait that ends early when
	// ctx is done.
	Sleep func(ctx context.Context, d time.Duration) error
}

// Request is one API call.
type Request struct {
	Method string

	// URL is complete, query string included.
	URL string

	Body []byte

	// ContentType is sent when Body is not empty, and defaults to
	// application/json.
	ContentType string

	// Header entries are set after the defaults, so they override them.
	Header http.Header

	// Redact removes credentials from a failed response's body before its error
	// repeats it. The Response keeps the body as received.
	Redact func([]byte) []byte
}

// Response is the result of a request's last attempt.
type Response struct {
	Status int
	Header http.Header
	Body   []byte
}

// Send sends a request. After a 401 it calls Refresh and sends once more; after
// a 429 it waits for Retry-After and sends again, within MaxRetries and
// MaxRetryWait. Nothing else is retried: a 5xx or a broken connection may follow
// a write that was carried out.
//
// A status of 400 or more returns the Response with an *APIError, or a
// *WrongRegionError for a 451. A nil Response means no response was received.
func (c *Client) Send(ctx context.Context, r Request) (*Response, error) {
	resp, err := c.attempt(ctx, r)
	if resp != nil && resp.Status == http.StatusUnauthorized && c.Refresh != nil {
		if refreshErr := c.Refresh(ctx); refreshErr != nil {
			if ctx.Err() != nil {
				return resp, ctx.Err()
			}
			return resp, err
		}
		resp, err = c.attempt(ctx, r)
	}

	var waited time.Duration
	for retries := 0; resp != nil && resp.Status == http.StatusTooManyRequests; retries++ {
		wait := RetryAfter(resp.Header.Get("Retry-After"))
		limit := ""
		if retries >= c.MaxRetries {
			limit = fmt.Sprintf("retry limit reached after %d attempts — try again after %s", retries, wait.Round(time.Second))
		} else if c.MaxRetryWait > 0 && waited+wait > c.MaxRetryWait {
			limit = fmt.Sprintf("retry wait of %s would exceed max-retry-timer of %s — try again after %s",
				(waited + wait).Round(time.Second), c.MaxRetryWait.Round(time.Second), wait.Round(time.Second))
		}
		if limit != "" {
			apiErr, ok := err.(*APIError)
			if !ok {
				// The 429's body could not be read; that error says why.
				return resp, err
			}
			return resp, &RateLimitError{err: apiErr, msg: "rate limited (429): " + limit}
		}
		if c.OnRetry != nil {
			c.OnRetry(wait, retries+1, c.MaxRetries)
		}
		if err := c.sleep(ctx, wait); err != nil {
			return resp, err
		}
		waited += wait
		resp, err = c.attempt(ctx, r)
	}

	if resp != nil && resp.Status == http.StatusUnavailableForLegalReasons {
		apiErr, ok := err.(*APIError)
		if !ok {
			return resp, err
		}
		return resp, newWrongRegionError(apiErr)
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
		for i, v := range values {
			if i == 0 {
				httpReq.Header.Set(k, v)
			} else {
				httpReq.Header.Add(k, v)
			}
		}
	}

	if c.BeforeSend != nil {
		if err := c.BeforeSend(httpReq); err != nil {
			return nil, err
		}
	}

	httpResp, err := c.httpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer httpResp.Body.Close()

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return &Response{Status: httpResp.StatusCode, Header: httpResp.Header}, fmt.Errorf("reading response: %w", err)
	}
	if c.AfterResponse != nil {
		c.AfterResponse(httpResp.StatusCode, len(data))
	}
	resp := &Response{Status: httpResp.StatusCode, Header: httpResp.Header, Body: data}
	if resp.Status >= 300 {
		apiErr := &APIError{Status: resp.Status, Body: data}
		if r.Redact != nil {
			// A copy, so a redactor that edits its argument leaves the response as received.
			apiErr.Body = r.Redact(bytes.Clone(data))
		}
		if resp.Status < 400 {
			return resp, &RedirectError{Location: resp.Header.Get("Location"), err: apiErr}
		}
		return resp, apiErr
	}
	return resp, nil
}

// httpClient is the HTTP client set never to follow a redirect.
func (c *Client) httpClient() *http.Client {
	base := c.HTTP
	if base == nil {
		base = http.DefaultClient
	}
	client := *base
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &client
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

// RetryAfter reads a Retry-After value given in seconds. An absent or unusable
// value means five seconds.
func RetryAfter(header string) time.Duration {
	if secs, err := strconv.Atoi(strings.TrimSpace(header)); err == nil && secs > 0 {
		return time.Duration(secs) * time.Second
	}
	return 5 * time.Second
}
