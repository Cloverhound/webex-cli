package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

var format = "json"

func SetFormat(f string) {
	if f != "" {
		format = f
	}
}

func Format() string { return format }

// Print formats and prints response data based on the output format.
func Print(data []byte, statusCode int) error {
	if len(data) == 0 {
		if statusCode == 204 {
			fmt.Fprintln(os.Stderr, "Success (no content)")
		}
		return nil
	}

	switch format {
	case "raw":
		_, err := os.Stdout.Write(data)
		return err
	case "table":
		return printTable(data)
	default:
		return printJSON(data)
	}
}

func printJSON(data []byte) error {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "  "); err != nil {
		// Not valid JSON — print raw
		_, err := os.Stdout.Write(data)
		return err
	}
	buf.WriteByte('\n')
	_, err := buf.WriteTo(os.Stdout)
	return err
}

func printTable(data []byte) error {
	// Try to parse as array of objects
	var items []map[string]interface{}
	if err := json.Unmarshal(data, &items); err != nil {
		// Try as single object wrapping an array
		var wrapper map[string]json.RawMessage
		if err2 := json.Unmarshal(data, &wrapper); err2 != nil {
			return printJSON(data) // fallback
		}
		// Find first array value
		for _, v := range wrapper {
			if err2 := json.Unmarshal(v, &items); err2 == nil && len(items) > 0 {
				break
			}
		}
		if len(items) == 0 {
			return printJSON(data) // fallback
		}
	}

	if len(items) == 0 {
		fmt.Println("(empty)")
		return nil
	}

	// Collect all keys
	keySet := map[string]bool{}
	for _, item := range items {
		for k := range item {
			keySet[k] = true
		}
	}
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Compute column widths
	widths := make(map[string]int, len(keys))
	for _, k := range keys {
		widths[k] = len(k)
	}
	rows := make([]map[string]string, len(items))
	for i, item := range items {
		rows[i] = map[string]string{}
		for _, k := range keys {
			v := ""
			if val, ok := item[k]; ok {
				v = formatValue(val)
			}
			rows[i][k] = v
			if len(v) > widths[k] {
				widths[k] = len(v)
			}
		}
	}

	// Cap column widths
	for k, w := range widths {
		if w > 50 {
			widths[k] = 50
		}
	}

	// Print header
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = padRight(strings.ToUpper(k), widths[k])
	}
	fmt.Println(strings.Join(parts, "  "))

	// Print separator
	for i, k := range keys {
		parts[i] = strings.Repeat("-", widths[k])
	}
	fmt.Println(strings.Join(parts, "  "))

	// Print rows
	for _, row := range rows {
		for i, k := range keys {
			v := row[k]
			if len(v) > widths[k] {
				v = v[:widths[k]-3] + "..."
			}
			parts[i] = padRight(v, widths[k])
		}
		fmt.Println(strings.Join(parts, "  "))
	}

	return nil
}

func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		b, _ := json.Marshal(val)
		return string(b)
	}
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
