// Package meetings publishes the Webex Meetings API commands, built from the
// operation table, plus a handwritten recording download.
package meetings

import (
	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/apicmd"
)

var groups = apicmd.Register(cmd.MeetingsCmd, "meetings")

// recordingsCmd is the group a handwritten command joins.
var recordingsCmd = groups["recordings"]
