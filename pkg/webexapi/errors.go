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
}

func (e *APIError) Error() string {
	return e.message
}

func newAPIError(resp *Response) *APIError {
	return &APIError{
		Status:     resp.Status,
		TrackingID: resp.Header.Get("Trackingid"),
		Body:       resp.Body,
		message:    fmt.Sprintf("API error %d: %s", resp.Status, truncate(string(resp.Body), errorBodyLimit)),
	}
}

func rateLimited(resp *Response, detail string) *APIError {
	e := newAPIError(resp)
	e.message = "rate limited (429): " + detail
	return e
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
