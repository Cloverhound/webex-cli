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

var resourceMembershipsCmd = &cobra.Command{
	Use:   "resource-memberships",
	Short: "ResourceMemberships commands",
}

func init() {
	cmd.AdminCmd.AddCommand(resourceMembershipsCmd)

	{ // list-group
		var licenseId string
		var personId string
		var personOrgId string
		var status string
		var max string
		cmd := &cobra.Command{
			Use:   "list-group",
			Short: "List Resource Group Memberships",
			Long: `Lists all resource group memberships for an organization.

Use query parameters to filter the response.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/resourceGroup/memberships")
				req.QueryParam("licenseId", licenseId)
				req.QueryParam("personId", personId)
				req.QueryParam("personOrgId", personOrgId)
				req.QueryParam("status", status)
				req.QueryParam("max", max)
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
		cmd.Flags().StringVar(&licenseId, "license-id", "", "List resource group memberships for a license, by ID.")
		cmd.Flags().StringVar(&personId, "person-id", "", "List resource group memberships for a person, by ID.")
		cmd.Flags().StringVar(&personOrgId, "person-org-id", "", "List resource group memberships for an organization, by ID.")
		cmd.Flags().StringVar(&status, "status", "", "Limit resource group memberships to a specific status.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of resource group memberships in the response.")
		resourceMembershipsCmd.AddCommand(cmd)
	}

	{ // get-group
		var resourceGroupMembershipId string
		cmd := &cobra.Command{
			Use:   "get-group",
			Short: "Get Resource Group Membership Details",
			Long:  "Shows details for a resource group membership, by ID.\n\nSpecify the resource group membership ID in the `resourceGroupMembershipId` URI parameter.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/resourceGroup/memberships/{resourceGroupMembershipId}")
				req.PathParam("resourceGroupMembershipId", resourceGroupMembershipId)
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
		cmd.Flags().StringVar(&resourceGroupMembershipId, "resource-group-membership-id", "", "The unique identifier for the resource group membership.")
		cmd.MarkFlagRequired("resource-group-membership-id")
		resourceMembershipsCmd.AddCommand(cmd)
	}

	{ // update-group
		var resourceGroupMembershipId string
		var resourceGroupId string
		var licenseId string
		var personId string
		var personOrgId string
		var status string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-group",
			Short: "Update a Resource Group Membership",
			Long:  "Updates a resource group membership, by ID.\n\nSpecify the resource group membership ID in the `resourceGroupMembershipId`\u00a0URI parameter.\n\nOnly the `resourceGroupId`\u00a0can be changed with this action. Resource group memberships with a `status`\u00a0of \"pending\" cannot be updated. For more information about resource group memberships, see the [Managing Hybrid Services](/docs/api/guides/managing-hybrid-services-licenses#webex-resource-groups) guide.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/resourceGroup/memberships/{resourceGroupMembershipId}")
				req.PathParam("resourceGroupMembershipId", resourceGroupMembershipId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("resourceGroupId", resourceGroupId)
					req.BodyString("licenseId", licenseId)
					req.BodyString("personId", personId)
					req.BodyString("personOrgId", personOrgId)
					req.BodyString("status", status)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&resourceGroupMembershipId, "resource-group-membership-id", "", "The unique identifier for the resource group membership.")
		cmd.MarkFlagRequired("resource-group-membership-id")
		cmd.Flags().StringVar(&resourceGroupId, "resource-group-id", "", "")
		cmd.Flags().StringVar(&licenseId, "license-id", "", "")
		cmd.Flags().StringVar(&personId, "person-id", "", "")
		cmd.Flags().StringVar(&personOrgId, "person-org-id", "", "")
		cmd.Flags().StringVar(&status, "status", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		resourceMembershipsCmd.AddCommand(cmd)
	}

	{ // list-group-v2
		var licenseId string
		var id string
		var orgId string
		var status string
		var typeVal string
		var max string
		cmd := &cobra.Command{
			Use:   "list-group-v2",
			Short: "List Resource Group Memberships V2",
			Long: `Lists all resource group memberships for an organization having filtering option based on entity type (User / Workspace).

Use query parameters to filter the response.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/resourceGroup/memberships/v2")
				req.QueryParam("licenseId", licenseId)
				req.QueryParam("id", id)
				req.QueryParam("orgId", orgId)
				req.QueryParam("status", status)
				req.QueryParam("type", typeVal)
				req.QueryParam("max", max)
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
		cmd.Flags().StringVar(&licenseId, "license-id", "", "List resource group memberships for a license, by ID.")
		cmd.Flags().StringVar(&id, "id", "", "List resource group memberships by ID.")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List resource group memberships for an organization, by ID.")
		cmd.Flags().StringVar(&status, "status", "", "Limit resource group memberships to a specific status.")
		cmd.Flags().StringVar(&typeVal, "type", "", "List resource group memberships for an organization, by type. If left blank it will include both User and Workspace type.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of resource group memberships in the response.")
		resourceMembershipsCmd.AddCommand(cmd)
	}

}
