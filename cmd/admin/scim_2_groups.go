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

var scim2GroupsCmd = &cobra.Command{
	Use:   "scim-2-groups",
	Short: "Scim2Groups commands",
}

func init() {
	cmd.AdminCmd.AddCommand(scim2GroupsCmd)

	{ // create
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a group",
			Long:  "Create a new group for a given organization. The group may optionally be created with group members.\n\n<br/>\n\n**Authorization**\n\nOAuth token returned by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_group_admin`\n\n<br/>\n\n**Usage**:\n\n1. The input JSON must conform to one of the following schemas:\n    - `urn:ietf:params:scim:schemas:core:2.0:Group`\n    - `urn:scim:schemas:extension:cisco:webexidentity:2.0:Group`\n\n1. Unrecognized schemas (ID/section) are ignored.\n\n1. Read-only attributes provided as input values are ignored.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/identity/scim/{orgId}/v2/Groups")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "The ID of the organization to which this group belongs. If not specified, the organization ID from the OAuth token is used.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		scim2GroupsCmd.AddCommand(cmd)
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
		var includeMembers string
		var memberType string
		cmd := &cobra.Command{
			Use:   "search",
			Short: "Search groups",
			Long:  "Retrieve a list of groups in the organization.\n\nLong result sets are split into [pages](/docs/basics#pagination).\n\n<br/>\n\n**Authorization**\n\nAn OAuth token rendered by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n- `identity:people_read`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_group_admin`\n\n- `id_readonly_admin`\n\n- `id_device_admin`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/identity/scim/{orgId}/v2/Groups")
				req.PathParam("orgId", orgId)
				req.QueryParam("filter", filter)
				req.QueryParam("attributes", attributes)
				req.QueryParam("excludedAttributes", excludedAttributes)
				req.QueryParam("sortBy", sortBy)
				req.QueryParam("sortOrder", sortOrder)
				req.QueryParam("startIndex", startIndex)
				req.QueryParam("count", count)
				req.QueryParam("includeMembers", includeMembers)
				req.QueryParam("memberType", memberType)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "The ID of the organization to which this group belongs. If not specified, the organization ID from the OAuth token is used.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&filter, "filter", "", "The url encoded filter. The example content is 'displayName Eq \"group1@example.com\" or displayName Eq \"group2@example.com\"'. For more filter patterns, see https://datatracker.ietf.org/doc/html/rfc7644#section-3.4.2.2. If the value is empty, the API returns all groups under the organization. ")
		cmd.Flags().StringVar(&attributes, "attributes", "", "The attributes to return.")
		cmd.Flags().StringVar(&excludedAttributes, "excluded-attributes", "", "Attributes to be excluded from the return.")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "A string indicating the attribute whose value be used to order the returned responses. Now we only allow `displayName, id, meta.lastModified` to sort.")
		cmd.Flags().StringVar(&sortOrder, "sort-order", "", "A string indicating the order in which the `sortBy` parameter is applied. Allowed values are `ascending` and `descending`.")
		cmd.Flags().StringVar(&startIndex, "start-index", "", "An integer indicating the 1-based index of the first query result. The default is 1.")
		cmd.Flags().StringVar(&count, "count", "", "An integer indicating the desired maximum number of query results per page. The default is 100.")
		cmd.Flags().StringVar(&includeMembers, "include-members", "", "Default \"false\". If false, no members returned.")
		cmd.Flags().StringVar(&memberType, "member-type", "", "Filter the members by member type. Sample data: `user`, `machine`, `group`.")
		scim2GroupsCmd.AddCommand(cmd)
	}

	{ // get
		var orgId string
		var groupId string
		var excludedAttributes string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get a group",
			Long:  "Retrieve details for a group, by ID.\n\nOptionally, members can be retrieved with this request. The maximum number of members returned is 500.\n\n<br/>\n\n**Authorization**\n\nOAuth token rendered by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n- `identity:people_read`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_group_admin`\n\n- `id_readonly_admin`\n\n- `id_device_admin`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/identity/scim/{orgId}/v2/Groups/{groupId}")
				req.PathParam("orgId", orgId)
				req.PathParam("groupId", groupId)
				req.QueryParam("excludedAttributes", excludedAttributes)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "The ID of the organization to which this group belongs. If not specified, the organization ID from the OAuth token is used.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&groupId, "group-id", "", "A unique identifier for the group.")
		cmd.MarkFlagRequired("group-id")
		cmd.Flags().StringVar(&excludedAttributes, "excluded-attributes", "", "Attributes to be excluded from the return.")
		scim2GroupsCmd.AddCommand(cmd)
	}

	{ // update-put
		var orgId string
		var groupId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-put",
			Short: "Update a group with PUT",
			Long:  "Replace the contents of the Group.\n\nSpecify the group ID in the `groupId` parameter in the URI.\n\n<br/>\n\n**Authorization**\n\nOAuth token returned by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_group_admin`\n\n<br/>\n\n**Usage**:\n\n1. The input JSON must conform to one of the following schemas:\n    - `urn:ietf:params:scim:schemas:core:2.0:Group`\n    - `urn:scim:schemas:extension:cisco:webexidentity:2.0:Group`\n\n1. Unrecognized schemas (ID/section) are ignored.\n\n1. Read-only attributes provided as input values are ignored.\n\n1. The group `id` is not changed.\n\n1. All attributes are cleaned up if a new value is not provided by the client.\n\n1. The values, `meta` and `created` are not changed.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/identity/scim/{orgId}/v2/Groups/{groupId}")
				req.PathParam("orgId", orgId)
				req.PathParam("groupId", groupId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "The ID of the organization to which this group belongs. If not specified, the organization ID from the OAuth token is used.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&groupId, "group-id", "", "A unique identifier for the group.")
		cmd.MarkFlagRequired("group-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		scim2GroupsCmd.AddCommand(cmd)
	}

	{ // update-patch
		var orgId string
		var groupId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-patch",
			Short: "Update a group with PATCH",
			Long:  "Update group attributes with PATCH.\n\nSpecify the group ID in the `groupId` parameter in the URI.\n\n<br/>\n\n**Authorization**\n\nOAuth token returned by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_group_admin`\n\n<br/>\n\n**Usage**:\n\n1. The input JSON must conform to one of the following schemas:\n    - `urn:ietf:params:scim:schemas:core:2.0:Group`\n    - `urn:scim:schemas:extension:cisco:webexidentity:2.0:Group`\n\n1. Unrecognized schemas (ID/section) are ignored.\n\n1. Read-only attributes provided as input values are ignored.\n\n1. Each operation on an attribute must be compatible with the attribute's mutability.\n\n1. Each PATCH operation represents a single action to be applied to the\n   same SCIM resource specified by the request URI. Operations are\n   applied sequentially in the order they appear in the array. Each\n   operation in the sequence is applied to the target resource; the\n   resulting resource becomes the target of the next operation.\n   Evaluation continues until all operations are successfully applied or\n   until an error condition is encountered.\n\n<br/>\n\n**Add operations**:\n\nThe `add` operation is used to add a new attribute value to an existing resource. The operation must contain a `value` member whose content specifies the value to be added. The value may be a quoted value, or it may be a JSON object containing the sub-attributes of the complex attribute specified in the operation's `path`. The result of the add operation depends upon the target location indicated by `path` references:\n\n<br/>\n\n- If omitted, the target location is assumed to be the resource itself. The `value` parameter contains a set of attributes to be added to the resource.\n\n- If the target location does not exist, the attribute and value are added.\n\n- If the target location specifies a complex attribute, a set of sub-attributes shall be specified in the `value` parameter.\n\n- If the target location specifies a multi-valued attribute, a new value is added to the attribute.\n\n- If the target location specifies a single-valued attribute, the existing value is replaced.\n\n- If the target location specifies an attribute that does not exist (has no value), the attribute is added with the new value.\n\n- If the target location exists, the value is replaced.\n\n- If the target location already contains the value specified, no changes should be made to the resource.\n\n<br/>\n\n**Replace operations**:\n\nThe `replace` operation replaces the value at the target location specified by the `path`. The operation performs the following functions, depending on the target location specified by `path`:\n\n<br/>\n\n- If the `path` parameter is omitted, the target is assumed to be the resource itself. In this case, the `value` attribute shall contain a list of one or more attributes that are to be replaced.\n\n- If the target location is a single-value attribute, the value of the attribute is replaced.\n\n- If the target location is a multi-valued attribute and no filter is specified, the attribute and all values are replaced.\n\n- If the target location path specifies an attribute that does not exist, the service provider shall treat the operation as an \"add\".\n\n- If the target location specifies a complex attribute, a set of sub-attributes SHALL be specified in the `value` parameter, which replaces any existing values or adds where an attribute did not previously exist. Sub-attributes that are not specified in the `value` parameters are left unchanged.\n\n- If the target location is a multi-valued attribute and a value selection (\"valuePath\") filter is specified that matches one or more values of the multi-valued attribute, then all matching record values will be replaced.\n\n- If the target location is a complex multi-valued attribute with a value selection filter (\"valuePath\") and a specific sub-attribute (e.g., \"addresses[type eq \"work\"].streetAddress\"), the matching sub-attribute of all matching records is replaced.\n\n- If the target location is a multi-valued attribute for which a value selection filter (\"valuePath\") has been supplied and no record match was made, the service provider will indicate the failure by returning HTTP status code 400 and a `scimType` error code of `noTarget`.\n\n<br/>\n\n**Remove operations**:\n\nThe `remove` operation removes the value at the target location specified by the required attribute `path`. The operation performs the following functions, depending on the target location specified by `path`:\n\n<br/>\n\n- If `path` is unspecified, the operation fails with HTTP status code 400 and a \"scimType\" error code of \"noTarget\".\n\n- If the target location is a single-value attribute, the attribute and its associated value is removed, and the attribute will be considered unassigned.\n\n- If the target location is a multi-valued attribute and no filter is specified, the attribute and all values are removed, and the attribute SHALL be considered unassigned.\n\n- If the target location is a multi-valued attribute and a complex filter is specified comparing a `value`, the values matched by the filter are removed. If no other values remain after the removal of the selected values, the multi-valued attribute will be considered unassigned.\n\n- If the target location is a complex multi-valued attribute and a complex filter is specified based on the attribute`s sub-attributes, the matching records are removed. Sub-attributes whose values have been removed will be considered unassigned. If the complex multi-valued attribute has no remaining records, the attribute will be considered unassigned.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PATCH", "/identity/scim/{orgId}/v2/Groups/{groupId}")
				req.PathParam("orgId", orgId)
				req.PathParam("groupId", groupId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "The ID of the organization to which this group belongs. If not specified, the organization ID from the OAuth token is used.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&groupId, "group-id", "", "A unique identifier for the group.")
		cmd.MarkFlagRequired("group-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		scim2GroupsCmd.AddCommand(cmd)
	}

	{ // delete
		var orgId string
		var groupId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a group",
			Long:  "Remove a group from the system.\n\nSpecify the group ID in the `groupId` parameter in the URI.\n\n<br/>\n\n**Authorization**\n\nOAuth token rendered by Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_group_admin`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/identity/scim/{orgId}/v2/Groups/{groupId}")
				req.PathParam("orgId", orgId)
				req.PathParam("groupId", groupId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "The ID of the organization to which this group belongs. If not specified, the organization ID from the OAuth token is used.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&groupId, "group-id", "", "A unique identifier for the group.")
		cmd.MarkFlagRequired("group-id")
		scim2GroupsCmd.AddCommand(cmd)
	}

	{ // get-members
		var orgId string
		var groupId string
		var startIndex string
		var count string
		var memberType string
		cmd := &cobra.Command{
			Use:   "get-members",
			Short: "Get Group Members",
			Long:  "Returns the members of a group.\n\n- The default maximum number of members returned is 500.\n\n- Control parameters are available to page through the members and to control the size of the results.\n\n- Long result sets are split into [pages](/docs/basics#pagination).\n\n**Note**\nLocation groups are different from SCIM groups. You cannot search for identities in a location via groups.\n\n<br/>\n\n**Authorization**\n\nOAuth token returned by the Identity Broker.\n\n<br/>\n\nOne of the following OAuth scopes is required:\n\n- `identity:people_rw`\n\n- `identity:people_read`\n\n<br/>\n\nThe following administrators can use this API:\n\n- `id_full_admin`\n\n- `id_group_admin`\n\n- `id_readonly_admin`\n\n- `id_device_admin`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/identity/scim/{orgId}/v2/Groups/{groupId}/Members")
				req.PathParam("orgId", orgId)
				req.PathParam("groupId", groupId)
				req.QueryParam("startIndex", startIndex)
				req.QueryParam("count", count)
				req.QueryParam("memberType", memberType)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "The ID of the organization to which this group belongs. If not specified, the organization ID from the OAuth token is used.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&groupId, "group-id", "", "A unique identifier for the group.")
		cmd.MarkFlagRequired("group-id")
		cmd.Flags().StringVar(&startIndex, "start-index", "", "The index to start for group pagination.")
		cmd.Flags().StringVar(&count, "count", "", "Non-negative integer that specifies the desired number of search results per page. The maximum value for the count is 500.")
		cmd.Flags().StringVar(&memberType, "member-type", "", "Filter the members by member type. Sample data: `user`, `machine`, `group`.")
		scim2GroupsCmd.AddCommand(cmd)
	}

}
