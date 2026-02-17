package calling

import (
	"fmt"
	"strconv"
	"strings"

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
var _ = strings.Join

var userCallCmd = &cobra.Command{
	Use:   "user-call",
	Short: "UserCall commands",
}

func init() {
	cmd.CallingCmd.AddCommand(userCallCmd)

	{ // get-person-application-services-settings
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-application-services-settings",
			Short: "Retrieve a person's Application Services Settings",
			Long: `Application services let you determine the ringing behavior for calls made to people in certain scenarios. You can also specify which devices can download the Webex Calling app.

This API requires a full, user, or read-only administrator or location administrator auth token with a scope of spark-admin:people_read.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/applications")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "(Required) Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-application-services-settings
		var personId string
		var orgId string
		var ringDevicesForClickToDialCallsEnabled bool
		var ringDevicesForGroupPageEnabled bool
		var ringDevicesForCallParkEnabled bool
		var browserClientEnabled bool
		var desktopClientEnabled bool
		var tabletClientEnabled bool
		var mobileClientEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-application-services-settings",
			Short: "Modify a person's Application Services Settings",
			Long: `Application services let you determine the ringing behavior for calls made to users in certain scenarios. You can also specify which devices users can download the Webex Calling app on.

This API requires a full or user administrator or location administrator auth token with the spark-admin:people_write scope.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/applications")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("ringDevicesForClickToDialCallsEnabled", ringDevicesForClickToDialCallsEnabled, cmd.Flags().Changed("ring-devices-for-click-to-dial-calls-enabled"))
					req.BodyBool("ringDevicesForGroupPageEnabled", ringDevicesForGroupPageEnabled, cmd.Flags().Changed("ring-devices-for-group-page-enabled"))
					req.BodyBool("ringDevicesForCallParkEnabled", ringDevicesForCallParkEnabled, cmd.Flags().Changed("ring-devices-for-call-park-enabled"))
					req.BodyBool("browserClientEnabled", browserClientEnabled, cmd.Flags().Changed("browser-client-enabled"))
					req.BodyBool("desktopClientEnabled", desktopClientEnabled, cmd.Flags().Changed("desktop-client-enabled"))
					req.BodyBool("tabletClientEnabled", tabletClientEnabled, cmd.Flags().Changed("tablet-client-enabled"))
					req.BodyBool("mobileClientEnabled", mobileClientEnabled, cmd.Flags().Changed("mobile-client-enabled"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&ringDevicesForClickToDialCallsEnabled, "ring-devices-for-click-to-dial-calls-enabled", false, "")
		cmd.Flags().BoolVar(&ringDevicesForGroupPageEnabled, "ring-devices-for-group-page-enabled", false, "")
		cmd.Flags().BoolVar(&ringDevicesForCallParkEnabled, "ring-devices-for-call-park-enabled", false, "")
		cmd.Flags().BoolVar(&browserClientEnabled, "browser-client-enabled", false, "")
		cmd.Flags().BoolVar(&desktopClientEnabled, "desktop-client-enabled", false, "")
		cmd.Flags().BoolVar(&tabletClientEnabled, "tablet-client-enabled", false, "")
		cmd.Flags().BoolVar(&mobileClientEnabled, "mobile-client-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-barge-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-barge-settings-person",
			Short: "Read Barge In Settings for a Person",
			Long:  "Retrieve a person's Barge In settings.\n\nThe Barge In feature enables you to use a Feature Access Code (FAC) to answer a call that was directed to another subscriber, or barge-in on the call if it was already answered. Barge In can be used across locations.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read` or a user auth token with `spark:people_read` scope can be used by a person to read their own settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/bargeIn")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-barge-settings-person
		var personId string
		var orgId string
		var enabled bool
		var toneEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-barge-settings-person",
			Short: "Configure Barge In Settings for a Person",
			Long:  "Configure a person's Barge In settings.\n\nThe Barge In feature enables you to use a Feature Access Code (FAC) to answer a call that was directed to another subscriber, or barge-in on the call if it was already answered. Barge In can be used across locations.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope or a user auth token with `spark:people_write` scope can be used by a person to update their own settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/bargeIn")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyBool("toneEnabled", toneEnabled, cmd.Flags().Changed("tone-enabled"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().BoolVar(&toneEnabled, "tone-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-forwarding-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-forwarding-settings-person",
			Short: "Read Forwarding Settings for a Person",
			Long:  "Retrieve a person's Call Forwarding settings.\n\nThree types of call forwarding are supported:\n\n+ Always - forwards all incoming calls to the destination you choose.\n\n+ When busy - forwards all incoming calls to the destination you chose while the phone is in use or the person is busy.\n\n+ When no answer - forwarding only occurs when you are away or not answering your phone.\n\nIn addition, the Business Continuity feature will send calls to a destination of your choice if your phone is not connected to the network for any reason, such as a power outage, failed Internet connection, or wiring problem.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read` or a user auth token with `spark:people_read` scope can be used by a person to read their own settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/callForwarding")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-forward-person
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-forward-person",
			Short: "Configure Call Forwarding Settings for a Person",
			Long:  "Configure a person's Call Forwarding settings.\n\nThree types of call forwarding are supported:\n\n+ Always - forwards all incoming calls to the destination you choose.\n\n+ When busy - forwards all incoming calls to the destination you chose while the phone is in use or the person is busy.\n\n+ When no answer - forwarding only occurs when you are away or not answering your phone.\n\nIn addition, the Business Continuity feature will send calls to a destination of your choice if your phone is not connected to the network for any reason, such as a power outage, failed Internet connection, or wiring problem.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope or a user auth token with `spark:people_write` scope can be used by a person to update their settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/callForwarding")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-intercept-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-intercept-settings-person",
			Short: "Read Call Intercept Settings for a Person",
			Long:  "Retrieves Person's Call Intercept settings.\n\nThe intercept feature gracefully takes a person's phone out of service, while providing callers with informative announcements and alternative routing options. Depending on the service configuration, none, some, or all incoming calls to the specified person are intercepted. Also depending on the service configuration, outgoing calls are intercepted or rerouted to another location.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/intercept")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-intercept-settings-person
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-intercept-settings-person",
			Short: "Configure Call Intercept Settings for a Person",
			Long:  "Configures a person's Call Intercept settings.\n\nThe intercept feature gracefully takes a person's phone out of service, while providing callers with informative announcements and alternative routing options. Depending on the service configuration, none, some, or all incoming calls to the specified person are intercepted. Also depending on the service configuration, outgoing calls are intercepted or rerouted to another location.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/intercept")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-intercept-greeting-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "update-intercept-greeting-person",
			Short: "Configure Call Intercept Greeting for a Person",
			Long:  "Configure a person's Call Intercept Greeting by uploading a Waveform Audio File Format, `.wav`, encoded audio file.\n\nYour request will need to be a `multipart/form-data` request rather than JSON, using the `audio/wav` Content-Type.\n\nThis API requires a full or user administrator auth token with the `spark-admin:people_write` scope or a user auth token with `spark:people_write` scope can be used by a person to update their settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/people/{personId}/features/intercept/actions/announcementUpload/invoke")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-recording-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-recording-person",
			Short: "Read Call Recording Settings for a Person",
			Long:  "Retrieve a person's Call Recording settings.\n\nThe Call Recording feature provides a hosted mechanism to record the calls placed and received on the Carrier platform for replay and archival. This feature is helpful for quality assurance, security, training, and more.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_read` scope.\n\n<div><Callout type=\"warning\">A person with a Webex Calling Standard license is eligible for the Call Recording feature only when the Call Recording vendor is Webex.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/callRecording")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-recording-person
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-recording-person",
			Short: "Configure Call Recording Settings for a Person",
			Long:  "Configure a person's Call Recording settings.\n\nThe Call Recording feature provides a hosted mechanism to record the calls placed and received on the Carrier platform for replay and archival. This feature is helpful for quality assurance, security, training, and more.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.\n\n<div><Callout type=\"warning\">A person with a Webex Calling Standard license is eligible for the Call Recording feature only when the Call Recording vendor is Webex.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/callRecording")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-waiting-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-waiting-settings-person",
			Short: "Read Call Waiting Settings for a Person",
			Long:  "Retrieve a person's Call Waiting settings.\n\nWith this feature, a person can place an active call on hold and answer an incoming call.  When enabled, while you are on an active call, a tone alerts you of an incoming call and you can choose to answer or ignore the call.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/callWaiting")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-waiting-settings-person
		var personId string
		var orgId string
		var enabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-waiting-settings-person",
			Short: "Configure Call Waiting Settings for a Person",
			Long:  "Configure a person's Call Waiting settings.\n\nWith this feature, a person can place an active call on hold and answer an incoming call.  When enabled, while you are on an active call, a tone alerts you of an incoming call and you can choose to answer or ignore the call.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/callWaiting")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-caller-id-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-caller-id-person",
			Short: "Read Caller ID Settings for a Person",
			Long:  "Retrieve a person's Caller ID settings.\n\nCaller ID settings control how a person's information is displayed when making outgoing calls.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, `dialByFirstName`, and `dialByLastName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/callerId")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-caller-id-person
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-caller-id-person",
			Short: "Configure Caller ID Settings for a Person",
			Long:  "Configure a person's Caller ID settings.\n\nCaller ID settings control how a person's information is displayed when making outgoing calls.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.<div><Callout type=\"warning\">The fields `directLineCallerIdName.selection`, `directLineCallerIdName.customName`, `dialByFirstName`, and `dialByLastName` are not supported in Webex for Government (FedRAMP). Instead, administrators must use the `firstName` and `lastName` fields to configure and view both caller ID and dial-by-name settings.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/callerId")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-calling-behavior
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-calling-behavior",
			Short: "Read Person's Calling Behavior",
			Long:  "Retrieves the calling behavior and UC Manager Profile settings for the person which includes overall calling behavior and calling UC Manager Profile ID.\n\nWebex Calling Behavior controls which Webex telephony application and which UC Manager Profile is to be used for a person.\n\nAn organization has an organization-wide default Calling Behavior that may be overridden for individual persons.\n\nUC Manager Profiles are applicable if your organization uses Jabber in Team Messaging mode or Calling in Webex (Unified CM).\n\nThe UC Manager Profile also has an organization-wide default and may be overridden for individual persons.\n\nThis API requires a full, user, or read-only administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/callingBehavior")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-calling-behavior
		var personId string
		var orgId string
		var behaviorType string
		var profileId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-calling-behavior",
			Short: "Configure a person's Calling Behavior",
			Long:  "Modifies the calling behavior settings for the person which includes calling behavior and UC Manager Profile ID.\n\nWebex Calling Behavior controls which Webex telephony application and which UC Manager Profile is to be used for a person.\n\nAn organization has an organization-wide default Calling Behavior that may be overridden for individual persons.\n\nUC Manager Profiles are applicable if your organization uses Jabber in Team Messaging mode or Calling in Webex (Unified CM).\n\nThe UC Manager Profile also has an organization-wide default and may be overridden for individual persons.\n\nThis API requires a full or user administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/callingBehavior")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("behaviorType", behaviorType)
					req.BodyString("profileId", profileId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&behaviorType, "behavior-type", "", "")
		cmd.Flags().StringVar(&profileId, "profile-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-do-not-disturb-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-do-not-disturb-settings-person",
			Short: "Read Do Not Disturb Settings for a Person",
			Long:  "Retrieve a person's Do Not Disturb settings.\n\nWhen enabled, this feature will give all incoming calls the busy treatment. Optionally, you can enable a Ring Reminder to play a brief tone on your desktop phone when you receive incoming calls.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read` or a user auth token with `spark:people_read` scope can be used by a person to read their settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/doNotDisturb")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-do-not-disturb-settings-person
		var personId string
		var orgId string
		var webexGoOverrideEnabled bool
		var enabled bool
		var ringSplashEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-do-not-disturb-settings-person",
			Short: "Configure Do Not Disturb Settings for a Person",
			Long:  "Configure a person's Do Not Disturb settings.\n\nWhen enabled, this feature will give all incoming calls the busy treatment. Optionally, you can enable a Ring Reminder to play a brief tone on your desktop phone when you receive incoming calls.\n\nThis API requires a full or user administrator auth token with the `spark-admin:people_write` scope or a user auth token with `spark:people_write` scope can be used by a person to update their settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/doNotDisturb")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("webexGoOverrideEnabled", webexGoOverrideEnabled, cmd.Flags().Changed("webex-go-override-enabled"))
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyBool("ringSplashEnabled", ringSplashEnabled, cmd.Flags().Changed("ring-splash-enabled"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&webexGoOverrideEnabled, "webex-go-override-enabled", false, "")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().BoolVar(&ringSplashEnabled, "ring-splash-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-executive-assistant-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-executive-assistant-settings-person",
			Short: "Retrieve Executive Assistant Settings for a Person",
			Long:  "Retrieve the executive assistant settings for the specified `personId`.\n\nPeople with the executive service enabled, can select from a pool of assistants who have been assigned the executive assistant service and who can answer or place calls on their behalf. Executive assistants can set the call forward destination and join or leave an executive's pool.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/executiveAssistant")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-executive-assistant-settings-person
		var personId string
		var orgId string
		var typeVal string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-executive-assistant-settings-person",
			Short: "Modify Executive Assistant Settings for a Person",
			Long:  "Modify the executive assistant settings for the specified personId.\n\nPeople with the executive service enabled, can select from a pool of assistants who have been assigned the executive assistant service and who can answer or place calls on their behalf. Executive assistants can set the call forward destination and join or leave an executive's pool.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/executiveAssistant")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("type", typeVal)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&typeVal, "type", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-hoteling-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-hoteling-settings-person",
			Short: "Read Hoteling Settings for a Person",
			Long:  "Retrieve a person's hoteling settings.\n\nAs an administrator, you can enable hoteling for people so that their phone profile (phone number, features, and calling plan) is temporarily loaded onto a shared (host) phone.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/hoteling")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-hoteling-settings-person
		var personId string
		var orgId string
		var enabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-hoteling-settings-person",
			Short: "Configure Hoteling Settings for a Person",
			Long:  "Configure a person's hoteling settings.\n\nAs an administrator, you can enable hoteling for people so that their phone profile (phone number, features, and calling plan) is temporarily loaded onto a shared (host) phone.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/hoteling")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-monitoring-settings
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-monitoring-settings",
			Short: "Retrieve a person's Monitoring Settings",
			Long:  "Retrieves the monitoring settings of the person, which shows specified people, places, virtual lines or call park extenions that are being monitored.\nMonitors the line status which indicates if a person, place or virtual line is on a call and if a call has been parked on that extension.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/monitoring")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-monitoring-settings
		var personId string
		var orgId string
		var enableCallParkNotification bool
		var monitoredElements []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-monitoring-settings",
			Short: "Modify a person's Monitoring Settings",
			Long:  "Modifies the monitoring settings of the person.\nMonitors the line status of specified people, places, virtual lines or call park extension. The line status indicates if a person, place or virtual line is on a call and if a call has been parked on that extension.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/monitoring")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enableCallParkNotification", enableCallParkNotification, cmd.Flags().Changed("enable-call-park-notification"))
					req.BodyStringSlice("monitoredElements", monitoredElements)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&enableCallParkNotification, "enable-call-park-notification", false, "")
		cmd.Flags().StringSliceVar(&monitoredElements, "monitored-elements", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-incoming-permission-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-incoming-permission-settings-person",
			Short: "Read Incoming Permission Settings for a Person",
			Long:  "Retrieve a person's Incoming Permission settings.\n\nYou can change the incoming calling permissions for a person if you want them to be different from your organization's default.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/incomingPermission")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-incoming-permission-settings-person
		var personId string
		var orgId string
		var useCustomEnabled bool
		var externalTransfer string
		var internalCallsEnabled bool
		var collectCallsEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-incoming-permission-settings-person",
			Short: "Configure Incoming Permission Settings for a Person",
			Long:  "Configure a person's Incoming Permission settings.\n\nYou can change the incoming calling permissions for a person if you want them to be different from your organization's default.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/incomingPermission")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("useCustomEnabled", useCustomEnabled, cmd.Flags().Changed("use-custom-enabled"))
					req.BodyString("externalTransfer", externalTransfer)
					req.BodyBool("internalCallsEnabled", internalCallsEnabled, cmd.Flags().Changed("internal-calls-enabled"))
					req.BodyBool("collectCallsEnabled", collectCallsEnabled, cmd.Flags().Changed("collect-calls-enabled"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&useCustomEnabled, "use-custom-enabled", false, "")
		cmd.Flags().StringVar(&externalTransfer, "external-transfer", "", "")
		cmd.Flags().BoolVar(&internalCallsEnabled, "internal-calls-enabled", false, "")
		cmd.Flags().BoolVar(&collectCallsEnabled, "collect-calls-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-outgoing-permissions-settings
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-outgoing-permissions-settings",
			Short: "Retrieve a person's Outgoing Calling Permissions Settings",
			Long:  "Retrieve a person's Outgoing Calling Permissions settings.\n\nOutgoing calling permissions regulate behavior for calls placed to various destinations and default to the local level settings. You can change the outgoing calling permissions for a person if you want them to be different from your organization's default.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/outgoingPermission")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-outgoing-permissions-settings
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-outgoing-permissions-settings",
			Short: "Modify a person's Outgoing Calling Permissions Settings",
			Long:  "Modify a person's Outgoing Calling Permissions settings.\n\nOutgoing calling permissions regulate behavior for calls placed to various destinations and default to the local level settings. You can change the outgoing calling permissions for a person if you want them to be different from your organization's default.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/outgoingPermission")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // list-phone-numbers-person
		var personId string
		var orgId string
		var preferE164Format string
		cmd := &cobra.Command{
			Use:   "list-phone-numbers-person",
			Short: "Get a List of Phone Numbers for a Person",
			Long:  "Get a person's phone numbers including alternate numbers.\n\nA person can have one or more phone numbers and/or extensions via which they can be called.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_read` scope.\n\n<br/>\n\n<div><Callout type=\"warning\">The `preferE164Format` query parameter can be used to get phone numbers either in E.164 format or in their legacy format. The support for getting phone numbers in non-E.164 format in some geographies will be removed in the future.</Callout></div>",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/numbers")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("preferE164Format", preferE164Format)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&preferE164Format, "prefer-e164-format", "", "Return phone numbers in E.164 format.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-privacy-settings
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-privacy-settings",
			Short: "Get a person's Privacy Settings",
			Long: `Get a person's privacy settings for the specified person ID.

The privacy feature enables the person's line to be monitored by others and determine if they can be reached by Auto Attendant services.

This API requires a full, user, or read-only administrator or location administrator auth token with a scope of spark-admin:people_read.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/privacy")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-privacy-settings
		var personId string
		var orgId string
		var aaExtensionDialingEnabled bool
		var aaNamingDialingEnabled bool
		var enablePhoneStatusDirectoryPrivacy bool
		var enablePhoneStatusPickupBargeInPrivacy bool
		var monitoringAgents []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-privacy-settings",
			Short: "Configure a person's Privacy Settings",
			Long: `Configure a person's privacy settings for the specified person ID.

The privacy feature enables the person's line to be monitored by others and determine if they can be reached by Auto Attendant services.

This API requires a full or user administrator or location administrator auth token with the spark-admin:people_write scope.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/privacy")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("aaExtensionDialingEnabled", aaExtensionDialingEnabled, cmd.Flags().Changed("aa-extension-dialing-enabled"))
					req.BodyBool("aaNamingDialingEnabled", aaNamingDialingEnabled, cmd.Flags().Changed("aa-naming-dialing-enabled"))
					req.BodyBool("enablePhoneStatusDirectoryPrivacy", enablePhoneStatusDirectoryPrivacy, cmd.Flags().Changed("enable-phone-status-directory-privacy"))
					req.BodyBool("enablePhoneStatusPickupBargeInPrivacy", enablePhoneStatusPickupBargeInPrivacy, cmd.Flags().Changed("enable-phone-status-pickup-barge-in-privacy"))
					req.BodyStringSlice("monitoringAgents", monitoringAgents)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&aaExtensionDialingEnabled, "aa-extension-dialing-enabled", false, "")
		cmd.Flags().BoolVar(&aaNamingDialingEnabled, "aa-naming-dialing-enabled", false, "")
		cmd.Flags().BoolVar(&enablePhoneStatusDirectoryPrivacy, "enable-phone-status-directory-privacy", false, "")
		cmd.Flags().BoolVar(&enablePhoneStatusPickupBargeInPrivacy, "enable-phone-status-pickup-barge-in-privacy", false, "")
		cmd.Flags().StringSliceVar(&monitoringAgents, "monitoring-agents", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-push-to-talk-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-push-to-talk-settings-person",
			Short: "Read Push-to-Talk Settings for a Person",
			Long:  "Retrieve a person's Push-to-Talk settings.\n\nPush-to-Talk allows the use of desk phones as either a one-way or two-way intercom that connects people in different parts of your organization.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/pushToTalk")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-push-to-talk-settings-person
		var personId string
		var orgId string
		var allowAutoAnswer bool
		var connectionType string
		var accessType string
		var members []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-push-to-talk-settings-person",
			Short: "Configure Push-to-Talk Settings for a Person",
			Long:  "Configure a person's Push-to-Talk settings.\n\nPush-to-Talk allows the use of desk phones as either a one-way or two-way intercom that connects people in different parts of your organization.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/pushToTalk")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("allowAutoAnswer", allowAutoAnswer, cmd.Flags().Changed("allow-auto-answer"))
					req.BodyString("connectionType", connectionType)
					req.BodyString("accessType", accessType)
					req.BodyStringSlice("members", members)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&allowAutoAnswer, "allow-auto-answer", false, "")
		cmd.Flags().StringVar(&connectionType, "connection-type", "", "")
		cmd.Flags().StringVar(&accessType, "access-type", "", "")
		cmd.Flags().StringSliceVar(&members, "members", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-receptionist-client-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-receptionist-client-settings-person",
			Short: "Read Receptionist Client Settings for a Person",
			Long:  "Retrieve a person's Receptionist Client settings.\n\nTo help support the needs of your front-office personnel, you can set up people, workspaces or virtual lines as telephone attendants so that they can screen all incoming calls to certain numbers within your organization.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/reception")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-receptionist-client-settings-person
		var personId string
		var orgId string
		var receptionEnabled bool
		var monitoredMembers []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-receptionist-client-settings-person",
			Short: "Configure Receptionist Client Settings for a Person",
			Long:  "Configure a person's Receptionist Client settings.\n\nTo help support the needs of your front-office personnel, you can set up people, workspaces or virtual lines as telephone attendants so that they can screen all incoming calls to certain numbers within your organization.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/reception")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("receptionEnabled", receptionEnabled, cmd.Flags().Changed("reception-enabled"))
					req.BodyStringSlice("monitoredMembers", monitoredMembers)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&receptionEnabled, "reception-enabled", false, "")
		cmd.Flags().StringSliceVar(&monitoredMembers, "monitored-members", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // list-schedules-person
		var personId string
		var orgId string
		var start string
		var max string
		var name string
		var typeVal string
		cmd := &cobra.Command{
			Use:   "list-schedules-person",
			Short: "List of Schedules for a Person",
			Long:  "List schedules for a person in an organization.\n\nSchedules are used to support calling features and can be defined at the location or person level. `businessHours` schedules allow you to apply specific call settings at different times of the day or week by defining one or more events. `holidays` schedules define exceptions to normal business hours by defining one or more events.\n\nThis API requires a full, user, or read-only administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/schedules")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
				req.QueryParam("max", max)
				req.QueryParam("name", name)
				req.QueryParam("type", typeVal)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&start, "start", "", "Specifies the offset from the first result that you want to fetch.")
		cmd.Flags().StringVar(&max, "max", "", "Specifies the maximum number of records that you want to fetch.")
		cmd.Flags().StringVar(&name, "name", "", "Specifies the case insensitive substring to be matched against the schedule names. The maximum length is 40.")
		cmd.Flags().StringVar(&typeVal, "type", "", "Specifies the schedule event type to be matched on the given type.")
		userCallCmd.AddCommand(cmd)
	}

	{ // create-schedule-person
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-schedule-person",
			Short: "Create Schedule for a Person",
			Long:  "Create a new schedule for a person.\n\nSchedules are used to support calling features and can be defined at the location or person level. `businessHours` schedules allow you to apply specific call settings at different times of the day or week by defining one or more events. `holidays` schedules define exceptions to normal business hours by defining one or more events.\n\nThis API requires a full or user administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/people/{personId}/features/schedules")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-schedule
		var personId string
		var scheduleType string
		var scheduleId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-schedule",
			Short: "Get a Schedule Details",
			Long:  "Retrieve a schedule by its schedule ID.\n\nSchedules are used to support calling features and can be defined at the location or person level. `businessHours` schedules allow you to apply specific call settings at different times of the day or week by defining one or more events. `holidays` schedules define exceptions to normal business hours by defining one or more events.\n\nThis API requires a full, user, or read-only administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/schedules/{scheduleType}/{scheduleId}")
				req.PathParam("personId", personId)
				req.PathParam("scheduleType", scheduleType)
				req.PathParam("scheduleId", scheduleId)
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
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "Type of schedule, either `businessHours` or `holidays`.")
		cmd.MarkFlagRequired("schedule-type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Unique identifier for the schedule.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-schedule
		var personId string
		var scheduleType string
		var scheduleId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-schedule",
			Short: "Update a Schedule",
			Long:  "Modify a schedule by its schedule ID.\n\nSchedules are used to support calling features and can be defined at the location or person level. `businessHours` schedules allow you to apply specific call settings at different times of the day or week by defining one or more events. `holidays` schedules define exceptions to normal business hours by defining one or more events.\n\nThis API requires a full or user administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/schedules/{scheduleType}/{scheduleId}")
				req.PathParam("personId", personId)
				req.PathParam("scheduleType", scheduleType)
				req.PathParam("scheduleId", scheduleId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "Type of schedule, either `businessHours` or `holidays`.")
		cmd.MarkFlagRequired("schedule-type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Unique identifier for the schedule.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // delete-schedule
		var personId string
		var scheduleType string
		var scheduleId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-schedule",
			Short: "Delete a Schedule",
			Long:  "Delete a schedule by its schedule ID.\n\nSchedules are used to support calling features and can be defined at the location or person level. `businessHours` schedules allow you to apply specific call settings at different times of the day or week by defining one or more events. `holidays` schedules define exceptions to normal business hours by defining one or more events.\n\nThis API requires a full or user administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/people/{personId}/features/schedules/{scheduleType}/{scheduleId}")
				req.PathParam("personId", personId)
				req.PathParam("scheduleType", scheduleType)
				req.PathParam("scheduleId", scheduleId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "Type of schedule, either `businessHours` or `holidays`.")
		cmd.MarkFlagRequired("schedule-type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Unique identifier for the schedule.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-event-person-schedule
		var personId string
		var scheduleType string
		var scheduleId string
		var eventId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-event-person-schedule",
			Short: "Fetch Event for a person's Schedule",
			Long:  "People can use shared location schedules or define personal schedules containing events.\n\n`businessHours` schedules allow you to apply specific call settings at different times of the day or week by defining one or more events. `holidays` schedules define exceptions to normal business hours by defining one or more events.\n\nThis API requires a full, user, or read-only administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/schedules/{scheduleType}/{scheduleId}/events/{eventId}")
				req.PathParam("personId", personId)
				req.PathParam("scheduleType", scheduleType)
				req.PathParam("scheduleId", scheduleId)
				req.PathParam("eventId", eventId)
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
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "Type of schedule, either `businessHours` or `holidays`.")
		cmd.MarkFlagRequired("schedule-type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Unique identifier for the schedule.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&eventId, "event-id", "", "Unique identifier for the event.")
		cmd.MarkFlagRequired("event-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-event-person-schedule
		var personId string
		var scheduleType string
		var scheduleId string
		var eventId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-event-person-schedule",
			Short: "Update an Event for a person's Schedule",
			Long:  "People can use shared location schedules or define personal schedules containing events.\n\n`businessHours` schedules allow you to apply specific call settings at different times of the day or week by defining one or more events. `holidays` schedules define exceptions to normal business hours by defining one or more events.\n\nThis API requires a full or user administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/schedules/{scheduleType}/{scheduleId}/events/{eventId}")
				req.PathParam("personId", personId)
				req.PathParam("scheduleType", scheduleType)
				req.PathParam("scheduleId", scheduleId)
				req.PathParam("eventId", eventId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "Type of schedule, either `businessHours` or `holidays`.")
		cmd.MarkFlagRequired("schedule-type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Unique identifier for the schedule.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&eventId, "event-id", "", "Unique identifier for the event.")
		cmd.MarkFlagRequired("event-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // delete-event-person-schedule
		var personId string
		var scheduleType string
		var scheduleId string
		var eventId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-event-person-schedule",
			Short: "Delete an Event for a person's Schedule",
			Long:  "People can use shared location schedules or define personal schedules containing events.\n\n`businessHours` schedules allow you to apply specific call settings at different times of the day or week by defining one or more events. `holidays` schedules define exceptions to normal business hours by defining one or more events.\n\nThis API requires a full or user administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/people/{personId}/features/schedules/{scheduleType}/{scheduleId}/events/{eventId}")
				req.PathParam("personId", personId)
				req.PathParam("scheduleType", scheduleType)
				req.PathParam("scheduleId", scheduleId)
				req.PathParam("eventId", eventId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "Type of schedule, either `businessHours` or `holidays`.")
		cmd.MarkFlagRequired("schedule-type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Unique identifier for the schedule.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&eventId, "event-id", "", "Unique identifier for the event.")
		cmd.MarkFlagRequired("event-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // create-event-person-schedule
		var personId string
		var scheduleType string
		var scheduleId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-event-person-schedule",
			Short: "Add a New Event for Person's Schedule",
			Long:  "People can use shared location schedules or define personal schedules containing events.\n\n`businessHours` schedules allow you to apply specific call settings at different times of the day or week by defining one or more events. `holidays` schedules define exceptions to normal business hours by defining one or more events.\n\nThis API requires a full or user administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/people/{personId}/features/schedules/{scheduleType}/{scheduleId}/events")
				req.PathParam("personId", personId)
				req.PathParam("scheduleType", scheduleType)
				req.PathParam("scheduleId", scheduleId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "Type of schedule, either `businessHours` or `holidays`.")
		cmd.MarkFlagRequired("schedule-type")
		cmd.Flags().StringVar(&scheduleId, "schedule-id", "", "Unique identifier for the schedule.")
		cmd.MarkFlagRequired("schedule-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-voicemail-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-voicemail-person",
			Short: "Read Voicemail Settings for a Person",
			Long:  "Retrieve a person's Voicemail settings.\n\nThe voicemail feature transfers callers to voicemail based on your settings. You can then retrieve voice messages via Voicemail. Voicemail audio is sent in Waveform Audio File Format, `.wav`, format.\n\nOptionally, notifications can be sent to a mobile phone via text or email. These notifications will not include the voicemail files.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read` or a user auth token with `spark:people_read` scope can be used by a person to read their settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/voicemail")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-voicemail-person
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-voicemail-person",
			Short: "Configure Voicemail Settings for a Person",
			Long:  "Configure a person's Voicemail settings.\n\nThe voicemail feature transfers callers to voicemail based on your settings. You can then retrieve voice messages via Voicemail. Voicemail audio is sent in Waveform Audio File Format, `.wav`, format.\n\nOptionally, notifications can be sent to a mobile phone via text or email. These notifications will not include the voicemail files.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope or a user auth token with `spark:people_write` scope can be used by a person to update their settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/people/{personId}/features/voicemail")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-busy-voicemail-greeting-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "update-busy-voicemail-greeting-person",
			Short: "Configure Busy Voicemail Greeting for a Person",
			Long:  "Configure a person's Busy Voicemail Greeting by uploading a Waveform Audio File Format, `.wav`, encoded audio file.\n\nYour request will need to be a `multipart/form-data` request rather than JSON, using the `audio/wav` Content-Type.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope or a user auth token with `spark:people_write` scope can be used by a person to update their settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/people/{personId}/features/voicemail/actions/uploadBusyGreeting/invoke")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-no-answer-voicemail-greeting-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "update-no-answer-voicemail-greeting-person",
			Short: "Configure No Answer Voicemail Greeting for a Person",
			Long:  "Configure a person's No Answer Voicemail Greeting by uploading a Waveform Audio File Format, `.wav`, encoded audio file.\n\nYour request will need to be a `multipart/form-data` request rather than JSON, using the `audio/wav` Content-Type.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope or a user auth token with `spark:people_write` scope can be used by a person to update their settings.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/people/{personId}/features/voicemail/actions/uploadNoAnswerGreeting/invoke")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // reset-voicemail-pin
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "reset-voicemail-pin",
			Short: "Reset Voicemail PIN",
			Long:  "Reset a voicemail PIN for a person.\n\nThe voicemail feature transfers callers to voicemail based on your settings. You can then retrieve voice messages via Voicemail.  A voicemail PIN is used to retrieve your voicemail messages.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.\n\n**NOTE**: This API is expected to have an empty request body and Content-Type header should be set to `application/json`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/people/{personId}/features/voicemail/actions/resetPin/invoke")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // validate-initiate-move-job
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "validate-initiate-move-job",
			Short: "Validate or Initiate Move Users Job",
			Long:  "This API allows the user to perform one of the following operations:\n\n* Setting the `validate` attribute to `true` validates the user move.\n\n* Setting the `validate` attribute to `false` performs the user move.\n\n<br/>\n\nNotes:\n\n* A maximum of `100` users can be moved at a time.\n\n* Setting the `validate` attribute to `true` only allowed for calling user.\n\n* When a single non calling user is moved, it will be moved synchronously without creating any job.\n\n<br/>\n\nErrors occurring during the initial API request validation are captured directly in the error response, along with the appropriate HTTP status code.\n\n<br/>\n\nBelow is a list of possible error `code` values and their associated `message`, which can be found in the `errors` array during initial API request validation, regardless of the `validate` attribute value:\n\n* BATCH-400 - Attribute 'User ID' is required.\n\n* BATCH-400 - Users list should not be empty.\n\n* BATCH-400 - Users should not be empty.\n\n* 1026006 - Attribute 'Validate' is required.\n\n* 1026010 - User is not a valid Calling User.\n\n* 1026011 - Users list should not be empty.\n\n* 1026012 - Users should not be empty.\n\n* 1026013 - The source and the target location cannot be the same.\n\n* 1026014 - Error occurred while processing the move users request.\n\n* 1026015 - Error occurred while moving user number to target location.\n\n* 1026016 - User should have either phone number or extension.\n\n* 1026017 - Phone number is not in e164 format.\n\n* 1026018 - Selected Users list exceeds the maximum limit.\n\n* 1026019 - Duplicate entry for user is not allowed.\n\n* 1026020 - Validate 'true' is supported only for single user.\n\n* 1026021 - Attribute location id is required for Calling user.\n\n* 1026022 - Validate 'true' is supported for calling users only.\n\n* 1026023 - Extension and phone number is supported for calling users only.\n\n* 2150012 - User was not found\n\n<br/>\n\nWhen the `validate` attribute is set to true, the API identifies and returns the `errors` and `impacts` associated with the user move in the response.\n\n<br/>\n\nBelow is a list of possible error `code` values and their associated `message`, which can be found in the `errors` array, when `validate` attribute is set to be true:\n\n* 4003 - `User Not Found`\n\n* 4007 - `User Not Found`\n\n* 4152 - `Location Not Found`\n\n* 5620 - `Location Not Found`\n\n* 4202 - `The extension is not available. It is already assigned to a user : {0}`\n\n* 8264 - `Routing profile is different with new group: {0}`\n\n* 19600 - `User has to be within an enterprise to be moved.`\n\n* 19601 - `User can only be moved to a different group within the same enterprise.`\n\n* 19602 - `Only regular end user can be moved. Service instance virtual user cannot be moved.`\n\n* 19603 - `New group already reaches maximum number of user limits.`\n\n* 19604 - `The {0} number of the user is the same as the calling line ID of the group.`\n\n* 19605 - `User is assigned services not authorized to the new group: {0}.`\n\n* 19606 - `User is in an active hoteling/flexible seating association.`\n\n* 19607 - `User is pilot user of a trunk group.`\n\n* 19608 - `User is using group level device profiles which is used by other users in current group. Following are the device profiles shared with other users: {0}.`\n\n* 19609 - `Following device profiles cannot be moved to the new group because there are already devices with the same name defined in the new group: {0}.`\n\n* 19610 - `The extension of the user is used as transfer to operator number for following Auto Attendent : {0}.`\n\n* 19611 - `Fail to move announcement file from {0} to {1}.`\n\n* 19612 - `Fail to move device management file from {0} to {1}.`\n\n* 19613 - `User is assigned service packs not authorized to the new group: {0}.`\n\n* 25008 - `Missing Mandatory field name: {0}`\n\n* 25110 - `{fieldName} cannot be less than {0} or greater than {1} characters.`\n\n* 25378 - `Target location is same as user's current location.`\n\n* 25379 - `Error Occurred while Fetching User's Current Location Id.`\n\n* 25381 - `Error Occurred while rolling back to Old Location Call recording Settings`\n\n* 25382 - `Error Occurred while Disabling Call Recording for user which is required Before User can be Moved`\n\n* 25383 - `OCI Error while moving user`\n\n* 25384 - `Error Occurred while checking for Possible Call Recording Impact.`\n\n* 25385 - `Error Occurred while getting Call Recording Settings`\n\n* 27559 - `The groupExternalId search criteria contains groups with different calling zone.`\n\n* 27960 - `Parameter isWebexCalling, newPhoneNumber, or newExtension can only be set in Webex Calling deployment mode.`\n\n* 27961 - `Parameter isWebexCalling shall be set if newPhoneNumber or newExtension is set.`\n\n* 27962 - `Work space cannot be moved.`\n\n* 27963 - `Virtual profile user cannot be moved.`\n\n* 27965 - `The user's phone number: {0}, is same as the current group charge number.`\n\n* 27966 - `The phone number, {0}, is not available in the new group.`\n\n* 27967 - `User is configured as the ECBN user for another user in the current group.`\n\n* 27968 - `User is configured as the ECBN user for the current group.`\n\n* 27969 - `User is associated with DECT handset(s): {0}`\n\n* 27970 - `User is using a customer managed device: {0}`\n\n* 27971 - `User is using an ATA device: {0}`\n\n* 27972 - `User is in an active hotdesking association.`\n\n* 27975 - `Need to unassign CLID number from group before moving the number to the new group. Phone number: {0}`\n\n* 27976 - `Local Gateway configuration is different with new group. Phone number: {0}`\n\n* 1026015 - `Error occurred while moving user number to target location`\n\n* 10010000 - `Total numbers exceeded maximum limit allowed`\n\n* 10010001 - `to-location and from-location cannot be same`\n\n* 10010002 - `to-location and from-location should belong to same customer`\n\n* 10010003 - `to-location must have a carrier`\n\n* 10010004 - `from-location must have a carrier`\n\n* 10010005 - `Different Carrier move is not supported for non-Cisco PSTN carriers.`\n\n* 10010006 - `Number move not supported for WEBEX_DIRECT carriers.`\n\n* 10010007 - `Numbers out of sync, missing on CPAPI`\n\n* 10010008 - `from-location not found or pstn connection missing in CPAPI`\n\n* 10010010 - `from-location is in transition`\n\n* 10010009 - `to-location not found or pstn connection missing in CPAPI`\n\n* 10010011 - `to-location is in transition`\n\n* 10010012 - `Numbers don't have a carrier Id`\n\n* 10010013 - `Location less numbers don't have a carrier Id`\n\n* 10010014 - `Different Carrier move is not supported for numbers with different country or region.`\n\n* 10010015 - `Numbers contain mobile and non-mobile types.`\n\n* 10010016 - `To/From location carriers must be same for mobile numbers.`\n\n* 10010017 - `Move request for location less number not supported`\n\n* 10010200 - `Move request for more than one block number is not supported`\n\n* 10010201 - `Cannot move block number as few numbers not from the block starting %s to %s`\n\n* 10010202 - `Cannot move block number as few numbers failed VERIFICATION from the block %s to %s`\n\n* 10010203 - `Cannot move block number as few numbers missing from the block %s to %s`\n\n* 10010204 - `Cannot move number as it is NOT a part of the block %s to %s`\n\n* 10010205 - `Move request for Cisco PSTN block order not supported.`\n\n* 10010299 - `Move order couldn't be created as no valid number to move`\n\n* 10030000 - `Number not found`\n\n* 10030001 - `Number does not belong to from-location`\n\n* 10030002 - `Number is not present in CPAPI`\n\n* 10030003 - `Number assigned to an user or device`\n\n* 10030004 - `Number not in Active status`\n\n* 10030005 - `Number is set as main number of the location`\n\n* 10030006 - `Number has pending order associated with it`\n\n* 10030007 - `Number belongs to a location but a from-location was not set`\n\n* 10030008 - `Numbers from multiple carrier ids are not supported`\n\n* 10030009 - `Location less number belongs to a location. from-location value is set to null or no location id`\n\n* 10030010 - `One or more numbers are not portable.`\n\n* 10030011 - `Mobile number carrier was not set`\n\n* 10030012 - `Number must be assigned for assigned move`\n\n* 10050000 - `Failed to update customer reference for phone numbers on carrier`\n\n* 10050001 - `Failed to update customer reference`\n\n* 10050002 - `Order is not of operation type MOVE`\n\n* 10050003 - `CPAPI delete call failed`\n\n* 10050004 - `Not found in database`\n\n* 10050005 - `Error sending notification to WxcBillingService`\n\n* 10050006 - `CPAPI provision number as active call failed with status %s ,reason %s`\n\n* 10050007 - `Failed to update E911 Service`\n\n* 10050008 - `Target location does not have Inbound Toll Free license`\n\n* 10050009 - `Source location or Target location subscription found cancelled or suspended`\n\n* 10050010 - `Moving On Premises or Non Integrated CCP numbers from one location to another is not supported.`\n\n* 10099999 - `{Error Code} - {Error Message}`\n\n<br/>\n\nBelow is a list of possible impact `code` values and their associated `message`, which can be found in the `impacts` array, when `validate` attribute is set to be true:\n\n* 19701 - `The identity/device profile the user is using is moved to the new group: {0}.`\n\n* 19702 - `The user level customized incoming digit string setting is removed from the user. User is set to use the new group setting.`\n\n* 19703 - `The user level customized outgoing digit plan setting is removed from the user. User is set to use the new group setting.`\n\n* 19704 - `The user level customized enhanced outgoing calling plan setting is removed from the user. User is set to use the new group setting.`\n\n* 19705 - `User is removed from following group services: {0}.`\n\n* 19706 - `The current group schedule used in any criteria is removed from the service settings.`\n\n* 19707 - `User is removed from the department of the old group.`\n\n* 19708 - `User is changed to use the default communication barring profile of the new group.`\n\n* 19709 - `The communication barring profile of the user is assigned to the new group: {0}.`\n\n* 19710 - `The charge number for the user is removed.`\n\n* 19711 - `The disabled FACs for the user are removed because they are not available in the new group.`\n\n* 19712 - `User is removed from trunk group.`\n\n* 19713 - `The extension of the user is reset to empty due to either the length is out of bounds of the new group, or the extension is already taken in new group.`\n\n* 19714 - `The extension of the following alternate number is reset to empty due to either the length out of bounds of the new group or the extension is already taken in new group: {0}.`\n\n* 19715 - `The collaborate room using current group default collaborate bridge is moved to the default collaborate bridge of the new group.`\n\n* 19716 - `Previously stored voice messages of the user are no longer available. The new voice message will be stored on the mail server of the new group.`\n\n* 19717 - `The primary number, alternate numbers or fax messaging number of the user are assigned to the new group: {0}.`\n\n* 19718 - `Following domains are assigned to the new group: {0}.`\n\n* 19719 - `The NCOS of the user is assigned to the new group: {0}.`\n\n* 19720 - `The office zone of the user is assigned to the new group: {0}.`\n\n* 19721 - `The announcement media files are relocated to the new group directory.`\n\n* 19722 - `User CLID number is set to use the new group CLID number: {0}.`\n\n* 19723 - `New group CLID number is not configured.`\n\n* 19724 - `The group level announcement file(s) are removed from the user's music on hold settings.`\n\n* 25388 - `Target Location Does not Have Vendor Configured. Call Recording for user will be disabled`\n\n* 25389 - `Call Recording Vendor for user will be changed from:{0} to:{1}`\n\n* 25390 - `Dub point of user is moved to new external group`\n\n* 25391 - `Error Occurred while moving Call recording Settings to new location`\n\n* 25392 - `Error Occurred while checking for Possible Call Recording Impact.`\n\n* 25393 - `Sending Billing Notification Failed`\n\nThis API requires a full administrator auth token with the scopes `spark-admin:telephony_config_write`, `spark-admin:people_write`, and `identity:groups_rw`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/jobs/person/moveLocation")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Create Move Users job for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // list-move-jobs
		var orgId string
		var start string
		var max string
		cmd := &cobra.Command{
			Use:   "list-move-jobs",
			Short: "List Move Users Jobs",
			Long:  "Lists all the Move Users jobs for the given organization in order of most recent job to oldest job irrespective of its status.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/jobs/person/moveLocation")
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve list of Move Users jobs for this organization.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of jobs. Default is 0.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of jobs returned to this maximum count. Default is 2000.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-move-job-status
		var jobId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-move-job-status",
			Short: "Get Move Users Job Status",
			Long:  "Returns the status and other details of the job.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/jobs/person/moveLocation/{jobId}")
				req.PathParam("jobId", jobId)
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
		cmd.Flags().StringVar(&jobId, "job-id", "", "Retrieve job details for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve job details for this organization.")
		userCallCmd.AddCommand(cmd)
	}

	{ // pause-move-job
		var jobId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "pause-move-job",
			Short: "Pause the Move Users Job",
			Long:  "Pause the running Move Users Job. A paused job can be resumed.\n\nThis API requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/jobs/person/moveLocation/{jobId}/actions/pause/invoke")
				req.PathParam("jobId", jobId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&jobId, "job-id", "", "Pause the Move Users job for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Pause the Move Users job for this organization.")
		userCallCmd.AddCommand(cmd)
	}

	{ // resume-move-job
		var jobId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "resume-move-job",
			Short: "Resume the Move Users Job",
			Long:  "Resume the paused Move Users Job that is in paused status.\n\nThis API requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/jobs/person/moveLocation/{jobId}/actions/resume/invoke")
				req.PathParam("jobId", jobId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&jobId, "job-id", "", "Resume the Move Users job for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Resume the Move Users job for this organization.")
		userCallCmd.AddCommand(cmd)
	}

	{ // list-move-job-errors
		var jobId string
		var orgId string
		var start string
		var max string
		cmd := &cobra.Command{
			Use:   "list-move-job-errors",
			Short: "List Move Users Job errors",
			Long:  "Lists all error details of Move Users job. This will not list any errors if `exitCode` is `COMPLETED`. If the status is `COMPLETED_WITH_ERRORS` then this lists the cause of failures.\n\nThis API requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/jobs/person/moveLocation/{jobId}/errors")
				req.PathParam("jobId", jobId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&jobId, "job-id", "", "Retrieve the error details for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve list of jobs for this organization.")
		cmd.Flags().StringVar(&start, "start", "", "Specifies the error offset from the first result that you want to fetch.")
		cmd.Flags().StringVar(&max, "max", "", "Specifies the maximum number of records that you want to fetch.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-music-hold-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-music-hold-settings-person",
			Short: "Retrieve Music On Hold Settings for a Person",
			Long:  "Retrieve the person's music on hold settings.\n\nMusic on hold is played when a caller is put on hold, or the call is parked.\n\nRetrieving a person's music on hold settings requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/musicOnHold")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-music-hold-settings-person
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-music-hold-settings-person",
			Short: "Configure Music On Hold Settings for a Person",
			Long:  "Configure a person's music on hold settings.\n\nMusic on hold is played when a caller is put on hold, or the call is parked.\n\nTo configure music on hold settings for a person, music on hold setting must be enabled for this location.\n\nUpdating a person's music on hold settings requires a full or user administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/musicOnHold")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-access-codes-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-access-codes-person",
			Short: "Retrieve Access Codes for a Person",
			Long:  "Retrieve the person's access codes.\n\nAccess codes are used to bypass permissions.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/outgoingPermission/accessCodes")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-access-codes-person
		var personId string
		var orgId string
		var deleteCodes []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-access-codes-person",
			Short: "Modify Access Codes for a Person",
			Long:  "Modify a person's access codes.\n\nAccess codes are used to bypass permissions.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/outgoingPermission/accessCodes")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("deleteCodes", deleteCodes)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		cmd.Flags().StringSliceVar(&deleteCodes, "delete-codes", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // create-access-codes-person
		var personId string
		var orgId string
		var code string
		var description string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-access-codes-person",
			Short: "Create Access Codes for a Person",
			Long:  "Create new Access codes for the person.\n\nAccess codes are used to bypass permissions.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/people/{personId}/outgoingPermission/accessCodes")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("code", code)
					req.BodyString("description", description)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		cmd.Flags().StringVar(&code, "code", "", "")
		cmd.Flags().StringVar(&description, "description", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // delete-access-codes-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-access-codes-person",
			Short: "Delete Access Codes for a Person",
			Long:  "Deletes all Access codes for the person.\n\nAccess codes are used to bypass permissions.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/people/{personId}/outgoingPermission/accessCodes")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-transfer-numbers-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-transfer-numbers-person",
			Short: "Retrieve Transfer Numbers for a Person",
			Long:  "Retrieve the person's transfer numbers.\n\nWhen calling a specific call type, this person will be automatically transferred to another number. The person assigned the Auto Transfer Number can then approve the call and send it through or reject the call type. You can add up to 3 numbers.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/outgoingPermission/autoTransferNumbers")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-transfer-numbers-person
		var personId string
		var orgId string
		var useCustomTransferNumbers bool
		var autoTransferNumber1 string
		var autoTransferNumber2 string
		var autoTransferNumber3 string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-transfer-numbers-person",
			Short: "Modify Transfer Numbers for a Person",
			Long:  "Modify a person's transfer numbers.\n\nWhen calling a specific call type, this person will be automatically transferred to another number. The person assigned the Auto Transfer Number can then approve the call and send it through or reject the call type. You can add up to 3 numbers.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/outgoingPermission/autoTransferNumbers")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("useCustomTransferNumbers", useCustomTransferNumbers, cmd.Flags().Changed("use-custom-transfer-numbers"))
					req.BodyString("autoTransferNumber1", autoTransferNumber1)
					req.BodyString("autoTransferNumber2", autoTransferNumber2)
					req.BodyString("autoTransferNumber3", autoTransferNumber3)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		cmd.Flags().BoolVar(&useCustomTransferNumbers, "use-custom-transfer-numbers", false, "")
		cmd.Flags().StringVar(&autoTransferNumber1, "auto-transfer-number1", "", "")
		cmd.Flags().StringVar(&autoTransferNumber2, "auto-transfer-number2", "", "")
		cmd.Flags().StringVar(&autoTransferNumber3, "auto-transfer-number3", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-digit-patterns-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-digit-patterns-person",
			Short: "Retrieve Digit Patterns for a Person",
			Long:  "Retrieve the person's digit patterns.\n\nDigit patterns are used to bypass permissions.\n\nRetrieving digit patterns requires a full or user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/outgoingPermission/digitPatterns")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // create-digit-patterns-person
		var personId string
		var orgId string
		var name string
		var pattern string
		var action string
		var transferEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-digit-patterns-person",
			Short: "Create Digit Patterns for a Person",
			Long:  "Create a new digit pattern for the given person.\n\nDigit patterns are used to bypass permissions.\n\nCreating the digit pattern requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/people/{personId}/outgoingPermission/digitPatterns")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("pattern", pattern)
					req.BodyString("action", action)
					req.BodyBool("transferEnabled", transferEnabled, cmd.Flags().Changed("transfer-enabled"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&pattern, "pattern", "", "")
		cmd.Flags().StringVar(&action, "action", "", "")
		cmd.Flags().BoolVar(&transferEnabled, "transfer-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-digit-pattern-control-person
		var personId string
		var orgId string
		var useCustomDigitPatterns bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-digit-pattern-control-person",
			Short: "Modify the Digit Pattern Category Control Settings for a Person",
			Long:  "Modifies whether this user uses the specified digit patterns when placing outbound calls or not.\n\nUpdating the digit pattern category control settings requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/outgoingPermission/digitPatterns")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("useCustomDigitPatterns", useCustomDigitPatterns, cmd.Flags().Changed("use-custom-digit-patterns"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		cmd.Flags().BoolVar(&useCustomDigitPatterns, "use-custom-digit-patterns", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // delete-all-digit-patterns-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-all-digit-patterns-person",
			Short: "Delete all Digit Patterns for a Person",
			Long:  "Delete all digit patterns for a Person.\n\nDigit patterns are used to bypass permissions.\n\nDeleting the digit patterns requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/people/{personId}/outgoingPermission/digitPatterns")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-digit-pattern-person
		var personId string
		var digitPatternId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-digit-pattern-person",
			Short: "Retrieve Digit Pattern Details for a Person",
			Long:  "Retrieve the digit pattern details for a person.\n\nDigit patterns are used to bypass permissions.\n\nRetrieving the digit pattern details requires a full or user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/outgoingPermission/digitPatterns/{digitPatternId}")
				req.PathParam("personId", personId)
				req.PathParam("digitPatternId", digitPatternId)
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
		cmd.Flags().StringVar(&digitPatternId, "digit-pattern-id", "", "Unique identifier for the digit pattern.")
		cmd.MarkFlagRequired("digit-pattern-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-digit-pattern-person
		var personId string
		var digitPatternId string
		var orgId string
		var name string
		var pattern string
		var action string
		var transferEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-digit-pattern-person",
			Short: "Modify a Digit Pattern for a Person",
			Long:  "Modify a digit pattern for a Person.\n\nDigit patterns are used to bypass permissions.\n\nUpdating the digit pattern requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/outgoingPermission/digitPatterns/{digitPatternId}")
				req.PathParam("personId", personId)
				req.PathParam("digitPatternId", digitPatternId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("pattern", pattern)
					req.BodyString("action", action)
					req.BodyBool("transferEnabled", transferEnabled, cmd.Flags().Changed("transfer-enabled"))
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
		cmd.Flags().StringVar(&digitPatternId, "digit-pattern-id", "", "Unique identifier for the digit pattern.")
		cmd.MarkFlagRequired("digit-pattern-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&pattern, "pattern", "", "")
		cmd.Flags().StringVar(&action, "action", "", "")
		cmd.Flags().BoolVar(&transferEnabled, "transfer-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // delete-digit-pattern-person
		var personId string
		var digitPatternId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-digit-pattern-person",
			Short: "Delete a Digit Pattern for a Person",
			Long:  "Delete a digit pattern for a Person.\n\nDigit patterns are used to bypass permissions.\n\nDeleting the digit pattern requires a full or user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/people/{personId}/outgoingPermission/digitPatterns/{digitPatternId}")
				req.PathParam("personId", personId)
				req.PathParam("digitPatternId", digitPatternId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&digitPatternId, "digit-pattern-id", "", "Unique identifier for the digit pattern.")
		cmd.MarkFlagRequired("digit-pattern-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access the API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // assign-unassign-numbers-person
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "assign-unassign-numbers-person",
			Short: "Assign or Unassign numbers to a person",
			Long:  "Assign or unassign alternate phone numbers to a person.\n\nEach location has a set of phone numbers that can be assigned to people, workspaces, or features. Phone numbers must follow the E.164 format for all countries, except for the United States, which can also follow the National format. Active phone numbers are in service.\n\nAssigning or unassigning an alternate phone number to a person requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/numbers")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "Unique identitfier of the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization of the Route Group.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-preferred-endpoint
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-preferred-endpoint",
			Short: "Get Preferred Answer Endpoint",
			Long:  "Get the person's preferred answer endpoint and the list of endpoints available for selection. The preferred answer endpoint is null if one has not been selected. The list of endpoints is empty if the person has no endpoints assigned which support the preferred answer endpoint functionality. These endpoints can be used by the following Call Control API's that allow the person to specify an endpointId to use for the call:<br>\n\n+ [/v1/telephony/calls/dial](/docs/api/v1/call-controls/dial)<br>\n\n+ [/v1/telephony/calls/retrieve](/docs/api/v1/call-controls/retrieve)<br>\n\n+ [/v1/telephony/calls/pickup](/docs/api/v1/call-controls/pickup)<br>\n\n+ [/v1/telephony/calls/barge-in](/docs/api/v1/call-controls/barge-in)<br>\n\n+ [/v1/telephony/calls/answer](/docs/api/v1/call-controls/answer)<br>\n\nThis API requires `spark:telephony_config_read` or `spark-admin:telephony_config_read` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/preferredAnswerEndpoint")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-preferred-endpoint
		var personId string
		var orgId string
		var preferredAnswerEndpointId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-preferred-endpoint",
			Short: "Modify Preferred Answer Endpoint",
			Long:  "Sets or clears the person\u2019s preferred answer endpoint. To clear the preferred answer endpoint the `preferredAnswerEndpointId` attribute must be set to null.<br>\nThis API requires `spark:telephony_config_write` or `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/preferredAnswerEndpoint")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("preferredAnswerEndpointId", preferredAnswerEndpointId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&preferredAnswerEndpointId, "preferred-answer-endpoint-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // search-shared-line-appearance-members
		var personId string
		var applicationId string
		var max string
		var start string
		var location string
		var name string
		var number string
		var order string
		var extension string
		cmd := &cobra.Command{
			Use:   "search-shared-line-appearance-members",
			Short: "Search Shared-Line Appearance Members",
			Long:  "Get members available for shared-line assignment to a Webex Calling Apps Desktop device.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_read` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/applications/{applicationId}/availableMembers")
				req.PathParam("personId", personId)
				req.PathParam("applicationId", applicationId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("location", location)
				req.QueryParam("name", name)
				req.QueryParam("number", number)
				req.QueryParam("order", order)
				req.QueryParam("extension", extension)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&applicationId, "application-id", "", "A unique identifier for the application.")
		cmd.MarkFlagRequired("application-id")
		cmd.Flags().StringVar(&max, "max", "", "Number of records per page.")
		cmd.Flags().StringVar(&start, "start", "", "Page number.")
		cmd.Flags().StringVar(&location, "location", "", "Location ID for the user.")
		cmd.Flags().StringVar(&name, "name", "", "Search for users with names that match the query.")
		cmd.Flags().StringVar(&number, "number", "", "Search for users with numbers that match the query.")
		cmd.Flags().StringVar(&order, "order", "", "Sort by first name (`fname`) or last name (`lname`).")
		cmd.Flags().StringVar(&extension, "extension", "", "Search for users with extensions that match the query.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-shared-line-appearance-members
		var personId string
		var applicationId string
		cmd := &cobra.Command{
			Use:   "get-shared-line-appearance-members",
			Short: "Get Shared-Line Appearance Members",
			Long:  "Get primary and secondary members assigned to a shared line on a Webex Calling Apps Desktop device.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_read` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/applications/{applicationId}/members")
				req.PathParam("personId", personId)
				req.PathParam("applicationId", applicationId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&applicationId, "application-id", "", "A unique identifier for the application.")
		cmd.MarkFlagRequired("application-id")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-shared-line-appearance-members
		var personId string
		var applicationId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-shared-line-appearance-members",
			Short: "Put Shared-Line Appearance Members",
			Long:  "Add or modify primary and secondary users assigned to shared-lines on a Webex Calling Apps Desktop device.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/applications/{applicationId}/members")
				req.PathParam("personId", personId)
				req.PathParam("applicationId", applicationId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&applicationId, "application-id", "", "A unique identifier for the application.")
		cmd.MarkFlagRequired("application-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-voicemail-passcode
		var personId string
		var orgId string
		var passcode string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-voicemail-passcode",
			Short: "Modify a person's voicemail passcode",
			Long:  "Modify a person's voicemail passcode.\n\nModifying a person's voicemail passcode requires a full administrator, user administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/voicemail/passcode")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("passcode", passcode)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "Modify voicemail passcode for this person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Modify voicemail passcode for a person in this organization.")
		cmd.Flags().StringVar(&passcode, "passcode", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-message-summary
		cmd := &cobra.Command{
			Use:   "get-message-summary",
			Short: "Get Message Summary",
			Long:  `Get a summary of the voicemail messages for the user.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/voiceMessages/summary")
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
		userCallCmd.AddCommand(cmd)
	}

	{ // list-messages
		cmd := &cobra.Command{
			Use:   "list-messages",
			Short: "List Messages",
			Long:  `Get the list of all voicemail messages for the user.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/voiceMessages")
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
		userCallCmd.AddCommand(cmd)
	}

	{ // delete-message
		var messageId string
		cmd := &cobra.Command{
			Use:   "delete-message",
			Short: "Delete Message",
			Long:  `Delete a specfic voicemail message for the user.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/voiceMessages/{messageId}")
				req.PathParam("messageId", messageId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&messageId, "message-id", "", "The message identifer of the voicemail message to delete")
		cmd.MarkFlagRequired("message-id")
		userCallCmd.AddCommand(cmd)
	}

	{ // mark-read
		var messageId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "mark-read",
			Short: "Mark As Read",
			Long:  "Update the voicemail message(s) as read for the user.\n\nIf the `messageId` is provided, then only mark that message as read.  Otherwise, all messages for the user are marked as read.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/voiceMessages/markAsRead")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("messageId", messageId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&messageId, "message-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // mark-unread
		var messageId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "mark-unread",
			Short: "Mark As Unread",
			Long:  "Update the voicemail message(s) as unread for the user.\n\nIf the `messageId` is provided, then only mark that message as unread.  Otherwise, all messages for the user are marked as unread.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/voiceMessages/markAsUnread")
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("messageId", messageId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&messageId, "message-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-agent-list-available-caller-ids
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-agent-list-available-caller-ids",
			Short: "Retrieve Agent's List of Available Caller IDs",
			Long:  "Get the list of call queues and hunt groups available for caller ID use by this person as an agent.\n\nThis API requires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/agent/availableCallerIds")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-agent-caller-id
		var personId string
		cmd := &cobra.Command{
			Use:   "get-agent-caller-id",
			Short: "Retrieve Agent's Caller ID Information",
			Long:  "Retrieve the Agent's Caller ID Information.\n\nEach agent will be able to set their outgoing Caller ID as either the Call Queue's Caller ID, Hunt Group's Caller ID or their own configured Caller ID.\n\nThis API requires a full admin or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/agent/callerId")
				req.PathParam("personId", personId)
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
		userCallCmd.AddCommand(cmd)
	}

	{ // update-agent-caller-id
		var personId string
		var selectedCallerId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-agent-caller-id",
			Short: "Modify Agent's Caller ID Information",
			Long:  "Modify Agent's Caller ID Information.\n\nEach Agent will be able to set their outgoing Caller ID as either the designated Call Queue's Caller ID or Hunt Group's Caller ID or their own configured Caller ID\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/agent/callerId")
				req.PathParam("personId", personId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("selectedCallerId", selectedCallerId)
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
		cmd.Flags().StringVar(&selectedCallerId, "selected-caller-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-bridge-settings-person
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-bridge-settings-person",
			Short: "Read Call Bridge Settings for a Person",
			Long:  "Retrieve a person's Call Bridge settings.\n\nThis API requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:people_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/features/callBridge")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-bridge-settings-person
		var personId string
		var orgId string
		var warningToneEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-bridge-settings-person",
			Short: "Configure Call Bridge Settings for a Person",
			Long:  "Configure a person's Call Bridge settings.\n\nThis API requires a full or user administrator or location administrator auth token with the `spark-admin:people_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/features/callBridge")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("warningToneEnabled", warningToneEnabled, cmd.Flags().Changed("warning-tone-enabled"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		cmd.Flags().BoolVar(&warningToneEnabled, "warning-tone-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-secondary-available-numbers
		var personId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-person-secondary-available-numbers",
			Short: "Get Person Secondary Available Phone Numbers",
			Long:  "List standard numbers that are available to be assigned as a person's secondary phone number.\nThese numbers are associated with the location of the person specified in the request URL, can be active or inactive, and are unassigned.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/secondary/availableNumbers")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-fax-available-numbers
		var personId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-person-fax-available-numbers",
			Short: "Get Person Fax Message Available Phone Numbers",
			Long:  "List standard numbers that are available to be assigned as a person's FAX message number.\nThese numbers are associated with the location of the person specified in the request URL, can be active or inactive, and are unassigned.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/faxMessage/availableNumbers")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-forward-available-numbers
		var personId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		var ownerName string
		var extension string
		cmd := &cobra.Command{
			Use:   "get-person-forward-available-numbers",
			Short: "Get Person Call Forward Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as a person's call forward number.\nThese numbers are associated with the location of the person specified in the request URL, can be active or inactive, and are assigned to an owning entity.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/callForwarding/availableNumbers")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("ownerName", ownerName)
				req.QueryParam("extension", extension)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		cmd.Flags().StringVar(&ownerName, "owner-name", "", "Return the list of phone numbers that are owned by the given `ownerName`. Maximum length is 255.")
		cmd.Flags().StringVar(&extension, "extension", "", "Returns the list of PSTN phone numbers with the given `extension`.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-primary-numbers
		var orgId string
		var locationId string
		var max string
		var start string
		var phoneNumber string
		var licenseType string
		cmd := &cobra.Command{
			Use:   "get-person-primary-numbers",
			Short: "Get Person Primary Available Phone Numbers",
			Long:  "List numbers that are available to be assigned as a person's primary phone number.\nBy default, this API returns standard and mobile numbers from all locations that are unassigned. The parameters `licenseType` and `locationId` must align with the person's settings to determine the appropriate number for assignment.\nFailure to provide these parameters may result in the unsuccessful assignment of the returned number.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/primary/availableNumbers")
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationId", locationId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("licenseType", licenseType)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the list of phone numbers for this location within the given organization. The maximum length is 36.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		cmd.Flags().StringVar(&licenseType, "license-type", "", "Used to search numbers according to the person's `licenseType` to which the number will be assigned.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-ecbn-available-numbers
		var personId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		var ownerName string
		cmd := &cobra.Command{
			Use:   "get-person-ecbn-available-numbers",
			Short: "Get Person ECBN Available Phone Numbers",
			Long:  "List standard numbers that are available to be assigned as a person's emergency callback number.\nThese numbers are associated with the location of the person specified in the request URL, can be active or inactive, and are assigned to an owning entity.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/emergencyCallbackNumber/availableNumbers")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("ownerName", ownerName)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		cmd.Flags().StringVar(&ownerName, "owner-name", "", "Return the list of phone numbers that are owned by the given `ownerName`. Maximum length is 255.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-intercept-available-numbers
		var personId string
		var orgId string
		var max string
		var start string
		var phoneNumber string
		var ownerName string
		var extension string
		cmd := &cobra.Command{
			Use:   "get-person-intercept-available-numbers",
			Short: "Get Person Call Intercept Available Phone Numbers",
			Long:  "List the service and standard PSTN numbers that are available to be assigned as a person's call intercept number.\nThese numbers are associated with the location specified in the request URL, can be active or inactive, and are assigned to an owning entity.\n\nThe available numbers APIs help identify candidate numbers and their owning entities to simplify the assignment or association of these numbers to members or features.\n\nRetrieving this list requires a full, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/callIntercept/availableNumbers")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("ownerName", ownerName)
				req.QueryParam("extension", extension)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List numbers for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of phone numbers returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching phone numbers. The default is 0.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Filter phone numbers based on the comma-separated list provided in the `phoneNumber` array.")
		cmd.Flags().StringVar(&ownerName, "owner-name", "", "Return the list of phone numbers that are owned by the given `ownerName`. Maximum length is 255.")
		cmd.Flags().StringVar(&extension, "extension", "", "Returns the list of PSTN phone numbers with the given `extension`.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-ms-teams-settings
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-ms-teams-settings",
			Short: "Retrieve a Person's MS Teams Settings",
			Long:  "<div><Callout type=\"warning\">Not supported for Webex for Government (FedRAMP)</Callout></div>\n\nRetrieve a person's MS Teams settings.\n\nAt a person level, MS Teams settings allow access to retrieving the `HIDE WEBEX APP` and `PRESENCE SYNC` settings.\n\nTo retrieve a person's MS Teams settings requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/settings/msTeams")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter since the default is the same organization as the token used to access the API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-ms-teams-setting
		var personId string
		var orgId string
		var settingName string
		var value bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-ms-teams-setting",
			Short: "Configure a Person's MS Teams Setting",
			Long:  "<div><Callout type=\"warning\">Not supported for Webex for Government (FedRAMP)</Callout></div>\n\nConfigure a Person's MS Teams setting.\n\nMS Teams settings can be configured at the person level.\n\nRequires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/settings/msTeams")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("settingName", settingName)
					req.BodyBool("value", value, cmd.Flags().Changed("value"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter since the default is the same organization as the token used to access the API.")
		cmd.Flags().StringVar(&settingName, "setting-name", "", "")
		cmd.Flags().BoolVar(&value, "value", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-personal-assistant
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-personal-assistant",
			Short: "Get Personal Assistant",
			Long:  "Retrieve Personal Assistant details for a specific user.\n\nPersonal Assistant is used to manage a user's incoming calls when they are away.\n\nRetrieving Personal Assistant details requires a full, user, or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/features/personalAssistant")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Get Personal Assistant details for the organization.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-personal-assistant
		var personId string
		var orgId string
		var enabled bool
		var presence string
		var untilDateTime string
		var transferEnabled bool
		var transferNumber string
		var alerting string
		var alertMeFirstNumberOfRings int64
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-personal-assistant",
			Short: "Update Personal Assistant",
			Long:  "Update Personal Assistant details for a specific user.\n\nPersonal Assistant is used to manage a user's incoming calls when they are away.\n\nUpdating Personal Assistant details requires a full or user administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/features/personalAssistant")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyString("presence", presence)
					req.BodyString("untilDateTime", untilDateTime)
					req.BodyBool("transferEnabled", transferEnabled, cmd.Flags().Changed("transfer-enabled"))
					req.BodyString("transferNumber", transferNumber)
					req.BodyString("alerting", alerting)
					req.BodyInt("alertMeFirstNumberOfRings", alertMeFirstNumberOfRings, cmd.Flags().Changed("alert-me-first-number-of-rings"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update Personal Assistant details for the organization.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&presence, "presence", "", "")
		cmd.Flags().StringVar(&untilDateTime, "until-date-time", "", "")
		cmd.Flags().BoolVar(&transferEnabled, "transfer-enabled", false, "")
		cmd.Flags().StringVar(&transferNumber, "transfer-number", "", "")
		cmd.Flags().StringVar(&alerting, "alerting", "", "")
		cmd.Flags().Int64Var(&alertMeFirstNumberOfRings, "alert-me-first-number-of-rings", 0, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // list-available-features
		var personId string
		var name string
		var phoneNumber string
		var extension string
		var max string
		var start string
		var order string
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-available-features",
			Short: "Retrieve the List of Available Features",
			Long:  "Retrieve a list of feature identifiers that can be assigned to a user for `Mode Management`. Feature identifiers reference feature instances like `Auto Attendants`, `Call Queues`, and `Hunt Groups`.\n\nFeatures with mode-based call forwarding enabled can be assigned to a user for `Mode Management`.\n\nRetrieving this list requires a full, read-only, or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/modeManagement/availableFeatures")
				req.PathParam("personId", personId)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("extension", extension)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("order", order)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the user.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&name, "name", "", "List features whose `name` contains this string.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "List features whose phoneNumber contains this matching string.")
		cmd.Flags().StringVar(&extension, "extension", "", "List features whose `extension` contains this matching string.")
		cmd.Flags().StringVar(&max, "max", "", "Maximum number of features to return in a single page.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&order, "order", "", "Sort the list of features based on `name`, `phoneNumber`, or `extension`, either `asc`, or `desc`.")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve features list from this organization.")
		userCallCmd.AddCommand(cmd)
	}

	{ // list-features-assigned-mode-management
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-features-assigned-mode-management",
			Short: "Retrieve the List of Features Assigned to a User for Mode Management",
			Long:  "Retrieve a list of feature identifiers that are already assigned to a user for `Mode Management`. Feature identifiers reference feature instances like `Auto Attendants`, `Call Queues`, and `Hunt Groups`.\nA maximum of 50 features can be assigned to a user for `Mode Management`.\n\nFeatures with mode-based call forwarding enabled can be assigned to a user for `Mode Management`.\n\nRetrieving this list requires a full, read-only, or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/modeManagement/features")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the user.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve features list from this organization.")
		userCallCmd.AddCommand(cmd)
	}

	{ // assign-list-features-mode-management
		var personId string
		var orgId string
		var featureIds []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "assign-list-features-mode-management",
			Short: "Assign a List of Features to a User for Mode Management",
			Long:  "Assign a user a list of feature identifiers for `Mode Management`. Feature identifiers reference feature instances like `Auto Attendants`, `Call Queues`, and `Hunt Groups`.\nA maximum of 50 features can be assigned to a user for `Mode Management`.\n\nUpdating mode management settings for a user requires a full, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/modeManagement/features")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("featureIds", featureIds)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the user.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve features list from this organization.")
		cmd.Flags().StringSliceVar(&featureIds, "feature-ids", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-users-selective-accept-list
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-users-selective-accept-list",
			Short: "Get the User\u2019s Selective Call Accept Criteria List",
			Long:  "Retrieve selective call accept criteria list for a user.\n\nWith the selective call accept feature, you can create different rules to accept specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, read-only, or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/selectiveAccept")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-users-selective-accept
		var personId string
		var orgId string
		var enabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-users-selective-accept",
			Short: "Update User\u2019s Selective Call Accept Criteria",
			Long:  "Modify selective call accept setting for a user.\n\nWith the Selective Call accept feature, you can create different rules to accept specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/selectiveAccept")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // create-users-selective-accept-service
		var personId string
		var orgId string
		var callsFrom string
		var acceptEnabled bool
		var scheduleName string
		var scheduleType string
		var scheduleLevel string
		var anonymousCallersEnabled bool
		var unavailableCallersEnabled bool
		var phoneNumbers []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-users-selective-accept-service",
			Short: "Create a Criteria to the User\u2019s Selective Call Accept Service",
			Long:  "Add a criteria to the user's selective call accept service.\n\nWith the Selective Call accept feature, you can create different rules to accept specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/people/{personId}/selectiveAccept/criteria")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callsFrom", callsFrom)
					req.BodyBool("acceptEnabled", acceptEnabled, cmd.Flags().Changed("accept-enabled"))
					req.BodyString("scheduleName", scheduleName)
					req.BodyString("scheduleType", scheduleType)
					req.BodyString("scheduleLevel", scheduleLevel)
					req.BodyBool("anonymousCallersEnabled", anonymousCallersEnabled, cmd.Flags().Changed("anonymous-callers-enabled"))
					req.BodyBool("unavailableCallersEnabled", unavailableCallersEnabled, cmd.Flags().Changed("unavailable-callers-enabled"))
					req.BodyStringSlice("phoneNumbers", phoneNumbers)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		cmd.Flags().StringVar(&callsFrom, "calls-from", "", "")
		cmd.Flags().BoolVar(&acceptEnabled, "accept-enabled", false, "")
		cmd.Flags().StringVar(&scheduleName, "schedule-name", "", "")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "")
		cmd.Flags().StringVar(&scheduleLevel, "schedule-level", "", "")
		cmd.Flags().BoolVar(&anonymousCallersEnabled, "anonymous-callers-enabled", false, "")
		cmd.Flags().BoolVar(&unavailableCallersEnabled, "unavailable-callers-enabled", false, "")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-users-selective-accept-service
		var personId string
		var id string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-users-selective-accept-service",
			Short: "Get a Criteria for the User\u2019s Selective Call Accept Service",
			Long:  "Get the criteria details for the user's selective call accept service.\n\nWith the Selective Call accept feature, you can create different rules to accept specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, read-only, or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/selectiveAccept/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Criteria ID.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-users-selective-accept-service
		var personId string
		var id string
		var orgId string
		var callsFrom string
		var acceptEnabled bool
		var scheduleName string
		var scheduleType string
		var scheduleLevel string
		var anonymousCallersEnabled bool
		var unavailableCallersEnabled bool
		var phoneNumbers []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-users-selective-accept-service",
			Short: "Modify a Criteria From the User\u2019s Selective Call Accept Service",
			Long:  "Modify a criteria for the user's selective call accept service.\n\nWith the Selective Call Accept feature, you can create different rules to accept specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/selectiveAccept/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callsFrom", callsFrom)
					req.BodyBool("acceptEnabled", acceptEnabled, cmd.Flags().Changed("accept-enabled"))
					req.BodyString("scheduleName", scheduleName)
					req.BodyString("scheduleType", scheduleType)
					req.BodyString("scheduleLevel", scheduleLevel)
					req.BodyBool("anonymousCallersEnabled", anonymousCallersEnabled, cmd.Flags().Changed("anonymous-callers-enabled"))
					req.BodyBool("unavailableCallersEnabled", unavailableCallersEnabled, cmd.Flags().Changed("unavailable-callers-enabled"))
					req.BodyStringSlice("phoneNumbers", phoneNumbers)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Criteria ID.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		cmd.Flags().StringVar(&callsFrom, "calls-from", "", "")
		cmd.Flags().BoolVar(&acceptEnabled, "accept-enabled", false, "")
		cmd.Flags().StringVar(&scheduleName, "schedule-name", "", "")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "")
		cmd.Flags().StringVar(&scheduleLevel, "schedule-level", "", "")
		cmd.Flags().BoolVar(&anonymousCallersEnabled, "anonymous-callers-enabled", false, "")
		cmd.Flags().BoolVar(&unavailableCallersEnabled, "unavailable-callers-enabled", false, "")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // delete-users-selective-accept-service
		var personId string
		var id string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-users-selective-accept-service",
			Short: "Delete a Criteria From the User\u2019s Selective Call Accept service",
			Long:  "Delete a criteria from the user's selective call accept criteria list.\n\nWith the Selective Call Accept feature, you can create different rules to accept specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/people/{personId}/selectiveAccept/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Criteria ID.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-users-selective-reject-listing
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-users-selective-reject-listing",
			Short: "Get the User\u2019s Selective Call Rejection Criteria Listing",
			Long:  "Retrieve selective call rejection criteria for a user.\n\nWith the Selective Call Rejection feature, you can create different rules to reject specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, read-only, or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/selectiveReject")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-users-selective-reject-list
		var personId string
		var orgId string
		var enabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-users-selective-reject-list",
			Short: "Update User\u2019s Selective Call Rejection Criteria List",
			Long:  "Modify selective call rejection setting for a user.\n\nWith the Selective Call Rejection feature, you can create different rules to reject specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/selectiveReject")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // create-users-selective-reject-service
		var personId string
		var orgId string
		var callsFrom string
		var rejectEnabled bool
		var scheduleName string
		var scheduleType string
		var scheduleLevel string
		var anonymousCallersEnabled bool
		var unavailableCallersEnabled bool
		var phoneNumbers []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-users-selective-reject-service",
			Short: "Create a Criteria to the User\u2019s Selective Call Rejection Service",
			Long:  "Add a criteria to the user's selective call rejection service.\n\nWith the Selective Call Rejection feature, you can create different rules to reject specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/people/{personId}/selectiveReject/criteria")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callsFrom", callsFrom)
					req.BodyBool("rejectEnabled", rejectEnabled, cmd.Flags().Changed("reject-enabled"))
					req.BodyString("scheduleName", scheduleName)
					req.BodyString("scheduleType", scheduleType)
					req.BodyString("scheduleLevel", scheduleLevel)
					req.BodyBool("anonymousCallersEnabled", anonymousCallersEnabled, cmd.Flags().Changed("anonymous-callers-enabled"))
					req.BodyBool("unavailableCallersEnabled", unavailableCallersEnabled, cmd.Flags().Changed("unavailable-callers-enabled"))
					req.BodyStringSlice("phoneNumbers", phoneNumbers)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		cmd.Flags().StringVar(&callsFrom, "calls-from", "", "")
		cmd.Flags().BoolVar(&rejectEnabled, "reject-enabled", false, "")
		cmd.Flags().StringVar(&scheduleName, "schedule-name", "", "")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "")
		cmd.Flags().StringVar(&scheduleLevel, "schedule-level", "", "")
		cmd.Flags().BoolVar(&anonymousCallersEnabled, "anonymous-callers-enabled", false, "")
		cmd.Flags().BoolVar(&unavailableCallersEnabled, "unavailable-callers-enabled", false, "")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-users-selective-reject-service
		var personId string
		var id string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-users-selective-reject-service",
			Short: "Get a Criteria for the User\u2019s Selective Call Rejection Service",
			Long:  "Get a criteria for the user's selective call rejection service.\n\nWith the Selective Call Rejection feature, you can create different rules to reject specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, read-only, or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/selectiveReject/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Criteria ID.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-users-selective-reject-service
		var personId string
		var id string
		var orgId string
		var callsFrom string
		var rejectEnabled bool
		var scheduleName string
		var scheduleType string
		var scheduleLevel string
		var anonymousCallersEnabled bool
		var unavailableCallersEnabled bool
		var phoneNumbers []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-users-selective-reject-service",
			Short: "Modify a Criteria for the User\u2019s Selective Call Rejection Service",
			Long:  "Modify a criteria for the user's selective call rejection service.\n\nWith the Selective Call Rejection feature, you can create different rules to reject specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/selectiveReject/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("callsFrom", callsFrom)
					req.BodyBool("rejectEnabled", rejectEnabled, cmd.Flags().Changed("reject-enabled"))
					req.BodyString("scheduleName", scheduleName)
					req.BodyString("scheduleType", scheduleType)
					req.BodyString("scheduleLevel", scheduleLevel)
					req.BodyBool("anonymousCallersEnabled", anonymousCallersEnabled, cmd.Flags().Changed("anonymous-callers-enabled"))
					req.BodyBool("unavailableCallersEnabled", unavailableCallersEnabled, cmd.Flags().Changed("unavailable-callers-enabled"))
					req.BodyStringSlice("phoneNumbers", phoneNumbers)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Criteria ID.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		cmd.Flags().StringVar(&callsFrom, "calls-from", "", "")
		cmd.Flags().BoolVar(&rejectEnabled, "reject-enabled", false, "")
		cmd.Flags().StringVar(&scheduleName, "schedule-name", "", "")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "")
		cmd.Flags().StringVar(&scheduleLevel, "schedule-level", "", "")
		cmd.Flags().BoolVar(&anonymousCallersEnabled, "anonymous-callers-enabled", false, "")
		cmd.Flags().BoolVar(&unavailableCallersEnabled, "unavailable-callers-enabled", false, "")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // delete-users-selective-reject-service
		var personId string
		var id string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-users-selective-reject-service",
			Short: "Delete a Criteria From the User\u2019s Selective Call Rejection Service",
			Long:  "Delete a criteria from the user's selective call rejection service.\n\nWith the Selective Call Rejection feature, you can create different rules to reject specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/people/{personId}/selectiveReject/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Criteria ID.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-users-selective-forward
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-users-selective-forward",
			Short: "Get the User\u2019s Selective Call Forwarding",
			Long:  "Retrieve selective call forwarding criteria for a user.\n\nWith the Selective Call Forwarding feature, you can create different rules to forward specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, read-only, or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/selectiveForward")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-users-selective-forward-list
		var personId string
		var orgId string
		var enabled bool
		var defaultPhoneNumberToForward string
		var ringReminderEnabled bool
		var destinationVoicemailEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-users-selective-forward-list",
			Short: "Update User\u2019s Selective Call Forwarding Criteria List",
			Long:  "Modify selective call forwarding setting for a user.\n\nWith the Selective Call Forwarding feature, you can create different rules to forward specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/selectiveForward")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyString("defaultPhoneNumberToForward", defaultPhoneNumberToForward)
					req.BodyBool("ringReminderEnabled", ringReminderEnabled, cmd.Flags().Changed("ring-reminder-enabled"))
					req.BodyBool("destinationVoicemailEnabled", destinationVoicemailEnabled, cmd.Flags().Changed("destination-voicemail-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&defaultPhoneNumberToForward, "default-phone-number-to-forward", "", "")
		cmd.Flags().BoolVar(&ringReminderEnabled, "ring-reminder-enabled", false, "")
		cmd.Flags().BoolVar(&destinationVoicemailEnabled, "destination-voicemail-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // create-users-selective-forward-service
		var personId string
		var orgId string
		var forwardToPhoneNumber string
		var sendToVoicemailEnabled bool
		var callsFrom string
		var scheduleName string
		var scheduleType string
		var scheduleLevel string
		var anonymousCallersEnabled bool
		var unavailableCallersEnabled bool
		var phoneNumbers []string
		var forwardEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-users-selective-forward-service",
			Short: "Create a Criteria to the User\u2019s Selective Call Forwarding Service",
			Long:  "Add a criteria to the user's selective call forwarding service.\n\nWith the Selective Call Forwarding feature, you can create different rules to forward specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/people/{personId}/selectiveForward/criteria")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("forwardToPhoneNumber", forwardToPhoneNumber)
					req.BodyBool("sendToVoicemailEnabled", sendToVoicemailEnabled, cmd.Flags().Changed("send-to-voicemail-enabled"))
					req.BodyString("callsFrom", callsFrom)
					req.BodyString("scheduleName", scheduleName)
					req.BodyString("scheduleType", scheduleType)
					req.BodyString("scheduleLevel", scheduleLevel)
					req.BodyBool("anonymousCallersEnabled", anonymousCallersEnabled, cmd.Flags().Changed("anonymous-callers-enabled"))
					req.BodyBool("unavailableCallersEnabled", unavailableCallersEnabled, cmd.Flags().Changed("unavailable-callers-enabled"))
					req.BodyStringSlice("phoneNumbers", phoneNumbers)
					req.BodyBool("forwardEnabled", forwardEnabled, cmd.Flags().Changed("forward-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		cmd.Flags().StringVar(&forwardToPhoneNumber, "forward-to-phone-number", "", "")
		cmd.Flags().BoolVar(&sendToVoicemailEnabled, "send-to-voicemail-enabled", false, "")
		cmd.Flags().StringVar(&callsFrom, "calls-from", "", "")
		cmd.Flags().StringVar(&scheduleName, "schedule-name", "", "")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "")
		cmd.Flags().StringVar(&scheduleLevel, "schedule-level", "", "")
		cmd.Flags().BoolVar(&anonymousCallersEnabled, "anonymous-callers-enabled", false, "")
		cmd.Flags().BoolVar(&unavailableCallersEnabled, "unavailable-callers-enabled", false, "")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().BoolVar(&forwardEnabled, "forward-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-users-selective-forward-service
		var personId string
		var id string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-users-selective-forward-service",
			Short: "Get a Criteria for the User\u2019s Selective Call Forwarding Service",
			Long:  "Get the criteria details for the user's selective call forwarding service.\n\nWith the Selective Call Forwarding feature, you can create different rules to forward specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, read-only, or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/selectiveForward/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Criteria ID.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-users-selective-forward-service
		var personId string
		var id string
		var orgId string
		var forwardToPhoneNumber string
		var sendToVoicemailEnabled bool
		var callsFrom string
		var scheduleName string
		var scheduleType string
		var scheduleLevel string
		var anonymousCallersEnabled bool
		var unavailableCallersEnabled bool
		var phoneNumbers []string
		var forwardEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-users-selective-forward-service",
			Short: "Modify a Criteria for the User\u2019s Selective Call Forwarding Service",
			Long:  "Modify a criteria for the user's selective call forwarding service.\n\nWith the Selective Call Forwarding feature, you can create different rules to forward specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/selectiveForward/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("forwardToPhoneNumber", forwardToPhoneNumber)
					req.BodyBool("sendToVoicemailEnabled", sendToVoicemailEnabled, cmd.Flags().Changed("send-to-voicemail-enabled"))
					req.BodyString("callsFrom", callsFrom)
					req.BodyString("scheduleName", scheduleName)
					req.BodyString("scheduleType", scheduleType)
					req.BodyString("scheduleLevel", scheduleLevel)
					req.BodyBool("anonymousCallersEnabled", anonymousCallersEnabled, cmd.Flags().Changed("anonymous-callers-enabled"))
					req.BodyBool("unavailableCallersEnabled", unavailableCallersEnabled, cmd.Flags().Changed("unavailable-callers-enabled"))
					req.BodyStringSlice("phoneNumbers", phoneNumbers)
					req.BodyBool("forwardEnabled", forwardEnabled, cmd.Flags().Changed("forward-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Criteria ID.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		cmd.Flags().StringVar(&forwardToPhoneNumber, "forward-to-phone-number", "", "")
		cmd.Flags().BoolVar(&sendToVoicemailEnabled, "send-to-voicemail-enabled", false, "")
		cmd.Flags().StringVar(&callsFrom, "calls-from", "", "")
		cmd.Flags().StringVar(&scheduleName, "schedule-name", "", "")
		cmd.Flags().StringVar(&scheduleType, "schedule-type", "", "")
		cmd.Flags().StringVar(&scheduleLevel, "schedule-level", "", "")
		cmd.Flags().BoolVar(&anonymousCallersEnabled, "anonymous-callers-enabled", false, "")
		cmd.Flags().BoolVar(&unavailableCallersEnabled, "unavailable-callers-enabled", false, "")
		cmd.Flags().StringSliceVar(&phoneNumbers, "phone-numbers", nil, "")
		cmd.Flags().BoolVar(&forwardEnabled, "forward-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // delete-users-selective-forward-service
		var personId string
		var id string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-users-selective-forward-service",
			Short: "Delete a Criteria From the User\u2019s Selective Call Forwarding Service",
			Long:  "Delete a criteria from the user's selective call forwarding service.\n\nWith the Selective Call Forwarding feature, you can create different rules to forward specific calls based on the phone number, who's calling, and/or the time and day of the call.\n\nRequires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/people/{personId}/selectiveForward/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "Criteria ID.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which the user resides.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-application-services-settings-2
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-application-services-settings-2",
			Short: "Retrieve a person's Application Services Settings New",
			Long:  "Gets mobile and PC applications settings for a user.\n\nApplication services let you determine the ringing behavior for calls made to people in certain scenarios. You can also specify which devices can download the Webex Calling app.\n\nRequires a full, user, or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/people/{personId}/features/applications")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization in which the person resides. Only admin users of another organization (such as partners) may use this parameter as the default is the same organization as the token used to access API.")
		userCallCmd.AddCommand(cmd)
	}

	{ // search-shared-line-appearance-members-2
		var personId string
		var max string
		var start string
		var order string
		var location string
		var name string
		var phoneNumber string
		var extension string
		cmd := &cobra.Command{
			Use:   "search-shared-line-appearance-members-2",
			Short: "Search Shared-Line Appearance Members New",
			Long:  "Get members available for shared-line assignment to a Webex Calling Apps.\n\nLike most hardware devices, applications support assigning additional shared lines which can monitored and utilized by the application.\n\nThis API requires a full, user, or location administrator auth token with the `spark-admin:telephony_config_read` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/applications/availableMembers")
				req.PathParam("personId", personId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("order", order)
				req.QueryParam("location", location)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
				req.QueryParam("extension", extension)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&order, "order", "", "Order the Route Lists according to number, ascending or descending.")
		cmd.Flags().StringVar(&location, "location", "", "Location ID for the user.")
		cmd.Flags().StringVar(&name, "name", "", "Search for users with names that match the query.")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Search for users with numbers that match the query.")
		cmd.Flags().StringVar(&extension, "extension", "", "Search for users with extensions that match the query.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-shared-line-appearance-members-2
		var personId string
		cmd := &cobra.Command{
			Use:   "get-shared-line-appearance-members-2",
			Short: "Get Shared-Line Appearance Members New",
			Long:  "Get primary and secondary members assigned to a shared line on a Webex Calling Apps.\n\nLike most hardware devices, applications support assigning additional shared lines which can monitored and utilized by the application.\n\nThis API requires a full, user, or location administrator auth token with the `spark-admin:telephony_config_read` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/applications/members")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-shared-line-appearance-members-2
		var personId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-shared-line-appearance-members-2",
			Short: "Put Shared-Line Appearance Members New",
			Long:  "Add or modify primary and secondary users assigned to shared-lines on a Webex Calling Apps.\n\nLike most hardware devices, applications support assigning additional shared lines which can monitored and utilized by the application.\n\nThis API requires a full, user, or location administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/applications/members")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-captions-settings
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-captions-settings",
			Short: "Get the user call captions settings",
			Long:  "Retrieve the user's call captions settings.\n\n**NOTE**: The call captions feature is not supported for Webex Calling Standard users or users assigned to locations in India.\n\nThe call caption feature allows the customer to enable and manage closed captions and transcript functionality (rolling caption panel) in Webex Calling, without requiring the user to escalate the call to a meeting.\n\nThis API requires a full, user, read-only, or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/callCaptions")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-captions-settings
		var personId string
		var orgId string
		var userClosedCaptionsEnabled bool
		var userTranscriptsEnabled bool
		var useLocationSettingsEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-captions-settings",
			Short: "Update the user call captions settings",
			Long:  "Update the user's call captions settings.\n\n**NOTE**: The call captions feature is not supported for Webex Calling Standard users or users assigned to locations in India.\n\nThe call caption feature allows the customer to enable and manage closed captions and transcript functionality (rolling caption panel) in Webex Calling, without requiring the user to escalate the call to a meeting.\n\nThis API requires a full, user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/callCaptions")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("userClosedCaptionsEnabled", userClosedCaptionsEnabled, cmd.Flags().Changed("user-closed-captions-enabled"))
					req.BodyBool("userTranscriptsEnabled", userTranscriptsEnabled, cmd.Flags().Changed("user-transcripts-enabled"))
					req.BodyBool("useLocationSettingsEnabled", useLocationSettingsEnabled, cmd.Flags().Changed("use-location-settings-enabled"))
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Unique identifier for the organization.")
		cmd.Flags().BoolVar(&userClosedCaptionsEnabled, "user-closed-captions-enabled", false, "")
		cmd.Flags().BoolVar(&userTranscriptsEnabled, "user-transcripts-enabled", false, "")
		cmd.Flags().BoolVar(&useLocationSettingsEnabled, "use-location-settings-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-executive-filtering-settings
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-executive-filtering-settings",
			Short: "Get Person Executive Call Filtering Settings",
			Long:  "Retrieve the executive call filtering settings for the specified person.\n\nExecutive Call Filtering in Webex allows you to control which calls are allowed to reach the executive assistant based on custom criteria, such as specific phone numbers or call types. You can enable or disable call filtering and configure filter rules to manage incoming calls.\n\nThis API requires a full, user, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/executive/callFiltering")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the user.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-executive-filtering-settings
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-executive-filtering-settings",
			Short: "Modify Person Executive Call Filtering Settings",
			Long:  "Update the executive call filtering settings for the specified person.\n\nExecutive Call Filtering in Webex allows you to control which calls are allowed to reach the executive assistant based on custom criteria, such as specific phone numbers or call types. You can enable or disable call filtering and configure filter rules to manage incoming calls.\n\nThis API requires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/executive/callFiltering")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the user.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-executive-filtering-settings-2
		var personId string
		var id string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-executive-filtering-settings-2",
			Short: "Get Person Executive Call Filtering Criteria Settings",
			Long:  "Retrieve the executive call filtering criteria settings for the specified person.\n\nExecutive Call Filtering Criteria in Webex allows you to retrieve the detailed configuration for a specific filter rule. This includes schedule settings, phone number filters, and call routing preferences for executive call filtering.\n\nThis API requires a full, user, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/executive/callFiltering/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
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
		cmd.Flags().StringVar(&id, "id", "", "The `id` parameter specifies the unique identifier for the executive call filtering criteria.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the user.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-executive-filtering-settings-2
		var personId string
		var id string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-executive-filtering-settings-2",
			Short: "Modify Person Executive Call Filtering Criteria Settings",
			Long:  "Update the executive call filtering settings for the specified person.\n\nExecutive Call Filtering in Webex allows you to control which calls are allowed to reach the executive assistant based on custom criteria, such as specific phone numbers or call types. You can enable or disable call filtering and configure filter rules to manage incoming calls.\n\nThis API requires a full, user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/executive/callFiltering/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "The `id` parameter specifies the unique identifier for the executive call filtering criteria. Example: `Y2lzY29zcGFyazovL3VzL0NSSVRFUklBL2RHVnpkRjltYVd4MFpYST0`.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the user.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // delete-person-executive-filtering
		var personId string
		var id string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-person-executive-filtering",
			Short: "Delete Person Executive Call Filtering Criteria",
			Long:  "Delete a specific executive call filtering criteria configuration for the specified person.\n\nExecutive Call Filtering Criteria in Webex allows you to manage detailed filter rules for incoming calls. This API removes a specific filter rule by its unique identifier.\n\nThis API requires a full, user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/people/{personId}/executive/callFiltering/criteria/{id}")
				req.PathParam("personId", personId)
				req.PathParam("id", id)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&id, "id", "", "The `id` parameter specifies the unique identifier for the executive call filtering criteria. Example: `Y2lzY29zcGFyazovL3VzL0NSSVRFUklBL2RHVnpkRjltYVd4MFpYST0`.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the user.")
		userCallCmd.AddCommand(cmd)
	}

	{ // create-person-executive-filtering
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-person-executive-filtering",
			Short: "Add Person Executive Call Filtering Criteria",
			Long:  "Create a new executive call filtering criteria configuration for the specified person.\n\nExecutive Call Filtering Criteria in Webex allows you to define detailed filter rules for incoming calls. This API creates a new filter rule with the specified configuration, including schedule, phone numbers, and call routing preferences.\n\nThis API requires a full, user or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/people/{personId}/executive/callFiltering/criteria")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the user.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-executive-alert-settings
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-executive-alert-settings",
			Short: "Get Person Executive Alert Settings",
			Long:  "Get executive alert settings for the specified person.\n\nExecutive Alert settings in Webex allow you to control how calls are routed to executive assistants, including alerting mode, rollover options, and caller ID presentation. You can configure settings such as sequential or simultaneous alerting, and specify what happens when calls aren't answered.\n\nThis API requires a full, user, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/executive/alert")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the person.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-executive-alert-settings
		var personId string
		var orgId string
		var alertingMode string
		var nextAssistantNumberOfRings int64
		var rolloverEnabled bool
		var rolloverAction string
		var rolloverForwardToPhoneNumber string
		var rolloverWaitTimeInSecs int64
		var clidNameMode string
		var customClidname string
		var customClidnameInUnicode string
		var clidPhoneNumberMode string
		var customClidphoneNumber string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-executive-alert-settings",
			Short: "Modify Person Executive Alert Settings",
			Long:  "Update executive alert settings for the specified person.\n\nExecutive Alert settings in Webex allow you to control how calls are routed to executive assistants, including alerting mode, rollover options, and caller ID presentation. You can configure settings such as sequential or simultaneous alerting, and specify what happens when calls aren't answered.\n\nThis API requires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/executive/alert")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("alertingMode", alertingMode)
					req.BodyInt("nextAssistantNumberOfRings", nextAssistantNumberOfRings, cmd.Flags().Changed("next-assistant-number-of-rings"))
					req.BodyBool("rolloverEnabled", rolloverEnabled, cmd.Flags().Changed("rollover-enabled"))
					req.BodyString("rolloverAction", rolloverAction)
					req.BodyString("rolloverForwardToPhoneNumber", rolloverForwardToPhoneNumber)
					req.BodyInt("rolloverWaitTimeInSecs", rolloverWaitTimeInSecs, cmd.Flags().Changed("rollover-wait-time-in-secs"))
					req.BodyString("clidNameMode", clidNameMode)
					req.BodyString("customCLIDName", customClidname)
					req.BodyString("customCLIDNameInUnicode", customClidnameInUnicode)
					req.BodyString("clidPhoneNumberMode", clidPhoneNumberMode)
					req.BodyString("customCLIDPhoneNumber", customClidphoneNumber)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the person.")
		cmd.Flags().StringVar(&alertingMode, "alerting-mode", "", "")
		cmd.Flags().Int64Var(&nextAssistantNumberOfRings, "next-assistant-number-of-rings", 0, "")
		cmd.Flags().BoolVar(&rolloverEnabled, "rollover-enabled", false, "")
		cmd.Flags().StringVar(&rolloverAction, "rollover-action", "", "")
		cmd.Flags().StringVar(&rolloverForwardToPhoneNumber, "rollover-forward-to-phone-number", "", "")
		cmd.Flags().Int64Var(&rolloverWaitTimeInSecs, "rollover-wait-time-in-secs", 0, "")
		cmd.Flags().StringVar(&clidNameMode, "clid-name-mode", "", "")
		cmd.Flags().StringVar(&customClidname, "custom-clidname", "", "")
		cmd.Flags().StringVar(&customClidnameInUnicode, "custom-clidname-in-unicode", "", "")
		cmd.Flags().StringVar(&clidPhoneNumberMode, "clid-phone-number-mode", "", "")
		cmd.Flags().StringVar(&customClidphoneNumber, "custom-clidphone-number", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-executive-assigned-assistants
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-executive-assigned-assistants",
			Short: "Get Person Executive Assigned Assistants",
			Long:  "Get list of assigned executive assistants for the specified person.\n\nAs an executive, you can add assistants to your executive pool to manage calls for you. You can set when and which types of calls they can handle. Assistants can opt in when needed or opt out when not required.\n\nThis API requires a full, user, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/executive/assignedAssistants")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the person.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-executive-assigned-assistants
		var personId string
		var orgId string
		var assistantIds []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-executive-assigned-assistants",
			Short: "Modify Person Executive Assigned Assistants",
			Long:  "Update assigned executive assistants for the specified person.\n\nAs an executive, you can add assistants to your executive pool to manage calls for you. You can set when and which types of calls they can handle. Assistants can opt in when needed or opt out when not required.\n\nThis API requires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/executive/assignedAssistants")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("assistantIds", assistantIds)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the person.")
		cmd.Flags().StringSliceVar(&assistantIds, "assistant-ids", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-executive-available-assistants
		var personId string
		var orgId string
		var max string
		var start string
		var name string
		var phoneNumber string
		cmd := &cobra.Command{
			Use:   "get-person-executive-available-assistants",
			Short: "Get Person Executive Available Assistants",
			Long:  "Retrieves a list of people available for assignment as executive assistants to the specified person.\n\nAs an executive, you can add assistants to your executive pool to manage calls for you. You can set when and which types of calls they can handle. Assistants can opt in when needed or opt out when not required.\n\nThis API requires a full, user, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/executive/availableAssistants")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("name", name)
				req.QueryParam("phoneNumber", phoneNumber)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the person.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&name, "name", "", "Only return people with the matching name (person's first and last name combination).")
		cmd.Flags().StringVar(&phoneNumber, "phone-number", "", "Only return people with the matching phone number or extension.")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-executive-assistant-settings
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-executive-assistant-settings",
			Short: "Get Person Executive Assistant Settings",
			Long:  "Get executive assistant settings for the specified person when person is configured as executive assistant.\n\nExecutive assistants can make, answer, intercept, and route calls appropriately on behalf of their executive.\nAssistants can also set the call forwarding destination, and join or leave an executive\u2019s pool.\n\nThis API requires a full, user, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/executive/assistant")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the person.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-executive-assistant-settings
		var personId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-executive-assistant-settings",
			Short: "Modify Person Executive Assistant Settings",
			Long:  "Update executive assistant settings for the specified person when person is configured as executive assistant.\n\nExecutive assistants can make, answer, intercept, and route calls appropriately on behalf of their executive.\nAssistants can also set the call forwarding destination, and join or leave an executive\u2019s pool.\n\nThis API requires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/executive/assistant")
				req.PathParam("personId", personId)
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the person.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

	{ // get-person-executive-screening-settings
		var personId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-person-executive-screening-settings",
			Short: "Get Person Executive Screening Settings",
			Long:  "Get executive screening settings for the specified person.\n\nExecutive Screening in Webex allows you to manage how incoming calls are screened and alerted based on your preferences. You can enable or disable executive screening and configure alert types and locations for notifications.\n\nThis API requires a full, user, read-only or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/people/{personId}/executive/screening")
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
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the person.")
		userCallCmd.AddCommand(cmd)
	}

	{ // update-person-executive-screening-settings
		var personId string
		var orgId string
		var enabled bool
		var alertType string
		var alertAnywhereLocationEnabled bool
		var alertMobilityLocationEnabled bool
		var alertSharedCallAppearanceLocationEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-person-executive-screening-settings",
			Short: "Modify Person Executive Screening Settings",
			Long:  "Update executive screening settings for the specified person.\n\nExecutive Screening in Webex allows you to manage how incoming calls are screened and alerted based on your preferences. You can enable or disable executive screening and configure alert types and locations for notifications.\n\nThis API requires a full, user, or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/people/{personId}/executive/screening")
				req.PathParam("personId", personId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyString("alertType", alertType)
					req.BodyBool("alertAnywhereLocationEnabled", alertAnywhereLocationEnabled, cmd.Flags().Changed("alert-anywhere-location-enabled"))
					req.BodyBool("alertMobilityLocationEnabled", alertMobilityLocationEnabled, cmd.Flags().Changed("alert-mobility-location-enabled"))
					req.BodyBool("alertSharedCallAppearanceLocationEnabled", alertSharedCallAppearanceLocationEnabled, cmd.Flags().Changed("alert-shared-call-appearance-location-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&personId, "person-id", "", "A unique identifier for the person.")
		cmd.MarkFlagRequired("person-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization ID for the person.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&alertType, "alert-type", "", "")
		cmd.Flags().BoolVar(&alertAnywhereLocationEnabled, "alert-anywhere-location-enabled", false, "")
		cmd.Flags().BoolVar(&alertMobilityLocationEnabled, "alert-mobility-location-enabled", false, "")
		cmd.Flags().BoolVar(&alertSharedCallAppearanceLocationEnabled, "alert-shared-call-appearance-location-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		userCallCmd.AddCommand(cmd)
	}

}
