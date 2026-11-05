// Package cc publishes the Webex Contact Center API commands, built from the
// operation table, plus handwritten commands for search and file uploads.
package cc

import (
	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/apicmd"
)

var groups = apicmd.Register(cmd.CcCmd, "cc")

// Groups that handwritten commands join.
var (
	audioFilesCmd                 = groups["audio-files"]
	agentPersonalGreetingFilesCmd = groups["agent-personal-greeting-files"]
)
