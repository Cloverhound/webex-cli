package webexapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	ccPageSize     = 100
	offsetPageSize = 100
)

// ListCC reads a Contact Center list to its end, passing each page to visit.
//
// CC endpoints wrap items in {data, meta}, or return a bare array with no page
// metadata. Neither form proves the list ended, so a full page always leads to
// another request. An endpoint that repeats a full page, or a later page that is not
// a list, gives ErrIncompleteList: a silently short list would look like deletions.
func (c *Client) ListCC(ctx context.Context, listURL string, visit func([]json.RawMessage) error) error {
	var previousFirst json.RawMessage
	for page := 0; ; page++ {
		pageURL, err := withQuery(listURL, map[string]string{"page": strconv.Itoa(page), "pageSize": strconv.Itoa(ccPageSize)})
		if err != nil {
			return err
		}
		resp, err := c.Send(ctx, Request{Method: http.MethodGet, URL: pageURL})
		if err != nil {
			return err
		}
		items, totalPages, err := ccPage(resp.Body)
		if err != nil {
			return fmt.Errorf("%w: page %d of %s: %v", ErrIncompleteList, page, listURL, err)
		}
		if len(items) == ccPageSize && bytes.Equal(items[0], previousFirst) {
			return fmt.Errorf("%w: %s repeated page %d", ErrIncompleteList, listURL, page)
		}
		if len(items) > 0 {
			if err := visit(items); err != nil {
				return err
			}
			previousFirst = items[0]
		}
		switch {
		case totalPages > 0 && page+1 >= totalPages:
			return nil
		case totalPages == 0 && len(items) < ccPageSize:
			return nil
		}
	}
}

func ccPage(body []byte) (items []json.RawMessage, totalPages int, err error) {
	if err := json.Unmarshal(body, &items); err == nil {
		return items, 0, nil
	}
	var result struct {
		Data     *[]json.RawMessage `json:"data"`
		Meta     *ccMeta            `json:"meta"`
		PageInfo *ccMeta            `json:"pageInfo"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, err
	}
	if result.Data == nil {
		return nil, 0, fmt.Errorf("response has no data array")
	}
	meta := result.Meta
	if meta == nil {
		meta = result.PageInfo
	}
	if meta != nil {
		totalPages = meta.TotalPages
	}
	return *result.Data, totalPages, nil
}

type ccMeta struct {
	TotalPages int `json:"totalPages"`
}

// ListOffset reads a start/max list to its end, passing each page to visit. arrayKey
// names the array in the response, which differs per endpoint. A Link rel="next"
// header is followed when present. See ListCC for why failures are not silent.
func (c *Client) ListOffset(ctx context.Context, listURL, arrayKey string, visit func([]json.RawMessage) error) error {
	var previousFirst json.RawMessage
	next := ""
	for start := 0; ; {
		pageURL := next
		if pageURL == "" {
			var err error
			pageURL, err = withQuery(listURL, map[string]string{"start": strconv.Itoa(start), "max": strconv.Itoa(offsetPageSize)})
			if err != nil {
				return err
			}
		}
		resp, err := c.Send(ctx, Request{Method: http.MethodGet, URL: pageURL})
		if err != nil {
			return err
		}
		// Several Calling lists answer {} when they hold nothing at all, which only
		// means that on the first page.
		if start == 0 && next == "" && bytes.Equal(bytes.TrimSpace(resp.Body), []byte("{}")) {
			return nil
		}
		items, err := namedArray(resp.Body, arrayKey)
		if err != nil {
			return fmt.Errorf("%w: %s: %v", ErrIncompleteList, listURL, err)
		}
		if len(items) == offsetPageSize && bytes.Equal(items[0], previousFirst) {
			return fmt.Errorf("%w: %s repeated a page", ErrIncompleteList, listURL)
		}
		if len(items) > 0 {
			if err := visit(items); err != nil {
				return err
			}
			previousFirst = items[0]
		}
		if link := NextLink(resp.Header, listURL); link != "" {
			next = link
			continue
		}
		if next != "" || len(items) < offsetPageSize {
			return nil
		}
		start += len(items)
	}
}

func namedArray(body []byte, key string) ([]json.RawMessage, error) {
	var result map[string]json.RawMessage
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	raw, ok := result[key]
	if !ok {
		return nil, fmt.Errorf("response has no %q array", key)
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("%q is not an array: %w", key, err)
	}
	return items, nil
}

func withQuery(rawURL string, params map[string]string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parsing list URL: %w", err)
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
