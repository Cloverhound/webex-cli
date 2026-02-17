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

var orgContactsCmd = &cobra.Command{
	Use:   "org-contacts",
	Short: "OrgContacts commands",
}

func init() {
	cmd.AdminCmd.AddCommand(orgContactsCmd)

	{ // create
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Contact",
			Long:  "Creating a new contact for a given organization requires an org admin role.\n\nAt least one of the following body parameters: `phoneNumbers`, `emails`, `sipAddresses` is required to create a new contact for source \"CH\",\n`displayName` is required to create a new contact for source \"Webex4Broadworks\".\n\nUse the optional `groupIds` field to add group IDs in an array within the organisation contact. This will become a group contact.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/contacts/organizations/{orgId}/contacts")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for the user's organization or the organization he manages.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		orgContactsCmd.AddCommand(cmd)
	}

	{ // get
		var orgId string
		var contactId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get a Contact",
			Long:  "Shows details for an organization contact by ID.\nSpecify the organization ID in the `orgId` parameter in the URI, and specify the contact ID in the `contactId` parameter in the URI.\n\n**NOTE**:\nThe `orgId` used in the path for this API are the org UUIDs. They follow a xxxx-xxxx-xxxx-xxxx pattern. If you have an orgId in base64 encoded format (starting with Y2.....) you need to base64 decode the id and extract the UUID from the slug, before you use it in your API call.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/contacts/organizations/{orgId}/contacts/{contactId}")
				req.PathParam("orgId", orgId)
				req.PathParam("contactId", contactId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for the user's organization or the organization he manages.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&contactId, "contact-id", "", "The contact ID.")
		cmd.MarkFlagRequired("contact-id")
		orgContactsCmd.AddCommand(cmd)
	}

	{ // update
		var orgId string
		var contactId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Contact",
			Long:  "Update details for contact by ID. Only an admin can update a contact.\nSpecify the organization ID in the `orgId` parameter in the URI, and specify the contact ID in the `contactId` parameter in the URI.\n\nUse the optional `groupIds` field to update the group IDs by changing the existing array. You can add or remove one or all groups. To remove all associated groups, pass an empty array in the `groupIds` field.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PATCH", "/contacts/organizations/{orgId}/contacts/{contactId}")
				req.PathParam("orgId", orgId)
				req.PathParam("contactId", contactId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for the user's organization or the organization he manages.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&contactId, "contact-id", "", "The contact ID.")
		cmd.MarkFlagRequired("contact-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		orgContactsCmd.AddCommand(cmd)
	}

	{ // delete
		var orgId string
		var contactId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Contact",
			Long:  "Remove a contact from the organization. Only an admin can remove a contact.\n\nSpecify the organization ID in the `orgId` parameter in the URI, and specify the contact ID in the `contactId` parameter in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/contacts/organizations/{orgId}/contacts/{contactId}")
				req.PathParam("orgId", orgId)
				req.PathParam("contactId", contactId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for the user's organization or the organization he manages.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&contactId, "contact-id", "", "The contact ID.")
		cmd.MarkFlagRequired("contact-id")
		orgContactsCmd.AddCommand(cmd)
	}

	{ // list
		var orgId string
		var keyword string
		var source string
		var limit string
		var groupIds string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Contacts",
			Long:  "List contacts in the organization. The default limit is `100`.\n\n`keyword` can be the value of \"displayName\", \"firstName\", \"lastName\", \"email\". An empty string of `keyword` means get all contacts.\n\n`groupIds` is a comma separated list group IDs. Results are filtered based on those group IDs.\n\nLong result sets will be split into [pages](/docs/basics#pagination).",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/contacts/organizations/{orgId}/contacts/search")
				req.PathParam("orgId", orgId)
				req.QueryParam("keyword", keyword)
				req.QueryParam("source", source)
				req.QueryParam("limit", limit)
				req.QueryParam("groupIds", groupIds)
				req.QueryParam("groupIds", groupIds)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&keyword, "keyword", "", "List contacts with a keyword.")
		cmd.Flags().StringVar(&source, "source", "", "List contacts with source.")
		cmd.Flags().StringVar(&limit, "limit", "", "Limit the maximum number of contact in the response.         + Default: 100 ")
		cmd.Flags().StringVar(&groupIds, "group-ids", "", "Filter contacts based on groups.")
		orgContactsCmd.AddCommand(cmd)
	}

	{ // bulk-create-update
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "bulk-create-update",
			Short: "Bulk Create or Update Contacts",
			Long:  "Create or update contacts in bulk. Update an existing contact by specifying the contact ID in the `contactId` parameter in the request body.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/contacts/organizations/{orgId}/contacts/bulk")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for the user's organization or the organization he manages.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		orgContactsCmd.AddCommand(cmd)
	}

	{ // bulk-delete
		var orgId string
		var schemas string
		var objectIds []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "bulk-delete",
			Short: "Bulk Delete Contacts",
			Long:  `Delete contacts in bulk.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/contacts/organizations/{orgId}/contacts/bulk/delete")
				req.PathParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("schemas", schemas)
					req.BodyStringSlice("objectIds", objectIds)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Webex Identity assigned organization identifier for the user's organization or the organization he manages.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&schemas, "schemas", "", "")
		cmd.Flags().StringSliceVar(&objectIds, "object-ids", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		orgContactsCmd.AddCommand(cmd)
	}

}
