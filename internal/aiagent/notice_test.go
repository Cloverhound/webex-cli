package aiagent

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrivateAPINoticeShowsOnce(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	var out bytes.Buffer
	noticeOutput = &out
	t.Cleanup(func() { noticeOutput = defaultNoticeOutput })

	PrivateAPINotice()
	if !strings.Contains(out.String(), "not documented or supported by Cisco") {
		t.Fatalf("first run printed %q", out.String())
	}
	out.Reset()
	PrivateAPINotice()
	if out.Len() != 0 {
		t.Fatalf("second run printed %q", out.String())
	}
}
