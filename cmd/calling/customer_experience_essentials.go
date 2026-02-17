package calling

import (
	"fmt"
	"strconv"
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
var _ = strconv.Itoa
var _ = strings.Join

var customerExperienceEssentialsCmd = &cobra.Command{
	Use:   "customer-experience-essentials",
	Short: "CustomerExperienceEssentials commands",
}

func init() {
	cmd.CallingCmd.AddCommand(customerExperienceEssentialsCmd)

	{ // list-wrap-up-reasons
		cmd := &cobra.Command{
			Use:   "list-wrap-up-reasons",
			Short: "List Wrap Up Reasons",
			Long:  "Return the list of wrap-up reasons configured for a customer.\n\nAgents handling calls use wrap-up reasons to categorize the outcome after a call ends. The control hub admin can configure these reasons for customers and assign them to queues. Upon call completion, agents select a wrap-up reason from the queue's assigned list. Each wrap-up reason includes a name and description, and can be set as the default for a queue. Admins can also configure a timer, which dictates the time agents have to select a reason post-call, with a default of 60 seconds. This timer can be disabled if necessary.\n\nRetrieving the list of wrap-up reasons requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/cxEssentials/wrapup/reasons")
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
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // create-wrap-up-reason
		var name string
		var description string
		var queues []string
		var assignAllQueuesEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-wrap-up-reason",
			Short: "Create Wrap Up Reason",
			Long:  "Create a wrap-up reason.\n\nAgents handling calls use wrap-up reasons to categorize the outcome after a call ends. The control hub admin can configure these reasons for customers and assign them to queues.\nUpon call completion, agents select a wrap-up reason from the queue's assigned list. Each wrap-up reason includes a name and description, and can be set as the default for a queue.\nAdmins can also configure a timer, which dictates the time agents have to select a reason post-call, with a default of 60 seconds. This timer can be disabled if necessary.\n\nCreating a wrap-up reason requires a full or device administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/cxEssentials/wrapup/reasons")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("description", description)
					req.BodyStringSlice("queues", queues)
					req.BodyBool("assignAllQueuesEnabled", assignAllQueuesEnabled, cmd.Flags().Changed("assign-all-queues-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().StringSliceVar(&queues, "queues", nil, "")
		cmd.Flags().BoolVar(&assignAllQueuesEnabled, "assign-all-queues-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // get-wrap-up-reason
		var wrapupReasonId string
		cmd := &cobra.Command{
			Use:   "get-wrap-up-reason",
			Short: "Read Wrap Up Reason",
			Long:  "Return the wrap-up reason by ID.\n\nAgents handling calls use wrap-up reasons to categorize the outcome after a call ends. The control hub admin can configure these reasons for customers and assign them to queues.\nUpon call completion, agents select a wrap-up reason from the queue's assigned list. Each wrap-up reason includes a name and description, and can be set as the default for a queue.\nAdmins can also configure a timer, which dictates the time agents have to select a reason post-call, with a default of 60 seconds. This timer can be disabled if necessary.\n\nRetrieving the wrap-up reason by ID requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/cxEssentials/wrapup/reasons/{wrapupReasonId}")
				req.PathParam("wrapupReasonId", wrapupReasonId)
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
		cmd.Flags().StringVar(&wrapupReasonId, "wrapup-reason-id", "", "Wrap-up reason ID.")
		cmd.MarkFlagRequired("wrapup-reason-id")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // update-wrap-up-reason
		var wrapupReasonId string
		var name string
		var description string
		var queuesToAssign []string
		var queuesToUnassign []string
		var assignAllQueuesEnabled bool
		var unassignAllQueuesEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-wrap-up-reason",
			Short: "Update Wrap Up Reason",
			Long:  "Modify a wrap-up reason.\n\nAgents handling calls use wrap-up reasons to categorize the outcome after a call ends. The control hub admin can configure these reasons for customers and assign them to queues.\nUpon call completion, agents select a wrap-up reason from the queue's assigned list. Each wrap-up reason includes a name and description, and can be set as the default for a queue.\nAdmins can also configure a timer, which dictates the time agents have to select a reason post-call, with a default of 60 seconds. This timer can be disabled if necessary.\n\nModifying a wrap-up reason requires a full or device administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/cxEssentials/wrapup/reasons/{wrapupReasonId}")
				req.PathParam("wrapupReasonId", wrapupReasonId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("description", description)
					req.BodyStringSlice("queuesToAssign", queuesToAssign)
					req.BodyStringSlice("queuesToUnassign", queuesToUnassign)
					req.BodyBool("assignAllQueuesEnabled", assignAllQueuesEnabled, cmd.Flags().Changed("assign-all-queues-enabled"))
					req.BodyBool("unassignAllQueuesEnabled", unassignAllQueuesEnabled, cmd.Flags().Changed("unassign-all-queues-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&wrapupReasonId, "wrapup-reason-id", "", "Wrap-up reason ID.")
		cmd.MarkFlagRequired("wrapup-reason-id")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().StringSliceVar(&queuesToAssign, "queues-to-assign", nil, "")
		cmd.Flags().StringSliceVar(&queuesToUnassign, "queues-to-unassign", nil, "")
		cmd.Flags().BoolVar(&assignAllQueuesEnabled, "assign-all-queues-enabled", false, "")
		cmd.Flags().BoolVar(&unassignAllQueuesEnabled, "unassign-all-queues-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // delete-wrap-up-reason
		var wrapupReasonId string
		cmd := &cobra.Command{
			Use:   "delete-wrap-up-reason",
			Short: "Delete Wrap Up Reason",
			Long:  "Delete a wrap-up reason.\n\nAgents handling calls use wrap-up reasons to categorize the outcome after a call ends. The control hub admin can configure these reasons for customers and assign them to queues.\nUpon call completion, agents select a wrap-up reason from the queue's assigned list. Each wrap-up reason includes a name and description, and can be set as the default for a queue.\nAdmins can also configure a timer, which dictates the time agents have to select a reason post-call, with a default of 60 seconds. This timer can be disabled if necessary.\n\nDeleting the wrap-up reason requires a full or device administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/cxEssentials/wrapup/reasons/{wrapupReasonId}")
				req.PathParam("wrapupReasonId", wrapupReasonId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&wrapupReasonId, "wrapup-reason-id", "", "Wrap-up reason ID.")
		cmd.MarkFlagRequired("wrapup-reason-id")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // validate-wrap-up-reason
		var name string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "validate-wrap-up-reason",
			Short: "Validate Wrap Up Reason",
			Long:  "Validate the wrap-up reason name.\n\nAgents handling calls use wrap-up reasons to categorize the outcome after a call ends. The control hub admin can configure these reasons for customers and assign them to queues.\nUpon call completion, agents select a wrap-up reason from the queue's assigned list. Each wrap-up reason includes a name and description, and can be set as the default for a queue.\nAdmins can also configure a timer, which dictates the time agents have to select a reason post-call, with a default of 60 seconds. This timer can be disabled if necessary.\n\nValidating the wrap-up reason name requires a full or device administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/cxEssentials/wrapup/reasons/actions/validateName/invoke")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // get-available-queues
		var wrapupReasonId string
		cmd := &cobra.Command{
			Use:   "get-available-queues",
			Short: "Read Available Queues",
			Long:  "Return the available queues for a wrap-up reason.\n\nAgents handling calls use wrap-up reasons to categorize the outcome after a call ends. The control hub admin can configure these reasons for customers and assign them to queues.\nUpon call completion, agents select a wrap-up reason from the queue's assigned list. Each wrap-up reason includes a name and description, and can be set as the default for a queue.\nAdmins can also configure a timer, which dictates the time agents have to select a reason post-call, with a default of 60 seconds. This timer can be disabled if necessary.\n\nRetrieving the available queues for a wrap-up reason requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/cxEssentials/wrapup/reasons/{wrapupReasonId}/availableQueues")
				req.PathParam("wrapupReasonId", wrapupReasonId)
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
		cmd.Flags().StringVar(&wrapupReasonId, "wrapup-reason-id", "", "Wrap-up reason ID.")
		cmd.MarkFlagRequired("wrapup-reason-id")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // get-wrap-up-reason-settings
		var locationId string
		var queueId string
		cmd := &cobra.Command{
			Use:   "get-wrap-up-reason-settings",
			Short: "Read Wrap Up Reason Settings",
			Long:  "Return a wrap-up reason by location ID and queue ID.\n\nAgents handling calls use wrap-up reasons to categorize the outcome after a call ends. The control hub admin can configure these reasons for customers and assign them to queues.\nUpon call completion, agents select a wrap-up reason from the queue's assigned list. Each wrap-up reason includes a name and description, and can be set as the default for a queue.\nAdmins can also configure a timer, which dictates the time agents have to select a reason post-call, with a default of 60 seconds. This timer can be disabled if necessary.\n\nRetrieving the wrap-up reason by location ID and queue ID requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/cxEssentials/locations/{locationId}/queues/{queueId}/wrapup/settings")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "The location ID.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "The queue ID.")
		cmd.MarkFlagRequired("queue-id")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // update-wrap-up-reason-settings
		var locationId string
		var queueId string
		var wrapupReasons []string
		var defaultWrapupReasonId string
		var wrapupTimerEnabled bool
		var wrapupTimer int64
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-wrap-up-reason-settings",
			Short: "Update Wrap Up Reason Settings",
			Long:  "Modify a wrap-up reason by location ID and queue ID.\n\nAgents handling calls use wrap-up reasons to categorize the outcome after a call ends. The control hub admin can configure these reasons for customers and assign them to queues.\nUpon call completion, agents select a wrap-up reason from the queue's assigned list. Each wrap-up reason includes a name and description, and can be set as the default for a queue.\nAdmins can also configure a timer, which dictates the time agents have to select a reason post-call, with a default of 60 seconds. This timer can be disabled if necessary.\n\nModifying a wrap-up reason by location ID and queue ID requires a full or device administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/cxEssentials/locations/{locationId}/queues/{queueId}/wrapup/settings")
				req.PathParam("locationId", locationId)
				req.PathParam("queueId", queueId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("wrapupReasons", wrapupReasons)
					req.BodyString("defaultWrapupReasonId", defaultWrapupReasonId)
					req.BodyBool("wrapupTimerEnabled", wrapupTimerEnabled, cmd.Flags().Changed("wrapup-timer-enabled"))
					req.BodyInt("wrapupTimer", wrapupTimer, cmd.Flags().Changed("wrapup-timer"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "The location ID.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "The queue ID.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringSliceVar(&wrapupReasons, "wrapup-reasons", nil, "")
		cmd.Flags().StringVar(&defaultWrapupReasonId, "default-wrapup-reason-id", "", "")
		cmd.Flags().BoolVar(&wrapupTimerEnabled, "wrapup-timer-enabled", false, "")
		cmd.Flags().Int64Var(&wrapupTimer, "wrapup-timer", 0, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // get-screen-pop-configuration
		var locationId string
		var queueId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-screen-pop-configuration",
			Short: "Read Screen Pop Configuration",
			Long:  "Returns the screen pop configuration for a call queue in a location.\n\nScreen pop lets agents view customer-related info in a pop-up window.\n\nRetrieving the screen pop configuration requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/queues/{queueId}/cxEssentials/screenPop")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "The location ID where the call queue resides.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "The call queue ID for which screen pop configuration is modified.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID of the customer or partner's organization.")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // update-screen-pop-configuration
		var locationId string
		var queueId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-screen-pop-configuration",
			Short: "Update Screen Pop Configuration",
			Long:  "Modifies the screen pop configuration for a call queue in a location.\n\nScreen pop lets agents view customer-related info in a pop-up window.\n\nModifying the screen pop configuration requires a full or device administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/queues/{queueId}/cxEssentials/screenPop")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "The location ID where the call queue resides.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "The call queue ID for which screen pop configuration is modified.")
		cmd.MarkFlagRequired("queue-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID of the customer or partner's organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

	{ // list-available-agents
		var locationId string
		var orgId string
		var hasCxEssentials string
		cmd := &cobra.Command{
			Use:   "list-available-agents",
			Short: "List Available Agents",
			Long:  "Return a list of available agents with Customer Experience Essentials license in a location.\n\nRetrieving the list of available agents requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/cxEssentials/agents/availableAgents")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve the list of avaiilable agents in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID of the customer or partner's organization.")
		cmd.Flags().StringVar(&hasCxEssentials, "has-cx-essentials", "", "Returns only the list of available agents with Customer Experience Essentials license when `true`, otherwise returns the list of available agents with Customer Experience Basic license.")
		customerExperienceEssentialsCmd.AddCommand(cmd)
	}

}
