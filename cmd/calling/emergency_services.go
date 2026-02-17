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

var emergencyServicesCmd = &cobra.Command{
	Use:   "emergency-services",
	Short: "EmergencyServices commands",
}

func init() {
	cmd.CallingCmd.AddCommand(emergencyServicesCmd)

	{ // update-redsky-settings
		var orgId string
		var enabled bool
		var companyId string
		var secret string
		var externalTenantEnabled bool
		var email string
		var password string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-redsky-settings",
			Short: "Update RedSky Service Settings",
			Long:  "Update the RedSky service settings.\n\nThe Enhanced Emergency (E911) Service for Webex Calling provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada. E911 services are provided in conjunction with a RedSky account.\n\nUpdating the RedSky service settings requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/redSky/serviceSettings")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyString("companyId", companyId)
					req.BodyString("secret", secret)
					req.BodyBool("externalTenantEnabled", externalTenantEnabled, cmd.Flags().Changed("external-tenant-enabled"))
					req.BodyString("email", email)
					req.BodyString("password", password)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update E911 settings for the organization.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&companyId, "company-id", "", "")
		cmd.Flags().StringVar(&secret, "secret", "", "")
		cmd.Flags().BoolVar(&externalTenantEnabled, "external-tenant-enabled", false, "")
		cmd.Flags().StringVar(&email, "email", "", "")
		cmd.Flags().StringVar(&password, "password", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // create-account-admin-redsky
		var orgId string
		var email string
		var orgPrefix string
		var partnerRedskyOrgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-account-admin-redsky",
			Short: "Create an Account and Admin in RedSky",
			Long:  "Create an account and admin in RedSky.\n\nThe Enhanced Emergency (E911) Service for Webex Calling provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada. E911 services are provided in conjunction with a RedSky account.\n\nCreating a RedSky account requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/redSky")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("email", email)
					req.BodyString("orgPrefix", orgPrefix)
					req.BodyString("partnerRedskyOrgId", partnerRedskyOrgId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create RedSky account for the organization.")
		cmd.Flags().StringVar(&email, "email", "", "")
		cmd.Flags().StringVar(&orgPrefix, "org-prefix", "", "")
		cmd.Flags().StringVar(&partnerRedskyOrgId, "partner-redsky-org-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-redsky-account
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-redsky-account",
			Short: "Retrieve RedSky Account Details for an Organization",
			Long:  "Retrieve RedSky account details for an organization.\n\nThe Enhanced Emergency (E911) Service for Webex Calling provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada. E911 services are provided in conjunction with a RedSky account.\n\nTo retrieve the RedSky account details requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/redSky")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve RedSky account for the organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // update-org-redsky-account-compliance-status
		var orgId string
		var complianceStatus string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-org-redsky-account-compliance-status",
			Short: "Update the Organization RedSky Account's Compliance Status",
			Long:  "Update the compliance status for the customer's RedSky account.\n\nThe Enhanced Emergency (E911) Service for Webex Calling provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada. E911 services are provided in conjunction with a RedSky account.\n\nUpdating the RedSky account's compliance status requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/redSky/status")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("complianceStatus", complianceStatus)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update E911 compliance status for the organization.")
		cmd.Flags().StringVar(&complianceStatus, "compliance-status", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-org-compliance-status-redsky-account
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-org-compliance-status-redsky-account",
			Short: "Get the Organization Compliance Status for a RedSky Account",
			Long:  "Get the organization compliance status for a RedSky account. The `locationStatus.state` in the response will show the state for the location that is in the earliest stage of configuration.\n\nThe enhanced emergency (E911) service for Webex Calling provides an emergency service designed for organizations with a hybrid or nomadic workforce. It provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada.\n\nTo retrieve organization compliance status requires a full, user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/redSky/status")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve the compliance status for the organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-org-compliance-status-location-status-list
		var orgId string
		var start string
		var max string
		var order string
		cmd := &cobra.Command{
			Use:   "get-org-compliance-status-location-status-list",
			Short: "Get the Organization Compliance Status and the Location Status List",
			Long:  "Get the organization compliance status and the location status list for a RedSky account.\n\nThe enhanced emergency (E911) service for Webex Calling provides an emergency service designed for organizations with a hybrid or nomadic workforce. It provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada.\n\nTo retrieve organization compliance status requires a full, user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/redSky/complianceStatus")
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
				req.QueryParam("max", max)
				req.QueryParam("order", order)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve the compliance status and the list of location statuses for the organization.")
		cmd.Flags().StringVar(&start, "start", "", "Specifies the offset from the first result that you want to fetch.")
		cmd.Flags().StringVar(&max, "max", "", "Specifies the maximum number of records that you want to fetch.")
		cmd.Flags().StringVar(&order, "order", "", "Sort the list of locations in ascending or descending order. To sort in descending order append `-desc` to possible sort order values. Possible sort order values are `locationName` and `locationState`.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // login-redsky-admin-account
		var orgId string
		var email string
		var password string
		var redSkyOrgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "login-redsky-admin-account",
			Short: "Login to a RedSky Admin Account",
			Long:  "Login to Redsky for an existing account admin user to retrieve the `companyId` and verify the status of `externalTenantEnabled`. The password provided will not be stored.\n\nThe enhanced emergency (E911) service for Webex Calling provides an emergency service designed for organizations with a hybrid or nomadic workforce. It provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada.\n\nLogging in requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/redSky/actions/login/invoke")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("email", email)
					req.BodyString("password", password)
					req.BodyString("redSkyOrgId", redSkyOrgId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Login to a RedSky account for the organization.")
		cmd.Flags().StringVar(&email, "email", "", "")
		cmd.Flags().StringVar(&password, "password", "", "")
		cmd.Flags().StringVar(&redSkyOrgId, "red-sky-org-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-location-redsky-calling-parameters
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-location-redsky-calling-parameters",
			Short: "Get a Location's RedSky Emergency Calling Parameters",
			Long:  "Get the Emergency Calling Parameters for a specific location.\n\nThe enhanced emergency (E911) service for Webex Calling provides an emergency service designed for organizations with a hybrid or nomadic workforce. It provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada.\n\nTo retrieve location calling parameters requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/redSky")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve Calling Parameters for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve Calling Parameters for the location in this organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-location-redsky-compliance-status
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-location-redsky-compliance-status",
			Short: "Get a Location's RedSky Compliance Status",
			Long:  "Get RedSky compliance status for a specific location.\n\nThe enhanced emergency (E911) service for Webex Calling provides an emergency service designed for organizations with a hybrid or nomadic workforce. It provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada.\n\nRetrieving the location's compliance status requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/redSky/status")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve the compliance status for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve compliance status for the location in this organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // update-location-redsky-compliance-status
		var locationId string
		var orgId string
		var complianceStatus string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-location-redsky-compliance-status",
			Short: "Update a Location's RedSky Compliance Status",
			Long:  "Update the compliance status for a specific location.\n\nThe Enhanced Emergency (E911) Service for Webex Calling provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada. E911 services are provided in conjunction with a RedSky account.\n\nUpdating the RedSky account's compliance status requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/redSky/status")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("complianceStatus", complianceStatus)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update the E911 compliance status for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update the E911 compliance status for the location in this organization.")
		cmd.Flags().StringVar(&complianceStatus, "compliance-status", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // create-redsky-building-address-location
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-redsky-building-address-location",
			Short: "Create a RedSky Building Address and Alert Email for a Location",
			Long:  "Add a RedSky building address and alert email for a specified location.\n\nThe Enhanced Emergency (E911) Service for Webex Calling provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada. E911 services are provided in conjunction with a RedSky account.\n\nCreating a building address and alert email requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/redSky/building")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Create the building address and alert email for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization in which the location exists.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // update-redsky-building-address-location
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-redsky-building-address-location",
			Short: "Update a RedSky Building Address for a Location",
			Long:  "Update a RedSky building address for a specified location.\n\nThe Enhanced Emergency (E911) Service for Webex Calling provides dynamic location support and a network that routes emergency calls to Public Safety Answering Points (PSAP) around the US, its territories, and Canada. E911 services are provided in conjunction with a RedSky account.\n\nUpdating a building address requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/redSky/building")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update the building address for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "The organization in which the location exists.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-org-call-notification
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-org-call-notification",
			Short: "Get an Organization Emergency Call Notification",
			Long:  "Get organization emergency call notification.\n\nEmergency Call Notifications can be enabled at the organization level, allowing specified email addresses to receive email notifications when an emergency call is made. To comply with U.S. Public Law 115-127, also known as Kari\u2019s Law, any call that's made from within your organization to emergency services must generate an email notification.\n\nTo retrieve organization call notifications requires a full, user or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/emergencyCallNotification")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve Emergency Call Notification attributes for the organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // update-org-call-notification
		var orgId string
		var emergencyCallNotificationEnabled bool
		var allowEmailNotificationAllLocationEnabled bool
		var emailAddress string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-org-call-notification",
			Short: "Update an Organization Emergency Call Notification",
			Long:  "Update an organization emergency call notification.\n\nOnce settings are enabled at the organization level, the configured email address will receive emergency call notifications for all locations.\n\nEmergency Call Notifications can be enabled at the organization level, allowing specified email addresses to receive email notifications when an emergency call is made. To comply with U.S. Public Law 115-127, also known as Kari\u2019s Law, any call that's made from within your organization to emergency services must generate an email notification.\n\nTo update organization call notification requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/emergencyCallNotification")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("emergencyCallNotificationEnabled", emergencyCallNotificationEnabled, cmd.Flags().Changed("emergency-call-notification-enabled"))
					req.BodyBool("allowEmailNotificationAllLocationEnabled", allowEmailNotificationAllLocationEnabled, cmd.Flags().Changed("allow-email-notification-all-location-enabled"))
					req.BodyString("emailAddress", emailAddress)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update Emergency Call Notification attributes for the organization.")
		cmd.Flags().BoolVar(&emergencyCallNotificationEnabled, "emergency-call-notification-enabled", false, "")
		cmd.Flags().BoolVar(&allowEmailNotificationAllLocationEnabled, "allow-email-notification-all-location-enabled", false, "")
		cmd.Flags().StringVar(&emailAddress, "email-address", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-location-call-notification
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-location-call-notification",
			Short: "Get a Location Emergency Call Notification",
			Long:  "Get location emergency call notification.\n\nEmergency Call Notifications can be enabled at the organization level, allowing specified email addresses to receive email notifications when an emergency call is made. Once activated at the organization level, individual locations can configure this setting to direct notifications to specific email addresses. To comply with U.S. Public Law 115-127, also known as Kari\u2019s Law, any call that's made from within your organization to emergency services must generate an email notification.\n\nTo retrieve location call notifications requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/emergencyCallNotification")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve Emergency Call Notification attributes for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve Emergency Call Notification attributes for the location in this organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // update-location-call-notification
		var locationId string
		var orgId string
		var emergencyCallNotificationEnabled bool
		var emailAddress string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-location-call-notification",
			Short: "Update a Location Emergency Call Notification",
			Long:  "Update a location emergency call notification.\n\nOnce settings enabled at the organization level, the configured email address will receive emergency call notifications for all locations; for specific location customization, users can navigate to Management > Locations, select the Calling tab, and update the Emergency Call Notification settings.\n\nEmergency Call Notifications can be enabled at the organization level, allowing specified email addresses to receive email notifications when an emergency call is made. Once activated at the organization level, individual locations can configure this setting to direct notifications to specific email addresses. To comply with U.S. Public Law 115-127, also known as Kari\u2019s Law, any call that's made from within your organization to emergency services must generate an email notification.\n\nTo update location call notification requires a full, user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/emergencyCallNotification")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("emergencyCallNotificationEnabled", emergencyCallNotificationEnabled, cmd.Flags().Changed("emergency-call-notification-enabled"))
					req.BodyString("emailAddress", emailAddress)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update Emergency Call Notification attributes for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update Emergency Call Notification attributes for a location in this organization.")
		cmd.Flags().BoolVar(&emergencyCallNotificationEnabled, "emergency-call-notification-enabled", false, "")
		cmd.Flags().StringVar(&emailAddress, "email-address", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-dependencies-hunt-group-callback
		var huntGroupId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-dependencies-hunt-group-callback",
			Short: "Get Dependencies for a Hunt Group Emergency Callback Number",
			Long:  "Retrieves the emergency callback number dependencies for a specific hunt group.\n\nHunt groups can route incoming calls to a group of people, workspaces or virtual lines. You can even configure a pattern to route to a whole group.\n\nRetrieving the dependencies requires a full, user, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/huntGroups/{huntGroupId}/emergencyCallbackNumber/dependencies")
				req.PathParam("huntGroupId", huntGroupId)
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
		cmd.Flags().StringVar(&huntGroupId, "hunt-group-id", "", "Unique identifier for the hunt group.")
		cmd.MarkFlagRequired("hunt-group-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve Emergency Callback Number attributes for the hunt group under this organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-person-callback
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-callback",
			Short: "Get a Person's Emergency Callback Number",
			Long:  "Retrieve a person's emergency callback number settings.\n\nEmergency Callback Configurations can be enabled at the organization level, Users without individual telephone numbers, such as extension-only users, must be set up with accurate Emergency Callback Numbers (ECBN) and Emergency Service Addresses to enable them to make emergency calls. These users can either utilize the default ECBN for their location or be assigned another specific telephone number from that location for emergency purposes.\n\nTo retrieve a person's callback number requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/emergencyCallbackNumber")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization within which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // update-person-callback
		var personId string
		var orgId string
		var selected string
		var locationMemberId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-callback",
			Short: "Update a Person's Emergency Callback Number",
			Long:  "Update a person's emergency callback number settings.\n\nEmergency Callback Configurations can be enabled at the organization level, Users without individual telephone numbers, such as extension-only users, must be set up with accurate Emergency Callback Numbers (ECBN) to enable them to make emergency calls. These users can either utilize the default ECBN for their location or be assigned another specific telephone number from that location for emergency purposes.\n\nTo update an emergency callback number requires a full, location, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/emergencyCallbackNumber")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("selected", selected)
					req.BodyString("locationMemberId", locationMemberId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization within which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&selected, "selected", "", "")
		cmd.Flags().StringVar(&locationMemberId, "location-member-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-person-callback-dependencies
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-callback-dependencies",
			Short: "Retrieve A Person's Emergency Callback Number Dependencies",
			Long:  "Retrieve Emergency Callback Number dependencies for a person.\n\nEmergency Callback Configurations can be enabled at the organization level, Users without individual telephone numbers, such as extension-only users, must be set up with accurate Emergency Call Back Numbers (ECBN) to enable them to make emergency calls. These users can either utilize the default ECBN for their location or be assigned another specific telephone number from that location for emergency purposes.\n\nRetrieving the dependencies requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/emergencyCallbackNumber/dependencies")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve Emergency Callback Number attributes for this organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-workspace-callback
		var workspaceId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-workspace-callback",
			Short: "Get a Workspace Emergency Callback Number",
			Long:  "Retrieve the emergency callback number setting associated with a specific workspace.\n\nEmergency Callback Configurations can be enabled at the organization level, Users without individual telephone numbers, such as extension-only users, must be set up with accurate Emergency Callback Numbers (ECBN) and Emergency Service Addresses to enable them to make emergency calls. These users can either utilize the default ECBN for their location or be assigned another specific telephone number from that location for emergency purposes.\n\nTo retrieve an emergency callback number, it requires a full, location, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/workspaces/{workspaceId}/emergencyCallbackNumber")
				req.PathParam("workspaceId", workspaceId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "Retrieve Emergency Callback Number attributes for this workspace.")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve Emergency Callback Number attributes for this organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // update-workspace-callback
		var workspaceId string
		var orgId string
		var selected string
		var locationMemberId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-workspace-callback",
			Short: "Update a Workspace Emergency Callback Number",
			Long:  "Update the emergency callback number settings for a workspace.\n\nEmergency Callback Configurations can be enabled at the organization level, Users without individual telephone numbers, such as extension-only users, must be set up with accurate Emergency Call Back Numbers (ECBN) to enable them to make emergency calls. These users can either utilize the default ECBN for their location or be assigned another specific telephone number from that location for emergency purposes.\n\nTo update an emergency callback number requires a full, location, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/workspaces/{workspaceId}/emergencyCallbackNumber")
				req.PathParam("workspaceId", workspaceId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("selected", selected)
					req.BodyString("locationMemberId", locationMemberId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "Updating Emergency Callback Number attributes for this workspace.")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Updating Emergency Callback Number attributes for this organization.")
		cmd.Flags().StringVar(&selected, "selected", "", "")
		cmd.Flags().StringVar(&locationMemberId, "location-member-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-workspace-callback-dependencies
		var workspaceId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-workspace-callback-dependencies",
			Short: "Retrieve Workspace Emergency Callback Number Dependencies",
			Long:  "Retrieve Emergency Callback Number dependencies for a workspace.\n\nEmergency Callback Configurations can be enabled at the organization level, Users without individual telephone numbers, such as extension-only users, must be set up with accurate Emergency Call Back Numbers (ECBN) to enable them to make emergency calls. These users can either utilize the default ECBN for their location or be assigned another specific telephone number from that location for emergency purposes.\n\nRetrieving the dependencies requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/workspaces/{workspaceId}/emergencyCallbackNumber/dependencies")
				req.PathParam("workspaceId", workspaceId)
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
		cmd.Flags().StringVar(&workspaceId, "workspace-id", "", "Retrieve Emergency Callback Number attributes for this workspace.")
		cmd.MarkFlagRequired("workspace-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve Emergency Callback Number attributes for this organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-dependencies-vline-callback
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-dependencies-vline-callback",
			Short: "Get Dependencies for a Virtual Line Emergency Callback Number",
			Long:  "Retrieves the emergency callback number dependencies for a specific virtual line.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines for Webex Calling users.\n\nRetrieving the dependencies requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/emergencyCallbackNumber/dependencies")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List virtual lines for this organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // get-vline-callback-settings
		var virtualLineId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-vline-callback-settings",
			Short: "Get the Virtual Line's Emergency Callback settings",
			Long:  "Retrieves the emergency callback number settings for a specific virtual line.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines for Webex Calling users.\n\nRetrieving the dependencies requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/virtualLines/{virtualLineId}/emergencyCallbackNumber")
				req.PathParam("virtualLineId", virtualLineId)
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
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List virtual lines for this organization.")
		emergencyServicesCmd.AddCommand(cmd)
	}

	{ // update-vline-callback-settings
		var virtualLineId string
		var orgId string
		var selected string
		var locationMemberId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-vline-callback-settings",
			Short: "Update a Virtual Line's Emergency Callback settings",
			Long:  "Update the emergency callback number settings for a specific virtual line.\n\nVirtual line is a capability in Webex Calling that allows administrators to configure multiple lines for Webex Calling users.\n\nTo update virtual line callback number requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/virtualLines/{virtualLineId}/emergencyCallbackNumber")
				req.PathParam("virtualLineId", virtualLineId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("selected", selected)
					req.BodyString("locationMemberId", locationMemberId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&virtualLineId, "virtual-line-id", "", "Unique identifier for the virtual line.")
		cmd.MarkFlagRequired("virtual-line-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List virtual lines for this organization.")
		cmd.Flags().StringVar(&selected, "selected", "", "")
		cmd.Flags().StringVar(&locationMemberId, "location-member-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		emergencyServicesCmd.AddCommand(cmd)
	}

}
