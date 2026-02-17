package meetings

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

var meetingSummaryCmd = &cobra.Command{
	Use:   "meeting-summary",
	Short: "MeetingSummary commands",
}

func init() {
	cmd.MeetingsCmd.AddCommand(meetingSummaryCmd)

	{ // list-usage-reports
		var siteUrl string
		var serviceType string
		var from string
		var to string
		var max string
		var timezone string
		cmd := &cobra.Command{
			Use:   "list-usage-reports",
			Short: "List Meeting Usage Reports",
			Long:  "List meeting usage reports of all the users on the specified site by an admin. You can specify a date range and the maximum number of meeting usage reports to return.\n\nThe list returned is sorted in descending order by the date and time the meetings were started.\n\nLong result sets are split into [pages](/docs/basics#pagination).\n\n* `siteUrl` is required, and the meeting usage reports of the specified site are listed. All available Webex sites can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\n#### Request Header\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default timezone is `UTC` if not defined.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetingReports/usage")
				req.QueryParam("siteUrl", siteUrl)
				req.QueryParam("serviceType", serviceType)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("max", max)
				req.Header("timezone", timezone)
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
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API lists meeting usage reports from. All available Webex sites can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		cmd.Flags().StringVar(&serviceType, "service-type", "", "Meeting usage report's service-type. If `serviceType` is specified, the API filters meeting usage reports by service-type. If `serviceType` is not specified, the API returns meeting usage reports by `MeetingCenter` by default. Valid values:  + `MeetingCenter`  + `EventCenter`  + `SupportCenter`  + `TrainingCenter`")
		cmd.Flags().StringVar(&from, "from", "", "Starting date and time for meeting usage reports to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `from` cannot be after `to`. The interval between `to` and `from` cannot exceed 30 days and `from` cannot be earlier than 90 days ago.")
		cmd.Flags().StringVar(&to, "to", "", "Ending date and time for meeting usage reports to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `to` cannot be before `from`. The interval between `to` and `from` cannot exceed 30 days.")
		cmd.Flags().StringVar(&max, "max", "", "Maximum number of meetings to include in the meetings usage report in a single page. `max` must be greater than 0 and equal to or less than `1000`.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. Asia/Shanghai")
		meetingSummaryCmd.AddCommand(cmd)
	}

	{ // list-attendee-reports
		var siteUrl string
		var from string
		var to string
		var max string
		var meetingId string
		var meetingNumber string
		var meetingTitle string
		var timezone string
		cmd := &cobra.Command{
			Use:   "list-attendee-reports",
			Short: "List Meeting Attendee Reports",
			Long:  "Lists of meeting attendee reports by a date range, the maximum number of meeting attendee reports, a meeting ID, a meeting number or a meeting title.\n\nIf the requesting user is an admin, the API returns meeting attendee reports of the meetings hosted by all the users on the specified site filtered by meeting ID, meeting number or meeting title.\n\nIf it's a normal meeting host, the API returns meeting attendee reports of the meetings hosted by the user himself on the specified site filtered by meeting ID, meeting number or meeting title.\n\nThe list returned is grouped by meeting instances. Both the groups and items of each group are sorted in descending order of `joinedTime`. For example, if `meetingId` is specified and it's a meeting series ID, the returned list is grouped by meeting instances of that series. The groups are sorted in descending order of `joinedTime`, and within each group the items are also sorted in descending order of `joinedTime`. Please refer to [Meetings Overview](/docs/meetings) for details of meeting series, scheduled meeting and meeting instance.\n\nLong result sets are split into [pages](/docs/basics#pagination).\n\n* `siteUrl` is required, and the meeting attendee reports of the specified site are listed. All available Webex sites can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\n* `meetingId`, `meetingNumber` and `meetingTitle` are optional parameters to query the meeting attendee reports, but at least one of them should be specified. If more than one parameter in the sequence of `meetingId`, `meetingNumber`, and `meetingTitle` are specified, the first one in the sequence is used. Currently, only ended meeting instance IDs and meeting series IDs are supported for `meetingId`. IDs of scheduled meetings or personal room meetings are not supported.\n\n#### Request Header\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default timezone is `UTC` if not defined.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetingReports/attendees")
				req.QueryParam("siteUrl", siteUrl)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("max", max)
				req.QueryParam("meetingId", meetingId)
				req.QueryParam("meetingNumber", meetingNumber)
				req.QueryParam("meetingTitle", meetingTitle)
				req.Header("timezone", timezone)
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
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API lists meeting attendee reports from. All available Webex sites can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		cmd.Flags().StringVar(&from, "from", "", "Starting date and time for the meeting attendee reports to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `from` cannot be after `to`. The interval between `to` and `from` cannot exceed 30 days and `from` cannot be earlier than 90 days ago.")
		cmd.Flags().StringVar(&to, "to", "", "Ending date and time for the meeting attendee reports to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `to` cannot be before `from`. The interval between `to` and `from` cannot exceed 30 days.")
		cmd.Flags().StringVar(&max, "max", "", "Maximum number of meeting attendees to include in the meeting attendee report in a single page. `max` must be greater than 0 and equal to or less than `1000`.")
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Meeting ID for the meeting attendee reports to return. If specified, return meeting attendee reports of the specified meeting; otherwise, return meeting attendee reports of all meetings. Currently, only ended meeting instance IDs are supported. IDs of meeting series, scheduled meetings or personal room meetings are not supported.")
		cmd.Flags().StringVar(&meetingNumber, "meeting-number", "", "Meeting number for the meeting attendee reports to return. If specified, return meeting attendee reports of the specified meeting; otherwise, return meeting attendee reports of all meetings.")
		cmd.Flags().StringVar(&meetingTitle, "meeting-title", "", "Meeting title for the meeting attendee reports to return. If specified, return meeting attendee reports of the specified meeting; otherwise, return meeting attendee reports of all meetings.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. Asia/Shanghai")
		meetingSummaryCmd.AddCommand(cmd)
	}

}
