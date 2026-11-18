// Package calling publishes the Webex Calling API commands, built from the
// operation table, and handwritten commands for file uploads and downloads.
package calling

import (
	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/apicmd"
)

var groups = apicmd.Register(cmd.CallingCmd, "calling")

// Groups that handwritten commands join.
var (
	announcementRepositoryCmd = groups["announcement-repository"]
	convergedRecordingsCmd    = groups["converged-recordings"]
	userCallCmd               = groups["user-call"]
	callSettingsForMeCmd      = groups["call-settings-for-me"]
	virtualLineCallCmd        = groups["virtual-line-call"]
	workspaceCallCmd          = groups["workspace-call"]
)
