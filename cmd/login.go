package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
	"github.com/Cloverhound/webex-cli/internal/localconfig"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Webex via OAuth",
	Long:  "Opens a browser for Webex OAuth login. Stores tokens in the OS keyring for the authenticated user.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		clientID := cfg.EffectiveClientID()
		clientSecret := cfg.EffectiveClientSecret()
		scopes := cfg.EffectiveScopes()

		result, err := auth.Login(clientID, clientSecret, scopes)
		if err != nil {
			return err
		}

		// Store token in keyring
		if err := auth.SaveToken(result.Email, &result.Token); err != nil {
			return fmt.Errorf("saving token: %w", err)
		}

		// Update config
		cfg.AddUser(result.Email, result.DisplayName, result.OrgID, result.OrgName)
		cfg.SetDefaultUser(result.Email)
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		orgInfo := ""
		if result.OrgName != "" {
			orgInfo = fmt.Sprintf(" — %s", result.OrgName)
		}
		fmt.Printf("Logged in as %s (%s)%s\n\n", result.DisplayName, result.Email, orgInfo)

		// Offer to associate this user with the current folder
		if cwd, err := os.Getwd(); err == nil {
			promptFolderAssociation(result.Email, cwd)
		}

		return nil
	},
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
	rootCmd.AddCommand(loginCmd)
}
