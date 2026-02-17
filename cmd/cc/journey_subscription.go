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

var journeySubscriptionCmd = &cobra.Command{
	Use:   "journey-subscription",
	Short: "JourneySubscription commands",
}

func init() {
	cmd.CcCmd.AddCommand(journeySubscriptionCmd)

	{ // get-wxcc
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "get-wxcc",
			Short: "Get WXCC Subscription",
			Long:  `Get WXCC Subscription in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/wxcc-subscription/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
				if config.Paginate() {
					resp, statusCode, err := req.DoPaginated(false)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		journeySubscriptionCmd.AddCommand(cmd)
	}

	{ // create-wxcc
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "create-wxcc",
			Short: "Create WXCC Subscription",
			Long:  `Create WXCC Subscription in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_write or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/admin/v1/api/wxcc-subscription/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		journeySubscriptionCmd.AddCommand(cmd)
	}

	{ // delete-wxcc
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "delete-wxcc",
			Short: "Delete WXCC Subscription",
			Long:  `Delete WXCC Subscription in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_write or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/admin/v1/api/wxcc-subscription/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		journeySubscriptionCmd.AddCommand(cmd)
	}

}
