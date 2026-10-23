package auth

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/zalando/go-keyring"
)

// HasReadOnlyToken reports whether the token store holds a read-only token for email.
func HasReadOnlyToken(email string) bool {
	tok, err := LoadToken(email)
	return err == nil && tok.ReadOnly
}

// PurgeWriteTokens deletes every stored token that is not read-only and removes
// cfg's users who lost their only login. It returns the removed emails, sorted.
// Read-only mode depends on this: the keyring and credentials file are readable
// by any process running as the user, so a write-capable token left in either
// would bypass the CLI's checks. Every entry in both stores is checked, not only
// cfg's users, whatever $WEBEX_TOKEN_STORE selects. That includes entries that
// cannot be decoded, since they may hold write access, and the
// $WEBEX_REFRESH_TOKEN cache, which read-only mode never uses. A keyring that
// cannot be read is skipped, because no other process can read it either.
//
// It holds refresh.lock so that a refresh already in progress saves its token
// before the purge runs, not after.
func PurgeWriteTokens(cfg *appconfig.Config) ([]string, error) {
	var removed []string
	err := withFileLock("refresh.lock", func() error {
		keptKeyring, deletedKeyring, keyringRead, err := purgeKeyring(cfg.UserEmails())
		if err != nil {
			return err
		}
		keptFile, deletedFile, err := purgeFile()
		if err != nil {
			return err
		}
		if err := DeleteEnvRefreshCache(); err != nil {
			return fmt.Errorf("could not remove cached $%s tokens from %s: %w", RefreshTokenEnv, CredentialsPath(), err)
		}
		if !keyringRead && keyringExpected(runtime.GOOS, os.Getenv, fmt.Sprintf("/run/user/%d", os.Getuid())) {
			fmt.Fprintln(os.Stderr, "Warning: the keyring could not be read, so any write-capable logins stored there were not removed")
		}
		for _, email := range cfg.UserEmails() {
			if (deletedKeyring[email] || deletedFile[email]) && !keptKeyring[email] && !keptFile[email] {
				cfg.RemoveUser(email)
				removed = append(removed, email)
			}
		}
		return nil
	})
	sort.Strings(removed)
	return removed, err
}

// purgeKeyring keeps the read-only keyring tokens of emails and deletes every
// other webex-cli keyring entry. go-keyring cannot list entries, so it clears
// the service and saves the kept tokens again. It reports read=false, and
// changes nothing, when the keyring cannot be read.
func purgeKeyring(emails []string) (kept, deleted map[string]bool, read bool, err error) {
	kept, deleted = map[string]bool{}, map[string]bool{}
	keep := map[string]string{}
	for _, email := range emails {
		data, err := keyring.Get(serviceName, email)
		if errors.Is(err, keyring.ErrNotFound) {
			read = true
			continue
		}
		if err != nil {
			return kept, deleted, false, nil
		}
		read = true
		if tok, err := decodeToken(data); err == nil && tok.ReadOnly {
			keep[email] = data
		} else {
			deleted[email] = true
		}
	}
	if err := keyring.DeleteAll(serviceName); err != nil {
		if !read {
			return kept, deleted, false, nil
		}
		return kept, deleted, true, fmt.Errorf("clearing write-capable logins from the keyring: %w", err)
	}
	for email, data := range keep {
		if err := keyring.Set(serviceName, email, data); err != nil {
			return kept, deleted, true, fmt.Errorf("restoring the read-only login for %s to the keyring (log in again with --read-only): %w", email, err)
		}
		kept[email] = true
	}
	return kept, deleted, true, nil
}

// purgeFile deletes every credentials-file entry that is not a read-only token.
func purgeFile() (kept, deleted map[string]bool, err error) {
	kept, deleted = map[string]bool{}, map[string]bool{}
	if _, err := os.Stat(CredentialsPath()); os.IsNotExist(err) {
		return kept, deleted, nil
	}
	err = withFileLock("credentials.lock", func() error {
		cf, err := readCredentials()
		if err != nil {
			return err
		}
		for key, data := range cf.Tokens {
			if tok, err := decodeToken(data); err == nil && tok.ReadOnly {
				kept[key] = true
				continue
			}
			delete(cf.Tokens, key)
			deleted[key] = true
		}
		if len(deleted) == 0 {
			return nil
		}
		return writeCredentials(cf)
	})
	return kept, deleted, err
}

// keyringExpected reports whether this host should have a keyring, so that one
// that cannot be read deserves a warning. Linux containers and servers without
// a D-Bus session bus have none; runDir is /run/user/<uid>, where systemd puts
// the bus that godbus finds.
func keyringExpected(goos string, getenv func(string) string, runDir string) bool {
	if strings.ToLower(strings.TrimSpace(getenv(TokenStoreEnv))) == StoreFile {
		return false
	}
	if goos == "darwin" || goos == "windows" {
		return true
	}
	if getenv("DBUS_SESSION_BUS_ADDRESS") != "" {
		return true
	}
	for _, name := range []string{"bus", "dbus-session"} {
		if _, err := os.Stat(filepath.Join(runDir, name)); err == nil {
			return true
		}
	}
	return false
}
