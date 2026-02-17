package meetings

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

var meetingsCmd = &cobra.Command{
	Use:   "meetings",
	Short: "Meetings commands",
}

func init() {
	cmd.MeetingsCmd.AddCommand(meetingsCmd)

	{ // get-admin
		var meetingId string
		var current string
		var timezone string
		cmd := &cobra.Command{
			Use:   "get-admin",
			Short: "Get a Meeting By an Admin",
			Long:  "Retrieves details for a meeting with a specified meeting ID by an admin. The following sensitive attributes are hidden from the response: `agenda`, `hostKey`, `password`, `phoneAndVideoSystemPassword`, `panelistPassword`, `phoneAndVideoSystemPanelistPassword`, `webLink`, `sipAddress`, `dialInIpAddress`, `registration` and `telephony.accessCode`.\n\n* If the `meetingId` value specified is for a meeting series and `current` is `true`, the operation returns details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series.\n\n* If the `meetingId` value specified is for a meeting series and `current` is `false` or `current` is not specified, the operation returns details for the entire meeting series.\n\n* If the `meetingId` value specified is for a scheduled meeting from a meeting series, the operation returns details for that scheduled meeting.\n\n* If the `meetingId` value specified is for a meeting instance which is happening or has happened, the operation returns details for that meeting instance.\n\n* `trackingCodes` is not supported for ended meeting instances.\n\n* To learn more about which attributes are available for different meeting states, please refer to [Available Meeting Attributes for Different Meeting States](/docs/meetings#available-meeting-attributes-for-different-meeting-states).\n\n#### Request Header\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) for time stamps in response body, defined in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default value is `UTC` if not specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/admin/meetings/{meetingId}")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("current", current)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting being requested.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series. The default value is `false`.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		meetingsCmd.AddCommand(cmd)
	}

	{ // list-admin
		var meetingNumber string
		var webLink string
		var current string
		var timezone string
		cmd := &cobra.Command{
			Use:   "list-admin",
			Short: "List Meetings By an Admin",
			Long:  "Retrieves details for meetings by an admin with a specified meeting number or web link. Please note that there are various products in the [Webex Suite](https://www.webex.com/collaboration-suite.html) such as `Meetings` and `Events`. Currently, only meetings of the `Meetings` product are supported by this API, meetings of others in the suite are not supported. Ad-hoc meetings created by [Create a Meeting](/docs/api/v1/meetings/create-a-meeting) with `adhoc` of `true` and a `roomId` will not be listed, but the ended and ongoing ad-hoc meeting instances will be listed. The following sensitive attributes are hidden from the response: `agenda`, `hostKey`, `password`, `phoneAndVideoSystemPassword`, `panelistPassword`, `phoneAndVideoSystemPanelistPassword`, `webLink`, `sipAddress`, `dialInIpAddress`, `registration` and `telephony.accessCode`.\n\n* If `meetingNumber` is specified, the operation returns an array of meeting objects specified by the `meetingNumber`. Each object in the array can be a scheduled meeting or a meeting series depending on whether the `current` parameter is `true` or `false`, and each object contains the simultaneous interpretation object. Please note that `meetingNumber` and `webLink` are mutually exclusive and they cannot be specified simultaneously.\n\n* If `webLink` is specified, the operation returns an array of meeting objects specified by the `webLink`. Each object in the array is a scheduled meeting, and each object contains the simultaneous interpretation object. Please note that `meetingNumber` and `webLink` are mutually exclusive and they cannot be specified simultaneously.\n\n* The `current` parameter only applies to meeting series. If it's `false`, the `start` and `end` attributes of each returned meeting series object are for the first scheduled meeting of that series. If it's `true` or not specified, the `start` and `end` attributes are for the scheduled meeting which is ready to start or join or the upcoming scheduled meeting of that series.\n\n* To learn more about which attributes are available for different meeting states, please refer to [Available Meeting Attributes for Different Meeting States](/docs/meetings#available-meeting-attributes-for-different-meeting-states).\n\n#### Request Header\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) for time stamps in response body, defined in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default value is `UTC` if not specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/admin/meetings")
				req.QueryParam("meetingNumber", meetingNumber)
				req.QueryParam("webLink", webLink)
				req.QueryParam("current", current)
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
		cmd.Flags().StringVar(&meetingNumber, "meeting-number", "", "Meeting number for the meeting objects being requested. `meetingNumber` and `webLink` are mutually exclusive. If it's an exceptional meeting from a meeting series, the exceptional meeting instead of the primary meeting series is returned.")
		cmd.Flags().StringVar(&webLink, "web-link", "", "URL encoded link to information page for the meeting objects being requested. `meetingNumber` and `webLink` are mutually exclusive.")
		cmd.Flags().StringVar(&current, "current", "", "Flag identifying to retrieve the current scheduled meeting of the meeting series or the entire meeting series. This parameter only applies to scenarios where `meetingNumber` is specified and the meeting is not an exceptional meeting from a meeting series. If it's `true`, return the scheduled meeting of the meeting series which is ready to join or start or the upcoming scheduled meeting of the meeting series; if it's `false`, return the entire meeting series. The default value is `true`.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		meetingsCmd.AddCommand(cmd)
	}

	{ // create
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Meeting",
			Long:  "Creates a new meeting. Regular users can schedule up to 100 meetings in 24 hours and admin users up to 3000 overall or 800 for a single user. Please note that the failed requests are also counted toward the limits.\n\n* The `spark:all` scope is required when `roomId` is specified.\n\n* If the parameter `adhoc` is `true` and `roomId` is specified, an ad-hoc meeting is created for the target room. An ad-hoc meeting is a non-recurring instant meeting for the target room which is supposed to be started immediately after being created for a quick collaboration. There's only one ad-hoc meeting for a room at the same time. So, if there's already an ongoing ad-hoc meeting for the room, the API returns this ongoing meeting instead of creating a new one. If it's a [direct](/docs/api/v1/rooms/get-room-details) room, both members of the room can create an ad-hoc meeting for the room. If it's a [group](/docs/api/v1/rooms/get-room-details) room, only room members that are in the same [organization](/docs/api/v1/organizations/get-organization-details) as the room can create an ad-hoc meeting for the room. Please note that an ad-hoc meeting is for the purpose of an instant collaboration with people in a room, user should not persist the `id` and `meetingNumber` of the ad-hoc meeting when it's been created since this meeting may become an inactive ad-hoc meeting for the room if it's not been started after being created for a while or it has been started and ended. Each time a user needs an ad-hoc meeting for a room, they should create one instead of reusing the previous persisted one. Moreover, for the same reason, no email will be sent when an ad-hoc meeting is created. Ad-hoc meetings cannot be updated by [Update a Meeting](/docs/api/v1/meetings/update-a-meeting) or deleted by [Delete a Meeting](/docs/api/v1/meetings/delete-a-meeting). Ad-hoc meetings cannot be listed by [List Meetings](/docs/api/v1/meetings/list-meetings) and the scheduled meetings of an ad-hoc meeting cannot be listed by [List Meetings of a Meeting Series](/docs/api/v1/meetings/list-meetings-of-a-meeting-series), but the ended and ongoing instances of ad-hoc meetings can be listed by [List Meetings](/docs/api/v1/meetings/list-meetings) and [List Meetings of a Meeting Series](/docs/api/v1/meetings/list-meetings-of-a-meeting-series).\n\n* If the parameter `adhoc` is `true`, `roomId` is required and the others are optional or ignored.\n\n* The default value of `title` for an ad-hoc meeting is the user's name if not specified. The following parameters for an ad-hoc meeting have default values and the user's input values will be ignored: `scheduledType` is always `meeting`; `start` and `end` are 5 minutes after the current time and 20 minutes after the current time respectively; `timezone` is `UTC`; `allowAnyUserToBeCoHost`, `allowAuthenticatedDevices`, `enabledJoinBeforeHost`, `enableConnectAudioBeforeHost` are always `true`; `allowFirstUserToBeCoHost`, `enableAutomaticLock`, `publicMeeting`, `sendEmail` are always `false`; `invitees` is the room members except \"me\"; `joinBeforeHostMinutes` is 5; `automaticLockMinutes` is null; `unlockedMeetingJoinSecurity` is `allowJoinWithLobby`. An ad-hoc meeting can be started immediately even if the `start` is 5 minutes after the current time.\n\n* The following parameters are not supported and will be ignored for an ad-hoc meeting: `templateId`, `recurrence`, `excludePassword`, `reminderTime`, `registration`, `integrationTags`, `enabledWebcastView`, and `panelistPassword`.\n\n* If the value of the parameter `recurrence` is null, a non-recurring meeting is created.\n\n* If the parameter `recurrence` has a value, a recurring meeting is created based on the rule defined by the value of `recurrence`. For a non-recurring meeting which has no `recurrence` value set, its `meetingType` is also `meetingSeries` which is a meeting series with only one occurrence in Webex meeting modeling. If you specify a `recurrence` like `FREQ=DAILY;INTERVAL=1` which never ends, the furthest date of the series is unlimited. You can also specify a `recurrence` with a very distant ending date in the future, e.g. `FREQ=DAILY;INTERVAL=1;UNTIL=21241001T000000Z`, but the actual furthest date accepted for the recurring meeting is five years from now. Specifically, if it has an ending date, there can be up to 5 occurrences for a yearly meeting, 60 occurrences for a monthly meeting, 261 occurrences for a weekly meeting, or 1826 occurrences for a daily meeting.\n\n* If the parameter `templateId` has a value, the meeting is created based on the meeting template specified by `templateId`. The list of meeting templates that is available for the authenticated user can be retrieved from [List Meeting Templates](/docs/api/v1/meetings/list-meeting-templates).\n\n* If the parameter `siteUrl` has a value, the meeting is created on the specified site. Otherwise, the meeting is created on the user's preferred site. All available Webex sites and preferred site of the user can be retrieved by `Get Site List` API.\n\n* If the parameter `scheduledType` equals \"personalRoomMeeting\", the meeting is created in the user's [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings).\n\n* If the parameter `roomId` has a value, the meeting is created for the Webex space specified by `roomId`. If `roomId` is specified but the user calling the API is not a member of the Webex space specified by `roomId`, the API will fail even if the user has the admin-level scopes or he is calling the API on behalf of another user which is specified by `hostEmail` and is a member of the Webex space.\n\n* If the parameter `enabledAudioWatermark` is `true`, `scheduledType` equals or defaults to `meeting`, and `audioConnectionOptions.audioConnectionType` equals `VoIP`, the audio for this meeting will have a watermark. In this case, a unique identifier is embedded into the audio that plays out of each Webex app and device. An administrator can use this watermark when analyzing an unauthorized recording to identify which Webex app or device was the source of the recording.\n\n* If the parameter `enabledVisualWatermark` is `true`, the video for this meeting will have a watermark. In this case, Webex superimposes a watermark image pattern on top of the meeting video and shared content to deter participants from leaking meeting information. Each participant viewing the meeting sees a watermark image pattern with their email address. If the participant is not signed in, the watermark image pattern includes their display name and email address.\n\n* The default value of `visualWatermarkOpacity` is 10 if not specified. The value must be between 5 and 80, inclusive. A smaller value means less distraction for meeting participants, while a larger value shows a clearer watermark. It's supported when `enabledVisualWatermark` is `true`.\n\n* When `enabledLiveStream` is `true`, `liveStream` must be specified. With these setting, the RTMP streaming specified by `liveStream.rtmpUrl` can be started and viewed during the meeting without any ad-hoc settings.\n\n* The `registration` can be specified when creating a meeting, but it can't be updated by [Update a Meeting](/docs/api/v1/meetings/update-a-meeting) or [Patch a Meeting](/docs/api/v1/meetings/patch-a-meeting). Create a registration form for a meeting that doesn't have one or update the registration form for a meeting that already has one by using [Update Meeting Registration Form](/docs/api/v1/meetings/update-meeting-registration-form). Delete the registration form for a meeting by using [Delete Meeting Registration Form](/docs/api/v1/meetings/delete-meeting-registration-form).\n\n* You can't create a meeting that starts 10 years or more in the future.\n\n* If all meeting invitees of a meeting should not receive emails, the host can create a meeting with invitees, and the parameter `sendEmail` is set to `false`. If only some meeting invitees should not receive emails and others can, the host should not invite these invitees along with creating a meeting request. Instead, the host should add the invitees by [Create a Meeting Invitee](/docs/api/v1/meeting-invitees/create-a-meeting-invitee) or [Create Meeting Invitees](/docs/api/v1/meeting-invitees/create-meeting-invitees) with the parameter `sendEmail` is set to `false` after the meeting has been created.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings")
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
		meetingsCmd.AddCommand(cmd)
	}

	{ // list
		var meetingNumber string
		var webLink string
		var roomId string
		var meetingSeriesId string
		var max string
		var from string
		var to string
		var meetingType string
		var state string
		var scheduledType string
		var isModified string
		var hasChat string
		var hasRecording string
		var hasTranscription string
		var hasSummary string
		var hasClosedCaption string
		var hasPolls string
		var hasQa string
		var hasSlido string
		var current string
		var hostEmail string
		var siteUrl string
		var integrationTag string
		var password string
		var timezone string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Meetings",
			Long:  "<div><Callout type=\"info\">The previous `List Meetings of a Meeting Series` API is merged into the [List Meetings](/docs/api/v1/meetings/list-meetings) API.</Callout></div>\n\nRetrieves details for meetings with a specified meeting number, web link, meeting type, etc. Please note that there are various products in the [Webex Suite](https://www.webex.com/collaboration-suite.html) such as `Meetings` and `Events`. Currently, only meetings of the `Meetings` product are supported by this API, meetings of others in the suite are not supported. Ad-hoc meetings created by [Create a Meeting](/docs/api/v1/meetings/create-a-meeting) with `adhoc` of `true` and a `roomId` will not be listed, but the ended and ongoing ad-hoc meeting instances will be listed.\n\nLists scheduled meeting and meeting instances of a meeting series identified by `meetingSeriesId`. Scheduled meetings of an ad-hoc meeting created by [Create a Meeting](/docs/api/v1/meetings/create-a-meeting) with `adhoc` of `true` and a `roomId` are not listed, but the ended and ongoing meeting instances of it are. Each _scheduled meeting_ or _meeting_ instance of a _meeting series_ has its own `start`, `end`, etc. Thus, for example, when a daily meeting has been scheduled from `2019-04-01` to `2019-04-10`, there are 10 scheduled meeting instances in this series, one instance for each day, and each one has its own attributes. When a scheduled meeting has been started and ended or is happening, there are even more ended or in-progress meeting instances.\n\nLong result sets are split into [pages](/docs/basics#pagination).\n\n* The default value of `meetingSeries` will be used if `meetingType` is not specified. When listing meetings with `meetingType=meetingSeries` implicitly or explicitly, the API returns all the recurring meeting series where the user is the meeting host or an invitee to the meeting. Please note that a meeting with no `recurrence` attribute is considered a meeting series with only one occurrence and it can also be listed with `meetingType=meetingSeries`. A recurring meeting series may have multiple occurrences which are scattered over weeks, months, or years. So, any meeting series that overlaps with the time range specified by `from` and `to` will be listed. For example, a monthly meeting series with `start=2024-01-01T10:00:00Z`, `end=2024-01-01T11:00:00Z` and `recurrence=FREQ=MONTHLY;INTERVAL=1;BYMONTHDAY=1;UNTIL=20250210T000000Z` can be listed with `meetingType=meetingSeries&from=2024-10-01&to=2024-11-01` even if the `start` and `end` are not in the specified time range.\n\n* If `meetingType` is specified and equals `meeting`, it lists ongoing and ended instances of the meeting series. The default value of `from` and `to` is: `from` equals the current date and time minus 7 days, and `to` equals the current date and time.\n\n* If `meetingSeriesId` is specified, `meetingNumber`, `webLink`, `roomId`, `current`, `integrationTag`, `scheduledType`, `siteUrl` will be ignored.\n\n* If `meetingNumber` is specified, the operation returns an array of meeting objects specified by the `meetingNumber`. Each object in the array can be a scheduled meeting or a meeting series depending on whether the `current` parameter is `true` or `false`, and each object contains the simultaneous interpretation object. When `meetingNumber` is specified, parameters of `from`, `to`, `meetingType`, `state`, `isModified` and `siteUrl` will be ignored. Please note that `meetingNumber`, `webLink` and `roomId` are mutually exclusive and they cannot be specified simultaneously.\n\n* If `webLink` is specified, the operation returns an array of meeting objects specified by the `webLink`. Each object in the array is a scheduled meeting, and each object contains the simultaneous interpretation object. When `webLink` is specified, parameters of `current`, `from`, `to`, `meetingType`, `state`, `isModified` and `siteUrl` will be ignored. Please note that `meetingNumber`, `webLink` and `roomId` are mutually exclusive and they cannot be specified simultaneously.\n\n* If `roomId` is specified, the operation returns an array of meeting objects of the Webex space specified by the `roomId`. When `roomId` is specified, parameters of `current`, `meetingType`, `state` and `isModified` will be ignored. The meeting objects are queried on the user's preferred site if no `siteUrl` is specified; otherwise, queried on the specified site. `meetingNumber`, `webLink` and `roomId` are mutually exclusive and they cannot be specified simultaneously.\n\n* If `state` parameter is specified, the returned array only has items in the specified state. If `state` is not specified, return items of all states.\n\n* If `meetingType` equals \"meetingSeries\", the `scheduledType` parameter can be \"meeting\", \"webinar\" or null. If `scheduledType` is specified, the returned array only has items of the specified scheduled type; otherwise, it has items of \"meeting\" and \"webinar\".\n\n* If `meetingType` equals \"scheduledMeeting\", the `scheduledType` parameter can be \"meeting\", \"webinar\", \"personalRoomMeeting\" or null. If `scheduledType` is specified, the returned array only has items of the specified scheduled type; otherwise, it has items of all scheduled types.\n\n* If `meetingType` equals \"meeting\", the `scheduledType` parameter can be \"meeting\", \"webinar\" or null. If `scheduledType` is specified, the returned array only has items of the specified scheduled type; otherwise, it has items of \"meeting\" and \"webinar\". Please note that ended or in-progress meeting instances of [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) also fall into the category of \"meeting\" `scheduledType`.\n\n* If `meetingType` equals \"meeting\", a maximum of 10000 meeting instances can be listed even if pagination is enabled.\n\n* If `isModified` parameter is specified, the returned array only has items which have been modified to exceptional meetings. This parameter only applies to scheduled meeting.\n\n* If any of the `hasChat`, `hasRecording`, `hasTranscription`, `hasSummary`, `hasClosedCaption`, `hasPolls `, `hasQA` and `hasSlido` parameters is specified, the `meetingType` must be \"meeting\" and `state` must be \"ended\". These parameters are null by default.\n\n* The `current` parameter only applies to meeting series. If it's `false`, the `start` and `end` attributes of each returned meeting series object are for the first scheduled meeting of that series. If it's `true` or not specified, the `start` and `end` attributes are for the scheduled meeting which is ready to start or join or the upcoming scheduled meeting of that series.\n\n* If `from` and `to` are specified, the operation returns an array of meeting objects in that specified time range.\n\n* If the parameter `siteUrl` has a value, the operation lists meetings on the specified site; otherwise, lists meetings on the user's all sites. All available Webex sites of the user can be retrieved by `Get Site List` API.\n\n* `trackingCodes` is not supported for ended meeting instances.\n\n* A full admin or a content admin can list all the ended and ongoing meeting instances of the organization he manages with the `meeting:admin_schedule_read` scope and `meetingType=meeting` parameter.\n\n* To learn more about which attributes are available for different meeting states, please refer to [Available Meeting Attributes for Different Meeting States](/docs/meetings#available-meeting-attributes-for-different-meeting-states).\n\n#### Request Header\n\n* `password`: Meeting password. Required when the meeting is protected by a password and the current user is not privileged to view it if they are not a host, cohost or invitee of the meeting.\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) for time stamps in response body, defined in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default value is `UTC` if not specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings")
				req.QueryParam("meetingNumber", meetingNumber)
				req.QueryParam("webLink", webLink)
				req.QueryParam("roomId", roomId)
				req.QueryParam("meetingSeriesId", meetingSeriesId)
				req.QueryParam("max", max)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("meetingType", meetingType)
				req.QueryParam("state", state)
				req.QueryParam("scheduledType", scheduledType)
				req.QueryParam("isModified", isModified)
				req.QueryParam("hasChat", hasChat)
				req.QueryParam("hasRecording", hasRecording)
				req.QueryParam("hasTranscription", hasTranscription)
				req.QueryParam("hasSummary", hasSummary)
				req.QueryParam("hasClosedCaption", hasClosedCaption)
				req.QueryParam("hasPolls", hasPolls)
				req.QueryParam("hasQA", hasQa)
				req.QueryParam("hasSlido", hasSlido)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("siteUrl", siteUrl)
				req.QueryParam("integrationTag", integrationTag)
				req.Header("password", password)
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
		cmd.Flags().StringVar(&meetingNumber, "meeting-number", "", "Meeting number for the meeting objects being requested. `meetingNumber`, `webLink` and `roomId` are mutually exclusive. If it's an exceptional meeting from a meeting series, the exceptional meeting instead of the primary meeting series is returned.")
		cmd.Flags().StringVar(&webLink, "web-link", "", "URL encoded link to information page for the meeting objects being requested. `meetingNumber`, `webLink` and `roomId` are mutually exclusive.")
		cmd.Flags().StringVar(&roomId, "room-id", "", "Associated Webex space ID for the meeting objects being requested. `meetingNumber`, `webLink` and `roomId` are mutually exclusive.")
		cmd.Flags().StringVar(&meetingSeriesId, "meeting-series-id", "", "Unique identifier for the meeting series. The meeting ID of a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting is not supported. If `meetingSeriesId` is specified, it lists all occurrences and instances of the meeting series by default; with `meetingType` of `scheduledMeeting`, it lists occurrences; with `meetingType` of `meeting`, it lists ongoing and ended instances.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of meetings in the response, up to 100. This parameter is ignored if `meetingNumber`, `webLink` or `roomId` is specified. The default value is 10.")
		cmd.Flags().StringVar(&from, "from", "", "Start date and time (inclusive) for the range for which meetings are to be returned in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `from` cannot be after `to`. This parameter will be ignored if `meetingNumber`, `webLink` or `roomId` is specified.  When `meetingType` is `meetingSeries`(either explicitly set or by default), if `to` is specified, the default value for `from` is `to` minus 7 days. If `to` is also not specified, the default value for `from` is the current date and time. When `meetingType` is `scheduledMeeting`, `from` is the same as above. When `meetingType` is `meeting`, if `to` is specified, the default value for `from` is `to` minus 7 days. If `to` is also not specified, the default value for `from` is 7 days before the current date and time.")
		cmd.Flags().StringVar(&to, "to", "", "End date and time (exclusive) for the range for which meetings are to be returned in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `to` cannot be before `from`. This parameter will be ignored if `meetingNumber`, `webLink` or `roomId` is specified.  When `meetingType` is `meetingSeries`(either explicitly set or by default), if `from` is specified, the default value for `to` is `from` plus 7 days. If `from` is also not specified, the default value for `to` is 7 days after the current date and time. When `meetingType` is `scheduledMeeting`, `to` is the same as above. When `meetingType` is `meeting`, if `from` is specified, the default value for `to` is `from` plus 7 days. If `from` is also not specified, the default value for `to` is the current date and time.")
		cmd.Flags().StringVar(&meetingType, "meeting-type", "", "Meeting type for the meeting objects being requested. This parameter will be ignored if `meetingNumber`, `webLink` or `roomId` is specified.")
		cmd.Flags().StringVar(&state, "state", "", "Meeting state for the meeting objects being requested. If not specified, return meetings of all states. This parameter will be ignored if `meetingNumber`, `webLink` or `roomId` is specified. Details of an `ended` meeting will only be available 15 minutes after the meeting has ended. `inProgress` meetings are not fully supported. The API will try to return details of an `inProgress` meeting 15 minutes after the meeting starts. However, it may take longer depending on the traffic. See the [Webex Meetings](/docs/meetings#meeting-states) guide for more information about the states of meetings.")
		cmd.Flags().StringVar(&scheduledType, "scheduled-type", "", "Scheduled type for the meeting objects being requested.")
		cmd.Flags().StringVar(&isModified, "is-modified", "", "Flag identifying whether a meeting has been modified. Only applies to scheduled meetings. If `true`, only return modified scheduled meetings; if `false`, only return unmodified scheduled meetings; if not specified, all scheduled meetings will be returned.")
		cmd.Flags().StringVar(&hasChat, "has-chat", "", "Flag identifying whether a meeting has a chat log. Only applies to ended meeting instances. If `true`, only return meeting instances which have chats; if `false`, only return meeting instances which have no chats; if not specified, all meeting instances will be returned.")
		cmd.Flags().StringVar(&hasRecording, "has-recording", "", "Flag identifying meetings which have been recorded. Only applies to ended meeting instances. If true, only return meeting instances which have been recorded; if false, only return meeting instances which have not been recorded; if not specified, all meeting instances will be returned.")
		cmd.Flags().StringVar(&hasTranscription, "has-transcription", "", "Flag identifying meetings with transcripts. Only applies to ended meeting instances. If `true`, only return meeting instances which have transcripts; if `false`, only return meeting instances which have no transcripts; if not specified, all meeting instances will be returned.")
		cmd.Flags().StringVar(&hasSummary, "has-summary", "", "Flag identifying meetings with summaries. Only applies to ended meeting instances. If `true`, only return meeting instances which have summaries; if `false`, only return meeting instances which have no summaries; if not specified, all meeting instances will be returned.")
		cmd.Flags().StringVar(&hasClosedCaption, "has-closed-caption", "", "Flag identifying meetings with closed captions. Only applies to ended meeting instances. If `true`, only return meeting instances which have closed captions; if `false`, only return meeting instances which have no closed captions; if not specified, all meeting instances will be returned.")
		cmd.Flags().StringVar(&hasPolls, "has-polls", "", "Flag identifying meetings with polls. Only applies to ended meeting instances. If `true`, only return meeting instances which have polls; if `false`, only return meeting instances which have no polls; if not specified, all meeting instances will be returned.")
		cmd.Flags().StringVar(&hasQa, "has-qa", "", "Flag identifying meetings with Q&A. Only applies to ended meeting instances. If `true`, only return meeting instances which have Q&A; if `false`, only return meeting instances which have no Q&A; if not specified, all meeting instances will be returned.")
		cmd.Flags().StringVar(&hasSlido, "has-slido", "", "Flag identifying meetings with Slido interactions. Only applies to ended meeting instances. If `true`, only return meeting instances which have Slido interactions like Q&A or polling; if `false`, only return meeting instances which have no Slido interactions; if not specified, all meeting instances will be returned.")
		cmd.Flags().StringVar(&current, "current", "", "Flag identifying to retrieve the current scheduled meeting of the meeting series or the entire meeting series. This parameter only applies to scenarios where `meetingNumber` is specified and the meeting is not an exceptional meeting from a meeting series. If it's `true`, return the scheduled meeting of the meeting series which is ready to join or start or the upcoming scheduled meeting of the meeting series; if it's `false`, return the entire meeting series.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API returns meetings as if the user calling the API were the user of `hostEmail` themself, and the meetings returned by the API include the meetings where the user of `hostEmail` is the meeting host and those where they are an invitee.")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API lists meetings from. If not specified, the API lists meetings from user's all sites. All available Webex sites of the user can be retrieved by `Get Site List` API.")
		cmd.Flags().StringVar(&integrationTag, "integration-tag", "", "External key created by an integration application. This parameter is used by the integration application to query meetings by a key in its own domain such as a Zendesk ticket ID, a Jira ID, a Salesforce Opportunity ID, etc. An integrationTag created by one client cannot be accessed or used as a filtering parameter by another client. For example, if a meeting has an `integrationTag` of \"Sales\" which is created by the client behind the developer portal, then this integrationTag can't be accessed on the meeting or its recordings by another client. Neither can it be used to filter meetings or recordings by a client other than the one that created the integrationTag of \"Sales\".")
		cmd.Flags().StringVar(&password, "password", "", "e.g. BgJep@4323")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		meetingsCmd.AddCommand(cmd)
	}

	{ // get
		var meetingId string
		var current string
		var hostEmail string
		var password string
		var timezone string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get a Meeting",
			Long:  "Retrieves details for a meeting with a specified meeting ID.\n\n* If the `meetingId` value specified is for a meeting series and `current` is `true`, the operation returns details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series.\n\n* If the `meetingId` value specified is for a meeting series and `current` is `false` or `current` is not specified, the operation returns details for the entire meeting series.\n\n* If the `meetingId` value specified is for a scheduled meeting from a meeting series, the operation returns details for that scheduled meeting.\n\n* If the `meetingId` value specified is for a meeting instance which is happening or has happened, the operation returns details for that meeting instance.\n\n* `trackingCodes` is not supported for ended meeting instances.\n\n* To learn more about which attributes are available for different meeting states, please refer to [Available Meeting Attributes for Different Meeting States](/docs/meetings#available-meeting-attributes-for-different-meeting-states).\n\n#### Request Header\n\n* `password`: Meeting password. Required when the meeting is protected by a password and the current user is not privileged to view it if they are not a host, cohost or invitee of the meeting.\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) for time stamps in response body, defined in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default value is `UTC` if not specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/{meetingId}")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
				req.Header("password", password)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting being requested.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		cmd.Flags().StringVar(&password, "password", "", "e.g. BgJep@4323")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		meetingsCmd.AddCommand(cmd)
	}

	{ // patch
		var meetingId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "patch",
			Short: "Patch a Meeting",
			Long:  "Updates details for a meeting with a specified meeting ID. This operation applies to meeting series and scheduled meetings. It doesn't apply to ended or in-progress meeting instances. Ad-hoc meetings created by [Create a Meeting](/docs/api/v1/meetings/create-a-meeting) with `adhoc` of `true` and a `roomId` cannot be updated.\n\n* If the `meetingId` value specified is for a scheduled meeting, the operation updates that scheduled meeting without impact on other scheduled meeting of the parent meeting series.\n\n* If the `meetingId` value specified is for a meeting series, the operation updates the entire meeting series. **Note**: If the value of `start`, `end`, or `recurrence` for the meeting series is changed, any exceptional scheduled meeting in this series is cancelled when the meeting series is updated.\n\n* The `agenda`, `recurrence`, and `trackingCodes` attributes can be specified as `null` so that these attributes become null and hidden from the response after the patch. Note that it's the keyword `null` not the string \"null\".\n\n* If the parameter `recurrence` has a value, a recurring meeting is created based on the rule defined by the value of `recurrence`. For a non-recurring meeting which has no `recurrence` value set, its `meetingType` is also `meetingSeries` which is a meeting series with only one occurrence in Webex meeting modeling. If you specify a `recurrence` like `FREQ=DAILY;INTERVAL=1` which never ends, the furthest date of the series is unlimited. You can also specify a `recurrence` with a very distant ending date in the future, e.g. `FREQ=DAILY;INTERVAL=1;UNTIL=21241001T000000Z`, but the actual furthest date accepted for the recurring meeting is five years from now. Specifically, if it has an ending date, there can be up to 5 occurrences for a yearly meeting, 60 occurrences for a monthly meeting, 261 occurrences for a weekly meeting, or 1826 occurrences for a daily meeting.\n\n* You can't update a meeting that starts 10 years or more in the future.\n\n* Updating a meeting in the API that was created via a calendar connector is not allowed. The meeting may be updated in Webex, but the calendar event may not be updated, resulting in duplicate entries. This action is therefore blocked. In case you must overwrite this behavior, please contact devsupport@webex.com.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PATCH", "/meetings/{meetingId}")
				req.PathParam("meetingId", meetingId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting to be updated. This parameter applies to meeting series and scheduled meetings. It doesn't apply to ended or in-progress meeting instances. Please note that currently meeting ID of a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting is not supported for this API.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // update
		var meetingId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Meeting",
			Long:  "<div>\n<Callout type=\"warning\">The PUT method is still supported and behaves the same as before, will be deprecated in the future. Use the PATCH method instead.</Callout>\n</div>\n\nUpdates details for a meeting with a specified meeting ID. This operation applies to meeting series and scheduled meetings. It doesn't apply to ended or in-progress meeting instances. Ad-hoc meetings created by [Create a Meeting](/docs/api/v1/meetings/create-a-meeting) with `adhoc` of `true` and a `roomId` cannot be updated.\n\n* If the `meetingId` value specified is for a scheduled meeting, the operation updates that scheduled meeting without impact on other scheduled meeting of the parent meeting series.\n\n* If the `meetingId` value specified is for a meeting series, the operation updates the entire meeting series. **Note**: If the value of `start`, `end`, or `recurrence` for the meeting series is changed, any exceptional scheduled meeting in this series is cancelled when the meeting series is updated.\n\n* If the parameter `recurrence` has a value, a recurring meeting is created based on the rule defined by the value of `recurrence`. For a non-recurring meeting which has no `recurrence` value set, its `meetingType` is also `meetingSeries` which is a meeting series with only one occurrence in Webex meeting modeling. If you specify a `recurrence` like `FREQ=DAILY;INTERVAL=1` which never ends, the furthest date of the series is unlimited. You can also specify a `recurrence` with a very distant ending date in the future, e.g. `FREQ=DAILY;INTERVAL=1;UNTIL=21241001T000000Z`, but the actual furthest date accepted for the recurring meeting is five years from now. Specifically, if it has an ending date, there can be up to 5 occurrences for a yearly meeting, 60 occurrences for a monthly meeting, 261 occurrences for a weekly meeting, or 1826 occurrences for a daily meeting.\n\n* You can't update a meeting that starts 10 years or more in the future.\n\n* Updating a meeting created using the API via a calendar connector is not supported. While the meeting may be updated in Webex, the calendar event may not be updated resulting in duplicate entries. If you need to override this behavior, contact devsupport@webex.com.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/meetings/{meetingId}")
				req.PathParam("meetingId", meetingId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting to be updated. This parameter applies to meeting series and scheduled meetings. It doesn't apply to ended or in-progress meeting instances. Please note that currently meeting ID of a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting is not supported for this API.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // delete
		var meetingId string
		var hostEmail string
		var sendEmail string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Meeting",
			Long:  "Deletes a meeting with a specified meeting ID. The deleted meeting cannot be recovered. This operation applies to meeting series and scheduled meetings. It doesn't apply to ended or in-progress meeting instances. Ad-hoc meetings created by [Create a Meeting](/docs/api/v1/meetings/create-a-meeting) with `adhoc` of `true` and a `roomId` cannot be deleted.\n\n* If the `meetingId` value specified is for a scheduled meeting, the operation deletes that scheduled meeting without impact on other scheduled meeting of the parent meeting series.\n\n* If the `meetingId` value specified is for a meeting series, the operation deletes the entire meeting series.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/meetings/{meetingId}")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("sendEmail", sendEmail)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting to be deleted. This parameter applies to meeting series and scheduled meetings. It doesn't apply to ended or in-progress meeting instances.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will delete a meeting that is hosted by that user.")
		cmd.Flags().StringVar(&sendEmail, "send-email", "", "Whether or not to send emails to host and invitees. It is an optional field and default value is true.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // join
		var meetingId string
		var meetingNumber string
		var webLink string
		var joinDirectly bool
		var email string
		var displayName string
		var password string
		var expirationMinutes int64
		var registrationId string
		var hostEmail string
		var createJoinLinkAsWebLink bool
		var createStartLinkAsWebLink bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "join",
			Short: "Join a Meeting",
			Long:  "Retrieves links for a meeting with a specified `meetingId`, `meetingNumber`, or `webLink` that allow users to start or join the meeting directly without logging in and entering a password.\n\n* Please note that `meetingId`, `meetingNumber` and `webLink` are mutually exclusive and they cannot be specified simultaneously.\n\n* If `joinDirectly` is true or not specified, the response will have HTTP response code 302 and the request will be redirected to `joinLink`; otherwise, the response will have HTTP response code 200 and `joinLink` will be returned in response body.\n\n* Only the meeting host or cohost can generate the `startLink`.\n\n* An admin user or a [Service App](/docs/service-apps) can generate the `startLink` and `joinLink` on behalf of another meeting host using the `hostEmail` parameter. When a [Service App](/docs/service-apps) generates the `startLink` and `joinLink`, the `hostEmail` parameter is required. The `hostEmail` parameter only applies to meetings, not webinars.\n\n* For Service Apps, `hostEmail` must be provided in the request.\n\n* Generating a join link or a start link before the time specified by `joinBeforeHostMinutes` for a webinar is not supported.\n\n* The `joinLink` and `startLink` generated by the API only work in a web browser, not in mobile apps. When the `joinLink` or `startLink` is used in a mobile app, the user can't join or start the meeting directly by the dial-in or dial-out phone numbers displayed in the app.\n\n* When `createJoinLinkAsWebLink` or `createStartLinkAsWebLink` is set to true, a user cannot join or start the meeting using the `joinLink` or `startLink` returned in the response, and must complete the login flow. Those options are typically useful if mandatory user login is configured in Control Hub.\n\n* Appending the parameter `&launchApp=true` to the weblink will launch directly into the browser experience without showing the user the download selection for the native client.\n\n* [Embedded Apps](/docs/embedded-apps) are not supported in a meeting which is started from the `startLink` generated by this API.\n\n<div><Callout type=\"warning\">When the `email` or `displayName` is omitted from the API request, the backend inserts the data from the user making the request. That may lead to situations where an attendee is shown in the roster with the hostname which should be avoided. In the future, we will make those two fields required.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/join")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("meetingId", meetingId)
					req.BodyString("meetingNumber", meetingNumber)
					req.BodyString("webLink", webLink)
					req.BodyBool("joinDirectly", joinDirectly, cmd.Flags().Changed("join-directly"))
					req.BodyString("email", email)
					req.BodyString("displayName", displayName)
					req.BodyString("password", password)
					req.BodyInt("expirationMinutes", expirationMinutes, cmd.Flags().Changed("expiration-minutes"))
					req.BodyString("registrationId", registrationId)
					req.BodyString("hostEmail", hostEmail)
					req.BodyBool("createJoinLinkAsWebLink", createJoinLinkAsWebLink, cmd.Flags().Changed("create-join-link-as-web-link"))
					req.BodyBool("createStartLinkAsWebLink", createStartLinkAsWebLink, cmd.Flags().Changed("create-start-link-as-web-link"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "")
		cmd.Flags().StringVar(&meetingNumber, "meeting-number", "", "")
		cmd.Flags().StringVar(&webLink, "web-link", "", "")
		cmd.Flags().BoolVar(&joinDirectly, "join-directly", false, "")
		cmd.Flags().StringVar(&email, "email", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&password, "password", "", "")
		cmd.Flags().Int64Var(&expirationMinutes, "expiration-minutes", 0, "")
		cmd.Flags().StringVar(&registrationId, "registration-id", "", "")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "")
		cmd.Flags().BoolVar(&createJoinLinkAsWebLink, "create-join-link-as-web-link", false, "")
		cmd.Flags().BoolVar(&createStartLinkAsWebLink, "create-start-link-as-web-link", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // list-templates
		var templateType string
		var locale string
		var isDefault string
		var isStandard string
		var hostEmail string
		var siteUrl string
		cmd := &cobra.Command{
			Use:   "list-templates",
			Short: "List Meeting Templates",
			Long:  "Retrieves the list of meeting templates that is available for the authenticated user.\n\nThere are separate lists of meeting templates for different `templateType`, `locale` and `siteUrl`.\n\n* If `templateType` is specified, the operation returns an array of meeting template objects specified by the `templateType`; otherwise, returns an array of meeting template objects of all template types.\n\n* If `locale` is specified, the operation returns an array of meeting template objects specified by the `locale`; otherwise, returns an array of meeting template objects of the default `en_US` locale. Refer to [Meeting Template Locales](/docs/meetings#meeting-template-locales) for all the locales supported by Webex.\n\n* If the parameter `siteUrl` has a value, the operation lists meeting templates on the specified site; otherwise, lists meeting templates on the user's preferred site. All available Webex sites and preferred site of the user can be retrieved by `Get Site List` API.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/templates")
				req.QueryParam("templateType", templateType)
				req.QueryParam("locale", locale)
				req.QueryParam("isDefault", isDefault)
				req.QueryParam("isStandard", isStandard)
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("siteUrl", siteUrl)
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
		cmd.Flags().StringVar(&templateType, "template-type", "", "Meeting template type for the meeting template objects being requested. If not specified, return meeting templates of all types.")
		cmd.Flags().StringVar(&locale, "locale", "", "Locale for the meeting template objects being requested. If not specified, return meeting templates of the default `en_US` locale. Refer to [Meeting Template Locales](/docs/meetings#meeting-template-locales) for all the locales supported by Webex.")
		cmd.Flags().StringVar(&isDefault, "is-default", "", "The value is `true` or `false`. If it's `true`, return the default meeting templates; if it's `false`, return the non-default meeting templates. If it's not specified, return both default and non-default meeting templates.")
		cmd.Flags().StringVar(&isStandard, "is-standard", "", "The value is `true` or `false`. If it's `true`, return the standard meeting templates; if it's `false`, return the non-standard meeting templates. If it's not specified, return both standard and non-standard meeting templates.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return meeting templates that are available for that user.")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API lists meeting templates from. If not specified, the API lists meeting templates from user's preferred site. All available Webex sites and preferred site of the user can be retrieved by `Get Site List` API.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // get-template
		var templateId string
		var hostEmail string
		var timezone string
		cmd := &cobra.Command{
			Use:   "get-template",
			Short: "Get a Meeting Template",
			Long:  "Retrieves details for a meeting template with a specified meeting template ID.\n\n#### Request Header\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) for time stamps in response body, defined in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default value is `UTC` if not specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/templates/{templateId}")
				req.PathParam("templateId", templateId)
				req.QueryParam("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&templateId, "template-id", "", "Unique identifier for the meeting template being requested.")
		cmd.MarkFlagRequired("template-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return the meeting template that is available for that user.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		meetingsCmd.AddCommand(cmd)
	}

	{ // get-control-status
		var meetingId string
		cmd := &cobra.Command{
			Use:   "get-control-status",
			Short: "Get Meeting Control Status",
			Long:  `Get the meeting control of a live meeting, which is consisted of meeting control status on "locked" and "recording" to reflect whether the meeting is currently locked and there is recording in progress.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/controls")
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Does not support meeting IDs for a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // update-control-status
		var meetingId string
		var recordingStarted bool
		var recordingPaused bool
		var locked bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-control-status",
			Short: "Update Meeting Control Status",
			Long:  `To start, pause, resume, or stop a meeting recording; To lock or unlock an on-going meeting.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/meetings/controls")
				req.QueryParam("meetingId", meetingId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("recordingStarted", recordingStarted, cmd.Flags().Changed("recording-started"))
					req.BodyBool("recordingPaused", recordingPaused, cmd.Flags().Changed("recording-paused"))
					req.BodyBool("locked", locked, cmd.Flags().Changed("locked"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Does not support meeting IDs for a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting.")
		cmd.Flags().BoolVar(&recordingStarted, "recording-started", false, "")
		cmd.Flags().BoolVar(&recordingPaused, "recording-paused", false, "")
		cmd.Flags().BoolVar(&locked, "locked", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // list-session-types
		var hostEmail string
		var siteUrl string
		cmd := &cobra.Command{
			Use:   "list-session-types",
			Short: "List Meeting Session Types",
			Long:  `List all the meeting session types enabled for a given user.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/sessionTypes")
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("siteUrl", siteUrl)
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
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the user. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will list all the meeting session types enabled for the user.")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "Webex site URL to query. If `siteUrl` is not specified, the users' preferred site will be used. If the authorization token has the admin-level scopes, the admin can set the Webex site URL on behalf of the user specified in the `hostEmail` parameter.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // get-session-type
		var sessionTypeId string
		var hostEmail string
		var siteUrl string
		cmd := &cobra.Command{
			Use:   "get-session-type",
			Short: "Get a Meeting Session Type",
			Long:  `Retrieves details for a meeting session type with a specified session type ID.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/sessionTypes/{sessionTypeId}")
				req.PathParam("sessionTypeId", sessionTypeId)
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("siteUrl", siteUrl)
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
		cmd.Flags().StringVar(&sessionTypeId, "session-type-id", "", "A unique identifier for the sessionType.")
		cmd.MarkFlagRequired("session-type-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the user. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will get a meeting session type with the specified session type ID enabled for the user.")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "Webex site URL to query. If `siteUrl` is not specified, the users' preferred site will be used. If the authorization token has the admin-level scopes, the admin can set the Webex site URL on behalf of the user specified in the `hostEmail` parameter.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // get-registration-form
		var meetingId string
		var current string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "get-registration-form",
			Short: "Get registration form for a meeting",
			Long:  `Get a meeting's registration form to understand which fields are required.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/{meetingId}/registration")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("current", current)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // update-registration-form
		var meetingId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-registration-form",
			Short: "Update Meeting Registration Form",
			Long:  `Enable or update a registration form for a meeting.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/meetings/{meetingId}/registration")
				req.PathParam("meetingId", meetingId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // delete-registration-form
		var meetingId string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "delete-registration-form",
			Short: "Delete Meeting Registration Form",
			Long:  `Disable the registration form for a meeting.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/meetings/{meetingId}/registration")
				req.PathParam("meetingId", meetingId)
				req.Header("hostEmail", hostEmail)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "e.g. `brenda.song@example.com` (string, optional) - Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for an interpreter of the meeting that is hosted by that user.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // register-registrant
		var meetingId string
		var current string
		var hostEmail string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "register-registrant",
			Short: "Register a Meeting Registrant",
			Long:  `Register a new registrant for a meeting. When a meeting or webinar is created, this API can only be used if Registration is checked on the page or the registration attribute is specified through the [Create a Meeting](/docs/api/v1/meetings/create-a-meeting) API.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/registrants")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // list-registrants
		var meetingId string
		var max string
		var hostEmail string
		var current string
		var email string
		var registrationTimeFrom string
		var registrationTimeTo string
		cmd := &cobra.Command{
			Use:   "list-registrants",
			Short: "List Meeting Registrants",
			Long:  `Meeting's host and cohost can retrieve the list of registrants for a meeting with a specified meeting Id.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/{meetingId}/registrants")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("max", max)
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("current", current)
				req.QueryParam("email", email)
				req.QueryParam("registrationTimeFrom", registrationTimeFrom)
				req.QueryParam("registrationTimeTo", registrationTimeTo)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of registrants in the response, up to 100. The default is 10.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.")
		cmd.Flags().StringVar(&email, "email", "", "Registrant's email to filter registrants.")
		cmd.Flags().StringVar(&registrationTimeFrom, "registration-time-from", "", "The time registrants register a meeting starts from the specified date and time (inclusive) in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. If `registrationTimeFrom` is not specified, it equals `registrationTimeTo` minus 7 days.")
		cmd.Flags().StringVar(&registrationTimeTo, "registration-time-to", "", "The time registrants register a meeting before the specified date and time (exclusive) in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. If `registrationTimeTo` is not specified, it equals `registrationTimeFrom` plus 7 days. The interval between `registrationTimeFrom` and `registrationTimeTo` must be within 90 days.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // batch-register-registrants
		var meetingId string
		var current string
		var hostEmail string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "batch-register-registrants",
			Short: "Batch register Meeting Registrants",
			Long:  `Bulk register new registrants for a meeting. When a meeting or webinar is created, this API can only be used if Registration is checked on the page or the registration attribute is specified through the [Create a Meeting](/docs/api/v1/meetings/create-a-meeting) API.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/registrants/bulkInsert")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // get-detailed-information-registrant
		var meetingId string
		var registrantId string
		var current string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "get-detailed-information-registrant",
			Short: "Get Detailed Information for a Meeting Registrant",
			Long:  `Retrieves details for a meeting registrant with a specified registrant Id.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/{meetingId}/registrants/{registrantId}")
				req.PathParam("meetingId", meetingId)
				req.PathParam("registrantId", registrantId)
				req.QueryParam("current", current)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&registrantId, "registrant-id", "", "Unique identifier for the registrant")
		cmd.MarkFlagRequired("registrant-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // delete-registrant
		var meetingId string
		var registrantId string
		var current string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "delete-registrant",
			Short: "Delete a Meeting Registrant",
			Long:  `Meeting's host or cohost can delete a registrant with a specified registrant ID.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/meetings/{meetingId}/registrants/{registrantId}")
				req.PathParam("meetingId", meetingId)
				req.PathParam("registrantId", registrantId)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&registrantId, "registrant-id", "", "Unique identifier for the registrant.")
		cmd.MarkFlagRequired("registrant-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // query-registrants
		var meetingId string
		var max string
		var current string
		var hostEmail string
		var emails []string
		var status string
		var orderType string
		var orderBy string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "query-registrants",
			Short: "Query Meeting Registrants",
			Long:  `Meeting's host and cohost can query the list of registrants for a meeting with a specified meeting ID and registrants email.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/registrants/query")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("max", max)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("emails", emails)
					req.BodyString("status", status)
					req.BodyString("orderType", orderType)
					req.BodyString("orderBy", orderBy)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of registrants in the response, up to 100. The default is 10.")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		cmd.Flags().StringSliceVar(&emails, "emails", nil, "")
		cmd.Flags().StringVar(&status, "status", "", "")
		cmd.Flags().StringVar(&orderType, "order-type", "", "")
		cmd.Flags().StringVar(&orderBy, "order-by", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // batch-update-registrants-status
		var meetingId string
		var statusOpType string
		var current string
		var hostEmail string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "batch-update-registrants-status",
			Short: "Batch Update Meeting Registrants status",
			Long:  "Meeting's host or cohost can update the set of registrants for a meeting. `cancel` means the registrant(s) will be moved back to the registration list. `bulkDelete` means the registrant(s) will be deleted.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/registrants/{statusOpType}")
				req.PathParam("meetingId", meetingId)
				req.PathParam("statusOpType", statusOpType)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&statusOpType, "status-op-type", "", "Update registrant's status.")
		cmd.MarkFlagRequired("status-op-type")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.     + Default: `false` ")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // update-simultaneous-interpretation
		var meetingId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-simultaneous-interpretation",
			Short: "Update Meeting Simultaneous interpretation",
			Long:  `Updates simultaneous interpretation options of a meeting with a specified meeting ID. This operation applies to meeting series and scheduled meetings. It doesn't apply to ended or in-progress meeting instances.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/meetings/{meetingId}/simultaneousInterpretation")
				req.PathParam("meetingId", meetingId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Does not support meeting IDs for a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // create-interpreter
		var meetingId string
		var languageCode1 string
		var languageCode2 string
		var email string
		var displayName string
		var hostEmail string
		var sendEmail bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-interpreter",
			Short: "Create a Meeting Interpreter",
			Long:  `Assign an interpreter to a bi-directional simultaneous interpretation language channel for a meeting.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/interpreters")
				req.PathParam("meetingId", meetingId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("languageCode1", languageCode1)
					req.BodyString("languageCode2", languageCode2)
					req.BodyString("email", email)
					req.BodyString("displayName", displayName)
					req.BodyString("hostEmail", hostEmail)
					req.BodyBool("sendEmail", sendEmail, cmd.Flags().Changed("send-email"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting to which the interpreter is to be assigned.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&languageCode1, "language-code1", "", "")
		cmd.Flags().StringVar(&languageCode2, "language-code2", "", "")
		cmd.Flags().StringVar(&email, "email", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "")
		cmd.Flags().BoolVar(&sendEmail, "send-email", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // list-interpreters
		var meetingId string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "list-interpreters",
			Short: "List Meeting Interpreters",
			Long:  "Lists meeting interpreters for a meeting with a specified `meetingId`.\n\nThis operation can be used for meeting series, scheduled meeting and ended or ongoing meeting instance objects. If the specified `meetingId` is for a meeting series, the interpreters for the series will be listed; if the `meetingId` is for a scheduled meeting, the interpreters for the particular scheduled meeting will be listed; if the `meetingId` is for an ended or ongoing meeting instance, the interpreters for the particular meeting instance will be listed. See the [Webex Meetings](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) guide for more information about the types of meetings.\n\nThe list returned is sorted in descending order by when interpreters were created.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/{meetingId}/interpreters")
				req.PathParam("meetingId", meetingId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting for which interpreters are being requested. The meeting can be meeting series, scheduled meeting or meeting instance which has ended or is ongoing. Please note that currently meeting ID of a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting is not supported for this API.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will return interpreters of the meeting that is hosted by that user.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // get-interpreter
		var meetingId string
		var interpreterId string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "get-interpreter",
			Short: "Get a Meeting Interpreter",
			Long:  "Retrieves details for a meeting interpreter identified by `meetingId` and `interpreterId` in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/{meetingId}/interpreters/{interpreterId}")
				req.PathParam("meetingId", meetingId)
				req.PathParam("interpreterId", interpreterId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting to which the interpreter has been assigned.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&interpreterId, "interpreter-id", "", "Unique identifier for the interpreter whose details are being requested.")
		cmd.MarkFlagRequired("interpreter-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for an interpreter of the meeting that is hosted by that user.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // update-interpreter
		var meetingId string
		var interpreterId string
		var languageCode1 string
		var languageCode2 string
		var email string
		var displayName string
		var hostEmail string
		var sendEmail bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-interpreter",
			Short: "Update a Meeting Interpreter",
			Long:  "Updates details for a meeting interpreter identified by `meetingId` and `interpreterId` in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/meetings/{meetingId}/interpreters/{interpreterId}")
				req.PathParam("meetingId", meetingId)
				req.PathParam("interpreterId", interpreterId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("languageCode1", languageCode1)
					req.BodyString("languageCode2", languageCode2)
					req.BodyString("email", email)
					req.BodyString("displayName", displayName)
					req.BodyString("hostEmail", hostEmail)
					req.BodyBool("sendEmail", sendEmail, cmd.Flags().Changed("send-email"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting whose interpreters were belong to.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&interpreterId, "interpreter-id", "", "Unique identifier for the interpreter whose details are being requested.")
		cmd.MarkFlagRequired("interpreter-id")
		cmd.Flags().StringVar(&languageCode1, "language-code1", "", "")
		cmd.Flags().StringVar(&languageCode2, "language-code2", "", "")
		cmd.Flags().StringVar(&email, "email", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "")
		cmd.Flags().BoolVar(&sendEmail, "send-email", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // delete-interpreter
		var meetingId string
		var interpreterId string
		var hostEmail string
		var sendEmail string
		cmd := &cobra.Command{
			Use:   "delete-interpreter",
			Short: "Delete a Meeting Interpreter",
			Long:  "Removes a meeting interpreter identified by `meetingId` and `interpreterId` in the URI. The deleted meeting interpreter cannot be recovered.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/meetings/{meetingId}/interpreters/{interpreterId}")
				req.PathParam("meetingId", meetingId)
				req.PathParam("interpreterId", interpreterId)
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("sendEmail", sendEmail)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting whose interpreters were belong to.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&interpreterId, "interpreter-id", "", "Unique identifier for the interpreter to be removed.")
		cmd.MarkFlagRequired("interpreter-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will delete an interpreter of the meeting that is hosted by that user.")
		cmd.Flags().StringVar(&sendEmail, "send-email", "", "If `true`, send email to the interpreter.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // list-breakout-sessions
		var meetingId string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "list-breakout-sessions",
			Short: "List Meeting Breakout Sessions",
			Long:  "Lists meeting breakout sessions for a meeting with a specified `meetingId`.\n\nThis operation can be used for meeting series, scheduled meeting and ended or ongoing meeting instance objects. See the [Webex Meetings](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) guide for more information about the types of meetings.\n\n* If the meeting of `meetingId` is in progress, the operation returns the breakout sessions which are currently held in the meeting.\n\n* If the meeting of `meetingId` is ended meeting instance, the operation returns the breakout sessions which happended during the meeting instance.\n\n* Otherwise, if it's a meeting series or scheduled meeting and it's not in progress, the operation returns the breakout sessions which are created when the meeting is scheduled.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/{meetingId}/breakoutSessions")
				req.PathParam("meetingId", meetingId)
				req.Header("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. This parameter applies to meeting series, scheduled meeting and ended or ongoing meeting instance objects. Please note that currently meeting ID of a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting is not supported for this API.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "e.g. `john.andersen@example.com` (string, optional) - Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // update-breakout-sessions
		var meetingId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-breakout-sessions",
			Short: "Update Meeting Breakout Sessions",
			Long:  `Updates breakout sessions of a meeting with a specified meeting ID in the pre-meeting state. This operation applies to meeting series and scheduled meetings.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/meetings/{meetingId}/breakoutSessions")
				req.PathParam("meetingId", meetingId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Does not support meeting IDs for a scheduled [personal room](https://help.webex.com/en-us/article/nul0wut/Webex-Personal-Rooms-in-Webex-Meetings) meeting.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // delete-breakout-sessions
		var meetingId string
		var sendEmail string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "delete-breakout-sessions",
			Short: "Delete Meeting Breakout Sessions",
			Long:  "Deletes breakout sessions with a specified meeting ID. The deleted breakout sessions cannot be recovered. The value of `enabledBreakoutSessions` attribute is set to `false` automatically.\nThis operation applies to meeting series and scheduled meetings. It doesn't apply to ended or in-progress meeting instances.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/meetings/{meetingId}/breakoutSessions")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("sendEmail", sendEmail)
				req.Header("hostEmail", hostEmail)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. This parameter applies to meeting series and scheduled meetings. It doesn't apply to ended or in-progress meeting instances.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&sendEmail, "send-email", "", "Whether or not to send emails to host and invitees. It is an optional field and default value is true.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "e.g. `john.andersen@example.com` (string, optional) - Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will delete breakout sessions that are created by that user.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // get-survey
		var meetingId string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "get-survey",
			Short: "Get a Meeting Survey",
			Long:  "Retrieves details for a meeting survey identified by `meetingId`.\n\n#### Request Header\n\n* `hostEmail`: Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will return survey details of that user.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/{meetingId}/survey")
				req.PathParam("meetingId", meetingId)
				req.Header("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Please note that only the meeting ID of a scheduled webinar is supported for this API.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "e.g. john.andersen@example.com")
		meetingsCmd.AddCommand(cmd)
	}

	{ // list-survey-results
		var meetingId string
		var meetingStartTimeFrom string
		var meetingStartTimeTo string
		var max string
		var timezone string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "list-survey-results",
			Short: "List Meeting Survey Results",
			Long:  "Retrieves results for a meeting survey identified by `meetingId`.\n\n#### Request Header\n\n* `timezone`: Time zone for time stamps in response body, defined in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default value is `UTC` if not specified.\n\n* `hostEmail`: Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will return the survey results of that user.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/{meetingId}/surveyResults")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("meetingStartTimeFrom", meetingStartTimeFrom)
				req.QueryParam("meetingStartTimeTo", meetingStartTimeTo)
				req.QueryParam("max", max)
				req.Header("timezone", timezone)
				req.Header("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Please note that only the meeting ID of a scheduled webinar is supported for this API.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&meetingStartTimeFrom, "meeting-start-time-from", "", "Start date and time (inclusive) in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format for the meeting objects being requested. `meetingStartTimeFrom` cannot be after `meetingStartTimeTo`. This parameter will be ignored if `meetingId` is the unique identifier for the specific meeting instance. When `meetingId` is not the unique identifier for the specific meeting instance, the `meetingStartTimeFrom`, if not specified, equals `meetingStartTimeTo` minus `1` month; if `meetingStartTimeTo` is also not specified, the default value for `meetingStartTimeFrom` is `1` month before the current date and time. ")
		cmd.Flags().StringVar(&meetingStartTimeTo, "meeting-start-time-to", "", "End date and time (exclusive) in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format for the meeting objects being requested. `meetingStartTimeTo` cannot be prior to `meetingStartTimeFrom`. This parameter will be ignored if `meetingId` is the unique identifier for the specific meeting instance. When `meetingId` is not the unique identifier for the specific meeting instance, if `meetingStartTimeFrom` is also not specified, the default value for `meetingStartTimeTo` is the current date and time;For example,if `meetingStartTimeFrom` is a month ago, the default value for `meetingStartTimeTo` is `1` month after `meetingStartTimeFrom`.Otherwise it is the current date and time. ")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of meetings in the response, up to 100. The default is 10.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "e.g. john.andersen@example.com")
		meetingsCmd.AddCommand(cmd)
	}

	{ // get-survey-links
		var meetingId string
		var timezone string
		var hostEmail string
		var meetingStartTimeFrom string
		var meetingStartTimeTo string
		var emails []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "get-survey-links",
			Short: "Get Meeting Survey Links",
			Long:  "Get survey links of a meeting for different users.\n\n#### Request Header\n\n* `timezone`: Time zone for the `meetingStartTimeFrom` and `meetingStartTimeTo` parameters and defined in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default value is `UTC` if not specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/surveyLinks")
				req.PathParam("meetingId", meetingId)
				req.Header("timezone", timezone)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("hostEmail", hostEmail)
					req.BodyString("meetingStartTimeFrom", meetingStartTimeFrom)
					req.BodyString("meetingStartTimeTo", meetingStartTimeTo)
					req.BodyStringSlice("emails", emails)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only applies to webinars. Meetings and personal room meetings are not supported.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "")
		cmd.Flags().StringVar(&meetingStartTimeFrom, "meeting-start-time-from", "", "")
		cmd.Flags().StringVar(&meetingStartTimeTo, "meeting-start-time-to", "", "")
		cmd.Flags().StringSliceVar(&emails, "emails", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // create-invitation-sources
		var meetingId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-invitation-sources",
			Short: "Create Invitation Sources",
			Long:  `Creates one or more invitation sources for a meeting.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/invitationSources")
				req.PathParam("meetingId", meetingId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the meeting ID of a scheduled webinar is supported for this API.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // list-invitation-sources
		var meetingId string
		var hostEmail string
		var personId string
		cmd := &cobra.Command{
			Use:   "list-invitation-sources",
			Short: "List Invitation Sources",
			Long:  "Lists invitation sources for a meeting.\n\n#### Request Header\n\n* `hostEmail`: Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin on-behalf-of scopes. If set, the admin may specify the email of a user in a site they manage and the API will return recording details of that user.\n\n* `personId`:  Unique identifier for the meeting host. This attribute should only be set if the user or application calling the API has the admin-level scopes. When used, the admin may specify the email of a user in a site they manage to be the meeting host.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/{meetingId}/invitationSources")
				req.PathParam("meetingId", meetingId)
				req.Header("hostEmail", hostEmail)
				req.Header("personId", personId)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the meeting ID of a scheduled webinar is supported for this API.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "e.g. john.andersen@example.com")
		cmd.Flags().StringVar(&personId, "person-id", "", "e.g. Y2lzY29zcGFyazovL3VzL1BFT1BMRS8yNWJiZjgzMS01YmU5LTRjMjUtYjRiMC05YjU5MmM4YTA4NmI")
		meetingsCmd.AddCommand(cmd)
	}

	{ // list-tracking-codes
		var siteUrl string
		var service string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "list-tracking-codes",
			Short: "List Meeting Tracking Codes",
			Long:  "Lists tracking codes on a site by a meeting host. The result indicates which tracking codes and what options can be used to create or update a meeting on the specified site.\n\n* The `options` here differ from those in the [site-level tracking codes](/docs/api/v1/tracking-codes/get-a-tracking-code) and the [user-level tracking codes](/docs/api/v1/tracking-codes/get-user-tracking-codes). It is the result of a selective combination of the two.\n\n* For a tracking code, if there is no user-level tracking code, the API returns the site-level options, and the `defaultValue` of the site-level default option is `true`. If there is a user-level tracking code, it is merged into the `options`. Meanwhile, the `defaultValue` of this user-level option is `true` and the site-level default option becomes non default.\n\n* If `siteUrl` is specified, tracking codes of the specified site will be listed; otherwise, tracking codes of the user's preferred site will be listed. All available Webex sites and the preferred sites of a user can be retrieved by [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/meetings/trackingCodes")
				req.QueryParam("siteUrl", siteUrl)
				req.QueryParam("service", service)
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
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API retrieves the tracking code from. If not specified, the API retrieves the tracking code from the user's preferred site. All available Webex sites and preferred sites of a user can be retrieved by [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		cmd.Flags().StringVar(&service, "service", "", "Service for schedule or sign-up pages.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if a user or application calling the API has the admin-level scopes. The admin may specify the email of a user on a site they manage and the API will return meeting participants of the meetings that are hosted by that user.")
		meetingsCmd.AddCommand(cmd)
	}

	{ // reassign-host
		var siteUrl string
		var hostEmail string
		var meetingIds []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "reassign-host",
			Short: "Reassign Meetings to a New Host",
			Long:  "Reassigns a list of meetings to a new host by an admin user.\n\nAll the meetings of `meetingIds` should belong to the same site, which is the `siteUrl` in the request header, if specified, or the admin user's preferred site, if not specified. All available Webex sites and the preferred sites of a user can be retrieved by [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\nIf the user of `hostEmail` is not qualified to be a host of the target site, the API returns an error with the HTTP status code `403`. If all the meetings referenced by `meetingIds` have been reassigned the new host successfully, the API returns an empty response with the HTTP status code `204`. Otherwise, if all the meetings of `meetingIds` fail or some of them fail, the API returns a \"Multi-Status\" response with status code of `207`, and individual errors for each meeting in the response body.\n\nIf a meeting already has several ended meeting instances before it's assigned to a new host, the existing ended instances are accessible to the original host and not accessible to the new host, but the new meeting instances which happen after the reassignment are only accessible to the new host.\n\nAfter the reassignment, the original host will receive an email of meeting cancellation. However, the meeting in the original host's calendar will not necessarily be removed because user's calendar, e.g. Outlook calendar, is not managed by the meeting system directly.\n\n**Note**: Only IDs of meeting series are supported for the `meetingIds`. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. To learn more about different types of meetings, please refer to [Meeting Series, Scheduled Meetings, and Meeting Instances](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances).\n\nThere are several limitations when reassigning meetings:\n\n* Users cannot assign an in-progress meeting.\n\n* Users cannot assign a meeting to a user who is not a Webex user, or an attendee who does not have host privilege.\n\n* Users cannot assign a meeting with calling/callback to a host user who does not have calling/callback privileges\n\n* Users cannot assign a meeting with session type A to a host user who does not have session type A privileges.\n\n* Users cannot assign an MC or Webinar to a new host who does not have an MC license or a Webinar license.\n\n* Users cannot assign a TC/EC1.0/SC meeting, or a meeting that is created by on-behalf to a new host.\n\n* Users can reassign hosts for meetings from third-party integrations, such as Outlook or Google. Note that this is not recommended because it may result in inconsistent data between both parties.\n\n#### Request Header\n\n* `siteUrl`: Optional request header parameter. All the meetings of `meetingIds` should belong to the site referenced by siteUrl if specified. Otherwise, the meetings should belong to the admin user's preferred sites. All available Webex sites and the preferred sites of a user can be retrieved by [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/reassignHost")
				req.Header("siteUrl", siteUrl)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("hostEmail", hostEmail)
					req.BodyStringSlice("meetingIds", meetingIds)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "e.g. example.webex.com")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "")
		cmd.Flags().StringSliceVar(&meetingIds, "meeting-ids", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // end
		var meetingId string
		var reason string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "end",
			Short: "End a Meeting",
			Long:  "Ends a meeting with a specified meeting ID. This operation applies to meeting series, scheduled meetings, and in-progress meetings. Only the meeting host, cohost, or compliance officer can end a meeting with this API.\n\n* If the `meetingId` value specified is for a scheduled meeting, the operation ends that meeting without impacting other scheduled meetings of the parent meeting series.\n\n* If the `meetingId` value specified is for a meeting series, the operation ends the current meeting occurrence.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/end")
				req.PathParam("meetingId", meetingId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("reason", reason)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting to be ended.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&reason, "reason", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // batch-approve-registrants
		var meetingId string
		var current string
		var hostEmail string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "batch-approve-registrants",
			Short: "Batch Approve Meeting Registrants",
			Long:  `Batch approve a set of registrants for a meeting by the meeting host or cohost.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/registrants/approve")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.     + Default: `false` ")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // batch-reject-registrants
		var meetingId string
		var current string
		var hostEmail string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "batch-reject-registrants",
			Short: "Batch Reject Meeting Registrants",
			Long:  `Batch reject a set of registrants for a meeting by the meeting host or cohost.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/registrants/reject")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.     + Default: `false` ")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // batch-cancel-registrants
		var meetingId string
		var current string
		var hostEmail string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "batch-cancel-registrants",
			Short: "Batch Cancel Meeting Registrants",
			Long:  `Batch cancel a set of registrants for a meeting by the meeting host or cohost.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/registrants/cancel")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.     + Default: `false` ")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

	{ // batch-delete-registrants
		var meetingId string
		var current string
		var hostEmail string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "batch-delete-registrants",
			Short: "Batch Delete Meeting Registrants",
			Long:  `Batch delete a set of registrants for a meeting by the meeting host or cohost.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/meetings/{meetingId}/registrants/bulkDelete")
				req.PathParam("meetingId", meetingId)
				req.QueryParam("current", current)
				req.QueryParam("hostEmail", hostEmail)
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
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the meeting. Only the ID of the meeting series is supported for meetingId. IDs of scheduled meetings, meeting instances, or scheduled personal room meetings are not supported. See the [Meetings Overview](/docs/meetings#meeting-series-scheduled-meetings-and-meeting-instances) for more information about meeting types.")
		cmd.MarkFlagRequired("meeting-id")
		cmd.Flags().StringVar(&current, "current", "", "Whether or not to retrieve only the current scheduled meeting of the meeting series, i.e. the meeting ready to join or start or the upcoming meeting of the meeting series. If it's `true`, return details for the current scheduled meeting of the series, i.e. the scheduled meeting ready to join or start or the upcoming scheduled meeting of the meeting series. If it's `false` or not specified, return details for the entire meeting series. This parameter only applies to meeting series.     + Default: `false` ")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the admin-level scopes. If set, the admin may specify the email of a user in a site they manage and the API will return details for a meeting that is hosted by that user.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		meetingsCmd.AddCommand(cmd)
	}

}
