package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/config"
)

// PaginateCalling auto-paginates a Calling API list endpoint.
// Calling uses startIndex/itemsPerPage with items array.
func PaginateCalling(baseURL, method, path string, pathParams, queryParams, headers map[string]string) ([]json.RawMessage, error) {
	return paginateOffset(baseURL, method, path, "start", "max", 0, pathParams, queryParams, headers)
}

// paginateOffset pages a non-CC list endpoint. A Link rel="next" header is
// followed whenever present, since most core Webex APIs page only that way and
// ignore offsets. Without one, it advances startParam by the items received,
// or stops when startParam is empty.
func paginateOffset(baseURL, method, path, startParam, sizeParam string, firstIndex int, pathParams, queryParams, headers map[string]string) ([]json.RawMessage, error) {
	var allItems []json.RawMessage
	var prevFirst json.RawMessage
	start := firstIndex
	pageSize := 100
	next := ""

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
		if next != "" {
			r.rawURL = next
		} else {
			if startParam != "" {
				r.QueryParam(startParam, strconv.Itoa(start))
			}
			if sizeParam != "" {
				r.QueryParam(sizeParam, strconv.Itoa(pageSize))
			}
		}

		r.Finalize()
		body, _, respHeaders, err := doWithHeaders(r)
		if err != nil {
			return allItems, err
		}

		items, total, err := listItems(body)
		if err != nil {
			return allItems, err
		}
		if len(items) == 0 {
			return allItems, nil
		}
		// An endpoint that ignores the offset returns the same page again.
		if prevFirst != nil && bytes.Equal(items[0], prevFirst) {
			return allItems, nil
		}
		prevFirst = items[0]
		allItems = append(allItems, items...)

		if link := nextLink(respHeaders, baseURL); link != "" {
			next = link
			continue
		}
		if next != "" || startParam == "" {
			return allItems, nil
		}
		start += len(items)
		if len(items) < pageSize || (total > 0 && len(allItems) >= total) {
			return allItems, nil
		}
	}
}

// listItems extracts a page's items: a bare array, or the "items" (Webex) or
// "Resources" (SCIM) array, or else the first array field by name. total is
// SCIM's totalResults when present.
func listItems(body []byte) ([]json.RawMessage, int, error) {
	var bare []json.RawMessage
	if err := json.Unmarshal(body, &bare); err == nil {
		return bare, 0, nil
	}
	var result map[string]json.RawMessage
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, fmt.Errorf("parsing page: %w", err)
	}
	var total int
	if raw, ok := result["totalResults"]; ok {
		json.Unmarshal(raw, &total)
	}

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	keys = append([]string{"items", "Resources"}, keys...)
	for _, k := range keys {
		var items []json.RawMessage
		if v, ok := result[k]; ok && json.Unmarshal(v, &items) == nil && len(items) > 0 {
			return items, total, nil
		}
	}
	return nil, total, nil
}

// nextLink returns the rel="next" URL from a Link header, but only when it
// points at the same host as baseURL, so the bearer token is never sent elsewhere.
func nextLink(h http.Header, baseURL string) string {
	base, err := urlpkg.Parse(baseURL)
	if err != nil {
		return ""
	}
	for _, value := range h.Values("Link") {
		for _, part := range strings.Split(value, ",") {
			target, params, ok := strings.Cut(part, ";")
			if !ok || !strings.Contains(strings.ReplaceAll(params, " ", ""), `rel="next"`) {
				continue
			}
			target = strings.Trim(strings.TrimSpace(target), "<>")
			u, err := urlpkg.Parse(target)
			if err != nil || u.Scheme != base.Scheme || u.Host != base.Host {
				return ""
			}
			return target
		}
	}
	return ""
}

// PaginateCC auto-paginates a Contact Center API list endpoint.
// CC uses page/pageSize with meta.totalPages.
func PaginateCC(baseURL, method, path string, pathParams, queryParams, headers map[string]string) ([]json.RawMessage, error) {
	items, _, err := paginateCC(baseURL, method, path, "pageSize", pathParams, queryParams, headers)
	return items, err
}

// paginateCC pages with page and sizeParam. Most CC endpoints wrap items in
// {data, meta}; the flow store and functions use {data, pageInfo}, and the
// flow store returns a bare array unless includePagination=true. A first
// response that is neither is returned as single, to be printed unchanged.
func paginateCC(baseURL, method, path, sizeParam string, pathParams, queryParams, headers map[string]string) (items []json.RawMessage, single []byte, err error) {
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
		r.QueryParam(sizeParam, strconv.Itoa(pageSize))

		body, _, err := r.Do()
		if err != nil {
			return allItems, nil, err
		}

		// A bare array carries no page metadata, so it is the only page.
		var bare []json.RawMessage
		if err := json.Unmarshal(body, &bare); err == nil {
			return append(allItems, bare...), nil, nil
		}

		var result map[string]json.RawMessage
		if err := json.Unmarshal(body, &result); err != nil {
			return allItems, nil, fmt.Errorf("parsing page: %w", err)
		}

		var items []json.RawMessage
		data, hasData := result["data"]
		if hasData && json.Unmarshal(data, &items) != nil {
			hasData = false
		}
		if !hasData {
			if page == 0 {
				return nil, body, nil
			}
			return allItems, nil, nil
		}
		allItems = append(allItems, items...)

		metaRaw, ok := result["meta"]
		if !ok {
			metaRaw, ok = result["pageInfo"]
		}
		if !ok || len(items) == 0 {
			return allItems, nil, nil
		}
		var meta struct {
			TotalPages int `json:"totalPages"`
		}
		if err := json.Unmarshal(metaRaw, &meta); err != nil || page+1 >= meta.TotalPages {
			return allItems, nil, nil
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
