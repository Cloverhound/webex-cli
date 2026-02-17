package cmd

import (
	"github.com/spf13/cobra"
)

var MeetingsCmd = &cobra.Command{
	Use:     "meetings",
	Aliases: []string{"meeting"},
	Short:   "Webex Meetings APIs",
	Long:    `Manage Webex Meetings — scheduling, participants, recordings, and more.`,
}

func init() {
	rootCmd.AddCommand(MeetingsCmd)
}
