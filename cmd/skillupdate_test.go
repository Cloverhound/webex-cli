package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Cloverhound/webex-cli/skill"
)

func TestCheckSkillUpdatesNonInteractive(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	claudeDir := filepath.Join(home, ".claude", "skills", "webex-cli")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "SKILL.md"), []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	// Detected but never installed: left alone without --yes.
	if err := os.MkdirAll(filepath.Join(home, ".codex"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := checkSkillUpdates(setupOptions{skills: true}); err != nil {
		t.Fatal(err)
	}

	want, _ := skill.Files()
	got, err := os.ReadFile(filepath.Join(claudeDir, "admin", "SKILL.md"))
	if err != nil || string(got) != string(want["admin/SKILL.md"]) {
		t.Errorf("outdated Claude Code skill not updated: %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, ".codex", "skills", "webex-cli", "SKILL.md")); !os.IsNotExist(err) {
		t.Errorf("Codex skill installed without --yes: %v", err)
	}
}
