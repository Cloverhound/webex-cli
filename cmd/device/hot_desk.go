package device

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

var hotDeskCmd = &cobra.Command{
	Use:   "hot-desk",
	Short: "HotDesk commands",
}

func init() {
	cmd.DeviceCmd.AddCommand(hotDeskCmd)

	{ // list-sessions
		var orgId string
		var personId string
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "list-sessions",
			Short: "List Sessions",
			Long:  "List hot desk sessions.\n\nUse query parameters to filter the response.\nThe `orgId` parameter is for use by partner administrators acting on a managed organization.\nThe `personId` and `workspaceId` parameters are optional and are used to filter the response to only include sessions for a specific person or workspace.\nWhen used together they are used as an AND filter.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/hotdesk/sessions")
				req.QueryParam("orgId", orgId)
				req.QueryParam("personId", personId)
				req.QueryParam("workspaceId", workspaceId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List sessions in this organization. Only admin users of another organization (such as partners) may use this parameter.")
		cmd.Flags().StringVar(&personId, "person-id", "", "List sessions for this person.")
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "List sessions for this workspace.")
		hotDeskCmd.AddCommand(cmd)
	}

	{ // delete-session
		var sessionId string
		cmd := &cobra.Command{
			Use:   "delete-session",
			Short: "Delete Session",
			Long:  `Delete a hot desk session.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/hotdesk/sessions/{sessionId}")
				req.PathParam("sessionId", sessionId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&sessionId, "session-id", "", "The unique identifier for the hot desk session.")
		cmd.MarkFlagRequired("session-id")
		hotDeskCmd.AddCommand(cmd)
	}

}
