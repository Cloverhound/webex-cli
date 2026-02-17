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

var inviteesCmd = &cobra.Command{
	Use:   "invitees",
	Short: "Invitees commands",
}

func init() {
	cmd.MeetingsCmd.AddCommand(inviteesCmd)

	{ // list-meeting
		var meetingId string
		var max string
		var hostEmail string
		var panelist string
		cmd := &cobra.Command{
			Use:   "list-meeting",
			Short: "List Meeting Invitees",
			Long:  "Lists meeting invitees for a meeting with a specified `meetingId`. You can set a maximum number of invitees to return.\n\nThis operation can be used for meeting series, scheduled meetings, and ended or ongoing meeting instance objects. If the specified `meetingId` is for a meeting series, the invitees for the series will be listed; if the `meetingId` is for a scheduled meeting, the invitees for the particular scheduled meeting will be listed; if the `meetingId` is for an ended or ongoing meeting instance, the invitees for the particular meeting instance will be listed. See the [Webex Meetings](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) guide for more information about the types of meetings.\n\nThe list returned is sorted in ascending order by email address.\n\nLong result sets are split into [pages](/docs/basics#pagination).",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetingInvitees")
				req.QueryParam("meetingId", meetingId)
				req.QueryParam("max", max)
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("panelist", panelist)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting for which invitees are being requested. The meeting can be a meeting series, a scheduled meeting, or a meeting instance which has ended or is ongoing. The meeting ID of a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting is not supported for this API.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of meeting invitees in the response, up to 100.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will return meeting invitees that are hosted by that user.")
		cmd.Flags().StringVar(&panelist, "panelist", "", "Filter invitees or attendees for webinars only. If `true`, returns invitees. If `false`, returns attendees. If `null`, returns both invitees and attendees.")
		inviteesCmd.AddCommand(cmd)
	}

	{ // create-meeting
		var meetingId string
		var email string
		var displayName string
		var coHost bool
		var hostEmail string
		var sendEmail bool
		var panelist bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-meeting",
			Short: "Create a Meeting Invitee",
			Long:  "* Invite a person to attend a meeting.\n\n* Identify the invitee in the request body, by email address.\n\n* The `sendEmail` parameter is `true` by default and the meeting emails will be sent to the invitee's `email`. Please set `sendEmail` to `false` to prevent the invitee from receiving emails.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetingInvitees")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("meetingId", meetingId)
					req.BodyString("email", email)
					req.BodyString("displayName", displayName)
					req.BodyBool("coHost", coHost, cmd.Flags().Changed("co-host"))
					req.BodyString("hostEmail", hostEmail)
					req.BodyBool("sendEmail", sendEmail, cmd.Flags().Changed("send-email"))
					req.BodyBool("panelist", panelist, cmd.Flags().Changed("panelist"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "")
		cmd.Flags().StringVar(&email, "email", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().BoolVar(&coHost, "co-host", false, "")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "")
		cmd.Flags().BoolVar(&sendEmail, "send-email", false, "")
		cmd.Flags().BoolVar(&panelist, "panelist", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		inviteesCmd.AddCommand(cmd)
	}

	{ // create-meeting-2
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-meeting-2",
			Short: "Create Meeting Invitees",
			Long:  "* Invite people to attend a meeting in bulk.\n\n* Identify each invitee by the email address of each item in the `items` of the request body.\n\n* Each invitee should have a unique `email`.\n\n* This API limits the maximum size of `items` in the request body to 100.\n\n* The `sendEmail` parameter for each invitee is `true` by default and the meeting emails will be sent to the invitee's `email`. Please set `sendEmail` to `false` to prevent an invitee from receiving emails.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetingInvitees/bulkInsert")
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
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		inviteesCmd.AddCommand(cmd)
	}

	{ // get-meeting
		var meetingInviteeId string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "get-meeting",
			Short: "Get a Meeting Invitee",
			Long:  "Retrieve details for a meeting invitee identified by a `meetingInviteeId` in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetingInvitees/{meetingInviteeId}")
				req.PathParam("meetingInviteeId", meetingInviteeId)
				req.QueryParam("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingInviteeId, "meeting-invitee-id", "", "Unique identifier for the invitee whose details are being requested.")
		cmd.MarkFlagRequired("meeting-invitee-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting invitee that is hosted by that user.")
		inviteesCmd.AddCommand(cmd)
	}

	{ // update-meeting
		var meetingInviteeId string
		var email string
		var displayName string
		var coHost bool
		var hostEmail string
		var sendEmail bool
		var panelist bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-meeting",
			Short: "Update a Meeting Invitee",
			Long:  "Update details for a meeting invitee identified by a `meetingInviteeId` in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/meetingInvitees/{meetingInviteeId}")
				req.PathParam("meetingInviteeId", meetingInviteeId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("email", email)
					req.BodyString("displayName", displayName)
					req.BodyBool("coHost", coHost, cmd.Flags().Changed("co-host"))
					req.BodyString("hostEmail", hostEmail)
					req.BodyBool("sendEmail", sendEmail, cmd.Flags().Changed("send-email"))
					req.BodyBool("panelist", panelist, cmd.Flags().Changed("panelist"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingInviteeId, "meeting-invitee-id", "", "Unique identifier for the invitee to be updated. This parameter only applies to an invitee to a meeting series or a scheduled meeting. It doesn't apply to an invitee to an ended or ongoing meeting instance.")
		cmd.MarkFlagRequired("meeting-invitee-id")
		cmd.Flags().StringVar(&email, "email", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().BoolVar(&coHost, "co-host", false, "")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "")
		cmd.Flags().BoolVar(&sendEmail, "send-email", false, "")
		cmd.Flags().BoolVar(&panelist, "panelist", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		inviteesCmd.AddCommand(cmd)
	}

	{ // delete-meeting
		var meetingInviteeId string
		var hostEmail string
		var sendEmail string
		cmd := &cobra.Command{
			Use:   "delete-meeting",
			Short: "Delete a Meeting Invitee",
			Long:  "Removes a meeting invitee identified by a `meetingInviteeId` specified in the URI. The deleted meeting invitee cannot be recovered.\n\nIf the meeting invitee is associated with a meeting series, the invitee will be removed from the entire meeting series. If the invitee is associated with a scheduled meeting, the invitee will be removed from only that scheduled meeting.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/meetingInvitees/{meetingInviteeId}")
				req.PathParam("meetingInviteeId", meetingInviteeId)
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("sendEmail", sendEmail)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingInviteeId, "meeting-invitee-id", "", "Unique identifier for the invitee to be removed. This parameter only applies to an invitee to a meeting series or a scheduled meeting. It doesn't apply to an invitee to an ended or ongoing meeting instance.")
		cmd.MarkFlagRequired("meeting-invitee-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will delete a meeting invitee that is hosted by that user.")
		cmd.Flags().StringVar(&sendEmail, "send-email", "", "If `true`, send an email to the invitee.")
		inviteesCmd.AddCommand(cmd)
	}

}
