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

var dectDevicesCmd = &cobra.Command{
	Use:   "dect-devices",
	Short: "DectDevices commands",
}

func init() {
	cmd.CallingCmd.AddCommand(dectDevicesCmd)

	{ // create-network
		var locationId string
		var orgId string
		var name string
		var model string
		var defaultAccessCodeEnabled bool
		var defaultAccessCode string
		var displayName string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-network",
			Short: "Create a DECT Network",
			Long:  "Create a multi-cell DECT network for a given location.\n\nCreating a DECT network requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/dectNetworks")
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
					req.BodyString("model", model)
					req.BodyBool("defaultAccessCodeEnabled", defaultAccessCodeEnabled, cmd.Flags().Changed("default-access-code-enabled"))
					req.BodyString("defaultAccessCode", defaultAccessCode)
					req.BodyString("displayName", displayName)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Create a DECT network in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create a DECT network in this organization.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&model, "model", "", "")
		cmd.Flags().BoolVar(&defaultAccessCodeEnabled, "default-access-code-enabled", false, "")
		cmd.Flags().StringVar(&defaultAccessCode, "default-access-code", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // list-networks
		var orgId string
		var name string
		var locationId string
		cmd := &cobra.Command{
			Use:   "list-networks",
			Short: "Get the List of DECT Networks for an organization",
			Long:  "Retrieves the list of DECT networks for an organization.\n\nDECT Networks provide roaming voice services via base stations and wireless handsets. A DECT network can be provisioned up to 1000 lines across up to 254 base stations.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/dectNetworks")
				req.QueryParam("orgId", orgId)
				req.QueryParam("name", name)
				req.QueryParam("locationId", locationId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List of DECT networks in this organization.")
		cmd.Flags().StringVar(&name, "name", "", "List of DECT networks with this name.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "List of DECT networks at this location.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // get-network
		var locationId string
		var dectNetworkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-network",
			Short: "Get DECT Network Details",
			Long:  "Retrieves the details of a DECT network.\n\nDECT Networks provide roaming voice services via base stations and wireless handsets. A DECT network can be provisioned up to 1000 lines across up to 254 base stations.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Details of the DECT network at this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Details of the specified DECT network.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Details of the DECT network in this organization.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // update-network
		var locationId string
		var dectNetworkId string
		var orgId string
		var name string
		var defaultAccessCodeEnabled bool
		var defaultAccessCode string
		var displayName string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-network",
			Short: "Update DECT Network",
			Long:  "Update the details of a DECT network.\n\nDECT Networks provide roaming voice services via base stations and wireless handsets. A DECT network can be provisioned up to 1000 lines across up to 254 base stations.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyBool("defaultAccessCodeEnabled", defaultAccessCodeEnabled, cmd.Flags().Changed("default-access-code-enabled"))
					req.BodyString("defaultAccessCode", defaultAccessCode)
					req.BodyString("displayName", displayName)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update DECT network details in the specified location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Update DECT network details in the specified DECT network.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update DECT network details in the specified organization.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().BoolVar(&defaultAccessCodeEnabled, "default-access-code-enabled", false, "")
		cmd.Flags().StringVar(&defaultAccessCode, "default-access-code", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // delete-network
		var locationId string
		var dectNetworkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-network",
			Short: "Delete DECT Network",
			Long:  "Delete a DECT network.\n\nDECT Networks provide roaming voice services via base stations and wireless handsets. A DECT network can be provisioned up to 1000 lines across up to 254 base stations.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Delete the DECT network in the specified location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Delete the specified DECT network.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the DECT network in the specified organization.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // create-multiple-base-stations
		var locationId string
		var dectNetworkId string
		var orgId string
		var baseStationMacs []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-multiple-base-stations",
			Short: "Create Multiple Base Stations",
			Long:  "This API is used to create multiple base stations in a DECT network in an organization.\n\nCreating base stations in a DECT network requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/baseStations")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("baseStationMacs", baseStationMacs)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Create a base station in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Create a base station for the DECT network.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create a base station for a DECT network in this organization.")
		cmd.Flags().StringSliceVar(&baseStationMacs, "base-station-macs", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // list-network-base-stations
		var locationId string
		var dectNetworkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-network-base-stations",
			Short: "Get a list of DECT Network Base Stations",
			Long:  "Retrieve a list of base stations in a DECT Network.\n\nA DECT network supports 2 types of base stations, DECT DBS-110 Single-Cell and DECT DBS-210 Multi-Cell.\nA DECT DBS-110 allows up to 30 lines of registration and supports 1 base station only. A DECT DBS-210 can have up to 254 base stations and supports up to 1000 lines of registration.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/baseStations")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location containing the DECT network.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Retrieve the list of base stations in the specified DECT network ID.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization containing the DECT network.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // delete-bulk-network-base-stations
		var locationId string
		var dectNetworkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-bulk-network-base-stations",
			Short: "Delete bulk DECT Network Base Stations",
			Long:  "Delete all the base stations in the DECT Network.\n\nA DECT network supports 2 types of base stations, DECT DBS-110 Single-Cell and DECT DBS-210 Multi-Cell.\nA DECT DBS-110 allows up to 30 lines of registration and supports 1 base station only. A DECT DBS-210 can have up to 254 base stations and supports up to 1000 lines of registration.\n\nThis API requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/baseStations")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location containing the DECT network.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Delete all the base stations in the specified DECT network ID.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization containing the DECT network.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // get-network-base-station
		var locationId string
		var dectNetworkId string
		var baseStationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-network-base-station",
			Short: "Get the details of a specific DECT Network Base Station",
			Long:  "Retrieve details of a specific base station in the DECT Network.\n\nA DECT network supports 2 types of base stations, DECT DBS-110 Single-Cell and DECT DBS-210 Multi-Cell.\nA DECT DBS-110 allows up to 30 lines of registration and supports 1 base station only. A DECT DBS-210 can have up to 254 base stations and supports up to 1000 lines of registration.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/baseStations/{baseStationId}")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.PathParam("baseStationId", baseStationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location containing the DECT network.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Retrieve details of a specific base station in the specified DECT network ID.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&baseStationId, "base-station-id", "", "Retrieve details of the specific DECT base station ID.")
		cmd.MarkFlagRequired("base-station-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization containing the DECT network.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // delete-network-base-station
		var locationId string
		var dectNetworkId string
		var baseStationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-network-base-station",
			Short: "Delete a specific DECT Network Base Station",
			Long:  "Delete a specific base station in the DECT Network.\n\nA DECT network supports 2 types of base stations, DECT DBS-110 Single-Cell and DECT DBS-210 Multi-Cell.\nA DECT DBS-110 allows up to 30 lines of registration and supports 1 base station only. A DECT DBS-210 can have up to 254 base stations and supports up to 1000 lines of registration.\n\nThis API requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/baseStations/{baseStationId}")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.PathParam("baseStationId", baseStationId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location containing the DECT network.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Delete a specific base station in the specified DECT network ID.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&baseStationId, "base-station-id", "", "Delete the specific DECT base station ID.")
		cmd.MarkFlagRequired("base-station-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization containing the DECT network.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // create-handset-network
		var locationId string
		var dectNetworkId string
		var orgId string
		var line1MemberId string
		var customDisplayName string
		var line2MemberId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-handset-network",
			Short: "Add a Handset to a DECT Network",
			Long:  "Add a handset to a DECT network in a location in an organization.\n\nAdding a handset to a DECT network requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n<div><Callout type=\"warning\">Adding a DECT handset to a person with a Webex Calling Standard license will disable Webex Calling across their Webex mobile, tablet, desktop, and browser applications.</Callout></div>\n\n<div><Callout type=\"warning\">Adding or removing handsets to the DECT network in less than 90 seconds may result in base station not having the latest configuration until the base station is rebooted.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/handsets")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("line1MemberId", line1MemberId)
					req.BodyString("customDisplayName", customDisplayName)
					req.BodyString("line2MemberId", line2MemberId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Add handset in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "A unique identifier for the DECT network.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Add handset in this organization.")
		cmd.Flags().StringVar(&line1MemberId, "line1-member-id", "", "")
		cmd.Flags().StringVar(&customDisplayName, "custom-display-name", "", "")
		cmd.Flags().StringVar(&line2MemberId, "line2-member-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // list-handsets-network-id
		var locationId string
		var dectNetworkId string
		var orgId string
		var basestationId string
		var memberId string
		cmd := &cobra.Command{
			Use:   "list-handsets-network-id",
			Short: "Get List of Handsets for a DECT Network ID",
			Long:  "List all the handsets associated with a DECT Network ID.\n\nA handset can have up to two lines, and a DECT network supports a total of 120 lines across all handsets.\nA member on line1 of a DECT handset can be of type PEOPLE or PLACE while a member on line2 of a DECT handset can be of type PEOPLE, PLACE, or VIRTUAL_LINE.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/handsets")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("basestationId", basestationId)
				req.QueryParam("memberId", memberId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location containing the DECT network.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Search handset details in the specified DECT network ID.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization containing the DECT network.")
		cmd.Flags().StringVar(&basestationId, "basestation-id", "", "Search handset details in the specified DECT base station ID.")
		cmd.Flags().StringVar(&memberId, "member-id", "", "ID of the member of the handset. Members can be of type PEOPLE, PLACE, or VIRTUAL_LINE.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // get-network-handset
		var locationId string
		var dectNetworkId string
		var handsetId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-network-handset",
			Short: "Get Specific DECT Network Handset Details",
			Long:  "List the specific DECT Network handset details.\n\nA handset can have up to two lines, and a DECT network supports a total of 120 lines across all handsets.\nA member on line1 of a DECT handset can be of type PEOPLE or PLACE while a member on line2 of a DECT handset can be of type PEOPLE, PLACE, or VIRTUAL_LINE.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/handsets/{handsetId}")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.PathParam("handsetId", handsetId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location containing the DECT network.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Search handset details in the specified DECT network ID.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&handsetId, "handset-id", "", "A unique identifier for the handset.")
		cmd.MarkFlagRequired("handset-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization containing the DECT network.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // update-network-handset
		var locationId string
		var dectNetworkId string
		var handsetId string
		var orgId string
		var line1MemberId string
		var customDisplayName string
		var line2MemberId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-network-handset",
			Short: "Update DECT Network Handset",
			Long:  "Update the line assignment on a handset.\n\nA handset can have up to two lines, and a DECT network supports a total of 120 lines across all handsets.\nA member on line1 of a DECT handset can be of type PEOPLE or PLACE while a member on line2 of a DECT handset can be of type PEOPLE, PLACE, or VIRTUAL_LINE.\n\nUpdating a DECT Network handset requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n<div><Callout type=\"warning\">Adding a person with a Webex Calling Standard license to the DECT handset line1 will disable Webex Calling across their Webex mobile, tablet, desktop, and browser applications. \nRemoving a person with a Webex Calling Standard license from the DECT handset line1 will enable Webex Calling across their Webex mobile, tablet, desktop, and browser applications.</Callout></div>\n\n<div><Callout type=\"warning\">Adding or removing handsets to the DECT network in less than 90 seconds may result in base station not having the latest configuration until the base station is rebooted.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/handsets/{handsetId}")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.PathParam("handsetId", handsetId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("line1MemberId", line1MemberId)
					req.BodyString("customDisplayName", customDisplayName)
					req.BodyString("line2MemberId", line2MemberId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location containing the DECT network.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Update handset details in the specified DECT network.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&handsetId, "handset-id", "", "A unique identifier for the handset.")
		cmd.MarkFlagRequired("handset-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization containing the DECT network.")
		cmd.Flags().StringVar(&line1MemberId, "line1-member-id", "", "")
		cmd.Flags().StringVar(&customDisplayName, "custom-display-name", "", "")
		cmd.Flags().StringVar(&line2MemberId, "line2-member-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // delete-network-handset
		var locationId string
		var dectNetworkId string
		var handsetId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-network-handset",
			Short: "Delete specific DECT Network Handset Details",
			Long:  "Delete a specific DECT Network handset.\n\nA handset can have up to two lines, and a DECT network supports a total of 120 lines across all handsets.\nA member on line1 of a DECT handset can be of type PEOPLE or PLACE while a member on line2 of a DECT handset can be of type PEOPLE, PLACE, or VIRTUAL_LINE.\n\nThis API requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n<div><Callout type=\"warning\">Deleting a DECT handset from a person with a Webex Calling Standard license will enable Webex Calling across their Webex mobile, tablet, desktop, and browser applications.</Callout></div>\n\n<div><Callout type=\"warning\">Adding or removing handsets to the DECT network in less than 90 seconds may result in base station not having the latest configuration until the base station is rebooted.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/handsets/{handsetId}")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.PathParam("handsetId", handsetId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location containing the DECT network.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Delete handset details in the specified DECT network ID.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&handsetId, "handset-id", "", "A unique identifier for the handset.")
		cmd.MarkFlagRequired("handset-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization containing the DECT network.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // delete-multiple-handsets
		var locationId string
		var dectNetworkId string
		var orgId string
		var handsetIds []string
		var deleteAll bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "delete-multiple-handsets",
			Short: "Delete multiple handsets",
			Long:  "Delete multiple handsets or all of them.\n\nA handset can have up to two lines, and a DECT network supports a total of 120 lines across all handsets.\nA member on line1 of a DECT handset can be of type PEOPLE or PLACE while a member on line2 of a DECT handset can be of type PEOPLE, PLACE, or VIRTUAL_LINE.\n\nThis API requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n<div><Callout type=\"warning\">Deleting a DECT handset from a person with a Webex Calling Standard license will enable Webex Calling across their Webex mobile, tablet, desktop, and browser applications.</Callout></div>\n\n<div><Callout type=\"warning\">Adding or removing handsets to the DECT network in less than 90 seconds may result in base station not having the latest configuration until the base station is rebooted.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/handsets")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("handsetIds", handsetIds)
					req.BodyBool("deleteAll", deleteAll, cmd.Flags().Changed("delete-all"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location containing the DECT network.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Delete handset details in the specified DECT network ID.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization containing the DECT network.")
		cmd.Flags().StringSliceVar(&handsetIds, "handset-ids", nil, "")
		cmd.Flags().BoolVar(&deleteAll, "delete-all", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // create-list-handsets-network
		var locationId string
		var dectNetworkId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-list-handsets-network",
			Short: "Add a List of Handsets to a DECT Network",
			Long:  "Add a list of up to 50 handsets to a DECT network in a location.\n\nA DECT network acts as a container that can support up to 1,000 lines across all handsets, with each handset capable of handling up to two lines. Once the network is created, you can add bases, handsets, and assign users or lines as needed.\n\nAdding a list of handsets to a DECT network requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n<div><Callout type=\"warning\">Adding a DECT handset to a person with a Webex Calling Standard license will disable Webex Calling across their Webex mobile, tablet, desktop, and browser applications.\n\n</Callout></div>\n\n<div><Callout type=\"warning\">Adding or removing handsets to the DECT network in less than 90 seconds may result in base station not having the latest configuration until the base station is rebooted.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/handsets/bulk")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Add handsets in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "A unique identifier for the DECT network.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Add handsets in this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // list-networks-associated-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-networks-associated-person",
			Short: "GET List of DECT networks associated with a Person",
			Long:  "Retrieves the list of DECT networks for a person in an organization.\n\nDECT Network provides roaming voice services via base stations and wireless handsets. DECT network can be provisioned up to 1000 lines across up to 254 base stations.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/dectNetworks")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "List of DECT networks associated with this person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List of DECT networks associated with a person in this organization.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // list-networks-associated-workspace
		var workspaceId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-networks-associated-workspace",
			Short: "GET List of DECT networks associated with a workspace",
			Long:  "Retrieves the list of DECT networks for a workspace in an organization.\n\nDECT Network provides roaming voice services via base stations and wireless handsets. DECT network can be provisioned up to 1000 lines across up to 254 base stations.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/workspaces/{workspaceId}/dectNetworks")
				req.PathParam("workspaceId", workspaceId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "List of DECT networks associated with this workspace.")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List of DECT networks associated with a workspace in this organization.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // search-available-members
		var orgId string
		var start string
		var max string
		var memberName string
		var phoneNumber string
		var extension string
		var order string
		var locationId string
		var excludeVirtualLine string
		var usageType string
		cmd := &cobra.Command{
			Use:   "search-available-members",
			Short: "Search Available Members",
			Long:  "List the members that are available to be assigned to DECT handset lines.\n\nThis requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/devices/availableMembers")
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
				req.QueryParam("max", max)
				req.QueryParam("memberName", memberName)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("extension", extension)
				req.QueryParam("order", order)
				req.QueryParam("locationId", locationId)
				req.QueryParam("excludeVirtualLine", excludeVirtualLine)
				req.QueryParam("usageType", usageType)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Search members in this organization.")
		cmd.Flags().StringVar(&start, "start", "", "Specifies the offset from the first result that you want to fetch.")
		cmd.Flags().StringVar(&max, "max", "", "Specifies the maximum number of records that you want to fetch.")
		cmd.Flags().StringVar(&memberName, "member-name", "", "Search (Contains) numbers based on member name.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Search (Contains) based on number.")
		cmd.Flags().StringVar(&extension, "extension", "", "Search (Contains) based on extension.")
		cmd.Flags().StringVar(&order, "order", "", "Sort the list of available members on the device in ascending order by name, using either last name `lname` or first name `fname`. Default sort is the last name in ascending order.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "List members for the location ID.")
		cmd.Flags().StringVar(&excludeVirtualLine, "exclude-virtual-line", "", "If true, search results will exclude virtual lines in the member list. NOTE: Virtual lines cannot be assigned as the primary line.")
		cmd.Flags().StringVar(&usageType, "usage-type", "", "Search for members eligible to become the owner of the device, or share line on the device.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // generate-enable-service-password
		var locationId string
		var dectNetworkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "generate-enable-service-password",
			Short: "Generate and Enable DECT Serviceability Password",
			Long:  "Generates and enables a 16-character DECT serviceability password.\n\n<div><Callout type=\"warning\">Generating a password and transmitting it to the DECT network can reboot the entire network. Be sure you choose an appropriate time to generate a new password.</Callout></div>\n\nThe DECT serviceability password, also known as the admin override password, provides read/write access to DECT base stations for performing system serviceability and troubleshooting functions.\n\nThis API requires either a full administrator auth token with the scope `spark-admin:telephony_config_write`, or a device administrator token with the scope of `spark-admin:devices_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/serviceabilityPassword/actions/generate/invoke")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Unique identifier for the DECT network.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // get-service-password-status
		var locationId string
		var dectNetworkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-service-password-status",
			Short: "Get DECT Serviceability Password status",
			Long:  "Retrieves the DECT serviceability password status.\n\n<div><Callout type=\"info\">If the serviceability password is enabled but has not been generated, the `enabled` status will be returned as `true` even though there is no active serviceability password.</Callout></div>\n\nThe DECT serviceability password, also known as the admin override password, provides read/write access to DECT base stations for performing system serviceability and troubleshooting functions.\n\nThis API requires an auth token with either a full, read-only token with the scope of `spark-admin:telephony_config_read`, or a device administrator token with the scope of `spark-admin:devices_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/serviceabilityPassword")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Unique identifier for the DECT network.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		dectDevicesCmd.AddCommand(cmd)
	}

	{ // update-service-password-status
		var locationId string
		var dectNetworkId string
		var orgId string
		var enabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-service-password-status",
			Short: "Update DECT Serviceability Password status",
			Long:  "Enables or disables the DECT serviceability password.\n\n<div><Callout type=\"warning\">Enabling or disabling the password and transmitting it to the DECT network can reboot the entire network. Be sure you choose an appropriate time for this action.</Callout></div>\n\n<div><Callout type=\"info\">If enabling is requested, but the serviceability password has not been generated, we will not actively reject the request even though there is no serviceability password.</Callout></div>\n\nThe DECT serviceability password, also known as the admin override password, provides read/write access to DECT base stations for performing system serviceability and troubleshooting functions.\n\nThis API requires either a full administrator auth token with the scope `spark-admin:telephony_config_write`, or a device administrator token with the scope of `spark-admin:devices_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/dectNetworks/{dectNetworkId}/serviceabilityPassword")
				req.PathParam("locationId", locationId)
				req.PathParam("dectNetworkId", dectNetworkId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&dectNetworkId, "dect-network-id", "", "Unique identifier for the DECT network.")
		cmd.MarkFlagRequired("dect-network-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dectDevicesCmd.AddCommand(cmd)
	}

}
