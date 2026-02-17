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

var callRecordingCmd = &cobra.Command{
	Use:   "call-recording",
	Short: "CallRecording commands",
}

func init() {
	cmd.CallingCmd.AddCommand(callRecordingCmd)

	{ // get
		var orgId string
		cmd := &cobra.Command{
			Use:   "get",
			Short: "Get Call Recording Settings",
			Long:  "Retrieve call recording settings for the organization.\n\nThe Call Recording feature enables authorized agents to record any active call that Webex Contact Center manages.\n\nRetrieving call recording settings requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/callRecording")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call recording settings from this organization.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // update
		var orgId string
		var enabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update",
			Short: "Update Call Recording Settings",
			Long:  "Update call recording settings for the organization.\n\nThe Call Recording feature enables authorized agents to record any active call that Webex Contact Center manages.\n\nUpdating call recording settings requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.\n\n**NOTE**: This API is for Cisco partners only.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/callRecording")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call recording settings from this organization.")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // get-terms-settings
		var vendorId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-terms-settings",
			Short: "Get Call Recording Terms Of Service Settings",
			Long:  "Retrieve call recording terms of service settings for the organization.\n\nThe Call Recording feature enables authorized agents to record any active call that Webex Contact Center manages.\n\nRetrieving call recording terms of service settings requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/callRecording/vendors/{vendorId}/termsOfService")
				req.PathParam("vendorId", vendorId)
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
		cmd.Flags().StringVar(&vendorId, "vendor-id", "", "Retrieve call recording terms of service details for the given vendor.")
		cmd.MarkFlagRequired("vendor-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call recording terms of service details from this organization.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // update-terms-settings
		var vendorId string
		var orgId string
		var termsOfServiceEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-terms-settings",
			Short: "Update Call Recording Terms Of Service Settings",
			Long:  "Update call recording terms of service settings for the given vendor.\n\nThe Call Recording feature enables authorized agents to record any active call that Webex Contact Center manages.\n\nUpdating call recording terms of service settings requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/callRecording/vendors/{vendorId}/termsOfService")
				req.PathParam("vendorId", vendorId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("termsOfServiceEnabled", termsOfServiceEnabled, cmd.Flags().Changed("terms-of-service-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&vendorId, "vendor-id", "", "Update call recording terms of service settings for the given vendor.")
		cmd.MarkFlagRequired("vendor-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update call recording terms of service settings from this organization.")
		cmd.Flags().BoolVar(&termsOfServiceEnabled, "terms-of-service-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // get-org-compliance-announcement
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-org-compliance-announcement",
			Short: "Get details for the organization Compliance Announcement Setting",
			Long:  "Retrieve the organization compliance announcement settings.\n\nThe Compliance Announcement feature interacts with the Call Recording feature, specifically with the playback of the start/stop announcement. When the compliance announcement is played to the PSTN party, and the PSTN party is connected to a party with call recording enabled, then the start/stop announcement is inhibited.\n\nRetrieving organization compliance announcement setting requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/callRecording/complianceAnnouncement")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve compliance announcement setting from this organization.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // update-org-compliance-announcement
		var orgId string
		var inboundPstncallsEnabled bool
		var outboundPstncallsEnabled bool
		var outboundPstncallsDelayEnabled bool
		var delayInSeconds int64
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-org-compliance-announcement",
			Short: "Update the organization Compliance Announcement",
			Long:  "Update the organization compliance announcement.\n\nThe Compliance Announcement feature interacts with the Call Recording feature, specifically with the playback of the start/stop announcement. When the compliance announcement is played to the PSTN party, and the PSTN party is connected to a party with call recording enabled, then the start/stop announcement is inhibited.\n\nUpdating the organization compliance announcement requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/callRecording/complianceAnnouncement")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("inboundPSTNCallsEnabled", inboundPstncallsEnabled, cmd.Flags().Changed("inbound-pstncalls-enabled"))
					req.BodyBool("outboundPSTNCallsEnabled", outboundPstncallsEnabled, cmd.Flags().Changed("outbound-pstncalls-enabled"))
					req.BodyBool("outboundPSTNCallsDelayEnabled", outboundPstncallsDelayEnabled, cmd.Flags().Changed("outbound-pstncalls-delay-enabled"))
					req.BodyInt("delayInSeconds", delayInSeconds, cmd.Flags().Changed("delay-in-seconds"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update the compliance announcement setting from this organization.")
		cmd.Flags().BoolVar(&inboundPstncallsEnabled, "inbound-pstncalls-enabled", false, "")
		cmd.Flags().BoolVar(&outboundPstncallsEnabled, "outbound-pstncalls-enabled", false, "")
		cmd.Flags().BoolVar(&outboundPstncallsDelayEnabled, "outbound-pstncalls-delay-enabled", false, "")
		cmd.Flags().Int64Var(&delayInSeconds, "delay-in-seconds", 0, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // get-location-compliance-announcement
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-location-compliance-announcement",
			Short: "Get details for the Location Compliance Announcement Setting",
			Long:  "Retrieve the location compliance announcement settings.\n\nThe Compliance Announcement feature interacts with the Call Recording feature, specifically with the playback of the start/stop announcement. When the compliance announcement is played to the PSTN party, and the PSTN party is connected to a party with call recording enabled, then the start/stop announcement is inhibited.\n\nRetrieving location compliance announcement setting requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/callRecording/complianceAnnouncement")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve compliance announcement settings for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve compliance announcement setting from this organization.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // update-location-compliance-announcement
		var locationId string
		var orgId string
		var inboundPstncallsEnabled bool
		var useOrgSettingsEnabled bool
		var outboundPstncallsEnabled bool
		var outboundPstncallsDelayEnabled bool
		var delayInSeconds int64
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-location-compliance-announcement",
			Short: "Update the Location Compliance Announcement",
			Long:  "Update the location compliance announcement.\n\nThe Compliance Announcement feature interacts with the Call Recording feature, specifically with the playback of the start/stop announcement. When the compliance announcement is played to the PSTN party, and the PSTN party is connected to a party with call recording enabled, then the start/stop announcement is inhibited.\n\nUpdating the location compliance announcement requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/callRecording/complianceAnnouncement")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyBool("inboundPSTNCallsEnabled", inboundPstncallsEnabled, cmd.Flags().Changed("inbound-pstncalls-enabled"))
					req.BodyBool("useOrgSettingsEnabled", useOrgSettingsEnabled, cmd.Flags().Changed("use-org-settings-enabled"))
					req.BodyBool("outboundPSTNCallsEnabled", outboundPstncallsEnabled, cmd.Flags().Changed("outbound-pstncalls-enabled"))
					req.BodyBool("outboundPSTNCallsDelayEnabled", outboundPstncallsDelayEnabled, cmd.Flags().Changed("outbound-pstncalls-delay-enabled"))
					req.BodyInt("delayInSeconds", delayInSeconds, cmd.Flags().Changed("delay-in-seconds"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update the compliance announcement settings for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update the compliance announcement setting from this organization.")
		cmd.Flags().BoolVar(&inboundPstncallsEnabled, "inbound-pstncalls-enabled", false, "")
		cmd.Flags().BoolVar(&useOrgSettingsEnabled, "use-org-settings-enabled", false, "")
		cmd.Flags().BoolVar(&outboundPstncallsEnabled, "outbound-pstncalls-enabled", false, "")
		cmd.Flags().BoolVar(&outboundPstncallsDelayEnabled, "outbound-pstncalls-delay-enabled", false, "")
		cmd.Flags().Int64Var(&delayInSeconds, "delay-in-seconds", 0, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // get-regions
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-regions",
			Short: "Get Call Recording Regions",
			Long:  "Retrieve all the call recording regions that are available for an organization.\n\nThe Call Recording feature supports multiple third-party call recording providers, or vendors, to capture and manage call recordings. An organization is configured with an overall provider, but locations can be configured to use a different vendor than the overall organization default.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/callRecording/regions")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call recording regions for this organization.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // get-vendor-users
		var orgId string
		var max string
		var start string
		var standardUserOnly string
		cmd := &cobra.Command{
			Use:   "get-vendor-users",
			Short: "Get Call Recording Vendor Users",
			Long:  "Retrieve call recording vendor users of an organization. This API is used to get the list of users who are assigned to the default call-recording vendor of the organization.\n\nThe Call Recording feature supports multiple third-party call recording providers, or vendors, to capture and manage call recordings. An organization is configured with an overall provider, but locations can be configured to use a different vendor than the overall organization default.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/callRecording/vendorUsers")
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("standardUserOnly", standardUserOnly)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call recording vendor users for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of vendor users returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects. The default is 0.")
		cmd.Flags().StringVar(&standardUserOnly, "standard-user-only", "", "If true, results only include Webex Calling standard users.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // update-vendor-location
		var locationId string
		var orgId string
		var id string
		var orgDefaultEnabled bool
		var storageRegion string
		var orgStorageRegionEnabled bool
		var failureBehavior string
		var orgFailureBehaviorEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-vendor-location",
			Short: "Set Call Recording Vendor for a Location",
			Long:  "Assign a call recording vendor to a location of an organization. Response will be `204` if the changes can be applied immediatley otherwise `200` with a job ID is returned.\n\nThe Call Recording feature supports multiple third-party call recording providers, or vendors, to capture and manage call recordings. An organization is configured with an overall provider, but locations can be configured to use a different vendor than the overall organization default.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/callRecording/vendor")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("id", id)
					req.BodyBool("orgDefaultEnabled", orgDefaultEnabled, cmd.Flags().Changed("org-default-enabled"))
					req.BodyString("storageRegion", storageRegion)
					req.BodyBool("orgStorageRegionEnabled", orgStorageRegionEnabled, cmd.Flags().Changed("org-storage-region-enabled"))
					req.BodyString("failureBehavior", failureBehavior)
					req.BodyBool("orgFailureBehaviorEnabled", orgFailureBehaviorEnabled, cmd.Flags().Changed("org-failure-behavior-enabled"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update the call recording vendor for this location")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update the call recording vendor for this organization.")
		cmd.Flags().StringVar(&id, "id", "", "")
		cmd.Flags().BoolVar(&orgDefaultEnabled, "org-default-enabled", false, "")
		cmd.Flags().StringVar(&storageRegion, "storage-region", "", "")
		cmd.Flags().BoolVar(&orgStorageRegionEnabled, "org-storage-region-enabled", false, "")
		cmd.Flags().StringVar(&failureBehavior, "failure-behavior", "", "")
		cmd.Flags().BoolVar(&orgFailureBehaviorEnabled, "org-failure-behavior-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // get-location-vendors
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-location-vendors",
			Short: "Get Location Call Recording Vendors",
			Long:  "Retrieve details of the call recording vendor that the location is assigned and also a list of vendors.\n\nThe Call Recording feature supports multiple third-party call recording providers, or vendors, to capture and manage call recordings. An organization is configured with an overall provider, but locations can be configured to use a different vendor than the overall organization default.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/callRecording/vendors")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve vendor details for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve vendor details for this organization.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // get-vendor-users-location
		var locationId string
		var orgId string
		var max string
		var start string
		var standardUserOnly string
		cmd := &cobra.Command{
			Use:   "get-vendor-users-location",
			Short: "Get Call Recording Vendor Users for a Location",
			Long:  "Retrieve call recording vendor users of a location. This API is used to get the list of users assigned to the call recording vendor of the location.\n\nThe Call Recording feature supports multiple third-party call recording providers, or vendors, to capture and manage call recordings. An organization is configured with an overall provider, but locations can be configured to use a different vendor than the overall organization default.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/callRecording/vendorUsers")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("standardUserOnly", standardUserOnly)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve vendor users for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve vendor users for this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of vendor users returned to this maximum count. The default is 2000.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects. The default is 0.")
		cmd.Flags().StringVar(&standardUserOnly, "standard-user-only", "", "If true, results only include Webex Calling standard users.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // list-jobs
		var orgId string
		var max string
		var start string
		cmd := &cobra.Command{
			Use:   "list-jobs",
			Short: "List Call Recording Jobs",
			Long:  "Get the list of all call recording jobs in an organization.\n\nThe Call Recording feature supports multiple third-party call recording providers, or vendors, to capture and manage call recordings. An organization is configured with an overall provider, but locations can be configured to use a different vendor than the overall organization default.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/jobs/callRecording")
				req.QueryParam("orgId", orgId)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List call recording jobs in this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of jobs returned to this maximum count. The default is 50.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects. The default is 0.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // get-job-status-job
		var jobId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-job-status-job",
			Short: "Get the Job Status of a Call Recording Job",
			Long:  "Get the details of a call recording job by its job ID.\n\nThe Call Recording feature supports multiple third-party call recording providers, or vendors, to capture and manage call recordings. An organization is configured with an overall provider, but locations can be configured to use a different vendor than the overall organization default.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/jobs/callRecording/{jobId}")
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
		cmd.Flags().StringVar(&jobId, "job-id", "", "Retrieve job status for this `jobId`.")
		cmd.MarkFlagRequired("job-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve job status in this organization.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // get-job-errors-job
		var jobId string
		var orgId string
		var max string
		cmd := &cobra.Command{
			Use:   "get-job-errors-job",
			Short: "Get Job Errors for a Call Recording Job",
			Long:  "Get errors for a call recording job in an organization.\n\nThe Call Recording feature supports multiple third-party call recording providers, or vendors, to capture and manage call recordings. An organization is configured with an overall provider, but locations can be configured to use a different vendor than the overall organization default.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/jobs/callRecording/{jobId}/errors")
				req.PathParam("jobId", jobId)
				req.QueryParam("orgId", orgId)
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
		cmd.Flags().StringVar(&jobId, "job-id", "", "Retrieve job errors for this job.")
		cmd.MarkFlagRequired("job-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve job errors for a call recording job in this organization.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of errors returned to this maximum count. The default is 50.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // get-org-vendors
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-org-vendors",
			Short: "Get Organization Call Recording Vendors",
			Long:  "Returns what the current vendor is as well as a list of all the available vendors.\n\nThe Call Recording feature supports multiple third-party call recording providers, or vendors, to capture and manage call recordings. An organization is configured with an overall provider, but locations can be configured to use a different vendor than the overall organization default.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/callRecording/vendors")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve call recording settings from this organization.")
		callRecordingCmd.AddCommand(cmd)
	}

	{ // update-org-vendor
		var orgId string
		var vendorId string
		var storageRegion string
		var failureBehavior string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-org-vendor",
			Short: "Set Organization Call Recording Vendor",
			Long:  "Returns a Job ID that you can use to get the status of the job.\n\nThe Call Recording feature supports multiple third-party call recording providers, or vendors, to capture and manage call recordings. An organization is configured with an overall provider, but locations can be configured to use a different vendor than the overall organization default.\n\nRequires a full administrator auth token with a scope of `spark-admin:telephony_config_write`, `spark-admin:telephony_config_read`, and `spark-admin:people_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/callRecording/vendor")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("vendorId", vendorId)
					req.BodyString("storageRegion", storageRegion)
					req.BodyString("failureBehavior", failureBehavior)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Modify call recording settings from this organization.")
		cmd.Flags().StringVar(&vendorId, "vendor-id", "", "")
		cmd.Flags().StringVar(&storageRegion, "storage-region", "", "")
		cmd.Flags().StringVar(&failureBehavior, "failure-behavior", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRecordingCmd.AddCommand(cmd)
	}

}
