// Package messaging publishes the Webex Messaging API commands, built from the
// operation table, and a handwritten command for messages with attachments.
package messaging

import (
	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/apicmd"
)

var groups = apicmd.Register(cmd.MessagingCmd, "messaging")

// Groups that handwritten commands join.
var messagesCmd = groups["messages"]
