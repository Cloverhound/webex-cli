package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// Request builds an HTTP request for the Webex API.
type Request struct {
	method      string
	baseURL     string
	path        string
	pathParams  map[string]string
	queryParams map[string]string
	headers     map[string]string
	bodyJSON    map[string]any
	bodyRaw     string

	// Multipart bodies override the JSON content type and log a safe summary.
	bodyContentType, bodySummary string

	// redactError removes credentials from a failed response before its error
	// repeats the body.
	redactError func([]byte) []byte
}

// NewRequest creates a new request with the given method and path.
func NewRequest(baseURL, method, path string) *Request {
	return &Request{
		method:      method,
		baseURL:     baseURL,
		path:        path,
		pathParams:  make(map[string]string),
		queryParams: make(map[string]string),
		headers:     make(map[string]string),
		bodyJSON:    make(map[string]any),
	}
}

func (r *Request) PathParam(key, value string) {
	if value != "" {
		r.pathParams[key] = value
	}
}

func (r *Request) QueryParam(key, value string) {
	if value != "" {
		r.queryParams[key] = value
	}
}

func (r *Request) Header(key, value string) {
	if value != "" {
		r.headers[key] = value
	}
}

func (r *Request) BodyString(key, value string) {
	if value != "" {
		r.bodyJSON[key] = value
	}
}

func (r *Request) BodyInt(key string, value int64, set bool) {
	if set {
		r.bodyJSON[key] = value
	}
}

func (r *Request) BodyBool(key string, value bool, set bool) {
	if set {
		r.bodyJSON[key] = value
	}
}

func (r *Request) BodyFloat(key string, value float64, set bool) {
	if set {
		r.bodyJSON[key] = value
	}
}

func (r *Request) BodyStringSlice(key string, values []string) {
	if len(values) > 0 {
		r.bodyJSON[key] = values
	}
}

// SetBodyRaw sets the raw JSON body, overriding individual fields.
func (r *Request) SetBodyRaw(raw string) {
	r.bodyRaw = raw
	r.bodyContentType = ""
	r.bodySummary = ""
}

// SetBodyFile reads a file and sets the raw body from its contents.
func (r *Request) SetBodyFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading body file: %w", err)
	}
	r.SetBodyRaw(string(data))
	return nil
}

// RedactErrors sets how a failed response's body is redacted before the error
// message repeats it.
func (r *Request) RedactErrors(redact func([]byte) []byte) {
	r.redactError = redact
}

// Finalize builds the body JSON if no raw body is set.
func (r *Request) Finalize() {
	if r.bodyRaw == "" && len(r.bodyJSON) > 0 {
		data, err := json.Marshal(r.bodyJSON)
		if err == nil {
			r.bodyRaw = string(data)
		}
	}
}

// Do executes the request and returns the response.
func (r *Request) Do() ([]byte, int, error) {
	r.Finalize()
	return Do(r)
}

// url is the request URL with its path parameters filled in and no query.
func (r *Request) url() string {
	return webexapi.FillPath(r.baseURL+r.path, r.pathParams)
}

func (r *Request) query() url.Values {
	q := url.Values{}
	for k, v := range r.queryParams {
		q.Set(k, v)
	}
	return q
}

func (r *Request) header() http.Header {
	h := http.Header{}
	for k, v := range r.headers {
		h.Set(k, v)
	}
	return h
}

// logBody is how diagnostics show the body: a multipart body's summary, never
// its bytes.
func (r *Request) logBody() string {
	if r.bodySummary != "" {
		return r.bodySummary
	}
	return r.bodyRaw
}

func (r *Request) sdkRequest() webexapi.Request {
	return webexapi.Request{
		Method:      r.method,
		URL:         webexapi.WithQuery(r.url(), r.query()),
		Body:        []byte(r.bodyRaw),
		ContentType: r.bodyContentType,
		Header:      r.header(),
		Redact:      r.redactError,
	}
}

func (r *Request) sdkList() webexapi.List {
	return webexapi.List{
		Method:      r.method,
		URL:         r.url(),
		Query:       r.query(),
		Body:        []byte(r.bodyRaw),
		ContentType: r.bodyContentType,
		Header:      r.header(),
		Redact:      r.redactError,
		Name:        r.path,
	}
}
