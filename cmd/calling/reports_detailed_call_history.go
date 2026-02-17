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

var reportsDetailedCallHistoryCmd = &cobra.Command{
	Use:   "reports-detailed-call-history",
	Short: "ReportsDetailedCallHistory commands",
}

func init() {
	cmd.CallingCmd.AddCommand(reportsDetailedCallHistoryCmd)

	{ // get
		var startTime string
		var endTime string
		var locations string
		var max string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Detailed Call History",
			Long:  "Provides Webex Calling Detailed Call History data for your organization.\n\nResults can be filtered with the `startTime`, `endTime` and `locations` request parameters. The `startTime` and `endTime` parameters specify the start and end of the time period for the Detailed Call History reports you wish to collect. The API will return all reports that were created between `startTime` and `endTime`.\n\n<br/><br/>\nResponse entries may be added as more information is made available for the reports.\nValues in response items may be extended as more capabilities are added to Webex Calling.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/cdr_feed")
				req.QueryParam("startTime", startTime)
				req.QueryParam("endTime", endTime)
				req.QueryParam("locations", locations)
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
		cmd.Flags().StringVar(&startTime, "start-time", "", "Time of the first report you wish to collect. (Report time is the time the call finished). **Note:** The specified time must be between 5 minutes ago and 48 hours ago, and formatted as `YYYY-MM-DDTHH:MM:SS.mmmZ`.")
		cmd.Flags().StringVar(&endTime, "end-time", "", "Time of the last report you wish to collect. (Report time is the time the call finished). **Note:** The specified time should be later than `startTime` but no later than 48 hours, and formatted as `YYYY-MM-DDTHH:MM:SS.mmmZ`.")
		cmd.Flags().StringVar(&locations, "locations", "", "Name of the location (as shown in Control Hub). Up to 10 comma-separated locations can be provided. Allows you to query reports by location.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of reports per page of the response. The range is 500 to 5000. Values below 500 are automatically adjusted up to 500, and values above 5000 are automatically adjusted down to 5000. When the API has more reports to return than the max value, the API response will be paginated. Follow the next link contained in the “Link” header within a response to request the next page of results. If there is no next link, all reports for the selected time range have been collected.  For instance, let's say the initial API request is  https://analytics-calling.webexapis.com/v1/cdr_feed?endTime=2025-08-15T10:00:00.000Z&startTime=2025-08-15T08:00:00.000Z&max=5000  The link header in the response would look something like  <<https://analytics-calling.webexapis.com/v1/cdr_feed?endTime=2025-08-15T10:00:00.000Z&startTime=2025-08-15T08:00:00.000Z&startTimeForNextFetch=2025-08-15T09:30:00.000Z&totalCount=20000&max=5000&orgId=zzzzzzzz-yyyy-zzzz-xxxx-yyyyyyyyyyyy>;rel=\"next\">")
		reportsDetailedCallHistoryCmd.AddCommand(cmd)
	}

	{ // get-live-stream
		var startTime string
		var endTime string
		var locations string
		var max string
		cmd := &cobra.Command{
			Use:   "get-live-stream",
			Short: "Get Live Stream Detailed Call History",
			Long:  "Provides Webex Calling Detailed Call History data for your organization.\n\nResults can be filtered with the `startTime`, `endTime` and `locations` request parameters. The `startTime` and `endTime` parameters specify the time window during which Detailed Call History data was inserted into the Webex Calling cloud. The API will return all reports whose insertion time into the Webex Calling cloud falls between `startTime` and `endTime`.\n\n<br/><br/>\nResponse entries may be added as more information is made available for the reports.\nValues in response items may be extended as more capabilities are added to Webex Calling.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/cdr_stream")
				req.QueryParam("startTime", startTime)
				req.QueryParam("endTime", endTime)
				req.QueryParam("locations", locations)
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
		cmd.Flags().StringVar(&startTime, "start-time", "", "The start date-time of the first record you wish to collect in UTC time. It would be the earliest time at which the data was inserted into the Webex Calling cloud for the records you wish to collect. Format must be as `YYYY-MM-DDTHH:MM:SS.mmmZ`. `startTime` can't be older than 12 hours from your current UTC time. The window period between `startTime` and `endTime` must not exceed 2 hours in a single API request.")
		cmd.Flags().StringVar(&endTime, "end-time", "", "The end date-time of the last record you wish to collect in UTC time. It would be the latest time at which the data was inserted into the Webex Calling cloud for the records you wish to collect. Format must be as `YYYY-MM-DDTHH:MM:SS.mmmZ`. `endTime` must be 1 minute ago from your current UTC time and can’t be older than 12 hours. `endTime` must be greater than `startTime`. The window period between `startTime` and `endTime` must not exceed 2 hours in a single API request.")
		cmd.Flags().StringVar(&locations, "locations", "", "Name of the location (as shown in Control Hub). Up to 10 comma-separated locations can be provided. Allows you to query reports by location.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of reports per page of the response. The range is 500 to 5000. Values below 500 are automatically adjusted up to 500, and values above 5000 are automatically adjusted down to 5000. When the API has more reports to return than the max value, the API response will be paginated. Follow the next link contained in the “Link” header within a response to request the next page of results. If there is no next link, all reports for the selected time range have been collected.  For instance, let's say the initial API request is  https://analytics-calling.webexapis.com/v1/cdr_stream?endTime=2025-08-15T10:00:00.000Z&startTime=2025-08-15T08:00:00.000Z&max=5000  The link header in the response would look something like  <<https://analytics-calling.webexapis.com/v1/cdr_stream?endTime=2025-08-15T10:00:00.000Z&startTime=2025-08-15T08:00:00.000Z&startTimeForNextFetch=2025-08-15T09:30:00.000Z&totalCount=20000&max=5000&orgId=zzzzzzzz-yyyy-zzzz-xxxx-yyyyyyyyyyyy>;rel=\"next\">")
		reportsDetailedCallHistoryCmd.AddCommand(cmd)
	}

}
