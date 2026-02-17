package cc

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

var autoCsatCmd = &cobra.Command{
	Use:   "auto-csat",
	Short: "AutoCsat commands",
}

func init() {
	cmd.CcCmd.AddCommand(autoCsatCmd)

	{ // create-mapped-question
		var orgid string
		var autoCsatId string
		var questionId string
		var questionnaireId string
		var organizationId string
		var id string
		var version int64
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "create-mapped-question",
			Short: "Create a new Auto CSAT mapped Question",
			Long:  `Create a new Auto CSAT mapped Question in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/auto-csat/{autoCsatId}/question")
				req.PathParam("orgid", orgid)
				req.PathParam("autoCsatId", autoCsatId)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("questionId", questionId)
					req.BodyString("questionnaireId", questionnaireId)
					req.BodyString("organizationId", organizationId)
					req.BodyString("id", id)
					req.BodyInt("version", version, cmd.Flags().Changed("version"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&autoCsatId, "auto-csat-id", "", "Resource ID of the Auto CSAT resource")
		cmd.MarkFlagRequired("auto-csat-id")
		cmd.Flags().StringVar(&questionId, "question-id", "", "")
		cmd.Flags().StringVar(&questionnaireId, "questionnaire-id", "", "")
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "")
		cmd.Flags().StringVar(&id, "id", "", "")
		cmd.Flags().Int64Var(&version, "version", 0, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		autoCsatCmd.AddCommand(cmd)
	}

	{ // bulk-save-mapped-question
		var orgid string
		var autoCsatId string
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "bulk-save-mapped-question",
			Short: "Bulk save Auto CSAT mapped Question(s)",
			Long:  `Create, Update or delete Auto CSAT mapped Question(s) in bulk for Auto CSAT resource in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "POST", "/organization/{orgid}/auto-csat/{autoCsatId}/question/bulk")
				req.PathParam("orgid", orgid)
				req.PathParam("autoCsatId", autoCsatId)
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
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&autoCsatId, "auto-csat-id", "", "Resource ID of the Auto CSAT resource")
		cmd.MarkFlagRequired("auto-csat-id")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		autoCsatCmd.AddCommand(cmd)
	}

	{ // get-mapped-question-id
		var orgid string
		var autoCsatId string
		var id string
		cmd := &cobra.Command{
			Use:   "get-mapped-question-id",
			Short: "Get specific Auto CSAT mapped Question by ID",
			Long:  `Retrieve an existing Auto CSAT mapped Question by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/auto-csat/{autoCsatId}/question/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("autoCsatId", autoCsatId)
				req.PathParam("id", id)
				if config.Paginate() {
					resp, statusCode, err := req.DoPaginated(false)
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
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&autoCsatId, "auto-csat-id", "", "Resource ID of the Auto CSAT resource")
		cmd.MarkFlagRequired("auto-csat-id")
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Auto CSAT mapped Question.")
		cmd.MarkFlagRequired("id")
		autoCsatCmd.AddCommand(cmd)
	}

	{ // delete-mapped-question-id
		var orgid string
		var autoCsatId string
		var id string
		cmd := &cobra.Command{
			Use:   "delete-mapped-question-id",
			Short: "Delete specific Auto CSAT mapped Question by ID",
			Long:  `Delete an existing Auto CSAT mapped Question by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "DELETE", "/organization/{orgid}/auto-csat/{autoCsatId}/question/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("autoCsatId", autoCsatId)
				req.PathParam("id", id)
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&autoCsatId, "auto-csat-id", "", "Resource ID of the Auto CSAT resource")
		cmd.MarkFlagRequired("auto-csat-id")
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Auto CSAT mapped Question.")
		cmd.MarkFlagRequired("id")
		autoCsatCmd.AddCommand(cmd)
	}

	{ // get-id
		var orgid string
		var id string
		cmd := &cobra.Command{
			Use:   "get-id",
			Short: "Get specific Auto CSAT resource by ID",
			Long:  `Retrieve an existing Auto CSAT resource by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/auto-csat/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
				if config.Paginate() {
					resp, statusCode, err := req.DoPaginated(false)
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
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Auto CSAT resource.")
		cmd.MarkFlagRequired("id")
		autoCsatCmd.AddCommand(cmd)
	}

	{ // update-id
		var orgid string
		var id string
		var agentInclusionType string
		var enabled bool
		var selectedGlobalVariableId string
		var surveyDataSource string
		var organizationId string
		var version int64
		var bodyRaw string
		var bodyFile string
		cmd := &cobra.Command{
			Use:   "update-id",
			Short: "Update specific Auto CSAT resource by ID",
			Long:  `Update an existing Auto CSAT resource by ID in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "PUT", "/organization/{orgid}/auto-csat/{id}")
				req.PathParam("orgid", orgid)
				req.PathParam("id", id)
				if bodyFile != "" {
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				} else if bodyRaw != "" {
					req.SetBodyRaw(bodyRaw)
				} else {
					req.BodyString("agentInclusionType", agentInclusionType)
					req.BodyBool("enabled", enabled, cmd.Flags().Changed("enabled"))
					req.BodyString("selectedGlobalVariableId", selectedGlobalVariableId)
					req.BodyString("surveyDataSource", surveyDataSource)
					req.BodyString("organizationId", organizationId)
					req.BodyString("id", id)
					req.BodyInt("version", version, cmd.Flags().Changed("version"))
				}
				resp, statusCode, err := req.Do()
				if err != nil {
					return err
				}
				return output.Print(resp, statusCode)
			},
		}
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&id, "id", "", "Resource ID of the Auto CSAT resource.")
		cmd.MarkFlagRequired("id")
		cmd.Flags().StringVar(&agentInclusionType, "agent-inclusion-type", "", "")
		cmd.Flags().BoolVar(&enabled, "enabled", false, "")
		cmd.Flags().StringVar(&selectedGlobalVariableId, "selected-global-variable-id", "", "")
		cmd.Flags().StringVar(&surveyDataSource, "survey-data-source", "", "")
		cmd.Flags().StringVar(&organizationId, "organization-id", "", "")
		cmd.Flags().Int64Var(&version, "version", 0, "")
		cmd.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
		cmd.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
		autoCsatCmd.AddCommand(cmd)
	}

	{ // list
		var orgid string
		var filter string
		var attributes string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "list",
			Short: "List Auto CSAT resource(s)",
			Long:  `Retrieve a list of Auto CSAT resource(s) in a given organization.Only one entry per organization can exist for Auto CSAT resource.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v2/auto-csat")
				req.PathParam("orgid", orgid)
				req.QueryParam("filter", filter)
				req.QueryParam("attributes", attributes)
				req.QueryParam("page", page)
				req.QueryParam("pageSize", pageSize)
				if config.Paginate() {
					resp, statusCode, err := req.DoPaginated(false)
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
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&filter, "filter", "", "Specify a filter based on which the results will be fetched. All the fields are supported except: organizationId, createdTime, lastUpdatedTime   The examples below show some search queries - id==\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id!=\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id=in=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") - id=out=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see  <a href=\"https://www.here.com/docs/bundle/data-client-library-developer-guide-java-scala/page/client/rsql.html\">this reference</a>. For a list of supported operators, see <a href=\"https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference\">this syntax guide</a>.  Note: values to be used in the filter syntax should not contain space, and if so kindly bound it with quotes to apply filter. ")
		cmd.Flags().StringVar(&attributes, "attributes", "", "Specify the attributes to be returned.Default all attributes are returned along with specified columns. All Attributes are supported")
		cmd.Flags().StringVar(&page, "page", "", "Defines the number of displayed page. The page number starts from 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Defines the number of items to be displayed on a page. If the number specified is more than allowed max page size, the API will automatically adjust the page size to the max page size.")
		autoCsatCmd.AddCommand(cmd)
	}

	{ // list-mapped-question
		var orgid string
		var autoCsatId string
		var filter string
		var attributes string
		var page string
		var pageSize string
		cmd := &cobra.Command{
			Use:   "list-mapped-question",
			Short: "List Auto CSAT mapped Question(s)",
			Long:  `Retrieve a list of Auto CSAT mapped Question(s) in a given organization.`,
			RunE: func(cmd *cobra.Command, args []string) error {
				req := client.NewRequest(config.CcBaseURL, "GET", "/organization/{orgid}/v2/auto-csat/{autoCsatId}/question")
				req.PathParam("orgid", orgid)
				req.PathParam("autoCsatId", autoCsatId)
				req.QueryParam("filter", filter)
				req.QueryParam("attributes", attributes)
				req.QueryParam("page", page)
				req.QueryParam("pageSize", pageSize)
				if config.Paginate() {
					resp, statusCode, err := req.DoPaginated(false)
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
		cmd.Flags().StringVar(&orgid, "orgid", "", "Organization ID to be used for this operation. The specified security token must have permission to interact with the organization.")
		cmd.MarkFlagRequired("orgid")
		cmd.Flags().StringVar(&autoCsatId, "auto-csat-id", "", "Resource ID of the Auto CSAT resource")
		cmd.MarkFlagRequired("auto-csat-id")
		cmd.Flags().StringVar(&filter, "filter", "", "Specify a filter based on which the results will be fetched. All the fields are supported except: organizationId, createdTime, lastUpdatedTime   The examples below show some search queries - id==\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id!=\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\" - id=in=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") - id=out=(\"57efb0e6-5af0-4245-a67d-d3c5045cdb6e\",\"a421e0b2-732e-46f3-a057-39160a53afb9\") This parameter uses the RSQL query syntax, a URI-friendly format for expressing criteria for filtering REST entities. For more information about RSQL in general, see  <a href=\"https://www.here.com/docs/bundle/data-client-library-developer-guide-java-scala/page/client/rsql.html\">this reference</a>. For a list of supported operators, see <a href=\"https://github.com/perplexhub/rsql-jpa-specification#rsql-syntax-reference\">this syntax guide</a>.  Note: values to be used in the filter syntax should not contain space, and if so kindly bound it with quotes to apply filter. ")
		cmd.Flags().StringVar(&attributes, "attributes", "", "Specify the attributes to be returned.Default all attributes are returned along with specified columns. All Attributes are supported(id, questionId, questionnaireId)")
		cmd.Flags().StringVar(&page, "page", "", "Defines the number of displayed page. The page number starts from 0.")
		cmd.Flags().StringVar(&pageSize, "page-size", "", "Defines the number of items to be displayed on a page. If the number specified is more than allowed max page size, the API will automatically adjust the page size to the max page size.")
		autoCsatCmd.AddCommand(cmd)
	}

}
