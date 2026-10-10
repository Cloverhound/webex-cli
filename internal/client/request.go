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
	bodyJSON    map[string]any
	bodyRaw     string
	// pageSizeParam is the query parameter Contact Center auto-pagination
	// sets the page size with; empty means "pageSize".
	pageSizeParam string
	// offset paging for every other area; see OffsetPaging.
	offsetSet                     bool
	offsetStart, offsetSize       string
	offsetFirst                   int
	// rawURL, when set, is sent verbatim in place of path and query params.
	rawURL string
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

// PageSizeParam names the page-size query parameter for Contact Center
// endpoints that do not use "pageSize" (the flow store and functions use "size").
func (r *Request) PageSizeParam(name string) {
	r.pageSizeParam = name
}

// OffsetPaging names the offset and page-size query parameters for non-CC
// auto-pagination, which otherwise uses start/max. An empty startParam means
// the endpoint pages only through Link headers; firstIndex is the offset of
// the first item (1 for SCIM).
func (r *Request) OffsetPaging(startParam, sizeParam string, firstIndex int) {
	r.offsetSet = true
	r.offsetStart, r.offsetSize, r.offsetFirst = startParam, sizeParam, firstIndex
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
		start, size, first := "start", "max", 0
		if r.offsetSet {
			start, size, first = r.offsetStart, r.offsetSize, r.offsetFirst
		}
		items, err := paginateOffset(r.baseURL, r.method, r.path, start, size, first, r.pathParams, r.queryParams, r.headers)
		if err != nil {
			return nil, 0, err
		}
		data, err := json.Marshal(items)
		return data, 200, err
	}
	sizeParam := r.pageSizeParam
	if sizeParam == "" {
		sizeParam = "pageSize"
	}
	items, single, err := paginateCC(r.baseURL, r.method, r.path, sizeParam, r.pathParams, r.queryParams, r.headers)
	if err != nil {
		return nil, 0, err
	}
	if single != nil {
		return single, 200, nil
	}
	data, err := json.Marshal(items)
	return data, 200, err
}
