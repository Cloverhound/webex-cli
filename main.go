package main

import (
	"os"

	"github.com/Cloverhound/webex-cli/cmd"
	_ "github.com/Cloverhound/webex-cli/cmd/calling"
	_ "github.com/Cloverhound/webex-cli/cmd/cc"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
