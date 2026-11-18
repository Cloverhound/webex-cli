package client

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// ErrIncompleteList means --paginate could not read a list to its end. It fails
// rather than returning the pages it has, which would look like the whole list.
var ErrIncompleteList = webexapi.ErrIncompleteList

// Paging is how --paginate reads one list endpoint; see webexapi.Paging.
type Paging = webexapi.Paging

// PagingMode is how a list endpoint pages.
type PagingMode = webexapi.PagingMode

const (
	Single = webexapi.Single
	Link   = webexapi.Link
	Pages  = webexapi.Pages
	Offset = webexapi.Offset
)

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
	r.Finalize()
	return CLIError(API(r.path, r.logBody()).Visit(context.Background(), r.sdkList(), p, visit))
}
