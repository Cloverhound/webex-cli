package webexapi

import (
	"fmt"
	"regexp"
)

// APIError is a response with a status of 400 or more. Body is the response body
// after the request's Redact, and the message repeats up to 500 bytes of it; the
// Response keeps the body as received.
type APIError struct {
	Status int
	Body   []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.Status, truncate(string(e.Body), 500))
}

// RateLimitError is a 429 that was not retried further, because of MaxRetries or
// MaxRetryWait. It unwraps to the response's *APIError.
type RateLimitError struct {
	msg string
	err *APIError
}

func (e *RateLimitError) Error() string { return e.msg }
func (e *RateLimitError) Unwrap() error {
	if e.err == nil {
		return nil
	}
	return e.err
}

// WrongRegionError is a 451: the organization's data is held in another data
// region. Endpoint is the regional endpoint Webex names, when it names one, and
// Body is the response body after the request's Redact. It unwraps to the
// *APIError.
type WrongRegionError struct {
	Endpoint string
	Body     []byte

	err *APIError
}

func (e *WrongRegionError) Unwrap() error { return e.err }

// RedirectError is a 3xx. Redirects are never followed, so the request did not
// reach what it asked for. It unwraps to the *APIError.
type RedirectError struct {
	// Location is where the response points, as Webex gave it.
	Location string

	err *APIError
}

func (e *RedirectError) Error() string {
	return fmt.Sprintf("redirect (%d) to %s not followed", e.err.Status, e.Location)
}

func (e *RedirectError) Unwrap() error { return e.err }

func (e *WrongRegionError) Error() string {
	msg := "wrong data region (451): this organization's data is held in another region"
	if e.Endpoint != "" {
		msg += "; Webex says to use " + e.Endpoint
	}
	return msg
}

// endpointURL finds the regional endpoint Webex returns in a 451 body.
var endpointURL = regexp.MustCompile(`https?://[a-zA-Z0-9.\-]+`)

// newWrongRegionError finds the endpoint in the error's body, which is already
// redacted.
func newWrongRegionError(err *APIError) *WrongRegionError {
	return &WrongRegionError{Endpoint: string(endpointURL.Find(err.Body)), Body: err.Body, err: err}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
