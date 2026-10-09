package main

import (
	"os"

	"github.com/Cloverhound/webex-cli/cmd"
	_ "github.com/Cloverhound/webex-cli/cmd/admin"
	_ "github.com/Cloverhound/webex-cli/cmd/calling"
	_ "github.com/Cloverhound/webex-cli/cmd/cc"
	_ "github.com/Cloverhound/webex-cli/cmd/device"
	_ "github.com/Cloverhound/webex-cli/cmd/meetings"
	_ "github.com/Cloverhound/webex-cli/cmd/mcp"
	_ "github.com/Cloverhound/webex-cli/cmd/messaging"

	// Used only when the OS has no CA bundle, as in slim container images.
	_ "golang.org/x/crypto/x509roots/fallback"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
