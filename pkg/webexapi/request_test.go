package webexapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// testRequest is a list request as the paging tests build it.
type testRequest struct {
	base, method, path string
	query              url.Values
}

func NewRequest(base, method, path string) *testRequest {
	return &testRequest{base: base, method: method, path: path, query: url.Values{}}
}

func (r *testRequest) QueryParam(k, v string) { r.query.Set(k, v) }

func (r *testRequest) list() List {
	return List{Method: r.method, URL: r.base + r.path, Query: r.query, Name: r.path}
}

func (r *testRequest) DoPaginated(p Paging) ([]byte, int, error) {
	items, err := (&Client{}).All(context.Background(), r.list(), p)
	if err != nil {
		return nil, 0, err
	}
	data, err := json.Marshal(items)
	return data, http.StatusOK, err
}

func (r *testRequest) Visit(p Paging, visit func([]json.RawMessage) bool) error {
	return (&Client{}).Visit(context.Background(), r.list(), p, visit)
}
