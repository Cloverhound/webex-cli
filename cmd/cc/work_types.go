package cc

import (
	"fmt"
	"strconv"

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

var workTypesCmd = &cobra.Command{
	Use:   "work-types",
	Short: "WorkTypes commands",
}

func init() {
	cmd.CcCmd.AddCommand(workTypesCmd)

	{ // create
		var orgid string
		var active bool
		var name string
		var workTypeCode string
		var organizationId string
		var id string
		var version int64
		var description string
		var systemDefault bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a new Work Type",
			Long:  `Create a new Work Type in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/work-type")
				req.PathParam("orgid", orgid)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("active", active, cmd.Flags().Changed("active"))
					req.BodyString("name", name)
					req.BodyString("workTypeCode", workTypeCode)
					req.BodyString("organizationId", organizationId)
					req.BodyString("id", id)
					req.BodyInt("version", version, cmd.Flags().Changed("version"))
					req.BodyString("description", description)
					req.BodyBool("systemDefault", systemDefault, cmd.Flags().Changed("system-default"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().BoolVar(&active, "active", false, "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&workTypeCode, "work-type-code", "", "")
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "")
		cmd.Flags().StringVar(&id, "id", "", "")
		cmd.Flags().Int64Var(&version, "version", 0, "")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().BoolVar(&systemDefault, "system-default", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		workTypesCmd.AddCommand(cmd)
	}

	{ // bulk-save
		var orgid string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "bulk-save",
			Short: "Bulk save Work Type(s)",
			Long:  `Create, Update or delete Work Type(s) in bulk in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/work-type/bulk")
				req.PathParam("orgid", orgid)
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
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		workTypesCmd.AddCommand(cmd)
	}

	{ // bulk-export
		var orgid string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "bulk-export",
			Short: "Bulk export Work Type(s)",
			Long:  `Export all Work Type(s) in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/work-type/bulk-export")
				req.PathParam("orgid", orgid)
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
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&page, "page", "", "Defines the number of displayed page. The page number starts from 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Defines the number of items to be displayed on a page. If the number specified is more than allowed max page size, the API will automatically adjust the page size to the max page size.")
		workTypesCmd.AddCommand(cmd)
	}

	{ // purge-inactive
		var orgid string
		var nextStartId string
		cmd := &cobra.Command{
			Use:   "purge-inactive",
			Short: "Purge inactive Work Type(s)",
			Long:  `Purge inactive Work Type(s) older than the configured interval for a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/work-type/purge-inactive-entities")
				req.PathParam("orgid", orgid)
				req.QueryParam("nextStartId", nextStartId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&nextStartId, "next-start-id", "", "This is the entity ID from which items for the next purge batch with be selected.")
		workTypesCmd.AddCommand(cmd)
	}

	{ // get-id
		var orgid string
		var id string
		cmd := &cobra.Command{
			Use:   "get-id",
			Short: "Get specific Work Type by ID",
			Long:  `Retrieve an existing Work Type by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/work-type/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
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
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&id, "id", "", "ID of the work_type.")
		cmd.MarkFlagRequired("id")
		workTypesCmd.AddCommand(cmd)
	}

	{ // update-id
		var orgid string
		var id string
		var active bool
		var name string
		var workTypeCode string
		var organizationId string
		var version int64
		var description string
		var systemDefault bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-id",
			Short: "Update specific Work Type by ID",
			Long:  `Update an existing Work Type by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PUT", "/organization/{orgid}/work-type/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("active", active, cmd.Flags().Changed("active"))
					req.BodyString("name", name)
					req.BodyString("workTypeCode", workTypeCode)
					req.BodyString("organizationId", organizationId)
					req.BodyString("id", id)
					req.BodyInt("version", version, cmd.Flags().Changed("version"))
					req.BodyString("description", description)
					req.BodyBool("systemDefault", systemDefault, cmd.Flags().Changed("system-default"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&id, "id", "", "ID of the work_type.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().BoolVar(&active, "active", false, "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&workTypeCode, "work-type-code", "", "")
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "")
		cmd.Flags().Int64Var(&version, "version", 0, "")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().BoolVar(&systemDefault, "system-default", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		workTypesCmd.AddCommand(cmd)
	}

	{ // delete-id
		var orgid string
		var id string
		cmd := &cobra.Command{
			Use:   "delete-id",
			Short: "Delete specific Work Type by ID",
			Long:  `Delete an existing Work Type by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/organization/{orgid}/work-type/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&id, "id", "", "ID of the work_type.")
		cmd.MarkFlagRequired("id")
		workTypesCmd.AddCommand(cmd)
	}

	{ // list-references
		var orgid string
		var id string
		var typeVal string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "list-references",
			Short: "List references for a specific Work Type",
			Long:  `Retrieve a list of all entities that have reference to an existing Work Type by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/work-type/{id}/incoming-references")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
				req.QueryParam("type", typeVal)
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
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&id, "id", "", "ID of this contact center resource.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&typeVal, "type", "", "Entity type of the other entity that has a reference to this specific entity.")
		cmd.Flags().StringVar(&page, "page", "", "Defines the number of displayed page. The page number starts from 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Defines the number of items to be displayed on a page. If the number specified is more than allowed max page size, the API will automatically adjust the page size to the max page size.")
		workTypesCmd.AddCommand(cmd)
	}

	{ // list
		var orgid string
		var filter string
		var attributes string
		var search string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Work Type(s)",
			Long:  `Retrieve a list of Work Type(s) in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v2/work-type")
				req.PathParam("orgid", orgid)
				req.QueryParam("filter", filter)
				req.QueryParam("attributes", attributes)
				req.QueryParam("search", search)
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
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&filter, "filter", "", "Specify a filter based on which the results will be fetched. All the fields are supported except: organizationId, createdTime, lastUpdatedTime   The examples below show some search queries - id==\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id!=\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id=in=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") - id=out=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see  <a href=\"https://www.here.com/docs/bundle/data-client-library-developer-guide-java-scala/page/client/rsql.html\">this reference</a>. For a list of supported operators, see <a href=\"https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference\">this syntax guide</a>.  Note: values to be used in the filter syntax should not contain space, and if so kindly bound it with quotes to apply filter. ")
		cmd.Flags().StringVar(&attributes, "attributes", "", "Specify the attributes to be returned.Default all attributes are returned along with specified columns. All Attributes are supported")
		cmd.Flags().StringVar(&search, "search", "", "Filter data based on the search keyword.Supported search columns(name)  The examples below show some search queries - \"Cisco\" - field==\"name\";value==\"Cisco\" - fields=in=(\"name\");value==\"Cisco\" ")
		cmd.Flags().StringVar(&page, "page", "", "Defines the number of displayed page. The page number starts from 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Defines the number of items to be displayed on a page. If the number specified is more than allowed max page size, the API will automatically adjust the page size to the max page size.")
		workTypesCmd.AddCommand(cmd)
	}

}
