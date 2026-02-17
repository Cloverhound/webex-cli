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

var sessionTypesCmd = &cobra.Command{
	Use:   "session-types",
	Short: "SessionTypes commands",
}

func init() {
	cmd.MeetingsCmd.AddCommand(sessionTypesCmd)

	{ // list-site
		var siteUrl string
		cmd := &cobra.Command{
			Use:   "list-site",
			Short: "List Site Session Types",
			Long:  `List session types for a specific site.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/admin/meeting/config/sessionTypes")
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
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site to query. If siteUrl is not specified, the query will use the default site for the admin's authorization token used to make the call.")
		sessionTypesCmd.AddCommand(cmd)
	}

	{ // list-user
		var siteUrl string
		var personId string
		var email string
		cmd := &cobra.Command{
			Use:   "list-user",
			Short: "List User Session Type",
			Long:  `List session types for a specific user.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/admin/meeting/userconfig/sessionTypes")
				req.QueryParam("siteUrl", siteUrl)
				req.QueryParam("personId", personId)
				req.Header("email", email)
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
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site to query.")
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the user.")
		cmd.Flags().StringVar(&email, "email", "", "e.g. `john.andersen@example.com` (string, optional) - The email of the user.")
		sessionTypesCmd.AddCommand(cmd)
	}

	{ // update-user
		var siteUrl string
		var sessionTypeIds []string
		var personId string
		var email string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-user",
			Short: "Update User Session Types",
			Long:  "Assign session types to specific users.\n\n* At least one of the following body parameters is required to update a specific user session type: `personId`, `email`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/admin/meeting/userconfig/sessionTypes")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("siteUrl", siteUrl)
					req.BodyStringSlice("sessionTypeIds", sessionTypeIds)
					req.BodyString("personId", personId)
					req.BodyString("email", email)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "")
		cmd.Flags().StringSliceVar(&sessionTypeIds, "session-type-ids", nil, "")
		cmd.Flags().StringVar(&personId, "person-id", "", "")
		cmd.Flags().StringVar(&email, "email", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		sessionTypesCmd.AddCommand(cmd)
	}

}
