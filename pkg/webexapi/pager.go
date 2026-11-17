package webexapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// ErrIncompleteList means a list could not be read to its end. The pager fails
// rather than return the pages it has, which would look like the whole list.
var ErrIncompleteList = errors.New("incomplete list")

// PagingMode is how a list endpoint pages.
type PagingMode int

const (
	// Single: one request returns the whole list.
	Single PagingMode = iota
	// Link: follow Link rel="next" until there is none.
	Link
	// Pages: a page index and size, until a stated total.
	Pages
	// Offset: a start index and size, until a stated total, or a page shorter
	// than asked when the endpoint states none.
	Offset
)

// Paging is how one list endpoint is read to its end.
type Paging struct {
	Mode PagingMode
	// Page (Pages) and Start (Offset) name the index parameter; First is the
	// index of the first page or item. In Link mode, Start names the endpoint's
	// index parameter, which is left out so the list is read from its start.
	Page, Start string
	First       int
	// Size names the page-size parameter; empty sends none. PageSize is the size
	// sent, 100 when zero.
	Size     string
	PageSize int
	// Items is the dotted path to the item array; empty is a bare array.
	// OmitEmpty means an empty list may leave it out.
	Items     string
	OmitEmpty bool
	// TotalPages and TotalItems are the dotted paths of the field stating the
	// list's size, in pages or in items.
	TotalPages, TotalItems string
	// LinkQueryOnly means the endpoint's next link has the wrong path, and its
	// query is applied to the request path.
	LinkQueryOnly bool
	// Query holds parameters the paging needs, such as one that makes the
	// endpoint state its total.
	Query map[string]string
	// NotFoundEnds means the endpoint answers 404 for a page past its last one.
	NotFoundEnds bool
}

const defaultPageSize = 100

// List is a list request before paging.
type List struct {
	Method string

	// URL has the path filled in and no query.
	URL string

	// Query holds the caller's parameters. The paging's parameters are set over
	// them.
	Query url.Values

	Body        []byte
	ContentType string
	Header      http.Header
	Redact      func([]byte) []byte

	// Name is how errors name the list, such as its path template.
	Name string
}

// All reads every page of a list and returns its items.
func (c *Client) All(ctx context.Context, l List, p Paging) ([]json.RawMessage, error) {
	items := []json.RawMessage{}
	err := c.Visit(ctx, l, p, func(page []json.RawMessage) bool {
		items = append(items, page...)
		return true
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

// Visit reads a list page by page, passing each page's items to visit, and stops
// as soon as visit returns false.
func (c *Client) Visit(ctx context.Context, l List, p Paging, visit func(items []json.RawMessage) bool) error {
	pg := pager{c: c, ctx: ctx, l: l, p: p}
	size := p.PageSize
	if size == 0 {
		size = defaultPageSize
	}
	switch p.Mode {
	case Single:
		body, _, err := pg.fetch(nil, "")
		if err != nil {
			return err
		}
		items, err := pg.items(body)
		if err != nil {
			return err
		}
		visit(items)
		return nil
	case Link:
		return pg.visitLinks(size, visit)
	case Pages, Offset:
		return pg.visitIndexed(size, visit)
	}
	return fmt.Errorf("unknown paging mode %d", p.Mode)
}

type pager struct {
	c   *Client
	ctx context.Context
	l   List
	p   Paging
}

func (pg *pager) visitLinks(size int, visit func([]json.RawMessage) bool) error {
	query := map[string]string{}
	if pg.p.Size != "" {
		query[pg.p.Size] = strconv.Itoa(size)
	}
	first := *pg
	if pg.p.Start != "" {
		first.l.Query = url.Values{}
		for k, v := range pg.l.Query {
			if k != pg.p.Start {
				first.l.Query[k] = v
			}
		}
	}
	next := ""
	seen := map[string]bool{}
	for {
		body, header, err := first.fetch(query, next)
		if err != nil {
			return err
		}
		items, err := pg.items(body)
		if err != nil {
			return err
		}
		if !visit(items) {
			return nil
		}
		link, err := pg.nextLink(header)
		if err != nil || link == "" {
			return err
		}
		if len(items) == 0 {
			return pg.incomplete("it links to another page after an empty one")
		}
		if pg.p.LinkQueryOnly {
			u, err := url.Parse(link)
			if err != nil {
				return pg.incomplete("its next link %q is not a URL", link)
			}
			link = pg.l.URL + "?" + u.RawQuery
		}
		if seen[link] {
			return pg.incomplete("its next link repeats a page")
		}
		seen[link] = true
		next = link
	}
}

func (pg *pager) visitIndexed(size int, visit func([]json.RawMessage) bool) error {
	param := pg.p.Page
	if pg.p.Mode == Offset {
		param = pg.p.Start
	}
	index, read := pg.p.First, 0
	var prev []byte
	for {
		query := map[string]string{param: strconv.Itoa(index)}
		if pg.p.Size != "" {
			query[pg.p.Size] = strconv.Itoa(size)
		}
		resp, err := pg.send(query, "")
		if err != nil {
			if pg.p.NotFoundEnds && resp != nil && resp.Status == http.StatusNotFound && index != pg.p.First {
				return nil
			}
			return err
		}
		items, err := pg.items(resp.Body)
		if err != nil {
			return err
		}
		if !visit(items) {
			return nil
		}
		read += len(items)
		switch {
		case pg.p.TotalPages != "":
			total, ok := intField(resp.Body, pg.p.TotalPages)
			if !ok {
				return pg.incomplete("its response has no %s", pg.p.TotalPages)
			}
			if index+1-pg.p.First >= total {
				return nil
			}
			if len(items) == 0 {
				return pg.incomplete("page %d of %d is empty", index+1-pg.p.First, total)
			}
		case pg.p.TotalItems != "":
			total, ok := intField(resp.Body, pg.p.TotalItems)
			if !ok {
				return pg.incomplete("its response has no %s", pg.p.TotalItems)
			}
			if read >= total {
				return nil
			}
			if len(items) == 0 {
				return pg.incomplete("it ended at %d of %d items", read, total)
			}
		default:
			if len(items) < size {
				return nil
			}
			// Without a total only a page shorter than asked ends the list, so
			// an endpoint that ignored the index would never end.
			page, _ := json.Marshal(items)
			if bytes.Equal(page, prev) {
				return pg.incomplete("it returned the same page again")
			}
			prev = page
		}
		if pg.p.Mode == Pages {
			index++
		} else {
			index += len(items)
		}
	}
}

func (pg *pager) fetch(query map[string]string, rawURL string) ([]byte, http.Header, error) {
	resp, err := pg.send(query, rawURL)
	if err != nil {
		return nil, nil, err
	}
	return resp.Body, resp.Header, nil
}

// send sends the list request with the paging's query and then query set over
// its own query, or sends rawURL instead when it is not empty.
func (pg *pager) send(query map[string]string, rawURL string) (*Response, error) {
	target := rawURL
	if target == "" {
		values := url.Values{}
		for k, v := range pg.l.Query {
			values[k] = v
		}
		for _, set := range []map[string]string{pg.p.Query, query} {
			for k, v := range set {
				values.Set(k, v)
			}
		}
		target = WithQuery(pg.l.URL, values)
	}
	return pg.c.Send(pg.ctx, Request{
		Method: pg.l.Method, URL: target, Body: pg.l.Body, ContentType: pg.l.ContentType,
		Header: pg.l.Header, Redact: pg.l.Redact,
	})
}

// WithQuery adds the query to a URL that has none. Parameters with an empty value
// are left out.
func WithQuery(rawURL string, query url.Values) string {
	values := url.Values{}
	for k, vs := range query {
		for _, v := range vs {
			if v != "" {
				values.Add(k, v)
			}
		}
	}
	if encoded := values.Encode(); encoded != "" {
		return rawURL + "?" + encoded
	}
	return rawURL
}

// items returns the item array of one page.
func (pg *pager) items(body []byte) ([]json.RawMessage, error) {
	value := bytes.TrimSpace(body)
	if pg.p.Items != "" {
		for _, key := range strings.Split(pg.p.Items, ".") {
			var object map[string]json.RawMessage
			if len(value) == 0 || bytes.Equal(value, []byte("null")) {
				value = nil
				break
			}
			if err := json.Unmarshal(value, &object); err != nil {
				return nil, pg.incomplete("its response is not an object")
			}
			value = object[key]
		}
	}
	if len(value) == 0 || bytes.Equal(value, []byte("null")) {
		if pg.p.OmitEmpty {
			return nil, nil
		}
		return nil, pg.incomplete("its response has no %s", itemsName(pg.p.Items))
	}
	var items []json.RawMessage
	if err := json.Unmarshal(value, &items); err != nil {
		return nil, pg.incomplete("%s is not an array", itemsName(pg.p.Items))
	}
	return items, nil
}

func itemsName(path string) string {
	if path == "" {
		return "item array"
	}
	return `"` + path + `"`
}

func intField(body []byte, path string) (int, bool) {
	value := json.RawMessage(body)
	for _, key := range strings.Split(path, ".") {
		var object map[string]json.RawMessage
		if json.Unmarshal(value, &object) != nil {
			return 0, false
		}
		value = object[key]
	}
	var n int
	return n, json.Unmarshal(value, &n) == nil
}

func (pg *pager) incomplete(format string, args ...any) error {
	return fmt.Errorf("%w: %s: %s", ErrIncompleteList, pg.l.Name, fmt.Sprintf(format, args...))
}

var linkPart = regexp.MustCompile(`<([^>]*)>([^<]*)`)

// nextLink returns the rel="next" URL of a Link header, resolved against the
// list URL. A next page on another host is an error, so the token is never sent
// there.
func (pg *pager) nextLink(h http.Header) (string, error) {
	base, err := url.Parse(pg.l.URL)
	if err != nil {
		return "", err
	}
	for _, value := range h.Values("Link") {
		for _, m := range linkPart.FindAllStringSubmatch(value, -1) {
			params := strings.ReplaceAll(m[2], " ", "")
			if !strings.Contains(params, `rel="next"`) && !strings.Contains(params, "rel=next") {
				continue
			}
			ref, err := url.Parse(m[1])
			if err != nil {
				return "", pg.incomplete("its next link %q is not a URL", m[1])
			}
			next := base.ResolveReference(ref)
			if !sameHost(next, base) {
				return "", pg.incomplete("its next page is on another host")
			}
			return next.String(), nil
		}
	}
	return "", nil
}

func sameHost(a, b *url.URL) bool {
	return a.Scheme == b.Scheme && strings.EqualFold(a.Hostname(), b.Hostname()) && port(a) == port(b)
}

func port(u *url.URL) string {
	if p := u.Port(); p != "" {
		return p
	}
	if u.Scheme == "http" {
		return "80"
	}
	return "443"
}
