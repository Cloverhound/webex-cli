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

var locationSchedulesCmd = &cobra.Command{
	Use:   "location-schedules",
	Short: "LocationSchedules commands",
}

func init() {
	cmd.CallingCmd.AddCommand(locationSchedulesCmd)

	{ // list
		var locationId string
		var orgId string
		var max string
		var start string
		var name string
		var typeVal string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "Read the List of Schedules",
			Long:  "List all schedules for the given location of the organization.\n\nA time schedule establishes a set of times during the day or holidays in the year in which a feature, for example auto attendants, can perform a specific action.\n\nRetrieving this list requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/schedules")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("name", name)
				req.QueryParam("type", typeVal)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the list of schedules for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List schedules for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&name, "name", "", "Only return schedules with the matching name.")
		cmd.Flags().StringVar(&typeVal, "type", "", "Type of the schedule.")
		locationSchedulesCmd.AddCommand(cmd)
	}

	{ // create
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Schedule",
			Long:  "Create new Schedule for the given location.\n\nA time schedule establishes a set of times during the day or holidays in the year in which a feature, for example auto attendants, can perform a specific action.\n\nCreating a schedule requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/schedules")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Create the schedule for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create the schedule for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationSchedulesCmd.AddCommand(cmd)
	}

	{ // get
		var locationId string
		var typeVal string
		var scheduleId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Details for a Schedule",
			Long:  "Retrieve Schedule details.\n\nA time schedule establishes a set of times during the day or holidays in the year in which a feature, for example auto attendants, can perform a specific action.\n\nRetrieving schedule details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/schedules/{type}/{scheduleId}")
				req.PathParam("locationId", locationId)
				req.PathParam("type", typeVal)
				req.PathParam("scheduleId", scheduleId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve schedule details in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&typeVal, "type", "", "Type of the schedule.")
		cmd.MarkFlagRequired("type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Retrieve the schedule with the matching ID.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve schedule details from this organization.")
		locationSchedulesCmd.AddCommand(cmd)
	}

	{ // update
		var locationId string
		var typeVal string
		var scheduleId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Schedule",
			Long:  "Update the designated schedule.\n\nA time schedule establishes a set of times during the day or holidays in the year in which a feature, for example auto attendants, can perform a specific action.\n\nUpdating a schedule requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The Schedule ID will change upon modification of the Schedule name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/schedules/{type}/{scheduleId}")
				req.PathParam("locationId", locationId)
				req.PathParam("type", typeVal)
				req.PathParam("scheduleId", scheduleId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this schedule exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&typeVal, "type", "", "Type of schedule.")
		cmd.MarkFlagRequired("type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Update schedule with the matching ID.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update schedule from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationSchedulesCmd.AddCommand(cmd)
	}

	{ // delete
		var locationId string
		var typeVal string
		var scheduleId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Schedule",
			Long:  "Delete the designated Schedule.\n\nA time schedule establishes a set of times during the day or holidays in the year in which a feature, for example auto attendants, can perform a specific action.\n\nDeleting a schedule requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/schedules/{type}/{scheduleId}")
				req.PathParam("locationId", locationId)
				req.PathParam("type", typeVal)
				req.PathParam("scheduleId", scheduleId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location from which to delete a schedule.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&typeVal, "type", "", "Type of the schedule.")
		cmd.MarkFlagRequired("type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Delete the schedule with the matching ID.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the schedule from this organization.")
		locationSchedulesCmd.AddCommand(cmd)
	}

	{ // get-event
		var locationId string
		var typeVal string
		var scheduleId string
		var eventId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-event",
			Short: "Get Details for a Schedule Event",
			Long:  "Retrieve Schedule Event details.\n\nA time schedule establishes a set of times during the day or holidays in the year in which a feature, for example auto attendants, can perform a specific action.\n\nRetrieving a schedule event's details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/schedules/{type}/{scheduleId}/events/{eventId}")
				req.PathParam("locationId", locationId)
				req.PathParam("type", typeVal)
				req.PathParam("scheduleId", scheduleId)
				req.PathParam("eventId", eventId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve schedule event details in this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&typeVal, "type", "", "Type of schedule.")
		cmd.MarkFlagRequired("type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Retrieve the schedule event with the matching schedule ID.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&eventId, "event-id", "", "Retrieve the schedule event with the matching schedule event ID.")
		cmd.MarkFlagRequired("event-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve schedule event details from this organization.")
		locationSchedulesCmd.AddCommand(cmd)
	}

	{ // update-event
		var locationId string
		var typeVal string
		var scheduleId string
		var eventId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-event",
			Short: "Update a Schedule Event",
			Long:  "Update the designated Schedule Event.\n\nA time schedule establishes a set of times during the day or holidays in the year in which a feature, for example auto attendants, can perform a specific action.\n\nUpdating a schedule event requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: The schedule event ID will change upon modification of the schedule event name.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/schedules/{type}/{scheduleId}/events/{eventId}")
				req.PathParam("locationId", locationId)
				req.PathParam("type", typeVal)
				req.PathParam("scheduleId", scheduleId)
				req.PathParam("eventId", eventId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location in which this schedule event exists.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&typeVal, "type", "", "Type of schedule.  + `businessHours` - Business hours schedule type.  + `holidays` - Holidays schedule type.")
		cmd.MarkFlagRequired("type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Update schedule event with the matching schedule ID.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&eventId, "event-id", "", "Update the schedule event with the matching schedule event ID.")
		cmd.MarkFlagRequired("event-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update schedule from this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationSchedulesCmd.AddCommand(cmd)
	}

	{ // delete-event
		var locationId string
		var typeVal string
		var scheduleId string
		var eventId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-event",
			Short: "Delete a Schedule Event",
			Long:  "Delete the designated Schedule Event.\n\nA time schedule establishes a set of times during the day or holidays in the year in which a feature, for example auto attendants, can perform a specific action.\n\nDeleting a schedule event requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/schedules/{type}/{scheduleId}/events/{eventId}")
				req.PathParam("locationId", locationId)
				req.PathParam("type", typeVal)
				req.PathParam("scheduleId", scheduleId)
				req.PathParam("eventId", eventId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location from which to delete a schedule.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&typeVal, "type", "", "Type of schedule.")
		cmd.MarkFlagRequired("type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Delete the schedule with the matching ID.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&eventId, "event-id", "", "Delete the schedule event with the matching schedule event ID.")
		cmd.MarkFlagRequired("event-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the schedule from this organization.")
		locationSchedulesCmd.AddCommand(cmd)
	}

	{ // create-event
		var locationId string
		var typeVal string
		var scheduleId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-event",
			Short: "Create a Schedule Event",
			Long:  "Create new Event for the given location Schedule.\n\nA time schedule establishes a set of times during the day or holidays in the year in which a feature, for example auto attendants, can perform a specific action.\n\nCreating a schedule event requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/schedules/{type}/{scheduleId}/events")
				req.PathParam("locationId", locationId)
				req.PathParam("type", typeVal)
				req.PathParam("scheduleId", scheduleId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Create the schedule for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&typeVal, "type", "", "Type of schedule.")
		cmd.MarkFlagRequired("type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Create event for a given schedule ID.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create the schedule for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationSchedulesCmd.AddCommand(cmd)
	}

}
