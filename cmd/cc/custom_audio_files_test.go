package cc

import (
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/config"
)

// Uploads go through the shared client, so an expired token is refreshed and
// the same multipart body is sent again.
func TestAudioFileUploadRefreshesTokenAndReplaysBody(t *testing.T) {
	file := filepath.Join(t.TempDir(), "prompt.wav")
	if err := os.WriteFile(file, []byte("RIFFwave"), 0600); err != nil {
		t.Fatal(err)
	}
	var tokens, paths []string
	var bodies []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies = append(bodies, string(b))
		tokens = append(tokens, r.Header.Get("Authorization"))
		paths = append(paths, r.URL.Path)
		if len(bodies) == 1 {
			w.WriteHeader(401)
			return
		}
		mediaType, params, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		form, err := multipart.NewReader(strings.NewReader(string(b)), params["boundary"]).ReadForm(1 << 20)
		if mediaType != "multipart/form-data" || err != nil || len(form.File["audioFile"]) != 1 ||
			!strings.Contains(form.Value["audioFileInfo"][0], `"name":"prompt.wav"`) {
			t.Errorf("form = %v / %v (err=%v)", form, mediaType, err)
		}
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	token, refresher := config.Token(), config.TokenRefresher
	config.SetToken("old-token")
	config.TokenRefresher = func() (string, error) { return "new-token", nil }
	t.Cleanup(func() {
		config.SetToken(token)
		config.TokenRefresher = refresher
	})

	if err := uploadAudioFile(srv.URL, "org 1", file, "", ""); err != nil {
		t.Fatal(err)
	}
	if len(bodies) != 2 || bodies[0] != bodies[1] {
		t.Fatalf("requests=%d, body replayed unchanged=%v", len(bodies), len(bodies) == 2 && bodies[0] == bodies[1])
	}
	if tokens[1] != "Bearer new-token" || paths[1] != "/organization/org 1/audio-file" {
		t.Errorf("tokens=%v paths=%v", tokens, paths)
	}
}
