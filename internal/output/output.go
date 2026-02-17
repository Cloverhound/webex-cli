package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
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

// valueKind classifies a JSON value for table display.
type valueKind int

const (
	kindScalar      valueKind = iota // string, number, bool, nil
	kindSimpleArray                  // []string or []number
	kindComplex                      // nested object or array of objects
)

// classifyValue returns a formatted string and a kind for a JSON value.
func classifyValue(v interface{}) (string, valueKind) {
	switch val := v.(type) {
	case string:
		return val, kindScalar
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val)), kindScalar
		}
		return fmt.Sprintf("%g", val), kindScalar
	case bool:
		if val {
			return "true", kindScalar
		}
		return "false", kindScalar
	case nil:
		return "", kindScalar
	case []interface{}:
		if len(val) == 0 {
			return "", kindScalar
		}
		// Check if all elements are scalars (string or number)
		parts := make([]string, 0, len(val))
		for _, elem := range val {
			switch e := elem.(type) {
			case string:
				parts = append(parts, e)
			case float64:
				if e == float64(int64(e)) {
					parts = append(parts, fmt.Sprintf("%d", int64(e)))
				} else {
					parts = append(parts, fmt.Sprintf("%g", e))
				}
			default:
				return "", kindComplex
			}
		}
		return strings.Join(parts, ", "), kindSimpleArray
	default:
		// map[string]interface{} or anything else
		return "", kindComplex
	}
}

// acronyms that should be uppercased in headers.
var acronyms = map[string]bool{
	"id": true, "url": true, "api": true, "sip": true,
	"mac": true, "ip": true, "uri": true,
	"esn": true, "pstn": true, "uuid": true, "dn": true,
	"did": true, "ani": true, "dnis": true, "ata": true,
	"http": true, "https": true, "html": true,
}

// camelToHeader converts a camelCase key to a human-readable header.
// Examples: "displayName" → "Display Name", "orgId" → "Org ID", "sipAddress" → "SIP Address"
func camelToHeader(s string) string {
	if s == "" {
		return s
	}

	// Split on camelCase boundaries
	var words []string
	runes := []rune(s)
	start := 0
	for i := 1; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) && unicode.IsLower(runes[i-1]) {
			words = append(words, string(runes[start:i]))
			start = i
		} else if unicode.IsUpper(runes[i-1]) && unicode.IsUpper(runes[i]) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
			// Handles transitions like "orgID" → ["org", "ID"] at the boundary
			words = append(words, string(runes[start:i]))
			start = i
		}
	}
	words = append(words, string(runes[start:]))

	// Title-case each word, uppercase known acronyms
	for i, w := range words {
		lower := strings.ToLower(w)
		if acronyms[lower] {
			words[i] = strings.ToUpper(w)
		} else {
			words[i] = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
	}

	return strings.Join(words, " ")
}

// getTerminalWidth returns the current terminal width, or 120 for non-TTY.
func getTerminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 120
	}
	return w
}

// preferredArrayKeys is checked in order when unwrapping a wrapper object.
var preferredArrayKeys = []string{"items", "data"}

func printTable(data []byte) error {
	var items []map[string]interface{}

	// Try to parse as array of objects directly
	if err := json.Unmarshal(data, &items); err != nil {
		// Try as single object
		var obj map[string]interface{}
		if err2 := json.Unmarshal(data, &obj); err2 != nil {
			return printJSON(data) // fallback
		}

		// Try preferred keys first
		found := false
		foundEmpty := false
		for _, key := range preferredArrayKeys {
			if v, ok := obj[key]; ok {
				if raw, err3 := json.Marshal(v); err3 == nil {
					var arr []map[string]interface{}
					if err3 = json.Unmarshal(raw, &arr); err3 == nil {
						items = arr
						found = true
						if len(arr) == 0 {
							foundEmpty = true
						}
						break
					}
				}
			}
		}

		// If not found, try first array value (sorted keys for determinism)
		if !found {
			sortedKeys := make([]string, 0, len(obj))
			for k := range obj {
				sortedKeys = append(sortedKeys, k)
			}
			sort.Strings(sortedKeys)

			for _, k := range sortedKeys {
				v := obj[k]
				if raw, err3 := json.Marshal(v); err3 == nil {
					var arr []map[string]interface{}
					if err3 = json.Unmarshal(raw, &arr); err3 == nil {
						items = arr
						found = true
						if len(arr) == 0 {
							foundEmpty = true
						}
						break
					}
				}
			}
		}

		// If an array was found but empty, show empty message
		if foundEmpty {
			fmt.Println("(empty)")
			return nil
		}

		// If no array found at all, treat the single object as a one-row table
		if !found {
			items = []map[string]interface{}{obj}
		}
	}

	if len(items) == 0 {
		fmt.Println("(empty)")
		return nil
	}

	// Classify all values and determine which columns to include.
	// A column is included if at least one row has a non-complex value for it.
	keySet := map[string]bool{}
	for _, item := range items {
		for k := range item {
			keySet[k] = true
		}
	}

	allKeys := make([]string, 0, len(keySet))
	for k := range keySet {
		allKeys = append(allKeys, k)
	}
	sort.Strings(allKeys)

	// Build formatted rows and track which columns have at least one non-complex value
	type cellInfo struct {
		text string
		kind valueKind
	}
	rowCells := make([]map[string]cellInfo, len(items))
	columnHasScalar := map[string]bool{}

	for i, item := range items {
		rowCells[i] = make(map[string]cellInfo, len(allKeys))
		for _, k := range allKeys {
			val, ok := item[k]
			if !ok {
				rowCells[i][k] = cellInfo{"", kindScalar}
				continue
			}
			text, kind := classifyValue(val)
			rowCells[i][k] = cellInfo{text, kind}
			if kind != kindComplex {
				columnHasScalar[k] = true
			}
		}
	}

	// Filter to only included columns
	var cols []string
	for _, k := range allKeys {
		if columnHasScalar[k] {
			cols = append(cols, k)
		}
	}

	if len(cols) == 0 {
		// All columns are complex — fall back to JSON
		return printJSON(data)
	}

	// Build headers
	headers := make([]string, len(cols))
	for i, k := range cols {
		headers[i] = camelToHeader(k)
	}

	// Build string rows
	rows := make([][]string, len(items))
	for i, cells := range rowCells {
		row := make([]string, len(cols))
		for j, k := range cols {
			row[j] = cells[k].text
		}
		rows[i] = row
	}

	// Render with tablewriter
	termWidth := getTerminalWidth()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(termWidth)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	// Truncate wide cells to keep table within terminal width
	// Account for borders: | col | col | = len(cols)+1 pipe chars + 2*len(cols) spaces
	borderOverhead := len(cols) + 1 + 2*len(cols)
	availWidth := termWidth - borderOverhead
	if availWidth < len(cols)*4 {
		availWidth = len(cols) * 4
	}
	maxColWidth := availWidth / len(cols)
	if maxColWidth < 4 {
		maxColWidth = 4
	}

	for i, row := range rows {
		for j, cell := range row {
			if len(cell) > maxColWidth {
				rows[i][j] = cell[:maxColWidth-3] + "..."
			}
		}
	}

	table.AppendBulk(rows)
	table.Render()

	return nil
}
