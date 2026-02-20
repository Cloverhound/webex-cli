package search

import (
	"encoding/json"
	"fmt"
	"strings"
)

// QueryType enumerates the supported search query types.
type QueryType string

const (
	QueryTask             QueryType = "task"
	QueryTaskDetails      QueryType = "taskDetails"
	QueryTaskLegDetails   QueryType = "taskLegDetails"
	QueryAgentSession     QueryType = "agentSession"
	QueryFlowInteractions QueryType = "flowInteractions"
	QueryFlowTraceEvents  QueryType = "flowTraceEvents"
)

// ResultsKey returns the top-level key in the query response for each query type.
var ResultsKey = map[QueryType]string{
	QueryTask:             "tasks",
	QueryTaskDetails:      "tasks",
	QueryTaskLegDetails:   "tasks",
	QueryAgentSession:     "agentSessions",
	QueryFlowInteractions: "flowInteractions",
	QueryFlowTraceEvents:  "flowTraceEvents",
}

// Aggregation defines a single aggregation request.
type Aggregation struct {
	Field string
	Type  string
	Name  string
}

// Params holds the common parameters for all search queries.
type Params struct {
	From           int64
	To             int64
	TimeComparator string
	Filter         string // pre-built GraphQL filter string (from BuildFilter)
	Fields         string // raw GraphQL field selection (overrides default)
	Cursor         string
	PageSize       int
	Aggregations   []Aggregation
	Interval       string
	Timezone       string
}

// BuildQuery constructs the full JSON request body with the GraphQL query.
func BuildQuery(queryType QueryType, params Params) (string, error) {
	var qb strings.Builder

	qb.WriteString("{ ")
	qb.WriteString(string(queryType))
	qb.WriteString("(")

	// Required: from and to
	qb.WriteString(fmt.Sprintf("from: %d, to: %d", params.From, params.To))

	// Optional: timeComparator
	if params.TimeComparator != "" {
		qb.WriteString(fmt.Sprintf(", timeComparator: %s", params.TimeComparator))
	}

	// Optional: filter
	if params.Filter != "" {
		qb.WriteString(", ")
		qb.WriteString(params.Filter)
	}

	// Optional: pagination
	if queryType == QueryFlowInteractions || queryType == QueryFlowTraceEvents {
		if params.PageSize > 0 {
			qb.WriteString(fmt.Sprintf(", pagination: { pageSize: %d", params.PageSize))
			qb.WriteString(", currentPage: 1 }")
		}
	} else if params.Cursor != "" {
		qb.WriteString(fmt.Sprintf(`, pagination: { cursor: "%s" }`, params.Cursor))
	}

	// Optional: aggregations
	if len(params.Aggregations) > 0 {
		qb.WriteString(", aggregations: [")
		for i, agg := range params.Aggregations {
			if i > 0 {
				qb.WriteString(", ")
			}
			qb.WriteString(fmt.Sprintf(`{ field: "%s", type: %s, name: "%s" }`, agg.Field, agg.Type, agg.Name))
		}
		qb.WriteString("]")
	}

	// Optional: aggregationInterval
	if params.Interval != "" {
		qb.WriteString(fmt.Sprintf(`, aggregationInterval: { interval: %s`, params.Interval))
		if params.Timezone != "" {
			qb.WriteString(fmt.Sprintf(`, timezone: "%s"`, params.Timezone))
		}
		qb.WriteString(" }")
	}

	qb.WriteString(") { ")

	// Fields
	fields := params.Fields
	if fields == "" {
		fields = DefaultFields[queryType]
	}
	qb.WriteString(fields)

	qb.WriteString(" } }")

	// Wrap in JSON body
	body := map[string]string{"query": qb.String()}
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal query: %w", err)
	}
	return string(jsonBytes), nil
}

// ParseAggregation parses a string like "totalDuration:sum:Total Duration" into an Aggregation.
// Format: field:type[:name] — name defaults to "field_type" if omitted.
func ParseAggregation(s string) (Aggregation, error) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) < 2 {
		return Aggregation{}, fmt.Errorf("invalid aggregation %q: expected field:type[:name]", s)
	}

	field := parts[0]
	aggType := strings.ToLower(parts[1])

	validTypes := map[string]bool{"sum": true, "average": true, "count": true, "min": true, "max": true, "cardinality": true}
	if !validTypes[aggType] {
		return Aggregation{}, fmt.Errorf("invalid aggregation type %q: must be sum, average, count, min, max, or cardinality", aggType)
	}

	name := field + "_" + aggType
	if len(parts) == 3 {
		name = parts[2]
	}

	return Aggregation{Field: field, Type: aggType, Name: name}, nil
}
