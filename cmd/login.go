package cmd

import (
	"fmt"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
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
		fmt.Printf("Logged in as %s (%s)%s\n", result.DisplayName, result.Email, orgInfo)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
