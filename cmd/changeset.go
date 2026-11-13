package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/Cloverhound/webex-cli/internal/client"
	"github.com/Cloverhound/webex-cli/internal/config"
	"github.com/Cloverhound/webex-cli/pkg/catalog"
	"github.com/Cloverhound/webex-cli/pkg/changeset"
	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/spf13/cobra"
)

var changesetCmd = &cobra.Command{
	Use:   "changeset",
	Short: "Propose changes for Webex Admin Studio in a change set file (beta)",
	Long: `Beta. A change set file holds proposed Webex changes instead of sending them.
Webex Admin Studio imports it for review and apply.

Record changes by running ordinary write commands with --changeset <file>
(or $WEBEX_CHANGESET). Nothing is sent to Webex, so this also works in
read-only mode. Only creates, updates and deletes of the types Studio applies
can be recorded; 'webex changeset types' lists them and their commands.

With --changeset, an update body describes proposed changes. Send only the
fields you intend to change: omitted fields stay as they are, and lists are
replaced whole. Values equal to Webex add no change, and fields Studio does
not model are ignored. Do not apply a proposal by removing --changeset: a
direct Contact Center update needs the complete object.

A create returns a placeholder ID such as "new:boulder-desk". Use it in later
commands in this change set, in flags or bodies, to refer to the new object.`,
}

func init() {
	newCmd := &cobra.Command{
		Use:   "new <file>",
		Short: "Start a change set for the current org",
		Args:  cobra.ExactArgs(1),
		RunE:  runChangesetNew,
	}
	newCmd.Flags().String("title", "", "Short title")
	newCmd.Flags().String("rationale", "", "Why these changes are needed")
	newCmd.Flags().String("author", "", "Who or what is proposing the changes")
	newCmd.Flags().String("author-type", "agent", "agent or human")
	newCmd.Flags().String("on-behalf-of", "", "The person an agent is working for")
	newCmd.Flags().StringArray("link", nil, "A link as name=url, such as ticket=https://...")

	showCmd := &cobra.Command{
		Use:   "show <file>",
		Short: "Summarize a change set",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cs, err := changeset.Load(args[0])
			if err != nil {
				return err
			}
			showChangeSet(cmd.OutOrStdout(), cs)
			return nil
		},
	}

	validateCmd := &cobra.Command{
		Use:   "validate <file>",
		Short: "Check a change set, and with --live compare its preconditions with Webex",
		Args:  cobra.ExactArgs(1),
		RunE:  runChangesetValidate,
	}
	validateCmd.Flags().Bool("live", false, "Re-read every updated or deleted object and report drift")

	intentCmd := &cobra.Command{
		Use:   "intent <file> <intent-id> <summary> [op-id...]",
		Short: "Add or reword an intent and assign operations to it",
		Args:  cobra.MinimumNArgs(3),
		RunE:  runChangesetIntent,
	}

	dropCmd := &cobra.Command{
		Use:   "drop <file> <op-id> [<path>...]",
		Short: "Remove an operation, or some field changes from an update",
		Long: `Remove an operation from a change set, which reverts the create, update or
delete it proposes. With paths, remove only those field changes from an update;
an update with none left is removed.

A create that another operation refers to by its new: ID cannot be dropped
until that operation is dropped too.`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return editChangeSet(args[0], func(cs *changeset.ChangeSet) error {
				return cs.Drop(args[1], args[2:]...)
			})
		},
	}

	metaCmd := &cobra.Command{
		Use:   "meta <file>",
		Short: "Change a change set's title, rationale, author or links",
		Long: `Change the details a change set was started with. Only the flags given change;
an empty value clears one, and --link name= removes that link.`,
		Args: cobra.ExactArgs(1),
		RunE: runChangesetMeta,
	}
	metaCmd.Flags().String("title", "", "Short title")
	metaCmd.Flags().String("rationale", "", "Why these changes are needed")
	metaCmd.Flags().String("author", "", "Who or what is proposing the changes")
	metaCmd.Flags().String("author-type", "", "agent or human")
	metaCmd.Flags().String("on-behalf-of", "", "The person an agent is working for")
	metaCmd.Flags().StringArray("link", nil, "A link as name=url; name= removes it")

	typesCmd := &cobra.Command{
		Use:   "types [<type>]",
		Short: "List the object types a change set can record",
		Long: `List the object types Webex Admin Studio can apply from a change set, with the
commands that record each one's creates, updates and deletes. Given a type, list
its fields: the JSON paths an update records, which ones a create needs, and which
are fixed once the object exists. Paths outside these are ignored.`,
		Args: cobra.MaximumNArgs(1),
		RunE: runChangesetTypes,
	}

	schemaCmd := &cobra.Command{
		Use:   "schema",
		Short: "Print the change set JSON Schema",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := cmd.OutOrStdout().Write(changeset.Schema)
			return err
		},
	}

	changesetCmd.AddCommand(newCmd, showCmd, validateCmd, intentCmd, dropCmd, metaCmd, typesCmd, schemaCmd)
	rootCmd.AddCommand(changesetCmd)
}

func runChangesetNew(cmd *cobra.Command, args []string) error {
	path := args[0]
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", path)
	}
	if config.OrgID() == "" {
		return errors.New("no organization: log in or pass --organization")
	}
	cs := changeset.New(config.OrgID())
	cs.Metadata.Title, _ = cmd.Flags().GetString("title")
	cs.Metadata.Rationale, _ = cmd.Flags().GetString("rationale")
	author, _ := cmd.Flags().GetString("author")
	authorType, _ := cmd.Flags().GetString("author-type")
	onBehalfOf, _ := cmd.Flags().GetString("on-behalf-of")
	if author != "" || onBehalfOf != "" {
		cs.Metadata.Author = &changeset.Author{Type: authorType, ID: author, OnBehalfOf: onBehalfOf}
	}
	links, _ := cmd.Flags().GetStringArray("link")
	for _, link := range links {
		name, url, ok := strings.Cut(link, "=")
		if !ok || name == "" || url == "" {
			return fmt.Errorf("--link %q: expected name=url", link)
		}
		if cs.Metadata.Links == nil {
			cs.Metadata.Links = map[string]string{}
		}
		cs.Metadata.Links[name] = url
	}
	if err := changeset.Save(path, cs); err != nil {
		return err
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "Started %s for org %s. Record changes with --changeset %s\n", path, cs.Scope.OrgID, path)
	fmt.Fprintln(cmd.ErrOrStderr(), "Change sets are in beta: Studio may refuse operations it cannot apply.")
	return nil
}

// editChangeSet changes an existing change set file under its lock.
func editChangeSet(path string, change func(*changeset.ChangeSet) error) error {
	return changeset.Edit(path, func() (*changeset.ChangeSet, error) {
		return nil, fmt.Errorf("%s does not exist", path)
	}, change)
}

func runChangesetMeta(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()
	links, _ := flags.GetStringArray("link")
	for _, link := range links {
		if name, _, ok := strings.Cut(link, "="); !ok || name == "" {
			return fmt.Errorf("--link %q: expected name=url, or name= to remove it", link)
		}
	}
	return editChangeSet(args[0], func(cs *changeset.ChangeSet) error {
		m := &cs.Metadata
		if flags.Changed("title") {
			m.Title, _ = flags.GetString("title")
		}
		if flags.Changed("rationale") {
			m.Rationale, _ = flags.GetString("rationale")
		}
		if flags.Changed("author") || flags.Changed("author-type") || flags.Changed("on-behalf-of") {
			a := changeset.Author{Type: "agent"}
			if m.Author != nil {
				a = *m.Author
			}
			if flags.Changed("author") {
				a.ID, _ = flags.GetString("author")
			}
			if flags.Changed("author-type") {
				a.Type, _ = flags.GetString("author-type")
			}
			if flags.Changed("on-behalf-of") {
				a.OnBehalfOf, _ = flags.GetString("on-behalf-of")
			}
			m.Author = &a
			if a.ID == "" && a.OnBehalfOf == "" {
				m.Author = nil
			}
		}
		for _, link := range links {
			name, url, _ := strings.Cut(link, "=")
			if url == "" {
				delete(m.Links, name)
				continue
			}
			if m.Links == nil {
				m.Links = map[string]string{}
			}
			m.Links[name] = url
		}
		if len(m.Links) == 0 {
			m.Links = nil
		}
		return nil
	})
}

type recordableType struct {
	Name      string            `json:"name"`
	Title     string            `json:"title"`
	WriteMode string            `json:"writeMode"`
	Commands  map[string]string `json:"commands"`
	Fields    []recordableField `json:"fields,omitempty"`
}

type recordableField struct {
	Path     string            `json:"path"`
	Kind     string            `json:"kind"`
	Required bool              `json:"required,omitempty"`
	Fixed    bool              `json:"fixed,omitempty"`
	Options  []string          `json:"options,omitempty"`
	Item     []recordableField `json:"item,omitempty"`
}

func runChangesetTypes(cmd *cobra.Command, args []string) error {
	var out []recordableType
	for _, t := range catalog.Managed() {
		if len(args) == 1 && t.Name != args[0] {
			continue
		}
		rt := recordableType{Name: t.Name, Title: t.Title, WriteMode: writeModes[t.Write], Commands: map[string]string{}}
		for action, r := range map[string]resource.Route{"create": t.Routes.Create, "update": t.Routes.Update, "delete": t.Routes.Delete} {
			if r.Op == nil || (action == "create" && t.NoCreate) || (action == "delete" && t.NoDelete) {
				continue
			}
			rt.Commands[action] = "webex " + r.Op.Name()
		}
		if len(args) == 1 {
			rt.Fields = recordableFields(t.Fields)
		}
		out = append(out, rt)
	}
	if len(args) == 1 && len(out) == 0 {
		return fmt.Errorf("%s is not a type a change set can record; run 'webex changeset types' for the list", args[0])
	}
	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	if len(args) == 1 {
		return enc.Encode(out[0])
	}
	return enc.Encode(out)
}

var writeModes = map[resource.WriteMode]string{
	resource.WriteDelta: "changed fields",
	resource.WriteFull:  "whole object",
	resource.WritePatch: "changed fields, by PATCH",
}

var fieldKinds = map[resource.Kind]string{resource.String: "string", resource.Bool: "bool", resource.Int: "int",
	resource.StringList: "list", resource.ObjectList: "rows"}

func recordableFields(fields []resource.Field) []recordableField {
	var out []recordableField
	for _, f := range fields {
		if f.Path == "" {
			continue
		}
		out = append(out, recordableField{Path: f.Path, Kind: fieldKinds[f.Kind], Required: f.Required,
			Fixed: f.Immutable, Options: f.Options, Item: recordableFields(f.Item)})
	}
	return out
}

func runChangesetIntent(cmd *cobra.Command, args []string) error {
	path, id, summary, ops := args[0], args[1], args[2], args[3:]
	return editChangeSet(path, func(cs *changeset.ChangeSet) error {
		found := false
		for i := range cs.Intents {
			if cs.Intents[i].ID == id {
				cs.Intents[i].Summary, found = summary, true
			}
		}
		if !found {
			cs.Intents = append(cs.Intents, changeset.Intent{ID: id, Summary: summary})
		}
		for _, opID := range ops {
			assigned := false
			for i := range cs.Operations {
				if cs.Operations[i].ID == opID {
					cs.Operations[i].Intent, assigned = id, true
				}
			}
			if !assigned {
				return fmt.Errorf("no operation %s", opID)
			}
		}
		return nil
	})
}

type validateResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
	*client.LiveReport
}

func runChangesetValidate(cmd *cobra.Command, args []string) error {
	cs, err := changeset.Load(args[0])
	if err != nil {
		return err
	}
	problems := changeset.Validate(cs)
	result := validateResult{Errors: problems.Errors, Warnings: problems.Warnings}

	if live, _ := cmd.Flags().GetBool("live"); live && len(problems.Errors) == 0 {
		if cs.Scope.OrgID != config.OrgID() {
			return fmt.Errorf("the change set is for org %s, but you are using org %s", cs.Scope.OrgID, config.OrgID())
		}
		report := client.CheckLive(cs)
		result.LiveReport = &report
	}
	result.Valid = len(result.Errors) == 0 && (result.LiveReport == nil || result.LiveReport.OK())

	enc := json.NewEncoder(cmd.OutOrStdout())
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		return err
	}
	if !result.Valid {
		return errors.New("the change set is not valid")
	}
	return nil
}

func showChangeSet(w io.Writer, cs *changeset.ChangeSet) {
	title := cs.Metadata.Title
	if title == "" {
		title = "(untitled change set)"
	}
	fmt.Fprintf(w, "%s\n%s · org %s", title, cs.Metadata.ID, cs.Scope.OrgID)
	if cs.Metadata.CreatedAt != "" {
		fmt.Fprintf(w, " · %s", cs.Metadata.CreatedAt)
	}
	fmt.Fprintln(w)
	if a := cs.Metadata.Author; a != nil {
		fmt.Fprintf(w, "Author: %s %s", a.Type, a.ID)
		if a.OnBehalfOf != "" {
			fmt.Fprintf(w, " for %s", a.OnBehalfOf)
		}
		fmt.Fprintln(w)
	}
	if cs.Metadata.Rationale != "" {
		fmt.Fprintf(w, "Rationale: %s\n", cs.Metadata.Rationale)
	}
	links := make([]string, 0, len(cs.Metadata.Links))
	for name, url := range cs.Metadata.Links {
		links = append(links, name+": "+url)
	}
	sort.Strings(links)
	for _, l := range links {
		fmt.Fprintf(w, "Link %s\n", l)
	}

	groups := []changeset.Intent{{}}
	groups = append(groups, cs.Intents...)
	for _, intent := range groups {
		var ops []changeset.Operation
		for _, op := range cs.Operations {
			if op.Intent == intent.ID {
				ops = append(ops, op)
			}
		}
		if len(ops) == 0 {
			continue
		}
		fmt.Fprintln(w)
		if intent.ID != "" {
			fmt.Fprintf(w, "[%s] %s\n", intent.ID, intent.Summary)
		}
		for _, op := range ops {
			showOperation(w, op)
		}
	}
	if len(cs.Operations) == 0 {
		fmt.Fprintln(w, "\nNo operations yet.")
	}
}

func showOperation(w io.Writer, op changeset.Operation) {
	label := op.Target.Display
	if label == "" {
		label = op.Target.ID()
	}
	action := string(op.Action)
	if op.Action == changeset.Invoke {
		action = "invoke " + op.Method
	}
	fmt.Fprintf(w, "  %s %s %s %s", op.ID, action, op.Target.API, op.Target.Path())
	if label != "" {
		fmt.Fprintf(w, "  %q", label)
	}
	if op.Handle != "" {
		fmt.Fprintf(w, "  as new:%s", op.Handle)
	}
	fmt.Fprintln(w)
	switch op.Action {
	case changeset.Update:
		for _, c := range op.Changes {
			if c.HasBefore {
				fmt.Fprintf(w, "      %s: %s → %s\n", c.Path, compact(c.Before), compact(c.Value))
			} else {
				fmt.Fprintf(w, "      %s = %s\n", c.Path, compact(c.Value))
			}
		}
	case changeset.Create, changeset.Invoke:
		if body, ok := op.Body.(map[string]any); ok {
			for _, c := range changeset.Leaves(body) {
				fmt.Fprintf(w, "      %s = %s\n", c.Path, compact(c.Value))
			}
		} else if op.Body != nil {
			fmt.Fprintf(w, "      body = %s\n", compact(op.Body))
		}
	}
}

func compact(v any) string {
	data, err := json.Marshal(placeholderIDs(v))
	if err != nil {
		return fmt.Sprint(v)
	}
	s := string(data)
	if len(s) > 120 {
		s = s[:117] + "..."
	}
	return s
}

// placeholderIDs shows references the way commands write them, as new:<handle>.
func placeholderIDs(v any) any {
	if h, ok := changeset.Ref(v); ok {
		return changeset.NewIDPrefix + h
	}
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, e := range t {
			out[k] = placeholderIDs(e)
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, e := range t {
			out[i] = placeholderIDs(e)
		}
		return out
	}
	return v
}
