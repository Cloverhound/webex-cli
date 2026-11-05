// Package apicmd builds the CLI's API commands from the operation table in
// pkg/webexops. Each operation becomes `webex <area> <group> <command>`, with a
// flag per parameter and body field.
package apicmd

import (
	"strings"

	"github.com/Cloverhound/webex-cli/internal/aiagent"
	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/Cloverhound/webex-cli/internal/timeutil"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
	"github.com/spf13/cobra"
)

// PrivateAPINote ends the help of commands on an API Cisco does not support.
const PrivateAPINote = "Caution: this command uses Webex AI Agent Studio's private API, which is not " +
	"documented or supported by Cisco and can change or stop working without notice. " +
	"Do not build production automation on it."

// Register adds an area's groups and commands under parent and returns the group
// commands by name, for handwritten commands to join.
func Register(parent *cobra.Command, area string) map[string]*cobra.Command {
	groups := map[string]*cobra.Command{}
	for _, g := range webexops.Groups(area) {
		if g.Handwritten {
			continue
		}
		gc := &cobra.Command{Use: g.Name, Aliases: g.Aliases, Short: g.Summary}
		parent.AddCommand(gc)
		groups[g.Name] = gc
	}
	for _, op := range webexops.Operations(area) {
		if op.Handwritten {
			continue
		}
		groups[op.Group].AddCommand(NewCommand(op))
	}
	return groups
}

// BaseURL is the base URL of an operation's host for this run.
func BaseURL(h webexops.Host) string { return h.URL(config.Hosts()) }

// flagValues holds a command's flag values by flag name.
type flagValues struct {
	strings map[string]*string
	ints    map[string]*int64
	floats  map[string]*float64
	bools   map[string]*bool
	slices  map[string]*[]string
}

func (v *flagValues) text(flag string) string {
	if p := v.strings[flag]; p != nil {
		return *p
	}
	return ""
}

// NewCommand builds the command for one operation.
func NewCommand(op *webexops.Operation) *cobra.Command {
	values := &flagValues{
		strings: map[string]*string{}, ints: map[string]*int64{}, floats: map[string]*float64{},
		bools: map[string]*bool{}, slices: map[string]*[]string{},
	}
	var last, bodyRaw, bodyFile string

	description := op.Description
	if op.Private {
		if description != "" {
			description += "\n\n" + PrivateAPINote
		} else {
			description = PrivateAPINote
		}
	}
	isCalling := op.Area != "cc"

	c := &cobra.Command{
		Use:     op.Command,
		Aliases: op.Aliases,
		Short:   op.Summary,
		Long:    description,
		RunE: func(cmd *cobra.Command, args []string) error {
			if op.Private {
				aiagent.PrivateAPINotice()
			}
			req := client.NewRequest(BaseURL(op.Host), op.Method, op.Path)
			if op.Accept != "" {
				req.Header("Accept", op.Accept)
			}
			if op.TimeRange != "" && last != "" {
				parse := timeutil.ParseLastISO
				if op.TimeRange == webexops.TimeEpochMs {
					parse = timeutil.ParseLastEpochMs
				}
				from, to, err := parse(last)
				if err != nil {
					return err
				}
				for flag, v := range map[string]string{"from": from, "to": to} {
					if p := values.strings[flag]; p != nil {
						*p = v
					}
				}
			}
			for _, p := range op.Params {
				switch p.In {
				case webexops.InPath:
					v := p.Value
					if v == "" {
						v = values.text(p.Flag)
					}
					req.PathParam(p.Name, v)
				case webexops.InQuery:
					req.QueryParam(p.Name, values.text(p.Flag))
				}
			}
			for _, p := range op.Params {
				switch p.In {
				case webexops.InAgentToken:
					token, err := aiagent.BotToken(values.text(p.Flag))
					if err != nil {
						return err
					}
					req.Header(p.Name, token)
				case webexops.InHeader:
					req.Header(p.Name, values.text(p.Flag))
				}
			}
			if op.Body != nil {
				switch {
				case bodyFile != "":
					if err := req.SetBodyFile(bodyFile); err != nil {
						return err
					}
				case bodyRaw != "":
					req.SetBodyRaw(bodyRaw)
				default:
					setBodyFields(cmd, req, op.Body.Fields, values)
				}
			}
			if op.Paging != nil && config.Paginate() {
				if o := op.Paging.Offset; o != nil {
					req.OffsetPaging(o.Start, o.Size, o.First)
				}
				if op.Paging.SizeParam != "" {
					req.PageSizeParam(op.Paging.SizeParam)
				}
				for k, v := range op.Paging.Query {
					req.QueryParam(k, v)
				}
				resp, statusCode, err := req.DoPaginated(isCalling)
				if err != nil {
					return err
				}
				if op.Redact {
					resp = aiagent.RedactSecrets(resp)
				}
				return output.Print(resp, statusCode)
			}
			resp, statusCode, err := req.Do()
			if err != nil {
				return err
			}
			if op.Redact {
				resp = aiagent.RedactSecrets(resp)
			}
			return output.Print(resp, statusCode)
		},
	}

	flags := c.Flags()
	for _, p := range op.Params {
		if p.Value != "" || values.strings[p.Flag] != nil {
			continue
		}
		values.strings[p.Flag] = new(string)
		flags.StringVar(values.strings[p.Flag], p.Flag, "", strings.ReplaceAll(p.Usage, "\n", " "))
		if p.In == webexops.InPath || p.In == webexops.InAgentToken {
			c.MarkFlagRequired(p.Flag)
		}
	}
	if op.Body != nil {
		for _, f := range op.Body.Fields {
			if flags.Lookup(f.Flag) != nil {
				continue
			}
			switch f.Type {
			case "int64":
				values.ints[f.Flag] = new(int64)
				flags.Int64Var(values.ints[f.Flag], f.Flag, 0, "")
			case "float64":
				values.floats[f.Flag] = new(float64)
				flags.Float64Var(values.floats[f.Flag], f.Flag, 0, "")
			case "bool":
				values.bools[f.Flag] = new(bool)
				flags.BoolVar(values.bools[f.Flag], f.Flag, false, "")
			case "[]string":
				values.slices[f.Flag] = new([]string)
				flags.StringSliceVar(values.slices[f.Flag], f.Flag, nil, "")
			default:
				values.strings[f.Flag] = new(string)
				flags.StringVar(values.strings[f.Flag], f.Flag, "", "")
			}
		}
		flags.StringVar(&bodyRaw, "body", "", "Raw JSON body")
		flags.StringVar(&bodyFile, "body-file", "", "Path to JSON body file")
	}
	if op.TimeRange != "" {
		flags.StringVar(&last, "last", "", "Time range shorthand (e.g. 1h, 30m, 24h). Sets --from automatically.")
	}
	return c
}

func setBodyFields(cmd *cobra.Command, req *client.Request, fields []webexops.BodyField, values *flagValues) {
	for _, f := range fields {
		changed := cmd.Flags().Changed(f.Flag)
		switch f.Type {
		case "int64":
			req.BodyInt(f.Name, *values.ints[f.Flag], changed)
		case "float64":
			req.BodyFloat(f.Name, *values.floats[f.Flag], changed)
		case "bool":
			req.BodyBool(f.Name, *values.bools[f.Flag], changed)
		case "[]string":
			req.BodyStringSlice(f.Name, *values.slices[f.Flag])
		default:
			req.BodyString(f.Name, values.text(f.Flag))
		}
	}
}
