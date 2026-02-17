package cc

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

var callbacksCmd = &cobra.Command{
	Use:   "callbacks",
	Short: "Callbacks commands",
}

func init() {
	cmd.CcCmd.AddCommand(callbacksCmd)

	{ // schedule
		var orgId string
		var customerName string
		var callbackNumber string
		var timezone string
		var scheduleDate string
		var startTime string
		var endTime string
		var queueId string
		var callbackReason string
		var sourceInteraction string
		var assigneeAgent string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "schedule",
			Short: "Schedule a Callback",
			Long:  `Creates a new callback request for a customer. Authorization requires the cjp:user scope. The callback default endpoint (EP) and default ANI must be configured as mandatory settings to successfully make API calls.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/v1/callbacks/organization/{orgId}/scheduled-callback")
				req.PathParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("customerName", customerName)
					req.BodyString("callbackNumber", callbackNumber)
					req.BodyString("timezone", timezone)
					req.BodyString("scheduleDate", scheduleDate)
					req.BodyString("startTime", startTime)
					req.BodyString("endTime", endTime)
					req.BodyString("queueId", queueId)
					req.BodyString("callbackReason", callbackReason)
					req.BodyString("sourceInteraction", sourceInteraction)
					req.BodyString("assigneeAgent", assigneeAgent)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID for which the callback is being scheduled. This should be a valid UUID.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&customerName, "customer-name", "", "")
		cmd.Flags().StringVar(&callbackNumber, "callback-number", "", "")
		cmd.Flags().StringVar(&timezone, "timezone", "", "")
		cmd.Flags().StringVar(&scheduleDate, "schedule-date", "", "")
		cmd.Flags().StringVar(&startTime, "start-time", "", "")
		cmd.Flags().StringVar(&endTime, "end-time", "", "")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "")
		cmd.Flags().StringVar(&callbackReason, "callback-reason", "", "")
		cmd.Flags().StringVar(&sourceInteraction, "source-interaction", "", "")
		cmd.Flags().StringVar(&assigneeAgent, "assignee-agent", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callbacksCmd.AddCommand(cmd)
	}

	{ // get-scheduled
		var orgId string
		var callbackNumber string
		var assigneeAgent string
		var page string
		var pageSize string
		var sortBy string
		var sortOrder string
		cmd := &cobra.Command{
			Use:   "get-scheduled",
			Short: "Get scheduled callbacks",
			Long:  `Allows the user to list scheduled callbacks for a given customer number or assignee agent, excluding those whose scheduled trigger time has already passed. Requires 'cjp:user' scope for authorization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/v1/callbacks/organization/{orgId}/scheduled-callback")
				req.PathParam("orgId", orgId)
				req.QueryParam("callbackNumber", callbackNumber)
				req.QueryParam("assigneeAgent", assigneeAgent)
				req.QueryParam("page", page)
				req.QueryParam("pageSize", pageSize)
				req.QueryParam("sortBy", sortBy)
				req.QueryParam("sortOrder", sortOrder)
				if config.Paginate() {
					resp, statusCode, err := req.DoPaginated(false)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID for which the callback is being scheduled. This should be a valid UUID.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&callbackNumber, "callback-number", "", "The callback customer number to filter the scheduled callbacks. Only an exact match will yield the result. Allows an optional country code followed by digits (0-9) and the special characters: space, hyphen -, parentheses ( and ), and period ., ensuring the total length is between 7 and 15 characters.")
		cmd.Flags().StringVar(&assigneeAgent, "assignee-agent", "", "The unique identifier of the agent assigned to handle the callback. Must be in UUID format. This parameter is optional, but at least one of assigneeAgent or callbackNumber must be provided.")
		cmd.Flags().StringVar(&page, "page", "", "The page number to retrieve.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "The number of items per page.")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "The field to sort the results by. If `sortBy` is set to `assignedTime`, the `assigneeAgent` parameter must also be provided.")
		cmd.Flags().StringVar(&sortOrder, "sort-order", "", "The order to sort the results in.")
		callbacksCmd.AddCommand(cmd)
	}

	{ // get-scheduled-id
		var orgId string
		var id string
		cmd := &cobra.Command{
			Use:   "get-scheduled-id",
			Short: "Get scheduled callback by Id",
			Long:  `Retrieve an existing scheduled callback by Id, excluding those whose scheduled trigger time has already passed. Requires 'cjp:user' scope for authorization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/v1/callbacks/organization/{orgId}/scheduled-callback/{id}")
				req.PathParam("orgId", orgId)
				req.PathParam("id", id)
				if config.Paginate() {
					resp, statusCode, err := req.DoPaginated(false)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID for which the callback is being scheduled. This should be a valid UUID.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&id, "id", "", "The id with which the Scheduled Callback has been created.")
		cmd.MarkFlagRequired("id")
		callbacksCmd.AddCommand(cmd)
	}

	{ // update-scheduled-id
		var orgId string
		var id string
		var customerName string
		var callbackNumber string
		var timezone string
		var scheduleDate string
		var startTime string
		var endTime string
		var queueId string
		var callbackReason string
		var sourceInteraction string
		var assigneeAgent string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-scheduled-id",
			Short: "Update scheduled callback by Id",
			Long:  `Update an existing scheduled callback by Id,those whose scheduled trigger time has already passed cannot be updated. Requires 'cjp:user' scope for authorization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PUT", "/v1/callbacks/organization/{orgId}/scheduled-callback/{id}")
				req.PathParam("orgId", orgId)
				req.PathParam("id", id)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("id", id)
					req.BodyString("customerName", customerName)
					req.BodyString("callbackNumber", callbackNumber)
					req.BodyString("timezone", timezone)
					req.BodyString("scheduleDate", scheduleDate)
					req.BodyString("startTime", startTime)
					req.BodyString("endTime", endTime)
					req.BodyString("queueId", queueId)
					req.BodyString("callbackReason", callbackReason)
					req.BodyString("sourceInteraction", sourceInteraction)
					req.BodyString("assigneeAgent", assigneeAgent)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID for which the callback is being scheduled. This should be a valid UUID.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&id, "id", "", "The id with which the Scheduled Callback has been created.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&customerName, "customer-name", "", "")
		cmd.Flags().StringVar(&callbackNumber, "callback-number", "", "")
		cmd.Flags().StringVar(&timezone, "timezone", "", "")
		cmd.Flags().StringVar(&scheduleDate, "schedule-date", "", "")
		cmd.Flags().StringVar(&startTime, "start-time", "", "")
		cmd.Flags().StringVar(&endTime, "end-time", "", "")
		cmd.Flags().StringVar(&queueId, "queue-id", "", "")
		cmd.Flags().StringVar(&callbackReason, "callback-reason", "", "")
		cmd.Flags().StringVar(&sourceInteraction, "source-interaction", "", "")
		cmd.Flags().StringVar(&assigneeAgent, "assignee-agent", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callbacksCmd.AddCommand(cmd)
	}

	{ // delete-scheduled-id
		var orgId string
		var id string
		cmd := &cobra.Command{
			Use:   "delete-scheduled-id",
			Short: "Delete scheduled callback by Id",
			Long:  `Delete an existing scheduled callback by Id, those whose scheduled trigger time has already passed cannot be deleted. Requires 'cjp:user' scope for authorization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/v1/callbacks/organization/{orgId}/scheduled-callback/{id}")
				req.PathParam("orgId", orgId)
				req.PathParam("id", id)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization ID for which the callback is being scheduled. This should be a valid UUID.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&id, "id", "", "The id with which the Scheduled Callback has been created.")
		cmd.MarkFlagRequired("id")
		callbacksCmd.AddCommand(cmd)
	}

}
