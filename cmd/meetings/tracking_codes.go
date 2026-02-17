package meetings

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

var trackingCodesCmd = &cobra.Command{
	Use:   "tracking-codes",
	Short: "TrackingCodes commands",
}

func init() {
	cmd.MeetingsCmd.AddCommand(trackingCodesCmd)

	{ // list
		var siteUrl string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Tracking Codes",
			Long:  "Lists tracking codes on a site by an admin user.\n\n* If `siteUrl` is specified, tracking codes of the specified site will be listed; otherwise, tracking codes of the user's preferred site will be listed. All available Webex sites and the preferred sites of a user can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\n* Admins can switch any Control Hub managed site from using classic tracking codes to mapped tracking codes in Control Hub. This is a one-time irreversible operation. Once the tracking codes are mapped to custom or user profile attributes, the response returns the mapped tracking codes.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/admin/meeting/config/trackingCodes")
				req.QueryParam("siteUrl", siteUrl)
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
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API retrieves the tracking code from. If not specified, the API retrieves the tracking code from the user's preferred site. All available Webex sites and preferred sites of a user can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		trackingCodesCmd.AddCommand(cmd)
	}

	{ // create
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create",
			Short: "Create a Tracking Code",
			Long:  "Create a new tracking code by an admin user.\n\n* The `siteUrl` is required. The operation creates a tracking code for the specified site. All or a user's available Webex sites can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\n* The `inputMode` of `hostProfileSelect` is only available for a host profile and sign-up pages and does not apply to the meeting scheduler page or the meeting start page. The value for `scheduleStartCodes` must be `null` or the value for all services must be `notUsed` when the `inputMode` is `hostProfileSelect`.\n\n* The `hostProfileCode` of `required` is only allowed for a Site Admin managed site, and not for a Control Hub managed site.\n\n* When the `hostProfileCode` is `adminSet`, only `adminSet`, `notUsed`, and `notApplicable` are available for the types of `scheduleStartCodes`. When the `hostProfileCode` is not `adminSet`, only `optional`, `required`, `notUsed`, and `notApplicable` are available for `scheduleStartCodes`.\n\n* If the type of the `All` service has a value other than `notApplicable`, and another service, e.g. `EventCenter`, is missing from the `scheduleStartCodes`, then the type of this missing `EventCenter` service shares the same type as the `All` service. If the type of `All` service has a value other than `notApplicable`, and another service, e.g. `EventCenter`, has a type, then the type specified should be the same as the `All` service.\n\n* If the `All` service is missing from the `scheduleStartCodes`, any of the other four services, e.g. `EventCenter`, have a default type of `notUsed` if it is also missing from the `scheduleStartCodes`.\n\n* Admins can switch any Control Hub managed site from using classic tracking codes to mapped tracking codes in Control Hub, this is a one-time irreversible operation. Once the tracking codes are mapped to custom or user profile attributes, they cannot create tracking codes when the mapping process is in progress or the mapping process is completed.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/admin/meeting/config/trackingCodes")
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
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		trackingCodesCmd.AddCommand(cmd)
	}

	{ // get
		var trackingCodeId string
		var siteUrl string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get a Tracking Code",
			Long:  "Retrieves details for a tracking code by an admin user.\n\n* If `siteUrl` is specified, the tracking code is retrieved from the specified site; otherwise, the tracking code is retrieved from the user's preferred site. All available Webex sites and the preferred sites of a user can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\n* Admins can switch any Control Hub managed site from using classic tracking codes to mapped tracking codes in Control Hub, this is a one-time irreversible operation. Once the tracking codes are mapped to custom or user profile attributes, the response returns details for a mapped tracking code.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/admin/meeting/config/trackingCodes/{trackingCodeId}")
				req.PathParam("trackingCodeId", trackingCodeId)
				req.QueryParam("siteUrl", siteUrl)
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
		cmd.Flags().StringVar(&trackingCodeId, "tracking-code-id", "", "Unique identifier for the tracking code whose details are being requested.")
		cmd.MarkFlagRequired("tracking-code-id")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site which the API retrieves the tracking code from. If not specified, the API retrieves the tracking code from the user's preferred site. All available Webex sites and the preferred sites of a user can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		trackingCodesCmd.AddCommand(cmd)
	}

	{ // update
		var trackingCodeId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update a Tracking Code",
			Long:  "Updates details for a tracking code by an admin user.\n\n* The `siteUrl` is required. The operation updates a tracking code for the specified site. All of a user's available Webex sites can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\n* The `inputMode` of `hostProfileSelect` is only available for the host profile and sign-up pages and it doesn't apply to the meeting scheduler page or meeting start page. Therefore, `scheduleStartCodes` must be `null` or type of all services must be `notUsed` when the `inputMode` is `hostProfileSelect`.\n\n* Currently, the `hostProfileCode` of `required` is only allowed for a Site Admin managed site, and not allowed for a Control Hub managed site.\n\n* When the `hostProfileCode` is `adminSet`, only `adminSet`, `notUsed` and `notApplicable` are available for the types of `scheduleStartCodes`. When the `hostProfileCode` is not `adminSet`, only `optional`, `required`, `notUsed` and `notApplicable` are available for types of `scheduleStartCodes`.\n\n* If the type of the `All` service has a value other than `notApplicable`, and another service, e.g. `EventCenter`, is missing from the `scheduleStartCodes`, then the type of this missing `EventCenter` service shares the same type as the `All` service silently. If the type of `All` service has a value other than `notApplicable`, and another service, e.g. `EventCenter`, has a type, then the type specified should be the same as the `All` service.\n\n* If the `All` service is missing from the `scheduleStartCodes`, any of the other four services, e.g. `EventCenter`, has a default type of `notUsed` if that service is also missing from the `scheduleStartCodes`.\n\n* Admins can switch any Control Hub managed site from using classic tracking codes to mapped tracking codes in Control Hub, this is a one-time irreversible operation. Once the tracking codes are mapped to custom or user profile attributes, they cannot update tracking codes when the mapping process is in progress or the mapping process is completed.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/admin/meeting/config/trackingCodes/{trackingCodeId}")
				req.PathParam("trackingCodeId", trackingCodeId)
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
		cmd.Flags().StringVar(&trackingCodeId, "tracking-code-id", "", "")
		cmd.MarkFlagRequired("tracking-code-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		trackingCodesCmd.AddCommand(cmd)
	}

	{ // delete
		var trackingCodeId string
		var siteUrl string
		cmd := &cobra.Command{
			Use:   "delete",
			Short: "Delete a Tracking Code",
			Long:  "Deletes a tracking code by an admin user.\n\n* The `siteUrl` is required. The operation deletes a tracking code for the specified site. All of a user's available Webex sites can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.\n\n* Admins can switch any Control Hub managed site from using classic tracking codes to mapped tracking codes in Control Hub, this is a one-time irreversible operation. Once the tracking codes are mapped to custom or user profile attributes, they cannot delete tracking codes when the mapping process is in progress or the mapping process is completed.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/admin/meeting/config/trackingCodes/{trackingCodeId}")
				req.PathParam("trackingCodeId", trackingCodeId)
				req.QueryParam("siteUrl", siteUrl)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&trackingCodeId, "tracking-code-id", "", "Unique identifier for the tracking code to be deleted.")
		cmd.MarkFlagRequired("tracking-code-id")
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site from which the API deletes the tracking code. All available Webex sites and preferred sites of a user can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		trackingCodesCmd.AddCommand(cmd)
	}

	{ // get-user
		var siteUrl string
		var personId string
		var email string
		cmd := &cobra.Command{
			Use:   "get-user",
			Short: "Get User Tracking Codes",
			Long:  "Lists user's tracking codes by an admin user.\n\n* At least one parameter, either `personId`, or `email` is required. `personId` must come before `email` if both are specified. Please note that `email` is specified in the request header.\n\n* If `siteUrl` is specified, the tracking codes of the specified site will be listed; otherwise, the tracking codes of a user's preferred site are listed. All available Webex sites and preferred sites of a user can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API. Please note that the user here is the admin user who invokes the API, not the user specified by `personId` or email.\n\n* Admins can switch any Control Hub managed site from using classic tracking codes to mapped tracking codes in Control Hub, this is a one-time irreversible operation. Once the tracking codes are mapped to custom or user profile attributes, the response returns the user's mapped tracking codes.\n\n#### Request Header\n\n* `email`: Email address for the user whose tracking codes are being retrieved. The admin users can specify the email of a user on a site they manage and the API returns details for the user's tracking codes. At least one parameter of `personId` or `email` is required.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/admin/meeting/userconfig/trackingCodes")
				req.QueryParam("siteUrl", siteUrl)
				req.QueryParam("personId", personId)
				req.Header("email", email)
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
		cmd.Flags().StringVar(&siteUrl, "site-url", "", "URL of the Webex site from which the API retrieves the tracking code. If not specified, the API retrieves the tracking code from the user's preferred site. All available Webex sites and preferred sites of a user can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API.")
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the user whose tracking codes are being retrieved. The admin user can specify the `personId` of a user on a site they manage and the API returns details for the user's tracking codes. At least one parameter of `personId` or `email` is required.")
		cmd.Flags().StringVar(&email, "email", "", "e.g. john.andersen@example.com")
		trackingCodesCmd.AddCommand(cmd)
	}

	{ // update-user
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-user",
			Short: "Update User Tracking Codes",
			Long:  "Updates tracking codes for a specified user by an admin user.\n\n* The `siteUrl` is required. The operation updates a user's tracking code on the specified site. All a user's available Webex sites can be retrieved by the [Get Site List](/docs/api/v1/meeting-preferences/get-site-list) API. Please note that the user here is the admin user who invokes the API, not the user specified by `personId` or `email`.\n\n* A name that is not found in the site-level tracking codes cannot be set for a user's tracking codes. All available site-level tracking codes for a site can be retrieved by the [List Tracking Codes](/docs/api/v1/tracking-codes/list-tracking-codes) API.\n\n* If the `inputMode` of a user's tracking code is `select` or `hostProfileSelect`, its value must be one of the site-level options of that tracking code. All available site-level tracking codes for a site can be retrieved by the [List Tracking Codes](/docs/api/v1/tracking-codes/list-tracking-codes) API.\n\n* Admins can switch any Control Hub managed site from using classic tracking codes to mapped tracking codes in Control Hub, this is a one-time irreversible operation. Once the tracking codes are mapped to custom or user profile attributes, they cannot update user's tracking codes when the mapping process is in progress or the mapping process is completed.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/admin/meeting/userconfig/trackingCodes")
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
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		trackingCodesCmd.AddCommand(cmd)
	}

}
