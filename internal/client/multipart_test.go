package client

import (
	"bytes"
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

func writeTempFile(t *testing.T, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

// RFC 7578 forbids filename*= in multipart/form-data; names stay quoted UTF-8.
func TestMultipartAttachmentUsesQuotedUTF8Filename(t *testing.T) {
	for name, header := range map[string]string{
		"日本.md":             `filename="日本.md"`,
		"guide #% 日本語.md":   `filename="guide #% 日本語.md"`,
		"guide;v1.md":       `filename="guide;v1.md"`,
		`guide "quoted".md`: `filename="guide \"quoted\".md"`,
		`guide\part.md`:     `filename="guide\\part.md"`,
	} {
		if runtime.GOOS == "windows" && strings.ContainsAny(name, `"\`) {
			continue // not valid Windows filename characters
		}
		req := NewRequest("http://example.invalid", "POST", "/messages")
		if err := req.SetMultipartFile(nil, "files", writeTempFile(t, name, []byte("guide"))); err != nil {
			t.Fatal(err)
		}
		_, params, _ := mime.ParseMediaType(req.bodyContentType)
		part, err := multipart.NewReader(strings.NewReader(req.bodyRaw), params["boundary"]).NextPart()
		if err != nil {
			t.Fatal(err)
		}
		got := part.Header.Get("Content-Disposition")
		if !strings.Contains(got, header) || strings.Contains(got, "filename*=") || part.FileName() != name {
			t.Errorf("%s: non-compliant form-data filename header: %s", name, got)
		}
	}
}

func TestMultipartAttachmentRejectsUnsafeInputs(t *testing.T) {
	path := writeTempFile(t, "guide.md", []byte("guide"))
	large := writeTempFile(t, "large.bin", nil)
	if err := os.Truncate(large, MaxMessageFileSize+1); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		fieldName, path, want string
		fields                map[string]string
	}{
		{"files\r\nX-Injected: true", path, "invalid", nil},
		{"files", path, "invalid", map[string]string{"text\r\nX-Injected: true": "guide"}},
		{"files", large, "100 MiB", nil},
		{"files", filepath.Dir(path), "regular file", nil},
	} {
		req := NewRequest("http://example.invalid", "POST", "/messages")
		if err := req.SetMultipartFile(tc.fields, tc.fieldName, tc.path); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Errorf("%q / %v / %s: err=%v, want %q", tc.fieldName, tc.fields, tc.path, err, tc.want)
		}
	}
}

func TestMultipartAttachmentAuthenticationRetryReplaysBody(t *testing.T) {
	data := []byte("one unchanged attachment\x00\xff")
	var bodies [][]byte
	var tokens []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies = append(bodies, b)
		tokens = append(tokens, r.Header.Get("Authorization"))
		if len(bodies) == 1 {
			w.WriteHeader(401)
			return
		}
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	token, refresher := config.Token(), config.TokenRefresher
	config.SetToken("old-token")
	config.TokenRefresher = func() (string, error) { return "new-token", nil }
	t.Cleanup(func() { config.SetToken(token); config.TokenRefresher = refresher })

	req := NewRequest(srv.URL, "POST", "/messages")
	if err := req.SetMultipartFile(map[string]string{"roomId": "room"}, "files", writeTempFile(t, "guide.bin", data)); err != nil {
		t.Fatal(err)
	}
	if _, _, err := req.Do(); err != nil {
		t.Fatal(err)
	}
	if len(bodies) != 2 || !bytes.Equal(bodies[0], bodies[1]) || !bytes.Contains(bodies[1], data) {
		t.Fatal("body not replayed unchanged")
	}
	if tokens[0] != "Bearer old-token" || tokens[1] != "Bearer new-token" {
		t.Fatalf("tokens=%v", tokens)
	}
}

func TestSetMultipartWritesTypedPartsAndSummarizesFiles(t *testing.T) {
	data := []byte("RIFF\x00\xffWAVE_BYTES_MUST_NOT_BE_LOGGED")
	req := NewRequest("http://example.invalid", "POST", "/audio-file")
	err := req.SetMultipart([]MultipartPart{
		{FieldName: "audioFileInfo", ContentType: "application/json", Data: []byte(`{"name":"prompt.wav"}`)},
		{FieldName: "audioFile", FileName: "prompt.wav", ContentType: "audio/wav", Data: data},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, params, _ := mime.ParseMediaType(req.bodyContentType)
	reader := multipart.NewReader(strings.NewReader(req.bodyRaw), params["boundary"])
	for _, want := range []struct{ field, file, contentType, body string }{
		{"audioFileInfo", "", "application/json", `{"name":"prompt.wav"}`},
		{"audioFile", "prompt.wav", "audio/wav", string(data)},
	} {
		part, err := reader.NextPart()
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(part)
		if part.FormName() != want.field || part.FileName() != want.file ||
			part.Header.Get("Content-Type") != want.contentType || string(body) != want.body {
			t.Errorf("part %s = %v %q", want.field, part.Header, body)
		}
	}
	if strings.Contains(req.bodySummary, "WAVE_BYTES") || !strings.Contains(req.bodySummary, "prompt.wav") ||
		!strings.Contains(req.bodySummary, `{\"name\":\"prompt.wav\"}`) {
		t.Errorf("summary must show text parts and name files without their bytes: %s", req.bodySummary)
	}
}

// Diagnostics must be stable across runs even though fields arrive as a map.
func TestSetMultipartFileSummarizesPartsInFieldOrder(t *testing.T) {
	path := writeTempFile(t, "guide.md", []byte("guide"))
	want := `[{"field":"markdown","value":"**hi**"},{"field":"roomId","value":"room"},{"field":"text","value":"hi"},` +
		`{"field":"files","file":"guide.md","contentType":"text/markdown; charset=utf-8","size":5}]`
	for i := 0; i < 20; i++ {
		req := NewRequest("http://example.invalid", "POST", "/messages")
		fields := map[string]string{"roomId": "room", "text": "hi", "markdown": "**hi**"}
		if err := req.SetMultipartFile(fields, "files", path); err != nil {
			t.Fatal(err)
		}
		if req.bodySummary != want {
			t.Fatalf("summary = %s", req.bodySummary)
		}
	}
}

func TestSetMultipartRejectsHeaderInjection(t *testing.T) {
	for _, part := range []MultipartPart{
		{FieldName: ""},
		{FieldName: "file\r\nX-Injected: true"},
		{FieldName: "file", FileName: "a\nb.wav"},
		{FieldName: "file", ContentType: "audio/wav\r\nX-Injected: true"},
	} {
		req := NewRequest("http://example.invalid", "POST", "/upload")
		if err := req.SetMultipart([]MultipartPart{part}); err == nil {
			t.Errorf("accepted %+v", part)
		}
	}
}

// SetBodyRaw replaces the whole body, including multipart content type and summary.
func TestSetBodyRawAfterMultipartRestoresJSON(t *testing.T) {
	req := NewRequest("http://example.invalid", "POST", "/messages")
	if err := req.SetMultipartFile(nil, "files", writeTempFile(t, "guide.md", []byte("guide"))); err != nil {
		t.Fatal(err)
	}
	req.SetBodyRaw(`{"text":"new JSON"}`)
	if req.bodyContentType != "" || req.bodySummary != "" {
		t.Fatal("stale multipart state retained")
	}
}
