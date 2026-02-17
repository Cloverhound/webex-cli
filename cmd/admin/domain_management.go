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

var domainManagementCmd = &cobra.Command{
	Use:   "domain-management",
	Short: "DomainManagement commands",
}

func init() {
	cmd.AdminCmd.AddCommand(domainManagementCmd)

	{ // get-verification-token
		var orgId string
		var domain string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "get-verification-token",
			Short: "Get Domain Verification Token",
			Long:  "This endpoint helps generate a token for a given domain within the specified organization. The user needs to add this token as a 'TXT' record to the DNS server.\n\n**Possible Error:**\n\n- 409: The request encountered a resource conflict. This error occurs if the domain is either claimed by another organization or by the same organization.\n\n**Authorization:**\n\nAn 'OAuth' token issued by the 'Identity Broker' is required to access this endpoint. The token must include one of the following scopes:\n\n- `Identity:Organization`\n\n- `identity:organizations_rw`\n\n**Administrator Roles:**\n\nThe following administrators can use this API:\n\n- `id_full_admin`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/identity/organizations/{orgId}/actions/getDomainVerificationToken")
				req.PathParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("domain", domain)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "The Webex Identity-assigned organization identifier for a user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&domain, "domain", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		domainManagementCmd.AddCommand(cmd)
	}

	{ // verify
		var orgId string
		var domain string
		var claimDomain bool
		var reserveDomain bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "verify",
			Short: "Verify Domain",
			Long:  "This endpoint helps verify a given domain within the specified organization. This API verifies domain ownership by looking up and validating the 'TXT' record for the domain.\nOnce verified, domain enforcement will be applied to the organization. Any users in the organization whose email domain doesn't match one of the verified domains will be marked as transient.\n\nIf you want to verify and claim the domain, just set the 'claimDomain' parameter to true. By default, it's set to false, which will only verify the domain.\n\n**Possible Errors:**\n\n- 400: The request was a Bad Request. The domain can't be verified. This error happens if the user didn't request a token before trying to verify the domain.\n\n- 409: The request resulted in a resource conflict. This error occurs if the domain has already been claimed by another organization.\n\n**Authorization:**\n\nAn 'OAuth' token issued by the 'Identity Broker' is required to access this endpoint. The token must include one of the following scopes:\n\n- `Identity:Organization`\n\n- `identity:organizations_rw`\n\n**Administrator Roles:**\n\nThe following administrators can use this API:\n\n- `id_full_admin`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/identity/organizations/{orgId}/actions/verifyDomain")
				req.PathParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("domain", domain)
					req.BodyBool("claimDomain", claimDomain, cmd.Flags().Changed("claim-domain"))
					req.BodyBool("reserveDomain", reserveDomain, cmd.Flags().Changed("reserve-domain"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "The Webex Identity-assigned organization identifier for a user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&domain, "domain", "", "")
		cmd.Flags().BoolVar(&claimDomain, "claim-domain", false, "")
		cmd.Flags().BoolVar(&reserveDomain, "reserve-domain", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		domainManagementCmd.AddCommand(cmd)
	}

	{ // claim
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "claim",
			Short: "Claim Domain",
			Long:  "This endpoint helps claim the given domain within the specified organization. The domain needs to be verified before it can be claimed.\n\n**Note**\n<callout type=\"warning\">\n\nThere's an organization-level boolean flag called 'enforceVerifiedDomains'. If this flag is set to false, we won't put any user in the organization into a transient state when verifying or claiming a domain.\nCustomers can still create users within the organization who don't use the verified domains as their email. However, if the flag is set to true, all users in the organization must use one of the verified domains as their email.\nThis flag defines whether the organization enforces user email verification within the organization. If set to true, all users inside the organization must use one of the verified domains.\nThis flag is effective only after the admin has verified at least one email domain.\n</callout>\n\n**Possible Error:**\n\n- 400: The request was a Bad Request. This error occurs if the domain is not verified.\n\n**Authorization:**\n\nAn 'OAuth' token issued by the 'Identity Broker' is required to access this endpoint. The token must include one of the following scopes:\n\n- `Identity:Organization`\n\n- `identity:organizations_rw`\n\n**Administrator Roles:**\n\nThe following administrators can use this API:\n\n- `id_full_admin`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/identity/organizations/{orgId}/actions/claimDomain")
				req.PathParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "The Webex Identity-assigned organization identifier for a user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		domainManagementCmd.AddCommand(cmd)
	}

	{ // unverify
		var orgId string
		var domain string
		var removePending bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "unverify",
			Short: "Unverify Domain",
			Long:  "After you unclaim the domain, it will still be verified. Domain enforcement will still apply to the organization. The unverify endpoint helps to remove the domain ownership verification for the organization.\n\n**Possible Error:**\n\n- 400: The request was a Bad Request. The domain cannot be unverified. This error occurs if the domain is still claimed.\n\n- 404: The request was Not Found. This error occurs if the domain is not associated with the organization.\n\n**Authorization:**\n\nAn 'OAuth' token issued by the 'Identity Broker' is required to access this endpoint. The token must include one of the following scopes:\n\n- `Identity:Organization`\n\n- `identity:organizations_rw`\n\n**Administrator Roles:**\n\nThe following administrators can use this API:\n\n- `id_full_admin`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/identity/organizations/{orgId}/actions/unverifyDomain")
				req.PathParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("domain", domain)
					req.BodyBool("removePending", removePending, cmd.Flags().Changed("remove-pending"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "The Webex Identity-assigned organization identifier for a user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&domain, "domain", "", "")
		cmd.Flags().BoolVar(&removePending, "remove-pending", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		domainManagementCmd.AddCommand(cmd)
	}

	{ // unclaim
		var orgId string
		var domain string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "unclaim",
			Short: "Unclaim Domain",
			Long:  "This API is used to unclaim a domain for the organization. The domain will remain verified, and domain enforcement will still apply to the given organization.\n\n**Possible Error:**\n\n- 400: The request was a Bad Request. The domain cannot be unclaimed. This error occurs if the requested parameter is invalid.\n\n**Authorization:**\n\nAn 'OAuth' token issued by the 'Identity Broker' is required to access this endpoint. The token must include one of the following scopes:\n\n- `Identity:Organization`\n\n- `identity:organizations_rw`\n\n**Administrator Roles:**\n\nThe following administrators can use this API:\n\n- `id_full_admin`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/identity/organizations/{orgId}/actions/unclaimDomain")
				req.PathParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("domain", domain)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "The Webex Identity-assigned organization identifier for a user's organization.")
		cmd.MarkFlagRequired("org-id")
		cmd.Flags().StringVar(&domain, "domain", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		domainManagementCmd.AddCommand(cmd)
	}

}
