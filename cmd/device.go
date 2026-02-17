package cmd

import (
	"github.com/spf13/cobra"
)

var DeviceCmd = &cobra.Command{
	Use:     "device",
	Aliases: []string{"devices"},
	Short:   "Webex Device APIs",
	Long:    `Manage Webex devices — workspaces, device configurations, xAPI, and more.`,
}

func init() {
	rootCmd.AddCommand(DeviceCmd)
}
