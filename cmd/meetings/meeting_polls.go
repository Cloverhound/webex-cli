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

var meetingPollsCmd = &cobra.Command{
	Use:   "meeting-polls",
	Short: "MeetingPolls commands",
}

func init() {
	cmd.MeetingsCmd.AddCommand(meetingPollsCmd)

	{ // list
		var meetingId string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Meeting Polls",
			Long:  "Lists all the polls and the poll questions in a meeting when ready.\n\n* Only [meeting instances](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) in state `ended` or `inProgress` are supported for `meetingId`.\n\n* No pagination for this API because we don't expect a large number of questions for each meeting.\n\n<div><Callout type=\"info\">Polls are available within 15 minutes following the meeting.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/polls")
				req.QueryParam("meetingId", meetingId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "A unique identifier for the [meeting instance](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) to which the polls belong.")
		meetingPollsCmd.AddCommand(cmd)
	}

	{ // get-pollresults
		var meetingId string
		var max string
		cmd := &cobra.Command{
			Use:   "get-pollresults",
			Short: "Get Meeting PollResults",
			Long:  "List the meeting polls, the poll's questions, and answers from the meeting when ready.\n\n* Only [meeting instances](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) in state `ended` or `inProgress` are supported for `meetingId`.\n\n* Long result sets will be split into [pages](/docs/basics#pagination).\n\n* This API is paginated by the sum of respondents from all questions in a meeting, these pagination links are returned in the response header.\n\n<div><Callout type=\"info\">Polls results are available within 15 minutes following the meeting.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/pollResults")
				req.QueryParam("meetingId", meetingId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "A unique identifier for the [meeting instance](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) to which the polls belong.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of respondents in a meeting in the response, up to 100.")
		meetingPollsCmd.AddCommand(cmd)
	}

	{ // list-respondents-question
		var pollId string
		var questionId string
		var meetingId string
		var max string
		cmd := &cobra.Command{
			Use:   "list-respondents-question",
			Short: "List Respondents of a Question",
			Long:  "Lists the respondents to a specific questions in a poll.\n\n* Only [meeting instances](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) in state `ended` or `inProgress` are supported for `meetingId`.\n\n* Long result sets are split into [pages](/docs/basics#pagination).\n\n<div><Callout type=\"info\">The list of poll respondents are available within 15 minutes following the meeting.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/polls/{pollId}/questions/{questionId}/respondents")
				req.PathParam("pollId", pollId)
				req.PathParam("questionId", questionId)
				req.QueryParam("meetingId", meetingId)
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
		cmd.Flags().StringVar(&pollId, "poll-id", "", "A unique identifier for the poll to which the respondents belong.")
		cmd.MarkFlagRequired("poll-id")
		cmd.Flags().StringVar(&questionId, "question-id", "", "A unique identifier for the question to which the respondents belong.")
		cmd.MarkFlagRequired("question-id")
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "A unique identifier for the [meeting instance](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) to which the respondents belong.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of respondents in a specified question in the response, up to 100.")
		meetingPollsCmd.AddCommand(cmd)
	}

}
