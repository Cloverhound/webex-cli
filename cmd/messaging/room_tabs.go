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

var roomTabsCmd = &cobra.Command{
	Use:   "room-tabs",
	Short: "RoomTabs commands",
}

func init() {
	cmd.MessagingCmd.AddCommand(roomTabsCmd)

	{ // list
		var roomId string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Room Tabs",
			Long:  "Lists all Room Tabs of a room specified by the `roomId` query parameter.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/room/tabs")
				req.QueryParam("roomId", roomId)
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
		cmd.Flags().StringVar(&roomId, "room-id", "", "ID of the room for which to list room tabs.")
		roomTabsCmd.AddCommand(cmd)
	}

	{ // create
		var roomId string
		var contentUrl string
		var displayName string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Room Tab",
			Long:  `Add a tab with a specified URL to a room.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/room/tabs")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("roomId", roomId)
					req.BodyString("contentUrl", contentUrl)
					req.BodyString("displayName", displayName)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&roomId, "room-id", "", "")
		cmd.Flags().StringVar(&contentUrl, "content-url", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		roomTabsCmd.AddCommand(cmd)
	}

	{ // get
		var id string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Room Tab Details",
			Long:  `Get details for a Room Tab with the specified room tab ID.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/room/tabs/{id}")
				req.PathParam("id", id)
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
		cmd.Flags().StringVar(&id, "id", "", "The unique identifier for the Room Tab.")
		cmd.MarkFlagRequired("id")
		roomTabsCmd.AddCommand(cmd)
	}

	{ // update
		var id string
		var roomId string
		var contentUrl string
		var displayName string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Room Tab",
			Long:  `Updates the content URL of the specified Room Tab ID.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/room/tabs/{id}")
				req.PathParam("id", id)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("roomId", roomId)
					req.BodyString("contentUrl", contentUrl)
					req.BodyString("displayName", displayName)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&id, "id", "", "The unique identifier for the Room Tab.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&roomId, "room-id", "", "")
		cmd.Flags().StringVar(&contentUrl, "content-url", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		roomTabsCmd.AddCommand(cmd)
	}

	{ // delete
		var id string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Room Tab",
			Long:  `Deletes a Room Tab with the specified ID.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/room/tabs/{id}")
				req.PathParam("id", id)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&id, "id", "", "The unique identifier for the Room Tab to delete.")
		cmd.MarkFlagRequired("id")
		roomTabsCmd.AddCommand(cmd)
	}

}
