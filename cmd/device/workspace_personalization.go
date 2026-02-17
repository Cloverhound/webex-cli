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

var workspacePersonalizationCmd = &cobra.Command{
	Use:   "workspace-personalization",
	Short: "WorkspacePersonalization commands",
}

func init() {
	cmd.DeviceCmd.AddCommand(workspacePersonalizationCmd)

	{ // personalize
		var workspaceId string
		var email string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "personalize",
			Short: "Personalize a Workspace",
			Long:  "Initializes the personalization for a given workspace for the user email provided.\n\nThe personalization process is asynchronous and thus a background task is created when this endpoint is invoked.\nAfter successful invocation of this endpoint a personalization task status URL will be returned in the `Location` header, which will point to the [Get Personalization Task](/docs/api/v1/workspace-personalization/get-personalization-task) endpoint for this workspace.\nThe task should be completed in approximately 30 seconds.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/workspaces/{workspaceId}/personalize")
				req.PathParam("workspaceId", workspaceId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("email", email)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "A unique identifier for the workspace.")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&email, "email", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		workspacePersonalizationCmd.AddCommand(cmd)
	}

	{ // get-task
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "get-task",
			Short: "Get Personalization Task",
			Long:  "Returns the status of a personalization task for a given workspace.\n\nWhilst in progress the endpoint will return `Accepted` and provide a `Retry-After` header indicating the number of seconds a client should wait before retrying.\n\nUpon completion of the task, the endpoint will return `OK` with a body detailing if the personalization was successful and an error description if appropriate.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/workspaces/{workspaceId}/personalizationTask")
				req.PathParam("workspaceId", workspaceId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "A unique identifier for the workspace.")
		cmd.MarkFlagRequired("workspace-id")
		workspacePersonalizationCmd.AddCommand(cmd)
	}

}
