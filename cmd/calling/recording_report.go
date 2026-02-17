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

var recordingReportCmd = &cobra.Command{
	Use:   "recording-report",
	Short: "RecordingReport commands",
}

func init() {
	cmd.CallingCmd.AddCommand(recordingReportCmd)

	{ // list-audit-summaries
		var max string
		var from string
		var to string
		var hostEmail string
		var siteUrl string
		var timezone string
		cmd := &cobra.Command{
			Use:   "list-audit-summaries",
			Short: "List of Recording Audit Report Summaries",
			Long:  "Lists of recording audit report summaries. You can specify a date range and the maximum number of recording audit report summaries to return.\n\nOnly recording audit report summaries of meetings hosted by or shared with the authenticated user will be listed.\n\nThe list returned is sorted in descending order by the date and time that the recordings were created.\n\nLong result sets are split into [pages](/docs/basics#pagination).\n\n* If `siteUrl` is specified, the recording audit report summaries of the specified site will be listed; otherwise, recording audit report summaries of the user's preferred site will be listed. All available Webex sites and the preferred site of the user can be retrieved by the `Get Site List` API.\n\n#### Request Header\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/recordingReport/accessSummary")
				req.QueryParam("max", max)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("siteUrl", siteUrl)
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
		cmd.Flags().StringVar(&max, "max", "", "Maximum number of recording audit report summaries to return in a single page. `max` must be equal to or greater than `1` and equal to or less than `100`.")
		cmd.Flags().StringVar(&from, "from", "", "Starting date and time (inclusive) for recording audit report summaries to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `from` cannot be after `to`. Please note that the interval between `to` and `from` cannot exceed 90 days and the interval between the current time and `from` cannot exceed 365 days.")
		cmd.Flags().StringVar(&to, "to", "", "Ending date and time (exclusive) for recording audit report summaries to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `to` cannot be before `from`. Please note that the interval between `to` and `from` cannot exceed 90 days and the interval between the current time and `from` cannot exceed 365 days.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will return recording audit report summaries of that user. If a special value of `all` is set for `hostEmail`, the admin can list recording audit report summaries of all users on the target site, not of a single user.")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API lists recording audit report summaries from. If not specified, the API lists summary audit report for recordings from the user's preferred site. All available Webex sites and the preferred site of the user can be retrieved by `Get Site List` API.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		recordingReportCmd.AddCommand(cmd)
	}

	{ // get-audit
		var recordingId string
		var hostEmail string
		var max string
		var timezone string
		cmd := &cobra.Command{
			Use:   "get-audit",
			Short: "Get Recording Audit Report Details",
			Long:  "Retrieves details for a recording audit report with a specified recording ID.\n\nOnly recording audit report details of meetings hosted by or shared with the authenticated user may be retrieved.\n\n#### Request Header\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/recordingReport/accessDetail")
				req.QueryParam("recordingId", recordingId)
				req.QueryParam("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&recordingId, "recording-id", "", "A unique identifier for the recording.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will return recording details of that user.")
		cmd.Flags().StringVar(&max, "max", "", "Maximum number of recording audit report details to return in a single page. `max` must be equal to or greater than `1` and equal to or less than `100`.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		recordingReportCmd.AddCommand(cmd)
	}

	{ // list-meeting-archive-summaries
		var max string
		var from string
		var to string
		var siteUrl string
		var timezone string
		cmd := &cobra.Command{
			Use:   "list-meeting-archive-summaries",
			Short: "List Meeting Archive Summaries",
			Long:  "Lists of meeting archive summaries. You can specify a date range and the maximum number of meeting archive summaries to return.\n\nMeeting archive summaries are only available to full administrators, not even the meeting host.\n\nThe list returned is sorted in descending order by the date and time that the archives were created.\n\nLong result sets are split into [pages](/docs/basics#pagination).\n\n* If `siteUrl` is specified, the meeting archive summaries of the specified site will be listed; otherwise, meeting archive summaries of the user's preferred site will be listed. All available Webex sites and the preferred site of the user can be retrieved by the `Get Site List` API.\n\n#### Request Header\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/recordingReport/meetingArchiveSummaries")
				req.QueryParam("max", max)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("siteUrl", siteUrl)
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
		cmd.Flags().StringVar(&max, "max", "", "Maximum number of meeting archive summaries to return in a single page. `max` must be equal to or greater than `1` and equal to or less than `100`.")
		cmd.Flags().StringVar(&from, "from", "", "Starting date and time (inclusive) for meeting archive summaries to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `from` cannot be after `to`. Please note that the interval between `to` and `from` cannot exceed 30 days.")
		cmd.Flags().StringVar(&to, "to", "", "Ending date and time (exclusive) for meeting archive summaries to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `to` cannot be before `from`. Please note that the interval between `to` and `from` cannot exceed 30 days.")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API lists meeting archive summaries from. If not specified, the API lists meeting archive summaries for recordings from the user's preferred site. All available Webex sites and the preferred site of the user can be retrieved by `Get Site List` API.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		recordingReportCmd.AddCommand(cmd)
	}

	{ // get-meeting-archive
		var archiveId string
		var timezone string
		cmd := &cobra.Command{
			Use:   "get-meeting-archive",
			Short: "Get Meeting Archive Details",
			Long:  "Retrieves details for a meeting archive report with a specified archive ID, which contains recording metadata.\n\nMeeting archive details are only available to full administrators, not even the meeting host.\n\n#### Request Header\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/recordingReport/meetingArchives/{archiveId}")
				req.PathParam("archiveId", archiveId)
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
		cmd.Flags().StringVar(&archiveId, "archive-id", "", "A unique identifier for the meeting archive summary.")
		cmd.MarkFlagRequired("archive-id")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		recordingReportCmd.AddCommand(cmd)
	}

}
