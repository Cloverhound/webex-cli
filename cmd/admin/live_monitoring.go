package admin

import (
	"fmt"
	"strings"

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
var _ = strings.Join

var liveMonitoringCmd = &cobra.Command{
	Use:   "live-monitoring",
	Short: "LiveMonitoring commands",
}

func init() {
	cmd.AdminCmd.AddCommand(liveMonitoringCmd)

	{ // get-meeting-metrics-categorized-country
		var siteIds []string
		var siteUrl string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "get-meeting-metrics-categorized-country",
			Short: "Get Live Meeting metrics categorized by Country",
			Long:  "Retrieve live meeting metrics categorized by country for a specific meeting site or for all meeting sites owned by the customer. \n\nTo retrieve live monitoring information, you must use an administrator token with the `analytics:read_all` [scope](/docs/integrations#scopes). The authenticated user must be a read-only or full administrator of the organization to which the meeting belongs and must not be an external administrator.\n\nTo use this endpoint, the org needs to be licensed for the Webex Pro Pack.\n\nA rate limit of one API call per minute applies to each customer organization",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/livemonitoring/liveMeetingsByCountry")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("siteIds", siteIds)
					req.BodyString("siteUrl", siteUrl)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringSliceVar(&siteIds, "site-ids", nil, "")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		liveMonitoringCmd.AddCommand(cmd)
	}

}
