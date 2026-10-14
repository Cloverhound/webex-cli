package client

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/config"
)

func TestMultipartAttachmentUsesQuotedUTF8Filename(t *testing.T) {
	for _, name := range []string{"日本.md", "guide;v1.md", `guide "quoted".md`, `guide\part.md`} {
		t.Run(name, func(t *testing.T) {
			if runtime.GOOS == "windows" && strings.ContainsAny(name, `"\`) {
				t.Skip("quotes and backslashes are not valid Windows filename characters")
			}
			path := filepath.Join(t.TempDir(), name)
			if err := os.WriteFile(path, []byte("guide"), 0600); err != nil {
				t.Fatal(err)
			}
			req := NewRequest("http://example.invalid", "POST", "/messages")
			if err := req.SetMultipartFile(nil, "files", path); err != nil {
				t.Fatal(err)
			}
			_, params, err := mime.ParseMediaType(req.bodyContentType)
			if err != nil {
				t.Fatal(err)
			}
			part, err := multipart.NewReader(strings.NewReader(req.bodyRaw), params["boundary"]).NextPart()
			if err != nil {
				t.Fatal(err)
			}
			defer part.Close()
			header := part.Header.Get("Content-Disposition")
			escaped := strings.NewReplacer(`\`, `\\`, `"`, `\"`).Replace(name)
			if !strings.Contains(header, `filename="`+escaped+`"`) || strings.Contains(header, "filename*=") {
				t.Fatalf("non-compliant form-data filename header: %s", header)
			}
		})
	}
}

func TestMultipartAttachmentRejectsOversizeFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "large.bin")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(MaxMessageFileSize + 1); err != nil {
		f.Close()
		t.Fatal(err)
	}
	f.Close()
	req := NewRequest("http://example.invalid", "POST", "/messages")
	if err := req.SetMultipartFile(nil, "files", path); err == nil {
		t.Fatal("oversize file accepted")
	}
}

func TestMultipartAttachmentRejectsHeaderInjection(t *testing.T) {
	path := filepath.Join(t.TempDir(), "guide.md")
	if err := os.WriteFile(path, []byte("guide"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		fieldName string
		fields    map[string]string
	}{
		{"files\r\nX-Injected: true", nil},
		{"files", map[string]string{"text\r\nX-Injected: true": "guide"}},
	} {
		req := NewRequest("http://example.invalid", "POST", "/messages")
		if err := req.SetMultipartFile(tc.fields, tc.fieldName, path); err == nil {
			t.Errorf("accepted unsafe field name %q / %v", tc.fieldName, tc.fields)
		}
	}
}

func TestMultipartAttachmentAuthenticationRetryReplaysBody(t *testing.T) {
	data := []byte("one unchanged attachment\x00\xff")
	path := filepath.Join(t.TempDir(), "guide.bin")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	var bodies [][]byte
	var tokens []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		bodies = append(bodies, b)
		tokens = append(tokens, r.Header.Get("Authorization"))
		if len(bodies) == 1 {
			w.WriteHeader(401)
			return
		}
		_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Error(err)
			return
		}
		form, err := multipart.NewReader(bytes.NewReader(b), params["boundary"]).ReadForm(1 << 20)
		if err != nil {
			t.Error(err)
			return
		}
		defer form.RemoveAll()
		f, err := form.File["files"][0].Open()
		if err != nil {
			t.Error(err)
			return
		}
		defer f.Close()
		got, _ := io.ReadAll(f)
		if !bytes.Equal(got, data) {
			t.Error("retried file content differs")
		}
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	token := config.Token()
	refresher := config.TokenRefresher
	config.SetToken("old-token")
	config.TokenRefresher = func() (string, error) { return "new-token", nil }
	t.Cleanup(func() { config.SetToken(token); config.TokenRefresher = refresher })
	req := NewRequest(srv.URL, "POST", "/messages")
	if err := req.SetMultipartFile(map[string]string{"roomId": "room"}, "files", path); err != nil {
		t.Fatal(err)
	}
	if _, _, err := req.Do(); err != nil {
		t.Fatal(err)
	}
	if len(bodies) != 2 || !bytes.Equal(bodies[0], bodies[1]) {
		t.Fatal("body not replayed unchanged")
	}
	if tokens[0] != "Bearer old-token" || tokens[1] != "Bearer new-token" {
		t.Fatalf("tokens=%v", tokens)
	}
}

func TestSetBodyRawAfterMultipartRestoresJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "guide.md")
	if err := os.WriteFile(path, []byte("guide"), 0600); err != nil {
		t.Fatal(err)
	}
	req := NewRequest("http://example.invalid", "POST", "/messages")
	if err := req.SetMultipartFile(nil, "files", path); err != nil {
		t.Fatal(err)
	}
	req.SetBodyRaw(`{"text":"new JSON"}`)
	if req.bodyContentType != "" || req.bodySummary != "" {
		t.Fatal("stale multipart state retained")
	}
	config.SetDryRun(true)
	defer config.SetDryRun(false)
	if _, _, err := req.Do(); !errors.Is(err, ErrDryRun) {
		t.Fatal(err)
	}
}
