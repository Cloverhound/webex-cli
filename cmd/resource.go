package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/internal/output"
	"github.com/Cloverhound/webex-cli/pkg/catalog"
	"github.com/Cloverhound/webex-cli/pkg/changeset"
	"github.com/Cloverhound/webex-cli/pkg/guard"
	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webex"
	"github.com/Cloverhound/webex-cli/pkg/webexops"
	"github.com/spf13/cobra"
)

var resourceCmd = &cobra.Command{
	Use:     "resource",
	Aliases: []string{"resources"},
	Short:   "Read and safely change Webex objects by type",
	Long: `Read and change Webex objects by type, using the SDK catalog's knowledge of
each type: where its fields live, how Webex lists it, and how a write must be
sent.

Writes are safe where raw API commands are not. A Contact Center PUT replaces the
whole object, so "set" reads the object first and sends it back whole with only
the named fields changed. Every write is read back to confirm Webex holds the
result, and nothing is retried after a failure that may have committed.

  webex resource list-types
  webex resource describe cc_queue
  webex resource list cc_team
  webex resource get cc_team <id>
  webex resource set cc_queue <id> max_time_in_queue=900 recording_permitted=true
  webex resource create cc_skill name=zz-test skill_type=BOOLEAN
  webex resource delete cc_skill <id>

Objects of types with a parent, such as hunt groups in a location, are named
<parent>/<id>; a bare ID is looked up. Field values are text, true or false,
numbers, a comma-separated or JSON list, or JSON rows. --values takes them as
one JSON object instead, which suits values with spaces:

  webex resource set cc_team <id> --values '{"name":"Front Desk","capacity":4}'

With --changeset, set, create and delete are recorded in the change set instead
of being sent, like any other write command.`,
}

var resourceListTypesCmd = &cobra.Command{
	Use:   "list-types",
	Short: "List the object types",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		type typeInfo struct {
			Name      string `json:"name"`
			Title     string `json:"title"`
			Area      string `json:"area"`
			Writable  bool   `json:"writable"`
			CanCreate bool   `json:"canCreate"`
			CanDelete bool   `json:"canDelete"`
			Parent    string `json:"parent,omitempty"`
		}
		var out []typeInfo
		for _, t := range catalog.All() {
			out = append(out, typeInfo{Name: t.Name, Title: t.Title, Area: areaName(t), Writable: t.Managed,
				CanCreate: t.Managed && !t.NoCreate, CanDelete: t.Managed && !t.NoDelete, Parent: parentType(t)})
		}
		return printValue(out)
	},
}

var resourceDescribeCmd = &cobra.Command{
	Use:   "describe <type>",
	Short: "Show a type's fields and operations",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := catalog.ByName(args[0])
		if err != nil {
			return err
		}
		return printValue(describe(t))
	},
}

var resourceListCmd = &cobra.Command{
	Use:   "list <type>",
	Short: "List every object of a type",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := catalog.ByName(args[0])
		if err != nil {
			return err
		}
		objects, err := objectClient().List(cmd.Context(), t)
		if err != nil {
			return err
		}
		raw, _ := cmd.Flags().GetBool("raw")
		out := make([]objectView, 0, len(objects))
		for _, o := range objects {
			out = append(out, view(o, raw))
		}
		return printValue(out)
	},
}

var resourceGetCmd = &cobra.Command{
	Use:   "get <type> <id>",
	Short: "Read one object",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := catalog.ByName(args[0])
		if err != nil {
			return err
		}
		c := objectClient()
		key, err := objectKey(cmd.Context(), c, t, args[1])
		if err != nil {
			return err
		}
		o, err := c.Read(cmd.Context(), t, key)
		if err != nil {
			return readError(t, key, err)
		}
		raw, _ := cmd.Flags().GetBool("raw")
		return printValue(view(o, raw))
	},
}

var resourceSetCmd = &cobra.Command{
	Use:   "set <type> <id> [<field>=<value>...]",
	Short: "Change fields of one object, keeping the rest",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := writableType(args[0])
		if err != nil {
			return err
		}
		changes, err := fieldValues(cmd, t, args[2:], false)
		if err != nil {
			return err
		}
		ctx := cmd.Context()
		c := objectClient()
		key, err := objectKey(ctx, c, t, args[1])
		if err != nil {
			return err
		}
		recording := config.ChangesetPath() != ""
		if _, isNew := changeset.HandleOf(key.ID); isNew && recording {
			// The object exists only in the change set, which folds this into its create.
			return sendRaw(t.Routes.Update, key, t.PatchBody(changes))
		}
		fresh, err := c.Read(ctx, t, key)
		if err != nil {
			return readError(t, key, err)
		}
		planned := resource.Values{}
		for k, v := range fresh.Values {
			planned[k] = v
		}
		for k, v := range changes {
			planned[k] = v
		}
		changed := t.Unmet(fresh.Values, planned)
		// A field set back to its value in Webex still goes to the change set, which
		// withdraws the change it holds for that field.
		if len(changed) == 0 && !recording {
			return printValue(writeResult{Type: t.Name, Key: key, Label: fresh.Label, Changed: []string{}, Values: fresh.Values})
		}
		delta := resource.Values{}
		for _, name := range changed {
			delta[name] = planned[name]
		}
		body := t.PatchBody(delta)
		if t.Write == resource.WriteFull {
			body = t.FullBody(fresh.Raw, delta)
		}
		if recording || config.DryRun() {
			var sets []string
			for _, c := range changeset.Leaves(t.PatchBody(changes)) {
				sets = append(sets, c.Path)
			}
			return sendRaw(t.Routes.Update, key, body, sets...)
		}
		if err := refuseReadOnly(); err != nil {
			return err
		}
		values, err := writeGuard(c).Update(ctx, t, key, fresh.Values, planned)
		if err != nil {
			return err
		}
		return printValue(writeResult{Type: t.Name, Key: key, Label: t.Label(values), Changed: changed, Values: values})
	},
}

var resourceCreateCmd = &cobra.Command{
	Use:   "create <type> [<field>=<value>...]",
	Short: "Create an object",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := writableType(args[0])
		if err != nil {
			return err
		}
		if t.NoCreate {
			return fmt.Errorf("%s cannot be created through the API", t.Plural)
		}
		desired, err := fieldValues(cmd, t, args[1:], true)
		if err != nil {
			return err
		}
		for _, f := range t.Fields {
			if f.Required || f.Name == t.ParentField {
				if desired[f.Name] == nil {
					return fmt.Errorf("%s needs %s to be created", t.Title, f.Name)
				}
			}
		}
		if config.ChangesetPath() != "" || config.DryRun() {
			body, err := objectClient().CreateBody(cmd.Context(), t, desired)
			if err != nil {
				return err
			}
			parent, _ := desired[t.ParentField].(string)
			return sendRaw(t.Routes.Create, resource.Key{Parent: parent}, body)
		}
		if err := refuseReadOnly(); err != nil {
			return err
		}
		key, values, err := writeGuard(objectClient()).Create(cmd.Context(), t, desired)
		if err != nil {
			if key.ID != "" {
				return fmt.Errorf("%w (the object exists as %s)", err, key.ImportID())
			}
			return err
		}
		return printValue(writeResult{Type: t.Name, Key: key, Label: t.Label(values), Changed: t.FieldNames(), Values: values})
	},
}

var resourceDeleteCmd = &cobra.Command{
	Use:   "delete <type> <id>",
	Short: "Delete an object",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := writableType(args[0])
		if err != nil {
			return err
		}
		if t.NoDelete {
			return fmt.Errorf("%s cannot be deleted through the API", t.Plural)
		}
		ctx := cmd.Context()
		c := objectClient()
		key, err := objectKey(ctx, c, t, args[1])
		if err != nil {
			return err
		}
		if config.ChangesetPath() != "" || config.DryRun() {
			return sendRaw(t.Routes.Delete, key, nil)
		}
		if err := refuseReadOnly(); err != nil {
			return err
		}
		fresh, err := c.Read(ctx, t, key)
		if err != nil {
			return readError(t, key, err)
		}
		if err := writeGuard(c).Delete(ctx, t, key, fresh.Values); err != nil {
			return err
		}
		return printValue(map[string]any{"type": t.Name, "key": key, "label": fresh.Label, "deleted": true})
	},
}

func init() {
	for _, c := range []*cobra.Command{resourceListCmd, resourceGetCmd} {
		c.Flags().Bool("raw", false, "Include the object as Webex returns it")
	}
	for _, c := range []*cobra.Command{resourceSetCmd, resourceCreateCmd} {
		c.Flags().String("values", "", `Field values as a JSON object, such as {"name":"Front Desk","capacity":4}`)
	}
	resourceCmd.AddCommand(resourceListTypesCmd, resourceDescribeCmd, resourceListCmd, resourceGetCmd,
		resourceSetCmd, resourceCreateCmd, resourceDeleteCmd)
	rootCmd.AddCommand(resourceCmd)
}

// verifyDelay separates read-backs of a write, for Webex's read-after-write lag.
const verifyDelay = 400 * time.Millisecond

func objectClient() *webex.Client {
	return &webex.Client{API: client.SDK(), Hosts: config.Hosts(), OrgUUID: config.OrgID(), OrgID: config.OrgIDBase64()}
}

func writeGuard(c *webex.Client) *guard.Guard {
	return &guard.Guard{API: c, Writable: true, VerifyDelay: verifyDelay}
}

func refuseReadOnly() error {
	if config.ReadOnly() {
		return fmt.Errorf("%w: run 'webex login' to leave read-only mode, or record the change with --changeset", client.ErrReadOnly)
	}
	return nil
}

func areaName(t resource.Type) string {
	if t.Area == resource.ContactCenter {
		return "cc"
	}
	return "calling"
}

func parentType(t resource.Type) string {
	if f, ok := t.Field(t.ParentField); ok {
		return f.Ref
	}
	return ""
}

func writableType(name string) (resource.Type, error) {
	t, err := catalog.ByName(name)
	if err != nil {
		return t, err
	}
	if !t.Managed {
		return t, fmt.Errorf("%s are read-only: their writes have not been proven on a live tenant", t.Plural)
	}
	return t, nil
}

// objectKey parses <parent>/<id>, or finds the parent of a bare ID by listing.
func objectKey(ctx context.Context, c *webex.Client, t resource.Type, s string) (resource.Key, error) {
	if t.ParentField == "" || strings.Contains(s, "/") {
		return t.ParseImportID(s)
	}
	keys, err := c.ListKeys(ctx, t)
	if err != nil {
		return resource.Key{}, err
	}
	for _, k := range keys {
		if k.ID == s {
			return k, nil
		}
	}
	return resource.Key{}, fmt.Errorf("no %s has ID %s", t.Title, s)
}

func readError(t resource.Type, key resource.Key, err error) error {
	if errors.Is(err, webex.ErrNotFound) {
		return fmt.Errorf("%s %s not found", t.Title, key.ImportID())
	}
	return err
}

// fieldValues reads field=value arguments and the --values JSON object. Immutable
// fields can only be given on create, because changing one means deleting and
// recreating the object.
func fieldValues(cmd *cobra.Command, t resource.Type, args []string, create bool) (resource.Values, error) {
	field := func(name string) (resource.Field, error) {
		f, ok := t.Field(name)
		if !ok {
			return f, fmt.Errorf("%s has no field %q (see: webex resource describe %s)", t.Title, name, t.Name)
		}
		if f.Immutable && !create {
			return f, fmt.Errorf("%s.%s cannot be changed after the object is created", t.Name, name)
		}
		return f, nil
	}
	out := resource.Values{}
	if raw, _ := cmd.Flags().GetString("values"); raw != "" {
		var given map[string]any
		if err := json.Unmarshal([]byte(raw), &given); err != nil {
			return nil, fmt.Errorf("--values must be a JSON object: %w", err)
		}
		for name, v := range given {
			f, err := field(name)
			if err != nil {
				return nil, err
			}
			if out[name], err = jsonValue(f, v); err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
		}
	}
	for _, a := range args {
		name, text, ok := strings.Cut(a, "=")
		if !ok {
			return nil, fmt.Errorf("%q is not field=value", a)
		}
		f, err := field(name)
		if err != nil {
			return nil, err
		}
		if out[name], err = parseValue(f, text); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
	}
	if len(out) == 0 {
		return nil, errors.New("name at least one field, as field=value or with --values")
	}
	return out, nil
}

// jsonValue checks a --values entry against its field. Lists and rows are checked
// as text, the same way as field=value.
func jsonValue(f resource.Field, v any) (any, error) {
	switch f.Kind {
	case resource.Bool:
		if b, ok := v.(bool); ok {
			return b, nil
		}
		return nil, errors.New("want true or false")
	case resource.Int:
		if n, ok := v.(float64); ok && n == float64(int64(n)) {
			return int64(n), nil
		}
		return nil, errors.New("want a whole number")
	case resource.StringList, resource.ObjectList:
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return parseValue(f, string(data))
	default:
		s, ok := v.(string)
		if !ok {
			return nil, errors.New("want text")
		}
		return parseValue(f, s)
	}
}

func parseValue(f resource.Field, text string) (any, error) {
	switch f.Kind {
	case resource.Bool:
		b, err := strconv.ParseBool(text)
		if err != nil {
			return nil, fmt.Errorf("want true or false, not %q", text)
		}
		return b, nil
	case resource.Int:
		n, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("want a whole number, not %q", text)
		}
		return n, nil
	case resource.StringList:
		var list []any
		if strings.HasPrefix(strings.TrimSpace(text), "[") {
			if err := json.Unmarshal([]byte(text), &list); err != nil {
				return nil, fmt.Errorf("want a JSON list: %w", err)
			}
		} else {
			list = []any{}
			for _, s := range strings.Split(text, ",") {
				if s = strings.TrimSpace(s); s != "" {
					list = append(list, s)
				}
			}
		}
		for _, e := range list {
			s, ok := e.(string)
			if !ok {
				return nil, fmt.Errorf("list entries must be text")
			}
			if err := checkOption(f, s); err != nil {
				return nil, err
			}
		}
		return f.Normalize(list), nil
	case resource.ObjectList:
		var rows []any
		if err := json.Unmarshal([]byte(text), &rows); err != nil {
			return nil, fmt.Errorf("want JSON rows such as [{\"%s\":...}]: %w", firstItem(f), err)
		}
		for _, r := range rows {
			row, ok := r.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("each row must be a JSON object")
			}
			for k := range row {
				if _, ok := (resource.Type{Fields: f.Item}).Field(k); !ok {
					return nil, fmt.Errorf("rows have no field %q", k)
				}
			}
		}
		return f.Normalize(rows), nil
	default:
		if err := checkOption(f, text); err != nil {
			return nil, err
		}
		return text, nil
	}
}

func checkOption(f resource.Field, s string) error {
	if len(f.Options) > 0 && s != "" && !slices.Contains(f.Options, s) {
		return fmt.Errorf("%q is not one of %s", s, strings.Join(f.Options, ", "))
	}
	return nil
}

func firstItem(f resource.Field) string {
	if len(f.Item) > 0 {
		return f.Item[0].Name
	}
	return "field"
}

// sendRaw sends one write through the CLI's request path, so --changeset records
// it and --dry-run prints it like any other write command. sets names the body
// paths the write means to set (see client.Request.Sets).
func sendRaw(r resource.Route, key resource.Key, body map[string]any, sets ...string) error {
	req := client.NewRequest(r.Op.Host.URL(config.Hosts()), r.Op.Method, r.Op.Path)
	req.Sets(sets...)
	for _, p := range r.Op.Params {
		switch {
		case p.Value != "":
			req.PathParam(p.Name, p.Value)
		case p.Name == r.Parent:
			req.PathParam(p.Name, key.Parent)
		case p.Name == r.ID:
			req.PathParam(p.Name, key.ID)
		case p.Name == "orgid":
			req.PathParam(p.Name, config.OrgID())
		case p.Name == "orgId" && p.In == webexops.InQuery:
			req.QueryParam(p.Name, config.OrgIDBase64())
		}
	}
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		req.SetBodyRaw(string(data))
	}
	resp, status, err := req.Do()
	if err != nil {
		return err
	}
	return output.Print(resp, status)
}

type objectView struct {
	Type   string          `json:"type"`
	Key    resource.Key    `json:"key"`
	Label  string          `json:"label"`
	Values resource.Values `json:"values"`
	Raw    map[string]any  `json:"raw,omitempty"`
}

func view(o webex.Object, raw bool) objectView {
	v := objectView{Type: o.Type, Key: o.Key, Label: o.Label, Values: o.Values}
	if raw {
		v.Raw = o.Raw
	}
	return v
}

type writeResult struct {
	Type    string          `json:"type"`
	Key     resource.Key    `json:"key"`
	Label   string          `json:"label"`
	Changed []string        `json:"changed"`
	Values  resource.Values `json:"values"`
}

type fieldDescription struct {
	Name         string             `json:"name"`
	Path         string             `json:"path,omitempty"`
	Kind         string             `json:"kind"`
	Required     bool               `json:"required,omitempty"`
	Immutable    bool               `json:"immutable,omitempty"`
	Ref          string             `json:"ref,omitempty"`
	Options      []string           `json:"options,omitempty"`
	OptionSource string             `json:"optionSource,omitempty"`
	Format       string             `json:"format,omitempty"`
	Needs        string             `json:"needs,omitempty"`
	Item         []fieldDescription `json:"item,omitempty"`
}

func describeField(f resource.Field) fieldDescription {
	kinds := map[resource.Kind]string{resource.String: "string", resource.Bool: "bool", resource.Int: "int",
		resource.StringList: "list", resource.ObjectList: "rows"}
	d := fieldDescription{Name: f.Name, Path: f.Path, Kind: kinds[f.Kind], Required: f.Required, Immutable: f.Immutable,
		Ref: f.Ref, Options: f.Options, OptionSource: f.OptionSource, Format: f.Format, Needs: f.Needs}
	for _, sub := range f.Item {
		d.Item = append(d.Item, describeField(sub))
	}
	return d
}

func describe(t resource.Type) map[string]any {
	ops := map[string]string{}
	for use, r := range map[string]resource.Route{"list": t.Routes.List, "create": t.Routes.Create, "read": t.Routes.Read,
		"update": t.Routes.Update, "delete": t.Routes.Delete} {
		if r.Defined() {
			ops[use] = "webex " + r.Op.Name()
		}
	}
	fields := []fieldDescription{}
	for _, f := range t.Fields {
		fields = append(fields, describeField(f))
	}
	write := map[resource.WriteMode]string{resource.WriteDelta: "changed fields", resource.WriteFull: "whole object",
		resource.WritePatch: "changed fields, by PATCH"}
	return map[string]any{
		"name": t.Name, "title": t.Title, "area": areaName(t), "writable": t.Managed,
		"canCreate": t.Managed && !t.NoCreate, "canDelete": t.Managed && !t.NoDelete,
		"parent": parentType(t), "labelField": t.LabelField, "writes": write[t.Write],
		"operations": ops, "fields": fields,
	}
}

func printValue(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return output.Print(data, 200)
}
