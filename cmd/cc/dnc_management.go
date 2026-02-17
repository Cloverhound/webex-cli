package cc

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

var dncManagementCmd = &cobra.Command{
	Use:   "dnc-management",
	Short: "DncManagement commands",
}

func init() {
	cmd.CcCmd.AddCommand(dncManagementCmd)

	{ // create-phone-number-list
		var dncListName string
		var phoneNumber string
		var source string
		var reason string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-phone-number-list",
			Short: "Add Phone Number to DNC List",
			Long: `Adds a phone number to the specified Do Not Contact (DNC) list. This operation helps ensure compliance with applicable rules & regulations by preventing outbound calls to specific numbers.

**Note:** Phone numbers must be in E.164 format (e.g., +1234567890). Duplicate entries will be rejected with a 409 error.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/v3/campaign-management/dncList/{dncListName}/phoneNumber")
				req.PathParam("dncListName", dncListName)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("phoneNumber", phoneNumber)
					req.BodyString("source", source)
					req.BodyString("reason", reason)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&dncListName, "dnc-list-name", "", "This is the Name of the DNC list to which you want to add a phone number. List names are case-sensitive and must be URL-encoded if they contain special characters.")
		cmd.MarkFlagRequired("dnc-list-name")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "")
		cmd.Flags().StringVar(&source, "source", "", "")
		cmd.Flags().StringVar(&reason, "reason", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		dncManagementCmd.AddCommand(cmd)
	}

	{ // get-phone-number-list
		var dncListName string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-phone-number-list",
			Short: "Get Phone Number from DNC List",
			Long: `Retrieves details of a specific phone number from the specified Do Not Contact (DNC) list. This operation allows you to check if a number is present in the list and view associated metadata.

**Note:** Phone numbers must be URL-encoded and in E.164 format (e.g., %2B1234567890 for +1234567890).`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/v3/campaign-management/dncList/{dncListName}/phoneNumber/{phoneNumber}")
				req.PathParam("dncListName", dncListName)
				req.PathParam("phoneNumber", phoneNumber)
				if config.Paginate() {
					resp, statusCode, err := req.DoPaginated(false)
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
		cmd.Flags().StringVar(&dncListName, "dnc-list-name", "", "The name of the DNC list to search in. List names are case-sensitive and must be URL-encoded if they contain special characters.")
		cmd.MarkFlagRequired("dnc-list-name")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "The phone number to retrieve from the DNC list. Must be URL-encoded and in E.164 format (e.g., %2B1234567890 for +1234567890).")
		cmd.MarkFlagRequired("phone-number")
		dncManagementCmd.AddCommand(cmd)
	}

	{ // delete-phone-number-list
		var dncListName string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "delete-phone-number-list",
			Short: "Remove Phone Number from DNC List",
			Long: `Removes a phone number from the specified Do Not Contact (DNC) list. This operation allows administrators to remove numbers that should no longer be blocked from receiving calls.

**Note:** Phone numbers must be URL-encoded and in E.164 format (e.g., +1234567890 for +1234567890). If the number doesn't exist in the list, a 404 error will be returned.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/v3/campaign-management/dncList/{dncListName}/phoneNumber/{phoneNumber}")
				req.PathParam("dncListName", dncListName)
				req.PathParam("phoneNumber", phoneNumber)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&dncListName, "dnc-list-name", "", "The name of the DNC list to remove the phone number from. List names are case-sensitive and must be URL-encoded if they contain special characters.")
		cmd.MarkFlagRequired("dnc-list-name")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "The phone number to remove from the DNC list. Must be URL-encoded and in E.164 format (e.g., %2B1234567890 for +1234567890).")
		cmd.MarkFlagRequired("phone-number")
		dncManagementCmd.AddCommand(cmd)
	}

}
