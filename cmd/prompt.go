package cmd

import (
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"golang.org/x/term"
)

const nonInteractiveEnv = "WEBEX_CLI_NONINTERACTIVE"

// isTerminal is a variable so tests can simulate a terminal.
var isTerminal = func(f *os.File) bool { return term.IsTerminal(int(f.Fd())) }

// canPrompt reports whether a person can answer a form. Forms are drawn on
// stderr, so stdout can be piped or carry JSON. Installers and agents usually
// run without a terminal, where a form would fail or hang.
func canPrompt() bool {
	if off, _ := strconv.ParseBool(os.Getenv(nonInteractiveEnv)); off {
		return false
	}
	return isTerminal(os.Stdin) && isTerminal(os.Stderr)
}

// confirm asks a yes/no question on stderr.
func confirm(title, description, yes, no string) (bool, error) {
	var ok bool
	err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title(title).
			Description(description).
			Affirmative(yes).
			Negative(no).
			Value(&ok),
	)).WithOutput(os.Stderr).Run()
	return ok, err
}
