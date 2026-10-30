package skill

import (
	"strings"
	"testing"
)

func TestSetupIsAStandaloneSkill(t *testing.T) {
	if !strings.HasPrefix(Setup, "---\nname: webex-cli-setup\n") {
		t.Errorf("setup skill has no webex-cli-setup front matter: %.40q", Setup)
	}
	files, err := Files()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := files["SETUP.md"]; ok {
		t.Error("SETUP.md would be installed inside the webex-cli skill")
	}
}

func TestFilesIncludesEveryArea(t *testing.T) {
	files, err := Files()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"SKILL.md", "admin/SKILL.md", "calling/SKILL.md", "cc/SKILL.md", "device/SKILL.md", "meetings/SKILL.md", "messaging/SKILL.md"} {
		if len(files[p]) == 0 {
			t.Errorf("missing or empty %s", p)
		}
	}
}
