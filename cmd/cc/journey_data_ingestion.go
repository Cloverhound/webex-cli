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

var journeyDataIngestionCmd = &cobra.Command{
	Use:   "journey-data-ingestion",
	Short: "JourneyDataIngestion commands",
}

func init() {
	cmd.CcCmd.AddCommand(journeyDataIngestionCmd)

	{ // event-posting
		var workspaceId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "event-posting",
			Short: "Journey Event Posting",
			Long: `Journey Event Posting API accepts events that describe what occurred, when, and by whom on every interaction across touch points and applications. Data Ingestion is based on Cloud Events specification for describing event data in a common way. API accepts data in the form of POST with support for Header based authorization. 

Role and Scope: Requires id full admin role with cjds:admin_org_write or any role with cjp:config_write scope.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/publish/v1/api/event")
				req.QueryParam("workspaceId", workspaceId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "Workspace ID")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyDataIngestionCmd.AddCommand(cmd)
	}

}
