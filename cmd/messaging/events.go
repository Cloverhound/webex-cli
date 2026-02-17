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

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Events commands",
}

func init() {
	cmd.MessagingCmd.AddCommand(eventsCmd)

	{ // list
		var resource string
		var typeVal string
		var actorId string
		var from string
		var to string
		var max string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Events",
			Long: `List events in your organization. Several query parameters are available to filter the events returned in the response.

Long result sets will be split into [pages](/docs/basics#pagination).`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/events")
				req.QueryParam("resource", resource)
				req.QueryParam("type", typeVal)
				req.QueryParam("actorId", actorId)
				req.QueryParam("from", from)
				req.QueryParam("to", to)
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
		cmd.Flags().StringVar(&resource, "resource", "", "List events with a specific resource type.")
		cmd.Flags().StringVar(&typeVal, "type", "", "List events with a specific event type.")
		cmd.Flags().StringVar(&actorId, "actor-id", "", "List events performed by this person, by person ID.")
		cmd.Flags().StringVar(&from, "from", "", "List events which occurred after a specific date and time.")
		cmd.Flags().StringVar(&to, "to", "", "List events that occurred before a specific date and time. If not specified, events up to the present time will be listed. Cannot be set to a future date relative to the current time.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of events in the response. Value must be between 1 and 1000, inclusive.")
		eventsCmd.AddCommand(cmd)
	}

	{ // get
		var eventId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Event Details",
			Long:  "Shows details for an event, by event ID.\n\nSpecify the event ID in the `eventId` parameter in the URI.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/events/{eventId}")
				req.PathParam("eventId", eventId)
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
		cmd.Flags().StringVar(&eventId, "event-id", "", "The unique identifier for the event.")
		cmd.MarkFlagRequired("event-id")
		eventsCmd.AddCommand(cmd)
	}

}
