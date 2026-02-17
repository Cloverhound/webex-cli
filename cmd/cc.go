package cmd

import (
	"github.com/spf13/cobra"
)

var CcCmd = &cobra.Command{
	Use:     "cc",
	Aliases: []string{"contact-center"},
	Short:   "Webex Contact Center APIs",
	Long:    `Manage Webex Contact Center resources — agents, queues, tasks, routing, and more.`,
}

func init() {
	rootCmd.AddCommand(CcCmd)
}
