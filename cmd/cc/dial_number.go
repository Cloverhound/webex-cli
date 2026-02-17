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

var dialNumberCmd = &cobra.Command{
	Use:   "dial-number",
	Short: "DialNumber commands",
}

func init() {
	cmd.CcCmd.AddCommand(dialNumberCmd)

	{ // list-dialed-mapping
		var orgid string
		var filter string
		var attributes string
		var search string
		var page string
		var pageSize string
		var includeEntryPointName string
		cmd := &cobra.Command{
			Use:   "list-dialed-mapping",
			Short: "List Dialed Number Mapping(s)",
			Long:  `Retrieve a list of Dialed Number Mapping(s) in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v3/dial-number")
				req.PathParam("orgid", orgid)
				req.QueryParam("filter", filter)
				req.QueryParam("attributes", attributes)
				req.QueryParam("search", search)
				req.QueryParam("page", page)
				req.QueryParam("pageSize", pageSize)
				req.QueryParam("includeEntryPointName", includeEntryPointName)
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
		cmd.Flags().StringVar(&filter, "filter", "", "Specify a filter based on which the results will be fetched. All the fields are supported except: organizationId, createdTime, lastUpdatedTime   The examples below show some search queries - id==\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id!=\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id=in=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") - id=out=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see  <a href=\"https://www.here.com/docs/bundle/data-client-library-developer-guide-java-scala/page/client/rsql.html\">this reference</a>. For a list of supported operators, see <a href=\"https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference\">this syntax guide</a>.  Note: values to be used in the filter syntax should not contain spaces. If they do, please enclose them in quotes to apply the filter. ")
		cmd.Flags().StringVar(&attributes, "attributes", "", "Specify the attributes to be returned. By default, all attributes are returned along with the specified columns. All attributes are supported. except (links)")
		cmd.Flags().StringVar(&search, "search", "", "Filter data based on the search keyword.Supported search columns(dialledNumber)  The examples below show some search queries - \"Cisco\" - field==\"dialledNumber\";value==\"Cisco\" - fields=in=(\"dialledNumber\");value==\"Cisco\" ")
		cmd.Flags().StringVar(&page, "page", "", "Defines the number of displayed page. The page number starts from 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Defines the number of items to be displayed on a page. If the number specified is more than allowed max page size, the API will automatically adjust the page size to the max page size.")
		cmd.Flags().StringVar(&includeEntryPointName, "include-entry-point-name", "", "If includeEntryPointName is set to true and entryPointName is in the attributes, the API will return entryPointName in the Get All response, and filtering, searching, and sorting on entryPointName will also be enabled.")
		dialNumberCmd.AddCommand(cmd)
	}

	{ // create-dialed-mapping
		var orgid string
		var entryPointId string
		var entryPointName string
		var organizationId string
		var id string
		var version int64
		var dialledNumber string
		var extension string
		var routingPrefix string
		var esn string
		var routePointId string
		var defaultAni bool
		var location string
		var regionId string
		var createdTime int64
		var lastUpdatedTime int64
		var dialledNumberDigits string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-dialed-mapping",
			Short: "Create a new Dialed Number Mapping",
			Long:  `Create a new Dialed Number Mapping in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/dial-number")
				req.PathParam("orgid", orgid)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("entryPointId", entryPointId)
					req.BodyString("entryPointName", entryPointName)
					req.BodyString("organizationId", organizationId)
					req.BodyString("id", id)
					req.BodyInt("version", version, cmd.Flags().Changed("version"))
					req.BodyString("dialledNumber", dialledNumber)
					req.BodyString("extension", extension)
					req.BodyString("routingPrefix", routingPrefix)
					req.BodyString("esn", esn)
					req.BodyString("routePointId", routePointId)
					req.BodyBool("defaultAni", defaultAni, cmd.Flags().Changed("default-ani"))
					req.BodyString("location", location)
					req.BodyString("regionId", regionId)
					req.BodyInt("createdTime", createdTime, cmd.Flags().Changed("created-time"))
					req.BodyInt("lastUpdatedTime", lastUpdatedTime, cmd.Flags().Changed("last-updated-time"))
					req.BodyString("dialledNumberDigits", dialledNumberDigits)
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
		cmd.Flags().StringVar(&entryPointId, "entry-point-id", "", "")
		cmd.Flags().StringVar(&entryPointName, "entry-point-name", "", "")
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "")
		cmd.Flags().StringVar(&id, "id", "", "")
		cmd.Flags().Int64Var(&version, "version", 0, "")
		cmd.Flags().StringVar(&dialledNumber, "dialled-number", "", "")
		cmd.Flags().StringVar(&extension, "extension", "", "")
		cmd.Flags().StringVar(&routingPrefix, "routing-prefix", "", "")
		cmd.Flags().StringVar(&esn, "esn", "", "")
		cmd.Flags().StringVar(&routePointId, "route-point-id", "", "")
		cmd.Flags().BoolVar(&defaultAni, "default-ani", false, "")
		cmd.Flags().StringVar(&location, "location", "", "")
		cmd.Flags().StringVar(&regionId, "region-id", "", "")
		cmd.Flags().Int64Var(&createdTime, "created-time", 0, "")
		cmd.Flags().Int64Var(&lastUpdatedTime, "last-updated-time", 0, "")
		cmd.Flags().StringVar(&dialledNumberDigits, "dialled-number-digits", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dialNumberCmd.AddCommand(cmd)
	}

	{ // delete-all-dialed-mapping
		var orgid string
		cmd := &cobra.Command{
			Use:   "delete-all-dialed-mapping",
			Short: "Delete all Dialed Number Mapping(s)",
			Long:  `Delete all Dialed Number Mapping(s) in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/organization/{orgid}/dial-number")
				req.PathParam("orgid", orgid)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		dialNumberCmd.AddCommand(cmd)
	}

	{ // bulk-save-dialed-mapping
		var orgid string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "bulk-save-dialed-mapping",
			Short: "Bulk save Dialed Number Mapping(s)",
			Long:  `Create, Update or delete Dialed Number Mapping(s) in bulk in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/dial-number/bulk")
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
		dialNumberCmd.AddCommand(cmd)
	}

	{ // bulk-export-dialed-mapping
		var orgid string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "bulk-export-dialed-mapping",
			Short: "Bulk export Dialed Number Mapping(s)",
			Long:  `Export all Dialed Number Mapping(s) in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/dial-number/bulk-export")
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
		dialNumberCmd.AddCommand(cmd)
	}

	{ // list-dialed-dialed-mapping
		var orgid string
		cmd := &cobra.Command{
			Use:   "list-dialed-dialed-mapping",
			Short: "List  only dialed numbers(property - dialledNumber) from Dialed Number Mapping(s)",
			Long:  `Retrieve a list of  only dialed numbers(property - dialledNumber) from Dialed Number Mapping(s) without pagination in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/dial-number/numbers-only")
				req.PathParam("orgid", orgid)
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
		dialNumberCmd.AddCommand(cmd)
	}

	{ // get-dialed-mapping-id
		var orgid string
		var id string
		cmd := &cobra.Command{
			Use:   "get-dialed-mapping-id",
			Short: "Get specific Dialed Number Mapping by ID",
			Long:  `Retrieve an existing Dialed Number Mapping by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/dial-number/{id}")
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Dialed Number Mapping.")
		cmd.MarkFlagRequired("id")
		dialNumberCmd.AddCommand(cmd)
	}

	{ // update-dialed-mapping-id
		var orgid string
		var id string
		var entryPointId string
		var entryPointName string
		var organizationId string
		var version int64
		var dialledNumber string
		var extension string
		var routingPrefix string
		var esn string
		var routePointId string
		var defaultAni bool
		var location string
		var regionId string
		var createdTime int64
		var lastUpdatedTime int64
		var dialledNumberDigits string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-dialed-mapping-id",
			Short: "Update specific Dialed Number Mapping by ID",
			Long:  `Update an existing Dialed Number Mapping by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PUT", "/organization/{orgid}/dial-number/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("entryPointId", entryPointId)
					req.BodyString("entryPointName", entryPointName)
					req.BodyString("organizationId", organizationId)
					req.BodyString("id", id)
					req.BodyInt("version", version, cmd.Flags().Changed("version"))
					req.BodyString("dialledNumber", dialledNumber)
					req.BodyString("extension", extension)
					req.BodyString("routingPrefix", routingPrefix)
					req.BodyString("esn", esn)
					req.BodyString("routePointId", routePointId)
					req.BodyBool("defaultAni", defaultAni, cmd.Flags().Changed("default-ani"))
					req.BodyString("location", location)
					req.BodyString("regionId", regionId)
					req.BodyInt("createdTime", createdTime, cmd.Flags().Changed("created-time"))
					req.BodyInt("lastUpdatedTime", lastUpdatedTime, cmd.Flags().Changed("last-updated-time"))
					req.BodyString("dialledNumberDigits", dialledNumberDigits)
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Dialed Number Mapping.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&entryPointId, "entry-point-id", "", "")
		cmd.Flags().StringVar(&entryPointName, "entry-point-name", "", "")
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "")
		cmd.Flags().Int64Var(&version, "version", 0, "")
		cmd.Flags().StringVar(&dialledNumber, "dialled-number", "", "")
		cmd.Flags().StringVar(&extension, "extension", "", "")
		cmd.Flags().StringVar(&routingPrefix, "routing-prefix", "", "")
		cmd.Flags().StringVar(&esn, "esn", "", "")
		cmd.Flags().StringVar(&routePointId, "route-point-id", "", "")
		cmd.Flags().BoolVar(&defaultAni, "default-ani", false, "")
		cmd.Flags().StringVar(&location, "location", "", "")
		cmd.Flags().StringVar(&regionId, "region-id", "", "")
		cmd.Flags().Int64Var(&createdTime, "created-time", 0, "")
		cmd.Flags().Int64Var(&lastUpdatedTime, "last-updated-time", 0, "")
		cmd.Flags().StringVar(&dialledNumberDigits, "dialled-number-digits", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dialNumberCmd.AddCommand(cmd)
	}

	{ // delete-dialed-mapping-id
		var orgid string
		var id string
		cmd := &cobra.Command{
			Use:   "delete-dialed-mapping-id",
			Short: "Delete specific Dialed Number Mapping by ID",
			Long:  `Delete an existing Dialed Number Mapping by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/organization/{orgid}/dial-number/{id}")
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Dialed Number Mapping.")
		cmd.MarkFlagRequired("id")
		dialNumberCmd.AddCommand(cmd)
	}

	{ // list-references-dialed-mapping
		var orgid string
		var id string
		var typeVal string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "list-references-dialed-mapping",
			Short: "List references for a specific Dialed Number Mapping",
			Long:  `Retrieve a list of all entities that have reference to an existing Dialed Number Mapping by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/dial-number/{id}/incoming-references")
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
		dialNumberCmd.AddCommand(cmd)
	}

}
