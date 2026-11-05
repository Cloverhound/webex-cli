package webex

import (
	"net/http"
	"strings"
	"sync"
	"time"
)

// CCInterval spaces Contact Center requests. Contact Center refuses bursts above about
// 25 requests a second with a Retry-After that runs to the next minute boundary, often
// 20 to 30 seconds, while 16 a second has never been refused. Pacing costs milliseconds;
// one refusal stalls a whole sync.
const CCInterval = time.Second / 15

// Pacer spaces requests under one URL prefix evenly across every goroutine that shares
// it, which covers a whole snapshot or every resource an OpenTofu run refreshes.
type Pacer struct {
	Prefix   string
	Interval time.Duration

	mu   sync.Mutex
	next time.Time
}

// Wait holds a request under the prefix until its turn, and passes others at once.
// It suits webexapi.Client.BeforeSend.
func (p *Pacer) Wait(r *http.Request) error {
	if !strings.HasPrefix(r.URL.String(), p.Prefix) {
		return nil
	}
	p.mu.Lock()
	now := time.Now()
	slot := p.next
	if slot.Before(now) {
		slot = now
	}
	p.next = slot.Add(p.Interval)
	p.mu.Unlock()

	delay := time.Until(slot)
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-r.Context().Done():
		return r.Context().Err()
	}
}
