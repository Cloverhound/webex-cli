package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	mcplib "github.com/mark3labs/mcp-go/mcp"
)

type usageEntry struct {
	Time    string `json:"time"`
	Tool    string `json:"tool"`
	Command string `json:"command"`
	Flags   string `json:"flags"`
	Status  string `json:"status"`
	Ms      int64  `json:"ms"`
}

type usageLogger struct {
	path     string
	maxBytes int64
	maxFiles int
	file     *os.File
	mu       sync.Mutex
}

func openLogger(path string, maxBytes int64, maxFiles int) (*usageLogger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log: %w", err)
	}
	return &usageLogger{path: path, maxBytes: maxBytes, maxFiles: maxFiles, file: f}, nil
}

func (l *usageLogger) Write(e usageEntry) error {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}
	if _, err := l.file.Write(append(b, '\n')); err != nil {
		return err
	}
	info, err := l.file.Stat()
	if err == nil && info.Size() >= l.maxBytes {
		l.rotateLocked()
	}
	return nil
}

// rotateLocked performs log rotation. Must be called with l.mu held.
func (l *usageLogger) rotateLocked() {
	l.file.Close()
	for i := l.maxFiles - 1; i >= 1; i-- {
		os.Rename(rotatedPath(l.path, i), rotatedPath(l.path, i+1))
	}
	os.Rename(l.path, rotatedPath(l.path, 1))
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err == nil {
		l.file = f
	}
}

func (l *usageLogger) Close() {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.file != nil {
		l.file.Close()
		l.file = nil
	}
}

func rotatedPath(base string, n int) string {
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)
	return fmt.Sprintf("%s.%d%s", stem, n, ext)
}

func readUsage(path string, limit int, cmdFilter string) []usageEntry {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if line := scanner.Text(); line != "" {
			lines = append(lines, line)
		}
	}

	filter := strings.ToLower(cmdFilter)
	var result []usageEntry
	for i := len(lines) - 1; i >= 0 && len(result) < limit; i-- {
		var e usageEntry
		if err := json.Unmarshal([]byte(lines[i]), &e); err != nil {
			continue
		}
		if filter != "" && !strings.Contains(strings.ToLower(e.Command), filter) {
			continue
		}
		result = append(result, e)
	}
	return result
}

func jsonText(v any) *mcplib.CallToolResult {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcplib.NewToolResultText(fmt.Sprintf("error serialising result: %v", err))
	}
	return mcplib.NewToolResultText(string(b))
}
