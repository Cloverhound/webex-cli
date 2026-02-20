package cc

import (
	"fmt"

	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/Cloverhound/webex-cli/internal/search"
	"github.com/Cloverhound/webex-cli/internal/timeutil"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Contact Center historical data (GraphQL)",
	Long: `Query Webex Contact Center Search API for tasks, agent sessions, flow interactions, and more.

The Search API is a GraphQL endpoint that queries historical contact center data.
All queries require a time window (--from/--to or --last). The FROM parameter
cannot be older than 36 months from the current time.

Subcommands provide pre-built queries with convenience flags for filtering,
aggregation, and field selection. Use "raw" for full GraphQL control.

Examples:
  webex cc search tasks --last 24h
  webex cc search tasks --last 7d --channel telephony --direction inbound
  webex cc search agent-sessions --last 24h
  webex cc search tasks --last 7d --aggregate totalDuration:sum --aggregate id:count
  webex cc search raw --body '{"query":"{ task(from:1700000000000 to:1700100000000) { tasks { id } } }"}'`,
}

func init() {
	cmd.CcCmd.AddCommand(searchCmd)

	registerSearchRaw(searchCmd)
	registerSearchTasks(searchCmd)
	registerSearchTaskDetails(searchCmd)
	registerSearchTaskLegDetails(searchCmd)
	registerSearchAgentSessions(searchCmd)
	registerSearchFlowInteractions(searchCmd)
	registerSearchFlowTrace(searchCmd)
}

// searchFlags holds common search flag values.
type searchFlags struct {
	from, to, last           string
	timeComparator           string
	filter                   string
	channel, direction       string
	status, agentID, queueID string
	fields                   string
	cursor                   string
	pageSize                 int
	aggregations             []string
	interval, timezone       string
	bodyRaw, bodyFile        string
	orgID, trackingID        string
}

// registerCommonFlags adds all common search flags to a cobra command.
func registerCommonFlags(c *cobra.Command, f *searchFlags) {
	c.Flags().StringVar(&f.from, "from", "", "Start time (epoch ms, ISO 8601, YYYY-MM-DD, or relative like -24h)")
	c.Flags().StringVar(&f.to, "to", "", "End time (same formats as --from, or \"now\")")
	c.Flags().StringVar(&f.last, "last", "", "Time range shorthand (e.g. 24h, 7d). Sets --from/--to automatically")
	c.Flags().StringVar(&f.timeComparator, "time-comparator", "", "Time field to filter on: createdTime (default) or endedTime")
	c.Flags().StringVar(&f.filter, "filter", "", "Raw JSON filter (AND-merged with convenience flags)")
	c.Flags().StringVar(&f.channel, "channel", "", "Filter by channel type: telephony, email, chat, social")
	c.Flags().StringVar(&f.direction, "direction", "", "Filter by direction: inbound, outbound")
	c.Flags().StringVar(&f.status, "status", "", "Filter by status: new, parked, connected, ended")
	c.Flags().StringVar(&f.agentID, "agent-id", "", "Filter by agent/owner ID")
	c.Flags().StringVar(&f.queueID, "queue-id", "", "Filter by queue ID")
	c.Flags().StringVar(&f.fields, "fields", "", "GraphQL field selection (overrides defaults)")
	c.Flags().StringVar(&f.cursor, "cursor", "", "Pagination cursor (from previous response's endCursor)")
	c.Flags().IntVar(&f.pageSize, "page-size", 0, "Page size (for flow queries)")
	c.Flags().StringSliceVar(&f.aggregations, "aggregate", nil, "Aggregation in field:type[:name] format (repeatable). Types: sum, average, count, min, max, cardinality")
	c.Flags().StringVar(&f.interval, "interval", "", "Aggregation interval: FIFTEEN_MINUTES, THIRTY_MINUTES, HOURLY, DAILY, WEEKLY, MONTHLY")
	c.Flags().StringVar(&f.timezone, "timezone", "", "Timezone for interval aggregation (e.g. America/Los_Angeles)")
	c.Flags().StringVar(&f.bodyRaw, "body", "", "Raw JSON body (bypasses flag-based query building)")
	c.Flags().StringVar(&f.bodyFile, "body-file", "", "Path to JSON body file (bypasses flag-based query building)")
	c.Flags().StringVar(&f.orgID, "orgid", "", "Organization ID (defaults to token org)")
	c.Flags().StringVar(&f.trackingID, "tracking-id", "", "Tracking ID for traceability")
}

// buildAndExecute parses flags, builds the GraphQL query, and executes it.
func buildAndExecute(f *searchFlags, queryType search.QueryType) error {
	req := client.NewRequest(config.CcBaseURL, "POST", "/search")
	req.QueryParam("orgId", f.orgID)
	req.Header("TrackingId", f.trackingID)

	// If raw body provided, skip flag-based query building
	if f.bodyFile != "" {
		if err := req.SetBodyFile(f.bodyFile); err != nil {
			return err
		}
		resp, statusCode, err := req.Do()
		if err != nil {
			return err
		}
		return output.Print(resp, statusCode)
	}
	if f.bodyRaw != "" {
		req.SetBodyRaw(f.bodyRaw)
		resp, statusCode, err := req.Do()
		if err != nil {
			return err
		}
		return output.Print(resp, statusCode)
	}

	// Resolve time window
	var fromMs, toMs int64
	var err error

	if f.last != "" {
		fromStr, toStr, parseErr := timeutil.ParseLastEpochMs(f.last)
		if parseErr != nil {
			return parseErr
		}
		fmt.Sscanf(fromStr, "%d", &fromMs)
		fmt.Sscanf(toStr, "%d", &toMs)
	} else {
		if f.from == "" || f.to == "" {
			return fmt.Errorf("--from and --to are required (or use --last)")
		}
		fromMs, err = search.ParseFlexibleTime(f.from)
		if err != nil {
			return fmt.Errorf("--from: %w", err)
		}
		toMs, err = search.ParseFlexibleTime(f.to)
		if err != nil {
			return fmt.Errorf("--to: %w", err)
		}
	}

	// Build filter
	filterStr, err := search.BuildFilter(f.channel, f.direction, f.status, f.agentID, f.queueID, f.filter)
	if err != nil {
		return err
	}

	// Parse aggregations
	var aggs []search.Aggregation
	for _, a := range f.aggregations {
		agg, err := search.ParseAggregation(a)
		if err != nil {
			return err
		}
		aggs = append(aggs, agg)
	}

	params := search.Params{
		From:           fromMs,
		To:             toMs,
		TimeComparator: f.timeComparator,
		Filter:         filterStr,
		Fields:         f.fields,
		Cursor:         f.cursor,
		PageSize:       f.pageSize,
		Aggregations:   aggs,
		Interval:       f.interval,
		Timezone:       f.timezone,
	}

	// Auto-paginate if --paginate is set
	if config.Paginate() {
		resp, statusCode, err := search.PaginateGraphQL(queryType, params, f.orgID, f.trackingID)
		if err != nil {
			return err
		}
		return output.Print(resp, statusCode)
	}

	body, err := search.BuildQuery(queryType, params)
	if err != nil {
		return err
	}

	req.SetBodyRaw(body)
	resp, statusCode, err := req.Do()
	if err != nil {
		return err
	}
	return output.Print(resp, statusCode)
}

// --- Subcommand registrations ---

func registerSearchRaw(parent *cobra.Command) {
	var orgID string
	var trackingID string
	var bodyRaw string
	var bodyFile string
	raw := &cobra.Command{
		Use:   "raw",
		Short: "Send a raw GraphQL query",
		Long: `Send a raw GraphQL query to the Search API.

This is the low-level escape hatch for queries not covered by the other
subcommands. You must provide the full GraphQL query as JSON.

Examples:
  webex cc search raw --body '{"query":"{ task(from:1700000000000 to:1700100000000) { tasks { id status } } }"}'
  webex cc search raw --body-file my_query.json`,
		RunE: func(c *cobra.Command, args []string) error {
			req := client.NewRequest(config.CcBaseURL, "POST", "/search")
			req.QueryParam("orgId", orgID)
			req.Header("TrackingId", trackingID)
			if bodyFile != "" {
				if err := req.SetBodyFile(bodyFile); err != nil {
					return err
				}
			} else if bodyRaw != "" {
				req.SetBodyRaw(bodyRaw)
			} else {
				return fmt.Errorf("--body or --body-file is required")
			}
			resp, statusCode, err := req.Do()
			if err != nil {
				return err
			}
			return output.Print(resp, statusCode)
		},
	}
	raw.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
	raw.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
	raw.Flags().StringVar(&orgID, "orgid", "", "Organization ID (defaults to token org)")
	raw.Flags().StringVar(&trackingID, "tracking-id", "", "Tracking ID for traceability")
	parent.AddCommand(raw)

	// Backward-compatible alias
	alias := &cobra.Command{
		Use:    "search-tasks",
		Short:  "Send a raw GraphQL query (alias for 'raw')",
		Hidden: true,
		RunE:   raw.RunE,
	}
	alias.Flags().StringVar(&bodyRaw, "body", "", "Raw JSON body")
	alias.Flags().StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
	alias.Flags().StringVar(&orgID, "orgid", "", "Organization ID (defaults to token org)")
	alias.Flags().StringVar(&trackingID, "tracking-id", "", "Tracking ID for traceability")
	parent.AddCommand(alias)
}

func registerSearchTasks(parent *cobra.Command) {
	f := &searchFlags{}
	c := &cobra.Command{
		Use:   "tasks",
		Short: "Query task/contact summary records",
		Long: `Query summary-level task (contact) records from the Search API.

Returns one record per customer interaction with key metrics like duration,
queue time, and agent information.

Common fields: id, status, channelType, direction, origin, destination,
  createdTime, endedTime, totalDuration, connectedDuration, queueDuration,
  holdDuration, wrapupDuration, ringingDuration, selfserviceDuration,
  owner {id name}, lastEntryPoint {id name}, lastQueue {id name},
  lastTeam {id name}, lastWrapupCodeName, csatScore

Filter fields: channelType (telephony/email/chat/social), direction (inbound/outbound),
  status (new/parked/connected/ended), isActive, origin, destination,
  owner.id, lastQueue.id, lastWrapupCodeName, totalDuration, connectedDuration

Aggregation fields: id, totalDuration, connectedDuration, holdDuration,
  queueDuration, wrapupDuration, ringingDuration, selfserviceDuration,
  csatScore, connectedCount, holdCount, conferenceCount, consultCount,
  blindTransferCount, transferCount

Examples:
  webex cc search tasks --last 24h
  webex cc search tasks --last 7d --channel telephony --direction inbound
  webex cc search tasks --from 2024-01-15 --to 2024-01-16
  webex cc search tasks --last 7d --aggregate totalDuration:sum --aggregate id:count
  webex cc search tasks --last 7d --aggregate connectedDuration:average --interval DAILY
  webex cc search tasks --last 24h --status ended --queue-id UUID`,
		RunE: func(c *cobra.Command, args []string) error {
			return buildAndExecute(f, search.QueryTask)
		},
	}
	registerCommonFlags(c, f)
	parent.AddCommand(c)
}

func registerSearchTaskDetails(parent *cobra.Command) {
	f := &searchFlags{}
	c := &cobra.Command{
		Use:   "task-details",
		Short: "Query extended task records with activity drill-down",
		Long: `Query enriched task records with full activity history, email fields,
campaign data, recording info, and AI analytics.

Includes everything from 'tasks' plus: activities (step-by-step workflow),
contributors (all agents), email fields, campaign fields, monitoring/recording
fields, and AI analytics (sentiment, CSAT, topic classification).

Additional fields: activities {nodes {activityName eventName agentName duration}},
  contributors {id name}, csatScore, sentiment, autoCsat, terminationType,
  emailBody, isCampaign, campaignName, isMonitored, recordingCount,
  topicName, contactDriver

Examples:
  webex cc search task-details --last 24h
  webex cc search task-details --last 7d --channel telephony --agent-id UUID
  webex cc search task-details --last 1h --filter '{"id":{"equals":"task-uuid-here"}}'`,
		RunE: func(c *cobra.Command, args []string) error {
			return buildAndExecute(f, search.QueryTaskDetails)
		},
	}
	registerCommonFlags(c, f)
	parent.AddCommand(c)
}

func registerSearchTaskLegDetails(parent *cobra.Command) {
	f := &searchFlags{}
	c := &cobra.Command{
		Use:   "task-leg-details",
		Short: "Query per-leg/per-queue breakdown of tasks",
		Long: `Query per-leg breakdown of task interactions. Each leg represents the
contact's interaction with a specific queue or agent.

Key fields: id, taskId, status, channelType, direction, origin, destination,
  entryPoint {id name}, queue {id name}, owner {id name}, site {id name},
  callLegType (main/consult/conference), handleType (short/abandoned/normal/dequeued),
  isTaskLegHandled, isWithinServiceLevel, connectedDuration, holdDuration,
  queueDuration, wrapupDuration

Examples:
  webex cc search task-leg-details --last 24h
  webex cc search task-leg-details --last 7d --queue-id UUID
  webex cc search task-leg-details --last 7d --channel telephony`,
		RunE: func(c *cobra.Command, args []string) error {
			return buildAndExecute(f, search.QueryTaskLegDetails)
		},
	}
	registerCommonFlags(c, f)
	parent.AddCommand(c)
}

func registerSearchAgentSessions(parent *cobra.Command) {
	f := &searchFlags{}
	c := &cobra.Command{
		Use:   "agent-sessions",
		Short: "Query agent login sessions with per-channel stats",
		Long: `Query agent login session records with per-channel performance breakdown.

Key fields: agentSessionId, agentId, agentName, userLoginId, isActive,
  state (logged_in/logged_out), startTime, endTime, siteName, teamName

Per-channel stats (channelInfo): channelType, totalDuration, idleDuration/Count,
  availableDuration/Count, connectedDuration/Count, holdDuration/Count,
  wrapupDuration/Count, notRespondedCount, ringingDuration/Count

Filter fields: agentId, teamId, siteId, state, isActive,
  channelInfo.channelType, channelInfo.currentState

Examples:
  webex cc search agent-sessions --last 24h
  webex cc search agent-sessions --last 7d --agent-id UUID
  webex cc search agent-sessions --last 7d --aggregate agentId:cardinality:Unique_Agents
  webex cc search agent-sessions --last 7d --filter '{"state":{"equals":"logged_in"}}'`,
		RunE: func(c *cobra.Command, args []string) error {
			return buildAndExecute(f, search.QueryAgentSession)
		},
	}
	registerCommonFlags(c, f)
	parent.AddCommand(c)
}

func registerSearchFlowInteractions(parent *cobra.Command) {
	f := &searchFlags{}
	c := &cobra.Command{
		Use:   "flow-interactions",
		Short: "Query flow execution metadata",
		Long: `Query flow execution summaries showing which flows ran and their outcomes.

Key fields: interactionId, flowId, flowVersionId, entryPointName,
  lastExecutedActivity, mainFlowActivities, eventFlowActivities,
  flowStartTime, flowEndTime, status

Note: Flow queries typically require a filter (e.g. by interactionId or flowId)
and use page-based pagination (--page-size) instead of cursor-based.

Examples:
  webex cc search flow-interactions --last 24h --filter '{"interactionId":{"equals":"task-uuid"}}'
  webex cc search flow-interactions --last 7d --page-size 100`,
		RunE: func(c *cobra.Command, args []string) error {
			return buildAndExecute(f, search.QueryFlowInteractions)
		},
	}
	registerCommonFlags(c, f)
	parent.AddCommand(c)
}

func registerSearchFlowTrace(parent *cobra.Command) {
	f := &searchFlags{}
	c := &cobra.Command{
		Use:   "flow-trace",
		Short: "Query granular flow trace events",
		Long: `Query per-activity flow trace records showing individual activity inputs,
outputs, outcomes, and variable changes.

Key fields: interactionId, flowId, activityName, activityProcessId, outcome,
  activityInputs {name value}, activityOutput {name type value},
  modifiedFlowVariables {name type value}, flowStartTime, flowEndTime

Note: Flow queries typically require a filter (e.g. by interactionId) and
use page-based pagination (--page-size).

Examples:
  webex cc search flow-trace --last 24h --filter '{"interactionId":{"equals":"task-uuid"}}'
  webex cc search flow-trace --last 7d --page-size 500`,
		RunE: func(c *cobra.Command, args []string) error {
			return buildAndExecute(f, search.QueryFlowTraceEvents)
		},
	}
	registerCommonFlags(c, f)
	parent.AddCommand(c)
}
