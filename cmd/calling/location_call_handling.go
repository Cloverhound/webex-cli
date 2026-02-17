package calling

import (
	"fmt"
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
var _ = strings.Join

var locationCallHandlingCmd = &cobra.Command{
	Use:   "location-call-handling",
	Short: "LocationCallHandling commands",
}

func init() {
	cmd.CallingCmd.AddCommand(locationCallHandlingCmd)

	{ // generate-example-password
		var locationId string
		var orgId string
		var generate []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "generate-example-password",
			Short: "Generate example password for Location",
			Long:  "Generates an example password using the effective password settings for the location. If you don't specify anything in the `generate` field or don't provide a request body, then you will receive a SIP password by default.\n\nUsed while creating a trunk and shouldn't be used anywhere else.\n\nGenerating an example password requires a full or write-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/actions/generatePassword/invoke")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("generate", generate)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Location for which example password has to be generated.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the location belongs.")
		cmd.Flags().StringSliceVar(&generate, "generate", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // get-internal-dialing-configuration
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-internal-dialing-configuration",
			Short: "Read the Internal Dialing configuration for a location",
			Long:  "Get current configuration for routing unknown extensions to the Premises as internal calls\n\nIf some users in a location are registered to a PBX, retrieve the setting to route unknown extensions (digits that match the extension length) to the PBX.\n\nRetrieving the internal dialing configuration requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/internalDialing")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "location for which internal calling configuration is being requested")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List route identities for this organization.")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // update-internal-dialing-configuration
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-internal-dialing-configuration",
			Short: "Modify the Internal Dialing configuration for a location",
			Long:  "Modify current configuration for routing unknown extensions to the premise as internal calls\n\nIf some users in a location are registered to a PBX, enable the setting to route unknown extensions (digits that match the extension length) to the PBX.\n\nEditing the internal dialing configuration requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/internalDialing")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "location for which internal calling configuration is being requested")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "List route identities for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // get-intercept
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-intercept",
			Short: "Get Location Intercept",
			Long:  "Retrieve intercept location details for a customer location.\n\nIntercept incoming or outgoing calls for persons in your organization. If this is enabled, calls are either routed to a designated number the person chooses, or to the person's voicemail.\n\nRetrieving intercept location details requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/intercept")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve intercept details for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve intercept location details for a customer location.")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // update-intercept
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-intercept",
			Short: "Put Location Intercept",
			Long:  "Modifies the intercept location details for a customer location.\n\nIntercept incoming or outgoing calls for users in your organization. If this is enabled, calls are either routed to a designated number the user chooses, or to the user's voicemail.\n\nModifying the intercept location details requires a full, user administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/intercept")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Modifies the intercept details for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Modifies the intercept location details for a customer location.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // get-outgoing-permission
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-outgoing-permission",
			Short: "Get Location Outgoing Permission",
			Long: `Retrieve the location's outgoing call settings.

A location's outgoing call settings allow you to determine the types of calls the people/workspaces at the location are allowed to make, as well as configure the default calling permission for each call type at the location.

Retrieving a location's outgoing call settings requires a full, user or read-only administrator or location administrator auth token with a scope of spark-admin:telephony_config_read.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/outgoingPermission")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve outgoing call settings for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve outgoing call settings for this organization.")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // update-outgoing-permission
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-outgoing-permission",
			Short: "Update Location Outgoing Permission",
			Long: `Update the location's outgoing call settings.

Location's outgoing call settings allows you to determine the types of calls the people/workspaces at this location are allowed to make and configure the default calling permission for each call type at a location.

Updating a location's outgoing call settings requires a full administrator or location administrator auth token with a scope of spark-admin:telephony_config_write.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/outgoingPermission")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update outgoing call settings for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update outgoing call settings for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // get-outgoing-auto-transfer
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-outgoing-auto-transfer",
			Short: "Get Outgoing Permission Auto Transfer Number",
			Long:  "Get the transfer numbers for the outbound permission in a location.\n\nOutbound permissions can specify which transfer number an outbound call should transfer to via the `action` field.\n\nRetrieving an auto transfer number requires a full, user or read-only administrator or location administrator auth token with a scope of spark-admin:telephony_config_read.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/outgoingPermission/autoTransferNumbers")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve auto transfer number for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve auto transfer number for this organization.")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // update-outgoing-auto-transfer
		var locationId string
		var orgId string
		var autoTransferNumber1 string
		var autoTransferNumber2 string
		var autoTransferNumber3 string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-outgoing-auto-transfer",
			Short: "Put Outgoing Permission Auto Transfer Number",
			Long:  "Modifies the transfer numbers for the outbound permission in a location.\n\nOutbound permissions can specify which transfer number an outbound call should transfer to via the `action` field.\n\nUpdating auto transfer number requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/outgoingPermission/autoTransferNumbers")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Updating auto transfer number for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Updating auto transfer number for this organization.")
		cmd.Flags().StringVar(&autoTransferNumber1, "auto-transfer-number1", "", "")
		cmd.Flags().StringVar(&autoTransferNumber2, "auto-transfer-number2", "", "")
		cmd.Flags().StringVar(&autoTransferNumber3, "auto-transfer-number3", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // get-outgoing-permission-access-code
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-outgoing-permission-access-code",
			Short: "Get Outgoing Permission Location Access Code",
			Long:  "Retrieve access codes details for a customer location.\n\nUse Access Codes to bypass the set permissions for all persons/workspaces at this location.\n\nRetrieving access codes details requires a full, user or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/outgoingPermission/accessCodes")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve access codes details for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve access codes details for a customer location.")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // create-outgoing-permission-access-code-customer
		var locationId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-outgoing-permission-access-code-customer",
			Short: "Create Outgoing Permission a new access code for a customer location",
			Long: `Add a new access code for the given location for a customer.

Use Access Codes to bypass the set permissions for all persons/workspaces at this location.

Creating an access code for the given location requires a full or user administrator or location administrator auth token with a scope of spark-admin:telephony_config_write.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/outgoingPermission/accessCodes")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Add new access code for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Add new access code for this organization.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // delete-outgoing-access-code
		var locationId string
		var orgId string
		var deleteCodes []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "delete-outgoing-access-code",
			Short: "Delete Outgoing Permission Access Code Location",
			Long:  "Deletes the access code details for a particular location for a customer.\n\nUse Access Codes to bypass the set permissions for all persons/workspaces at this location.\n\nModifying the access code location details requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/outgoingPermission/accessCodes")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Deletes the access code details for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Deletes the access code details for a customer location.")
		cmd.Flags().StringSliceVar(&deleteCodes, "delete-codes", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // delete-all-outgoing-access-code
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-all-outgoing-access-code",
			Short: "Delete all Outgoing Permission Access Code for a Location",
			Long:  "Deletes all the access codes for a particular location for a customer.\n\nUse Access Codes to bypass the set permissions for all persons/workspaces at this location.\n\nDeleting the access codes requires a full or user administrator or location administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/outgoingPermission/accessCodes")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Deletes all the access codes for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Deletes the access codes for a customer location.")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // get-outgoing-digit-pattern
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-outgoing-digit-pattern",
			Short: "Get Outgoing Permission Digit Pattern for a Location",
			Long:  "Get the digit patterns for the outbound permission in a location.\n\nUse Digit Patterns to bypass the set permissions for all persons/workspaces at this location.\n\nRetrieving digit patterns requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/outgoingPermission/digitPatterns")
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve the digit patterns for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve the digit patterns for this organization.")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // create-outgoing-permission-digit-pattern
		var locationId string
		var orgId string
		var name string
		var pattern string
		var action string
		var transferEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-outgoing-permission-digit-pattern",
			Short: "Create Outgoing Permission a new Digit Pattern for a location",
			Long:  "Add a new digit pattern for the given location for a customer.\n\nUse Digit Patterns to bypass the set permissions for all persons/workspaces at this location.\n\nCreating a digit pattern requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/outgoingPermission/digitPatterns")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Add a new digit pattern for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Add a new digit pattern for this organization.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&pattern, "pattern", "", "")
		cmd.Flags().StringVar(&action, "action", "", "")
		cmd.Flags().BoolVar(&transferEnabled, "transfer-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // delete-all-outgoing-digit-patterns
		var locationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-all-outgoing-digit-patterns",
			Short: "Delete all Outgoing Permission Digit Patterns for a Location",
			Long:  "Deletes all the digit patterns for a particular location for a customer.\n\nUse Digit Patterns to bypass the set permissions for all persons/workspaces at this location.\n\nDeleting the digit patterns requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/outgoingPermission/digitPatterns")
				req.PathParam("locationId", locationId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Delete the digit patterns for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the digit patterns for this organization.")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // get-outgoing-digit-pattern-2
		var locationId string
		var digitPatternId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-outgoing-digit-pattern-2",
			Short: "Get Details for a Outgoing Permission Digit Pattern for a Location",
			Long:  "Get the digit pattern details.\n\nUse Digit Patterns to bypass the set permissions for all persons/workspaces at this location.\n\nRetrieving digit pattern details requires a full or read-only administrator or location administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/outgoingPermission/digitPatterns/{digitPatternId}")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Retrieve the digit pattern details for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&digitPatternId, "digit-pattern-id", "", "Retrieve the digit pattern with the matching ID.")
		cmd.MarkFlagRequired("digit-pattern-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Retrieve the digit pattern details for this organization.")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // update-outgoing-digit-pattern
		var locationId string
		var digitPatternId string
		var orgId string
		var name string
		var pattern string
		var action string
		var transferEnabled bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-outgoing-digit-pattern",
			Short: "Update a Outgoing Permission Digit Pattern for a Location",
			Long:  "Update the designated digit pattern.\n\nUse Digit Patterns to bypass the set permissions for all persons/workspaces at this location.\n\nUpdating a digit pattern requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/outgoingPermission/digitPatterns/{digitPatternId}")
				req.PathParam("locationId", locationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Update the digit pattern for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&digitPatternId, "digit-pattern-id", "", "Update the digit pattern with the matching ID.")
		cmd.MarkFlagRequired("digit-pattern-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Update the digit pattern for this organization.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&pattern, "pattern", "", "")
		cmd.Flags().StringVar(&action, "action", "", "")
		cmd.Flags().BoolVar(&transferEnabled, "transfer-enabled", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		locationCallHandlingCmd.AddCommand(cmd)
	}

	{ // delete-outgoing-digit-pattern
		var locationId string
		var digitPatternId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-outgoing-digit-pattern",
			Short: "Delete a Outgoing Permission Digit Pattern for a Location",
			Long:  "Delete the designated digit pattern.\n\nUse Digit Patterns to bypass the set permissions for all persons/workspaces at this location.\n\nDeleting a digit pattern requires a full administrator or location administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/outgoingPermission/digitPatterns/{digitPatternId}")
				req.PathParam("locationId", locationId)
				req.PathParam("digitPatternId", digitPatternId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Delete the digit pattern for this location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&digitPatternId, "digit-pattern-id", "", "Delete the digit pattern with the matching ID.")
		cmd.MarkFlagRequired("digit-pattern-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Delete the digit pattern for this organization.")
		locationCallHandlingCmd.AddCommand(cmd)
	}

}
