package client

import (
	"encoding/json"
	"fmt"
	"os"
)

// Request builds an HTTP request for the Webex API.
type Request struct {
	method      string
	baseURL     string
	path        string
	pathParams  map[string]string
	queryParams map[string]string
	headers     map[string]string
	bodyJSON    map[string]interface{}
	bodyRaw     string
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
		bodyJSON:    make(map[string]interface{}),
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
}

// SetBodyFile reads a file and sets the raw body from its contents.
func (r *Request) SetBodyFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading body file: %w", err)
	}
	r.bodyRaw = string(data)
	return nil
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

// DoPaginated fetches all pages and returns the merged items array.
// isCalling selects the pagination strategy (Calling: start/max, CC: page/pageSize).
func (r *Request) DoPaginated(isCalling bool) ([]byte, int, error) {
	if isCalling {
		items, err := PaginateCalling(r.baseURL, r.method, r.path, r.pathParams, r.queryParams, r.headers)
		if err != nil {
			return nil, 0, err
		}
		data, err := json.Marshal(items)
		return data, 200, err
	}
	items, err := PaginateCC(r.baseURL, r.method, r.path, r.pathParams, r.queryParams, r.headers)
	if err != nil {
		return nil, 0, err
	}
	data, err := json.Marshal(items)
	return data, 200, err
}
