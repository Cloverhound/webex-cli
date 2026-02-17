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

var chatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Chats commands",
}

func init() {
	cmd.MeetingsCmd.AddCommand(chatsCmd)

	{ // list-meeting
		var meetingId string
		var max string
		var offset string
		cmd := &cobra.Command{
			Use:   "list-meeting",
			Short: "List Meeting Chats",
			Long:  "Lists the meeting chats of a finished [meeting instance](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) specified by `meetingId`. You can set a maximum number of chats to return.\n\nUse this operation to list the chats of a finished [meeting instance](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) when they are ready. Please note that only **meeting instances** in state `ended` are supported for `meetingId`. **Meeting series**, **scheduled meetings** and `in-progress` **meeting instances** are not supported.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/postMeetingChats")
				req.QueryParam("meetingId", meetingId)
				req.QueryParam("max", max)
				req.QueryParam("offset", offset)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "A unique identifier for the [meeting instance](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) to which the chats belong. The meeting ID of a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting is not supported.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of meeting chats in the response, up to 100.")
		cmd.Flags().StringVar(&offset, "offset", "", "Offset from the first result that you want to fetch.")
		chatsCmd.AddCommand(cmd)
	}

	{ // delete-meeting
		var meetingId string
		cmd := &cobra.Command{
			Use:   "delete-meeting",
			Short: "Delete Meeting Chats",
			Long:  "Deletes the meeting chats of a finished [meeting instance](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) specified by `meetingId`.\n\nUse this operation to delete the chats of a finished [meeting instance](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) when they are ready. Please note that only **meeting instances** in state `ended` are supported for `meetingId`. **Meeting series**, **scheduled meetings** and `in-progress` **meeting instances** are not supported.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/meetings/postMeetingChats")
				req.QueryParam("meetingId", meetingId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "A unique identifier for the [meeting instance](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) to which the chats belong. Meeting IDs of a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting are not supported.")
		chatsCmd.AddCommand(cmd)
	}

}
