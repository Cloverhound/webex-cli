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

var callQueueCmd = &cobra.Command{
	Use:   "call-queue",
	Short: "CallQueue commands",
}

func init() {
	cmd.CallingCmd.AddCommand(callQueueCmd)

	{ // list-cxe
		var orgId string
		var locationId string
		var max string
		var start string
		var name string
		var phoneNumber string
		var departmentId string
		var departmentName string
		var hasCxEssentials string
		cmd := &cobra.Command{
			Use:   "list-cxe",
			Short: "Read the List of Call Queues with Customer Experience Essentials",
			Long:  "List all Call Queues for the organization.\n\nCall queues temporarily hold calls in the cloud, when all agents\nassigned to receive calls from the queue are unavailable. Queued calls are routed to \nan available agent, when not on an active call. Each call queue is assigned a lead number, which is a telephone\nnumber that external callers can dial to reach the users assigned to the call queue.\nCall queues are also assigned an internal extension, which can be dialed\ninternally to reach the users assigned to the call queue.\n\nRetrieving this list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/queues")
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationId", locationId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("departmentId", departmentId)
				req.QueryParam("departmentName", departmentName)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Returns the list of call queues in this organization.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Returns the list of call queues in this location.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&name, "name", "", "Returns only the call queues matching the given name.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Returns only the call queues matching the given primary phone number or extension.")
		cmd.Flags().StringVar(&departmentId, "department-id", "", "Returns only call queues matching the given department ID.")
		cmd.Flags().StringVar(&departmentName, "department-name", "", "Returns only call queues matching the given department name.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Returns only the list of call queues with Customer Experience Essentials license when `true`, otherwise returns the list of Customer Experience Basic call queues.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // create-cxe
		var locationId string
		var orgId string
		var hasCxEssentials string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-cxe",
			Short: "Create a Call Queue with Customer Experience Essentials",
			Long:  "Create new Call Queues for the given location.\n\nCall queues temporarily hold calls in the cloud, when all agents assigned to receive calls from the queue are unavailable.\nQueued calls are routed to an available agent, when not on an active call. Each call queue is assigned a lead number, which is a telephone\nnumber that external callers can dial to reach the users assigned to the call queue. Call queues are also assigned an internal extension,\nwhich can be dialed internally to reach the users assigned to the call queue.\n\nCreating a call queue requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, and `dialByName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/queues")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "The location ID where the call queue needs to be created.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID where the call queue needs to be created.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Creates a Customer Experience Essentials call queue, when `true`. This requires Customer Experience Essentials licensed agents.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // delete
		var locationId string
		var queueId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Call Queue",
			Long:  "Delete the designated Call Queue.\n\nCall queues temporarily hold calls in the cloud when all agents, which\ncan be users or agents, assigned to receive calls from the queue are\nunavailable. Queued calls are routed to an available agent when not on an\nactive call. Each call queue is assigned a Lead Number, which is a telephone\nnumber outside callers can dial to reach users assigned to the call queue.\nCall queues are also assigned an internal extension, which can be dialed\ninternally to reach users assigned to the call queue.\n\nDeleting a call queue requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/queues/{queueId}")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location from which to delete a call queue.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Delete the call queue with the matching ID.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the call queue from this organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-cxe
		var locationId string
		var queueId string
		var orgId string
		var hasCxEssentials string
		cmd := &cobra.Command{
			Use:   "get-cxe",
			Short: "Get Details for a Call Queue with Customer Experience Essentials",
			Long:  "Retrieve Call Queue details.\n\nCall queues temporarily hold calls in the cloud, when all agents assigned to receive calls from the queue are unavailable.\nQueued calls are routed to an available agent, when not on an active call. Each call queue is assigned a lead number, which is a telephone\nnumber that external callers can dial to reach the users assigned to the call queue. Call queues are also assigned an internal extension,\nwhich can be dialed internally to reach the users assigned to the call queue.\n\nRetrieving call queue details requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, and `dialByName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/{queueId}")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieves the details of a call queue in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Retrieves the details of call queue with this identifier.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieves the details of a call queue in this organization.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Must be set to `true`, to view the details of a call queue with Customer Experience Essentials license. This can otherwise be ommited or set to `false`.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // update
		var locationId string
		var queueId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Call Queue",
			Long:  "Update the designated Call Queue.\n\nCall queues temporarily hold calls in the cloud when all agents, which\ncan be users or agents, assigned to receive calls from the queue are\nunavailable. Queued calls are routed to an available agent when not on an\nactive call. Each call queue is assigned a Lead Number, which is a telephone\nnumber outside callers can dial to reach users assigned to the call queue.\nCall queues are also assigned an internal extension, which can be dialed\ninternally to reach users assigned to the call queue.\n\nUpdating a call queue requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, and `dialByName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/queues/{queueId}")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Update setting for the call queue with the matching ID.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update call queue settings from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // list-announcement-files
		var locationId string
		var queueId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-announcement-files",
			Short: "Read the List of Call Queue Announcement Files",
			Long:  "List file info for all Call Queue announcement files associated with this Call Queue.\n\nCall Queue announcement files contain messages and music that callers hear while waiting in the queue. A call queue can be configured to play whatever subset of these announcement files is desired.\n\nRetrieving this list of files requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.\n\nNote that uploading of announcement files via API is not currently supported, but is available via Webex Control Hub.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/{queueId}/announcements")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Retrieve anouncement files for the call queue with this identifier.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve announcement files for a call queue from this organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // delete-announcement-file
		var locationId string
		var queueId string
		var fileName string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-announcement-file",
			Short: "Delete a Call Queue Announcement File",
			Long:  "Delete an announcement file for the designated Call Queue.\n\nCall Queue announcement files contain messages and music that callers hear while waiting in the queue. A call queue can be configured to play whatever subset of these announcement files is desired.\n\nDeleting an announcement file for a call queue requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/queues/{queueId}/announcements/{fileName}")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
				req.PathParam("fileName", fileName)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Delete an announcement for a call queue in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Delete an announcement for the call queue with this identifier.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&fileName, "file-name", "", "")
		cmd.MarkFlagRequired("file-name")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete call queue announcement from this organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-forward
		var locationId string
		var queueId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-forward",
			Short: "Get Call Forwarding Settings for a Call Queue",
			Long:  "Retrieve Call Forwarding settings for the specified Call Queue, including the list of call forwarding rules.\n\nThe call forwarding feature allows you to direct all incoming calls based on specific criteria that you define.\nBelow are the available options for configuring your call forwarding:\n1. Always forward calls to a designated number.\n2. Forward calls to a designated number based on certain criteria.\n3. Forward calls using different modes.\n\nRetrieving call forwarding settings for a call queue requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/{queueId}/callForwarding")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Retrieve the call forwarding settings for this call queue.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call queue forwarding settings from this organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // update-forward
		var locationId string
		var queueId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-forward",
			Short: "Update Call Forwarding Settings for a Call Queue",
			Long:  "Update Call Forwarding settings for the designated Call Queue.\n\nThe call forwarding feature allows you to direct all incoming calls based on specific criteria that you define.\nBelow are the available options for configuring your call forwarding:\n1. Always forward calls to a designated number.\n2. Forward calls to a designated number based on certain criteria.\n3. Forward calls using different modes.\n\nUpdating call forwarding settings for a call queue requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/queues/{queueId}/callForwarding")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Update call forwarding settings for this call queue.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update call queue forwarding settings from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // create-selective-forward-rule
		var locationId string
		var queueId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-selective-forward-rule",
			Short: "Create a Selective Call Forwarding Rule for a Call Queue",
			Long:  "Create a Selective Call Forwarding Rule for the designated Call Queue.\n\nA selective call forwarding rule for a call queue allows calls to be forwarded or not forwarded to the designated number, based on the defined criteria.\n\nNote that the list of existing call forward rules is available in the call queue's call forwarding settings.\n\nCreating a selective call forwarding rule for a call queue requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The Call Forwarding Rule ID will change upon modification of the Call Forwarding Rule name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/queues/{queueId}/callForwarding/selectiveRules")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which the call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Create the rule for this call queue.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create the call queue rule for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-selective-forward-rule
		var locationId string
		var queueId string
		var ruleId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-selective-forward-rule",
			Short: "Get Selective Call Forwarding Rule for a Call Queue",
			Long:  "Retrieve a Selective Call Forwarding Rule's settings for the designated Call Queue.\n\nA selective call forwarding rule for a call queue allows calls to be forwarded or not forwarded to the designated number, based on the defined criteria.\n\nNote that the list of existing call forward rules is available in the call queue's call forwarding settings.\n\nRetrieving a selective call forwarding rule's settings for a call queue requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.\n\n**NOTE**: The Call Forwarding Rule ID will change upon modification of the Call Forwarding Rule name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/{queueId}/callForwarding/selectiveRules/{ruleId}")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which to call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Retrieve setting for a rule for this call queue.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&ruleId, "rule-id", "", "Call queue rule you are retrieving settings for.")
		cmd.MarkFlagRequired("rule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call queue rule settings for this organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // update-selective-forward-rule
		var locationId string
		var queueId string
		var ruleId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-selective-forward-rule",
			Short: "Update a Selective Call Forwarding Rule for a Call Queue",
			Long:  "Update a Selective Call Forwarding Rule's settings for the designated Call Queue.\n\nA selective call forwarding rule for a call queue allows calls to be forwarded or not forwarded to the designated number, based on the defined criteria.\n\nNote that the list of existing call forward rules is available in the call queue's call forwarding settings.\n\nUpdating a selective call forwarding rule's settings for a call queue requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The Call Forwarding Rule ID will change upon modification of the Call Forwarding Rule name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/queues/{queueId}/callForwarding/selectiveRules/{ruleId}")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Update settings for a rule for this call queue.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&ruleId, "rule-id", "", "Call queue rule you are updating settings for.")
		cmd.MarkFlagRequired("rule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update call queue rule settings for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // delete-selective-forward-rule
		var locationId string
		var queueId string
		var ruleId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-selective-forward-rule",
			Short: "Delete a Selective Call Forwarding Rule for a Call Queue",
			Long:  "Delete a Selective Call Forwarding Rule for the designated Call Queue.\n\nA selective call forwarding rule for a call queue allows calls to be forwarded or not forwarded to the designated number, based on the defined criteria.\n\nNote that the list of existing call forward rules is available in the call queue's call forwarding settings.\n\nDeleting a selective call forwarding rule for a call queue requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The Call Forwarding Rule ID will change upon modification of the Call Forwarding Rule name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/queues/{queueId}/callForwarding/selectiveRules/{ruleId}")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
				req.PathParam("ruleId", ruleId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Delete the rule for this call queue.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&ruleId, "rule-id", "", "Call queue rule you are deleting.")
		cmd.MarkFlagRequired("rule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete call queue rule from this organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-holiday-service
		var locationId string
		var queueId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-holiday-service",
			Short: "Get Details for a Call Queue Holiday Service",
			Long:  "Retrieve Call Queue Holiday Service details.\n\nConfigure the call queue to route calls differently during the holidays.\n\nRetrieving call queue holiday service details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/{queueId}/holidayService")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve settings for a call queue in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Retrieve settings for the call queue with this identifier.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call queue settings from this organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // update-holiday-service
		var locationId string
		var queueId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-holiday-service",
			Short: "Update a Call Queue Holiday Service",
			Long:  "Update the designated Call Queue Holiday Service.\n\nConfigure the call queue to route calls differently during the holidays.\n\nUpdating a call queue holiday service requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/queues/{queueId}/holidayService")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Update setting for the call queue with the matching ID.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update call queue settings from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-night-service
		var locationId string
		var queueId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-night-service",
			Short: "Get Details for a Call Queue Night Service",
			Long:  "Retrieve Call Queue Night service details.\n\nConfigure the call queue to route calls differently during the hours when the queue is not in service. This is\ndetermined by a schedule that defines the business hours of the queue.\n\nRetrieving call queue details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/{queueId}/nightService")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve settings for a call queue in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Retrieve settings for the call queue night service with this identifier.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call queue night service settings from this organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // update-night-service
		var locationId string
		var queueId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-night-service",
			Short: "Update a Call Queue Night Service",
			Long:  "Update Call Queue Night Service details.\n\nConfigure the call queue to route calls differently during the hours when the queue is not in service. This is\ndetermined by a schedule that defines the business hours of the queue.\n\nUpdating call queue night service details requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/queues/{queueId}/nightService")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update settings for a call queue in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Update settings for the call queue night service with this identifier.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update call queue night service settings from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-forced-forward
		var locationId string
		var queueId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-forced-forward",
			Short: "Get Details for a Call Queue Forced Forward",
			Long:  "Retrieve Call Queue policy Forced Forward details.\n\nThis policy allows calls to be temporarily diverted to a configured destination.\n\nRetrieving call queue Forced Forward details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/{queueId}/forcedForward")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve settings for a call queue in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Retrieve settings for the call queue with this identifier.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call queue settings from this organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // update-forced-forward-service
		var locationId string
		var queueId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-forced-forward-service",
			Short: "Update a Call Queue Forced Forward service",
			Long:  "Update the designated Forced Forward Service.\n\nIf the option is enabled, then incoming calls to the queue are forwarded to the configured destination. Calls that are already in the queue remain queued.\nThe policy can be configured to play an announcement prior to proceeding with the forward.\n\nUpdating a call queue Forced Forward service requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/queues/{queueId}/forcedForward")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Update setting for the call queue with the matching ID.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update call queue settings from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-stranded
		var locationId string
		var queueId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-stranded",
			Short: "Get Details for a Call Queue Stranded Calls",
			Long:  "Allow admin to view default/configured Stranded Calls settings.\n\nStranded-All agents logoff Policy: If the last agent staffing a queue \u201cunjoins\u201d the queue or signs out, then all calls in the queue become stranded.\nStranded-Unavailable Policy: This policy allows for the configuration of the processing of calls that are in a staffed queue when all agents are unavailable.\n\nRetrieving call queue Stranded Calls details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/{queueId}/strandedCalls")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve settings for a call queue in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Retrieve settings for the call queue with this identifier.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call queue settings from this organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // update-stranded-service
		var locationId string
		var queueId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-stranded-service",
			Short: "Update a Call Queue Stranded Calls service",
			Long:  "Update the designated Call Stranded Calls Service.\n\nAllow admin to modify configured Stranded Calls settings.\n\nUpdating a call queue stranded calls requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/queues/{queueId}/strandedCalls")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this call queue exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Update setting for the call queue with the matching ID.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update call queue settings from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-primary-numbers
		var locationId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-primary-numbers",
			Short: "Get Call Queue Primary Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as the call queue's primary phone number.\nThese numbers are associated with the location specified in the request URL, can be active or inactive, and are unassigned.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/availableNumbers")
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
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-alternate-numbers
		var locationId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-alternate-numbers",
			Short: "Get Call Queue Alternate Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as the call queue's alternate phone number.\nThese numbers are associated with the location specified in the request URL, can be active or inactive, and are unassigned.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/alternate/availableNumbers")
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
		callQueueCmd.AddCommand(cmd)
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
			Short: "Get Call Queue Call Forward Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as the call queue's call forward number.\nThese numbers are associated with the location specified in the request URL, can be active or inactive, and are assigned to an owning entity.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/callForwarding/availableNumbers")
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
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-available-agents
		var locationId string
		var orgId string
		var max string
		var start string
		var name string
		var phoneNumber string
		var order string
		cmd := &cobra.Command{
			Use:   "get-available-agents",
			Short: "Get Call Queue Available Agents",
			Long:  "List all available users, workspaces, or virtual lines that can be assigned as call queue agents.\n\nAvailable agents are users (excluding users with Webex Calling Standard license), workspaces, or virtual lines that can be assigned to a call queue. \nCalls from the call queue are routed to assigned agents based on configuration. \nAn agent can be assigned to one or more call queues and can be managed by supervisors.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/queues/agents/availableAgents")
				req.QueryParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "The location ID of the call queue. Temporary mandatory query parameter, used for performance reasons only and not a filter.")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List available agents for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&name, "name", "", "Search based on name (user first and last name combination).")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Search based on number or extension.")
		cmd.Flags().StringVar(&order, "order", "", "Order the available agents according to the designated fields. Up to three comma-separated sort order fields may be specified. Available sort fields are: `userId`, `fname`, `firstname`, `lname`, `lastname`, `dn`, and `extension`. Sort order can be added together with each field using a hyphen, `-`. Available sort orders are: `asc`, and `desc`.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // list-supervisors-cxe
		var orgId string
		var max string
		var start string
		var name string
		var phoneNumber string
		var order string
		var hasCxEssentials string
		cmd := &cobra.Command{
			Use:   "list-supervisors-cxe",
			Short: "Get List of Supervisors with Customer Experience Essentials",
			Long:  "Get list of supervisors for an organization.\n\nAgents in a call queue can be associated with a supervisor who can silently monitor, coach, barge in or to take over calls that their assigned agents are currently handling.\n\nRequires a full, location, user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/supervisors")
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("order", order)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List the supervisors in this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&name, "name", "", "Only return the supervisors that match the given name.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Only return the supervisors that match the given phone number, extension, or ESN.")
		cmd.Flags().StringVar(&order, "order", "", "Sort results alphabetically by supervisor name, in ascending or descending order.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Returns only the list of supervisors with Customer Experience Essentials license, when `true`. Otherwise returns the list of supervisors with Customer Experience Basic license.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // create-supervisor-cxe
		var orgId string
		var hasCxEssentials string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-supervisor-cxe",
			Short: "Create a Supervisor with Customer Experience Essentials",
			Long:  "Create a new supervisor. The supervisor must be created with at least one agent.\n\nAgents in a call queue can be associated with a supervisor who can silently monitor, coach, barge in or to take over calls that their assigned agents are currently handling.\n\nThis operation requires a full or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/supervisors")
				req.QueryParam("orgId", orgId)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID where the supervisor needs to be created.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Creates a Customer Experience Essentials queue supervisor, when `true`. Customer Experience Essentials queue supervisors must have a Customer Experience Essentials license.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // delete-bulk-supervisors
		var orgId string
		var supervisorIds []string
		var deleteAll bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "delete-bulk-supervisors",
			Short: "Delete Bulk supervisors",
			Long:  "Deletes supervisors in bulk from an organization.\n\nSupervisors are users who manage agents and who perform functions including monitoring, coaching, and more.\n\nRequires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/supervisors")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("supervisorIds", supervisorIds)
					req.BodyBool("deleteAll", deleteAll, cmd.Flags().Changed("delete-all"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete supervisors in bulk for this organization.")
		cmd.Flags().StringSliceVar(&supervisorIds, "supervisor-ids", nil, "")
		cmd.Flags().BoolVar(&deleteAll, "delete-all", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // delete-supervisor
		var supervisorId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-supervisor",
			Short: "Delete A Supervisor",
			Long:  "Deletes the supervisor from an organization.\n\nSupervisors are users who manage agents and who perform functions including monitoring, coaching, and more.\n\nRequires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/supervisors/{supervisorId}")
				req.PathParam("supervisorId", supervisorId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&supervisorId, "supervisor-id", "", "Delete the specified supervisor.")
		cmd.MarkFlagRequired("supervisor-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the supervisor in the specified organization.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // supervisor-detail-cxe
		var supervisorId string
		var orgId string
		var max string
		var start string
		var name string
		var phoneNumber string
		var order string
		var hasCxEssentials string
		cmd := &cobra.Command{
			Use:   "supervisor-detail-cxe",
			Short: "GET Supervisor Detail with Customer Experience Essentials",
			Long:  "Get details of a specific supervisor, which includes the agents associated agents with the supervisor, in an organization.\n\nAgents in a call queue can be associated with a supervisor who can silently monitor, coach, barge in or to take over calls that their assigned agents are currently handling.\n\nThis operation requires a full, user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/supervisors/{supervisorId}")
				req.PathParam("supervisorId", supervisorId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("order", order)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&supervisorId, "supervisor-id", "", "List the agents assigned to this supervisor.")
		cmd.MarkFlagRequired("supervisor-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List the agents assigned to a supervisor in this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&name, "name", "", "Only return the agents that match the given name.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Only return agents that match the given phone number, extension, or ESN.")
		cmd.Flags().StringVar(&order, "order", "", "Sort results alphabetically by supervisor name, in ascending or descending order.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Must be set to `true`, to view the details of a supervisor with Customer Experience Essentials license. This can otherwise be ommited or set to `false`.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // assign-unassign-agents-supervisor-cxe
		var supervisorId string
		var orgId string
		var hasCxEssentials string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "assign-unassign-agents-supervisor-cxe",
			Short: "Assign or Unassign Agents to Supervisor with Customer Experience Essentials",
			Long:  "Assign or unassign agents to the supervisor for an organization.\n\nAgents in a call queue can be associated with a supervisor who can silently monitor, coach, barge in or to take over calls that their assigned agents are currently handling.\n\nThis operation requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/supervisors/{supervisorId}")
				req.PathParam("supervisorId", supervisorId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&supervisorId, "supervisor-id", "", "Identifier of the supervisor to be updated.")
		cmd.MarkFlagRequired("supervisor-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Assign or unassign agents to a supervisor in this organization.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Must be set to `true` to modify a supervisor with Customer Experience Essentials license. This can otherwise be ommited or set to `false`.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // list-available-supervisors-cxe
		var orgId string
		var max string
		var start string
		var name string
		var phoneNumber string
		var order string
		var hasCxEssentials string
		cmd := &cobra.Command{
			Use:   "list-available-supervisors-cxe",
			Short: "List Available Supervisors with Customer Experience Essentials",
			Long:  "Get list of available supervisors for an organization.\n\nAgents in a call queue can be associated with a supervisor who can silently monitor, coach, barge in or to take over calls that their assigned agents are currently handling.\n\nThis operation requires a full, user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/supervisors/availableSupervisors")
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("order", order)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List the available supervisors in this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&name, "name", "", "Only return the supervisors that match the given name.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Only return the supervisors that match the given phone number, extension, or ESN.")
		cmd.Flags().StringVar(&order, "order", "", "Sort results alphabetically by supervisor name, in ascending or descending order.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Returns only the list of available supervisors with Customer Experience Essentials license, when `true`. When ommited or set to 'false', will return the list of available supervisors with Customer Experience Basic license.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // list-available-agents-cxe
		var orgId string
		var max string
		var start string
		var name string
		var phoneNumber string
		var order string
		var hasCxEssentials string
		cmd := &cobra.Command{
			Use:   "list-available-agents-cxe",
			Short: "List Available Agents with Customer Experience Essentials",
			Long:  "Get list of available agents for an organization.\n\nAgents in a call queue can be associated with a supervisor who can silently monitor, coach, barge in or to take over calls that their assigned agents are currently handling.\n\nThis operation requires a full, user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/supervisors/availableAgents")
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("order", order)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List of available agents in a supervisor's list for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&name, "name", "", "Returns only the agents that match the given name.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Returns only the agents that match the phone number, extension, or ESN.")
		cmd.Flags().StringVar(&order, "order", "", "Sort results alphabetically by supervisor name, in ascending or descending order.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Returns only the list of available agents with Customer Experience Essentials license, when `true`. When ommited or set to `false`, will return the list of available agents with Customer Experience Basic license.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // list-agents-cxe
		var orgId string
		var locationId string
		var queueId string
		var max string
		var start string
		var name string
		var phoneNumber string
		var joinEnabled string
		var hasCxEssentials string
		var order string
		cmd := &cobra.Command{
			Use:   "list-agents-cxe",
			Short: "Read the List of Call Queue Agents with Customer Experience Essentials",
			Long:  "List all Call Queues Agents for the organization.\n\nAgents can be users, workplace or virtual lines assigned to a call queue. Calls from the call queue are routed to agents based on configuration. \nAn agent can be assigned to one or more call queues and can be managed by supervisors.\n\nRetrieving this list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.\n\n**Note**: The decoded value of the agent's `id`, and the `type` returned in the response, are always returned as `PEOPLE`, even when the agent is a workspace or virtual line. This will be addressed in a future release.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/queues/agents")
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationId", locationId)
				req.QueryParam("queueId", queueId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("joinEnabled", joinEnabled)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List call queues agents in this organization.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return only the call queue agents in this location.")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Only return call queue agents with the matching queue ID.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&name, "name", "", "Returns only the list of call queue agents that match the given name.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Returns only the list of call queue agents that match the given phone number or extension.")
		cmd.Flags().StringVar(&joinEnabled, "join-enabled", "", "Returns only the list of call queue agents that match the given `joinEnabled` value.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Returns only the list of call queues with Customer Experience Essentials license when `true`, otherwise returns the list of Customer Experience Basic call queues.")
		cmd.Flags().StringVar(&order, "order", "", "Sort results alphabetically by call queue agent's name, in ascending or descending order.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-agent-cxe
		var id string
		var orgId string
		var hasCxEssentials string
		var max string
		var start string
		cmd := &cobra.Command{
			Use:   "get-agent-cxe",
			Short: "Get Details for a Call Queue Agent with Customer Experience Essentials",
			Long:  "Retrieve details of a particular Call queue agent based on the agent ID.\n\nAgents can be users, workplace or virtual lines assigned to a call queue. Calls from the call queue are routed to agents based on configuration. \nAn agent can be assigned to one or more call queues and can be managed by supervisors.\n\nRetrieving a call queue agent's details require a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.\n\n**Note**: The agent's `type` returned in the response and in the decoded value of the agent's `id`, is always of type `PEOPLE`, even if the agent is a workspace or virtual line. This` will be corrected in a future release.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/queues/agents/{id}")
				req.PathParam("id", id)
				req.QueryParam("orgId", orgId)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&id, "id", "", "Retrieve call queue agents with this identifier.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call queue agents from this organization.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Must be set to `true` to view the details of an agent with Customer Experience Essentials license. This can otherwise be ommited or set to `false`.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		callQueueCmd.AddCommand(cmd)
	}

	{ // update-agent-settings-one-more-cxe
		var id string
		var orgId string
		var hasCxEssentials string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-agent-settings-one-more-cxe",
			Short: "Update an Agent's Settings of One or More Call Queues with Customer Experience Essentials",
			Long:  "Modify an agent's call queue settings for an organization.\n\nCalls from the call queue are routed to agents based on configuration. \nAn agent can be assigned to one or more call queues and can be managed by supervisors.\n\nThis operation requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/queues/agents/{id}/settings")
				req.PathParam("id", id)
				req.QueryParam("orgId", orgId)
				req.QueryParam("hasCxEssentials", hasCxEssentials)
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
		cmd.Flags().StringVar(&id, "id", "", "Identifier of the agent to be updated.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update the settings of an agent in this organization.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Must be set to `true` to modify an agent that has Customer Experience Essentials license. This can otherwise be ommited or set to `false`.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // switch-mode-forward
		var locationId string
		var queueId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "switch-mode-forward",
			Short: "Switch Mode for Call Forwarding Settings for a Call Queue",
			Long:  "Switches the current operating mode of the `Call Queue` to the mode as per normal operations.\n\nSwitching operating mode for a `call queue` requires a full, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/queues/{queueId}/callForwarding/actions/switchMode/invoke")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "`Location` in which this `call queue` exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "Switch operating mode to normal operations for this `call queue`.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Switch operating mode as per normal operations for the `call queue` from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callQueueCmd.AddCommand(cmd)
	}

	{ // get-playlist-usage
		var playListId string
		var playlistUsageType string
		cmd := &cobra.Command{
			Use:   "get-playlist-usage",
			Short: "Get Playlist Usage",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/announcements/playlists/{playListId}/usage")
				req.PathParam("playListId", playListId)
				req.QueryParam("playlistUsageType", playlistUsageType)
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
		cmd.Flags().StringVar(&playListId, "play-list-id", "", "Unique identifier of the playlist.")
		cmd.MarkFlagRequired("play-list-id")
		cmd.Flags().StringVar(&playlistUsageType, "playlist-usage-type", "", "Filter usage by type.")
		callQueueCmd.AddCommand(cmd)
	}

}
