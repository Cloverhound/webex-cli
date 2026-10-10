package auth

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
)

// HasReadOnlyToken reports whether the token store holds a read-only token for email.
func HasReadOnlyToken(email string) bool {
	tok, err := LoadToken(email)
	return err == nil && tok.ReadOnly
}

// PurgeWriteTokens deletes every stored token that is not read-only and removes
// users left without a login from cfg. It returns the removed emails, sorted.
// Read-only mode depends on this: the keyring and credentials file are readable
// by any process running as the user, so a write-capable token left in either
// would bypass the CLI's checks. Both stores are checked whatever
// $WEBEX_TOKEN_STORE selects. An entry that cannot be decoded is deleted too,
// since it may hold write access, and so is the $WEBEX_REFRESH_TOKEN cache,
// which read-only mode never uses.
func PurgeWriteTokens(cfg *appconfig.Config) ([]string, error) {
	var removed, failed []string
	for _, email := range cfg.UserEmails() {
		copies, err := readCopies(email, []string{StoreKeyring, StoreFile})
		if errors.Is(err, ErrTokenNotFound) {
			continue
		}
		if err != nil {
			failed = append(failed, email)
			continue
		}
		kept, ok := 0, true
		for _, c := range copies {
			if tok, err := decodeToken(c.data); err == nil && tok.ReadOnly {
				kept++
				continue
			}
			if err := deleteFrom(c.store, email); err != nil {
				ok = false
			}
		}
		if !ok {
			failed = append(failed, email)
			continue
		}
		if kept == 0 {
			cfg.RemoveUser(email)
			removed = append(removed, email)
		}
	}
	sort.Strings(removed)
	if len(failed) > 0 {
		sort.Strings(failed)
		return removed, fmt.Errorf("could not remove write-capable tokens for %v from the token store; remove them manually", failed)
	}
	if err := DeleteEnvRefreshCache(); err != nil {
		return removed, fmt.Errorf("could not remove cached $%s tokens from %s: %w", RefreshTokenEnv, CredentialsPath(), err)
	}
	return removed, nil
}
