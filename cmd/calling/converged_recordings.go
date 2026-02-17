package calling

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

var convergedRecordingsCmd = &cobra.Command{
	Use:   "converged-recordings",
	Short: "ConvergedRecordings commands",
}

func init() {
	cmd.CallingCmd.AddCommand(convergedRecordingsCmd)

	{ // list
		var max string
		var from string
		var to string
		var status string
		var serviceType string
		var format string
		var ownerType string
		var storageRegion string
		var locationId string
		var topic string
		var timezone string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Recordings",
			Long:  "List recordings. You can specify a date range, and the maximum number of recordings to return.\n\nThe list returned is sorted in descending order by the date and time that the recordings were created.\n\nLong result sets are split into [pages](/docs/basics#pagination).\n\nList recordings requires the `spark:recordings_read` scope.\n\nPlease use `List Recordings for Admin or Compliance Officer` API to list all recordings for a user with the role Compliance officer or Admin\n\nRequest Header\n\n* `timezone`: *[Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.*",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/convergedRecordings")
				req.QueryParam("max", max)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("status", status)
				req.QueryParam("serviceType", serviceType)
				req.QueryParam("format", format)
				req.QueryParam("ownerType", ownerType)
				req.QueryParam("storageRegion", storageRegion)
				req.QueryParam("locationId", locationId)
				req.QueryParam("topic", topic)
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
		cmd.Flags().StringVar(&status, "status", "", "Recording's status. If not specified or `available`, retrieves recordings that are available. Otherwise, if specified as `deleted`, retrieves recordings that have been moved into the recycle bin.")
		cmd.Flags().StringVar(&serviceType, "service-type", "", "Recording's service-type. If specified, the API filters recordings by service-type. Valid values are `calling` and `customerAssist`.")
		cmd.Flags().StringVar(&format, "format", "", "Recording's file format. If specified, the API filters recordings by format. Valid values are `MP3`.")
		cmd.Flags().StringVar(&ownerType, "owner-type", "", "Recording based on type of user.")
		cmd.Flags().StringVar(&storageRegion, "storage-region", "", "Recording stored in certain Webex locations.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Fetch recordings for users in a particular Webex Calling location (as configured in Control Hub).")
		cmd.Flags().StringVar(&topic, "topic", "", "Recording's topic. If specified, the API filters recordings by topic in a case-insensitive manner.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		convergedRecordingsCmd.AddCommand(cmd)
	}

	{ // list-admin-compliance-officer
		var max string
		var from string
		var to string
		var status string
		var serviceType string
		var format string
		var ownerId string
		var ownerEmail string
		var ownerType string
		var storageRegion string
		var locationId string
		var topic string
		var timezone string
		cmd := &cobra.Command{
			Use:   "list-admin-compliance-officer",
			Short: "List Recordings for Admin or Compliance officer",
			Long:  "List recordings for an admin or compliance officer. You can specify a date range, and the maximum number of recordings to return.\n\nThe list returned is sorted in descending order by the date and time that the recordings were created.\n\nLong result sets are split into [pages](/docs/basics#pagination).\n\nList recordings requires the `spark-compliance:recordings_read` scope for compliance officer and `spark-admin:recordings_read` scope for admin.\n\n#### Request Header\n\n* `timezone`: *[Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.*",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/admin/convergedRecordings")
				req.QueryParam("max", max)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("status", status)
				req.QueryParam("serviceType", serviceType)
				req.QueryParam("format", format)
				req.QueryParam("ownerId", ownerId)
				req.QueryParam("ownerEmail", ownerEmail)
				req.QueryParam("ownerType", ownerType)
				req.QueryParam("storageRegion", storageRegion)
				req.QueryParam("locationId", locationId)
				req.QueryParam("topic", topic)
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
		cmd.Flags().StringVar(&from, "from", "", "Starting date and time (inclusive) for recordings to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `from` cannot be after `to`. The interval between `from` and `to` must be within 30 days.")
		cmd.Flags().StringVar(&to, "to", "", "Ending date and time (exclusive) for List recordings to return, in any [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) compliant format. `to` cannot be before `from`. The interval between `from` and `to` must be within 30 days.")
		cmd.Flags().StringVar(&status, "status", "", "Recording's status. If not specified or `available`, retrieves recordings that are available. if specified as `purged`, retrieves recordings those are deleted but available to a compliance officer. Otherwise, if specified as `deleted`, retrieves recordings that have been moved into the recycle bin.")
		cmd.Flags().StringVar(&serviceType, "service-type", "", "Recording's service-type. If specified, the API filters recordings by service-type. Valid values are `calling` and `customerAssist`.")
		cmd.Flags().StringVar(&format, "format", "", "Recording's file format. If specified, the API filters recordings by format. Valid values are `MP3`.")
		cmd.Flags().StringVar(&ownerId, "owner-id", "", "Webex user Id to fetch recordings for a particular user.")
		cmd.Flags().StringVar(&ownerEmail, "owner-email", "", "Webex email address to fetch recordings for a particular user.")
		cmd.Flags().StringVar(&ownerType, "owner-type", "", "Recording based on type of user.")
		cmd.Flags().StringVar(&storageRegion, "storage-region", "", "Recording stored in certain Webex locations.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Fetch recordings for users in a particular Webex Calling location (as configured in Control Hub).")
		cmd.Flags().StringVar(&topic, "topic", "", "Recording's topic. If specified, the API filters recordings by topic in a case-insensitive manner.")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		convergedRecordingsCmd.AddCommand(cmd)
	}

	{ // get
		var recordingId string
		var timezone string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Recording Details",
			Long:  "Retrieves details for a recording with a specified recording ID.\n\nOnly recordings of owner with the authenticated user may be retrieved.\n\nGet Recording Details requires the `spark-compliance:recordings_read` scope for compliance officer, `spark-admin:recordings_read` scope for admin and `spark:recordings_read` scope for user.\n\n#### Request Header\n\n* `timezone`: *[Time zone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) in conformance with the [IANA time zone database](https://www.iana.org/time-zones). The default is UTC if `timezone` is not defined.*",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/convergedRecordings/{recordingId}")
				req.PathParam("recordingId", recordingId)
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
		cmd.MarkFlagRequired("recording-id")
		cmd.Flags().StringVar(&timezone, "timezone", "", "e.g. UTC")
		convergedRecordingsCmd.AddCommand(cmd)
	}

	{ // delete
		var recordingId string
		var reason string
		var comment string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Recording",
			Long:  "Removes a recording with a specified recording ID. The deleted recording cannot be recovered.\n\nIf a Compliance Officer deletes another user's recording, the recording will be inaccessible to regular users (host, attendees and shared), and to the Compliance officer as well. This action purges the recordings from Webex.\n\nDelete a Recording requires the `spark-compliance:recordings_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/convergedRecordings/{recordingId}")
				req.PathParam("recordingId", recordingId)
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
		cmd.Flags().StringVar(&reason, "reason", "", "")
		cmd.Flags().StringVar(&comment, "comment", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		convergedRecordingsCmd.AddCommand(cmd)
	}

	{ // get-metadata
		var recordingId string
		var showAllTypes string
		cmd := &cobra.Command{
			Use:   "get-metadata",
			Short: "Get Recording metadata",
			Long:  "Retrieves metadata details for a recording with a specified recording ID. The recording must be owned by the authenticated user.\n\nFor information on the metadata fields, refer to [Metadata Guide](https://developer.webex.com/docs/api/guides/consolidated-metadata-documentation-and-samples-guide)\n\nGet Recording metadata requires the `spark-compliance:recordings_read` scope for compliance officer, `spark-admin:recordings_read` for admin and `spark:recordings_read` for user.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/convergedRecordings/{recordingId}/metadata")
				req.PathParam("recordingId", recordingId)
				req.QueryParam("showAllTypes", showAllTypes)
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
		cmd.Flags().StringVar(&showAllTypes, "show-all-types", "", "If `showAllTypes` is `true`, all attributes will be shown. If it's `false` or not specified, the following attributes of the metadata will be hidden.                                           serviceData.callActivity.mediaStreams                                           serviceData.callActivity.participants                                           serviceData.callActivity.redirectInfo                                           serviceData.callActivity.redirectedCall ")
		convergedRecordingsCmd.AddCommand(cmd)
	}

	{ // reassign
		var reassignOwnerEmail string
		var ownerEmail string
		var ownerId string
		var recordingIds []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "reassign",
			Short: "Reassign Recordings",
			Long:  "Reassigns recordings to a new user. As an administrator, you can reassign a list of recordings or all recordings of a particular user to a new user.\nThe recordings can belong to an org user, a virtual line, or a workspace, but the destination user should only be a valid org user.\n\n* For a org user either `ownerEmail` or `recordingIds` or both must be provided.\n\n* For a virtual line or a workspace, `ownerID` or `recordingIds` or both must be provided.\n\n* If `recordingIds` and `ownerID` is empty but `ownerEmail` is provided, all recordings owned by the `ownerEmail` are reassigned to `reassignOwnerEmail`.\n\n* If `recordingIds` is provided and `ownerEmail` or `ownerID` is also provided, only the recordings specified by `recordingIds` that are owned by `ownerEmail` or `ownerID` are reassigned to `reassignOwnerEmail`.\n\n* If `ownerEmail` and `ownerID` is empty but `recordingIds` is provided, the recordings specified by `recordingIds` are reassigned to `reassignOwnerEmail` regardless of the current owner.\n\n* If both `ownerId` and `ownerEmail` are passed along with `recordingIds`, only the recordings specified by `recordingIds` that are owned by `ownerEmail` are reassigned to `reassignOwnerEmail`.\n\n* If `recordingIds` is empty but both `ownerId` and `ownerEmail` is provided, all recordings owned by the `ownerEmail` are reassigned to `reassignOwnerEmail`.\n\nThe `spark-admin:recordings_write` scope is required to reassign recordings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/convergedRecordings/reassign")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("reassignOwnerEmail", reassignOwnerEmail)
					req.BodyString("ownerEmail", ownerEmail)
					req.BodyString("ownerID", ownerId)
					req.BodyStringSlice("recordingIds", recordingIds)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&reassignOwnerEmail, "reassign-owner-email", "", "")
		cmd.Flags().StringVar(&ownerEmail, "owner-email", "", "")
		cmd.Flags().StringVar(&ownerId, "owner-id", "", "")
		cmd.Flags().StringSliceVar(&recordingIds, "recording-ids", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		convergedRecordingsCmd.AddCommand(cmd)
	}

	{ // move-recycle-bin
		var trashAll bool
		var ownerEmail string
		var recordingIds []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "move-recycle-bin",
			Short: "Move Recordings into the Recycle Bin",
			Long:  "Move recordings into the recycle bin with recording IDs or move all the recordings to the recycle bin.\n\nOnly the following two entities can use this API\n\n* Administrator: A user or an application with the scope `spark-admin:recordings_write`.\n\n* User: An authenticated user who does not have the scope `spark-admin:recordings_write` but has `spark:recordings_write`.\n\nAs an `administrator`, you can move a list of recordings or all recordings of a particular user within the org you manage to the recycle bin.\n\nAs a `user`, you can move a list of your own recordings or all your recordings to the recycle bin.\n\nRecordings in the recycle bin can be recovered by [Restore Recordings from Recycle Bin](/docs/api/v1/converged-recordings/restore-recordings-from-recycle-bin) API. If you'd like to empty recordings from the recycle bin, you can use [Purge Recordings from Recycle Bin](/docs/api/v1/converged-recordings/purge-recordings-from-recycle-bin) API to purge all or some of them.\n\n* If `trashAll` is `true`:\n  * `recordingIds` should be empty.\n  * If the caller of this API is an `administrator`, `ownerEmail` should not be empty and all recordings owned by the `ownerEmail` will be moved to the recycle bin.\n  * If the caller of this API is a `user`, `ownerEmail` should be empty and all recordings owned by the caller will be moved to the recycle bin.\n\n* If `trashAll` is `false`:\n  * `ownerEmail` should be empty.\n  * `recordingIds` should not be empty and its maximum size is `100`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/convergedRecordings/softDelete")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("trashAll", trashAll, cmd.Flags().Changed("trash-all"))
					req.BodyString("ownerEmail", ownerEmail)
					req.BodyStringSlice("recordingIds", recordingIds)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().BoolVar(&trashAll, "trash-all", false, "")
		cmd.Flags().StringVar(&ownerEmail, "owner-email", "", "")
		cmd.Flags().StringSliceVar(&recordingIds, "recording-ids", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		convergedRecordingsCmd.AddCommand(cmd)
	}

	{ // restore-recycle-bin
		var restoreAll bool
		var ownerEmail string
		var recordingIds []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "restore-recycle-bin",
			Short: "Restore Recordings from Recycle Bin",
			Long:  "Restore recordings from the recycle bin with recording IDs or restore all the recordings that are in the recycle bin.\n\nOnly the following two entities can use this API\n\n* Administrator: A user or an application with the scope `spark-admin:recordings_write`.\n\n* User: An authenticated user who does not have the scope `spark-admin:recordings_write` but has `spark:recordings_write`.\n\nAs an `administrator`, you can restore a list of recordings or all recordings of a particular user within the org you manage from the recycle bin.\n\nAs a `user`, you can restore a list of your own recordings or all your recordings from the recycle bin.\n\n* If `restoreAll` is `true`:\n  * `recordingIds` should be empty.\n  * If the caller of this API is an `administrator`, `ownerEmail` should not be empty and all recordings owned by the `ownerEmail` will be restored from the recycle bin.\n  * If the caller of this API is a `user`, `ownerEmail` should be empty and all recordings owned by the caller will be restored from the recycle bin.\n\n* If `restoreAll` is `false`:\n  * `ownerEmail` should be empty.\n  * `recordingIds` should not be empty and its maximum size is `100`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/convergedRecordings/restore")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("restoreAll", restoreAll, cmd.Flags().Changed("restore-all"))
					req.BodyString("ownerEmail", ownerEmail)
					req.BodyStringSlice("recordingIds", recordingIds)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().BoolVar(&restoreAll, "restore-all", false, "")
		cmd.Flags().StringVar(&ownerEmail, "owner-email", "", "")
		cmd.Flags().StringSliceVar(&recordingIds, "recording-ids", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		convergedRecordingsCmd.AddCommand(cmd)
	}

	{ // purge-recycle-bin
		var purgeAll bool
		var ownerEmail string
		var recordingIds []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "purge-recycle-bin",
			Short: "Purge Recordings from Recycle Bin",
			Long:  "Purge recordings from the recycle bin matching the supplied recording IDs, or purge all the recordings that are in the recycle bin. A recording, once purged, cannot be restored.\n\nOnly the following two entities can use this API\n\n* Administrator: A user or an application with the scope `spark-admin:recordings_write`.\n\n* User: An authenticated user who does not have the scope `spark-admin:recordings_write` but has `spark:recordings_write`.\n\nAs an `administrator`, you can purge a list of recordings or all recordings of a particular user within the org you manage from the recycle bin.\n\nAs a `user`, you can purge a list of your own recordings or all your recordings from the recycle bin.\n\n* If `purgeAll` is `true`:\n  * `recordingIds` should be empty.\n  * If the caller of this API is an `administrator`, `ownerEmail` should not be empty and all recordings owned the `ownerEmail` will be purged from the recycle bin.\n  * If the caller of this API is a `user`, `ownerEmail` should be empty and all recordings owned by the caller will be purged from the recycle bin.\n\n* If `purgeAll` is `false`:\n  * `ownerEmail` should be empty.\n  * `recordingIds` should not be empty and its maximum size is `100`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/convergedRecordings/purge")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("purgeAll", purgeAll, cmd.Flags().Changed("purge-all"))
					req.BodyString("ownerEmail", ownerEmail)
					req.BodyStringSlice("recordingIds", recordingIds)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().BoolVar(&purgeAll, "purge-all", false, "")
		cmd.Flags().StringVar(&ownerEmail, "owner-email", "", "")
		cmd.Flags().StringSliceVar(&recordingIds, "recording-ids", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		convergedRecordingsCmd.AddCommand(cmd)
	}

}
