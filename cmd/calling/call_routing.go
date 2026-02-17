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

var callRoutingCmd = &cobra.Command{
	Use:   "call-routing",
	Short: "CallRouting commands",
}

func init() {
	cmd.CallingCmd.AddCommand(callRoutingCmd)

	{ // test
		var orgId string
		var originatorId string
		var originatorType string
		var destination string
		var originatorNumber string
		var includeAppliedServices bool
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "test",
			Short: "Test Call Routing",
			Long:  "Validates that an incoming call can be routed.\n\nDial plans route calls to on-premises destinations by use of trunks or route groups.\nThey are configured globally for an enterprise and apply to all users, regardless of location.\nA dial plan also specifies the routing choice (trunk or route group) for calls that match any of its dial patterns.\nSpecific dial patterns can be defined as part of your dial plan.\n\nTest call routing requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/actions/testCallRouting/invoke")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("originatorId", originatorId)
					req.BodyString("originatorType", originatorType)
					req.BodyString("destination", destination)
					req.BodyString("originatorNumber", originatorNumber)
					req.BodyBool("includeAppliedServices", includeAppliedServices, cmd.Flags().Changed("include-applied-services"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization in which we are validating a call routing.")
		cmd.Flags().StringVar(&originatorId, "originator-id", "", "")
		cmd.Flags().StringVar(&originatorType, "originator-type", "", "")
		cmd.Flags().StringVar(&destination, "destination", "", "")
		cmd.Flags().StringVar(&originatorNumber, "originator-number", "", "")
		cmd.Flags().BoolVar(&includeAppliedServices, "include-applied-services", false, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-lgw-dial-plan-usage-trunk
		var trunkId string
		var orgId string
		var start string
		var max string
		var order string
		var name string
		cmd := &cobra.Command{
			Use:   "get-lgw-dial-plan-usage-trunk",
			Short: "Get Local Gateway Dial Plan Usage for a Trunk",
			Long:  "Get Local Gateway Dial Plan Usage for a Trunk.\n\nA trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group which is a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nRetrieving this information requires a full administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/trunks/{trunkId}/usageDialPlan")
				req.PathParam("trunkId", trunkId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
				req.QueryParam("max", max)
				req.QueryParam("order", order)
				req.QueryParam("name", name)
				req.QueryParam("name", name)
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
		cmd.Flags().StringVar(&trunkId, "trunk-id", "", "ID of the trunk.")
		cmd.MarkFlagRequired("trunk-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the trunk belongs.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&order, "order", "", "Order the trunks according to the designated fields.  Available sort fields are `name`, and `locationName`. Sort order is ascending by default")
		cmd.Flags().StringVar(&name, "name", "", "Return the list of trunks matching the local gateway names")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-locations-lgw-pstn-connection
		var trunkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-locations-lgw-pstn-connection",
			Short: "Get Locations Using the Local Gateway as PSTN Connection Routing",
			Long:  "Get Locations Using the Local Gateway as PSTN Connection Routing.\n\nA trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group which is a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nRetrieving this information requires a full administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/trunks/{trunkId}/usagePstnConnection")
				req.PathParam("trunkId", trunkId)
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
		cmd.Flags().StringVar(&trunkId, "trunk-id", "", "ID of the trunk.")
		cmd.MarkFlagRequired("trunk-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the trunk belongs.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-route-groups-lgw
		var trunkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-route-groups-lgw",
			Short: "Get Route Groups Using the Local Gateway",
			Long:  "Get Route Groups Using the Local Gateway.\n\nA trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group which is a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nRetrieving this information requires a full administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/trunks/{trunkId}/usageRouteGroup")
				req.PathParam("trunkId", trunkId)
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
		cmd.Flags().StringVar(&trunkId, "trunk-id", "", "ID of the trunk.")
		cmd.MarkFlagRequired("trunk-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the trunk belongs.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-lgw-usage-count
		var trunkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-lgw-usage-count",
			Short: "Get Local Gateway Usage Count",
			Long:  "Get Local Gateway Usage Count\n\nA trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group which is a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nRetrieving this information requires a full administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/trunks/{trunkId}/usage")
				req.PathParam("trunkId", trunkId)
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
		cmd.Flags().StringVar(&trunkId, "trunk-id", "", "ID of the trunk.")
		cmd.MarkFlagRequired("trunk-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the trunk belongs.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // update-dial-patterns
		var dialPlanId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-dial-patterns",
			Short: "Modify Dial Patterns",
			Long:  "Modify dial patterns for the Dial Plan.\n\nDial plans route calls to on-premises destinations by use of trunks or route groups.\nThey are configured globally for an enterprise and apply to all users, regardless of location.\nA dial plan also specifies the routing choice (trunk or route group) for calls that match any of its dial patterns.\nSpecific dial patterns can be defined as part of your dial plan.\n\nModifying a dial pattern requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/premisePstn/dialPlans/{dialPlanId}/dialPatterns")
				req.PathParam("dialPlanId", dialPlanId)
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
		cmd.Flags().StringVar(&dialPlanId, "dial-plan-id", "", "ID of the dial plan being modified.")
		cmd.MarkFlagRequired("dial-plan-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which dial plan belongs.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // validate-dial-pattern
		var orgId string
		var dialPatterns []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "validate-dial-pattern",
			Short: "Validate a Dial Pattern",
			Long:  "Validate a Dial Pattern.\n\nDial plans route calls to on-premises destinations by use of trunks or route groups.\nThey are configured globally for an enterprise and apply to all users, regardless of location.\nA dial plan also specifies the routing choice (trunk or route group) for calls that match any of its dial patterns.\nSpecific dial patterns can be defined as part of your dial plan.\n\nValidating a dial pattern requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/premisePstn/actions/validateDialPatterns/invoke")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyStringSlice("dialPatterns", dialPatterns)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which dial plan belongs.")
		cmd.Flags().StringSliceVar(&dialPatterns, "dial-patterns", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // list-dial-plans
		var orgId string
		var dialPlanName string
		var routeGroupName string
		var trunkName string
		var max string
		var start string
		var order string
		cmd := &cobra.Command{
			Use:   "list-dial-plans",
			Short: "Read the List of Dial Plans",
			Long:  "List all Dial Plans for the organization.\n\nDial plans route calls to on-premises destinations by use of the trunks or route groups with which the dial plan is associated. Multiple dial patterns can be defined as part of your dial plan.  Dial plans are configured globally for an enterprise and apply to all users, regardless of location.\n\nRetrieving this list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/dialPlans")
				req.QueryParam("orgId", orgId)
				req.QueryParam("dialPlanName", dialPlanName)
				req.QueryParam("routeGroupName", routeGroupName)
				req.QueryParam("trunkName", trunkName)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List dial plans for this organization.")
		cmd.Flags().StringVar(&dialPlanName, "dial-plan-name", "", "Return the list of dial plans matching the dial plan name.")
		cmd.Flags().StringVar(&routeGroupName, "route-group-name", "", "Return the list of dial plans matching the Route group name..")
		cmd.Flags().StringVar(&trunkName, "trunk-name", "", "Return the list of dial plans matching the Trunk name..")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&order, "order", "", "Order the dial plans according to the designated fields.  Available sort fields: `name`, `routeName`, `routeType`. Sort order is ascending by default")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // create-dial-plan
		var orgId string
		var name string
		var routeId string
		var routeType string
		var dialPatterns []string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-dial-plan",
			Short: "Create a Dial Plan",
			Long:  "Create a Dial Plan for the organization.\n\nDial plans route calls to on-premises destinations by use of trunks or route groups.\nThey are configured globally for an enterprise and apply to all users, regardless of location.\nA dial plan also specifies the routing choice (trunk or route group) for calls that match any of its dial patterns.\nSpecific dial patterns can be defined as part of your dial plan.\n\nCreating a dial plan requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/premisePstn/dialPlans")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("routeId", routeId)
					req.BodyString("routeType", routeType)
					req.BodyStringSlice("dialPatterns", dialPatterns)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which dial plan belongs.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&routeId, "route-id", "", "")
		cmd.Flags().StringVar(&routeType, "route-type", "", "")
		cmd.Flags().StringSliceVar(&dialPatterns, "dial-patterns", nil, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-dial-plan
		var dialPlanId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-dial-plan",
			Short: "Get a Dial Plan",
			Long:  "Get a Dial Plan for the organization.\n\nDial plans route calls to on-premises destinations by use of trunks or route groups.\nThey are configured globally for an enterprise and apply to all users, regardless of location.\nA dial plan also specifies the routing choice (trunk or route group) for calls that match any of its dial patterns.\nSpecific dial patterns can be defined as part of your dial plan.\n\nRetrieving a dial plan requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/dialPlans/{dialPlanId}")
				req.PathParam("dialPlanId", dialPlanId)
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
		cmd.Flags().StringVar(&dialPlanId, "dial-plan-id", "", "ID of the dial plan.")
		cmd.MarkFlagRequired("dial-plan-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which dial plan belongs.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // update-dial-plan
		var dialPlanId string
		var orgId string
		var name string
		var routeId string
		var routeType string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-dial-plan",
			Short: "Modify a Dial Plan",
			Long:  "Modify a Dial Plan for the organization.\n\nDial plans route calls to on-premises destinations by use of trunks or route groups.\nThey are configured globally for an enterprise and apply to all users, regardless of location.\nA dial plan also specifies the routing choice (trunk or route group) for calls that match any of its dial patterns.\nSpecific dial patterns can be defined as part of your dial plan.\n\nModifying a dial plan requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/premisePstn/dialPlans/{dialPlanId}")
				req.PathParam("dialPlanId", dialPlanId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("routeId", routeId)
					req.BodyString("routeType", routeType)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&dialPlanId, "dial-plan-id", "", "ID of the dial plan being modified.")
		cmd.MarkFlagRequired("dial-plan-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which dial plan belongs.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&routeId, "route-id", "", "")
		cmd.Flags().StringVar(&routeType, "route-type", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // delete-dial-plan
		var dialPlanId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-dial-plan",
			Short: "Delete a Dial Plan",
			Long:  "Delete a Dial Plan for the organization.\n\nDial plans route calls to on-premises destinations by use of trunks or route groups.\nThey are configured globally for an enterprise and apply to all users, regardless of location.\nA dial plan also specifies the routing choice (trunk or route group) for calls that match any of its dial patterns.\nSpecific dial patterns can be defined as part of your dial plan.\n\nDeleting a dial plan requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/premisePstn/dialPlans/{dialPlanId}")
				req.PathParam("dialPlanId", dialPlanId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&dialPlanId, "dial-plan-id", "", "ID of the dial plan.")
		cmd.MarkFlagRequired("dial-plan-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which dial plan belongs.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // validate-lgw-fqdn-domain-trunk
		var orgId string
		var address string
		var domain string
		var port int64
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "validate-lgw-fqdn-domain-trunk",
			Short: "Validate Local Gateway FQDN and Domain for a Trunk",
			Long:  "Validate Local Gateway FQDN and Domain for the organization trunks.\n\nA Trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group - a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nValidating Local Gateway FQDN and Domain requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/premisePstn/trunks/actions/fqdnValidation/invoke")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("address", address)
					req.BodyString("domain", domain)
					req.BodyInt("port", port, cmd.Flags().Changed("port"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which trunk types belongs.")
		cmd.Flags().StringVar(&address, "address", "", "")
		cmd.Flags().StringVar(&domain, "domain", "", "")
		cmd.Flags().Int64Var(&port, "port", 0, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // list-trunks
		var orgId string
		var name string
		var locationName string
		var trunkType string
		var max string
		var start string
		var order string
		cmd := &cobra.Command{
			Use:   "list-trunks",
			Short: "Read the List of Trunks",
			Long:  "List all Trunks for the organization.\n\nA Trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group - a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nRetrieving this list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/trunks")
				req.QueryParam("orgId", orgId)
				req.QueryParam("name", name)
				req.QueryParam("name", name)
				req.QueryParam("locationName", locationName)
				req.QueryParam("locationName", locationName)
				req.QueryParam("trunkType", trunkType)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List trunks for this organization.")
		cmd.Flags().StringVar(&name, "name", "", "Return the list of trunks matching the local gateway names.")
		cmd.Flags().StringVar(&locationName, "location-name", "", "Return the list of trunks matching the location names.")
		cmd.Flags().StringVar(&trunkType, "trunk-type", "", "Return the list of trunks matching the trunk type.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&order, "order", "", "Order the trunks according to the designated fields.  Available sort fields: name, locationName. Sort order is ascending by default")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // create-trunk
		var orgId string
		var name string
		var locationId string
		var password string
		var trunkType string
		var dualIdentitySupportEnabled bool
		var deviceType string
		var address string
		var domain string
		var port int64
		var maxConcurrentCalls int64
		var pChargeInfoSupportPolicy string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-trunk",
			Short: "Create a Trunk",
			Long:  "Create a Trunk for the organization.\n\nA Trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group which is a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nCreating a trunk requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/premisePstn/trunks")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("locationId", locationId)
					req.BodyString("password", password)
					req.BodyString("trunkType", trunkType)
					req.BodyBool("dualIdentitySupportEnabled", dualIdentitySupportEnabled, cmd.Flags().Changed("dual-identity-support-enabled"))
					req.BodyString("deviceType", deviceType)
					req.BodyString("address", address)
					req.BodyString("domain", domain)
					req.BodyInt("port", port, cmd.Flags().Changed("port"))
					req.BodyInt("maxConcurrentCalls", maxConcurrentCalls, cmd.Flags().Changed("max-concurrent-calls"))
					req.BodyString("pChargeInfoSupportPolicy", pChargeInfoSupportPolicy)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the trunk belongs.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&locationId, "location-id", "", "")
		cmd.Flags().StringVar(&password, "password", "", "")
		cmd.Flags().StringVar(&trunkType, "trunk-type", "", "")
		cmd.Flags().BoolVar(&dualIdentitySupportEnabled, "dual-identity-support-enabled", false, "")
		cmd.Flags().StringVar(&deviceType, "device-type", "", "")
		cmd.Flags().StringVar(&address, "address", "", "")
		cmd.Flags().StringVar(&domain, "domain", "", "")
		cmd.Flags().Int64Var(&port, "port", 0, "")
		cmd.Flags().Int64Var(&maxConcurrentCalls, "max-concurrent-calls", 0, "")
		cmd.Flags().StringVar(&pChargeInfoSupportPolicy, "p-charge-info-support-policy", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-trunk
		var trunkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-trunk",
			Short: "Get a Trunk",
			Long:  "Get a Trunk for the organization.\n\nA Trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group - a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nRetrieving a trunk requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/trunks/{trunkId}")
				req.PathParam("trunkId", trunkId)
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
		cmd.Flags().StringVar(&trunkId, "trunk-id", "", "ID of the trunk.")
		cmd.MarkFlagRequired("trunk-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which trunk belongs.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // update-trunk
		var trunkId string
		var orgId string
		var name string
		var password string
		var dualIdentitySupportEnabled bool
		var maxConcurrentCalls int64
		var pChargeInfoSupportPolicy string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-trunk",
			Short: "Modify a Trunk",
			Long:  "Modify a Trunk for the organization.\n\nA Trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group - a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nModifying a trunk requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/premisePstn/trunks/{trunkId}")
				req.PathParam("trunkId", trunkId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("password", password)
					req.BodyBool("dualIdentitySupportEnabled", dualIdentitySupportEnabled, cmd.Flags().Changed("dual-identity-support-enabled"))
					req.BodyInt("maxConcurrentCalls", maxConcurrentCalls, cmd.Flags().Changed("max-concurrent-calls"))
					req.BodyString("pChargeInfoSupportPolicy", pChargeInfoSupportPolicy)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&trunkId, "trunk-id", "", "ID of the trunk being modified.")
		cmd.MarkFlagRequired("trunk-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which trunk belongs.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&password, "password", "", "")
		cmd.Flags().BoolVar(&dualIdentitySupportEnabled, "dual-identity-support-enabled", false, "")
		cmd.Flags().Int64Var(&maxConcurrentCalls, "max-concurrent-calls", 0, "")
		cmd.Flags().StringVar(&pChargeInfoSupportPolicy, "p-charge-info-support-policy", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // delete-trunk
		var trunkId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-trunk",
			Short: "Delete a Trunk",
			Long:  "Delete a Trunk for the organization.\n\nA Trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group - a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nDeleting a trunk requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/premisePstn/trunks/{trunkId}")
				req.PathParam("trunkId", trunkId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&trunkId, "trunk-id", "", "ID of the trunk.")
		cmd.MarkFlagRequired("trunk-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which trunk belongs.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // list-trunk-types
		var orgId string
		cmd := &cobra.Command{
			Use:   "list-trunk-types",
			Short: "Read the List of Trunk Types",
			Long:  "List all Trunk Types with Device Types for the organization.\n\nA Trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group which is a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy. Trunk Types are Registering or Certificate Based and are configured in Call Manager.\n\nRetrieving trunk types requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/trunks/trunkTypes")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the trunk types belong.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // list-groups
		var orgId string
		var name string
		var max string
		var start string
		var order string
		cmd := &cobra.Command{
			Use:   "list-groups",
			Short: "Read the List of Routing Groups",
			Long:  "List all Route Groups for an organization. A Route Group is a group of trunks that allows further scale and redundancy with the connection to the premises.\n\nRetrieving this route group list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/routeGroups")
				req.QueryParam("orgId", orgId)
				req.QueryParam("name", name)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List route groups for this organization.")
		cmd.Flags().StringVar(&name, "name", "", "Return the list of route groups matching the Route group name..")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&order, "order", "", "Order the route groups according to designated fields.  Available sort orders are `asc` and `desc`.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // create-route-group
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-route-group",
			Short: "Create Route Group for a Organization",
			Long:  "Creates a Route Group for the organization.\n\nA Route Group is a collection of trunks that allows further scale and redundancy with the connection to the premises. Route groups can include up to 10 trunks from different locations.\n\nCreating a Route Group requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/premisePstn/routeGroups")
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the Route Group belongs.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-route-group
		var routeGroupId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-route-group",
			Short: "Read a Route Group for a Organization",
			Long:  "Reads a Route Group for the organization based on id.\n\nA Route Group is a collection of trunks that allows further scale and redundancy with the connection to the premises. Route groups can include up to 10 trunks from different locations.\n\nReading a Route Group requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/routeGroups/{routeGroupId}")
				req.PathParam("routeGroupId", routeGroupId)
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
		cmd.Flags().StringVar(&routeGroupId, "route-group-id", "", "Route Group for which details are being requested.")
		cmd.MarkFlagRequired("route-group-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization of the Route Group.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // update-route-group
		var routeGroupId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-route-group",
			Short: "Modify a Route Group for a Organization",
			Long:  "Modifies an existing Route Group for an organization based on id.\n\nA Route Group is a collection of trunks that allows further scale and redundancy with the connection to the premises. Route groups can include up to 10 trunks from different locations.\n\nModifying a Route Group requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/premisePstn/routeGroups/{routeGroupId}")
				req.PathParam("routeGroupId", routeGroupId)
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
		cmd.Flags().StringVar(&routeGroupId, "route-group-id", "", "Route Group for which details are being requested.")
		cmd.MarkFlagRequired("route-group-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization of the Route Group.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // delete-route-group-org
		var routeGroupId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-route-group-org",
			Short: "Remove a Route Group from an Organization",
			Long:  "Remove a Route Group from an Organization based on id.\n\nA Route Group is a collection of trunks that allows further scale and redundancy with the connection to the premises. Route groups can include up to 10 trunks from different locations.\n\nRemoving a Route Group requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/premisePstn/routeGroups/{routeGroupId}")
				req.PathParam("routeGroupId", routeGroupId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&routeGroupId, "route-group-id", "", "Route Group for which details are being requested.")
		cmd.MarkFlagRequired("route-group-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization of the Route Group.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-usage-group
		var routeGroupId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-usage-group",
			Short: "Read the Usage of a Routing Group",
			Long:  "List the number of \"Call to\" on-premises Extensions, Dial Plans, PSTN Connections, and Route Lists used by a specific Route Group.\nUsers within Call to Extension locations are registered to a PBX which allows you to route unknown extensions (calling number length of 2-6 digits) to the PBX using an existing Trunk or Route Group.\nPSTN Connections may be a Cisco PSTN, a cloud-connected PSTN, or a premises-based PSTN (local gateway).\nDial Plans allow you to route calls to on-premises extensions via your trunk or route group.\nRoute Lists are a list of numbers that can be reached via a route group and can be used to provide cloud PSTN connectivity to Webex Calling Dedicated Instance.\n\nRetrieving usage information requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/routeGroups/{routeGroupId}/usage")
				req.PathParam("routeGroupId", routeGroupId)
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
		cmd.Flags().StringVar(&routeGroupId, "route-group-id", "", "ID of the requested Route group.")
		cmd.MarkFlagRequired("route-group-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization associated with the specific route group.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-extension-locations-group
		var routeGroupId string
		var orgId string
		var locationName string
		var max string
		var start string
		var order string
		cmd := &cobra.Command{
			Use:   "get-extension-locations-group",
			Short: "Read the Call to Extension Locations of a Routing Group",
			Long:  "List \"Call to\" on-premises Extension Locations for a specific route group. Users within these locations are registered to a PBX which allows you to route unknown extensions (calling number length of 2-6 digits) to the PBX using an existing trunk or route group.\n\nRetrieving this location list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/routeGroups/{routeGroupId}/usageCallToExtension")
				req.PathParam("routeGroupId", routeGroupId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationName", locationName)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&routeGroupId, "route-group-id", "", "ID of the requested Route group.")
		cmd.MarkFlagRequired("route-group-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization associated with specific route group.")
		cmd.Flags().StringVar(&locationName, "location-name", "", "Return the list of locations matching the location name.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&order, "order", "", "Order the locations according to designated fields.  Available sort orders are `asc`, and `desc`.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-dial-plan-locations-group
		var routeGroupId string
		var orgId string
		var locationName string
		var max string
		var start string
		var order string
		cmd := &cobra.Command{
			Use:   "get-dial-plan-locations-group",
			Short: "Read the Dial Plan Locations of a Routing Group",
			Long:  "List Dial Plan Locations for a specific route group.\n\nDial Plans allow you to route calls to on-premises destinations by use of trunks or route groups. They are configured globally for an enterprise and apply to all users, regardless of location.\nA Dial Plan also specifies the routing choice (trunk or route group) for calls that match any of its dial patterns. Specific dial patterns can be defined as part of your dial plan.\n\nRetrieving this location list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/routeGroups/{routeGroupId}/usageDialPlan")
				req.PathParam("routeGroupId", routeGroupId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationName", locationName)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&routeGroupId, "route-group-id", "", "ID of the requested Route group.")
		cmd.MarkFlagRequired("route-group-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization associated with specific route group.")
		cmd.Flags().StringVar(&locationName, "location-name", "", "Return the list of locations matching the location name.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&order, "order", "", "Order the locations according to designated fields.  Available sort orders are `asc`, and `desc`.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-pstn-connection-locations-group
		var routeGroupId string
		var orgId string
		var locationName string
		var max string
		var start string
		var order string
		cmd := &cobra.Command{
			Use:   "get-pstn-connection-locations-group",
			Short: "Read the PSTN Connection Locations of a Routing Group",
			Long:  "List PSTN Connection Locations for a specific route group. This solution lets you configure users to use Cloud PSTN (CCP or Cisco PSTN) or Premises-based PSTN.\n\nRetrieving this Location list requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/routeGroups/{routeGroupId}/usagePstnConnection")
				req.PathParam("routeGroupId", routeGroupId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("locationName", locationName)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&routeGroupId, "route-group-id", "", "ID of the requested Route group.")
		cmd.MarkFlagRequired("route-group-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization associated with specific route group.")
		cmd.Flags().StringVar(&locationName, "location-name", "", "Return the list of locations matching the location name.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&order, "order", "", "Order the locations according to designated fields.  Available sort orders are `asc`, and `desc`.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-route-lists-group
		var routeGroupId string
		var orgId string
		var name string
		var max string
		var start string
		var order string
		cmd := &cobra.Command{
			Use:   "get-route-lists-group",
			Short: "Read the Route Lists of a Routing Group",
			Long:  "List Route Lists for a specific route group. Route Lists are a list of numbers that can be reached via a Route Group. It can be used to provide cloud PSTN connectivity to Webex Calling Dedicated Instance.\n\nRetrieving this list of Route Lists requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/routeGroups/{routeGroupId}/usageRouteList")
				req.PathParam("routeGroupId", routeGroupId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("name", name)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
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
		cmd.Flags().StringVar(&routeGroupId, "route-group-id", "", "ID of the requested Route group.")
		cmd.MarkFlagRequired("route-group-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization associated with specific route group.")
		cmd.Flags().StringVar(&name, "name", "", "Return the list of locations matching the location name.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&order, "order", "", "Order the locations according to designated fields.  Available sort orders are `asc`, and `desc`.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // list-route-lists
		var orgId string
		var start string
		var max string
		var order string
		var name string
		var locationId string
		cmd := &cobra.Command{
			Use:   "list-route-lists",
			Short: "Read the List of Route Lists",
			Long:  "List all Route Lists for the organization.\n\nA Route List is a list of numbers that can be reached via a Route Group. It can be used to provide cloud PSTN connectivity to Webex Calling Dedicated Instance.\n\nRetrieving the Route List requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/routeLists")
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
				req.QueryParam("max", max)
				req.QueryParam("order", order)
				req.QueryParam("name", name)
				req.QueryParam("name", name)
				req.QueryParam("locationId", locationId)
				req.QueryParam("locationId", locationId)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "List all Route List for this organization.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&order, "order", "", "Order the Route List according to the designated fields. Available sort fields are `name`, and `locationId`. Sort order is ascending by default")
		cmd.Flags().StringVar(&name, "name", "", "Return the list of Route List matching the route list name.")
		cmd.Flags().StringVar(&locationId, "location-id", "", "Return the list of Route Lists matching the location id.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // create-route-list
		var orgId string
		var name string
		var locationId string
		var routeGroupId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-route-list",
			Short: "Create a Route List",
			Long:  "Create a Route List for the organization.\n\nA Route List is a list of numbers that can be reached via a Route Group. It can be used to provide cloud PSTN connectivity to Webex Calling Dedicated Instance.\n\nCreating a Route List requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/premisePstn/routeLists")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("locationId", locationId)
					req.BodyString("routeGroupId", routeGroupId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the Route List belongs.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&locationId, "location-id", "", "")
		cmd.Flags().StringVar(&routeGroupId, "route-group-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // delete-route-list
		var routeListId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-route-list",
			Short: "Delete a Route List",
			Long:  "Delete a route list for a customer.\n\nA Route List is a list of numbers that can be reached via a Route Group. It can be used to provide cloud PSTN connectivity to Webex Calling Dedicated Instance.\n\nDeleting a Route List requires a full administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/premisePstn/routeLists/{routeListId}")
				req.PathParam("routeListId", routeListId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&routeListId, "route-list-id", "", "ID of the Route List.")
		cmd.MarkFlagRequired("route-list-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the Route List belongs.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-route-list
		var routeListId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-route-list",
			Short: "Get a Route List",
			Long:  "Get a rout list details.\n\nA Route List is a list of numbers that can be reached via a Route Group. It can be used to provide cloud PSTN connectivity to Webex Calling Dedicated Instance.\n\nRetrieving a Route List requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/routeLists/{routeListId}")
				req.PathParam("routeListId", routeListId)
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
		cmd.Flags().StringVar(&routeListId, "route-list-id", "", "ID of the Route List.")
		cmd.MarkFlagRequired("route-list-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the Route List belongs.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // update-route-list
		var routeListId string
		var orgId string
		var name string
		var routeGroupId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-route-list",
			Short: "Modify a Route List",
			Long:  "Modify the details for a Route List.\n\nA Route List is a list of numbers that can be reached via a Route Group. It can be used to provide cloud PSTN connectivity to Webex Calling Dedicated Instance.\n\nRetrieving a Route List requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/premisePstn/routeLists/{routeListId}")
				req.PathParam("routeListId", routeListId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("routeGroupId", routeGroupId)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&routeListId, "route-list-id", "", "ID of the Route List.")
		cmd.MarkFlagRequired("route-list-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the Route List belongs.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&routeGroupId, "route-group-id", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // update-numbers-route-list
		var routeListId string
		var orgId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-numbers-route-list",
			Short: "Modify Numbers for Route List",
			Long:  "Modify numbers for a specific Route List of a Customer.\n\nA Route List is a list of numbers that can be reached via a Route Group. It can be used to provide cloud PSTN connectivity to Webex Calling Dedicated Instance.\n\nRetrieving a Route List requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/premisePstn/routeLists/{routeListId}/numbers")
				req.PathParam("routeListId", routeListId)
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
		cmd.Flags().StringVar(&routeListId, "route-list-id", "", "ID of the Route List.")
		cmd.MarkFlagRequired("route-list-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the Route List belongs.")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-numbers-assigned-route-list
		var routeListId string
		var orgId string
		var start string
		var max string
		var number string
		var order string
		cmd := &cobra.Command{
			Use:   "get-numbers-assigned-route-list",
			Short: "Get Numbers assigned to a Route List",
			Long:  "Get numbers assigned to a Route List\n\nA Route List is a list of numbers that can be reached via a Route Group. It can be used to provide cloud PSTN connectivity to Webex Calling Dedicated Instance.\n\nRetrieving a Route List requires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_write`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/routeLists/{routeListId}/numbers")
				req.PathParam("routeListId", routeListId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
				req.QueryParam("max", max)
				req.QueryParam("number", number)
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
		cmd.Flags().StringVar(&routeListId, "route-list-id", "", "ID of the Route List.")
		cmd.MarkFlagRequired("route-list-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the Route List belongs.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&number, "number", "", "Number assigned to the route list.")
		cmd.Flags().StringVar(&order, "order", "", "Order the Route Lists according to number, ascending or descending.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-lgw-on-premises-extension-usage-trunk
		var trunkId string
		var orgId string
		var start string
		var max string
		var order string
		var name string
		cmd := &cobra.Command{
			Use:   "get-lgw-on-premises-extension-usage-trunk",
			Short: "Get Local Gateway Call to On-Premises Extension Usage for a Trunk",
			Long:  "Get local gateway call to on-premises extension usage for a trunk.\n\nA trunk is a connection between Webex Calling and the premises, which terminates on the premises with a local gateway or other supported device.\nThe trunk can be assigned to a Route Group which is a group of trunks that allow Webex Calling to distribute calls over multiple trunks or to provide redundancy.\n\nRetrieving this information requires a full administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/premisePstn/trunks/{trunkId}/usageCallToExtension")
				req.PathParam("trunkId", trunkId)
				req.QueryParam("orgId", orgId)
				req.QueryParam("start", start)
				req.QueryParam("max", max)
				req.QueryParam("order", order)
				req.QueryParam("name", name)
				req.QueryParam("name", name)
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
		cmd.Flags().StringVar(&trunkId, "trunk-id", "", "ID of the trunk.")
		cmd.MarkFlagRequired("trunk-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Organization to which the trunk belongs.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&order, "order", "", "Order the trunks according to the designated fields.  Available sort fields are `name`, and `locationName`. Sort order is ascending by default")
		cmd.Flags().StringVar(&name, "name", "", "Return the list of trunks matching the local gateway names")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // create-translation-pattern
		var orgId string
		var name string
		var matchingPattern string
		var replacementPattern string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-translation-pattern",
			Short: "Create a Translation Pattern for an Organization",
			Long:  "Create a translation pattern for a given organization.\n\nA translation pattern lets you manipulate dialed digits before routing a call and applies to outbound calls only. See [this article](https://help.webex.com/en-us/article/nib9o6h/Translation-patterns-for-outbound-calls) for details about the translation pattern syntax.\n\nRequires a full administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/callRouting/translationPatterns")
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("matchingPattern", matchingPattern)
					req.BodyString("replacementPattern", replacementPattern)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization containing the translation pattern.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&matchingPattern, "matching-pattern", "", "")
		cmd.Flags().StringVar(&replacementPattern, "replacement-pattern", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // list-translation-patterns
		var orgId string
		var limitToLocationId string
		var limitToOrgLevelEnabled string
		var max string
		var start string
		var order string
		var name string
		var matchingPattern string
		cmd := &cobra.Command{
			Use:   "list-translation-patterns",
			Short: "Retrieve the list of Translation Patterns",
			Long:  "Retrieve a list of translation patterns for a given organization.\n\nA translation pattern lets you manipulate dialed digits before routing a call and applies to outbound calls only. See [this article](https://help.webex.com/en-us/article/nib9o6h/Translation-patterns-for-outbound-calls) for details about the translation pattern syntax.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/callRouting/translationPatterns")
				req.QueryParam("orgId", orgId)
				req.QueryParam("limitToLocationId", limitToLocationId)
				req.QueryParam("limitToOrgLevelEnabled", limitToOrgLevelEnabled)
				req.QueryParam("max", max)
				req.QueryParam("start", start)
				req.QueryParam("order", order)
				req.QueryParam("name", name)
				req.QueryParam("matchingPattern", matchingPattern)
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
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization containing the translation patterns.")
		cmd.Flags().StringVar(&limitToLocationId, "limit-to-location-id", "", "When a location ID is passed, then return only the corresponding location level translation patterns.")
		cmd.Flags().StringVar(&limitToOrgLevelEnabled, "limit-to-org-level-enabled", "", "When set to be `true`, then return only the organization-level translation patterns.")
		cmd.Flags().StringVar(&max, "max", "", "Limit the number of objects returned to this maximum count.")
		cmd.Flags().StringVar(&start, "start", "", "Start at the zero-based offset in the list of matching objects.")
		cmd.Flags().StringVar(&order, "order", "", "Sort the list of translation patterns according to translation pattern name, ascending or descending.")
		cmd.Flags().StringVar(&name, "name", "", "Only return translation patterns with the matching `name`.")
		cmd.Flags().StringVar(&matchingPattern, "matching-pattern", "", "Only return translation patterns with the matching `matchingPattern`.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-translation-pattern
		var translationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-translation-pattern",
			Short: "Retrieve a specific Translation Pattern for an Organization",
			Long:  "Retrieve the details of a translation pattern for a given organization.\n\nA translation pattern lets you manipulate dialed digits before routing a call and applies to outbound calls only. See [this article](https://help.webex.com/en-us/article/nib9o6h/Translation-patterns-for-outbound-calls) for details about the translation pattern syntax.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/callRouting/translationPatterns/{translationId}")
				req.PathParam("translationId", translationId)
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
		cmd.Flags().StringVar(&translationId, "translation-id", "", "Retrieve the translation pattern with the matching ID.")
		cmd.MarkFlagRequired("translation-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization containing the translation pattern.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // update-translation-pattern
		var translationId string
		var orgId string
		var name string
		var matchingPattern string
		var replacementPattern string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-translation-pattern",
			Short: "Modify a specific Translation Pattern for an Organization",
			Long:  "Modify a translation pattern for a given organization.\n\nA translation pattern lets you manipulate dialed digits before routing a call and applies to outbound calls only. See [this article](https://help.webex.com/en-us/article/nib9o6h/Translation-patterns-for-outbound-calls) for details about the translation pattern syntax.\n\nRequires a full administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/callRouting/translationPatterns/{translationId}")
				req.PathParam("translationId", translationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("matchingPattern", matchingPattern)
					req.BodyString("replacementPattern", replacementPattern)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&translationId, "translation-id", "", "Modify translation pattern with the matching ID.")
		cmd.MarkFlagRequired("translation-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization containing the translation pattern.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&matchingPattern, "matching-pattern", "", "")
		cmd.Flags().StringVar(&replacementPattern, "replacement-pattern", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // delete-translation-pattern
		var translationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-translation-pattern",
			Short: "Delete a specific Translation Pattern",
			Long:  "Delete a translation pattern for a given organization.\n\nA translation pattern lets you manipulate dialed digits before routing a call and applies to outbound calls only. See [this article](https://help.webex.com/en-us/article/nib9o6h/Translation-patterns-for-outbound-calls) for details about the translation pattern syntax.\n\nRequires a full administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/callRouting/translationPatterns/{translationId}")
				req.PathParam("translationId", translationId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&translationId, "translation-id", "", "Delete a translation pattern with the matching ID.")
		cmd.MarkFlagRequired("translation-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "ID of the organization containing the translation pattern.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // create-translation-pattern-location
		var locationId string
		var orgId string
		var name string
		var matchingPattern string
		var replacementPattern string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-translation-pattern-location",
			Short: "Create a Translation Pattern for a Location",
			Long:  "Create a translation pattern for a given location.\n\nA translation pattern lets you manipulate dialed digits before routing a call and applies to outbound calls only. See [this article](https://help.webex.com/en-us/article/nib9o6h/Translation-patterns-for-outbound-calls) for details about the translation pattern syntax.\n\nRequires a full administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "POST", "/telephony/config/locations/{locationId}/callRouting/translationPatterns")
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
					req.BodyString("matchingPattern", matchingPattern)
					req.BodyString("replacementPattern", replacementPattern)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Only admin users of another organization (such as partners) may use this parameter since the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&matchingPattern, "matching-pattern", "", "")
		cmd.Flags().StringVar(&replacementPattern, "replacement-pattern", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // get-translation-pattern-location
		var locationId string
		var translationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "get-translation-pattern-location",
			Short: "Retrieve a specific Translation Pattern for a Location",
			Long:  "Retrieve a specific translation pattern for a given location.\n\nA translation pattern lets you manipulate dialed digits before routing a call and applies to outbound calls only. See [this article](https://help.webex.com/en-us/article/nib9o6h/Translation-patterns-for-outbound-calls) for details about the translation pattern syntax.\n\nRequires a full or read-only administrator auth token with a scope of `spark-admin:telephony_config_read`.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "GET", "/telephony/config/locations/{locationId}/callRouting/translationPatterns/{translationId}")
				req.PathParam("locationId", locationId)
				req.PathParam("translationId", translationId)
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
		cmd.Flags().StringVar(&locationId, "location-id", "", "Unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&translationId, "translation-id", "", "Unique identifier for the translation pattern.")
		cmd.MarkFlagRequired("translation-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Only admin users of another organization (such as partners) may use this parameter since the default is the same organization as the token used to access API.")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // update-translation-pattern-location
		var locationId string
		var translationId string
		var orgId string
		var name string
		var matchingPattern string
		var replacementPattern string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-translation-pattern-location",
			Short: "Modify a specific Translation Pattern for a Location",
			Long:  "Modify a specific translation pattern for a given location.\n\nA translation pattern lets you manipulate dialed digits before routing a call and applies to outbound calls only. See [this article](https://help.webex.com/en-us/article/nib9o6h/Translation-patterns-for-outbound-calls) for details about the translation pattern syntax.\n\nRequires a full administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "PUT", "/telephony/config/locations/{locationId}/callRouting/translationPatterns/{translationId}")
				req.PathParam("locationId", locationId)
				req.PathParam("translationId", translationId)
				req.QueryParam("orgId", orgId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("name", name)
					req.BodyString("matchingPattern", matchingPattern)
					req.BodyString("replacementPattern", replacementPattern)
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&translationId, "translation-id", "", "Unique identifier for the translation pattern.")
		cmd.MarkFlagRequired("translation-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Only admin users of another organization (such as partners) may use this parameter since the default is the same organization as the token used to access API.")
		cmd.Flags().StringVar(&name, "name", "", "")
		cmd.Flags().StringVar(&matchingPattern, "matching-pattern", "", "")
		cmd.Flags().StringVar(&replacementPattern, "replacement-pattern", "", "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		callRoutingCmd.AddCommand(cmd)
	}

	{ // delete-translation-pattern-location
		var locationId string
		var translationId string
		var orgId string
		cmd := &cobra.Command{
			Use:   "delete-translation-pattern-location",
			Short: "Delete a specific Translation Pattern for a Location",
			Long:  "Delete a specific translation pattern for a given location.\n\nA translation pattern lets you manipulate dialed digits before routing a call and applies to outbound calls only. See [this article](https://help.webex.com/en-us/article/nib9o6h/Translation-patterns-for-outbound-calls) for details about the translation pattern syntax.\n\nRequires a full administrator auth token with the `spark-admin:telephony_config_write` scope.",
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CallingBaseURL, "DELETE", "/telephony/config/locations/{locationId}/callRouting/translationPatterns/{translationId}")
				req.PathParam("locationId", locationId)
				req.PathParam("translationId", translationId)
				req.QueryParam("orgId", orgId)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&locationId, "location-id", "", "Unique identifier for the location.")
		cmd.MarkFlagRequired("location-id")
		cmd.Flags().StringVar(&translationId, "translation-id", "", "Unique identifier for the translation pattern.")
		cmd.MarkFlagRequired("translation-id")
		cmd.Flags().StringVar(&orgId, "org-id", "", "Only admin users of another organization (such as partners) may use this parameter since the default is the same organization as the token used to access API.")
		callRoutingCmd.AddCommand(cmd)
	}

}
