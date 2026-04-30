package mcp

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mcplib "github.com/mark3labs/mcp-go/mcp"
)

// parseItems unmarshals a Webex API response body and returns a flat item list.
func parseItems(body []byte) []map[string]any {
	var raw any
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}
	return items(raw)
}

// items normalises a parsed JSON value into a flat []map[string]any.
// Handles {"items":[...]}, bare arrays, and single objects.
func items(data any) []map[string]any {
	if data == nil {
		return nil
	}
	switch v := data.(type) {
	case []any:
		out := make([]map[string]any, 0, len(v))
		for _, elem := range v {
			if m, ok := elem.(map[string]any); ok {
				out = append(out, m)
			}
		}
		return out
	case map[string]any:
		if raw, ok := v["items"]; ok {
			return items(raw)
		}
		if len(v) > 0 {
			return []map[string]any{v}
		}
	}
	return nil
}

// str safely extracts a string from a map value.
func str(v any) string {
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}

// fmtTS formats an ISO 8601 timestamp for human reading.
func fmtTS(ts string) string {
	if ts == "" {
		return "—"
	}
	formats := []string{
		"2006-01-02T15:04:05.999999999Z",
		"2006-01-02T15:04:05Z",
		time.RFC3339Nano,
		time.RFC3339,
	}
	normalised := strings.Replace(ts, "+00:00", "Z", 1)
	for _, layout := range formats {
		if t, err := time.Parse(layout, normalised); err == nil {
			return t.UTC().Format("2006-01-02 15:04 UTC")
		}
	}
	return ts
}

// jsonText marshals a value to indented JSON and wraps it as an MCP tool result.
func jsonText(v any) *mcplib.CallToolResult {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcplib.NewToolResultText(fmt.Sprintf("Error serialising result: %v", err))
	}
	return mcplib.NewToolResultText(string(b))
}
