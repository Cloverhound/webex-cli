package cmd

import (
	"github.com/spf13/cobra"
)

var MessagingCmd = &cobra.Command{
	Use:     "messaging",
	Aliases: []string{"msg"},
	Short:   "Webex Messaging APIs",
	Long:    `Manage Webex Messaging — rooms, messages, teams, webhooks, and more.`,
}

func init() {
	rootCmd.AddCommand(MessagingCmd)
}
