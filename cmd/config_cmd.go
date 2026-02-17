package cmd

import (
	"fmt"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long: `Get and set CLI configuration values like OAuth client credentials.

To use your own Webex Integration instead of the built-in default:

  1. Create an integration at https://developer.webex.com/my-apps
  2. Set the redirect URI to: http://localhost:8085/callback
  3. Add the OAuth scopes required for the APIs you plan to use
  4. Run:  webex config set client-id <your-client-id>
           webex config set client-secret <your-client-secret>
           webex config set scopes "spark:all spark-admin:all <additional-scopes>"`,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  "Set a configuration value. Valid keys: client-id, client-secret, scopes",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		switch key {
		case "client-id":
			cfg.ClientID = value
		case "client-secret":
			cfg.ClientSecret = value
		case "scopes":
			cfg.Scopes = value
		default:
			return fmt.Errorf("unknown config key: %s (valid: client-id, client-secret, scopes)", key)
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("saving config: %w", err)
		}

		fmt.Printf("Set %s\n", key)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long:  "Get a configuration value. Valid keys: client-id, client-secret, scopes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		cfg, err := appconfig.Load()
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		switch key {
		case "client-id":
			value := cfg.EffectiveClientID()
			if value == "" {
				fmt.Println("(not set)")
			} else if cfg.ClientID != "" {
				fmt.Printf("%s (custom)\n", value)
			} else {
				fmt.Printf("%s (default)\n", value)
			}
		case "client-secret":
			value := cfg.EffectiveClientSecret()
			if value == "" {
				fmt.Println("(not set)")
			} else {
				masked := strings.Repeat("*", len(value))
				if len(value) > 4 {
					masked = strings.Repeat("*", len(value)-4) + value[len(value)-4:]
				}
				if cfg.ClientSecret != "" {
					fmt.Printf("%s (custom)\n", masked)
				} else {
					fmt.Printf("%s (default)\n", masked)
				}
			}
		case "scopes":
			value := cfg.EffectiveScopes()
			if cfg.Scopes != "" {
				fmt.Printf("%s (custom)\n", value)
			} else {
				fmt.Printf("%s (default)\n", value)
			}
		default:
			return fmt.Errorf("unknown config key: %s (valid: client-id, client-secret, scopes)", key)
		}

		return nil
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	rootCmd.AddCommand(configCmd)
}
