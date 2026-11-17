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

	shown []byte
}

func (e *APIError) Error() string {
	shown := e.shown
	if shown == nil {
		shown = e.Body
	}
	return fmt.Sprintf("API error %d: %s", e.Status, truncate(string(shown), 500))
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
