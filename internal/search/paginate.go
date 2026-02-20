package search

import (
	"encoding/json"
	"fmt"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
)

// PaginateGraphQL fetches all pages of a GraphQL search query and merges the results.
// It rebuilds the query with updated cursor/page for each iteration.
func PaginateGraphQL(queryType QueryType, params Params, orgID, trackingID string) ([]byte, int, error) {
	resultsKey := ResultsKey[queryType]
	isFlowQuery := queryType == QueryFlowInteractions || queryType == QueryFlowTraceEvents

	var allItems []json.RawMessage

	if isFlowQuery && params.PageSize == 0 {
		params.PageSize = 100
	}
	if isFlowQuery && params.CurrentPage == 0 {
		params.CurrentPage = 1
	}

	for {
		body, err := BuildQuery(queryType, params)
		if err != nil {
			return nil, 0, err
		}

		req := client.NewRequest(config.CcBaseURL, "POST", "/search")
		req.QueryParam("orgId", orgID)
		req.Header("TrackingId", trackingID)
		req.SetBodyRaw(body)

		resp, statusCode, err := req.Do()
		if err != nil {
			return marshalItems(allItems)
		}

		// Parse response: { "data": { "<queryType>": { "<resultsKey>": [...], "pageInfo": {...} } } }
		var envelope struct {
			Data map[string]json.RawMessage `json:"data"`
		}
		if err := json.Unmarshal(resp, &envelope); err != nil {
			return resp, statusCode, nil // can't parse, return as-is
		}

		queryData, ok := envelope.Data[string(queryType)]
		if !ok {
			return resp, statusCode, nil
		}

		var inner map[string]json.RawMessage
		if err := json.Unmarshal(queryData, &inner); err != nil {
			return resp, statusCode, nil
		}

		// Extract items
		if itemsRaw, ok := inner[resultsKey]; ok {
			var items []json.RawMessage
			if err := json.Unmarshal(itemsRaw, &items); err == nil {
				allItems = append(allItems, items...)
			}
		}

		// Check pageInfo
		pageInfoRaw, ok := inner["pageInfo"]
		if !ok {
			break
		}
		var pageInfo struct {
			HasNextPage bool   `json:"hasNextPage"`
			EndCursor   string `json:"endCursor"`
		}
		if err := json.Unmarshal(pageInfoRaw, &pageInfo); err != nil {
			break
		}

		if !pageInfo.HasNextPage {
			break
		}

		// Update pagination for next request
		if isFlowQuery {
			params.CurrentPage++
		} else {
			if pageInfo.EndCursor == "" {
				break
			}
			params.Cursor = pageInfo.EndCursor
		}
	}

	return marshalItems(allItems)
}

func marshalItems(items []json.RawMessage) ([]byte, int, error) {
	if items == nil {
		items = []json.RawMessage{}
	}
	data, err := json.Marshal(items)
	if err != nil {
		return nil, 0, fmt.Errorf("marshaling paginated results: %w", err)
	}
	return data, 200, nil
}
