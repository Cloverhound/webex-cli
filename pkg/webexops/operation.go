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

// Upload says an operation's body is a multipart form. The collections describe
// these forms only in words, if at all, so each was measured against Webex.
type Upload struct {
	// Parts are the form's parts, in the order they are sent.
	Parts []Part
}

// PartKind says what a form part carries.
type PartKind string

const (
	// PartFile is a file's contents, sent with the file's name.
	PartFile PartKind = "file"

	// PartJSON is a JSON document. The CLI reads it from --body or --body-file.
	PartJSON PartKind = "json"

	// PartFileName is a name for the upload's file part. The CLI takes it from
	// that file's path unless the part's flag gives one.
	PartFileName PartKind = "file-name"

	// PartText is a plain value, such as an ID, which the part's flag gives.
	PartText PartKind = "text"
)

// Part is one part of an upload's form.
type Part struct {
	// Name is the part's form field name.
	Name string
	Kind PartKind

	// Flag is the CLI flag that gives a file's path, a file name or a text
	// value. A JSON part has none.
	Flag     string
	Usage    string
	Required bool
}

// Group is a set of operations published under one CLI command.
type Group struct {
	Area, Name string
	Aliases    []string
	Summary    string
}

// Operation is one API call.
type Operation struct {
	// ID is the operation's Go variable name: its area, method and path, such as
	// CcGetOrganizationByOrgidV2Team. A route the collection lists more than once
	// with different requests is exported as one operation with every listing's
	// parameters; that operation is not a CLI command, and its listings are found
	// with Operations.
	ID string

	// Area, Group and Command are the CLI names: webex <area> <group> <command>.
	// Command is empty for an operation that is not a CLI command.
	Area, Group, Command string
	Aliases              []string

	Summary, Description string

	Method string
	Host   Host

	// Path is the path template, with each path parameter as {name}.
	Path string

	Params []Param

	// Body is the JSON body the collection gives the operation, nil when it gives
	// none. A body in Args is sent either way, such as a multipart upload the
	// collection describes only in words.
	Body *Body

	// Upload is the multipart form the operation sends, nil when it sends none.
	// Call sends the form that Args gives as its Body and ContentType.
	Upload *Upload

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

	// Body is sent when it is not empty; ContentType defaults to
	// application/json.
	Body        []byte
	ContentType string

	// AllowEmptyPath sends a path parameter with no value as its {name}, as the
	// CLI does for an empty flag. Without it, a missing path value is an error.
	AllowEmptyPath bool
}

// Request is the request that sends an operation with its arguments. Arguments
// the operation cannot send are an error: a name it has no parameter for, a value
// for a fixed parameter, or a missing path value.
func (op *Operation) Request(hosts webexapi.Hosts, args Args) (webexapi.Request, error) {
	l, err := op.List(hosts, args)
	if err != nil {
		return webexapi.Request{}, err
	}
	return webexapi.Request{
		Method: l.Method, URL: webexapi.WithQuery(l.URL, l.Query), Body: l.Body,
		ContentType: l.ContentType, Header: l.Header, Redact: l.Redact,
	}, nil
}

// List is the request that reads an operation's list with its arguments, before
// paging. Its arguments are checked as Request's are.
func (op *Operation) List(hosts webexapi.Hosts, args Args) (webexapi.List, error) {
	if err := op.check(args); err != nil {
		return webexapi.List{}, err
	}
	return op.list(hosts, args), nil
}

// check refuses arguments the operation cannot send as given.
func (op *Operation) check(args Args) error {
	params := map[string]Param{}
	for _, p := range op.Params {
		params[p.Name] = p
	}
	for name, v := range args.Params {
		p, ok := params[name]
		switch {
		case !ok:
			return fmt.Errorf("%s has no parameter %q", op.ID, name)
		case p.Value != "" && v != "":
			return fmt.Errorf("%s's parameter %q is fixed", op.ID, name)
		}
	}
	for name := range args.Query {
		if p, ok := params[name]; !ok || p.In != InQuery {
			return fmt.Errorf("%s has no query parameter %q", op.ID, name)
		}
	}
	var missing []string
	for _, p := range op.Params {
		if p.In == InPath && p.Value == "" && args.Params[p.Name] == "" && !args.AllowEmptyPath {
			missing = append(missing, p.Name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("%s needs the path parameter %s", op.ID, strings.Join(missing, ", "))
	}
	return nil
}

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
	l.Body, l.ContentType = args.Body, args.ContentType
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

// Call sends an operation. Arguments it cannot send are an error, and nothing is
// sent; see Request. An operation marked Redact returns its body with AI agent
// credentials removed; send op.Request with the client for the raw body.
func Call(ctx context.Context, c *webexapi.Client, hosts webexapi.Hosts, op *Operation, args Args) (*webexapi.Response, error) {
	r, err := op.Request(hosts, args)
	if err != nil {
		return nil, err
	}
	resp, err := c.Send(ctx, r)
	if resp != nil && op.Redact {
		resp.Body = aiagent.RedactSecrets(resp.Body)
	}
	return resp, err
}

// All reads every page of an operation's list and returns its items, checked
// and redacted as Call's are.
func All(ctx context.Context, c *webexapi.Client, hosts webexapi.Hosts, op *Operation, args Args) ([]json.RawMessage, error) {
	if op.Paging == nil {
		return nil, fmt.Errorf("%s is not a list that can be read to its end", op.ID)
	}
	l, err := op.List(hosts, args)
	if err != nil {
		return nil, err
	}
	items, err := c.All(ctx, l, *op.Paging)
	if op.Redact {
		for i := range items {
			items[i] = aiagent.RedactSecrets(items[i])
		}
	}
	return items, err
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

// Operations lists an area's operations as the CLI's commands, in collection
// order; operations that combine a route's listings are not in it. The list is a
// copy; the operations it points to are shared and must not be changed.
func Operations(areaName string) []*Operation { return slices.Clone(areas[areaName].operations) }
