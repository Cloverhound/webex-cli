package calling

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

var locationsCmd = &cobra.Command{
	Use:   "locations",
	Short: "Locations commands",
}

func init() {
	cmd.CallingCmd.AddCommand(locationsCmd)

	{ // list
		var name string
		var id string
		var orgId string
		var max string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Locations",
			Long:  "List locations for an organization.\n\n* Use query parameters to filter the result set by location name, ID, or organization.\n\n* Long result sets will be split into [pages](/docs/basics#pagination).\n\n* Searching and viewing locations in your organization requires an administrator or location administrator auth token with any of the following scopes: `spark-admin:locations_read`, `spark-admin:people_read` or `spark-admin:device_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/locations")
				req.QueryParam("name", name)
				req.QueryParam("id", id)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
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
		cmd.Flags().StringVar(&name, "name", "", "List locations whose name contains this string (case-insensitive).")
		cmd.Flags().StringVar(&id, "id", "", "List locations by ID.")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List locations in this organization. Only admin users of another organization (such as partners) may use this parameter.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of location in the response.")
		locationsCmd.AddCommand(cmd)
	}

	{ // create
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Location",
			Long:  "Create a new Location for a given organization. Only an admin in the organization can create a new Location.\n\n* Creating a location in your organization requires a full administrator auth token with a scope of `spark-admin:locations_write`.\n\n* Partners may specify `orgId` query parameter to create location in managed organization.\n\n* The following body parameters are required to create a new location: \n    * `name`\n    * `timeZone`\n    * `preferredLanguage`\n    * `address`\n    * `announcementLanguage`.\n\n* `latitude`, `longitude` and `notes` are optional parameters to create a new location.\n\n* **Important:** While the `name` field supports up to 256 characters, locations that will be enabled for Webex Calling must have names with a maximum of 80 characters. If you plan to enable calling for this location, ensure the name does not exceed 80 characters to maintain compatibility with Control Hub and calling features.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/locations")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create a location common attribute for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationsCmd.AddCommand(cmd)
	}

	{ // get
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Location Details",
			Long:  "Shows details for a location, by ID.\n\n* Specify the location ID in the `locationId` parameter in the URI.\n\n* Use query parameter `orgId` to filter the result set by organization(optional).\n\n* Searching and viewing location in your organization requires an administrator or location administrator auth token with any of the following scopes:\n\n    * `spark-admin:locations_read`\n    * `spark-admin:people_read`\n    * `spark-admin:device_read`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/locations/{locationId}")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "A unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Get location common attributes for this organization.")
		locationsCmd.AddCommand(cmd)
	}

	{ // update
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Location",
			Long:  "Update details for a location, by ID.\n\n* Updating a location in your organization requires a full administrator or location administrator auth token with a scope of `spark-admin:locations_write`.\n\n* Specify the location ID in the `locationId` parameter in the URI.\n\n* Partners may specify `orgId` query parameter to update location in managed organization.\n\n* **Important:** While the `name` field supports up to 256 characters, locations that are enabled for Webex Calling must have names with a maximum of 80 characters. If the location is enabled for calling, ensure the name does not exceed 80 characters to maintain compatibility with Control Hub and calling features.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/locations/{locationId}")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update location common attributes for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update location common attributes for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationsCmd.AddCommand(cmd)
	}

	{ // list-floors
		var locationId string
		cmd := &cobra.Command{
			Use:   "list-floors",
			Short: "List Location Floors",
			Long:  "List location floors.\nRequires an administrator auth token with the `spark-admin:locations_read` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/locations/{locationId}/floors")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "A unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		locationsCmd.AddCommand(cmd)
	}

	{ // create-floor
		var locationId string
		var floorNumber int64
		var displayName string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-floor",
			Short: "Create a Location Floor",
			Long:  "Create a new floor in the given location. The `displayName` parameter is optional, and omitting it will result in the creation of a floor without that value set.\nRequires an administrator auth token with the `spark-admin:locations_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/locations/{locationId}/floors")
				req.PathParam("locationId", locationId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyInt("floorNumber", floorNumber, cmd.Flags().Changed("floor-number"))
					req.BodyString("displayName", displayName)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "A unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().Int64Var(&floorNumber, "floor-number", 0, "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationsCmd.AddCommand(cmd)
	}

	{ // get-floor
		var locationId string
		var floorId string
		cmd := &cobra.Command{
			Use:   "get-floor",
			Short: "Get Location Floor Details",
			Long:  "Shows details for a floor, by ID. Specify the floor ID in the `floorId` parameter in the URI.\nRequires an administrator auth token with the `spark-admin:locations_read` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/locations/{locationId}/floors/{floorId}")
				req.PathParam("locationId", locationId)
				req.PathParam("floorId", floorId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "A unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&floorId, "floor-id", "", "A unique identifier for the floor.")
		cmd.MarkFlagRequired("floor-id")
		locationsCmd.AddCommand(cmd)
	}

	{ // update-floor
		var locationId string
		var floorId string
		var floorNumber int64
		var displayName string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-floor",
			Short: "Update a Location Floor",
			Long:  "Updates details for a floor, by ID. Specify the floor ID in the `floorId` parameter in the URI. Include all details for the floor returned by a previous call to [Get Location Floor Details](/docs/api/v1/locations/get-location-floor-details). Omitting the optional `displayName` field will result in that field no longer being defined for the floor.\nRequires an administrator auth token with the `spark-admin:locations_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/locations/{locationId}/floors/{floorId}")
				req.PathParam("locationId", locationId)
				req.PathParam("floorId", floorId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyInt("floorNumber", floorNumber, cmd.Flags().Changed("floor-number"))
					req.BodyString("displayName", displayName)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "A unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&floorId, "floor-id", "", "A unique identifier for the floor.")
		cmd.MarkFlagRequired("floor-id")
		cmd.Flags().Int64Var(&floorNumber, "floor-number", 0, "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationsCmd.AddCommand(cmd)
	}

	{ // delete-floor
		var locationId string
		var floorId string
		cmd := &cobra.Command{
			Use:   "delete-floor",
			Short: "Delete a Location Floor",
			Long:  "Deletes a floor, by ID.\nRequires an administrator auth token with the `spark-admin:locations_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/locations/{locationId}/floors/{floorId}")
				req.PathParam("locationId", locationId)
				req.PathParam("floorId", floorId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "A unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&floorId, "floor-id", "", "A unique identifier for the floor.")
		cmd.MarkFlagRequired("floor-id")
		locationsCmd.AddCommand(cmd)
	}

	{ // delete
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete Location",
			Long:  "Delete a location, by ID.\n\n* Specify the location ID in the `locationId` parameter in the URI.\n\n* Deleting a location in your organization requires a full administrator auth token with a scope of `spark-admin:locations_write`.\n\n* NOTE: Disabling Webex Calling for a Webex Calling enabled location is required prior to deleting a location.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/locations/{locationId}")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "A unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Specify the organization for the location to be deleted.")
		locationsCmd.AddCommand(cmd)
	}

}
