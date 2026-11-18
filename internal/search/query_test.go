package search

import (
	"encoding/json"
	"strings"
	"testing"
)

// graphQL is the query a request body carries.
func graphQL(t *testing.T, body string) string {
	t.Helper()
	var v struct{ Query string }
	if err := json.Unmarshal([]byte(body), &v); err != nil {
		t.Fatal(err)
	}
	return v.Query
}

func TestFlowQueriesAskForFieldsWebexHas(t *testing.T) {
	for queryType, list := range map[QueryType]string{
		QueryFlowTraceEvents:  "traces {",
		QueryFlowInteractions: "interactions {",
	} {
		body, err := BuildQuery(queryType, Params{From: 1, To: 2})
		if err != nil {
			t.Fatal(err)
		}
		q := graphQL(t, body)
		if !strings.Contains(q, list) {
			t.Errorf("%s: query has no %q list: %s", queryType, list, q)
		}
		if strings.Contains(q, string(queryType)+" {") {
			t.Errorf("%s: query asks for itself as its list: %s", queryType, q)
		}
		if !strings.Contains(q, "pageInfo { pageSize currentPage totalPages totalRecords }") {
			t.Errorf("%s: query asks for page fields Webex does not have: %s", queryType, q)
		}
		if !strings.Contains(q, "filter: {}") {
			t.Errorf("%s: query has no filter, which Webex requires: %s", queryType, q)
		}
		if ResultsKey[queryType] != strings.TrimSuffix(list, " {") {
			t.Errorf("%s: results are read from %q", queryType, ResultsKey[queryType])
		}
	}
}

func TestFlowPagesStartAtZero(t *testing.T) {
	body, err := BuildQuery(QueryFlowTraceEvents, Params{From: 1, To: 2, PageSize: 5})
	if err != nil {
		t.Fatal(err)
	}
	if q := graphQL(t, body); !strings.Contains(q, "pagination: { pageSize: 5, currentPage: 0 }") {
		t.Errorf("first page is not page 0: %s", q)
	}
}

func TestAFlowFilterIsSentAsGiven(t *testing.T) {
	body, err := BuildQuery(QueryFlowTraceEvents, Params{From: 1, To: 2, Filter: `filter: { interactionId: { equals: "i1" } }`})
	if err != nil {
		t.Fatal(err)
	}
	q := graphQL(t, body)
	if !strings.Contains(q, `filter: { interactionId: { equals: "i1" } }`) || strings.Contains(q, "filter: {}") {
		t.Errorf("filter not sent as given: %s", q)
	}
}
