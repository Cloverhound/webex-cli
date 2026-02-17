package cmd

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
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

		cfg.SetDefaultUser(email)
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Default user set to %s\n", email)
		return nil
	},
}

func init() {
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authListCmd)
	authCmd.AddCommand(authSwitchCmd)
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
