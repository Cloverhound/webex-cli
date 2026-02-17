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

var numbersCmd = &cobra.Command{
	Use:   "numbers",
	Short: "Numbers commands",
}

func init() {
	cmd.CallingCmd.AddCommand(numbersCmd)

	{ // create-phone-location
		var locationId string
		var orgId string
		var phoneNumbers []string
		var numberType string
		var numberUsageType string
		var state string
		var subscriptionId string
		var carrierId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-phone-location",
			Short: "Add Phone Numbers to a location",
			Long:  "Adds a specified set of phone numbers to a location for an organization. Phone numbers must follow the E.164 format.\n\nEach location has a set of phone numbers that can be assigned to people, workspaces, or features. Active phone numbers are in service.\n\nAdding a phone number to a location requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\nPhone numbers included in the request that already exist in the location will simply be ignored.\n\n<br/>\n\n<div><Callout type=\"warning\">This API is only supported for adding DID and Toll-free numbers to non-integrated PSTN connection types such as Local Gateway (LGW) and Non-integrated CPP. It should never be used for locations with integrated PSTN connection types like Cisco Calling Plans or Integrated CCP because backend data issues may occur.\n</Callout></div>\n<div><Callout type=\"warning\">Mobile numbers can be added to any location that has PSTN connection setup. Only 20 mobile numbers can be added per request.\n</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/numbers")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("phoneNumbers", phoneNumbers)
					req.BodyString("numberType", numberType)
					req.BodyString("numberUsageType", numberUsageType)
					req.BodyString("state", state)
					req.BodyString("subscriptionId", subscriptionId)
					req.BodyString("carrierId", carrierId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "LocationId to which numbers should be added.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization of the Route Group.")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().StringVar(&numberType, "number-type", "", "")
		cmd.Flags().StringVar(&numberUsageType, "number-usage-type", "", "")
		cmd.Flags().StringVar(&state, "state", "", "")
		cmd.Flags().StringVar(&subscriptionId, "subscription-id", "", "")
		cmd.Flags().StringVar(&carrierId, "carrier-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		numbersCmd.AddCommand(cmd)
	}

	{ // manage-state-location
		var locationId string
		var orgId string
		var phoneNumbers []string
		var action string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "manage-state-location",
			Short: "Manage Number State in a location",
			Long:  "Activate or deactivate the specified set of phone numbers in a location for an organization.\n\nEach location has a set of phone numbers that can be assigned to people, workspaces, or features. Phone numbers must follow the E.164 format.\n\nActive phone numbers are in service. A mobile number is activated when assigned to a user. This API will not activate or deactivate mobile numbers.\n\nManaging phone number state in a location requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n<br/>\n\n<div><Callout type=\"warning\">This API is only supported for non-integrated PSTN connection types of Local Gateway (LGW) and Non-integrated CCP.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/numbers")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("phoneNumbers", phoneNumbers)
					req.BodyString("action", action)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "`LocationId` to which numbers should be added.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization of the Route Group.")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().StringVar(&action, "action", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		numbersCmd.AddCommand(cmd)
	}

	{ // delete-phone-location
		var locationId string
		var orgId string
		var phoneNumbers []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "delete-phone-location",
			Short: "Remove phone numbers from a location",
			Long:  "Remove the specified set of phone numbers from a location for an organization.\n\nPhone numbers must follow the E.164 format.\n\nRemoving a mobile number may require more time depending on mobile carrier capabilities.\n\nRemoving a phone number from a location requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\nA location's main number cannot be removed.\n\n<br/>\n\n<div><Callout type=\"warning\">This API is only supported for non-integrated PSTN connection types of Local Gateway (LGW) and Non-integrated CPP. It should never be used for locations with integrated PSTN connection types like Cisco Calling Plans or Integrated CCP because backend data issues may occur.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/numbers")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "`LocationId` to which numbers should be added.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization of the Route Group.")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		numbersCmd.AddCommand(cmd)
	}

	{ // validate-phone
		var orgId string
		var phoneNumbers []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "validate-phone",
			Short: "Validate phone numbers",
			Long:  "Validate the list of phone numbers in an organization. Each phone number's availability is indicated in the response.\n\nEach location has a set of phone numbers that can be assigned to people, workspaces, or features. Phone numbers must follow the E.164 format for all countries, except for the United States, which can also follow the National format. Active phone numbers are in service.\n\nValidating a phone number in an organization requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/actions/validateNumbers/invoke")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization of the Route Group.")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		numbersCmd.AddCommand(cmd)
	}

	{ // get-phone-org
		var orgId string
		var locationId string
		var max string
		var start string
		var phoneNumber string
		var available string
		var order string
		var ownerName string
		var ownerId string
		var ownerType string
		var extension string
		var numberType string
		var phoneNumberType string
		var state string
		var details string
		var tollFreeNumbers string
		var restrictedNonGeoNumbers string
		var includedTelephonyTypes string
		var serviceNumber string
		cmd := &cobra.Command{
			Use:   "get-phone-org",
			Short: "Get Phone Numbers for an Organization with Given Criteria",
			Long:  "List all the phone numbers for the given organization along with the status and owner (if any).\n\nNumbers can be standard, service, or mobile. Both standard and service numbers are PSTN numbers.\nService numbers are considered high-utilization or high-concurrency phone numbers and can be assigned to features like auto-attendants, call queues, and hunt groups.\nPhone numbers can be linked to a specific location, be active or inactive, and be assigned or unassigned.\nThe owner of a number is the person, workspace, or feature to which the number is assigned.\nOnly a person can own a mobile number.\n\nRetrieving this list requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/numbers")
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationId", locationId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("available", available)
				req.QueryParam("order", order)
				req.QueryParam("ownerName", ownerName)
				req.QueryParam("ownerId", ownerId)
				req.QueryParam("ownerType", ownerType)
				req.QueryParam("extension", extension)
				req.QueryParam("numberType", numberType)
				req.QueryParam("phoneNumberType", phoneNumberType)
				req.QueryParam("state", state)
				req.QueryParam("details", details)
				req.QueryParam("tollFreeNumbers", tollFreeNumbers)
				req.QueryParam("restrictedNonGeoNumbers", restrictedNonGeoNumbers)
				req.QueryParam("includedTelephonyTypes", includedTelephonyTypes)
				req.QueryParam("serviceNumber", serviceNumber)
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
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Search for this `phoneNumber`.")
		cmd.Flags().StringVar(&available, "available", "", "Search among the available phone numbers. This parameter cannot be used along with `ownerType` parameter when set to `true`.")
		cmd.Flags().StringVar(&order, "order", "", "Sort the list of phone numbers based on the following:`lastName`,`dn`,`extension`. Sorted by number and extension in ascending order.")
		cmd.Flags().StringVar(&ownerName, "owner-name", "", "Return the list of phone numbers that are owned by the given `ownerName`. Maximum length is 255.")
		cmd.Flags().StringVar(&ownerId, "owner-id", "", "Returns only the matched number/extension entries assigned to the feature with the specified UUID or `broadsoftId`.")
		cmd.Flags().StringVar(&ownerType, "owner-type", "", "Returns the list of phone numbers of the given `ownerType`. Possible input values:")
		cmd.Flags().StringVar(&extension, "extension", "", "Returns the list of phone numbers with the given extension.")
		cmd.Flags().StringVar(&numberType, "number-type", "", "Returns the filtered list of phone numbers that contain a given type of number. `available` or `state` query parameters cannot be used when `numberType=EXTENSION`. Possible input values:")
		cmd.Flags().StringVar(&phoneNumberType, "phone-number-type", "", "Returns the filtered list of phone numbers of the given `phoneNumberType`. Response excludes any extensions without numbers. Possible input values:")
		cmd.Flags().StringVar(&state, "state", "", "Returns the list of phone numbers with the matching state. Response excludes any extensions without numbers. Possible input values:")
		cmd.Flags().StringVar(&details, "details", "", "Returns the overall count of the phone numbers along with other details for a given organization.")
		cmd.Flags().StringVar(&tollFreeNumbers, "toll-free-numbers", "", "Returns the list of toll-free phone numbers.")
		cmd.Flags().StringVar(&restrictedNonGeoNumbers, "restricted-non-geo-numbers", "", "Returns the list of restricted non-geographical numbers.")
		cmd.Flags().StringVar(&includedTelephonyTypes, "included-telephony-types", "", "Returns the list of phone numbers that are of given `includedTelephonyTypes`. By default, if this query parameter is not provided, it will list both PSTN and Mobile Numbers. Possible input values are PSTN_NUMBER or MOBILE_NUMBER.")
		cmd.Flags().StringVar(&serviceNumber, "service-number", "", "Returns the list of service phone numbers.")
		numbersCmd.AddCommand(cmd)
	}

	{ // list-manage-jobs
		var orgId string
		var start string
		var max string
		cmd := &cobra.Command{
			Use:   "list-manage-jobs",
			Short: "List Manage Numbers Jobs",
			Long:  "Lists all Manage Numbers jobs for the given organization in order of most recent one to oldest one irrespective of its status.\n\nThe public API only supports initiating jobs which move numbers between locations.\n\nVia Control Hub they can initiate both the move and delete, so this listing can show both.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/jobs/numbers/manageNumbers")
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve list of Manage Number jobs for this organization.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of jobs. Default is 0.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of jobs returned to this maximum count. Default is 2000.")
		numbersCmd.AddCommand(cmd)
	}

	{ // initiate-jobs
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "initiate-jobs",
			Short: "Initiate Number Jobs",
			Long:  "Starts the execution of an operation on a set of numbers. Supported operations are: `MOVE`, `NUMBER_USAGE_CHANGE`.\n\nUp to 1000 numbers can be given in `MOVE` operation type and `NUMBER_USAGE_CHANGE` operation type per request.\nIf another move number job request is initiated while a move job is in progress, the API call will receive a `409` HTTP status code.\n\nIn order to move a number the following is required:\n\n* The number must be unassigned.\n\n* Both locations must have the same PSTN Connection Type.\n\n* Both locations must have the same PSTN Provider.\n\n* Both locations have to be in the same country.\n\nFor example, you can move from Cisco Calling Plan to Cisco Calling Plan, but you cannot move from Cisco Calling Plan to a location with Cloud Connected PSTN.\n\nIn order to change the number usage the following is required:\n\n* The number must be unassigned.\n\n* Number Usage Type can be set to `NONE` if carrier has the PSTN service `GEOGRAPHIC_NUMBERS`.\n\n* Number Usage Type can be set to `SERVICE` if carrier has the PSTN service `SERVICE_NUMBERS`.\n\nFor example, you can initiate a `NUMBER_USAGE_CHANGE` job to change the number type from Standard number to Service number, or the other way around.\n\nThis API requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/jobs/numbers/manageNumbers")
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
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		numbersCmd.AddCommand(cmd)
	}

	{ // get-manage-job-status
		var jobId string
		cmd := &cobra.Command{
			Use:   "get-manage-job-status",
			Short: "Get Manage Numbers Job Status",
			Long:  "Returns the status and other details of the job.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/jobs/numbers/manageNumbers/{jobId}")
				req.PathParam("jobId", jobId)
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
		cmd.Flags().StringVar(&jobId, "job-id", "", "Retrieve job details for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		numbersCmd.AddCommand(cmd)
	}

	{ // pause-manage-job
		var jobId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "pause-manage-job",
			Short: "Pause the Manage Numbers Job",
			Long:  "Pause the running Manage Numbers Job. A paused job can be resumed.\n\nThis API requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/jobs/numbers/manageNumbers/{jobId}/actions/pause/invoke")
				req.PathParam("jobId", jobId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&jobId, "job-id", "", "Pause the Manage Numbers job for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Pause the Manage Numbers job for this organization.")
		numbersCmd.AddCommand(cmd)
	}

	{ // resume-manage-job
		var jobId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "resume-manage-job",
			Short: "Resume the Manage Numbers Job",
			Long:  "Resume the paused Manage Numbers Job. A paused job can be resumed.\n\nThis API requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/jobs/numbers/manageNumbers/{jobId}/actions/resume/invoke")
				req.PathParam("jobId", jobId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&jobId, "job-id", "", "Resume the Manage Numbers job for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Resume the Manage Numbers job for this organization.")
		numbersCmd.AddCommand(cmd)
	}

	{ // list-manage-job-errors
		var jobId string
		var orgId string
		var start string
		var max string
		cmd := &cobra.Command{
			Use:   "list-manage-job-errors",
			Short: "List Manage Numbers Job errors",
			Long:  "Lists all error details of Manage Numbers job. This will not list any errors if `exitCode` is `COMPLETED`. If the status is `COMPLETED_WITH_ERRORS` then this lists the cause of failures.\n\nList of possible Errors:\n\n+ BATCH-1017021 - Failed to move because it is an inactive number.\n\n+ BATCH-1017022 - Failed to move because the source location and target location have different CCP providers.\n\n+ BATCH-1017023 - Failed because it is not an unassigned number.\n\n+ BATCH-1017024 - Failed because it is a main number.\n\n+ BATCH-1017027 - Manage Numbers Move Operation is not supported.\n\n+ BATCH-1017031 - Hydra request is supported only for single number move job.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/jobs/numbers/manageNumbers/{jobId}/errors")
				req.PathParam("jobId", jobId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&jobId, "job-id", "", "Retrieve the error details for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve list of jobs for this organization.")
		cmd.Flags().StringVar(&start, "start", "", "Specifies the error offset from the first result that you want to fetch.")
		cmd.Flags().StringVar(&max, "max", "", "Specifies the maximum number of records that you want to fetch.")
		numbersCmd.AddCommand(cmd)
	}

}
