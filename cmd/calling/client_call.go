package calling

import (
	"fmt"

	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/spf13/cobra"
)

// Ensure imports are used.
var _ = fmt.Sprintf
var _ = config.Token
var _ = output.Print

var clientCallCmd = &cobra.Command{
	Use:   "client-call",
	Short: "ClientCall commands",
}

func init() {
	cmd.CallingCmd.AddCommand(clientCallCmd)

	{ // get-org-ms-teams-settings
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-org-ms-teams-settings",
			Short: "Get an Organization's MS Teams Settings",
			Long:  "<div><Callout type=\"warning\">Not supported for Webex for Government (FedRAMP)</Callout></div>\n\nGet organization MS Teams settings.\n\nAt an organization level, MS Teams settings allow access to viewing the `HIDE WEBEX APP` and `PRESENCE SYNC` settings.\n\nTo retrieve an organization's MS Teams settings requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/settings/msTeams")
				req.QueryParam("orgId", orgId)
				if config.Paginate() {
					resp, statusCode, err := req.DoPaginated(true)
					if err != nil {
						return err
					}
					return output.Print(resp, statusCode)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve MS Teams settings for the organization.")
		clientCallCmd.AddCommand(cmd)
	}

	{ // update-org-ms-teams-setting
		var orgId string
		var settingName string
		var value bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-org-ms-teams-setting",
			Short: "Update an Organization's MS Teams Setting",
			Long:  "<div><Callout type=\"warning\">Not supported for Webex for Government (FedRAMP)</Callout></div>\n\nUpdate an MS Teams setting.\n\nMS Teams setting can be updated at the organization level.\n\nRequires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/settings/msTeams")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("settingName", settingName)
					req.BodyBool("value", value, cmd.Flags().Changed("value"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update MS Teams setting value for the organization.")
		cmd.Flags().StringVar(&settingName, "setting-name", "", "")
		cmd.Flags().BoolVar(&value, "value", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		clientCallCmd.AddCommand(cmd)
	}

}
