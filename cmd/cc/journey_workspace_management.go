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

var journeyWorkspaceManagementCmd = &cobra.Command{
	Use:   "journey-workspace-management",
	Short: "JourneyWorkspaceManagement commands",
}

func init() {
	cmd.CcCmd.AddCommand(journeyWorkspaceManagementCmd)

	{ // get
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Workspace",
			Long:  `Get workspace details. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/workspace/workspace-id/{workspaceId}")
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
		journeyWorkspaceManagementCmd.AddCommand(cmd)
	}

	{ // update
		var workspaceId string
		var description string
		var name string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update Workspace",
			Long:  `Update workspace by Id. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_write or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PUT", "/admin/v1/api/workspace/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("description", description)
					req.BodyString("name", name)
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
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyWorkspaceManagementCmd.AddCommand(cmd)
	}

	{ // delete
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete Workspace",
			Long:  `Delete Workspace By Id. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_write or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/admin/v1/api/workspace/workspace-id/{workspaceId}")
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
		journeyWorkspaceManagementCmd.AddCommand(cmd)
	}

	{ // get-all
		var filter string
		var sortBy string
		var sort string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "get-all",
			Short: "Get All Workspaces",
			Long:  `Get All Workspaces. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/workspace")
				req.QueryParam("filter", filter)
				req.QueryParam("sortBy", sortBy)
				req.QueryParam("sort", sort)
				req.QueryParam("page", page)
				req.QueryParam("pageSize", pageSize)
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
		cmd.Flags().StringVar(&filter, "filter", "", "Optional filter which can be applied to the elements to be fetched.   This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see [this reference](https://developer.here.com/documentation/data-client-library/dev_guide/client/rsql.html). For a list of supported operators, see this [syntax guide](https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference).")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort By Field")
		cmd.Flags().StringVar(&sort, "sort", "", "Sort direction")
		cmd.Flags().StringVar(&page, "page", "", "Index of the page of results to be fetched.  Results are returned in blocks of pageSize elements. This parameter specifies which page number to retrieve.The page numbering starts with 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Number of items to be displayed on a page.")
		journeyWorkspaceManagementCmd.AddCommand(cmd)
	}

	{ // create
		var description string
		var name string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create Workspace",
			Long:  `Create Workspace in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope.It requires the appropriate cjds:admin_org_write or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/admin/v1/api/workspace")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("description", description)
					req.BodyString("name", name)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyWorkspaceManagementCmd.AddCommand(cmd)
	}

}
