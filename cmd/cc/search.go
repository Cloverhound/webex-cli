package cc

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

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search commands",
}

func init() {
	cmd.CcCmd.AddCommand(searchCmd)

	{ // search-tasks
		var orgId string
		var trackingId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "search-tasks",
			Short: "Search tasks",
			Long: `The /search API is a GraphQL endpoint that enables customers to fetch data from Webex Contact Center.

Mandatory parameters are FROM and TO, which accept datetime in epoch format. The FROM parameter cannot be older than 36 months from the current time. The TO parameter, if given as a future time, will be set to the current time. Optional parameters such as filter and aggregation are accepted for each query.

Response Compression: For this API, response compression using gzip can be enabled by including the 'Accept-Encoding' header in the request with its value as 'gzip'. The response will be compressed only if its size exceeds 1 MB. If the header is not present in the request or if gzip is not listed as one of the encodings in the header's value (comma-separated encodings), then the API response will not be compressed, impacting latency as observed from clients.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/search")
				req.QueryParam("orgId", orgId)
				req.Header("TrackingId", trackingId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID to use for this operation. If unspecified, inferred from token. Token must have permission to interact with this organization.")
		cmd.Flags().StringVar(&trackingId, "tracking-id", "", "Tracking ID to use for this operation, for traceability, debugging, and error reporting purposes. ")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		searchCmd.AddCommand(cmd)
	}

}
