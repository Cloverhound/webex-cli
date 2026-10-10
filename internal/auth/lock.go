package auth

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var lockTimeout = 30 * time.Second

// withFileLock runs fn while holding an exclusive lock on the named file in
// CredentialsDir, so parallel CLI processes serialize token writes and refreshes.
// When locking is impossible (read-only home, no lock support), fn runs anyway,
// since no other process can hold the lock either. When another process holds
// it past the timeout, fn does not run: it could reuse a single-use refresh
// token or overwrite that process's write.
func withFileLock(name string, fn func() error) error {
	dir := CredentialsDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fn()
	}
	f, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return fn()
	}
	defer f.Close()

	deadline := time.Now().Add(lockTimeout)
	for {
		locked, err := tryLock(f)
		if err != nil {
			return fn()
		}
		if locked {
			break
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timed out after %s waiting for %s; another webex process is still using it", lockTimeout, f.Name())
		}
		time.Sleep(50 * time.Millisecond)
	}
	defer unlock(f)
	return fn()
}
