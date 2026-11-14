package webexapi

import (
	"errors"
	"fmt"
)

const errorBodyLimit = 500

// ErrIncompleteList reports a list that could not be read to its end. Callers must
// not treat the items seen so far as the whole list.
var ErrIncompleteList = errors.New("incomplete list")

// APIError is a response with a status of 400 or above.
type APIError struct {
	Status     int
	TrackingID string
	Body       []byte

	message string
	// fromBody marks a message built from Body, which WithBody rebuilds.
	fromBody bool
}

func (e *APIError) Error() string {
	return e.message
}

// WithBody returns a copy of the error that carries body instead, with its message
// rebuilt from it, as when a caller must redact what the response held.
func (e *APIError) WithBody(body []byte) *APIError {
	c := *e
	c.Body = body
	if c.fromBody {
		c.message = bodyMessage(c.Status, body)
	}
	return &c
}

func newAPIError(resp *Response) *APIError {
	return &APIError{
		Status:     resp.Status,
		TrackingID: resp.Header.Get("Trackingid"),
		Body:       resp.Body,
		message:    bodyMessage(resp.Status, resp.Body),
		fromBody:   true,
	}
}

func bodyMessage(status int, body []byte) string {
	return fmt.Sprintf("API error %d: %s", status, truncate(string(body), errorBodyLimit))
}

func rateLimited(resp *Response, detail string) *APIError {
	e := newAPIError(resp)
	e.message, e.fromBody = "rate limited (429): "+detail, false
	return e
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
