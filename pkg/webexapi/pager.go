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

// Paging says how a list is read to its end.
type Paging struct {
	// CC lists page with page and SizeParam, pageSize when empty, and wrap items in
	// {data, meta} or {data, pageInfo}, or send a bare array.
	CC        bool
	SizeParam string

	// Start and Size name an offset list's parameters. An empty Start means the list
	// pages only through Link headers, and an empty Size sends no page size. First is
	// the first item's offset, 1 for SCIM.
	Start, Size string
	First       int

	// ArrayKey names the array in an offset list's response. When empty, a bare
	// array, "items" or "Resources" is read, or else the response's only array.
	ArrayKey string

	// Header is sent with every page request.
	Header http.Header
}

// List reads a list to its end, passing each page to visit. A list that cannot be
// read to its end gives ErrIncompleteList: a silently short list would look like
// deletions.
func (c *Client) List(ctx context.Context, listURL string, p Paging, visit func([]json.RawMessage) error) error {
	if p.CC {
		if p.SizeParam == "" {
			p.SizeParam = "pageSize"
		}
		return c.listCC(ctx, listURL, p, visit)
	}
	if p.Start != "" && p.Size == "" {
		return fmt.Errorf("offset paging by %s needs a page size parameter", p.Start)
	}
	return c.listOffset(ctx, listURL, p, visit)
}

// ListCC reads a Contact Center list to its end, passing each page to visit.
//
// CC endpoints wrap items in {data, meta}, or return a bare array with no page
// metadata. Neither form proves the list ended, so a full page always leads to
// another request. An endpoint that repeats a full page, or a later page that is not
// a list, gives ErrIncompleteList: a silently short list would look like deletions.
func (c *Client) ListCC(ctx context.Context, listURL string, visit func([]json.RawMessage) error) error {
	return c.listCC(ctx, listURL, Paging{CC: true, SizeParam: "pageSize"}, visit)
}

func (c *Client) listCC(ctx context.Context, listURL string, p Paging, visit func([]json.RawMessage) error) error {
	var previousFirst json.RawMessage
	for page := 0; ; page++ {
		pageURL, err := withQuery(listURL, map[string]string{"page": strconv.Itoa(page), p.SizeParam: strconv.Itoa(ccPageSize)})
		if err != nil {
			return err
		}
		resp, err := c.Send(ctx, Request{Method: http.MethodGet, URL: pageURL, Header: p.Header.Clone()})
		if err != nil {
			return err
		}
		items, totalPages, err := ccPage(resp.Body)
		if err != nil {
			return fmt.Errorf("%w: page %d of %s: %v", ErrIncompleteList, page, listURL, err)
		}
		if repeats(items, previousFirst) {
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

// array decodes a JSON array. null decodes into a nil slice without error, but
// reading it as no items would make a failed page look like the end of the list.
func array(raw []byte) ([]json.RawMessage, error) {
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil, fmt.Errorf("null, not an array")
	}
	var items []json.RawMessage
	err := json.Unmarshal(raw, &items)
	return items, err
}

func ccPage(body []byte) (items []json.RawMessage, totalPages int, err error) {
	if items, err := array(body); err == nil {
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
	return c.listOffset(ctx, listURL, Paging{Start: "start", Size: "max", ArrayKey: arrayKey}, visit)
}

func (c *Client) listOffset(ctx context.Context, listURL string, p Paging, visit func([]json.RawMessage) error) error {
	var previousFirst json.RawMessage
	next := ""
	seen := 0
	for start := p.First; ; {
		pageURL := next
		if pageURL == "" {
			params := map[string]string{}
			if p.Start != "" {
				params[p.Start] = strconv.Itoa(start)
			}
			if p.Size != "" {
				params[p.Size] = strconv.Itoa(offsetPageSize)
			}
			var err error
			if pageURL, err = withQuery(listURL, params); err != nil {
				return err
			}
		}
		resp, err := c.Send(ctx, Request{Method: http.MethodGet, URL: pageURL, Header: p.Header.Clone()})
		if err != nil {
			return err
		}
		// Several Calling lists answer {} when they hold nothing at all, which only
		// means that on the first page.
		if start == p.First && next == "" && bytes.Equal(bytes.TrimSpace(resp.Body), []byte("{}")) {
			return nil
		}
		items, total, err := pageArray(resp.Body, p.ArrayKey)
		if err != nil {
			return fmt.Errorf("%w: %s: %v", ErrIncompleteList, listURL, err)
		}
		if repeats(items, previousFirst) {
			return fmt.Errorf("%w: %s repeated a page", ErrIncompleteList, listURL)
		}
		if len(items) > 0 {
			if err := visit(items); err != nil {
				return err
			}
			previousFirst = items[0]
		}
		seen += len(items)
		if link := NextLink(resp.Header, listURL); link != "" {
			next = link
			continue
		}
		if next != "" || p.Start == "" {
			return nil
		}
		// A stated total ends the list only when reached; a short page can be a
		// server's own page size.
		switch {
		case total > 0 && seen >= total:
			return nil
		case total > 0 && len(items) == 0:
			return fmt.Errorf("%w: %s ended at %d of %d items", ErrIncompleteList, listURL, seen, total)
		case total == 0 && len(items) < offsetPageSize:
			return nil
		}
		start += len(items)
	}
}

// repeats reports whether a page starts with the previous page's first item: the
// endpoint ignored the paging parameters, whatever size of page it answers with.
func repeats(items []json.RawMessage, previousFirst json.RawMessage) bool {
	return len(items) > 0 && previousFirst != nil && bytes.Equal(items[0], previousFirst)
}

// pageArray finds a page's items. total is SCIM's totalResults when present.
func pageArray(body []byte, key string) (items []json.RawMessage, total int, err error) {
	if key == "" {
		if items, err := array(body); err == nil {
			return items, 0, nil
		}
	}
	var result map[string]json.RawMessage
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, 0, err
	}
	if raw, ok := result["totalResults"]; ok {
		json.Unmarshal(raw, &total)
	}
	if key != "" {
		items, err := namedArray(result, key)
		return items, total, err
	}
	for _, k := range []string{"items", "Resources"} {
		if _, ok := result[k]; ok {
			items, err := namedArray(result, k)
			return items, total, err
		}
	}
	var found []string
	for k, raw := range result {
		var probe []json.RawMessage
		if json.Unmarshal(raw, &probe) == nil && probe != nil {
			found = append(found, k)
			items = probe
		}
	}
	switch len(found) {
	case 1:
		return items, total, nil
	case 0:
		return nil, 0, fmt.Errorf("response has no array")
	default:
		return nil, 0, fmt.Errorf("response has several arrays; name one")
	}
}

func namedArray(result map[string]json.RawMessage, key string) ([]json.RawMessage, error) {
	raw, ok := result[key]
	if !ok {
		return nil, fmt.Errorf("response has no %q array", key)
	}
	items, err := array(raw)
	if err != nil {
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
