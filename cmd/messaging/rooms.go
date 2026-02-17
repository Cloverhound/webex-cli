package messaging

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

var roomsCmd = &cobra.Command{
	Use:   "rooms",
	Short: "Rooms commands",
}

func init() {
	cmd.MessagingCmd.AddCommand(roomsCmd)

	{ // list
		var teamId string
		var typeVal string
		var orgPublicSpaces string
		var from string
		var to string
		var sortBy string
		var max string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Rooms",
			Long:  "List rooms to which the authenticated user belongs to.\n\nThe `title` of the room for 1:1 rooms will be the display name of the other person. Please use the [memberships API](https://developer.webex.com/docs/api/v1/memberships) to list the people in the space.\n\nLong result sets will be split into [pages](/docs/basics#pagination).\n\nKnown Limitations:\nThe underlying database does not support natural sorting by `lastactivity` and will only sort on limited set of results, which are pulled from the database in order of `roomId`. For users or bots in more than 3000 spaces this can result in anomalies such as spaces that have had recent activity not being returned in the results when sorting by `lastacivity`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/rooms")
				req.QueryParam("teamId", teamId)
				req.QueryParam("type", typeVal)
				req.QueryParam("orgPublicSpaces", orgPublicSpaces)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
				req.QueryParam("sortBy", sortBy)
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
		cmd.Flags().StringVar(&teamId, "team-id", "", "List rooms associated with a team, by ID. Cannot be set in combination with `orgPublicSpaces`.")
		cmd.Flags().StringVar(&typeVal, "type", "", "List rooms by type. Cannot be set in combination with `orgPublicSpaces`.")
		cmd.Flags().StringVar(&orgPublicSpaces, "org-public-spaces", "", "Shows the org's public spaces joined and unjoined. When set the result list is sorted by the `madePublic` timestamp.")
		cmd.Flags().StringVar(&from, "from", "", "Filters rooms, that were made public after this time. See `madePublic` timestamp")
		cmd.Flags().StringVar(&to, "to", "", "Filters rooms, that were made public before this time. See `maePublic` timestamp")
		cmd.Flags().StringVar(&sortBy, "sort-by", "", "Sort results. Cannot be set in combination with `orgPublicSpaces`.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of rooms in the response. Value must be between 1 and 1000, inclusive.")
		roomsCmd.AddCommand(cmd)
	}

	{ // create
		var title string
		var teamId string
		var classificationId string
		var isLocked bool
		var isPublic bool
		var description string
		var isAnnouncementOnly bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Room",
			Long:  "Creates a room. The authenticated user is automatically added as a member of the room. See the [Memberships API](/docs/api/v1/memberships) to learn how to add more people to the room.\n\nTo create a 1:1 room, use the [Create Messages](/docs/api/v1/messages/create-a-message) endpoint to send a message directly to another person by using the `toPersonId` or `toPersonEmail` parameters.\n\nBots are not able to create and simultaneously classify a room. A bot may update a space classification after a person of the same owning organization joined the space as the first human user.\nA space can only be put into announcement mode when it is locked.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/rooms")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("title", title)
					req.BodyString("teamId", teamId)
					req.BodyString("classificationId", classificationId)
					req.BodyBool("isLocked", isLocked, cmd.Flags().Changed("is-locked"))
					req.BodyBool("isPublic", isPublic, cmd.Flags().Changed("is-public"))
					req.BodyString("description", description)
					req.BodyBool("isAnnouncementOnly", isAnnouncementOnly, cmd.Flags().Changed("is-announcement-only"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&title, "title", "", "")
		cmd.Flags().StringVar(&teamId, "team-id", "", "")
		cmd.Flags().StringVar(&classificationId, "classification-id", "", "")
		cmd.Flags().BoolVar(&isLocked, "is-locked", false, "")
		cmd.Flags().BoolVar(&isPublic, "is-public", false, "")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().BoolVar(&isAnnouncementOnly, "is-announcement-only", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		roomsCmd.AddCommand(cmd)
	}

	{ // get
		var roomId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Room Details",
			Long:  "Shows details for a room, by ID.\n\nThe `title` of the room for 1:1 rooms will be the display name of the other person. When a Compliance Officer lists 1:1 rooms, the \"other\" person cannot be determined. This means that the room's title may not be filled in and instead shows \"Empty Title\". Please use the [memberships API](https://developer.webex.com/docs/api/v1/memberships) to list the other person in the space.\n\nSpecify the room ID in the `roomId` parameter in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/rooms/{roomId}")
				req.PathParam("roomId", roomId)
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
		cmd.Flags().StringVar(&roomId, "room-id", "", "The unique identifier for the room.")
		cmd.MarkFlagRequired("room-id")
		roomsCmd.AddCommand(cmd)
	}

	{ // update
		var roomId string
		var title string
		var classificationId string
		var teamId string
		var isLocked bool
		var isPublic bool
		var description string
		var isAnnouncementOnly bool
		var isReadOnly bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Room",
			Long:  "Updates details for a room, by ID.\n\nSpecify the room ID in the `roomId` parameter in the URI.\nA space can only be put into announcement mode when it is locked.\nAny space participant or compliance officer can convert a space from public to private. Only a compliance officer can convert a space from private to public and only if the space is classified with the lowest category (usually `public`), and the space has a description.\nTo remove a `description` please use a space character ` ` by itself.\n\n<div><Callout type=\"info\">When using this method for moving a space under a team, ensure that all moderators in the space are also team members. If a moderator is not part of the team, demote or remove them as a moderator. Alternatively, add the non-team moderators to the team. This ensures compliance with the requirement that all space moderators must be team members for successful operation execution.\n</Callout></div>\n\n<div><Callout type=\"info\">A Compliance Officer who is not a member of a space can only update the `classificationId`, `isAnnouncementOnly`, `description`, and `isPublic` fields.\n</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/rooms/{roomId}")
				req.PathParam("roomId", roomId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("title", title)
					req.BodyString("classificationId", classificationId)
					req.BodyString("teamId", teamId)
					req.BodyBool("isLocked", isLocked, cmd.Flags().Changed("is-locked"))
					req.BodyBool("isPublic", isPublic, cmd.Flags().Changed("is-public"))
					req.BodyString("description", description)
					req.BodyBool("isAnnouncementOnly", isAnnouncementOnly, cmd.Flags().Changed("is-announcement-only"))
					req.BodyBool("isReadOnly", isReadOnly, cmd.Flags().Changed("is-read-only"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&roomId, "room-id", "", "The unique identifier for the room.")
		cmd.MarkFlagRequired("room-id")
		cmd.Flags().StringVar(&title, "title", "", "")
		cmd.Flags().StringVar(&classificationId, "classification-id", "", "")
		cmd.Flags().StringVar(&teamId, "team-id", "", "")
		cmd.Flags().BoolVar(&isLocked, "is-locked", false, "")
		cmd.Flags().BoolVar(&isPublic, "is-public", false, "")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().BoolVar(&isAnnouncementOnly, "is-announcement-only", false, "")
		cmd.Flags().BoolVar(&isReadOnly, "is-read-only", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		roomsCmd.AddCommand(cmd)
	}

	{ // delete
		var roomId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Room",
			Long:  "Deletes a room, by ID. Deleted rooms cannot be recovered.\nAs a security measure to prevent accidental deletion, when a non moderator deletes the room they are removed from the room instead.\n\nDeleting a room that is part of a team will archive the room instead.\n\nA Compliance Officer has no special privileges, i.e. they cannot delete rooms they are not part of.\n\nSpecify the room ID in the `roomId` parameter in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/rooms/{roomId}")
				req.PathParam("roomId", roomId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&roomId, "room-id", "", "The unique identifier for the room.")
		cmd.MarkFlagRequired("room-id")
		roomsCmd.AddCommand(cmd)
	}

	{ // get-meeting
		var roomId string
		cmd := &cobra.Command{
			Use:   "get-meeting",
			Short: "Get Room Meeting Details",
			Long:  "<div>\n<callout type=\"warning\">\nThe meetingInfo API is deprecated and will be EOL on Jan 31, 2025. Meetings in the WSMP must be scheduled and licensed via the meetings backend.\nThe [Create a Meeting](/docs/api/v1/meetings/create-a-meeting) endpoint will provide the SIP address for the meeting to call.\n</callout>\n</div>\n\nShows Webex meeting details for a room such as the SIP address, meeting URL, toll-free and toll dial-in numbers.\n\nSpecify the room ID in the `roomId` parameter in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/rooms/{roomId}/meetingInfo")
				req.PathParam("roomId", roomId)
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
		cmd.Flags().StringVar(&roomId, "room-id", "", "The unique identifier for the room.")
		cmd.MarkFlagRequired("room-id")
		roomsCmd.AddCommand(cmd)
	}

}
