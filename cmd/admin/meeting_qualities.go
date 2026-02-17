package admin

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

var meetingQualitiesCmd = &cobra.Command{
	Use:   "meeting-qualities",
	Short: "MeetingQualities commands",
}

func init() {
	cmd.AdminCmd.AddCommand(meetingQualitiesCmd)

	{ // get
		var meetingId string
		var max string
		var offset string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Meeting Qualities",
			Long:  "Get quality data for a meeting, by `meetingId`. Only organization administrators can retrieve meeting quality data.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meeting/qualities")
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the specific meeting instance. **Note:** The `meetingId` can be obtained via the Meeting List API when meetingType=meeting. The `id` attribute in the Meeting List Response is what is needed, for example, `e5dba9613a9d455aa49f6ffdafb6e7db_I_191395283063545470`.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of media sessions in the response.")
		cmd.Flags().StringVar(&offset, "offset", "", "Offset from the first result that you want to fetch.")
		meetingQualitiesCmd.AddCommand(cmd)
	}

}
