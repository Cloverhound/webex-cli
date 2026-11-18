package webexapi

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMultipartFileWithTheLargestLimitKeepsTheFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a.txt")
	if err := os.WriteFile(path, []byte("file contents"), 0o600); err != nil {
		t.Fatal(err)
	}
	body, err := MultipartFile(nil, "file", path, math.MaxInt64)
	if err != nil || !strings.Contains(string(body.Body), "file contents") {
		t.Fatalf("body %q, err %v", body.Body, err)
	}
}
