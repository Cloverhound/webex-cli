package client

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Cloverhound/webex-cli/internal/config"
)

func TestReadOnlyBlocksWritesBeforeSending(t *testing.T) {
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()
	defer config.SetReadOnly(false, false)

	config.SetReadOnly(true, true)

	for _, method := range []string{"PUT", "PATCH", "DELETE"} {
		_, _, err := Do(NewRequest(srv.URL, method, "/things/1"))
		if !errors.Is(err, ErrReadOnly) {
			t.Errorf("%s: got err %v, want ErrReadOnly", method, err)
		}
	}
	if hits != 0 {
		t.Fatalf("blocked requests reached the server %d times", hits)
	}

	if _, _, err := Do(NewRequest(srv.URL, "GET", "/things")); err != nil {
		t.Fatalf("GET: %v", err)
	}
	if _, _, err := Do(NewRequest(srv.URL, "POST", "/search")); err != nil {
		t.Fatalf("query POST: %v", err)
	}

	config.SetReadOnly(true, false)
	if _, _, err := Do(NewRequest(srv.URL, "POST", "/things")); !errors.Is(err, ErrReadOnly) {
		t.Fatalf("write POST: got err %v, want ErrReadOnly", err)
	}
	if hits != 2 {
		t.Fatalf("server hits = %d, want 2", hits)
	}
}
