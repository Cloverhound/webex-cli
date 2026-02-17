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

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Skill commands",
}

func init() {
	cmd.CcCmd.AddCommand(skillCmd)

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
			Short: "List Skill(s)",
			Long: `Retrieve a list of Skill(s) in a given organization.
 Note: Array fields are removed from List API. If all fields are required please fetch Id's and use get-by-id API.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v2/skill")
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
		cmd.Flags().StringVar(&filter, "filter", "", "Specify a filter based on which the results will be fetched. All the fields are supported except: organizationId, enumSkillValues, createdTime, lastUpdatedTime   The examples below show some search queries - id==\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id!=\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id=in=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") - id=out=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see  <a href=\"https://www.here.com/docs/bundle/data-client-library-developer-guide-java-scala/page/client/rsql.html\">this reference</a>. For a list of supported operators, see <a href=\"https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference\">this syntax guide</a>.  Note: values to be used in the filter syntax should not contain space, and if so kindly bound it with quotes to apply filter. ")
		cmd.Flags().StringVar(&attributes, "attributes", "", "Specify the attributes to be returned.Default all attributes are returned along with specified columns. All Attributes are supported except (enumSkillValues)")
		cmd.Flags().StringVar(&search, "search", "", "Filter data based on the search keyword.Supported search columns(name)  The examples below show some search queries - \"Cisco\" - field==\"name\";value==\"Cisco\" - fields=in=(\"name\");value==\"Cisco\" ")
		cmd.Flags().StringVar(&page, "page", "", "Defines the number of displayed page. The page number starts from 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Defines the number of items to be displayed on a page. If the number specified is more than allowed max page size, the API will automatically adjust the page size to the max page size.")
		cmd.Flags().StringVar(&singleObjectResponse, "single-object-response", "", "Specifiy whether to include array fields in the response, This query param should use only if the response contain single record, if we are using for multiple objects response query param not supported and throws an exception.")
		skillCmd.AddCommand(cmd)
	}

	{ // create
		var orgid string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a new Skill",
			Long:  `Create a new Skill in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/skill")
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
		skillCmd.AddCommand(cmd)
	}

	{ // bulk-save
		var orgid string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "bulk-save",
			Short: "Bulk save Skill(s)",
			Long:  `Create, Update or delete Skill(s) in bulk in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/skill/bulk")
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
		skillCmd.AddCommand(cmd)
	}

	{ // bulk-export
		var orgid string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "bulk-export",
			Short: "Bulk export Skill(s)",
			Long:  `Export all Skill(s) in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/skill/bulk-export")
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
		skillCmd.AddCommand(cmd)
	}

	{ // purge-inactive
		var orgid string
		var nextStartId string
		cmd := &cobra.Command{
			Use:   "purge-inactive",
			Short: "Purge inactive Skill(s)",
			Long:  `Purge inactive Skill(s) older than the configured interval for a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/skill/purge-inactive-entities")
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
		skillCmd.AddCommand(cmd)
	}

	{ // get-id
		var orgid string
		var id string
		cmd := &cobra.Command{
			Use:   "get-id",
			Short: "Get specific Skill by ID",
			Long:  `Retrieve an existing Skill by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/skill/{id}")
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Skill.")
		cmd.MarkFlagRequired("id")
		skillCmd.AddCommand(cmd)
	}

	{ // update-id
		var orgid string
		var id string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-id",
			Short: "Update specific Skill by ID",
			Long:  `Update an existing Skill by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PUT", "/organization/{orgid}/skill/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Skill.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		skillCmd.AddCommand(cmd)
	}

	{ // delete-id
		var orgid string
		var id string
		cmd := &cobra.Command{
			Use:   "delete-id",
			Short: "Delete specific Skill by ID",
			Long:  `Delete an existing Skill by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/organization/{orgid}/skill/{id}")
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Skill.")
		cmd.MarkFlagRequired("id")
		skillCmd.AddCommand(cmd)
	}

	{ // list-references
		var orgid string
		var id string
		var typeVal string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "list-references",
			Short: "List references for a specific Skill",
			Long:  `Retrieve a list of all entities that have reference to an existing Skill by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/skill/{id}/incoming-references")
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
		skillCmd.AddCommand(cmd)
	}

}
