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

var userProfilesCmd = &cobra.Command{
	Use:   "user-profiles",
	Short: "UserProfiles commands",
}

func init() {
	cmd.CcCmd.AddCommand(userProfilesCmd)

	{ // list
		var orgid string
		var filter string
		var attributes string
		var search string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List user profiles",
			Long: `Retrieve a list of user profiles in a given organization.
 Note: Array fields are removed from List API. If all fields are required please fetch Id's and use get-by-id API.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v3/user-profile")
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
		cmd.Flags().StringVar(&filter, "filter", "", "Specify a filter based on which the results will be fetched. All the fields are supported except: organizationId, userProfileAppModules, entryPoints, sites, queues, teams, editableFolderIds, viewableFolderIds, nonViewableFolderIds, createdTime, lastUpdatedTime   The examples below show some search queries - id==\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id!=\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id=in=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") - id=out=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see  <a href=\"https://www.here.com/docs/bundle/data-client-library-developer-guide-java-scala/page/client/rsql.html\">this reference</a>. For a list of supported operators, see <a href=\"https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference\">this syntax guide</a>.  Note: values to be used in the filter syntax should not contain spaces. If they do, please enclose them in quotes to apply the filter. ")
		cmd.Flags().StringVar(&attributes, "attributes", "", "Specify the attributes to be returned. By default, all attributes are returned along with the specified columns. All attributes are supported. except (entryPoints,sites, queues, teams, userProfileAppModules,editableFolderIds, viewableFolderIds, nonViewableFolderIds)")
		cmd.Flags().StringVar(&search, "search", "", "Filter data based on the search keyword.Supported search columns(name, profileType, description)  The examples below show some search queries - \"Cisco\" - field==\"name\";value==\"Cisco\" - fields=in=(\"name\",\"description\");value==\"Cisco\" ")
		cmd.Flags().StringVar(&page, "page", "", "Defines the number of displayed page. The page number starts from 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Defines the number of items to be displayed on a page. If the number specified is more than allowed max page size, the API will automatically adjust the page size to the max page size.")
		userProfilesCmd.AddCommand(cmd)
	}

	{ // bulk-save
		var orgid string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "bulk-save",
			Short: "Bulk save User Profiles",
			Long:  `Create, Update or delete user profiles in bulk in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/v3/user-profile/bulk")
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
		userProfilesCmd.AddCommand(cmd)
	}

	{ // bulk-export
		var orgid string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "bulk-export",
			Short: "Bulk export User Profiles",
			Long:  `Export all user profiles in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v3/user-profile/bulk-export")
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
		userProfilesCmd.AddCommand(cmd)
	}

	{ // purge-inactive
		var orgid string
		var nextStartId string
		cmd := &cobra.Command{
			Use:   "purge-inactive",
			Short: "Purge inactive User Profile(s)",
			Long:  `Purge inactive User Profile(s) older than the configured interval for a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/user-profile/purge-inactive-entities")
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
		userProfilesCmd.AddCommand(cmd)
	}

	{ // get-id
		var orgid string
		var id string
		var includeNames string
		cmd := &cobra.Command{
			Use:   "get-id",
			Short: "Get specific User Profile by ID",
			Long:  `Retrieve an existing user profile by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v3/user-profile/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
				req.QueryParam("includeNames", includeNames)
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the User Profile.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&includeNames, "include-names", "", "Flag to include resource names in the response.")
		userProfilesCmd.AddCommand(cmd)
	}

	{ // update-id
		var orgid string
		var id string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-id",
			Short: "Update specific User Profile by ID",
			Long:  `Update an existing user profile by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PUT", "/organization/{orgid}/v3/user-profile/{id}")
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the User Profile.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userProfilesCmd.AddCommand(cmd)
	}

	{ // delete-id
		var orgid string
		var id string
		cmd := &cobra.Command{
			Use:   "delete-id",
			Short: "Delete specific User Profile by ID",
			Long:  `Delete an existing user profile by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/organization/{orgid}/v3/user-profile/{id}")
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the User Profile.")
		cmd.MarkFlagRequired("id")
		userProfilesCmd.AddCommand(cmd)
	}

	{ // list-references
		var orgid string
		var id string
		var typeVal string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "list-references",
			Short: "List references for a specific User Profile",
			Long:  `Retrieve a list of all entities that have reference to an existing user profile by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/user-profile/{id}/incoming-references")
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
		userProfilesCmd.AddCommand(cmd)
	}

	{ // create
		var orgid string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a new User Profile",
			Long:  `Create a new user profile in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/v3/user-profile")
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
		userProfilesCmd.AddCommand(cmd)
	}

	{ // get-acl-id
		var orgid string
		var id string
		var names string
		cmd := &cobra.Command{
			Use:   "get-acl-id",
			Short: "Get specific User Profile ACL by ID",
			Long:  `Retrieve an existing User Profile ACL by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v3/user-profile/{id}/acl")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
				req.QueryParam("names", names)
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
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the User Profile.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&names, "names", "", "    Default all resources are returned in the ACL.     If you want to filter the ACL by specific resources,     provide a comma-separated list of resource names to filter the ACL. Ex: /url?names=site,team ")
		userProfilesCmd.AddCommand(cmd)
	}

}
