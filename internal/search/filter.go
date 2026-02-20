package search

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BuildFilter constructs a GraphQL filter string from convenience flags and an optional raw JSON filter.
// Convenience flags and raw filter are AND-merged when both are provided.
func BuildFilter(channel, direction, status, agentID, queueID, rawJSON string) (string, error) {
	var parts []string

	if channel != "" {
		parts = append(parts, fmt.Sprintf(`{ channelType: { equals: %s } }`, channel))
	}
	if direction != "" {
		parts = append(parts, fmt.Sprintf(`{ direction: { equals: "%s" } }`, direction))
	}
	if status != "" {
		parts = append(parts, fmt.Sprintf(`{ status: { equals: "%s" } }`, status))
	}
	if agentID != "" {
		parts = append(parts, fmt.Sprintf(`{ owner: { id: { equals: "%s" } } }`, agentID))
	}
	if queueID != "" {
		parts = append(parts, fmt.Sprintf(`{ lastQueue: { id: { equals: "%s" } } }`, queueID))
	}

	// Parse raw JSON filter if provided
	if rawJSON != "" {
		rawJSON = strings.TrimSpace(rawJSON)
		// Validate it's valid JSON
		var js json.RawMessage
		if err := json.Unmarshal([]byte(rawJSON), &js); err != nil {
			return "", fmt.Errorf("invalid --filter JSON: %w", err)
		}
		// Convert JSON filter to GraphQL syntax
		gql, err := jsonFilterToGraphQL(js)
		if err != nil {
			return "", err
		}
		parts = append(parts, gql)
	}

	if len(parts) == 0 {
		return "", nil
	}
	if len(parts) == 1 {
		return "filter: " + parts[0], nil
	}
	return "filter: { and: [" + strings.Join(parts, ", ") + "] }", nil
}

// jsonFilterToGraphQL converts a JSON filter object to GraphQL syntax.
// The JSON uses standard JSON format; we convert to GraphQL object notation.
func jsonFilterToGraphQL(raw json.RawMessage) (string, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &obj); err != nil {
		return "", fmt.Errorf("filter must be a JSON object: %w", err)
	}
	return jsonObjToGraphQL(obj)
}

func jsonObjToGraphQL(obj map[string]json.RawMessage) (string, error) {
	var parts []string
	for key, val := range obj {
		gql, err := jsonValToGraphQL(val)
		if err != nil {
			return "", err
		}
		parts = append(parts, key+": "+gql)
	}
	return "{ " + strings.Join(parts, ", ") + " }", nil
}

func jsonValToGraphQL(raw json.RawMessage) (string, error) {
	s := strings.TrimSpace(string(raw))

	// Array
	if s[0] == '[' {
		var arr []json.RawMessage
		if err := json.Unmarshal(raw, &arr); err != nil {
			return "", err
		}
		var items []string
		for _, item := range arr {
			gql, err := jsonValToGraphQL(item)
			if err != nil {
				return "", err
			}
			items = append(items, gql)
		}
		return "[" + strings.Join(items, ", ") + "]", nil
	}

	// Object
	if s[0] == '{' {
		var obj map[string]json.RawMessage
		if err := json.Unmarshal(raw, &obj); err != nil {
			return "", err
		}
		return jsonObjToGraphQL(obj)
	}

	// String — check if it looks like an enum value (no spaces, lowercase/underscored)
	if s[0] == '"' {
		var str string
		if err := json.Unmarshal(raw, &str); err != nil {
			return "", err
		}
		// Pass through as quoted string
		return fmt.Sprintf("%q", str), nil
	}

	// Number, boolean, null — pass through as-is
	return s, nil
}
