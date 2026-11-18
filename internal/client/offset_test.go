package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// TestOffsetPagingAppliesOutsideCalling covers AI Agent Studio lists, which are
// Contact Center commands that page by limit/offset under "objects".
func TestOffsetPagingAppliesOutsideCalling(t *testing.T) {
	var offsets []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		offsets = append(offsets, q.Get("offset"))
		offset, _ := strconv.Atoi(q.Get("offset"))
		limit, _ := strconv.Atoi(q.Get("limit"))
		var objects []map[string]string
		for i := offset; i < offset+limit && i < 150; i++ {
			objects = append(objects, map[string]string{"name": fmt.Sprintf("a%d", i)})
		}
		json.NewEncoder(w).Encode(map[string]any{
			"meta":    map[string]any{"limit": limit, "offset": offset, "total_count": 150, "next": nil},
			"objects": objects,
		})
	}))
	defer srv.Close()

	req := NewRequest(srv.URL, "GET", "/api/v2/organizations/o/agents/")
	body, _, err := req.DoPaginated(Paging{Mode: Offset, Start: "offset", Size: "limit", Items: "objects", TotalItems: "meta.total_count"})
	if err != nil {
		t.Fatal(err)
	}
	var items []json.RawMessage
	if err := json.Unmarshal(body, &items); err != nil || len(items) != 150 {
		t.Fatalf("got %d items (err %v), want 150", len(items), err)
	}
	if fmt.Sprint(offsets) != "[0 100]" {
		t.Fatalf("offsets %v, want [0 100]", offsets)
	}
}

func TestDebugOutputHidesAgentTokens(t *testing.T) {
	if got := redactHeader("Bot-Access-Token", "secret-value"); got != "<12 chars>" {
		t.Fatalf("redactHeader = %q", got)
	}
	if got := redactHeader("Accept", "application/json"); got != "application/json" {
		t.Fatalf("redactHeader changed a plain header: %q", got)
	}
}
