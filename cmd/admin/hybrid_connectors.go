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

var hybridConnectorsCmd = &cobra.Command{
	Use:   "hybrid-connectors",
	Short: "HybridConnectors commands",
}

func init() {
	cmd.AdminCmd.AddCommand(hybridConnectorsCmd)

	{ // list
		var orgId string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Hybrid Connectors",
			Long:  "List hybrid connectors for an organization. If no `orgId` is specified, the default is the organization of the authenticated user.\n\nOnly an admin auth token with the `spark-admin:hybrid_connectors_read` scope can list connectors.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/hybrid/connectors")
				req.QueryParam("orgId", orgId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List hybrid connectors in this organization. If an organization is not specified, the organization of the caller will be used.")
		hybridConnectorsCmd.AddCommand(cmd)
	}

	{ // get
		var connectorId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Hybrid Connector Details",
			Long:  "Shows details for a hybrid connector, by ID.\n\nOnly an admin auth token with the `spark-admin:hybrid_connectors_read` scope can see connector details.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/hybrid/connectors/{connectorId}")
				req.PathParam("connectorId", connectorId)
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
		cmd.Flags().StringVar(&connectorId, "connector-id", "", "The ID of the connector.")
		cmd.MarkFlagRequired("connector-id")
		hybridConnectorsCmd.AddCommand(cmd)
	}

}
