package admin

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

var serviceAppsCmd = &cobra.Command{
	Use:   "service-apps",
	Short: "ServiceApps commands",
}

func init() {
	cmd.AdminCmd.AddCommand(serviceAppsCmd)

	{ // create-access-token
		var applicationId string
		var clientId string
		var clientSecret string
		var targetOrgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-access-token",
			Short: "Create Service App Access Token",
			Long: `Retrieves an organization specific token pair for an already authorized Service App. Service Apps use machine accounts to make API calls on behalf of an organization, independent of individual user life cycles.

This endpoint allows you to programmatically retrieve access and refresh tokens after a Full Admin has authorized your Service App in Control Hub.

To call this endpoint, you need an integration with the spark:applications_token scope.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/applications/{applicationId}/token")
				req.PathParam("applicationId", applicationId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("clientId", clientId)
					req.BodyString("clientSecret", clientSecret)
					req.BodyString("targetOrgId", targetOrgId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&applicationId, "application-id", "", "The unique identifier of the Service App")
		cmd.MarkFlagRequired("application-id")
		cmd.Flags().StringVar(&clientId, "client-id", "", "")
		cmd.Flags().StringVar(&clientSecret, "client-secret", "", "")
		cmd.Flags().StringVar(&targetOrgId, "target-org-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		serviceAppsCmd.AddCommand(cmd)
	}

}
