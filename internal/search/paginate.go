package search

import (
	"encoding/json"
	"fmt"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// PaginateGraphQL reads every page of a search query and returns their items as
// one JSON array. It fails rather than return part of a list.
func PaginateGraphQL(queryType QueryType, params Params, orgID, trackingID string) ([]byte, int, error) {
	flow := isFlowQuery(queryType)
	if flow && params.PageSize == 0 {
		params.PageSize = 100
	}

	allItems := []json.RawMessage{}
	for page := 1; ; page++ {
		body, err := BuildQuery(queryType, params)
		if err != nil {
			return nil, 0, err
		}

		req := client.NewRequest(config.CcBaseURL, "POST", "/search")
		req.QueryParam("orgId", orgID)
		req.Header("TrackingId", trackingID)
		req.SetBodyRaw(body)

		resp, _, err := req.Do()
		if err != nil {
			if page == 1 {
				return nil, 0, err
			}
			return nil, 0, incomplete("page %d: %v", page, err)
		}
		items, info, err := readPage(resp, queryType)
		if err != nil {
			return nil, 0, incomplete("page %d: %v", page, err)
		}
		allItems = append(allItems, items...)

		if flow {
			// totalPages is the record count divided by the page size, rounded
			// down, so the end of a list is found by counting records.
			if info.TotalRecords == nil {
				return nil, 0, incomplete("page %d has no pageInfo.totalRecords", page)
			}
			if len(allItems) >= *info.TotalRecords {
				break
			}
			if len(items) == 0 {
				return nil, 0, incomplete("page %d is empty after %d of %d records", page, len(allItems), *info.TotalRecords)
			}
			params.CurrentPage++
			continue
		}
		if !info.HasNextPage {
			break
		}
		if info.EndCursor == "" {
			return nil, 0, incomplete("page %d has a next page but no endCursor", page)
		}
		params.Cursor = info.EndCursor
	}

	data, err := json.Marshal(allItems)
	if err != nil {
		return nil, 0, fmt.Errorf("marshaling paginated results: %w", err)
	}
	return data, 200, nil
}

// pageInfo holds the paging fields of both kinds of search query.
type pageInfo struct {
	HasNextPage  bool   `json:"hasNextPage"`
	EndCursor    string `json:"endCursor"`
	TotalRecords *int   `json:"totalRecords"`
}

// readPage is the items and paging of one page of search results:
// { "data": { "<queryType>": { "<resultsKey>": [...], "pageInfo": {...} } } }.
func readPage(resp []byte, queryType QueryType) ([]json.RawMessage, pageInfo, error) {
	var envelope struct {
		Data   map[string]json.RawMessage `json:"data"`
		Errors []json.RawMessage          `json:"errors"`
	}
	if err := json.Unmarshal(resp, &envelope); err != nil {
		return nil, pageInfo{}, fmt.Errorf("reading search response: %w", err)
	}
	if len(envelope.Errors) > 0 {
		errs, _ := json.Marshal(envelope.Errors)
		return nil, pageInfo{}, fmt.Errorf("search returned errors: %s", errs)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(envelope.Data[string(queryType)], &result); err != nil || result == nil {
		return nil, pageInfo{}, fmt.Errorf("search response has no %s results", queryType)
	}
	resultsKey := ResultsKey[queryType]
	var items []json.RawMessage
	if err := json.Unmarshal(result[resultsKey], &items); err != nil {
		return nil, pageInfo{}, fmt.Errorf("search response has no %s list", resultsKey)
	}
	var info pageInfo
	if err := json.Unmarshal(result["pageInfo"], &info); err != nil || result["pageInfo"] == nil {
		return nil, pageInfo{}, fmt.Errorf("search response has no pageInfo; --fields must include it to page")
	}
	return items, info, nil
}

// incomplete is the error for a list that could not be read to its end.
func incomplete(format string, args ...any) error {
	return fmt.Errorf("%w: search: %s", webexapi.ErrIncompleteList, fmt.Sprintf(format, args...))
}
