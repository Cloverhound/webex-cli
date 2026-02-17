package output

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestCamelToHeader(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"displayName", "Display Name"},
		{"orgId", "Org ID"},
		{"sipAddress", "SIP Address"},
		{"firstName", "First Name"},
		{"isActive", "Is Active"},
		{"id", "ID"},
		{"email", "Email"},
		{"macAddress", "MAC Address"},
		{"ipAddress", "IP Address"},
		{"created", "Created"},
	}
	for _, tt := range tests {
		got := camelToHeader(tt.input)
		if got != tt.want {
			t.Errorf("camelToHeader(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestClassifyValue(t *testing.T) {
	text, kind := classifyValue("hello")
	if kind != kindScalar || text != "hello" {
		t.Errorf("string: got %q %v", text, kind)
	}

	text, kind = classifyValue(float64(42))
	if kind != kindScalar || text != "42" {
		t.Errorf("int: got %q %v", text, kind)
	}

	text, kind = classifyValue(true)
	if kind != kindScalar || text != "true" {
		t.Errorf("bool: got %q %v", text, kind)
	}

	text, kind = classifyValue(nil)
	if kind != kindScalar || text != "" {
		t.Errorf("nil: got %q %v", text, kind)
	}

	text, kind = classifyValue([]interface{}{"a", "b", "c"})
	if kind != kindSimpleArray || text != "a, b, c" {
		t.Errorf("string array: got %q %v", text, kind)
	}

	text, kind = classifyValue([]interface{}{float64(1), float64(2)})
	if kind != kindSimpleArray || text != "1, 2" {
		t.Errorf("number array: got %q %v", text, kind)
	}

	_, kind = classifyValue(map[string]interface{}{"key": "val"})
	if kind != kindComplex {
		t.Errorf("object: got kind %v, want kindComplex", kind)
	}

	_, kind = classifyValue([]interface{}{map[string]interface{}{"key": "val"}})
	if kind != kindComplex {
		t.Errorf("array of objects: got kind %v, want kindComplex", kind)
	}
}

func TestPrintTableWrapperObject(t *testing.T) {
	SetFormat("table")
	data := []byte(`{"items":[{"displayName":"Alice","orgId":"abc-123","created":"2024-01-01"},{"displayName":"Bob","orgId":"def-456","created":"2024-02-15"}]}`)
	out := captureStdout(func() {
		Print(data, 200)
	})
	// Should contain human-readable headers
	if !strings.Contains(out, "Display Name") {
		t.Errorf("expected 'Display Name' header, got:\n%s", out)
	}
	if !strings.Contains(out, "Org ID") {
		t.Errorf("expected 'Org ID' header, got:\n%s", out)
	}
	// Should contain data
	if !strings.Contains(out, "Alice") {
		t.Errorf("expected 'Alice' in output, got:\n%s", out)
	}
	// Should have table borders
	if !strings.Contains(out, "|") {
		t.Errorf("expected table borders '|', got:\n%s", out)
	}
}

func TestPrintTableSingleObject(t *testing.T) {
	SetFormat("table")
	data := []byte(`{"displayName":"Charlie","email":"charlie@example.com","isActive":true}`)
	out := captureStdout(func() {
		Print(data, 200)
	})
	if !strings.Contains(out, "Display Name") {
		t.Errorf("expected 'Display Name' header, got:\n%s", out)
	}
	if !strings.Contains(out, "Charlie") {
		t.Errorf("expected 'Charlie' in output, got:\n%s", out)
	}
}

func TestPrintTableExcludesComplexColumns(t *testing.T) {
	SetFormat("table")
	data := []byte(`{"items":[{"name":"Alice","details":{"nested":"obj"}}]}`)
	out := captureStdout(func() {
		Print(data, 200)
	})
	if !strings.Contains(out, "Name") {
		t.Errorf("expected 'Name' header, got:\n%s", out)
	}
	if strings.Contains(out, "Details") {
		t.Errorf("should not contain complex column 'Details', got:\n%s", out)
	}
}

func TestPrintTableEmpty(t *testing.T) {
	SetFormat("table")
	data := []byte(`{"items":[]}`)
	out := captureStdout(func() {
		Print(data, 200)
	})
	if !strings.Contains(out, "(empty)") {
		t.Errorf("expected '(empty)' for empty items, got:\n%s", out)
	}
}

func TestPrintTableSimpleArrayIncluded(t *testing.T) {
	SetFormat("table")
	data := []byte(`[{"name":"Alice","emails":["a@b.com","c@d.com"]}]`)
	out := captureStdout(func() {
		Print(data, 200)
	})
	if !strings.Contains(out, "Emails") {
		t.Errorf("expected 'Emails' header for simple array, got:\n%s", out)
	}
	if !strings.Contains(out, "a@b.com, c@d.com") {
		t.Errorf("expected joined emails, got:\n%s", out)
	}
}

// --- CSV tests ---

func TestPrintCSVWrapperObject(t *testing.T) {
	SetFormat("csv")
	data := []byte(`{"items":[{"displayName":"Alice","orgId":"abc-123","created":"2024-01-01"},{"displayName":"Bob","orgId":"def-456","created":"2024-02-15"}]}`)
	out := captureStdout(func() {
		Print(data, 200)
	})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines (header + 2 rows), got %d:\n%s", len(lines), out)
	}
	// Header should contain human-readable column names
	if !strings.Contains(lines[0], "Display Name") {
		t.Errorf("expected 'Display Name' in header, got: %s", lines[0])
	}
	if !strings.Contains(lines[0], "Org ID") {
		t.Errorf("expected 'Org ID' in header, got: %s", lines[0])
	}
	// Data rows
	if !strings.Contains(lines[1], "Alice") {
		t.Errorf("expected 'Alice' in first data row, got: %s", lines[1])
	}
	if !strings.Contains(lines[2], "Bob") {
		t.Errorf("expected 'Bob' in second data row, got: %s", lines[2])
	}
}

func TestPrintCSVSingleObject(t *testing.T) {
	SetFormat("csv")
	data := []byte(`{"displayName":"Charlie","email":"charlie@example.com","isActive":true}`)
	out := captureStdout(func() {
		Print(data, 200)
	})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines (header + 1 row), got %d:\n%s", len(lines), out)
	}
	if !strings.Contains(lines[0], "Display Name") {
		t.Errorf("expected 'Display Name' in header, got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "Charlie") {
		t.Errorf("expected 'Charlie' in data row, got: %s", lines[1])
	}
}

func TestPrintCSVExcludesComplexColumns(t *testing.T) {
	SetFormat("csv")
	data := []byte(`{"items":[{"name":"Alice","details":{"nested":"obj"}}]}`)
	out := captureStdout(func() {
		Print(data, 200)
	})
	if !strings.Contains(out, "Name") {
		t.Errorf("expected 'Name' header, got:\n%s", out)
	}
	if strings.Contains(out, "Details") {
		t.Errorf("should not contain complex column 'Details', got:\n%s", out)
	}
}

func TestPrintCSVEmpty(t *testing.T) {
	SetFormat("csv")
	data := []byte(`{"items":[]}`)
	out := captureStdout(func() {
		Print(data, 200)
	})
	if out != "" {
		t.Errorf("expected empty output for empty items, got:\n%s", out)
	}
}

func TestPrintCSVQuoting(t *testing.T) {
	SetFormat("csv")
	data := []byte(`[{"name":"O'Brien, James","note":"He said \"hello\""}]`)
	out := captureStdout(func() {
		Print(data, 200)
	})
	// CSV should quote fields with commas
	if !strings.Contains(out, `"O'Brien, James"`) {
		t.Errorf("expected quoted field with comma, got:\n%s", out)
	}
	// CSV should escape double quotes by doubling them
	if !strings.Contains(out, `"He said ""hello"""`) {
		t.Errorf("expected escaped double quotes, got:\n%s", out)
	}
}
