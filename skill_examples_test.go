package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/skill"
	"github.com/spf13/cobra"
)

// Agents run the skill's examples as written, and a collection refresh can
// rename a command or flag under them. Every `webex` command in a code block of a
// skill file or a documentation page must name a real command and only flags
// that command takes. The generated Command Reference is skipped: codegen
// writes it from the commands.
func TestDocumentedExamplesUseRealCommandsAndFlags(t *testing.T) {
	files, err := skill.Files()
	if err != nil {
		t.Fatal(err)
	}
	docs := map[string]string{"skill/SETUP.md": skill.Setup}
	for name, content := range files {
		docs["skill/"+name] = string(content)
	}
	pages, err := filepath.Glob("site/src/content/docs/*.mdx")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range append(pages, "README.md") {
		content, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		docs[name] = string(content)
	}
	names := make([]string, 0, len(docs))
	for name := range docs {
		names = append(names, name)
	}
	sort.Strings(names)

	root := cmd.RootCommand()
	for _, name := range names {
		for _, problem := range documentProblems(root, name, docs[name]) {
			t.Error(problem)
		}
	}
}

// The checker must itself catch each kind of stale example.
func TestDocumentCheckFindsStaleExamples(t *testing.T) {
	root := cmd.RootCommand()
	for _, tc := range []struct{ name, doc, want string }{
		{"table row", "| `site` | list, nosuch |", `"webex cc site" has no command "nosuch"`},
		{"leading global flag", "```bash\nwebex --token \"$T\" calling people list --nosuch\n```", "has no flag --nosuch"},
		{"argument placeholder", "```bash\nwebex auth switch <email> --nosuch\n```", "has no flag --nosuch"},
		{"after &&", "```bash\nmkdir -p x && webex post-install -y --nosuch\n```", "has no flag --nosuch"},
		{"after nohup", "```bash\nnohup webex login --device --nosuch > log 2>&1 &\n```", "has no flag --nosuch"},
		{"flag prefix", "```bash\nwebex calling people list --max_results 5\n```", "--max_results"},
		{"unknown shorthand", "```bash\nwebex post-install -z\n```", "has no flag -z"},
		{"continuation line", "```bash\nwebex calling people list\n  --nosuch <x>\n```", "has no flag --nosuch"},
	} {
		problems := documentProblems(root, "skill/cc/SKILL.md", tc.doc)
		if !slices.ContainsFunc(problems, func(p string) bool { return strings.Contains(p, tc.want) }) {
			t.Errorf("%s: problems %q, want one containing %q", tc.name, problems, tc.want)
		}
	}
	for _, ok := range []string{
		"```bash\nwebex cc <resource> list --filter x\n```",
		"```bash\nmkdir -p x && webex post-install -y\n```",
		"```bash\nwebex --token \"$T\" calling people list --max 5\n```",
		"```bash\nwebex auth switch user@example.com\n```",
		"```bash\nwebex cc site list --help\n```",
		"```bash\nwebex cc search tasks --from -24h\n```",
	} {
		if problems := documentProblems(root, "skill/cc/SKILL.md", ok); len(problems) > 0 {
			t.Errorf("%q: unexpected problems %q", ok, problems)
		}
	}
}

// documentProblems checks a document's examples and, for an area skill, its
// table of groups and their commands.
func documentProblems(root *cobra.Command, name, content string) []string {
	var problems []string
	for _, ex := range skillExamples(root, content) {
		for _, p := range checkExample(root, ex) {
			problems = append(problems, fmt.Sprintf("%s:%d: %s\n\t%s", name, ex.line, p, ex.text))
		}
	}
	area, file, ok := strings.Cut(strings.TrimPrefix(name, "skill/"), "/")
	if !ok || file != "SKILL.md" {
		return problems
	}
	for _, row := range commandTableRows(content) {
		for _, command := range row.commands {
			for _, p := range checkExample(root, skillExample{words: []string{area, row.group, command}}) {
				problems = append(problems, fmt.Sprintf("%s:%d: %s\n\t%s", name, row.line, p, row.text))
			}
		}
	}
	return problems
}

// commandTableRow is a row of an area skill's table of groups and their
// commands: | `group` | command, command |
type commandTableRow struct {
	line     int
	text     string
	group    string
	commands []string
}

var commandTableRowPattern = regexp.MustCompile("^\\| `([a-z0-9-]+)` \\| ([a-z0-9-]+(?:, [a-z0-9-]+)*) \\|$")

func commandTableRows(content string) []commandTableRow {
	var rows []commandTableRow
	for i, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "<!-- codegen:start") {
			break
		}
		if m := commandTableRowPattern.FindStringSubmatch(line); m != nil {
			rows = append(rows, commandTableRow{line: i + 1, text: line, group: m[1], commands: strings.Split(m[2], ", ")})
		}
	}
	return rows
}

// skillExample is one `webex` command in a document, with the flags on its
// continuation lines.
type skillExample struct {
	line  int
	text  string
	words []string // command words after "webex" and any global flags
	// generic is set when a placeholder follows the words, as in
	// "webex cc <resource> list": it stands for many commands where the words
	// name a group, and is an argument where they name a command.
	generic bool
	flags   []string // flags as written, "--name" or "-x", without values
}

var commandWord = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)

var (
	// shellOperators end one command; the word after one is a command.
	shellOperators = map[string]bool{"&&": true, "||": true, ";": true, "|": true, "&": true}
	// commandPrefixes run the command that follows them.
	commandPrefixes = map[string]bool{"nohup": true, "sudo": true, "exec": true}
	redirection     = regexp.MustCompile(`^(\d?>|<$)`)
)

// skillExamples finds the `webex` commands in a document's fenced code blocks,
// outside the generated block. A line in the same block that starts with a flag
// continues the command above it.
func skillExamples(root *cobra.Command, content string) []skillExample {
	var out []skillExample
	inCode, generated := false, false
	last := -1
	for i, raw := range strings.Split(content, "\n") {
		line := strings.TrimSpace(raw)
		switch {
		case strings.HasPrefix(line, "<!-- codegen:start"):
			generated = true
		case strings.HasPrefix(line, "<!-- codegen:end"):
			generated = false
		case strings.HasPrefix(line, "```"):
			inCode, last = !inCode, -1
		}
		if !inCode || generated {
			continue
		}
		tokens := shellWords(line)
		if len(tokens) > 0 && strings.HasPrefix(tokens[0], "-") && last >= 0 {
			out[last].flags = append(out[last].flags, flagsIn(shellCommand(tokens))...)
			continue
		}
		if line == "" {
			last = -1
		}
		for start, tok := range tokens {
			if tok != "webex" || start > 0 && !shellOperators[tokens[start-1]] && !commandPrefixes[tokens[start-1]] {
				continue
			}
			ex := skillExample{line: i + 1, text: line}
			args := shellCommand(tokens[start+1:])
			for len(args) > 0 && strings.HasPrefix(args[0], "-") {
				ex.flags = append(ex.flags, flagsIn(args[:1])...)
				name, _, hasValue := strings.Cut(strings.TrimLeft(args[0], "-"), "=")
				f := root.PersistentFlags().Lookup(name)
				args = args[1:]
				if f != nil && f.NoOptDefVal == "" && !hasValue && len(args) > 0 {
					args = args[1:]
				}
			}
			for len(args) > 0 && commandWord.MatchString(args[0]) {
				ex.words, args = append(ex.words, args[0]), args[1:]
			}
			if len(ex.words) == 0 {
				continue
			}
			ex.generic = len(args) > 0 && strings.HasPrefix(args[0], "<")
			ex.flags = append(ex.flags, flagsIn(args)...)
			out = append(out, ex)
			last = len(out) - 1
		}
	}
	return out
}

// shellCommand is the words of one command: those before a shell operator or
// redirection.
func shellCommand(tokens []string) []string {
	for i, tok := range tokens {
		if shellOperators[tok] || redirection.MatchString(tok) {
			return tokens[:i]
		}
	}
	return tokens
}

// flagsIn is the flags among a command's words, as written up to any "=value".
// "[--flag <value>]" marks an optional flag, and "[--a, --b]" several.
func flagsIn(tokens []string) []string {
	var flags []string
	for _, tok := range tokens {
		tok = strings.TrimRight(strings.TrimLeft(tok, "["), "],")
		tok, _, _ = strings.Cut(tok, "=")
		switch {
		case strings.HasPrefix(tok, "--") && len(tok) > 2:
			flags = append(flags, tok)
		case len(tok) == 2 && tok[0] == '-' && tok[1] >= 'a' && tok[1] <= 'z':
			flags = append(flags, tok)
		}
	}
	return flags
}

// shellWords splits a line as a shell would, enough for examples: quotes group
// words, and a # that starts a word begins a comment.
func shellWords(line string) []string {
	var words []string
	var cur strings.Builder
	inWord := false
	var quote rune
	for _, r := range line {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				cur.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote, inWord = r, true
		case r == ' ' || r == '\t':
			if inWord {
				words, inWord = append(words, cur.String()), false
				cur.Reset()
			}
		case r == '#' && !inWord:
			return words
		default:
			cur.WriteRune(r)
			inWord = true
		}
	}
	if inWord {
		words = append(words, cur.String())
	}
	return words
}

var longFlag = regexp.MustCompile(`^--[a-z0-9][a-z0-9-]*$`)

// checkExample reports a word that names no command, and a flag the command
// does not take. Words after a command with no subcommands are its arguments.
func checkExample(root *cobra.Command, ex skillExample) []string {
	c := root
	for _, w := range ex.words {
		if !c.HasSubCommands() {
			break
		}
		next := subcommand(c, w)
		if next == nil {
			return []string{fmt.Sprintf("%q has no command %q", c.CommandPath(), w)}
		}
		c = next
	}
	if c.HasSubCommands() && ex.generic {
		return nil
	}
	if c.HasSubCommands() && !c.Runnable() && !slices.Contains(ex.flags, "--help") {
		return []string{fmt.Sprintf("%q needs a command after it", c.CommandPath())}
	}
	var problems []string
	for _, f := range ex.flags {
		switch {
		case f == "--help" || f == "-h":
		case strings.HasPrefix(f, "--"):
			name := f[2:]
			if !longFlag.MatchString(f) || c.Flags().Lookup(name) == nil && c.InheritedFlags().Lookup(name) == nil {
				problems = append(problems, fmt.Sprintf("%q has no flag %s", c.CommandPath(), f))
			}
		default:
			if c.Flags().ShorthandLookup(f[1:]) == nil && c.InheritedFlags().ShorthandLookup(f[1:]) == nil {
				problems = append(problems, fmt.Sprintf("%q has no flag %s", c.CommandPath(), f))
			}
		}
	}
	return problems
}

func subcommand(c *cobra.Command, name string) *cobra.Command {
	for _, child := range c.Commands() {
		if child.Name() == name || child.HasAlias(name) {
			return child
		}
	}
	return nil
}
