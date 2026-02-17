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

var virtualExtensionsCmd = &cobra.Command{
	Use:   "virtual-extensions",
	Short: "VirtualExtensions commands",
}

func init() {
	cmd.CallingCmd.AddCommand(virtualExtensionsCmd)

	{ // create
		var orgId string
		var displayName string
		var phoneNumber string
		var extension string
		var firstName string
		var lastName string
		var locationId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Virtual Extension",
			Long:  "Create new Virtual Extension for the given organization or location.\n\nYou can set up virtual extensions at the organization or location level. The organization level enables everyone across your organization to dial the same extension number to reach someone.\nYou can use the location level virtual extension like any other extension assigned to the specific location.\nUsers at the specific location can dial the extension. However, users at other locations can reach the virtual extension by dialing the ESN.\n\nCreating a virtual extension requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write` and `identity:contacts_rw`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualExtensions")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("displayName", displayName)
					req.BodyString("phoneNumber", phoneNumber)
					req.BodyString("extension", extension)
					req.BodyString("firstName", firstName)
					req.BodyString("lastName", lastName)
					req.BodyString("locationId", locationId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "")
		cmd.Flags().StringVar(&extension, "extension", "", "")
		cmd.Flags().StringVar(&firstName, "first-name", "", "")
		cmd.Flags().StringVar(&lastName, "last-name", "", "")
		cmd.Flags().StringVar(&locationId, "location-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // list
		var orgId string
		var max string
		var start string
		var order string
		var extension string
		var phoneNumber string
		var name string
		var locationName string
		var locationId string
		var orgLevelOnly string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "Read the List of Virtual Extensions",
			Long:  "Retrieve virtual extensions associated with a specific customer.\n\nThe GET Virtual Extensions API allows administrators to retrieve a list of virtual extensions configured within their organization. Virtual extensions enable users to dial extension numbers that route to external phone numbers, such as those of remote workers or frequently contacted clients.\nThis API returns key information including the  extension, associated  phone number (in E.164 format), display name, and the location to which the virtual extension belongs\nThe API supports filtering by various parameters, such as extension number, phone number, and location name. The results can be paginated using the `max` and `start` parameters, and the order of the results can be specified using the `order` parameter.\n\nRetrieving a Virtual Extension requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualExtensions")
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("order", order)
				req.QueryParam("extension", extension)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("name", name)
				req.QueryParam("locationName", locationName)
				req.QueryParam("locationId", locationId)
				req.QueryParam("orgLevelOnly", orgLevelOnly)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of virtual extensions returned to this maximum count. Default is 10.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching virtual extensions. Default is 0.")
		cmd.Flags().StringVar(&order, "order", "", "Order the list of virtual extensions in ascending or descending order. Default is ascending.")
		cmd.Flags().StringVar(&extension, "extension", "", "Filter the list of virtual extensions by extension number.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter the list of virtual extensions by phone number.")
		cmd.Flags().StringVar(&name, "name", "", "Filter the list of virtual extensions by name. This can be either first name or last name.")
		cmd.Flags().StringVar(&locationName, "location-name", "", "Filter the list of virtual extensions by location name.(Only one of the locationName, locationId, and OrgLevelOnly query parameters is allowed at the same time.)")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Filter the list of virtual extensions by location ID.")
		cmd.Flags().StringVar(&orgLevelOnly, "org-level-only", "", "Filter the list of virtual extensions by organization level. If orgLevelOnly is true, return only the organization level virtual extensions.")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // get
		var extensionId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get a Virtual Extension",
			Long:  "Retrieve Virtual Extension details for the given extension ID.\n\nVirtual extensions integrate remote workers on separate telephony systems into Webex Calling, enabling users to reach them via extension dialing.\nThis endpoint allows administrators to retrieve configuration details for a specific virtual extension, ensuring visibility into the mapping between extensions and external phone numbers.\n\nRetrieving a Virtual Extension requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualExtensions/{extensionId}")
				req.PathParam("extensionId", extensionId)
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
		cmd.Flags().StringVar(&extensionId, "extension-id", "", "ID of the virtual extension.")
		cmd.MarkFlagRequired("extension-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // update
		var extensionId string
		var orgId string
		var firstName string
		var lastName string
		var displayName string
		var phoneNumber string
		var extension string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Virtual Extension",
			Long:  "Update Virtual Extension details for the given extension ID.\n\nThis API updates the configuration of an existing virtual extension identified by its unique extension ID. Administrators can modify fields such as the extension, associated phone number (in E.164 format), display name, and location etc.\n\nUpdating a Virtual Extension requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write` and `identity:contacts_rw`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualExtensions/{extensionId}")
				req.PathParam("extensionId", extensionId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("firstName", firstName)
					req.BodyString("lastName", lastName)
					req.BodyString("displayName", displayName)
					req.BodyString("phoneNumber", phoneNumber)
					req.BodyString("extension", extension)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&extensionId, "extension-id", "", "ID of the virtual extension.")
		cmd.MarkFlagRequired("extension-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().StringVar(&firstName, "first-name", "", "")
		cmd.Flags().StringVar(&lastName, "last-name", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "")
		cmd.Flags().StringVar(&extension, "extension", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // delete
		var extensionId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Virtual Extension",
			Long:  "Delete Virtual Extension using the extension ID.\n\nThis API permanently deletes a virtual extension from the organization. Once deleted, the extension will no longer route calls to the external phone number, and users won\u2019t be able to reach it via the assigned extension.\n\nDeleting a Virtual Extension requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write` and `identity:contacts_rw`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/virtualExtensions/{extensionId}")
				req.PathParam("extensionId", extensionId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&extensionId, "extension-id", "", "ID of the virtual extension.")
		cmd.MarkFlagRequired("extension-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // get-settings
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-settings",
			Short: "Get Virtual extension settings",
			Long:  "Retrieve Virtual Extension settings for the given Org.\n\nThis API retrieves the virtual extension mode settings configured for a given organization. Virtual extensions can operate in two modes: STANDARD and ENHANCED. The selected mode determines how the system handles routing and signaling for virtual extensions.\nBy default, the virtual extensions that you create use the Standard mode. Another mode, enhanced signaling mode, is available to all customers, however, virtual extensions won't function properly in this mode unless your PSTN provider supports special network signaling extensions and there aren't many PSTN providers that do.\n\nRetrieving a Virtual Extension settings requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualExtensions/settings")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // update-settings
		var orgId string
		var mode string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-settings",
			Short: "Modify Virtual Extension Settings",
			Long:  "Update Virtual Extension details for the given extension ID.\n\nThis endpoint updates the virtual extension settings for an organization. It is primarily used to configure the operating mode for virtual extensions.\nModes determine how virtual extensions are assigned or managed within the system.\n\nUpdating a Virtual Extension requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualExtensions/settings")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("mode", mode)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().StringVar(&mode, "mode", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // validate-external-phone-number
		var orgId string
		var phoneNumbers []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "validate-external-phone-number",
			Short: "Validate an external phone number",
			Long:  "Validate external phone number for the given organization.\n\nThis API is designed to validate external phone numbers before they are assigned as virtual extensions for a customer.\nIt ensures that the provided numbers are properly formatted, eligible for use, and not already in use within the system.\nThis validation is typically part of a pre-check process during provisioning or number assignment workflows, helping administrators or systems prevent conflicts or errors related to number reuse or format issues.\n\nCreating a virtual extension requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualExtensions/actions/validateNumbers/invoke")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("phoneNumbers", phoneNumbers)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // create-range
		var orgId string
		var name string
		var prefix string
		var patterns []string
		var locationId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-range",
			Short: "Create a Virtual Extension Range",
			Long:  "Create a new Virtual Extension Range for the given organization or location.\n\nVirtual extension ranges integrate remote workers on a separate telephony system into Webex Calling and enable extension dialing. Using these ranges, you can define patterns that can be used to route calls at a location level or an organization level. You are allowed to define virtual extensions ranges in addition to individual virtual extensions.\nThis works in both Standard and Enhanced modes\n\nVirtual extension range can be set up at the organization or location level.\n\nCreating a virtual extension range requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualExtensionRanges")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("prefix", prefix)
					req.BodyStringSlice("patterns", patterns)
					req.BodyString("locationId", locationId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&prefix, "prefix", "", "")
		cmd.Flags().StringSliceVar(&patterns, "patterns", nil, "")
		cmd.Flags().StringVar(&locationId, "location-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // list-range
		var orgId string
		var max string
		var start string
		var order string
		var name string
		var prefix string
		var locationId string
		var orgLevelOnly string
		cmd := &cobra.Command{
			Use:   "list-range",
			Short: "Get a list of a Virtual Extension Range",
			Long:  "Retrieves the list of Virtual Extension Ranges.\n\nVirtual extension ranges integrate remote workers on a separate telephony system into Webex Calling and enable extension dialing. Using these ranges, you can define patterns that can be used to route calls at a location level or an organization level. You are allowed to define virtual extensions ranges in addition to individual virtual extensions.\nThis works in both Standard and Enhanced modes\n\nRetrieving a virtual extension range requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualExtensionRanges")
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("order", order)
				req.QueryParam("name", name)
				req.QueryParam("prefix", prefix)
				req.QueryParam("locationId", locationId)
				req.QueryParam("orgLevelOnly", orgLevelOnly)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().StringVar(&max, "max", "", "Maximum number of results to return.")
		cmd.Flags().StringVar(&start, "start", "", "The starting index of the results to return.")
		cmd.Flags().StringVar(&order, "order", "", "Sort the list of virtual extension ranges by name or prefix, either ASC or DSC. Default sort order is ASC.")
		cmd.Flags().StringVar(&name, "name", "", "Filter the list of virtual extension ranges by name.")
		cmd.Flags().StringVar(&prefix, "prefix", "", "Filter the list of virtual extension ranges by prefix.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Filter the list of virtual extension ranges by location ID. Only one of the `locationId` and `OrgLevelOnly` query parameters is allowed at the same time.")
		cmd.Flags().StringVar(&orgLevelOnly, "org-level-only", "", "Filter the list of virtual extension ranges by organization level. If `orgLevelOnly` is true, return only the organization level virtual extension ranges.")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // get-range
		var extensionRangeId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-range",
			Short: "Get details of a Virtual Extension Range",
			Long:  "Retrieve virtual extension range details for the given extension range ID.\n\nVirtual extension ranges integrate remote workers on a separate telephony system into Webex Calling and enable extension dialing. Using these ranges, you can define patterns that can be used to route calls at a location level or an organization level. You are allowed to define virtual extensions ranges in addition to individual virtual extensions.\nThis works in both Standard and Enhanced modes\n\nRetrieving a virtual extension range requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualExtensionRanges/{extensionRangeId}")
				req.PathParam("extensionRangeId", extensionRangeId)
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
		cmd.Flags().StringVar(&extensionRangeId, "extension-range-id", "", "ID of the virtual extension range.")
		cmd.MarkFlagRequired("extension-range-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // delete-range
		var extensionRangeId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-range",
			Short: "Delete a Virtual Extension Range",
			Long:  "Delete a virtual extension range for the given extension range ID.\n\nVirtual extension ranges integrate remote workers on a separate telephony system into Webex Calling and enable extension dialing. Using these ranges, you can define patterns that can be used to route calls at a location level or an organization level. You are allowed to define virtual extensions ranges in addition to individual virtual extensions.\nThis works in both Standard and Enhanced modes\n\nDeleting a virtual extension range requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/virtualExtensionRanges/{extensionRangeId}")
				req.PathParam("extensionRangeId", extensionRangeId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&extensionRangeId, "extension-range-id", "", "ID of the virtual extension range.")
		cmd.MarkFlagRequired("extension-range-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // update-range
		var extensionRangeId string
		var orgId string
		var name string
		var prefix string
		var patterns []string
		var action string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-range",
			Short: "Modify Virtual Extension Range",
			Long:  "Modify virtual extension range for the given extension range ID.\n\nVirtual extension ranges integrate remote workers on a separate telephony system into Webex Calling and enable extension dialing. Using these ranges, you can define patterns that can be used to route calls at a location level or an organization level. You are allowed to define virtual extensions ranges in addition to individual virtual extensions.\nThis works in both Standard and Enhanced modes\n\nModifying a virtual extension range requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualExtensionRanges/{extensionRangeId}")
				req.PathParam("extensionRangeId", extensionRangeId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("prefix", prefix)
					req.BodyStringSlice("patterns", patterns)
					req.BodyString("action", action)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&extensionRangeId, "extension-range-id", "", "ID of the virtual extension range.")
		cmd.MarkFlagRequired("extension-range-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&prefix, "prefix", "", "")
		cmd.Flags().StringSliceVar(&patterns, "patterns", nil, "")
		cmd.Flags().StringVar(&action, "action", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualExtensionsCmd.AddCommand(cmd)
	}

	{ // validate-prefix-pattern-range
		var orgId string
		var locationId string
		var name string
		var prefix string
		var patterns []string
		var rangeId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "validate-prefix-pattern-range",
			Short: "Validate the prefix and extension pattern for a Virtual Extension Range",
			Long:  "Validate the prefix and extension pattern for a Virtual Extension Range.\n\nVirtual extension ranges integrate remote workers on a separate telephony system into Webex Calling and enable extension dialing. Using these ranges, you can define patterns that can be used to route calls at a location level or an organization level. You are allowed to define virtual extensions ranges in addition to individual virtual extensions.\nThis works in both Standard and Enhanced modes\n\nValidating a prefix and extension pattern for a Virtual Extension Range requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualExtensionRanges/actions/validate/invoke")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("locationId", locationId)
					req.BodyString("name", name)
					req.BodyString("prefix", prefix)
					req.BodyStringSlice("patterns", patterns)
					req.BodyString("rangeId", rangeId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&prefix, "prefix", "", "")
		cmd.Flags().StringSliceVar(&patterns, "patterns", nil, "")
		cmd.Flags().StringVar(&rangeId, "range-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualExtensionsCmd.AddCommand(cmd)
	}

}
