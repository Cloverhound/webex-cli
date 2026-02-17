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

var virtualLineCallCmd = &cobra.Command{
	Use:   "virtual-line-call",
	Short: "VirtualLineCall commands",
}

func init() {
	cmd.CallingCmd.AddCommand(virtualLineCallCmd)

	{ // list
		var orgId string
		var locationId string
		var max string
		var start string
		var id string
		var ownerName string
		var phoneNumber string
		var locationName string
		var order string
		var hasDeviceAssigned string
		var hasExtensionAssigned string
		var hasDnAssigned string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "Read the List of Virtual Lines",
			Long:  "List all Virtual Lines for the organization.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines to Webex Calling users.\n\nRetrieving this list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines")
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationId", locationId)
				req.QueryParam("locationId", locationId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("id", id)
				req.QueryParam("id", id)
				req.QueryParam("ownerName", ownerName)
				req.QueryParam("ownerName", ownerName)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("locationName", locationName)
				req.QueryParam("locationName", locationName)
				req.QueryParam("order", order)
				req.QueryParam("order", order)
				req.QueryParam("hasDeviceAssigned", hasDeviceAssigned)
				req.QueryParam("hasExtensionAssigned", hasExtensionAssigned)
				req.QueryParam("hasDnAssigned", hasDnAssigned)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List virtual lines for this organization.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the list of virtual lines matching these location ids. Example for multiple values - `?locationId=locId1&locationId=locId2`.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&id, "id", "", "Return the list of virtual lines matching these virtualLineIds. Example for multiple values - `?id=id1&id=id2`.")
		cmd.Flags().StringVar(&ownerName, "owner-name", "", "Return the list of virtual lines matching these owner names. Example for multiple values - `?ownerName=name1&ownerName=name2`.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Return the list of virtual lines matching these phone numbers. Example for multiple values - `?phoneNumber=number1&phoneNumber=number2`.")
		cmd.Flags().StringVar(&locationName, "location-name", "", "Return the list of virtual lines matching the location names. Example for multiple values - `?locationName=loc1&locationName=loc2`.")
		cmd.Flags().StringVar(&order, "order", "", "Return the list of virtual lines based on the order. Default sort will be in an Ascending order. Maximum 3 orders allowed at a time. Example for multiple values - `?order=order1&order=order2`.")
		cmd.Flags().StringVar(&hasDeviceAssigned, "has-device-assigned", "", "If `true`, includes only virtual lines with devices assigned. When not explicitly specified, the default includes both virtual lines with devices assigned and not assigned.")
		cmd.Flags().StringVar(&hasExtensionAssigned, "has-extension-assigned", "", "If `true`, includes only virtual lines with an extension assigned. When not explicitly specified, the default includes both virtual lines with extension assigned and not assigned.")
		cmd.Flags().StringVar(&hasDnAssigned, "has-dn-assigned", "", "If `true`, includes only virtual lines with an assigned directory number, also known as a Dn. When not explicitly specified, the default includes both virtual lines with a Dn assigned and not assigned.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // create
		var orgId string
		var firstName string
		var lastName string
		var locationId string
		var displayName string
		var phoneNumber string
		var extension string
		var callerIdLastName string
		var callerIdFirstName string
		var callerIdNumber string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Virtual Line",
			Long:  "Create new Virtual Line for the given location.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines to Webex Calling users.\n\nCreating a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualLines")
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
					req.BodyString("locationId", locationId)
					req.BodyString("displayName", displayName)
					req.BodyString("phoneNumber", phoneNumber)
					req.BodyString("extension", extension)
					req.BodyString("callerIdLastName", callerIdLastName)
					req.BodyString("callerIdFirstName", callerIdFirstName)
					req.BodyString("callerIdNumber", callerIdNumber)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create the virtual line for this organization.")
		cmd.Flags().StringVar(&firstName, "first-name", "", "")
		cmd.Flags().StringVar(&lastName, "last-name", "", "")
		cmd.Flags().StringVar(&locationId, "location-id", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "")
		cmd.Flags().StringVar(&extension, "extension", "", "")
		cmd.Flags().StringVar(&callerIdLastName, "caller-id-last-name", "", "")
		cmd.Flags().StringVar(&callerIdFirstName, "caller-id-first-name", "", "")
		cmd.Flags().StringVar(&callerIdNumber, "caller-id-number", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-recording
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-recording",
			Short: "Read Call Recording Settings for a Virtual Line",
			Long:  "Retrieve Virtual Line's Call Recording settings.\n\nThe Call Recording feature provides a hosted mechanism to record the calls placed and received on the Carrier platform for replay and archival. This feature is helpful for quality assurance, security, training, and more.\n\nThis API requires a full or user administrator auth token with the `spark-admin:telephony_config_read` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/callRecording")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-recording
		var virtualLineId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-recording",
			Short: "Configure Call Recording Settings for a Virtual Line",
			Long:  "Configure virtual line's Call Recording settings.\n\nThe Call Recording feature provides a hosted mechanism to record the calls placed and received on the Carrier platform for replay and archival. This feature is helpful for quality assurance, security, training, and more.\n\nThis API requires a full or user administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/callRecording")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual profile resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // delete
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Virtual Line",
			Long:  "Delete the designated Virtual Line.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines to Webex Calling users.\n\nDeleting a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write` and `identity:contacts_rw`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/virtualLines/{virtualLineId}")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Delete the virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the virtual line from this organization.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Details for a Virtual Line",
			Long:  "Retrieve Virtual Line details.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines to Webex Calling users.\n\nRetrieving virtual line details requires a full or user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve virtual line settings from this organization.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update
		var virtualLineId string
		var orgId string
		var firstName string
		var lastName string
		var displayName string
		var phoneNumber string
		var extension string
		var announcementLanguage string
		var callerIdLastName string
		var callerIdFirstName string
		var callerIdNumber string
		var timeZone string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Virtual Line",
			Long:  "Update the designated Virtual Line.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines to Webex Calling users.\n\nUpdating a virtual line requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write` and `identity:contacts_rw`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}")
				req.PathParam("virtualLineId", virtualLineId)
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
					req.BodyString("announcementLanguage", announcementLanguage)
					req.BodyString("callerIdLastName", callerIdLastName)
					req.BodyString("callerIdFirstName", callerIdFirstName)
					req.BodyString("callerIdNumber", callerIdNumber)
					req.BodyString("timeZone", timeZone)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Update settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update virtual line settings from this organization.")
		cmd.Flags().StringVar(&firstName, "first-name", "", "")
		cmd.Flags().StringVar(&lastName, "last-name", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "")
		cmd.Flags().StringVar(&extension, "extension", "", "")
		cmd.Flags().StringVar(&announcementLanguage, "announcement-language", "", "")
		cmd.Flags().StringVar(&callerIdLastName, "caller-id-last-name", "", "")
		cmd.Flags().StringVar(&callerIdFirstName, "caller-id-first-name", "", "")
		cmd.Flags().StringVar(&callerIdNumber, "caller-id-number", "", "")
		cmd.Flags().StringVar(&timeZone, "time-zone", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-phone-number-assigned
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-phone-number-assigned",
			Short: "Get Phone Number assigned for a Virtual Line",
			Long:  "Get details on the assigned phone number and extension for the virtual line.\n\nRetrieving virtual line phone number details requires a full or user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/number")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve virtual line settings from this organization.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-directory-search
		var virtualLineId string
		var orgId string
		var enabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-directory-search",
			Short: "Update Directory search for a Virtual Line",
			Long:  "Update the directory search for a designated Virtual Line.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines to Webex Calling users.\n\nUpdating Directory search for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write` and `identity:contacts_rw`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/directorySearch")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Update settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update virtual line settings from this organization.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // list-devices-assigned
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-devices-assigned",
			Short: "Get List of Devices assigned for a Virtual Line",
			Long:  "Retrieve Device details assigned for a virtual line.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines to Webex Calling users.\n\nRetrieving the assigned device detials for a virtual line requires a full or user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/devices")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve virtual line settings from this organization.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // list-networks-handsets
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-networks-handsets",
			Short: "Get List of DECT Networks Handsets for a Virtual Line",
			Long:  "<div><Callout type=\"warning\">Not supported for Webex for Government (FedRAMP)</Callout></div>\n\nRetrieve DECT Network details assigned for a virtual line.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines to Webex Calling users.\n\nRetrieving the assigned device detials for a virtual line requires a full or user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/dectNetworks")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve virtual line settings from this organization.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-caller-id
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-caller-id",
			Short: "Read Caller ID Settings for a Virtual Line",
			Long:  "Retrieve a virtual line's Caller ID settings.\n\nCaller ID settings control how a virtual line's information is displayed when making outgoing calls.\n\nRetrieving the caller ID settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, `dialByFirstName`, and `dialByLastName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/callerId")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter, as the default is the same organization as the token used to access the API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-caller-id
		var virtualLineId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-caller-id",
			Short: "Configure Caller ID Settings for a Virtual Line",
			Long:  "Configure a virtual line's Caller ID settings.\n\nCaller ID settings control how a virtual line's information is displayed when making outgoing calls.\n\nUpdating the caller ID settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, `dialByFirstName`, and `dialByLastName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/callerId")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Update settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter, as the default is the same organization as the token used to access the API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-waiting-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-waiting-settings",
			Short: "Read Call Waiting Settings for a Virtual Line",
			Long:  "Retrieve a virtual line's Call Waiting settings.\n\nWith this feature, a virtual line can place an active call on hold and answer an incoming call.  When enabled, while you are on an active call, a tone alerts you of an incoming call and you can choose to answer or ignore the call.\n\nRetrieving the call waiting settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/callWaiting")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-waiting-settings
		var virtualLineId string
		var orgId string
		var enabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-waiting-settings",
			Short: "Configure Call Waiting Settings for a Virtual Line",
			Long:  "Configure a virtual line's Call Waiting settings.\n\nWith this feature, a virtual line can place an active call on hold and answer an incoming call.  When enabled, while you are on an active call, a tone alerts you of an incoming call and you can choose to answer or ignore the call.\n\nUpdating the call waiting settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/callWaiting")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Update settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-forward
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-forward",
			Short: "Read Call Forwarding Settings for a Virtual Line",
			Long:  "Retrieve a virtual line's Call Forwarding settings.\n\nThree types of call forwarding are supported:\n\n+ Always - forwards all incoming calls to the destination you choose.\n\n+ When busy - forwards all incoming calls to the destination you chose while the phone is in use or the virtual line is busy.\n\n+ When no answer - forwarding only occurs when you are away or not answering your phone.\n\nIn addition, the Business Continuity feature will send calls to a destination of your choice if your phone is not connected to the network for any reason, such as a power outage, failed Internet connection, or wiring problem.\n\nRetrieving the call forwarding settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/callForwarding")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-forward
		var virtualLineId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-forward",
			Short: "Configure Call Forwarding Settings for a Virtual Line",
			Long:  "Configure a virtual line's Call Forwarding settings.\n\nThree types of call forwarding are supported:\n\n+ Always - forwards all incoming calls to the destination you choose.\n\n+ When busy - forwards all incoming calls to the destination you chose while the phone is in use or the virtual line is busy.\n\n+ When no answer - forwarding only occurs when you are away or not answering your phone.\n\nIn addition, the Business Continuity feature will send calls to a destination of your choice if your phone is not connected to the network for any reason, such as a power outage, failed Internet connection, or wiring problem.\n\nUpdating the call forwarding settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/callForwarding")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Update settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-incoming-permission-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-incoming-permission-settings",
			Short: "Read Incoming Permission Settings for a Virtual Line",
			Long:  "Retrieve a virtual line's Incoming Permission settings.\n\nYou can change the incoming calling permissions for a virtual line if you want them to be different from your organization's default.\n\nRetrieving the incoming permission settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/incomingPermission")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-incoming-permission-settings
		var virtualLineId string
		var orgId string
		var useCustomEnabled bool
		var externalTransfer string
		var internalCallsEnabled bool
		var collectCallsEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-incoming-permission-settings",
			Short: "Configure Incoming Permission Settings for a Virtual Line",
			Long:  "Configure a virtual line's Incoming Permission settings.\n\nYou can change the incoming calling permissions for a virtual line if you want them to be different from your organization's default.\n\nUpdating the incoming permission settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/incomingPermission")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("useCustomEnabled", useCustomEnabled, cmd.Flags().Changed("use-custom-enabled"))
					req.BodyString("externalTransfer", externalTransfer)
					req.BodyBool("internalCallsEnabled", internalCallsEnabled, cmd.Flags().Changed("internal-calls-enabled"))
					req.BodyBool("collectCallsEnabled", collectCallsEnabled, cmd.Flags().Changed("collect-calls-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Update settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&useCustomEnabled, "use-custom-enabled", false, "")
		cmd.Flags().StringVar(&externalTransfer, "external-transfer", "", "")
		cmd.Flags().BoolVar(&internalCallsEnabled, "internal-calls-enabled", false, "")
		cmd.Flags().BoolVar(&collectCallsEnabled, "collect-calls-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-outgoing-permissions-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-outgoing-permissions-settings",
			Short: "Retrieve a virtual line's Outgoing Calling Permissions Settings",
			Long:  "Retrieve a virtual line's Outgoing Calling Permissions settings.\n\nOutgoing calling permissions regulate behavior for calls placed to various destinations and default to the local level settings. You can change the outgoing calling permissions for a virtual line if you want them to be different from your organization's default.\n\nRetrieving the outgoing permission settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-outgoing-permissions-settings
		var virtualLineId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-outgoing-permissions-settings",
			Short: "Modify a virtual line's Outgoing Calling Permissions Settings",
			Long:  "Modify a virtual line's Outgoing Calling Permissions settings.\n\nOutgoing calling permissions regulate behavior for calls placed to various destinations and default to the local level settings. You can change the outgoing calling permissions for a virtual line if you want them to be different from your organization's default.\n\nUpdating the outgoing permission settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Update settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-access-codes
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-access-codes",
			Short: "Retrieve Access Codes for a Virtual Line",
			Long: `Retrieve the virtual line's access codes.

Access codes are used to bypass permissions.

This API requires a full, user or read-only administrator auth token with a scope of spark-admin:telephony_config_read`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/accessCodes")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-access-codes
		var virtualLineId string
		var orgId string
		var deleteCodes []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-access-codes",
			Short: "Modify Access Codes for a Virtual Line",
			Long:  "Modify a virtual line's access codes.\n\nAccess codes are used to bypass permissions.\n\nThis API requires a full or user administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/accessCodes")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("deleteCodes", deleteCodes)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		cmd.Flags().StringSliceVar(&deleteCodes, "delete-codes", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // create-access-codes
		var virtualLineId string
		var orgId string
		var code string
		var description string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-access-codes",
			Short: "Create Access Codes for a Virtual Line",
			Long:  "Create a new access codes for the virtual line.\n\nAccess codes are used to bypass permissions.\n\nThis API requires a full or user administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/accessCodes")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("code", code)
					req.BodyString("description", description)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		cmd.Flags().StringVar(&code, "code", "", "")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // delete-access-codes
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-access-codes",
			Short: "Delete Access Codes for a Virtual Line",
			Long:  "Deletes all access codes for the virtual line.\n\nAccess codes are used to bypass permissions.\n\nThis API requires a full or user administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/accessCodes")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-transfer-numbers
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-transfer-numbers",
			Short: "Retrieve Transfer Numbers for a Virtual Line",
			Long: `Retrieve the virtual line's transfer numbers.

When calling a specific call type, this virtual line will be automatically transferred to another number. The virtual line assigned to the Auto Transfer Number can then approve the call and send it through or reject the call type. You can add up to 3 numbers.

This API requires a full, user or read-only administrator auth token with a scope of spark-admin:telephony_config_read`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/autoTransferNumbers")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-transfer-numbers
		var virtualLineId string
		var orgId string
		var useCustomTransferNumbers bool
		var autoTransferNumber1 string
		var autoTransferNumber2 string
		var autoTransferNumber3 string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-transfer-numbers",
			Short: "Modify Transfer Numbers for a Virtual Line",
			Long:  "Modify a virtual line's transfer numbers.\n\nWhen calling a specific call type, this virtual line will be automatically transferred to another number. The virtual line assigned the Auto Transfer Number can then approve the call and send it through or reject the call type. You can add up to 3 numbers.\n\nThis API requires a full or user administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/autoTransferNumbers")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("useCustomTransferNumbers", useCustomTransferNumbers, cmd.Flags().Changed("use-custom-transfer-numbers"))
					req.BodyString("autoTransferNumber1", autoTransferNumber1)
					req.BodyString("autoTransferNumber2", autoTransferNumber2)
					req.BodyString("autoTransferNumber3", autoTransferNumber3)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		cmd.Flags().BoolVar(&useCustomTransferNumbers, "use-custom-transfer-numbers", false, "")
		cmd.Flags().StringVar(&autoTransferNumber1, "auto-transfer-number1", "", "")
		cmd.Flags().StringVar(&autoTransferNumber2, "auto-transfer-number2", "", "")
		cmd.Flags().StringVar(&autoTransferNumber3, "auto-transfer-number3", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-digit-patterns-profile
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-digit-patterns-profile",
			Short: "Retrieve Digit Patterns for a Virtual Profile",
			Long:  "Get list of digit patterns for the virtual profile.\n\nDigit patterns are used to bypass permissions.\n\nRetrieving this list requires a full, user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/digitPatterns")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // create-digit-pattern-profile
		var virtualLineId string
		var orgId string
		var name string
		var pattern string
		var action string
		var transferEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-digit-pattern-profile",
			Short: "Create Digit Pattern for a Virtual Profile",
			Long:  "Create a new digit pattern for a virtual profile.\n\nDigit patterns are used to bypass permissions.\n\nCreating the digit pattern requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/digitPatterns")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("pattern", pattern)
					req.BodyString("action", action)
					req.BodyBool("transferEnabled", transferEnabled, cmd.Flags().Changed("transfer-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&pattern, "pattern", "", "")
		cmd.Flags().StringVar(&action, "action", "", "")
		cmd.Flags().BoolVar(&transferEnabled, "transfer-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-digit-pattern-control-profile
		var virtualLineId string
		var orgId string
		var useCustomDigitPatterns bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-digit-pattern-control-profile",
			Short: "Modify the Digit Pattern Category Control Settings for a Virtual Profile",
			Long:  "Modifies whether this virtual profile uses the specified digit patterns when placing outbound calls or not.\n\nUpdating the digit pattern category control settings requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/digitPatterns")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("useCustomDigitPatterns", useCustomDigitPatterns, cmd.Flags().Changed("use-custom-digit-patterns"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&useCustomDigitPatterns, "use-custom-digit-patterns", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // delete-all-digit-patterns-profile
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-all-digit-patterns-profile",
			Short: "Delete all Digit Patterns for a Virtual Profile",
			Long:  "Delete all digit patterns for a virtual profile.\n\nDigit patterns are used to bypass permissions.\n\nDeleting the digit patterns requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/digitPatterns")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-specified-digit-pattern-profile
		var virtualLineId string
		var digitPatternId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-specified-digit-pattern-profile",
			Short: "Retrieve Specified Digit Pattern Details for a Virtual Profile",
			Long:  "Get the specified digit pattern for the virtual profile\u200b.\n\nDigit patterns are used to bypass permissions.\n\nRetrieving the digit pattern details requires a full, user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/digitPatterns/{digitPatternId}")
				req.PathParam("virtualLineId", virtualLineId)
				req.PathParam("digitPatternId", digitPatternId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&digitPatternId, "digit-pattern-id", "", "Unique identifier for the digit pattern.")
		cmd.MarkFlagRequired("digit-pattern-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-digit-pattern-profile
		var virtualLineId string
		var digitPatternId string
		var orgId string
		var name string
		var pattern string
		var action string
		var transferEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-digit-pattern-profile",
			Short: "Modify a Digit Pattern for a Virtual Profile",
			Long:  "Modify a digit patterns for a virtual profile.\n\nDigit patterns are used to bypass permissions.\n\nUpdating the digit pattern requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/digitPatterns/{digitPatternId}")
				req.PathParam("virtualLineId", virtualLineId)
				req.PathParam("digitPatternId", digitPatternId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("pattern", pattern)
					req.BodyString("action", action)
					req.BodyBool("transferEnabled", transferEnabled, cmd.Flags().Changed("transfer-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&digitPatternId, "digit-pattern-id", "", "Unique identifier for the digit pattern.")
		cmd.MarkFlagRequired("digit-pattern-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&pattern, "pattern", "", "")
		cmd.Flags().StringVar(&action, "action", "", "")
		cmd.Flags().BoolVar(&transferEnabled, "transfer-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // delete-digit-pattern-profile
		var virtualLineId string
		var digitPatternId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-digit-pattern-profile",
			Short: "Delete a Digit Pattern for a Virtual Profile",
			Long:  "Delete a digit pattern for a virtual profile.\n\nDigit patterns are used to bypass permissions.\n\nDeleting the digit pattern requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/virtualLines/{virtualLineId}/outgoingPermission/digitPatterns/{digitPatternId}")
				req.PathParam("virtualLineId", virtualLineId)
				req.PathParam("digitPatternId", digitPatternId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&digitPatternId, "digit-pattern-id", "", "Unique identifier for the digit pattern.")
		cmd.MarkFlagRequired("digit-pattern-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-intercept-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-intercept-settings",
			Short: "Read Call Intercept Settings for a Virtual Line",
			Long:  "Retrieves Virtual Line's Call Intercept settings.\n\nThe intercept feature gracefully takes a virtual line's phone out of service, while providing callers with informative announcements and alternative routing options. Depending on the service configuration, none, some, or all incoming calls to the specified virtual line are intercepted. Also depending on the service configuration, outgoing calls are intercepted or rerouted to another location.\n\nRetrieving the intercept settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/intercept")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-intercept-settings
		var virtualLineId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-intercept-settings",
			Short: "Configure Call Intercept Settings for a Virtual Line",
			Long:  "Configures a virtual line's Call Intercept settings.\n\nThe intercept feature gracefully takes a virtual line's phone out of service, while providing callers with informative announcements and alternative routing options. Depending on the service configuration, none, some, or all incoming calls to the specified virtual line are intercepted. Also depending on the service configuration, outgoing calls are intercepted or rerouted to another location.\n\nUpdating the intercept settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/intercept")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Update settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-intercept-greeting
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "update-intercept-greeting",
			Short: "Configure Call Intercept Greeting for a Virtual Line",
			Long:  "Configure a virtual line's Call Intercept Greeting by uploading a Waveform Audio File Format, `.wav`, encoded audio file.\n\nYour request will need to be a `multipart/form-data` request rather than JSON, using the `audio/wav` Content-Type.\n\nUploading the intercept greeting announcement for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**WARNING:** This API is not callable using the developer portal web interface due to the lack of support for multipart POST. This API can be utilized using other tools that support multipart POST, such as Postman.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualLines/{virtualLineId}/intercept/actions/announcementUpload/invoke")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Update settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-agent-list-available-caller-ids
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-agent-list-available-caller-ids",
			Short: "Retrieve Agent's List of Available Caller IDs",
			Long:  "Get the list of call queues and hunt groups available for caller ID use by this virtual line as an agent.\n\nThis API requires a full, user, or read-only administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/agent/availableCallerIds")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the Virtual Line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the Virtual Line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-agent-caller-id
		var virtualLineId string
		cmd := &cobra.Command{
			Use:   "get-agent-caller-id",
			Short: "Retrieve Agent's Caller ID Information",
			Long:  "Retrieve the Agent's Caller ID Information.\n\nEach agent will be able to set their outgoing Caller ID as either the Call Queue's Caller ID, Hunt Group's Caller ID or their own configured Caller ID.\n\nThis API requires a full admin or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/agent/callerId")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the Virtual Line.")
		cmd.MarkFlagRequired("virtual-line-id")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-agent-caller-id
		var virtualLineId string
		var selectedCallerId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-agent-caller-id",
			Short: "Modify Agent's Caller ID Information",
			Long:  "Modify Agent's Caller ID Information.\n\nEach Agent is able to set their outgoing Caller ID as either the designated Call Queue's Caller ID or the Hunt Group's Caller ID or their own configured Caller ID.\nThis API requires a full or user administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/agent/callerId")
				req.PathParam("virtualLineId", virtualLineId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("selectedCallerId", selectedCallerId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the Virtual Line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&selectedCallerId, "selected-caller-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-voicemail
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-voicemail",
			Short: "Read Voicemail Settings for a Virtual Line",
			Long:  "Retrieve a virtual line's voicemail settings.\n\nThe voicemail feature transfers callers to voicemail based on your settings. You can then retrieve voice messages via voicemail.\n\nOptionally, notifications can be sent to a mobile phone via text or email. These notifications will not include the voicemail files.\n\nRetrieving the voicemail settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/voicemail")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-voicemail
		var virtualLineId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-voicemail",
			Short: "Configure Voicemail Settings for a Virtual Line",
			Long:  "Configure a virtual line's voicemail settings.\n\nThe voicemail feature transfers callers to voicemail based on your settings. You can then retrieve voice messages via voicemail.\n\nOptionally, notifications can be sent to a mobile phone via text or email. These notifications will not include the voicemail files.\n\nUpdating the voicemail settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/voicemail")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-busy-voicemail-greeting
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "update-busy-voicemail-greeting",
			Short: "Configure Busy Voicemail Greeting for a Virtual Line",
			Long:  "Configure a virtual line's Busy Voicemail Greeting by uploading a Waveform Audio File Format, `.wav`, encoded audio file.\n\nYour request will need to be a `multipart/form-data` request rather than JSON, using the `audio/wav` Content-Type.\n\nUploading the voicemail busy greeting announcement for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**WARNING:** This API is not callable using the developer portal web interface due to the lack of support for multipart POST. This API can be utilized using other tools that support multipart POST, such as Postman.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualLines/{virtualLineId}/voicemail/actions/uploadBusyGreeting/invoke")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-no-answer-voicemail-greeting
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "update-no-answer-voicemail-greeting",
			Short: "Configure No Answer Voicemail Greeting for a Virtual Line",
			Long:  "Configure a virtual line's No Answer Voicemail Greeting by uploading a Waveform Audio File Format, `.wav`, encoded audio file.\n\nYour request will need to be a `multipart/form-data` request rather than JSON, using the `audio/wav` Content-Type.\n\nUploading the voicemail no answer greeting announcement for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**WARNING:** This API is not callable using the developer portal web interface due to the lack of support for multipart POST. This API can be utilized using other tools that support multipart POST, such as Postman.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualLines/{virtualLineId}/voicemail/actions/uploadNoAnswerGreeting/invoke")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // reset-voicemail-pin
		var virtualLineId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "reset-voicemail-pin",
			Short: "Reset Voicemail PIN for a Virtual Line",
			Long:  "Reset a voicemail PIN for a virtual line.\n\nThe voicemail feature transfers callers to voicemail based on your settings. You can then retrieve voice messages via Voicemail.  A voicemail PIN is used to retrieve your voicemail messages.\n\nUpdating the voicemail pin for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: This API is expected to have an empty request body and Content-Type header should be set to `application/json`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/virtualLines/{virtualLineId}/voicemail/actions/resetPin/invoke")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-voicemail-passcode
		var virtualLineId string
		var orgId string
		var passcode string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-voicemail-passcode",
			Short: "Modify a virtual line's voicemail passcode",
			Long:  "Modify a virtual line's voicemail passcode.\n\nModifying a virtual line's voicemail passcode requires a full administrator, user administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/voicemail/passcode")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("passcode", passcode)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Modify voicemail passcode for this virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Modify voicemail passcode for a virtual line in this organization.")
		cmd.Flags().StringVar(&passcode, "passcode", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-music-hold-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-music-hold-settings",
			Short: "Retrieve Music On Hold Settings for a Virtual Line",
			Long:  "Retrieve the virtual line's music on hold settings.\n\nMusic on hold is played when a caller is put on hold, or the call is parked.\n\nRetrieving the music on hold settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/musicOnHold")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-music-hold-settings
		var virtualLineId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-music-hold-settings",
			Short: "Configure Music On Hold Settings for a Virtual Line",
			Long:  "Configure a virtual line's music on hold settings.\n\nMusic on hold is played when a caller is put on hold, or the call is parked.\n\nTo configure music on hold settings for a virtual line, music on hold setting must be enabled for this location.\n\nUpdating the music on hold settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/musicOnHold")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-push-to-talk-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-push-to-talk-settings",
			Short: "Read Push-to-Talk Settings for a Virtual Line",
			Long:  "Retrieve a virtual line's Push-to-Talk settings.\n\nPush-to-Talk allows the use of desk phones as either a one-way or two-way intercom that connects people in different parts of your organization.\n\nRetrieving the Push-to-Talk settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/pushToTalk")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-push-to-talk-settings
		var virtualLineId string
		var orgId string
		var allowAutoAnswer bool
		var connectionType string
		var accessType string
		var members []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-push-to-talk-settings",
			Short: "Configure Push-to-Talk Settings for a Virtual Line",
			Long:  "Configure a virtual line's Push-to-Talk settings.\n\nPush-to-Talk allows the use of desk phones as either a one-way or two-way intercom that connects people in different parts of your organization.\n\nUpdating the Push-to-Talk settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/pushToTalk")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("allowAutoAnswer", allowAutoAnswer, cmd.Flags().Changed("allow-auto-answer"))
					req.BodyString("connectionType", connectionType)
					req.BodyString("accessType", accessType)
					req.BodyStringSlice("members", members)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&allowAutoAnswer, "allow-auto-answer", false, "")
		cmd.Flags().StringVar(&connectionType, "connection-type", "", "")
		cmd.Flags().StringVar(&accessType, "access-type", "", "")
		cmd.Flags().StringSliceVar(&members, "members", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-bridge-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-bridge-settings",
			Short: "Read Call Bridge Settings for a Virtual Line",
			Long:  "Retrieve a virtual line's call bridge settings.\n\nRetrieving the call bridge settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/callBridge")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-bridge-settings
		var virtualLineId string
		var orgId string
		var warningToneEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-bridge-settings",
			Short: "Configure Call Bridge Settings for a Virtual Line",
			Long:  "Configure a virtual line's call bridge settings.\n\nUpdating the call bridge settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/callBridge")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("warningToneEnabled", warningToneEnabled, cmd.Flags().Changed("warning-tone-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&warningToneEnabled, "warning-tone-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-barge-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-barge-settings",
			Short: "Read Barge In Settings for a Virtual Line",
			Long:  "Retrieve a virtual line's barge in settings.\n\nThe Barge In feature enables you to use a Feature Access Code (FAC) to answer a call that was directed to another subscriber, or barge-in on the call if it was already answered. Barge In can be used across locations.\n\nRetrieving the barge in settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/bargeIn")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-barge-settings
		var virtualLineId string
		var orgId string
		var enabled bool
		var toneEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-barge-settings",
			Short: "Configure Barge In Settings for a Virtual Line",
			Long:  "Configure a virtual line's barge in settings.\n\nThe Barge In feature enables you to use a Feature Access Code (FAC) to answer a call that was directed to another subscriber, or barge-in on the call if it was already answered. Barge In can be used across locations.\n\nUpdating the barge in settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/bargeIn")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyBool("toneEnabled", toneEnabled, cmd.Flags().Changed("tone-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().BoolVar(&toneEnabled, "tone-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-privacy-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-privacy-settings",
			Short: "Get a Virtual Line's Privacy Settings",
			Long:  "Get a virtual line's privacy settings for the specified virtual line ID.\n\nThe privacy feature enables the virtual line's line to be monitored by others and determine if they can be reached by Auto Attendant services.\n\nRetrieving the privacy settings for a virtual line requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/privacy")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-privacy-settings
		var virtualLineId string
		var orgId string
		var aaExtensionDialingEnabled bool
		var aaNamingDialingEnabled bool
		var enablePhoneStatusDirectoryPrivacy bool
		var enablePhoneStatusPickupBargeInPrivacy bool
		var monitoringAgents []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-privacy-settings",
			Short: "Configure a Virtual Line's Privacy Settings",
			Long:  "Configure a virtual line's privacy settings for the specified virtual line ID.\n\nThe privacy feature enables the virtual line's line to be monitored by others and determine if they can be reached by Auto Attendant services.\n\nUpdating the privacy settings for a virtual line requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/privacy")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("aaExtensionDialingEnabled", aaExtensionDialingEnabled, cmd.Flags().Changed("aa-extension-dialing-enabled"))
					req.BodyBool("aaNamingDialingEnabled", aaNamingDialingEnabled, cmd.Flags().Changed("aa-naming-dialing-enabled"))
					req.BodyBool("enablePhoneStatusDirectoryPrivacy", enablePhoneStatusDirectoryPrivacy, cmd.Flags().Changed("enable-phone-status-directory-privacy"))
					req.BodyBool("enablePhoneStatusPickupBargeInPrivacy", enablePhoneStatusPickupBargeInPrivacy, cmd.Flags().Changed("enable-phone-status-pickup-barge-in-privacy"))
					req.BodyStringSlice("monitoringAgents", monitoringAgents)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Retrieve settings for a virtual line with the matching ID.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the virtual line resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&aaExtensionDialingEnabled, "aa-extension-dialing-enabled", false, "")
		cmd.Flags().BoolVar(&aaNamingDialingEnabled, "aa-naming-dialing-enabled", false, "")
		cmd.Flags().BoolVar(&enablePhoneStatusDirectoryPrivacy, "enable-phone-status-directory-privacy", false, "")
		cmd.Flags().BoolVar(&enablePhoneStatusPickupBargeInPrivacy, "enable-phone-status-pickup-barge-in-privacy", false, "")
		cmd.Flags().StringSliceVar(&monitoringAgents, "monitoring-agents", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-fax-available-numbers
		var virtualLineId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-fax-available-numbers",
			Short: "Get Virtual Line Fax Message Available Phone Numbers",
			Long:  "List standard numbers that are available to be assigned as a virtual line's FAX message number.\nThese numbers are associated with the location of the virtual line specified in the request URL, can be active or inactive, and are unassigned.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/faxMessage/availableNumbers")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-forward-available-numbers
		var virtualLineId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		var ownerName string
		var extension string
		cmd := &cobra.Command{
			Use:   "get-forward-available-numbers",
			Short: "Get Virtual Line Call Forward Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as a virtual line's call forward number.\nThese numbers are associated with the location of the virtual line specified in the request URL, can be active or inactive, and are assigned to an owning entity.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/callForwarding/availableNumbers")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("ownerName", ownerName)
				req.QueryParam("extension", extension)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		cmd.Flags().StringVar(&ownerName, "owner-name", "", "Return the list of phone numbers that are owned by the given `ownerName`. Maximum length is 255.")
		cmd.Flags().StringVar(&extension, "extension", "", "Returns the list of PSTN phone numbers with the given `extension`.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-available-numbers
		var orgId string
		var locationId string
		var max string
		var start string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-available-numbers",
			Short: "Get Virtual Line Available Phone Numbers",
			Long:  "List standard numbers that are available to be assigned as a virtual line's phone number.\nBy default, this API returns unassigned numbers from all locations. To select the suitable number for assignment, ensure the virtual line's location ID is provided as the `locationId` request parameter.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/availableNumbers")
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationId", locationId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the list of phone numbers for this location within the given organization. The maximum length is 36.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-ecbn-available-numbers
		var virtualLineId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		var ownerName string
		cmd := &cobra.Command{
			Use:   "get-ecbn-available-numbers",
			Short: "Get Virtual Line ECBN Available Phone Numbers",
			Long:  "List standard numbers that can be assigned as a virtual line's call forward number.\nThese numbers are associated with the location of the virtual line specified in the request URL, can be active or inactive, and are assigned to an owning entity.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/emergencyCallbackNumber/availableNumbers")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("ownerName", ownerName)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		cmd.Flags().StringVar(&ownerName, "owner-name", "", "Return the list of phone numbers that are owned by the given `ownerName`. Maximum length is 255.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-intercept-available-numbers
		var virtualLineId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		var ownerName string
		var extension string
		cmd := &cobra.Command{
			Use:   "get-intercept-available-numbers",
			Short: "Get Virtual Line Call Intercept Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as a virtual line's call intercept number.\nThese numbers are associated with the location of the virtual line specified in the request URL, can be active or inactive, and are assigned to an owning entity.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/callIntercept/availableNumbers")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("ownerName", ownerName)
				req.QueryParam("extension", extension)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		cmd.Flags().StringVar(&ownerName, "owner-name", "", "Return the list of phone numbers that are owned by the given `ownerName`. Maximum length is 255.")
		cmd.Flags().StringVar(&extension, "extension", "", "Returns the list of PSTN phone numbers with the given `extension`.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // get-donotdisturb-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-donotdisturb-settings",
			Short: "Retrieve DoNotDisturb Settings for a Virtual Line",
			Long:  "Retrieve DoNotDisturb Settings for a Virtual Line.\n\nSilence incoming calls with the Do Not Disturb feature.\nWhen enabled, callers hear the busy signal.\n\nThis API requires a full, read-only or location administrator auth token with a scope of `telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/doNotDisturb")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization within which the virtual line resides.")
		virtualLineCallCmd.AddCommand(cmd)
	}

	{ // update-donotdisturb-settings
		var virtualLineId string
		var orgId string
		var enabled bool
		var ringSplashEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-donotdisturb-settings",
			Short: "Modify DoNotDisturb Settings for a Virtual Line",
			Long:  "Modify DoNotDisturb Settings for a Virtual Line.\n\nSilence incoming calls with the Do Not Disturb feature.\nWhen enabled, callers hear the busy signal.\n\nThis API requires a full, user or location administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/doNotDisturb")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyBool("ringSplashEnabled", ringSplashEnabled, cmd.Flags().Changed("ring-splash-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization within which the virtual line resides.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().BoolVar(&ringSplashEnabled, "ring-splash-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		virtualLineCallCmd.AddCommand(cmd)
	}

}
