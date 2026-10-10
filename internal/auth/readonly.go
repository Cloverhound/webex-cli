package auth

import (
	"errors"
	"fmt"
	"sort"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/zalando/go-keyring"
)

// HasReadOnlyToken reports whether the keyring holds a read-only token for email.
func HasReadOnlyToken(email string) bool {
	tok, err := LoadToken(email)
	return err == nil && tok.ReadOnly
}

// PurgeWriteTokens deletes every stored token that is not read-only and removes
// its user from cfg. It returns the removed emails, sorted. Read-only mode
// depends on this: the keyring is readable by any process running as the user,
// so a write-capable token left there would bypass the CLI's checks. An entry
// that cannot be read or decoded is deleted too, since it may hold write access.
func PurgeWriteTokens(cfg *appconfig.Config) ([]string, error) {
	var removed, failed []string
	for _, email := range cfg.UserEmails() {
		tok, err := LoadToken(email)
		if errors.Is(err, keyring.ErrNotFound) || (err == nil && tok.ReadOnly) {
			continue
		}
		if err := DeleteToken(email); err != nil {
			failed = append(failed, email)
			continue
		}
		cfg.RemoveUser(email)
		removed = append(removed, email)
	}
	sort.Strings(removed)
	if len(failed) > 0 {
		sort.Strings(failed)
		return removed, fmt.Errorf("could not remove write-capable tokens for %v from the keyring; remove them manually", failed)
	}
	return removed, nil
}
