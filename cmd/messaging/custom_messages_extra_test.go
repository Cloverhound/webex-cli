package messaging

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/config"
)

func TestMessageCreateRegisteredOnce(t *testing.T) {
	var count int
	for _, c := range messagesCmd.Commands() {
		if c.Name() == "create" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("create registered %d times", count)
	}
}

func TestMessageCreateRateLimitReplaysMultipart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "guide.md")
	if err := os.WriteFile(path, []byte("guide"), 0600); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(map[string]any{"roomId": "room", "files": []string{path}})
	var bodies [][]byte
	messageServer(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies = append(bodies, b)
		if len(bodies) == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(429)
			return
		}
		w.Write([]byte(`{}`))
	})
	oldRetry := config.MaxRetry()
	config.SetMaxRetry(1)
	t.Cleanup(func() { config.SetMaxRetry(oldRetry) })
	if err := runCreate(t, "body", string(raw)); err != nil {
		t.Fatal(err)
	}
	if len(bodies) != 2 || !bytes.Equal(bodies[0], bodies[1]) {
		t.Fatal("multipart body was not replayed unchanged")
	}
}

func TestMessageCreateLocalPathsWithURLCharacters(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	path := "guide #% 日本語.md"
	if err := os.WriteFile(path, []byte("guide"), 0600); err != nil {
		t.Fatal(err)
	}
	raw, _ := json.Marshal(map[string]any{"roomId": "room", "files": []string{path}})
	messageServer(t, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Error(err)
			return
		}
		defer r.MultipartForm.RemoveAll()
		f, header, err := r.FormFile("files")
		if err != nil {
			t.Error(err)
			return
		}
		defer f.Close()
		if header.Filename != path {
			t.Errorf("filename = %q", header.Filename)
		}
		w.Write([]byte(`{}`))
	})
	if err := runCreate(t, "body", string(raw)); err != nil {
		t.Fatal(err)
	}
}
