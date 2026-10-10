package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
	"github.com/charmbracelet/huh"
	"golang.org/x/term"
)

const readOnlyEnv = "WEBEX_READ_ONLY"

// readOnlyEnvSet reports whether $WEBEX_READ_ONLY asks for read-only mode.
func readOnlyEnvSet() bool {
	on, _ := strconv.ParseBool(os.Getenv(readOnlyEnv))
	return on
}

// readOnlyMode reports whether read-only mode is on. The environment can turn
// it on for one process but never off.
func readOnlyMode(cfg *appconfig.Config) bool {
	return cfg.ReadOnly || readOnlyEnvSet()
}

// confirmLeaveReadOnly asks the person at the terminal to approve a login with
// write access. Agents usually run commands without a terminal, so requiring
// one keeps them from leaving read-only mode on their own.
func confirmLeaveReadOnly() error {
	if readOnlyEnvSet() {
		return fmt.Errorf("$%s is set; unset it to log in with write access", readOnlyEnv)
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return fmt.Errorf("read-only mode is on; leaving it requires 'webex login' from an interactive terminal")
	}

	var leave bool
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Leave read-only mode?").
				Description("This logs in with write access. Only continue if you started this login yourself.").
				Affirmative("Leave read-only mode").
				Negative("Cancel").
				Value(&leave),
		),
	).Run()
	if err != nil {
		return err
	}
	if !leave {
		return fmt.Errorf("cancelled; still in read-only mode")
	}
	return nil
}
