package cmd

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/localconfig"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "webex",
	Short: "Webex CLI — manage Webex APIs",
	Long:  `A command-line interface for Webex APIs — Admin, Calling, Contact Center, Devices, Meetings, and Messaging.`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Debug mode (set early so auth debug works)
		debug, _ := cmd.Flags().GetBool("debug")
		config.SetDebug(debug)

		// Dry-run mode
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		config.SetDryRun(dryRun)

		// Output format
		format, _ := cmd.Flags().GetString("output")
		output.SetFormat(format)

		// Pagination
		paginate, _ := cmd.Flags().GetBool("paginate")
		config.SetPaginate(paginate)

		// Rate-limit retry controls
		maxRetry, _ := cmd.Flags().GetInt("max-retry")
		config.SetMaxRetry(maxRetry)
		maxRetryTimer, _ := cmd.Flags().GetInt("max-retry-timer")
		config.SetMaxRetryTimer(maxRetryTimer)

		// Load app config early (safe local file read, needed by skipAuth-exempt commands like set-org)
		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		// Data region for region-specific hosts (Detailed Call History):
		// --region flag > WEBEX_REGION env > config file.
		regionFlag, _ := cmd.Flags().GetString("region")
		if regionFlag == "" {
			regionFlag = os.Getenv("WEBEX_REGION")
		}
		if regionFlag == "" {
			regionFlag = cfg.Region
		}
		if err := setRegion(regionFlag); err != nil {
			return err
		}

		// Skip auth for certain commands
		if skipAuth(cmd) {
			return nil
		}

		// Resolve token
		flagToken, _ := cmd.Flags().GetString("token")
		envToken := os.Getenv("WEBEX_TOKEN")
		userFlag, _ := cmd.Flags().GetString("user")
		envUser := os.Getenv("WEBEX_USER")

		// Check local folder config before falling back to global default
		if userFlag == "" && envUser == "" {
			if cwd, err := os.Getwd(); err == nil {
				if lcfg, err := localconfig.Load(cwd); err == nil && lcfg != nil && lcfg.User != "" {
					envUser = lcfg.User
				}
			}
		}

		result, err := auth.ResolveToken(flagToken, envToken, userFlag, envUser, cfg)
		if err != nil {
			return err
		}

		config.SetToken(result.Token)

		// Wire up token refresher for keyring-based auth
		if result.Source == auth.SourceKeyring && result.UserEmail != "" {
			config.TokenRefresher = auth.MakeRefresher(result.UserEmail, cfg)
		}

		// Organization: --organization flag > config default org > resolved user's org > default user's org.
		// SetOrgID stores both UUID and base64 formats for downstream use.
		orgFlag, _ := cmd.Flags().GetString("organization")
		if orgFlag != "" {
			config.SetOrgID(orgFlag)
		} else if cfg.DefaultOrgID != "" {
			config.SetOrgID(cfg.DefaultOrgID)
		} else if result.OrgID != "" {
			config.SetOrgID(result.OrgID)
		} else if cfg.DefaultUser != "" {
			// Fallback: use default user's org from config (covers env/flag token sources)
			if userInfo, ok := cfg.Users[cfg.DefaultUser]; ok && userInfo.OrgID != "" {
				config.SetOrgID(userInfo.OrgID)
			}
		}

		// Auto-populate --orgid on CC commands from the resolved org.
		// CC API uses UUID format, so decode any base64 values.
		if f := cmd.Flags().Lookup("orgid"); f != nil {
			if f.Value.String() == "" {
				if config.OrgID() != "" {
					cmd.Flags().Set("orgid", config.OrgID())
				}
			} else {
				cmd.Flags().Set("orgid", config.DecodeOrgID(f.Value.String()))
			}
		}

		// Auto-populate --org-id on non-CC commands from base64 org ID.
		// These APIs require base64 format, so convert if needed.
		if f := cmd.Flags().Lookup("org-id"); f != nil {
			if f.Value.String() == "" {
				if config.OrgIDBase64() != "" {
					cmd.Flags().Set("org-id", config.OrgIDBase64())
				}
			} else {
				// User passed --org-id directly; normalize to base64
				decoded := config.DecodeOrgID(f.Value.String())
				cmd.Flags().Set("org-id", config.EncodeOrgID(decoded))
			}
		}

		return nil
	},
}

// setRegion validates a region value and stores it for region-specific hosts.
func setRegion(region string) error {
	if region == "" {
		config.SetRegion("")
		return nil
	}
	normalized := strings.ToLower(strings.TrimSpace(region))
	if !slices.Contains(config.Regions, normalized) {
		return fmt.Errorf("unknown region %q (valid: %s)", region, strings.Join(config.Regions, ", "))
	}
	config.SetRegion(normalized)
	return nil
}

// RootCommand returns the root cobra command for external introspection.
func RootCommand() *cobra.Command { return rootCmd }

func Execute() error {
	// Remove "required" from --orgid/--org-id flags on all subcommands.
	// This allows PersistentPreRunE to auto-populate from config/login.
	stripRequiredOrgID(rootCmd)
	// Hide per-command org flags so only --organization is visible in help.
	hideOrgFlags(rootCmd)
	err := rootCmd.Execute()
	if errors.Is(err, client.ErrDryRun) {
		return nil
	}
	return err
}

// stripRequiredOrgID recursively removes the "required" annotation from --orgid and --org-id flags.
func stripRequiredOrgID(cmd *cobra.Command) {
	for _, name := range []string{"orgid", "org-id"} {
		if f := cmd.Flags().Lookup(name); f != nil {
			cmd.Flags().SetAnnotation(name, cobra.BashCompOneRequiredFlag, []string{"false"})
			annotations := f.Annotations
			if annotations != nil {
				delete(annotations, cobra.BashCompOneRequiredFlag)
			}
		}
	}
	for _, child := range cmd.Commands() {
		stripRequiredOrgID(child)
	}
}

// hideOrgFlags recursively hides --orgid and --org-id flags so only --organization is visible.
func hideOrgFlags(cmd *cobra.Command) {
	for _, name := range []string{"orgid", "org-id"} {
		if f := cmd.Flags().Lookup(name); f != nil {
			f.Hidden = true
		}
	}
	for _, child := range cmd.Commands() {
		hideOrgFlags(child)
	}
}

func init() {
	rootCmd.PersistentFlags().String("token", "", "Webex API token (overrides keyring)")
	rootCmd.PersistentFlags().String("output", "json", "Output format: json, table, csv, raw")
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug logging of HTTP requests")
	rootCmd.PersistentFlags().Bool("paginate", false, "Auto-paginate list results")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Print write requests without executing them")
	rootCmd.PersistentFlags().String("user", "", "Use a specific authenticated user (email)")
	rootCmd.PersistentFlags().String("organization", "", "Override organization ID for this command")
	rootCmd.PersistentFlags().String("region", "", "Data region for region-specific APIs: "+strings.Join(config.Regions, ", ")+" (default us)")
	rootCmd.PersistentFlags().Int("max-retry", 3, "Max number of 429 retries before giving up (0 = no retries)")
	rootCmd.PersistentFlags().Int("max-retry-timer", 60, "Max total seconds to wait across all 429 retries (0 = unlimited)")

	// On flag errors (unknown flag, bad value), print usage with valid flags.
	// SilenceUsage suppresses Cobra's automatic usage, so we print it ourselves.
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.PrintErrf("Error: %s\n\n", err)
		cmd.PrintErr(cmd.UsageString())
		cmd.Root().SilenceErrors = true
		return err
	})
}

// skipAuth returns true for commands that don't need authentication. Only the
// top-level name is matched: generated API subcommands reuse names from the
// exempt list ("update", "login", "logout") and do need a token.
func skipAuth(cmd *cobra.Command) bool {
	if !cmd.HasParent() {
		return true
	}
	switch cmd.Name() {
	case "help", "completion", cobra.ShellCompRequestCmd, cobra.ShellCompNoDescRequestCmd:
		return true
	}

	top := cmd
	for top.Parent().HasParent() {
		top = top.Parent()
	}

	switch top.Name() {
	case "login", "logout", "config", "version", "update", "post-install", "help", "completion":
		return true
	case "auth":
		// set-org validates the org against the API; token prints it.
		return cmd.Name() != "set-org" && cmd.Name() != "token"
	}

	return false
}
