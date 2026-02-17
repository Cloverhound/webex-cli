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

var agentSummariesCmd = &cobra.Command{
	Use:   "agent-summaries",
	Short: "AgentSummaries commands",
}

func init() {
	cmd.CcCmd.AddCommand(agentSummariesCmd)

	{ // list
		var orgId string
		var searchType string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List summaries",
			Long:  `Lists summaries based on the requested search type.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/summary/list")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("orgId", orgId)
					req.BodyString("searchType", searchType)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "")
		cmd.Flags().StringVar(&searchType, "search-type", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		agentSummariesCmd.AddCommand(cmd)
	}

}
