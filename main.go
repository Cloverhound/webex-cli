package main

import (
	"os"

	"github.com/Cloverhound/webex-cli/cmd"
	_ "github.com/Cloverhound/webex-cli/cmd/admin"
	_ "github.com/Cloverhound/webex-cli/cmd/calling"
	_ "github.com/Cloverhound/webex-cli/cmd/cc"
	_ "github.com/Cloverhound/webex-cli/cmd/device"
	_ "github.com/Cloverhound/webex-cli/cmd/meetings"
	_ "github.com/Cloverhound/webex-cli/cmd/messaging"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
