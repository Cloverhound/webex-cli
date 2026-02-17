package client

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Cloverhound/webex-cli/internal/config"
)

// PaginateCalling auto-paginates a Calling API list endpoint.
// Calling uses startIndex/itemsPerPage with items array.
func PaginateCalling(baseURL, method, path string, pathParams, queryParams, headers map[string]string) ([]json.RawMessage, error) {
	var allItems []json.RawMessage
	start := 0
	max := 100

	for {
		r := NewRequest(baseURL, method, path)
		for k, v := range pathParams {
			r.PathParam(k, v)
		}
		for k, v := range queryParams {
			r.QueryParam(k, v)
		}
		for k, v := range headers {
			r.Header(k, v)
		}
		r.QueryParam("start", strconv.Itoa(start))
		r.QueryParam("max", strconv.Itoa(max))

		body, _, err := r.Do()
		if err != nil {
			return allItems, err
		}

		var result map[string]json.RawMessage
		if err := json.Unmarshal(body, &result); err != nil {
			return allItems, fmt.Errorf("parsing page: %w", err)
		}

		// Find the items array (first array field)
		found := false
		for _, v := range result {
			var items []json.RawMessage
			if err := json.Unmarshal(v, &items); err == nil && len(items) > 0 {
				allItems = append(allItems, items...)
				found = true
				if len(items) < max {
					return allItems, nil
				}
				break
			}
		}

		if !found {
			return allItems, nil
		}

		start += max
	}
}

// PaginateCC auto-paginates a Contact Center API list endpoint.
// CC uses page/pageSize with meta.totalPages.
func PaginateCC(baseURL, method, path string, pathParams, queryParams, headers map[string]string) ([]json.RawMessage, error) {
	var allItems []json.RawMessage
	page := 0
	pageSize := 100

	for {
		r := NewRequest(baseURL, method, path)
		for k, v := range pathParams {
			r.PathParam(k, v)
		}
		for k, v := range queryParams {
			r.QueryParam(k, v)
		}
		for k, v := range headers {
			r.Header(k, v)
		}
		r.QueryParam("page", strconv.Itoa(page))
		r.QueryParam("pageSize", strconv.Itoa(pageSize))

		body, _, err := r.Do()
		if err != nil {
			return allItems, err
		}

		var result map[string]json.RawMessage
		if err := json.Unmarshal(body, &result); err != nil {
			return allItems, fmt.Errorf("parsing page: %w", err)
		}

		// Find data array
		if data, ok := result["data"]; ok {
			var items []json.RawMessage
			if err := json.Unmarshal(data, &items); err == nil {
				allItems = append(allItems, items...)
			}
		}

		// Check meta for total pages
		if metaRaw, ok := result["meta"]; ok {
			var meta struct {
				TotalPages int `json:"totalPages"`
			}
			if err := json.Unmarshal(metaRaw, &meta); err == nil {
				if page+1 >= meta.TotalPages {
					return allItems, nil
				}
			}
		} else {
			return allItems, nil
		}

		page++
	}
}

// AutoPaginate is a convenience wrapper that selects the right pagination strategy.
func AutoPaginate(isCalling bool, baseURL, method, path string, pathParams, queryParams, headers map[string]string) ([]byte, error) {
	if !config.Paginate() {
		return nil, fmt.Errorf("pagination not enabled")
	}

	var items []json.RawMessage
	var err error
	if isCalling {
		items, err = PaginateCalling(baseURL, method, path, pathParams, queryParams, headers)
	} else {
		items, err = PaginateCC(baseURL, method, path, pathParams, queryParams, headers)
	}
	if err != nil {
		return nil, err
	}

	return json.Marshal(items)
}
