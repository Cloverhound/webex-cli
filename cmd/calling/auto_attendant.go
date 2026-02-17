package calling

import (
	"fmt"

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

var autoAttendantCmd = &cobra.Command{
	Use:   "auto-attendant",
	Short: "AutoAttendant commands",
}

func init() {
	cmd.CallingCmd.AddCommand(autoAttendantCmd)

	{ // list
		var orgId string
		var locationId string
		var max string
		var start string
		var name string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "Read the List of Auto Attendants",
			Long:  "List all Auto Attendants for the organization.\n\nAuto attendants play customized prompts and provide callers with menu options for routing their calls through your system.\n\nRetrieving this list requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/autoAttendants")
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationId", locationId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List auto attendants for this organization.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the list of auto attendants for this location.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&name, "name", "", "Only return auto attendants with the matching name.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Only return auto attendants with the matching phone number.")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // get
		var locationId string
		var autoAttendantId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Details for an Auto Attendant",
			Long:  "Retrieve an Auto Attendant details.\n\nAuto attendants play customized prompts and provide callers with menu options for routing their calls through your system.\n\nRetrieving an auto attendant details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, and `dialByName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve an auto attendant details in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Retrieve the auto attendant with the matching ID.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve auto attendant details from this organization.")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // update
		var locationId string
		var autoAttendantId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update an Auto Attendant",
			Long:  "Update the designated Auto Attendant.\n\nAuto attendants play customized prompts and provide callers with menu options for routing their calls through your system.\n\nUpdating an auto attendant requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, and `dialByName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this auto attendant exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Update an auto attendant with the matching ID.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update an auto attendant from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // delete
		var locationId string
		var autoAttendantId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete an Auto Attendant",
			Long:  "Delete the designated Auto Attendant.\n\nAuto attendants play customized prompts and provide callers with menu options for routing their calls through your system.\n\nDeleting an auto attendant requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location from which to delete an auto attendant.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Delete the auto attendant with the matching ID.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the auto attendant from this organization.")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // create
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create an Auto Attendant",
			Long:  "Create new Auto Attendant for the given location.\n\nAuto attendants play customized prompts and provide callers with menu options for routing their calls through your system.\n\nCreating an auto attendant requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, and `dialByName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/autoAttendants")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Create the auto attendant for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create the auto attendant for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // get-call-forward
		var locationId string
		var autoAttendantId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-call-forward",
			Short: "Get Call Forwarding Settings for an Auto Attendant",
			Long:  "Retrieve Call Forwarding settings for the designated Auto Attendant including the list of call forwarding rules.\n\nThe call forwarding feature allows you to direct all incoming calls based on specific criteria that you define.\nBelow are the available options for configuring your call forwarding:\n1. Always forward calls to a designated number.\n2. Forward calls to a designated number based on certain criteria.\n3. Forward calls using different modes.\n\nRetrieving call forwarding settings for an auto attendant requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}/callForwarding")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this auto attendant exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Retrieve the call forwarding settings for this auto attendant.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve auto attendant forwarding settings from this organization.")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // update-call-forward
		var locationId string
		var autoAttendantId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-call-forward",
			Short: "Update Call Forwarding Settings for an Auto Attendant",
			Long:  "Update Call Forwarding settings for the designated Auto Attendant.\n\nThe call forwarding feature allows you to direct all incoming calls based on specific criteria that you define.\nBelow are the available options for configuring your call forwarding:\n1. Always forward calls to a designated number.\n2. Forward calls to a designated number based on certain criteria.\n3. Forward calls using different modes.\n\nUpdating call forwarding settings for an auto attendant requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}/callForwarding")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this auto attendant exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Update call forwarding settings for this auto attendant.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update auto attendant forwarding settings from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // create-selective-forward-rule
		var locationId string
		var autoAttendantId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-selective-forward-rule",
			Short: "Create a Selective Call Forwarding Rule for an Auto Attendant",
			Long:  "Create a Selective Call Forwarding Rule for the designated Auto Attendant.\n\nA selective call forwarding rule for an auto attendant allows calls to be forwarded or not forwarded to the designated number, based on the defined criteria.\n\nNote that the list of existing call forward rules is available in the auto attendant's call forwarding settings.\n\nCreating a selective call forwarding rule for an auto attendant requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The Call Forwarding Rule ID will change upon modification of the Call Forwarding Rule name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}/callForwarding/selectiveRules")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which the auto attendant exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Create the rule for this auto attendant.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create the auto attendant rule for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // get-selective-forward-rule
		var locationId string
		var autoAttendantId string
		var ruleId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-selective-forward-rule",
			Short: "Get Selective Call Forwarding Rule for an Auto Attendant",
			Long: `Retrieve a Selective Call Forwarding Rule's settings for the designated Auto Attendant.

A selective call forwarding rule for an auto attendant allows calls to be forwarded or not forwarded to the designated number, based on the defined criteria.

Note that the list of existing call forward rules is available in the auto attendant's call forwarding settings.

Retrieving a selective call forwarding rule's settings for an auto attendant requires a full or read-only administrator or location administrator

**NOTE**: The Call Forwarding Rule ID will change upon modification of the Call Forwarding Rule name.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}/callForwarding/selectiveRules/{ruleId}")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
				req.PathParam("ruleId", ruleId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this auto attendant exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Retrieve settings for a rule for this auto attendant.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&ruleId, "rule-id", "", "Auto attendant rule you are retrieving settings for.")
		cmd.MarkFlagRequired("rule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve auto attendant rule settings for this organization.")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // update-selective-forward-rule
		var locationId string
		var autoAttendantId string
		var ruleId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-selective-forward-rule",
			Short: "Update Selective Call Forwarding Rule for an Auto Attendant",
			Long:  "Update a Selective Call Forwarding Rule's settings for the designated Auto Attendant.\n\nA selective call forwarding rule for an auto attendant allows calls to be forwarded or not forwarded to the designated number, based on the defined criteria.\n\nNote that the list of existing call forward rules is available in the auto attendant's call forwarding settings.\n\nUpdating a selective call forwarding rule's settings for an auto attendant requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The Call Forwarding Rule ID will change upon modification of the Call Forwarding Rule name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}/callForwarding/selectiveRules/{ruleId}")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
				req.PathParam("ruleId", ruleId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this auto attendant exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Update settings for a rule for this auto attendant.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&ruleId, "rule-id", "", "Auto attendant rule you are updating settings for.")
		cmd.MarkFlagRequired("rule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update auto attendant rule settings for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // delete-selective-forward-rule
		var locationId string
		var autoAttendantId string
		var ruleId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-selective-forward-rule",
			Short: "Delete a Selective Call Forwarding Rule for an Auto Attendant",
			Long:  "Delete a Selective Call Forwarding Rule for the designated Auto Attendant.\n\nA selective call forwarding rule for an auto attendant allows calls to be forwarded or not forwarded to the designated number, based on the defined criteria.\n\nNote that the list of existing call forward rules is available in the auto attendant's call forwarding settings.\n\nDeleting a selective call forwarding rule for an auto attendant requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The Call Forwarding Rule ID will change upon modification of the Call Forwarding Rule name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}/callForwarding/selectiveRules/{ruleId}")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
				req.PathParam("ruleId", ruleId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this auto attendant exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Delete the rule for this auto attendant.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&ruleId, "rule-id", "", "Auto attendant rule you are deleting.")
		cmd.MarkFlagRequired("rule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete auto attendant rule from this organization.")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // get-primary-numbers
		var locationId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-primary-numbers",
			Short: "Get Auto Attendant Primary Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as the auto attendant's primary phone number.\nThese numbers are associated with the location specified in the request URL, can be active or inactive, and are unassigned.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/autoAttendants/availableNumbers")
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
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // get-alternate-numbers
		var locationId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-alternate-numbers",
			Short: "Get Auto Attendant Alternate Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as the auto attendant's alternate number.\nThese numbers are associated with the location specified in the request URL, can be active or inactive, and are unassigned.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/autoAttendants/alternate/availableNumbers")
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
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // get-forward-available-numbers
		var locationId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		var ownerName string
		var extension string
		cmd := &cobra.Command{
			Use:   "get-forward-available-numbers",
			Short: "Get Auto Attendant Call Forward Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as the auto attendant's call forward number.\nThese numbers are associated with the location specified in the request URL, can be active or inactive, and are assigned to an owning entity.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/autoAttendants/callForwarding/availableNumbers")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the list of phone numbers for this location within the given organization. The maximum length is 36.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		cmd.Flags().StringVar(&ownerName, "owner-name", "", "Return the list of phone numbers that are owned by the given `ownerName`. Maximum length is 255.")
		cmd.Flags().StringVar(&extension, "extension", "", "Returns the list of PSTN phone numbers with the given `extension`.")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // switch-mode-call-forward
		var locationId string
		var autoAttendantId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "switch-mode-call-forward",
			Short: "Switch Mode for Call Forwarding Settings for an Auto Attendant",
			Long:  "Switches the current operating mode of the `Auto Attendant` to the mode as per normal operations.\n\nSwitching operating mode for an `auto attendant` requires a full, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}/callForwarding/actions/switchMode/invoke")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "`Location` in which this `auto attendant` exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Switch operating mode to normal operations for this `auto attendant`.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Switch operating mode as per normal operations for the `auto attendant` from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // delete-announcement-file
		var locationId string
		var autoAttendantId string
		var fileName string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-announcement-file",
			Short: "Delete a Auto Attendant Announcement File",
			Long:  "Delete an announcement file for the designated auto attendant.\n\nAuto Attendant announcement files contain messages and music that callers hear while waiting in the queue. A auto attendant can be configured to play whatever subset of these announcement files is desired.\n\nDeleting an announcement file for a auto attendant requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}/announcements/{fileName}")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
				req.PathParam("fileName", fileName)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Delete an announcement for a auto attendant in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Delete an announcement for the auto attendant with this identifier.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&fileName, "file-name", "", "")
		cmd.MarkFlagRequired("file-name")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete auto attendant announcement from this organization.")
		autoAttendantCmd.AddCommand(cmd)
	}

	{ // list-announcement-files
		var locationId string
		var autoAttendantId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-announcement-files",
			Short: "Read the List of Auto Attendant Announcement Files",
			Long:  "List file info for all auto attendant announcement files associated with this auto attendant.\n\nAuto attendant announcement files contain messages and music that callers hear while waiting in the queue. A auto attendant can be configured to play whatever subset of these announcement files is desired.\n\nRetrieving this list of files requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.\n\nNote that uploading of announcement files via API is not currently supported, but is available via Webex Control Hub.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/autoAttendants/{autoAttendantId}/announcements")
				req.PathParam("locationId", locationId)
				req.PathParam("autoAttendantId", autoAttendantId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this auto attendant exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&autoAttendantId, "auto-attendant-id", "", "Retrieve announcement files for the auto attendant with this identifier.")
		cmd.MarkFlagRequired("auto-attendant-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve announcement files for a auto attendant from this organization.")
		autoAttendantCmd.AddCommand(cmd)
	}

}
