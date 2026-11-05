// Package messaging publishes the Webex Messaging API commands, built from the
// operation table, plus a handwritten message create that uploads files.
package messaging

import (
	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/apicmd"
)

var groups = apicmd.Register(cmd.MessagingCmd, "messaging")

// messagesCmd is the group a handwritten command joins.
var messagesCmd = groups["messages"]
