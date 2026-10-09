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
// If the lock cannot be taken (read-only home, timeout), fn runs anyway: a rare
// race is better than refusing to work.
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
			fmt.Fprintf(os.Stderr, "Warning: timed out waiting for %s; continuing without it\n", f.Name())
			return fn()
		}
		time.Sleep(50 * time.Millisecond)
	}
	defer unlock(f)
	return fn()
}
