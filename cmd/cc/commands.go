// Package cc publishes the Webex Contact Center API commands, built from the
// operation table, and handwritten commands for search and audio files.
package cc

import (
	cmd "github.com/Cloverhound/webex-cli/cmd"
	"github.com/Cloverhound/webex-cli/internal/apicmd"
)

var groups = apicmd.Register(cmd.CcCmd, "cc")

// Groups that handwritten commands join.
var (
	agentPersonalGreetingFilesCmd = groups["agent-personal-greeting-files"]
	audioFilesCmd                 = groups["audio-files"]
)
