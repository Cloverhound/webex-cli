package cmd

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  "View and manage authenticated users, check token status, and switch between users.",
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		email := cfg.DefaultUser
		if email == "" {
			fmt.Println("No authenticated user. Run: webex login")
			return nil
		}

		userInfo := cfg.Users[email]
		tok, err := auth.LoadToken(email)
		if err != nil {
			fmt.Printf("User:   %s\n", email)
			fmt.Println("Status: token not found in keyring")
			return nil
		}

		fmt.Printf("User:    %s", email)
		if userInfo.DisplayName != "" {
			fmt.Printf(" (%s)", userInfo.DisplayName)
		}
		fmt.Println()

		if userInfo.OrgName != "" {
			fmt.Printf("Org:     %s (%s)\n", userInfo.OrgName, userInfo.OrgID)
		} else if userInfo.OrgID != "" {
			fmt.Printf("Org:     %s\n", userInfo.OrgID)
		}

		if cfg.DefaultOrgName != "" {
			fmt.Printf("Org override: %s (%s)\n", cfg.DefaultOrgName, cfg.DefaultOrgID)
		} else if cfg.DefaultOrgID != "" {
			fmt.Printf("Org override: %s\n", cfg.DefaultOrgID)
		}

		fmt.Printf("Source:  keyring\n")

		if tok.IsExpired() {
			fmt.Printf("Token:   expired (at %s)\n", tok.ExpiresAt.Format(time.RFC3339))
			if tok.IsRefreshExpired() {
				fmt.Println("Refresh: expired — run: webex login")
			} else {
				fmt.Println("Refresh: available (will auto-refresh on next command)")
			}
		} else {
			remaining := time.Until(tok.ExpiresAt).Round(time.Second)
			fmt.Printf("Token:   valid (expires in %s)\n", remaining)
		}

		// Live check
		if !tok.IsExpired() {
			if ok := liveCheck(tok.AccessToken); ok {
				fmt.Println("Live:    verified")
			} else {
				fmt.Println("Live:    failed (token may have been revoked)")
			}
		}

		return nil
	},
}

var authListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all authenticated users",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if len(cfg.Users) == 0 {
			fmt.Println("No authenticated users. Run: webex login")
			return nil
		}

		fmt.Printf("%-30s  %-20s  %-25s  %-10s  %s\n", "EMAIL", "NAME", "ORG", "STATUS", "")
		fmt.Printf("%-30s  %-20s  %-25s  %-10s  %s\n", "-----", "----", "---", "------", "")

		for email, info := range cfg.Users {
			status := "unknown"
			tok, err := auth.LoadToken(email)
			if err != nil {
				status = "no token"
			} else if tok.IsExpired() {
				if tok.IsRefreshExpired() {
					status = "expired"
				} else {
					status = "refreshable"
				}
			} else {
				status = "active"
			}

			marker := ""
			if email == cfg.DefaultUser {
				marker = "(default)"
			}

			name := info.DisplayName
			if len(name) > 20 {
				name = name[:17] + "..."
			}
			org := info.OrgName
			if len(org) > 25 {
				org = org[:22] + "..."
			}

			fmt.Printf("%-30s  %-20s  %-25s  %-10s  %s\n", email, name, org, status, marker)
		}

		return nil
	},
}

var authSwitchCmd = &cobra.Command{
	Use:   "switch <email>",
	Short: "Switch the default user",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]

		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if _, ok := cfg.Users[email]; !ok {
			return fmt.Errorf("user %s not found — run: webex login", email)
		}

		clearedOrg := cfg.DefaultOrgName
		if clearedOrg == "" {
			clearedOrg = cfg.DefaultOrgID
		}

		cfg.SetDefaultUser(email)
		cfg.DefaultOrgID = ""
		cfg.DefaultOrgName = ""
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Default user set to %s\n", email)
		if clearedOrg != "" {
			fmt.Printf("Cleared org override (was: %s)\n", clearedOrg)
		}
		return nil
	},
}

var authSetOrgCmd = &cobra.Command{
	Use:   "set-org <org-id>",
	Short: "Set a persistent organization override",
	Long:  "Set a default organization ID that will be used for all commands unless overridden by --organization.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		orgInput := args[0]

		// Normalize to base64 for the API call (Webex org IDs are base64-encoded)
		uuid := config.DecodeOrgID(orgInput)
		base64ID := config.EncodeOrgID(uuid)

		// Validate by fetching the org name
		orgName, err := auth.FetchOrgName(config.Token(), base64ID)
		if err != nil {
			fmt.Printf("Error: could not validate org %s\n", orgInput)
			fmt.Printf("  %s\n", err)
			fmt.Println()
			fmt.Println("To see available users and their orgs: webex auth list")
			fmt.Println("To switch to a different user:         webex auth switch <email>")
			return nil
		}

		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		cfg.DefaultOrgID = base64ID
		cfg.DefaultOrgName = orgName
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Org override set: %s (%s)\n", orgName, uuid)
		return nil
	},
}

var authClearOrgCmd = &cobra.Command{
	Use:   "clear-org",
	Short: "Clear the persistent organization override",
	Long:  "Remove the default organization override, reverting to the login user's home org.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		if cfg.DefaultOrgID == "" {
			fmt.Println("No org override is set.")
			return nil
		}

		was := cfg.DefaultOrgName
		if was == "" {
			was = cfg.DefaultOrgID
		}

		cfg.DefaultOrgID = ""
		cfg.DefaultOrgName = ""
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Org override cleared (was: %s)\n", was)

		// Show what org will now be used
		if cfg.DefaultUser != "" {
			if userInfo, ok := cfg.Users[cfg.DefaultUser]; ok && userInfo.OrgName != "" {
				fmt.Printf("Now using: %s (%s)\n", userInfo.OrgName, cfg.DefaultUser)
			}
		}

		return nil
	},
}

func init() {
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authListCmd)
	authCmd.AddCommand(authSwitchCmd)
	authCmd.AddCommand(authSetOrgCmd)
	authCmd.AddCommand(authClearOrgCmd)
	rootCmd.AddCommand(authCmd)
}

func liveCheck(accessToken string) bool {
	req, _ := http.NewRequest("GET", "https://webexapis.com/v1/people/me", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	return resp.StatusCode == 200
}
