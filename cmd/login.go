package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
	"github.com/Cloverhound/webex-cli/internal/localconfig"
	"github.com/Cloverhound/webex-cli/internal/readonly"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Webex via OAuth",
	Long: `Opens a browser for Webex OAuth login. Stores tokens in the OS keyring for the authenticated user.

With --read-only, the login requests only read scopes and turns on read-only
mode: stored logins with write access are deleted from the keyring, write
requests are refused, and only read-only logins can be used. Leaving read-only
mode requires a plain 'webex login' from an interactive terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		readOnlyLogin, _ := cmd.Flags().GetBool("read-only")

		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		clientID := cfg.EffectiveClientID()
		clientSecret := cfg.EffectiveClientSecret()
		scopes := cfg.EffectiveScopes()
		if readOnlyLogin {
			scopes = cfg.EffectiveReadOnlyScopes()
			if err := readonly.ValidateScopes(scopes); err != nil {
				return fmt.Errorf("read-only-scopes config: %w", err)
			}
			fmt.Println("Read-only login: stored logins with write access will be removed.")
		} else if readOnlyMode(cfg) {
			if err := confirmLeaveReadOnly(); err != nil {
				return err
			}
		}

		result, err := auth.Login(clientID, clientSecret, scopes)
		if err != nil {
			return err
		}
		if readOnlyLogin {
			if err := readonly.ValidateScopes(result.Token.Scopes); err != nil {
				return fmt.Errorf("granted scopes include write access, token not saved: %w", err)
			}
			result.Token.ReadOnly = true
		}

		// Offer to store OAuth credentials with this user account before saving the token.
		// This allows each user to have independent client credentials for token refresh.
		if clientID != "" {
			promptStoreCredentials(&result.Token, clientID, clientSecret)
		}

		// Store token in keyring
		if err := auth.SaveToken(result.Email, &result.Token); err != nil {
			return fmt.Errorf("saving token: %w", err)
		}

		// Update config
		cfg.AddUser(result.Email, result.DisplayName, result.OrgID, result.OrgName)
		cfg.SetDefaultUser(result.Email)
		cfg.ReadOnly = readOnlyLogin
		var purged []string
		var purgeErr error
		if readOnlyLogin {
			purged, purgeErr = auth.PurgeWriteTokens(cfg)
		}
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		orgInfo := ""
		if result.OrgName != "" {
			orgInfo = fmt.Sprintf(" — %s", result.OrgName)
		}
		fmt.Printf("Logged in as %s (%s)%s\n", result.DisplayName, result.Email, orgInfo)
		if readOnlyLogin {
			fmt.Println("Read-only mode is on.")
			for _, email := range purged {
				fmt.Printf("Removed write-capable login: %s\n", email)
			}
			if purgeErr != nil {
				return purgeErr
			}
		}
		fmt.Println()

		// Offer to associate this user with the current folder
		if cwd, err := os.Getwd(); err == nil {
			promptFolderAssociation(result.Email, cwd)
		}

		return nil
	},
}

func promptStoreCredentials(tok *auth.StoredToken, clientID, clientSecret string) {
	var choice string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Save OAuth credentials with this account?").
				Description(
					"Storing credentials per-account lets each user refresh tokens independently,\n" +
						"even when multiple accounts are configured with different OAuth apps.",
				).
				Options(
					huh.NewOption("Yes — save with this account (recommended for multiple users)", "yes"),
					huh.NewOption("No — use global credentials only", "no"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return
	}

	if choice == "yes" {
		tok.ClientID = clientID
		tok.ClientSecret = clientSecret
	}
}

func promptFolderAssociation(email, dir string) {
	folderName := filepath.Base(dir)

	var choice string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Associate this user with the current folder?").
				Description(
					fmt.Sprintf(
						"This saves \"%s\" as the default user for %s/.\n"+
							"Useful when different folders connect to different Webex orgs,\n"+
							"so the right credentials are used automatically.",
						email, folderName,
					),
				).
				Options(
					huh.NewOption(fmt.Sprintf("Yes — use \"%s\" whenever I'm in %s/", email, folderName), "yes"),
					huh.NewOption("No — don't set folder default", "no"),
				).
				Value(&choice),
		),
	)

	if err := form.Run(); err != nil {
		return
	}

	if choice == "yes" {
		if err := localconfig.Save(dir, email); err != nil {
			fmt.Printf("Warning: could not save folder config: %v\n", err)
			return
		}
		fmt.Printf("Saved to %s/.webex-cli/config.json\n", folderName)
	}
}

func init() {
	loginCmd.Flags().Bool("read-only", false, "Log in with read-only scopes and turn on read-only mode")
	rootCmd.AddCommand(loginCmd)
}
