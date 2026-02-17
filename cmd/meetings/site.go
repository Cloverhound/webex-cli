package meetings

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

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Site commands",
}

func init() {
	cmd.MeetingsCmd.AddCommand(siteCmd)

	{ // get-meeting-common-settings-configuration
		var siteUrl string
		cmd := &cobra.Command{
			Use:   "get-meeting-common-settings-configuration",
			Short: "Get Meeting Common Settings Configuration",
			Long:  "Site administrators can use this API to get a list of functions, options, and privileges that are configured for their Webex service sites.\n\n* If `siteUrl` is specified, common settings of the meeting's configuration of the specified site will be queried; otherwise, the API will query from the site administrator's preferred site. All available Webex sites and preferred site of the user can be retrieved by [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/admin/meeting/config/commonSettings")
				req.QueryParam("siteUrl", siteUrl)
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
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API queries common settings of the meeting's configuration from. If not specified, the API will query from the site administrator's preferred site. All available Webex sites and the preferred site of the user can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		siteCmd.AddCommand(cmd)
	}

	{ // update-meeting-common-settings-configuration
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-meeting-common-settings-configuration",
			Short: "Update Meeting Common Settings Configuration",
			Long:  `Site administrators can use this API to update the option of features, options and privileges that are configured for their WebEx service sites.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PATCH", "/admin/meeting/config/commonSettings")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		siteCmd.AddCommand(cmd)
	}

}
