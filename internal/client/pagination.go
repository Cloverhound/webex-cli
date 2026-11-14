package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/config"
)

// ErrIncompleteList means --paginate could not read a list to its end. It fails
// rather than returning the pages it has, which would look like the whole list.
var ErrIncompleteList = errors.New("incomplete list")

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
	var allItems, prev []json.RawMessage
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
			if total > len(allItems) {
				return allItems, fmt.Errorf("%w: %s ended at %d of %d items", ErrIncompleteList, path, len(allItems), total)
			}
			return allItems, nil
		}
		if done, err := repeated(path, items, prev, pageSize); done || err != nil {
			return allItems, err
		}
		prev = items
		allItems = append(allItems, items...)

		link, refused := nextLink(respHeaders, baseURL)
		if refused {
			return allItems, fmt.Errorf("%w: %s gave its next page on another host", ErrIncompleteList, path)
		}
		if link != "" {
			next = link
			continue
		}
		if next != "" || startParam == "" {
			return allItems, nil
		}
		start += len(items)
		// A stated total (SCIM) ends the list only once reached, since a short page
		// can be the server's own page size.
		if total > 0 {
			if len(allItems) >= total {
				return allItems, nil
			}
			continue
		}
		if len(items) < pageSize {
			return allItems, nil
		}
	}
}

// repeated checks a page that starts like the previous one, from an endpoint that
// ignores the paging parameters. If the previous page was longer than the page
// asked for, the endpoint ignored the size as well and sent the whole list, so the
// list is done. Otherwise it may hold more than one page that cannot be reached.
func repeated(path string, items, prev []json.RawMessage, pageSize int) (done bool, err error) {
	if len(prev) == 0 || len(items) == 0 || !bytes.Equal(items[0], prev[0]) {
		return false, nil
	}
	if len(prev) > pageSize {
		return true, nil
	}
	return true, fmt.Errorf("%w: %s returned the same page again", ErrIncompleteList, path)
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
// refused reports a next page that is somewhere else.
func nextLink(h http.Header, baseURL string) (link string, refused bool) {
	base, err := urlpkg.Parse(baseURL)
	if err != nil {
		return "", false
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
				return "", true
			}
			return target, false
		}
	}
	return "", false
}

// PaginateCC auto-paginates a Contact Center API list endpoint.
// CC uses page/pageSize with meta.totalPages.
func PaginateCC(baseURL, method, path string, pathParams, queryParams, headers map[string]string) ([]json.RawMessage, error) {
	items, _, err := paginateCC(baseURL, method, path, "pageSize", pathParams, queryParams, headers)
	return items, err
}

// VisitCC pages a Contact Center list endpoint like PaginateCC, passing each
// page's items to visit and stopping as soon as visit returns false.
func VisitCC(baseURL, method, path string, pathParams, queryParams map[string]string, visit func(items []json.RawMessage) bool) error {
	_, err := visitCC(baseURL, method, path, "pageSize", pathParams, queryParams, nil, visit)
	return err
}

// paginateCC pages with page and sizeParam. A first response that is not a
// list is returned as single, to be printed unchanged.
func paginateCC(baseURL, method, path, sizeParam string, pathParams, queryParams, headers map[string]string) (items []json.RawMessage, single []byte, err error) {
	single, err = visitCC(baseURL, method, path, sizeParam, pathParams, queryParams, headers, func(page []json.RawMessage) bool {
		items = append(items, page...)
		return true
	})
	return items, single, err
}

// visitCC pages with page and sizeParam. Most CC endpoints wrap items in
// {data, meta}; the flow store and functions use {data, pageInfo}, and the
// flow store returns a bare array unless includePagination=true. A first
// response that is neither is returned as single and not visited. Without page
// metadata nothing shows the list ended, so a full page leads to another request.
func visitCC(baseURL, method, path, sizeParam string, pathParams, queryParams, headers map[string]string, visit func(items []json.RawMessage) bool) (single []byte, err error) {
	page := 0
	pageSize := 100
	var prev []json.RawMessage

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
			return nil, err
		}

		var items []json.RawMessage
		var metaRaw json.RawMessage
		hasMeta := false
		if err := json.Unmarshal(body, &items); err != nil {
			var result map[string]json.RawMessage
			if err := json.Unmarshal(body, &result); err != nil {
				return nil, fmt.Errorf("parsing page: %w", err)
			}
			data, hasData := result["data"]
			if hasData && json.Unmarshal(data, &items) != nil {
				hasData = false
			}
			if !hasData {
				if page == 0 {
					return body, nil
				}
				return nil, nil
			}
			metaRaw, hasMeta = result["meta"]
			if !hasMeta {
				metaRaw, hasMeta = result["pageInfo"]
			}
		}

		if done, err := repeated(path, items, prev, pageSize); done || err != nil {
			return nil, err
		}
		prev = items
		if !visit(items) || len(items) == 0 {
			return nil, nil
		}

		if hasMeta {
			var meta struct {
				TotalPages int `json:"totalPages"`
			}
			if err := json.Unmarshal(metaRaw, &meta); err != nil || page+1 >= meta.TotalPages {
				return nil, nil
			}
		} else if len(items) < pageSize {
			return nil, nil
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
