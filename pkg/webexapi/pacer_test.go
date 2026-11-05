package webexapi

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestPacerSpacesRequestsToItsHostOnly(t *testing.T) {
	p := &Pacer{Prefix: "https://cc.example/", Interval: 20 * time.Millisecond}
	request := func(host string) *http.Request {
		r, _ := http.NewRequest(http.MethodGet, "https://"+host+"/x", nil)
		return r
	}

	start := time.Now()
	var wg sync.WaitGroup
	for range 6 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p.Wait(request("cc.example")); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if elapsed := time.Since(start); elapsed < 100*time.Millisecond {
		t.Fatalf("6 paced requests took %v, want at least 5 intervals", elapsed)
	}

	start = time.Now()
	for range 6 {
		if err := p.Wait(request("calling.example")); err != nil {
			t.Fatal(err)
		}
	}
	if elapsed := time.Since(start); elapsed > 10*time.Millisecond {
		t.Fatalf("unpaced host waited %v", elapsed)
	}
}

func TestPacerGivesUpWhenTheRequestIsCancelled(t *testing.T) {
	p := &Pacer{Prefix: "https://cc.example/", Interval: time.Hour}
	ctx, cancel := context.WithCancel(context.Background())
	first, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://cc.example/x", nil)
	if err := p.Wait(first); err != nil {
		t.Fatal(err)
	}
	cancel()
	second, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://cc.example/x", nil)
	if err := p.Wait(second); err == nil {
		t.Fatal("a cancelled request waited for its turn")
	}
}

func TestNilOrZeroPacerDoesNotWait(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "https://cc.example/x", nil)
	var p *Pacer
	if err := p.Wait(r); err != nil {
		t.Fatal(err)
	}
	if err := (&Pacer{Prefix: "https://cc.example/"}).Wait(r); err != nil {
		t.Fatal(err)
	}
}
