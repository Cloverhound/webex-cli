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

var singleNumberReachCmd = &cobra.Command{
	Use:   "single-number-reach",
	Short: "SingleNumberReach commands",
}

func init() {
	cmd.CallingCmd.AddCommand(singleNumberReachCmd)

	{ // get-primary
		var locationId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-primary",
			Short: "Get Single Number Reach Primary Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as the single number reach's primary phone number.\nThese numbers are associated with the location specified in the request URL, can be active or inactive, and are unassigned.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/singleNumberReach/availableNumbers")
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
		singleNumberReachCmd.AddCommand(cmd)
	}

	{ // create-person
		var personId string
		var phoneNumber string
		var enabled bool
		var name string
		var doNotForwardCallsEnabled bool
		var answerConfirmationEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-person",
			Short: "Create Single Number Reach For a Person",
			Long:  "Create a single number reach for a person in an organization.\n\nSingle number reach allows you to setup your work calls ring any phone number. This means that your office phone, mobile phone, or any other designated devices can ring at the same time, ensuring you don't miss important calls.\n\nCreating a single number reach for a person requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/people/{personId}/singleNumberReach/numbers")
				req.PathParam("personId", personId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("phoneNumber", phoneNumber)
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyString("name", name)
					req.BodyBool("doNotForwardCallsEnabled", doNotForwardCallsEnabled, cmd.Flags().Changed("do-not-forward-calls-enabled"))
					req.BodyBool("answerConfirmationEnabled", answerConfirmationEnabled, cmd.Flags().Changed("answer-confirmation-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().BoolVar(&doNotForwardCallsEnabled, "do-not-forward-calls-enabled", false, "")
		cmd.Flags().BoolVar(&answerConfirmationEnabled, "answer-confirmation-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		singleNumberReachCmd.AddCommand(cmd)
	}

	{ // get-settings-person
		var personId string
		cmd := &cobra.Command{
			Use:   "get-settings-person",
			Short: "Get Single Number Reach Settings For A Person",
			Long:  "Retrieve Single Number Reach settings for the given person.\n\nSingle number reach allows you to setup your work calls ring any phone number. This means that your office phone, mobile phone, or any other designated devices can ring at the same time, ensuring you don't miss important calls.\n\nRetrieving Single number reach settings requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/singleNumberReach")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		singleNumberReachCmd.AddCommand(cmd)
	}

	{ // update-settings-person
		var personId string
		var alertAllNumbersForClickToDialCallsEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-settings-person",
			Short: "Update Single Number Reach Settings For A Person",
			Long:  "Update Single number reach settings for a person.\n\nSingle number reach allows you to setup your work calls ring any phone number. This means that your office phone, mobile phone, or any other designated devices can ring at the same time, ensuring you don't miss important calls.\n\nUpdating a single number reach settings for a person requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/singleNumberReach")
				req.PathParam("personId", personId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("alertAllNumbersForClickToDialCallsEnabled", alertAllNumbersForClickToDialCallsEnabled, cmd.Flags().Changed("alert-all-numbers-for-click-to-dial-calls-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().BoolVar(&alertAllNumbersForClickToDialCallsEnabled, "alert-all-numbers-for-click-to-dial-calls-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		singleNumberReachCmd.AddCommand(cmd)
	}

	{ // update-settings
		var personId string
		var id string
		var phoneNumber string
		var enabled bool
		var name string
		var doNotForwardCallsEnabled bool
		var answerConfirmationEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-settings",
			Short: "Update Single Number Reach Settings For A Number",
			Long:  "Update Single number reach settings for a number.\n\nSingle number reach allows you to setup your work calls ring any phone number. This means that your office phone, mobile phone, or any other designated devices can ring at the same time, ensuring you don't miss important calls.\n\nThe response returns an ID that can change if the phoneNumber is modified, as the ID contains base64 encoded phone number data.\n\nUpdating a single number reach settings for a number requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/singleNumberReach/numbers/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("phoneNumber", phoneNumber)
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyString("name", name)
					req.BodyBool("doNotForwardCallsEnabled", doNotForwardCallsEnabled, cmd.Flags().Changed("do-not-forward-calls-enabled"))
					req.BodyBool("answerConfirmationEnabled", answerConfirmationEnabled, cmd.Flags().Changed("answer-confirmation-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Unique identifier for single number reach.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().BoolVar(&doNotForwardCallsEnabled, "do-not-forward-calls-enabled", false, "")
		cmd.Flags().BoolVar(&answerConfirmationEnabled, "answer-confirmation-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		singleNumberReachCmd.AddCommand(cmd)
	}

	{ // delete
		var personId string
		var id string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete A Single Number Reach Number",
			Long:  "Delete Single number reach number for a person.\n\nSingle number reach allows you to setup your work calls ring any phone number. This means that your office phone, mobile phone, or any other designated devices can ring at the same time, ensuring you don't miss important calls.\n\nDeleting a Single number reach number requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/people/{personId}/singleNumberReach/numbers/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Unique identifier for single number reach.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		singleNumberReachCmd.AddCommand(cmd)
	}

}
