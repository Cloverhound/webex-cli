package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/localconfig"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  "View and manage authenticated users, check token status, and switch between users.",
}

var authTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Print the current access token",
	Long:  "Print the access token for the active user, refreshing it if expired. Useful for making manual API calls with curl.",
	RunE: func(cmd *cobra.Command, args []string) error {
		tok := config.Token()
		if tok == "" {
			return fmt.Errorf("no access token available; run: webex login")
		}
		fmt.Println(tok)
		return nil
	},
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
		envSet := false
		for _, env := range []string{"WEBEX_TOKEN", auth.RefreshTokenEnv} {
			if os.Getenv(env) == "" {
				continue
			}
			envSet = true
			switch {
			case readOnlyMode(cfg):
				fmt.Printf("Env:     $%s is set; read-only mode refuses it, so commands fail until it is unset\n", env)
			case email == "":
				fmt.Printf("Env:     $%s is set; commands authenticate with it\n", env)
			default:
				fmt.Printf("Env:     $%s is set and overrides the stored login below\n", env)
			}
			break
		}

		if email == "" {
			if !envSet {
				fmt.Println("No authenticated user. Run: webex login")
			}
			return nil
		}

		userInfo := cfg.Users[email]
		tok, store, err := auth.LoadTokenWithStore(email)
		if err != nil {
			fmt.Printf("User:    %s\n", email)
			if errors.Is(err, auth.ErrTokenNotFound) {
				fmt.Println("Status:  token not found in the keyring or credentials file")
			} else {
				fmt.Printf("Status:  %v\n", err)
			}
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

		if store == auth.StoreFile {
			fmt.Printf("Store:   file, plain text (%s)\n", auth.CredentialsPath())
		} else {
			fmt.Println("Store:   OS keyring")
		}

		switch {
		case cfg.ReadOnly:
			fmt.Println("Mode:    read-only")
		case readOnlyEnvSet():
			fmt.Printf("Mode:    read-only ($%s)\n", readOnlyEnv)
		}
		if tok.ReadOnly {
			fmt.Printf("Scopes:  %s\n", tok.Scopes)
		}

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

		// Show folder default if set
		if cwd, err := os.Getwd(); err == nil {
			if lcfg, err := localconfig.Load(cwd); err == nil && lcfg != nil && lcfg.User != "" {
				fmt.Printf("Folder:  %s (for this directory)\n", lcfg.User)
			}
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

var authExportCmd = &cobra.Command{
	Use:   "export [email]",
	Short: "Print a stored refresh token for use as $WEBEX_REFRESH_TOKEN",
	Long: `Prints the refresh token of a stored login so a login made on one machine can
seed a CI secret or an unattended agent through $WEBEX_REFRESH_TOKEN. The user
is the email given, else the one --user, $WEBEX_USER, the folder default, or the
default user selects. Anyone holding the token can act as that user until it
expires (about 90 days after its last use), so treat it as a password.

In read-only mode only read-only logins can be exported, and the consumer must
not be in read-only mode, which refuses $WEBEX_REFRESH_TOKEN.

The token is printed alone on stdout. Confirm at the prompt, or pass --yes
when there is no terminal.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		yes, _ := cmd.Flags().GetBool("yes")

		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		flagUser, envUser := requestedUser(cmd)
		email := cfg.DefaultUser
		switch {
		case len(args) > 0:
			email = args[0]
		case flagUser != "":
			email = flagUser
		case envUser != "":
			email = envUser
		}
		if email == "" {
			return fmt.Errorf("no authenticated user — run: webex login")
		}
		tok, err := auth.LoadToken(email)
		if err != nil {
			return fmt.Errorf("no stored login for %s — run: webex login", email)
		}
		// The command skips token resolution, which is where read-only mode
		// normally rejects write-capable logins.
		if readOnlyMode(cfg) && !tok.ReadOnly {
			return fmt.Errorf("read-only mode: %s has no read-only login; only read-only logins can be exported", email)
		}
		if tok.RefreshToken == "" || tok.IsRefreshExpired() {
			return fmt.Errorf("the stored login for %s has no usable refresh token — run: webex login", email)
		}

		if !yes {
			if !canPrompt() {
				return fmt.Errorf("refusing to print a refresh token without confirmation; pass --yes")
			}
			ok, err := confirm(fmt.Sprintf("Print the refresh token for %s?", email),
				"Anyone with this token can act as you until it expires.",
				"Print it", "Cancel")
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("cancelled")
			}
		}

		if tok.ClientID != "" && tok.ClientID != appconfig.DefaultClientID {
			fmt.Fprintf(os.Stderr, "This login uses its own OAuth integration; also set $%s and $%s.\n",
				auth.ClientIDEnv, auth.ClientSecretEnv)
		}
		fmt.Fprintf(os.Stderr, "Webex may rotate refresh tokens on use. If this machine and the consumer of the export\n"+
			"both refresh it, one of them can be signed out; prefer a separate login for each.\n")
		fmt.Println(tok.RefreshToken)
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

			var markers []string
			if email == cfg.DefaultUser {
				markers = append(markers, "(default)")
			}
			if err == nil && tok.ReadOnly {
				markers = append(markers, "(read-only)")
			}
			marker := strings.Join(markers, " ")

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
		if readOnlyMode(cfg) && !auth.HasReadOnlyToken(email) {
			return fmt.Errorf("read-only mode: %s has no read-only login — run: webex login --read-only", email)
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

var authSetFolderDefaultCmd = &cobra.Command{
	Use:   "set-folder-default <email>",
	Short: "Set the default user for the current folder",
	Long:  "Associates a Webex user with the current working directory. When running commands from this folder, this user's credentials will be used automatically.",
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

		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		if err := localconfig.Save(cwd, email); err != nil {
			return fmt.Errorf("saving folder default: %w", err)
		}

		fmt.Printf("Set folder default to %s for %s\n", email, cwd)
		return nil
	},
}

var authClearFolderDefaultCmd = &cobra.Command{
	Use:   "clear-folder-default",
	Short: "Remove the folder default user for the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		lcfg, err := localconfig.Load(cwd)
		if err != nil || lcfg == nil {
			fmt.Println("No folder default set for this directory.")
			return nil
		}

		configPath := filepath.Join(cwd, ".webex-cli", "config.json")
		if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("removing folder config: %w", err)
		}

		configDir := filepath.Join(cwd, ".webex-cli")
		os.Remove(configDir) // ignore error — dir may not be empty

		fmt.Printf("Cleared folder default (was %s)\n", lcfg.User)
		return nil
	},
}

func init() {
	authCmd.AddCommand(authTokenCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authListCmd)
	authExportCmd.Flags().Bool("yes", false, "Print without asking for confirmation")
	authCmd.AddCommand(authExportCmd)
	authCmd.AddCommand(authSwitchCmd)
	authCmd.AddCommand(authSetOrgCmd)
	authCmd.AddCommand(authClearOrgCmd)
	authCmd.AddCommand(authSetFolderDefaultCmd)
	authCmd.AddCommand(authClearFolderDefaultCmd)
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
