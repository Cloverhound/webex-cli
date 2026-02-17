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

var peopleCmd = &cobra.Command{
	Use:   "people",
	Short: "People commands",
}

func init() {
	cmd.AdminCmd.AddCommand(peopleCmd)

	{ // list
		var email string
		var displayName string
		var id string
		var orgId string
		var roles string
		var callingData string
		var locationId string
		var max string
		var excludeStatus string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List People",
			Long:  "List people in your organization. For most users, either the `email` or `displayName` parameter is required. Admin users can omit these fields and list all users in their organization.\n\nResponse properties associated with a user's presence status, such as `status` or `lastActivity`, will only be returned for people within your organization or an organization you manage. Presence information will not be returned if the authenticated user has [disabled status sharing](https://help.webex.com/nkzs6wl/). Calling /people frequently to poll `status` information for a large set of users will quickly lead to `429` errors and throttling of such requests and is therefore discouraged.\n\nAdmin users can include `Webex Calling` (BroadCloud) user details in the response by specifying `callingData` parameter as `true`. Admin users can list all users in a location or with a specific phone number. Admin users will receive an enriched payload with additional administrative fields like `licenses`,`roles`, `locations` etc. These fields are shown when accessing a user via GET /people/{id}, not when doing a GET /people?id=\n\nLookup by `email` is only supported for people within the same org or where a partner admin relationship is in place.\n\nLookup by `roles` is only supported for Admin users for the people within the same org.\n\nLong result sets will be split into [pages](/docs/basics#pagination).",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people")
				req.QueryParam("email", email)
				req.QueryParam("displayName", displayName)
				req.QueryParam("id", id)
				req.QueryParam("orgId", orgId)
				req.QueryParam("roles", roles)
				req.QueryParam("callingData", callingData)
				req.QueryParam("locationId", locationId)
				req.QueryParam("max", max)
				req.QueryParam("excludeStatus", excludeStatus)
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
		cmd.Flags().StringVar(&email, "email", "", "List people with this email address. For non-admin requests, either this or `displayName` are required. With the exception of partner admins and a managed org relationship, people lookup by email is only available for users in the same org.")
		cmd.Flags().StringVar(&displayName, "display-name", "", "List people whose name starts with this string. For non-admin requests, either this or email are required.")
		cmd.Flags().StringVar(&id, "id", "", "List people by ID. Accepts up to 85 person IDs separated by commas. If this parameter is provided then presence information (such as the `lastActivity` or `status` properties) will not be included in the response.")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List people in this organization. Only admin users of another organization (such as partners) may use this parameter.")
		cmd.Flags().StringVar(&roles, "roles", "", "List of roleIds separated by commas.")
		cmd.Flags().StringVar(&callingData, "calling-data", "", "Include Webex Calling user details in the response.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "List people present in this location.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of people in the response. If `callingData`=true, then `max` will not be more than 100. If `locationId` is specified then `max` will not be more than 50.")
		cmd.Flags().StringVar(&excludeStatus, "exclude-status", "", "Omit people status/availability to enhance query performance.")
		peopleCmd.AddCommand(cmd)
	}

	{ // create-person
		var callingData string
		var minResponse string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-person",
			Short: "Create a Person",
			Long:  "Create a new user account for a given organization. Only an admin can create a new user account.\n\nAt least one of the following body parameters is required to create a new user: `displayName`, `firstName`, `lastName`.\n\nCurrently, users may have only one email address associated with their account. The `emails` parameter is an array, which accepts multiple values to allow for future expansion, but currently only one email address will be used for the new user.\n\nAdmin users can include `Webex calling` (BroadCloud) user details in the response by specifying `callingData` parameter as true. It may happen that the POST request with calling data returns a 400 status, but the person was created still. One way to get into this state is if an invalid phone number is assigned to a user. The people API aggregates calls to several other microservices, and one may have failed. A best practice is to check if the user exists before retrying. This can be done with the user's email address and a GET /people.\n\nWhen doing attendee management, append `#attendee` to the `siteUrl` parameter (e.g. `mysite.webex.com#attendee`) to make the new user an attendee for a site.\n\n**NOTES**:\n\n* For creating a `Webex Calling` user, you must provide `phoneNumbers` or `extension`, `locationId`, and `licenses` string in the same request.\n\n* `SipAddresses` are asigned via an asynchronous process. This means that the POST response may not show the SIPAddresses immediately. Instead you can verify them with a separate GET to /people, after they were newly configured.\n\n* When assigning multiple licenses in a single request, the system will assign all valid and available licenses. If any requested licenses cannot be assigned, the operation will continue with the remaining licenses. As a result, it is possible that not all requested licenses are assigned to the user.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/people")
				req.QueryParam("callingData", callingData)
				req.QueryParam("minResponse", minResponse)
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
		cmd.Flags().StringVar(&callingData, "calling-data", "", "Include Webex Calling user details in the response.")
		cmd.Flags().StringVar(&minResponse, "min-response", "", "Set to `true` to improve performance by omitting person details and returning only the ID in the response when successful. If unsuccessful the response will have optional error details.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		peopleCmd.AddCommand(cmd)
	}

	{ // get-person
		var personId string
		var callingData string
		cmd := &cobra.Command{
			Use:   "get-person",
			Short: "Get Person Details",
			Long:  "Shows details for a person, by ID.\n\nResponse properties associated with a user's presence status, such as `status` or `lastActivity`, will only be displayed for people within your organization or an organization you manage. Presence information will not be shown if the authenticated user has [disabled status sharing](https://help.webex.com/nkzs6wl/).\n\nAdmin users can include `Webex Calling` (BroadCloud) user details in the response by specifying `callingData` parameter as `true`.\n\nSpecify the person ID in the `personId` parameter in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}")
				req.PathParam("personId", personId)
				req.QueryParam("callingData", callingData)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&callingData, "calling-data", "", "Include Webex Calling user details in the response.")
		peopleCmd.AddCommand(cmd)
	}

	{ // update-person
		var personId string
		var callingData string
		var showAllTypes string
		var minResponse string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person",
			Short: "Update a Person",
			Long:  "Update details for a person, by ID.\n\nSpecify the person ID in the `personId` parameter in the URI. Only an admin can update a person details.\n\nInclude all details for the person. This action expects all user details to be present in the request. A common approach is to first [GET the person's details](/docs/api/v1/people/get-person-details), make changes, then PUT both the changed and unchanged values.\n\nAdmin users can include `Webex Calling` (BroadCloud) user details in the response by specifying `callingData` parameter as true.\n\nWhen doing attendee management, to update a user from host role to an attendee for a site append `#attendee` to the respective `siteUrl` and remove the meeting host license for this site from the license array.\nTo update a person from an attendee role to a host for a site, add the meeting license for this site in the meeting array, and remove that site from the `siteurl` parameter.\n\nTo remove the attendee privilege for a user on a meeting site, remove the `sitename#attendee` from the `siteUrl`s array. The `showAllTypes` parameter must be set to `true`.\n\n**NOTE**:\n\n* The `locationId` can only be set when assigning a calling license to a user. It cannot be changed if a user is already an existing calling user.\n\n* The `extension` field should be used to update the Webex Calling extension for a person. The extension value should not include the location routing prefix. The `work_extension` type in the `phoneNumbers` object as seen in the response payload of [List People](/docs/api/v1/people/list-people) or [Get Person Details](/docs/api/v1/people/get-person-details), cannot be used to set the Webex Calling extension for a person.\n\n* When updating a user with multiple email addresses using a PUT request, ensure that the primary email address is listed first in the array. Note that the order of email addresses returned by a GET request is not guaranteed..\n\n* The People API is a combination of several microservices, each responsible for specific attributes of a person. As a result, a PUT request that returns an error response code may still have altered some values of the person's data. Therefore, it is recommended to perform a GET request after encountering an error to verify the current state of the resource. \n\n* Some licenses are implicitly assigned by the system and cannot be admin controlled. They are necessary for the baseline function of the Webex system. If you get an error about implicitly assigned licensed that cannot be removed, please ensure you have the corresponding license in your PUT request.\n\n* When assigning multiple licenses in a single request, the system will assign all valid and available licenses. If any requested licenses cannot be assigned, the operation will continue with the remaining licenses. As a result, it is possible that not all requested licenses are assigned to the user.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}")
				req.PathParam("personId", personId)
				req.QueryParam("callingData", callingData)
				req.QueryParam("showAllTypes", showAllTypes)
				req.QueryParam("minResponse", minResponse)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&callingData, "calling-data", "", "Include Webex Calling user details in the response.")
		cmd.Flags().StringVar(&showAllTypes, "show-all-types", "", "Include additional user data like `#attendee` role.")
		cmd.Flags().StringVar(&minResponse, "min-response", "", "Set to `true` to improve performance by omitting person details in the response. If unsuccessful the response will have optional error details.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		peopleCmd.AddCommand(cmd)
	}

	{ // delete-person
		var personId string
		cmd := &cobra.Command{
			Use:   "delete-person",
			Short: "Delete a Person",
			Long:  "Remove a person from the system.\n\n**Required Administrator Roles:**\n\nThe following administrators have permission to use this API:\n\n**Customer Organization:**\n- Full administrator\n- User administrator\n\n**Partner/External Access:**\n- External full administrator\n\n**Note:** External read-only administrators, provisioning administrators, and device administrators cannot delete users.\n\nSpecify the person ID in the `personId` parameter in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/people/{personId}")
				req.PathParam("personId", personId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		peopleCmd.AddCommand(cmd)
	}

	{ // get-my-own
		var callingData string
		cmd := &cobra.Command{
			Use:   "get-my-own",
			Short: "Get My Own Details",
			Long:  "Get profile details for the authenticated user. This is the same as GET `/people/{personId}` using the Person ID associated with your Auth token.\n\nAdmin users can include `Webex Calling` (BroadCloud) user details in the response by specifying `callingData` parameter as true.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/me")
				req.QueryParam("callingData", callingData)
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
		cmd.Flags().StringVar(&callingData, "calling-data", "", "Include Webex Calling user details in the response.")
		peopleCmd.AddCommand(cmd)
	}

}
