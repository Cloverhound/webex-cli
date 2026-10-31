package auth

import (
	"testing"
	"time"
)

func TestWithFileLockTimesOutWithoutRunning(t *testing.T) {
	old := lockTimeout
	lockTimeout = 100 * time.Millisecond
	t.Cleanup(func() { lockTimeout = old })

	held, release, done := make(chan struct{}), make(chan struct{}), make(chan struct{})
	go func() {
		defer close(done)
		_ = withFileLock("test.lock", func() error {
			close(held)
			<-release
			return nil
		})
	}()
	<-held

	ran := false
	err := withFileLock("test.lock", func() error {
		ran = true
		return nil
	})
	close(release)
	<-done

	if err == nil {
		t.Error("expected a timeout error while another holder has the lock")
	}
	if ran {
		t.Error("callback ran without the lock")
	}
}
