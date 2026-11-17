package webexapi

import (
	"fmt"
	"regexp"
)

// APIError is a response with a status of 400 or more. Its message repeats up to
// 500 bytes of the body, after the request's Redact.
type APIError struct {
	Status int
	Body   []byte

	// shown is the body as Redact returned it, when the request had one.
	shown    []byte
	redacted bool
}

func (e *APIError) Error() string {
	shown := e.Body
	if e.redacted {
		shown = e.shown
	}
	return fmt.Sprintf("API error %d: %s", e.Status, truncate(string(shown), 500))
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
// region. Endpoint is the regional endpoint Webex names, when it names one.
type WrongRegionError struct {
	Endpoint string
	Body     []byte
}

func (e *WrongRegionError) Error() string {
	msg := "wrong data region (451): this organization's data is held in another region"
	if e.Endpoint != "" {
		msg += "; Webex says to use " + e.Endpoint
	}
	return msg
}

// endpointURL finds the regional endpoint Webex returns in a 451 body.
var endpointURL = regexp.MustCompile(`https?://[a-zA-Z0-9.\-]+`)

func newWrongRegionError(body []byte) *WrongRegionError {
	return &WrongRegionError{Endpoint: string(endpointURL.Find(body)), Body: body}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
