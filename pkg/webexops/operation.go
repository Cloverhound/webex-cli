// Package webexops is the table of Webex API operations, generated from Cisco's
// Postman collections by codegen/generate_cli.py. Each operation is an exported
// *Operation, so a program that names one fails to build when a collection
// refresh drops or renames it. Call sends an operation through a webexapi.Client.
//
// Bodies and responses are JSON documents, not Go types: the collections' example
// bodies do not match what Webex returns, so generated structs would be wrong.
package webexops

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strings"

	"github.com/Cloverhound/webex-cli/pkg/aiagent"
	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// Host names the API host an operation is served from.
type Host string

const (
	HostWebex            Host = "webex"
	HostIdentity         Host = "identity"
	HostAnalytics        Host = "analytics"
	HostAnalyticsCalling Host = "analytics-calling"
	HostCC               Host = "cc"
	HostAIAgent          Host = "ai-agent"
)

// URL returns the host's base URL.
func (h Host) URL(hosts webexapi.Hosts) string {
	switch h {
	case HostIdentity:
		return hosts.Identity
	case HostAnalytics:
		return hosts.Analytics
	case HostAnalyticsCalling:
		return hosts.AnalyticsCalling
	case HostCC:
		return hosts.CC
	case HostAIAgent:
		return hosts.AIAgent
	default:
		return hosts.Webex
	}
}

// In says where a parameter is sent.
type In string

const (
	InPath   In = "path"
	InQuery  In = "query"
	InHeader In = "header"

	// InAgentToken is AI Agent Studio's bot-access-token header, which selects
	// the agent of a v1 call. Its value is the agent's runtime token; see
	// aiagent.BotToken. The CLI looks it up from --agent-id.
	InAgentToken In = "agent-token"
)

// TimeRange says how the CLI's --last fills an operation's from and to
// parameters.
type TimeRange string

const (
	TimeISO     TimeRange = "iso"
	TimeEpochMs TimeRange = "epoch-ms"
)

// Param is one path, query or header parameter.
type Param struct {
	// Name is the parameter as the API names it.
	Name string
	In   In

	// Value is fixed; no flag gives it.
	Value string

	// Flag is the CLI flag that gives the value, in kebab case. Parameters that
	// share a flag share a value, and the first of them declares the flag with
	// its Usage and Required.
	Flag     string
	Usage    string
	Required bool

	// OrgUUID means the API takes the organization's UUID. Call converts a base64
	// organization ID to it.
	OrgUUID bool
}

// Body says an operation sends a body.
type Body struct {
	// Fields are top-level body fields the CLI offers as flags. Bodies that are
	// not flat have none.
	Fields []BodyField
}

// BodyField is one top-level body field.
type BodyField struct {
	// Name is the field as the API names it.
	Name string

	// Flag gives the field's value. A field that repeats a parameter's name uses
	// that parameter's flag.
	Flag string

	// Type is string, bool, int64, float64 or []string.
	Type string
}

// Group is a set of operations published under one CLI command.
type Group struct {
	Area, Name string
	Aliases    []string
	Summary    string
}

// Operation is one API call.
type Operation struct {
	// ID is the operation's Go variable name, such as CcTeamList.
	ID string

	// Area, Group and Command are the CLI names: webex <area> <group> <command>.
	Area, Group, Command string
	Aliases              []string

	Summary, Description string

	Method string
	Host   Host

	// Path is the path template, with each path parameter as {name}.
	Path string

	Params []Param

	// Body is nil when the operation sends none.
	Body *Body

	// Paging is how the operation's list is read to its end; nil when it is not
	// a list, or its paging is unknown.
	Paging *webexapi.Paging

	TimeRange TimeRange

	// Accept, when set, replaces application/json as the accepted type.
	Accept string

	// Redact means responses carry AI agent credentials; see
	// aiagent.RedactSecrets. Call redacts them from errors, and the caller from
	// what it shows.
	Redact bool

	// Private means the API is AI Agent Studio's, which Cisco does not document or
	// support.
	Private bool
}

// Args are an operation's values.
type Args struct {
	// Params are the parameter values by API name. Empty values are not sent.
	Params map[string]string

	// Query gives a query parameter's values when it is sent more than once, in
	// place of its value in Params.
	Query url.Values

	// Body is sent when the operation has a Body; ContentType defaults to
	// application/json.
	Body        []byte
	ContentType string
}

// Request is the request that sends an operation with its arguments.
func (op *Operation) Request(hosts webexapi.Hosts, args Args) webexapi.Request {
	l := op.list(hosts, args)
	return webexapi.Request{
		Method: l.Method, URL: webexapi.WithQuery(l.URL, l.Query), Body: l.Body,
		ContentType: l.ContentType, Header: l.Header, Redact: l.Redact,
	}
}

// List is the request that reads an operation's list with its arguments, before
// paging.
func (op *Operation) List(hosts webexapi.Hosts, args Args) webexapi.List { return op.list(hosts, args) }

func (op *Operation) list(hosts webexapi.Hosts, args Args) webexapi.List {
	path := map[string]string{}
	query := url.Values{}
	header := http.Header{}
	if op.Accept != "" {
		header.Set("Accept", op.Accept)
	}
	for _, p := range op.Params {
		v := p.Value
		if v == "" {
			v = args.Params[p.Name]
		}
		if v == "" {
			continue
		}
		if p.OrgUUID {
			v = webexapi.DecodeOrgID(v)
		}
		switch p.In {
		case InPath:
			path[p.Name] = v
		case InQuery:
			query.Set(p.Name, v)
		case InHeader, InAgentToken:
			header.Set(p.Name, v)
		}
	}
	for k, vs := range args.Query {
		vs = slices.Clone(vs)
		if op.takesOrgUUID(k) {
			for i := range vs {
				vs[i] = webexapi.DecodeOrgID(vs[i])
			}
		}
		query[k] = vs
	}
	l := webexapi.List{
		Method: op.Method,
		URL:    webexapi.FillPath(op.Host.URL(hosts)+op.Path, path),
		Query:  query,
		Header: header,
		Name:   op.Path,
	}
	if op.Body != nil {
		l.Body, l.ContentType = args.Body, args.ContentType
	}
	if op.Redact {
		l.Redact = aiagent.RedactSecrets
	}
	return l
}

// takesOrgUUID reports whether the query parameter name takes the org UUID.
func (op *Operation) takesOrgUUID(name string) bool {
	for _, p := range op.Params {
		if p.In == InQuery && p.Name == name && p.OrgUUID {
			return true
		}
	}
	return false
}

// missingPath reports the path parameters args gives no value for.
func (op *Operation) missingPath(args Args) error {
	var missing []string
	for _, p := range op.Params {
		if p.In == InPath && p.Value == "" && args.Params[p.Name] == "" {
			missing = append(missing, p.Name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%s needs the path parameter %s", op.ID, strings.Join(missing, ", "))
	}
	return nil
}

// Call sends an operation. A path parameter without a value is an error, and
// nothing is sent.
func Call(ctx context.Context, c *webexapi.Client, hosts webexapi.Hosts, op *Operation, args Args) (*webexapi.Response, error) {
	if err := op.missingPath(args); err != nil {
		return nil, err
	}
	return c.Send(ctx, op.Request(hosts, args))
}

// All reads every page of an operation's list and returns its items. A path
// parameter without a value is an error, and nothing is sent.
func All(ctx context.Context, c *webexapi.Client, hosts webexapi.Hosts, op *Operation, args Args) ([]json.RawMessage, error) {
	if op.Paging == nil {
		return nil, fmt.Errorf("%s is not a list that can be read to its end", op.ID)
	}
	if err := op.missingPath(args); err != nil {
		return nil, err
	}
	return c.All(ctx, op.List(hosts, args), *op.Paging)
}

type area struct {
	groups     []Group
	operations []*Operation
}

// Areas lists the areas, such as calling and cc.
func Areas() []string {
	names := make([]string, 0, len(areas))
	for name := range areas {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Groups lists an area's groups, in collection order. The list is a copy.
func Groups(areaName string) []Group {
	groups := slices.Clone(areas[areaName].groups)
	for i := range groups {
		groups[i].Aliases = slices.Clone(groups[i].Aliases)
	}
	return groups
}

// Operations lists an area's operations, in collection order. The list is a copy;
// the operations it points to are shared and must not be changed.
func Operations(areaName string) []*Operation { return slices.Clone(areas[areaName].operations) }
