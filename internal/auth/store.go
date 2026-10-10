package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zalando/go-keyring"
)

// TokenStoreEnv forces a token store: "keyring" or "file". Unset, the keyring
// is used and the credentials file is the fallback when the keyring fails.
const TokenStoreEnv = "WEBEX_TOKEN_STORE"

const (
	StoreKeyring = "keyring"
	StoreFile    = "file"
)

// ErrTokenNotFound means no store holds a token for the key.
var ErrTokenNotFound = errors.New("token not found")

// forcedStore returns the store named by $WEBEX_TOKEN_STORE, or "" for automatic selection.
func forcedStore() (string, error) {
	switch v := strings.ToLower(strings.TrimSpace(os.Getenv(TokenStoreEnv))); v {
	case "", "auto":
		return "", nil
	case StoreKeyring, StoreFile:
		return v, nil
	default:
		return "", fmt.Errorf("$%s must be %q or %q, not %q", TokenStoreEnv, StoreKeyring, StoreFile, v)
	}
}

// selectedStores returns the stores $WEBEX_TOKEN_STORE allows, keyring first.
func selectedStores() ([]string, error) {
	forced, err := forcedStore()
	if err != nil {
		return nil, err
	}
	if forced != "" {
		return []string{forced}, nil
	}
	return []string{StoreKeyring, StoreFile}, nil
}

// storedCopy is the value one store holds for a key.
type storedCopy struct {
	store string
	data  string
}

// readCopies returns the copies of key held by each of stores. A missing entry
// is not an error. A keyring that fails for any other reason may still hold the
// key, so its error is returned when no store produced a copy: reporting "not
// found" would be a guess, and read-only purging relies on that.
func readCopies(key string, stores []string) ([]storedCopy, error) {
	var copies []storedCopy
	var kerr error
	for _, store := range stores {
		var data string
		var err error
		if store == StoreKeyring {
			data, err = keyring.Get(serviceName, key)
			if errors.Is(err, keyring.ErrNotFound) {
				continue
			}
			if err != nil {
				kerr = err
				continue
			}
		} else {
			data, err = fileGet(key)
			if errors.Is(err, ErrTokenNotFound) {
				continue
			}
			if err != nil {
				return nil, err
			}
		}
		copies = append(copies, storedCopy{store: store, data: data})
	}
	if len(copies) == 0 {
		if kerr != nil {
			return nil, fmt.Errorf("the keyring could not be read: %w", kerr)
		}
		return nil, ErrTokenNotFound
	}
	return copies, nil
}

// deleteFrom removes key from one store.
func deleteFrom(store, key string) error {
	if store == StoreFile {
		return fileDelete(key)
	}
	if err := keyring.Delete(serviceName, key); err != nil && !errors.Is(err, keyring.ErrNotFound) {
		return err
	}
	return nil
}

// storeSet saves data and returns the store that now holds it.
func storeSet(key, data string) (string, error) {
	forced, err := forcedStore()
	if err != nil {
		return "", err
	}
	if forced == StoreFile {
		return StoreFile, fileSetUser(key, data)
	}
	kerr := keyring.Set(serviceName, key, data)
	if kerr == nil {
		// A copy left in the file from an earlier fallback would shadow nothing
		// but would keep a stale secret on disk.
		_ = fileDelete(key)
		return StoreKeyring, nil
	}
	if forced == StoreKeyring {
		return "", kerr
	}
	if err := fileSetUser(key, data); err != nil {
		return "", fmt.Errorf("keyring unavailable (%v) and file store failed: %w", kerr, err)
	}
	// An older keyring entry would be a stale secret. Loading prefers the newest
	// copy, so this matters only for cleanup and may fail while the keyring is down.
	_ = keyring.Delete(serviceName, key)
	return StoreFile, nil
}

func storeDelete(key string) error {
	forced, err := forcedStore()
	if err != nil {
		return err
	}
	var kerr error
	if forced != StoreFile {
		if kerr = keyring.Delete(serviceName, key); errors.Is(kerr, keyring.ErrNotFound) {
			kerr = nil
		}
		if forced == StoreKeyring {
			return kerr
		}
	}
	hadFile := false
	if _, err := fileGet(key); err == nil {
		hadFile = true
	}
	if err := fileDelete(key); err != nil {
		return err
	}
	// Without a keyring, every delete fails there; that only matters when the
	// token was not in the file.
	if hadFile {
		return nil
	}
	return kerr
}

// CredentialsDir returns the directory for the credentials file and lock files:
// $XDG_CONFIG_HOME/webex-cli, %APPDATA%\webex-cli on Windows, or ~/.config/webex-cli.
func CredentialsDir() string {
	if d := os.Getenv("XDG_CONFIG_HOME"); d != "" {
		return filepath.Join(d, "webex-cli")
	}
	if runtime.GOOS == "windows" {
		if d := os.Getenv("APPDATA"); d != "" {
			return filepath.Join(d, "webex-cli")
		}
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "webex-cli")
}

// CredentialsPath returns the path of the plain-text token file.
func CredentialsPath() string {
	return filepath.Join(CredentialsDir(), "credentials.json")
}

type credentialsFile struct {
	Tokens map[string]string `json:"tokens"`
}

func readCredentials() (*credentialsFile, error) {
	cf := &credentialsFile{Tokens: map[string]string{}}
	data, err := os.ReadFile(CredentialsPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cf, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cf); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", CredentialsPath(), err)
	}
	if cf.Tokens == nil {
		cf.Tokens = map[string]string{}
	}
	return cf, nil
}

func fileGet(key string) (string, error) {
	cf, err := readCredentials()
	if err != nil {
		return "", err
	}
	data, ok := cf.Tokens[key]
	if !ok {
		return "", ErrTokenNotFound
	}
	return data, nil
}

// fileDeletePrefix removes every file entry whose key starts with prefix.
func fileDeletePrefix(prefix string) error {
	if _, err := os.Stat(CredentialsPath()); os.IsNotExist(err) {
		return nil
	}
	return withFileLock("credentials.lock", func() error {
		cf, err := readCredentials()
		if err != nil {
			return err
		}
		n := len(cf.Tokens)
		for key := range cf.Tokens {
			if strings.HasPrefix(key, prefix) {
				delete(cf.Tokens, key)
			}
		}
		if len(cf.Tokens) == n {
			return nil
		}
		return writeCredentials(cf)
	})
}

func fileSet(key, data string) error {
	return withFileLock("credentials.lock", func() error {
		cf, err := readCredentials()
		if err != nil {
			return err
		}
		cf.Tokens[key] = data
		return writeCredentials(cf)
	})
}

// fileSetUser saves a login to the file and warns the first time the file is
// created, since it holds tokens in plain text.
func fileSetUser(key, data string) error {
	_, statErr := os.Stat(CredentialsPath())
	if err := fileSet(key, data); err != nil {
		return err
	}
	if os.IsNotExist(statErr) {
		fmt.Fprintf(os.Stderr, "Warning: tokens are stored in plain text at %s (mode 0600)\n", CredentialsPath())
	}
	return nil
}

func fileDelete(key string) error {
	if _, err := os.Stat(CredentialsPath()); os.IsNotExist(err) {
		return nil
	}
	return withFileLock("credentials.lock", func() error {
		cf, err := readCredentials()
		if err != nil {
			return err
		}
		if _, ok := cf.Tokens[key]; !ok {
			return nil
		}
		delete(cf.Tokens, key)
		return writeCredentials(cf)
	})
}

// writeCredentials replaces the file through a rename so a reader never sees a
// partial write.
func writeCredentials(cf *credentialsFile) error {
	dir := CredentialsDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}
	data, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, "credentials-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	if err := tmp.Chmod(0600); err != nil && runtime.GOOS != "windows" {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, CredentialsPath()); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return nil
}
