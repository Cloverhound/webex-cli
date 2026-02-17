package cc

import (
	"fmt"
	"strconv"
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
var _ = strconv.Itoa
var _ = strings.Join

var journeyCmd = &cobra.Command{
	Use:   "journey",
	Short: "Journey commands",
}

func init() {
	cmd.CcCmd.AddCommand(journeyCmd)

	{ // get-workspace
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "get-workspace",
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // update-workspace
		var workspaceId string
		var description string
		var name string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-workspace",
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // delete-workspace
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "delete-workspace",
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-template-searched-template-id
		var workspaceId string
		var templateId string
		cmd := &cobra.Command{
			Use:   "get-template-searched-template-id",
			Short: "Get A specific Template searched by template id",
			Long:  `Get Template details by template Id in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/profile-view-template/workspace-id/{workspaceId}/template-id/{templateId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("templateId", templateId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // update-profileviewtemplate
		var workspaceId string
		var templateId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-profileviewtemplate",
			Short: "Update existing ProfileViewTemplate",
			Long:  `Update existing Profile View Template in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PUT", "/admin/v1/api/profile-view-template/workspace-id/{workspaceId}/template-id/{templateId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("templateId", templateId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // delete-template-template-id
		var workspaceId string
		var templateId string
		cmd := &cobra.Command{
			Use:   "delete-template-template-id",
			Short: "Delete Template by template Id",
			Long:  `Delete Template By template id in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/admin/v1/api/profile-view-template/workspace-id/{workspaceId}/template-id/{templateId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("templateId", templateId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // delete-person-id
		var workspaceId string
		var personId string
		cmd := &cobra.Command{
			Use:   "delete-person-id",
			Short: "Delete specific Person by id",
			Long:  `Delete Person Details searched by Person id in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/admin/v1/api/person/workspace-id/{workspaceId}/person-id/{personId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("personId", personId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&personId, "person-id", "", "(Required) Person ID")
		cmd.MarkFlagRequired("person-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // create-remove-replace-person
		var workspaceId string
		var personId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-remove-replace-person",
			Short: "Add/Remove/Replace details of a Person.",
			Long: `The Patch Api can be used to add/remove identities(email, phone, customerId) or replace firstName and lastName of an Individual. We support only add, replace and remove operations. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope. 

For a more information on Patch Requests, see this  [JSON PATCH guide](https://jsonpatch.com)`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PATCH", "/admin/v1/api/person/workspace-id/{workspaceId}/person-id/{personId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&personId, "person-id", "", "(Required) Person ID")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // create-one-more-identities-person
		var workspaceId string
		var personId string
		var phone []string
		var email []string
		var temporaryId []string
		var customerId []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-one-more-identities-person",
			Short: "Add one/more Identities to a person",
			Long:  `This Patch Api can be used to add identities(email, phone, customerId) to a person. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PATCH", "/admin/v1/api/person/add-identities/workspace-id/{workspaceId}/person-id/{personId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("personId", personId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("phone", phone)
					req.BodyStringSlice("email", email)
					req.BodyStringSlice("temporaryId", temporaryId)
					req.BodyStringSlice("customerId", customerId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&personId, "person-id", "", "(Required) Person ID")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringSliceVar(&phone, "phone", nil, "")
		cmd.Flags().StringSliceVar(&email, "email", nil, "")
		cmd.Flags().StringSliceVar(&temporaryId, "temporary-id", nil, "")
		cmd.Flags().StringSliceVar(&customerId, "customer-id", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-wxcc-subscription
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "get-wxcc-subscription",
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // create-wxcc-subscription
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "create-wxcc-subscription",
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // delete-wxcc-subscription
		var workspaceId string
		cmd := &cobra.Command{
			Use:   "delete-wxcc-subscription",
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-all-workspaces
		var filter string
		var sortBy string
		var sort string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "get-all-workspaces",
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
		cmd.Flags().StringVar(&filter, "filter", "", "Optional filter which can be applied to the elements to be fetched.   This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see [this reference](https://developer.here.com/docs/data-client-library/dev_guide/client/rsql.html). For a list of supported operators, see this [syntax guide](https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference).")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort By Field")
		cmd.Flags().StringVar(&sort, "sort", "", "Sort direction")
		cmd.Flags().StringVar(&page, "page", "", "Index of the page of results to be fetched.  Results are returned in blocks of pageSize elements. This parameter specifies which page number to retrieve.The page numbering starts with 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Number of items to be displayed on a page.")
		journeyCmd.AddCommand(cmd)
	}

	{ // create-workspace
		var organizationId string
		var description string
		var name string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-workspace",
			Short: "Create Workspace",
			Long:  `Create Workspace in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_write or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/admin/v1/api/workspace")
				req.QueryParam("organizationId", organizationId)
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
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "Organization ID")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-all-template
		var workspaceId string
		var filter string
		var sort string
		var sortBy string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "get-all-template",
			Short: "Get All Template Details.",
			Long:  `Get Template details by Organization Id and workspaceId in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/profile-view-template/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
				req.QueryParam("filter", filter)
				req.QueryParam("sort", sort)
				req.QueryParam("sortBy", sortBy)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&filter, "filter", "", "Optional filter which can be applied to the elements to be fetched.   This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see [this reference](https://developer.here.com/docs/data-client-library/dev_guide/client/rsql.html). For a list of supported operators, see this [syntax guide](https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference).")
		cmd.Flags().StringVar(&sort, "sort", "", "Sort direction")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort By Field")
		cmd.Flags().StringVar(&page, "page", "", "Index of the page of results to be fetched.  Results are returned in blocks of pageSize elements. This parameter specifies which page number to retrieve.The page numbering starts with 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Number of items to be displayed on a page.")
		journeyCmd.AddCommand(cmd)
	}

	{ // create-template
		var workspaceId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-template",
			Short: "Create Template",
			Long:  `Creates a Profile View Template in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/admin/v1/api/profile-view-template/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-all-person
		var workspaceId string
		var personId string
		var filter string
		var sortBy string
		var sort string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "get-all-person",
			Short: "Get all or a specific Person Details",
			Long: `Get Person Details in JDS. If personId is provided by query parameter, this returns the Person whose personId matches the parameter.
If not, this will return ALL the Persons within the organization and workspace. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/person/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
				req.QueryParam("personId", personId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&personId, "person-id", "", "Person ID")
		cmd.Flags().StringVar(&filter, "filter", "", "Optional filter which can be applied to the elements to be fetched.   This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see [this reference](https://developer.here.com/docs/data-client-library/dev_guide/client/rsql.html). For a list of supported operators, see this [syntax guide](https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference).")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort By Field")
		cmd.Flags().StringVar(&sort, "sort", "", "Sort direction")
		cmd.Flags().StringVar(&page, "page", "", "Index of the page of results to be fetched.  Results are returned in blocks of pageSize elements. This parameter specifies which page number to retrieve. The page numbering starts with 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Number of items to be displayed on a page")
		journeyCmd.AddCommand(cmd)
	}

	{ // create-person
		var workspaceId string
		var firstName string
		var lastName string
		var phone []string
		var email []string
		var temporaryId []string
		var customerId []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-person",
			Short: "Create a Person",
			Long:  `Creates a Person in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/admin/v1/api/person/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("firstName", firstName)
					req.BodyString("lastName", lastName)
					req.BodyStringSlice("phone", phone)
					req.BodyStringSlice("email", email)
					req.BodyStringSlice("temporaryId", temporaryId)
					req.BodyStringSlice("customerId", customerId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&firstName, "first-name", "", "")
		cmd.Flags().StringVar(&lastName, "last-name", "", "")
		cmd.Flags().StringSliceVar(&phone, "phone", nil, "")
		cmd.Flags().StringSliceVar(&email, "email", nil, "")
		cmd.Flags().StringSliceVar(&temporaryId, "temporary-id", nil, "")
		cmd.Flags().StringSliceVar(&customerId, "customer-id", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // merges-identities-primary-identity
		var workspaceId string
		var primaryPersonId string
		var personIdsToMerge []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "merges-identities-primary-identity",
			Short: "Merges Identities to a Primary Identity.",
			Long:  `Merges one/more Identities to a **Primary** Individual in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/admin/v1/api/person/merge/workspace-id/{workspaceId}/primary-person-id/{primaryPersonId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("primaryPersonId", primaryPersonId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("personIdsToMerge", personIdsToMerge)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&primaryPersonId, "primary-person-id", "", "(Required) Primary Person ID")
		cmd.MarkFlagRequired("primary-person-id")
		cmd.Flags().StringSliceVar(&personIdsToMerge, "person-ids-to-merge", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // creates-merges-aliases-individual-jds
		var workspaceId string
		var firstName string
		var lastName string
		var phone []string
		var email []string
		var temporaryId []string
		var customerId []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "creates-merges-aliases-individual-jds",
			Short: "Creates or merges aliases to an Individual in JDS.",
			Long:  `Merges one/more aliases in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/admin/v1/api/person/merge-identities/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("firstName", firstName)
					req.BodyString("lastName", lastName)
					req.BodyStringSlice("phone", phone)
					req.BodyStringSlice("email", email)
					req.BodyStringSlice("temporaryId", temporaryId)
					req.BodyStringSlice("customerId", customerId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&firstName, "first-name", "", "")
		cmd.Flags().StringVar(&lastName, "last-name", "", "")
		cmd.Flags().StringSliceVar(&phone, "phone", nil, "")
		cmd.Flags().StringSliceVar(&email, "email", nil, "")
		cmd.Flags().StringSliceVar(&temporaryId, "temporary-id", nil, "")
		cmd.Flags().StringSliceVar(&customerId, "customer-id", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-all-actions
		var workspaceId string
		var sortBy string
		var sort string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "get-all-actions",
			Short: "Get all Journey Actions",
			Long:  `Get all Journey Actions in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/journey-actions/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort By Field")
		cmd.Flags().StringVar(&sort, "sort", "", "Sort direction")
		cmd.Flags().StringVar(&page, "page", "", "Index of the page of results to be fetched.  Results are returned in blocks of pageSize elements. This parameter specifies which page number to retrieve. The page numbering starts with 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Number of items to be displayed on a page.")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-template-searched-template-name
		var workspaceId string
		var templateName string
		cmd := &cobra.Command{
			Use:   "get-template-searched-template-name",
			Short: "Get A specific Template searched by template name",
			Long:  `Get Template details by template Name in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/profile-view-template/workspace-id/{workspaceId}/template-name/{templateName}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("templateName", templateName)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&templateName, "template-name", "", "(Required) Template Name")
		cmd.MarkFlagRequired("template-name")
		journeyCmd.AddCommand(cmd)
	}

	{ // search-identity-aliases
		var workspaceId string
		var aliases string
		var sortBy string
		var sort string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "search-identity-aliases",
			Short: "Search for an Identity via aliases",
			Long:  `Get one or more Person Details searched by aliases in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_read or cjp:config_write scopes. Multiple aliases should be separated by a comma.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/person/workspace-id/{workspaceId}/aliases/{aliases}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("aliases", aliases)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&aliases, "aliases", "", "(Required) Aliases to search for. Multiple aliases should be separated by a comma.    In case the alias(es) contain(s) non-uri-encodable characters, eg: '+', '>' etc, you can URL-encode the same and then pass it as parameter.")
		cmd.MarkFlagRequired("aliases")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort By Field")
		cmd.Flags().StringVar(&sort, "sort", "", "Sort direction")
		cmd.Flags().StringVar(&page, "page", "", "Index of the page of results to be fetched.  Results are returned in blocks of pageSize elements. This parameter specifies which page number to retrieve.The page numbering starts with 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Number of items to be displayed on a page.")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-all-actions-template
		var workspaceId string
		var templateId string
		cmd := &cobra.Command{
			Use:   "get-all-actions-template",
			Short: "Get all Journey Actions for a template",
			Long:  `Get all Journey Actions for a template in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/journey-actions/workspace-id/{workspaceId}/template-id/{templateId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("templateId", templateId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // create-action
		var workspaceId string
		var templateId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-action",
			Short: "Create a new  Journey Action.",
			Long:  `Create a new Journey Action in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/admin/v1/api/journey-actions/workspace-id/{workspaceId}/template-id/{templateId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("templateId", templateId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-action-name
		var workspaceId string
		var templateId string
		var actionName string
		cmd := &cobra.Command{
			Use:   "get-action-name",
			Short: "Get specific Journey Action By Name",
			Long:  `Get specific Journey Action By Name in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/journey-actions/workspace-id/{workspaceId}/template-id/{templateId}/action-name/{actionName}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("templateId", templateId)
				req.PathParam("actionName", actionName)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		cmd.Flags().StringVar(&actionName, "action-name", "", "(Required) Action Name")
		cmd.MarkFlagRequired("action-name")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-action-actionid
		var workspaceId string
		var templateId string
		var actionId string
		cmd := &cobra.Command{
			Use:   "get-action-actionid",
			Short: "Get specific Journey Action By ActionId",
			Long:  `Get specific Journey Action By ActionId in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/journey-actions/workspace-id/{workspaceId}/template-id/{templateId}/action-id/{actionId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("templateId", templateId)
				req.PathParam("actionId", actionId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		cmd.Flags().StringVar(&actionId, "action-id", "", "(Required) Action ID")
		cmd.MarkFlagRequired("action-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // update-action
		var workspaceId string
		var templateId string
		var actionId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-action",
			Short: "Update existing Journey Action.",
			Long:  `Update existing Journey Action in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PUT", "/admin/v1/api/journey-actions/workspace-id/{workspaceId}/template-id/{templateId}/action-id/{actionId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("templateId", templateId)
				req.PathParam("actionId", actionId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		cmd.Flags().StringVar(&actionId, "action-id", "", "(Required) Action ID")
		cmd.MarkFlagRequired("action-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // delete-action-configuration-actionid
		var workspaceId string
		var templateId string
		var actionId string
		cmd := &cobra.Command{
			Use:   "delete-action-configuration-actionid",
			Short: "Delete Journey Action configuration By ActionId.",
			Long:  `Delete Journey Action configuration By ActionId in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/admin/v1/api/journey-actions/workspace-id/{workspaceId}/template-id/{templateId}/action-id/{actionId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("templateId", templateId)
				req.PathParam("actionId", actionId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		cmd.Flags().StringVar(&actionId, "action-id", "", "(Required) Action ID")
		cmd.MarkFlagRequired("action-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-historic-profile-view-template-name
		var workspaceId string
		var personId string
		var templateName string
		cmd := &cobra.Command{
			Use:   "get-historic-profile-view-template-name",
			Short: "Historic Progressive Profile View using Template Name.",
			Long:  `Get Historic Progressive Profile View in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/admin/v1/api/progressive-profile-view/workspace-id/{workspaceId}/person-id/{personId}/template-name/{templateName}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("personId", personId)
				req.PathParam("templateName", templateName)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&personId, "person-id", "", "(Required) Person ID")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&templateName, "template-name", "", "(Required) Template Name")
		cmd.MarkFlagRequired("template-name")
		journeyCmd.AddCommand(cmd)
	}

	{ // delete-one-more-identities-person
		var workspaceId string
		var personId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "delete-one-more-identities-person",
			Short: "Remove one/more Identities from a person",
			Long:  `This Patch Api can be used to remove identities(email, phone, customerId) from a person. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjp:config_write scope.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PATCH", "/v1/api/person/remove-identities/workspace-id/{workspaceId}/person-id/{personId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&personId, "person-id", "", "(Required) Person ID")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-historic-profile-view
		var workspaceId string
		var personId string
		var templateId string
		cmd := &cobra.Command{
			Use:   "get-historic-profile-view",
			Short: "Historic Progressive Profile View.",
			Long:  `Get Historic Progressive Profile View in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/v1/api/progressive-profile-view/workspace-id/{workspaceId}/person-id/{personId}/template-id/{templateId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("personId", personId)
				req.PathParam("templateId", templateId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&personId, "person-id", "", "(Required) Person ID")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-historic-profile-view-identity-template-name
		var workspaceId string
		var identity string
		var templateName string
		cmd := &cobra.Command{
			Use:   "get-historic-profile-view-identity-template-name",
			Short: "Historic Progressive Profile View By Identity and Template Name.",
			Long:  `Get Historic Progressive Profile View in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/v1/api/progressive-profile-view/workspace-id/{workspaceId}/identity/{identity}/template-name/{templateName}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("identity", identity)
				req.PathParam("templateName", templateName)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&identity, "identity", "", "(Required) Identity")
		cmd.MarkFlagRequired("identity")
		cmd.Flags().StringVar(&templateName, "template-name", "", "(Required) Template Name")
		cmd.MarkFlagRequired("template-name")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-historic-profile-view-identity-template-id
		var workspaceId string
		var identity string
		var templateId string
		cmd := &cobra.Command{
			Use:   "get-historic-profile-view-identity-template-id",
			Short: "Historic Progressive Profile View By Identity and Template Id",
			Long:  `Get Historic Progressive Profile View in JDS. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/v1/api/progressive-profile-view/workspace-id/{workspaceId}/identity/{identity}/template-id/{templateId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("identity", identity)
				req.PathParam("templateId", templateId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&identity, "identity", "", "(Required) identity")
		cmd.MarkFlagRequired("identity")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // stream-profile-views-template-name
		var workspaceId string
		var identity string
		var templateName string
		cmd := &cobra.Command{
			Use:   "stream-profile-views-template-name",
			Short: "Stream Progressive profile Views By Template Name",
			Long:  `Real-time streaming enables API consumers to listen for Progressive profile Views as it created/updated as part of the Journey; these may be transformed, value-added/enriched, and ready to be consumed or forwarded to another destination. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/v1/api/progressive-profile-view/stream/workspace-id/{workspaceId}/identity/{identity}/template-name/{templateName}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("identity", identity)
				req.PathParam("templateName", templateName)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&identity, "identity", "", "(Required) Identity to search Progressive Profile View for.    In case the identity contains non-uri-encodable characters, eg: '+', '>' etc, you can URL-encode the same and then pass it as parameter.")
		cmd.MarkFlagRequired("identity")
		cmd.Flags().StringVar(&templateName, "template-name", "", "(Required) Template Name")
		cmd.MarkFlagRequired("template-name")
		journeyCmd.AddCommand(cmd)
	}

	{ // stream-profile-views
		var workspaceId string
		var identity string
		var templateId string
		cmd := &cobra.Command{
			Use:   "stream-profile-views",
			Short: "Stream Progressive profile Views",
			Long:  `Real-time streaming enables API consumers to listen for Progressive profile Views as it created/updated as part of the Journey; these may be transformed, value-added/enriched, and ready to be consumed or forwarded to another destination. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/v1/api/progressive-profile-view/stream/workspace-id/{workspaceId}/identity/{identity}/template-id/{templateId}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("identity", identity)
				req.PathParam("templateId", templateId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&identity, "identity", "", "(Required) Identity to search Progressive Profile View for.    In case the identity contains non-uri-encodable characters, eg: '+', '>' etc, you can URL-encode the same and then pass it as parameter.")
		cmd.MarkFlagRequired("identity")
		cmd.Flags().StringVar(&templateId, "template-id", "", "(Required) Template ID")
		cmd.MarkFlagRequired("template-id")
		journeyCmd.AddCommand(cmd)
	}

	{ // get-historic-events
		var workspaceId string
		var identity string
		var sortBy string
		var sort string
		var filter string
		var data string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "get-historic-events",
			Short: "Historic Journey Events",
			Long:  `Getting Historic Customer Journey Events from Pinot. These events are append-only, immutable data ledger that can be queried to retrieve snapshot of latest events that moment in time or historically to play-back events as they occurred to understand or analyze Journeys using ML/AI models. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes `,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/v1/api/events/workspace-id/{workspaceId}")
				req.PathParam("workspaceId", workspaceId)
				req.QueryParam("identity", identity)
				req.QueryParam("sortBy", sortBy)
				req.QueryParam("sort", sort)
				req.QueryParam("filter", filter)
				req.QueryParam("data", data)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&identity, "identity", "", "Identity to search events for.    In case the identity contains non-uri-encodable characters, eg: '+', '>' etc, you can URL-encode the same and then pass it as parameter.")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "sort By Field")
		cmd.Flags().StringVar(&sort, "sort", "", "sort direction")
		cmd.Flags().StringVar(&filter, "filter", "", "Optional filter which can be applied to the elements to be fetched.  This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see [this reference](https://developer.here.com/docs/data-client-library/dev_guide/client/rsql.html). For a list of supported operators, see this  [syntax guide](https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference).")
		cmd.Flags().StringVar(&data, "data", "", "Optional filter on data filed which can be applied to the elements to be fetched.  This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see [this reference](https://developer.here.com/docs/data-client-library/dev_guide/client/rsql.html). For a list of supported operators, see this  [syntax guide](https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference).")
		cmd.Flags().StringVar(&page, "page", "", "Index of the page of results to be fetched.  Results are returned in blocks of pageSize elements. This parameter specifies which page number to retrieve.The page numbering starts with 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Number of items to be displayed on a page.")
		journeyCmd.AddCommand(cmd)
	}

	{ // stream-events-identity
		var workspaceId string
		var identity string
		var filter string
		var data string
		cmd := &cobra.Command{
			Use:   "stream-events-identity",
			Short: "Stream Events By Identity",
			Long:  `Real-time streaming enables API consumers to listen for events as it arrives as part of the Journey; these may be transformed, value-added/enriched, and ready to be consumed or forwarded to another destination. Optionally accepts filter and data parameters slice/dice further. It  requires the appropriate cjds:admin_org_read or cjds:admin_org_write scopes or cjp:config_read or cjp:config_write scopes`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/v1/api/events/stream/workspace-id/{workspaceId}/identity/{identity}")
				req.PathParam("workspaceId", workspaceId)
				req.PathParam("identity", identity)
				req.QueryParam("filter", filter)
				req.QueryParam("data", data)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&identity, "identity", "", "(Required) Person Identity.    In case the identity contains non-uri-encodable characters, eg: '+', '>' etc, you can URL-encode the same and then pass it as parameter.")
		cmd.MarkFlagRequired("identity")
		cmd.Flags().StringVar(&filter, "filter", "", "Optional filter which can be applied to the elements to be fetched.  This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see [this reference](https://developer.here.com/docs/data-client-library/dev_guide/client/rsql.html). For a list of supported operators, see this  [syntax guide](https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference).")
		cmd.Flags().StringVar(&data, "data", "", "Optional filter on data filed which can be applied to the elements to be fetched.  This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see [this reference](https://developer.here.com/docs/data-client-library/dev_guide/client/rsql.html). For a list of supported operators, see this [syntax guide](https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference).")
		journeyCmd.AddCommand(cmd)
	}

	{ // event-posting
		var workspaceId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "event-posting",
			Short: "Journey Event Posting",
			Long:  `Journey Event Posting Api.API accepts events that describe what occurred, when, and by whom on every interaction across touch points and applications. Data Ingestion is based on Cloud Events specification for describing event data in a common way. API accepts data in the form of POST with support for Header based authorization. Use the cjp scope if you have a contact center license; otherwise, use the cjds scope. It requires the appropriate cjds:admin_org_write or cjp:config_write scopes`,
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "(Required) Workspace ID")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		journeyCmd.AddCommand(cmd)
	}

}
