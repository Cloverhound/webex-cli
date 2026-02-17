package calling

import (
	"fmt"
	"strconv"

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
var _ = strconv.Itoa

var callingServiceCmd = &cobra.Command{
	Use:   "calling-service",
	Short: "CallingService commands",
}

func init() {
	cmd.CallingCmd.AddCommand(callingServiceCmd)

	{ // list-announcement-languages
		cmd := &cobra.Command{
			Use:   "list-announcement-languages",
			Short: "Read the List of Announcement Languages",
			Long:  "List all languages supported by Webex Calling for announcements and voice prompts.\n\nRetrieving announcement languages requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/announcementLanguages")
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
		callingServiceCmd.AddCommand(cmd)
	}

	{ // get-voicemail
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-voicemail",
			Short: "Get Voicemail Settings",
			Long:  "Retrieve the organization's voicemail settings.\n\nOrganizational voicemail settings determines what voicemail features a person can configure and automatic message expiration.\n\nRetrieving organization's voicemail settings requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/voicemail/settings")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve voicemail settings for this organization.")
		callingServiceCmd.AddCommand(cmd)
	}

	{ // update-voicemail
		var orgId string
		var messageExpiryEnabled bool
		var numberOfDaysForMessageExpiry int64
		var strictDeletionEnabled bool
		var voiceMessageForwardingEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-voicemail",
			Short: "Update Voicemail Settings",
			Long:  "Update the organization's voicemail settings.\n\nOrganizational voicemail settings determines what voicemail features a person can configure and automatic message expiration.\n\nUpdating an organization's voicemail settings requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/voicemail/settings")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("messageExpiryEnabled", messageExpiryEnabled, cmd.Flags().Changed("message-expiry-enabled"))
					req.BodyInt("numberOfDaysForMessageExpiry", numberOfDaysForMessageExpiry, cmd.Flags().Changed("number-of-days-for-message-expiry"))
					req.BodyBool("strictDeletionEnabled", strictDeletionEnabled, cmd.Flags().Changed("strict-deletion-enabled"))
					req.BodyBool("voiceMessageForwardingEnabled", voiceMessageForwardingEnabled, cmd.Flags().Changed("voice-message-forwarding-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update voicemail settings for this organization.")
		cmd.Flags().BoolVar(&messageExpiryEnabled, "message-expiry-enabled", false, "")
		cmd.Flags().Int64Var(&numberOfDaysForMessageExpiry, "number-of-days-for-message-expiry", 0, "")
		cmd.Flags().BoolVar(&strictDeletionEnabled, "strict-deletion-enabled", false, "")
		cmd.Flags().BoolVar(&voiceMessageForwardingEnabled, "voice-message-forwarding-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callingServiceCmd.AddCommand(cmd)
	}

	{ // get-voicemail-rules
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-voicemail-rules",
			Short: "Get Voicemail Rules",
			Long:  "Retrieve the organization's voicemail rules.\n\nOrganizational voicemail rules specify the default passcode requirements. They are provided for informational purposes only and cannot be modified.\n\nRetrieving the organization's voicemail rules requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/voicemail/rules")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve voicemail rules for this organization.")
		callingServiceCmd.AddCommand(cmd)
	}

	{ // update-voicemail-rules
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-voicemail-rules",
			Short: "Update Voicemail Rules",
			Long:  "Update the organization's default voicemail passcode and/or rules.\n\nOrganizational voicemail rules specify the default passcode requirements.\n\nIf you choose to set a default passcode for new people added to your organization, communicate to your people what that passcode is, and that it must be reset before they can access their voicemail. If this feature is not turned on, each new person must initially set their own passcode.\n\nUpdating an organization's voicemail passcode and/or rules requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/voicemail/rules")
				req.QueryParam("orgId", orgId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update voicemail rules for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callingServiceCmd.AddCommand(cmd)
	}

	{ // get-org-music-hold-configuration
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-org-music-hold-configuration",
			Short: "Get the organization Music on Hold configuration",
			Long:  `Retrieve the organization's Music on Hold settings.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/moh/settings")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve Music on Hold settings for this organization.")
		callingServiceCmd.AddCommand(cmd)
	}

	{ // update-org-music-hold-configuration
		var orgId string
		var defaultOrgMoh string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-org-music-hold-configuration",
			Short: "Update the organization Music on Hold configuration",
			Long:  `Update the organization's Music on Hold settings.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/moh/settings")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("defaultOrgMoh", defaultOrgMoh)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Patch Music on Hold for this organization.")
		cmd.Flags().StringVar(&defaultOrgMoh, "default-org-moh", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callingServiceCmd.AddCommand(cmd)
	}

	{ // get-org-call-captions-settings
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-org-call-captions-settings",
			Short: "Get the organization call captions settings",
			Long:  "Retrieve the organization's call captions settings.\n\nThe call caption feature allows the customer to enable and manage closed captions and transcript functionality (rolling caption panel) in Webex Calling, without requiring the user to escalate the call to a meeting.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/callCaptions")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		callingServiceCmd.AddCommand(cmd)
	}

	{ // update-org-call-captions-settings
		var orgId string
		var orgClosedCaptionsEnabled bool
		var orgTranscriptsEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-org-call-captions-settings",
			Short: "Update the organization call captions settings",
			Long:  "Update the organization's call captions settings.\n\nThe call caption feature allows the customer to enable and manage closed captions and transcript functionality (rolling caption panel) in Webex Calling, without requiring the user to escalate the call to a meeting.\n\nThis API requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/callCaptions")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("orgClosedCaptionsEnabled", orgClosedCaptionsEnabled, cmd.Flags().Changed("org-closed-captions-enabled"))
					req.BodyBool("orgTranscriptsEnabled", orgTranscriptsEnabled, cmd.Flags().Changed("org-transcripts-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().BoolVar(&orgClosedCaptionsEnabled, "org-closed-captions-enabled", false, "")
		cmd.Flags().BoolVar(&orgTranscriptsEnabled, "org-transcripts-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callingServiceCmd.AddCommand(cmd)
	}

}
