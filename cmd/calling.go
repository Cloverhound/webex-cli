package cmd

import (
	"github.com/spf13/cobra"
)

var CallingCmd = &cobra.Command{
	Use:   "calling",
	Short: "Webex Cloud Calling APIs",
	Long:  `Manage Webex Cloud Calling resources — telephony config, devices, users, locations, and more.`,
}

func init() {
	rootCmd.AddCommand(CallingCmd)
}
