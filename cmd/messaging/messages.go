package messaging

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

var messagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Messages commands",
}

func init() {
	cmd.MessagingCmd.AddCommand(messagesCmd)

	{ // list
		var roomId string
		var parentId string
		var mentionedPeople string
		var before string
		var beforeMessage string
		var max string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Messages",
			Long: `Lists all messages in a room.  Each message will include content attachments if present.

The list sorts the messages in descending order by creation date.

Long result sets will be split into [pages](/docs/basics#pagination).`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/messages")
				req.QueryParam("roomId", roomId)
				req.QueryParam("parentId", parentId)
				req.QueryParam("mentionedPeople", mentionedPeople)
				req.QueryParam("mentionedPeople", mentionedPeople)
				req.QueryParam("before", before)
				req.QueryParam("beforeMessage", beforeMessage)
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
		cmd.Flags().StringVar(&roomId, "room-id", "", "List messages in a room, by ID.")
		cmd.Flags().StringVar(&parentId, "parent-id", "", "List messages with a parent, by ID.")
		cmd.Flags().StringVar(&mentionedPeople, "mentioned-people", "", "List messages with these people mentioned, by ID. Use `me` as a shorthand for the current API user. Only `me` or the person ID of the current user may be specified. Bots must include this parameter to list messages in group rooms (spaces).")
		cmd.Flags().StringVar(&before, "before", "", "List messages sent before a date and time.")
		cmd.Flags().StringVar(&beforeMessage, "before-message", "", "List messages sent before a message, by ID.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of messages in the response. Cannot exceed 100 if used with `mentionedPeople`.")
		messagesCmd.AddCommand(cmd)
	}

	{ // create
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Message",
			Long:  "Post a plain text or [rich text](/docs/basics#formatting-messages) message, and optionally, a [file attachment](/docs/basics#message-attachments) attachment, to a room.\n\nThe `files` parameter is an array, which accepts multiple values to allow for future expansion, but currently only one file may be included with the message. File previews are only rendered for attachments of 1MB or less.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/messages")
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
		messagesCmd.AddCommand(cmd)
	}

	{ // list-direct
		var parentId string
		var personId string
		var personEmail string
		cmd := &cobra.Command{
			Use:   "list-direct",
			Short: "List Direct Messages",
			Long:  "List all messages in a 1:1 (direct) room. Use the `personId` or `personEmail` query parameter to specify the room. Each message will include content attachments if present.\n\nThe list sorts the messages in descending order by creation date.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/messages/direct")
				req.QueryParam("parentId", parentId)
				req.QueryParam("personId", personId)
				req.QueryParam("personEmail", personEmail)
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
		cmd.Flags().StringVar(&parentId, "parent-id", "", "List messages with a parent, by ID.")
		cmd.Flags().StringVar(&personId, "person-id", "", "List messages in a 1:1 room, by person ID.")
		cmd.Flags().StringVar(&personEmail, "person-email", "", "List messages in a 1:1 room, by person email.")
		messagesCmd.AddCommand(cmd)
	}

	{ // edit
		var messageId string
		var roomId string
		var text string
		var markdown string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "edit",
			Short: "Edit a Message",
			Long:  "Update a message you have posted not more than 10 times.\n\nSpecify the `messageId` of the message you want to edit.\n\nEdits of messages containing files or attachments are not currently supported.\nIf a user attempts to edit a message containing files or attachments a `400 Bad Request` will be returned by the API with a message stating that the feature is currently unsupported.\n\nThere is also a maximum number of times a user can edit a message. The maximum currently supported is 10 edits per message.\n    If a user attempts to edit a message greater that the maximum times allowed the API will return 400 Bad Request with a message stating the edit limit has been reached.\n\nWhile only the `roomId` and `text` or `markdown` attributes are *required* in the request body, a common pattern for editing message is to first call `GET /messages/{id}` for the message you wish to edit and to then update the `text` or `markdown` attribute accordingly, passing the updated message object in the request body of the `PUT /messages/{id}` request.\nWhen this pattern is used on a message that included markdown, the `html` attribute must be deleted prior to making the `PUT` request.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/messages/{messageId}")
				req.PathParam("messageId", messageId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("roomId", roomId)
					req.BodyString("text", text)
					req.BodyString("markdown", markdown)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&messageId, "message-id", "", "The unique identifier for the message.")
		cmd.MarkFlagRequired("message-id")
		cmd.Flags().StringVar(&roomId, "room-id", "", "")
		cmd.Flags().StringVar(&text, "text", "", "")
		cmd.Flags().StringVar(&markdown, "markdown", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		messagesCmd.AddCommand(cmd)
	}

	{ // get
		var messageId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Message Details",
			Long:  "Show details for a message, by message ID.\n\nSpecify the message ID in the `messageId` parameter in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/messages/{messageId}")
				req.PathParam("messageId", messageId)
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
		cmd.Flags().StringVar(&messageId, "message-id", "", "The unique identifier for the message.")
		cmd.MarkFlagRequired("message-id")
		messagesCmd.AddCommand(cmd)
	}

	{ // delete
		var messageId string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Message",
			Long:  "Delete a message, by message ID.\n\nSpecify the message ID in the `messageId` parameter in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/messages/{messageId}")
				req.PathParam("messageId", messageId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&messageId, "message-id", "", "The unique identifier for the message.")
		cmd.MarkFlagRequired("message-id")
		messagesCmd.AddCommand(cmd)
	}

}
