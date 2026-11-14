package cmd

import (
	"github.com/spf13/cobra"
)

var CcCmd = &cobra.Command{
	Use:     "cc",
	Aliases: []string{"contact-center"},
	Short:   "Webex Contact Center APIs",
	Long: `Manage Webex Contact Center resources — agents, queues, tasks, routing, and more.

Requests go to the organization's Contact Center data center: --cc-region, else
$WEBEX_CC_REGION, else the region saved with 'webex config set cc-region <region>',
else us1. Regions: us1, eu1, eu2, anz1, ca1, jp1, sg1, in1.`,
}

func init() {
	rootCmd.AddCommand(CcCmd)
}
