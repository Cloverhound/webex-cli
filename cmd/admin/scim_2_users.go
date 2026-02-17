package admin

import (
	"fmt"
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
var _ = strings.Join

var scim2UsersCmd = &cobra.Command{
	Use:   "scim-2-users",
	Short: "Scim2Users commands",
}

func init() {
	cmd.AdminCmd.AddCommand(scim2UsersCmd)

	{ // create
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a user",
			Long:  "The SCIM 2 /Users API provides a programmatic way to manage users in Webex Identity using The Internet Engineering Task Force standard SCIM 2.0 standard as specified by [RFC 7643 SCIM 2.0 Core Schema ](https://datatracker.ietf.org/doc/html/rfc7643) and [RFC 7644 SCIM 2.0 Core Protocol](https://datatracker.ietf.org/doc/html/rfc7644).  The WebEx SCIM 2.0  APIs allow clients supporting the SCIM 2.0 standard to manage users, and groups within Webex.  Webex supports the following SCIM 2.0 Schemas:\n\n\u2022 urn:ietf:params:scim:schemas:core:2.0:User\n\n\u2022 urn:ietf:params:scim:schemas:extension:enterprise:2.0:User\n\n\u2022 urn:scim:schemas:extension:cisco:webexidentity:2.0:User\n\n<br/>\n\n**Authorization**\n\nOAuth token rendered by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_user_admin`\n\n<br/>\n\n**Usage**:\n\n1. Input JSON must contain schema: \"urn:ietf:params:scim:schemas:core:2.0:User\".\n\n2. Support 3 schemas :\n    - \"urn:ietf:params:scim:schemas:core:2.0:User\"\n    - \"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User\"\n    - \"urn:scim:schemas:extension:cisco:webexidentity:2.0:User\"\n\n3. Unrecognized schemas (ID/section) are ignored.\n\n4. Read-only attributes provided as input values are ignored.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/identity/scim/{orgId}/v2/Users")
				req.PathParam("orgId", orgId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		scim2UsersCmd.AddCommand(cmd)
	}

	{ // search
		var orgId string
		var filter string
		var attributes string
		var excludedAttributes string
		var sortBy string
		var sortOrder string
		var startIndex string
		var count string
		var returnGroups string
		var includeGroupDetails string
		var groupUsageTypes string
		cmd := &cobra.Command{
			Use:   "search",
			Short: "Search users",
			Long:  "<br/>\n\n**Authorization**\n\nOAuth token rendered by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n- `identity:people_read`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_user_admin`\n\n- `id_readonly_admin`\n\n- `id_device_admin`\n\n<br/>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/identity/scim/{orgId}/v2/Users")
				req.PathParam("orgId", orgId)
				req.QueryParam("filter", filter)
				req.QueryParam("attributes", attributes)
				req.QueryParam("excludedAttributes", excludedAttributes)
				req.QueryParam("sortBy", sortBy)
				req.QueryParam("sortOrder", sortOrder)
				req.QueryParam("startIndex", startIndex)
				req.QueryParam("count", count)
				req.QueryParam("returnGroups", returnGroups)
				req.QueryParam("includeGroupDetails", includeGroupDetails)
				req.QueryParam("groupUsageTypes", groupUsageTypes)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&filter, "filter", "", "The URL encoded filter. If the value is empty, the API will return all users under the organization.  The examples below show some search filters:  - `userName` eq \"user1@example.com\"  - `userName` sw \"user1@example\"  - `userName` ew \"example\"  - `phoneNumbers` [ type eq \"mobile\" and value eq \"14170120\"]  - `urn:scim:schemas:extension:cisco:webexidentity:2.0:User:meta.organizationId` eq \"0ae87ade-8c8a-4952-af08-318798958d0c\"  - For more filter patterns, please check [filtering](https://datatracker.ietf.org/doc/html/rfc7644#section-3.4.2.2).  | **Attributes** | **Operators** | |-----|-----| | **SCIM Core** &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; | ---- | | `id` | eq | | `userName` | eq sw ew | | `name.familyName` | eq sw ew | | `name.givenName` | eq sw | | `name.middleName` | eq sw | | `name.formatted` | eq sw | | `displayName` | eq sw ew | | `nickName` | eq sw ew | | `emails.display` | eq sw ew | | `emails.value` | eq sw ew | | `phoneNumbers.value` | eq sw ew | | `phoneNumbers.display` | eq sw ew | | **Enterprise Extensions** | ---- | | `employeeNumber` | eq sw ew | | `costCenter` | eq sw ew | | `organization` | eq sw ew | | `division` | eq sw ew | | `department` | eq sw ew | | `manager.value` | eq | | `manager.displayName` | eq sw ew |")
		cmd.Flags().StringVar(&attributes, "attributes", "", "A multi-valued list of string names for resource attributes to return in the response, like 'userName,department,emails'. It supports the SCIM id 'urn:ietf:params:scim:schemas:extension:enterprise:2.0:User,userName'. The default is empty, all attributes will be returned")
		cmd.Flags().StringVar(&excludedAttributes, "excluded-attributes", "", "A multi-valued list of strings names for resource attributes to be removed from the default set of attributes to return. The default is empty, all attributes will be returned")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "A string for the attribute whose value can be used to order the returned responses. Now we only allow `userName`, `id`, `meta.lastModified` to sort.")
		cmd.Flags().StringVar(&sortOrder, "sort-order", "", "A string for the order in which the 'sortBy' parameter is applied. Allowed values are 'ascending' and 'descending'.")
		cmd.Flags().StringVar(&startIndex, "start-index", "", "An integer for the 1-based index of the first query result. The default is 1.")
		cmd.Flags().StringVar(&count, "count", "", "An integer for the maximum number of query results per page.  The default is 100.")
		cmd.Flags().StringVar(&returnGroups, "return-groups", "", "Define whether the group information needs to be returned.  The default is false.")
		cmd.Flags().StringVar(&includeGroupDetails, "include-group-details", "", "Define whether the group information with details needs to be returned. The default is false.")
		cmd.Flags().StringVar(&groupUsageTypes, "group-usage-types", "", "Returns groups with details of the specified group type.")
		scim2UsersCmd.AddCommand(cmd)
	}

	{ // get
		var orgId string
		var userId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get a user",
			Long:  "<br/>\n\n**Authorization**\n\nOAuth token rendered by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n- `identity:people_read`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_user_admin`\n\n- `id_readonly_admin`\n\n- `id_device_admin`\n\n<br/>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/identity/scim/{orgId}/v2/Users/{userId}")
				req.PathParam("orgId", orgId)
				req.PathParam("userId", userId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&userId, "user-id", "", "Webex Identity assigned user identifier.")
		cmd.MarkFlagRequired("user-id")
		scim2UsersCmd.AddCommand(cmd)
	}

	{ // update-put
		var orgId string
		var userId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-put",
			Short: "Update a user with PUT",
			Long:  "<br/>\n\n**Authorization**\n\nOAuth token rendered by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_user_admin`\n\n<br/>\n\n**Usage**:\n\n1. Input JSON must contain schema: \"urn:ietf:params:scim:schemas:core:2.0:User\".\n\n2. Support 3 schemas :\n    - \"urn:ietf:params:scim:schemas:core:2.0:User\"\n    - \"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User\"\n    - \"urn:scim:schemas:extension:cisco:webexidentity:2.0:User\"\n\n3. Unrecognized schemas (ID/section) are ignored.\n\n4. Read-only attributes provided as input values are ignored.\n\n5. User `id` will not be changed.\n\n6. `meta`.`created` will not be changed.\n\n7. The PUT API replaces the contents of the user's data with the data in the request body.  All attributes specified in the request body will replace all existing attributes for the `userId` specified in the URL.  Should you wish to replace or change some attributes as opposed to all attributes please refer to the SCIM PATCH operation https://developer.webex.com/docs/api/v1/scim2-user/update-a-user-with-patch.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/identity/scim/{orgId}/v2/Users/{userId}")
				req.PathParam("orgId", orgId)
				req.PathParam("userId", userId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&userId, "user-id", "", "Webex Identity assigned user identifier.")
		cmd.MarkFlagRequired("user-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		scim2UsersCmd.AddCommand(cmd)
	}

	{ // update-patch
		var orgId string
		var userId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-patch",
			Short: "Update a user with PATCH",
			Long:  "<br/>\n\n**Authorization**\n\nOAuth token rendered by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_user_admin`\n\n<br/>\n\n**Usage**:\n\n1. The PATCH API replaces individual attributes and roles of the user's data in the request body.\n   The PATCH API supports `add`, `remove`, and `replace` operations on any individual\n   attribute or role allowing only specific attributes of the user's object to be modified.\n\n2. Each operation against an attribute must be compatible with the attribute's mutability.\n\n3. Each PATCH operation represents a single action to be applied to the\n   same SCIM resource specified by the request URI.  Operations are\n   applied sequentially in the order they appear in the array.  Each\n   operation in the sequence is applied to the target resource; the\n   resulting resource becomes the target of the next operation.\n   Evaluation continues until all operations are successfully applied or\n   until an error condition is encountered.\n\n<br/>\n\n**Add operations**:\n\nThe `add` operation adds a new attribute value to an existing resource.\nThe operation must contain a `value` member whose content specifies the value to be added.\nThe value may be a quoted value, or it may be a JSON object containing the sub-attributes of the complex attribute specified in the operation's `path`.\nThe result of the add operation depends upon the target `path` reference locations:\n\n<br/>\n\n- If omitted, the target location is assumed to be the resource itself.  The `value` parameter contains a set of attributes to be added to the resource.\n\n- If the target location does not exist, the attribute and value are added.\n\n- If the target location specifies a complex attribute, a set of sub-attributes shall be specified in the `value` parameter.\n\n- If the target location specifies a multi-valued attribute, a new value is added to the attribute.\n\n- If the target location specifies a single-valued attribute, the existing value is replaced.\n\n- If the target location specifies an attribute that does not exist (has no value), the attribute is added with the new value.\n\n- If the target location exists, the value is replaced.\n\n- If the target location already contains the value specified, no changes should be made to the resource.\n\n<br/>\n\n**Replace operations**:\n\nThe `replace` operation replaces the value at the target location specified by the `path`.\nThe operation performs the following functions, depending on the target location specified by `path`:\n\n<br/>\n\n- If the `path` parameter is omitted, the target is assumed to be the resource itself.  In this case, the `value` attribute shall contain a list of one or more attributes to be replaced.\n\n- If the target location is a single-value attribute, the value of the attribute is replaced.\n\n- If the target location is a multi-valued attribute and no filter is specified, the attribute and all values are replaced.\n\n- If the target location path specifies an attribute that does not exist, the service provider shall treat the operation as an \"add\".\n\n- If the target location specifies a complex attribute, a set of sub-attributes SHALL be specified in the `value` parameter, which replaces any existing values or adds where an attribute did not previously exist.  Sub-attributes not specified in the `value` parameters are left unchanged.\n\n- If the target location is a multi-valued attribute and a value selection (\"valuePath\") filter is specified that matches one or more values of the multi-valued attribute, then all matching record values will be replaced.\n\n- If the target location is a complex multi-valued attribute with a value selection filter (\"valuePath\") and a specific sub-attribute (e.g., \"addresses[type eq \"work\"].streetAddress\"), the matching sub-attribute of all matching records is replaced.\n\n- If the target location is a multi-valued attribute for which a value selection filter (\"valuePath\") has been supplied and no record match was made, the service provider will return failure as HTTP status code 400 and a `scimType` error code of \"noTarget\".\n\n<br/>\n\n**Remove operations**:\n\nThe `remove` operation removes the value at the target location specified by the required attribute `path`.  The operation performs the following functions, depending on the target location specified by `path`:\n\n<br/>\n\n- If `path` is unspecified, the operation fails with HTTP status code 400 and a \"scimType\" error code of \"noTarget\".\n\n- If the target location is a single-value attribute, the attribute and its associated value is removed, and the attribute will be considered unassigned.\n\n- If the target location is a multi-valued attribute and no filter is specified, the attribute and all values are removed, and the attribute SHALL be considered unassigned.\n\n- If the target location is a multi-valued attribute and a complex filter is specified comparing a `value`, the values matched by the filter are removed.  If no other values remain after the removal of the selected values, the multi-valued attribute will be considered unassigned.\n\n- If the target location is a complex multi-valued attribute and a complex filter is specified based on the attribute's sub-attributes, the matching records are removed.  Sub-attributes whose values have been removed will be considered unassigned.  If the complex multi-valued attribute has no remaining records, the attribute will be considered unassigned.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PATCH", "/identity/scim/{orgId}/v2/Users/{userId}")
				req.PathParam("orgId", orgId)
				req.PathParam("userId", userId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&userId, "user-id", "", "Webex Identity assigned user identifier.")
		cmd.MarkFlagRequired("user-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		scim2UsersCmd.AddCommand(cmd)
	}

	{ // delete
		var orgId string
		var userId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a user",
			Long:  "<br/>\n\n**Authorization**\n\nOAuth token rendered by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_user_admin`\n\n<br/>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/identity/scim/{orgId}/v2/Users/{userId}")
				req.PathParam("orgId", orgId)
				req.PathParam("userId", userId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&userId, "user-id", "", "Webex Identity assigned user identifier.")
		cmd.MarkFlagRequired("user-id")
		scim2UsersCmd.AddCommand(cmd)
	}

	{ // get-me
		cmd := &cobra.Command{
			Use:   "get-me",
			Short: "Get Me",
			Long:  "<br/>\n\n**Authorization**\n\nOAuth token rendered by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n- `identity:people_read`\n\n<br/>\n\nThe API can be used by any user to retrieve user information using their own access token.\n\n<br/>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/identity/scim/v2/Users/me")
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
		scim2UsersCmd.AddCommand(cmd)
	}

}
