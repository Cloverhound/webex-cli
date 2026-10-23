package skill

import "testing"

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
