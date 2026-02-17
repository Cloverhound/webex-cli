package calling

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

var pagingGroupCmd = &cobra.Command{
	Use:   "paging-group",
	Short: "PagingGroup commands",
}

func init() {
	cmd.CallingCmd.AddCommand(pagingGroupCmd)

	{ // list
		var orgId string
		var max string
		var start string
		var locationId string
		var name string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "Read the List of Paging Groups",
			Long:  "List all Paging Groups for the organization.\n\nGroup Paging allows a person to place a one-way call or group page to up to 75 people and/or workspaces by\ndialing a number or extension assigned to a specific paging group. The Group Paging service makes a simultaneous call to all the assigned targets.\n\nRetrieving this list requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/paging")
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("locationId", locationId)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List paging groups for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count. Default is 2000")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects. Default is 0")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return only paging groups with matching location ID. Default is all locations")
		cmd.Flags().StringVar(&name, "name", "", "Return only paging groups with the matching name.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Return only paging groups with matching primary phone number or extension.")
		pagingGroupCmd.AddCommand(cmd)
	}

	{ // create
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a new Paging Group",
			Long:  "Create a new Paging Group for the given location.\n\nGroup Paging allows a one-way call or group page to up to 75 people, workspaces and virtual lines by\ndialing a number or extension assigned to a specific paging group. The Group Paging service makes a simultaneous call to all the assigned targets.\n\nCreating a paging group requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, and `dialByName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/paging")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Create the paging group for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create the paging group for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		pagingGroupCmd.AddCommand(cmd)
	}

	{ // delete
		var locationId string
		var pagingId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Paging Group",
			Long:  "Delete the designated Paging Group.\n\nGroup Paging allows a person to place a one-way call or group page to up to 75 people and/or workspaces by\ndialing a number or extension assigned to a specific paging group. The Group Paging service makes a simultaneous call to all the assigned targets.\n\nDeleting a paging group requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/paging/{pagingId}")
				req.PathParam("locationId", locationId)
				req.PathParam("pagingId", pagingId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location from which to delete a paging group.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&pagingId, "paging-id", "", "Delete the paging group with the matching ID.")
		cmd.MarkFlagRequired("paging-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the paging group from this organization.")
		pagingGroupCmd.AddCommand(cmd)
	}

	{ // get
		var locationId string
		var pagingId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Details for a Paging Group",
			Long:  "Retrieve Paging Group details.\n\nGroup Paging allows a person, place or virtual line a one-way call or group page to up to 75 people and/or workspaces and/or virtual line by\ndialing a number or extension assigned to a specific paging group. The Group Paging service makes a simultaneous call to all the assigned targets.\n\nRetrieving paging group details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, and `dialByName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/paging/{pagingId}")
				req.PathParam("locationId", locationId)
				req.PathParam("pagingId", pagingId)
				req.QueryParam("orgId", orgId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve settings for a paging group in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&pagingId, "paging-id", "", "Retrieve settings for the paging group with this identifier.")
		cmd.MarkFlagRequired("paging-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve paging group settings from this organization.")
		pagingGroupCmd.AddCommand(cmd)
	}

	{ // update
		var locationId string
		var pagingId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Paging Group",
			Long:  "Update the designated Paging Group.\n\nGroup Paging allows a person to place a one-way call or group page to up to 75 people, workspaces and virtual lines by\ndialing a number or extension assigned to a specific paging group. The Group Paging service makes a simultaneous call to all the assigned targets.\n\nUpdating a paging group requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, and `dialByName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/paging/{pagingId}")
				req.PathParam("locationId", locationId)
				req.PathParam("pagingId", pagingId)
				req.QueryParam("orgId", orgId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update settings for a paging group in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&pagingId, "paging-id", "", "Update settings for the paging group with this identifier.")
		cmd.MarkFlagRequired("paging-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update paging group settings from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		pagingGroupCmd.AddCommand(cmd)
	}

	{ // get-primary-numbers
		var locationId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-primary-numbers",
			Short: "Get Paging Group Primary Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as the paging group's primary phone number.\nThese numbers are associated with the location specified in the request URL, can be active or inactive, and are unassigned.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/paging/availableNumbers")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the list of phone numbers for this location within the given organization. The maximum length is 36.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		pagingGroupCmd.AddCommand(cmd)
	}

}
