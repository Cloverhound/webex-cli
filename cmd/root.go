package cmd

import (
	"fmt"
	"os"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/Cloverhound/webex-cli/internal/auth"
	"github.com/Cloverhound/webex-cli/internal/config"
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

		// Output format
		format, _ := cmd.Flags().GetString("output")
		output.SetFormat(format)

		// Pagination
		paginate, _ := cmd.Flags().GetBool("paginate")
		config.SetPaginate(paginate)

		// Skip auth for certain commands
		if skipAuth(cmd) {
			return nil
		}

		// Load app config
		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		// Resolve token
		flagToken, _ := cmd.Flags().GetString("token")
		envToken := os.Getenv("WEBEX_TOKEN")
		userFlag, _ := cmd.Flags().GetString("user")
		envUser := os.Getenv("WEBEX_USER")

		result, err := auth.ResolveToken(flagToken, envToken, userFlag, envUser, cfg)
		if err != nil {
			return err
		}

		config.SetToken(result.Token)

		// Wire up token refresher for keyring-based auth
		if result.Source == auth.SourceKeyring && result.UserEmail != "" {
			config.TokenRefresher = auth.MakeRefresher(result.UserEmail, cfg)
		}

		// Organization: --organization flag > resolved user's org > default user's org from config.
		// SetOrgID stores both UUID and base64 formats for downstream use.
		orgFlag, _ := cmd.Flags().GetString("organization")
		if orgFlag != "" {
			config.SetOrgID(orgFlag)
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

func Execute() error {
	// Remove "required" from --orgid/--org-id flags on all subcommands.
	// This allows PersistentPreRunE to auto-populate from config/login.
	stripRequiredOrgID(rootCmd)
	// Hide per-command org flags so only --organization is visible in help.
	hideOrgFlags(rootCmd)
	return rootCmd.Execute()
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
	rootCmd.PersistentFlags().String("user", "", "Use a specific authenticated user (email)")
	rootCmd.PersistentFlags().String("organization", "", "Override organization ID for this command")

	// On flag errors (unknown flag, bad value), print usage with valid flags.
	// SilenceUsage suppresses Cobra's automatic usage, so we print it ourselves.
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.PrintErrf("Error: %s\n\n", err)
		cmd.PrintErr(cmd.UsageString())
		cmd.Root().SilenceErrors = true
		return err
	})
}

// skipAuth returns true for commands that don't need authentication.
func skipAuth(cmd *cobra.Command) bool {
	// Walk up to find the root-level command name
	name := cmd.Name()

	// Check the command itself and all parents
	for c := cmd; c != nil; c = c.Parent() {
		switch c.Name() {
		case "login", "logout", "auth", "config", "version", "update", "help", "webex":
			// "webex" is the root — only skip if it's the actual command being run (bare `webex`)
			if c.Name() == "webex" {
				continue
			}
			return true
		}
	}

	// Also skip bare root command and help
	if name == "help" || name == "webex" {
		return true
	}

	return false
}
