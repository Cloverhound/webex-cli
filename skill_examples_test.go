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
// rename a command or flag under them. Every `webex` line in a code block of a
// skill file or a documentation page must name a real command and only flags
// that command takes. The generated Command Reference is skipped: codegen
// writes it from the commands.
func TestDocumentedExamplesUseRealCommandsAndFlags(t *testing.T) {
	files, err := skill.Files()
	if err != nil {
		t.Fatal(err)
	}
	docs := map[string][]byte{"skill/SETUP.md": []byte(skill.Setup)}
	for name, content := range files {
		docs["skill/"+name] = content
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
		docs[name] = content
	}
	names := make([]string, 0, len(docs))
	for name := range docs {
		names = append(names, name)
	}
	sort.Strings(names)

	root := cmd.RootCommand()
	for _, name := range names {
		for _, ex := range skillExamples(string(docs[name])) {
			for _, problem := range checkExample(root, ex.words, ex.flags) {
				t.Errorf("%s:%d: %s\n\t%s", name, ex.line, problem, ex.text)
			}
		}
		area, file, ok := strings.Cut(strings.TrimPrefix(name, "skill/"), "/")
		if !ok || file != "SKILL.md" {
			continue
		}
		for _, row := range commandTableRows(string(files[name])) {
			for _, command := range row.commands {
				for _, problem := range checkExample(root, []string{area, row.group, command}, nil) {
					t.Errorf("%s:%d: %s\n\t%s", name, row.line, problem, row.text)
				}
			}
		}
	}
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

// skillExample is one `webex` command line in a skill file, with the flags on
// its continuation lines.
type skillExample struct {
	line  int
	text  string
	words []string // command words after "webex", up to the first argument or flag
	flags []string // flag names, without "--"
}

var (
	commandWord = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	flagName    = regexp.MustCompile(`^--([a-z0-9][a-z0-9-]*)`)
)

// skillExamples finds the `webex` lines in a file's fenced code blocks, outside
// the generated block. A line in the same block that starts with a flag
// continues the command above it.
func skillExamples(content string) []skillExample {
	var out []skillExample
	inCode, generated := false, false
	var last *skillExample
	for i, raw := range strings.Split(content, "\n") {
		line := strings.TrimSpace(raw)
		switch {
		case strings.HasPrefix(line, "<!-- codegen:start"):
			generated = true
		case strings.HasPrefix(line, "<!-- codegen:end"):
			generated = false
		case strings.HasPrefix(line, "```"):
			inCode, last = !inCode, nil
		}
		if !inCode || generated {
			continue
		}
		tokens := shellWords(line)
		switch {
		case len(tokens) > 0 && tokens[0] == "webex":
			ex := skillExample{line: i + 1, text: line}
			args := tokens[1:]
			for len(args) > 0 && commandWord.MatchString(args[0]) {
				ex.words, args = append(ex.words, args[0]), args[1:]
			}
			// "webex <area> <resource> --help" stands for many commands.
			if len(ex.words) == 0 || len(args) > 0 && strings.HasPrefix(args[0], "<") {
				last = nil
				continue
			}
			ex.flags = flagsIn(args)
			out = append(out, ex)
			last = &out[len(out)-1]
		case len(tokens) > 0 && strings.HasPrefix(tokens[0], "--") && last != nil:
			last.flags = append(last.flags, flagsIn(tokens)...)
		case line == "":
			last = nil
		}
	}
	return out
}

func flagsIn(tokens []string) []string {
	var flags []string
	for _, tok := range tokens {
		// "[--flag <value>]" marks an optional flag.
		if m := flagName.FindStringSubmatch(strings.TrimLeft(tok, "[")); m != nil {
			flags = append(flags, m[1])
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

// checkExample reports a word that names no command, and a flag the command
// does not take. Words after a command with no subcommands are its arguments.
func checkExample(root *cobra.Command, words, flags []string) []string {
	c := root
	for _, w := range words {
		if !c.HasSubCommands() {
			break
		}
		next := subcommand(c, w)
		if next == nil {
			return []string{fmt.Sprintf("%q has no command %q", c.CommandPath(), w)}
		}
		c = next
	}
	if c.HasSubCommands() && !c.Runnable() && !slices.Contains(flags, "help") {
		return []string{fmt.Sprintf("%q needs a command after it", c.CommandPath())}
	}
	var problems []string
	for _, f := range flags {
		if f != "help" && c.Flags().Lookup(f) == nil && c.InheritedFlags().Lookup(f) == nil {
			problems = append(problems, fmt.Sprintf("%q has no flag --%s", c.CommandPath(), f))
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
