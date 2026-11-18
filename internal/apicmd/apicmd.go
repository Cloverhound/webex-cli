// Package apicmd builds the CLI's API commands from the operation table in
// pkg/webexops. Each operation becomes `webex <area> <group> <command>`, with a
// flag per parameter and body field.
package apicmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/aiagent"
	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/Cloverhound/webex-cli/internal/timeutil"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
	"github.com/spf13/cobra"
)

// privateAPINote ends the help of commands on an API Cisco does not support.
const privateAPINote = "Caution: this command uses Webex AI Agent Studio's private API, which is not " +
	"documented or supported by Cisco and can change or stop working without notice. " +
	"Do not build production automation on it."

// Register adds an area's groups and commands under parent and returns the group
// commands by name, for handwritten commands to join.
func Register(parent *cobra.Command, area string) map[string]*cobra.Command {
	groups := map[string]*cobra.Command{}
	for _, g := range webexops.Groups(area) {
		gc := &cobra.Command{Use: g.Name, Aliases: g.Aliases, Short: g.Summary}
		parent.AddCommand(gc)
		groups[g.Name] = gc
	}
	for _, op := range webexops.Operations(area) {
		groups[op.Group].AddCommand(NewCommand(op))
	}
	return groups
}

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
			description += "\n\n" + privateAPINote
		} else {
			description = privateAPINote
		}
	}

	c := &cobra.Command{
		Use:     op.Command,
		Aliases: op.Aliases,
		Short:   op.Summary,
		Long:    description,
		RunE: func(cmd *cobra.Command, args []string) error {
			if op.Private {
				aiagent.PrivateAPINotice()
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

			call := webexops.Args{Params: map[string]string{}}
			for _, p := range op.Params {
				switch {
				case p.Value != "":
				case p.In == webexops.InAgentToken:
					token, err := aiagent.BotToken(values.text(p.Flag))
					if err != nil {
						return err
					}
					call.Params[p.Name] = token
				default:
					call.Params[p.Name] = values.text(p.Flag)
				}
			}
			if op.Body != nil {
				switch {
				case bodyFile != "":
					data, err := os.ReadFile(bodyFile)
					if err != nil {
						return fmt.Errorf("reading body file: %w", err)
					}
					call.Body = data
				case bodyRaw != "":
					call.Body = []byte(bodyRaw)
				default:
					body, err := bodyFromFlags(cmd, op.Body.Fields, values)
					if err != nil {
						return err
					}
					call.Body = body
				}
			}

			// The CLI sends an empty path flag as it always has, so it goes to the
			// client directly rather than through webexops.Call and All.
			api := client.API(op.Path, string(call.Body))
			ctx := context.Background()
			if op.Paging != nil && config.Paginate() {
				items, err := api.All(ctx, op.List(config.Hosts(), call), *op.Paging)
				if err != nil {
					return client.CLIError(err)
				}
				resp, err := json.Marshal(items)
				if err != nil {
					return err
				}
				return show(op, resp, http.StatusOK)
			}
			resp, err := api.Send(ctx, op.Request(config.Hosts(), call))
			if err != nil {
				return client.CLIError(err)
			}
			return show(op, resp.Body, resp.Status)
		},
	}

	flags := c.Flags()
	for _, p := range op.Params {
		if p.Value != "" || values.strings[p.Flag] != nil {
			continue
		}
		values.strings[p.Flag] = new(string)
		flags.StringVar(values.strings[p.Flag], p.Flag, "", strings.ReplaceAll(p.Usage, "\n", " "))
		if p.Required {
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

// bodyFromFlags is the JSON body of the fields given by flags, or none when no
// field is given. A string or list field is sent when it is not empty, and any
// other field when its flag is set.
func bodyFromFlags(cmd *cobra.Command, fields []webexops.BodyField, values *flagValues) ([]byte, error) {
	body := map[string]any{}
	for _, f := range fields {
		changed := cmd.Flags().Changed(f.Flag)
		switch f.Type {
		case "int64":
			if changed {
				body[f.Name] = *values.ints[f.Flag]
			}
		case "float64":
			if changed {
				body[f.Name] = *values.floats[f.Flag]
			}
		case "bool":
			if changed {
				body[f.Name] = *values.bools[f.Flag]
			}
		case "[]string":
			if v := *values.slices[f.Flag]; len(v) > 0 {
				body[f.Name] = v
			}
		default:
			if v := values.text(f.Flag); v != "" {
				body[f.Name] = v
			}
		}
	}
	if len(body) == 0 {
		return nil, nil
	}
	return json.Marshal(body)
}

// show prints a response, without the AI agent credentials some carry.
func show(op *webexops.Operation, body []byte, status int) error {
	if op.Redact {
		body = aiagent.RedactSecrets(body)
	}
	return output.Print(body, status)
}
