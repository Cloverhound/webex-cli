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

var callParkCmd = &cobra.Command{
	Use:   "call-park",
	Short: "CallPark commands",
}

func init() {
	cmd.CallingCmd.AddCommand(callParkCmd)

	{ // list
		var locationId string
		var orgId string
		var max string
		var start string
		var order string
		var name string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "Read the List of Call Parks",
			Long:  "List all Call Parks for the organization.\n\nCall Park allows call recipients to place a call on hold so that it can be retrieved from another device.\n\nRetrieving this list requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.\n\n**NOTE**: The Call Park ID will change upon modification of the Call Park name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/callParks")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("order", order)
				req.QueryParam("name", name)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the list of call parks for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List call parks for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of call parks returned to this maximum count. Default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching call parks. Default is 0.")
		cmd.Flags().StringVar(&order, "order", "", "Sort the list of call parks by name, either ASC or DSC. Default is ASC.")
		cmd.Flags().StringVar(&name, "name", "", "Return the list of call parks that contains the given name. The maximum length is 80.")
		callParkCmd.AddCommand(cmd)
	}

	{ // create
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Call Park",
			Long:  "Create new Call Parks for the given location.\n\nCall Park allows call recipients to place a call on hold so that it can be retrieved from another device.\n\nCreating a call park requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The Call Park ID will change upon modification of the Call Park name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/callParks")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Create the call park for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create the call park for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callParkCmd.AddCommand(cmd)
	}

	{ // delete
		var locationId string
		var callParkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Call Park",
			Long:  "Delete the designated Call Park.\n\nCall Park allows call recipients to place a call on hold so that it can be retrieved from another device.\n\nDeleting a call park requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The Call Park ID will change upon modification of the Call Park name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/callParks/{callParkId}")
				req.PathParam("locationId", locationId)
				req.PathParam("callParkId", callParkId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location from which to delete a call park.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&callParkId, "call-park-id", "", "Delete the call park with the matching ID.")
		cmd.MarkFlagRequired("call-park-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the call park from this organization.")
		callParkCmd.AddCommand(cmd)
	}

	{ // get
		var locationId string
		var callParkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Details for a Call Park",
			Long:  "Retrieve Call Park details.\n\nCall Park allows call recipients to place a call on hold so that it can be retrieved from another device.\n\nRetrieving call park details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.\n\n**NOTE**: The Call Park ID will change upon modification of the Call Park name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/callParks/{callParkId}")
				req.PathParam("locationId", locationId)
				req.PathParam("callParkId", callParkId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve settings for a call park in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&callParkId, "call-park-id", "", "Retrieve settings for a call park with the matching ID.")
		cmd.MarkFlagRequired("call-park-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call park settings from this organization.")
		callParkCmd.AddCommand(cmd)
	}

	{ // update
		var locationId string
		var callParkId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Call Park",
			Long:  "Update the designated Call Park.\n\nCall Park allows call recipients to place a call on hold so that it can be retrieved from another device.\n\nUpdating a call park requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The Call Park ID will change upon modification of the Call Park name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/callParks/{callParkId}")
				req.PathParam("locationId", locationId)
				req.PathParam("callParkId", callParkId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call park exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&callParkId, "call-park-id", "", "Update settings for a call park with the matching ID.")
		cmd.MarkFlagRequired("call-park-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update call park settings from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callParkCmd.AddCommand(cmd)
	}

	{ // get-available-agents
		var locationId string
		var orgId string
		var callParkName string
		var max string
		var start string
		var name string
		var phoneNumber string
		var order string
		cmd := &cobra.Command{
			Use:   "get-available-agents",
			Short: "Get available agents from Call Parks",
			Long:  "Retrieve available agents from call parks for a given location.\n\nCall Park allows call recipients to place a call on hold so that it can be retrieved from another device.\n\nRetrieving available agents from call parks requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/callParks/availableUsers")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("callParkName", callParkName)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("order", order)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the available agents for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Return the available agents for this organization.")
		cmd.Flags().StringVar(&callParkName, "call-park-name", "", "Only return available agents from call parks with the matching name.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of available agents returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching available agents.")
		cmd.Flags().StringVar(&name, "name", "", "Only return available agents with the matching name.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Only return available agents with the matching primary number.")
		cmd.Flags().StringVar(&order, "order", "", "Order the available agents according to the designated fields. Up to three vertical bar (|) separated sort order fields may be specified. Available sort fields: fname, lname, number and extension. The maximum supported sort order value is 3.")
		callParkCmd.AddCommand(cmd)
	}

	{ // get-available-recall-hunt-groups
		var locationId string
		var orgId string
		var max string
		var start string
		var name string
		var order string
		cmd := &cobra.Command{
			Use:   "get-available-recall-hunt-groups",
			Short: "Get available recall hunt groups from Call Parks",
			Long:  "Retrieve available recall hunt groups from call parks for a given location.\n\nCall Park allows call recipients to place a call on hold so that it can be retrieved from another device.\n\nRetrieving available recall hunt groups from call parks requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/callParks/availableRecallHuntGroups")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("name", name)
				req.QueryParam("order", order)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the available recall hunt groups for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Return the available recall hunt groups for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of available recall hunt groups returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching available recall hunt groups.")
		cmd.Flags().StringVar(&name, "name", "", "Only return available recall hunt groups with the matching name.")
		cmd.Flags().StringVar(&order, "order", "", "Order the available recall hunt groups according to the designated fields. Available sort fields: lname.")
		callParkCmd.AddCommand(cmd)
	}

	{ // get-2
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-2",
			Short: "Get Call Park Settings",
			Long:  "Retrieve Call Park Settings from call parks for a given location.\n\nCall Park allows call recipients to place a call on hold so that it can be retrieved from another device.\n\nRetrieving settings from call parks requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/callParks/settings")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the call park settings for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Return the call park settings for this organization.")
		callParkCmd.AddCommand(cmd)
	}

	{ // update-2
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-2",
			Short: "Update Call Park settings",
			Long:  "Update Call Park settings for the designated location.\n\nCall Park allows call recipients to place a call on hold so that it can be retrieved from another device.\n\nUpdating call park settings requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/callParks/settings")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location for which call park settings will be updated.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update call park settings from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callParkCmd.AddCommand(cmd)
	}

	{ // list-extensions
		var orgId string
		var locationId string
		var max string
		var start string
		var extension string
		var locationName string
		var name string
		var order string
		cmd := &cobra.Command{
			Use:   "list-extensions",
			Short: "Read the List of Call Park Extensions",
			Long:  "List all Call Park Extensions for the organization.\n\nThe Call Park service, enabled for all users by default, allows a user to park a call against an available user's extension or to a Call Park Extension. Call Park Extensions are extensions defined within the Call Park service for holding parked calls.\n\nRetrieving this list requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/callParkExtensions")
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationId", locationId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("extension", extension)
				req.QueryParam("locationName", locationName)
				req.QueryParam("name", name)
				req.QueryParam("order", order)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List call park extensions for this organization.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Only return call park extensions with matching location ID.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&extension, "extension", "", "Only return call park extensions with the matching extension.")
		cmd.Flags().StringVar(&locationName, "location-name", "", "Only return call park extensions with the matching extension.")
		cmd.Flags().StringVar(&name, "name", "", "Only return call park extensions with the matching name.")
		cmd.Flags().StringVar(&order, "order", "", "Order the available agents according to the designated fields.  Available sort fields: `groupName`, `callParkExtension`, `callParkExtensionName`, `callParkExtensionExternalId`.")
		callParkCmd.AddCommand(cmd)
	}

	{ // get-extension
		var locationId string
		var callParkExtensionId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-extension",
			Short: "Get Details for a Call Park Extension",
			Long:  "Retrieve Call Park Extension details.\n\nThe Call Park service, enabled for all users by default, allows a user to park a call against an available user's extension or to a Call Park Extension. Call Park Extensions are extensions defined within the Call Park service for holding parked calls.\n\nRetrieving call park extension details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/callParkExtensions/{callParkExtensionId}")
				req.PathParam("locationId", locationId)
				req.PathParam("callParkExtensionId", callParkExtensionId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve details for a call park extension in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&callParkExtensionId, "call-park-extension-id", "", "Retrieve details for a call park extension with the matching ID.")
		cmd.MarkFlagRequired("call-park-extension-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call park extension details from this organization.")
		callParkCmd.AddCommand(cmd)
	}

	{ // delete-extension
		var locationId string
		var callParkExtensionId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-extension",
			Short: "Delete a Call Park Extension",
			Long:  "Delete the designated Call Park Extension.\n\nCall Park Extension enables a call recipient to park a call to an extension, so someone else within the same Organization can retrieve the parked call by dialing that extension. Call Park Extensions can be added as monitored lines by users' Cisco phones, so users can park and retrieve calls by pressing the associated phone line key.\n\nDeleting a call park extension requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/callParkExtensions/{callParkExtensionId}")
				req.PathParam("locationId", locationId)
				req.PathParam("callParkExtensionId", callParkExtensionId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location from which to delete a call park extension.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&callParkExtensionId, "call-park-extension-id", "", "Delete the call park extension with the matching ID.")
		cmd.MarkFlagRequired("call-park-extension-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the call park extension from this organization.")
		callParkCmd.AddCommand(cmd)
	}

	{ // update-extension
		var locationId string
		var callParkExtensionId string
		var orgId string
		var name string
		var extension string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-extension",
			Short: "Update a Call Park Extension",
			Long:  "Update the designated Call Park Extension.\n\nCall Park Extension enables a call recipient to park a call to an extension, so someone else within the same Organization can retrieve the parked call by dialing that extension. Call Park Extensions can be added as monitored lines by users' Cisco phones, so users can park and retrieve calls by pressing the associated phone line key.\n\nUpdating a call park extension requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/callParkExtensions/{callParkExtensionId}")
				req.PathParam("locationId", locationId)
				req.PathParam("callParkExtensionId", callParkExtensionId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("extension", extension)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call park extension exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&callParkExtensionId, "call-park-extension-id", "", "Update a call park extension with the matching ID.")
		cmd.MarkFlagRequired("call-park-extension-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update a call park extension from this organization.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&extension, "extension", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callParkCmd.AddCommand(cmd)
	}

	{ // create-extension
		var locationId string
		var orgId string
		var name string
		var extension string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-extension",
			Short: "Create a Call Park Extension",
			Long:  "Create new Call Park Extensions for the given location.\n\nCall Park Extension enables a call recipient to park a call to an extension, so someone else within the same Organization can retrieve the parked call by dialing that extension. Call Park Extensions can be added as monitored lines by users' Cisco phones, so users can park and retrieve calls by pressing the associated phone line key.\n\nCreating a call park extension requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/callParkExtensions")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("extension", extension)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Create the call park extension for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create the call park extension for this organization.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&extension, "extension", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callParkCmd.AddCommand(cmd)
	}

}
