package calling

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

var callerReputationProviderCmd = &cobra.Command{
	Use:   "caller-reputation-provider",
	Short: "CallerReputationProvider commands",
}

func init() {
	cmd.CallingCmd.AddCommand(callerReputationProviderCmd)

	{ // get-settings
		var organizationId string
		cmd := &cobra.Command{
			Use:   "get-settings",
			Short: "Get Caller Reputation Provider Service Settings",
			Long:  `Retrieves the configuration and status of the caller reputation provider service for Webex Calling.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/serviceSettings/callerReputationProvider")
				req.QueryParam("organizationId", organizationId)
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
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "Unique identifier for the organization.")
		callerReputationProviderCmd.AddCommand(cmd)
	}

	{ // update-settings
		var organizationId string
		var enabled bool
		var id string
		var name string
		var clientId string
		var clientSecret string
		var callBlockScoreThreshold string
		var callAllowScoreThreshold string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-settings",
			Short: "Update Caller Reputation Provider Service Settings",
			Long:  `Updates the configuration of the caller reputation provider service for Webex Calling.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/serviceSettings/callerReputationProvider")
				req.QueryParam("organizationId", organizationId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyString("id", id)
					req.BodyString("name", name)
					req.BodyString("clientId", clientId)
					req.BodyString("clientSecret", clientSecret)
					req.BodyString("callBlockScoreThreshold", callBlockScoreThreshold)
					req.BodyString("callAllowScoreThreshold", callAllowScoreThreshold)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "Unique identifier for the organization.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&id, "id", "", "")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&clientId, "client-id", "", "")
		cmd.Flags().StringVar(&clientSecret, "client-secret", "", "")
		cmd.Flags().StringVar(&callBlockScoreThreshold, "call-block-score-threshold", "", "")
		cmd.Flags().StringVar(&callAllowScoreThreshold, "call-allow-score-threshold", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callerReputationProviderCmd.AddCommand(cmd)
	}

	{ // get-status
		var organizationId string
		cmd := &cobra.Command{
			Use:   "get-status",
			Short: "Get Caller Reputation Provider Status",
			Long:  `Retrieves the current status of the caller reputation provider integration for Webex Calling.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/serviceSettings/callerReputationProvider/status")
				req.QueryParam("organizationId", organizationId)
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
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "Unique identifier for the organization.")
		callerReputationProviderCmd.AddCommand(cmd)
	}

	{ // unlock
		var organizationId string
		var id string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "unlock",
			Short: "Unlock Caller Reputation Provider",
			Long:  `Unlocks the caller reputation provider for Webex Calling.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/serviceSettings/callerReputationProvider/actions/unlock/invoke")
				req.QueryParam("organizationId", organizationId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("id", id)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "Unique identifier for the organization.")
		cmd.Flags().StringVar(&id, "id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callerReputationProviderCmd.AddCommand(cmd)
	}

	{ // get
		var organizationId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Caller Reputation Provider Providers",
			Long:  `Retrieves the list of available caller reputation providers and their regions for Webex Calling.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/serviceSettings/callerReputationProvider/providers")
				req.QueryParam("organizationId", organizationId)
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
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "Unique identifier for the organization.")
		callerReputationProviderCmd.AddCommand(cmd)
	}

}
