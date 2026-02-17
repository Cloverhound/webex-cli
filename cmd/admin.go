package cmd

import (
	"github.com/spf13/cobra"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Webex Admin APIs",
	Long:  `Manage Webex Admin resources — people, organizations, licenses, roles, and more.`,
}

func init() {
	rootCmd.AddCommand(AdminCmd)
}
