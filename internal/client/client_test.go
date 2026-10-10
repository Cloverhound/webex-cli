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
	defer config.SetReadOnly(false)

	config.SetReadOnly(true)

	blocked := []struct{ method, path string }{
		{"PUT", "/things/1"},
		{"PATCH", "/things/1"},
		{"DELETE", "/things/1"},
		{"POST", "/things"},
		{"POST", "/telephony/calls/retrieve"},
	}
	for _, b := range blocked {
		_, _, err := Do(NewRequest(srv.URL, b.method, b.path))
		if !errors.Is(err, ErrReadOnly) {
			t.Errorf("%s %s: got err %v, want ErrReadOnly", b.method, b.path, err)
		}
	}
	if hits != 0 {
		t.Fatalf("blocked requests reached the server %d times", hits)
	}

	if _, _, err := Do(NewRequest(srv.URL, "GET", "/things")); err != nil {
		t.Fatalf("GET: %v", err)
	}
	req := NewRequest(srv.URL, "POST", "/meetings/{meetingId}/registrants/query")
	req.PathParam("meetingId", "m1")
	if _, _, err := Do(req); err != nil {
		t.Fatalf("query POST: %v", err)
	}
	if hits != 2 {
		t.Fatalf("server hits = %d, want 2", hits)
	}

	config.SetReadOnly(false)
	if _, _, err := Do(NewRequest(srv.URL, "DELETE", "/things/1")); err != nil {
		t.Fatalf("DELETE outside read-only mode: %v", err)
	}
}
