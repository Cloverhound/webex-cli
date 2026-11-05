// Package device publishes the Webex Device API commands.
package device

import (
	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/apicmd"
)

var _ = apicmd.Register(cmd.DeviceCmd, "device")
