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

var guestManagementCmd = &cobra.Command{
	Use:   "guest-management",
	Short: "GuestManagement commands",
}

func init() {
	cmd.AdminCmd.AddCommand(guestManagementCmd)

	{ // create
		var subject string
		var displayName string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Guest",
			Long:  "Create a new token for a single guest user. The Service App that creates the guest must have the scope `guest-issuer:write`.\n\nGuests are implicitly created by retrieving the guest access token.\n\nRepeated calls to this API with the same `subject` will create additional tokens without invalidating previous ones. Tokens are valid until the `expiresIn`.\n\nGuests can be renamed by supplying the same `subject` and changing the `displayName.`\n\nTo retrieve a new token for an existing guest, please provide the existing guest's `subject`. Tokens are valid until `expiresIn`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/guests/token")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("subject", subject)
					req.BodyString("displayName", displayName)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&subject, "subject", "", "")
		cmd.Flags().StringVar(&displayName, "display-name", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		guestManagementCmd.AddCommand(cmd)
	}

	{ // get-count
		cmd := &cobra.Command{
			Use:   "get-count",
			Short: "Get Guest Count",
			Long:  "To retrieve the number of guests, the scopes `guest-issuer:read` or `guest-issuer:write` are needed.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/guests/count")
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
		guestManagementCmd.AddCommand(cmd)
	}

}
