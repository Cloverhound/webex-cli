package aiagent

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Cloverhound/webex-cli/internal/appconfig"
)

var defaultNoticeOutput io.Writer = os.Stderr

// noticeOutput is a variable so tests can capture the notice.
var noticeOutput = defaultNoticeOutput

// noticeMarker records that this user has seen the notice, so scripts and
// agents that run these commands repeatedly are not flooded with it.
const noticeMarker = "private-api-notice-shown"

// PrivateAPINotice tells the user, once, that the command uses an API Cisco
// does not document or support. Every such command repeats the caution in
// its --help.
func PrivateAPINotice() {
	marker := filepath.Join(appconfig.Dir(), noticeMarker)
	if _, err := os.Stat(marker); err == nil {
		return
	}
	fmt.Fprintln(noticeOutput, "Caution: webex cc ai-agent commands use Webex AI Agent Studio's private API, "+
		"which is not documented or supported by Cisco and can change or stop working without notice. "+
		"Do not build production automation on it. (Shown once.)")
	if err := os.MkdirAll(filepath.Dir(marker), 0o700); err == nil {
		os.WriteFile(marker, nil, 0o600)
	}
}
