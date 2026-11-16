package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	urlpkg "net/url"
	"regexp"
	"strconv"
	"strings"
)

// ErrIncompleteList means --paginate could not read a list to its end. It fails
// rather than returning the pages it has, which would look like the whole list.
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

// Paging is how --paginate reads one list endpoint. Codegen sets it for each
// endpoint from codegen/paging.py: its API family's paging, or what the endpoint
// was measured to do where that differs.
type Paging struct {
	Mode PagingMode
	// Page (Pages) and Start (Offset) name the index parameter; First is the
	// index of the first page or item.
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

// DoPaginated reads every page of a list endpoint and returns its items as one
// JSON array.
func (r *Request) DoPaginated(p Paging) ([]byte, int, error) {
	items := []json.RawMessage{}
	err := r.Visit(p, func(page []json.RawMessage) bool {
		items = append(items, page...)
		return true
	})
	if err != nil {
		return nil, 0, err
	}
	data, err := json.Marshal(items)
	return data, http.StatusOK, err
}

// Visit reads a list endpoint page by page, passing each page's items to visit,
// and stops as soon as visit returns false.
func (r *Request) Visit(p Paging, visit func(items []json.RawMessage) bool) error {
	size := p.PageSize
	if size == 0 {
		size = defaultPageSize
	}
	switch p.Mode {
	case Single:
		body, _, err := r.fetch(p, nil, "")
		if err != nil {
			return err
		}
		items, err := r.items(p, body)
		if err != nil {
			return err
		}
		visit(items)
		return nil
	case Link:
		return r.visitLinks(p, size, visit)
	case Pages, Offset:
		return r.visitIndexed(p, size, visit)
	}
	return fmt.Errorf("unknown paging mode %d", p.Mode)
}

func (r *Request) visitLinks(p Paging, size int, visit func([]json.RawMessage) bool) error {
	query := map[string]string{}
	if p.Size != "" {
		query[p.Size] = strconv.Itoa(size)
	}
	next := ""
	seen := map[string]bool{}
	for {
		body, header, err := r.fetch(p, query, next)
		if err != nil {
			return err
		}
		items, err := r.items(p, body)
		if err != nil {
			return err
		}
		if !visit(items) {
			return nil
		}
		link, err := r.nextLink(header)
		if err != nil || link == "" {
			return err
		}
		if len(items) == 0 {
			return r.incomplete("it links to another page after an empty one")
		}
		if p.LinkQueryOnly {
			u, err := urlpkg.Parse(link)
			if err != nil {
				return r.incomplete("its next link %q is not a URL", link)
			}
			link = r.url() + "?" + u.RawQuery
		}
		if seen[link] {
			return r.incomplete("its next link repeats a page")
		}
		seen[link] = true
		next = link
	}
}

func (r *Request) visitIndexed(p Paging, size int, visit func([]json.RawMessage) bool) error {
	param := p.Page
	if p.Mode == Offset {
		param = p.Start
	}
	index, read := p.First, 0
	var prev []byte
	for {
		query := map[string]string{param: strconv.Itoa(index)}
		if p.Size != "" {
			query[p.Size] = strconv.Itoa(size)
		}
		body, status, _, err := r.fetchStatus(p, query, "")
		if err != nil {
			if p.NotFoundEnds && status == http.StatusNotFound && index != p.First {
				return nil
			}
			return err
		}
		items, err := r.items(p, body)
		if err != nil {
			return err
		}
		if !visit(items) {
			return nil
		}
		read += len(items)
		switch {
		case p.TotalPages != "":
			total, ok := intField(body, p.TotalPages)
			if !ok {
				return r.incomplete("its response has no %s", p.TotalPages)
			}
			if index+1-p.First >= total {
				return nil
			}
			if len(items) == 0 {
				return r.incomplete("page %d of %d is empty", index+1-p.First, total)
			}
		case p.TotalItems != "":
			total, ok := intField(body, p.TotalItems)
			if !ok {
				return r.incomplete("its response has no %s", p.TotalItems)
			}
			if read >= total {
				return nil
			}
			if len(items) == 0 {
				return r.incomplete("it ended at %d of %d items", read, total)
			}
		default:
			if len(items) < size {
				return nil
			}
			// Without a total only a page shorter than asked ends the list, so
			// an endpoint that ignored the index would never end.
			page, _ := json.Marshal(items)
			if bytes.Equal(page, prev) {
				return r.incomplete("it returned the same page again")
			}
			prev = page
		}
		if p.Mode == Pages {
			index++
		} else {
			index += len(items)
		}
	}
}

// fetch sends the request with the paging's query and then query set over its
// own query parameters, or sends rawURL instead when it is not empty.
func (r *Request) fetch(p Paging, query map[string]string, rawURL string) ([]byte, http.Header, error) {
	body, _, header, err := r.fetchStatus(p, query, rawURL)
	return body, header, err
}

func (r *Request) fetchStatus(p Paging, query map[string]string, rawURL string) ([]byte, int, http.Header, error) {
	page := *r
	page.queryParams = make(map[string]string, len(r.queryParams)+len(p.Query)+len(query))
	for _, set := range []map[string]string{r.queryParams, p.Query, query} {
		for k, v := range set {
			page.queryParams[k] = v
		}
	}
	page.rawURL = rawURL
	page.Finalize()
	return doWithHeaders(&page)
}

// items returns the item array of one page.
func (r *Request) items(p Paging, body []byte) ([]json.RawMessage, error) {
	value := bytes.TrimSpace(body)
	if p.Items != "" {
		for _, key := range strings.Split(p.Items, ".") {
			var object map[string]json.RawMessage
			if len(value) == 0 || bytes.Equal(value, []byte("null")) {
				value = nil
				break
			}
			if err := json.Unmarshal(value, &object); err != nil {
				return nil, r.incomplete("its response is not an object")
			}
			value = object[key]
		}
	}
	if len(value) == 0 || bytes.Equal(value, []byte("null")) {
		if p.OmitEmpty {
			return nil, nil
		}
		return nil, r.incomplete("its response has no %s", itemsName(p.Items))
	}
	var items []json.RawMessage
	if err := json.Unmarshal(value, &items); err != nil {
		return nil, r.incomplete("%s is not an array", itemsName(p.Items))
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

func (r *Request) incomplete(format string, args ...any) error {
	return fmt.Errorf("%w: %s: %s", ErrIncompleteList, r.path, fmt.Sprintf(format, args...))
}

var linkPart = regexp.MustCompile(`<([^>]*)>([^<]*)`)

// nextLink returns the rel="next" URL of a Link header, resolved against the
// request URL. A next page on another host is an error, so the token is never
// sent there.
func (r *Request) nextLink(h http.Header) (string, error) {
	base, err := urlpkg.Parse(r.url())
	if err != nil {
		return "", err
	}
	for _, value := range h.Values("Link") {
		for _, m := range linkPart.FindAllStringSubmatch(value, -1) {
			params := strings.ReplaceAll(m[2], " ", "")
			if !strings.Contains(params, `rel="next"`) && !strings.Contains(params, "rel=next") {
				continue
			}
			ref, err := urlpkg.Parse(m[1])
			if err != nil {
				return "", r.incomplete("its next link %q is not a URL", m[1])
			}
			next := base.ResolveReference(ref)
			if !sameHost(next, base) {
				return "", r.incomplete("its next page is on another host")
			}
			return next.String(), nil
		}
	}
	return "", nil
}

func sameHost(a, b *urlpkg.URL) bool {
	return a.Scheme == b.Scheme && strings.EqualFold(a.Hostname(), b.Hostname()) && port(a) == port(b)
}

func port(u *urlpkg.URL) string {
	if p := u.Port(); p != "" {
		return p
	}
	if u.Scheme == "http" {
		return "80"
	}
	return "443"
}

// url is the request URL with its path parameters filled in and no query.
func (r *Request) url() string {
	url := r.baseURL + r.path
	for k, v := range r.pathParams {
		url = strings.ReplaceAll(url, "{"+k+"}", urlpkg.PathEscape(v))
	}
	return url
}
