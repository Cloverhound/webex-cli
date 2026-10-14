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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
)

// messageServer records requests on a local server, never Webex. A nil
// handler replies with an empty JSON object.
func messageServer(t *testing.T, handle func(http.ResponseWriter, *http.Request, []byte)) (string, *[][]byte) {
	t.Helper()
	var bodies [][]byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies = append(bodies, b)
		if handle == nil {
			w.Write([]byte(`{}`))
			return
		}
		handle(w, r, b)
	}))
	oldToken := config.Token()
	config.SetToken("test-token")
	t.Cleanup(func() {
		srv.Close()
		config.SetToken(oldToken)
		config.SetDryRun(false)
		config.SetReadOnly(false)
		config.SetDebug(false)
	})
	return srv.URL, &bodies
}

func TestMessageCreateCommandRegisteredOnce(t *testing.T) {
	var found int
	for _, c := range messagesCmd.Commands() {
		if c.Name() == "create" && c.Flags().Lookup("body") != nil && c.Flags().Lookup("body-file") != nil {
			found++
		}
	}
	if found != 1 || len(messagesCmd.Commands()) < 2 {
		t.Fatalf("create with --body/--body-file registered %d times", found)
	}
}

func TestMessageCreateLocalAttachment(t *testing.T) {
	t.Chdir(t.TempDir())
	data := []byte("# Guide\n\x00\xffbinary bytes\n")
	if err := os.WriteFile("guide.md", data, 0600); err != nil {
		t.Fatal(err)
	}
	url, bodies := messageServer(t, func(w http.ResponseWriter, r *http.Request, body []byte) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("normal auth missing")
		}
		mediaType, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		form, err := multipart.NewReader(bytes.NewReader(body), params["boundary"]).ReadForm(1 << 20)
		if mediaType != "multipart/form-data" || err != nil {
			t.Errorf("Content-Type=%s err=%v", mediaType, err)
			return
		}
		defer form.RemoveAll()
		for k, v := range map[string]string{"toPersonEmail": "recipient@example.invalid", "markdown": "**Guide**", "parentId": "parent"} {
			if len(form.Value[k]) != 1 || form.Value[k][0] != v {
				t.Errorf("field %s = %v", k, form.Value[k])
			}
		}
		if len(form.Value) != 3 || len(form.File["files"]) != 1 || form.File["files"][0].Filename != "guide.md" {
			t.Errorf("form = %v / %v", form.Value, form.File)
			return
		}
		f, _ := form.File["files"][0].Open()
		defer f.Close()
		if got, _ := io.ReadAll(f); !bytes.Equal(got, data) {
			t.Errorf("file bytes differ: %q", got)
		}
		w.Write([]byte(`{"id":"test-message"}`))
	})
	raw := `{"toPersonEmail":"recipient@example.invalid","markdown":"**Guide**","parentId":"parent","files":["guide.md"]}`
	if err := createMessage(url, raw, ""); err != nil || len(*bodies) != 1 {
		t.Fatalf("err=%v requests=%d", err, len(*bodies))
	}
}

func TestMessageCreateBodyFilePathsResolveFromWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.WriteFile("guide.md", []byte("working directory"), 0600); err != nil {
		t.Fatal(err)
	}
	bodyPath := filepath.Join(dir, "config", "message.json")
	if err := os.Mkdir(filepath.Dir(bodyPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(bodyPath, []byte(`{"roomId":"room","files":["./guide.md"]}`), 0600); err != nil {
		t.Fatal(err)
	}
	url, bodies := messageServer(t, nil)
	if err := createMessage(url, "", bodyPath); err != nil {
		t.Fatal(err)
	}
	if len(*bodies) != 1 || !bytes.Contains((*bodies)[0], []byte("working directory")) {
		t.Fatalf("body-file attachment not resolved from working directory: %q", *bodies)
	}
}

func TestMessageCreateJSONPassthrough(t *testing.T) {
	var contentTypes []string
	url, _ := messageServer(t, func(w http.ResponseWriter, r *http.Request, body []byte) {
		contentTypes = append(contentTypes, r.Header.Get("Content-Type"))
		w.Write(body)
	})
	for _, raw := range []string{
		`{"roomId":"room","text":"Hello"}`,
		`{"roomId":"room","files":["https://example.invalid/file.pdf"],"attachments":[{"contentType":"example"}]}`,
		`{"roomId":"room","files":["http://example.invalid/file.pdf"]}`,
		`{"roomId":"room","files":["HTTPS://example.invalid/file.pdf"]}`,
		`{"roomId":"room","files":["https:///missing-host"]}`, // Webex validates URLs, as before
		`{"roomId":"room","files":[]}`,
		`{"roomId":"room","files":null}`, // no local path: sent unchanged, as before
	} {
		req, err := messageCreateRequest(url, raw)
		if err != nil {
			t.Fatalf("%s: %v", raw, err)
		}
		if got, _, err := req.Do(); err != nil || string(got) != raw {
			t.Errorf("JSON changed: %s (err=%v)", got, err)
		}
	}
	for _, ct := range contentTypes {
		if ct != "application/json" {
			t.Errorf("content type = %s", ct)
		}
	}
}

func TestMessageCreateInvalidInputsSendNothing(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.WriteFile("guide.md", []byte("guide"), 0600); err != nil {
		t.Fatal(err)
	}
	url, bodies := messageServer(t, nil)
	for _, raw := range []string{
		``, `not-json`, `[]`, `null`,
		`{"files":"guide.md"}`, `{"files":[42]}`, `{"files":[null]}`, `{"files":[""]}`, `{"files":[" "]}`,
		`{"files":["guide.md","guide.md"]}`, `{"files":["https://example.invalid/a","guide.md"]}`,
		`{"files":["missing.md"]}`, `{"files":["file:///tmp/guide.md"]}`, `{"files":["."]}`,
		`{"files":["guide.md"],"attachments":[{"type":"card"}]}`,
		`{"files":["guide.md"],"text":null}`, `{"files":["guide.md"],"text":1}`,
	} {
		if err := createMessage(url, raw, ""); err == nil {
			t.Errorf("accepted invalid body %q", raw)
		}
	}
	if err := createMessage(url, "", filepath.Join(dir, "missing.json")); err == nil {
		t.Error("accepted missing body file")
	}
	if len(*bodies) != 0 {
		t.Errorf("invalid bodies sent %d requests", len(*bodies))
	}
}

func TestMessageCreateLocalAttachmentSafetyModes(t *testing.T) {
	file := filepath.Join(t.TempDir(), "secret.md")
	if err := os.WriteFile(file, []byte("FILE_BYTES_MUST_NOT_BE_LOGGED"), 0600); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(map[string]any{"roomId": "room", "files": []string{file}})
	url, bodies := messageServer(t, nil)
	log, err := os.CreateTemp(t.TempDir(), "log")
	if err != nil {
		t.Fatal(err)
	}
	defer log.Close()
	oldStderr := os.Stderr
	os.Stderr = log
	t.Cleanup(func() { os.Stderr = oldStderr })

	config.SetDryRun(true)
	if err := createMessage(url, string(raw), ""); !errors.Is(err, client.ErrDryRun) {
		t.Errorf("dry run err=%v", err)
	}
	config.SetDryRun(false)
	config.SetReadOnly(true)
	if err := createMessage(url, string(raw), ""); !errors.Is(err, client.ErrReadOnly) {
		t.Errorf("read-only err=%v", err)
	}
	config.SetReadOnly(false)
	config.SetDebug(true)
	if err := createMessage(url, string(raw), ""); err != nil {
		t.Fatal(err)
	}
	logged, _ := os.ReadFile(log.Name())
	if strings.Contains(string(logged), "FILE_BYTES_MUST_NOT_BE_LOGGED") || !strings.Contains(string(logged), "secret.md") {
		t.Errorf("diagnostics must name the file without its bytes:\n%s", logged)
	}
	if len(*bodies) != 1 {
		t.Errorf("requests=%d, want only the debug-mode request", len(*bodies))
	}
}
