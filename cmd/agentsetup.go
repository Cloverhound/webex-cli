package cmd

import (
	"fmt"

	"github.com/Cloverhound/webex-cli/skill"
	"github.com/spf13/cobra"
)

// agentSetupCmd serves the setup skill from the npm package, because the npm
// registry is the one host cloud sandboxes reliably allow.
var agentSetupCmd = &cobra.Command{
	Use:   "agent-setup",
	Short: "Print setup instructions for an AI agent",
	Long:  "Print a skill that tells an AI agent how to install the CLI and sign in. Run it with npx -y @cloverhound/webex-cli agent-setup.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Fprint(cmd.OutOrStdout(), skill.Setup)
		return err
	},
}

func init() {
	rootCmd.AddCommand(agentSetupCmd)
}
