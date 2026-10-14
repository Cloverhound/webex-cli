package messaging

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
)

type messageTransport func(*http.Request) (*http.Response, error)

func (f messageTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// Redirect the real command's HTTP requests to a local server, never Webex.
func messageServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	target, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	old := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: messageTransport(func(r *http.Request) (*http.Response, error) {
		clone := r.Clone(r.Context())
		clone.URL.Scheme = target.Scheme
		clone.URL.Host = target.Host
		return http.DefaultTransport.RoundTrip(clone)
	})}
	oldToken := config.Token()
	config.SetToken("test-token")
	t.Cleanup(func() {
		http.DefaultClient = old
		config.SetToken(oldToken)
		config.SetDryRun(false)
		config.SetReadOnly(false)
		config.SetDebug(false)
	})
}

func runCreate(t *testing.T, flag, body string) error {
	t.Helper()
	c, _, err := messagesCmd.Find([]string{"create"})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"body", "body-file"} {
		if err := c.Flags().Set(name, ""); err != nil {
			t.Fatal(err)
		}
	}
	if err := c.Flags().Set(flag, body); err != nil {
		t.Fatal(err)
	}
	return c.RunE(c, nil)
}

func TestMessageCreateLocalAttachment(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	name := `guide "quoted".md`
	if runtime.GOOS == "windows" {
		name = "guide unicode.md"
	} // Quotes are not valid Windows filenames.
	data := []byte("# Guide\n\x00\xffbinary bytes\n")
	if err := os.WriteFile(name, data, 0600); err != nil {
		t.Fatal(err)
	}
	var hits int
	messageServer(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Error(err)
			return
		}
		if mediaType != "multipart/form-data" {
			t.Errorf("Content-Type = %s, want multipart/form-data", mediaType)
			return
		}
		form, err := multipart.NewReader(r.Body, params["boundary"]).ReadForm(1 << 20)
		if err != nil {
			t.Error(err)
			return
		}
		defer form.RemoveAll()
		for k, v := range map[string]string{"toPersonEmail": "recipient@example.invalid", "text": "Guide", "markdown": "**Guide**", "parentId": "parent"} {
			if len(form.Value[k]) != 1 || form.Value[k][0] != v {
				t.Errorf("field %s = %v", k, form.Value[k])
			}
		}
		files := form.File["files"]
		if len(files) != 1 {
			t.Errorf("files = %v", files)
			return
		}
		if files[0].Filename != name {
			t.Errorf("filename = %q", files[0].Filename)
		}
		f, err := files[0].Open()
		if err != nil {
			t.Error(err)
			return
		}
		defer f.Close()
		got, err := io.ReadAll(f)
		if err != nil {
			t.Error(err)
		}
		if !bytes.Equal(got, data) {
			t.Errorf("file bytes differ: %q", got)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("normal auth missing")
		}
		w.Write([]byte(`{"id":"test-message"}`))
	})
	raw, _ := json.Marshal(map[string]any{"toPersonEmail": "recipient@example.invalid", "text": "Guide", "markdown": "**Guide**", "parentId": "parent", "files": []string{name}})
	if err := runCreate(t, "body", string(raw)); err != nil {
		t.Fatal(err)
	}
	if hits != 1 {
		t.Fatalf("hits = %d", hits)
	}
}

func TestMessageCreateBodyFilePathsResolveFromWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.WriteFile("guide.md", []byte("working directory"), 0600); err != nil {
		t.Fatal(err)
	}
	bodyDir := filepath.Join(dir, "config")
	if err := os.Mkdir(bodyDir, 0700); err != nil {
		t.Fatal(err)
	}
	bodyPath := filepath.Join(bodyDir, "message.json")
	if err := os.WriteFile(bodyPath, []byte(`{"roomId":"room","files":["./guide.md"]}`), 0600); err != nil {
		t.Fatal(err)
	}
	messageServer(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Error(err)
			return
		}
		defer r.MultipartForm.RemoveAll()
		f, _, err := r.FormFile("files")
		if err != nil {
			t.Error(err)
			return
		}
		defer f.Close()
		b, _ := io.ReadAll(f)
		if string(b) != "working directory" {
			t.Errorf("bytes=%q", b)
		}
		w.Write([]byte(`{}`))
	})
	if err := runCreate(t, "body-file", bodyPath); err != nil {
		t.Fatal(err)
	}
}

func TestMessageCreateJSONCompatibility(t *testing.T) {
	var got []byte
	messageServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("content type = %s", r.Header.Get("Content-Type"))
		}
		got, _ = io.ReadAll(r.Body)
		w.Write([]byte(`{}`))
	})
	for _, raw := range []string{
		`{"roomId":"room","text":"Hello"}`,
		`{"roomId":"room","files":["https://example.invalid/file.pdf"],"attachments":[{"contentType":"example"}]}`,
		`{"roomId":"room","files":["http://example.invalid/file.pdf"]}`,
		`{"roomId":"room","files":[]}`,
	} {
		if err := runCreate(t, "body", raw); err != nil {
			t.Fatal(err)
		}
		if string(got) != raw {
			t.Errorf("JSON changed: %s", got)
		}
	}
}

func TestMessageCreateInvalidAttachmentsFailBeforeSending(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "guide.md")
	if err := os.WriteFile(file, []byte("guide"), 0600); err != nil {
		t.Fatal(err)
	}
	var hits int
	messageServer(t, func(w http.ResponseWriter, r *http.Request) { hits++; w.Write([]byte(`{}`)) })
	inputs := []string{
		`{"files": "guide.md"}`, `{"files":null}`, `{"files":[null]}`, `{"files":[42]}`, `{"files":[""]}`,
		`{"files":["a.md","b.md"]}`, `{"files":["https://example.invalid/a","b.md"]}`,
		`{"files":["file:///tmp/guide.md"]}`, `{"files":["ftp://example.invalid/file"]}`,
		`{"files":["missing-file.md"]}`, `{"files":["https:///missing-host"]}`,
		`not-json`, `[]`, `null`,
	}
	extra, _ := json.Marshal(map[string]any{"files": []string{dir}})
	inputs = append(inputs, string(extra))
	extra, _ = json.Marshal(map[string]any{"files": []string{file}, "attachments": []any{map[string]any{"type": "card"}}})
	inputs = append(inputs, string(extra))
	for _, raw := range inputs {
		if err := runCreate(t, "body", raw); err == nil {
			t.Errorf("accepted invalid body %s", raw)
		}
	}
	if hits != 0 {
		t.Errorf("invalid bodies sent %d requests", hits)
	}
}

func TestMessageCreateLocalAttachmentSafetyModes(t *testing.T) {
	file := filepath.Join(t.TempDir(), "secret.md")
	if err := os.WriteFile(file, []byte("FILE_BYTES_MUST_NOT_BE_LOGGED"), 0600); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(map[string]any{"roomId": "room", "files": []string{file}})
	var hits int
	messageServer(t, func(w http.ResponseWriter, r *http.Request) { hits++; w.Write([]byte(`{}`)) })
	oldStderr := os.Stderr
	log, err := os.CreateTemp(t.TempDir(), "log")
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	os.Stderr = log
	t.Cleanup(func() { os.Stderr = oldStderr })
	config.SetDryRun(true)
	if err := runCreate(t, "body", string(raw)); !errors.Is(err, client.ErrDryRun) {
		t.Errorf("dry run err=%v", err)
	}
	config.SetDryRun(false)
	config.SetReadOnly(true)
	if err := runCreate(t, "body", string(raw)); !errors.Is(err, client.ErrReadOnly) {
		t.Errorf("read-only err=%v", err)
	}
	config.SetReadOnly(false)
	config.SetDebug(true)
	if err := runCreate(t, "body", string(raw)); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(log.Name())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "FILE_BYTES_MUST_NOT_BE_LOGGED") {
		t.Error("file bytes leaked into logs")
	}
	if !strings.Contains(string(b), "secret.md") {
		t.Error("preview does not identify the file")
	}
	if hits != 1 {
		t.Errorf("hits=%d, want only debug-mode request", hits)
	}
}
