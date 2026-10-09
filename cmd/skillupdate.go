package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Cloverhound/webex-cli/skill"
	"github.com/charmbracelet/huh"
)

// checkSkillUpdates compares installed agent skills with the skill embedded in
// this binary. Without prompts, outdated skills are updated and agents without
// the skill are left alone, unless opts.assumeYes is set.
func checkSkillUpdates(opts setupOptions) error {
	if !opts.skills {
		return nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	type skillAction struct {
		Platform agentPlatform
		Path     string
		Status   string // "outdated", "new"
	}

	latest, err := skill.Files()
	if err != nil {
		return fmt.Errorf("reading embedded skill: %w", err)
	}

	var actions []skillAction
	for _, p := range agentPlatforms {
		if p.SkillDir == "" {
			continue // skip Cowork (no filesystem path)
		}
		skillDir := filepath.Join(home, p.SkillDir)
		if _, err := os.Stat(filepath.Join(skillDir, "SKILL.md")); err == nil {
			outdated := false
			for subPath, content := range latest {
				current, rerr := os.ReadFile(filepath.Join(skillDir, filepath.FromSlash(subPath)))
				if rerr != nil || !bytes.Equal(current, content) {
					outdated = true
					break
				}
			}
			if outdated {
				actions = append(actions, skillAction{p, skillDir, "outdated"})
			}
		} else if agentDetected(home, p) {
			actions = append(actions, skillAction{p, skillDir, "new"})
		}
	}

	if len(actions) == 0 {
		fmt.Println("All installed skills are up to date.")
		return nil
	}

	var selected []string
	if !opts.prompt {
		for _, a := range actions {
			if a.Status == "outdated" || opts.assumeYes {
				selected = append(selected, a.Platform.Name)
			} else {
				fmt.Printf("  %s: skill not installed (run: webex post-install --yes)\n", a.Platform.Name)
			}
		}
	} else {
		// Step 1: Ask if user wants to review skill updates
		var proceed bool
		confirmForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Agent skill updates available").
					Description(fmt.Sprintf("%d agent skill(s) can be updated or installed.", len(actions))).
					Affirmative("Review changes").
					Negative("Skip").
					Value(&proceed),
			),
		)

		if err := confirmForm.Run(); err != nil || !proceed {
			return nil
		}

		// Step 2: Show multiselect with specific skills
		var options []huh.Option[string]
		for _, a := range actions {
			label := a.Platform.Name
			if a.Status == "outdated" {
				label += "  (update available)"
			} else {
				label += "  (not yet installed)"
			}
			options = append(options, huh.NewOption(label, a.Platform.Name).Selected(true))
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Skill updates available").
					Description("The following agents can be updated or have the Webex skill\ninstalled. Deselect any you want to skip.").
					Options(options...).
					Value(&selected),
			),
		)

		if err := form.Run(); err != nil {
			return nil
		}
	}

	if len(selected) == 0 {
		return nil
	}

	for _, name := range selected {
		for _, a := range actions {
			if a.Platform.Name == name {
				failed := false
				for subPath, content := range latest {
					dest := filepath.Join(a.Path, filepath.FromSlash(subPath))
					if err := installSkill(dest, content); err != nil {
						fmt.Printf("  %s: failed installing %s (%v)\n", a.Platform.Name, subPath, err)
						failed = true
						break
					}
				}
				if !failed {
					if a.Status == "outdated" {
						fmt.Printf("  %s: updated %s\n", a.Platform.Name, a.Path)
					} else {
						fmt.Printf("  %s: installed to %s\n", a.Platform.Name, a.Path)
					}
				}
			}
		}
	}

	return nil
}
