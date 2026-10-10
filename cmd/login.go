package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
	"github.com/Cloverhound/webex-cli/internal/localconfig"
	"github.com/Cloverhound/webex-cli/internal/readonly"
	"github.com/charmbracelet/huh"
	"github.com/mdp/qrterminal/v3"
	"github.com/spf13/cobra"
)

const loginModeEnv = "WEBEX_LOGIN_MODE"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Webex via OAuth",
	Long: `Logs in to Webex with OAuth and stores the tokens for the authenticated user.

On a desktop, login opens a browser. With --device, or automatically over SSH,
in CI, on Linux without a display, or when no browser can be opened, login
prints a URL and a code to approve from any browser instead. Set
$` + loginModeEnv + `=device or =browser to choose the flow without a flag.
With --output json, the device code is printed as JSON on stdout.

Tokens are stored in the OS keyring. When no keyring is available they are
stored in plain text in ` + "`~/.config/webex-cli/credentials.json`" + ` (mode 0600);
set $` + auth.TokenStoreEnv + `=keyring or =file to choose the store.

With --read-only, the login requests only read scopes and turns on read-only
mode: stored logins with write access are deleted, write requests are refused,
and only read-only logins can be used. Leaving read-only mode requires a plain
'webex login' from an interactive terminal.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		readOnlyLogin, _ := cmd.Flags().GetBool("read-only")
		deviceFlag, _ := cmd.Flags().GetBool("device")
		browserFlag, _ := cmd.Flags().GetBool("browser")
		format, _ := cmd.Flags().GetString("output")
		jsonOut := cmd.Flags().Changed("output") && format == "json"

		envHeadless := auth.HeadlessReason(os.Getenv, runtime.GOOS)
		mode, headless, err := chooseLoginMode(deviceFlag, browserFlag, os.Getenv(loginModeEnv), envHeadless)
		if err != nil {
			return err
		}

		// Machine-readable output owns stdout; everything for people goes to stderr.
		var msg io.Writer = os.Stdout
		if jsonOut {
			msg = os.Stderr
		}
		interactive := !jsonOut && canPrompt()

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
			fmt.Fprintln(msg, "Read-only login: stored logins with write access will be removed.")
		} else if readOnlyMode(cfg) {
			if err := confirmLeaveReadOnly(); err != nil {
				return err
			}
		}

		deviceLogin := func() (*auth.LoginResult, error) {
			return auth.DeviceLogin(clientID, clientSecret, scopes, func(dc *auth.DeviceCode) {
				showDeviceCode(dc, jsonOut)
			})
		}

		var result *auth.LoginResult
		switch mode {
		case "device":
			if headless != "" {
				fmt.Fprintf(os.Stderr, "Using device login (%s). Use --browser to force the browser flow.\n", headless)
			}
			result, err = deviceLogin()
		case "browser":
			result, err = auth.Login(clientID, clientSecret, scopes, false)
		default:
			result, err = auth.Login(clientID, clientSecret, scopes, true)
			if errors.Is(err, auth.ErrNoBrowser) {
				fmt.Fprintf(os.Stderr, "Browser login unavailable (%v); using device login.\n", err)
				result, err = deviceLogin()
			}
		}
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
		if clientID != "" && interactive {
			promptStoreCredentials(&result.Token, clientID, clientSecret)
		}

		store, err := auth.SaveTokenWithStore(result.Email, &result.Token)
		if err != nil {
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
		fmt.Fprintf(msg, "Logged in as %s (%s)%s\n", result.DisplayName, result.Email, orgInfo)
		if readOnlyLogin {
			fmt.Fprintln(msg, "Read-only mode is on.")
			for _, email := range purged {
				fmt.Fprintf(msg, "Removed write-capable login: %s\n", email)
			}
		}
		if jsonOut {
			out := map[string]any{
				"status":       "logged_in",
				"email":        result.Email,
				"display_name": result.DisplayName,
				"org_id":       result.OrgID,
				"org_name":     result.OrgName,
				"read_only":    readOnlyLogin,
				"token_store":  store,
			}
			if readOnlyLogin {
				out["removed_logins"] = purged
				if purgeErr != nil {
					out["purge_error"] = purgeErr.Error()
				}
			}
			_ = json.NewEncoder(os.Stdout).Encode(out)
		}
		if purgeErr != nil {
			return purgeErr
		}

		// Offer to associate this user with the current folder
		if interactive {
			fmt.Println()
			if cwd, err := os.Getwd(); err == nil {
				promptFolderAssociation(result.Email, cwd)
			}
		}

		return nil
	},
}

// chooseLoginMode returns "device", "browser" or "auto" (browser, falling back
// to device when no browser can be used), plus the reason when the environment
// looks headless.
func chooseLoginMode(device, browser bool, envMode, headless string) (mode, reason string, err error) {
	if device && browser {
		return "", "", fmt.Errorf("--device and --browser cannot be used together")
	}
	switch {
	case device:
		return "device", "", nil
	case browser:
		return "browser", "", nil
	}
	switch strings.ToLower(strings.TrimSpace(envMode)) {
	case "":
	case "device":
		return "device", "", nil
	case "browser":
		return "browser", "", nil
	default:
		return "", "", fmt.Errorf("$%s must be \"device\" or \"browser\", not %q", loginModeEnv, envMode)
	}
	if headless != "" {
		return "device", headless, nil
	}
	return "auto", "", nil
}

// showDeviceCode tells the user how to approve the login. The first line is
// meant to be relayed word for word by an agent running the CLI on the user's
// behalf; the user approves in their own browser, not where the CLI runs, so
// no browser is opened here.
func showDeviceCode(dc *auth.DeviceCode, jsonOut bool) {
	link := dc.VerificationURIComplete
	if link == "" {
		link = dc.VerificationURI
	}
	if jsonOut {
		_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
			"status":                    "pending",
			"verification_uri":          dc.VerificationURI,
			"verification_uri_complete": dc.VerificationURIComplete,
			"user_code":                 dc.UserCode,
			"expires_in":                dc.ExpiresIn,
			"expires_at":                time.Now().Add(time.Duration(dc.ExpiresIn) * time.Second).UTC().Format(time.RFC3339),
		})
	}

	fmt.Fprintf(os.Stderr, "Open %s and enter %s\n", dc.VerificationURI, dc.UserCode)
	if isTerminal(os.Stderr) {
		fmt.Fprintln(os.Stderr, "Or scan this code with a phone:")
		qrterminal.GenerateHalfBlock(link, qrterminal.L, os.Stderr)
	}
	fmt.Fprintf(os.Stderr, "Waiting for approval (the code expires in %d minutes)...\n", (dc.ExpiresIn+59)/60)
}

func promptStoreCredentials(tok *auth.StoredToken, clientID, clientSecret string) {
	var choice string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Save OAuth credentials with this account?").
				Description(
					"Storing credentials per-account lets each user refresh tokens independently,\n"+
						"even when multiple accounts are configured with different OAuth apps.",
				).
				Options(
					huh.NewOption("Yes — save with this account (recommended for multiple users)", "yes"),
					huh.NewOption("No — use global credentials only", "no"),
				).
				Value(&choice),
		),
	).WithOutput(os.Stderr)

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
	).WithOutput(os.Stderr)

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
	loginCmd.Flags().Bool("device", false, "Log in with a code approved from any browser (no local browser needed)")
	loginCmd.Flags().Bool("browser", false, "Force the local browser login even when no browser is detected")
	rootCmd.AddCommand(loginCmd)
}
