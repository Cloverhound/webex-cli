package cmd

import (
	"fmt"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout [email]",
	Short: "Log out and remove stored credentials",
	Long:  "Removes stored tokens from the OS keyring. Without arguments, removes the default user. Specify an email to remove a specific user, or --all to remove all users.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		allFlag, _ := cmd.Flags().GetBool("all")

		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if allFlag {
			for _, email := range cfg.UserEmails() {
				if err := auth.DeleteToken(email); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Warning: could not remove token for %s: %v\n", email, err)
				}
			}
			cfg.Users = make(map[string]appconfig.UserInfo)
			cfg.DefaultUser = ""
			if err := cfg.Save(); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}
			fmt.Println("Logged out all users")
			return nil
		}

		email := ""
		if len(args) > 0 {
			email = args[0]
		} else {
			email = cfg.DefaultUser
		}

		if email == "" {
			return fmt.Errorf("no user specified and no default user — nothing to log out")
		}

		if _, ok := cfg.Users[email]; !ok {
			return fmt.Errorf("user %s not found", email)
		}

		if err := auth.DeleteToken(email); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Warning: could not remove token for %s: %v\n", email, err)
		}
		cfg.RemoveUser(email)
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Logged out %s\n", email)
		return nil
	},
}

func init() {
	logoutCmd.Flags().Bool("all", false, "Remove all stored users and tokens")
	rootCmd.AddCommand(logoutCmd)
}
