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

var sendActivationEmailCmd = &cobra.Command{
	Use:   "send-activation-email",
	Short: "SendActivationEmail commands",
}

func init() {
	cmd.AdminCmd.AddCommand(sendActivationEmailCmd)

	{ // initiate-bulk-resend-job
		var orgId string
		cmd := &cobra.Command{
			Use:   "initiate-bulk-resend-job",
			Short: "Initiate Bulk Activation Email Resend Job",
			Long:  "Initiate a bulk activation email resend job that sends an activation email to all eligible users in an organization. Only a single instance of the job can be running for an organization.\n\nRequires a full or user administrator auth token with a scope of `spark-admin:people_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/identity/organizations/{orgId}/jobs/sendActivationEmails")
				req.PathParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Initiate job for this organization.")
		cmd.MarkFlagRequired("org-id")
		sendActivationEmailCmd.AddCommand(cmd)
	}

	{ // get-bulk-resend-job-status
		var orgId string
		var jobId string
		cmd := &cobra.Command{
			Use:   "get-bulk-resend-job-status",
			Short: "Get Bulk Activation Email Resend Job Status",
			Long:  "Get the details of an activation email resend job by its job ID.\n\nRequires a full or user administrator auth token with a scope of `spark-admin:people_write` or read-only administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/identity/organizations/{orgId}/jobs/sendActivationEmails/{jobId}/status")
				req.PathParam("orgId", orgId)
				req.PathParam("jobId", jobId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Check job status for this organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&jobId, "job-id", "", "Retrieve job status for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		sendActivationEmailCmd.AddCommand(cmd)
	}

	{ // get-bulk-resend-job-errors
		var orgId string
		var jobId string
		var max string
		cmd := &cobra.Command{
			Use:   "get-bulk-resend-job-errors",
			Short: "Get Bulk Activation Email Resend Job Errors",
			Long:  "Get errors of an activation email resend job by its job ID.\n\nRequires a full or user administrator auth token with a scope of `spark-admin:people_write` or read-only administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/identity/organizations/{orgId}/jobs/sendActivationEmails/{jobId}/errors")
				req.PathParam("orgId", orgId)
				req.PathParam("jobId", jobId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Check job status for this organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&jobId, "job-id", "", "Retrieve job status for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		cmd.Flags().StringVar(&max, "max", "", "Limit the maximum number of errors in the response.")
		sendActivationEmailCmd.AddCommand(cmd)
	}

}
