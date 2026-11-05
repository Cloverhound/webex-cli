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
	"fmt"
	"net/url"
	"strings"

	"github.com/Cloverhound/webex-cli/pkg/webexapi"
)

// Host names the API host an operation is served from.
type Host string

const (
	HostWebex            Host = "webex"
	HostCC               Host = "cc"
	HostAnalytics        Host = "analytics"
	HostAnalyticsCalling Host = "analytics-calling"
	HostAIAgent          Host = "ai-agent"
)

// URL returns the host's base URL.
func (h Host) URL(hosts webexapi.Hosts) string {
	switch h {
	case HostCC:
		return hosts.CC
	case HostAnalytics:
		return hosts.Analytics
	case HostAnalyticsCalling:
		return hosts.AnalyticsCalling
	case HostAIAgent:
		return hosts.AIAgent
	default:
		return hosts.Webex
	}
}

// In says where a parameter travels.
type In string

const (
	InPath   In = "path"
	InQuery  In = "query"
	InHeader In = "header"

	// InAgentToken is AI Agent Studio's bot-access-token header. Callers supply the
	// agent's ID, and the token is looked up from it.
	InAgentToken In = "agent-token"
)

// TimeRange says how a --last duration fills an operation's from and to.
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

	// Flag is the CLI flag that supplies the value, in kebab case. Parameters that
	// share a flag share a value.
	Flag  string
	Usage string

	// Value is fixed, so no flag supplies it.
	Value string
}

// Body says an operation sends a JSON body.
type Body struct {
	// Fields are top-level body fields the CLI offers as flags. Bodies that are
	// not flat have none.
	Fields []BodyField
}

// BodyField is one top-level body field.
type BodyField struct {
	Name string
	Flag string

	// Type is string, bool, int64, float64 or []string.
	Type string
}

// Paging says how a list reads to its end.
type Paging struct {
	// Offset names offset paging's parameters. Nil means the area's default: start
	// and max outside Contact Center, page and pageSize in it.
	Offset *OffsetPaging

	// SizeParam is Contact Center's page size parameter when it is not pageSize.
	SizeParam string

	// Query is set on every page request.
	Query map[string]string
}

// OffsetPaging names an offset and page size parameter. An empty Start means the
// list pages only through Link headers. First is the first item's offset.
type OffsetPaging struct {
	Start, Size string
	First       int
}

// Group is a set of operations published under one CLI command.
type Group struct {
	Name    string
	Aliases []string
	Summary string

	// Handwritten groups are published by handwritten CLI commands.
	Handwritten bool
}

// Operation is one API call.
type Operation struct {
	// Area, Group and Command are the CLI names: webex <area> <group> <command>.
	Area, Group, Command string
	Aliases              []string

	Summary, Description string

	Method string
	// Path is relative to Host, with {name} placeholders for path parameters.
	Path string
	Host Host

	Params []Param
	Body   *Body

	// Paging is set on reads that can page.
	Paging *Paging

	TimeRange TimeRange

	// Private operations are on an API that Cisco neither documents nor supports.
	Private bool

	// Redact means responses can carry credentials, which must be removed before
	// they are shown.
	Redact bool

	// Accept overrides the default application/json Accept header.
	Accept string

	// Handwritten operations are published by a handwritten CLI command instead,
	// because the request needs something the table cannot describe, such as a
	// file upload.
	Handwritten bool
}

// Name is the operation's CLI command path without "webex".
func (o *Operation) Name() string { return o.Area + " " + o.Group + " " + o.Command }

// Param returns the parameter with this API name.
func (o *Operation) Param(name string) (Param, bool) {
	for _, p := range o.Params {
		if p.Name == name {
			return p, true
		}
	}
	return Param{}, false
}

// Areas lists the operation areas in CLI order.
var Areas = []string{"admin", "calling", "cc", "device", "meetings", "messaging"}

// Operations lists an area's operations in collection order.
func Operations(area string) []*Operation {
	switch area {
	case "admin":
		return adminOperations
	case "calling":
		return callingOperations
	case "cc":
		return ccOperations
	case "device":
		return deviceOperations
	case "meetings":
		return meetingsOperations
	case "messaging":
		return messagingOperations
	}
	return nil
}

// Groups lists an area's groups in collection order.
func Groups(area string) []Group {
	switch area {
	case "admin":
		return adminGroups
	case "calling":
		return callingGroups
	case "cc":
		return ccGroups
	case "device":
		return deviceGroups
	case "meetings":
		return meetingsGroups
	case "messaging":
		return messagingGroups
	}
	return nil
}

// Find returns the operation published as webex <area> <group> <command>, the
// first one when a collection repeats a name.
func Find(area, group, command string) (*Operation, bool) {
	for _, op := range Operations(area) {
		if op.Group == group && op.Command == command {
			return op, true
		}
	}
	return nil, false
}

// Args are the values for one call, keyed by API parameter name. Fixed parameters
// need no value.
type Args struct {
	Path   map[string]string
	Query  map[string]string
	Header map[string]string
	Body   []byte
}

// URL builds the request URL. Every path parameter needs a value; an argument
// the operation does not declare is an error, so a misspelled name is never
// silently dropped.
func (o *Operation) URL(hosts webexapi.Hosts, args Args) (string, error) {
	declared := map[In]map[string]bool{InPath: {}, InQuery: {}, InHeader: {}, InAgentToken: {}}
	path := o.Path
	for _, p := range o.Params {
		declared[p.In][p.Name] = true
		if p.In != InPath {
			continue
		}
		v := p.Value
		if v == "" {
			v = args.Path[p.Name]
		}
		if v == "" {
			return "", fmt.Errorf("%s: path parameter %s is required", o.Name(), p.Name)
		}
		path = strings.ReplaceAll(path, "{"+p.Name+"}", url.PathEscape(v))
	}
	for name := range args.Path {
		if !declared[InPath][name] {
			return "", fmt.Errorf("%s has no path parameter %s", o.Name(), name)
		}
	}
	query := url.Values{}
	for name, v := range args.Query {
		if !declared[InQuery][name] {
			return "", fmt.Errorf("%s has no query parameter %s", o.Name(), name)
		}
		if v != "" {
			query.Set(name, v)
		}
	}
	u := o.Host.URL(hosts) + path
	if encoded := query.Encode(); encoded != "" {
		u += "?" + encoded
	}
	return u, nil
}

// Call sends the operation. Header values in args are sent as given; resolving an
// AI agent's token is the caller's job, since it needs another call.
func Call(ctx context.Context, c *webexapi.Client, hosts webexapi.Hosts, op *Operation, args Args) (*webexapi.Response, error) {
	u, err := op.URL(hosts, args)
	if err != nil {
		return nil, err
	}
	if len(args.Body) > 0 && op.Body == nil {
		return nil, fmt.Errorf("%s takes no body", op.Name())
	}
	req := webexapi.Request{Method: op.Method, URL: u, Body: args.Body, Header: map[string][]string{}}
	for name, v := range args.Header {
		req.Header.Set(name, v)
	}
	if op.Accept != "" {
		req.Header.Set("Accept", op.Accept)
	}
	return c.Send(ctx, req)
}
