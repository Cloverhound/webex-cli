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

var resourceGroupsCmd = &cobra.Command{
	Use:   "resource-groups",
	Short: "ResourceGroups commands",
}

func init() {
	cmd.AdminCmd.AddCommand(resourceGroupsCmd)

	{ // list
		var orgId string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Resource Groups",
			Long: `List resource groups.

Use query parameters to filter the response.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/resourceGroups")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List resource groups in this organization. Only admin users of another organization (such as partners) may use this parameter.")
		resourceGroupsCmd.AddCommand(cmd)
	}

	{ // get
		var resourceGroupId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Resource Group Details",
			Long:  "Shows details for a resource group, by ID.\n\nSpecify the resource group ID in the `resourceGroupId` parameter in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/resourceGroups/{resourceGroupId}")
				req.PathParam("resourceGroupId", resourceGroupId)
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
		cmd.Flags().StringVar(&resourceGroupId, "resource-group-id", "", "The unique identifier for the resource group.")
		cmd.MarkFlagRequired("resource-group-id")
		resourceGroupsCmd.AddCommand(cmd)
	}

}
