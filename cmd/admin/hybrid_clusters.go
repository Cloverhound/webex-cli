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

var hybridClustersCmd = &cobra.Command{
	Use:   "hybrid-clusters",
	Short: "HybridClusters commands",
}

func init() {
	cmd.AdminCmd.AddCommand(hybridClustersCmd)

	{ // list
		var orgId string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Hybrid Clusters",
			Long:  "List hybrid clusters for an organization. If no `orgId` is specified, the default is the organization of the authenticated user.\n\nOnly an admin auth token with the `spark-admin:hybrid_clusters_read` scope can list clusters.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/hybrid/clusters")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List hybrid clusters in this organization. If an organization is not specified, the organization of the caller will be used.")
		hybridClustersCmd.AddCommand(cmd)
	}

	{ // get
		var hybridClusterId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Hybrid Cluster Details",
			Long:  "Shows details for a hybrid cluster, by ID.\n\nOnly an admin auth token with the `spark-admin:hybrid_clusters_read` scope can see cluster details.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/hybrid/clusters/{hybridClusterId}")
				req.PathParam("hybridClusterId", hybridClusterId)
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
		cmd.Flags().StringVar(&hybridClusterId, "hybrid-cluster-id", "", "The ID of the cluster.")
		cmd.MarkFlagRequired("hybrid-cluster-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", " Find the cluster in this specific organization. If this is not specified, the organization of the caller will be used. ")
		hybridClustersCmd.AddCommand(cmd)
	}

}
