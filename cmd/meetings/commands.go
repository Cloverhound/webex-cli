// Package meetings publishes the Webex Meetings API commands, built from the
// operation table, and a handwritten command for recording downloads.
package meetings

import (
	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/apicmd"
)

var groups = apicmd.Register(cmd.MeetingsCmd, "meetings")

// Groups that handwritten commands join.
var recordingsCmd = groups["recordings"]
