package meetings

import (
	"fmt"
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
var _ = strings.Join

var recordingsCmd = &cobra.Command{
	Use:   "recordings",
	Short: "Recordings commands",
}

func init() {
	cmd.MeetingsCmd.AddCommand(recordingsCmd)

	{ // list
		var max string
		var from string
		var to string
		var meetingId string
		var hostEmail string
		var siteUrl string
		var integrationTag string
		var topic string
		var format string
		var serviceType string
		var status string
		var timezone string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Recordings",
			Long:  "Lists recordings. You can specify a date range, a parent meeting ID, and the maximum number of recordings to return.\n\nOnly recordings of meetings hosted by or shared with the authenticated user will be listed.\n\nThe list returned is sorted in descending order by the date and time that the recordings were created.\n\nLong result sets are split into [pages](/docs/basics#pagination).\n\n* If `meetingId` is specified, only recordings associated with the specified meeting will be listed. **NOTE**: when `meetingId` is specified, parameter of `siteUrl` will be ignored.\n\n* If `siteUrl` is specified, recordings of the specified site will be listed; otherwise, the API lists recordings of all the user's sites. All available Webex sites and preferred site of the user can be retrieved by [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\n#### Request Header\n\n* `timezone`: *[Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.*",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/recordings")
				req.QueryParam("max", max)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("meetingId", meetingId)
				req.QueryParam("hostEmail", hostEmail)
				req.QueryParam("siteUrl", siteUrl)
				req.QueryParam("integrationTag", integrationTag)
				req.QueryParam("topic", topic)
				req.QueryParam("format", format)
				req.QueryParam("serviceType", serviceType)
				req.QueryParam("status", status)
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
		cmd.Flags().StringVar(&max, "max", "", "Maximum number of recordings to return in a single page. `max` must be equal to or greater than `1` and equal to or less than `100`.")
		cmd.Flags().StringVar(&from, "from", "", "Starting date and time (inclusive) for recordings to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `from` cannot be after `to`.")
		cmd.Flags().StringVar(&to, "to", "", "Ending date and time (exclusive) for List recordings to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `to` cannot be before `from`.")
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the parent meeting series, scheduled meeting, or meeting instance for which recordings are being requested. If a meeting series ID is specified, the operation returns an array of recordings for the specified meeting series. If a scheduled meeting ID is specified, the operation returns an array of recordings for the specified scheduled meeting. If a meeting instance ID is specified, the operation returns an array of recordings for the specified meeting instance. If no ID is specified, the operation returns an array of recordings for all meetings of the current user. When `meetingId` is specified, the `siteUrl` parameter is ignored.")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the required [admin-level meeting scopes](/docs/meetings#adminorganization-level-authentication-and-scopes). If set, the admin may specify the email of a user in a site they manage and the API will return recordings of that user.")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site from which the API lists recordings. If not specified, the API lists recordings from all of a user's sites. All available Webex sites and the preferred site of the user can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		cmd.Flags().StringVar(&integrationTag, "integration-tag", "", "External key of the parent meeting created by an integration application. This parameter is used by the integration application to query recordings by a key in its own domain, such as a Zendesk ticket ID, a Jira ID, a Salesforce Opportunity ID, etc. An integrationTag created by one client cannot be accessed or used as a filtering parameter by another client. For example, if a meeting has an `integrationTag` of \"Sales\" which is created by the client behind the developer portal, then this integrationTag can't be accessed on the meeting or its recordings by another client. Neither can it be used to filter meetings or recordings by a client other than the one that created the integrationTag of \"Sales\".")
		cmd.Flags().StringVar(&topic, "topic", "", "Recording's topic. If specified, the API filters recordings by topic in a case-insensitive manner.")
		cmd.Flags().StringVar(&format, "format", "", "Recording's file format. If specified, the API filters recordings by format.")
		cmd.Flags().StringVar(&serviceType, "service-type", "", "The service type for recordings. If this item is specified, the API filters recordings by service-type.")
		cmd.Flags().StringVar(&status, "status", "", "Recording's status. If not specified or `available`, retrieves recordings that are available. Otherwise, if specified as `deleted`, retrieves recordings that have been moved into the recycle bin. The `purged` status only applies if the user calling the API is a Compliance Officer and `meetingId` is specified.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		recordingsCmd.AddCommand(cmd)
	}

	{ // list-admin-compliance-officer
		var max string
		var from string
		var to string
		var meetingId string
		var siteUrl string
		var integrationTag string
		var topic string
		var format string
		var serviceType string
		var status string
		var timezone string
		cmd := &cobra.Command{
			Use:   "list-admin-compliance-officer",
			Short: "List Recordings For an Admin or Compliance Officer",
			Long:  "List recordings for an admin or compliance officer. You can specify a date range, a parent meeting ID, and the maximum number of recordings to return.\n\nThe list returned is sorted in descending order by the date and time that the recordings were created.\n\nLong result sets are split into [pages](/docs/basics#pagination).\n\n* If `meetingId` is specified, only recordings associated with the specified meeting will be listed. Please note that when `meetingId` is specified, parameters of `siteUrl`, `from`, and `to` will be ignored.\n\n* If `siteUrl` is specified, all the recordings on the specified site are listed; otherwise, all the recordings on the admin user's or compliance officer's preferred site are listed. All the available Webex sites and the admin user's or compliance officer's preferred site can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\n#### Request Header\n\n* `timezone`: *[Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.*",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/admin/recordings")
				req.QueryParam("max", max)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("meetingId", meetingId)
				req.QueryParam("siteUrl", siteUrl)
				req.QueryParam("integrationTag", integrationTag)
				req.QueryParam("topic", topic)
				req.QueryParam("format", format)
				req.QueryParam("serviceType", serviceType)
				req.QueryParam("status", status)
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
		cmd.Flags().StringVar(&max, "max", "", "Maximum number of recordings to return in a single page. `max` must be equal to or greater than `1` and equal to or less than `100`.")
		cmd.Flags().StringVar(&from, "from", "", "Starting date and time (inclusive) for recordings to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `from` cannot be after `to`. The interval between `from` and `to` must be within 30 days. If `to` is specified, the default value for `from` is `to` minus 7 days. If `to` is also not specified, the default value for `from` is current date and time minus 7 days.")
		cmd.Flags().StringVar(&to, "to", "", "Ending date and time (exclusive) for List recordings to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `to` cannot be before `from`. The interval between `from` and `to` must be within 30 days. If `from` is specified, the default value for `to` is `from` plus 7 days. If `from` is also not specified, the default value for `to` is the current date and time.")
		cmd.Flags().StringVar(&meetingId, "meeting-id", "", "Unique identifier for the parent meeting series, scheduled meeting, or meeting instance for which recordings are being requested. If a meeting series ID is specified, the operation returns an array of recordings for the specified meeting series. If a scheduled meeting ID is specified, the operation returns an array of recordings for the specified scheduled meeting. If a meeting instance ID is specified, the operation returns an array of recordings for the specified meeting instance. If not specified, the operation returns an array of recordings for all the current user's meetings. When `meetingId` is specified, the `siteUrl` parameter is ignored.")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API lists recordings from. If not specified, the API lists recordings from user's preferred site. All available Webex sites and preferred site of the user can be retrieved by [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		cmd.Flags().StringVar(&integrationTag, "integration-tag", "", "External key of the parent meeting created by an integration application. This parameter is used by the integration application to query recordings by a key in its own domain such as a Zendesk ticket ID, a Jira ID, a Salesforce Opportunity ID, etc. An integrationTag created by one client cannot be accessed or used as a filtering parameter by another client. For example, if a meeting has an `integrationTag` of \"Sales\" which is created by the client behind the developer portal, then this integrationTag can't be accessed on the meeting or its recordings by another client. Neither can it be used to filter meetings or recordings by a client other than the one that created the integrationTag of \"Sales\".")
		cmd.Flags().StringVar(&topic, "topic", "", "Recording topic. If specified, the API filters recordings by topic in a case-insensitive manner.")
		cmd.Flags().StringVar(&format, "format", "", "Recording's file format. If specified, the API filters recordings by format.")
		cmd.Flags().StringVar(&serviceType, "service-type", "", "The service type for recordings. If specified, the API filters recordings by service type.")
		cmd.Flags().StringVar(&status, "status", "", "Recording's status. If not specified or `available`, retrieves recordings that are available. If specified as `deleted`, retrieves recordings that have been moved to the recycle bin. Otherwise, if specified as `purged`, retrieves recordings that have been purged from the recycle bin.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		recordingsCmd.AddCommand(cmd)
	}

	{ // delete-admin
		var recordingId string
		cmd := &cobra.Command{
			Use:   "delete-admin",
			Short: "Delete a Recording By an Admin",
			Long:  "Removes a recording with a specified recording ID by an admin. The deleted recording cannot be recovered. It will be inaccessible to regular users (host, attendees and shared) or an admin, but it will be still available to the Compliance Officer.\n\nAny recording on a site which is managed by the admin can be deleted by him.\n\nThe `temporaryDirectDownloadLinks` of a recording which are retrieved by the [Get Recording Details](/docs/api/v1/recordings/get-recording-details) API are still available to Compliance Officers even if the recording has been deleted.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/admin/recordings/{recordingId}")
				req.PathParam("recordingId", recordingId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&recordingId, "recording-id", "", "A unique identifier for the recording.")
		cmd.MarkFlagRequired("recording-id")
		recordingsCmd.AddCommand(cmd)
	}

	{ // get
		var recordingId string
		var hostEmail string
		var timezone string
		var siteUrl string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Recording Details",
			Long:  "Retrieves details for a recording with a specified recording ID.\n\nOnly recordings of meetings hosted by or shared with the authenticated user may be retrieved.\n\n#### Request Header\n\n* `timezone`: *[Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.*\n\n* `siteUrl`: Optional request header parameter. If specified, retrieve the recording details from that site; otherwise, retrieve it from the site which is implied based on the recording ID.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/recordings/{recordingId}")
				req.PathParam("recordingId", recordingId)
				req.QueryParam("hostEmail", hostEmail)
				req.Header("timezone", timezone)
				req.Header("siteUrl", siteUrl)
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
		cmd.MarkFlagRequired("recording-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. Only used if the user or application calling the API has required [admin-level meeting scopes](/docs/meetings#adminorganization-level-authentication-and-scopes). If set, the admin may specify the email of a user in a site they manage, and the API will return recording details of that user.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "e.g. example.webex.com")
		recordingsCmd.AddCommand(cmd)
	}

	{ // delete
		var recordingId string
		var hostEmail string
		var reason string
		var comment string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Recording",
			Long:  "Removes a recording with a specified recording ID. The deleted recording cannot be recovered. If a Compliance Officer deletes another user's recording, the recording will be inaccessible to regular users (host, attendees and shared), but will be still available to the Compliance Officer.\n\nOnly recordings of meetings hosted by the authenticated user can be deleted.\n\nThe `temporaryDirectDownloadLinks` of a recording which are retrieved by the [Get Recording Details](/docs/api/v1/recordings/get-recording-details) API are still available to Compliance Officers even if the recording has been deleted.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/recordings/{recordingId}")
				req.PathParam("recordingId", recordingId)
				req.QueryParam("hostEmail", hostEmail)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("reason", reason)
					req.BodyString("comment", comment)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&recordingId, "recording-id", "", "A unique identifier for the recording.")
		cmd.MarkFlagRequired("recording-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. Only used if the user or application calling the API has the required [admin-level meeting scopes](/docs/meetings#adminorganization-level-authentication-and-scopes). If set, the admin may specify the email of a user in a site they manage and the API will delete a recording of that user.")
		cmd.Flags().StringVar(&reason, "reason", "", "")
		cmd.Flags().StringVar(&comment, "comment", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		recordingsCmd.AddCommand(cmd)
	}

	{ // move-recycle-bin
		var hostEmail string
		var recordingIds []string
		var siteUrl string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "move-recycle-bin",
			Short: "Move Recordings into the Recycle Bin",
			Long:  "Move recordings into the recycle bin with recording IDs. Recordings in the recycle bin can be recovered by [Restore Recordings from Recycle Bin](/docs/api/v1/recordings/restore-recordings-from-recycle-bin) API. If you'd like to empty recordings from the recycle bin, you can use [Purge Recordings from Recycle Bin](/docs/api/v1/recordings/purge-recordings-from-recycle-bin) API to purge all or some of them.\n\nOnly recordings of meetings hosted by the authenticated user can be moved into the recycle bin.\n\n* `recordingIds` should not be empty and its maximum size is `100`.\n\n* All the IDs of `recordingIds` should belong to the site of `siteUrl` or the user's preferred site if `siteUrl` is not specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/recordings/softDelete")
				req.QueryParam("hostEmail", hostEmail)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("recordingIds", recordingIds)
					req.BodyString("siteUrl", siteUrl)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. Only used if the user or application calling the API has the required [admin-level meeting scopes](/docs/meetings#adminorganization-level-authentication-and-scopes). If set, the admin may specify the email of a user in a site they manage and the API will move recordings into recycle bin of that user")
		cmd.Flags().StringSliceVar(&recordingIds, "recording-ids", nil, "")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		recordingsCmd.AddCommand(cmd)
	}

	{ // restore-recycle-bin
		var hostEmail string
		var restoreAll bool
		var recordingIds []string
		var siteUrl string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "restore-recycle-bin",
			Short: "Restore Recordings from Recycle Bin",
			Long:  "Restore all or some recordings from the recycle bin. Only recordings of meetings hosted by the authenticated user can be restored from recycle bin.\n\n* If `restoreAll` is `true`, `recordingIds` should be empty.\n\n* If `restoreAll` is `false`, `recordingIds` should not be empty and its maximum size is `100`.\n\n* All the IDs of `recordingIds` should belong to the site of `siteUrl` or the user's preferred site if `siteUrl` is not specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/recordings/restore")
				req.QueryParam("hostEmail", hostEmail)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("restoreAll", restoreAll, cmd.Flags().Changed("restore-all"))
					req.BodyStringSlice("recordingIds", recordingIds)
					req.BodyString("siteUrl", siteUrl)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. This parameter is only used if the user or application calling the API has the required [admin-level meeting scopes](/docs/meetings#adminorganization-level-authentication-and-scopes). If set, the admin may specify the email of a user in a site they manage and the API will restore recordings of that user.")
		cmd.Flags().BoolVar(&restoreAll, "restore-all", false, "")
		cmd.Flags().StringSliceVar(&recordingIds, "recording-ids", nil, "")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		recordingsCmd.AddCommand(cmd)
	}

	{ // purge-recycle-bin
		var hostEmail string
		var purgeAll bool
		var recordingIds []string
		var siteUrl string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "purge-recycle-bin",
			Short: "Purge Recordings from Recycle Bin",
			Long:  "Purge recordings from recycle bin with recording IDs or purge all the recordings that are in the recycle bin.\n\nOnly recordings of meetings hosted by the authenticated user can be purged from recycle bin.\n\n* If `purgeAll` is `true`, `recordingIds` should be empty.\n\n* If `purgeAll` is `false`, `recordingIds` should not be empty and its maximum size is `100`.\n\n* All the IDs of `recordingIds` should belong to the site of `siteUrl` or the user's preferred site if `siteUrl` is not specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/recordings/purge")
				req.QueryParam("hostEmail", hostEmail)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("purgeAll", purgeAll, cmd.Flags().Changed("purge-all"))
					req.BodyStringSlice("recordingIds", recordingIds)
					req.BodyString("siteUrl", siteUrl)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email address for the meeting host. Only used if the user or application calling the API has the required [admin-level meeting scopes](/docs/meetings#adminorganization-level-authentication-and-scopes). If set, the admin may specify the email of a user in a site they manage and the API will purge recordings from recycle bin of that user.")
		cmd.Flags().BoolVar(&purgeAll, "purge-all", false, "")
		cmd.Flags().StringSliceVar(&recordingIds, "recording-ids", nil, "")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		recordingsCmd.AddCommand(cmd)
	}

	{ // share
		var recordingId string
		var hostEmail string
		var addEmails []string
		var removeEmails []string
		var sendEmail bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "share",
			Short: "Share a Recording",
			Long:  `Share or unshare a recording with other users by recording ID and email addresses.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/recordings/{recordingId}/accessList")
				req.PathParam("recordingId", recordingId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("hostEmail", hostEmail)
					req.BodyStringSlice("addEmails", addEmails)
					req.BodyStringSlice("removeEmails", removeEmails)
					req.BodyBool("sendEmail", sendEmail, cmd.Flags().Changed("send-email"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&recordingId, "recording-id", "", "A unique identifier for the recording.")
		cmd.MarkFlagRequired("recording-id")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "")
		cmd.Flags().StringSliceVar(&addEmails, "add-emails", nil, "")
		cmd.Flags().StringSliceVar(&removeEmails, "remove-emails", nil, "")
		cmd.Flags().BoolVar(&sendEmail, "send-email", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		recordingsCmd.AddCommand(cmd)
	}

	{ // list-group
		var personId string
		var max string
		var from string
		var to string
		var siteUrl string
		var integrationTag string
		var topic string
		var format string
		var serviceType string
		var timezone string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "list-group",
			Short: "List Group Recordings",
			Long:  "List group recordings for a service app which has group recording access. You can specify a date range, the maximum number of recordings to return, and `personId` or `hostEmail` of whom the recordings will be retrieved.\n\n* The list returned is sorted in descending order by the date and time that the recordings were created.\n\n* Only recordings which are in the `available` status and not shared by others can be listed. Those in the `deleted` or `purged` status or shared by others can't be listed.\n\n* Long result sets are split into [pages](/docs/basics#pagination).\n\n* If `siteUrl` is specified, the API lists group recordings on the specified site; otherwise, the API lists group recordings on all the sites managed by the service app. All the sites managed by a service app can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\n* One of the `personId` parameter and `hostEmail` header must be specified so that only recordings of meetings hosted by the person of `personId` or `hostEmail` will be returned.\n\n#### Request Header\n\n* `timezone`: [Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.\n\n* `hostEmail`: Email of the user whose recordings will be retrieved. The `hostEmail` parameter is optional, but one of the `personId` parameter and `hostEmail` header must be specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/group/recordings")
				req.QueryParam("personId", personId)
				req.QueryParam("max", max)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("siteUrl", siteUrl)
				req.QueryParam("integrationTag", integrationTag)
				req.QueryParam("topic", topic)
				req.QueryParam("format", format)
				req.QueryParam("serviceType", serviceType)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Person ID of the user whose recordings will be retrieved. The person ID can be retrieved from the [People APIs](/docs/api/v1/people), e.g. [Lit People](/docs/api/v1/people/list-people). Note that a person ID retrieved from the People APIs is a Base64-encoded string, e.g. `Y2lzY29zcGFyazovL3VzL1BFT1BMRS9kNDdiMmU3ZC01ZTBmLTRmNjktYWVmNC1lNGZmOTBhZWE3Yzk`. The person ID in the raw UUID format which is the last part of the Base64-decoded string, e.g. `d47b2e7d-5e0f-4f69-aef4-e4ff90aea7c9`, is also supported. The `personId` parameter is optional, but one of the `personId` parameter and `hostEmail` header must be specified.")
		cmd.Flags().StringVar(&max, "max", "", "Maximum number of recordings to return in a single page. `max` must be equal to or greater than `1` and equal to or less than `100`.")
		cmd.Flags().StringVar(&from, "from", "", "Starting date and time (inclusive) for recordings to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `from` cannot be after `to`. The interval between `from` and `to` must be within 30 days.")
		cmd.Flags().StringVar(&to, "to", "", "Ending date and time (exclusive) for List recordings to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `to` cannot be before `from`. The interval between `from` and `to` must be within 30 days.")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API lists recordings from. If not specified, the API lists recordings from user's preferred site. All available Webex sites and preferred site of the user can be retrieved by [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		cmd.Flags().StringVar(&integrationTag, "integration-tag", "", "External key of the parent meeting created by an integration application. This parameter is used by the integration application to query recordings by a key in its own domain such as a Zendesk ticket ID, a Jira ID, a Salesforce Opportunity ID, etc. An integrationTag created by one client cannot be accessed or used as a filtering parameter by another client. For example, if a meeting has an `integrationTag` of \"Sales\" which is created by the client behind the developer portal, then this integrationTag can't be accessed on the meeting or its recordings by another client. Neither can it be used to filter meetings or recordings by a client other than the one that created the integrationTag of \"Sales\".")
		cmd.Flags().StringVar(&topic, "topic", "", "Recording topic. If specified, the API filters recordings by topic in a case-insensitive manner.")
		cmd.Flags().StringVar(&format, "format", "", "Recording's file format. If specified, the API filters recordings by format.")
		cmd.Flags().StringVar(&serviceType, "service-type", "", "The service type for recordings. If specified, the API filters recordings by service type.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email of the user whose recordings will be retrieved. The `hostEmail` parameter is optional, but one of the `personId` parameter and `hostEmail` header must be specified.")
		recordingsCmd.AddCommand(cmd)
	}

	{ // get-group
		var recordingId string
		var personId string
		var timezone string
		var hostEmail string
		cmd := &cobra.Command{
			Use:   "get-group",
			Short: "Get Group Recording Details",
			Long:  "Retrieves details for a group recording for a service app which has group recording access.\n\n* Only recordings which are in the `available` status and not shared by others can be listed. Those in the `deleted` or `purged` status or shared by others can't be listed.\n\n* One of the `personId` parameter and `hostEmail` header must be specified so that only recordings of meetings hosted by the person of `personId` or `hostEmail` will be returned.\n\n#### Request Header\n\n* `timezone`: *[Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.*\n\n* `hostEmail`: Email of the user whose recordings will be retrieved. The `hostEmail` parameter is optional, but one of the `personId` parameter and `hostEmail` header must be specified.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/group/recordings/{recordingId}")
				req.PathParam("recordingId", recordingId)
				req.QueryParam("personId", personId)
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
		cmd.Flags().StringVar(&recordingId, "recording-id", "", "A unique identifier for the recording.")
		cmd.MarkFlagRequired("recording-id")
		cmd.Flags().StringVar(&personId, "person-id", "", "Person ID of the user whose recordings will be retrieved. The person ID can be retrieved from the [People APIs](/docs/api/v1/people), e.g. [Lit People](/docs/api/v1/people/list-people). Note that a person ID retrieved from the People APIs is a Base64-encoded string, e.g. `Y2lzY29zcGFyazovL3VzL1BFT1BMRS9kNDdiMmU3ZC01ZTBmLTRmNjktYWVmNC1lNGZmOTBhZWE3Yzk`. The person ID in the raw UUID format which is the last part of the Base64-decoded string, e.g. `d47b2e7d-5e0f-4f69-aef4-e4ff90aea7c9`, is also supported. The `personId` parameter is optional, but one of the `personId` parameter and `hostEmail` header must be specified.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "Email of the user whose recordings will be retrieved. The `hostEmail` parameter is optional, but one of the `personId` parameter and `hostEmail` header must be specified.")
		recordingsCmd.AddCommand(cmd)
	}

	{ // share-link
		var hostEmail string
		var webShareLink string
		var addEmails []string
		var removeEmails []string
		var sendEmail bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "share-link",
			Short: "Share a Recording Link",
			Long:  `Share or unshare a recording with other users by recording link and email addresses.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/recordings/accessList")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("hostEmail", hostEmail)
					req.BodyString("webShareLink", webShareLink)
					req.BodyStringSlice("addEmails", addEmails)
					req.BodyStringSlice("removeEmails", removeEmails)
					req.BodyBool("sendEmail", sendEmail, cmd.Flags().Changed("send-email"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&hostEmail, "host-email", "", "")
		cmd.Flags().StringVar(&webShareLink, "web-share-link", "", "")
		cmd.Flags().StringSliceVar(&addEmails, "add-emails", nil, "")
		cmd.Flags().StringSliceVar(&removeEmails, "remove-emails", nil, "")
		cmd.Flags().BoolVar(&sendEmail, "send-email", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		recordingsCmd.AddCommand(cmd)
	}

}
