package auth

import (
	"fmt"
	"sort"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
)

// HasReadOnlyToken reports whether the keyring holds a read-only token for email.
func HasReadOnlyToken(email string) bool {
	tok, err := LoadToken(email)
	return err == nil && tok.ReadOnly
}

// PurgeWriteTokens deletes every stored token that is not read-only and removes
// its user from cfg. It returns the removed emails, sorted. Read-only mode
// depends on this: the keyring is readable by any process running as the user,
// so a write-capable token left there would bypass the CLI's checks.
func PurgeWriteTokens(cfg *appconfig.Config) ([]string, error) {
	var removed, failed []string
	for _, email := range cfg.UserEmails() {
		tok, err := LoadToken(email)
		if err != nil || tok.ReadOnly {
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
