// Package admin publishes the Webex Admin API commands.
package admin

import (
	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/apicmd"
)

var _ = apicmd.Register(cmd.AdminCmd, "admin")
