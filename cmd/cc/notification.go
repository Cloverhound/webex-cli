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

var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Notification commands",
}

func init() {
	cmd.CcCmd.AddCommand(notificationCmd)

	{ // subscribe
		var isKeepAliveEnabled bool
		var clientType string
		var allowMultiLogin bool
		var force bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "subscribe",
			Short: "Subscribe Notification",
			Long:  `Access this endpoint when the user has to register for a WebSocket Session. Requires 'cjp:user' scope or roles 'id_full_admin', 'id_readonly_admin', 'atlas-portal.partner.salesadmin', 'cjp.supervisor', 'cjp.admin', 'atlas-portal.partner.provision_admin', 'cloud-contact-center:pod_conv' for authorization`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/v1/notification/subscribe")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("isKeepAliveEnabled", isKeepAliveEnabled, cmd.Flags().Changed("is-keep-alive-enabled"))
					req.BodyString("clientType", clientType)
					req.BodyBool("allowMultiLogin", allowMultiLogin, cmd.Flags().Changed("allow-multi-login"))
					req.BodyBool("force", force, cmd.Flags().Changed("force"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().BoolVar(&isKeepAliveEnabled, "is-keep-alive-enabled", false, "")
		cmd.Flags().StringVar(&clientType, "client-type", "", "")
		cmd.Flags().BoolVar(&allowMultiLogin, "allow-multi-login", false, "")
		cmd.Flags().BoolVar(&force, "force", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		notificationCmd.AddCommand(cmd)
	}

}
