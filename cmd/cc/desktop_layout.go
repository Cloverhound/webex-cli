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

var desktopLayoutCmd = &cobra.Command{
	Use:   "desktop-layout",
	Short: "DesktopLayout commands",
}

func init() {
	cmd.CcCmd.AddCommand(desktopLayoutCmd)

	{ // create
		var orgid string
		var defaultJsonModified bool
		var editedBy string
		var global bool
		var jsonFileContent string
		var jsonFileName string
		var name string
		var status bool
		var validated bool
		var organizationId string
		var id string
		var version int64
		var description string
		var validatedTime int64
		var defaultJsonModifiedTime int64
		var modifiedTime int64
		var teamIds []string
		var systemDefault bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a new Desktop Layout",
			Long:  `Create a new Desktop Layout in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/desktop-layout")
				req.PathParam("orgid", orgid)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("defaultJsonModified", defaultJsonModified, cmd.Flags().Changed("default-json-modified"))
					req.BodyString("editedBy", editedBy)
					req.BodyBool("global", global, cmd.Flags().Changed("global"))
					req.BodyString("jsonFileContent", jsonFileContent)
					req.BodyString("jsonFileName", jsonFileName)
					req.BodyString("name", name)
					req.BodyBool("status", status, cmd.Flags().Changed("status"))
					req.BodyBool("validated", validated, cmd.Flags().Changed("validated"))
					req.BodyString("organizationId", organizationId)
					req.BodyString("id", id)
					req.BodyInt("version", version, cmd.Flags().Changed("version"))
					req.BodyString("description", description)
					req.BodyInt("validatedTime", validatedTime, cmd.Flags().Changed("validated-time"))
					req.BodyInt("defaultJsonModifiedTime", defaultJsonModifiedTime, cmd.Flags().Changed("default-json-modified-time"))
					req.BodyInt("modifiedTime", modifiedTime, cmd.Flags().Changed("modified-time"))
					req.BodyStringSlice("teamIds", teamIds)
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
		cmd.Flags().BoolVar(&defaultJsonModified, "default-json-modified", false, "")
		cmd.Flags().StringVar(&editedBy, "edited-by", "", "")
		cmd.Flags().BoolVar(&global, "global", false, "")
		cmd.Flags().StringVar(&jsonFileContent, "json-file-content", "", "")
		cmd.Flags().StringVar(&jsonFileName, "json-file-name", "", "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().BoolVar(&status, "status", false, "")
		cmd.Flags().BoolVar(&validated, "validated", false, "")
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "")
		cmd.Flags().StringVar(&id, "id", "", "")
		cmd.Flags().Int64Var(&version, "version", 0, "")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().Int64Var(&validatedTime, "validated-time", 0, "")
		cmd.Flags().Int64Var(&defaultJsonModifiedTime, "default-json-modified-time", 0, "")
		cmd.Flags().Int64Var(&modifiedTime, "modified-time", 0, "")
		cmd.Flags().StringSliceVar(&teamIds, "team-ids", nil, "")
		cmd.Flags().BoolVar(&systemDefault, "system-default", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		desktopLayoutCmd.AddCommand(cmd)
	}

	{ // bulk-save
		var orgid string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "bulk-save",
			Short: "Bulk save Desktop Layout(s)",
			Long:  `Create, Update or delete Desktop Layout(s) in bulk in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/desktop-layout/bulk")
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
		desktopLayoutCmd.AddCommand(cmd)
	}

	{ // bulk-export
		var orgid string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "bulk-export",
			Short: "Bulk export Desktop Layout(s)",
			Long:  `Export all Desktop Layout(s) in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/desktop-layout/bulk-export")
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
		desktopLayoutCmd.AddCommand(cmd)
	}

	{ // purge-inactive
		var orgid string
		var nextStartId string
		cmd := &cobra.Command{
			Use:   "purge-inactive",
			Short: "Purge inactive Desktop Layout(s)",
			Long:  `Purge inactive Desktop Layout(s) older than the configured interval for a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/desktop-layout/purge-inactive-entities")
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
		desktopLayoutCmd.AddCommand(cmd)
	}

	{ // get-id
		var orgid string
		var id string
		cmd := &cobra.Command{
			Use:   "get-id",
			Short: "Get specific Desktop Layout by ID",
			Long:  `Retrieve an existing Desktop Layout by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/desktop-layout/{id}")
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Desktop Layout.")
		cmd.MarkFlagRequired("id")
		desktopLayoutCmd.AddCommand(cmd)
	}

	{ // update-id
		var orgid string
		var id string
		var defaultJsonModified bool
		var editedBy string
		var global bool
		var jsonFileContent string
		var jsonFileName string
		var name string
		var status bool
		var validated bool
		var organizationId string
		var version int64
		var description string
		var validatedTime int64
		var defaultJsonModifiedTime int64
		var modifiedTime int64
		var teamIds []string
		var systemDefault bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-id",
			Short: "Update specific Desktop Layout by ID",
			Long:  `Update an existing Desktop Layout by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PUT", "/organization/{orgid}/desktop-layout/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("defaultJsonModified", defaultJsonModified, cmd.Flags().Changed("default-json-modified"))
					req.BodyString("editedBy", editedBy)
					req.BodyBool("global", global, cmd.Flags().Changed("global"))
					req.BodyString("jsonFileContent", jsonFileContent)
					req.BodyString("jsonFileName", jsonFileName)
					req.BodyString("name", name)
					req.BodyBool("status", status, cmd.Flags().Changed("status"))
					req.BodyBool("validated", validated, cmd.Flags().Changed("validated"))
					req.BodyString("organizationId", organizationId)
					req.BodyString("id", id)
					req.BodyInt("version", version, cmd.Flags().Changed("version"))
					req.BodyString("description", description)
					req.BodyInt("validatedTime", validatedTime, cmd.Flags().Changed("validated-time"))
					req.BodyInt("defaultJsonModifiedTime", defaultJsonModifiedTime, cmd.Flags().Changed("default-json-modified-time"))
					req.BodyInt("modifiedTime", modifiedTime, cmd.Flags().Changed("modified-time"))
					req.BodyStringSlice("teamIds", teamIds)
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Desktop Layout.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().BoolVar(&defaultJsonModified, "default-json-modified", false, "")
		cmd.Flags().StringVar(&editedBy, "edited-by", "", "")
		cmd.Flags().BoolVar(&global, "global", false, "")
		cmd.Flags().StringVar(&jsonFileContent, "json-file-content", "", "")
		cmd.Flags().StringVar(&jsonFileName, "json-file-name", "", "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().BoolVar(&status, "status", false, "")
		cmd.Flags().BoolVar(&validated, "validated", false, "")
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "")
		cmd.Flags().Int64Var(&version, "version", 0, "")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().Int64Var(&validatedTime, "validated-time", 0, "")
		cmd.Flags().Int64Var(&defaultJsonModifiedTime, "default-json-modified-time", 0, "")
		cmd.Flags().Int64Var(&modifiedTime, "modified-time", 0, "")
		cmd.Flags().StringSliceVar(&teamIds, "team-ids", nil, "")
		cmd.Flags().BoolVar(&systemDefault, "system-default", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		desktopLayoutCmd.AddCommand(cmd)
	}

	{ // delete-id
		var orgid string
		var id string
		cmd := &cobra.Command{
			Use:   "delete-id",
			Short: "Delete specific Desktop Layout by ID",
			Long:  `Delete an existing Desktop Layout by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/organization/{orgid}/desktop-layout/{id}")
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Desktop Layout.")
		cmd.MarkFlagRequired("id")
		desktopLayoutCmd.AddCommand(cmd)
	}

	{ // list-references
		var orgid string
		var id string
		var typeVal string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "list-references",
			Short: "List references for a specific Desktop Layout",
			Long:  `Retrieve a list of all entities that have reference to an existing Desktop Layout by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/desktop-layout/{id}/incoming-references")
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
		desktopLayoutCmd.AddCommand(cmd)
	}

	{ // list
		var orgid string
		var filter string
		var attributes string
		var search string
		var page string
		var pageSize string
		var singleObjectResponse string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Desktop Layout(s)",
			Long: `Retrieve a list of Desktop Layout(s) in a given organization. Json file content field won't be avalible in get all even though it is showing in sample response structure. and it will be avialable only in get by id.
 Note: Array fields are removed from List API. If all fields are required please fetch Id's and use get-by-id API.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v2/desktop-layout")
				req.PathParam("orgid", orgid)
				req.QueryParam("filter", filter)
				req.QueryParam("attributes", attributes)
				req.QueryParam("search", search)
				req.QueryParam("page", page)
				req.QueryParam("pageSize", pageSize)
				req.QueryParam("singleObjectResponse", singleObjectResponse)
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
		cmd.Flags().StringVar(&filter, "filter", "", "Specify a filter based on which the results will be fetched. All the fields are supported except: organizationId, validatedTime, defaultJsonModifiedTime, modifiedTime, teamIds, createdTime, lastUpdatedTime   The examples below show some search queries - id==\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id!=\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id=in=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") - id=out=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see  <a href=\"https://www.here.com/docs/bundle/data-client-library-developer-guide-java-scala/page/client/rsql.html\">this reference</a>. For a list of supported operators, see <a href=\"https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference\">this syntax guide</a>.  Note: values to be used in the filter syntax should not contain space, and if so kindly bound it with quotes to apply filter. ")
		cmd.Flags().StringVar(&attributes, "attributes", "", "Specify the attributes to be returned.Default all attributes are returned along with specified columns. All Attributes are supported")
		cmd.Flags().StringVar(&search, "search", "", "Filter data based on the search keyword.Supported search columns(name)  The examples below show some search queries - \"Cisco\" - field==\"name\";value==\"Cisco\" - fields=in=(\"name\");value==\"Cisco\" ")
		cmd.Flags().StringVar(&page, "page", "", "Defines the number of displayed page. The page number starts from 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Defines the number of items to be displayed on a page. If the number specified is more than allowed max page size, the API will automatically adjust the page size to the max page size.")
		cmd.Flags().StringVar(&singleObjectResponse, "single-object-response", "", "Specifiy whether to include array fields in the response, This query param should use only if the response contain single record, if we are using for multiple objects response query param not supported and throws an exception.")
		desktopLayoutCmd.AddCommand(cmd)
	}

}
